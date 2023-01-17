package warlock

import (
	"testing"

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

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     AfflictionTalents,
		Glyphs:      AfflictionGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Affliction Warlock", SpecOptions: DefaultAfflictionWarlock},

		ItemFilter: ItemFilter,
	}))
}

func TestDemonology(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassWarlock,
		Race:  proto.Race_RaceOrc,

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     DemonologyTalents,
		Glyphs:      DemonologyGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Demonology Warlock", SpecOptions: DefaultDemonologyWarlock},

		ItemFilter: ItemFilter,
	}))
}

func TestDestruction(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassWarlock,
		Race:  proto.Race_RaceOrc,

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     DestructionTalents,
		Glyphs:      DestructionGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Destruction Warlock", SpecOptions: DefaultDestroWarlock},

		ItemFilter: ItemFilter,
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

var defaultDestroRotation = &proto.Warlock_Rotation{
	Type:         proto.Warlock_Rotation_Destruction,
	PrimarySpell: proto.Warlock_Rotation_Incinerate,
	SecondaryDot: proto.Warlock_Rotation_Immolate,
	SpecSpell:    proto.Warlock_Rotation_ChaosBolt,
	Curse:        proto.Warlock_Rotation_Doom,
	Corruption:   false,
}

var defaultDestroOptions = &proto.Warlock_Options{
	Armor:       proto.Warlock_Options_FelArmor,
	Summon:      proto.Warlock_Options_Imp,
	WeaponImbue: proto.Warlock_Options_GrandFirestone,
}

var DefaultDestroWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Options:  defaultDestroOptions,
		Rotation: defaultDestroRotation,
	},
}

// ---------------------------------------
var DefaultAfflictionWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Options:  defaultAfflictionOptions,
		Rotation: defaultAfflictionRotation,
	},
}

var defaultAfflictionOptions = &proto.Warlock_Options{
	Armor:       proto.Warlock_Options_FelArmor,
	Summon:      proto.Warlock_Options_Felhunter,
	WeaponImbue: proto.Warlock_Options_GrandSpellstone,
}

var defaultAfflictionRotation = &proto.Warlock_Rotation{
	Type:         proto.Warlock_Rotation_Affliction,
	PrimarySpell: proto.Warlock_Rotation_ShadowBolt,
	SecondaryDot: proto.Warlock_Rotation_UnstableAffliction,
	SpecSpell:    proto.Warlock_Rotation_Haunt,
	Curse:        proto.Warlock_Rotation_Agony,
	Corruption:   true,
	DetonateSeed: true,
}

// ---------------------------------------
var DefaultDemonologyWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Options:  defaultDemonologyOptions,
		Rotation: defaultDemonologyRotation,
	},
}

var defaultDemonologyOptions = &proto.Warlock_Options{
	Armor:       proto.Warlock_Options_FelArmor,
	Summon:      proto.Warlock_Options_Felguard,
	WeaponImbue: proto.Warlock_Options_GrandSpellstone,
}

var defaultDemonologyRotation = &proto.Warlock_Rotation{
	Type:         proto.Warlock_Rotation_Demonology,
	PrimarySpell: proto.Warlock_Rotation_ShadowBolt,
	SecondaryDot: proto.Warlock_Rotation_Immolate,
	Curse:        proto.Warlock_Rotation_Doom,
	Corruption:   true,
	DetonateSeed: true,
}

// ---------------------------------------------------------

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	DefaultPotion: proto.Potions_PotionOfWildMagic,
	PrepopPotion:  proto.Potions_PotionOfWildMagic,
	Food:          proto.Food_FoodFishFeast,
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40421,
		"enchant": 3820,
		"gems": [
			41285,
			40051
		]
	},
	{
		"id": 44661,
		"gems": [
			40026
		]
	},
	{
		"id": 40424,
		"enchant": 3810,
		"gems": [
			39998
		]
	},
	{
		"id": 44005,
		"enchant": 3722,
		"gems": [
			40026
		]
	},
	{
		"id": 40423,
		"enchant": 3832,
		"gems": [
			39998,
			40051
		]
	},
	{
		"id": 44008,
		"enchant": 2332,
		"gems": [
			39998,
			0
		]
	},
	{
		"id": 40420,
		"enchant": 3604,
		"gems": [
			39998,
			0
		]
	},
	{
		"id": 40561,
		"gems": [
			39998
		]
	},
	{
		"id": 40560,
		"enchant": 3719
	},
	{
		"id": 40558,
		"enchant": 3606
	},
	{
		"id": 40399
	},
	{
		"id": 40719
	},
	{
		"id": 40432
	},
	{
		"id": 40255
	},
	{
		"id": 40396,
		"enchant": 3834
	},
	{
		"id": 39766
	},
	{
		"id": 39712
	}
]}`)
