package bulk

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

func Sim(input *proto.RaidSimRequest, replaceFile string) {
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
		runComboReplacements(input, replaceBySlots)
		return
	}
	runSingleReplacements(input, replaceInput)
}

func runSingleReplacements(input *proto.RaidSimRequest, replaceInput *ItemReplacementInput) {
	results := make([]*proto.RaidSimResult, len(replaceInput.Items)+1)

	waits := &sync.WaitGroup{}
	var totalIters int64 = int64(input.SimOptions.Iterations) * int64(len(results))
	var completedIters int64

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

	go func() {
		for i := 0; i < 100; i++ {
			fmt.Printf("Sim Progress: %d / %d\n", atomic.LoadInt64(&completedIters), totalIters)
			time.Sleep(time.Second * 2)
		}
	}()
	waits.Wait()

	fmt.Printf("\nBASE RESULT,,,%0.1f\n", results[0].RaidMetrics.Dps.Avg)
	for i := 1; i < len(results); i++ {
		if results[i].ErrorResult != "" {
			fmt.Printf("ERROR RESULT: %s\n", results[i].ErrorResult)
		}
		slot := replaceInput.replaceSlots[i-1]
		oldItemID := input.Raid.Parties[0].Players[0].Equipment.Items[slot].Id
		fmt.Printf("%s,%s,%s,%0.1f\n", slot.String(), core.ItemsByID[oldItemID].Name, core.ItemsByID[replaceInput.Items[i-1].Id].Name, results[i].RaidMetrics.Dps.Avg)
	}
}

func generateCombos(offset int, baseRepl ReplaceIter, replaceBySlot [][]proto.ItemSpec) []ReplaceIter {
	// TODO STILL
	//  1. Handle items with multiple slots?
	//   Or pre-process those into each possible slot.
	//		Fingers
	//		Trinkets
	//		Hands
	//  2. Handle MH / OH items mixed with 2H
	//		Don't generate combos with a 2H and OH together.

	result := []ReplaceIter{}

	for i := offset; i < len(replaceBySlot); i++ {
		for _, item := range replaceBySlot[i] {
			newRepl := baseRepl

			newItems := make([]proto.ItemSpec, len(newRepl.Items))
			copy(newItems, baseRepl.Items)
			newSlots := make([]core.ItemSlot, len(newRepl.Slots))
			copy(newSlots, baseRepl.Slots)

			newItems = append(newItems, item)
			newSlots = append(newSlots, core.ItemSlot(i))

			newRepl.Items = newItems
			newRepl.Slots = newSlots

			result = append(result, newRepl)

			for j := i + 1; j < len(replaceBySlot); j++ {
				if len(replaceBySlot[j]) == 0 {
					continue
				}
				result = append(result, generateCombos(j, newRepl, replaceBySlot)...)
			}
		}
	}

	return result
}

func runComboReplacements(baseInput *proto.RaidSimRequest, replaceBySlot [][]proto.ItemSpec) {
	totalCombo := 1
	for _, v := range replaceBySlot {
		if len(v) == 0 {
			continue
		}
		totalCombo *= len(v)
	}
	fmt.Printf("\nTotal Combinations: %d\n", totalCombo)

	combos := generateCombos(0, ReplaceIter{}, replaceBySlot)

	results := make([]*proto.RaidSimResult, len(combos)+1)

	waits := &sync.WaitGroup{}
	var totalIters int64 = int64(baseInput.SimOptions.Iterations) * int64(len(results))
	var completedIters int64

	for i := 0; i < len(results); i++ {
		waits.Add(1)
		go func(iter int) {
			newInput := *baseInput
			if iter > 0 {
				newRaid := *baseInput.Raid
				newParty := *newRaid.Parties[0]
				newPlayer := *newParty.Players[0]
				newEquip := *newPlayer.Equipment
				newItems := make([]*proto.ItemSpec, len(newEquip.Items))
				copy(newItems, newEquip.Items)

				for i, repl := range combos[iter-1].Items {
					replCpy := repl
					newItemInfo := core.ItemsByID[replCpy.Id]
					slot := combos[iter-1].Slots[i]
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

	go func() {
		for {
			compl := atomic.LoadInt64(&completedIters)
			fmt.Printf("Sim Progress: %d / %d\n", compl, totalIters)
			if completedIters == totalIters {
				break
			}
			time.Sleep(time.Second * 2)
		}
	}()

	waits.Wait()
	// baseOutput, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(results[0])
	// if err != nil {
	// 	log.Fatalf("failed to marshal final results: %s", err)
	// }

	// compareOutput, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(results[0])
	// if err != nil {
	// 	log.Fatalf("failed to marshal final results: %s", err)
	// }

	fmt.Printf("\n[BASE RESULT],%0.1f\n", results[0].RaidMetrics.Dps.Avg)
	for i := 1; i < len(results); i++ {
		if results[i].ErrorResult != "" {
			fmt.Printf("ERROR RESULT: %s\n", results[i].ErrorResult)
		}
		combo := combos[i-1]

		itemtext := "["
		for j, item := range combo.Items {
			if j != 0 {
				itemtext += ";"
			}
			slot := combo.Slots[j]
			itemtext += fmt.Sprintf("%s@%s", core.ItemsByID[item.Id].Name, slot.String())
		}
		itemtext += "]"
		fmt.Printf("%s,%0.1f\n", itemtext, results[i].RaidMetrics.Dps.Avg)
	}
}
