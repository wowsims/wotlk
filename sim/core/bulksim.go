package core

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sort"

	goproto "github.com/golang/protobuf/proto"

	"github.com/wowsims/wotlk/sim/core/proto"
)

// RunBulkSim runs a bulk simulation. The original sim request must contain exactly one player
// (i.e. not a full raid sim) with a specified bulk equipment.
func RunBulkSim(ctx context.Context, request *proto.RaidSimRequest, progress chan *proto.ProgressMetrics) (*proto.RaidSimResult, error) {
	// Bulk simming is only supported for the single-player use (i.e. not whole raid-wide simming).
	// Verify that we have exactly 1 player.
	var playerCount int
	var player *proto.Player
	for _, p := range request.GetRaid().GetParties() {
		for _, pl := range p.GetPlayers() {
			// TODO(Riotdog-GehennasEU): Is this a reasonable check?
			if pl.Name != "" {
				player = pl
				playerCount++
			}
		}
	}
	if playerCount != 1 || player == nil {
		return nil, fmt.Errorf("bulksim: expected exactly 1 player, found %d", playerCount)
	}

	// TODO(Riotdog-GehennasEU): Expose this as a setting?
	concurrency := (runtime.NumCPU() - 1) * 2
	if concurrency <= 0 {
		concurrency = 1
	}

	results := make(chan *bulkRaidSimResult, concurrency)
	go func() {
		var numCombinations int32
		for sub := range generateAllEquipmentSubstitutions(ctx, player.GetBulkEquipment()) {
			substitutedRequest, changeLog := createNewRequestWithSubstitution(request, sub)
			if isValidEquipment(substitutedRequest.Raid.Parties[0].Players[0].Equipment) {
				singleSimProgress := make(chan *proto.ProgressMetrics)
				go func(i int32) {
					for p := range singleSimProgress {
						// Do not forward the final message, since we have to do multiple invidiual sims and this
						// would confuse the UI.
						if p.FinalRaidResult != nil {
							continue
						}

						p.CompletedIterations += i * request.SimOptions.Iterations
						p.TotalIterations += i * request.SimOptions.Iterations
						select {
						case progress <- p:
						default:
							// We tried. Do not block here because it could slow down the sim.
						}
					}
				}(numCombinations)

				results <- &bulkRaidSimResult{
					Request:      substitutedRequest,
					Result:       runSim(substitutedRequest, singleSimProgress, false),
					Substitution: sub,
					ChangeLog:    changeLog,
				}
			}
			numCombinations++
		}
		close(results)
	}()

	var rankedResults []*bulkRaidSimResult
	for r := range results {
		rankedResults = append(rankedResults, r)
	}
	sort.Slice(rankedResults, func(i, j int) bool {
		return rankedResults[i].Score() > rankedResults[j].Score()
	})

	totalResultCount := int32(len(rankedResults))
	var baseResult *bulkRaidSimResult
	for _, r := range rankedResults {
		if !r.Substitution.HasItemReplacements() {
			baseResult = r
		}
	}

	// TODO(Riotdog-GehennasEU): Make this configurable?
	maxResults := 10
	if len(rankedResults) > maxResults {
		rankedResults = rankedResults[:maxResults]
	}

	for _, r := range rankedResults {
		baseResult.Result.BulkResults = append(baseResult.Result.BulkResults, &proto.BulkSimResultWithSubstitutions{
			ItemsAdded:       itemWithSlotToBulkSpec(r.ChangeLog.AddedItems),
			ItemsRemoved:     itemWithSlotToBulkSpec(r.ChangeLog.RemovedItems),
			RaidMetrics:      r.Result.RaidMetrics,
			EncounterMetrics: r.Result.EncounterMetrics,
			ErrorResult:      r.Result.ErrorResult,
		})
	}

	if progress != nil {
		progress <- &proto.ProgressMetrics{
			TotalIterations:     request.SimOptions.Iterations * totalResultCount,
			CompletedIterations: request.SimOptions.Iterations * totalResultCount,
			Dps:                 baseResult.Result.RaidMetrics.Dps.Avg,
			FinalRaidResult:     baseResult.Result,
		}
	}

	return baseResult.Result, nil
}

// bulkRaidSimResult stores the request and response of a simulation, along with the used equipment
// susbstitution and a changelog of which items were added and removed from the base equipment set.
type bulkRaidSimResult struct {
	Request      *proto.RaidSimRequest
	Result       *proto.RaidSimResult
	Substitution *equipmentSubstitution
	ChangeLog    *raidSimRequestChangeLog
}

// Score used to rank results.
func (r *bulkRaidSimResult) Score() float64 {
	return r.Result.RaidMetrics.Dps.Avg
}

// equipmentSubstitution specifies all items to be used as replacements for the equipped gear.
type equipmentSubstitution struct {
	Items []*itemWithSlot
}

// HasChanges returns true if the equipment substitution has any item replacmenets.
func (es *equipmentSubstitution) HasItemReplacements() bool {
	return len(es.Items) > 0
}

// IsValid returns true if the equipment substituion is valid. A valid substition can only
// reference an item or slot once.
func (es *equipmentSubstitution) IsValid() bool {
	slotReuseTracker := map[proto.ItemSlot]bool{}
	itemReuseTracker := map[int]bool{}

	for _, it := range es.Items {
		if itemReuseTracker[it.Index] {
			return false
		}
		itemReuseTracker[it.Index] = true

		if slotReuseTracker[it.Slot] {
			return false
		}
		slotReuseTracker[it.Slot] = true
	}

	return true
}

