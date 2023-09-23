package balance

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get caster sets included. (we use spellfire here)
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterBalanceDruid()
}

func TestBalance(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDruid,
		Race:  proto.Race_RaceTauren,

		GearSet: core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		OtherGearSets: []core.GearSetCombo{
			{Label: "P2", GearSet: P2Gear},
			{Label: "P2-4P", GearSet: P2Gear4P},
			{Label: "P3", GearSet: P3Gear},
		},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsAdaptive},
		Rotation:    core.RotationCombo{Label: "Default", Rotation: DefaultRotation},

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
				proto.WeaponType_WeaponTypePolearm,
			},
			ArmorType: proto.ArmorType_ArmorTypeLeather,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeIdol,
			},
		},
	}))
}

func TestBalancePhase3(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDruid,
		Race:  proto.Race_RaceTauren,

		GearSet: core.GearSetCombo{Label: "P3", GearSet: P3Gear},
		Talents: "5102233115331303213305311031--205003002",
		Glyphs: &proto.Glyphs{
			Major1: int32(proto.DruidMajorGlyph_GlyphOfStarfire),
			Major2: int32(proto.DruidMajorGlyph_GlyphOfMoonfire),
			Major3: int32(proto.DruidMajorGlyph_GlyphOfStarfall),
		},
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsAdaptive},
		Rotation:    core.RotationCombo{Label: "Default", Rotation: DefaultRotation},

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
				proto.WeaponType_WeaponTypePolearm,
			},
			ArmorType: proto.ArmorType_ArmorTypeLeather,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeIdol,
			},
		},
	}))
}

var StandardTalents = "5012203115331303213305311231--205003012"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.DruidMajorGlyph_GlyphOfStarfire),
	Major2: int32(proto.DruidMajorGlyph_GlyphOfInsectSwarm),
	Major3: int32(proto.DruidMajorGlyph_GlyphOfStarfall),
	Minor1: int32(proto.DruidMinorGlyph_GlyphOfTyphoon),
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	Food:          proto.Food_FoodFishFeast,
	DefaultPotion: proto.Potions_PotionOfSpeed,
	PrepopPotion:  proto.Potions_PotionOfWildMagic,
}

var PlayerOptionsAdaptive = &proto.Player_BalanceDruid{
	BalanceDruid: &proto.BalanceDruid{
		Options:  &proto.BalanceDruid_Options{},
		Rotation: &proto.BalanceDruid_Rotation{},
	},
}

