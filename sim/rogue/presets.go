package rogue

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var CombatTalents = &proto.RogueTalents{
	Malice:          5,
	Ruthlessness:    3,
	BloodSpatter:    2,
	Lethality:       5,
	VilePoisons:     2,
	ImprovedPoisons: 3,

	ImprovedSinisterStrike:  2,
	DualWieldSpecialization: 5,
	ImprovedSliceAndDice:    2,
	Precision:               5,
	Endurance:               1,
	CloseQuartersCombat:     5,
	LightningReflexes:       3,
	Aggression:              5,
	BladeFlurry:             true,
	WeaponExpertise:         2,
	BladeTwisting:           2,
	Vitality:                3,
	AdrenalineRush:          true,
	CombatPotency:           5,
	SurpriseAttacks:         true,
	SavageCombat:            2,
	PreyOnTheWeak:           5,
	KillingSpree:            true,
}

var CombatNoLethalityTalents = &proto.RogueTalents{
	Malice:          5,
	Ruthlessness:    3,
	BloodSpatter:    2,
	VilePoisons:     2,
	ImprovedPoisons: 3,

	ImprovedSinisterStrike:  2,
	DualWieldSpecialization: 5,
	ImprovedSliceAndDice:    2,
	Precision:               5,
	Endurance:               1,
	CloseQuartersCombat:     5,
	LightningReflexes:       3,
	Aggression:              5,
	BladeFlurry:             true,
	WeaponExpertise:         2,
	BladeTwisting:           2,
	Vitality:                3,
	AdrenalineRush:          true,
	CombatPotency:           5,
	SurpriseAttacks:         true,
	SavageCombat:            2,
	PreyOnTheWeak:           5,
	KillingSpree:            true,
}

var CombatNoPotWTalents = &proto.RogueTalents{
	Malice:          5,
	Ruthlessness:    3,
	BloodSpatter:    2,
	Lethality:       5,
	VilePoisons:     2,
	ImprovedPoisons: 3,

	ImprovedSinisterStrike:  2,
	DualWieldSpecialization: 5,
	ImprovedSliceAndDice:    2,
	Precision:               5,
	Endurance:               1,
	CloseQuartersCombat:     5,
	LightningReflexes:       3,
	Aggression:              5,
	BladeFlurry:             true,
	WeaponExpertise:         2,
	BladeTwisting:           2,
	Vitality:                3,
	AdrenalineRush:          true,
	CombatPotency:           5,
	SurpriseAttacks:         true,
	SavageCombat:            2,
	KillingSpree:            true,
}

var CombatNoLethalityNoPotWTalents = &proto.RogueTalents{
	Malice:          5,
	Ruthlessness:    3,
	BloodSpatter:    2,
	VilePoisons:     2,
	ImprovedPoisons: 3,

	ImprovedSinisterStrike:  2,
	DualWieldSpecialization: 5,
	ImprovedSliceAndDice:    2,
	Precision:               5,
	Endurance:               1,
	CloseQuartersCombat:     5,
	LightningReflexes:       3,
	Aggression:              5,
	BladeFlurry:             true,
	WeaponExpertise:         2,
	BladeTwisting:           2,
	Vitality:                3,
	AdrenalineRush:          true,
	CombatPotency:           5,
	SurpriseAttacks:         true,
	SavageCombat:            2,
	KillingSpree:            true,
}

var MutilateTalents = &proto.RogueTalents{
	Malice:              5,
	Ruthlessness:        3,
	Murder:              2,
	PuncturingWounds:    3,
	RelentlessStrikes:   5,
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
	HackAndSlash:            5,
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

var PlayerOptionsNoLethality = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  CombatNoLethalityTalents,
		Options:  basicOptions,
		Rotation: basicRotation,
	},
}

var PlayerOptionsNoPotW = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  CombatNoPotWTalents,
		Options:  basicOptions,
		Rotation: basicRotation,
	},
}

