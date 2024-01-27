package retribution

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get item effects included.
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterRetributionPaladin()
}

func TestRetribution(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassPaladin,
		Race:       proto.Race_RaceHuman,
		OtherRaces: []proto.Race{proto.Race_RaceHuman, proto.Race_RaceDwarf},

		GearSet:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p1"),
		Talents:     StandardTalents,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Retribution Paladin SOV", SpecOptions: DefaultOptions},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{
				Label: "Retribution Paladin SOC",
				SpecOptions: &proto.Player_RetributionPaladin{
					RetributionPaladin: &proto.RetributionPaladin{
						Options: &proto.RetributionPaladin_Options{
							Judgement: proto.PaladinJudgement_JudgementOfWisdom,
							Seal:      proto.PaladinSeal_Command,
							Aura:      proto.PaladinAura_RetributionAura,
						},
						Rotation: &proto.RetributionPaladin_Rotation{},
					},
				},
			},
			{
				Label: "Retribution Paladin SOR",
				SpecOptions: &proto.Player_RetributionPaladin{
					RetributionPaladin: &proto.RetributionPaladin{
						Options: &proto.RetributionPaladin_Options{
							Judgement: proto.PaladinJudgement_JudgementOfWisdom,
							Seal:      proto.PaladinSeal_Righteousness,
							Aura:      proto.PaladinAura_RetributionAura,
						},
						Rotation: &proto.RetributionPaladin_Rotation{},
					},
				},
			},
			{
				Label: "Retribution Paladin SOV 2 Target Swapping",
				SpecOptions: &proto.Player_RetributionPaladin{
					RetributionPaladin: &proto.RetributionPaladin{
						Options: &proto.RetributionPaladin_Options{
							Judgement: proto.PaladinJudgement_JudgementOfWisdom,
							Seal:      proto.PaladinSeal_Vengeance,
							Aura:      proto.PaladinAura_RetributionAura,
						},
						Rotation: &proto.RetributionPaladin_Rotation{
							ConsSlack:            500,
							ExoSlack:             500,
							UseDivinePlea:        true,
							DivinePleaPercentage: 0.75,
							HolyWrathThreshold:   4,
							SovTargets:           2,
						},
					},
				},
			},
		},
		Rotation: core.GetAplRotation("../../../ui/retribution_paladin/apls", "default"),

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
				Race:          proto.Race_RaceHuman,
				Class:         proto.Class_ClassPaladin,
				TalentsString: StandardTalents,
				Equipment:     core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "p1").GearSet,
				Consumes:      FullConsumes,
				Spec:          DefaultOptions,
				Buffs:         core.FullIndividualBuffs,
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

var StandardTalents = "050501-05-05232051203331302133231331"

var DefaultOptions = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options: &proto.RetributionPaladin_Options{
			Judgement: proto.PaladinJudgement_JudgementOfWisdom,
			Seal:      proto.PaladinSeal_Vengeance,
			Aura:      proto.PaladinAura_RetributionAura,
		},
		Rotation: &proto.RetributionPaladin_Rotation{},
	},
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfRelentlessAssault,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
	Food:            proto.Food_FoodRoastedClefthoof,
	ThermalSapper:   true,
}
