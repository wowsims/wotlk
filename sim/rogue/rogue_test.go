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
		Talents:     CombatTalents,
		Glyphs:      CombatGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "MH Deadly OH Instant", SpecOptions: PlayerOptionsCombatDI},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "MH Instant OH Deadly", SpecOptions: PlayerOptionsCombatID},
			{Label: "MH Instant OH Instant", SpecOptions: PlayerOptionsCombatII},
			{Label: "MH Deadly OH Deadly", SpecOptions: PlayerOptionsCombatDD},
		},
		Rotation: core.RotationCombo{Label: "Combat", Rotation: CombatExposeRotation},
		OtherRotations: []core.RotationCombo{
			{Label: "CleaveSND", Rotation: CleaveSNDRotation},
			{Label: "AOE", Rotation: AOERotation},
		},
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
		Talents:     AssassinationTalents,
		Glyphs:      AssassinationGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Assassination", SpecOptions: PlayerOptionsAssassinationDI},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "MH Instant OH Deadly", SpecOptions: PlayerOptionsAssassinationID},
			{Label: "MH Instant OH Instant", SpecOptions: PlayerOptionsAssassinationII},
			{Label: "MH Deadly OH Deadly", SpecOptions: PlayerOptionsAssassinationDD},
		},
		Rotation: core.RotationCombo{Label: "Mutilate", Rotation: MutilateRuptureExposeRotation},

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

func TestSubtlety(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:       proto.Class_ClassRogue,
		Race:        proto.Race_RaceBloodElf,
		OtherRaces:  []proto.Race{proto.Race_RaceOrc},
		GearSet:     core.GearSetCombo{Label: "P2 Subtlety", GearSet: SubtletyP2Gear},
		Talents:     SubtletyTalents,
		Glyphs:      SubtletyGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Subtlety", SpecOptions: PlayerOptionsSubtletyID},
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
	talents string,
	spec *proto.Player_Rogue,
	attackType AttackType,
	expectedMultiplier float64) {
	raid := core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
		Class:         proto.Class_ClassRogue,
		Race:          proto.Race_RaceOrc,
		Equipment:     equipment,
		TalentsString: talents,
	}, spec), nil, nil, nil)
	encounter := core.MakeSingleTargetEncounter(0.0)
	env, _, _ := core.NewEnvironment(raid, encounter)
	agent := env.Raid.Parties[0].Players[0]
	rog := agent.(RogueAgent).GetRogue()
	actualMultiplier := 0.0
	switch attackType {
	case Poison:
		actualMultiplier = rog.SpellCritMultiplier()
	case MHAuto:
		actualMultiplier = rog.MeleeCritMultiplier(false)
	case OHAuto:
		actualMultiplier = rog.MeleeCritMultiplier(false)
	case Builder:
		actualMultiplier = rog.MeleeCritMultiplier(true)
	case Finisher:
		actualMultiplier = rog.MeleeCritMultiplier(false)
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
	GenerateCriticalDamageMultiplierTestCase(t, "Poison", GearWithoutRED, CombatNoPotWTalents, PlayerOptionsNoPotW, Poison, 1.5)
	// Poison, with RED
	GenerateCriticalDamageMultiplierTestCase(t, "PoisonRED", GearWithRED, CombatNoPotWTalents, PlayerOptionsNoPotW, Poison, 1.545000)
	// Poison, with RED & PotW
	GenerateCriticalDamageMultiplierTestCase(t, "PoisonREDPotW", GearWithRED, CombatTalents, PlayerOptionsCombatDI, Poison, 1.854000)
	// Auto, no RED, no Lethality, no PotW
	GenerateCriticalDamageMultiplierTestCase(t, "Auto", GearWithoutRED, CombatNoLethalityNoPotWTalents, PlayerOptionsNoLethalityNoPotW, MHAuto, 2.0)
	// Auto, RED, no Lethality, no PotW
	GenerateCriticalDamageMultiplierTestCase(t, "AutoRED", GearWithRED, CombatNoLethalityNoPotWTalents, PlayerOptionsNoLethalityNoPotW, MHAuto, 2.06)
	// Auto, RED, no Lethality, PotW
	GenerateCriticalDamageMultiplierTestCase(t, "AutoREDPotW", GearWithRED, CombatNoLethalityTalents, PlayerOptionsNoLethality, MHAuto, 2.472)
	// Builder, no RED, Lethality, no PotW
	GenerateCriticalDamageMultiplierTestCase(t, "BuilderLethality", GearWithoutRED, CombatNoPotWTalents, PlayerOptionsNoPotW, Builder, 2.3)
	// Builder, RED, Lethality, no PotW
	GenerateCriticalDamageMultiplierTestCase(t, "BuilderREDLethality", GearWithRED, CombatNoPotWTalents, PlayerOptionsNoPotW, Builder, 2.378000)
	// Builder, no RED, Lethality, PotW
	GenerateCriticalDamageMultiplierTestCase(t, "BuilderLethalityPotW", GearWithoutRED, CombatTalents, PlayerOptionsCombatDI, Builder, 2.820000)
	// Builder, RED, Lethality, PotW
	GenerateCriticalDamageMultiplierTestCase(t, "BuilderREDLethalityPotW", GearWithRED, CombatTalents, PlayerOptionsCombatDI, Builder, 2.913600)
	// Finisher, no RED, Lethality, PotW
	GenerateCriticalDamageMultiplierTestCase(t, "FinisherLethalityPotW", GearWithoutRED, CombatTalents, PlayerOptionsCombatDI, Finisher, 2.4)
	// Finisher, no RED, Lethality, PotW
	GenerateCriticalDamageMultiplierTestCase(t, "FinisherREDLethalityPotW", GearWithRED, CombatTalents, PlayerOptionsCombatDI, Finisher, 2.472)
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
				Buffs:     core.FullIndividualBuffs,
			},
			core.FullPartyBuffs,
			core.FullRaidBuffs,
			core.FullDebuffs),
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

