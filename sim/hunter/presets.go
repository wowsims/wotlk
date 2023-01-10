package hunter

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var BMTalents = "51200201515012233110531351-005305-5"
var MMTalents = "502-035335131030013233035031051-5000002"
var SVTalents = "-015305101-5000032500033330532135301311"
var BMGlyphs = &proto.Glyphs{
	Major1: int32(proto.HunterMajorGlyph_GlyphOfBestialWrath),
	Major2: int32(proto.HunterMajorGlyph_GlyphOfSteadyShot),
	Major3: int32(proto.HunterMajorGlyph_GlyphOfSerpentSting),
}
var MMGlyphs = &proto.Glyphs{
	Major1: int32(proto.HunterMajorGlyph_GlyphOfSerpentSting),
	Major2: int32(proto.HunterMajorGlyph_GlyphOfSteadyShot),
	Major3: int32(proto.HunterMajorGlyph_GlyphOfChimeraShot),
}
var SVGlyphs = &proto.Glyphs{
	Major1: int32(proto.HunterMajorGlyph_GlyphOfSerpentSting),
	Major2: int32(proto.HunterMajorGlyph_GlyphOfExplosiveShot),
	Major3: int32(proto.HunterMajorGlyph_GlyphOfKillShot),
}

var FerocityTalents = &proto.HunterPetTalents{
	CobraReflexes:  2,
	Dive:           true,
	SpikedCollar:   3,
	BoarsSpeed:     true,
	CullingTheHerd: 3,
	SpidersBite:    3,
	Rabid:          true,
	CallOfTheWild:  true,
	WildHunt:       1,
}

var PlayerOptionsMM = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Options:  basicOptions,
		Rotation: basicRotation,
	},
}

var PlayerOptionsBM = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Options:  basicOptions,
		Rotation: basicRotation,
	},
}

var PlayerOptionsSV = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Options:  basicOptions,
		Rotation: basicRotation,
	},
}

var PlayerOptionsAOE = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Options:  basicOptions,
		Rotation: aoeRotation,
	},
}

var basicRotation = &proto.Hunter_Rotation{
	Sting: proto.Hunter_Rotation_SerpentSting,

	ViperStartManaPercent: 0.2,
	ViperStopManaPercent:  0.3,
}

var aoeRotation = &proto.Hunter_Rotation{
	TrapWeave:         true,
	TimeToTrapWeaveMs: 2000,

	ViperStartManaPercent: 0.2,
	ViperStopManaPercent:  0.3,
}

var basicOptions = &proto.Hunter_Options{
	Ammo:       proto.Hunter_Options_SaroniteRazorheads,
	PetType:    proto.Hunter_Options_Wolf,
	PetTalents: FerocityTalents,
	PetUptime:  0.9,

	SniperTrainingUptime: 0.8,
	UseHuntersMark:       true,
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfRelentlessAssault,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredFlameCap,
	PetFood:         proto.PetFood_PetFoodKiblersBits,
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40505,
		"enchant": 3817,
		"gems": [
			41398,
			42143
		]
	},
	{
		"id": 44664,
		"gems": [
			42143
		]
	},
	{
		"id": 40507,
		"enchant": 3808,
		"gems": [
			39997
		]
	},
	{
		"id": 40403,
		"enchant": 3605
	},
	{
		"id": 43998,
		"enchant": 3832,
		"gems": [
			42143,
			39997
		]
	},
	{
		"id": 40282,
		"enchant": 3845,
		"gems": [
			39997,
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
		"id": 39762,
		"enchant": 3601,
		"gems": [
			39997
		]
	},
	{
		"id": 40331,
		"enchant": 3823,
		"gems": [
			39997,
			49110
		]
	},
	{
		"id": 40549,
		"enchant": 3606
	},
	{
		"id": 40074
	},
	{
		"id": 40474
	},
	{
		"id": 40684
	},
	{
		"id": 44253
	},
	{
		"id": 40388,
		"enchant": 3827
	},
	{},
	{
		"id": 40385,
		"enchant": 3608
	}
]}`)
