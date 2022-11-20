package hunter

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var BMTalents = &proto.HunterTalents{
	ImprovedAspectOfTheHawk: 5,
	EnduranceTraining:       1,
	FocusedFire:             2,
	ImprovedRevivePet:       2,
	AspectMastery:           true,
	UnleashedFury:           5,
	Ferocity:                5,
	SpiritBond:              1,
	Intimidation:            true,
	BestialDiscipline:       2,
	AnimalHandler:           2,
	Frenzy:                  4,
	FerociousInspiration:    3,
	BestialWrath:            true,
	CatlikeReflexes:         2,
	SerpentsSwiftness:       5,
	Longevity:               3,
	TheBeastWithin:          true,
	CobraStrikes:            2,
	KindredSpirits:          5,
	BeastMastery:            true,

	LethalShots:    5,
	CarefulAim:     3,
	MortalShots:    5,
	GoForTheThroat: 2,
	AimedShot:      true,

	ImprovedTracking: 1,
}

var MMTalents = &proto.HunterTalents{
	ImprovedAspectOfTheHawk: 5,
	FocusedFire:             2,

	FocusedAim:                 3,
	LethalShots:                5,
	CarefulAim:                 3,
	MortalShots:                5,
	GoForTheThroat:             1,
	AimedShot:                  true,
	RapidKilling:               2,
	ImprovedStings:             3,
	Readiness:                  true,
	Barrage:                    3,
	CombatExperience:           2,
	RangedWeaponSpecialization: 3,
	PiercingShots:              3,
	TrueshotAura:               true,
	ImprovedBarrage:            3,
	MasterMarksman:             5,
	WildQuiver:                 3,
	SilencingShot:              true,
	ImprovedSteadyShot:         3,
	MarkedForDeath:             5,
	ChimeraShot:                true,

	ImprovedTracking:  5,
	SurvivalInstincts: 2,
}

var SVTalents = &proto.HunterTalents{
	FocusedAim:     2,
	LethalShots:    5,
	CarefulAim:     3,
	MortalShots:    5,
	GoForTheThroat: 1,
	AimedShot:      true,

	ImprovedTracking:  5,
	TrapMastery:       3,
	SurvivalInstincts: 2,
	Survivalist:       5,
	TNT:               3,
	LockAndLoad:       3,
	HunterVsWild:      3,
	KillerInstinct:    3,
	LightningReflexes: 5,
	Resourcefulness:   2,
	ExposeWeakness:    2,
	WyvernSting:       true,
	ThrillOfTheHunt:   3,
	MasterTactician:   5,
	NoxiousStings:     3,
	BlackArrow:        true,
	SniperTraining:    3,
	HuntingParty:      1,
	ExplosiveShot:     true,
}

var FerocityTalents = &proto.HunterPetTalents{
	CobraReflexes:  2,
	Dive:           true,
	SpikedCollar:   3,
	BoarsSpeed:     true,
	CullingTheHerd: 3,
	SpidersBite:    3,
	Rabid:          true,
	CallOfTheWild:  true,
	WildHunt:       1,
}

var DefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.HunterMajorGlyph_GlyphOfSteadyShot),
	Major2: int32(proto.HunterMajorGlyph_GlyphOfSerpentSting),
	Major3: int32(proto.HunterMajorGlyph_GlyphOfKillShot),
	// No interesting minor glyphs.
}

var PlayerOptionsMM = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Talents:  MMTalents,
		Options:  basicOptions,
		Rotation: basicRotation,
	},
}

var PlayerOptionsBM = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Talents:  BMTalents,
		Options:  basicOptions,
		Rotation: basicRotation,
	},
}

var PlayerOptionsSV = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Talents:  SVTalents,
		Options:  basicOptions,
		Rotation: basicRotation,
	},
}

var PlayerOptionsAOE = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Talents:  SVTalents,
		Options:  basicOptions,
		Rotation: aoeRotation,
	},
}

var basicRotation = &proto.Hunter_Rotation{
	Sting: proto.Hunter_Rotation_SerpentSting,

	ViperStartManaPercent: 0.2,
	ViperStopManaPercent:  0.3,
}

var aoeRotation = &proto.Hunter_Rotation{
	TrapWeave:         true,
	TimeToTrapWeaveMs: 2000,

	ViperStartManaPercent: 0.2,
	ViperStopManaPercent:  0.3,
}

var basicOptions = &proto.Hunter_Options{
	Ammo:       proto.Hunter_Options_SaroniteRazorheads,
	PetType:    proto.Hunter_Options_Wolf,
	PetTalents: FerocityTalents,
	PetUptime:  0.9,

	SniperTrainingUptime: 0.8,
	UseHuntersMark:       true,
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfRelentlessAssault,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredFlameCap,
	PetFood:         proto.PetFood_PetFoodKiblersBits,
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40505,
		"enchant": 3817,
		"gems": [
			41398,
			42143
		]
	},
	{
		"id": 44664,
		"gems": [
			42143
		]
	},
	{
		"id": 40507,
		"enchant": 3808,
		"gems": [
			39997
		]
	},
	{
		"id": 40403,
		"enchant": 3605
	},
	{
		"id": 43998,
		"enchant": 3832,
		"gems": [
			42143,
			39997
		]
	},
	{
		"id": 40282,
		"enchant": 3845,
		"gems": [
			39997,
			0
		]
	},
	{
		"id": 40541,
		"enchant": 3604,
		"gems": [
			0
		]
	},
	{
		"id": 39762,
		"enchant": 3601,
		"gems": [
			39997
		]
	},
	{
		"id": 40331,
		"enchant": 3823,
		"gems": [
			39997,
			49110
		]
	},
	{
		"id": 40549,
		"enchant": 3606
	},
	{
		"id": 40074
	},
	{
		"id": 40474
	},
	{
		"id": 40684
	},
	{
		"id": 44253
	},
	{
		"id": 40388,
		"enchant": 3827
	},
	{},
	{
		"id": 40385,
		"enchant": 3608
	}
]}`)