var PlayerOptionsNoLethalityNoPotW = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  CombatNoLethalityNoPotWTalents,
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
	Builder:                  proto.Rogue_Rotation_Auto,
	MaintainExposeArmor:      false,
	MaintainTricksOfTheTrade: true,
}

var basicOptions = &proto.Rogue_Options{
	MhImbue: proto.Rogue_Options_DeadlyPoison,
	OhImbue: proto.Rogue_Options_InstantPoison,
}

var FullRaidBuffs = &proto.RaidBuffs{
	GiftOfTheWild:   proto.TristateEffect_TristateEffectImproved,
	Bloodlust:       true,
	BattleShout:     proto.TristateEffect_TristateEffectImproved,
	LeaderOfThePack: proto.TristateEffect_TristateEffectImproved,
	WindfuryTotem:   proto.TristateEffect_TristateEffectImproved,
}
var FullPartyBuffs = &proto.PartyBuffs{}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings: true,
	BlessingOfMight: proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfRelentlessAssault,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredRogueThistleTea,
	ThermalSapper:   true,
	FillerExplosive: proto.Explosive_ExplosiveSaroniteBomb,
}

var FullDebuffs = &proto.Debuffs{
	BloodFrenzy: true,
	Mangle:      true,
	SunderArmor: true,
	FaerieFire:  proto.TristateEffect_TristateEffectImproved,
	Misery:      true,
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
var GearWithoutRED = items.EquipmentSpecFromJsonString(`{"items": [
	{
	  "id": 37293,
	  "enchant": 44879,
	  "gems": [
		41339,
		40088
	  ]
	},
	{
	  "id": 37861
	},
	{
	  "id": 37139,
	  "enchant": 44871,
	  "gems": [
		36766
	  ]
	},
	{
	  "id": 36947,
	  "enchant": 55002
	},
	{
	  "id": 37165,
	  "enchant": 44489,
	  "gems": [
		40044,
		36766
	  ]
	},
	{
	  "id": 44203,
	  "enchant": 44484,
	  "gems": [
		0
	  ]
	},
	{
	  "id": 37409,
	  "enchant": 54999,
	  "gems": [
		0
	  ]
	},
	{
	  "id": 37194,
	  "gems": [
		40014,
		40157
	  ]
	},
	{
	  "id": 37644,
	  "enchant": 38374
	},
	{
	  "id": 44297,
	  "enchant": 55016
	},
	{
	  "id": 43251,
	  "gems": [
		40136
	  ]
	},
	{
	  "id": 37642
	},
	{
	  "id": 37390
	},
	{
	  "id": 37166
	},
	{
	  "id": 37693,
	  "enchant": 44492
	},
	{
	  "id": 37856,
	  "enchant": 44492
	},
	{
	  "id": 37191
	}
  ]}`)
var GearWithRED = items.EquipmentSpecFromJsonString(`{"items": [
	{
	  "id": 37293,
	  "enchant": 44879,
	  "gems": [
		41398,
		40088
	  ]
	},
	{
	  "id": 37861
	},
	{
	  "id": 37139,
	  "enchant": 44871,
	  "gems": [
		36766
	  ]
	},
	{
	  "id": 36947,
	  "enchant": 55002
	},
	{
	  "id": 37165,
	  "enchant": 44489,
	  "gems": [
		40044,
		36766
	  ]
	},
	{
	  "id": 44203,
	  "enchant": 44484,
	  "gems": [
		0
	  ]
	},
	{
	  "id": 37409,
	  "enchant": 54999,
	  "gems": [
		0
	  ]
	},
	{
	  "id": 37194,
	  "gems": [
		40014,
		40157
	  ]
	},
	{
	  "id": 37644,
	  "enchant": 38374
	},
	{
	  "id": 44297,
	  "enchant": 55016
	},
	{
	  "id": 43251,
	  "gems": [
		40136
	  ]
	},
	{
	  "id": 37642
	},
	{
	  "id": 37390
	},
	{
	  "id": 37166
	},
	{
	  "id": 37693,
	  "enchant": 44492
	},
	{
	  "id": 37856,
	  "enchant": 44492
	},
	{
	  "id": 37191
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
