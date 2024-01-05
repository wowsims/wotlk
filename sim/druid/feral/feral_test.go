package feral

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterFeralDruid()
}

var FeralItemFilter = core.ItemFilter{
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
}

func TestFeral(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDruid,
		Race:  proto.Race_RaceTauren,

		GearSet:     core.GetGearSet("../../../ui/feral_druid/gear_sets", "p1"),
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsMonoCat},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "Default-NoBleed", SpecOptions: PlayerOptionsMonoCatNoBleed},
			{Label: "Flower-Aoe", SpecOptions: PlayerOptionsFlowerCatAoe},
		},
		Rotation:   core.GetAplRotation("../../../ui/feral_druid/apls", "default"),
		ItemFilter: FeralItemFilter,
	}))
}

func TestFeralApl(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDruid,
		Race:  proto.Race_RaceTauren,

		GearSet:     core.GetGearSet("../../../ui/feral_druid/gear_sets", "p3"),
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsMonoCat},
		Rotation:    core.GetAplRotation("../../../ui/feral_druid/apls", "default"),
		ItemFilter:  FeralItemFilter,
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:      proto.Race_RaceTauren,
				Class:     proto.Class_ClassDruid,
				Equipment: core.GetGearSet("../../../ui/feral_druid/gear_sets", "p1").GearSet,
				Consumes:  FullConsumes,
				Spec:      PlayerOptionsMonoCat,
				Buffs:     core.FullIndividualBuffs,
				Glyphs:    StandardGlyphs,

				InFrontOfTarget: true,
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

var StandardTalents = "-503202132322010053120230310511-205503012"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.DruidMajorGlyph_GlyphOfOmenOfClarity),
	Major2: int32(proto.DruidMajorGlyph_GlyphOfShred),
	Major3: int32(proto.DruidMajorGlyph_GlyphOfBerserk),
	Minor1: int32(proto.DruidMinorGlyph_GlyphOfTheWild),
}

var PlayerOptionsMonoCat = &proto.Player_FeralDruid{
	FeralDruid: &proto.FeralDruid{
		Options: &proto.FeralDruid_Options{
			InnervateTarget:   &proto.UnitReference{}, // no Innervate
			LatencyMs:         100,
			AssumeBleedActive: true,
		},
		Rotation: &proto.FeralDruid_Rotation{
			RotationType:       proto.FeralDruid_Rotation_SingleTarget,
			BearWeaveType:      proto.FeralDruid_Rotation_None,
			UseRake:            true,
			UseBite:            true,
			MinCombosForRip:    5,
			MinCombosForBite:   5,
			BiteTime:           4.0,
			MaintainFaerieFire: true,
			BerserkBiteThresh:  25.0,
			BerserkFfThresh:    15.0,
			MaxFfDelay:         0.7,
			MinRoarOffset:      24.0,
			RipLeeway:          3,
			SnekWeave:          false,
			FlowerWeave:        false,
			RaidTargets:        30,
			PrePopOoc:          true,
		},
	},
}

var PlayerOptionsMonoCatNoBleed = &proto.Player_FeralDruid{
	FeralDruid: &proto.FeralDruid{
		Options: &proto.FeralDruid_Options{
			InnervateTarget:   &proto.UnitReference{}, // no Innervate
			LatencyMs:         100,
			AssumeBleedActive: false,
		},
		Rotation: &proto.FeralDruid_Rotation{
			RotationType:       proto.FeralDruid_Rotation_SingleTarget,
			BearWeaveType:      proto.FeralDruid_Rotation_None,
			UseRake:            true,
			UseBite:            true,
			MinCombosForRip:    5,
			MinCombosForBite:   5,
			BiteTime:           4.0,
			MaintainFaerieFire: true,
			BerserkBiteThresh:  25.0,
			BerserkFfThresh:    15.0,
			MaxFfDelay:         0.7,
			MinRoarOffset:      24.0,
			RipLeeway:          3,
			SnekWeave:          false,
			FlowerWeave:        false,
			RaidTargets:        30,
			PrePopOoc:          true,
		},
	},
}

var PlayerOptionsFlowerCatAoe = &proto.Player_FeralDruid{
	FeralDruid: &proto.FeralDruid{
		Options: &proto.FeralDruid_Options{
			InnervateTarget:   &proto.UnitReference{}, // no Innervate
			LatencyMs:         100,
			AssumeBleedActive: false,
		},
		Rotation: &proto.FeralDruid_Rotation{
			RotationType:       proto.FeralDruid_Rotation_Aoe,
			BearWeaveType:      proto.FeralDruid_Rotation_None,
			UseRake:            true,
			UseBite:            true,
			MinCombosForRip:    5,
			MinCombosForBite:   5,
			BiteTime:           4.0,
			MaintainFaerieFire: true,
			BerserkBiteThresh:  25.0,
			BerserkFfThresh:    15.0,
			MaxFfDelay:         0.7,
			MinRoarOffset:      24.0,
			RipLeeway:          3,
			SnekWeave:          false,
			FlowerWeave:        true,
			RaidTargets:        30,
			PrePopOoc:          true,
		},
	},
}

var FullConsumes = &proto.Consumes{
	BattleElixir:    proto.BattleElixir_ElixirOfMajorAgility,
	GuardianElixir:  proto.GuardianElixir_ElixirOfMajorMageblood,
	Food:            proto.Food_FoodGrilledMudfish,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}
