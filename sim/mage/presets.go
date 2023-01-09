package mage

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var FireTalents = &proto.MageTalents{
	ArcaneSubtlety:      2,
	ArcaneFocus:         3,
	ArcaneConcentration: 5,
	SpellImpact:         3,
	StudentOfTheMind:    1,
	FocusMagic:          true,
	TormentTheWeak:      3,

	ImprovedFireball: 5,
	Ignite:           5,
	WorldInFlames:    3,
	Pyroblast:        true,
	BurningSoul:      1,
	ImprovedScorch:   3,
	MasterOfElements: 2,
	PlayingWithFire:  3,
	CriticalMass:     3,
	BlastWave:        true,
	FirePower:        5,
	Pyromaniac:       3,
	Combustion:       true,
	MoltenFury:       2,
	EmpoweredFire:    3,
	DragonsBreath:    true,
	Firestarter:      2,
	HotStreak:        3,
	Burnout:          5,
	LivingBomb:       true,
}

var FrostTalents = &proto.MageTalents{
	ArcaneSubtlety:      2,
	ArcaneFocus:         3,
	ArcaneConcentration: 5,
	SpellImpact:         3,
	StudentOfTheMind:    1,
	FocusMagic:          true,
	TormentTheWeak:      3,

	ImprovedFrostbolt:    5,
	IceFloes:             3,
	IceShards:            3,
	Precision:            3,
	PiercingIce:          3,
	IcyVeins:             true,
	ArcticReach:          2,
	FrostChanneling:      3,
	Shatter:              3,
	ColdSnap:             true,
	WintersChill:         3,
	IceBarrier:           true,
	ArcticWinds:          5,
	EmpoweredFrostbolt:   2,
	FingersOfFrost:       2,
	BrainFreeze:          3,
	SummonWaterElemental: true,
	EnduringWinter:       3,
	ChilledToTheBone:     5,
	DeepFreeze:           true,
}

var ArcaneTalents = &proto.MageTalents{
	ArcaneFocus:         3,
	ArcaneSubtlety:      2,
	ArcaneConcentration: 5,
	SpellImpact:         3,
	FocusMagic:          true,
	MagicAttunement:     1,
	StudentOfTheMind:    3,
	ArcaneMeditation:    3,
	TormentTheWeak:      3,
	PresenceOfMind:      true,
	ArcaneMind:          5,
	ArcaneInstability:   3,
	ArcanePotency:       2,
	ArcaneEmpowerment:   3,
	ArcanePower:         true,
	ArcaneFlows:         2,
	MindMastery:         5,
	MissileBarrage:      5,
	NetherwindPresence:  3,
	SpellPower:          2,
	ArcaneBarrage:       true,

	Incineration: 3,

	ImprovedFrostbolt: 2,
	IceFloes:          3,
	IceShards:         2,
	Precision:         3,
	IcyVeins:          true,
}

var fireMageOptions = &proto.Mage_Options{
	Armor:          proto.Mage_Options_MoltenArmor,
	ReactionTimeMs: 300,
	IgniteMunching: true,
}
var PlayerOptionsFire = &proto.Player_Mage{
	Mage: &proto.Mage{
		Talents: FireTalents,
		Options: fireMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type:                   proto.Mage_Rotation_Fire,
			PrimaryFireSpell:       proto.Mage_Rotation_Fireball,
			MaintainImprovedScorch: false,
			PyroblastDelayMs:       50,
		},
	},
}
var PlayerOptionsFireAOE = &proto.Player_Mage{
	Mage: &proto.Mage{
		Talents: FireTalents,
		Options: fireMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type: proto.Mage_Rotation_Aoe,
			Aoe:  proto.Mage_Rotation_Flamestrike,
		},
	},
}

var frostMageOptions = &proto.Mage_Options{
	Armor:          proto.Mage_Options_MageArmor,
	ReactionTimeMs: 300,
}
var PlayerOptionsFrost = &proto.Player_Mage{
	Mage: &proto.Mage{
		Talents: FrostTalents,
		Options: frostMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type: proto.Mage_Rotation_Frost,
		},
	},
}
var PlayerOptionsFrostAOE = &proto.Player_Mage{
	Mage: &proto.Mage{
		Talents: FrostTalents,
		Options: frostMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type: proto.Mage_Rotation_Aoe,
			Aoe:  proto.Mage_Rotation_Blizzard,
		},
	},
}

var arcaneMageOptions = &proto.Mage_Options{
	Armor:          proto.Mage_Options_MoltenArmor,
	ReactionTimeMs: 300,
}
var PlayerOptionsArcane = &proto.Player_Mage{
	Mage: &proto.Mage{
		Talents: ArcaneTalents,
		Options: arcaneMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type:                                       proto.Mage_Rotation_Arcane,
			ExtraBlastsDuringFirstAp:                   2,
			MissileBarrageBelowArcaneBlastStacks:       0,
			MissileBarrageBelowManaPercent:             0.1,
			BlastWithoutMissileBarrageAboveManaPercent: 0.2,
			Only_3ArcaneBlastStacksBelowManaPercent:    0.15,
		},
	},
}
var PlayerOptionsArcaneAOE = &proto.Player_Mage{
	Mage: &proto.Mage{
		Talents: ArcaneTalents,
		Options: arcaneMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type: proto.Mage_Rotation_Aoe,
			Aoe:  proto.Mage_Rotation_ArcaneExplosion,
		},
	},
}

var FullFireConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	Food:          proto.Food_FoodFirecrackerSalmon,
	DefaultPotion: proto.Potions_PotionOfSpeed,
	// DefaultConjured: proto.Conjured_ConjuredFlameCap,
}
var FullFrostConsumes = FullFireConsumes

var FullArcaneConsumes = FullFireConsumes

var P1ArcaneGear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40416,
		"enchant": 3820,
		"gems": [
			41285,
			39998
		]
	},
	{
		"id": 44661,
		"gems": [
			40026
		]
	},
	{
		"id": 40419,
		"enchant": 3810,
		"gems": [
			40051
		]
	},
	{
		"id": 44005,
		"enchant": 3722,
		"gems": [
			40026
		]
	},
	{
		"id": 44002,
		"enchant": 3832,
		"gems": [
			39998,
			39998
		]
	},
	{
		"id": 44008,
		"enchant": 2332,
		"gems": [
			39998,
			0
		]
	},
	{
		"id": 40415,
		"enchant": 3604,
		"gems": [
			39998,
			0
		]
	},
	{
		"id": 40561,
		"gems": [
			39998
		]
	},
	{
		"id": 40417,
		"enchant": 3719,
		"gems": [
			39998,
			40051
		]
	},
	{
		"id": 40558,
		"enchant": 3606
	},
	{
		"id": 40719
	},
	{
		"id": 40399
	},
	{
		"id": 39229
	},
	{
		"id": 40255
	},
	{
		"id": 40396,
		"enchant": 3834
	},
	{
		"id": 40273
	},
	{
		"id": 39426
	}
]}`)
var P1FrostGear = P1ArcaneGear
var P1FireGear = P1ArcaneGear
