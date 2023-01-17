package mage

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterMage()
}

func TestArcane(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,
		Race:  proto.Race_RaceTroll,

		GearSet:     core.GearSetCombo{Label: "P1Arcane", GearSet: P1ArcaneGear},
		Talents:     ArcaneTalents,
		Glyphs:      ArcaneGlyphs,
		Consumes:    FullArcaneConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "ArcaneRotation", SpecOptions: PlayerOptionsArcane},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "AOE", SpecOptions: PlayerOptionsArcaneAOE},
		},

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeCloth,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeWand,
			},
		},
	}))
}

func TestFire(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,
		Race:  proto.Race_RaceTroll,

		GearSet:     core.GearSetCombo{Label: "P1Fire", GearSet: P1FireGear},
		Talents:     FireTalents,
		Glyphs:      FireGlyphs,
		Consumes:    FullFireConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "FireRotation", SpecOptions: PlayerOptionsFire},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "AOE", SpecOptions: PlayerOptionsFireAOE},
		},

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeCloth,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeWand,
			},
		},
	}))
}

func TestFrost(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,
		Race:  proto.Race_RaceTroll,

		GearSet:     core.GearSetCombo{Label: "P1Frost", GearSet: P1FrostGear},
		Talents:     FrostTalents,
		Glyphs:      FrostGlyphs,
		Consumes:    FullFrostConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "FrostRotation", SpecOptions: PlayerOptionsFrost},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "AOE", SpecOptions: PlayerOptionsFrostAOE},
		},

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeCloth,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeWand,
			},
		},
	}))
}

var ArcaneTalents = "23000513310033015032310250532-03-023303001"
var FireTalents = "23000503110003-0055030012303331053120301351"
var FrostTalents = "23000503110003--0533030310233100030152231351"
var ArcaneGlyphs = &proto.Glyphs{
	Major1: int32(proto.MageMajorGlyph_GlyphOfArcaneBlast),
	Major2: int32(proto.MageMajorGlyph_GlyphOfArcaneMissiles),
	Major3: int32(proto.MageMajorGlyph_GlyphOfMoltenArmor),
}
var FireGlyphs = &proto.Glyphs{
	Major1: int32(proto.MageMajorGlyph_GlyphOfFireball),
	Major2: int32(proto.MageMajorGlyph_GlyphOfMoltenArmor),
	Major3: int32(proto.MageMajorGlyph_GlyphOfLivingBomb),
}
var FrostGlyphs = &proto.Glyphs{
	Major1: int32(proto.MageMajorGlyph_GlyphOfFrostbolt),
	Major3: int32(proto.MageMajorGlyph_GlyphOfEternalWater),
	Major2: int32(proto.MageMajorGlyph_GlyphOfMoltenArmor),
}

var fireMageOptions = &proto.Mage_Options{
	Armor:          proto.Mage_Options_MoltenArmor,
	ReactionTimeMs: 300,
	IgniteMunching: true,
}
var PlayerOptionsFire = &proto.Player_Mage{
	Mage: &proto.Mage{
		Options: fireMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type:                   proto.Mage_Rotation_Fire,
			PrimaryFireSpell:       proto.Mage_Rotation_Fireball,
			MaintainImprovedScorch: false,
			PyroblastDelayMs:       50,
		},
	},
}
var PlayerOptionsFireAOE = &proto.Player_Mage{
	Mage: &proto.Mage{
		Options: fireMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type: proto.Mage_Rotation_Aoe,
			Aoe:  proto.Mage_Rotation_Flamestrike,
		},
	},
}

var frostMageOptions = &proto.Mage_Options{
	Armor:          proto.Mage_Options_MageArmor,
	ReactionTimeMs: 300,
}
var PlayerOptionsFrost = &proto.Player_Mage{
	Mage: &proto.Mage{
		Options: frostMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type: proto.Mage_Rotation_Frost,
		},
	},
}
var PlayerOptionsFrostAOE = &proto.Player_Mage{
	Mage: &proto.Mage{
		Options: frostMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type: proto.Mage_Rotation_Aoe,
			Aoe:  proto.Mage_Rotation_Blizzard,
		},
	},
}

var arcaneMageOptions = &proto.Mage_Options{
	Armor:          proto.Mage_Options_MoltenArmor,
	ReactionTimeMs: 300,
}
var PlayerOptionsArcane = &proto.Player_Mage{
	Mage: &proto.Mage{
		Options: arcaneMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type:                                       proto.Mage_Rotation_Arcane,
			ExtraBlastsDuringFirstAp:                   2,
			MissileBarrageBelowArcaneBlastStacks:       0,
			MissileBarrageBelowManaPercent:             0.1,
			BlastWithoutMissileBarrageAboveManaPercent: 0.2,
			Only_3ArcaneBlastStacksBelowManaPercent:    0.15,
		},
	},
}
var PlayerOptionsArcaneAOE = &proto.Player_Mage{
	Mage: &proto.Mage{
		Options: arcaneMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type: proto.Mage_Rotation_Aoe,
			Aoe:  proto.Mage_Rotation_ArcaneExplosion,
		},
	},
}

var FullFireConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	Food:          proto.Food_FoodFirecrackerSalmon,
	DefaultPotion: proto.Potions_PotionOfSpeed,
	// DefaultConjured: proto.Conjured_ConjuredFlameCap,
}
var FullFrostConsumes = FullFireConsumes

var FullArcaneConsumes = FullFireConsumes

var P1ArcaneGear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40416,
		"enchant": 3820,
		"gems": [
			41285,
			39998
		]
	},
	{
		"id": 44661,
		"gems": [
			40026
		]
	},
	{
		"id": 40419,
		"enchant": 3810,
		"gems": [
			40051
		]
	},
	{
		"id": 44005,
		"enchant": 3722,
		"gems": [
			40026
		]
	},
	{
		"id": 44002,
		"enchant": 3832,
		"gems": [
			39998,
			39998
		]
	},
	{
		"id": 44008,
		"enchant": 2332,
		"gems": [
			39998,
			0
		]
	},
	{
		"id": 40415,
		"enchant": 3604,
		"gems": [
			39998,
			0
		]
	},
	{
		"id": 40561,
		"gems": [
			39998
		]
	},
	{
		"id": 40417,
		"enchant": 3719,
		"gems": [
			39998,
			40051
		]
	},
	{
		"id": 40558,
		"enchant": 3606
	},
	{
		"id": 40719
	},
	{
		"id": 40399
	},
	{
		"id": 39229
	},
	{
		"id": 40255
	},
	{
		"id": 40396,
		"enchant": 3834
	},
	{
		"id": 40273
	},
	{
		"id": 39426
	}
]}`)
var P1FrostGear = P1ArcaneGear
var P1FireGear = P1ArcaneGear
