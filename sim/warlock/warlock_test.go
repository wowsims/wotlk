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
