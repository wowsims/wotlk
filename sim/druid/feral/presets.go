package feral

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var StandardTalents = &proto.DruidTalents{
	Ferocity:                5,
	FeralInstinct:           3,
	SavageFury:              2,
	FeralSwiftness:          2,
	SurvivalInstincts:       true,
	SharpenedClaws:          3,
	ShreddingAttacks:        2,
	PredatoryStrikes:        3,
	PrimalFury:              2,
	PrimalPrecision:         2,
	HeartOfTheWild:          5,
	SurvivalOfTheFittest:    3,
	LeaderOfThePack:         true,
	ImprovedLeaderOfThePack: 2,
	ProtectorOfThePack:      2,
	PredatoryInstincts:      3,
	KingOfTheJungle:         3,
	Mangle:                  true,
	RendAndTear:             5,
	PrimalGore:              true,
	Berserk:                 true,
	ImprovedMarkOfTheWild:   2,
	Furor:                   5,
	Naturalist:              5,
	NaturalShapeshifter:     3,
	MasterShapeshifter:      2,
	OmenOfClarity:           true,
}

var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.DruidMajorGlyph_GlyphOfRip),
	Major2: int32(proto.DruidMajorGlyph_GlyphOfShred),
	Major3: int32(proto.DruidMajorGlyph_GlyphOfBerserk),
	// No interesting minor glyphs.
}

var PlayerOptionsBearweaveLacerate = &proto.Player_FeralDruid{
	FeralDruid: &proto.FeralDruid{
		Talents: StandardTalents,
		Options: &proto.FeralDruid_Options{
			InnervateTarget: &proto.RaidTarget{TargetIndex: -1}, // no Innervate
			LatencyMs:       100,
		},
		Rotation: &proto.FeralDruid_Rotation{
			BearWeaveType:      proto.FeralDruid_Rotation_Lacerate,
			UseRake:            true,
			UseBite:            true,
			MinCombosForRip:    5,
			MinCombosForBite:   5,
			BiteTime:           10.0,
			MaintainFaerieFire: true,
			BerserkBiteThresh:  30.0,
			MaxRoarOffset:      14.0,
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

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40473,
		"enchant": 44879,
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
		"enchant": 44871,
		"gems": [
			39996
		]
	},
	{
		"id": 40403,
		"enchant": 55002
	},
	{
		"id": 40539,
		"enchant": 44489,
		"gems": [
			39996
		]
	},
	{
		"id": 39765,
		"enchant": 44484,
		"gems": [
			39996,
			0
		]
	},
	{
		"id": 40541,
		"enchant": 54999,
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
		"enchant": 38374,
		"gems": [
			39996,
			49110
		]
	},
	{
		"id": 40243,
		"enchant": 55016,
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
		"enchant": 44492
	},
	{},
	{
		"id": 39757
	}
]}`)
