// This file is for functions that read and parse already-materialzed asset files.
// It is NOT for scrapers.
package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/wowsims/wotlk/tools"
)

// Returns the prefetched list of all wowhead tooltips.
// Maps item IDs to tooltip strings.
func getWowheadTooltipsDB(filepath string) map[int32]WowheadItemResponse {
	lines := tools.ReadFileLines(filepath)
	db := make(map[int32]WowheadItemResponse)
	for _, line := range lines {
		itemIDStr := line[:strings.Index(line, ",")]
		itemID, err := strconv.Atoi(itemIDStr)
		if err != nil {
			log.Fatal("Invalid item ID: " + itemIDStr)
		}

		tooltip := line[strings.Index(line, "{"):]
		db[int32(itemID)] = WowheadItemResponseFromBytes([]byte(tooltip))
	}

	fmt.Printf("\n--\nTOOLTIPS LOADED: %d\n--\n", len(db))
	return db
}
