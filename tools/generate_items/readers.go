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

func getGemDeclarations() []GemDeclaration {
	gemsData := readCsvFile("./assets/item_data/all_gem_ids.csv")

	// Ignore first line
	gemsData = gemsData[1:]

	gemDeclarations := make([]GemDeclaration, 0, len(gemsData))
	for _, gemsDataRow := range gemsData {
		gemID, err := strconv.Atoi(gemsDataRow[0])
		if err != nil {
			log.Fatal("Invalid gem ID: " + gemsDataRow[0])
		}
		declaration := GemDeclaration{
			ID: gemID,
		}

		for _, override := range GemDeclarationOverrides {
			if override.ID == gemID {
				declaration = override
				break
			}
		}

		gemDeclarations = append(gemDeclarations, declaration)
	}

	// Add any declarations that were missing from the csv file.
	for _, overrideGemDeclaration := range GemDeclarationOverrides {
		found := false
		for _, gemDecl := range gemDeclarations {
			if gemDecl.ID == overrideGemDeclaration.ID {
				found = true
				break
			}
		}
		if !found {
			gemDeclarations = append(gemDeclarations, overrideGemDeclaration)
		}
	}

	return gemDeclarations
}

func getItemDeclarations() []ItemDeclaration {
	itemsData := readCsvFile("./assets/item_data/all_item_ids.csv")

	// Ignore first line
	itemsData = itemsData[1:]

	// Create an empty declaration (just the ID) for all the items.
	itemDeclarations := make([]ItemDeclaration, 0, len(itemsData))
	for _, itemsDataRow := range itemsData {
		itemID, err := strconv.Atoi(itemsDataRow[0])
		if err != nil {
			log.Fatal("Invalid item ID: " + itemsDataRow[0])
		}

		itemDeclarations = append(itemDeclarations, ItemDeclaration{
			ID: itemID,
		})
	}

	// Apply declarations overrides.
	for _, overrideItemDeclaration := range ItemDeclarationOverrides {
		found := false
		for i, itemDecl := range itemDeclarations {
			if itemDecl.ID == overrideItemDeclaration.ID {
				found = true
				itemDeclarations[i] = overrideItemDeclaration
				break
			}
		}
		if !found {
			itemDeclarations = append(itemDeclarations, overrideItemDeclaration)
		}
	}

	return itemDeclarations
}

// Returns the prefetched list of all wowhead tooltips.
// Maps item IDs to tooltip strings.
func getWowheadTooltipsDB() map[int]WowheadItemResponse {
	file, err := os.Open("./assets/item_data/all_item_tooltips.csv")
	if err != nil {
		log.Fatalf("Failed to open all_item_tooltips.csv: %s", err)
	}
	defer file.Close()

	db := make(map[int]WowheadItemResponse)
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		i++
		if i == 1 {
			// Ignore first line
			continue
		}

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
