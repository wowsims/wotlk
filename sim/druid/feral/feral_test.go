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
		Race:  proto.Race_RaceTauren,

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsBearweaveLacerate},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "FlowerWeave", SpecOptions: PlayerOptionsMonoCatFlowerWeave},
		},

		ItemFilter: core.ItemFilter{
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
				Buffs:     core.FullIndividualBuffs,
				Glyphs:    StandardGlyphs,

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
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.DruidMajorGlyph_GlyphOfRip),
	Major2: int32(proto.DruidMajorGlyph_GlyphOfShred),
	Major3: int32(proto.DruidMajorGlyph_GlyphOfBerserk),
	Minor1: int32(proto.DruidMinorGlyph_GlyphOfTheWild),
}

var PlayerOptionsBearweaveLacerate = &proto.Player_FeralDruid{
	FeralDruid: &proto.FeralDruid{
		Options: &proto.FeralDruid_Options{
			InnervateTarget: &proto.RaidTarget{TargetIndex: -1}, // no Innervate
			LatencyMs:       100,
			PrepopOoc:       true,
		},
		Rotation: &proto.FeralDruid_Rotation{
			BearWeaveType:      proto.FeralDruid_Rotation_Lacerate,
			UseRake:            true,
			UseBite:            true,
			MinCombosForRip:    5,
			MinCombosForBite:   5,
			BiteTime:           10.0,
			MaintainFaerieFire: true,
			BerserkBiteThresh:  25.0,
			MinRoarOffset:      14.0,
			SnekWeave:          true,
			FlowerWeave:        false,
			RaidTargets:        30,
		},
	},
}

var PlayerOptionsMonoCatFlowerWeave = &proto.Player_FeralDruid{
	FeralDruid: &proto.FeralDruid{
		Options: &proto.FeralDruid_Options{
			InnervateTarget: &proto.RaidTarget{TargetIndex: -1}, // no Innervate
			LatencyMs:       100,
			PrepopOoc:       true,
		},
		Rotation: &proto.FeralDruid_Rotation{
			BearWeaveType:      proto.FeralDruid_Rotation_None,
			UseRake:            true,
			UseBite:            true,
			MinCombosForRip:    5,
			MinCombosForBite:   5,
			BiteTime:           10.0,
			MaintainFaerieFire: true,
			BerserkBiteThresh:  25.0,
			MinRoarOffset:      14.0,
			SnekWeave:          false,
			FlowerWeave:        true,
			RaidTargets:        30,
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

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40473,
		"enchant": 3817,
		"gems": [
			41398,
			39996
		]
	},
	{
		"id": 44664,
		"gems": [
			39996
		]
	},
	{
		"id": 40494,
		"enchant": 3808,
		"gems": [
			39996
		]
	},
	{
		"id": 40403,
		"enchant": 3605
	},
	{
		"id": 40539,
		"enchant": 3832,
		"gems": [
			39996
		]
	},
	{
		"id": 39765,
		"enchant": 3845,
		"gems": [
			39996,
			0
		]
	},
	{
		"id": 40541,
		"enchant": 3604,
		"gems": [
			0
		]
	},
	{
		"id": 40205,
		"gems": [
			39996
		]
	},
	{
		"id": 44011,
		"enchant": 3823,
		"gems": [
			39996,
			49110
		]
	},
	{
		"id": 40243,
		"enchant": 3606,
		"gems": [
			40014
		]
	},
	{
		"id": 40474
	},
	{
		"id": 40717
	},
	{
		"id": 42987
	},
	{
		"id": 40256
	},
	{
		"id": 40388,
		"enchant": 3789
	},
	{},
	{
		"id": 39757
	}
]}`)
