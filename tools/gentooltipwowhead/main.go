package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	var minId = flag.Int("minid", 1, "ID of item to start scan at")
	var maxId = flag.Int("maxid", 57000, "maximum ID to scan for")
	var idList = flag.String("ids", "", "Comma-separated list of IDs to fetch.")
	var numThreads = flag.Int("threads", 8, "number of parallel workers to fetch tooltips with")
	var output = flag.String("output", "all_item_tooltips.csv", "name of file to output results to")

	flag.Parse()

	database := map[int]string{}
	if existingFile, err := os.Open(*output); err == nil {
		scanner := bufio.NewScanner(existingFile)
		for scanner.Scan() {
			line := scanner.Text()

			itemIDStr := line[:strings.Index(line, ",")]
			itemID, err := strconv.Atoi(itemIDStr)
			if err != nil {
				log.Fatal("Invalid item ID: " + itemIDStr)
			}

			tooltip := line[strings.Index(line, "{"):]
			database[itemID] = tooltip
		}
		existingFile.Close()

		fmt.Printf("Found %d existing items in tooltip database.\n", len(database))
	}

	var idsToFetch []int
	if *idList == "" {
		for i := *minId; i <= *maxId; i++ {
			idsToFetch = append(idsToFetch, i)
		}

		// Filter out IDs that are already in the database.
		numIDsRemaining := 0
		for _, id := range idsToFetch {
			if _, ok := database[id]; !ok {
				idsToFetch[numIDsRemaining] = id
				numIDsRemaining++
			}
		}
		idsToFetch = idsToFetch[:numIDsRemaining]
	} else {
		idStrings := strings.Split(*idList, ",")
		for _, str := range idStrings {
			id, err := strconv.Atoi(str)
			if err != nil {
				log.Fatal("Invalid item ID: " + str)
			}
			idsToFetch = append(idsToFetch, id)
		}
	}

	fmt.Printf("Will fetch %d core...\n", len(idsToFetch))

	threads := *numThreads
	if threads > len(idsToFetch) {
		threads = len(idsToFetch)
	}

	type result struct {
		id    int
		value string
	}
	results := make(chan result, 10)
	wg := &sync.WaitGroup{}

	for thread := 0; thread < threads; thread++ {
		startIdx := len(idsToFetch) * thread / threads
		endIdx := len(idsToFetch) * (thread + 1) / threads
		wg.Add(1)
		go func(min, max int) {
			fmt.Printf("Starting worker for id block %d to %d\n", min, max-1)
			client := http.Client{}
			for i := min; i < max; i++ {
				id := idsToFetch[i]
				url := fmt.Sprintf("https://nether.wowhead.com/wotlk/tooltip/item/%d", id)
				resp, err := client.Get(url)
				if err != nil {
					fmt.Printf("Error fetching %d: %s\n", id, err)
					continue
				}
				body, _ := io.ReadAll(resp.Body)
				bstr := string(body)
				// fmt.Printf("Found tooltip for %d\n%s\n", i, bstr)
				if len(bstr) < 2 {
					fmt.Printf("Missing tooltip data for %d", id)
				}
				results <- result{id: id, value: bstr}
			}
			wg.Done()
		}(startIdx, endIdx)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	tempFile, err := os.Create("temp_tooltips.csv")
	if err != nil {
		panic("failed to create temp file to write tooltips to: " + err.Error())
	}

	totalComplete := 0
	var lastUpdate time.Time
	for res := range results {
		totalComplete++

		if time.Since(lastUpdate).Seconds() > 2 {
			lastUpdate = time.Now()
			fmt.Printf("Tooltips %d/%d complete\n", totalComplete, len(idsToFetch))
		}

		if strings.Contains(res.value, "\"error\":") {
			// fmt.Printf("Error in tooltip for %d: %s\n", i, bstr)
			continue
		}

		database[res.id] = res.value // replaces existing tooltips
		url := fmt.Sprintf("https://nether.wowhead.com/wotlk/tooltip/item/%d", res.id)
		tempFile.WriteString(fmt.Sprintf("%d, %s, %s\n", res.id, url, res.value))
	}
	fmt.Printf("Tooltips %d/%d complete\nNow writing tooltip file...", totalComplete, len(idsToFetch))

	finalOutput, err := os.Create(*output)
	if err != nil {
		log.Fatalf("failed to open output file to write: %s", err)
	}

	keys := []int{}
	for k, _ := range database {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		url := fmt.Sprintf("https://nether.wowhead.com/wotlk/tooltip/item/%d", k)
		finalOutput.WriteString(fmt.Sprintf("%d, %s, %s\n", k, url, database[k]))
	}

	fmt.Printf("Complete.\n")
}
