package feral

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterFeralDruid()
}

var FeralItemFilter = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeOffHand,
		proto.WeaponType_WeaponTypeStaff,
		proto.WeaponType_WeaponTypePolearm,
	},
	ArmorType: proto.ArmorType_ArmorTypeLeather,
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeIdol,
	},
}

func TestFeral(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDruid,
		Race:  proto.Race_RaceTauren,

		GearSet:     core.GetGearSet("../../../ui/feral_druid/gear_sets", "p1"),
		Talents:     StandardTalents,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsMonoCat},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "Default-NoBleed", SpecOptions: PlayerOptionsMonoCatNoBleed},
			{Label: "Flower-Aoe", SpecOptions: PlayerOptionsFlowerCatAoe},
		},
		ItemFilter: FeralItemFilter,
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:      proto.Race_RaceTauren,
				Class:     proto.Class_ClassDruid,
				Equipment: core.GetGearSet("../../../ui/feral_druid/gear_sets", "p1").GearSet,
				Consumes:  FullConsumes,
				Spec:      PlayerOptionsMonoCat,
				Buffs:     core.FullIndividualBuffs,

				InFrontOfTarget: true,
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

var StandardTalents = "-503202132322010053120230310511-205503012"

var PlayerOptionsMonoCat = &proto.Player_FeralDruid{
	FeralDruid: &proto.FeralDruid{
		Options: &proto.FeralDruid_Options{
			InnervateTarget:   &proto.UnitReference{}, // no Innervate
			LatencyMs:         100,
			AssumeBleedActive: true,
		},
		Rotation: &proto.FeralDruid_Rotation{
			RotationType:       proto.FeralDruid_Rotation_SingleTarget,
			BearWeaveType:      proto.FeralDruid_Rotation_None,
			UseRake:            true,
			UseBite:            true,
			MinCombosForRip:    5,
			MinCombosForBite:   5,
			BiteTime:           4.0,
			MaintainFaerieFire: true,
			BerserkBiteThresh:  25.0,
			BerserkFfThresh:    15.0,
			MaxFfDelay:         0.7,
			MinRoarOffset:      24.0,
			RipLeeway:          3,
			SnekWeave:          false,
			FlowerWeave:        false,
			RaidTargets:        30,
			PrePopOoc:          true,
		},
	},
}

var PlayerOptionsMonoCatNoBleed = &proto.Player_FeralDruid{
	FeralDruid: &proto.FeralDruid{
		Options: &proto.FeralDruid_Options{
			InnervateTarget:   &proto.UnitReference{}, // no Innervate
			LatencyMs:         100,
			AssumeBleedActive: false,
		},
		Rotation: &proto.FeralDruid_Rotation{
			RotationType:       proto.FeralDruid_Rotation_SingleTarget,
			BearWeaveType:      proto.FeralDruid_Rotation_None,
			UseRake:            true,
			UseBite:            true,
			MinCombosForRip:    5,
			MinCombosForBite:   5,
			BiteTime:           4.0,
			MaintainFaerieFire: true,
			BerserkBiteThresh:  25.0,
			BerserkFfThresh:    15.0,
			MaxFfDelay:         0.7,
			MinRoarOffset:      24.0,
			RipLeeway:          3,
			SnekWeave:          false,
			FlowerWeave:        false,
			RaidTargets:        30,
			PrePopOoc:          true,
		},
	},
}

var PlayerOptionsFlowerCatAoe = &proto.Player_FeralDruid{
	FeralDruid: &proto.FeralDruid{
		Options: &proto.FeralDruid_Options{
			InnervateTarget:   &proto.UnitReference{}, // no Innervate
			LatencyMs:         100,
			AssumeBleedActive: false,
		},
		Rotation: &proto.FeralDruid_Rotation{
			RotationType:       proto.FeralDruid_Rotation_Aoe,
			BearWeaveType:      proto.FeralDruid_Rotation_None,
			UseRake:            true,
			UseBite:            true,
			MinCombosForRip:    5,
			MinCombosForBite:   5,
			BiteTime:           4.0,
			MaintainFaerieFire: true,
			BerserkBiteThresh:  25.0,
			BerserkFfThresh:    15.0,
			MaxFfDelay:         0.7,
			MinRoarOffset:      24.0,
			RipLeeway:          3,
			SnekWeave:          false,
			FlowerWeave:        true,
			RaidTargets:        30,
			PrePopOoc:          true,
		},
	},
}

var FullConsumes = &proto.Consumes{
	BattleElixir:    proto.BattleElixir_ElixirOfMajorAgility,
	GuardianElixir:  proto.GuardianElixir_ElixirOfMajorMageblood,
	Food:            proto.Food_FoodGrilledMudfish,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var P2GearDoubleArmorPenTrinkets = core.EquipmentSpecFromJsonString(`
{
	"items": [
	{"id":46161,"enchant":3817},
	{"id":45517},
	{"id":45245,"enchant":3808},
	{"id":46032,"enchant":3605},
	{"id":45473,"enchant":3832},
	{"id":45869,"enchant":3845},
	{"id":46158,"enchant":3604},
	{"id":46095},
	{"id":45536,"enchant":3823},
	{"id":45564,"enchant":3606},
	{"id":46048},
	{"id":45608},
	{"id":45931},
	{"id":40256},
	{"id":45613,"enchant":3789},
	{},
	{"id":40713}
  ]
}`)
