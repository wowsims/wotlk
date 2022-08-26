package rogue

import (
	"github.com/wowsims/wotlk/sim/core"
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

var AssassinationTalents = &proto.RogueTalents{
	Malice:           5,
	Ruthlessness:     3,
	PuncturingWounds: 3,
	Lethality:        5,
	VilePoisons:      3,
	ImprovedPoisons:  5,
	FleetFooted:      2,
	ColdBlood:        true,
	SealFate:         5,
	Murder:           2,
	Overkill:         true,
	FocusedAttacks:   3,
	FindWeakness:     3,
	MasterPoisoner:   3,
	Mutilate:         true,
	CutToTheChase:    5,
	HungerForBlood:   true,

	DualWieldSpecialization: 5,
	Precision:               5,
	CloseQuartersCombat:     3,

	RelentlessStrikes: 5,
	Opportunity:       2,
}

var AssassinationRotationOptions = []*proto.Rogue_Rotation{
	{
		ExposeArmorFrequency:          proto.Rogue_Rotation_Never,
		TricksOfTheTradeFrequency:     proto.Rogue_Rotation_Never,
		AssassinationFinisherPriority: proto.Rogue_Rotation_EnvenomRupture,
	},
	{
		ExposeArmorFrequency:          proto.Rogue_Rotation_Never,
		TricksOfTheTradeFrequency:     proto.Rogue_Rotation_Once,
		AssassinationFinisherPriority: proto.Rogue_Rotation_EnvenomRupture,
	},
	{
		ExposeArmorFrequency:          proto.Rogue_Rotation_Once,
		MinimumComboPointsExposeArmor: 3,
		TricksOfTheTradeFrequency:     proto.Rogue_Rotation_Never,
		AssassinationFinisherPriority: proto.Rogue_Rotation_EnvenomRupture,
	},
	{
		ExposeArmorFrequency:          proto.Rogue_Rotation_Maintain,
		TricksOfTheTradeFrequency:     proto.Rogue_Rotation_Never,
		AssassinationFinisherPriority: proto.Rogue_Rotation_EnvenomRupture,
	},
	{
		ExposeArmorFrequency:          proto.Rogue_Rotation_Never,
		TricksOfTheTradeFrequency:     proto.Rogue_Rotation_Maintain,
		AssassinationFinisherPriority: proto.Rogue_Rotation_EnvenomRupture,
	},
	{
		ExposeArmorFrequency:          proto.Rogue_Rotation_Never,
		TricksOfTheTradeFrequency:     proto.Rogue_Rotation_Maintain,
		AssassinationFinisherPriority: proto.Rogue_Rotation_RuptureEnvenom,
	},
}

var CombatRotationOptions = []*proto.Rogue_Rotation{
	{
		ExposeArmorFrequency:      proto.Rogue_Rotation_Never,
		TricksOfTheTradeFrequency: proto.Rogue_Rotation_Never,
		CombatFinisherPriority:    proto.Rogue_Rotation_RuptureEviscerate,
	},
	{
		ExposeArmorFrequency:      proto.Rogue_Rotation_Never,
		TricksOfTheTradeFrequency: proto.Rogue_Rotation_Once,
		CombatFinisherPriority:    proto.Rogue_Rotation_RuptureEviscerate,
	},
	{
		ExposeArmorFrequency:          proto.Rogue_Rotation_Once,
		MinimumComboPointsExposeArmor: 3,
		TricksOfTheTradeFrequency:     proto.Rogue_Rotation_Never,
		CombatFinisherPriority:        proto.Rogue_Rotation_RuptureEviscerate,
	},
	{
		ExposeArmorFrequency:      proto.Rogue_Rotation_Maintain,
		TricksOfTheTradeFrequency: proto.Rogue_Rotation_Never,
		CombatFinisherPriority:    proto.Rogue_Rotation_RuptureEviscerate,
	},
	{
		ExposeArmorFrequency:      proto.Rogue_Rotation_Never,
		TricksOfTheTradeFrequency: proto.Rogue_Rotation_Maintain,
		CombatFinisherPriority:    proto.Rogue_Rotation_RuptureEviscerate,
	},
	{
		ExposeArmorFrequency:      proto.Rogue_Rotation_Never,
		TricksOfTheTradeFrequency: proto.Rogue_Rotation_Maintain,
		CombatFinisherPriority:    proto.Rogue_Rotation_EviscerateRupture,
	},
}

var RotationNames = []string{
	"No Expose No Tricks Primary",
	"No Expose One Tricks Primary",
	"One Expose No Tricks Primary",
	"Maintain Expose No Tricks Primary",
	"No Expose Maintain Tricks Primary",
	"No Expose Maintain Tricks Secondary",
}

func RotationSpecOptions(talents *proto.RogueTalents, options *proto.Rogue_Options) []core.SpecOptionsCombo {
	specs := make([]core.SpecOptionsCombo, 0)
	specLabel := "Assassination"
	rotationOptions := AssassinationRotationOptions
	if !talents.Mutilate {
		specLabel = "Combat"
		rotationOptions = CombatRotationOptions
	}
	for idx, rotation := range rotationOptions {
		specs = append(specs, core.SpecOptionsCombo{
			Label: specLabel + " " + RotationNames[idx],
			SpecOptions: &proto.Player_Rogue{
				Rogue: &proto.Rogue{
					Talents:  talents,
					Options:  options,
					Rotation: rotation,
				},
			},
		})
	}
	return specs
}

var PlayerOptionsCombatDI = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  CombatTalents,
		Options:  DeadlyInstant,
		Rotation: basicRotation,
	},
}
var PlayerOptionsCombatDD = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  CombatTalents,
		Options:  DeadlyDeadly,
		Rotation: basicRotation,
	},
}
var PlayerOptionsCombatID = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  CombatTalents,
		Options:  InstantDeadly,
		Rotation: basicRotation,
	},
}
var PlayerOptionsCombatII = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  CombatTalents,
		Options:  InstantInstant,
		Rotation: basicRotation,
	},
}