var CombatTalents = "00532000523-0252051050035010223100501251"
var CombatNoLethalityTalents = "00532000023-0252051050035010223100501251"
var CombatNoPotWTalents = "00532000523-0252051050035010223100501201"
var CombatNoLethalityNoPotWTalents = "00532000023-0252051050035010223100501201"
var AssassinationTalents = "005303005352100520103331051-005005003-502"
var SubtletyTalents = "30532000235--512003203032012135011503113"
var CombatGlyphs = &proto.Glyphs{
	Major1: int32(proto.RogueMajorGlyph_GlyphOfKillingSpree),
	Major2: int32(proto.RogueMajorGlyph_GlyphOfTricksOfTheTrade),
	Major3: int32(proto.RogueMajorGlyph_GlyphOfRupture),
}
var AssassinationGlyphs = &proto.Glyphs{
	Major1: int32(proto.RogueMajorGlyph_GlyphOfMutilate),
	Major2: int32(proto.RogueMajorGlyph_GlyphOfTricksOfTheTrade),
	Major3: int32(proto.RogueMajorGlyph_GlyphOfHungerForBlood),
}
var SubtletyGlyphs = &proto.Glyphs{
	Major1: int32(proto.RogueMajorGlyph_GlyphOfHemorrhage),
	Major2: int32(proto.RogueMajorGlyph_GlyphOfEviscerate),
	Major3: int32(proto.RogueMajorGlyph_GlyphOfRupture),
}

var PlayerOptionsCombatDI = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyInstant,
		Rotation: &proto.Rogue_Rotation{},
	},
}
var PlayerOptionsCombatDD = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyDeadly,
		Rotation: &proto.Rogue_Rotation{},
	},
}
var PlayerOptionsCombatID = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  InstantDeadly,
		Rotation: &proto.Rogue_Rotation{},
	},
}
var PlayerOptionsCombatII = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  InstantInstant,
		Rotation: &proto.Rogue_Rotation{},
	},
}

var PlayerOptionsNoLethality = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyInstant,
		Rotation: &proto.Rogue_Rotation{},
	},
}

var PlayerOptionsNoPotW = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyInstant,
		Rotation: &proto.Rogue_Rotation{},
	},
}

