package mage

import (
	"testing"

	_ "github.com/wowsims/tbc/sim/common"
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

func init() {
	RegisterMage()
}

func TestArcane(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,

		Race: proto.Race_RaceTroll10,

		GearSet: core.GearSetCombo{Label: "P1Arcane", GearSet: P1ArcaneGear},

		SpecOptions: core.SpecOptionsCombo{Label: "ArcaneRotation", SpecOptions: PlayerOptionsArcane},
		OtherSpecOptions: []core.SpecOptionsCombo{
			core.SpecOptionsCombo{Label: "AOE", SpecOptions: PlayerOptionsArcaneAOE},
		},

		RaidBuffs:   FullRaidBuffs,
		PartyBuffs:  FullArcanePartyBuffs,
		PlayerBuffs: FullIndividualBuffs,
		Consumes:    FullArcaneConsumes,
		Debuffs:     FullDebuffs,

		ItemFilter: core.ItemFilter{
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
		},
	}))
}

func TestFire(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,

		Race: proto.Race_RaceTroll10,

		GearSet: core.GearSetCombo{Label: "P1Fire", GearSet: P1FireGear},

		SpecOptions: core.SpecOptionsCombo{Label: "FireRotation", SpecOptions: PlayerOptionsFire},
		OtherSpecOptions: []core.SpecOptionsCombo{
			core.SpecOptionsCombo{Label: "AOE", SpecOptions: PlayerOptionsFireAOE},
		},

		RaidBuffs:   FullRaidBuffs,
		PartyBuffs:  FullFirePartyBuffs,
		PlayerBuffs: FullIndividualBuffs,
		Consumes:    FullFireConsumes,
		Debuffs:     FullDebuffs,

		ItemFilter: core.ItemFilter{
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
		},
	}))
}

func TestFrost(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,

		Race: proto.Race_RaceTroll10,

		GearSet: core.GearSetCombo{Label: "P1Frost", GearSet: P1FrostGear},

		SpecOptions: core.SpecOptionsCombo{Label: "FrostRotation", SpecOptions: PlayerOptionsFrost},
		OtherSpecOptions: []core.SpecOptionsCombo{
			core.SpecOptionsCombo{Label: "AOE", SpecOptions: PlayerOptionsFrostAOE},
		},

		RaidBuffs:   FullRaidBuffs,
		PartyBuffs:  FullFrostPartyBuffs,
		PlayerBuffs: FullIndividualBuffs,
		Consumes:    FullFrostConsumes,
		Debuffs:     FullDebuffs,

		ItemFilter: core.ItemFilter{
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
		},
	}))
}
