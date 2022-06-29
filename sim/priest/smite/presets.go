package smite

import (
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
)

var StandardTalents = &proto.PriestTalents{
	InnerFocus:           true,
	Meditation:           3,
	SilentResolve:        1,
	MentalAgility:        5,
	MentalStrength:       5,
	DivineSpirit:         true,
	ImprovedDivineSpirit: 2,
	ForceOfWill:          5,
	PowerInfusion:        true,
	HolySpecialization:   5,
	DivineFury:           5,
	SearingLight:         2,
	SpiritualGuidance:    5,
	SurgeOfLight:         2,
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance: true,
	GiftOfTheWild:    proto.TristateEffect_TristateEffectImproved,
}
var FullPartyBuffs = &proto.PartyBuffs{
	MoonkinAura:     proto.TristateEffect_TristateEffectRegular,
	TotemOfWrath:    1,
	WrathOfAirTotem: proto.TristateEffect_TristateEffectImproved,
	ManaSpringTotem: proto.TristateEffect_TristateEffectRegular,
}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:  true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	Flask:              proto.Flask_FlaskOfBlindingLight,
	Food:               proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:      proto.Potions_SuperManaPotion,
	NumStartingPotions: 1,
	DefaultConjured:    proto.Conjured_ConjuredDarkRune,
	MainHandImbue:      proto.WeaponImbue_WeaponImbueBrilliantWizardOil,
}

var FullDebuffs = &proto.Debuffs{
	JudgementOfWisdom: true,
	CurseOfElements:   proto.TristateEffect_TristateEffectImproved,
}

var PlayerOptionsBasic = &proto.Player_SmitePriest{
	SmitePriest: &proto.SmitePriest{
		Talents: StandardTalents,
		Options: &proto.SmitePriest_Options{
			UseShadowfiend: true,
		},
		Rotation: &proto.SmitePriest_Rotation{
			RotationType: proto.SmitePriest_Rotation_Basic,
		},
	},
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 24266,
		"enchant": 29191,
		"gems": [
			28118,
			24030,
			24030
		]
	},
	{
		"id": 28530
	},
	{
		"id": 29060,
		"enchant": 28886,
		"gems": [
			24030,
			24030
		]
	},
	{
		"id": 28766,
		"enchant": 33150
	},
	{
		"id": 29056,
		"enchant": 24003,
		"gems": [
			24030,
			24030,
			24030
		]
	},
	{
		"id": 24250,
		"enchant": 22534,
		"gems": [
			24030
		]
	},
	{
		"id": 30725,
		"enchant": 28272,
		"gems": [
			24030,
			24030
		]
	},
	{
		"id": 24256,
		"gems": [
			24030,
			24030
		]
	},
	{
		"id": 30734,
		"enchant": 24274,
		"gems": [
			24030,
			24030,
			24030
		]
	},
	{
		"id": 28517,
		"enchant": 35297,
		"gems": [
			24030,
			24030
		]
	},
	{
		"id": 28793,
		"enchant": 22536
	},
	{
		"id": 29172,
		"enchant": 22536
	},
	{
		"id": 27683
	},
	{
		"id": 29370
	},
	{
		"id": 30723,
		"enchant": 22555,
		"gems": [
			30564,
			31867
		]
	},
	{
		"id": 28734
	},
	{
		"id": 28673
	}
]}`)
var P3Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 32525,
		"enchant": 29191,
		"gems": [
			34220,
			30600
		]
	},
	{
		"id": 32349
	},
	{
		"id": 31070,
		"enchant": 28886,
		"gems": [
			32218,
			32215
		]
	},
	{
		"id": 32524,
		"enchant": 33150
	},
	{
		"id": 30107,
		"enchant": 33990,
		"gems": [
			32196,
			32196,
			32196
		]
	},
	{
		"id": 32586,
		"enchant": 22534
	},
	{
		"id": 31061,
		"enchant": 28272,
		"gems": [
			32196
		]
	},
	{
		"id": 30038,
		"gems": [
			32196,
			32196
		]
	},
	{
		"id": 30916,
		"enchant": 24274,
		"gems": [
			32196,
			32196,
			32196
		]
	},
	{
		"id": 32239,
		"enchant": 35297,
		"gems": [
			32196,
			32196
		]
	},
	{
		"id": 32527,
		"enchant": 22536
	},
	{
		"id": 32247,
		"enchant": 22536
	},
	{
		"id": 29370
	},
	{
		"id": 32483
	},
	{
		"id": 32374,
		"enchant": 22555
	},
	{},
	{
		"id": 29982
	}
]}`)
