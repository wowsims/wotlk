package dps

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var PlayerOptionsArms = &proto.Player_Warrior{
	Warrior: &proto.Warrior{
		Talents:  ArmsTalents,
		Options:  warriorOptions,
		Rotation: armsRotation,
	},
}

var PlayerOptionsFury = &proto.Player_Warrior{
	Warrior: &proto.Warrior{
		Talents:  FuryTalents,
		Options:  warriorOptions,
		Rotation: furyRotation,
	},
}

var ArmsTalents = &proto.WarriorTalents{
	ImprovedHeroicStrike:          3,
	ImprovedRend:                  2,
	TacticalMastery:               3,
	ImprovedOverpower:             2,
	AngerManagement:               true,
	Impale:                        2,
	DeepWounds:                    3,
	TwoHandedWeaponSpecialization: 3,
	TasteForBlood:                 3,
	PoleaxeSpecialization:         5,
	SweepingStrikes:               true,
	WeaponMastery:                 1,
	MortalStrike:                  true,
	StrengthOfArms:                2,
	ImprovedSlam:                  2,
	ImprovedMortalStrike:          3,
	UnrelentingAssault:            2,
	SuddenDeath:                   3,
	EndlessRage:                   true,
	BloodFrenzy:                   2,
	WreckingCrew:                  5,
	Bladestorm:                    true,

	ArmoredToTheTeeth: 3,
	Cruelty:           5,

	ImprovedBloodrage:   2,
	ImprovedThunderClap: 3,
	Incite:              3,
}

var FuryTalents = &proto.WarriorTalents{
	ImprovedHeroicStrike:          3,
	ImprovedRend:                  2,
	IronWill:                      2,
	TacticalMastery:               3,
	AngerManagement:               true,
	Impale:                        2,
	DeepWounds:                    3,
	TwoHandedWeaponSpecialization: 3,

	ArmoredToTheTeeth:       3,
	Cruelty:                 5,
	UnbridledWrath:          2,
	ImprovedCleave:          3,
	PiercingHowl:            true,
	CommandingPresence:      1,
	DualWieldSpecialization: 5,
	ImprovedExecute:         2,
	Precision:               3,
	DeathWish:               true,
	ImprovedBerserkerRage:   1,
	Flurry:                  5,
	IntensifyRage:           3,
	Bloodthirst:             true,
	ImprovedWhirlwind:       2,
	ImprovedBerserkerStance: 5,
	Rampage:                 true,
	Bloodsurge:              3,
	UnendingFury:            5,
	TitansGrip:              true,
}

var armsRotation = &proto.Warrior_Rotation{
	UseRend:   true,
	UseMs:     true,
	UseCleave: false,

	HsRageThreshold:   50,
	MsRageThreshold:   35,
	SlamRageThreshold: 25,
	RendCdThreshold:   0,

	SpamExecute: false,

	UseHsDuringExecute: true,

	MaintainDemoShout:   false,
	MaintainThunderClap: false,

	StanceOption: proto.Warrior_Rotation_DefaultStance,
}

var furyRotation = &proto.Warrior_Rotation{
	UseRend:   false,
	UseCleave: false,

	HsRageThreshold:          30,
	RendRageThresholdBelow:   100,
	SlamRageThreshold:        25,
	RendCdThreshold:          0,
	RendHealthThresholdAbove: 20,

	UseHsDuringExecute: true,
	UseWwDuringExecute: true,
	UseBtDuringExecute: true,
	UseSlamOverExecute: true,

	MaintainDemoShout:   false,
	MaintainThunderClap: false,

	StanceOption: proto.Warrior_Rotation_DefaultStance,
}

var warriorOptions = &proto.Warrior_Options{
	StartingRage:    50,
	UseRecklessness: true,
	Shout:           proto.WarriorShout_WarriorShoutBattle,
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfEndlessRage,
	DefaultPotion: proto.Potions_IndestructiblePotion,
	PrepopPotion:  proto.Potions_PotionOfSpeed,
	Food:          proto.Food_FoodDragonfinFilet,
}

var FuryP1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 44006,
		"enchant": 3817,
		"gems": [
			41285,
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
		"id": 40530,
		"enchant": 3808,
		"gems": [
			40037
		]
	},
	{
		"id": 40403,
		"enchant": 3605
	},
	{
		"id": 40539,
		"enchant": 3832,
		"gems": [
			42142
		]
	},
	{
		"id": 39765,
		"enchant": 3845,
		"gems": [
			39996,
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
		"id": 40205,
		"gems": [
			42142
		]
	},
	{
		"id": 40529,
		"enchant": 3823,
		"gems": [
			39996,
			40022
		]
	},
	{
		"id": 40591,
		"enchant": 3606
	},
	{
		"id": 43993,
		"gems": [
			42142
		]
	},
	{
		"id": 40717
	},
	{
		"id": 42987
	},
	{
		"id": 40256
	},
	{
		"id": 40384,
		"enchant": 3789
	},
	{
		"id": 40384,
		"enchant": 3789
	},
	{
		"id": 40385
	}
]}`)
