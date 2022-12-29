package rogue

import (
	"github.com/wowsims/wotlk/sim/core"
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

var SubtletyTalents = &proto.RogueTalents{
	Malice:                  5,
	BloodSpatter:            2,
	PuncturingWounds:        3,
	Lethality:               5,
	VilePoisons:             2,
	ImprovedPoisons:         3,
	RelentlessStrikes:       5,
	Opportunity:             2,
	Camouflage:              3,
	Elusiveness:             2,
	GhostlyStrike:           true,
	SerratedBlades:          3,
	Initiative:              3,
	ImprovedAmbush:          2,
	Preparation:             true,
	DirtyDeeds:              2,
	Hemorrhage:              true,
	MasterOfSubtlety:        3,
	Deadliness:              5,
	Premeditation:           true,
	CheatDeath:              1,
	SinisterCalling:         5,
	HonorAmongThieves:       3,
	Shadowstep:              true,
	FilthyTricks:            1,
	SlaughterFromTheShadows: 5,
	ShadowDance:             true,
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

var SubtletyRotationOptions = []*proto.Rogue_Rotation{
	{
		ExposeArmorFrequency:      proto.Rogue_Rotation_Never,
		TricksOfTheTradeFrequency: proto.Rogue_Rotation_Never,
		SubtletyFinisherPriority:  proto.Rogue_Rotation_Rupture,
	},
	{
		ExposeArmorFrequency:      proto.Rogue_Rotation_Never,
		TricksOfTheTradeFrequency: proto.Rogue_Rotation_Once,
		SubtletyFinisherPriority:  proto.Rogue_Rotation_Rupture,
	},
	{
		ExposeArmorFrequency:          proto.Rogue_Rotation_Once,
		MinimumComboPointsExposeArmor: 3,
		TricksOfTheTradeFrequency:     proto.Rogue_Rotation_Never,
		SubtletyFinisherPriority:      proto.Rogue_Rotation_Rupture,
	},
	{
		ExposeArmorFrequency:      proto.Rogue_Rotation_Maintain,
		TricksOfTheTradeFrequency: proto.Rogue_Rotation_Never,
		SubtletyFinisherPriority:  proto.Rogue_Rotation_Rupture,
	},
	{
		ExposeArmorFrequency:      proto.Rogue_Rotation_Never,
		TricksOfTheTradeFrequency: proto.Rogue_Rotation_Maintain,
		SubtletyFinisherPriority:  proto.Rogue_Rotation_Rupture,
	},
	{
		ExposeArmorFrequency:      proto.Rogue_Rotation_Never,
		TricksOfTheTradeFrequency: proto.Rogue_Rotation_Maintain,
		SubtletyFinisherPriority:  proto.Rogue_Rotation_Eviscerate,
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

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfEndlessRage,
	DefaultPotion:   proto.Potions_PotionOfSpeed,
	DefaultConjured: proto.Conjured_ConjuredRogueThistleTea,
}

var PreRaidGear = core.EquipmentSpecFromJsonString(`{"items": [
    {
      "id": 42550,
      "enchant": 3817,
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
      "enchant": 3808
    },
    {
      "id": 38614,
      "enchant": 3605
    },
    {
      "id": 39558,
      "enchant": 3832,
      "gems": [
        40003,
        42702
      ]
    },
    {
      "id": 34448,
      "enchant": 3845,
      "gems": [
        40003,
        0
      ]
    },
    {
      "id": 39560,
      "enchant": 3604,
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
      "enchant": 3823
    },
    {
      "id": 34575,
      "enchant": 3606,
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
      "enchant": 3789
    },
    {
      "id": 37667,
      "enchant": 3789
    },
    {
      "id": 43612
    }
  ]}`)

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
    {
      "id": 40499,
      "enchant": 3817,
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
      "enchant": 3808,
      "gems": [
        36766
      ]
    },
    {
      "id": 40403,
      "enchant": 3605
    },
    {
      "id": 40539,
      "enchant": 3832,
      "gems": [
        36766
      ]
    },
    {
      "id": 39765,
      "enchant": 3845,
      "gems": [
        40003,
        0
      ]
    },
    {
      "id": 40496,
      "enchant": 3604,
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
      "enchant": 3823,
      "gems": [
        40003,
        40003
      ]
    },
    {
      "id": 39701,
      "enchant": 3606
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
      "enchant": 3789
    },
    {
      "id": 40386,
      "enchant": 3789
    },
    {
      "id": 40385
    }
  ]}`)
var GearWithoutRED = core.EquipmentSpecFromJsonString(`{"items": [
	{
	  "id": 37293,
	  "enchant": 3817,
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
	  "enchant": 3808,
	  "gems": [
		36766
	  ]
	},
	{
	  "id": 36947,
	  "enchant": 3605
	},
	{
	  "id": 37165,
	  "enchant": 3832,
	  "gems": [
		40044,
		36766
	  ]
	},
	{
	  "id": 44203,
	  "enchant": 3845,
	  "gems": [
		0
	  ]
	},
	{
	  "id": 37409,
	  "enchant": 3604,
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
	  "enchant": 3823
	},
	{
	  "id": 44297,
	  "enchant": 3606
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
	  "enchant": 3789
	},
	{
	  "id": 37856,
	  "enchant": 3789
	},
	{
	  "id": 37191
	}
  ]}`)
var GearWithRED = core.EquipmentSpecFromJsonString(`{"items": [
	{
	  "id": 37293,
	  "enchant": 3817,
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
	  "enchant": 3808,
	  "gems": [
		36766
	  ]
	},
	{
	  "id": 36947,
	  "enchant": 3605
	},
	{
	  "id": 37165,
	  "enchant": 3832,
	  "gems": [
		40044,
		36766
	  ]
	},
	{
	  "id": 44203,
	  "enchant": 3845,
	  "gems": [
		0
	  ]
	},
	{
	  "id": 37409,
	  "enchant": 3604,
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
	  "enchant": 3823
	},
	{
	  "id": 44297,
	  "enchant": 3606
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
	  "enchant": 3789
	},
	{
	  "id": 37856,
	  "enchant": 3789
	},
	{
	  "id": 37191
	}
  ]}`)
var MutilateP1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 29044,
		"enchant": 3003,
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
		"enchant": 2986,
		"gems": [
			24061,
			24055
		]
	},
	{
		"id": 28672,
		"enchant": 368
	},
	{
		"id": 29045,
		"enchant": 2661,
		"gems": [
			24061,
			24051,
			24055
		]
	},
	{
		"id": 29246,
		"enchant": 1593
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
		"enchant": 3012,
		"gems": [
			24051,
			24051,
			24051
		]
	},
	{
		"id": 28545,
		"enchant": 2939,
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
		"enchant": 2673
	},
	{
		"id": 29182,
		"enchant": 2673
	},
	{
		"id": 28772
	}
]}`)
