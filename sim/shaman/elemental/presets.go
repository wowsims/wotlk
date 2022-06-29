package elemental

import (
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
)

var StandardTalents = &proto.ShamanTalents{
	Convection:         5,
	Concussion:         5,
	ElementalFocus:     true,
	CallOfThunder:      5,
	ElementalFury:      true,
	UnrelentingStorm:   3,
	ElementalPrecision: 3,
	LightningMastery:   5,
	ElementalMastery:   true,
	LightningOverload:  5,
	TotemOfWrath:       true,

	TotemicFocus:    5,
	NaturesGuidance: 3,
	TidalMastery:    5,
}

var eleShamOptionsNoBuffs = &proto.ElementalShaman_Options{
	WaterShield: true,
}

var NoTotems = &proto.ShamanTotems{}
var BasicTotems = &proto.ShamanTotems{
	Earth: proto.EarthTotem_TremorTotem,
	Air:   proto.AirTotem_WrathOfAirTotem,
	Water: proto.WaterTotem_ManaSpringTotem,
	Fire:  proto.FireTotem_TotemOfWrath,
}

var PlayerOptionsCLOnClearcastNoBuffs = &proto.Player_ElementalShaman{
	ElementalShaman: &proto.ElementalShaman{
		Talents: StandardTalents,
		Options: eleShamOptionsNoBuffs,
		Rotation: &proto.ElementalShaman_Rotation{
			Totems: NoTotems,
			Type:   proto.ElementalShaman_Rotation_CLOnClearcast,
		},
	},
}

var eleShamOptions = &proto.ElementalShaman_Options{
	WaterShield: true,
	Bloodlust:   true,
}
var PlayerOptionsAdaptive = &proto.Player_ElementalShaman{
	ElementalShaman: &proto.ElementalShaman{
		Talents: StandardTalents,
		Options: eleShamOptions,
		Rotation: &proto.ElementalShaman_Rotation{
			Totems: BasicTotems,
			Type:   proto.ElementalShaman_Rotation_Adaptive,
		},
	},
}

var PlayerOptionsLBOnly = &proto.Player_ElementalShaman{
	ElementalShaman: &proto.ElementalShaman{
		Talents: StandardTalents,
		Options: eleShamOptions,
		Rotation: &proto.ElementalShaman_Rotation{
			Totems: BasicTotems,
			Type:   proto.ElementalShaman_Rotation_LBOnly,
		},
	},
}

var PlayerOptionsFixed3LBCL = &proto.Player_ElementalShaman{
	ElementalShaman: &proto.ElementalShaman{
		Talents: StandardTalents,
		Options: eleShamOptions,
		Rotation: &proto.ElementalShaman_Rotation{
			Totems:   BasicTotems,
			Type:     proto.ElementalShaman_Rotation_FixedLBCL,
			LbsPerCl: 3,
		},
	},
}

var PlayerOptionsCLOnClearcast = &proto.Player_ElementalShaman{
	ElementalShaman: &proto.ElementalShaman{
		Talents: StandardTalents,
		Options: eleShamOptions,
		Rotation: &proto.ElementalShaman_Rotation{
			Totems: BasicTotems,
			Type:   proto.ElementalShaman_Rotation_CLOnClearcast,
		},
	},
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance: true,
	GiftOfTheWild:    proto.TristateEffect_TristateEffectImproved,
}
var FullPartyBuffs = &proto.PartyBuffs{
	MoonkinAura: proto.TristateEffect_TristateEffectRegular,
}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:  true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
	ShadowPriestDps:  500,
}

var FullConsumes = &proto.Consumes{
	Flask:              proto.Flask_FlaskOfBlindingLight,
	Food:               proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:      proto.Potions_SuperManaPotion,
	StartingPotion:     proto.Potions_DestructionPotion,
	MainHandImbue:      proto.WeaponImbue_WeaponImbueBrilliantWizardOil,
	NumStartingPotions: 1,
	DefaultConjured:    proto.Conjured_ConjuredDarkRune,
	Drums:              proto.Drums_DrumsOfBattle,
}

var FullDebuffs = &proto.Debuffs{
	ImprovedSealOfTheCrusader: true,
	JudgementOfWisdom:         true,
	Misery:                    true,
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 29035,
		"enchant": 29191,
		"gems": [
			34220,
			24059
		]
	},
	{
		"id": 28762
	},
	{
		"id": 29037,
		"enchant": 28886,
		"gems": [
			24059,
			24059
		]
	},
	{
		"id": 28797,
		"enchant": 33150
	},
	{
		"id": 29519,
		"enchant": 24003,
		"gems": [
			24030,
			24030,
			24030
		]
	},
	{
		"id": 29521,
		"enchant": 22534,
		"gems": [
			24059
		]
	},
	{
		"id": 28780,
		"enchant": 28272,
		"gems": [
			24059,
			24056
		]
	},
	{
		"id": 29520,
		"gems": [
			24056,
			24059
		]
	},
	{
		"id": 24262,
		"enchant": 24274,
		"gems": [
			24030,
			24030,
			24030
		]
	},
	{
		"id": 28517,
		"enchant": 35297,
		"gems": [
			24030,
			24030
		]
	},
	{
		"id": 30667,
		"enchant": 22536
	},
	{
		"id": 28753,
		"enchant": 22536
	},
	{
		"id": 29370
	},
	{
		"id": 28785
	},
	{
		"id": 28770,
		"enchant": 22555
	},
	{
		"id": 29273
	},
	{
		"id": 28248
	}
]}`)
