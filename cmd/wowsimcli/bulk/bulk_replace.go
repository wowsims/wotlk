package bulk

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

type ItemReplacementInput struct {
	Combinations bool
	Items        []proto.ItemSpec // spec for replacement
	replaceSlots []core.ItemSlot
}

type ReplaceIter struct {
	Items []proto.ItemSpec // List of items to substitute
	Slots []core.ItemSlot  // Slots for each sub item
}

func Sim(input *proto.RaidSimRequest, replaceFile string, verbose bool) string {
	replaceData, err := os.ReadFile(replaceFile)
	if err != nil {
		log.Fatalf("failed to load input json file: %s", err)
	}
	replaceInput := &ItemReplacementInput{}
	err = json.Unmarshal(replaceData, replaceInput)
	if err != nil {
		log.Fatalf("failed to load input json file: %s", err)
	}
	replaceInput.replaceSlots = make([]core.ItemSlot, len(replaceInput.Items))

	numInput := len(replaceInput.Items)
	for i := 0; i < numInput; i++ {
		repl := replaceInput.Items[i]
		newItem := core.ItemsByID[repl.Id]
		slot := core.ItemTypeToSlot(newItem.Type)

		if repl.Enchant == 0 {
			oldItem := input.Raid.Parties[0].Players[0].Equipment.Items[slot]
			replaceInput.Items[i].Enchant = oldItem.Enchant
		}

		replaceInput.replaceSlots[i] = slot

		var extraSlot core.ItemSlot
		switch slot {
		case core.ItemSlotFinger1:
			extraSlot = core.ItemSlotFinger2
		case core.ItemSlotTrinket1:
			extraSlot = core.ItemSlotTrinket2
		case core.ItemSlotMainHand:
			if newItem.HandType == proto.HandType_HandTypeOneHand {
				extraSlot = core.ItemSlotOffHand
			}
		}
		if extraSlot != 0 {
			replCopy := repl
			oldItem := input.Raid.Parties[0].Players[0].Equipment.Items[extraSlot]
			replCopy.Enchant = oldItem.Enchant
			replaceInput.Items = append(replaceInput.Items, replCopy)
			replaceInput.replaceSlots = append(replaceInput.replaceSlots, extraSlot)
		}
	}

	if replaceInput.Combinations {
		replaceBySlots := make([][]proto.ItemSpec, 17)
		for i, repl := range replaceInput.Items {
			slot := replaceInput.replaceSlots[i]
			replaceBySlots[slot] = append(replaceBySlots[slot], repl)
		}
		return runComboControl(input, replaceBySlots, verbose)
	}
	return runSingleReplacements(input, replaceInput, verbose)
}

func runSingleReplacements(input *proto.RaidSimRequest, replaceInput *ItemReplacementInput, verbose bool) string {
	results := make([]*proto.RaidSimResult, len(replaceInput.Items)+1)

	waits := &sync.WaitGroup{}
	var totalIters int64 = int64(input.SimOptions.Iterations) * int64(len(results))
	var completedIters int64

	startTime := time.Now()
	for i := 0; i < len(results); i++ {
		waits.Add(1)
		go func(iter int) {
			newInput := *input
			if iter > 0 {
				newRaid := *input.Raid
				newParty := *newRaid.Parties[0]
				newPlayer := *newParty.Players[0]
				newEquip := *newPlayer.Equipment
				newItems := make([]*proto.ItemSpec, len(newEquip.Items))
				copy(newItems, newEquip.Items)

				repl := replaceInput.Items[iter-1]
				newItemInfo := core.ItemsByID[repl.Id]

				slot := replaceInput.replaceSlots[iter-1]
				newItems[slot] = &repl

				// if new item is 2h and existing items are 1h/offhand remove the OH
				if slot == core.ItemSlotMainHand && newItemInfo.HandType == proto.HandType_HandTypeTwoHand && newItems[core.ItemSlotOffHand].Id != 0 {
					newItems[core.ItemSlotOffHand].Id = 0
				}

				newEquip.Items = newItems
				newPlayer.Equipment = &newEquip
				newParty.Players = []*proto.Player{&newPlayer}
				newRaid.Parties = []*proto.Party{&newParty}
				newInput.Raid = &newRaid
			}
			reporter := make(chan *proto.ProgressMetrics, 10)
			core.RunRaidSimAsync(&newInput, reporter)
			var finalResult *proto.RaidSimResult
			var lastComplete int64
			for v := range reporter {
				if v.FinalRaidResult != nil {
					finalResult = v.FinalRaidResult
					break
				}
				atomic.AddInt64(&completedIters, int64(v.CompletedIterations)-lastComplete)
				lastComplete = int64(v.CompletedIterations)
			}
			results[iter] = finalResult
			waits.Done()
		}(i)
	}

	if verbose {
		go func() {
			for {
				compl := atomic.LoadInt64(&completedIters)
				perDone := float64(compl) / float64(totalIters)
				elapsed := time.Since(startTime)
				totalTime := time.Duration(float64(elapsed) / perDone)
				fmt.Printf("Sim Progress: %d / %d | Estimated Time: %0.1f / %0.1f minutes\n", compl, totalIters, elapsed.Minutes(), totalTime.Minutes())
				if compl == totalIters {
					break
				}
				time.Sleep(time.Second * 2)
			}
		}()
	}
	waits.Wait()

	resultText := fmt.Sprintf("BASE RESULT,,,%0.1f\n", results[0].RaidMetrics.Dps.Avg)
	for i := 1; i < len(results); i++ {
		if results[i].ErrorResult != "" {
			panic("ERROR RESULT: " + results[i].ErrorResult)
		}
		slot := replaceInput.replaceSlots[i-1]
		oldItemID := input.Raid.Parties[0].Players[0].Equipment.Items[slot].Id
		resultText += fmt.Sprintf("%s,%s,%s,%0.1f\n", slot.String(), core.ItemsByID[oldItemID].Name, core.ItemsByID[replaceInput.Items[i-1].Id].Name, results[i].RaidMetrics.Dps.Avg)
	}
	return resultText
}

