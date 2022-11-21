package main

import (
	"flag"

	"github.com/wowsims/wotlk/tools/database"
)

func main() {
	var minId = flag.Int("minid", 1, "ID of item to start scan at")
	var maxId = flag.Int("maxid", 57000, "maximum ID to scan for")
	var output = flag.String("output", "all_item_tooltips.csv", "name of file to output results to")

	flag.Parse()

	wtm := database.NewWotlkItemTooltipManager(*output)
	wtm.Fetch(int32(*minId), int32(*maxId))
}