var DefaultRotation = core.APLRotationFromJsonString(`{
  "type": "TypeAPL",
  "prepullActions": [
    {"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAtValue":{"const":{"val":"-1.5s"}}},
    {"action":{"castSpell":{"spellId":{"spellId":48461}}},"doAtValue":{"const":{"val":"-1.5s"}}}
  ],
  "priorityList": [
    {"action":{"condition":{"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"5"}}}},"castSpell":{"spellId":{"tag":-1,"spellId":2825}}}},
    {"action":{"castSpell":{"spellId":{"itemId":41119}}}},
    {"action":{"multidot":{"spellId":{"spellId":48463},"maxDots":1,"maxOverlap":{"const":{"val":"0ms"}}}}},
    {"action":{"castSpell":{"spellId":{"spellId":53201}}}},
    {"action":{"castSpell":{"spellId":{"spellId":65861}}}},
    {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"auraRemainingTime":{"sourceUnit":{},"auraId":{"spellId":48518}}},"rhs":{"const":{"val":"10s"}}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{},"auraId":{"spellId":48518}}},"rhs":{"const":{"val":"14.8"}}}}]}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"12s"}}}}]}},"castSpell":{"spellId":{"spellId":54758}}}},
    {"action":{"condition":{"or":{"vals":[{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"auraRemainingTime":{"sourceUnit":{},"auraId":{"spellId":48518}}},"rhs":{"const":{"val":"10s"}}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{},"auraId":{"spellId":48518}}},"rhs":{"const":{"val":"14.8"}}}}]}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"15s"}}}}]}},"castSpell":{"spellId":{"itemId":40211}}}},
    {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"sourceUnit":{},"auraId":{"spellId":48518}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{},"auraId":{"spellId":48518}}},"rhs":{"const":{"val":"14.8s"}}}}]}},"castSpell":{"spellId":{"spellId":48465}}}},
    {"action":{"condition":{"and":{"vals":[{"auraIsActive":{"sourceUnit":{},"auraId":{"spellId":48517}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{},"auraId":{"spellId":48517}}},"rhs":{"const":{"val":"14.8s"}}}}]}},"castSpell":{"spellId":{"spellId":48461}}}},
    {"action":{"condition":{"and":{"vals":[{"not":{"val":{"dotIsActive":{"spellId":{"spellId":48468}}}}},{"auraInternalCooldown":{"auraId":{"spellId":48518}}}]}},"castSpell":{"spellId":{"spellId":48468}}}},
    {"action":{"condition":{"auraInternalCooldown":{"sourceUnit":{},"auraId":{"spellId":48518}}},"castSpell":{"spellId":{"spellId":48465}}}},
    {"action":{"castSpell":{"spellId":{"spellId":48461}}}}
  ]
}`)

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
  {"id":40467,"enchant":3820,"gems":[41285,42144]},
  {"id":44661,"gems":[40026]},
  {"id":40470,"enchant":3810,"gems":[42144]},
  {"id":44005,"enchant":3859,"gems":[40026]},
  {"id":40469,"enchant":3832,"gems":[42144,39998]},
  {"id":44008,"enchant":2332,"gems":[39998,0]},
  {"id":40466,"enchant":3604,"gems":[39998,0]},
  {"id":40561,"enchant":3601,"gems":[39998]},
  {"id":40560,"enchant":3719},
  {"id":40519,"enchant":3606},
  {"id":40399},
  {"id":40080},
  {"id":40255},
  {"id":40432},
  {"id":40395,"enchant":3834},
  {"id":40192},
  {"id":40321}
]}`)

var P2Gear = core.EquipmentSpecFromJsonString(`{"items": [
  {"id":45497,"enchant":3820,"gems":[41285,42144]},
  {"id":45133,"gems":[40048]},
  {"id":46196,"enchant":3810,"gems":[39998]},
  {"id":45242,"enchant":3859,"gems":[40048]},
  {"id":45519,"enchant":3832,"gems":[40051,42144,40026]},
  {"id":45446,"enchant":2332,"gems":[42144,0]},
  {"id":45665,"enchant":3604,"gems":[39998,39998,0]},
  {"id":45619,"gems":[39998,39998,39998]},
  {"id":46192,"enchant":3719,"gems":[39998,39998]},
  {"id":45537,"enchant":3606,"gems":[39998,40026]},
  {"id":46046,"gems":[39998]},
  {"id":45495,"gems":[39998]},
  {"id":45466},
  {"id":45518},
  {"id":45620,"enchant":3834,"gems":[39998]},
  {"id":45617},
  {"id":40321}
]}`)

var P2Gear4P = core.EquipmentSpecFromJsonString(`{"items": [
    {
      "id": 46191,
      "enchant": 3820,
      "gems": [
        41285,
        42144
      ]
    },
    {
      "id": 45933,
      "gems": [
        39998
      ]
    },
    {
      "id": 46196,
      "enchant": 3810,
      "gems": [
        40026
      ]
    },
    {
      "id": 45242,
      "enchant": 3859,
      "gems": [
        39998
      ]
    },
    {
      "id": 46194,
      "enchant": 3832,
      "gems": [
        39998,
        42144
      ]
    },
    {
      "id": 45446,
      "enchant": 2332,
      "gems": [
        42144,
        0
      ]
    },
    {
      "id": 45665,
      "enchant": 3604,
      "gems": [
        39998,
        39998,
        0
      ]
    },
    {
      "id": 45616,
      "gems": [
        39998,
        39998,
        39998
      ]
    },
    {
      "id": 46192,
      "enchant": 3719,
      "gems": [
        39998,
        39998
      ]
    },
    {
      "id": 45537,
      "enchant": 3606,
      "gems": [
        39998,
        40026
      ]
    },
    {
      "id": 46046,
      "gems": [
        39998
      ]
    },
    {
      "id": 45495,
      "gems": [
        39998
      ]
    },
    {
      "id": 45466
    },
    {
      "id": 45518
    },
    {
      "id": 45620,
      "enchant": 3834,
      "gems": [
        39998
      ]
    },
    {
      "id": 45617
    },
    {
      "id": 40321
    }
]}`)

var P3Gear = core.EquipmentSpecFromJsonString(`{"items": [
  {"id":48171,"enchant":3820,"gems":[41285,40153]},
  {"id":47144,"gems":[40153]},
  {"id":48168,"enchant":3810,"gems":[40153]},
  {"id":47552,"enchant":3722,"gems":[40113]},
  {"id":48169,"enchant":3832,"gems":[40113,40113]},
  {"id":47066,"enchant":2332,"gems":[40113,0]},
  {"id":48172,"enchant":3604,"gems":[40113,0]},
  {"id":47084,"gems":[40133,40113,40113]},
  {"id":47190,"enchant":3719,"gems":[40113,40113,40113]},
  {"id":47097,"enchant":3606,"gems":[40133,40113]},
  {"id":47237,"gems":[40113]},
  {"id":46046,"gems":[40113]},
  {"id":45518},
  {"id":47188},
  {"id":47206,"enchant":3834},
  {"id":47064},
  {"id":47670}
]}`)
