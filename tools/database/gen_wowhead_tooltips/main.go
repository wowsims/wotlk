package main

import (
	"flag"

	"github.com/wowsims/wotlk/tools/database"
)

// go run ./tools/database/gen_wowhead_tooltips
// go run ./tools/database/gen_wowhead_tooltips -fetchSpells=true

var minId = flag.Int("minid", 1, "ID of item to start scan at")
var maxId = flag.Int("maxid", 57000, "maximum ID to scan for")
var fetchSpells = flag.Bool("fetchSpells", false, "If true, fetch spell tooltips. Otherwise fetch item tooltips.")
var output = flag.String("output", "all_item_tooltips.csv", "name of file to output results to")

func main() {
	flag.Parse()

	if *fetchSpells {
		database.NewWowheadSpellTooltipManager(*output).Fetch(int32(*minId), int32(*maxId))
	} else {
		database.NewWowheadItemTooltipManager(*output).Fetch(int32(*minId), int32(*maxId))
	}
}
