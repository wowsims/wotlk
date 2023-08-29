package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	infile      string
	replacefile string
	outfile     string
	verbose     bool
)

var bulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "bulk simulate item replacements and combinations",
	Long:  "bulk simulate item replacements and combinations",
	Run:   bulkSimMain,
}

func init() {
	bulkCmd.Flags().StringVar(&infile, "infile", "input.json", "location of input file (RaidSimRequest in protojson format)")
	bulkCmd.Flags().StringVar(&replacefile, "replacefile", "", "location of replacement items file. Writes a CSV result of the items replaced instead of JSON")
	bulkCmd.Flags().StringVar(&outfile, "output", "", "location of output file, defaults to stdout")
	bulkCmd.Flags().BoolVar(&verbose, "verbose", false, "print information during runtime")
	bulkCmd.MarkFlagRequired("infile")
	bulkCmd.MarkFlagRequired("replacefile")
}

func bulkSimMain(cmd *cobra.Command, args []string) {
	data, err := os.ReadFile(infile)
	if err != nil {
		log.Fatalf("failed to load input json file %q: %v", infile, err)
	}
	input := &proto.RaidSimRequest{}

	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, input)
	if err != nil {
		log.Fatalf("failed to load input json file: %s", err)
	}

	output := BulkSim(input, replacefile, verbose)

	if outfile == "" {
		print(string(output))
	} else {
		err = os.WriteFile(outfile, []byte(output), 0666)
		if err != nil {
			log.Fatalf("failed to write output file:: %s", err)
		}
		if verbose {
			fmt.Printf("Wrote output file: `%s` successfully.\n", outfile)
		}
	}
}

type ItemReplacementInput struct {
	Combinations bool              `json:"combinations"`
	FastMode     bool              `json:"fast_mode"`
	Items        []*proto.ItemSpec // spec for replacement
}

type ReplaceIter struct {
	Items []proto.ItemSpec // List of items to substitute
	Slots []proto.ItemSlot // Slots for each sub item
}

func BulkSim(input *proto.RaidSimRequest, replaceFile string, verbose bool) string {
	// 1. Load up all the sim data we need
	replaceData, err := os.ReadFile(replaceFile)
	if err != nil {
		log.Fatalf("failed to load replace json file: %s", err)
	}
	replaceInput := &ItemReplacementInput{}
	err = json.Unmarshal(replaceData, replaceInput)
	if err != nil {
		log.Fatalf("failed to parse replace json file: %s", err)
	}

	bsr := &proto.BulkSimRequest{
		BaseSettings: input,
		BulkSettings: &proto.BulkSettings{
			Combinations:       replaceInput.Combinations,
			Items:              replaceInput.Items,
			IterationsPerCombo: input.SimOptions.Iterations,
			FastMode:           replaceInput.FastMode,
		},
	}
	progress := make(chan *proto.ProgressMetrics, 100)
	core.RunBulkSimAsync(context.Background(), bsr, progress)

	startTime := time.Now()

	c := time.After(time.Minute)
	var lastTotal int32
	for {
		select {
		case status, ok := <-progress:
			if !ok {
				return ""
			}
			if status.FinalBulkResult != nil {
				if status.FinalBulkResult.ErrorResult != "" {
					fmt.Printf("Failed: %s\n", status.FinalBulkResult.ErrorResult)
				} else {
					return printCombos(status.FinalBulkResult)
				}
			}

			if verbose {
				if lastTotal != status.TotalSims {
					if lastTotal > status.TotalSims {
						fmt.Printf("Refining results, running the best combos with more iterations...\n")
					}
					lastTotal = status.TotalSims
				}
				compl := status.CompletedIterations
				if compl == 0 {
					break
				}
				elapsed := time.Since(startTime)
				perDone := float64(compl) / float64(status.TotalIterations)
				totalTime := time.Duration(float64(elapsed) / perDone)
				var timeEst string
				if totalTime.Hours() > 48 {
					// use days
					timeEst = fmt.Sprintf("Estimated Time: %0.1f / %0.1f days", elapsed.Hours()/24, totalTime.Hours()/24)
				} else if totalTime.Minutes() > 120 {
					// use hours
					timeEst = fmt.Sprintf("Estimated Time: %0.1f / %0.1f hours", elapsed.Hours(), totalTime.Hours())
				} else {
					timeEst = fmt.Sprintf("Estimated Time: %0.1f / %0.1f minutes", elapsed.Minutes(), totalTime.Minutes())
				}
				totalStr := strconv.Itoa(int(status.TotalIterations))
				fmtStr := "%" + strconv.Itoa(len(totalStr)) + ".f"
				fmt.Printf("Sim Progress: "+fmtStr+" / %d | %s  (completed %d / %d)\n", float64(compl), status.TotalIterations, timeEst, status.CompletedSims, status.TotalSims)
			}
		case <-c:

		}
	}
}

func printCombos(results *proto.BulkSimResult) string {
	result := ""
	foundBase := false
	for i := 0; i < len(results.Results); i++ {
		if len(results.Results[i].ItemsAdded) == 0 {
			foundBase = true
		}
		result += printCombo(results.Results[i])
	}
	if !foundBase {
		result += fmt.Sprintf("[BASE RESULT],%0.1f\n", results.EquippedGearResult.UnitMetrics.Dps.Avg)
	}
	return result
}

func printCombo(combo *proto.BulkComboResult) string {
	itemtext := "["
	if len(combo.ItemsAdded) == 0 {
		itemtext += "BASE RESULT"
	}
	for j, item := range combo.ItemsAdded {
		if j != 0 {
			itemtext += ";"
		}
		itemtext += fmt.Sprintf("%s@%s", core.ItemsByID[item.Item.Id].Name, item.Slot.String())
	}
	itemtext += "]"
	return fmt.Sprintf("%s,%0.1f\n", itemtext, combo.UnitMetrics.Dps.Avg)
}
