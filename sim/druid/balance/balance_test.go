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

		GearSet: core.GetGearSet("../../../ui/balance_druid/gear_sets", "p1"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/balance_druid/gear_sets", "p2"),
			core.GetGearSet("../../../ui/balance_druid/gear_sets", "p3_alliance"),
		},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsAdaptive},
		Rotation:    core.GetAplRotation("../../../ui/balance_druid/apls", "basic_p3"),

		ItemFilter: ItemFilter,
	}))
}

func TestBalancePhase3(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDruid,
		Race:  proto.Race_RaceTauren,

		GearSet: core.GetGearSet("../../../ui/balance_druid/gear_sets", "p3_alliance"),
		Talents: "5102233115331303213305311031--205003002",
		Glyphs: &proto.Glyphs{
			Major1: int32(proto.DruidMajorGlyph_GlyphOfStarfire),
			Major2: int32(proto.DruidMajorGlyph_GlyphOfMoonfire),
			Major3: int32(proto.DruidMajorGlyph_GlyphOfStarfall),
		},
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsAdaptive},
		Rotation:    core.GetAplRotation("../../../ui/balance_druid/apls", "basic_p3"),

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

// Extra point in Owlkin Frenzy for testing
var StandardTalents = "5012203115331303213315311231--205003012"
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
		Options: &proto.BalanceDruid_Options{
			OkfUptime: 0.2,
		},
	},
}

var ItemFilter = core.ItemFilter{
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
}
