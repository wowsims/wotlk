package database

import (
	"github.com/wowsims/classic/sim/core/proto"
)

// Overrides for runes as needed
// Attempt to run "scrape_runes.py" prior as wowhead may update with datamining over time
// E.g. "python3 tools/scrape_runes.py assets/db_inputs/wowhead_rune_tooltips.csv"
// Then re-gen db "go run ./tools/database/gen_db -outDir=assets -gen=db"
var RuneOverrides = []*proto.UIRune{
	// Chest
	// Priest
	{Id: 425210, Name: "Engrave Chest - Twisted Faith", Type: proto.ItemType_ItemTypeChest, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 425211, Name: "Engrave Chest - Void Plague", Type: proto.ItemType_ItemTypeChest, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 415737, Name: "Engrave Chest - Serendipity", Type: proto.ItemType_ItemTypeChest, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 415740, Name: "Engrave Chest - Strength of Soul", Type: proto.ItemType_ItemTypeChest, Class: proto.Class_ClassPriest, RequiresLevel: 1},

	// Gloves
	// Priest
	{Id: 415738, Name: "Engrave Gloves - Mind Sear", Type: proto.ItemType_ItemTypeHands, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 402844, Name: "Engrave Gloves - Penance", Type: proto.ItemType_ItemTypeHands, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 402833, Name: "Engrave Gloves - Shadow Word: Death", Type: proto.ItemType_ItemTypeHands, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 402842, Name: "Engrave Gloves - Circle of Healing", Type: proto.ItemType_ItemTypeHands, Class: proto.Class_ClassPriest, RequiresLevel: 1},

	// Legs
	// Priest
	{Id: 402836, Name: "Engrave Legs - Homunculi", Type: proto.ItemType_ItemTypeLegs, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 402838, Name: "Engrave Chest - Shared Pain", Type: proto.ItemType_ItemTypeLegs, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 425212, Name: "Engrave Chest - Power Word: Barrier", Type: proto.ItemType_ItemTypeLegs, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 402832, Name: "Engrave Chest - Prayer of Mending", Type: proto.ItemType_ItemTypeLegs, Class: proto.Class_ClassPriest, RequiresLevel: 1},
}
