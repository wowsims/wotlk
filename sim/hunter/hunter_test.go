package hunter

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get item effects included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterHunter()
}

func TestHunter(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassHunter,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceDwarf},

		GearSet: core.GearSetCombo{Label: "P1", GearSet: P1Gear},

		SpecOptions: core.SpecOptionsCombo{Label: "Marksman", SpecOptions: PlayerOptionsMM},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "BM", SpecOptions: PlayerOptionsBM},
			{Label: "SV", SpecOptions: PlayerOptionsSV},
			{Label: "AOE", SpecOptions: PlayerOptionsAOE},
		},

		Glyphs:   DefaultGlyphs,
		Consumes: FullConsumes,

		ItemFilter: core.ItemFilter{
			ArmorType: proto.ArmorType_ArmorTypeMail,
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeFist,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypePolearm,
				proto.WeaponType_WeaponTypeStaff,
				proto.WeaponType_WeaponTypeSword,
			},
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeBow,
				proto.RangedWeaponType_RangedWeaponTypeCrossbow,
				proto.RangedWeaponType_RangedWeaponTypeGun,
			},
		},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:      proto.Race_RaceOrc,
				Class:     proto.Class_ClassHunter,
				Equipment: P1Gear,
				Consumes:  FullConsumes,
				Spec:      PlayerOptionsMM,
				Buffs:     core.FullIndividualBuffs,
			},
			core.FullPartyBuffs,
			core.FullRaidBuffs,
			core.FullDebuffs),
		Encounter: &proto.Encounter{
			Duration: 300,
			Targets: []*proto.Target{
				core.NewDefaultTarget(),
			},
		},
		SimOptions: core.AverageDefaultSimTestOptions,
	}

	core.RaidBenchmark(b, rsr)
}
