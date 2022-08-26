package rogue

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get item effects included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterRogue()
}

func TestCombat(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:       proto.Class_ClassRogue,
		Race:        proto.Race_RaceHuman,
		OtherRaces:  []proto.Race{proto.Race_RaceOrc},
		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		SpecOptions: core.SpecOptionsCombo{Label: "MH Deadly OH Instant", SpecOptions: PlayerOptionsCombatDI},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "MH Instant OH Deadly", SpecOptions: PlayerOptionsCombatID},
			{Label: "MH Instant OH Instant", SpecOptions: PlayerOptionsCombatII},
			{Label: "MH Deadly OH Deadly", SpecOptions: PlayerOptionsCombatDD},
		},
		RaidBuffs:   FullRaidBuffs,
		PartyBuffs:  FullPartyBuffs,
		PlayerBuffs: FullIndividualBuffs,
		Consumes:    FullConsumes,
		Debuffs:     FullDebuffs,
		ItemFilter: core.ItemFilter{
			ArmorType: proto.ArmorType_ArmorTypeLeather,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeBow,
				proto.RangedWeaponType_RangedWeaponTypeCrossbow,
				proto.RangedWeaponType_RangedWeaponTypeGun,
			},
		},
	}))
}

func TestAssassination(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:       proto.Class_ClassRogue,
		Race:        proto.Race_RaceHuman,
		OtherRaces:  []proto.Race{proto.Race_RaceOrc},
		GearSet:     core.GearSetCombo{Label: "P1 Assassination", GearSet: P1Gear},
		SpecOptions: core.SpecOptionsCombo{Label: "Assassination", SpecOptions: PlayerOptionsAssassinationDI},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "MH Instant OH Deadly", SpecOptions: PlayerOptionsAssassinationID},
			{Label: "MH Instant OH Instant", SpecOptions: PlayerOptionsAssassinationII},
			{Label: "MH Deadly OH Deadly", SpecOptions: PlayerOptionsAssassinationDD},
		},
		RaidBuffs:   FullRaidBuffs,
		PartyBuffs:  FullPartyBuffs,
		PlayerBuffs: FullIndividualBuffs,
		Consumes:    FullConsumes,
		Debuffs:     FullDebuffs,
		ItemFilter: core.ItemFilter{
			ArmorType: proto.ArmorType_ArmorTypeLeather,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeBow,
				proto.RangedWeaponType_RangedWeaponTypeCrossbow,
				proto.RangedWeaponType_RangedWeaponTypeGun,
			},
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeFist,
			},
		},
	}))
}

func TestRotation(t *testing.T) {
	specOptions := RotationSpecOptions(PlayerOptionsCombatDI.Rogue.Talents, PlayerOptionsCombatDI.Rogue.Options)
	specOptions = append(specOptions, RotationSpecOptions(PlayerOptionsAssassinationDI.Rogue.Talents, PlayerOptionsAssassinationDI.Rogue.Options)...)
	core.RunTestSuite(t, t.Name(), core.RotationTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:            proto.Class_ClassRogue,
		Race:             proto.Race_RaceHuman,
		OtherRaces:       []proto.Race{proto.Race_RaceOrc},
		GearSet:          core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		SpecOptions:      specOptions[0],
		OtherSpecOptions: specOptions[1:],
		RaidBuffs:        FullRaidBuffs,
		PartyBuffs:       FullPartyBuffs,
		PlayerBuffs:      FullIndividualBuffs,
		Consumes:         FullConsumes,
		Debuffs:          FullDebuffs,
	}))
}

type AttackType int

const (
	Poison AttackType = iota
	MHAuto
	OHAuto
	Builder
	Finisher
)

