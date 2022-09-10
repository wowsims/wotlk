package feral

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterFeralDruid()
}

func TestFeral(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDruid,

		Race: proto.Race_RaceTauren,

		GearSet: core.GearSetCombo{Label: "P1", GearSet: P1Gear},

		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsBearweaveLacerate},

		RaidBuffs:   FullRaidBuffs,
		PartyBuffs:  FullPartyBuffs,
		PlayerBuffs: FullIndividualBuffs,
		Consumes:    FullConsumes,
		Debuffs:     FullDebuffs,
		Glyphs:      StandardGlyphs,

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeLeather,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeIdol,
			},
		},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:      proto.Race_RaceTauren,
				Class:     proto.Class_ClassDruid,
				Equipment: P1Gear,
				Consumes:  FullConsumes,
				Spec:      PlayerOptionsBearweaveLacerate,
				Buffs:     FullIndividualBuffs,
				Glyphs:    StandardGlyphs,

				InFrontOfTarget: true,
			},
			FullPartyBuffs,
			FullRaidBuffs,
			FullDebuffs),
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
