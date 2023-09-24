package tank

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get item effects included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterTankDeathknight()
}

func TestBloodTank(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathknight,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet:     core.GearSetCombo{Label: "Blood Tank P1", GearSet: BloodP1Gear},
		Talents:     BloodTankTalents,
		Glyphs:      Glyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBloodTank},
		Rotation:    core.RotationCombo{Label: "BloodIT", Rotation: BloodITRotation},

		IsTank:          true,
		InFrontOfTarget: true,

		ItemFilter: core.ItemFilter{
			ArmorType: proto.ArmorType_ArmorTypePlate,

			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeMace,
			},
		},
	}))
}

var BloodTankTalents = "005510153330330220102013-3050505100023101-002"
var Glyphs = &proto.Glyphs{
	Major1: int32(proto.DeathknightMajorGlyph_GlyphOfDarkCommand),
	Major2: int32(proto.DeathknightMajorGlyph_GlyphOfObliterate),
	Major3: int32(proto.DeathknightMajorGlyph_GlyphOfVampiricBlood),
}

var PlayerOptionsBloodTank = &proto.Player_TankDeathknight{
	TankDeathknight: &proto.TankDeathknight{
		Options: &proto.TankDeathknight_Options{
			StartingRunicPower: 0,
		},
		Rotation: &proto.TankDeathknight_Rotation{},
	},
}

var FullRaidBuffs = &proto.RaidBuffs{
	GiftOfTheWild:         proto.TristateEffect_TristateEffectImproved,
	PowerWordFortitude:    proto.TristateEffect_TristateEffectImproved,
	AbominationsMight:     true,
	SwiftRetribution:      true,
	Bloodlust:             true,
	StrengthOfEarthTotem:  proto.TristateEffect_TristateEffectImproved,
	LeaderOfThePack:       proto.TristateEffect_TristateEffectImproved,
	SanctifiedRetribution: true,
	DevotionAura:          proto.TristateEffect_TristateEffectImproved,
	RetributionAura:       true,
	IcyTalons:             true,
}
var FullPartyBuffs = &proto.PartyBuffs{
	HeroicPresence: true,
}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:     true,
	BlessingOfMight:     proto.TristateEffect_TristateEffectImproved,
	BlessingOfSanctuary: true,
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfStoneblood,
	DefaultPotion: proto.Potions_IndestructiblePotion,
	PrepopPotion:  proto.Potions_IndestructiblePotion,
	Food:          proto.Food_FoodDragonfinFilet,
}

var FullDebuffs = &proto.Debuffs{
	SunderArmor:        true,
	Mangle:             true,
	DemoralizingShout:  proto.TristateEffect_TristateEffectImproved,
	JudgementOfLight:   true,
	FaerieFire:         proto.TristateEffect_TristateEffectRegular,
	Misery:             true,
	FrostFever:         proto.TristateEffect_TristateEffectImproved,
	BloodFrenzy:        true,
	EbonPlaguebringer:  true,
	HeartOfTheCrusader: true,
}

var BloodITRotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"prepullActions": [
		{"action":{"castSpell":{"spellId":{"spellId":48263}}},"doAtValue":{"const":{"val":"-10s"}}},
		{"action":{"castSpell":{"spellId":{"spellId":42650}}},"doAtValue":{"const":{"val":"-6s"}}}
	],
	"priorityList": [
		{"action":{"autocastOtherCooldowns":{}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentHealthPercent":{}},"rhs":{"const":{"val":"40%"}}}},"castSpell":{"spellId":{"spellId":48792}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentHealthPercent":{}},"rhs":{"const":{"val":"40%"}}}},"castSpell":{"spellId":{"spellId":55233}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentHealthPercent":{}},"rhs":{"const":{"val":"60%"}}}},"castSpell":{"spellId":{"spellId":48982}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentHealthPercent":{}},"rhs":{"const":{"val":"60%"}}}},"castSpell":{"spellId":{"spellId":48707}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"currentHealthPercent":{}},"rhs":{"const":{"val":"60%"}}}},"castSpell":{"spellId":{"spellId":48743}}}},
		{"action":{"condition":{"cmp":{"op":"OpGe","lhs":{"currentRunicPower":{}},"rhs":{"const":{"val":"40"}}}},"castSpell":{"spellId":{"spellId":56815}}}},
		{"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55095}}}}},"castSpell":{"spellId":{"spellId":59131}}}},
		{"action":{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":55078}}}}},"castSpell":{"spellId":{"tag":1,"spellId":49921}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":55078}}},"rhs":{"const":{"val":"3s"}}}},"castSpell":{"spellId":{"spellId":50842}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"currentNonDeathRuneCount":{"runeType":"RuneFrost"}},"rhs":{"const":{"val":"0"}}}},{"cmp":{"op":"OpGt","lhs":{"currentNonDeathRuneCount":{"runeType":"RuneUnholy"}},"rhs":{"const":{"val":"0"}}}}]}},"castSpell":{"spellId":{"tag":1,"spellId":49924}}}},
		{"action":{"condition":{"cmp":{"op":"OpGt","lhs":{"currentRuneCount":{"runeType":"RuneDeath"}},"rhs":{"const":{"val":"0"}}}},"castSpell":{"spellId":{"spellId":59131}}}},
		{"action":{"condition":{"or":{"vals":[{"cmp":{"op":"OpGt","lhs":{"currentNonDeathRuneCount":{"runeType":"RuneBlood"}},"rhs":{"const":{"val":"1"}}}},{"spellIsReady":{"spellId":{"spellId":47568}}}]}},"castSpell":{"spellId":{"tag":1,"spellId":49930}}}},
		{"action":{"castSpell":{"spellId":{"spellId":46584}}}},
		{"action":{"castSpell":{"spellId":{"spellId":47568}}}},
		{"action":{"condition":{"cmp":{"op":"OpGe","lhs":{"currentRunicPower":{}},"rhs":{"const":{"val":"80"}}}},"castSpell":{"spellId":{"spellId":49895}}}}
	]
}`)

var BloodP1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{"id":40565,"enchant":3878,"gems":[41380,36767]},
	{"id":40387},
	{"id":39704,"enchant":3852,"gems":[40008]},
	{"id":40252,"enchant":3605},
	{"id":40559,"gems":[40008,40022]},
	{"id":40306,"enchant":3850,"gems":[40008,0]},
	{"id":40563,"enchant":3860,"gems":[40008,0]},
	{"id":39759,"gems":[40008,40008]},
	{"id":40567,"enchant":3822,"gems":[40008,40008]},
	{"id":40297,"enchant":3232},
	{"id":40718},
	{"id":40107},
	{"id":44063,"gems":[36767,36767]},
	{"id":42341,"gems":[40008,40008]},
	{"id":40406,"enchant":3847},
	{},
	{"id":40207}
]}`)
