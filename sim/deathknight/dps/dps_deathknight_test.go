package dps

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get item effects included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterDpsDeathknight()
}

func TestBlood(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathknight,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet:     core.GearSetCombo{Label: "Blood P1 ", GearSet: BloodP1Gear},
		Talents:     BloodTalents,
		Glyphs:      BloodDefaultGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBlood},

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

func TestUnholy(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathknight,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet:     core.GearSetCombo{Label: "Unholy P1 ", GearSet: UnholyDwP1Gear},
		Talents:     UnholyTalents,
		Glyphs:      UnholyDefaultGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsUnholy},

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

func TestFrost(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathknight,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet:     core.GearSetCombo{Label: "Frost P1", GearSet: FrostP1Gear},
		Talents:     FrostTalents,
		Glyphs:      FrostDefaultGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFrost},

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

var BloodTalents = "2305120530003303231023001351--230220305003"
var BloodDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.DeathknightMajorGlyph_GlyphOfDancingRuneWeapon),
	Major2: int32(proto.DeathknightMajorGlyph_GlyphOfDeathStrike),
	Major3: int32(proto.DeathknightMajorGlyph_GlyphOfDisease),
	// No interesting minor glyphs.
}

var FrostTalents = "23050005-32005350352203012300033101351"
var FrostDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.DeathknightMajorGlyph_GlyphOfFrostStrike),
	Major2: int32(proto.DeathknightMajorGlyph_GlyphOfObliterate),
	Major3: int32(proto.DeathknightMajorGlyph_GlyphOfDisease),
	// No interesting minor glyphs.
}

var UnholyTalents = "-320043500002-2300303050032152000150013133051"
var UnholyDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.DeathknightMajorGlyph_GlyphOfTheGhoul),
	Major2: int32(proto.DeathknightMajorGlyph_GlyphOfDarkDeath),
	Major3: int32(proto.DeathknightMajorGlyph_GlyphOfDeathAndDecay),
	// No interesting minor glyphs.
}

var PlayerOptionsBlood = &proto.Player_Deathknight{
	Deathknight: &proto.Deathknight{
		Options:  deathKnightOptions,
		Rotation: bloodRotation,
	},
}

var PlayerOptionsUnholy = &proto.Player_Deathknight{
	Deathknight: &proto.Deathknight{
		Options:  deathKnightOptions,
		Rotation: unholyRotation,
	},
}

var PlayerOptionsFrost = &proto.Player_Deathknight{
	Deathknight: &proto.Deathknight{
		Options:  deathKnightOptions,
		Rotation: frostRotation,
	},
}

var bloodRotation = &proto.Deathknight_Rotation{
	ArmyOfTheDead:        proto.Deathknight_Rotation_PreCast,
	DrwDiseases:          proto.Deathknight_Rotation_Pestilence,
	UseEmpowerRuneWeapon: true,
	PreNerfedGargoyle:    false,
}

var unholyRotation = &proto.Deathknight_Rotation{
	UseDeathAndDecay:     true,
	StartingPresence:     proto.Deathknight_Rotation_Unholy,
	BlPresence:           proto.Deathknight_Rotation_Blood,
	Presence:             proto.Deathknight_Rotation_Blood,
	GargoylePresence:     proto.Deathknight_Rotation_Unholy,
	UseEmpowerRuneWeapon: true,
	UseGargoyle:          true,
	BtGhoulFrenzy:        false,
	HoldErwArmy:          false,
	PreNerfedGargoyle:    false,
	BloodRuneFiller:      proto.Deathknight_Rotation_BloodBoil,
	ArmyOfTheDead:        proto.Deathknight_Rotation_AsMajorCd,
	BloodTap:             proto.Deathknight_Rotation_GhoulFrenzy,
}

var frostRotation = &proto.Deathknight_Rotation{}

var deathKnightOptions = &proto.Deathknight_Options{
	UnholyFrenzyTarget:  &proto.RaidTarget{TargetIndex: 0},
	DrwPestiApply:       true,
	StartingRunicPower:  0,
	PetUptime:           1,
	PrecastGhoulFrenzy:  false,
	PrecastHornOfWinter: true,
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfEndlessRage,
	DefaultPotion: proto.Potions_PotionOfSpeed,
	PrepopPotion:  proto.Potions_PotionOfSpeed,
	Food:          proto.Food_FoodDragonfinFilet,
}

