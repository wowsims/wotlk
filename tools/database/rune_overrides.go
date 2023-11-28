package database

import (
	"github.com/wowsims/classic/sod/sim/core/proto"
)

// Overrides for runes as needed
// Regen db with "go run ./tools/database/gen_db -outDir=assets -gen=db"
// And ensure db files are copied from assets/db into dist/classic/sod/database
var RuneOverrides = []*proto.UIRune{
	// {Id: 415460, Name: "Engrave Chest - Burnout", Icon: "ability_mage_burnout", Type: proto.ItemType_ItemTypeChest, Class: proto.Class_ClassMage, RequiresLevel: 1},
}
