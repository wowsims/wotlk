package tank

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
	CurseOfRecklessness:         true,
	HuntersMark:                 proto.TristateEffect_TristateEffectImproved,
	ExposeWeaknessUptime:        1.0,
	ExposeWeaknessHunterAgility: 1000,
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 29098,
		"enchant": 29192,
		"gems": [
			24067,
			32409
		]
	},
	{
		"id": 28509
	},
	{
		"id": 29100,
		"enchant": 28911,
		"gems": [
			24033,
			24033
		]
	},
	{
		"id": 28660,
		"enchant": 34004
	},
	{
		"id": 29096,
		"enchant": 24003,
		"gems": [
			24067,
			24055,
			24055
		]
	},
	{
		"id": 28978,
		"enchant": 22533,
		"gems": [
			24033
		]
	},
	{
		"id": 29097,
		"enchant": 33153
	},
	{
		"id": 28986
	},
	{
		"id": 29099,
		"enchant": 29536
	},
	{
		"id": 30674,
		"enchant": 35297
	},
	{
		"id": 29279,
		"enchant": 22535
	},
	{
		"id": 28792,
		"enchant": 22535
	},
	{
		"id": 28830
	},
	{
		"id": 23836
	},
	{
		"id": 28476,
		"enchant": 22556
	},
	{
		"id": 23198
	}
]}`)
