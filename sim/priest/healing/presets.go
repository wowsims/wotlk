package healing

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var DiscTalents = "0503203130300512301313231251-2351010303"
var DiscGlyphs = &proto.Glyphs{
	Major1: int32(proto.PriestMajorGlyph_GlyphOfPowerWordShield),
	Major2: int32(proto.PriestMajorGlyph_GlyphOfFlashHeal),
	Major3: int32(proto.PriestMajorGlyph_GlyphOfPenance),
	// No interesting minor glyphs.
}

var HolyTalents = "05032031103-234051032002152530004311051"
var HolyGlyphs = &proto.Glyphs{
	Major1: int32(proto.PriestMajorGlyph_GlyphOfPrayerOfHealing),
	Major2: int32(proto.PriestMajorGlyph_GlyphOfRenew),
	Major3: int32(proto.PriestMajorGlyph_GlyphOfCircleOfHealing),
	// No interesting minor glyphs.
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	Food:          proto.Food_FoodFishFeast,
	DefaultPotion: proto.Potions_RunicManaInjector,
	PrepopPotion:  proto.Potions_PotionOfWildMagic,
}

var PlayerOptionsDisc = &proto.Player_HealingPriest{
	HealingPriest: &proto.HealingPriest{
		Options: &proto.HealingPriest_Options{
			UseInnerFire:      true,
			UseShadowfiend:    true,
			RapturesPerMinute: 5,
		},
		Rotation: &proto.HealingPriest_Rotation{},
	},
}

var PlayerOptionsHoly = &proto.Player_HealingPriest{
	HealingPriest: &proto.HealingPriest{
		Options: &proto.HealingPriest_Options{
			UseInnerFire:   true,
			UseShadowfiend: true,
		},
		Rotation: &proto.HealingPriest_Rotation{},
	},
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40456,
		"enchant": 3819,
		"gems": [
			41401,
			39998
		]
	},
	{
		"id": 44657,
		"gems": [
			40047
		]
	},
	{
		"id": 40450,
		"enchant": 3809,
		"gems": [
			42144
		]
	},
	{
		"id": 40724,
		"enchant": 3859
	},
	{
		"id": 40194,
		"enchant": 3832,
		"gems": [
			42144
		]
	},
	{
		"id": 40741,
		"enchant": 2332,
		"gems": [
			0
		]
	},
	{
		"id": 40445,
		"enchant": 3246,
		"gems": [
			42144,
			0
		]
	},
	{
		"id": 40271,
		"enchant": 3601,
		"gems": [
			40027,
			39998
		]
	},
	{
		"id": 40398,
		"enchant": 3719,
		"gems": [
			39998,
			39998
		]
	},
	{
		"id": 40236,
		"enchant": 3606
	},
	{
		"id": 40108
	},
	{
		"id": 40433
	},
	{
		"id": 37835
	},
	{
		"id": 40258
	},
	{
		"id": 40395,
		"enchant": 3834
	},
	{
		"id": 40350
	},
	{
		"id": 40245
	}
]}`)
