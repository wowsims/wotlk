package bulk

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

type ItemReplacementInput struct {
	Combinations bool
	Items        []*proto.ItemSpec // spec for replacement
	replaceSlots []core.ItemSlot
}

type ReplaceIter struct {
	Items []proto.ItemSpec // List of items to substitute
	Slots []core.ItemSlot  // Slots for each sub item
}

func Sim(input *proto.RaidSimRequest, replaceFile string, verbose bool) string {
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
			FastMode:           replaceInput.Combinations,
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
	for i := 1; i < len(results.Results); i++ {
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
