package warlock

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterWarlock()
}

func TestAffliction(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassWarlock,
		Race:  proto.Race_RaceOrc,

		GearSet:     core.GetGearSet("../../ui/warlock/gear_sets", "p4_affliction"),
		Talents:     AfflictionTalents,
		Glyphs:      AfflictionGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Affliction Warlock", SpecOptions: DefaultAfflictionWarlock},
		Rotation:    core.GetAplRotation("../../ui/warlock/apls", "affliction"),

		ItemFilter: ItemFilter,
	}))
}

func TestDemonology(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassWarlock,
		Race:  proto.Race_RaceOrc,

		GearSet:     core.GetGearSet("../../ui/warlock/gear_sets", "p4_demo"),
		Talents:     DemonologyTalents,
		Glyphs:      DemonologyGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Demonology Warlock", SpecOptions: DefaultDemonologyWarlock},
		Rotation:    core.GetAplRotation("../../ui/warlock/apls", "demo"),

		ItemFilter: ItemFilter,
	}))
}

func TestDestruction(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassWarlock,
		Race:  proto.Race_RaceOrc,

		GearSet:     core.GetGearSet("../../ui/warlock/gear_sets", "p4_destro"),
		Talents:     DestructionTalents,
		Glyphs:      DestructionGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Destruction Warlock", SpecOptions: DefaultDestroWarlock},
		Rotation:    core.GetAplRotation("../../ui/warlock/apls", "destro"),
		ItemFilter:  ItemFilter,
	}))
}

var ItemFilter = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeDagger,
	},
	HandTypes: []proto.HandType{
		proto.HandType_HandTypeOffHand,
	},
	ArmorType: proto.ArmorType_ArmorTypeCloth,
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeWand,
	},
}

var AfflictionTalents = "2350002030023510253500331151--550000051"
var DemonologyTalents = "-203203301035012530135201351-550000052"
var DestructionTalents = "-03310030003-05203205210331051335230351"
var AfflictionGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarlockMajorGlyph_GlyphOfQuickDecay),
	Major2: int32(proto.WarlockMajorGlyph_GlyphOfLifeTap),
	Major3: int32(proto.WarlockMajorGlyph_GlyphOfHaunt),
}
var DemonologyGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarlockMajorGlyph_GlyphOfQuickDecay),
	Major2: int32(proto.WarlockMajorGlyph_GlyphOfLifeTap),
	Major3: int32(proto.WarlockMajorGlyph_GlyphOfFelguard),
}
var DestructionGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarlockMajorGlyph_GlyphOfConflagrate),
	Major2: int32(proto.WarlockMajorGlyph_GlyphOfLifeTap),
	Major3: int32(proto.WarlockMajorGlyph_GlyphOfIncinerate),
}

var defaultDestroOptions = &proto.Warlock_Options{
	Armor:        proto.Warlock_Options_FelArmor,
	Summon:       proto.Warlock_Options_Imp,
	WeaponImbue:  proto.Warlock_Options_GrandFirestone,
	DetonateSeed: true,
}

var DefaultDestroWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Options: defaultDestroOptions,
	},
}

// ---------------------------------------
var DefaultAfflictionWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Options: defaultAfflictionOptions,
	},
}

var defaultAfflictionOptions = &proto.Warlock_Options{
	Armor:        proto.Warlock_Options_FelArmor,
	Summon:       proto.Warlock_Options_Felhunter,
	WeaponImbue:  proto.Warlock_Options_GrandSpellstone,
	DetonateSeed: true,
}

// ---------------------------------------
var DefaultDemonologyWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Options: defaultDemonologyOptions,
	},
}

var defaultDemonologyOptions = &proto.Warlock_Options{
	Armor:        proto.Warlock_Options_FelArmor,
	Summon:       proto.Warlock_Options_Felguard,
	WeaponImbue:  proto.Warlock_Options_GrandSpellstone,
	DetonateSeed: true,
}

// ---------------------------------------------------------

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	DefaultPotion: proto.Potions_PotionOfWildMagic,
	PrepopPotion:  proto.Potions_PotionOfWildMagic,
	Food:          proto.Food_FoodFishFeast,
}
