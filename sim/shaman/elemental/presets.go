package elemental

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var StandardTalents = "0532001523212351322301351-005052031"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.ShamanMajorGlyph_GlyphOfLava),
	Major2: int32(proto.ShamanMajorGlyph_GlyphOfTotemOfWrath),
	Major3: int32(proto.ShamanMajorGlyph_GlyphOfLightningBolt),
}

var NoTotems = &proto.ShamanTotems{}
var BasicTotems = &proto.ShamanTotems{
	Earth: proto.EarthTotem_TremorTotem,
	Air:   proto.AirTotem_WrathOfAirTotem,
	Water: proto.WaterTotem_ManaSpringTotem,
	Fire:  proto.FireTotem_TotemOfWrath,
}

var FireElementalBasicTotems = &proto.ShamanTotems{
	Earth:            proto.EarthTotem_TremorTotem,
	Air:              proto.AirTotem_WrathOfAirTotem,
	Water:            proto.WaterTotem_ManaSpringTotem,
	Fire:             proto.FireTotem_TotemOfWrath,
	UseFireElemental: true,
}

var eleShamOptions = &proto.ElementalShaman_Options{
	Shield:    proto.ShamanShield_WaterShield,
	Bloodlust: true,
}
var PlayerOptionsAdaptive = &proto.Player_ElementalShaman{
	ElementalShaman: &proto.ElementalShaman{
		Options: eleShamOptions,
		Rotation: &proto.ElementalShaman_Rotation{
			Totems: BasicTotems,
			Type:   proto.ElementalShaman_Rotation_Adaptive,
		},
	},
}

var PlayerOptionsAdaptiveFireElemental = &proto.Player_ElementalShaman{
	ElementalShaman: &proto.ElementalShaman{
		Options: eleShamOptions,
		Rotation: &proto.ElementalShaman_Rotation{
			Totems: FireElementalBasicTotems,
			Type:   proto.ElementalShaman_Rotation_Adaptive,
		},
	},
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfBlindingLight,
	Food:            proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:   proto.Potions_SuperManaPotion,
	PrepopPotion:    proto.Potions_DestructionPotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40516,
		"enchant": 3820,
		"gems": [
			41285,
			40027
		]
	},
	{
		"id": 44661,
		"gems": [
			39998
		]
	},
	{
		"id": 40286,
		"enchant": 3810
	},
	{
		"id": 44005,
		"enchant": 3722,
		"gems": [
			40027
		]
	},
	{
		"id": 40514,
		"enchant": 3832,
		"gems": [
			42144,
			42144
		]
	},
	{
		"id": 40324,
		"enchant": 2332,
		"gems": [
			42144,
			0
		]
	},
	{
		"id": 40302,
		"enchant": 3246,
		"gems": [
			0
		]
	},
	{
		"id": 40301,
		"gems": [
			40014
		]
	},
	{
		"id": 40560,
		"enchant": 3721
	},
	{
		"id": 40519,
		"enchant": 3826
	},
	{
		"id": 37694
	},
	{
		"id": 40399
	},
	{
		"id": 40432
	},
	{
		"id": 40255
	},
	{
		"id": 40395,
		"enchant": 3834
	},
	{
		"id": 40401,
		"enchant": 1128
	},
	{
		"id": 40267
	}
]}`)
