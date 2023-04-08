package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"runtime"
	"runtime/debug"
	"sort"
	"sync/atomic"
	"time"

	goproto "github.com/golang/protobuf/proto"

	"github.com/wowsims/wotlk/sim/core/proto"
)

const (
	maxItemCount              = 22
	defaultIterationsPerCombo = 1000
)

// raidSimRunner runs a standard raid simulation.
type raidSimRunner func(*proto.RaidSimRequest, chan *proto.ProgressMetrics, bool) *proto.RaidSimResult

// bulkSimRunner runs a bulk simulation.
type bulkSimRunner struct {
	// SingleRaidSimRunner used to run one simulation of the bulk.
	SingleRaidSimRunner raidSimRunner
	// Request used for this bulk simulation.
	Request *proto.BulkSimRequest
}

func BulkSim(ctx context.Context, request *proto.BulkSimRequest, progress chan *proto.ProgressMetrics) *proto.BulkSimResult {
	bulk := &bulkSimRunner{
		SingleRaidSimRunner: runSim,
		Request:             request,
	}

	result, err := bulk.Run(ctx, progress)
	if err != nil {
		result = &proto.BulkSimResult{
			ErrorResult: err.Error(),
		}
	}

	if progress != nil {
		progress <- &proto.ProgressMetrics{
			FinalBulkResult: result,
		}
		close(progress)
	}

	return result
}

type singleBulkSim struct {
	req *proto.RaidSimRequest
	cl  *raidSimRequestChangeLog
	eq  *equipmentSubstitution
}

func (b *bulkSimRunner) Run(pctx context.Context, progress chan *proto.ProgressMetrics) (result *proto.BulkSimResult, resultErr error) {
	ctx, cancel := context.WithCancel(pctx)
	defer func() {
		if err := recover(); err != nil {
			result = &proto.BulkSimResult{
				ErrorResult: fmt.Sprintf("%v\nStack Trace:\n%s", err, string(debug.Stack())),
			}
		}
		cancel()
	}()

	// Bulk simming is only supported for the single-player use (i.e. not whole raid-wide simming).
	// Verify that we have exactly 1 player.
	var playerCount int
	var player *proto.Player
	for _, p := range b.Request.GetBaseSettings().GetRaid().GetParties() {
		for _, pl := range p.GetPlayers() {
			// TODO(Riotdog-GehennasEU): Better way to check if a player is valid/set?
			if pl.Name != "" {
				player = pl
				playerCount++
			}
		}
	}
	if playerCount != 1 || player == nil {
		return nil, fmt.Errorf("bulksim: expected exactly 1 player, found %d", playerCount)
	}

	iterations := b.Request.GetBulkSettings().GetIterationsPerCombo()
	if iterations <= 0 {
		iterations = defaultIterationsPerCombo
	}

	items := b.Request.GetBulkSettings().GetItems()
	// numItems := len(items)
	// if b.Request.BulkSettings.Combinations && numItems > maxItemCount {
	// 	return nil, fmt.Errorf("too many items specified (%d > %d), not computationally feasible", numItems, maxItemCount)
	// }

	// Create all distinct combinations of (item, slot). For example, let's say the only item we
	// want to bulk sim is a one-handed item that can be worn both as an off-hand or a main-hand weapon.
	// For each slot, we will create one itemWithSlot pair, so (item, off-hand) and (item, main-hand).
	// We verify later that we are not emitting any invalid equipment set.
	if player.GetDatabase() != nil {
		addToDatabase(player.GetDatabase())
	}

	var distinctItemSlotCombos []*itemWithSlot
	for index, is := range items {
		item, ok := ItemsByID[is.Id]
		if !ok {
			return nil, fmt.Errorf("unknown item with id %d in bulk settings", is.Id)
		}
		for _, slot := range eligibleSlotsForItem(item) {
			distinctItemSlotCombos = append(distinctItemSlotCombos, &itemWithSlot{
				Item:  is,
				Slot:  ItemSlot(slot),
				Index: index,
			})
		}
	}
	baseItems := player.Equipment.Items

	allCombos := generateAllEquipmentSubstitutions(ctx, baseItems, b.Request.BulkSettings.Combinations, distinctItemSlotCombos)

	validCombos := []singleBulkSim{}
	for sub := range allCombos {
		substitutedRequest, changeLog := createNewRequestWithSubstitution(b.Request.BaseSettings, sub)
		if isValidEquipment(substitutedRequest.Raid.Parties[0].Players[0].Equipment) {
			validCombos = append(validCombos, singleBulkSim{req: substitutedRequest, cl: changeLog, eq: sub})
		}
	}

	// TODO(Riotdog-GehennasEU): Make this configurable?
	maxResults := 30

	var rankedResults []*itemSubstitutionSimResult
	var baseResult *itemSubstitutionSimResult
	newIters := int64(iterations)
	if b.Request.BulkSettings.FastMode {
		newIters /= 100

		// In fast mode try to keep starting iterations between 50 and 1000.
		if newIters < 50 {
			newIters = 50
		}
		if newIters > 1000 {
			newIters = 1000
		}
	}

	maxIterations := newIters * int64(len(validCombos))
	if maxIterations > math.MaxInt32 {
		return nil, fmt.Errorf("number of total iterations %d too large", maxIterations)
	}

	for {
		var tempBase *itemSubstitutionSimResult
		var err error
		// TODO: we could theoretically make getRankedResults accept a channel of validCombos that stream in to it and launches sims as it gets them...
		rankedResults, tempBase, err = b.getRankedResults(ctx, validCombos, newIters, progress)

		if err != nil {
			return nil, err
		}
		// keep replacing the base result with more refined base until we don't have base in the ranked results any more.
		if tempBase != nil {
			baseResult = tempBase
		}

		if !b.Request.BulkSettings.FastMode || len(rankedResults) <= maxResults {
			break
		}

		// we have reached max accuracy now
		if newIters >= int64(iterations) {
			break
		}

		// Increase accuracy
		newIters *= 2
		newNumCombos := len(rankedResults) / 2
		validCombos = validCombos[:newNumCombos]
		rankedResults = rankedResults[:newNumCombos]
		for i, comb := range rankedResults {
			validCombos[i] = singleBulkSim{
				req: comb.Request,
				cl:  comb.ChangeLog,
				eq:  comb.Substitution,
			}
		}
	}

	if baseResult == nil {
		return nil, fmt.Errorf("no base result for equipped gear found in bulk sim")
	}

	if len(rankedResults) > maxResults {
		rankedResults = rankedResults[:maxResults]
	}

	bum := baseResult.Result.GetRaidMetrics().GetParties()[0].GetPlayers()[0]
	bum.Actions = nil
	bum.Auras = nil
	bum.Resources = nil
	bum.Pets = nil

	result = &proto.BulkSimResult{
		EquippedGearResult: &proto.BulkComboResult{
			UnitMetrics: bum,
		},
	}

	for _, r := range rankedResults {
		um := r.Result.GetRaidMetrics().GetParties()[0].GetPlayers()[0]
		um.Actions = nil
		um.Auras = nil
		um.Resources = nil
		um.Pets = nil

		result.Results = append(result.Results, &proto.BulkComboResult{
			ItemsAdded:  r.ChangeLog.AddedItems,
			UnitMetrics: um,
		})
	}

	if progress != nil {
		progress <- &proto.ProgressMetrics{
			FinalBulkResult: result,
		}
	}

	return result, nil
}

