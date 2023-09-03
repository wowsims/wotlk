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
			{Label: "EnhItemSwap", SpecOptions: PlayerOptionsItemSwap},
			{Label: "Wf", SpecOptions: PlayerOptionsWFWF},
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
		Rotation: enhShamRotation,
	},
}

var PlayerOptionsFTFT = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options:  enhShamFTFT,
		Rotation: enhShamRotation,
	},
}

var PlayerOptionsItemSwap = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options:  enhShamWFFT,
		Rotation: enhShamRotationItemSwap,
	},
}

var enhShamRotation = &proto.EnhancementShaman_Rotation{
	RotationType:                 proto.EnhancementShaman_Rotation_Priority,
	FirenovaManaThreshold:        3000,
	ShamanisticRageManaThreshold: 25,
	PrimaryShock:                 proto.EnhancementShaman_Rotation_Earth,
	WeaveFlameShock:              true,
}

var enhShamRotationItemSwap = &proto.EnhancementShaman_Rotation{
	RotationType:                 proto.EnhancementShaman_Rotation_Priority,
	FirenovaManaThreshold:        3000,
	ShamanisticRageManaThreshold: 25,
	PrimaryShock:                 proto.EnhancementShaman_Rotation_Earth,
	WeaveFlameShock:              true,
	//Temp to test Item Swap, will switch to a more realistic swap with Phase 2 gear.
	EnableItemSwap: true,
	ItemSwap: &proto.ItemSwap{
		MhItem: &proto.ItemSpec{
			Id: 41752,
		},
		OhItem: &proto.ItemSpec{
			Id:      41752,
			Enchant: 3790,
		},
	},
}

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

var Phase1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40543,
		"enchant": 3817,
		"gems": [
			41398,
			40014
		]
	},
	{
		"id": 44661,
		"gems": [
			40014
		]
	},
	{
		"id": 40524,
		"enchant": 3808,
		"gems": [
			40014
		]
	},
	{
		"id": 40403,
		"enchant": 3605
	},
	{
		"id": 40523,
		"enchant": 3832,
		"gems": [
			40003,
			40014
		]
	},
	{
		"id": 40282,
		"enchant": 3845,
		"gems": [
			42702,
			0
		]
	},
	{
		"id": 40520,
		"enchant": 3604,
		"gems": [
			42154,
			0
		]
	},
	{
		"id": 40275,
		"gems": [
			42156
		]
	},
	{
		"id": 40522,
		"enchant": 3823,
		"gems": [
			39999,
			42156
		]
	},
	{
		"id": 40367,
		"enchant": 3606,
		"gems": [
			40058
		]
	},
	{
		"id": 40474
	},
	{
		"id": 40074
	},
	{
		"id": 40684
	},
	{
		"id": 37390
	},
	{
		"id": 39763,
		"enchant": 3789
	},
	{
		"id": 39468,
		"enchant": 3789
	},
	{
		"id": 40322
	}
]}`)
