package dps

import (
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
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
	Deflection:                    2,
	ImprovedThunderClap:           3,
	AngerManagement:               true,
	DeepWounds:                    3,
	TwoHandedWeaponSpecialization: 5,
	Impale:                        2,
	DeathWish:                     true,
	SwordSpecialization:           5,
	ImprovedDisciplines:           2,
	BloodFrenzy:                   2,
	MortalStrike:                  true,

	Cruelty:                   5,
	ImprovedDemoralizingShout: 5,
	CommandingPresence:        5,
	ImprovedSlam:              2,
	SweepingStrikes:           true,
	WeaponMastery:             2,
	Flurry:                    3,
}

var FuryTalents = &proto.WarriorTalents{
	ImprovedHeroicStrike: 3,
	AngerManagement:      true,
	DeepWounds:           3,
	Impale:               2,

	Cruelty:                 5,
	UnbridledWrath:          5,
	CommandingPresence:      5,
	DualWieldSpecialization: 5,
	SweepingStrikes:         true,
	WeaponMastery:           2,
	Flurry:                  5,
	Precision:               3,
	Bloodthirst:             true,
	ImprovedWhirlwind:       1,
	ImprovedBerserkerStance: 5,
	Rampage:                 true,
}

var armsSlamRotation = &proto.Warrior_Rotation{
	UseOverpower: true,
	UseHamstring: true,
	UseSlam:      true,

	HsRageThreshold:        70,
	HamstringRageThreshold: 75,
	OverpowerRageThreshold: 20,
	SlamLatency:            100,
	SlamGcdDelay:           400,
	SlamMsWwDelay:          2000,

	UseSlamDuringExecute: true,
	UseWwDuringExecute:   true,
	UseMsDuringExecute:   true,
	UseHsDuringExecute:   true,

	MaintainDemoShout:   true,
	MaintainThunderClap: true,
}

var warriorRotation = &proto.Warrior_Rotation{
	UseOverpower: true,
	UseHamstring: true,

	HsRageThreshold:        70,
	HamstringRageThreshold: 75,
	OverpowerRageThreshold: 20,
	RampageCdThreshold:     5,

	UseHsDuringExecute: true,
	UseWwDuringExecute: true,
	UseBtDuringExecute: true,
}

var warriorOptions = &proto.Warrior_Options{
	StartingRage:         50,
	UseRecklessness:      true,
	Shout:                proto.WarriorShout_WarriorShoutBattle,
	PrecastShout:         false,
	PrecastShoutT2:       false,
	PrecastShoutSapphire: false,
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
