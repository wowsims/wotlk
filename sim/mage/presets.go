package mage

import (
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
)

var FireTalents = &proto.MageTalents{
	ArcaneSubtlety: 2,

	ImprovedFireball:  5,
	Ignite:            5,
	Incineration:      2,
	Pyroblast:         true,
	ImprovedScorch:    3,
	MasterOfElements:  3,
	PlayingWithFire:   3,
	CriticalMass:      3,
	BlastWave:         true,
	FirePower:         5,
	Pyromaniac:        3,
	Combustion:        true,
	MoltenFury:        2,
	EmpoweredFireball: 5,

	ImprovedFrostbolt:  4,
	ElementalPrecision: 3,
	IceShards:          5,
	IcyVeins:           true,
}

var FrostTalents = &proto.MageTalents{
	ArcaneFocus:         5,
	ArcaneConcentration: 5,
	ArcaneImpact:        3,
	ArcaneMeditation:    3,

	ImprovedFrostbolt:    5,
	ElementalPrecision:   3,
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
	ArcaneFocus:         5,
	ArcaneConcentration: 5,
	ArcaneImpact:        3,
	ArcaneMeditation:    3,
	PresenceOfMind:      true,
	ArcaneMind:          5,
	ArcaneInstability:   3,
	ArcanePotency:       3,
	ArcanePower:         true,
	SpellPower:          2,
	MindMastery:         5,

	ImprovedFrostbolt:  5,
	ElementalPrecision: 3,
	IceShards:          5,
	IcyVeins:           true,
	PiercingIce:        5,
	FrostChanneling:    3,
	ColdSnap:           true,
}

var fireMageOptions = &proto.Mage_Options{
	Armor: proto.Mage_Options_MageArmor,
}
var PlayerOptionsFire = &proto.Player_Mage{
	Mage: &proto.Mage{
		Talents: FireTalents,
		Options: fireMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type: proto.Mage_Rotation_Fire,
			Fire: &proto.Mage_Rotation_FireRotation{
				PrimarySpell:           proto.Mage_Rotation_FireRotation_Fireball,
				MaintainImprovedScorch: true,
				WeaveFireBlast:         true,
			},
		},
	},
}
var PlayerOptionsFireAOE = &proto.Player_Mage{
	Mage: &proto.Mage{
		Talents: FireTalents,
		Options: fireMageOptions,
		Rotation: &proto.Mage_Rotation{
			MultiTargetRotation: true,
			Aoe: &proto.Mage_Rotation_AoeRotation{
				Rotation: proto.Mage_Rotation_AoeRotation_Flamestrike,
			},
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
			Type:  proto.Mage_Rotation_Frost,
			Frost: &proto.Mage_Rotation_FrostRotation{},
		},
	},
}
var PlayerOptionsFrostAOE = &proto.Player_Mage{
	Mage: &proto.Mage{
		Talents: FrostTalents,
		Options: frostMageOptions,
		Rotation: &proto.Mage_Rotation{
			MultiTargetRotation: true,
			Aoe: &proto.Mage_Rotation_AoeRotation{
				Rotation: proto.Mage_Rotation_AoeRotation_Blizzard,
			},
		},
	},
}

var arcaneMageOptions = &proto.Mage_Options{
	Armor: proto.Mage_Options_MageArmor,
}
var PlayerOptionsArcane = &proto.Player_Mage{
	Mage: &proto.Mage{
		Talents: ArcaneTalents,
		Options: arcaneMageOptions,
		Rotation: &proto.Mage_Rotation{
			Type: proto.Mage_Rotation_Arcane,
			Arcane: &proto.Mage_Rotation_ArcaneRotation{
				Filler:                     proto.Mage_Rotation_ArcaneRotation_ArcaneMissilesFrostbolt,
				ArcaneBlastsBetweenFillers: 3,
				StartRegenRotationPercent:  0.2,
				StopRegenRotationPercent:   0.3,
			},
		},
	},
}
var PlayerOptionsArcaneAOE = &proto.Player_Mage{
	Mage: &proto.Mage{
		Talents: ArcaneTalents,
		Options: arcaneMageOptions,
		Rotation: &proto.Mage_Rotation{
			MultiTargetRotation: true,
			Aoe: &proto.Mage_Rotation_AoeRotation{
				Rotation: proto.Mage_Rotation_AoeRotation_ArcaneExplosion,
			},
		},
	},
}

