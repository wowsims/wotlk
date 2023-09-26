package smite

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get caster sets included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterSmitePriest()
}

func TestSmite(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassPriest,
		Race:  proto.Race_RaceUndead,

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     DefaultTalents,
		Glyphs:      DefaultGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
		Rotation:    core.RotationCombo{Label: "Default", Rotation: DefaultRotation},

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeMace,
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

var DefaultTalents = "05332031013005023310001-005551002020152-00502"
var DefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.PriestMajorGlyph_GlyphOfSmite),
	Major2: int32(proto.PriestMajorGlyph_GlyphOfHolyNova),
	Major3: int32(proto.PriestMajorGlyph_GlyphOfShadowWordDeath),
	// No interesting minor glyphs.
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	Food:          proto.Food_FoodFishFeast,
	DefaultPotion: proto.Potions_RunicManaInjector,
	PrepopPotion:  proto.Potions_PotionOfWildMagic,
}

var PlayerOptionsBasic = &proto.Player_SmitePriest{
	SmitePriest: &proto.SmitePriest{
		Options: &proto.SmitePriest_Options{
			UseInnerFire:   true,
			UseShadowfiend: true,
		},
		Rotation: &proto.SmitePriest_Rotation{},
	},
}

var DefaultRotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"prepullActions": [
		{"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAtValue":{"const":{"val":"-1s"}}}
	],
	"priorityList": [
		{"action":{"autocastOtherCooldowns":{}}},
		{"action":{"condition":{"and":{"vals":[{"dotIsActive":{"spellId":{"spellId":48135}}},{"cmp":{"op":"OpLe","lhs":{"spellCastTime":{"spellId":{"spellId":48123}}},"rhs":{"dotRemainingTime":{"spellId":{"spellId":48135}}}}}]}},"castSpell":{"spellId":{"spellId":14751}}}},
		{"action":{"condition":{"and":{"vals":[{"dotIsActive":{"spellId":{"spellId":48135}}},{"cmp":{"op":"OpLe","lhs":{"spellCastTime":{"spellId":{"spellId":48123}}},"rhs":{"dotRemainingTime":{"spellId":{"spellId":48135}}}}}]}},"castSpell":{"spellId":{"spellId":48123}}}},
		{"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":48300}}}}},"castSpell":{"spellId":{"spellId":48300}}}},
		{"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":48125}}}}},"castSpell":{"spellId":{"spellId":48125}}}},
		{"action":{"castSpell":{"spellId":{"spellId":48135}}}},
		{"action":{"condition":{"and":{"vals":[{"not":{"val":{"spellIsReady":{"spellId":{"spellId":48135}}}}},{"cmp":{"op":"OpLe","lhs":{"spellTimeToReady":{"spellId":{"spellId":48135}}},"rhs":{"const":{"val":"50ms"}}}}]}},"wait":{"duration":{"spellTimeToReady":{"spellId":{"spellId":48135}}}}}},
		{"hide":true,"action":{"condition":{"auraIsActive":{"auraId":{"spellId":59000}}},"castSpell":{"spellId":{"spellId":48123}}}},
		{"hide":true,"action":{"castSpell":{"spellId":{"spellId":53007}}}},
		{"hide":true,"action":{"castSpell":{"spellId":{"spellId":48158}}}},
		{"hide":true,"action":{"castSpell":{"spellId":{"spellId":48127}}}},
		{"hide":true,"action":{"castSpell":{"spellId":{"tag":3,"spellId":48156}}}},
		{"action":{"castSpell":{"spellId":{"spellId":48123}}}}
	]
}`)

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{"id":40562,"enchant":3820,"gems":[41333,42144]},
	{"id":44661,"gems":[39998]},
	{"id":40459,"enchant":3810,"gems":[42144]},
	{"id":44005,"enchant":3859,"gems":[42144]},
	{"id":40234,"enchant":1144,"gems":[39998,39998]},
	{"id":44008,"enchant":2332,"gems":[39998,0]},
	{"id":40454,"enchant":3604,"gems":[40049,0]},
	{"id":40561,"enchant":3601,"gems":[39998]},
	{"id":40560,"enchant":3719},
	{"id":40558,"enchant":3826},
	{"id":40719},
	{"id":40399},
	{"id":40255},
	{"id":40432},
	{"id":40395,"enchant":3834},
	{"id":40273},
	{"id":39712}
]}`)
