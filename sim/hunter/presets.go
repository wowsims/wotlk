package hunter

import (
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
)

var BMTalents = &proto.HunterTalents{
	ImprovedAspectOfTheHawk: 5,
	FocusedFire:             2,
	UnleashedFury:           5,
	Ferocity:                5,
	BestialDiscipline:       2,
	AnimalHandler:           1,
	Frenzy:                  5,
	FerociousInspiration:    3,
	BestialWrath:            true,
	SerpentsSwiftness:       5,
	TheBeastWithin:          true,

	LethalShots:    5,
	Efficiency:     5,
	GoForTheThroat: 2,
	AimedShot:      true,
	RapidKilling:   2,
	MortalShots:    5,
}

var SVTalents = &proto.HunterTalents{
	ImprovedAspectOfTheHawk: 5,
	FocusedFire:             2,

	LethalShots:         5,
	ImprovedHuntersMark: 5,
	GoForTheThroat:      2,
	RapidKilling:        1,

	MonsterSlaying:    3,
	HumanoidSlaying:   3,
	SavageStrikes:     2,
	CleverTraps:       2,
	Survivalist:       2,
	Surefooted:        3,
	SurvivalInstincts: 2,
	KillerInstinct:    3,
	LightningReflexes: 5,
	ThrillOfTheHunt:   2,
	ExposeWeakness:    3,
	MasterTactician:   5,
	Readiness:         true,
}

var PlayerOptionsBasic = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Talents:  BMTalents,
		Options:  basicOptions,
		Rotation: basicRotation,
	},
}

var PlayerOptionsFrench = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Talents:  BMTalents,
		Options:  windSerpentOptions,
		Rotation: frenchRotation,
	},
}

var PlayerOptionsMeleeWeave = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Talents:  BMTalents,
		Options:  windSerpentOptions,
		Rotation: meleeWeaveRotation,
	},
}

var PlayerOptionsSV = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Talents:  SVTalents,
		Options:  windSerpentOptions,
		Rotation: meleeWeaveRotation,
	},
}

var basicRotation = &proto.Hunter_Rotation{
	UseMultiShot:     true,
	UseArcaneShot:    false,
	Sting:            proto.Hunter_Rotation_SerpentSting,
	PrecastAimedShot: true,
	LazyRotation:     true,

	ViperStartManaPercent: 0.2,
	ViperStopManaPercent:  0.3,
}
var frenchRotation = &proto.Hunter_Rotation{
	UseMultiShot:     true,
	UseArcaneShot:    true,
	Sting:            proto.Hunter_Rotation_SerpentSting,
	PrecastAimedShot: false,

	ViperStartManaPercent: 0.3,
	ViperStopManaPercent:  0.5,
}
var meleeWeaveRotation = &proto.Hunter_Rotation{
	UseMultiShot:  true,
	UseArcaneShot: true,
	Weave:         proto.Hunter_Rotation_WeaveFull,
	TimeToWeaveMs: 500,
	PercentWeaved: 0.8,

	ViperStartManaPercent: 0.3,
	ViperStopManaPercent:  0.5,
}

var basicOptions = &proto.Hunter_Options{
	QuiverBonus: proto.Hunter_Options_Speed15,
	Ammo:        proto.Hunter_Options_AdamantiteStinger,
	PetType:     proto.Hunter_Options_Ravager,
	PetUptime:   0.9,
	LatencyMs:   15,
}

var windSerpentOptions = &proto.Hunter_Options{
	QuiverBonus:      proto.Hunter_Options_Speed15,
	Ammo:             proto.Hunter_Options_AdamantiteStinger,
	PetType:          proto.Hunter_Options_WindSerpent,
	PetUptime:        0.9,
	PetSingleAbility: true,
	LatencyMs:        15,
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance: true,
	GiftOfTheWild:    proto.TristateEffect_TristateEffectImproved,
}
var FullPartyBuffs = &proto.PartyBuffs{
	Bloodlust: 1,
	Drums:     proto.Drums_DrumsOfBattle,

	BattleShout:       proto.TristateEffect_TristateEffectImproved,
	LeaderOfThePack:   proto.TristateEffect_TristateEffectImproved,
	ManaSpringTotem:   proto.TristateEffect_TristateEffectRegular,
	GraceOfAirTotem:   proto.TristateEffect_TristateEffectRegular,
	WindfuryTotemRank: 5,
}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:  true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
	BlessingOfMight:  proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfRelentlessAssault,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredFlameCap,
	PetFood:         proto.PetFood_PetFoodKiblersBits,
}

var FullDebuffs = &proto.Debuffs{
	BloodFrenzy:               true,
	FaerieFire:                proto.TristateEffect_TristateEffectImproved,
	ImprovedSealOfTheCrusader: true,
	JudgementOfWisdom:         true,
	Misery:                    true,
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