func (b *bulkSimRunner) getRankedResults(pctx context.Context, validCombos []singleBulkSim, iterations int64, progress chan *proto.ProgressMetrics) ([]*itemSubstitutionSimResult, *itemSubstitutionSimResult, error) {
	concurrency := (runtime.NumCPU() - 1) * 2
	if concurrency <= 0 {
		concurrency = 2
	}

	tickets := make(chan struct{}, concurrency)
	for i := 0; i < concurrency; i++ {
		tickets <- struct{}{}
	}

	results := make(chan *itemSubstitutionSimResult, 10)

	numCombinations := int32(len(validCombos))
	totalIterationsUpperBound := int64(numCombinations) * iterations

	var totalCompletedIterations int32
	var totalCompletedSims int32

	ctx, cancel := context.WithCancel(pctx)
	// reporter for all sims combined.
	go func() {
		for ctx.Err() == nil {
			complIters := atomic.LoadInt32(&totalCompletedIterations)
			complSims := atomic.LoadInt32(&totalCompletedSims)

			// stop reporting
			if complIters == int32(totalIterationsUpperBound) || numCombinations == complSims {
				return
			}

			progress <- &proto.ProgressMetrics{
				TotalSims:           numCombinations,
				CompletedSims:       complSims,
				CompletedIterations: complIters,
				TotalIterations:     int32(totalIterationsUpperBound),
			}
			time.Sleep(time.Second)
		}
	}()

	// launcher for all combos (limited by concurrency max)
	go func() {
		for _, singleCombo := range validCombos {
			<-tickets
			singleSimProgress := make(chan *proto.ProgressMetrics)
			// watches this progress and pushes up to main reporter.
			go func(prog chan *proto.ProgressMetrics) {
				var prevDone int32
				for p := range singleSimProgress {
					delta := p.CompletedIterations - prevDone
					atomic.AddInt32(&totalCompletedIterations, delta)
					prevDone = p.CompletedIterations
					if p.FinalRaidResult != nil {
						break
					}
				}
			}(singleSimProgress)
			// actually run the sim in here.
			go func(sub singleBulkSim) {
				// overwrite the requests iterations with the input for this function.
				sub.req.SimOptions.Iterations = int32(iterations)
				results <- &itemSubstitutionSimResult{
					Request:      sub.req,
					Result:       b.SingleRaidSimRunner(sub.req, singleSimProgress, false),
					Substitution: sub.eq,
					ChangeLog:    sub.cl,
				}
				atomic.AddInt32(&totalCompletedSims, 1)
				tickets <- struct{}{} // when done, allow for new sim to be launched.
			}(singleCombo)
		}
	}()

	rankedResults := make([]*itemSubstitutionSimResult, numCombinations)
	var baseResult *itemSubstitutionSimResult

	for i := range rankedResults {
		result := <-results
		if result.Result == nil || result.Result.ErrorResult != "" {
			cancel() // cancel reporter
			return nil, nil, errors.New("simulation failed: " + result.Result.ErrorResult)
		}
		if !result.Substitution.HasItemReplacements() {
			baseResult = result
		}
		rankedResults[i] = result
	}
	cancel() // cancel reporter

	sort.Slice(rankedResults, func(i, j int) bool {
		return rankedResults[i].Score() > rankedResults[j].Score()
	})
	return rankedResults, baseResult, nil
}

