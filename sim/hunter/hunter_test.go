package hunter

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get item effects included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterHunter()
}

func TestBM(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassHunter,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceDwarf},

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     BMTalents,
		Glyphs:      BMGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "BM", SpecOptions: PlayerOptionsBM},

		ItemFilter: ItemFilter,
	}))
}

func TestMM(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassHunter,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceDwarf},

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     MMTalents,
		Glyphs:      MMGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "MM", SpecOptions: PlayerOptionsMM},

		ItemFilter: ItemFilter,
	}))
}

func TestSV(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassHunter,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceDwarf},

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     SVTalents,
		Glyphs:      SVGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "SV", SpecOptions: PlayerOptionsSV},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "AOE", SpecOptions: PlayerOptionsAOE},
		},

		ItemFilter: ItemFilter,
	}))
}

func TestAPL(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassHunter,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceDwarf},

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     SVTalents,
		Glyphs:      SVGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "SV", SpecOptions: PlayerOptionsSV},
		Rotation:    core.RotationCombo{Label: "Default", Rotation: DefaultRotation},

		ItemFilter: ItemFilter,
	}))
}

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypeMail,
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeFist,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeOffHand,
		proto.WeaponType_WeaponTypePolearm,
		proto.WeaponType_WeaponTypeStaff,
		proto.WeaponType_WeaponTypeSword,
	},
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeBow,
		proto.RangedWeaponType_RangedWeaponTypeCrossbow,
		proto.RangedWeaponType_RangedWeaponTypeGun,
	},
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:          proto.Race_RaceOrc,
				Class:         proto.Class_ClassHunter,
				Equipment:     P1Gear,
				Consumes:      FullConsumes,
				Spec:          PlayerOptionsMM,
				Glyphs:        MMGlyphs,
				TalentsString: MMTalents,
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

var BMTalents = "51200201515012233110531351-005305-5"
var MMTalents = "502-035335131030013233035031051-5000002"
var SVTalents = "-015305101-5000032500033330532135301311"
var BMGlyphs = &proto.Glyphs{
	Major1: int32(proto.HunterMajorGlyph_GlyphOfBestialWrath),
	Major2: int32(proto.HunterMajorGlyph_GlyphOfSteadyShot),
	Major3: int32(proto.HunterMajorGlyph_GlyphOfSerpentSting),
}
var MMGlyphs = &proto.Glyphs{
	Major1: int32(proto.HunterMajorGlyph_GlyphOfSerpentSting),
	Major2: int32(proto.HunterMajorGlyph_GlyphOfSteadyShot),
	Major3: int32(proto.HunterMajorGlyph_GlyphOfChimeraShot),
}
var SVGlyphs = &proto.Glyphs{
	Major1: int32(proto.HunterMajorGlyph_GlyphOfSerpentSting),
	Major2: int32(proto.HunterMajorGlyph_GlyphOfExplosiveShot),
	Major3: int32(proto.HunterMajorGlyph_GlyphOfKillShot),
}

var FerocityTalents = &proto.HunterPetTalents{
	CobraReflexes:  2,
	Dive:           true,
	SpikedCollar:   3,
	BoarsSpeed:     true,
	CullingTheHerd: 3,
	SpidersBite:    3,
	Rabid:          true,
	CallOfTheWild:  true,
	WildHunt:       1,
}

var DefaultRotation = core.APLRotationFromJsonString(`{
	"enabled": true,
	"priorityList": [
		{"action": {
			"condition": {"not": {"val": {"dotIsActive": {"spellId": { "spellId": 49001 }}}}},
			"castSpell": {"spellId": { "spellId": 49001 }}
		}},
		{"action": {"castSpell": {"spellId": { "spellId": 61006 }}}},
		{"action": {"castSpell": {"spellId": { "spellId": 63672 }}}},
		{"action": {"castSpell": {"spellId": { "spellId": 60053 }}}},
		{"action": {"castSpell": {"spellId": { "spellId": 49050 }}}},
		{"action": {
			"condition": {"not": {"val": {"dotIsActive": {"spellId": { "spellId": 60053 }}}}},
			"castSpell": {"spellId": { "spellId": 49045 }}
		}},
		{"action": {"castSpell": {"spellId": { "spellId": 49052 }}}}
	]
}`)

var PlayerOptionsMM = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Options:  basicOptions,
		Rotation: basicRotation,
	},
}

var PlayerOptionsBM = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Options:  basicOptions,
		Rotation: basicRotation,
	},
}

var PlayerOptionsSV = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Options:  basicOptions,
		Rotation: basicRotation,
	},
}

var PlayerOptionsAOE = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Options:  basicOptions,
		Rotation: aoeRotation,
	},
}

var basicRotation = &proto.Hunter_Rotation{
	Sting: proto.Hunter_Rotation_SerpentSting,

	ViperStartManaPercent: 0.2,
	ViperStopManaPercent:  0.3,
}

var aoeRotation = &proto.Hunter_Rotation{
	TrapWeave:         true,
	TimeToTrapWeaveMs: 2000,

	ViperStartManaPercent: 0.2,
	ViperStopManaPercent:  0.3,
}

var basicOptions = &proto.Hunter_Options{
	Ammo:       proto.Hunter_Options_SaroniteRazorheads,
	PetType:    proto.Hunter_Options_Wolf,
	PetTalents: FerocityTalents,
	PetUptime:  0.9,

	SniperTrainingUptime: 0.8,
	UseHuntersMark:       true,
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfRelentlessAssault,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredFlameCap,
	PetFood:         proto.PetFood_PetFoodKiblersBits,
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40505,
		"enchant": 3817,
		"gems": [
			41398,
			42143
		]
	},
	{
		"id": 44664,
		"gems": [
			42143
		]
	},
	{
		"id": 40507,
		"enchant": 3808,
		"gems": [
			39997
		]
	},
	{
		"id": 40403,
		"enchant": 3605
	},
	{
		"id": 43998,
		"enchant": 3832,
		"gems": [
			42143,
			39997
		]
	},
	{
		"id": 40282,
		"enchant": 3845,
		"gems": [
			39997,
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
		"id": 39762,
		"enchant": 3601,
		"gems": [
			39997
		]
	},
	{
		"id": 40331,
		"enchant": 3823,
		"gems": [
			39997,
			49110
		]
	},
	{
		"id": 40549,
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
		"id": 40388,
		"enchant": 3827
	},
	{},
	{
		"id": 40385,
		"enchant": 3608
	}
]}`)