func generateCombos(offset int, baseRepl ReplaceIter, replaceBySlot [][]proto.ItemSpec) []ReplaceIter {
	result := []ReplaceIter{}

	var genOneSlotCombos func(slot int, specs []proto.ItemSpec, baseRepl ReplaceIter, replaceBySlot [][]proto.ItemSpec) []ReplaceIter

	genOneSlotCombos = func(slot int, specs []proto.ItemSpec, baseRepl ReplaceIter, replaceBySlot [][]proto.ItemSpec) []ReplaceIter {
		innerResult := []ReplaceIter{}
		for _, item := range specs {
			newRepl := baseRepl

			newItems := make([]proto.ItemSpec, len(newRepl.Items))
			copy(newItems, baseRepl.Items)
			newSlots := make([]core.ItemSlot, len(newRepl.Slots))
			copy(newSlots, baseRepl.Slots)

			newItems = append(newItems, item)
			newSlots = append(newSlots, core.ItemSlot(slot))

			newRepl.Items = newItems
			newRepl.Slots = newSlots

			innerResult = append(innerResult, newRepl)

			for j := slot + 1; j < len(replaceBySlot); j++ {
				if len(replaceBySlot[j]) == 0 {
					continue
				}
				innerResult = append(innerResult, genOneSlotCombos(j, replaceBySlot[j], newRepl, replaceBySlot)...)
			}
		}
		return innerResult
	}

	for i := offset; i < len(replaceBySlot); i++ {
		if len(replaceBySlot[i]) == 0 {
			continue
		}
		result = append(result, genOneSlotCombos(i, replaceBySlot[i], baseRepl, replaceBySlot)...)
	}

	return result
}

func runComboControl(baseInput *proto.RaidSimRequest, replaceBySlot [][]proto.ItemSpec, verbose bool) string {
	totalCombo := 1
	totalReplace := 0
	for _, v := range replaceBySlot {
		if len(v) == 0 {
			continue
		}
		totalCombo *= len(v) + 1
		totalReplace += len(v)
	}
	if verbose {
		for i, repl := range replaceBySlot {
			if len(repl) == 0 {
				continue
			}
			replStr := ""
			for _, v := range repl {
				replStr += "\thttps://wowhead.com/wotlk/item=" + strconv.Itoa(int(v.Id)) + "\n"
			}
			if verbose {
				fmt.Printf("Slot: %s -> Replacements:\n%s", core.ItemSlot(i).String(), replStr)
			}

		}
		if verbose {
			fmt.Printf("\nReplacements: %d, Total Combinations: %d\n", totalReplace, totalCombo)
		}
	}
	if totalCombo > 10000000 {
		panic("over a 10 million combinations. Unlikely this will even be able to run.")
	}

	combos := generateCombos(0, ReplaceIter{}, replaceBySlot)
	// Add base case (no replacements.)
	combos = append([]ReplaceIter{{}}, combos...)
	baseInput.SimOptions.Iterations = 50 // TODO: allow config of base iterations
	maxResults := 50                     // maximum number of combos to display. TODO: make configurable

	if verbose {
		fmt.Printf("Generated %d combos\n", len(combos))
	}
	var baseResult *proto.RaidSimResult
	var results []*proto.RaidSimResult
	for {
		results = runComboReplacements(baseInput, combos, verbose)
		if baseResult == nil {
			// First result of first combo run is the base one.
			baseResult = results[0]
		}
		if len(combos) <= maxResults {
			break
		}
		if verbose {
			fmt.Printf("Refining results...\n")
		}

		avg := 0.0
		for _, res := range results {
			avg += res.RaidMetrics.Dps.Avg
		}
		avg /= float64(len(results))

		newCombos := make([]ReplaceIter, 0, len(combos)/2)
		for i, res := range results {
			if res.ErrorResult != "" {
				panic("failed a simulation: " + res.ErrorResult)
			}
			if res.RaidMetrics.Dps.Avg > avg {
				newCombos = append(newCombos, combos[i])
			}
		}
		combos = newCombos
		baseInput.SimOptions.Iterations *= 2
	}

	return printCombos(baseResult, combos, results)
}

