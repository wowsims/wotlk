package shadow

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var StandardTalents = &proto.PriestTalents{
	InnerFocus:             true,
	Meditation:             3,
	ShadowAffinity:         3,
	ImprovedShadowWordPain: 2,
	ShadowFocus:            3,
	ImprovedMindBlast:      5,
	MindFlay:               true,
	ShadowWeaving:          5,
	VampiricEmbrace:        true,
	FocusedMind:            3,
	Darkness:               5,
	Shadowform:             true,
	ShadowPower:            4,
	Misery:                 5,
	VampiricTouch:          true,
}

var FullRaidBuffs = &proto.RaidBuffs{
	GiftOfTheWild:         2,
	PowerWordFortitude:    2,
	StrengthOfEarthTotem:  2,
	ArcaneBrilliance:      true,
	DivineSpirit:          true,
	TrueshotAura:          true,
	LeaderOfThePack:       2,
	IcyTalons:             true,
	TotemOfWrath:          true,
	MoonkinAura:           2,
	WrathOfAirTotem:       true,
	SanctifiedRetribution: true,
	Bloodlust:             true,
}
var FullPartyBuffs = &proto.PartyBuffs{}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:  true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
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

var FullDebuffs = &proto.Debuffs{
	SunderArmor:        true,
	FaerieFire:         2,
	BloodFrenzy:        true,
	EbonPlaguebringer:  true,
	HeartOfTheCrusader: true,
	JudgementOfWisdom:  true,
}

var PlayerOptionsBasic = &proto.Player_ShadowPriest{
	ShadowPriest: &proto.ShadowPriest{
		Talents: StandardTalents,
		Options: &proto.ShadowPriest_Options{
			Armor:          proto.ShadowPriest_Options_InnerFire,
			UseShadowfiend: true,
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
			Armor:          proto.ShadowPriest_Options_InnerFire,
			UseShadowfiend: true,
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
			Armor:          proto.ShadowPriest_Options_InnerFire,
			UseShadowfiend: true,
		},
		Rotation: &proto.ShadowPriest_Rotation{
			RotationType: proto.ShadowPriest_Rotation_Ideal,
			PrecastVt:    true,
			Latency:      50,
		},
	},
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40562,
		"enchant": 44877,
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
		"enchant": 44874,
		"gems": [
			39998
		]
	},
	{
		"id": 44005,
		"enchant": 55642,
		"gems": [
			40026
		]
	},
	{
		"id": 44002,
		"enchant": 33990,
		"gems": [
			39998,
			39998
		]
	},
	{
		"id": 44008,
		"enchant": 44498,
		"gems": [
			39998,
			0
		]
	},
	{
		"id": 40454,
		"enchant": 54999,
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
		"enchant": 41602
	},
	{
		"id": 40558,
		"enchant": 55016
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
		"enchant": 44487
	},
	{
		"id": 40273
	},
	{
		"id": 39712
	}
]}`)
