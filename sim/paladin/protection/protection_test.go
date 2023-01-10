package protection

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get item effects included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterProtectionPaladin()
}

func TestProtection(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassPaladin,
		Race:       proto.Race_RaceBloodElf,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Protection Paladin SOV", SpecOptions: DefaultOptions},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{
				Label: "Protection Paladin SOC",
				SpecOptions: &proto.Player_ProtectionPaladin{
					ProtectionPaladin: &proto.ProtectionPaladin{
						Options: &proto.ProtectionPaladin_Options{
							Judgement:            proto.PaladinJudgement_JudgementOfWisdom,
							Seal:                 proto.PaladinSeal_Command,
							Aura:                 proto.PaladinAura_RetributionAura,
							DamageTakenPerSecond: 0,
						},
						Rotation: defaultProtRotation,
					},
				},
			},
			{
				Label: "Protection Paladin SOR",
				SpecOptions: &proto.Player_ProtectionPaladin{
					ProtectionPaladin: &proto.ProtectionPaladin{
						Options: &proto.ProtectionPaladin_Options{
							Judgement:            proto.PaladinJudgement_JudgementOfWisdom,
							Seal:                 proto.PaladinSeal_Righteousness,
							Aura:                 proto.PaladinAura_RetributionAura,
							DamageTakenPerSecond: 0,
						},
						Rotation: defaultProtRotation,
					},
				},
			},
		},

		IsTank:          true,
		InFrontOfTarget: true,

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypePolearm,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeShield,
			},
			ArmorType: proto.ArmorType_ArmorTypePlate,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeLibram,
			},
		},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:      proto.Race_RaceBloodElf,
				Class:     proto.Class_ClassPaladin,
				Equipment: P1Gear,
				Consumes:  FullConsumes,
				Spec:      DefaultOptions,
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
