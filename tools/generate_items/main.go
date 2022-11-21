package main

import (
	"flag"
)

var outDir = flag.String("outDir", "", "Path to output directory for writing generated .go files.")

func main() {
	flag.Parse()

	if *outDir == "" {
		panic("outDir flag is required!")
	}

	db := NewWowDatabase(
		ItemOverrides,
		GemOverrides,
		EnchantOverrides,
		getWowheadTooltipsDB("./assets/item_data/all_item_tooltips.csv"),
		getWowheadTooltipsDB("./assets/spell_data/all_spell_tooltips.csv"))

	db.applyGlobalFilters()
	db.WriteBinaryAndJson("./assets/database/db.bin", "./assets/database/db.json")
}
