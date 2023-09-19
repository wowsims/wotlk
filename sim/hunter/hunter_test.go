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
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
		Rotation:    core.RotationCombo{Label: "BM", Rotation: BMRotation},

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
		Rotation:    core.RotationCombo{Label: "MM", Rotation: MMRotation},

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
		Rotation:    core.RotationCombo{Label: "SV", Rotation: SVRotation},
		OtherRotations: []core.RotationCombo{
			{Label: "AOE", Rotation: AOERotation},
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

var AOERotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"prepullActions": [
		{"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAtValue":{"const":{"val":"-1s"}}}
	],
	"priorityList": [
		{"action":{"condition":{"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"10s"}}}},"autocastOtherCooldowns":{}}},
		{"action":{"condition":{"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":34074}}}}},{"cmp":{"op":"OpLt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"10%"}}}}]}},"castSpell":{"spellId":{"spellId":34074}}}},
		{"action":{"condition":{"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":61847}}}}},{"cmp":{"op":"OpGt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"30%"}}}}]}},"castSpell":{"spellId":{"spellId":61847}}}},
		{"hide":true,"action":{"multidot":{"spellId":{"spellId":49001},"maxDots":3,"maxOverlap":{"const":{"val":"0ms"}}}}},
		{"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":49067}}}}},"castSpell":{"spellId":{"tag":1,"spellId":49067}}}},
		{"action":{"castSpell":{"spellId":{"spellId":58434}}}}
	]
}`)

var BMRotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"prepullActions": [
		{"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAtValue":{"const":{"val":"-1s"}}}
	],
	"priorityList": [
		{"action":{"condition":{"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"10s"}}}},"autocastOtherCooldowns":{}}},
		{"action":{"condition":{"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":34074}}}}},{"cmp":{"op":"OpLt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"10%"}}}}]}},"castSpell":{"spellId":{"spellId":34074}}}},
		{"action":{"condition":{"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":61847}}}}},{"cmp":{"op":"OpGt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"30%"}}}}]}},"castSpell":{"spellId":{"spellId":61847}}}},
		{"action":{"castSpell":{"spellId":{"spellId":61006}}}},
		{"hide":true,"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":49067}}}}},"castSpell":{"spellId":{"tag":1,"spellId":49067}}}},
		{"action":{"condition":{"and":{"vals":[{"not":{"val":{"dotIsActive":{"spellId":{"spellId":49001}}}}},{"cmp":{"op":"OpGt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"6s"}}}}]}},"castSpell":{"spellId":{"spellId":49001}}}},
		{"action":{"castSpell":{"spellId":{"spellId":49050}}}},
		{"action":{"castSpell":{"spellId":{"spellId":49048}}}},
		{"hide":true,"action":{"castSpell":{"spellId":{"spellId":49045}}}},
		{"action":{"castSpell":{"spellId":{"spellId":49052}}}}
	]
}`)

var MMRotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"prepullActions": [
		{"action":{"castSpell":{"spellId":{"spellId":61847}}},"doAtValue":{"const":{"val":"-25s"}}},
		{"action":{"castSpell":{"spellId":{"spellId":49067}}},"doAtValue":{"const":{"val":"-20s"}}},
		{"action":{"castSpell":{"spellId":{"spellId":53517}}},"doAtValue":{"const":{"val":"-3s"}}},
		{"action":{"castSpell":{"spellId":{"itemId":40211}}},"doAtValue":{"const":{"val":"-1.401s"}}},
		{"action":{"castSpell":{"spellId":{"spellId":49052}}},"doAtValue":{"const":{"val":"-1.4s"}}}
	],
	"priorityList": [
		{"action":{"condition":{"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"1.35s"}}}},"autocastOtherCooldowns":{}}},
		{"action":{"condition":{"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"1s"}}}},"castSpell":{"spellId":{"itemId":42641}}}},
		{"action":{"condition":{"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"61s"}}}},"castSpell":{"spellId":{"itemId":41119}}}},
		{"action":{"condition":{"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"10s"}}}},"castSpell":{"spellId":{"spellId":34026}}}},
		{"action":{"condition":{"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"1.35s"}}}},"castSpell":{"spellId":{"spellId":34490}}}},
		{"action":{"condition":{"cmp":{"op":"OpEq","lhs":{"currentTime":{}},"rhs":{"const":{"val":"0s"}}}},"castSpell":{"spellId":{"itemId":41119}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":53209}}},"rhs":{"const":{"val":"6s"}}}},{"or":{"vals":[{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49050}}},"rhs":{"const":{"val":"6s"}}}},{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49048}}},"rhs":{"const":{"val":"6s"}}}}]}},{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":3045}}},"rhs":{"const":{"val":"167s"}}}},{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":34490}}},"rhs":{"const":{"val":"13s"}}}},{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49067}}},"rhs":{"const":{"val":"23s"}}}}]}},"castSpell":{"spellId":{"spellId":23989}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"15%"}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":61847}}}}}]}},"castSpell":{"spellId":{"spellId":61847}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"10%"}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":34074}}}}}]}},"castSpell":{"spellId":{"spellId":34074}}}},
		{"action":{"condition":{"cmp":{"op":"OpLt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"40%"}}}},"castSpell":{"spellId":{"itemId":20520}}}},
		{"action":{"condition":{"cmp":{"op":"OpLt","lhs":{"spellTimeToReady":{"spellId":{"spellId":61006}}},"rhs":{"const":{"val":"0.21s"}}}},"castSpell":{"spellId":{"spellId":61006}}}},
		{"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":49001}}}}},"castSpell":{"spellId":{"spellId":49001}}}},
		{"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":49067}}}}},"castSpell":{"spellId":{"spellId":49067}}}},
		{"action":{"condition":{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"spellTimeToReady":{"spellId":{"spellId":53209}}},"rhs":{"const":{"val":"0.15s"}}}},{"spellCanCast":{"spellId":{"spellId":53209}}}]}},"castSpell":{"spellId":{"spellId":53209}}}},
		{"action":{"condition":{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49050}}},"rhs":{"const":{"val":"0.15s"}}}},{"spellCanCast":{"spellId":{"spellId":49050}}}]}},"castSpell":{"spellId":{"spellId":49050}}}},
		{"action":{"castSpell":{"spellId":{"spellId":49052}}}},
		{"hide":true,"action":{"condition":{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49045}}},"rhs":{"const":{"val":"0.2s"}}}},{"spellCanCast":{"spellId":{"spellId":49045}}}]}},"castSpell":{"spellId":{"spellId":49045}}}}
	]
}`)

var SVRotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"prepullActions": [
		{"action":{"castSpell":{"spellId":{"spellId":49067}}},"doAtValue":{"const":{"val":"-24s"}}},
		{"action":{"castSpell":{"spellId":{"spellId":61847}}},"doAtValue":{"const":{"val":"-20s"}}},
		{"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAtValue":{"const":{"val":"-1.4s"}}},
		{"action":{"castSpell":{"spellId":{"spellId":49052}}},"doAtValue":{"const":{"val":"-1.4s"}}}
	],
	"priorityList": [
		{"action":{"condition":{"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"1s"}}}},"autocastOtherCooldowns":{}}},
		{"action":{"condition":{"cmp":{"op":"OpLt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"40%"}}}},"castSpell":{"spellId":{"itemId":20520}}}},
		{"action":{"condition":{"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":61847}}}}},{"cmp":{"op":"OpGt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"30%"}}}}]}},"castSpell":{"spellId":{"spellId":61847}}}},
		{"action":{"condition":{"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":34074}}}}},{"cmp":{"op":"OpLt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"10%"}}}}]}},"castSpell":{"spellId":{"spellId":34074}}}},
		{"action":{"castSpell":{"spellId":{"spellId":61006}}}},
		{"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":60053}}}}},"castSpell":{"spellId":{"spellId":60053}}}},
		{"action":{"condition":{"dotIsActive":{"spellId":{"spellId":60053}}},"castSpell":{"spellId":{"spellId":60052}}}},
		{"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":49067}}}}},"castSpell":{"spellId":{"spellId":49067,"tag":1}}}},
		{"action":{"condition":{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"spellTimeToReady":{"spellId":{"spellId":63672}}},"rhs":{"const":{"val":"0.2s"}}}},{"spellCanCast":{"spellId":{"spellId":63672}}}]}},"castSpell":{"spellId":{"spellId":63672}}}},
		{"action":{"condition":{"and":{"vals":[{"not":{"val":{"dotIsActive":{"spellId":{"spellId":49001}}}}},{"cmp":{"op":"OpGt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"6s"}}}}]}},"castSpell":{"spellId":{"spellId":49001}}}},
		{"action":{"condition":{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"spellTimeToReady":{"spellId":{"spellId":49048}}},"rhs":{"const":{"val":"0.2s"}}}},{"spellCanCast":{"spellId":{"spellId":49048}}}]}},"castSpell":{"spellId":{"spellId":49048}}}},
		{"hide":true,"action":{"castSpell":{"spellId":{"spellId":49048}}}},
		{"action":{"castSpell":{"spellId":{"spellId":49052}}}}
	]
}`)

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
