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

		GearSet:     core.GearSetCombo{Label: "P3Arcane", GearSet: P3ArcaneGear},
		Talents:     ArcaneTalents,
		Glyphs:      ArcaneGlyphs,
		Consumes:    FullArcaneConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Arcane", SpecOptions: PlayerOptionsArcane},
		Rotation:    core.GetAplRotation("../../ui/mage/apls", "arcane"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../ui/mage/apls", "arcane_aoe"),
		},

		ItemFilter: ItemFilter,
	}))
}

func TestFire(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,
		Race:  proto.Race_RaceTroll,

		GearSet:     core.GearSetCombo{Label: "P3Fire", GearSet: P3FireGear},
		Talents:     FireTalents,
		Glyphs:      FireGlyphs,
		Consumes:    FullFireConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Fire", SpecOptions: PlayerOptionsFire},
		Rotation:    core.GetAplRotation("../../ui/mage/apls", "fire"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../ui/mage/apls", "fire_aoe"),
		},

		ItemFilter: ItemFilter,
	}))
}

func TestFrostFire(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,
		Race:  proto.Race_RaceTroll,

		GearSet:     core.GearSetCombo{Label: "P3FrostFire", GearSet: P3FireGear},
		Talents:     FrostFireTalents,
		Glyphs:      FrostFireGlyphs,
		Consumes:    FullFireConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Fire", SpecOptions: PlayerOptionsFire},
		Rotation:    core.GetAplRotation("../../ui/mage/apls", "frostfire"),

		ItemFilter: ItemFilter,
	}))
}

func TestFrost(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,
		Race:  proto.Race_RaceTroll,

		GearSet:     core.GearSetCombo{Label: "P3Frost", GearSet: P3FrostGear},
		Talents:     FrostTalents,
		Glyphs:      FrostGlyphs,
		Consumes:    FullFrostConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Frost", SpecOptions: PlayerOptionsFrost},
		Rotation:    core.GetAplRotation("../../ui/mage/apls", "frost"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../ui/mage/apls", "frost_aoe"),
		},

		ItemFilter: ItemFilter,
	}))
}

var ItemFilter = core.ItemFilter{
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
}

var ArcaneTalents = "23000513310033015032310250532-03-023303001"
var FireTalents = "23000503110003-0055030012303331053120301351"
var FrostFireTalents = "23000503110003-0055030012303331053120301351"
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
var FrostFireGlyphs = &proto.Glyphs{
	Major1: int32(proto.MageMajorGlyph_GlyphOfFrostfire),
	Major2: int32(proto.MageMajorGlyph_GlyphOfMoltenArmor),
	Major3: int32(proto.MageMajorGlyph_GlyphOfLivingBomb),
}
var FrostGlyphs = &proto.Glyphs{
	Major1: int32(proto.MageMajorGlyph_GlyphOfFrostbolt),
	Major3: int32(proto.MageMajorGlyph_GlyphOfEternalWater),
	Major2: int32(proto.MageMajorGlyph_GlyphOfMoltenArmor),
}

var PlayerOptionsFire = &proto.Player_Mage{
	Mage: &proto.Mage{
		Options: &proto.Mage_Options{
			Armor:          proto.Mage_Options_MoltenArmor,
			IgniteMunching: true,
		},
		Rotation: &proto.Mage_Rotation{},
	},
}

var PlayerOptionsFrost = &proto.Player_Mage{
	Mage: &proto.Mage{
		Options: &proto.Mage_Options{
			Armor: proto.Mage_Options_MageArmor,
		},
		Rotation: &proto.Mage_Rotation{},
	},
}

var PlayerOptionsArcane = &proto.Player_Mage{
	Mage: &proto.Mage{
		Options: &proto.Mage_Options{
			Armor: proto.Mage_Options_MoltenArmor,
		},
		Rotation: &proto.Mage_Rotation{},
	},
}

var FullFireConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	Food:          proto.Food_FoodFirecrackerSalmon,
	DefaultPotion: proto.Potions_PotionOfSpeed,
}
var FullFrostConsumes = FullFireConsumes

var FullArcaneConsumes = FullFireConsumes

var P3FireGear = core.EquipmentSpecFromJsonString(`{"items": [
		{"id":51281,"enchant":3820,"gems":[41285,40133]},
        {"id":50724,"gems":[40133]},
        {"id":51284,"enchant":3810,"gems":[40153]},
        {"id":50628,"enchant":3722,"gems":[40153]},
        {"id":51283,"enchant":3832,"gems":[40113,40133]},
        {"id":54582,"enchant":2332,"gems":[40155,0]},
        {"id":50722,"enchant":3604,"gems":[40153,40133,0]},
        {"id":50613,"enchant":3601,"gems":[40133,40113,40113]},
        {"id":51282,"enchant":3872,"gems":[40133,40153]},
        {"id":50699,"enchant":3606,"gems":[40133,40113]},
        {"id":50664,"gems":[40133]},
        {"id":50398,"gems":[40153]},
        {"id":47188},
        {"id":50348},
        {"id":50732,"enchant":3834,"gems":[40113]},
        {"id":50719},
        {"id":50684,"gems":[40153]}
]}`)
var P3FrostGear = P3FireGear
var P3ArcaneGear = P3FireGear
