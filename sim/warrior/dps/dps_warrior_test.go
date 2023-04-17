package dps

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get item effects included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterDpsWarrior()
}

func TestFury(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassWarrior,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		Talents:     FuryTalents,
		Glyphs:      FuryGlyphs,
		GearSet:     core.GearSetCombo{Label: "Fury P1", GearSet: FuryP1Gear},
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFury},

		ItemFilter: core.ItemFilter{
			ArmorType: proto.ArmorType_ArmorTypePlate,

			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeFist,
			},
		},
	}))
}

func TestArms(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassWarrior,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		Talents:     ArmsTalents,
		Glyphs:      ArmsGlyphs,
		GearSet:     core.GearSetCombo{Label: "Arms P1", GearSet: FuryP1Gear},
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsArms},

		ItemFilter: core.ItemFilter{
			ArmorType: proto.ArmorType_ArmorTypePlate,

			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeFist,
			},
		},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:          proto.Race_RaceOrc,
				Class:         proto.Class_ClassWarrior,
				Equipment:     FuryP1Gear,
				Consumes:      FullConsumes,
				Spec:          PlayerOptionsFury,
				TalentsString: FuryTalents,
				Buffs:         core.FullIndividualBuffs,
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

var FuryTalents = "302023102331-305053000520310053120500351"
var FuryGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarriorMajorGlyph_GlyphOfWhirlwind),
	Major2: int32(proto.WarriorMajorGlyph_GlyphOfHeroicStrike),
	Major3: int32(proto.WarriorMajorGlyph_GlyphOfRending),
}
var ArmsTalents = "3022032023335100102012213231251-305-2033"
var ArmsGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarriorMajorGlyph_GlyphOfRending),
	Major2: int32(proto.WarriorMajorGlyph_GlyphOfMortalStrike),
	Major3: int32(proto.WarriorMajorGlyph_GlyphOfExecution),
}

var PlayerOptionsArms = &proto.Player_Warrior{
	Warrior: &proto.Warrior{
		Options:  warriorOptions,
		Rotation: armsRotation,
	},
}

var PlayerOptionsFury = &proto.Player_Warrior{
	Warrior: &proto.Warrior{
		Options:  warriorOptions,
		Rotation: furyRotation,
	},
}

var armsRotation = &proto.Warrior_Rotation{
	UseRend:   true,
	UseMs:     true,
	UseCleave: false,

	HsRageThreshold:          50,
	MsRageThreshold:          35,
	SlamRageThreshold:        25,
	RendCdThreshold:          0,
	RendHealthThresholdAbove: 100,

	SpamExecute: false,

	UseHsDuringExecute: true,

	MaintainDemoShout:   false,
	MaintainThunderClap: false,

	StanceOption: proto.Warrior_Rotation_DefaultStance,
}

var furyRotation = &proto.Warrior_Rotation{
	UseRend:               true,
	UseCleave:             false,
	UseOverpower:          true,
	ExecutePhaseOverpower: false,

	HsRageThreshold:          30,
	RendRageThresholdBelow:   100,
	SlamRageThreshold:        25,
	RendCdThreshold:          0,
	RendHealthThresholdAbove: 20,

	UseHsDuringExecute: true,
	UseWwDuringExecute: true,
	UseBtDuringExecute: true,
	UseSlamOverExecute: true,

	MaintainDemoShout:   false,
	MaintainThunderClap: false,

	StanceOption: proto.Warrior_Rotation_DefaultStance,
}

var warriorOptions = &proto.Warrior_Options{
	StartingRage:       50,
	UseRecklessness:    true,
	UseShatteringThrow: true,
	Shout:              proto.WarriorShout_WarriorShoutBattle,
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfEndlessRage,
	DefaultPotion: proto.Potions_PotionOfSpeed,
	PrepopPotion:  proto.Potions_PotionOfSpeed,
	Food:          proto.Food_FoodFishFeast,
}

var FuryP1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 44006,
		"enchant": 3817,
		"gems": [
			41285,
			42702
		]
	},
	{
		"id": 44664,
		"gems": [
			39996
		]
	},
	{
		"id": 40530,
		"enchant": 3808,
		"gems": [
			40037
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
			42142
		]
	},
	{
		"id": 39765,
		"enchant": 3845,
		"gems": [
			39996,
			0
		]
	},
	{
		"id": 40541,
		"enchant": 3604,
		"gems": [
			0
		]
	},
	{
		"id": 40205,
		"gems": [
			42142
		]
	},
	{
		"id": 40529,
		"enchant": 3823,
		"gems": [
			39996,
			40022
		]
	},
	{
		"id": 40591,
		"enchant": 3606
	},
	{
		"id": 43993,
		"gems": [
			42142
		]
	},
	{
		"id": 40717
	},
	{
		"id": 42987
	},
	{
		"id": 40256
	},
	{
		"id": 40384,
		"enchant": 3789
	},
	{
		"id": 40384,
		"enchant": 3789
	},
	{
		"id": 40385
	}
]}`)
