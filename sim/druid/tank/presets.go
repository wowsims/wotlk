package tank

import (
	"github.com/wowsims/wotlk/sim/core/items"
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
	PredatoryInstincts:      5,
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
			MaulRageThreshold:        50,
			MaintainDemoralizingRoar: true,
			Swipe:                    proto.FeralTankDruid_Rotation_SwipeWithEnoughAP,
			SwipeApThreshold:         2700,
		},
	},
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance:     true,
	GiftOfTheWild:        proto.TristateEffect_TristateEffectImproved,
	Bloodlust:            true,
	BattleShout:          proto.TristateEffect_TristateEffectImproved,
	LeaderOfThePack:      proto.TristateEffect_TristateEffectImproved,
	ManaSpringTotem:      proto.TristateEffect_TristateEffectRegular,
	StrengthOfEarthTotem: proto.TristateEffect_TristateEffectImproved,
	UnleashedRage:        true,
}
var FullPartyBuffs = &proto.PartyBuffs{
	BraidedEterniumChain: true,
}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings: true,
	BlessingOfMight: proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	BattleElixir:    proto.BattleElixir_ElixirOfMajorAgility,
	Food:            proto.Food_FoodGrilledMudfish,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var FullDebuffs = &proto.Debuffs{
	JudgementOfWisdom: true,
	BloodFrenzy:       true,
	GiftOfArthas:      true,
	FaerieFire:        proto.TristateEffect_TristateEffectImproved,
	SunderArmor:       true,
	CurseOfWeakness:   proto.TristateEffect_TristateEffectImproved,
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 34404,
		"enchant": 29193,
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
		"enchant": 28911,
		"gems": [
			32200,
			32200
		]
	},
	{
		"id": 34190,
		"enchant": 34004
	},
	{
		"id": 34211,
		"enchant": 24003,
		"gems": [
			32200,
			32200,
			32200
		]
	},
	{
		"id": 34444,
		"enchant": 22533,
		"gems": [
			32200,
			0
		]
	},
	{
		"id": 34408,
		"enchant": 33153,
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
		"enchant": 29536,
		"gems": [
			32200,
			32200,
			32200
		]
	},
	{
		"id": 34573,
		"enchant": 35297,
		"gems": [
			32200
		]
	},
	{
		"id": 34213,
		"enchant": 22538
	},
	{
		"id": 34361,
		"enchant": 22538
	},
	{
		"id": 32501
	},
	{
		"id": 32658
	},
	{
		"id": 30883,
		"enchant": 22556
	},
	{},
	{
		"id": 32387
	}
]}`)
