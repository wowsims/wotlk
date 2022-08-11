package mage

import (
	"github.com/wowsims/wotlk/sim/core/items"
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
	ArcaneFocus:         5,
	ArcaneConcentration: 5,
	SpellImpact:         3,
	ArcaneMeditation:    3,

	ImprovedFrostbolt:    5,
	Precision:            3,
	IceShards:            5,
	IcyVeins:             true,
	PiercingIce:          5,
	FrostChanneling:      3,
	ColdSnap:             true,
	ImprovedConeOfCold:   2,
	IceFloes:             2,
	WintersChill:         4,
	ArcticWinds:          5,
	EmpoweredFrostbolt:   5,
	SummonWaterElemental: true,
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
	Armor: proto.Mage_Options_MoltenArmor,
}
var PlayerOptionsFire = &proto.Player_Mage{
	Mage: &proto.Mage{
		Talents: FireTalents,
		Options: fireMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type:                   proto.Mage_Rotation_Fire,
			PrimaryFireSpell:       proto.Mage_Rotation_Fireball,
			MaintainImprovedScorch: false,
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
	Armor: proto.Mage_Options_MageArmor,
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
	Armor: proto.Mage_Options_MoltenArmor,
}
var PlayerOptionsArcane = &proto.Player_Mage{
	Mage: &proto.Mage{
		Talents: ArcaneTalents,
		Options: arcaneMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type:                   proto.Mage_Rotation_Arcane,
			MinBlastBeforeMissiles: 4,
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

var FullRaidBuffs = &proto.RaidBuffs{
	GiftOfTheWild:     proto.TristateEffect_TristateEffectImproved,
	Bloodlust:         true,
	MoonkinAura:       proto.TristateEffect_TristateEffectRegular,
	ManaSpringTotem:   proto.TristateEffect_TristateEffectRegular,
	TotemOfWrath:      true,
	WrathOfAirTotem:   true,
	ArcaneBrilliance:  true,
	ArcaneEmpowerment: true,
	SwiftRetribution:  true,
	DivineSpirit:      true,
}
var FullFirePartyBuffs = &proto.PartyBuffs{}
var FullFrostPartyBuffs = FullFirePartyBuffs
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:     true,
	BlessingOfWisdom:    proto.TristateEffect_TristateEffectImproved,
	JudgementsOfTheWise: true,
}

var FullArcanePartyBuffs = &proto.PartyBuffs{
	ManaTideTotems: 1,
}
var FullArcaneIndividualBuffs = FullIndividualBuffs

var FullFireConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfTheFrostWyrm,
	Food:            proto.Food_FoodFirecrackerSalmon,
	DefaultPotion:   proto.Potions_PotionOfSpeed,
	DefaultConjured: proto.Conjured_ConjuredFlameCap,
}
var FullFrostConsumes = FullFireConsumes

var FullArcaneConsumes = FullFireConsumes

var FullDebuffs = &proto.Debuffs{
	JudgementOfWisdom: true,
	EarthAndMoon:      true,
	Misery:            true,
	ImprovedScorch:    true,
}

var P1ArcaneGear = items.EquipmentSpecFromJsonString(`{"items": [
	{
	  "id": 40416,
	  "enchant": 44877,
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
	  "enchant": 44874,
	  "gems": [
		40051
	  ]
	},
	{
	  "id": 44005,
	  "enchant": 55642,
	  "gems": [
		40026
	  ]
	},
	{
	  "id": 44002,
	  "enchant": 44489,
	  "gems": [
		39998,
		39998
	  ]
	},
	{
	  "id": 44008,
	  "enchant": 44498,
	  "gems": [
		39998,
		0
	  ]
	},
	{
	  "id": 40415,
	  "enchant": 54999,
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
	  "enchant": 41602,
	  "gems": [
		39998,
		40051
	  ]
	},
	{
	  "id": 40558,
	  "enchant": 55016
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
	  "enchant": 44487
	},
	{
	  "id": 40273
	},
	{
	  "id": 39426
	}
  ]
}`)
var P1FrostGear = P1ArcaneGear
var P1FireGear = P1ArcaneGear