var PlayerOptionsNoLethality = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  CombatNoLethalityTalents,
		Options:  DeadlyInstant,
		Rotation: basicRotation,
	},
}

var PlayerOptionsNoPotW = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  CombatNoPotWTalents,
		Options:  DeadlyInstant,
		Rotation: basicRotation,
	},
}

var PlayerOptionsNoLethalityNoPotW = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  CombatNoLethalityNoPotWTalents,
		Options:  DeadlyInstant,
		Rotation: basicRotation,
	},
}

var PlayerOptionsAssassinationDI = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  AssassinationTalents,
		Options:  DeadlyInstant,
		Rotation: basicRotation,
	},
}
var PlayerOptionsAssassinationDD = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  AssassinationTalents,
		Options:  DeadlyDeadly,
		Rotation: basicRotation,
	},
}
var PlayerOptionsAssassinationID = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  AssassinationTalents,
		Options:  InstantDeadly,
		Rotation: basicRotation,
	},
}
var PlayerOptionsAssassinationII = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Talents:  AssassinationTalents,
		Options:  InstantInstant,
		Rotation: basicRotation,
	},
}

var basicRotation = &proto.Rogue_Rotation{
	ExposeArmorFrequency:                proto.Rogue_Rotation_Never,
	TricksOfTheTradeFrequency:           proto.Rogue_Rotation_Maintain,
	AssassinationFinisherPriority:       proto.Rogue_Rotation_EnvenomRupture,
	CombatFinisherPriority:              proto.Rogue_Rotation_RuptureEviscerate,
	MinimumComboPointsExposeArmor:       4,
	MinimumComboPointsPrimaryFinisher:   3,
	MinimumComboPointsSecondaryFinisher: 2,
	MultiTargetSliceFrequency:           proto.Rogue_Rotation_Once,
	MinimumComboPointsMultiTargetSlice:  4,
}

