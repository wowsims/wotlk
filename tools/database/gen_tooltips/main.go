package main

import (
	"flag"

	"github.com/wowsims/wotlk/tools"
	"github.com/wowsims/wotlk/tools/database"
)

// To do a full re-scrape, delete the previous output file first.
// go run ./tools/database/gen_tooltips -source=wowhead-items -output=assets/db_inputs/wowhead_item_tooltips.csv
// go run ./tools/database/gen_tooltips -source=wowhead-spells -output=assets/db_inputs/wowhead_spell_tooltips.csv
// go run ./tools/database/gen_tooltips -source=wowhead-gearplannerdb -output=assets/db_inputs/wowhead_gearplannerdb.txt
// go run ./tools/database/gen_tooltips -source=wotlk-items -output=assets/db_inputs/wotlk_items_tooltips.csv

var minId = flag.Int("minid", 1, "Minimum ID to scan for")
var maxId = flag.Int("maxid", 57000, "Maximum ID to scan for")
var source = flag.String("source", "", "Which source to fetch tooltips from. Valid values are 'wowhead-items', 'wowhead-spells', 'wowhead-itemdb', and 'wotlk-items'")
var output = flag.String("output", "", "Output file name")

func main() {
	flag.Parse()

	if *output == "" {
		panic("output is required")
	}

	if *source == "wowhead-items" {
		database.NewWowheadItemTooltipManager(*output).Fetch(int32(*minId), int32(*maxId))
	} else if *source == "wowhead-spells" {
		database.NewWowheadSpellTooltipManager(*output).Fetch(int32(*minId), int32(*maxId))
	} else if *source == "wowhead-gearplannerdb" {
		tools.WriteFile(*output, tools.ReadWebRequired("https://nether.wowhead.com/wotlk/data/gear-planner?dv=100"))
	} else if *source == "wotlk-items" {
		database.NewWotlkItemTooltipManager(*output).Fetch(int32(*minId), int32(*maxId))
	} else {
		panic("Invalid source")
	}
}
