package restoration

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var StandardTalents = &proto.ShamanTalents{
	EnhancingTotems:    3,
	AncestralKnowledge: 2,
	ThunderingStrikes:  5,
	ImprovedShields:    3,

	ImprovedHealingWave:   5,
	TidalFocus:            5,
	ImprovedWaterShield:   3,
	HealingFocus:          3,
	TidalForce:            true,
	AncestralHealing:      3,
	RestorativeTotems:     3,
	TidalMastery:          5,
	HealingWay:            3,
	NaturesSwiftness:      true,
	Purification:          5,
	ManaTideTotem:         true,
	CleanseSpirit:         true,
	BlessingOfTheEternals: 2,
	ImprovedChainHeal:     2,
	NaturesBlessing:       3,
	AncestralAwakening:    3,
	EarthShield:           true,
	ImprovedEarthShield:   2,
	TidalWaves:            5,
	Riptide:               true,
}

var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.ShamanMajorGlyph_GlyphOfChainHeal),
	Major2: int32(proto.ShamanMajorGlyph_GlyphOfEarthShield),
	Major3: int32(proto.ShamanMajorGlyph_GlyphOfEarthlivingWeapon),
}

var BasicTotems = &proto.ShamanTotems{
	Earth: proto.EarthTotem_TremorTotem,
	Air:   proto.AirTotem_WrathOfAirTotem,
	Water: proto.WaterTotem_ManaSpringTotem,
	Fire:  proto.FireTotem_FlametongueTotem,
}

var restoShamOptions = &proto.RestorationShaman_Options{
	Shield:    proto.ShamanShield_WaterShield,
	Bloodlust: true,
}
var PlayerOptionsStandard = &proto.Player_RestorationShaman{
	RestorationShaman: &proto.RestorationShaman{
		Talents: StandardTalents,
		Options: restoShamOptions,
		Rotation: &proto.RestorationShaman_Rotation{
			Totems: BasicTotems,
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
		"id": 40510,
		"enchant": 3820,
		"gems": [
			41401,
			40051
		]
	},
	{
		"id": 44662,
		"gems": [
			42150
		]
	},
	{
		"id": 40513,
		"enchant": 3810,
		"gems": [
			40051
		]
	},
	{
		"id": 44005,
		"enchant": 3859,
		"gems": [
			40105
		]
	},
	{
		"id": 40508,
		"enchant": 2381,
		"gems": [
			42144,
			42150
		]
	},
	{
		"id": 40209,
		"enchant": 2332,
		"gems": [
			0
		]
	},
	{
		"id": 40564,
		"enchant": 3604,
		"gems": [
			0
		]
	},
	{
		"id": 40327,
		"gems": [
			0
		]
	},
	{
		"id": 40512,
		"enchant": 3721,
		"gems": [
			40051,
			40105
		]
	},
	{
		"id": 40237,
		"enchant": 3606,
		"gems": [
			40105
		]
	},
	{
		"id": 40399
	},
	{
		"id": 40375
	},
	{
		"id": 40432
	},
	{
		"id": 37111
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
		"id": 40709
	}
]}`)
