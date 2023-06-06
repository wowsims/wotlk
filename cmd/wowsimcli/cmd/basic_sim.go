package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

var simCmd = &cobra.Command{
	Use:   "sim",
	Short: "simulate items & settings",
	Run:   simMain,
}

func init() {
	simCmd.Flags().StringVar(&infile, "infile", "input.json", "location of input file (RaidSimRequest in protojson format)")
	simCmd.Flags().StringVar(&outfile, "outfile", "", "location of output file, defaults to stdout")
	simCmd.Flags().BoolVar(&verbose, "verbose", false, "print information during runtime")
	simCmd.MarkFlagRequired("infile")
}

func simMain(cmd *cobra.Command, args []string) {
	data, err := os.ReadFile(infile)
	if err != nil {
		log.Fatalf("failed to load input json file %q: %v", infile, err)
	}
	input := &proto.RaidSimRequest{}

	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, input)
	if err != nil {
		log.Fatalf("failed to load input json file: %s", err)
	}

	var output []byte
	reporter := make(chan *proto.ProgressMetrics, 10)
	core.RunRaidSimAsync(input, reporter)

	var finalResult *proto.RaidSimResult
	for v := range reporter {
		if v.FinalRaidResult != nil {
			finalResult = v.FinalRaidResult
			break
		}
		if verbose {
			fmt.Printf("Sim Progress: %d / %d\n", v.CompletedIterations, v.TotalIterations)
		}
	}

	output, err = protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(finalResult)
	if err != nil {
		log.Fatalf("failed to marshal final results: %s", err)
	}

	if outfile == "" {
		fmt.Print(string(output))
	} else {
		err = os.WriteFile(outfile, output, 0666)
		if err != nil {
			log.Fatalf("failed to write output file:: %s", err)
		}
		if verbose {
			fmt.Printf("Wrote output file: `%s` successfully.\n", outfile)
		}
	}
}
