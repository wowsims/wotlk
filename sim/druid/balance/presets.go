package balance

import (
	"github.com/wowsims/wotlk/sim/core"
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

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfBlindingLight,
	Food:            proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:   proto.Potions_SuperManaPotion,
	PrepopPotion:    proto.Potions_DestructionPotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
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

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40467,
		"enchant": 3820,
		"gems": [
			41285,
			42144
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
		"enchant": 3810,
		"gems": [
			42144
		]
	},
	{
		"id": 44005,
		"enchant": 3859,
		"gems": [
			40026
		]
	},
	{
		"id": 40469,
		"enchant": 3832,
		"gems": [
			42144,
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
		"id": 40466,
		"enchant": 3604,
		"gems": [
			39998,
			0
		]
	},
	{
		"id": 40561,
		"enchant": 3601,
		"gems": [
			39998
		]
	},
	{
		"id": 40560,
		"enchant": 3719
	},
	{
		"id": 40519,
		"enchant": 3606
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
		"enchant": 3834
	},
	{
		"id": 40192
	},
	{
		"id": 40321
	}
]}`)
