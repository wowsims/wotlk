package smite

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get caster sets included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterSmitePriest()
}

func TestSmite(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassPriest,
		Race:  proto.Race_RaceUndead,

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     DefaultTalents,
		Glyphs:      DefaultGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeMace,
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

var DefaultTalents = "05332031013005023310001-005551002020152-00502"
var DefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.PriestMajorGlyph_GlyphOfSmite),
	Major2: int32(proto.PriestMajorGlyph_GlyphOfHolyNova),
	Major3: int32(proto.PriestMajorGlyph_GlyphOfShadowWordDeath),
	// No interesting minor glyphs.
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	Food:          proto.Food_FoodFishFeast,
	DefaultPotion: proto.Potions_RunicManaInjector,
	PrepopPotion:  proto.Potions_PotionOfWildMagic,
}

var PlayerOptionsBasic = &proto.Player_SmitePriest{
	SmitePriest: &proto.SmitePriest{
		Options: &proto.SmitePriest_Options{
			UseInnerFire:   true,
			UseShadowfiend: true,
		},
		Rotation: &proto.SmitePriest_Rotation{
			UseDevouringPlague: true,
			UseShadowWordDeath: true,
			UseMindBlast:       true,

			AllowedHolyFireDelayMs: 50,
		},
	},
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40562,
		"enchant": 3820,
		"gems": [
			41333,
			42144
		]
	},
	{
		"id": 44661,
		"gems": [
			39998
		]
	},
	{
		"id": 40459,
		"enchant": 3810,
		"gems": [
			42144
		]
	},
	{
		"id": 44005,
		"enchant": 3859,
		"gems": [
			42144
		]
	},
	{
		"id": 40234,
		"enchant": 1144,
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
		"id": 40454,
		"enchant": 3604,
		"gems": [
			40049,
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
		"id": 40560,
		"enchant": 3719
	},
	{
		"id": 40558,
		"enchant": 3826
	},
	{
		"id": 40719
	},
	{
		"id": 40399
	},
	{
		"id": 40255
	},
	{
		"id": 40432
	},
	{
		"id": 40395,
		"enchant": 3834
	},
	{
		"id": 40273
	},
	{
		"id": 39712
	}
]}`)
