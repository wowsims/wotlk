package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/wowsims/wotlk/sim"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

func init() {
	sim.RegisterAll()
}

func main() {
	infile := flag.String("input", "input.json", "location of input file")
	//outfile := flag.String("output", "output.json", "location of output file")
	//verbose := flag.Bool("verbose", false, "print information during runtime")

	flag.Parse()

	data, err := os.ReadFile(*infile)
	if err != nil {
		log.Fatalf("failed to load input json file: %s", err)
	}
	input := &proto.RaidSimRequest{}
	err = protojson.Unmarshal(data, input)
	if err != nil {
		log.Fatalf("failed to load input json file: %s", err)
	}

	dk, _ := input.Raid.Parties[0].Players[0].Spec.(*proto.Player_Deathknight)
	dk.Deathknight.Rotation.NeatGenome = "in 0\nin 1\nin 2\nin 3\nin 4\nin 5\nin 6\nin 7\nin 8\nhidden 9\nhidden 10\nout 11\nout 12\nout 13\n 14\nout 15\n 16\n	out 17\nout 18\nconnection 2 9 0.333 t 0\nconnection 9 10 0.8 t 1\nconnection 3 9 0.4 t 2\nconnection 9 11 0.4 t 3"

	result := core.RunRaidSim(input)

	if result != nil {
		fmt.Printf("Avg Dps: %f\n", result.RaidMetrics.Dps.Avg)
	}

	//output, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(finalResult)
	//if err != nil {
	//	log.Fatalf("faield to marshal final results: %s", err)
	//}
	//
	//err = os.WriteFile(*outfile, output, 0666)
	//if err != nil {
	//	log.Fatalf("failed to write output file:: %s", err)
	//}
	//if *verbose {
	//	fmt.Printf("Wrote output file: `%s` successfully.\n", *outfile)
	//}
}
