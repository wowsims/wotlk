package balance

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get caster sets included. (we use spellfire here)
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterBalanceDruid()
}

func TestBalance(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDruid,
		Race:  proto.Race_RaceTauren,

		GearSet: core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		OtherGearSets: []core.GearSetCombo{
			{Label: "P2", GearSet: P2Gear},
			{Label: "P2-4P", GearSet: P2Gear4P},
		},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsAdaptive},

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

var StandardTalents = "5012203115331303213305311231--205003012"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.DruidMajorGlyph_GlyphOfStarfire),
	Major2: int32(proto.DruidMajorGlyph_GlyphOfInsectSwarm),
	Major3: int32(proto.DruidMajorGlyph_GlyphOfStarfall),
	Minor1: int32(proto.DruidMinorGlyph_GlyphOfTyphoon),
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	Food:          proto.Food_FoodFishFeast,
	DefaultPotion: proto.Potions_PotionOfSpeed,
	PrepopPotion:  proto.Potions_PotionOfWildMagic,
}

var PlayerOptionsAdaptive = &proto.Player_BalanceDruid{
	BalanceDruid: &proto.BalanceDruid{
		Options: &proto.BalanceDruid_Options{},
		Rotation: &proto.BalanceDruid_Rotation{
			MfUsage:            proto.BalanceDruid_Rotation_BeforeLunar,
			IsUsage:            proto.BalanceDruid_Rotation_MaximizeIs,
			WrathUsage:         proto.BalanceDruid_Rotation_RegularWrath,
			UseBattleRes:       false,
			UseStarfire:        true,
			UseTyphoon:         false,
			UseHurricane:       false,
			UseSmartCooldowns:  true,
			MaintainFaerieFire: true,
			PlayerLatency:      200,
		},
	},
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40467,
		"enchant": 3820,
		"gems": [
			41285,
			42144
		]
	},
	{
		"id": 44661,
		"gems": [
			40026
		]
	},
	{
		"id": 40470,
		"enchant": 3810,
		"gems": [
			42144
		]
	},
	{
		"id": 44005,
		"enchant": 3859,
		"gems": [
			40026
		]
	},
	{
		"id": 40469,
		"enchant": 3832,
		"gems": [
			42144,
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
		"id": 40466,
		"enchant": 3604,
		"gems": [
			39998,
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
		"id": 40519,
		"enchant": 3606
	},
	{
		"id": 40399
	},
	{
		"id": 40080
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
		"id": 40192
	},
	{
		"id": 40321
	}
]}`)

var P2Gear = core.EquipmentSpecFromJsonString(` {
      "items": [
        {
          "id": 45497,
          "enchant": 3820,
          "gems": [
            41285,
            42144
          ]
        },
        {
          "id": 45133,
          "gems": [
            40048
          ]
        },
        {
          "id": 46196,
          "enchant": 3810,
          "gems": [
            39998
          ]
        },
        {
          "id": 45242,
          "enchant": 3859,
          "gems": [
            40048
          ]
        },
        {
          "id": 45519,
          "enchant": 3832,
          "gems": [
            40051,
            42144,
            40026
          ]
        },
        {
          "id": 45446,
          "enchant": 2332,
          "gems": [
            42144,
            0
          ]
        },
        {
          "id": 45665,
          "enchant": 3604,
          "gems": [
            39998,
            39998,
            0
          ]
        },
        {
          "id": 45619,
          "gems": [
            39998,
            39998,
            39998
          ]
        },
        {
          "id": 46192,
          "enchant": 3719,
          "gems": [
            39998,
            39998
          ]
        },
        {
          "id": 45537,
          "enchant": 3606,
          "gems": [
            39998,
            40026
          ]
        },
        {
          "id": 46046,
          "gems": [
            39998
          ]
        },
        {
          "id": 45495,
          "gems": [
            39998
          ]
        },
        {
          "id": 45466
        },
        {
          "id": 45518
        },
        {
          "id": 45620,
          "enchant": 3834,
          "gems": [
            39998
          ]
        },
        {
          "id": 45617
        },
        {
          "id": 40321
        }
      ]
    }`)

var P2Gear4P = core.EquipmentSpecFromJsonString(` {
       "items": [
        {
          "id": 46191,
          "enchant": 3820,
          "gems": [
            41285,
            42144
          ]
        },
        {
          "id": 45933,
          "gems": [
            39998
          ]
        },
        {
          "id": 46196,
          "enchant": 3810,
          "gems": [
            40026
          ]
        },
        {
          "id": 45242,
          "enchant": 3859,
          "gems": [
            39998
          ]
        },
        {
          "id": 46194,
          "enchant": 3832,
          "gems": [
            39998,
            42144
          ]
        },
        {
          "id": 45446,
          "enchant": 2332,
          "gems": [
            42144,
            0
          ]
        },
        {
          "id": 45665,
          "enchant": 3604,
          "gems": [
            39998,
            39998,
            0
          ]
        },
        {
          "id": 45616,
          "gems": [
            39998,
            39998,
            39998
          ]
        },
        {
          "id": 46192,
          "enchant": 3719,
          "gems": [
            39998,
            39998
          ]
        },
        {
          "id": 45537,
          "enchant": 3606,
          "gems": [
            39998,
            40026
          ]
        },
        {
          "id": 46046,
          "gems": [
            39998
          ]
        },
        {
          "id": 45495,
          "gems": [
            39998
          ]
        },
        {
          "id": 45466
        },
        {
          "id": 45518
        },
        {
          "id": 45620,
          "enchant": 3834,
          "gems": [
            39998
          ]
        },
        {
          "id": 45617
        },
        {
          "id": 40321
        }
      ]
    }`)