func printCombos(baseResult *proto.RaidSimResult, combos []ReplaceIter, results []*proto.RaidSimResult) string {
	result := fmt.Sprintf("[BASE RESULT],%0.1f\n", baseResult.RaidMetrics.Dps.Avg)

	for i := 1; i < len(results); i++ {
		combo := combos[i]
		result += printCombo(combo, results[i])
	}
	return result
}

func printCombo(combo ReplaceIter, result *proto.RaidSimResult) string {
	itemtext := "["
	for j, item := range combo.Items {
		if j != 0 {
			itemtext += ";"
		}
		slot := combo.Slots[j]
		itemtext += fmt.Sprintf("%s@%s", core.ItemsByID[item.Id].Name, slot.String())
	}
	itemtext += "]"
	return fmt.Sprintf("%s,%0.1f\n", itemtext, result.RaidMetrics.Dps.Avg)
}

func runComboReplacements(baseInput *proto.RaidSimRequest, combos []ReplaceIter, verbose bool) []*proto.RaidSimResult {
	if verbose {
		fmt.Printf("Running %d combinations %d iterations each.\n", len(combos), baseInput.SimOptions.Iterations)
	}

	results := make([]*proto.RaidSimResult, len(combos))

	waits := &sync.WaitGroup{}
	var totalIters int64 = int64(baseInput.SimOptions.Iterations) * int64(len(results))
	var completedIters int64

	maxParallel := runtime.NumCPU() * 2
	tickets := make(chan struct{}, maxParallel)
	for i := 0; i < maxParallel; i++ {
		tickets <- struct{}{}
	}

	waits.Add(1)
	startTime := time.Now()
	launched := false
	go func() {
		for i := 0; i < len(results); i++ {
			waits.Add(1)
			<-tickets
			if !launched {
				launched = true
				waits.Done()
			}
			go func(iter int) {
				newInput := *baseInput

				if len(combos[iter].Items) > 0 {
					newRaid := *baseInput.Raid
					newParty := *newRaid.Parties[0]
					newPlayer := *newParty.Players[0]
					newEquip := *newPlayer.Equipment
					newItems := make([]*proto.ItemSpec, len(newEquip.Items))
					copy(newItems, newEquip.Items)

					for i, repl := range combos[iter].Items {
						replCpy := repl
						newItemInfo := core.ItemsByID[replCpy.Id]
						slot := combos[iter].Slots[i]
						newItems[slot] = &replCpy

						// if new item is 2h and existing items are 1h/offhand remove the OH
						if slot == core.ItemSlotMainHand && newItemInfo.HandType == proto.HandType_HandTypeTwoHand && newItems[core.ItemSlotOffHand].Id != 0 {
							newItems[core.ItemSlotOffHand].Id = 0
						}
					}

					newEquip.Items = newItems
					newPlayer.Equipment = &newEquip
					newParty.Players = []*proto.Player{&newPlayer}
					newRaid.Parties = []*proto.Party{&newParty}
					newInput.Raid = &newRaid
				}

				reporter := make(chan *proto.ProgressMetrics, 10)
				core.RunRaidSimAsync(&newInput, reporter)
				var finalResult *proto.RaidSimResult
				var lastComplete int64
				for v := range reporter {
					atomic.AddInt64(&completedIters, int64(v.CompletedIterations)-lastComplete)
					lastComplete = int64(v.CompletedIterations)
					if v.FinalRaidResult != nil {
						finalResult = v.FinalRaidResult
						break
					}
				}
				results[iter] = finalResult
				tickets <- struct{}{}
				waits.Done()
			}(i)
		}
	}()

	if verbose {
		go func() {
			for {
				compl := atomic.LoadInt64(&completedIters)
				if compl == 0 {
					time.Sleep(time.Second)
					continue
				}
				if compl == totalIters {
					break
				}
				elapsed := time.Since(startTime)
				perDone := float64(compl) / float64(totalIters)
				totalTime := time.Duration(float64(elapsed) / perDone)
				fmt.Printf("Sim Progress: %d / %d | Estimated Time: %0.1f / %0.1f minutes\n", compl, totalIters, elapsed.Minutes(), totalTime.Minutes())
				time.Sleep(time.Second * 2)
			}
		}()
	}

	waits.Wait()

	return results
}
