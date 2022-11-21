package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/tools"
)

func main() {
	var minId = flag.Int("minid", 1, "ID of item to start scan at")
	var maxId = flag.Int("maxid", 57000, "maximum ID to scan for")
	var idList = flag.String("ids", "", "Comma-separated list of IDs to fetch.")
	var output = flag.String("output", "all_item_tooltips.csv", "name of file to output results to")

	flag.Parse()

	database := map[int]string{}
	if lines := tools.ReadFileLinesOrNil(*output); lines != nil {
		for _, line := range lines {
			itemIDStr := line[:strings.Index(line, ",")]
			itemID, err := strconv.Atoi(itemIDStr)
			if err != nil {
				log.Fatal("Invalid item ID: " + itemIDStr)
			}

			tooltip := line[strings.Index(line, "{"):]
			database[itemID] = tooltip
		}

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

	newTooltips := tools.ReadWebMultiMap(idsToFetch, func(id int) string {
		return fmt.Sprintf("https://nether.wowhead.com/wotlk/tooltip/item/%d", id)
	})

	newTooltips = core.FilterMap(newTooltips, func(id int, body string) bool {
		if len(body) < 2 {
			fmt.Printf("Missing tooltip data for %d", id)
			return false
		}
		if strings.Contains(body, "\"error\":") {
			// fmt.Printf("Error in tooltip for %d: %s\n", id, body)
			return false
		}
		return true
	})

	for k, v := range newTooltips {
		database[k] = v
	}

	keys := []int{}
	for k, _ := range database {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	linesToWrite := core.MapSlice(keys, func(key int) string {
		url := fmt.Sprintf("https://nether.wowhead.com/wotlk/tooltip/item/%d", key)
		return fmt.Sprintf("%d, %s, %s\n", key, url, database[key])
	})

	tools.WriteFileLines(*output, linesToWrite)
}
