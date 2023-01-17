package protection

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get item effects included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterProtectionPaladin()
}

func TestProtection(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassPaladin,
		Race:       proto.Race_RaceBloodElf,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Protection Paladin SOV", SpecOptions: DefaultOptions},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{
				Label: "Protection Paladin SOC",
				SpecOptions: &proto.Player_ProtectionPaladin{
					ProtectionPaladin: &proto.ProtectionPaladin{
						Options: &proto.ProtectionPaladin_Options{
							Judgement:            proto.PaladinJudgement_JudgementOfWisdom,
							Seal:                 proto.PaladinSeal_Command,
							Aura:                 proto.PaladinAura_RetributionAura,
							DamageTakenPerSecond: 0,
						},
						Rotation: defaultProtRotation,
					},
				},
			},
			{
				Label: "Protection Paladin SOR",
				SpecOptions: &proto.Player_ProtectionPaladin{
					ProtectionPaladin: &proto.ProtectionPaladin{
						Options: &proto.ProtectionPaladin_Options{
							Judgement:            proto.PaladinJudgement_JudgementOfWisdom,
							Seal:                 proto.PaladinSeal_Righteousness,
							Aura:                 proto.PaladinAura_RetributionAura,
							DamageTakenPerSecond: 0,
						},
						Rotation: defaultProtRotation,
					},
				},
			},
		},

		IsTank:          true,
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
				Race:      proto.Race_RaceBloodElf,
				Class:     proto.Class_ClassPaladin,
				Equipment: P1Gear,
				Consumes:  FullConsumes,
				Spec:      DefaultOptions,
				Buffs:     core.FullIndividualBuffs,
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

var StandardTalents = "-05005135200132311333312321-511302012003"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.PaladinMajorGlyph_GlyphOfSealOfVengeance),
	Major2: int32(proto.PaladinMajorGlyph_GlyphOfRighteousDefense),
	Major3: int32(proto.PaladinMajorGlyph_GlyphOfDivinePlea),
	Minor1: int32(proto.PaladinMinorGlyph_GlyphOfLayOnHands),
	Minor2: int32(proto.PaladinMinorGlyph_GlyphOfSenseUndead),
}

var defaultProtRotation = &proto.ProtectionPaladin_Rotation{}

var defaultProtOptions = &proto.ProtectionPaladin_Options{
	Judgement: proto.PaladinJudgement_JudgementOfWisdom,
	Seal:      proto.PaladinSeal_Vengeance,
	Aura:      proto.PaladinAura_RetributionAura,
}

var DefaultOptions = &proto.Player_ProtectionPaladin{
	ProtectionPaladin: &proto.ProtectionPaladin{
		Options:  defaultProtOptions,
		Rotation: defaultProtRotation,
	},
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfStoneblood,
	Food:            proto.Food_FoodDragonfinFilet,
	DefaultPotion:   proto.Potions_IndestructiblePotion,
	PrepopPotion:    proto.Potions_IndestructiblePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40581,
		"enchant": 3818,
		"gems": [
			41396,
			36767
		]
	},
	{
		"id": 40387
	},
	{
		"id": 40584,
		"enchant": 3852,
		"gems": [
			49110
		]
	},
	{
		"id": 40410,
		"enchant": 3605
	},
	{
		"id": 40579,
		"enchant": 3832,
		"gems": [
			36767,
			40022
		]
	},
	{
		"id": 39764,
		"enchant": 3850,
		"gems": [
			0
		]
	},
	{
		"id": 40580,
		"enchant": 3860,
		"gems": [
			40008,
			0
		]
	},
	{
		"id": 39759,
		"enchant": 3601,
		"gems": [
			40008,
			40008
		]
	},
	{
		"id": 40589,
		"enchant": 3822
	},
	{
		"id": 39717,
		"enchant": 3606,
		"gems": [
			40089
		]
	},
	{
		"id": 40718
	},
	{
		"id": 40107
	},
	{
		"id": 44063,
		"gems": [
			36767,
			40089
		]
	},
	{
		"id": 37220
	},
	{
		"id": 40345,
		"enchant": 3788
	},
	{
		"id": 40400,
		"enchant": 3849
	},
	{
		"id": 40707
	}
]}`)