var PlayerOptionsNoLethalityNoPotW = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyInstant,
		Rotation: &proto.Rogue_Rotation{},
	},
}

var PlayerOptionsAssassinationDI = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyInstant,
		Rotation: &proto.Rogue_Rotation{},
	},
}
var PlayerOptionsAssassinationDD = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyDeadly,
		Rotation: &proto.Rogue_Rotation{},
	},
}
var PlayerOptionsAssassinationID = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  InstantDeadly,
		Rotation: &proto.Rogue_Rotation{},
	},
}
var PlayerOptionsAssassinationII = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  InstantInstant,
		Rotation: &proto.Rogue_Rotation{},
	},
}

var PlayerOptionsSubtletyID = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  InstantDeadly,
		Rotation: &proto.Rogue_Rotation{},
	},
}

var DeadlyInstant = &proto.Rogue_Options{
	MhImbue: proto.Rogue_Options_DeadlyPoison,
	OhImbue: proto.Rogue_Options_InstantPoison,
}
var InstantDeadly = &proto.Rogue_Options{
	MhImbue: proto.Rogue_Options_InstantPoison,
	OhImbue: proto.Rogue_Options_DeadlyPoison,
}
var InstantInstant = &proto.Rogue_Options{
	MhImbue: proto.Rogue_Options_InstantPoison,
	OhImbue: proto.Rogue_Options_InstantPoison,
}
var DeadlyDeadly = &proto.Rogue_Options{
	MhImbue: proto.Rogue_Options_DeadlyPoison,
	OhImbue: proto.Rogue_Options_DeadlyPoison,
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfEndlessRage,
	DefaultPotion:   proto.Potions_PotionOfSpeed,
	DefaultConjured: proto.Conjured_ConjuredRogueThistleTea,
}

var MutilateRuptureExposeRotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"prepullActions": [
		{"action":{"castSpell":{"spellId":{"itemId":40211}}},"doAtValue":{"const":{"val":"-1s"}}},
		{"action":{"activateAura":{"auraId":{"spellId":58426}}},"doAtValue":{"const":{"val":"-1s"}}}
	],
	"priorityList": [
		{"action":{"autocastOtherCooldowns":{}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":8647}}},"rhs":{"const":{"val":"1s"}}}},"castSpell":{"spellId":{"spellId":8647}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":6774}}},"rhs":{"const":{"val":"1s"}}}},"castSpell":{"spellId":{"spellId":6774}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":8647}}},"rhs":{"const":{"val":"1s"}}}},"sequence":{"actions":[{"castSpell":{"spellId":{"spellId":48666}}},{"castSpell":{"spellId":{"spellId":8647}}}]}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":6774}}},"rhs":{"const":{"val":"1s"}}}},"sequence":{"actions":[{"castSpell":{"spellId":{"spellId":48666}}},{"castSpell":{"spellId":{"spellId":6774}}}]}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":51662}}},"rhs":{"const":{"val":"1s"}}}},"castSpell":{"spellId":{"spellId":51662}}}},
		{"action":{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":58426}}}}},"castSpell":{"spellId":{"spellId":26889}}}},
		{"action":{"condition":{"auraIsActive":{"auraId":{"spellId":6774}}},"castSpell":{"spellId":{"itemId":40211}}}},
		{"action":{"condition":{"auraIsActive":{"auraId":{"spellId":6774}}},"castSpell":{"spellId":{"spellId":54758}}}},
		{"action":{"condition":{"auraIsActive":{"auraId":{"spellId":6774}}},"castSpell":{"spellId":{"spellId":13877}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"50"}}}},"castSpell":{"spellId":{"spellId":57934}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"50"}}}},"castSpell":{"spellId":{"itemId":7676}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGe","lhs":{"currentComboPoints":{}},"rhs":{"const":{"val":"4"}}}},{"not":{"val":{"auraIsActive":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":48672}}}}}]}},"castSpell":{"spellId":{"spellId":48672}}}},
		{"action":{"condition":{"or":{"vals":[{"cmp":{"op":"OpEq","lhs":{"currentComboPoints":{}},"rhs":{"const":{"val":"5"}}}},{"cmp":{"op":"OpLe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"5s"}}}}]}},"castSpell":{"spellId":{"spellId":14177}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGe","lhs":{"currentComboPoints":{}},"rhs":{"const":{"val":"4"}}}},{"or":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":57993}}}}},{"cmp":{"op":"OpGe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"85"}}}}]}}]}},"castSpell":{"spellId":{"spellId":57993}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentComboPoints":{}},"rhs":{"const":{"val":"3"}}}},"castSpell":{"spellId":{"spellId":48666}}}}
	]
}`)

var CombatExposeRotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"prepullActions": [
		{"action":{"castSpell":{"spellId":{"itemId":40211}}},"doAtValue":{"const":{"val":"-1s"}}}
	],
	"priorityList": [
		{"action":{"autocastOtherCooldowns":{}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":6774}}},"rhs":{"const":{"val":"1"}}}},{"cmp":{"op":"OpGe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"2s"}}}}]}},"castSpell":{"spellId":{"spellId":6774}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":8647}}},"rhs":{"const":{"val":"1"}}}},{"cmp":{"op":"OpGe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"2s"}}}}]}},"castSpell":{"spellId":{"spellId":8647}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":6774}}},"rhs":{"const":{"val":"1s"}}}},{"cmp":{"op":"OpGe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"2s"}}}}]}},"sequence":{"actions":[{"castSpell":{"spellId":{"spellId":48638}}},{"castSpell":{"spellId":{"spellId":6774}}}]}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":8647}}},"rhs":{"const":{"val":"1s"}}}},{"cmp":{"op":"OpGe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"2s"}}}}]}},"sequence":{"actions":[{"castSpell":{"spellId":{"spellId":48638}}},{"castSpell":{"spellId":{"spellId":8647}}}]}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":6774}}},"rhs":{"math":{"op":"OpAdd","lhs":{"dotRemainingTime":{"spellId":{"spellId":48672}}},"rhs":{"const":{"val":"2"}}}}}},{"cmp":{"op":"OpGe","lhs":{"currentComboPoints":{}},"rhs":{"const":{"val":"3"}}}},{"dotIsActive":{"spellId":{"spellId":48672}}},{"not":{"val":{"cmp":{"op":"OpLe","lhs":{"math":{"op":"OpAdd","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"2s"}}}},"rhs":{"auraRemainingTime":{"auraId":{"spellId":6774}}}}}}}]}},"castSpell":{"spellId":{"spellId":6774}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGe","lhs":{"auraRemainingTime":{"auraId":{"spellId":6774}}},"rhs":{"const":{"val":"4s"}}}},{"cmp":{"op":"OpEq","lhs":{"currentComboPoints":{}},"rhs":{"const":{"val":"5"}}}},{"not":{"val":{"dotIsActive":{"spellId":{"spellId":48672}}}}},{"cmp":{"op":"OpGt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"10s"}}}}]}},"castSpell":{"spellId":{"spellId":48672}}}},
		{"action":{"condition":{"and":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpGe","lhs":{"auraRemainingTime":{"auraId":{"spellId":6774}}},"rhs":{"const":{"val":"2s"}}}},{"cmp":{"op":"OpLt","lhs":{"auraRemainingTime":{"auraId":{"spellId":6774}}},"rhs":{"const":{"val":"4s"}}}}]}},{"cmp":{"op":"OpGe","lhs":{"currentComboPoints":{}},"rhs":{"const":{"val":"4"}}}},{"not":{"val":{"dotIsActive":{"spellId":{"spellId":48672}}}}},{"cmp":{"op":"OpGt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"8s"}}}}]}},"castSpell":{"spellId":{"spellId":48672}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGe","lhs":{"currentComboPoints":{}},"rhs":{"const":{"val":"4"}}}},{"cmp":{"op":"OpLe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"4s"}}}}]}},"castSpell":{"spellId":{"spellId":48668}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGe","lhs":{"auraRemainingTime":{"auraId":{"spellId":6774}}},"rhs":{"const":{"val":"4s"}}}},{"cmp":{"op":"OpGe","lhs":{"currentComboPoints":{}},"rhs":{"const":{"val":"4"}}}},{"cmp":{"op":"OpGe","lhs":{"auraRemainingTime":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":48672}}},"rhs":{"const":{"val":"6s"}}}},{"cmp":{"op":"OpGe","lhs":{"auraRemainingTime":{"auraId":{"spellId":13750}}},"rhs":{"const":{"val":"4s"}}}}]}},"castSpell":{"spellId":{"spellId":48668}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGe","lhs":{"auraRemainingTime":{"auraId":{"spellId":6774}}},"rhs":{"const":{"val":"4s"}}}},{"cmp":{"op":"OpGe","lhs":{"currentComboPoints":{}},"rhs":{"const":{"val":"4"}}}},{"cmp":{"op":"OpGe","lhs":{"auraRemainingTime":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":48672}}},"rhs":{"const":{"val":"10s"}}}}]}},"castSpell":{"spellId":{"spellId":48668}}}},
		{"action":{"condition":{"auraIsActive":{"auraId":{"spellId":6774}}},"castSpell":{"spellId":{"spellId":13877}}}},
		{"action":{"condition":{"or":{"vals":[{"auraIsActive":{"auraId":{"spellId":13877}}},{"cmp":{"op":"OpLe","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":13877}}},"rhs":{"const":{"val":"8s"}}}}}}]}},"castSpell":{"spellId":{"spellId":54758}}}},
		{"action":{"condition":{"and":{"vals":[{"not":{"val":{"spellIsReady":{"spellId":{"spellId":13877}}}}},{"cmp":{"op":"OpGe","lhs":{"spellTimeToReady":{"spellId":{"spellId":13877}}},"rhs":{"const":{"val":"57s"}}}}]}},"castSpell":{"spellId":{"spellId":54758}}}},
		{"action":{"condition":{"or":{"vals":[{"auraIsActive":{"auraId":{"spellId":13877}}},{"cmp":{"op":"OpLe","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":13877}}},"rhs":{"const":{"val":"10s"}}}}}}]}},"castSpell":{"spellId":{"itemId":40211}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"50"}}}},"castSpell":{"spellId":{"spellId":57934}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"50"}}}},"castSpell":{"spellId":{"spellId":51690}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"196s"}}}},{"cmp":{"op":"OpLe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"50"}}}}]}},"castSpell":{"spellId":{"spellId":13750}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"196s"}}}},{"cmp":{"op":"OpGe","lhs":{"spellTimeToReady":{"spellId":{"spellId":51690}}},"rhs":{"const":{"val":"15s"}}}},{"cmp":{"op":"OpLe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"50"}}}}]}},"castSpell":{"spellId":{"spellId":13750}}}},
		{"action":{"castSpell":{"spellId":{"spellId":48638}}}}
	]
}`)

var CleaveSNDRotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"prepullActions": [
		{"action":{"castSpell":{"spellId":{"itemId":40211}}},"doAtValue":{"const":{"val":"-1s"}}}
	],
	"priorityList": [
		{"action":{"autocastOtherCooldowns":{}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":8647}}},"rhs":{"const":{"val":"1"}}}},{"cmp":{"op":"OpGe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"2s"}}}}]}},"castSpell":{"spellId":{"spellId":8647}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":6774}}},"rhs":{"const":{"val":"1"}}}},{"cmp":{"op":"OpGe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"2s"}}}}]}},"castSpell":{"spellId":{"spellId":6774}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":8647}}},"rhs":{"const":{"val":"1s"}}}},{"cmp":{"op":"OpGe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"2s"}}}}]}},"sequence":{"actions":[{"castSpell":{"spellId":{"spellId":48638}}},{"castSpell":{"spellId":{"spellId":8647}}}]}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":6774}}},"rhs":{"const":{"val":"1s"}}}},{"cmp":{"op":"OpGe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"2s"}}}}]}},"sequence":{"actions":[{"castSpell":{"spellId":{"spellId":48638}}},{"castSpell":{"spellId":{"spellId":6774}}}]}}},
		{"action":{"condition":{"auraIsActive":{"auraId":{"spellId":6774}}},"castSpell":{"spellId":{"spellId":13877}}}},
		{"action":{"condition":{"or":{"vals":[{"auraIsActive":{"auraId":{"spellId":13877}}},{"cmp":{"op":"OpLe","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":13877}}},"rhs":{"const":{"val":"8s"}}}}}}]}},"castSpell":{"spellId":{"spellId":54758}}}},
		{"action":{"condition":{"and":{"vals":[{"not":{"val":{"spellIsReady":{"spellId":{"spellId":13877}}}}},{"cmp":{"op":"OpGe","lhs":{"spellTimeToReady":{"spellId":{"spellId":13877}}},"rhs":{"const":{"val":"57s"}}}}]}},"castSpell":{"spellId":{"spellId":54758}}}},
		{"action":{"condition":{"or":{"vals":[{"auraIsActive":{"auraId":{"spellId":13877}}},{"cmp":{"op":"OpLe","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":13877}}},"rhs":{"const":{"val":"10s"}}}}}}]}},"castSpell":{"spellId":{"itemId":40211}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"50"}}}},"castSpell":{"spellId":{"spellId":51690}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"196s"}}}},{"cmp":{"op":"OpLe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"50"}}}}]}},"castSpell":{"spellId":{"spellId":13750}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"196s"}}}},{"cmp":{"op":"OpGe","lhs":{"spellTimeToReady":{"spellId":{"spellId":51690}}},"rhs":{"const":{"val":"15s"}}}},{"cmp":{"op":"OpLe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"50"}}}}]}},"castSpell":{"spellId":{"spellId":13750}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"50"}}}},"castSpell":{"spellId":{"spellId":57934}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"50"}}}},"castSpell":{"spellId":{"itemId":7676}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLt","lhs":{"currentComboPoints":{}},"rhs":{"const":{"val":"1"}}}},{"or":{"vals":[{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":6774}}},"rhs":{"math":{"op":"OpSub","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"2s"}}}}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":8647}}},"rhs":{"math":{"op":"OpSub","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"2s"}}}}}}]}}]}},"castSpell":{"spellId":{"spellId":48638}}}},
		{"action":{"castSpell":{"spellId":{"spellId":51723}}}}
	]
}`)

var AOERotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"prepullActions": [
		{"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAtValue":{"const":{"val":"-1s"}}},
		{"action":{"activateAura":{"auraId":{"spellId":58426}}},"doAtValue":{"const":{"val":"-1s"}}}
	],
	"priorityList": [
		{"action":{"autocastOtherCooldowns":{}}},
		{"action":{"condition":{"not":{"val":{"spellIsReady":{"spellId":{"spellId":57934}}}}},"castSpell":{"spellId":{"spellId":57934}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"80"}}}},"castSpell":{"spellId":{"spellId":13750}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"65"}}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":58426}}},"rhs":{"const":{"val":"1s"}}}}]}},"castSpell":{"spellId":{"spellId":26889}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"50"}}}},"castSpell":{"spellId":{"itemId":7676}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentEnergy":{}},"rhs":{"const":{"val":"50"}}}},"castSpell":{"spellId":{"spellId":51690}}}},
		{"action":{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":16551}}}}},"castSpell":{"spellId":{"spellId":14177}}}},
		{"action":{"castSpell":{"spellId":{"spellId":51723}}}}
	]
}`)

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{"id":40499,"enchant":3817,"gems":[41398,42702]},
	{"id":44664,"gems":[40003]},
	{"id":40502,"enchant":3808,"gems":[40003]},
	{"id":40403,"enchant":3605},
	{"id":40539,"enchant":3832,"gems":[40003]},
	{"id":39765,"enchant":3845,"gems":[40003,0]},
	{"id":40496,"enchant":3604,"gems":[40053,0]},
	{"id":40260,"gems":[39999]},
	{"id":40500,"enchant":3823,"gems":[40003,40003]},
	{"id":39701,"enchant":3606},
	{"id":40074},
	{"id":40474},
	{"id":40684},
	{"id":44253},
	{"id":39714,"enchant":3789},
	{"id":40386,"enchant":3789},
	{"id":40385}
  ]}`)
var GearWithoutRED = core.EquipmentSpecFromJsonString(`{"items": [
	{
	  "id": 37293,
	  "enchant": 3817,
	  "gems": [
		41339,
		40088
	  ]
	},
	{
	  "id": 37861
	},
	{
	  "id": 37139,
	  "enchant": 3808,
	  "gems": [
		36766
	  ]
	},
	{
	  "id": 36947,
	  "enchant": 3605
	},
	{
	  "id": 37165,
	  "enchant": 3832,
	  "gems": [
		40044,
		36766
	  ]
	},
	{
	  "id": 44203,
	  "enchant": 3845,
	  "gems": [
		0
	  ]
	},
	{
	  "id": 37409,
	  "enchant": 3604,
	  "gems": [
		0
	  ]
	},
	{
	  "id": 37194,
	  "gems": [
		40014,
		40157
	  ]
	},
	{
	  "id": 37644,
	  "enchant": 3823
	},
	{
	  "id": 44297,
	  "enchant": 3606
	},
	{
	  "id": 43251,
	  "gems": [
		40136
	  ]
	},
	{
	  "id": 37642
	},
	{
	  "id": 37390
	},
	{
	  "id": 37166
	},
	{
	  "id": 37693,
	  "enchant": 3789
	},
	{
	  "id": 37856,
	  "enchant": 3789
	},
	{
	  "id": 37191
	}
  ]}`)
var GearWithRED = core.EquipmentSpecFromJsonString(`{"items": [
	{
	  "id": 37293,
	  "enchant": 3817,
	  "gems": [
		41398,
		40088
	  ]
	},
	{
	  "id": 37861
	},
	{
	  "id": 37139,
	  "enchant": 3808,
	  "gems": [
		36766
	  ]
	},
	{
	  "id": 36947,
	  "enchant": 3605
	},
	{
	  "id": 37165,
	  "enchant": 3832,
	  "gems": [
		40044,
		36766
	  ]
	},
	{
	  "id": 44203,
	  "enchant": 3845,
	  "gems": [
		0
	  ]
	},
	{
	  "id": 37409,
	  "enchant": 3604,
	  "gems": [
		0
	  ]
	},
	{
	  "id": 37194,
	  "gems": [
		40014,
		40157
	  ]
	},
	{
	  "id": 37644,
	  "enchant": 3823
	},
	{
	  "id": 44297,
	  "enchant": 3606
	},
	{
	  "id": 43251,
	  "gems": [
		40136
	  ]
	},
	{
	  "id": 37642
	},
	{
	  "id": 37390
	},
	{
	  "id": 37166
	},
	{
	  "id": 37693,
	  "enchant": 3789
	},
	{
	  "id": 37856,
	  "enchant": 3789
	},
	{
	  "id": 37191
	}
  ]}`)
var SubtletyP2Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{"id":46125,"enchant":3817,"gems":[41398,42143]},
	{"id":45517,"gems":[49110]},
	{"id":45245,"enchant":3808,"gems":[40023,40003]},
	{"id":45461,"enchant":3605,"gems":[40044]},
	{"id":45473,"enchant":3832,"gems":[40044,40023,40003]},
	{"id":45611,"enchant":3845,"gems":[40044,0]},
	{"id":46124,"enchant":3604,"gems":[39997,0]},
	{"id":46095,"enchant":3599,"gems":[42143,42143,39997]},
	{"id":45536,"enchant":3823,"gems":[40044,39997,40023]},
	{"id":45564,"enchant":3606,"gems":[40023,40003]},
	{"id":45608,"gems":[39997]},
	{"id":46048,"gems":[39997]},
	{"id":45609},
	{"id":45931},
	{"id":45132,"enchant":3789,"gems":[40044]},
	{"id":45484,"enchant":3789,"gems":[39997]},
	{"id":45296,"gems":[39997]}
]}`)
