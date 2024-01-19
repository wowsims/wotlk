package holy

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get item effects included.
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterHolyPaladin()
}

func TestHoly(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassPaladin,
		Race:       proto.Race_RaceHuman,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet:     core.GetGearSet("../../../ui/holy_paladin/gear_sets", "p1"),
		Talents:     StandardTalents,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: BasicOptions},
		Rotation:    core.RotationCombo{Label: "Default", Rotation: DefaultRotation},

		IsHealer:        true,
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
				Race:          proto.Race_RaceHuman,
				Class:         proto.Class_ClassPaladin,
				Equipment:     core.GetGearSet("../../../ui/holy_paladin/gear_sets", "p1").GearSet,
				Consumes:      FullConsumes,
				Spec:          BasicOptions,
				TalentsString: StandardTalents,
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

var StandardTalents = "50350151020013053100515221-50023131203"

var defaultProtOptions = &proto.HolyPaladin_Options{
	Judgement: proto.PaladinJudgement_JudgementOfWisdom,
	Aura:      proto.PaladinAura_DevotionAura,
}

var BasicOptions = &proto.Player_HolyPaladin{
	HolyPaladin: &proto.HolyPaladin{
		Options:  defaultProtOptions,
		Rotation: &proto.HolyPaladin_Rotation{},
	},
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfStoneblood,
	Food:            proto.Food_FoodDragonfinFilet,
	DefaultPotion:   proto.Potions_IndestructiblePotion,
	PrepopPotion:    proto.Potions_IndestructiblePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var DefaultRotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"priorityList": [
		{"action":{"autocastOtherCooldowns":{}}}
	]
}`)
