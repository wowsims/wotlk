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

	db := NewWowDatabase(getItemDeclarations(), getGemDeclarations(), getWowheadTooltipsDB())

	writeItemFile(*outDir, db.getSimmableItems())
	writeGemFile(*outDir, db.getSimmableGems())
	writeDatabaseFile(db)
}
