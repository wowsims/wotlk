package balance

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var StandardTalents = &proto.DruidTalents{
	StarlightWrath:        5,
	Moonglow:              1,
	NaturesMajesty:        2,
	ImprovedMoonfire:      2,
	NaturesGrace:          3,
	NaturesSplendor:       true,
	NaturesReach:          2,
	Vengeance:             5,
	CelestialFocus:        3,
	LunarGuidance:         3,
	InsectSwarm:           true,
	ImprovedInsectSwarm:   2,
	Moonfury:              3,
	BalanceOfPower:        2,
	MoonkinForm:           true,
	ImprovedMoonkinForm:   3,
	ImprovedFaerieFire:    3,
	WrathOfCenarius:       5,
	Eclipse:               3,
	Typhoon:               true,
	ForceOfNature:         true,
	GaleWinds:             2,
	EarthAndMoon:          3,
	Starfall:              true,
	ImprovedMarkOfTheWild: 2,
	Furor:                 5,
	NaturalShapeshifter:   3,
	MasterShapeshifter:    2,
	OmenOfClarity:         true,
}

var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.DruidMajorGlyph_GlyphOfFocus),
	Major2: int32(proto.DruidMajorGlyph_GlyphOfInsectSwarm),
	Major3: int32(proto.DruidMajorGlyph_GlyphOfStarfall),
	Minor1: int32(proto.DruidMinorGlyph_GlyphOfTyphoon),
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance: true,
	GiftOfTheWild:    proto.TristateEffect_TristateEffectImproved,
	MoonkinAura:      proto.TristateEffect_TristateEffectImproved,
}
var FullPartyBuffs = &proto.PartyBuffs{}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:  true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfBlindingLight,
	Food:            proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:   proto.Potions_SuperManaPotion,
	PrepopPotion:    proto.Potions_DestructionPotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var FullDebuffs = &proto.Debuffs{
	JudgementOfWisdom: true,
	Misery:            true,
	CurseOfElements:   true,
}

var PlayerOptionsAdaptive = &proto.Player_BalanceDruid{
	BalanceDruid: &proto.BalanceDruid{
		Talents: StandardTalents,
		Options: &proto.BalanceDruid_Options{
			InnervateTarget: &proto.RaidTarget{TargetIndex: 0}, // self innervate
		},
		Rotation: &proto.BalanceDruid_Rotation{
			Type:                     proto.BalanceDruid_Rotation_Adaptive,
			UseMf:                    true,
			UseIs:                    true,
			UseBattleRes:             true,
			IsInsideEclipseThreshold: 14.0,
			UseSmartCooldowns:        true,
			McdInsideLunarThreshold:  15.0,
			McdInsideSolarThreshold:  15.0,
		},
	},
}

var PlayerOptionsAOE = &proto.Player_BalanceDruid{
	BalanceDruid: &proto.BalanceDruid{
		Talents: StandardTalents,
		Options: &proto.BalanceDruid_Options{
			InnervateTarget: &proto.RaidTarget{TargetIndex: 0}, // self innervate
		},
		Rotation: &proto.BalanceDruid_Rotation{
			Type: proto.BalanceDruid_Rotation_Adaptive,
		},
	},
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40467,
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
		"id": 40470,
		"enchant": 44874,
		"gems": [
			39998
		]
	},
	{
		"id": 40405,
		"enchant": 44472
	},
	{
		"id": 40469,
		"enchant": 44489,
		"gems": [
			39998,
			40026
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
		"id": 40466,
		"enchant": 54999,
		"gems": [
			39998,
			0
		]
	},
	{
		"id": 40561,
		"enchant": 54793,
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
		"id": 40399
	},
	{
		"id": 40080
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
		"id": 40192
	},
	{
		"id": 40321
	}
]}`)
