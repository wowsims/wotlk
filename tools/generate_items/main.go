package main

import (
	"flag"
)

func main() {
	outDir := flag.String("outDir", "", "Path to output directory for writing generated .go files.")
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

	writeDatabaseFile(db)
}
