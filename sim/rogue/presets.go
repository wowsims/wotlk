package rogue

import (
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
)

var CombatTalents = &proto.RogueTalents{
	Malice:              5,
	Ruthlessness:        3,
	Murder:              2,
	RelentlessStrikes:   true,
	ImprovedExposeArmor: 2,
	Lethality:           5,
	VilePoisons:         2,

	ImprovedSinisterStrike:  2,
	ImprovedSliceAndDice:    3,
	Precision:               5,
	DualWieldSpecialization: 5,
	BladeFlurry:             true,
	SwordSpecialization:     5,
	WeaponExpertise:         2,
	Aggression:              3,
	Vitality:                2,
	AdrenalineRush:          true,
	CombatPotency:           5,
	SurpriseAttacks:         true,
}

var MutilateTalents = &proto.RogueTalents{
	Malice:              5,
	Ruthlessness:        3,
	Murder:              2,
	PuncturingWounds:    3,
	RelentlessStrikes:   true,
	ImprovedExposeArmor: 2,
	Lethality:           5,
	ImprovedPoisons:     5,
	ColdBlood:           true,
	QuickRecovery:       2,
	SealFate:            5,
	Vigor:               true,
	FindWeakness:        5,
	Mutilate:            true,

	ImprovedSinisterStrike:  2,
	ImprovedSliceAndDice:    3,
	Precision:               5,
	DualWieldSpecialization: 5,
}

var HemoTalents = &proto.RogueTalents{
	ImprovedSinisterStrike:  2,
	ImprovedSliceAndDice:    3,
	Precision:               5,
	DualWieldSpecialization: 5,
	BladeFlurry:             true,
	SwordSpecialization:     5,
	WeaponExpertise:         2,
	Aggression:              3,
	Vitality:                2,
	AdrenalineRush:          true,
	CombatPotency:           5,

	SerratedBlades: 3,
	Hemorrhage:     true,
}

var PlayerOptionsBasic = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  CombatTalents,
		Options:  basicOptions,
		Rotation: basicRotation,
	},
}

var PlayerOptionsMutilate = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  MutilateTalents,
		Options:  basicOptions,
		Rotation: basicRotation,
	},
}

var PlayerOptionsHemo = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  HemoTalents,
		Options:  basicOptions,
		Rotation: basicRotation,
	},
}

var basicRotation = &proto.Rogue_Rotation{
	Builder:             proto.Rogue_Rotation_Auto,
	MaintainExposeArmor: true,
	UseRupture:          true,
	UseShiv:             true,

	MinComboPointsForDamageFinisher: 3,
}

var basicOptions = &proto.Rogue_Options{}

var FullRaidBuffs = &proto.RaidBuffs{
	GiftOfTheWild: proto.TristateEffect_TristateEffectImproved,
}
var FullPartyBuffs = &proto.PartyBuffs{
	Bloodlust: 1,
	Drums:     proto.Drums_DrumsOfBattle,

	BattleShout:       proto.TristateEffect_TristateEffectImproved,
	LeaderOfThePack:   proto.TristateEffect_TristateEffectImproved,
	GraceOfAirTotem:   proto.TristateEffect_TristateEffectRegular,
	WindfuryTotemRank: 5,
}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings: true,
	BlessingOfMight: proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	MainHandImbue:   proto.WeaponImbue_WeaponImbueRogueInstantPoison,
	OffHandImbue:    proto.WeaponImbue_WeaponImbueRogueDeadlyPoison,
	Flask:           proto.Flask_FlaskOfRelentlessAssault,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredRogueThistleTea,
	SuperSapper:     true,
	FillerExplosive: proto.Explosive_ExplosiveGnomishFlameTurret,
}

var FullDebuffs = &proto.Debuffs{
	BloodFrenzy:               true,
	Mangle:                    true,
	SunderArmor:               true,
	FaerieFire:                proto.TristateEffect_TristateEffectImproved,
	ImprovedSealOfTheCrusader: true,
	Misery:                    true,
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 29044,
		"enchant": 29192,
		"gems": [
			32409,
			24061
		]
	},
	{
		"id": 29381
	},
	{
		"id": 27797,
		"enchant": 28888,
		"gems": [
			24061,
			24055
		]
	},
	{
		"id": 28672,
		"enchant": 34004
	},
	{
		"id": 29045,
		"enchant": 24003,
		"gems": [
			24061,
			24051,
			24055
		]
	},
	{
		"id": 29246,
		"enchant": 34002
	},
	{
		"id": 27531,
		"enchant": 19445,
		"gems": [
			24061,
			24061
		]
	},
	{
		"id": 29247
	},
	{
		"id": 28741,
		"enchant": 29535,
		"gems": [
			24051,
			24051,
			24051
		]
	},
	{
		"id": 28545,
		"enchant": 28279,
		"gems": [
			24061,
			24051
		]
	},
	{
		"id": 28757
	},
	{
		"id": 28649
	},
	{
		"id": 29383
	},
	{
		"id": 28830
	},
	{
		"id": 28729,
		"enchant": 22559
	},
	{
		"id": 28189,
		"enchant": 22559
	},
	{
		"id": 28772
	}
]}`)
var MutilateP1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 29044,
		"enchant": 29192,
		"gems": [
			32409,
			24061
		]
	},
	{
		"id": 29381
	},
	{
		"id": 27797,
		"enchant": 28888,
		"gems": [
			24061,
			24055
		]
	},
	{
		"id": 28672,
		"enchant": 34004
	},
	{
		"id": 29045,
		"enchant": 24003,
		"gems": [
			24061,
			24051,
			24055
		]
	},
	{
		"id": 29246,
		"enchant": 34002
	},
	{
		"id": 27531,
		"enchant": 19445,
		"gems": [
			24061,
			24061
		]
	},
	{
		"id": 29247
	},
	{
		"id": 28741,
		"enchant": 29535,
		"gems": [
			24051,
			24051,
			24051
		]
	},
	{
		"id": 28545,
		"enchant": 28279,
		"gems": [
			24061,
			24051
		]
	},
	{
		"id": 28757
	},
	{
		"id": 28649
	},
	{
		"id": 29383
	},
	{
		"id": 28830
	},
	{
		"id": 28768,
		"enchant": 22559
	},
	{
		"id": 29182,
		"enchant": 22559
	},
	{
		"id": 28772
	}
]}`)