// isValidEquipment returns true if the specified equipment spec is valid. A equipment spec
// is valid if it does not reference a two-hander and off-hand weapon combo.
func isValidEquipment(equipment *proto.EquipmentSpec) bool {
	var usesTwoHander, usesOffhand bool

	for _, it := range equipment.Items {
		if it.GetId() == 0 {
			continue
		}

		knownItem, ok := ItemsByID[it.Id]
		if !ok {
			// TODO(Riotdog-GehennasEU): Should we bother verifying that the item is in the database?
			// What about gems and enchants? Should we just expect that we will receive a valid database?
			log.Printf("Warning: bulk item %d not found in the provided database", it.Id)
			return false
		}

		if knownItem.HandType == proto.HandType_HandTypeTwoHand {
			usesTwoHander = true
		}
		if knownItem.HandType == proto.HandType_HandTypeOffHand {
			usesOffhand = true
		}
	}

	return !(usesTwoHander && usesOffhand)
}

// generateAllEquipmentSubstitutions generates all possible valid equipment substitutions for the
// given bulk sim request. Also returns the unchanged equipment ("base equipment set") set as the
// first result. This ensures that simming over all possible equipment substitutions includes the
// base case as well.
func generateAllEquipmentSubstitutions(ctx context.Context, spec *proto.BulkEquipmentSpec) chan *equipmentSubstitution {
	results := make(chan *equipmentSubstitution)
	go func() {
		defer close(results)

		// No substitutions (base case).
		results <- &equipmentSubstitution{}

		// Create all distinct combinations of (item, slot). For example, let's say the only item we
		// want to bulk sim is a one-handed item that can be worn both as an off-hand or a main-hand weapon.
		// For each slot, we will create one itemWithSlot pair, so (item, off-hand) and (item, main-hand).
		// We verify later that we are not emitting any item substitution that refers to the same item.
		var distinctItemSlotCombos []*itemWithSlot
		for i, is := range spec.GetItems() {
			for _, slot := range is.Slots {
				distinctItemSlotCombos = append(distinctItemSlotCombos, &itemWithSlot{
					Item:  is.Item,
					Slot:  slot,
					Index: i,
				})
			}
		}

		// Borrowed from https://github.com/mxschmitt/golang-combinations and adapted to
		// only emit valid combinations.
		count := uint(len(distinctItemSlotCombos))
		for bits := 1; bits < (1 << count); bits++ {
			combo := &equipmentSubstitution{}
			for idx := uint(0); idx < count; idx++ {
				if (bits>>idx)&1 == 1 {
					combo.Items = append(combo.Items, distinctItemSlotCombos[idx])
				}
			}

			if !combo.IsValid() {
				continue
			}

			select {
			case <-ctx.Done():
				return
			case results <- combo:
			}
		}
	}()

	return results
}

// itemWithSlot pairs an item with its fixed item slot.
type itemWithSlot struct {
	Item *proto.ItemSpec
	Slot proto.ItemSlot

	// This index refers to the item's position in the BulkEquipmentSpec of the player and serves as
	// a unique item ID. It is used to verify that a valid equipmentSubstitution only references an
	// item once.
	Index int
}

// itemWithSlotToBulkSpec converts the given slice of itemWithSlot to a BulkEquipmentSpec.
func itemWithSlotToBulkSpec(items []*itemWithSlot) *proto.BulkEquipmentSpec {
	r := &proto.BulkEquipmentSpec{}
	for _, it := range items {
		r.Items = append(r.Items, &proto.BulkEquipmentSpec_ItemSpecWithSlots{
			Item:  it.Item,
			Slots: []proto.ItemSlot{it.Slot},
		})
	}
	return r
}

// raidSimRequestChangeLog stores a change log of which items were added and removed from the base
// equipment set.
type raidSimRequestChangeLog struct {
	AddedItems   []*itemWithSlot
	RemovedItems []*itemWithSlot
}

// createNewRequestWithSubstitution creates a copy of the input RaidSimRequest and applis the given
// equipment susbstitution to the player's equipment.
func createNewRequestWithSubstitution(readonlyInputRequest *proto.RaidSimRequest, substitution *equipmentSubstitution) (*proto.RaidSimRequest, *raidSimRequestChangeLog) {
	request := goproto.Clone(readonlyInputRequest).(*proto.RaidSimRequest)
	changeLog := &raidSimRequestChangeLog{}
	player := request.Raid.Parties[0].Players[0]
	equipment := player.Equipment
	for _, is := range substitution.Items {
		equippedItem := equipment.Items[is.Slot]
		if equippedItem.GetId() > 0 {
			changeLog.RemovedItems = append(changeLog.RemovedItems, &itemWithSlot{
				Item: equippedItem,
				Slot: is.Slot,
			})
		}
		changeLog.AddedItems = append(changeLog.AddedItems, is)
		equipment.Items[is.Slot] = is.Item
	}

	// Clear the playuer's bulk equipment set since we don't want to recursively keep bulk simming.
	player.BulkEquipment = nil
	return request, changeLog
}
