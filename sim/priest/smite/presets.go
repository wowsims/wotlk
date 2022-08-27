package smite

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var StandardTalents = &proto.PriestTalents{
	TwinDisciplines:            5,
	SilentResolve:              3,
	ImprovedInnerFire:          3,
	ImprovedPowerWordFortitude: 2,
	Meditation:                 3,
	InnerFocus:                 true,
	MentalAgility:              3,
	MentalStrength:             5,
	FocusedPower:               2,
	Enlightenment:              3,
	FocusedWill:                3,
	PowerInfusion:              true,

	HolySpecialization: 5,
	SpellWarding:       5,
	DivineFury:         5,
	DesperatePrayer:    true,
	HolyReach:          2,
	SearingLight:       2,
	SpiritOfRedemption: true,
	SpiritualGuidance:  5,
	SurgeOfLight:       2,

	SpiritTap:         3,
	ImprovedSpiritTap: 2,
	Darkness:          4,
}

var DefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.PriestMajorGlyph_GlyphOfSmite),
	Major2: int32(proto.PriestMajorGlyph_GlyphOfHolyNova),
	Major3: int32(proto.PriestMajorGlyph_GlyphOfShadowWordDeath),
	// No interesting minor glyphs.
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance: true,
	GiftOfTheWild:    proto.TristateEffect_TristateEffectImproved,
	MoonkinAura:      proto.TristateEffect_TristateEffectRegular,
	TotemOfWrath:     true,
	WrathOfAirTotem:  true,
	ManaSpringTotem:  proto.TristateEffect_TristateEffectRegular,
}
var FullPartyBuffs = &proto.PartyBuffs{}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:  true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	Food:          proto.Food_FoodFishFeast,
	DefaultPotion: proto.Potions_RunicManaPotion,
	PrepopPotion:  proto.Potions_PotionOfWildMagic,
}

var FullDebuffs = &proto.Debuffs{
	JudgementOfWisdom: true,
	CurseOfElements:   true,
}

var PlayerOptionsBasic = &proto.Player_SmitePriest{
	SmitePriest: &proto.SmitePriest{
		Talents: StandardTalents,
		Options: &proto.SmitePriest_Options{
			UseInnerFire:   true,
			UseShadowfiend: true,
		},
		Rotation: &proto.SmitePriest_Rotation{
			UseDevouringPlague: true,
			UseShadowWordDeath: true,
			UseMindBlast:       true,

			AllowedHolyFireDelayMs: 50,
		},
	},
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40562,
		"enchant": 44877,
		"gems": [
			41307,
			40049
		]
	},
	{
		"id": 40374
	},
	{
		"id": 40555,
		"enchant": 44874
	},
	{
		"id": 41610,
		"enchant": 63765
	},
	{
		"id": 40526,
		"enchant": 33990,
		"gems": [
			40049
		]
	},
	{
		"id": 40325,
		"enchant": 44498,
		"gems": [
			0
		]
	},
	{
		"id": 40454,
		"enchant": 44592,
		"gems": [
			40049,
			0
		]
	},
	{
		"id": 40301,
		"gems": [
			40049
		]
	},
	{
		"id": 40560,
		"enchant": 41602
	},
	{
		"id": 40246,
		"enchant": 60623
	},
	{
		"id": 40399
	},
	{
		"id": 39389
	},
	{
		"id": 42129
	},
	{
		"id": 40382
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
