package rogue

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var CombatTalents = "00532000523-0252051050035010223100501251"
var CombatNoLethalityTalents = "00532000023-0252051050035010223100501251"
var CombatNoPotWTalents = "00532000523-0252051050035010223100501201"
var CombatNoLethalityNoPotWTalents = "00532000023-0252051050035010223100501201"
var AssassinationTalents = "005303005352100520103331051-005005003-502"
var CombatGlyphs = &proto.Glyphs{
	Major1: int32(proto.RogueMajorGlyph_GlyphOfKillingSpree),
	Major2: int32(proto.RogueMajorGlyph_GlyphOfTricksOfTheTrade),
	Major3: int32(proto.RogueMajorGlyph_GlyphOfRupture),
}
var AssassinationGlyphs = &proto.Glyphs{
	Major1: int32(proto.RogueMajorGlyph_GlyphOfMutilate),
	Major2: int32(proto.RogueMajorGlyph_GlyphOfTricksOfTheTrade),
	Major3: int32(proto.RogueMajorGlyph_GlyphOfHungerForBlood),
}

var PlayerOptionsCombatDI = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyInstant,
		Rotation: basicRotation,
	},
}
var PlayerOptionsCombatDD = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyDeadly,
		Rotation: basicRotation,
	},
}
var PlayerOptionsCombatID = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  InstantDeadly,
		Rotation: basicRotation,
	},
}
var PlayerOptionsCombatII = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  InstantInstant,
		Rotation: basicRotation,
	},
}

var PlayerOptionsNoLethality = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyInstant,
		Rotation: basicRotation,
	},
}

var PlayerOptionsNoPotW = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyInstant,
		Rotation: basicRotation,
	},
}

var PlayerOptionsNoLethalityNoPotW = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyInstant,
		Rotation: basicRotation,
	},
}

var PlayerOptionsAssassinationDI = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyInstant,
		Rotation: basicRotation,
	},
}
var PlayerOptionsAssassinationDD = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyDeadly,
		Rotation: basicRotation,
	},
}
var PlayerOptionsAssassinationID = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  InstantDeadly,
		Rotation: basicRotation,
	},
}
var PlayerOptionsAssassinationII = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
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
