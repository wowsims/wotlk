package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/wowsims/wotlk/sim/core/proto"
)

func main() {
	outDir := flag.String("outDir", "", "Path to output directory for writing generated .go files.")
	db := flag.String("db", "wowhead", "which database to use")
	flag.Parse()

	if *outDir == "" {
		panic("outDir flag is required!")
	}

	tooltipsDB := getWowheadTooltipsDB()

	// Generate all item/gem ids from the tooltips db

	// gems := &strings.Builder{}
	// items := &strings.Builder{}
	// for k := range tooltipsDB {
	// 	resp := getWowheadItemResponse(k, tooltipsDB)
	// 	if resp.Name == "" || strings.Contains(resp.Name, "QA Test") || strings.Contains(resp.Name, "zzOLD") || strings.Contains(resp.Name, "Monster -") {
	// 		continue
	// 	}
	// 	if resp.IsPattern() {
	// 		continue
	// 	}
	// 	// No socket color means that this isn't a gem
	// 	if resp.GetSocketColor() == proto.GemColor_GemColorUnknown {
	// 		itemLevel := resp.GetItemLevel()
	// 		qual := resp.GetQuality()
	// 		if qual < int(proto.ItemQuality_ItemQualityUncommon) {
	// 			continue
	// 		} else if qual > int(proto.ItemQuality_ItemQualityLegendary) {
	// 			continue
	// 		} else if qual < int(proto.ItemQuality_ItemQualityEpic) {
	// 			if itemLevel < 105 {
	// 				continue
	// 			}
	// 			if itemLevel < 110 && resp.GetItemSetName() == "" {
	// 				continue
	// 			}
	// 		} else if qual < int(proto.ItemQuality_ItemQualityEpic) {
	// 			if itemLevel < 110 {
	// 				continue
	// 			}
	// 			if itemLevel < 140 && resp.GetItemSetName() == "" {
	// 				continue
	// 			}
	// 		} else {
	// 			// Epic and legendary items might come from classic, so use a lower ilvl threshold.
	// 			if itemLevel < 75 {
	// 				continue
	// 			}
	// 		}

	// 		items.WriteString(fmt.Sprintf("%d\n", k))
	// 	} else {
	// 		qual := resp.GetQuality()
	// 		if qual <= int(proto.ItemQuality_ItemQualityUncommon) && k < 30000 {
	// 			continue
	// 		}
	// 		gems.WriteString(fmt.Sprintf("%d\n", k))
	// 	}
	// }

	// os.WriteFile("all_item_ids", []byte(items.String()), 0666)
	// os.WriteFile("all_gem_ids", []byte(gems.String()), 0666)

	// panic("done")

	var gemsData []GemData
	var itemsData []ItemData
	if *db == "wowhead" {
		gemDeclarations := getGemDeclarations()
		for _, gemDeclaration := range gemDeclarations {
			gemData := GemData{
				Declaration: gemDeclaration,
				Response:    getWowheadItemResponse(gemDeclaration.ID, tooltipsDB),
			}
			if gemData.Response.GetName() == "" {
				continue
			}
			//log.Printf("\n\n%+v\n", gemData.Response)
			gemsData = append(gemsData, gemData)
		}

		itemDeclarations := getItemDeclarations()
		// qualityModifiers := getItemQualityModifiers()
		for _, itemDeclaration := range itemDeclarations {
			itemData := ItemData{
				Declaration: itemDeclaration,
				Response:    getWowheadItemResponse(itemDeclaration.ID, tooltipsDB),
				// QualityModifier: qualityModifiers[itemDeclaration.ID],
			}
			if itemData.Response.GetName() == "" {
				continue
			}
			//fmt.Printf("\n\n%+v\n", itemData.Response)
			itemsData = append(itemsData, itemData)
		}
	} else if *db == "wotlkdb" {
		itemsData = make([]ItemData, 0, len(tooltipsDB))
		gemsData = make([]GemData, 0, len(tooltipsDB))
		for k := range tooltipsDB {
			resp := getWotlkItemResponse(k, tooltipsDB)
			if resp.Name == "" || strings.Contains(resp.Name, "zzOLD") {
				continue
			}
			if resp.IsPattern() {
				continue
			}
			// No socket color means that this isn't a gem
			if resp.GetSocketColor() == proto.GemColor_GemColorUnknown {
				itemsData = append(itemsData, ItemData{Response: resp, Declaration: ItemDeclaration{ID: k}})
			} else {
				gemsData = append(gemsData, GemData{Response: resp, Declaration: GemDeclaration{ID: k}})
			}
		}
	} else {
		panic("invalid item database source")
	}

	slices.SortStableFunc(gemsData, func(g1, g2 GemData) bool {
		if g1.Response.GetName() == g2.Response.GetName() {
			return g1.Declaration.ID < g2.Declaration.ID
		}
		return g1.Response.GetName() < g2.Response.GetName()
	})
	writeGemFile(*outDir, gemsData)

	slices.SortStableFunc(itemsData, func(i1, i2 ItemData) bool {
		if i1.Response.GetName() == i2.Response.GetName() {
			return i1.Declaration.ID < i2.Declaration.ID
		}
		return i1.Response.GetName() < i2.Response.GetName()
	})
	writeItemFile(*outDir, itemsData)

	writeDatabaseFile(gemsData, itemsData)
}

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
func getWowheadTooltipsDB() map[int]string {
	file, err := os.Open("./assets/item_data/all_item_tooltips.csv")
	if err != nil {
		log.Fatalf("Failed to open all_item_tooltips.csv: %s", err)
	}
	defer file.Close()

	db := make(map[int]string)
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
		db[itemID] = tooltip
	}

	fmt.Printf("\n--\nTOOLTIPS LOADED: %d\n--\n", len(db))
	return db
}