var FullRaidBuffs = &proto.RaidBuffs{
	GiftOfTheWild: proto.TristateEffect_TristateEffectImproved,
}
var FullFirePartyBuffs = &proto.PartyBuffs{
	Drums:           proto.Drums_DrumsOfBattle,
	Bloodlust:       1,
	MoonkinAura:     proto.TristateEffect_TristateEffectRegular,
	ManaSpringTotem: proto.TristateEffect_TristateEffectRegular,
	TotemOfWrath:    1,
	WrathOfAirTotem: proto.TristateEffect_TristateEffectRegular,
}
var FullFrostPartyBuffs = FullFirePartyBuffs
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:  true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
}

var FullArcanePartyBuffs = &proto.PartyBuffs{
	Drums:           proto.Drums_DrumsOfBattle,
	Bloodlust:       1,
	MoonkinAura:     proto.TristateEffect_TristateEffectRegular,
	ManaSpringTotem: proto.TristateEffect_TristateEffectImproved,
	ManaTideTotems:  1,
	WrathOfAirTotem: proto.TristateEffect_TristateEffectRegular,
}
var FullArcaneIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:  true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
	Innervates:       1,
}

var FullFireConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfPureDeath,
	Food:            proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:   proto.Potions_SuperManaPotion,
	DefaultConjured: proto.Conjured_ConjuredFlameCap,
	MainHandImbue:   proto.WeaponImbue_WeaponImbueBrilliantWizardOil,
}
var FullFrostConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfPureDeath,
	Food:            proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:   proto.Potions_SuperManaPotion,
	DefaultConjured: proto.Conjured_ConjuredMageManaEmerald,
	MainHandImbue:   proto.WeaponImbue_WeaponImbueBrilliantWizardOil,
}

var FullArcaneConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfBlindingLight,
	Food:            proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:   proto.Potions_SuperManaPotion,
	DefaultConjured: proto.Conjured_ConjuredMageManaEmerald,
	MainHandImbue:   proto.WeaponImbue_WeaponImbueBrilliantWizardOil,
}

var FullDebuffs = &proto.Debuffs{
	CurseOfElements:           proto.TristateEffect_TristateEffectImproved,
	ImprovedSealOfTheCrusader: true,
	JudgementOfWisdom:         true,
	Misery:                    true,
}

var P1FireGear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 29076,
		"enchant": 29191,
		"gems": [
			34220,
			24056
		]
	},
	{
		"id": 28134
	},
	{
		"id": 29079,
		"enchant": 28886,
		"gems": [
			31867,
			24030
		]
	},
	{
		"id": 28766,
		"enchant": 33150
	},
	{
		"id": 21848,
		"enchant": 24003,
		"gems": [
			31867,
			31867
		]
	},
	{
		"id": 28411,
		"enchant": 22534,
		"gems": [
			31867
		]
	},
	{
		"id": 21847,
		"enchant": 28272,
		"gems": [
			31867,
			24056
		]
	},
	{
		"id": 21846,
		"gems": [
			31867,
			31867
		]
	},
	{
		"id": 24262,
		"enchant": 24274,
		"gems": [
			31867,
			31867,
			31867
		]
	},
	{
		"id": 28517,
		"enchant": 35297,
		"gems": [
			31867,
			31867
		]
	},
	{
		"id": 28793,
		"enchant": 22536
	},
	{
		"id": 29172,
		"enchant": 22536
	},
	{
		"id": 29370
	},
	{
		"id": 27683
	},
	{
		"id": 28802,
		"enchant": 22560
	},
	{
		"id": 29270
	},
	{
		"id": 28673
	}
]}`)
var P1FrostGear = P1FireGear
var P1ArcaneGear = P1FireGear
