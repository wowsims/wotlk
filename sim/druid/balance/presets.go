package balance

import (
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
)

var StandardTalents = &proto.DruidTalents{
	StarlightWrath:        5,
	FocusedStarlight:      2,
	ImprovedMoonfire:      2,
	Brambles:              3,
	InsectSwarm:           true,
	Vengeance:             5,
	LunarGuidance:         3,
	NaturesGrace:          true,
	Moonglow:              3,
	Moonfury:              5,
	BalanceOfPower:        2,
	Dreamstate:            3,
	MoonkinForm:           true,
	ImprovedFaerieFire:    3,
	WrathOfCenarius:       5,
	ForceOfNature:         true,
	ImprovedMarkOfTheWild: 5,
	Furor:                 2,
	NaturalShapeshifter:   3,
	Intensity:             3,
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance: true,
	GiftOfTheWild:    proto.TristateEffect_TristateEffectImproved,
}
var FullPartyBuffs = &proto.PartyBuffs{
	MoonkinAura: proto.TristateEffect_TristateEffectRegular,
	Drums:       proto.Drums_DrumsOfBattle,
}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:  true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
	ShadowPriestDps:  500,
}

var FullConsumes = &proto.Consumes{
	Flask:              proto.Flask_FlaskOfBlindingLight,
	Food:               proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:      proto.Potions_SuperManaPotion,
	StartingPotion:     proto.Potions_DestructionPotion,
	MainHandImbue:      proto.WeaponImbue_WeaponImbueBrilliantWizardOil,
	NumStartingPotions: 1,
	DefaultConjured:    proto.Conjured_ConjuredDarkRune,
}

var FullDebuffs = &proto.Debuffs{
	JudgementOfWisdom: true,
	Misery:            true,
	CurseOfElements:   proto.TristateEffect_TristateEffectImproved,
}

var PlayerOptionsAdaptive = &proto.Player_BalanceDruid{
	BalanceDruid: &proto.BalanceDruid{
		Talents: StandardTalents,
		Options: &proto.BalanceDruid_Options{
			InnervateTarget: &proto.RaidTarget{TargetIndex: 0}, // self innervate
		},
		Rotation: &proto.BalanceDruid_Rotation{
			PrimarySpell: proto.BalanceDruid_Rotation_Adaptive,
			FaerieFire:   true,
		},
	},
}

var PlayerOptionsStarfire = &proto.Player_BalanceDruid{
	BalanceDruid: &proto.BalanceDruid{
		Talents: StandardTalents,
		Options: &proto.BalanceDruid_Options{
			InnervateTarget: &proto.RaidTarget{TargetIndex: 0}, // self innervate
		},
		Rotation: &proto.BalanceDruid_Rotation{
			PrimarySpell: proto.BalanceDruid_Rotation_Starfire,
			Moonfire:     true,
			FaerieFire:   true,
		},
	},
}

var PlayerOptionsWrath = &proto.Player_BalanceDruid{
	BalanceDruid: &proto.BalanceDruid{
		Talents: StandardTalents,
		Options: &proto.BalanceDruid_Options{
			InnervateTarget: &proto.RaidTarget{TargetIndex: 0}, // self innervate
		},
		Rotation: &proto.BalanceDruid_Rotation{
			PrimarySpell: proto.BalanceDruid_Rotation_Wrath,
			Moonfire:     true,
		},
	},
}

var PlayerOptionsAOE = &proto.Player_BalanceDruid{
	BalanceDruid: &proto.BalanceDruid{
		Talents: StandardTalents,
		Options: &proto.BalanceDruid_Options{
			InnervateTarget: &proto.RaidTarget{TargetIndex: 0}, // self innervate
		},
		Rotation: &proto.BalanceDruid_Rotation{
			PrimarySpell: proto.BalanceDruid_Rotation_Starfire,
			Moonfire:     true,
			FaerieFire:   true,
			Hurricane:    true,
		},
	},
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 29093,
		"enchant": 29191,
		"gems": [
			24030,
			34220
		]
	},
	{
		"id": 28762
	},
	{
		"id": 29095,
		"enchant": 28886,
		"gems": [
			24056,
			24059
		]
	},
	{
		"id": 28766,
		"enchant": 33150
	},
	{
		"id": 21848,
		"enchant": 24003,
		"gems": [
			24059,
			24056
		]
	},
	{
		"id": 24250,
		"enchant": 22534,
		"gems": [
			31867
		]
	},
	{
		"id": 21847,
		"enchant": 28272,
		"gems": [
			31867,
			31867
		]
	},
	{
		"id": 21846,
		"gems": [
			31867,
			31867
		]
	},
	{
		"id": 24262,
		"enchant": 24274,
		"gems": [
			31867,
			31867,
			31867
		]
	},
	{
		"id": 28517,
		"enchant": 35297,
		"gems": [
			31867,
			24059
		]
	},
	{
		"id": 28753,
		"enchant": 22536
	},
	{
		"id": 28793,
		"enchant": 22536
	},
	{
		"id": 29370
	},
	{
		"id": 27683
	},
	{
		"id": 28770,
		"enchant": 22560
	},
	{
		"id": 29271
	},
	{
		"id": 27518
	}
]}`)
var P2Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 30233,
		"enchant": 29191,
		"gems": [
			24059,
			34220
		]
	},
	{
		"id": 30015
	},
	{
		"id": 30235,
		"enchant": 28886,
		"gems": [
			24056,
			24059
		]
	},
	{
		"id": 28797,
		"enchant": 33150
	},
	{
		"id": 30231,
		"enchant": 24003,
		"gems": [
			24030,
			24030,
			24030
		]
	},
	{
		"id": 29918,
		"enchant": 22534
	},
	{
		"id": 30232,
		"enchant": 28272
	},
	{
		"id": 30038,
		"gems": [
			24056,
			24059
		]
	},
	{
		"id": 24262,
		"enchant": 24274,
		"gems": [
			24030,
			24030,
			24030
		]
	},
	{
		"id": 30067,
		"enchant": 35297
	},
	{
		"id": 28753,
		"enchant": 22536
	},
	{
		"id": 29302,
		"enchant": 22536
	},
	{
		"id": 29370
	},
	{
		"id": 27683
	},
	{
		"id": 29988,
		"enchant": 22560
	},
	{
		"id": 32387
	}
]}`)
