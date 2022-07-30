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
	outfile := flag.String("output", "output.json", "location of output file")
	verbose := flag.Bool("verbose", false, "print information during runtime")

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

	reporter := make(chan *proto.ProgressMetrics, 10)
	core.RunRaidSimAsync(input, reporter)

	var finalResult *proto.RaidSimResult
	for v := range reporter {
		if v.FinalRaidResult != nil {
			finalResult = v.FinalRaidResult
			break
		}
		if *verbose {
			fmt.Printf("Sim Progress: %d / %d\n", v.CompletedIterations, v.TotalIterations)
		}
	}

	output, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(finalResult)
	if err != nil {
		log.Fatalf("faield to marshal final results: %s", err)
	}

	err = os.WriteFile(*outfile, output, 0666)
	if err != nil {
		log.Fatalf("failed to write output file:: %s", err)
	}
	if *verbose {
		fmt.Printf("Wrote output file: `%s` successfully.\n", *outfile)
	}
}
