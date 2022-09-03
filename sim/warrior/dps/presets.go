package dps

import (
	"github.com/wowsims/wotlk/sim/core/items"
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
		Rotation: warriorRotation,
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
	ImprovedHeroicStrike:          3,
	ImprovedRend:                  2,
	IronWill:                      2,
	TacticalMastery:               3,
	AngerManagement:               true,
	Impale:                        2,
	DeepWounds:                    3,
	TwoHandedWeaponSpecialization: 3,

	Cruelty:                 5,
	ArmoredToTheTeeth:       3,
	UnbridledWrath:          2,
	ImprovedCleave:          3,
	PiercingHowl:            true,
	CommandingPresence:      1,
	DualWieldSpecialization: 5,
	Enrage:                  4,
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

var armsRotation = &proto.Warrior_Rotation{
	UseRend:   true,
	UseMs:     true,
	UseCleave: false,

	HsRageThreshold:   60,
	MsRageThreshold:   50,
	SlamRageThreshold: 15,
	RendCdThreshold:   1,

	SpamExecute: false,

	UseHsDuringExecute: true,

	MaintainDemoShout:   false,
	MaintainThunderClap: false,
}

var warriorRotation = &proto.Warrior_Rotation{
	UseRend:   true,
	UseCleave: false,

	HsRageThreshold:        60,
	RendRageThresholdBelow: 70,
	SlamRageThreshold:      15,
	RendCdThreshold:        1,

	UseHsDuringExecute: true,
	UseWwDuringExecute: true,
	UseBtDuringExecute: true,
	UseSlamOverExecute: true,

	MaintainDemoShout:   false,
	MaintainThunderClap: false,
}

var warriorOptions = &proto.Warrior_Options{
	StartingRage:    50,
	UseRecklessness: true,
	Shout:           proto.WarriorShout_WarriorShoutBattle,
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance: true,
	GiftOfTheWild:    proto.TristateEffect_TristateEffectImproved,
	BattleShout:      proto.TristateEffect_TristateEffectImproved,
	LeaderOfThePack:  proto.TristateEffect_TristateEffectImproved,
}
var FullPartyBuffs = &proto.PartyBuffs{}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:  true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
	BlessingOfMight:  proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfEndlessRage,
	DefaultPotion: proto.Potions_IndestructiblePotion,
	PrepopPotion:  proto.Potions_IndestructiblePotion,
	Food:          proto.Food_FoodDragonfinFilet,
}

var FullDebuffs = &proto.Debuffs{
	BloodFrenzy:       true,
	FaerieFire:        proto.TristateEffect_TristateEffectImproved,
	JudgementOfWisdom: true,
	Misery:            true,
}

var FuryP1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40528,
		"enchant": 44879,
		"gems": [
			41398,
			39996
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
		"enchant": 44871,
		"gems": [
			40058
		]
	},
	{
		"id": 40403,
		"enchant": 55002
	},
	{
		"id": 40525,
		"enchant": 44489,
		"gems": [
			42142,
			49110
		]
	},
	{
		"id": 40733,
		"enchant": 44484,
		"gems": [
			0
		]
	},
	{
		"id": 40541,
		"enchant": 54999,
		"gems": [
			0
		]
	},
	{
		"id": 40317,
		"gems": [
			42142
		]
	},
	{
		"id": 40529,
		"enchant": 38374,
		"gems": [
			39996,
			39996
		]
	},
	{
		"id": 40591,
		"enchant": 55016
	},
	{
		"id": 43993,
		"gems": [
			39996
		]
	},
	{
		"id": 40075
	},
	{
		"id": 42987
	},
	{
		"id": 40256
	},
	{
		"id": 40384,
		"enchant": 44492
	},
	{
		"id": 40384,
		"enchant": 44492
	},
	{
		"id": 40385
	}
]}`)
