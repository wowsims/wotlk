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

		GearSet:     core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "p1"),
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "FT", SpecOptions: PlayerOptionsFTFT},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "WF", SpecOptions: PlayerOptionsWFWF},
		},
		Rotation: core.GetAplRotation("../../../ui/enhancement_shaman/apls", "default_ft"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/enhancement_shaman/apls", "default_wf"),
			core.GetAplRotation("../../../ui/enhancement_shaman/apls", "phase_3"),
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
				Equipment:     core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "p1").GearSet,
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
		Options: enhShamWFWF,
	},
}

var PlayerOptionsFTFT = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options: enhShamFTFT,
	},
}

var enhShamWFWF = &proto.EnhancementShaman_Options{
	Shield:   proto.ShamanShield_WaterShield,
	SyncType: proto.ShamanSyncType_DelayOffhandSwings,
	ImbueMh:  proto.ShamanImbue_WindfuryWeapon,
	ImbueOh:  proto.ShamanImbue_WindfuryWeapon,
}

var enhShamFTFT = &proto.EnhancementShaman_Options{
	Shield:   proto.ShamanShield_LightningShield,
	SyncType: proto.ShamanSyncType_Auto,
	ImbueMh:  proto.ShamanImbue_FlametongueWeaponDownrank,
	ImbueOh:  proto.ShamanImbue_FlametongueWeapon,
	Totems: &proto.ShamanTotems{
		Earth:            proto.EarthTotem_StrengthOfEarthTotem,
		Air:              proto.AirTotem_WindfuryTotem,
		Water:            proto.WaterTotem_ManaSpringTotem,
		Fire:             proto.FireTotem_MagmaTotem,
		UseFireElemental: true,
	},
}

var enhShamWFFT = &proto.EnhancementShaman_Options{
	Shield:   proto.ShamanShield_LightningShield,
	SyncType: proto.ShamanSyncType_NoSync,
	ImbueMh:  proto.ShamanImbue_WindfuryWeapon,
	ImbueOh:  proto.ShamanImbue_FlametongueWeapon,
}

var FullConsumes = &proto.Consumes{
	DefaultConjured: proto.Conjured_ConjuredFlameCap,
}
