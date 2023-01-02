package holy

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var defaultProtTalents = &proto.PaladinTalents{
	SpiritualFocus:           5,
	HealingLight:             3,
	DivineIntellect:          5,
	AuraMastery:              true,
	Illumination:             5,
	ImprovedLayOnHands:       1,
	ImprovedBlessingOfWisdom: 2,
	DivineFavor:              true,
	SanctifiedLight:          3,
	HolyPower:                5,
	LightsGrace:              3,
	HolyShock:                true,
	HolyGuidance:             5,
	DivineIllumination:       true,
	JudgementsOfThePure:      5,
	InfusionOfLight:          2,
	EnlightenedJudgements:    2,
	BeaconOfLight:            true,

	Divinity:              5,
	GuardiansFavor:        2,
	Anticipation:          3,
	DivineSacrifice:       true,
	ImprovedRighteousFury: 3,
	Toughness:             1,
	DivineGuardian:        2,
	ImprovedDevotionAura:  3,
}

var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.PaladinMajorGlyph_GlyphOfHolyLight),
	Major2: int32(proto.PaladinMajorGlyph_GlyphOfSealOfWisdom),
	Major3: int32(proto.PaladinMajorGlyph_GlyphOfBeaconOfLight),
	Minor1: int32(proto.PaladinMinorGlyph_GlyphOfLayOnHands),
	Minor2: int32(proto.PaladinMinorGlyph_GlyphOfSenseUndead),
}

var defaultProtRotation = &proto.HolyPaladin_Rotation{}

var defaultProtOptions = &proto.HolyPaladin_Options{
	Judgement: proto.PaladinJudgement_JudgementOfWisdom,
	Aura:      proto.PaladinAura_DevotionAura,
}

var BasicOptions = &proto.Player_HolyPaladin{
	HolyPaladin: &proto.HolyPaladin{
		Talents:  defaultProtTalents,
		Options:  defaultProtOptions,
		Rotation: defaultProtRotation,
	},
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfStoneblood,
	Food:            proto.Food_FoodDragonfinFilet,
	DefaultPotion:   proto.Potions_IndestructiblePotion,
	PrepopPotion:    proto.Potions_IndestructiblePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40298,
		"enchant": 3819,
		"gems": [
			41401,
			40012
		]
	},
	{
		"id": 44662,
		"gems": [
			40012
		]
	},
	{
		"id": 40573,
		"enchant": 3809,
		"gems": [
			40012
		]
	},
	{
		"id": 44005,
		"enchant": 3831,
		"gems": [
			40012
		]
	},
	{
		"id": 40569,
		"enchant": 3832,
		"gems": [
			40012,
			40012
		]
	},
	{
		"id": 40332,
		"enchant": 1119,
		"gems": [
			40012,
			0
		]
	},
	{
		"id": 40570,
		"enchant": 3604,
		"gems": [
			40012,
			0
		]
	},
	{
		"id": 40259,
		"gems": [
			40012
		]
	},
	{
		"id": 40572,
		"enchant": 3721,
		"gems": [
			40027,
			40012
		]
	},
	{
		"id": 40592,
		"enchant": 3606
	},
	{
		"id": 40399
	},
	{
		"id": 40375
	},
	{
		"id": 44255
	},
	{
		"id": 37111
	},
	{
		"id": 40395,
		"enchant": 2666
	},
	{
		"id": 40401,
		"enchant": 1128
	},
	{
		"id": 40705
	}
]}`)