func GenerateCriticalDamageMultiplierTestCase(
	t *testing.T,
	testName string,
	equipment *proto.EquipmentSpec,
	spec *proto.Player_Rogue,
	attackType AttackType,
	expectedMultiplier float64) {
	raid := core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
		Class:     proto.Class_ClassRogue,
		Race:      proto.Race_RaceOrc,
		Equipment: equipment,
	}, spec), nil, nil, nil)
	encounter := core.MakeSingleTargetEncounter(0.0)
	env, _ := core.NewEnvironment(*raid, *encounter)
	agent := env.Raid.Parties[0].Players[0]
	rog := agent.(RogueAgent).GetRogue()
	actualMultiplier := 0.0
	switch attackType {
	case Poison:
		actualMultiplier = rog.SpellCritMultiplier()
	case MHAuto:
		actualMultiplier = rog.MeleeCritMultiplier(true, false)
	case OHAuto:
		actualMultiplier = rog.MeleeCritMultiplier(true, false)
	case Builder:
		actualMultiplier = rog.MeleeCritMultiplier(true, true)
	case Finisher:
		actualMultiplier = rog.MeleeCritMultiplier(true, false)
	}
	t.Run(testName, func(t *testing.T) {
		if !core.WithinToleranceFloat64(expectedMultiplier, actualMultiplier, 0.0001) {
			t.Logf("Crit damage multiplier for %s expected %f but was %f", testName, expectedMultiplier, actualMultiplier)
			t.Fail()
		}
	})

}

// Verifies the critical damage multipliers conform to
// https://github.com/where-fore/rogue-wotlk/issues/31
func TestCritDamageMultipliers(t *testing.T) {
	// Poison, no RED
	GenerateCriticalDamageMultiplierTestCase(t, "Poison", GearWithoutRED, PlayerOptionsNoPotW, Poison, 1.5)
	// Poison, with RED
	GenerateCriticalDamageMultiplierTestCase(t, "PoisonRED", GearWithRED, PlayerOptionsNoPotW, Poison, 1.545000)
	// Poison, with RED & PotW
	GenerateCriticalDamageMultiplierTestCase(t, "PoisonREDPotW", GearWithRED, PlayerOptionsCombatDI, Poison, 1.854000)
	// Auto, no RED, no Lethality, no PotW
	GenerateCriticalDamageMultiplierTestCase(t, "Auto", GearWithoutRED, PlayerOptionsNoLethalityNoPotW, MHAuto, 2.0)
	// Auto, RED, no Lethality, no PotW
	GenerateCriticalDamageMultiplierTestCase(t, "AutoRED", GearWithRED, PlayerOptionsNoLethalityNoPotW, MHAuto, 2.06)
	// Auto, RED, no Lethality, PotW
	GenerateCriticalDamageMultiplierTestCase(t, "AutoREDPotW", GearWithRED, PlayerOptionsNoLethality, MHAuto, 2.472)
	// Builder, no RED, Lethality, no PotW
	GenerateCriticalDamageMultiplierTestCase(t, "BuilderLethality", GearWithoutRED, PlayerOptionsNoPotW, Builder, 2.3)
	// Builder, RED, Lethality, no PotW
	GenerateCriticalDamageMultiplierTestCase(t, "BuilderREDLethality", GearWithRED, PlayerOptionsNoPotW, Builder, 2.378000)
	// Builder, no RED, Lethality, PotW
	GenerateCriticalDamageMultiplierTestCase(t, "BuilderLethalityPotW", GearWithoutRED, PlayerOptionsCombatDI, Builder, 2.820000)
	// Builder, RED, Lethality, PotW
	GenerateCriticalDamageMultiplierTestCase(t, "BuilderREDLethalityPotW", GearWithRED, PlayerOptionsCombatDI, Builder, 2.913600)
	// Finisher, no RED, Lethality, PotW
	GenerateCriticalDamageMultiplierTestCase(t, "FinisherLethalityPotW", GearWithoutRED, PlayerOptionsCombatDI, Finisher, 2.4)
	// Finisher, no RED, Lethality, PotW
	GenerateCriticalDamageMultiplierTestCase(t, "FinisherREDLethalityPotW", GearWithRED, PlayerOptionsCombatDI, Finisher, 2.472)
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:      proto.Race_RaceTroll,
				Class:     proto.Class_ClassRogue,
				Equipment: P1Gear,
				Consumes:  FullConsumes,
				Spec:      PlayerOptionsCombatDI,
				Buffs:     FullIndividualBuffs,
			},
			FullPartyBuffs,
			FullRaidBuffs,
			FullDebuffs),
		Encounter: &proto.Encounter{
			Duration: 300,
			Targets: []*proto.Target{
				core.NewDefaultTarget(),
			},
		},
		SimOptions: core.AverageDefaultSimTestOptions,
	}

	core.RaidBenchmark(b, rsr)
}
