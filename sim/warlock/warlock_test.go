package warlock

import (
	"testing"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterWarlock()
}

func TestWarlock(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassWarlock,

		Race: proto.Race_RaceOrc,

		GearSet: core.GearSetCombo{Label: "P1", GearSet: P1Gear},

		SpecOptions: core.SpecOptionsCombo{Label: "Affliction Warlock", SpecOptions: DefaultAfflictionWarlock},
		OtherSpecOptions: []core.SpecOptionsCombo{
			core.SpecOptionsCombo{Label: "Demonology Warlock", SpecOptions: DefaultDemonologyWarlock},
			core.SpecOptionsCombo{Label: "Destro Warlock", SpecOptions: DefaultDestroWarlock},
		},

		Glyphs: defaultAfflictionGlyphs,

		RaidBuffs:   FullRaidBuffs,
		PartyBuffs:  FullPartyBuffs,
		PlayerBuffs: FullIndividualBuffs,
		Consumes:    FullConsumes,
		Debuffs:     FullDebuffs,

		ItemFilter: core.ItemFilter{
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
		},
	}))
}

// func BenchmarkSimulate(b *testing.B) {
// 	rsr := &proto.RaidSimRequest{
// 		Raid: core.SinglePlayerRaidProto(
// 			&proto.Player{
// 				Race:      proto.Race_RaceBloodElf,
// 				Class:     proto.Class_ClassWarlock,
// 				Equipment: Phase4Gear,
// 				Consumes:  FullConsumes,
// 				Spec:      DefaultOptions,
// 				Buffs:     FullIndividualBuffs,
// 			},
// 			FullPartyBuffs,
// 			FullRaidBuffs),
// 		Encounter: &proto.Encounter{
// 			Duration: 300,
// 			Targets: []*proto.Target{
// 				FullDebuffTarget,
// 			},
// 		},
// 		SimOptions: core.AverageDefaultSimTestOptions,
// 	}

// 	core.RaidBenchmark(b, rsr)
// }
