package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/wowsims/wotlk/cmd/wowsimcli/bulk"
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
	replacefile := flag.String("replace", "replace.json", "location of replacement items file. Writes a CSV result of the items replaced instead of JSON")
	outfile := flag.String("output", "", "location of output file, defaults to stdout")
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

	var output []byte
	if *replacefile != "" {
		output = []byte(bulk.Sim(input, *replacefile, *verbose))
	} else {
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

		output, err = protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(finalResult)
		if err != nil {
			log.Fatalf("failed to marshal final results: %s", err)
		}
	}

	if *outfile == "" {
		print(string(output))
	} else {
		err = os.WriteFile(*outfile, output, 0666)
		if err != nil {
			log.Fatalf("failed to write output file:: %s", err)
		}
		if *verbose {
			fmt.Printf("Wrote output file: `%s` successfully.\n", *outfile)
		}
	}
}