// itemSubstitutionSimResult stores the request and response of a simulation, along with the used
// equipment susbstitution and a changelog of which items were added and removed from the base
// equipment set.
type itemSubstitutionSimResult struct {
	Request      *proto.RaidSimRequest
	Result       *proto.RaidSimResult
	Substitution *equipmentSubstitution
	ChangeLog    *raidSimRequestChangeLog
}

// Score used to rank results.
func (r *itemSubstitutionSimResult) Score() float64 {
	if r.Result == nil || r.Result.ErrorResult != "" {
		return 0
	}
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

	if equipment.Items[ItemSlotFinger1].Id == equipment.Items[ItemSlotFinger2].Id {
		return false
	} else if equipment.Items[ItemSlotTrinket1].Id == equipment.Items[ItemSlotTrinket2].Id {
		return false
	}

	return !(usesTwoHander && usesOffhand)
}

// generateAllEquipmentSubstitutions generates all possible valid equipment substitutions for the
// given bulk sim request. Also returns the unchanged equipment ("base equipment set") set as the
// first result. This ensures that simming over all possible equipment substitutions includes the
// base case as well.
func generateAllEquipmentSubstitutions(ctx context.Context, baseItems []*proto.ItemSpec, combinations bool, distinctItemSlotCombos []*itemWithSlot) chan *equipmentSubstitution {

	results := make(chan *equipmentSubstitution)
	go func() {
		defer close(results)

		// No substitutions (base case).
		results <- &equipmentSubstitution{}

		// seenCombos lets us deduplicate trinket/ring combos.
		comboChecker := ItemComboChecker{}

		// Pre-seed the existing item combos
		comboChecker.HasCombo(baseItems[ItemSlotFinger1].Id, baseItems[ItemSlotFinger2].Id)
		comboChecker.HasCombo(baseItems[ItemSlotTrinket1].Id, baseItems[ItemSlotTrinket2].Id)

		// Organize everything by slot.
		itemsBySlot := make([][]*proto.ItemSpec, 17)
		for _, is := range distinctItemSlotCombos {
			itemsBySlot[is.Slot] = append(itemsBySlot[is.Slot], is.Item)
		}

		if !combinations {
			for slotid, slot := range itemsBySlot {
				for _, item := range slot {
					sub := equipmentSubstitution{
						Items: []*itemWithSlot{{Item: item, Slot: ItemSlot(slotid)}},
					}
					// Handle finger/trinket specially to generate combos
					switch slotid {
					case int(ItemSlotFinger1), int(ItemSlotTrinket1):
						if !comboChecker.HasCombo(item.Id, baseItems[slotid+1].Id) {
							results <- &sub
						}
						// Generate extra combos
						subslot := slotid + 1
						for _, subitem := range itemsBySlot[subslot] {
							if shouldSkipCombo(baseItems, subitem, ItemSlot(subslot), comboChecker, sub) {
								continue
							}
							miniCombo := createReplacement(sub, &itemWithSlot{Item: subitem, Slot: ItemSlot(subslot)})
							results <- &miniCombo
						}
					case int(ItemSlotFinger2), int(ItemSlotTrinket2):
						// Ensure we don't have this combo with the base equipment.
						if !comboChecker.HasCombo(item.Id, baseItems[slotid-1].Id) {
							results <- &sub
						}
					default:
						results <- &sub
					}
				}
			}
			return
		}

		// Now generate combos by slot
		for i := 0; i < len(itemsBySlot); i++ {
			if len(itemsBySlot[i]) == 0 {
				continue
			}
			genSlotCombos(ItemSlot(i), baseItems, equipmentSubstitution{}, itemsBySlot, comboChecker, results)
		}
	}()

	return results
}

