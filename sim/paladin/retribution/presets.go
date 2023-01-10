package retribution

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var StandardTalents = "050501-05-05232051203331302133231331"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.PaladinMajorGlyph_GlyphOfSealOfVengeance),
	Major2: int32(proto.PaladinMajorGlyph_GlyphOfJudgement),
	Major3: int32(proto.PaladinMajorGlyph_GlyphOfConsecration),
	Minor1: int32(proto.PaladinMinorGlyph_GlyphOfSenseUndead),
	Minor2: int32(proto.PaladinMinorGlyph_GlyphOfLayOnHands),
	Minor3: int32(proto.PaladinMinorGlyph_GlyphOfBlessingOfKings),
}

var defaultRetRotation = &proto.RetributionPaladin_Rotation{
	ConsSlack:            500,
	ExoSlack:             500,
	UseDivinePlea:        true,
	DivinePleaPercentage: 0.75,
	HolyWrathThreshold:   4,
}

var defaultRetOptions = &proto.RetributionPaladin_Options{
	Judgement:            proto.PaladinJudgement_JudgementOfWisdom,
	Seal:                 proto.PaladinSeal_Vengeance,
	Aura:                 proto.PaladinAura_RetributionAura,
	DamageTakenPerSecond: 0,
}

var DefaultOptions = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options:  defaultRetOptions,
		Rotation: defaultRetRotation,
	},
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfRelentlessAssault,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
	Food:            proto.Food_FoodRoastedClefthoof,
	ThermalSapper:   true,
}

var Phase1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40576,
		"enchant": 3817,
		"gems": [
			41398,
			40037
		]
	},
	{
		"id": 44664,
		"gems": [
			42142
		]
	},
	{
		"id": 40578,
		"enchant": 3808,
		"gems": [
			49110
		]
	},
	{
		"id": 40403,
		"enchant": 3605
	},
	{
		"id": 40574,
		"enchant": 3832,
		"gems": [
			42142,
			39996
		]
	},
	{
		"id": 40186,
		"enchant": 3845,
		"gems": [
			0
		]
	},
	{
		"id": 40541,
		"enchant": 3604,
		"gems": [
			0
		]
	},
	{
		"id": 40205,
		"gems": [
			39996
		]
	},
	{
		"id": 40577,
		"enchant": 3823,
		"gems": [
			42142,
			40038
		]
	},
	{
		"id": 39701,
		"enchant": 3606
	},
	{
		"id": 40075
	},
	{
		"id": 40474
	},
	{
		"id": 42987
	},
	{
		"id": 40431
	},
	{
		"id": 40384,
		"enchant": 3789
	},
	{},
	{
		"id": 42852
	}
]}`)
