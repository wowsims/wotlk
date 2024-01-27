package protection

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get item effects included.
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterProtectionPaladin()
}

func TestProtection(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassPaladin,
		Race:       proto.Race_RaceHuman,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet:     core.GetGearSet("../../../ui/protection_paladin/gear_sets", "p1"),
		Talents:     StandardTalents,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Protection Paladin SOV", SpecOptions: DefaultOptions},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{
				Label: "Protection Paladin SOC",
				SpecOptions: &proto.Player_ProtectionPaladin{
					ProtectionPaladin: &proto.ProtectionPaladin{
						Options: &proto.ProtectionPaladin_Options{
							Judgement: proto.PaladinJudgement_JudgementOfWisdom,
							Seal:      proto.PaladinSeal_Command,
							Aura:      proto.PaladinAura_RetributionAura,
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
							Judgement: proto.PaladinJudgement_JudgementOfWisdom,
							Seal:      proto.PaladinSeal_Righteousness,
							Aura:      proto.PaladinAura_RetributionAura,
						},
						Rotation: defaultProtRotation,
					},
				},
			},
		},
		Rotation: core.GetAplRotation("../../../ui/protection_paladin/apls", "default"),

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
				Race:      proto.Race_RaceHuman,
				Class:     proto.Class_ClassPaladin,
				Equipment: core.GetGearSet("../../../ui/protection_paladin/gear_sets", "p1").GearSet,
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

var StandardTalents = "-05005135200132311333312321-511302012003"

var defaultProtRotation = &proto.ProtectionPaladin_Rotation{}

var defaultProtOptions = &proto.ProtectionPaladin_Options{
	Judgement: proto.PaladinJudgement_JudgementOfWisdom,
	Seal:      proto.PaladinSeal_Vengeance,
	Aura:      proto.PaladinAura_RetributionAura,
}

var DefaultOptions = &proto.Player_ProtectionPaladin{
	ProtectionPaladin: &proto.ProtectionPaladin{
		Options:  defaultProtOptions,
		Rotation: defaultProtRotation,
	},
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfStoneblood,
	Food:            proto.Food_FoodDragonfinFilet,
	DefaultPotion:   proto.Potions_IndestructiblePotion,
	PrepopPotion:    proto.Potions_IndestructiblePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}