var DeadlyInstant = &proto.Rogue_Options{
	MhImbue: proto.Rogue_Options_DeadlyPoison,
	OhImbue: proto.Rogue_Options_InstantPoison,
}
var InstantDeadly = &proto.Rogue_Options{
	MhImbue: proto.Rogue_Options_DeadlyPoison,
	OhImbue: proto.Rogue_Options_InstantPoison,
}
var InstantInstant = &proto.Rogue_Options{
	MhImbue: proto.Rogue_Options_DeadlyPoison,
	OhImbue: proto.Rogue_Options_DeadlyPoison,
}
var DeadlyDeadly = &proto.Rogue_Options{
	MhImbue: proto.Rogue_Options_InstantPoison,
	OhImbue: proto.Rogue_Options_InstantPoison,
}

var FullRaidBuffs = &proto.RaidBuffs{
	AbominationsMight:     true,
	Bloodlust:             true,
	ElementalOath:         true,
	GiftOfTheWild:         proto.TristateEffect_TristateEffectImproved,
	IcyTalons:             true,
	LeaderOfThePack:       proto.TristateEffect_TristateEffectImproved,
	SanctifiedRetribution: true,
	StrengthOfEarthTotem:  proto.TristateEffect_TristateEffectImproved,
	SwiftRetribution:      true,
}

var FullPartyBuffs = &proto.PartyBuffs{}

var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings: true,
	BlessingOfMight: proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfEndlessRage,
	DefaultPotion:   proto.Potions_PotionOfSpeed,
	DefaultConjured: proto.Conjured_ConjuredRogueThistleTea,
}

var FullDebuffs = &proto.Debuffs{
	BloodFrenzy:        true,
	EarthAndMoon:       true,
	FaerieFire:         proto.TristateEffect_TristateEffectImproved,
	HeartOfTheCrusader: true,
	Mangle:             true,
	ShadowMastery:      true,
	SunderArmor:        true,
}

var PreRaidGear = items.EquipmentSpecFromJsonString(`{"items": [
    {
      "id": 42550,
      "enchant": 44879,
      "gems": [
        41398,
        40058
      ]
    },
    {
      "id": 40678
    },
    {
      "id": 43481,
      "enchant": 44871
    },
    {
      "id": 38614,
      "enchant": 55002
    },
    {
      "id": 39558,
      "enchant": 44489,
      "gems": [
        40003,
        42702
      ]
    },
    {
      "id": 34448,
      "enchant": 44484,
      "gems": [
        40003,
        0
      ]
    },
    {
      "id": 39560,
      "enchant": 54999,
      "gems": [
        40058,
        0
      ]
    },
    {
      "id": 40694,
      "gems": [
        40003,
        40003
      ]
    },
    {
      "id": 37644,
      "enchant": 38374
    },
    {
      "id": 34575,
      "enchant": 55016,
      "gems": [
        40003
      ]
    },
    {
      "id": 40586
    },
    {
      "id": 37642
    },
    {
      "id": 40684
    },
    {
      "id": 44253
    },
    {
      "id": 37856,
      "enchant": 44492
    },
    {
      "id": 37667,
      "enchant": 44492
    },
    {
      "id": 43612
    }
  ]}`)

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
    {
      "id": 40499,
      "enchant": 44879,
      "gems": [
        41398,
        42702
      ]
    },
    {
      "id": 44664,
      "gems": [
        42154
      ]
    },
    {
      "id": 40502,
      "enchant": 44871,
      "gems": [
        36766
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
        36766
      ]
    },
    {
      "id": 39765,
      "enchant": 44484,
      "gems": [
        40003,
        0
      ]
    },
    {
      "id": 40496,
      "enchant": 54999,
      "gems": [
        40058,
        0
      ]
    },
    {
      "id": 40260,
      "gems": [
        39999
      ]
    },
    {
      "id": 40500,
      "enchant": 38374,
      "gems": [
        40003,
        40003
      ]
    },
    {
      "id": 39701,
      "enchant": 55016
    },
    {
      "id": 40074
    },
    {
      "id": 40474
    },
    {
        "id": 40684
    },
    {
      "id": 44253
    },
    {
      "id": 39714,
      "enchant": 44492
    },
    {
      "id": 40386,
      "enchant": 44492
    },
    {
      "id": 40385
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
