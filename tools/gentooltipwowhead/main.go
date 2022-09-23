package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	var maxParam = flag.Int("maxid", 52205, "maximum ID to scan for")
	var minParam = flag.Int("minid", 1, "ID of item to start scan at")
	var numThreads = flag.Int("threads", 8, "number of parallel workers to fetch tooltips with")
	var output = flag.String("output", "all_item_tooltips.csv", "name of file to output results to")

	flag.Parse()

	maxID := *maxParam
	startMinID := *minParam

	threads := *numThreads
	numPer := (maxID - startMinID) / threads

	fmt.Printf("Starting tooltip fetching.\n\tStart ID: %d, End ID: %d\n\tWorkers: %d\n\tOutput File: %s\n\n", startMinID, maxID, threads, *output)
	f, err := os.Create(*output)
	if err != nil {
		log.Fatalf("failed to open file to write: %s", err)
	}
	type result struct {
		id    int
		value string
	}
	results := make(chan result, 10)

	currentMinID := startMinID
	wg := &sync.WaitGroup{}
	for crawler := 0; crawler < threads; crawler++ {
		wg.Add(1)
		if crawler == threads-1 {
			numPer = maxID - currentMinID // last thread gets the odd extras
		}
		go func(min, max int) {
			client := http.Client{}
			for i := min; i < max; i++ {
				url := fmt.Sprintf("https://wowhead.com/wotlk/tooltip/item/%d?json", i)
				resp, err := client.Get(url)
				if err != nil {
					fmt.Printf("Error fetching %d: %s\n", i, err)
					continue
				}
				body, _ := ioutil.ReadAll(resp.Body)
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

	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0
	var lastUpdate time.Time
	for res := range results {
		total++

		if time.Since(lastUpdate).Seconds() > 2 {
			lastUpdate = time.Now()
			fmt.Printf("Tooltips %d/%d complete\n", total, (maxID - startMinID))
		}

		if strings.Contains(res.value, "\"error\":") {
			// fmt.Printf("Error in tooltip for %d: %s\n", i, bstr)
			continue
		}

		url := fmt.Sprintf("https://wowhead.com/wotlk/tooltip/item/%d?json", res.id)
		f.WriteString(fmt.Sprintf("%d, %s, %s\n", res.id, url, res.value))
	}
	fmt.Printf("Tooltips %d/%d complete\n", total, (maxID - startMinID))
}
