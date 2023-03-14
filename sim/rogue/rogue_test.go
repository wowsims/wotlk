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
	env, _ := core.NewEnvironment(raid, encounter)
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
		Rotation: basicRotation,
	},
}
var PlayerOptionsCombatDD = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyDeadly,
		Rotation: basicRotation,
	},
}
var PlayerOptionsCombatID = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  InstantDeadly,
		Rotation: basicRotation,
	},
}
var PlayerOptionsCombatII = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  InstantInstant,
		Rotation: basicRotation,
	},
}

var PlayerOptionsNoLethality = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyInstant,
		Rotation: basicRotation,
	},
}

var PlayerOptionsNoPotW = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyInstant,
		Rotation: basicRotation,
	},
}

var PlayerOptionsNoLethalityNoPotW = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyInstant,
		Rotation: basicRotation,
	},
}

var PlayerOptionsAssassinationDI = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyInstant,
		Rotation: basicRotation,
	},
}
var PlayerOptionsAssassinationDD = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  DeadlyDeadly,
		Rotation: basicRotation,
	},
}
var PlayerOptionsAssassinationID = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  InstantDeadly,
		Rotation: basicRotation,
	},
}
var PlayerOptionsAssassinationII = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  InstantInstant,
		Rotation: basicRotation,
	},
}

var PlayerOptionsSubtletyID = &proto.Player_Rogue{
	Rogue: &proto.Rogue{
		Options:  InstantDeadly,
		Rotation: subtletyRotation,
	},
}

var basicRotation = &proto.Rogue_Rotation{
	ExposeArmorFrequency:                proto.Rogue_Rotation_Never,
	TricksOfTheTradeFrequency:           proto.Rogue_Rotation_Maintain,
	AssassinationFinisherPriority:       proto.Rogue_Rotation_EnvenomRupture,
	CombatFinisherPriority:              proto.Rogue_Rotation_RuptureEviscerate,
	MinimumComboPointsExposeArmor:       4,
	MinimumComboPointsPrimaryFinisher:   3,
	MinimumComboPointsSecondaryFinisher: 2,
	MultiTargetSliceFrequency:           proto.Rogue_Rotation_Once,
	MinimumComboPointsMultiTargetSlice:  4,
}

var subtletyRotation = &proto.Rogue_Rotation{
	ExposeArmorFrequency:                proto.Rogue_Rotation_Never,
	TricksOfTheTradeFrequency:           proto.Rogue_Rotation_Never,
	SubtletyFinisherPriority:            proto.Rogue_Rotation_SubtletyEviscerate,
	MinimumComboPointsPrimaryFinisher:   5,
	MinimumComboPointsSecondaryFinisher: 5,
	OpenWithGarrote:                     true,
	OpenWithPremeditation:               true,
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

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
    {
      "id": 40499,
      "enchant": 3817,
      "gems": [
        41398,
        42702
      ]
    },
    {
      "id": 44664,
      "gems": [
        42154
      ]
    },
    {
      "id": 40502,
      "enchant": 3808,
      "gems": [
        36766
      ]
    },
    {
      "id": 40403,
      "enchant": 3605
    },
    {
      "id": 40539,
      "enchant": 3832,
      "gems": [
        36766
      ]
    },
    {
      "id": 39765,
      "enchant": 3845,
      "gems": [
        40003,
        0
      ]
    },
    {
      "id": 40496,
      "enchant": 3604,
      "gems": [
        40058,
        0
      ]
    },
    {
      "id": 40260,
      "gems": [
        39999
      ]
    },
    {
      "id": 40500,
      "enchant": 3823,
      "gems": [
        40003,
        40003
      ]
    },
    {
      "id": 39701,
      "enchant": 3606
    },
    {
      "id": 40074
    },
    {
      "id": 40474
    },
    {
        "id": 40684
    },
    {
      "id": 44253
    },
    {
      "id": 39714,
      "enchant": 3789
    },
    {
      "id": 40386,
      "enchant": 3789
    },
    {
      "id": 40385
    }
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
var MutilateP1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 29044,
		"enchant": 3003,
		"gems": [
			32409,
			24061
		]
	},
	{
		"id": 29381
	},
	{
		"id": 27797,
		"enchant": 2986,
		"gems": [
			24061,
			24055
		]
	},
	{
		"id": 28672,
		"enchant": 368
	},
	{
		"id": 29045,
		"enchant": 2661,
		"gems": [
			24061,
			24051,
			24055
		]
	},
	{
		"id": 29246,
		"enchant": 1593
	},
	{
		"id": 27531,
		"gems": [
			24061,
			24061
		]
	},
	{
		"id": 29247
	},
	{
		"id": 28741,
		"enchant": 3012,
		"gems": [
			24051,
			24051,
			24051
		]
	},
	{
		"id": 28545,
		"enchant": 2939,
		"gems": [
			24061,
			24051
		]
	},
	{
		"id": 28757
	},
	{
		"id": 28649
	},
	{
		"id": 29383
	},
	{
		"id": 28830
	},
	{
		"id": 28768,
		"enchant": 2673
	},
	{
		"id": 29182,
		"enchant": 2673
	},
	{
		"id": 28772
	}
]}`)
var SubtletyP2Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
	  "id": 46125,
	  "enchant": 3817,
	  "gems": [
		41398,
		42143
	  ]
	},
	{
	  "id": 45517,
	  "gems": [
		49110
	  ]
	},
	{
	  "id": 46127,
	  "enchant": 3808,
	  "gems": [
		39997
	  ]
	},
	{
	  "id": 45461,
	  "enchant": 3605,
	  "gems": [
		40044
	  ]
	},
	{
	  "id": 46123,
	  "enchant": 3832,
	  "gems": [
		39997,
		40044
	  ]
	},
	{
	  "id": 45611,
	  "enchant": 3845,
	  "gems": [
		40044,
		0
	  ]
	},
	{
	  "id": 46124,
	  "enchant": 3604,
	  "gems": [
		39997,
		0
	  ]
	},
	{
	  "id": 46095,
	  "enchant": 3599,
	  "gems": [
		42143,
		42143,
		39997
	  ]
	},
	{
	  "id": 45536,
	  "enchant": 3823,
	  "gems": [
		40044,
		39997,
		40023
	  ]
	},
	{
	  "id": 45564,
	  "enchant": 3606,
	  "gems": [
		40023,
		40003
	  ]
	},
	{
	  "id": 45608,
	  "gems": [
		39997
	  ]
	},
	{
	  "id": 46048,
	  "gems": [
		39997
	  ]
	},
	{
	  "id": 45609
	},
	{
	  "id": 45931
	},
	{
	  "id": 45132,
	  "enchant": 3789,
	  "gems": [
		40044
	  ]
	},
	{
	  "id": 45484,
	  "enchant": 3789,
	  "gems": [
		39997
	  ]
	},
	{
	  "id": 45296,
	  "gems": [
		39997
	  ]
	}
]}`)
