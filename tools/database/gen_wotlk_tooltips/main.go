package main

import (
	"flag"

	"github.com/wowsims/wotlk/tools/database"
)

// go run ./tools/database/gen_wotlk_tooltips

var minId = flag.Int("minid", 1, "ID of item to start scan at")
var maxId = flag.Int("maxid", 57000, "maximum ID to scan for")
var output = flag.String("output", "all_item_tooltips.csv", "name of file to output results to")

func main() {
	flag.Parse()

	database.NewWotlkItemTooltipManager(*output).Fetch(int32(*minId), int32(*maxId))
}
