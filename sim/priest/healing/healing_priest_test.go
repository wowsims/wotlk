package healing

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get caster sets included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterHealingPriest()
}

func TestDisc(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:    proto.Class_ClassPriest,
		Race:     proto.Race_RaceUndead,
		IsHealer: true,

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     DiscTalents,
		Glyphs:      DiscGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Disc", SpecOptions: PlayerOptionsDisc},
		Rotation:    core.RotationCombo{Label: "Disc", Rotation: DiscRotation},

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

		EPReferenceStat: proto.Stat_StatSpellPower,
		StatsToWeigh: []proto.Stat{
			proto.Stat_StatIntellect,
			proto.Stat_StatSpellPower,
			proto.Stat_StatSpellHaste,
			proto.Stat_StatSpellCrit,
		},
	}))
}

func TestHoly(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:    proto.Class_ClassPriest,
		Race:     proto.Race_RaceUndead,
		IsHealer: true,

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     HolyTalents,
		Glyphs:      HolyGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Holy", SpecOptions: PlayerOptionsHoly},
		Rotation:    core.RotationCombo{Label: "Holy", Rotation: HolyRotation},

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

var DiscTalents = "0503203130300512301313231251-2351010303"
var DiscGlyphs = &proto.Glyphs{
	Major1: int32(proto.PriestMajorGlyph_GlyphOfPowerWordShield),
	Major2: int32(proto.PriestMajorGlyph_GlyphOfFlashHeal),
	Major3: int32(proto.PriestMajorGlyph_GlyphOfPenance),
	// No interesting minor glyphs.
}

var HolyTalents = "05032031103-234051032002152530004311051"
var HolyGlyphs = &proto.Glyphs{
	Major1: int32(proto.PriestMajorGlyph_GlyphOfPrayerOfHealing),
	Major2: int32(proto.PriestMajorGlyph_GlyphOfRenew),
	Major3: int32(proto.PriestMajorGlyph_GlyphOfCircleOfHealing),
	// No interesting minor glyphs.
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	Food:          proto.Food_FoodFishFeast,
	DefaultPotion: proto.Potions_RunicManaInjector,
	PrepopPotion:  proto.Potions_PotionOfWildMagic,
}

var PlayerOptionsDisc = &proto.Player_HealingPriest{
	HealingPriest: &proto.HealingPriest{
		Options: &proto.HealingPriest_Options{
			UseInnerFire:      true,
			UseShadowfiend:    true,
			RapturesPerMinute: 5,
		},
		Rotation: &proto.HealingPriest_Rotation{},
	},
}

var PlayerOptionsHoly = &proto.Player_HealingPriest{
	HealingPriest: &proto.HealingPriest{
		Options: &proto.HealingPriest_Options{
			UseInnerFire:   true,
			UseShadowfiend: true,
		},
		Rotation: &proto.HealingPriest_Rotation{},
	},
}

var DiscRotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"priorityList": [
		{"action":{"autocastOtherCooldowns":{}}},
		{"action":{"condition":{"cmp":{"op":"OpLt","lhs":{"spellCpm":{"spellId":{"spellId":48066}}},"rhs":{"const":{"val":"18"}}}},"multishield":{"spellId":{"spellId":48066},"maxShields":10,"maxOverlap":{"const":{"val":"0ms"}}}}},
		{"action":{"condition":{"cmp":{"op":"OpLt","lhs":{"spellCpm":{"spellId":{"spellId":53007}}},"rhs":{"const":{"val":"4"}}}},"castSpell":{"spellId":{"spellId":53007}}}},
		{"action":{"condition":{"cmp":{"op":"OpLt","lhs":{"spellCpm":{"spellId":{"spellId":48113}}},"rhs":{"const":{"val":"2"}}}},"castSpell":{"spellId":{"spellId":48113}}}},
		{"action":{"condition":{"cmp":{"op":"OpLt","lhs":{"spellCpm":{"spellId":{"spellId":48063}}},"rhs":{"const":{"val":"1"}}}},"castSpell":{"spellId":{"spellId":48063}}}}
	]
}`)

var HolyRotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"priorityList": [
		{"action":{"autocastOtherCooldowns":{}}},
		{"action":{"condition":{"cmp":{"op":"OpLt","lhs":{"spellCpm":{"spellId":{"spellId":48063}}},"rhs":{"const":{"val":"10"}}}},"castSpell":{"spellId":{"spellId":48063}}}},
		{"action":{"condition":{"cmp":{"op":"OpLt","lhs":{"spellCpm":{"spellId":{"spellId":48089}}},"rhs":{"const":{"val":"5"}}}},"castSpell":{"spellId":{"spellId":48089}}}},
		{"action":{"condition":{"cmp":{"op":"OpLt","lhs":{"spellCpm":{"spellId":{"spellId":48068}}},"rhs":{"const":{"val":"10"}}}},"multidot":{"spellId":{"spellId":48068},"maxDots":10,"maxOverlap":{"const":{"val":"0ms"}}}}},
		{"action":{"condition":{"cmp":{"op":"OpLt","lhs":{"spellCpm":{"spellId":{"spellId":48113}}},"rhs":{"const":{"val":"2"}}}},"castSpell":{"spellId":{"spellId":48113}}}}
	]
}`)

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{"id":40456,"enchant":3819,"gems":[41401,39998]},
	{"id":44657,"gems":[40047]},
	{"id":40450,"enchant":3809,"gems":[42144]},
	{"id":40724,"enchant":3859},
	{"id":40194,"enchant":3832,"gems":[42144]},
	{"id":40741,"enchant":2332,"gems":[0]},
	{"id":40445,"enchant":3246,"gems":[42144,0]},
	{"id":40271,"enchant":3601,"gems":[40027,39998]},
	{"id":40398,"enchant":3719,"gems":[39998,39998]},
	{"id":40236,"enchant":3606},
	{"id":40108},
	{"id":40433},
	{"id":37835},
	{"id":40258},
	{"id":40395,"enchant":3834},
	{"id":40350},
	{"id":40245}
]}`)
