package restoration

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterRestorationShaman()
}

func TestRestoration(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassShaman,
		Race:  proto.Race_RaceTroll,

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Standard", SpecOptions: PlayerOptionsStandard},
		Rotation:    core.RotationCombo{Label: "Default", Rotation: DefaultRotation},

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeFist,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeShield,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeMail,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeTotem,
			},
		},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:          proto.Race_RaceOrc,
				Class:         proto.Class_ClassShaman,
				Equipment:     P1Gear,
				Consumes:      FullConsumes,
				Spec:          PlayerOptionsStandard,
				Buffs:         core.FullIndividualBuffs,
				TalentsString: StandardTalents,
				Glyphs:        StandardGlyphs,
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

var StandardTalents = "-3020503-50005331335310501122331251"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.ShamanMajorGlyph_GlyphOfChainHeal),
	Major2: int32(proto.ShamanMajorGlyph_GlyphOfEarthShield),
	Major3: int32(proto.ShamanMajorGlyph_GlyphOfEarthlivingWeapon),
}

var BasicTotems = &proto.ShamanTotems{
	Earth: proto.EarthTotem_TremorTotem,
	Air:   proto.AirTotem_WrathOfAirTotem,
	Water: proto.WaterTotem_ManaSpringTotem,
	Fire:  proto.FireTotem_FlametongueTotem,
}

var restoShamOptions = &proto.RestorationShaman_Options{
	Shield:    proto.ShamanShield_WaterShield,
	Bloodlust: true,
	Totems:    BasicTotems,
}
var PlayerOptionsStandard = &proto.Player_RestorationShaman{
	RestorationShaman: &proto.RestorationShaman{
		Options:  restoShamOptions,
		Rotation: &proto.RestorationShaman_Rotation{},
	},
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfBlindingLight,
	Food:            proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:   proto.Potions_SuperManaPotion,
	PrepopPotion:    proto.Potions_DestructionPotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var DefaultRotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"priorityList": [
		{"action":{"autocastOtherCooldowns":{}}}
	]
}`)

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{"id":40510,"enchant":3820,"gems":[41401,39998]},
	{"id":44662,"gems":[40051]},
	{"id":40513,"enchant":3810,"gems":[39998]},
	{"id":44005,"enchant":3831,"gems":[40027]},
	{"id":40508,"enchant":2381,"gems":[39998,40051]},
	{"id":40209,"enchant":2332,"gems":[0]},
	{"id":40564,"enchant":3246,"gems":[0]},
	{"id":40327,"gems":[39998]},
	{"id":40512,"enchant":3721,"gems":[39998,40027]},
	{"id":39734,"enchant":3244},
	{"id":40399},
	{"id":40375},
	{"id":37111},
	{"id":40685},
	{"id":40395,"enchant":3834},
	{"id":40401,"enchant":1128},
	{"id":40709}
]}`)
