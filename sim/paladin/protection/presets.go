package protection

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var StandardTalents = "-05005135200132311333312321-511302012003"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.PaladinMajorGlyph_GlyphOfSealOfVengeance),
	Major2: int32(proto.PaladinMajorGlyph_GlyphOfRighteousDefense),
	Major3: int32(proto.PaladinMajorGlyph_GlyphOfDivinePlea),
	Minor1: int32(proto.PaladinMinorGlyph_GlyphOfLayOnHands),
	Minor2: int32(proto.PaladinMinorGlyph_GlyphOfSenseUndead),
}

var defaultProtRotation = &proto.ProtectionPaladin_Rotation{}

var defaultProtOptions = &proto.ProtectionPaladin_Options{
	Judgement: proto.PaladinJudgement_JudgementOfWisdom,
	Seal:      proto.PaladinSeal_Vengeance,
	Aura:      proto.PaladinAura_RetributionAura,
}

var DefaultOptions = &proto.Player_ProtectionPaladin{
	ProtectionPaladin: &proto.ProtectionPaladin{
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
		"id": 40581,
		"enchant": 3818,
		"gems": [
			41396,
			36767
		]
	},
	{
		"id": 40387
	},
	{
		"id": 40584,
		"enchant": 3852,
		"gems": [
			49110
		]
	},
	{
		"id": 40410,
		"enchant": 3605
	},
	{
		"id": 40579,
		"enchant": 3832,
		"gems": [
			36767,
			40022
		]
	},
	{
		"id": 39764,
		"enchant": 3850,
		"gems": [
			0
		]
	},
	{
		"id": 40580,
		"enchant": 3860,
		"gems": [
			40008,
			0
		]
	},
	{
		"id": 39759,
		"enchant": 3601,
		"gems": [
			40008,
			40008
		]
	},
	{
		"id": 40589,
		"enchant": 3822
	},
	{
		"id": 39717,
		"enchant": 3606,
		"gems": [
			40089
		]
	},
	{
		"id": 40718
	},
	{
		"id": 40107
	},
	{
		"id": 44063,
		"gems": [
			36767,
			40089
		]
	},
	{
		"id": 37220
	},
	{
		"id": 40345,
		"enchant": 3788
	},
	{
		"id": 40400,
		"enchant": 3849
	},
	{
		"id": 40707
	}
]}`)
