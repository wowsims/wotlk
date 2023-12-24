package healing

import (
	_ "github.com/wowsims/sod/sim/common" // imported to get caster sets included.
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterHealingPriest()
}

// TODO: Classic
// func TestDisc(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
// 		Class:    proto.Class_ClassPriest,
// 		Race:     proto.Race_RaceUndead,
// 		IsHealer: true,

// 		GearSet:     core.GetGearSet("../../../ui/healing_priest/gear_sets", "p1_disc"),
// 		Talents:     DiscTalents,
// 		Consumes:    FullConsumes,
// 		SpecOptions: core.SpecOptionsCombo{Label: "Disc", SpecOptions: PlayerOptionsDisc},
// 		Rotation:    core.GetAplRotation("../../../ui/healing_priest/apls", "disc"),

// 		ItemFilter: core.ItemFilter{
// 			WeaponTypes: []proto.WeaponType{
// 				proto.WeaponType_WeaponTypeDagger,
// 				proto.WeaponType_WeaponTypeMace,
// 				proto.WeaponType_WeaponTypeOffHand,
// 				proto.WeaponType_WeaponTypeStaff,
// 			},
// 			ArmorType: proto.ArmorType_ArmorTypeCloth,
// 			RangedWeaponTypes: []proto.RangedWeaponType{
// 				proto.RangedWeaponType_RangedWeaponTypeWand,
// 			},
// 		},

// 		EPReferenceStat: proto.Stat_StatSpellPower,
// 		StatsToWeigh: []proto.Stat{
// 			proto.Stat_StatIntellect,
// 			proto.Stat_StatSpellPower,
// 			proto.Stat_StatSpellHaste,
// 			proto.Stat_StatSpellCrit,
// 		},
// 	}))
// }

// func TestHoly(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
// 		Class:    proto.Class_ClassPriest,
// 		Race:     proto.Race_RaceUndead,
// 		IsHealer: true,

// 		GearSet:     core.GetGearSet("../../../ui/healing_priest/gear_sets", "p1_holy"),
// 		Talents:     HolyTalents,
// 		Consumes:    FullConsumes,
// 		SpecOptions: core.SpecOptionsCombo{Label: "Holy", SpecOptions: PlayerOptionsHoly},
// 		Rotation:    core.GetAplRotation("../../../ui/healing_priest/apls", "holy"),

// 		ItemFilter: core.ItemFilter{
// 			WeaponTypes: []proto.WeaponType{
// 				proto.WeaponType_WeaponTypeDagger,
// 				proto.WeaponType_WeaponTypeMace,
// 				proto.WeaponType_WeaponTypeOffHand,
// 				proto.WeaponType_WeaponTypeStaff,
// 			},
// 			ArmorType: proto.ArmorType_ArmorTypeCloth,
// 			RangedWeaponTypes: []proto.RangedWeaponType{
// 				proto.RangedWeaponType_RangedWeaponTypeWand,
// 			},
// 		},
// 	}))
// }

var DiscTalents = "0503203130300512301313231251-2351010303"
var HolyTalents = "05032031103-234051032002152530004311051"

var FullConsumes = &proto.Consumes{
	Flask: proto.Flask_FlaskUnknown,
	Food:  proto.Food_FoodUnknown,
}

var PlayerOptionsDisc = &proto.Player_HealingPriest{
	HealingPriest: &proto.HealingPriest{
		Options: &proto.HealingPriest_Options{
			UseInnerFire:      true,
			UseShadowfiend:    true,
			RapturesPerMinute: 5,
		},
		Rotation: &proto.HealingPriest_Rotation{},
	},
}

var PlayerOptionsHoly = &proto.Player_HealingPriest{
	HealingPriest: &proto.HealingPriest{
		Options: &proto.HealingPriest_Options{
			UseInnerFire:   true,
			UseShadowfiend: true,
		},
		Rotation: &proto.HealingPriest_Rotation{},
	},
}
