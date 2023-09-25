package hunter

import (
	"log"
	"os"
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get item effects included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterHunter()
}

func GetAplRotation(dir string, file string) core.RotationCombo {
	filePath := dir + "/" + file + ".json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to load apl json file: %s, %s", filePath, err)
	}

	return core.RotationCombo{Label: file, Rotation: core.APLRotationFromJsonString(string(data))}
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
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
		Rotation:    GetAplRotation("../../ui/hunter/apls", "bm"),

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
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
		Rotation:    GetAplRotation("../../ui/hunter/apls", "mm"),
		OtherRotations: []core.RotationCombo{
			GetAplRotation("../../ui/hunter/apls", "mm_advanced"),
		},

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
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
		Rotation:    GetAplRotation("../../ui/hunter/apls", "sv"),
		OtherRotations: []core.RotationCombo{
			GetAplRotation("../../ui/hunter/apls", "sv_advanced"),
			GetAplRotation("../../ui/hunter/apls", "aoe"),
		},

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
				Spec:          PlayerOptionsBasic,
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

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfRelentlessAssault,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredFlameCap,
	PetFood:         proto.PetFood_PetFoodKiblersBits,
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

var PlayerOptionsBasic = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Options: &proto.Hunter_Options{
			Ammo:       proto.Hunter_Options_SaroniteRazorheads,
			PetType:    proto.Hunter_Options_Wolf,
			PetTalents: FerocityTalents,
			PetUptime:  0.9,

			TimeToTrapWeaveMs:    2000,
			SniperTrainingUptime: 0.8,
			UseHuntersMark:       true,
		},
	},
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{"id":40505,"enchant":3817,"gems":[41398,42143]},
	{"id":44664,"gems":[42143]},
	{"id":40507,"enchant":3808,"gems":[39997]},
	{"id":40403,"enchant":3605},
	{"id":43998,"enchant":3832,"gems":[42143,39997]},
	{"id":40282,"enchant":3845,"gems":[39997,0]},
	{"id":40541,"enchant":3604,"gems":[0]},
	{"id":39762,"enchant":3601,"gems":[39997]},
	{"id":40331,"enchant":3823,"gems":[39997,49110]},
	{"id":40549,"enchant":3606},
	{"id":40074},
	{"id":40474},
	{"id":40684},
	{"id":44253},
	{"id":40388,"enchant":3827},
	{},
	{"id":40385,"enchant":3608}
]}`)
