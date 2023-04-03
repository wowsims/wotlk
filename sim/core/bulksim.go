package core

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"strings"

	goproto "github.com/golang/protobuf/proto"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type ItemWithSlot struct {
	Item  *proto.ItemSpec
	Slot  proto.ItemSlot
	Index int
}

type EquipmentSubstitution struct {
	Items []*ItemWithSlot
}

func (es *EquipmentSubstitution) IsEquipped() bool {
	// No substitutions == equipped gear.
	return len(es.Items) == 0
}

func (es *EquipmentSubstitution) IsValid() bool {
	slotReuseTracker := map[int]bool{}
	itemReuseTracker := map[int]bool{}
	var usesTwoHander, usesOffhand bool

	for _, it := range es.Items {
		if itemReuseTracker[it.Index] {
			return false
		}
		itemReuseTracker[it.Index] = true

		if slotReuseTracker[it.Index] {
			return false
		}
		slotReuseTracker[it.Index] = true

		knownItem, ok := ItemsByID[it.Item.Id]
		if !ok {
			return false
		}

		if knownItem.Type == proto.ItemType_ItemTypeWeapon {
			if knownItem.HandType == proto.HandType_HandTypeTwoHand {
				usesTwoHander = true
			}
			if knownItem.HandType == proto.HandType_HandTypeOffHand {
				usesOffhand = true
			}
		}
	}

	return !(usesTwoHander && usesOffhand)
}

func GenerateAllEquipmentSubstitutions(ctx context.Context, spec *proto.BulkEquipmentSpec) chan *EquipmentSubstitution {
	results := make(chan *EquipmentSubstitution)
	go func() {
		defer close(results)

		// No substitutions (base case).
		results <- &EquipmentSubstitution{}

		// Create all distinct combinations of (item, slot). For example, let's say the only item we
		// want to bulk sim is a one-handed item that can be worn both as an off-hand or a main-hand weapon.
		// For each slot, we will create one ItemWithSlot pair, so (item, off-hand) and (item, main-hand).
		// We verify later that we are not emitting any item substitution that refers to the same item.
		var distinctItemSlotCombos []*ItemWithSlot
		for i, is := range spec.Items {
			for _, slot := range is.Slots {
				distinctItemSlotCombos = append(distinctItemSlotCombos, &ItemWithSlot{
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
			combo := &EquipmentSubstitution{}
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

type RaidSimRequestChangeLog struct {
	AddedItems   []*ItemWithSlot
	RemovedItems []*ItemWithSlot
}

func (cl *RaidSimRequestChangeLog) String() string {
	var buf strings.Builder
	for _, i := range cl.AddedItems {
		buf.WriteString(fmt.Sprintf("[+] %d:%d:%v\n", i.Item.Id, i.Item.Enchant, i.Item.Gems))
	}
	for _, i := range cl.RemovedItems {
		buf.WriteString(fmt.Sprintf("[-] %d:%d:%v\n", i.Item.Id, i.Item.Enchant, i.Item.Gems))
	}
	return strings.TrimSpace(buf.String())
}

func ItemWithSlotToBulkSpec(items []*ItemWithSlot) *proto.BulkEquipmentSpec {
	r := &proto.BulkEquipmentSpec{}
	for _, it := range items {
		r.Items = append(r.Items, &proto.BulkEquipmentSpec_ItemSpecWithSlots{
			Item:  it.Item,
			Slots: []proto.ItemSlot{it.Slot},
		})
	}
	return r
}

func CreateNewRequestWithSubstitution(readonlyInputRequest *proto.RaidSimRequest, substitution *EquipmentSubstitution) (*proto.RaidSimRequest, *RaidSimRequestChangeLog) {
	request := goproto.Clone(readonlyInputRequest).(*proto.RaidSimRequest)
	changeLog := &RaidSimRequestChangeLog{}
	player := request.Raid.Parties[0].Players[0]
	equipment := player.Equipment
	for _, is := range substitution.Items {
		equippedItem := equipment.Items[is.Slot]
		if equippedItem.GetId() > 0 {
			changeLog.RemovedItems = append(changeLog.RemovedItems, &ItemWithSlot{
				Item: equippedItem,
				Slot: is.Slot,
			})
		}
		changeLog.AddedItems = append(changeLog.AddedItems, is)
		equipment.Items[is.Slot] = is.Item
	}

	player.BulkEquipment = nil
	return request, changeLog
}

type BulkRaidSimResult struct {
	Request      *proto.RaidSimRequest
	Result       *proto.RaidSimResult
	Substitution *EquipmentSubstitution
	ChangeLog    *RaidSimRequestChangeLog
}

func (r *BulkRaidSimResult) Score() float64 {
	return r.Result.RaidMetrics.Dps.Avg
}

func RunBulkSim(request *proto.RaidSimRequest, progress chan *proto.ProgressMetrics) (*proto.RaidSimResult, error) {
	// Bulk simming is only supported for the single-player use (i.e. not whole raid-wide simming).
	// Verify that we have exactly 1 player.
	var playerCount int
	var player *proto.Player
	for _, p := range request.GetRaid().GetParties() {
		for _, pl := range p.GetPlayers() {
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

	results := make(chan *BulkRaidSimResult, concurrency)
	substitutions := GenerateAllEquipmentSubstitutions(context.Background(), player.GetBulkEquipment())
	go func() {
		var numCombinations int32
		for sub := range substitutions {
			substitutedRequest, changeLog := CreateNewRequestWithSubstitution(request, sub)
			// TODO(Riotdog-GehennasEU): We could do this a bit nicer: create a new progress reporter, and then forward
			// forward those results instead. That way we could accurately aggregate the progress metrics.
			results <- &BulkRaidSimResult{
				Request:      substitutedRequest,
				Result:       runSim(substitutedRequest, nil, false),
				Substitution: sub,
				ChangeLog:    changeLog,
			}

			if progress != nil {
				progress <- &proto.ProgressMetrics{
					TotalIterations:     (numCombinations + 1) * request.SimOptions.Iterations, // So close!
					CompletedIterations: numCombinations * request.SimOptions.Iterations,
				}
			}
			numCombinations++
		}
		close(results)
	}()

	var rankedResults []*BulkRaidSimResult
	for r := range results {
		rankedResults = append(rankedResults, r)
	}
	sort.Slice(rankedResults, func(i, j int) bool {
		return rankedResults[i].Score() > rankedResults[j].Score()
	})

	totalResultCount := int32(len(rankedResults))
	var baseResult *BulkRaidSimResult
	for _, r := range rankedResults {
		if r.Substitution.IsEquipped() {
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
			ItemsAdded:       ItemWithSlotToBulkSpec(r.ChangeLog.AddedItems),
			ItemsRemoved:     ItemWithSlotToBulkSpec(r.ChangeLog.RemovedItems),
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
