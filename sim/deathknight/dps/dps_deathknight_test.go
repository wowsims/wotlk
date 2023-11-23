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

		GearSet:     core.GetGearSet("../../../ui/deathknight/gear_sets", "p3_blood"),
		Talents:     BloodTalents,
		Glyphs:      BloodDefaultGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBlood},
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/deathknight/apls", "blood_dps"),
		},

		ItemFilter: ItemFilter,
	}))
}

func TestUnholy(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathknight,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet:     core.GetGearSet("../../../ui/deathknight/gear_sets", "p3_uh_dw"),
		Talents:     UnholyTalents,
		Glyphs:      UnholyDefaultGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsUnholy},

		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/deathknight/apls", "uh_2h_ss"),
			core.GetAplRotation("../../../ui/deathknight/apls", "uh_dnd_aoe"),
			core.GetAplRotation("../../../ui/deathknight/apls", "unholy_dw_ss"),
		},

		ItemFilter: ItemFilter,
	}))
}

func TestFrost(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathknight,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet:     core.GetGearSet("../../../ui/deathknight/gear_sets", "p3_frost"),
		Talents:     FrostTalents,
		Glyphs:      FrostDefaultGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFrost},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "Desync", SpecOptions: PlayerOptionsDesyncFrost},
		},

		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/deathknight/apls", "frost_bl_pesti"),
			core.GetAplRotation("../../../ui/deathknight/apls", "frost_uh_pesti"),
		},

		ItemFilter: ItemFilter,
	}))
}

func TestFrostUH(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathknight,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet:     core.GetGearSet("../../../ui/deathknight/gear_sets", "p3_frost"),
		Talents:     FrostUHTalents,
		Glyphs:      FrostUHDefaultGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFrost},
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/deathknight/apls", "frost_uh_pesti"),
		},

		ItemFilter: ItemFilter,
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

var FrostUHTalents = "01-32002350342203012300033101351-230200305003"
var FrostUHDefaultGlyphs = &proto.Glyphs{
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

var PlayerOptionsDesyncFrost = &proto.Player_Deathknight{
	Deathknight: &proto.Deathknight{
		Options:  deathKnightOptions,
		Rotation: frostDesyncRotation,
	},
}

var bloodRotation = &proto.Deathknight_Rotation{
	ArmyOfTheDead:        proto.Deathknight_Rotation_PreCast,
	DrwDiseases:          proto.Deathknight_Rotation_Pestilence,
	UseEmpowerRuneWeapon: true,
	PreNerfedGargoyle:    false,
	UseDancingRuneWeapon: true,
	BloodSpender:         proto.Deathknight_Rotation_HS,
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

var frostRotation = &proto.Deathknight_Rotation{
	UseEmpowerRuneWeapon: true,
}

var frostDesyncRotation = &proto.Deathknight_Rotation{
	UseEmpowerRuneWeapon: true,
	DesyncRotation:       true,
}

var deathKnightOptions = &proto.Deathknight_Options{
	UnholyFrenzyTarget:  &proto.UnitReference{Type: proto.UnitReference_Player, Index: 0},
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

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypePlate,

	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
	},
}
