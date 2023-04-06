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
					outputStr := ""
					return outputStr
				}
			}

			if verbose {
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
