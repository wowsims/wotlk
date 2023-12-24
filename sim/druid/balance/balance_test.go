package balance

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get caster sets included. (we use spellfire here)
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
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
		Consumes:      FullConsumes,
		SpecOptions:   core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsAdaptive},
		Rotation:      core.GetAplRotation("../../../ui/balance_druid/apls", "default"),

		ItemFilter: ItemFilter,
	}))
}

var StandardTalents = "5000500302551351--50050312"

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
