package balance

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var StandardTalents = &proto.DruidTalents{
	StarlightWrath:        5,
	ImprovedMoonfire:      2,
	Brambles:              3,
	InsectSwarm:           true,
	Vengeance:             5,
	LunarGuidance:         3,
	Moonglow:              3,
	Moonfury:              5,
	BalanceOfPower:        2,
	Dreamstate:            3,
	MoonkinForm:           true,
	ImprovedFaerieFire:    3,
	WrathOfCenarius:       5,
	ForceOfNature:         true,
	Genesis:               5,
	NaturesMajesty:        2,
	NaturesGrace:          3,
	NaturesSplendor:       true,
	CelestialFocus:        3,
	ImprovedInsectSwarm:   3,
	ImprovedMoonkinForm:   3,
	OwlkinFrenzy:          0,
	Eclipse:               3,
	Typhoon:               true,
	GaleWinds:             2,
	EarthAndMoon:          3,
	Starfall:              true,
	ImprovedMarkOfTheWild: 5,
	Furor:                 2,
	NaturalShapeshifter:   3,
	Intensity:             3,
}

var DefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.DruidMajorGlyph_GlyphOfStarfire),
	Major2: int32(proto.DruidMajorGlyph_GlyphOfMoonfire),
	Major3: int32(proto.DruidMajorGlyph_GlyphOfStarfall),
	Minor1: int32(proto.DruidMajorGlyph_GlyphOfUnburdenedRebirth),
	Minor2: int32(proto.DruidMajorGlyph_GlyphOfTyphoon),
	Minor3: int32(proto.DruidMajorGlyph_GlyphOfTheWild),
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
			Type: proto.BalanceDruid_Rotation_Adaptive,
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
		"id": 29093,
		"enchant": 29191,
		"gems": [
			24030,
			34220
		]
	},
	{
		"id": 28762
	},
	{
		"id": 29095,
		"enchant": 28886,
		"gems": [
			24056,
			24059
		]
	},
	{
		"id": 28766,
		"enchant": 33150
	},
	{
		"id": 21848,
		"enchant": 24003,
		"gems": [
			24059,
			24056
		]
	},
	{
		"id": 24250,
		"enchant": 22534,
		"gems": [
			31867
		]
	},
	{
		"id": 21847,
		"enchant": 28272,
		"gems": [
			31867,
			31867
		]
	},
	{
		"id": 21846,
		"gems": [
			31867,
			31867
		]
	},
	{
		"id": 24262,
		"enchant": 24274,
		"gems": [
			31867,
			31867,
			31867
		]
	},
	{
		"id": 28517,
		"enchant": 35297,
		"gems": [
			31867,
			24059
		]
	},
	{
		"id": 28753,
		"enchant": 22536
	},
	{
		"id": 28793,
		"enchant": 22536
	},
	{
		"id": 29370
	},
	{
		"id": 27683
	},
	{
		"id": 28770,
		"enchant": 22560
	},
	{
		"id": 29271
	},
	{
		"id": 27518
	}
]}`)
