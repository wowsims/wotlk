package feral

import (
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
)

var StandardTalents = &proto.DruidTalents{
	Ferocity:                5,
	SharpenedClaws:          3,
	ShreddingAttacks:        2,
	PredatoryStrikes:        3,
	PrimalFury:              2,
	SavageFury:              2,
	FaerieFire:              true,
	HeartOfTheWild:          5,
	SurvivalOfTheFittest:    3,
	LeaderOfThePack:         true,
	ImprovedLeaderOfThePack: 2,
	PredatoryInstincts:      5,
	Mangle:                  true,
	Furor:                   5,
	Naturalist:              5,
	NaturalShapeshifter:     3,
	Intensity:               3,
	OmenOfClarity:           true,
}

var PlayerOptionsBiteweave = &proto.Player_FeralDruid{
	FeralDruid: &proto.FeralDruid{
		Talents: StandardTalents,
		Options: &proto.FeralDruid_Options{
			InnervateTarget: &proto.RaidTarget{TargetIndex: -1}, // no Innervate
			LatencyMs:       100,
		},
		Rotation: &proto.FeralDruid_Rotation{
			FinishingMove:      proto.FeralDruid_Rotation_Rip,
			MangleTrick:        true,
			Biteweave:          true,
			RipMinComboPoints:  5,
			BiteMinComboPoints: 5,
			RakeTrick:          false,
			Ripweave:           false,
			MaintainFaerieFire: true,
		},
	},
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance: true,
	GiftOfTheWild:    proto.TristateEffect_TristateEffectImproved,
}
var FullPartyBuffs = &proto.PartyBuffs{
	Bloodlust:                  1,
	Drums:                      proto.Drums_DrumsOfBattle,
	BattleShout:                proto.TristateEffect_TristateEffectImproved,
	LeaderOfThePack:            proto.TristateEffect_TristateEffectImproved,
	GraceOfAirTotem:            proto.TristateEffect_TristateEffectImproved,
	ManaSpringTotem:            proto.TristateEffect_TristateEffectRegular,
	BraidedEterniumChain:       true,
	StrengthOfEarthTotem:       proto.StrengthOfEarthType_EnhancingTotems,
	SnapshotBsSolarianSapphire: true,
	SanctityAura:               proto.TristateEffect_TristateEffectImproved,
}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings: true,
	BlessingOfMight: proto.TristateEffect_TristateEffectImproved,
	UnleashedRage:   true,
}

var FullConsumes = &proto.Consumes{
	BattleElixir:     proto.BattleElixir_ElixirOfMajorAgility,
	GuardianElixir:   proto.GuardianElixir_ElixirOfMajorMageblood,
	Food:             proto.Food_FoodGrilledMudfish,
	DefaultPotion:    proto.Potions_HastePotion,
	MainHandImbue:    proto.WeaponImbue_WeaponImbueAdamantiteWeightstone,
	DefaultConjured:  proto.Conjured_ConjuredDarkRune,
	ScrollOfStrength: 5,
	ScrollOfAgility:  5,
}

var FullDebuffs = &proto.Debuffs{
	JudgementOfWisdom:           true,
	ImprovedSealOfTheCrusader:   true,
	BloodFrenzy:                 true,
	GiftOfArthas:                true,
	ExposeArmor:                 proto.TristateEffect_TristateEffectImproved,
	FaerieFire:                  proto.TristateEffect_TristateEffectImproved,
	SunderArmor:                 true,
	Mangle:                      true,
	CurseOfRecklessness:         true,
	HuntersMark:                 proto.TristateEffect_TristateEffectImproved,
	ExposeWeaknessUptime:        1.0,
	ExposeWeaknessHunterAgility: 1000,
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 8345,
		"enchant": 29192
	},
	{
		"id": 29381
	},
	{
		"id": 29100,
		"enchant": 28888,
		"gems": [
			24028,
			24028
		]
	},
	{
		"id": 28672,
		"enchant": 34004
	},
	{
		"id": 29096,
		"enchant": 24003,
		"gems": [
			24028,
			24028,
			24028
		]
	},
	{
		"id": 29246,
		"enchant": 27899
	},
	{
		"id": 28506,
		"gems": [
			24028,
			24028
		]
	},
	{
		"id": 28750,
		"gems": [
			24028,
			24028
		]
	},
	{
		"id": 28741,
		"enchant": 29535,
		"gems": [
			24028,
			24028,
			24028
		]
	},
	{
		"id": 28545,
		"enchant": 28279,
		"gems": [
			24028,
			24028
		]
	},
	{
		"id": 28649,
		"enchant": 22535
	},
	{
		"id": 30834,
		"enchant": 22535
	},
	{
		"id": 28830
	},
	{
		"id": 29383
	},
	{
		"id": 28658,
		"enchant": 22556
	},
	{
		"id": 29390
	}
]}`)
