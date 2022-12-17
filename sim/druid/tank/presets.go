package tank

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var StandardTalents = &proto.DruidTalents{
	Ferocity:                5,
	FeralInstinct:           3,
	SavageFury:              2,
	ThickHide:               3,
	FeralSwiftness:          2,
	SurvivalInstincts:       true,
	SharpenedClaws:          3,
	ShreddingAttacks:        2,
	PredatoryStrikes:        3,
	PrimalFury:              2,
	PrimalPrecision:         2,
	FeralCharge:             true,
	NaturalReaction:         3,
	HeartOfTheWild:          5,
	SurvivalOfTheFittest:    3,
	LeaderOfThePack:         true,
	ImprovedLeaderOfThePack: 2,
	ProtectorOfThePack:      3,
	KingOfTheJungle:         3,
	Mangle:                  true,
	ImprovedMangle:          3,
	RendAndTear:             5,
	PrimalGore:              true,
	Berserk:                 true,

	ImprovedMarkOfTheWild: 2,
	Furor:                 3,
	Naturalist:            5,
	OmenOfClarity:         true,
}

var PlayerOptionsDefault = &proto.Player_FeralTankDruid{
	FeralTankDruid: &proto.FeralTankDruid{
		Talents: StandardTalents,
		Options: &proto.FeralTankDruid_Options{
			InnervateTarget: &proto.RaidTarget{TargetIndex: -1}, // no Innervate
			StartingRage:    20,
		},
		Rotation: &proto.FeralTankDruid_Rotation{
			MaulRageThreshold:        25,
			MaintainDemoralizingRoar: true,
			LacerateTime:             8.0,
		},
	},
}

var FullConsumes = &proto.Consumes{
	BattleElixir:    proto.BattleElixir_GurusElixir,
	GuardianElixir:  proto.GuardianElixir_GiftOfArthas,
	Food:            proto.Food_FoodBlackenedDragonfin,
	DefaultPotion:   proto.Potions_IndestructiblePotion,
	DefaultConjured: proto.Conjured_ConjuredHealthstone,
	ThermalSapper:   true,
	FillerExplosive: proto.Explosive_ExplosiveSaroniteBomb,
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
		{
			"id": 40329,
			"enchant": 3878,
			"gems": [
				41339,
				40008
			]
		},
		{
			"id": 40387
		},
		{
			"id": 40494,
			"enchant": 3852,
			"gems": [
				40008
			]
		},
		{
			"id": 40252,
			"enchant": 3294
		},
		{
			"id": 40471,
			"enchant": 3832,
			"gems": [
				42702,
				40088
			]
		},
		{
			"id": 40186,
			"enchant": 3850,
			"gems": [
				40008,
				0
			]
		},
		{
			"id": 40472,
			"enchant": 3860,
			"gems": [
				40008,
				0
			]
		},
		{
			"id": 43591,
			"gems": [
				40008,
				40008,
				40008
			]
		},
		{
			"id": 44011,
			"enchant": 3822,
			"gems": [
				40008,
				40008
			]
		},
		{
			"id": 40243,
			"enchant": 3606,
			"gems": [
				40008
			]
		},
		{
			"id": 40370
		},
		{
			"id": 37784
		},
		{
			"id": 44253
		},
		{
			"id": 37220
		},
		{
			"id": 40280,
			"enchant": 2673
		},
		{},
		{
			"id": 38365
		}
]}`)
