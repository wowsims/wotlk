package protection

import (
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
)

var ImpaleProtTalents = &proto.WarriorTalents{
	ImprovedHeroicStrike: 3,
	Deflection:           5,
	ImprovedThunderClap:  3,
	AngerManagement:      true,
	DeepWounds:           3,
	Impale:               2,

	Cruelty: 3,

	Anticipation:                  5,
	ShieldSpecialization:          5,
	Toughness:                     5,
	ImprovedShieldBlock:           true,
	Defiance:                      3,
	ImprovedSunderArmor:           3,
	ShieldMastery:                 1,
	OneHandedWeaponSpecialization: 5,
	ShieldSlam:                    true,
	FocusedRage:                   3,
	Vitality:                      5,
	Devastate:                     true,
}

var PlayerOptionsBasic = &proto.Player_ProtectionWarrior{
	ProtectionWarrior: &proto.ProtectionWarrior{
		Talents:  ImpaleProtTalents,
		Options:  warriorOptions,
		Rotation: warriorRotation,
	},
}

var warriorRotation = &proto.ProtectionWarrior_Rotation{
	DemoShout:       proto.ProtectionWarrior_Rotation_DemoShoutMaintain,
	ThunderClap:     proto.ProtectionWarrior_Rotation_ThunderClapMaintain,
	ShieldBlock:     proto.ProtectionWarrior_Rotation_ShieldBlockToProcRevenge,
	HsRageThreshold: 30,
}

var warriorOptions = &proto.ProtectionWarrior_Options{
	Shout:                proto.WarriorShout_WarriorShoutCommanding,
	PrecastShout:         false,
	PrecastShoutT2:       false,
	PrecastShoutSapphire: false,

	StartingRage: 0,
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance: true,
	GiftOfTheWild:    proto.TristateEffect_TristateEffectImproved,
	Thorns:           proto.TristateEffect_TristateEffectImproved,
}
var FullPartyBuffs = &proto.PartyBuffs{
	Bloodlust:            1,
	Drums:                proto.Drums_DrumsOfBattle,
	BattleShout:          proto.TristateEffect_TristateEffectImproved,
	LeaderOfThePack:      proto.TristateEffect_TristateEffectImproved,
	GraceOfAirTotem:      proto.TristateEffect_TristateEffectRegular,
	StrengthOfEarthTotem: proto.StrengthOfEarthType_EnhancingTotems,
	WindfuryTotemRank:    5,
	WindfuryTotemIwt:     2,
}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings: true,
	BlessingOfMight: proto.TristateEffect_TristateEffectImproved,
	UnleashedRage:   true,
}

var FullConsumes = &proto.Consumes{
	BattleElixir:   proto.BattleElixir_ElixirOfMastery,
	GuardianElixir: proto.GuardianElixir_GiftOfArthas,
}

var FullDebuffs = &proto.Debuffs{
	BloodFrenzy:               true,
	FaerieFire:                proto.TristateEffect_TristateEffectImproved,
	ImprovedSealOfTheCrusader: true,
	Misery:                    true,
	ShadowEmbrace:             true,
	ScorpidSting:              true,
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 29011,
		"enchant": 29192,
		"gems": [
			25896,
			24033
		]
	},
	{
		"id": 28244,
		"gems": [
			33782
		]
	},
	{
		"id": 29023,
		"enchant": 28911,
		"gems": [
			24033,
			24033
		]
	},
	{
		"id": 28672,
		"enchant": 34004
	},
	{
		"id": 29012,
		"enchant": 24003,
		"gems": [
			24033,
			24033,
			24033
		]
	},
	{
		"id": 28996,
		"enchant": 22533,
		"gems": [
			33782
		]
	},
	{
		"id": 30644,
		"enchant": 33153
	},
	{
		"id": 28995
	},
	{
		"id": 28621,
		"enchant": 29536,
		"gems": [
			24033,
			24033,
			24033
		]
	},
	{
		"id": 28747,
		"enchant": 35297,
		"gems": [
			24033,
			24033
		]
	},
	{
		"id": 30834
	},
	{
		"id": 29279
	},
	{
		"id": 28121
	},
	{
		"id": 29387
	},
	{
		"id": 28749,
		"enchant": 22559
	},
	{
		"id": 28825,
		"enchant": 28282,
		"gems": [
			24033
		]
	},
	{
		"id": 28826
	}
]}`)
