package shadow

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get caster sets included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterShadowPriest()
}

func TestShadow(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassPriest,
		Race:       proto.Race_RaceUndead,
		OtherRaces: []proto.Race{proto.Race_RaceNightElf, proto.Race_RaceDraenei},

		GearSet:  core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:  DefaultTalents,
		Glyphs:   DefaultGlyphs,
		Consumes: FullConsumes,

		SpecOptions: core.SpecOptionsCombo{Label: "Ideal", SpecOptions: PlayerOptionsIdeal},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "Basic", SpecOptions: PlayerOptionsBasic},
			{Label: "Clipping", SpecOptions: PlayerOptionsClipping},
		},

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
