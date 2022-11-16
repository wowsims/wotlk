package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	var maxParam = flag.Int("maxid", 52205, "maximum ID to scan for")
	var minParam = flag.Int("minid", 1, "ID of item to start scan at")
	var idFile = flag.String("ids", "", "file with list of IDs to fetch")
	var numThreads = flag.Int("threads", 8, "number of parallel workers to fetch tooltips with")
	var output = flag.String("output", "all_item_tooltips.csv", "name of file to output results to")

	flag.Parse()

	maxID := *maxParam
	startMinID := *minParam

	threads := *numThreads

	fmt.Printf("Starting tooltip fetching.\n\tids file: %s\n\tStart ID: %d, End ID: %d\n\tWorkers: %d\n\tOutput File: %s\n\n", *idFile, startMinID, maxID, threads, *output)
	f, err := os.Create(*output)
	if err != nil {
		log.Fatalf("failed to open file to write: %s", err)
	}
	type result struct {
		id    int
		value string
	}
	results := make(chan result, 10)
	wg := &sync.WaitGroup{}

	total := 0
	if *idFile == "" {
		numPer := (maxID - startMinID) / threads
		total = maxID - startMinID
		currentMinID := startMinID
		for crawler := 0; crawler < threads; crawler++ {
			wg.Add(1)
			if crawler == threads-1 {
				numPer = maxID - currentMinID // last thread gets the odd extras
			}
			go func(min, max int) {
				client := http.Client{}
				for i := min; i < max; i++ {
					url := fmt.Sprintf("https://nether.wowhead.com/wotlk/tooltip/item/%d", i)
					resp, err := client.Get(url)
					if err != nil {
						fmt.Printf("Error fetching %d: %s\n", i, err)
						continue
					}
					body, _ := io.ReadAll(resp.Body)
					bstr := string(body)
					// fmt.Printf("Found tooltip for %d\n%s\n", i, bstr)
					if len(bstr) < 2 {
						fmt.Printf("Missing tooltip data for %d", i)
					}
					results <- result{id: i, value: bstr}
				}
				wg.Done()
			}(currentMinID, currentMinID+numPer)
			fmt.Printf("Started thread %d, %d to %d\n", crawler, currentMinID, currentMinID+numPer-1)
			currentMinID += numPer
		}
	} else {
		fmt.Printf("Fetching ids from file list...\n")
		ids := getItemDeclarations(*idFile)
		fmt.Printf("Found %d items...\n", len(ids))
		numPer := len(ids) / threads
		total = len(ids)

		idx := 0
		for crawler := 0; crawler < threads; crawler++ {
			wg.Add(1)
			if crawler == threads-1 {
				numPer = len(ids) - idx
			}
			go func(min, max int) {
				fmt.Printf("Starting worker for id block %d to %d\n", min, max-1)
				client := http.Client{}
				for i := min; i < max; i++ {
					url := fmt.Sprintf("https://wowhead.com/wotlk/tooltip/item/%d?json", ids[i])
					resp, err := client.Get(url)
					if err != nil {
						fmt.Printf("Error fetching %d: %s\n", ids[i], err)
						continue
					}
					body, _ := io.ReadAll(resp.Body)
					bstr := string(body)
					// fmt.Printf("Found tooltip for %d\n%s\n", i, bstr)
					if len(bstr) < 2 {
						fmt.Printf("Missing tooltip data for %d", ids[i])
					}
					results <- result{id: ids[i], value: bstr}
				}
				wg.Done()
			}(idx, idx+numPer)
			idx += numPer
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalComplete := 0
	var lastUpdate time.Time
	for res := range results {
		totalComplete++

		if time.Since(lastUpdate).Seconds() > 2 {
			lastUpdate = time.Now()
			fmt.Printf("Tooltips %d/%d complete\n", totalComplete, total)
		}

		if strings.Contains(res.value, "\"error\":") {
			// fmt.Printf("Error in tooltip for %d: %s\n", i, bstr)
			continue
		}

		url := fmt.Sprintf("https://wowhead.com/wotlk/tooltip/item/%d?json", res.id)
		f.WriteString(fmt.Sprintf("%d, %s, %s\n", res.id, url, res.value))
	}
	fmt.Printf("Tooltips %d/%d complete\n", totalComplete, total)
}

func getItemDeclarations(name string) []int {
	itemBytes, err := os.ReadFile(name)

	if err != nil {
		log.Fatalf("failed to read item declarations file: %s", err)
	}

	lines := strings.Split(string(itemBytes), "\n")
	itemDeclarations := make([]int, 0, len(lines))
	for _, line := range lines {
		itemID, err := strconv.Atoi(line)
		if err != nil || itemID == 0 {
			continue
		}
		itemDeclarations = append(itemDeclarations, itemID)
	}

	return itemDeclarations
}
