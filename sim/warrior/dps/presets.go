package dps

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var PlayerOptionsArmsSlam = &proto.Player_Warrior{
	Warrior: &proto.Warrior{
		Talents:  ArmsSlamTalents,
		Options:  warriorOptions,
		Rotation: armsSlamRotation,
	},
}

var PlayerOptionsFury = &proto.Player_Warrior{
	Warrior: &proto.Warrior{
		Talents:  FuryTalents,
		Options:  warriorOptions,
		Rotation: warriorRotation,
	},
}

var ArmsSlamTalents = &proto.WarriorTalents{
	ImprovedHeroicStrike:          3,
	ImprovedRend:                  2,
	TacticalMastery:               3,
	ImprovedOverpower:             2,
	AngerManagement:               true,
	Impale:                        2,
	DeepWounds:                    3,
	TwoHandedWeaponSpecialization: 3,
	PoleaxeSpecialization:         5,
	TasteForBlood:                 0,
	MaceSpecialization:            0,
	SwordSpecialization:           0,
	WeaponMastery:                 0,
	StrengthOfArms:                2,
	ImprovedSlam:                  2,
	ImprovedMortalStrike:          3,
	UnrelentingAssault:            2,
	SuddenDeath:                   3,
	EndlessRage:                   true,
	BloodFrenzy:                   2,
	WreckingCrew:                  5,
	Bladestorm:                    true,

	Cruelty:           5,
	ArmoredToTheTeeth: 3,

	ImprovedBloodrage:   2,
	Incite:              3,
	ImprovedThunderClap: 3,
}

var FuryTalents = &proto.WarriorTalents{
	ImprovedHeroicStrike: 3,
	ImprovedRend:         2,
	TacticalMastery:      3,
	AngerManagement:      true,
	Impale:               2,
	DeepWounds:           3,

	Cruelty:                 5,
	ArmoredToTheTeeth:       3,
	UnbridledWrath:          2,
	ImprovedCleave:          3,
	DualWieldSpecialization: 5,
	Precision:               3,
	DeathWish:               true,
	ImprovedBerserkerRage:   1,
	Flurry:                  5,
	IntensifyRage:           3,
	ImprovedWhirlwind:       2,
	ImprovedBerserkerStance: 5,
	Bloodsurge:              3,
	UnendingFury:            5,
	TitansGrip:              true,
}

var armsSlamRotation = &proto.Warrior_Rotation{
	UseSlam: true,
	UseMs:   true,

	HsRageThreshold: 60,
	MsRageThreshold: 60,
	RendCdThreshold: 3,

	SpamExecute: false,

	MaintainDemoShout:   true,
	MaintainThunderClap: true,
}

var warriorRotation = &proto.Warrior_Rotation{
	UseRend: true,

	HsRageThreshold:   40,
	RendRageThreshold: 70,

	UseHsDuringExecute: true,
	UseWwDuringExecute: true,
	UseBtDuringExecute: true,
	UseSlamOverExecute: true,
}

var warriorOptions = &proto.Warrior_Options{
	StartingRage:    50,
	UseRecklessness: true,
	Shout:           proto.WarriorShout_WarriorShoutBattle,
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance: true,
	GiftOfTheWild:    proto.TristateEffect_TristateEffectImproved,
}
var FullPartyBuffs = &proto.PartyBuffs{
	BattleShout:     proto.TristateEffect_TristateEffectImproved,
	LeaderOfThePack: proto.TristateEffect_TristateEffectImproved,
}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:  true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
	BlessingOfMight:  proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	Drums: proto.Drums_DrumsOfBattle,
}

var FullDebuffs = &proto.Debuffs{
	BloodFrenzy:               true,
	FaerieFire:                proto.TristateEffect_TristateEffectImproved,
	ImprovedSealOfTheCrusader: true,
	JudgementOfWisdom:         true,
	Misery:                    true,
}

var FuryP1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 29021,
		"enchant": 29192,
		"gems": [
			32409,
			24048
		]
	},
	{
		"id": 29381
	},
	{
		"id": 29023,
		"enchant": 28888,
		"gems": [
			24048,
			24067
		]
	},
	{
		"id": 24259,
		"enchant": 34004,
		"gems": [
			24058
		]
	},
	{
		"id": 29019,
		"enchant": 24003,
		"gems": [
			24048,
			24048,
			24048
		]
	},
	{
		"id": 28795,
		"enchant": 27899,
		"gems": [
			24067,
			24058
		]
	},
	{
		"id": 28824,
		"enchant": 33995,
		"gems": [
			24067,
			24048
		]
	},
	{
		"id": 28779,
		"gems": [
			24058,
			24067
		]
	},
	{
		"id": 28741,
		"enchant": 29535,
		"gems": [
			24048,
			24048,
			24048
		]
	},
	{
		"id": 28608,
		"enchant": 28279,
		"gems": [
			24058,
			24048
		]
	},
	{
		"id": 28757
	},
	{
		"id": 30834
	},
	{
		"id": 29383
	},
	{
		"id": 28830
	},
	{
		"id": 28438,
		"enchant": 22559
	},
	{
		"id": 28729,
		"enchant": 22559
	},
	{
		"id": 30279
	}
]}`)
