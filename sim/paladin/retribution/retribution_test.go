package retribution

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get item effects included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterRetributionPaladin()
}

func TestRetribution(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassPaladin,

		Race:       proto.Race_RaceBloodElf,
		OtherRaces: []proto.Race{proto.Race_RaceHuman, proto.Race_RaceDraenei, proto.Race_RaceDwarf},

		GearSet: core.GearSetCombo{Label: "P1", GearSet: Phase1Gear},

		SpecOptions: core.SpecOptionsCombo{Label: "Retribution Paladin SOV", SpecOptions: DefaultOptions},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{
				Label: "Retribution Paladin SOC",
				SpecOptions: &proto.Player_RetributionPaladin{
					RetributionPaladin: &proto.RetributionPaladin{
						Talents: defaultRetTalents,
						Options: &proto.RetributionPaladin_Options{
							Judgement:            proto.PaladinJudgement_JudgementOfWisdom,
							Seal:                 proto.PaladinSeal_Command,
							Aura:                 proto.PaladinAura_RetributionAura,
							DamageTakenPerSecond: 0,
						},
						Rotation: defaultRetRotation,
					},
				},
			},
			{
				Label: "Retribution Paladin SOR",
				SpecOptions: &proto.Player_RetributionPaladin{
					RetributionPaladin: &proto.RetributionPaladin{
						Talents: defaultRetTalents,
						Options: &proto.RetributionPaladin_Options{
							Judgement:            proto.PaladinJudgement_JudgementOfWisdom,
							Seal:                 proto.PaladinSeal_Righteousness,
							Aura:                 proto.PaladinAura_RetributionAura,
							DamageTakenPerSecond: 0,
						},
						Rotation: defaultRetRotation,
					},
				},
			},
		},
		Glyphs: defaultRetGlyphs,

		RaidBuffs:   FullRaidBuffs,
		PartyBuffs:  FullPartyBuffs,
		PlayerBuffs: FullIndividualBuffs,
		Consumes:    FullConsumes,
		Debuffs:     FullDebuffs,

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypePolearm,
				proto.WeaponType_WeaponTypeMace,
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
				Glyphs:    defaultRetGlyphs,
				Equipment: Phase1Gear,
				Consumes:  FullConsumes,
				Spec:      DefaultOptions,
				Buffs:     FullIndividualBuffs,
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
