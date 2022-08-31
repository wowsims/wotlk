package protection

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get item effects included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterProtectionWarrior()
}

func TestProtectionWarrior(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassWarrior,

		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet: core.GearSetCombo{Label: "P1", GearSet: P1Gear},

		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},

		Glyphs:      DefaultGlyphs,
		RaidBuffs:   FullRaidBuffs,
		PartyBuffs:  FullPartyBuffs,
		PlayerBuffs: FullIndividualBuffs,
		Consumes:    FullConsumes,
		Debuffs:     FullDebuffs,

		IsTank:          true,
		InFrontOfTarget: true,

		ItemFilter: core.ItemFilter{
			ArmorType: proto.ArmorType_ArmorTypePlate,

			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeFist,
				proto.WeaponType_WeaponTypeShield,
			},
		},

		EPReferenceStat: proto.Stat_StatAttackPower,
		StatsToWeigh: []proto.Stat{
			proto.Stat_StatStrength,
			proto.Stat_StatAttackPower,
			proto.Stat_StatArmor,
			proto.Stat_StatDodge,
			proto.Stat_StatBlockValue,
		},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:      proto.Race_RaceOrc,
				Class:     proto.Class_ClassWarrior,
				Equipment: P1Gear,
				Consumes:  FullConsumes,
				Spec:      PlayerOptionsBasic,
				Buffs:     FullIndividualBuffs,

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
