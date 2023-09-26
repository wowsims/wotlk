package holy

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get item effects included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterHolyPaladin()
}

func TestHoly(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassPaladin,
		Race:       proto.Race_RaceBloodElf,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: BasicOptions},
		Rotation:    core.RotationCombo{Label: "Default", Rotation: DefaultRotation},

		IsHealer:        true,
		InFrontOfTarget: true,

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypePolearm,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeShield,
			},
			ArmorType: proto.ArmorType_ArmorTypePlate,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeLibram,
			},
		},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:          proto.Race_RaceBloodElf,
				Class:         proto.Class_ClassPaladin,
				Equipment:     P1Gear,
				Consumes:      FullConsumes,
				Spec:          BasicOptions,
				TalentsString: StandardTalents,
				Glyphs:        StandardGlyphs,
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

var StandardTalents = "50350151020013053100515221-50023131203"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.PaladinMajorGlyph_GlyphOfHolyLight),
	Major2: int32(proto.PaladinMajorGlyph_GlyphOfSealOfWisdom),
	Major3: int32(proto.PaladinMajorGlyph_GlyphOfBeaconOfLight),
	Minor1: int32(proto.PaladinMinorGlyph_GlyphOfLayOnHands),
	Minor2: int32(proto.PaladinMinorGlyph_GlyphOfSenseUndead),
}

var defaultProtOptions = &proto.HolyPaladin_Options{
	Judgement: proto.PaladinJudgement_JudgementOfWisdom,
	Aura:      proto.PaladinAura_DevotionAura,
}

var BasicOptions = &proto.Player_HolyPaladin{
	HolyPaladin: &proto.HolyPaladin{
		Options:  defaultProtOptions,
		Rotation: &proto.HolyPaladin_Rotation{},
	},
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfStoneblood,
	Food:            proto.Food_FoodDragonfinFilet,
	DefaultPotion:   proto.Potions_IndestructiblePotion,
	PrepopPotion:    proto.Potions_IndestructiblePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var DefaultRotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"priorityList": [
		{"action":{"autocastOtherCooldowns":{}}}
	]
}`)

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{"id":40298,"enchant":3819,"gems":[41401,40012]},
	{"id":44662,"gems":[40012]},
	{"id":40573,"enchant":3809,"gems":[40012]},
	{"id":44005,"enchant":3831,"gems":[40012]},
	{"id":40569,"enchant":3832,"gems":[40012,40012]},
	{"id":40332,"enchant":1119,"gems":[40012,0]},
	{"id":40570,"enchant":3604,"gems":[40012,0]},
	{"id":40259,"gems":[40012]},
	{"id":40572,"enchant":3721,"gems":[40027,40012]},
	{"id":40592,"enchant":3606},
	{"id":40399},
	{"id":40375},
	{"id":44255},
	{"id":37111},
	{"id":40395,"enchant":2666},
	{"id":40401,"enchant":1128},
	{"id":40705}
]}`)
