package tank

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var StandardTalents = &proto.DruidTalents{
	Ferocity:                5,
	SharpenedClaws:          3,
	ShreddingAttacks:        2,
	PredatoryStrikes:        3,
	PrimalFury:              2,
	SavageFury:              2,
	HeartOfTheWild:          5,
	SurvivalOfTheFittest:    3,
	LeaderOfThePack:         true,
	ImprovedLeaderOfThePack: 2,
	PredatoryInstincts:      3,
	Mangle:                  true,

	Furor:               5,
	Naturalist:          5,
	NaturalShapeshifter: 3,
	Intensity:           3,
	OmenOfClarity:       true,
}

var PlayerOptionsDefault = &proto.Player_FeralTankDruid{
	FeralTankDruid: &proto.FeralTankDruid{
		Talents: StandardTalents,
		Options: &proto.FeralTankDruid_Options{
			InnervateTarget: &proto.RaidTarget{TargetIndex: -1}, // no Innervate
			StartingRage:    20,
		},
		Rotation: &proto.FeralTankDruid_Rotation{
			MaulRageThreshold:        60,
			MaintainDemoralizingRoar: true,
			LacerateTime:             5.0,
		},
	},
}

var FullConsumes = &proto.Consumes{
	BattleElixir:    proto.BattleElixir_ElixirOfMajorAgility,
	Food:            proto.Food_FoodGrilledMudfish,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 34404,
		"enchant": 3004,
		"gems": [
			32212,
			25896
		]
	},
	{
		"id": 34178
	},
	{
		"id": 34392,
		"enchant": 2991,
		"gems": [
			32200,
			32200
		]
	},
	{
		"id": 34190,
		"enchant": 368
	},
	{
		"id": 34211,
		"enchant": 2661,
		"gems": [
			32200,
			32200,
			32200
		]
	},
	{
		"id": 34444,
		"enchant": 2649,
		"gems": [
			32200,
			0
		]
	},
	{
		"id": 34408,
		"enchant": 2613,
		"gems": [
			32200,
			32200,
			0
		]
	},
	{
		"id": 35156,
		"gems": [
			0
		]
	},
	{
		"id": 34385,
		"enchant": 3013,
		"gems": [
			32200,
			32200,
			32200
		]
	},
	{
		"id": 34573,
		"enchant": 2940,
		"gems": [
			32200
		]
	},
	{
		"id": 34213,
		"enchant": 2931
	},
	{
		"id": 34361,
		"enchant": 2931
	},
	{
		"id": 32501
	},
	{
		"id": 32658
	},
	{
		"id": 30883,
		"enchant": 2670
	},
	{},
	{
		"id": 32387
	}
]}`)
