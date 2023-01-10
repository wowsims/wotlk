package shadow

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var DefaultTalents = "05032031--325023051223010323151301351"
var DefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.PriestMajorGlyph_GlyphOfShadow),
	Major2: int32(proto.PriestMajorGlyph_GlyphOfMindFlay),
	Major3: int32(proto.PriestMajorGlyph_GlyphOfDispersion),
	// No dps increasing minor glyphs.
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfPureDeath,
	Food:            proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:   proto.Potions_SuperManaPotion,
	PrepopPotion:    proto.Potions_PotionOfWildMagic,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var PlayerOptionsBasic = &proto.Player_ShadowPriest{
	ShadowPriest: &proto.ShadowPriest{
		Options: &proto.ShadowPriest_Options{
			Armor:              proto.ShadowPriest_Options_InnerFire,
			UseShadowfiend:     true,
			UseMindBlast:       true,
			UseShadowWordDeath: true,
		},
		Rotation: &proto.ShadowPriest_Rotation{
			RotationType: proto.ShadowPriest_Rotation_Basic,
			Latency:      50,
		},
	},
}
var PlayerOptionsClipping = &proto.Player_ShadowPriest{
	ShadowPriest: &proto.ShadowPriest{
		Options: &proto.ShadowPriest_Options{
			Armor:              proto.ShadowPriest_Options_InnerFire,
			UseShadowfiend:     true,
			UseMindBlast:       true,
			UseShadowWordDeath: true,
		},
		Rotation: &proto.ShadowPriest_Rotation{
			RotationType: proto.ShadowPriest_Rotation_Clipping,
			PrecastVt:    true,
			Latency:      50,
		},
	},
}
var PlayerOptionsIdeal = &proto.Player_ShadowPriest{
	ShadowPriest: &proto.ShadowPriest{
		Options: &proto.ShadowPriest_Options{
			Armor:              proto.ShadowPriest_Options_InnerFire,
			UseShadowfiend:     true,
			UseMindBlast:       true,
			UseShadowWordDeath: true,
		},
		Rotation: &proto.ShadowPriest_Rotation{
			RotationType: proto.ShadowPriest_Rotation_Ideal,
			PrecastVt:    true,
			Latency:      50,
		},
	},
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40562,
		"enchant": 3820,
		"gems": [
			41285,
			39998
		]
	},
	{
		"id": 44661,
		"gems": [
			40026
		]
	},
	{
		"id": 40459,
		"enchant": 3810,
		"gems": [
			39998
		]
	},
	{
		"id": 44005,
		"enchant": 3722,
		"gems": [
			40026
		]
	},
	{
		"id": 44002,
		"enchant": 1144,
		"gems": [
			39998,
			39998
		]
	},
	{
		"id": 44008,
		"enchant": 2332,
		"gems": [
			39998,
			0
		]
	},
	{
		"id": 40454,
		"enchant": 3604,
		"gems": [
			40049,
			0
		]
	},
	{
		"id": 40561,
		"gems": [
			39998
		]
	},
	{
		"id": 40560,
		"enchant": 3719
	},
	{
		"id": 40558,
		"enchant": 3606
	},
	{
		"id": 40719
	},
	{
		"id": 40399
	},
	{
		"id": 40255
	},
	{
		"id": 40432
	},
	{
		"id": 40395,
		"enchant": 3834
	},
	{
		"id": 40273
	},
	{
		"id": 39712
	}
]}`)
var P2Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
	  "id": 46172,
	  "enchant": 3820,
	  "gems": [
		41285,
		45883
	  ]
	},
	{
	  "id": 45243,
	  "gems": [
		39998
	  ]
	},
	{
	  "id": 46165,
	  "enchant": 3810,
	  "gems": [
		39998
	  ]
	},
	{
	  "id": 45242,
	  "enchant": 3722,
	  "gems": [
		40049
	  ]
	},
	{
	  "id": 46168,
	  "enchant": 1144,
	  "gems": [
		39998,
		39998
	  ]
	},
	{
	  "id": 45446,
	  "enchant": 2332,
	  "gems": [
		39998,
		0
	  ]
	},
	{
	  "id": 45665,
	  "enchant": 3604,
	  "gems": [
		39998,
		39998,
		0
	  ]
	},
	{
	  "id": 45619,
	  "enchant": 3601,
	  "gems": [
		39998,
		39998,
		39998
	  ]
	},
	{
	  "id": 46170,
	  "enchant": 3719,
	  "gems": [
		39998,
		40049
	  ]
	},
	{
	  "id": 45135,
	  "enchant": 3606,
	  "gems": [
		39998,
		40049
	  ]
	},
	{
	  "id": 45495,
	  "gems": [
		40026
	  ]
	},
	{
	  "id": 46046,
	  "gems": [
		39998
	  ]
	},
	{
	  "id": 45518
	},
	{
	  "id": 45466
	},
	{
	  "id": 45620,
	  "enchant": 3834,
	  "gems": [
		40026
	  ]
	},
	{
	  "id": 45617
	},
	{
	  "id": 45294,
	  "gems": [
		39998
	  ]
	}
  ]
}`)
