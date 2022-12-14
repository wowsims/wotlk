package healing

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get caster sets included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterHealingPriest()
}

func TestDisc(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:    proto.Class_ClassPriest,
		Race:     proto.Race_RaceUndead,
		IsHealer: true,

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		SpecOptions: core.SpecOptionsCombo{Label: "Disc", SpecOptions: PlayerOptionsDisc},
		Glyphs:      DiscGlyphs,
		Consumes:    FullConsumes,

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeCloth,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeWand,
			},
		},

		EPReferenceStat: proto.Stat_StatSpellPower,
		StatsToWeigh: []proto.Stat{
			proto.Stat_StatIntellect,
			proto.Stat_StatSpellPower,
			proto.Stat_StatSpellHaste,
			proto.Stat_StatSpellCrit,
		},
	}))
}

func TestHoly(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:    proto.Class_ClassPriest,
		Race:     proto.Race_RaceUndead,
		IsHealer: true,

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		SpecOptions: core.SpecOptionsCombo{Label: "Holy", SpecOptions: PlayerOptionsHoly},
		Glyphs:      HolyGlyphs,
		Consumes:    FullConsumes,

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeCloth,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeWand,
			},
		},
	}))
}
