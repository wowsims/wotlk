package database

import (
	"github.com/wowsims/classic/sim/core/proto"
)

// Overrides for runes as needed
// Regen db with "go run ./tools/database/gen_db -outDir=assets -gen=db"
// And ensure db files are copied from assets/db into dist/classic/database
var RuneOverrides = []*proto.UIRune{
	// Chest
	// Priest
	{Id: 425210, Name: "Engrave Chest - Twisted Faith", Type: proto.ItemType_ItemTypeChest, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 425211, Name: "Engrave Chest - Void Plague", Type: proto.ItemType_ItemTypeChest, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 415737, Name: "Engrave Chest - Serendipity", Type: proto.ItemType_ItemTypeChest, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 415740, Name: "Engrave Chest - Strength of Soul", Type: proto.ItemType_ItemTypeChest, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	// Mage
	{Id: 415460, Name: "Engrave Chest - Burnout", Icon: "ability_mage_burnout", Type: proto.ItemType_ItemTypeChest, Class: proto.Class_ClassMage, RequiresLevel: 1},
	{Id: 415729, Name: "Engrave Chest - Enlightenment", Icon: "spell_arcane_mindmastery", Type: proto.ItemType_ItemTypeChest, Class: proto.Class_ClassMage, RequiresLevel: 1},
	{Id: 401741, Name: "Engrave Chest - Fingers of Frost", Icon: "ability_mage_wintersgrasp", Type: proto.ItemType_ItemTypeChest, Class: proto.Class_ClassMage, RequiresLevel: 1},
	{Id: 401743, Name: "Engrave Chest - Regeneration", Icon: "inv_enchant_essencenethersmall", Type: proto.ItemType_ItemTypeChest, Class: proto.Class_ClassMage, RequiresLevel: 1},

	// Gloves
	// Priest
	{Id: 415738, Name: "Engrave Gloves - Mind Sear", Type: proto.ItemType_ItemTypeHands, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 402844, Name: "Engrave Gloves - Penance", Type: proto.ItemType_ItemTypeHands, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 402833, Name: "Engrave Gloves - Shadow Word: Death", Type: proto.ItemType_ItemTypeHands, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 402842, Name: "Engrave Gloves - Circle of Healing", Type: proto.ItemType_ItemTypeHands, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	// Mage
	{Id: 401729, Name: "Engrave Gloves - Arcane Blast", Icon: "spell_arcane_blast", Type: proto.ItemType_ItemTypeHands, Class: proto.Class_ClassMage, RequiresLevel: 1},
	{Id: 401732, Name: "Engrave Gloves - Ice Lance", Icon: "spell_frost_frostblast", Type: proto.ItemType_ItemTypeHands, Class: proto.Class_ClassMage, RequiresLevel: 1},
	{Id: 401731, Name: "Engrave Gloves - Living Bomb", Icon: "ability_mage_livingbomb", Type: proto.ItemType_ItemTypeHands, Class: proto.Class_ClassMage, RequiresLevel: 1},
	{Id: 401734, Name: "Engrave Gloves - Rewind Time", Icon: "spell_holy_borrowedtime", Type: proto.ItemType_ItemTypeHands, Class: proto.Class_ClassMage, RequiresLevel: 1},

	// Legs
	// Priest
	{Id: 402836, Name: "Engrave Pants - Homunculi", Type: proto.ItemType_ItemTypeLegs, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 402838, Name: "Engrave Pants - Shared Pain", Type: proto.ItemType_ItemTypeLegs, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 425212, Name: "Engrave Pants - Power Word: Barrier", Type: proto.ItemType_ItemTypeLegs, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	{Id: 402832, Name: "Engrave Pants - Prayer of Mending", Type: proto.ItemType_ItemTypeLegs, Class: proto.Class_ClassPriest, RequiresLevel: 1},
	// Mage
	{Id: 425168, Name: "Engrave Pants - Arcane Surge", Icon: "spell_arcane_arcanetorrent", Type: proto.ItemType_ItemTypeLegs, Class: proto.Class_ClassMage, RequiresLevel: 1},
	{Id: 425169, Name: "Engrave Pants - Icy Veins", Icon: "spell_frost_coldhearted", Type: proto.ItemType_ItemTypeLegs, Class: proto.Class_ClassMage, RequiresLevel: 1},
	{Id: 401744, Name: "Engrave Pants - Living Flame", Icon: "spell_fire_masterofelements", Type: proto.ItemType_ItemTypeLegs, Class: proto.Class_ClassMage, RequiresLevel: 1},
	{Id: 415467, Name: "Engrave Pants - Mass Regeneration", Icon: "inv_enchant_essencenetherlarge", Type: proto.ItemType_ItemTypeLegs, Class: proto.Class_ClassMage, RequiresLevel: 1},
}
