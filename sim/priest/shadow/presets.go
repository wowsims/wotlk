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
		"id": 24266,
		"enchant": 29191,
		"gems": [
			28118,
			24030,
			24030
		]
	},
	{
		"id": 30666
	},
	{
		"id": 21869,
		"enchant": 28886,
		"gems": [
			24030,
			24030
		]
	},
	{
		"id": 28570,
		"enchant": 33150
	},
	{
		"id": 21871,
		"enchant": 24003,
		"gems": [
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
		"id": 28507,
		"enchant": 28272,
		"gems": [
			24030,
			24030
		]
	},
	{
		"id": 28799,
		"gems": [
			24030,
			24030
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
		"id": 21870,
		"enchant": 35297,
		"gems": [
			24030,
			24030
		]
	},
	{
		"id": 29352,
		"enchant": 22536
	},
	{
		"id": 28793,
		"enchant": 22536
	},
	{
		"id": 28789
	},
	{
		"id": 29370
	},
	{
		"id": 28770,
		"enchant": 22561
	},
	{
		"id": 29272
	},
	{
		"id": 29350
	}
]}`)
var P3Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 31064,
		"enchant": 29191,
		"gems": [
			25893,
			32215
		]
	},
	{
		"id": 30666
	},
	{
		"id": 31070,
		"enchant": 28886,
		"gems": [
			32196,
			32196
		]
	},
	{
		"id": 32590,
		"enchant": 33150
	},
	{
		"id": 31065,
		"enchant": 24003,
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
		"id": 32256
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
		"id": 32527,
		"enchant": 22536
	},
	{
		"id": 32483
	},
	{
		"id": 29370
	},
	{
		"id": 32374,
		"enchant": 22561
	},
	{
		"id": 29982
	}
]}`)
