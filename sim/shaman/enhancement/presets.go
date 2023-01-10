package enhancement

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var StandardTalents = "053030152-30405003105021333031131031051"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.ShamanMajorGlyph_GlyphOfStormstrike),
	Major2: int32(proto.ShamanMajorGlyph_GlyphOfFlametongueWeapon),
	Major3: int32(proto.ShamanMajorGlyph_GlyphOfFeralSpirit),
}

var PlayerOptionsBasic = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options:  enhShamOptions,
		Rotation: enhShamRotation,
	},
}

var PlayerOptionsFireElemental = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options:  enhShamOptions,
		Rotation: enhShamRotationFireElemental,
	},
}

var PlayerOptionsItemSwap = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options:  enhShamOptions,
		Rotation: enhShamRotationItemSwap,
	},
}

var enhShamRotationFireElemental = &proto.EnhancementShaman_Rotation{
	Totems: &proto.ShamanTotems{
		Earth:            proto.EarthTotem_StrengthOfEarthTotem,
		Air:              proto.AirTotem_WindfuryTotem,
		Water:            proto.WaterTotem_ManaSpringTotem,
		Fire:             proto.FireTotem_MagmaTotem,
		UseFireElemental: true,
	},
	RotationType:                 proto.EnhancementShaman_Rotation_Priority,
	FirenovaManaThreshold:        3000,
	ShamanisticRageManaThreshold: 25,
	PrimaryShock:                 proto.EnhancementShaman_Rotation_Earth,
	WeaveFlameShock:              true,
}

var enhShamRotation = &proto.EnhancementShaman_Rotation{
	Totems: &proto.ShamanTotems{
		Earth: proto.EarthTotem_StrengthOfEarthTotem,
		Air:   proto.AirTotem_WindfuryTotem,
		Water: proto.WaterTotem_ManaSpringTotem,
		Fire:  proto.FireTotem_MagmaTotem,
	},
	RotationType:                 proto.EnhancementShaman_Rotation_Priority,
	FirenovaManaThreshold:        3000,
	ShamanisticRageManaThreshold: 25,
	PrimaryShock:                 proto.EnhancementShaman_Rotation_Earth,
	WeaveFlameShock:              true,
}

var enhShamRotationItemSwap = &proto.EnhancementShaman_Rotation{
	Totems: &proto.ShamanTotems{
		Earth:            proto.EarthTotem_StrengthOfEarthTotem,
		Air:              proto.AirTotem_WindfuryTotem,
		Water:            proto.WaterTotem_ManaSpringTotem,
		Fire:             proto.FireTotem_MagmaTotem,
		UseFireElemental: true,
	},
	RotationType:                 proto.EnhancementShaman_Rotation_Priority,
	FirenovaManaThreshold:        3000,
	ShamanisticRageManaThreshold: 25,
	PrimaryShock:                 proto.EnhancementShaman_Rotation_Earth,
	WeaveFlameShock:              true,
	//Temp to test Item Swap, will switch to a more realistic swap with Phase 2 gear.
	EnableItemSwap: true,
	ItemSwap: &proto.ItemSwap{
		MhItem: &proto.ItemSpec{
			Id: 41752,
		},
		OhItem: &proto.ItemSpec{
			Id:      41752,
			Enchant: 3790,
		},
	},
}

var enhShamOptions = &proto.EnhancementShaman_Options{
	Shield:    proto.ShamanShield_LightningShield,
	Bloodlust: true,
	SyncType:  proto.ShamanSyncType_SyncMainhandOffhandSwings,
	ImbueMh:   proto.ShamanImbue_FlametongueWeaponDownrank, //phase 1 (wraith strike) only
	ImbueOh:   proto.ShamanImbue_FlametongueWeapon,
}

var FullConsumes = &proto.Consumes{
	DefaultConjured: proto.Conjured_ConjuredFlameCap,
}

var Phase1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40543,
		"enchant": 3817,
		"gems": [
			41398,
			40014
		]
	},
	{
		"id": 44661,
		"gems": [
			40014
		]
	},
	{
		"id": 40524,
		"enchant": 3808,
		"gems": [
			40014
		]
	},
	{
		"id": 40403,
		"enchant": 3605
	},
	{
		"id": 40523,
		"enchant": 3832,
		"gems": [
			40003,
			40014
		]
	},
	{
		"id": 40282,
		"enchant": 3845,
		"gems": [
			42702,
			0
		]
	},
	{
		"id": 40520,
		"enchant": 3604,
		"gems": [
			42154,
			0
		]
	},
	{
		"id": 40275,
		"gems": [
			42156
		]
	},
	{
		"id": 40522,
		"enchant": 3823,
		"gems": [
			39999,
			42156
		]
	},
	{
		"id": 40367,
		"enchant": 3606,
		"gems": [
			40058
		]
	},
	{
		"id": 40474
	},
	{
		"id": 40074
	},
	{
		"id": 40684
	},
	{
		"id": 37390
	},
	{
		"id": 39763,
		"enchant": 3789
	},
	{
		"id": 39468,
		"enchant": 3789
	},
	{
		"id": 40322
	}
]}`)
