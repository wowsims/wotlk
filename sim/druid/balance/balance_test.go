package balance

import (
	"testing"

	_ "github.com/wowsims/classic/sim/common" // imported to get caster sets included. (we use spellfire here)
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
)

func init() {
	RegisterBalanceDruid()
}

func TestBalance(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDruid,
		Race:  proto.Race_RaceTauren,

		GearSet:       core.GetGearSet("../../../ui/balance_druid/gear_sets", "blank"),
		OtherGearSets: []core.GearSetCombo{},
		Talents:       StandardTalents,
		Glyphs:        StandardGlyphs,
		Consumes:      FullConsumes,
		SpecOptions:   core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsAdaptive},
		Rotation:      core.GetAplRotation("../../../ui/balance_druid/apls", "default"),

		ItemFilter: ItemFilter,
	}))
}

var StandardTalents = "5000500302551351--50050312"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.DruidMajorGlyph_DruidMajorGlyphNone),
	Major2: int32(proto.DruidMajorGlyph_DruidMajorGlyphNone),
	Major3: int32(proto.DruidMajorGlyph_DruidMajorGlyphNone),
	Minor1: int32(proto.DruidMinorGlyph_DruidMinorGlyphNone),
}

var FullConsumes = &proto.Consumes{
	Flask: proto.Flask_FlaskUnknown,
	Food:  proto.Food_FoodUnknown,
}

var PlayerOptionsAdaptive = &proto.Player_BalanceDruid{
	BalanceDruid: &proto.BalanceDruid{
		Options: &proto.BalanceDruid_Options{
			OkfUptime: 0.2,
		},
		Rotation: &proto.BalanceDruid_Rotation{},
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