func createReplacement(repl equipmentSubstitution, item *itemWithSlot) equipmentSubstitution {
	newItems := make([]*itemWithSlot, len(repl.Items))
	copy(newItems, repl.Items)
	newItems = append(newItems, item)
	repl.Items = newItems
	return repl
}

func shouldSkipCombo(baseItems []*proto.ItemSpec, item *proto.ItemSpec, slot ItemSlot, comboChecker ItemComboChecker, replacements equipmentSubstitution) bool {
	switch slot {
	case ItemSlotFinger1, ItemSlotTrinket1:
		return comboChecker.HasCombo(item.Id, baseItems[slot+1].Id)
	case ItemSlotFinger2, ItemSlotTrinket2:

		for _, repl := range replacements.Items {
			if slot == ItemSlotFinger2 && repl.Slot == ItemSlotFinger1 ||
				slot == ItemSlotTrinket2 && repl.Slot == ItemSlotTrinket1 {
				return comboChecker.HasCombo(repl.Item.Id, item.Id)
			}
		}
		// Since we didn't find an item in the opposite slot, check against base items.
		return comboChecker.HasCombo(item.Id, baseItems[slot-1].Id)
	}
	return false
}

func genSlotCombos(slot ItemSlot, baseItems []*proto.ItemSpec, baseRepl equipmentSubstitution, replaceBySlot [][]*proto.ItemSpec, comboChecker ItemComboChecker, results chan *equipmentSubstitution) {
	// iterate all items in this slot, add to the baseRepl, then descend to add all other item combos.
	for _, item := range replaceBySlot[slot] {

		// Make sure we don't generate invalid or duplicate ring/trinket combos
		if slot >= ItemSlotFinger1 && slot <= ItemSlotTrinket2 {
			if shouldSkipCombo(baseItems, item, slot, comboChecker, baseRepl) {
				continue
			}
		}

		combo := createReplacement(baseRepl, &itemWithSlot{Slot: slot, Item: item})
		results <- &combo

		// Now descend to each other slot to pair with this combo
		for j := slot + 1; int(j) < len(replaceBySlot); j++ {
			if len(replaceBySlot[j]) == 0 {
				continue
			}
			genSlotCombos(j, baseItems, combo, replaceBySlot, comboChecker, results)
		}
	}
}

// itemWithSlot pairs an item with its fixed item slot.
type itemWithSlot struct {
	Item *proto.ItemSpec
	Slot ItemSlot

	// This index refers to the item's position in the BulkEquipmentSpec of the player and serves as
	// a unique item ID. It is used to verify that a valid equipmentSubstitution only references an
	// item once.
	Index int
}

// raidSimRequestChangeLog stores a change log of which items were added and removed from the base
// equipment set.
type raidSimRequestChangeLog struct {
	AddedItems []*proto.ItemSpecWithSlot
}

// createNewRequestWithSubstitution creates a copy of the input RaidSimRequest and applis the given
// equipment susbstitution to the player's equipment.
func createNewRequestWithSubstitution(readonlyInputRequest *proto.RaidSimRequest, substitution *equipmentSubstitution) (*proto.RaidSimRequest, *raidSimRequestChangeLog) {
	request := goproto.Clone(readonlyInputRequest).(*proto.RaidSimRequest)
	changeLog := &raidSimRequestChangeLog{}
	player := request.Raid.Parties[0].Players[0]
	equipment := player.Equipment
	for _, is := range substitution.Items {
		changeLog.AddedItems = append(changeLog.AddedItems, &proto.ItemSpecWithSlot{
			Item: is.Item,
			Slot: proto.ItemSlot(is.Slot),
		})
		equipment.Items[is.Slot] = is.Item
	}
	return request, changeLog
}

type ItemComboChecker map[int64]struct{}

func (ic *ItemComboChecker) HasCombo(itema int32, itemb int32) bool {
	if itema == itemb {
		return true
	}
	key := ic.generateComboKey(itema, itemb)
	if _, ok := (*ic)[key]; ok {
		return true
	} else {
		(*ic)[key] = struct{}{}
	}
	return false
}

// put this function on ic just so it isn't in global namespace
func (ic *ItemComboChecker) generateComboKey(itemA int32, itemB int32) int64 {
	if itemA > itemB {
		return int64(itemA) + int64(itemB)<<4
	}
	return int64(itemB) + int64(itemA)<<4
}
