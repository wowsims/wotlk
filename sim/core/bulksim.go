package core

import (
	"context"
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

func (b *bulkSimRunner) Run(ctx context.Context, progress chan *proto.ProgressMetrics) (result *proto.BulkSimResult, resultErr error) {
	defer func() {
		if err := recover(); err != nil {
			result = &proto.BulkSimResult{
				ErrorResult: fmt.Sprintf("%v\nStack Trace:\n%s", err, string(debug.Stack())),
			}
		}
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
	numItems := len(items)
	if b.Request.BulkSettings.Combinations && numItems > maxItemCount {
		return nil, fmt.Errorf("too many items specified (%d > %d), not computationally feasible", numItems, maxItemCount)
	}

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
				Slot:  slot,
				Index: index,
			})
		}
	}

	allCombos := generateAllEquipmentSubstitutions(ctx, b.Request.BulkSettings.Combinations, distinctItemSlotCombos)

	validCombos := []singleBulkSim{}
	for sub := range allCombos {
		substitutedRequest, changeLog := createNewRequestWithSubstitution(b.Request.BaseSettings, sub)
		if isValidEquipment(substitutedRequest.Raid.Parties[0].Players[0].Equipment) {
			validCombos = append(validCombos, singleBulkSim{req: substitutedRequest, cl: changeLog, eq: sub})
		}
	}

	// TODO(Riotdog-GehennasEU): Make this configurable?
	maxResults := 20
	if !b.Request.BulkSettings.Combinations {
		maxResults = len(validCombos)
	}

	var rankedResults []*itemSubstitutionSimResult
	var baseResult *itemSubstitutionSimResult
	newIters := int64(iterations)
	if b.Request.BulkSettings.FastMode && b.Request.BulkSettings.Combinations {
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
		// TODO: we could theoretically make getRankedResults accept a channel of validCombos that stream in to it and launches sims as it gets them...
		rankedResults, tempBase = b.getRankedResults(validCombos, newIters, progress)

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

func (b *bulkSimRunner) getRankedResults(validCombos []singleBulkSim, iterations int64, progress chan *proto.ProgressMetrics) ([]*itemSubstitutionSimResult, *itemSubstitutionSimResult) {
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

	// reporter for all sims combined.
	go func() {
		for {
			time.Sleep(time.Second)
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
	for i := range rankedResults {
		rankedResults[i] = <-results
	}
	sort.Slice(rankedResults, func(i, j int) bool {
		return rankedResults[i].Score() > rankedResults[j].Score()
	})

	var baseResult *itemSubstitutionSimResult
	for _, r := range rankedResults {
		if !r.Substitution.HasItemReplacements() {
			baseResult = r
		}
	}

	return rankedResults, baseResult
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
func generateAllEquipmentSubstitutions(ctx context.Context, combinations bool, distinctItemSlotCombos []*itemWithSlot) chan *equipmentSubstitution {
	results := make(chan *equipmentSubstitution)
	go func() {
		defer close(results)

		// No substitutions (base case).
		results <- &equipmentSubstitution{}

		if !combinations {
			for _, slotItem := range distinctItemSlotCombos {
				results <- &equipmentSubstitution{
					Items: []*itemWithSlot{slotItem},
				}
			}
			return
		}
		// Borrowed from https://github.com/mxschmitt/golang-combinations and adapted to
		// only emit valid combinations.
		count := uint64(len(distinctItemSlotCombos))
		for bits := uint64(1); bits < (1 << count); bits++ {
			combo := &equipmentSubstitution{}
			for idx := uint64(0); idx < count; idx++ {
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
			Slot: is.Slot,
		})
		equipment.Items[is.Slot] = is.Item
	}
	return request, changeLog
}
