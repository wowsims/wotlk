package restoration

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get caster sets included. (we use spellfire here)
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterRestorationDruid()
}

func TestRestoration(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDruid,
		Race:  proto.Race_RaceTauren,

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Standard", SpecOptions: PlayerOptionsStandard},

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
				proto.WeaponType_WeaponTypePolearm,
			},
			ArmorType: proto.ArmorType_ArmorTypeLeather,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeIdol,
			},
		},
	}))
}

var StandardTalents = "05320031103--230023312131502331050313051"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.DruidMajorGlyph_GlyphOfWildGrowth),
	Major2: int32(proto.DruidMajorGlyph_GlyphOfSwiftmend),
	Major3: int32(proto.DruidMajorGlyph_GlyphOfNourish),
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfBlindingLight,
	Food:            proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:   proto.Potions_SuperManaPotion,
	PrepopPotion:    proto.Potions_DestructionPotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var PlayerOptionsStandard = &proto.Player_RestorationDruid{
	RestorationDruid: &proto.RestorationDruid{
		Options: &proto.RestorationDruid_Options{
			InnervateTarget: &proto.UnitReference{Type: proto.UnitReference_Player, Index: 0}, // self innervate
		},
		Rotation: &proto.RestorationDruid_Rotation{},
	},
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 44007,
		"enchant": 3819,
		"gems": [
			41401,
			40017
		]
	},
	{
		"id": 40071
	},
	{
		"id": 39719,
		"enchant": 3809,
		"gems": [
			39998
		]
	},
	{
		"id": 40723,
		"enchant": 3859
	},
	{
		"id": 44002,
		"enchant": 3832,
		"gems": [
			39998,
			40026
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
		"id": 40460,
		"enchant": 3246,
		"gems": [
			40017,
			0
		]
	},
	{
		"id": 40561,
		"enchant": 3601,
		"gems": [
			39998
		]
	},
	{
		"id": 40379,
		"enchant": 3719,
		"gems": [
			39998,
			40017
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
		"id": 40375
	},
	{
		"id": 37111
	},
	{
		"id": 40432
	},
	{
		"id": 40395,
		"enchant": 3834
	},
	{
		"id": 39766
	},
	{
		"id": 40342
	}
]}`)
