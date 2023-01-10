package protection

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var DefaultTalents = "2500030023-302-053351225000012521030113321"
var DefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarriorMajorGlyph_GlyphOfBlocking),
	Major2: int32(proto.WarriorMajorGlyph_GlyphOfDevastate),
	Major3: int32(proto.WarriorMajorGlyph_GlyphOfVigilance),
}

var PlayerOptionsBasic = &proto.Player_ProtectionWarrior{
	ProtectionWarrior: &proto.ProtectionWarrior{
		Options:  warriorOptions,
		Rotation: warriorRotation,
	},
}

var warriorRotation = &proto.ProtectionWarrior_Rotation{
	HsRageThreshold: 30,
}

var warriorOptions = &proto.ProtectionWarrior_Options{
	Shout:                proto.WarriorShout_WarriorShoutCommanding,
	PrecastShout:         false,
	PrecastShoutT2:       false,
	PrecastShoutSapphire: false,

	StartingRage: 0,
}

var FullConsumes = &proto.Consumes{
	BattleElixir:   proto.BattleElixir_ElixirOfMastery,
	GuardianElixir: proto.GuardianElixir_GiftOfArthas,
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40546,
		"enchant": 3818,
		"gems": [
			41380,
			40034
		]
	},
	{
		"id": 40387
	},
	{
		"id": 39704,
		"enchant": 3852,
		"gems": [
			40034
		]
	},
	{
		"id": 40722,
		"enchant": 3605
	},
	{
		"id": 44000,
		"enchant": 3832,
		"gems": [
			40034,
			40015
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
		"id": 40545,
		"enchant": 3860,
		"gems": [
			40034,
			0
		]
	},
	{
		"id": 39759,
		"enchant": 3601,
		"gems": [
			40008,
			36767
		]
	},
	{
		"id": 40589,
		"enchant": 3822
	},
	{
		"id": 39717,
		"enchant": 3232,
		"gems": [
			40089
		]
	},
	{
		"id": 40370
	},
	{
		"id": 40718
	},
	{
		"id": 40257
	},
	{
		"id": 44063,
		"gems": [
			36767,
			40089
		]
	},
	{
		"id": 40402,
		"enchant": 3788
	},
	{
		"id": 40400,
		"enchant": 3849
	},
	{
		"id": 41168,
		"gems": [
			36767
		]
	}
]}`)
