package shadow

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var StandardTalents = &proto.PriestTalents{
	TwinDisciplines:            5,
	ImprovedInnerFire:          3,
	ImprovedPowerWordFortitude: 2,
	Meditation:                 3,
	InnerFocus:                 true,

	SpiritTap:               3,
	ImprovedSpiritTap:       2,
	Darkness:                5,
	ImprovedShadowWordPain:  2,
	ShadowFocus:             3,
	ImprovedMindBlast:       5,
	MindFlay:                true,
	VeiledShadows:           2,
	ShadowReach:             2,
	ShadowWeaving:           3,
	VampiricEmbrace:         true,
	FocusedMind:             3,
	MindMelt:                2,
	ImprovedDevouringPlague: 3,
	Shadowform:              true,
	ShadowPower:             5,
	ImprovedShadowform:      1,
	Misery:                  3,
	VampiricTouch:           true,
	PainAndSuffering:        3,
	TwistedFaith:            5,
	Dispersion:              true,
}

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
		Talents: StandardTalents,
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
		Talents: StandardTalents,
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
		Talents: StandardTalents,
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