// func getItemQualityModifiers() map[int]float64 {
// 	file, err := os.Open("./assets/item_data/quality_modifiers.csv")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	qualityMods := make(map[int]float64)
// 	scanner := bufio.NewScanner(file)
// 	i := 0
// 	for scanner.Scan() {
// 		i++
// 		if i == 1 {
// 			// Ignore first line
// 			continue
// 		}

// 		line := scanner.Text()

// 		itemIDStr := line[:strings.Index(line, ",")]
// 		itemID, err := strconv.Atoi(itemIDStr)
// 		if err != nil {
// 			log.Fatal("Invalid item ID: " + itemIDStr)
// 		}

// 		qualityModStr := line[strings.LastIndex(line, ",")+1:]
// 		qualityMod, err := strconv.ParseFloat(qualityModStr, 64)
// 		if err != nil {
// 			log.Fatal("Invalid quality mod: ", qualityModStr)
// 		}

// 		qualityMods[itemID] = qualityMod
// 	}

// 	return qualityMods
// }

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

func writeDatabaseFile(gemsData []GemData, itemsData []ItemData) {
	type tempItemIcon struct {
		ID   int
		Name string
		Icon string
	}

	var tempItems []tempItemIcon
	for _, gemData := range gemsData {
		tempItems = append(tempItems, tempItemIcon{ID: gemData.Declaration.ID, Name: gemData.Response.GetName(), Icon: gemData.Response.GetIcon()})
	}
	for _, itemData := range itemsData {
		tempItems = append(tempItems, tempItemIcon{ID: itemData.Declaration.ID, Name: itemData.Response.GetName(), Icon: itemData.Response.GetIcon()})
	}

	// Write it out line-by-line so we can have 1 line / item, making it more human-readable.
	itemDB := &strings.Builder{}
	itemDB.WriteString("[\n")
	for i, item := range tempItems {
		itemJson, err := json.Marshal(item)
		if err != nil {
			log.Fatalf("failed to marshal: %s", err)
		}
		itemDB.WriteString(string(itemJson))
		if i != len(tempItems)-1 {
			itemDB.WriteString(",\n")
		}
	}
	itemDB.WriteString("]")
	os.WriteFile("./assets/item_data/all_items_db.json", []byte(itemDB.String()), 0666)
}
