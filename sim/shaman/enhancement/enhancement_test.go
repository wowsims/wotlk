package enhancement

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get item effects included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterEnhancementShaman()
}

func TestEnhancement(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassShaman,
		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceOrc},

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: Phase1Gear},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "FT", SpecOptions: PlayerOptionsFTFT},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "WF", SpecOptions: PlayerOptionsWFWF},
		},
		Rotation: core.RotationCombo{Label: "FT", Rotation: DefaultFTRotation},
		OtherRotations: []core.RotationCombo{
			{Label: "WF", Rotation: DefaultWFRotation},
		},

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
				Equipment:     Phase1Gear,
				TalentsString: StandardTalents,
				Glyphs:        StandardGlyphs,
				Consumes:      FullConsumes,
				Spec:          PlayerOptionsFTFT,
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

var StandardTalents = "053030152-30405003105021333031131031051"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.ShamanMajorGlyph_GlyphOfFireNova),
	Major2: int32(proto.ShamanMajorGlyph_GlyphOfFlametongueWeapon),
	Major3: int32(proto.ShamanMajorGlyph_GlyphOfFeralSpirit),
}

var PlayerOptionsWFWF = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options:  enhShamWFWF,
		Rotation: &proto.EnhancementShaman_Rotation{},
	},
}

var PlayerOptionsFTFT = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options:  enhShamFTFT,
		Rotation: &proto.EnhancementShaman_Rotation{},
	},
}

//var enhShamRotationItemSwap = &proto.EnhancementShaman_Rotation{
//	RotationType:                 proto.EnhancementShaman_Rotation_Priority,
//	FirenovaManaThreshold:        3000,
//	ShamanisticRageManaThreshold: 25,
//	PrimaryShock:                 proto.EnhancementShaman_Rotation_Earth,
//	WeaveFlameShock:              true,
//	//Temp to test Item Swap, will switch to a more realistic swap with Phase 2 gear.
//	EnableItemSwap: true,
//	ItemSwap: &proto.ItemSwap{
//		MhItem: &proto.ItemSpec{
//			Id: 41752,
//		},
//		OhItem: &proto.ItemSpec{
//			Id:      41752,
//			Enchant: 3790,
//		},
//	},
//}

var enhShamWFWF = &proto.EnhancementShaman_Options{
	Shield:    proto.ShamanShield_WaterShield,
	Bloodlust: true,
	SyncType:  proto.ShamanSyncType_DelayOffhandSwings,
	ImbueMh:   proto.ShamanImbue_WindfuryWeapon,
	ImbueOh:   proto.ShamanImbue_WindfuryWeapon,
}

var enhShamFTFT = &proto.EnhancementShaman_Options{
	Shield:    proto.ShamanShield_LightningShield,
	Bloodlust: true,
	SyncType:  proto.ShamanSyncType_Auto,
	ImbueMh:   proto.ShamanImbue_FlametongueWeaponDownrank,
	ImbueOh:   proto.ShamanImbue_FlametongueWeapon,
	Totems: &proto.ShamanTotems{
		Earth:            proto.EarthTotem_StrengthOfEarthTotem,
		Air:              proto.AirTotem_WindfuryTotem,
		Water:            proto.WaterTotem_ManaSpringTotem,
		Fire:             proto.FireTotem_MagmaTotem,
		UseFireElemental: true,
	},
}

var enhShamWFFT = &proto.EnhancementShaman_Options{
	Shield:    proto.ShamanShield_LightningShield,
	Bloodlust: true,
	SyncType:  proto.ShamanSyncType_NoSync,
	ImbueMh:   proto.ShamanImbue_WindfuryWeapon,
	ImbueOh:   proto.ShamanImbue_FlametongueWeapon,
}

var FullConsumes = &proto.Consumes{
	DefaultConjured: proto.Conjured_ConjuredFlameCap,
}

var DefaultFTRotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"prepullActions": [
		{"action":{"castSpell":{"spellId":{"spellId":66842}}},"doAtValue":{"const":{"val":"-3s"}}},
		{"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAtValue":{"const":{"val":"-1s"}}}
	],
	"priorityList": [
		{"action":{"autocastOtherCooldowns":{}}},
		{"action":{"condition":{"cmp":{"op":"OpEq","lhs":{"auraNumStacks":{"auraId":{"spellId":53817}}},"rhs":{"const":{"val":"5"}}}},"castSpell":{"spellId":{"spellId":49238}}}},
		{"action":{"castSpell":{"spellId":{"spellId":17364}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":49233}}},"rhs":{"const":{"val":"0s"}}}},"castSpell":{"spellId":{"spellId":49233}}}},
		{"action":{"condition":{"not":{"val":{"auraIsActive":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":17364}}}}},"castSpell":{"spellId":{"spellId":17364}}}},
		{"action":{"castSpell":{"spellId":{"spellId":49231}}}},
		{"action":{"condition":{"cmp":{"op":"OpLt","lhs":{"totemRemainingTime":{"totemType":"Water"}},"rhs":{"const":{"val":"20s"}}}},"castSpell":{"spellId":{"spellId":66842}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":58734}}},"rhs":{"const":{"val":"100ms"}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":2894}}}}}]}},"castSpell":{"spellId":{"spellId":58734}}}},
		{"action":{"castSpell":{"spellId":{"spellId":61657}}}},
		{"action":{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":49281}}}}},"castSpell":{"spellId":{"spellId":49281}}}},
		{"action":{"castSpell":{"spellId":{"spellId":60103}}}}
	]
}`)

var DefaultWFRotation = core.APLRotationFromJsonString(`{
	"type": "TypeAPL",
	"prepullActions": [
		{"action":{"castSpell":{"spellId":{"spellId":66842}}},"doAtValue":{"const":{"val":"-3s"}}},
		{"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAtValue":{"const":{"val":"-1s"}}}
	],
	"priorityList": [
		{"action":{"autocastOtherCooldowns":{}}},
		{"action":{"condition":{"cmp":{"op":"OpEq","lhs":{"auraNumStacks":{"auraId":{"spellId":53817}}},"rhs":{"const":{"val":"5"}}}},"castSpell":{"spellId":{"spellId":49238}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGe","lhs":{"auraNumStacks":{"auraId":{"spellId":53817}}},"rhs":{"const":{"val":"3"}}}},{"cmp":{"op":"OpLt","lhs":{"math":{"op":"OpAdd","lhs":{"const":{"val":"300ms"}},"rhs":{"spellCastTime":{"spellId":{"spellId":49238}}}}},"rhs":{"autoTimeToNext":{}}}}]}},"castSpell":{"spellId":{"spellId":49238}}}},
		{"action":{"castSpell":{"spellId":{"spellId":17364}}}},
		{"action":{"condition":{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":49233}}},"rhs":{"const":{"val":"0s"}}}},"castSpell":{"spellId":{"spellId":49233}}}},
		{"action":{"castSpell":{"spellId":{"spellId":49231}}}},
		{"action":{"condition":{"cmp":{"op":"OpLt","lhs":{"totemRemainingTime":{"totemType":"Water"}},"rhs":{"const":{"val":"20s"}}}},"castSpell":{"spellId":{"spellId":66842}}}},
		{"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":58734}}},"rhs":{"const":{"val":"100ms"}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":2894}}}}}]}},"castSpell":{"spellId":{"spellId":58734}}}},
		{"action":{"castSpell":{"spellId":{"spellId":61657}}}},
		{"action":{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":49281}}}}},"castSpell":{"spellId":{"spellId":49281}}}},
		{"action":{"castSpell":{"spellId":{"spellId":60103}}}}
	]
}`)

var Phase1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{"id":40543,"enchant":3817,"gems":[41398,40014]},
	{"id":44661,"gems":[40014]},
	{"id":40524,"enchant":3808,"gems":[40014]},
	{"id":40403,"enchant":3605},
	{"id":40523,"enchant":3832,"gems":[40003,40014]},
	{"id":40282,"enchant":3845,"gems":[42702,0]},
	{"id":40520,"enchant":3604,"gems":[42154,0]},
	{"id":40275,"gems":[42156]},
	{"id":40522,"enchant":3823,"gems":[39999,42156]},
	{"id":40367,"enchant":3606,"gems":[40058]},
	{"id":40474},
	{"id":40074},
	{"id":40684},
	{"id":37390},
	{"id":39763,"enchant":3789},
	{"id":39468,"enchant":3789},
	{"id":40322}
]}`)
