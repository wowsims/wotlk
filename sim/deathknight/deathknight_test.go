package deathknight

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get item effects included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	// Disabled until better rotations
	return
	RegisterDeathKnight()
}

func TestUnholy(t *testing.T) {
	// Disabled until better rotations
	return
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDeathKnight,

		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet: core.GearSetCombo{Label: "Unholy P1 ", GearSet: UnholyDwP1Gear},

		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsUnholy},

		RaidBuffs:   FullRaidBuffs,
		PartyBuffs:  FullPartyBuffs,
		PlayerBuffs: FullIndividualBuffs,
		Consumes:    FullConsumes,
		Debuffs:     FullDebuffs,

		ItemFilter: core.ItemFilter{
			ArmorType: proto.ArmorType_ArmorTypePlate,

			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeMace,
			},
		},
	}))
}

func TestFrost(t *testing.T) {
	// Disabled until better rotations
	return
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDeathKnight,

		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet: core.GearSetCombo{Label: "Frost P1", GearSet: UnholyDwP1Gear},

		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFrost},

		RaidBuffs:   FullRaidBuffs,
		PartyBuffs:  FullPartyBuffs,
		PlayerBuffs: FullIndividualBuffs,
		Consumes:    FullConsumes,
		Debuffs:     FullDebuffs,

		ItemFilter: core.ItemFilter{
			ArmorType: proto.ArmorType_ArmorTypePlate,

			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeMace,
			},
		},
	}))
}
