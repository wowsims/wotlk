package dps

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var FuryTalents = "32002301233-305053000520310053120500351"
var FuryGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarriorMajorGlyph_GlyphOfWhirlwind),
	Major2: int32(proto.WarriorMajorGlyph_GlyphOfHeroicStrike),
	Major3: int32(proto.WarriorMajorGlyph_GlyphOfExecution),
}
var ArmsTalents = "3022032023335100102012213231251-305-2033"
var ArmsGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarriorMajorGlyph_GlyphOfRending),
	Major2: int32(proto.WarriorMajorGlyph_GlyphOfMortalStrike),
	Major3: int32(proto.WarriorMajorGlyph_GlyphOfExecution),
}

var PlayerOptionsArms = &proto.Player_Warrior{
	Warrior: &proto.Warrior{
		Options:  warriorOptions,
		Rotation: armsRotation,
	},
}

var PlayerOptionsFury = &proto.Player_Warrior{
	Warrior: &proto.Warrior{
		Options:  warriorOptions,
		Rotation: furyRotation,
	},
}

var armsRotation = &proto.Warrior_Rotation{
	UseRend:   true,
	UseMs:     true,
	UseCleave: false,

	HsRageThreshold:   50,
	MsRageThreshold:   35,
	SlamRageThreshold: 25,
	RendCdThreshold:   0,

	SpamExecute: false,

	UseHsDuringExecute: true,

	MaintainDemoShout:   false,
	MaintainThunderClap: false,

	StanceOption: proto.Warrior_Rotation_DefaultStance,
}

var furyRotation = &proto.Warrior_Rotation{
	UseRend:   false,
	UseCleave: false,

	HsRageThreshold:          30,
	RendRageThresholdBelow:   100,
	SlamRageThreshold:        25,
	RendCdThreshold:          0,
	RendHealthThresholdAbove: 20,

	UseHsDuringExecute: true,
	UseWwDuringExecute: true,
	UseBtDuringExecute: true,
	UseSlamOverExecute: true,

	MaintainDemoShout:   false,
	MaintainThunderClap: false,

	StanceOption: proto.Warrior_Rotation_DefaultStance,
}

var warriorOptions = &proto.Warrior_Options{
	StartingRage:    50,
	UseRecklessness: true,
	Shout:           proto.WarriorShout_WarriorShoutBattle,
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfEndlessRage,
	DefaultPotion: proto.Potions_IndestructiblePotion,
	PrepopPotion:  proto.Potions_PotionOfSpeed,
	Food:          proto.Food_FoodDragonfinFilet,
}

var FuryP1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 44006,
		"enchant": 3817,
		"gems": [
			41285,
			42702
		]
	},
	{
		"id": 44664,
		"gems": [
			39996
		]
	},
	{
		"id": 40530,
		"enchant": 3808,
		"gems": [
			40037
		]
	},
	{
		"id": 40403,
		"enchant": 3605
	},
	{
		"id": 40539,
		"enchant": 3832,
		"gems": [
			42142
		]
	},
	{
		"id": 39765,
		"enchant": 3845,
		"gems": [
			39996,
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
			42142
		]
	},
	{
		"id": 40529,
		"enchant": 3823,
		"gems": [
			39996,
			40022
		]
	},
	{
		"id": 40591,
		"enchant": 3606
	},
	{
		"id": 43993,
		"gems": [
			42142
		]
	},
	{
		"id": 40717
	},
	{
		"id": 42987
	},
	{
		"id": 40256
	},
	{
		"id": 40384,
		"enchant": 3789
	},
	{
		"id": 40384,
		"enchant": 3789
	},
	{
		"id": 40385
	}
]}`)
