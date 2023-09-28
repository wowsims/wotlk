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

		GearSet:     core.GearSetCombo{Label: "Blood P3 ", GearSet: BloodP3Gear},
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

		GearSet:     core.GearSetCombo{Label: "Unholy P3 ", GearSet: UnholyDwP3Gear},
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

		GearSet:     core.GearSetCombo{Label: "Frost P3", GearSet: FrostP3Gear},
		Talents:     FrostTalents,
		Glyphs:      FrostDefaultGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFrost},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "Desync", SpecOptions: PlayerOptionsDesyncFrost},
		},

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

func TestFrostUH(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathknight,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet:     core.GearSetCombo{Label: "Frost P1", GearSet: FrostP3Gear},
		Talents:     FrostUHTalents,
		Glyphs:      FrostUHDefaultGlyphs,
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

var BloodP3Gear = core.EquipmentSpecFromJsonString(`{"items": [
		  {"id":48493,"enchant":3817,"gems":[41285,40142]},
		  {"id":47458,"gems":[40142]},
		  {"id":48495,"enchant":3808,"gems":[40111]},
		  {"id":47546,"enchant":3831,"gems":[42142]},
		  {"id":47449,"enchant":3832,"gems":[49110,42142,40142]},
		  {"id":48008,"enchant":3845,"gems":[40111,0]},
		  {"id":48492,"enchant":3604,"gems":[40142,0]},
		  {"id":47429,"gems":[40142,40142,40111]},
		  {"id":48494,"enchant":3823,"gems":[40142,40111]},
		  {"id":45599,"enchant":3606,"gems":[40111,40111]},
		  {"id":47993,"gems":[40111,45862]},
		  {"id":47413,"gems":[40142]},
		  {"id":45931},
		  {"id":47464},
		  {"id":47446,"enchant":3368,"gems":[42142,40141]},
		  {},
		  {"id":47673}
]}`)

var UnholyDwP3Gear = core.EquipmentSpecFromJsonString(`{"items": [
		{"id":48493,"enchant":3817,"gems":[41398,40146]},
		  {"id":47458,"gems":[40146]},
		  {"id":48495,"enchant":3808,"gems":[40111]},
		  {"id":47548,"enchant":3831,"gems":[40111]},
		  {"id":48491,"enchant":3832,"gems":[42142,42142]},
		  {"id":45663,"enchant":3845,"gems":[40111,0]},
		  {"id":48492,"enchant":3604,"gems":[40146,0]},
		  {"id":47429,"gems":[40111,45862,40111]},
		  {"id":47465,"enchant":3823,"gems":[49110,40111,40146]},
		  {"id":45599,"enchant":3606,"gems":[40111,40111]},
		  {"id":47413,"gems":[40146]},
		  {"id":45534,"gems":[42142]},
		  {"id":47464},
		  {"id":45609},
		  {"id":47528,"enchant":3368,"gems":[40111]},
		  {"id":47528,"enchant":3368,"gems":[40111]},
		  {"id":47673}
]}`)

var FrostP3Gear = core.EquipmentSpecFromJsonString(`{ "items": [
		{"id":48493,"enchant":3817,"gems":[41398,40142]},
		  {"id":45459,"gems":[40111]},
		  {"id":48495,"enchant":3808,"gems":[40111]},
		  {"id":47548,"enchant":3831,"gems":[40111]},
		  {"id":48491,"enchant":3832,"gems":[42142,42142]},
		  {"id":45663,"enchant":3845,"gems":[40111,0]},
		  {"id":47492,"enchant":3604,"gems":[49110,40111,0]},
		  {"id":45241,"gems":[40111,42142,40111]},
		  {"id":48494,"enchant":3823,"gems":[40142,40111]},
		  {"id":47473,"enchant":3606,"gems":[40142,40111]},
		  {"id":46966,"gems":[40111]},
		  {"id":45534,"gems":[40111]},
		  {"id":47464},
		  {"id":45931},
		  {"id":47528,"enchant":3370,"gems":[40111]},
		  {"id":47528,"enchant":3368,"gems":[40111]},
		  {"id":40207}
]}`)