var BloodP1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 44006,
		"enchant": 3817,
		"gems": [
		  41398,
		  42702
		]
	  },
	  {
		"id": 44664,
		"gems": [
		  39996
		]
	  },
	  {
		"id": 40557,
		"enchant": 3808,
		"gems": [
		  39996
		]
	  },
	  {
		"id": 40403,
		"enchant": 3831
	  },
	  {
		"id": 40550,
		"enchant": 3832,
		"gems": [
		  42142,
		  42142
		]
	  },
	  {
		"id": 40330,
		"enchant": 3845,
		"gems": [
		  42142,
		  0
		]
	  },
	  {
		"id": 40552,
		"enchant": 3604,
		"gems": [
		  39996,
		  0
		]
	  },
	  {
		"id": 40317,
		"gems": [
		  39996
		]
	  },
	  {
		"id": 40556,
		"enchant": 3823,
		"gems": [
		  39996,
		  39996
		]
	  },
	  {
		"id": 40591,
		"enchant": 3606
	  },
	  {
		"id": 40075
	  },
	  {
		"id": 39401
	  },
	  {
		"id": 40256
	  },
	  {
		"id": 42987
	  },
	  {
		"id": 40384,
		"enchant": 3368
	  },
	  {},
	  {
		"id": 40207
	  }
]}`)

var UnholyDwP1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 44006,
		"enchant": 3817,
		"gems": [
			41398,
			42702
		]
	},
	{
		"id": 39421
	},
	{
		"id": 40557,
		"enchant": 3808,
		"gems": [
			39996
		]
	},
	{
		"id": 40403,
		"enchant": 3831
	},
	{
		"id": 40550,
		"enchant": 3832,
		"gems": [
			42142,
			39996
		]
	},
	{
		"id": 40330,
		"enchant": 3845,
		"gems": [
			39996,
			0
		]
	},
	{
		"id": 40347,
		"enchant": 3604,
		"gems": [
			39996,
			0
		]
	},
	{
		"id": 40278,
		"gems": [
			42142,
			42142
		]
	},
	{
		"id": 40294,
		"enchant": 3823
	},
	{
		"id": 39706,
		"enchant": 3606,
		"gems": [
			39996
		]
	},
	{
		"id": 39401
	},
	{
		"id": 40075
	},
	{
		"id": 37390
	},
	{
		"id": 42987
	},
	{
		"id": 40402,
		"enchant": 3368
	},
	{
		"id": 40491,
		"enchant": 3368
	},
	{
		"id": 42620
	}
]}`)

var FrostP1Gear = core.EquipmentSpecFromJsonString(`{ "items": [
	{
		"id": 44006,
		"enchant": 3817,
		"gems": [
			41398,
			42702
		]
	},
	{
		"id": 44664,
		"gems": [
			39996
		]
	},
	{
		"id": 40557,
		"enchant": 3808,
		"gems": [
			39996
		]
	},
	{
		"id": 40403,
		"enchant": 3831
	},
	{
		"id": 40550,
		"enchant": 3832,
		"gems": [
			42142,
			39996
		]
	},
	{
		"id": 40330,
		"enchant": 3845,
		"gems": [
			39996,
			0
		]
	},
	{
		"id": 40552,
		"enchant": 3604,
		"gems": [
			39996,
			0
		]
	},
	{
		"id": 40278,
		"gems": [
			39996,
			42142
		]
	},
	{
		"id": 40556,
		"enchant": 3823,
		"gems": [
			42142,
			39996
		]
	},
	{
		"id": 40591,
		"enchant": 3606
	},
	{
		"id": 39401
	},
	{
		"id": 40075
	},
	{
		"id": 40256
	},
	{
		"id": 42987
	},
	{
		"id": 40189,
		"enchant": 3370
	},
	{
		"id": 40189,
		"enchant": 3368
	},
	{
		"id": 40207
	}
]}`)
