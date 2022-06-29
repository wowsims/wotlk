package warlock

import (
	"testing"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

func init() {
	RegisterWarlock()
}

func TestDestruction(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassWarlock,

		Race:       proto.Race_RaceBloodElf,
		OtherRaces: []proto.Race{proto.Race_RaceHuman, proto.Race_RaceGnome, proto.Race_RaceOrc, proto.Race_RaceUndead},

		GearSet: core.GearSetCombo{Label: "P4", GearSet: Phase4Gear},

		SpecOptions: core.SpecOptionsCombo{Label: "Destro Warlock", SpecOptions: DefaultDestroWarlock},

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
			ArmorType: proto.ArmorType_ArmorTypePlate,
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
