package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wowsims/wotlk/sim"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

func init() {
	sim.RegisterAll()
}

var (
	Version string
)

func main() {
	infile := flag.String("input", "input.json", "location of input file")
	// outfile := flag.String("output", "output.json", "location of output file")
	verbose := flag.Bool("verbose", false, "print information during runtime")
	printVersion := flag.Bool("version", false, "print version number and exit")

	flag.Parse()

	if Version == "" {
		Version = "development"
	}
	if *printVersion {
		fmt.Printf("Version: %s\n", Version)
		return
	}

	data, err := os.ReadFile(*infile)
	if err != nil {
		log.Fatalf("failed to load input json file: %s", err)
	}
	input := &proto.RaidSimRequest{}
	err = protojson.Unmarshal(data, input)
	if err != nil {
		log.Fatalf("failed to load input json file: %s", err)
	}

	type itemReplacement struct {
		proto.ItemSpec // spec for replacement
		core.ItemSlot  // slot to put it in
	}
	replacements := []itemReplacement{
		{ItemSpec: proto.ItemSpec{Id: 45559, Gems: []int32{39996, 40022}}, ItemSlot: core.ItemSlotFeet},
		{ItemSpec: proto.ItemSpec{Id: 45609}, ItemSlot: core.ItemSlotTrinket1},
		{ItemSpec: proto.ItemSpec{Id: 45609}, ItemSlot: core.ItemSlotTrinket2},
		{ItemSpec: proto.ItemSpec{Id: 45106}, ItemSlot: core.ItemSlotFinger1},
		{ItemSpec: proto.ItemSpec{Id: 45106}, ItemSlot: core.ItemSlotFinger2},
		{ItemSpec: proto.ItemSpec{Id: 45107, Gems: []int32{41398, 49110}}, ItemSlot: core.ItemSlotHead}, // meta + nightmare tear
		{ItemSpec: proto.ItemSpec{Id: 45134, Gems: []int32{39996, 39996, 39996}}, ItemSlot: core.ItemSlotLegs},
		{ItemSpec: proto.ItemSpec{Id: 45318}},
		{ItemSpec: proto.ItemSpec{Id: 45161, Gems: []int32{39996}}, ItemSlot: core.ItemSlotWaist},
		{ItemSpec: proto.ItemSpec{Id: 45157}, ItemSlot: core.ItemSlotFinger1},
		{ItemSpec: proto.ItemSpec{Id: 45157}, ItemSlot: core.ItemSlotFinger2},
		{ItemSpec: proto.ItemSpec{Id: 45298}, ItemSlot: core.ItemSlotMainHand},
		{ItemSpec: proto.ItemSpec{Id: 45298}, ItemSlot: core.ItemSlotOffHand},
		{ItemSpec: proto.ItemSpec{Id: 45299, Gems: []int32{41398, 49110}}, ItemSlot: core.ItemSlotHead}, // meta + nightmare tear
		{ItemSpec: proto.ItemSpec{Id: 45303}, ItemSlot: core.ItemSlotFinger1},
		{ItemSpec: proto.ItemSpec{Id: 45303}, ItemSlot: core.ItemSlotFinger2},
		{ItemSpec: proto.ItemSpec{Id: 45138}},
		{ItemSpec: proto.ItemSpec{Id: 45142}, ItemSlot: core.ItemSlotMainHand},
		{ItemSpec: proto.ItemSpec{Id: 45142}, ItemSlot: core.ItemSlotOffHand},
		{ItemSpec: proto.ItemSpec{Id: 45286}, ItemSlot: core.ItemSlotTrinket1},
		{ItemSpec: proto.ItemSpec{Id: 45286}, ItemSlot: core.ItemSlotTrinket2},
		{ItemSpec: proto.ItemSpec{Id: 45676, Gems: []int32{39996}}, ItemSlot: core.ItemSlotChest},
		{ItemSpec: proto.ItemSpec{Id: 45675}, ItemSlot: core.ItemSlotFinger1},
		{ItemSpec: proto.ItemSpec{Id: 45675}, ItemSlot: core.ItemSlotFinger2},
		{ItemSpec: proto.ItemSpec{Id: 45248, Gems: []int32{39996}}, ItemSlot: core.ItemSlotLegs},
		{ItemSpec: proto.ItemSpec{Id: 45250}, ItemSlot: core.ItemSlotFinger1},
		{ItemSpec: proto.ItemSpec{Id: 45250}, ItemSlot: core.ItemSlotFinger2},
		{ItemSpec: proto.ItemSpec{Id: 45254}, ItemSlot: core.ItemSlotRanged},
		{ItemSpec: proto.ItemSpec{Id: 45193}, ItemSlot: core.ItemSlotNeck},
		{ItemSpec: proto.ItemSpec{Id: 45224}},
		{ItemSpec: proto.ItemSpec{Id: 45225, Gems: []int32{39996}}, ItemSlot: core.ItemSlotChest},
	}
	results := make([]*proto.RaidSimResult, len(replacements)+1)

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
				repl := replacements[iter-1]
				if repl.ItemSlot == 0 {
					newItem := core.ItemsByID[repl.Id]
					repl.ItemSlot = core.ItemTypeToSlot(newItem.Type)
					replacements[iter-1].ItemSlot = repl.ItemSlot
				}
				if repl.Enchant == 0 {
					repl.Enchant = newItems[repl.ItemSlot].Enchant
					replacements[iter-1].Enchant = repl.Enchant
				}
				newItems[repl.ItemSlot] = &repl.ItemSpec
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

	if *verbose {
		go func() {
			for i := 0; i < 100; i++ {
				fmt.Printf("Sim Progress: %d / %d\n", atomic.LoadInt64(&completedIters), totalIters)
				time.Sleep(time.Second * 2)
			}
		}()
	}
	waits.Wait()
	// baseOutput, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(results[0])
	// if err != nil {
	// 	log.Fatalf("failed to marshal final results: %s", err)
	// }

	// compareOutput, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(results[0])
	// if err != nil {
	// 	log.Fatalf("failed to marshal final results: %s", err)
	// }

	fmt.Printf("\nBASE RESULT,,,%0.1f\n", results[0].RaidMetrics.Dps.Avg)
	for i := 1; i < len(results); i++ {
		if results[i].ErrorResult != "" {
			fmt.Printf("ERROR RESULT: %s\n", results[i].ErrorResult)
		}
		slot := replacements[i-1].ItemSlot
		oldItemID := input.Raid.Parties[0].Players[0].Equipment.Items[slot].Id
		fmt.Printf("%s,%s,%s,%0.1f\n", slot.String(), core.ItemsByID[oldItemID].Name, core.ItemsByID[replacements[i-1].Id].Name, results[i].RaidMetrics.Dps.Avg)
	}

	// err = os.WriteFile(*outfile, output, 0666)
	// if err != nil {
	// 	log.Fatalf("failed to write output file:: %s", err)
	// }
	// if *verbose {
	// 	fmt.Printf("Wrote output file: `%s` successfully.\n", *outfile)
	// }
}
