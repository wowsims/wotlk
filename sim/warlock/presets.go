package warlock

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var defaultDestroTalents = &proto.WarlockTalents{
	// destro
	ImprovedShadowBolt: 5,
	Bane:               5,
	Ruin:               5,
	Intensity:          2,
	DestructiveReach:   2,
	Devastation:        true,
	Aftermath:          2,
	ImprovedImmolate:   3,
	Emberstorm:         5,
	Conflagrate:        true,
	Backlash:           3,
	Shadowburn:         true,
	ShadowAndFlame:     5,
	Backdraft:          3,
	EmpoweredImp:       3,
	FireAndBrimstone:   5,
	ChaosBolt:          true,
	Shadowfury:         true,
	Pyroclasm:          3,
	DemonicPower:       2,
	Cataclysm:          3,
	SoulLeech:          3,
	ImprovedSoulLeech:  2,
	// demo
	FelSynergy:  2,
	ImprovedImp: 3,
}

var defaultDestroRotation = &proto.Warlock_Rotation{
	PrimarySpell: proto.Warlock_Rotation_Incinerate,
	SecondaryDot: proto.Warlock_Rotation_Immolate,
}

var defaultDestroOptions = &proto.Warlock_Options{
	Armor:       proto.Warlock_Options_FelArmor,
	Summon:      proto.Warlock_Options_Imp,
	WeaponImbue: proto.Warlock_Options_GrandFirestone,
}

var DefaultDestroWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Talents:  defaultDestroTalents,
		Options:  defaultDestroOptions,
		Rotation: defaultDestroRotation,
		// MajorGlyph: defaultDestroMajorGlyphs,
		// MinorGlyph: defaultDestroMinorGlyphs,
	},
}

var DefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarlockMajorGlyph_GlyphOfConflagrate),
	Major2: int32(proto.WarlockMajorGlyph_GlyphOfImp),
	Major3: int32(proto.WarlockMajorGlyph_GlyphOfIncinerate),
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance: true,
	GiftOfTheWild:    proto.TristateEffect_TristateEffectImproved,
	DivineSpirit:     true,
	Bloodlust:        true,
	ManaSpringTotem:  proto.TristateEffect_TristateEffectRegular,
	WrathOfAirTotem:  true,
	TotemOfWrath:     true,
}

var FullPartyBuffs = &proto.PartyBuffs{}

var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings: true,
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	DefaultPotion: proto.Potions_PotionOfWildMagic,
	PrepopPotion:  proto.Potions_PotionOfWildMagic,
	Food:          proto.Food_FoodFishFeast,
}

var FullDebuffs = &proto.Debuffs{
	JudgementOfWisdom: true,
	Misery:            true,
	BloodFrenzy:       true,
	SunderArmor:       true,
	FaerieFire:        proto.TristateEffect_TristateEffectImproved,
	CurseOfWeakness:   proto.TristateEffect_TristateEffectImproved,
}

var Phase4Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 31051,
		"enchant": 29191,
		"gems": [
		34220,
		32218
		]
	},
	{
		"id": 33281
	},
	{
		"id": 31054,
		"enchant": 28886,
		"gems": [
		32215,
		32218
		]
	},
	{
		"id": 32524,
		"enchant": 33150
	},
	{
		"id": 30107,
		"enchant": 24003,
		"gems": [
		32196,
		32196,
		32196
		]
	},
	{
		"id": 32586,
		"enchant": 22534
	},
	{
		"id": 31050,
		"enchant": 28272,
		"gems": [
		32196
		]
	},
	{
		"id": 30888,
		"gems": [
		32196,
		32196
		]
	},
	{
		"id": 31053,
		"enchant": 24274,
		"gems": [
		32196
		]
	},
	{
		"id": 32239,
		"enchant": 35297,
		"gems": [
		32218,
		32215
		]
	},
	{
		"id": 32527,
		"enchant": 22536
	},
	{
		"id": 33497,
		"enchant": 22536
	},
	{
		"id": 32483
	},
	{
		"id": 33829
	},
	{
		"id": 32374,
		"enchant": 22561
	},
	{},
	{
		"id": 33192,
		"gems": [
		32215
		]
	}
]}`)
