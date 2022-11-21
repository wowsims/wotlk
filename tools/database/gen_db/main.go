package main

import (
	"flag"
	"fmt"

	"github.com/wowsims/wotlk/tools/database"
)

var outDir = flag.String("outDir", "assets", "Path to output directory for writing generated .go files.")

func main() {
	flag.Parse()
	if *outDir == "" {
		panic("outDir flag is required!")
	}

	dbDir := fmt.Sprintf("%s/database", *outDir)
	inputsDir := fmt.Sprintf("%s/db_inputs", *outDir)

	itemTooltipsManager := database.NewWowheadItemTooltipManager(fmt.Sprintf("%s/all_item_tooltips.csv", inputsDir))
	spellTooltipsManager := database.NewWowheadSpellTooltipManager(fmt.Sprintf("%s/all_spell_tooltips.csv", inputsDir))

	db := database.NewWowDatabase(
		database.ItemOverrides,
		database.GemOverrides,
		database.EnchantOverrides,
		itemTooltipsManager.Read(),
		spellTooltipsManager.Read())

	db.ApplyGlobalFilters()
	db.WriteBinaryAndJson(fmt.Sprintf("%s/db.bin", dbDir), fmt.Sprintf("%s/db.json", dbDir))
}
