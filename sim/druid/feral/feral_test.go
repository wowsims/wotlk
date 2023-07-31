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

		GearSet: core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		OtherGearSets: []core.GearSetCombo{
			{Label: "P2", GearSet: P2Gear},
		},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsMonoCat},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "Default-NoBleed", SpecOptions: PlayerOptionsMonoCatNoBleed},
			{Label: "Flower-Aoe", SpecOptions: PlayerOptionsFlowerCatAoe},
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

func TestFeralDoubleArmorPenTrinketsNoDesync(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:       proto.Class_ClassDruid,
		Race:        proto.Race_RaceTauren,
		GearSet:     core.GearSetCombo{Label: "P2DoubleArmorPenTrinkets", GearSet: P2GearDoubleArmorPenTrinkets},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsMonoCat},

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

		Cooldowns: &proto.Cooldowns{
			DesyncProcTrinket1Seconds: 0,
			DesyncProcTrinket2Seconds: 0,
		},
	}))
}

func TestFeralDoubleArmorPenTrinketsWithDesync(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:       proto.Class_ClassDruid,
		Race:        proto.Race_RaceTauren,
		GearSet:     core.GearSetCombo{Label: "P2DoubleArmorPenTrinkets", GearSet: P2GearDoubleArmorPenTrinkets},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsMonoCat},

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

		Cooldowns: &proto.Cooldowns{
			DesyncProcTrinket1Seconds: 0,
			DesyncProcTrinket2Seconds: 10,
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
				Spec:      PlayerOptionsMonoCat,
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
	Major1: int32(proto.DruidMajorGlyph_GlyphOfOmenOfClarity),
	Major2: int32(proto.DruidMajorGlyph_GlyphOfShred),
	Major3: int32(proto.DruidMajorGlyph_GlyphOfBerserk),
	Minor1: int32(proto.DruidMinorGlyph_GlyphOfTheWild),
}

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

var P2Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 46161,
		"enchant": 3817,
		"gems": [
			41398,
			40002
		]
	},
	{
		"id": 45517,
		"gems": [
			40002
		]
	},
	{
		"id": 45245,
		"enchant": 3808,
		"gems": [
			40002,
			40002
		]
	},
	{
		"id": 46032,
		"enchant": 3605,
		"gems": [
			40002,
			40058
		]
	},
	{
		"id": 45473,
		"enchant": 3832,
		"gems": [
			40002,
			40002,
			40002
		]
	},
	{
		"id": 45869,
		"enchant": 3845,
		"gems": [
			40037
		]
	},
	{
		"id": 46158,
		"enchant": 3604,
		"gems": [
			40002
		]
	},
	{
		"id": 46095,
		"gems": [
			40002,
			40002,
			40002
		]
	},
	{
		"id": 45536,
		"enchant": 3823,
		"gems": [
			39996,
			39996,
			39996
		]
	},
	{
		"id": 45564,
		"enchant": 3606,
		"gems": [
			39996,
			39996
		]
	},
	{
		"id": 46048,
		"gems": [
			45862
		]
	},
	{
		"id": 45608,
		"gems": [
			39996
		]
	},
	{
		"id": 45931
	},
	{
		"id": 46038
	},
	{
		"id": 45613,
		"enchant": 3789,
		"gems": [
			40037,
			42702
		]
	},
	{},
	{
		"id": 40713
	}
]}`)

var P2GearDoubleArmorPenTrinkets = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 46161,
		"enchant": 3817,
		"gems": [
			41398,
			40002
		]
	},
	{
		"id": 45517,
		"gems": [
			40002
		]
	},
	{
		"id": 45245,
		"enchant": 3808,
		"gems": [
			40002,
			40002
		]
	},
	{
		"id": 46032,
		"enchant": 3605,
		"gems": [
			40002,
			40058
		]
	},
	{
		"id": 45473,
		"enchant": 3832,
		"gems": [
			40002,
			40002,
			40002
		]
	},
	{
		"id": 45869,
		"enchant": 3845,
		"gems": [
			40037
		]
	},
	{
		"id": 46158,
		"enchant": 3604,
		"gems": [
			40002
		]
	},
	{
		"id": 46095,
		"gems": [
			40002,
			40002,
			40002
		]
	},
	{
		"id": 45536,
		"enchant": 3823,
		"gems": [
			39996,
			39996,
			39996
		]
	},
	{
		"id": 45564,
		"enchant": 3606,
		"gems": [
			39996,
			39996
		]
	},
	{
		"id": 46048,
		"gems": [
			45862
		]
	},
	{
		"id": 45608,
		"gems": [
			39996
		]
	},
	{
		"id": 45931
	},
	{
		"id": 40256
	},
	{
		"id": 45613,
		"enchant": 3789,
		"gems": [
			40037,
			42702
		]
	},
	{},
	{
		"id": 40713
	}
]}`)
