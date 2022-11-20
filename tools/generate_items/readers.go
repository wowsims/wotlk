// This file is for functions that read and parse already-materialzed asset files.
// It is NOT for scrapers.
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func getItemOverrides() []ItemOverride {
	itemsData := readCsvFile("./assets/item_data/all_item_ids.csv")

	// Ignore first line
	itemsData = itemsData[1:]

	// Create an empty declaration (just the ID) for all the core.
	itemOverrides := make([]ItemOverride, 0, len(itemsData))
	for _, itemsDataRow := range itemsData {
		itemID, err := strconv.Atoi(itemsDataRow[0])
		if err != nil {
			log.Fatal("Invalid item ID: " + itemsDataRow[0])
		}

		itemOverrides = append(itemOverrides, ItemOverride{
			ID: itemID,
		})
	}

	// Apply declarations overrides.
	for _, overrideItemOverride := range ItemOverrideOverrides {
		found := false
		for i, itemDecl := range itemOverrides {
			if itemDecl.ID == overrideItemOverride.ID {
				found = true
				itemOverrides[i] = overrideItemOverride
				break
			}
		}
		if !found {
			itemOverrides = append(itemOverrides, overrideItemOverride)
		}
	}

	return itemOverrides
}

// Returns the prefetched list of all wowhead tooltips.
// Maps item IDs to tooltip strings.
func getWowheadTooltipsDB(filepath string) map[int]WowheadItemResponse {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Failed to open %s: %s", filepath, err)
	}
	defer file.Close()

	db := make(map[int]WowheadItemResponse)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		itemIDStr := line[:strings.Index(line, ",")]
		itemID, err := strconv.Atoi(itemIDStr)
		if err != nil {
			log.Fatal("Invalid item ID: " + itemIDStr)
		}

		tooltip := line[strings.Index(line, "{"):]
		db[itemID] = WowheadItemResponseFromBytes([]byte(tooltip))
	}

	fmt.Printf("\n--\nTOOLTIPS LOADED: %d\n--\n", len(db))
	return db
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}
