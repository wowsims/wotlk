package shadow

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get caster sets included.
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterShadowPriest()
}

func TestShadow(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassPriest,
		Race:       proto.Race_RaceUndead,
		OtherRaces: []proto.Race{proto.Race_RaceNightElf},

		GearSet:  core.GetGearSet("../../../ui/shadow_priest/gear_sets", "blank"),
		Talents:  DefaultTalents,
		Consumes: FullConsumes,

		SpecOptions: core.SpecOptionsCombo{Label: "APL", SpecOptions: PlayerOptionsBasic},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "Basic", SpecOptions: PlayerOptionsBasic},
		},

		Rotation: core.GetAplRotation("../../../ui/shadow_priest/apls", "default"),

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

var DefaultTalents = "5042001303--5002505103501051"

var FullConsumes = &proto.Consumes{
	Flask: proto.Flask_FlaskUnknown,
	Food:  proto.Food_FoodUnknown,
}

var PlayerOptionsBasic = &proto.Player_ShadowPriest{
	ShadowPriest: &proto.ShadowPriest{
		Options: &proto.ShadowPriest_Options{
			Armor: proto.ShadowPriest_Options_InnerFire,
		},
	},
}
