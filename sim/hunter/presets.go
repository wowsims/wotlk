package hunter

import (
	"github.com/wowsims/wotlk/sim/core/items"
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

var basicRotation = &proto.Hunter_Rotation{
	Sting: proto.Hunter_Rotation_SerpentSting,

	ViperStartManaPercent: 0.2,
	ViperStopManaPercent:  0.3,
}

var basicOptions = &proto.Hunter_Options{
	Ammo:       proto.Hunter_Options_SaroniteRazorheads,
	PetType:    proto.Hunter_Options_Wolf,
	PetTalents: FerocityTalents,
	PetUptime:  0.9,
	LatencyMs:  15,
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance:     true,
	GiftOfTheWild:        proto.TristateEffect_TristateEffectImproved,
	Bloodlust:            true,
	BattleShout:          proto.TristateEffect_TristateEffectImproved,
	LeaderOfThePack:      proto.TristateEffect_TristateEffectImproved,
	ManaSpringTotem:      proto.TristateEffect_TristateEffectRegular,
	StrengthOfEarthTotem: proto.TristateEffect_TristateEffectImproved,
	WindfuryTotem:        proto.TristateEffect_TristateEffectImproved,
	BlessingOfKings:  true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
	BlessingOfMight:  proto.TristateEffect_TristateEffectImproved,
}
var FullPartyBuffs = &proto.PartyBuffs{}
var FullIndividualBuffs = &proto.IndividualBuffs{
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfRelentlessAssault,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredFlameCap,
	PetFood:         proto.PetFood_PetFoodKiblersBits,
}

var FullDebuffs = &proto.Debuffs{
	BloodFrenzy:       true,
	FaerieFire:        proto.TristateEffect_TristateEffectImproved,
	JudgementOfWisdom: true,
	Misery:            true,
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 28275,
		"enchant": 29192,
		"gems": [
			24028,
			32409
		]
	},
	{
		"id": 29381
	},
	{
		"id": 27801,
		"enchant": 28888,
		"gems": [
			31868,
			24028
		]
	},
	{
		"id": 24259,
		"enchant": 34004,
		"gems": [
			24028
		]
	},
	{
		"id": 28228,
		"enchant": 24003,
		"gems": [
			24028,
			24028,
			24055
		]
	},
	{
		"id": 29246,
		"enchant": 34002
	},
	{
		"id": 27474,
		"enchant": 33152,
		"gems": [
			24028,
			24028
		]
	},
	{
		"id": 28828,
		"gems": [
			24055,
			31868
		]
	},
	{
		"id": 30739,
		"enchant": 29535,
		"gems": [
			24028,
			24028,
			24028
		]
	},
	{
		"id": 28545,
		"enchant": 28279,
		"gems": [
			24028,
			24061
		]
	},
	{
		"id": 28757
	},
	{
		"id": 28791
	},
	{
		"id": 28830
	},
	{
		"id": 29383
	},
	{
		"id": 28435,
		"enchant": 22556
	},
	{
		"id": 28772,
		"enchant": 23766
	}
]}`)
