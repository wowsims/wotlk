package mage

import (
	_ "github.com/wowsims/classic/sim/common"
)

func init() {
	RegisterMage()
}

// TODO: Classic mage tests
// func TestArcane(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
// 		Class: proto.Class_ClassMage,
// 		Race:  proto.Race_RaceTroll,

// 		GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p3_arcane_alliance"),
// 		Talents:     ArcaneTalents,
// 		Consumes:    FullArcaneConsumes,
// 		SpecOptions: core.SpecOptionsCombo{Label: "Arcane", SpecOptions: PlayerOptionsArcane},
// 		Rotation:    core.GetAplRotation("../../ui/mage/apls", "arcane"),
// 		OtherRotations: []core.RotationCombo{
// 			core.GetAplRotation("../../ui/mage/apls", "arcane_aoe"),
// 		},

// 		ItemFilter: ItemFilter,
// 	}))
// }

// func TestFire(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
// 		Class: proto.Class_ClassMage,
// 		Race:  proto.Race_RaceTroll,

// 		GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p3_fire_alliance"),
// 		Talents:     FireTalents,
// 		Consumes:    FullFireConsumes,
// 		SpecOptions: core.SpecOptionsCombo{Label: "Fire", SpecOptions: PlayerOptionsFire},
// 		Rotation:    core.GetAplRotation("../../ui/mage/apls", "fire"),
// 		OtherRotations: []core.RotationCombo{
// 			core.GetAplRotation("../../ui/mage/apls", "fire_aoe"),
// 		},

// 		ItemFilter: ItemFilter,
// 	}))
// }

// func TestFrostFire(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
// 		Class: proto.Class_ClassMage,
// 		Race:  proto.Race_RaceTroll,

// 		GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p3_ffb_alliance"),
// 		Talents:     FrostFireTalents,
// 		Consumes:    FullFireConsumes,
// 		SpecOptions: core.SpecOptionsCombo{Label: "Fire", SpecOptions: PlayerOptionsFire},
// 		Rotation:    core.GetAplRotation("../../ui/mage/apls", "frostfire"),

// 		ItemFilter: ItemFilter,
// 	}))
// }

// func TestFrost(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
// 		Class: proto.Class_ClassMage,
// 		Race:  proto.Race_RaceTroll,

// 		GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "p3_frost_alliance"),
// 		Talents:     FrostTalents,
// 		Consumes:    FullFrostConsumes,
// 		SpecOptions: core.SpecOptionsCombo{Label: "Frost", SpecOptions: PlayerOptionsFrost},
// 		Rotation:    core.GetAplRotation("../../ui/mage/apls", "frost"),
// 		OtherRotations: []core.RotationCombo{
// 			core.GetAplRotation("../../ui/mage/apls", "frost_aoe"),
// 		},

// 		ItemFilter: ItemFilter,
// 	}))
// }

// var ItemFilter = core.ItemFilter{
// 	WeaponTypes: []proto.WeaponType{
// 		proto.WeaponType_WeaponTypeDagger,
// 		proto.WeaponType_WeaponTypeSword,
// 		proto.WeaponType_WeaponTypeOffHand,
// 		proto.WeaponType_WeaponTypeStaff,
// 	},
// 	ArmorType: proto.ArmorType_ArmorTypeCloth,
// 	RangedWeaponTypes: []proto.RangedWeaponType{
// 		proto.RangedWeaponType_RangedWeaponTypeWand,
// 	},
// }

// var ArcaneTalents = "23000513310033015032310250532-03-023303001"
// var FireTalents = "23000503110003-0055030012303331053120301351"
// var FrostFireTalents = "23000503110003-0055030012303331053120301351"
// var FrostTalents = "23000503110003--0533030310233100030152231351"

// var PlayerOptionsFire = &proto.Player_Mage{
// 	Mage: &proto.Mage{
// 		Options: &proto.Mage_Options{
// 			Armor: proto.Mage_Options_MageArmor,
// 		},
// 		Rotation: &proto.Mage_Rotation{},
// 	},
// }

// var PlayerOptionsFrost = &proto.Player_Mage{
// 	Mage: &proto.Mage{
// 		Options: &proto.Mage_Options{
// 			Armor: proto.Mage_Options_MageArmor,
// 		},
// 		Rotation: &proto.Mage_Rotation{},
// 	},
// }

// var PlayerOptionsArcane = &proto.Player_Mage{
// 	Mage: &proto.Mage{
// 		Options: &proto.Mage_Options{
// 			Armor: proto.Mage_Options_MageArmor,
// 		},
// 		Rotation: &proto.Mage_Rotation{},
// 	},
// }

// var FullFireConsumes = &proto.Consumes{
// 	Flask:         proto.Flask_FlaskUnknown,
// 	Food:          proto.Food_FoodUnknown,
// }
// var FullFrostConsumes = FullFireConsumes
// var FullArcaneConsumes = FullFireConsumes
