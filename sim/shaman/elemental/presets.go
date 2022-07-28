package elemental

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var StandardTalents = &proto.ShamanTalents{
	Convection:         5,
	Concussion:         5,
	ElementalFocus:     true,
	CallOfThunder:      true,
	ElementalFury:      5,
	UnrelentingStorm:   3,
	ElementalPrecision: 3,
	LightningMastery:   5,
	ElementalMastery:   true,
	LightningOverload:  5,
	TotemOfWrath:       true,
}

var eleShamOptionsNoBuffs = &proto.ElementalShaman_Options{
	Shield: proto.ShamanShield_WaterShield,
}

var NoTotems = &proto.ShamanTotems{}
var BasicTotems = &proto.ShamanTotems{
	Earth: proto.EarthTotem_TremorTotem,
	Air:   proto.AirTotem_WrathOfAirTotem,
	Water: proto.WaterTotem_ManaSpringTotem,
	Fire:  proto.FireTotem_TotemOfWrath,
}

var eleShamOptions = &proto.ElementalShaman_Options{
	Shield:    proto.ShamanShield_WaterShield,
	Bloodlust: true,
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

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance: true,
	GiftOfTheWild:    proto.TristateEffect_TristateEffectImproved,
	MoonkinAura:      proto.TristateEffect_TristateEffectRegular,
}
var FullPartyBuffs = &proto.PartyBuffs{}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:  true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfBlindingLight,
	Food:            proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:   proto.Potions_SuperManaPotion,
	PrepopPotion:    proto.Potions_DestructionPotion,
	MainHandImbue:   proto.WeaponImbue_WeaponImbueBrilliantWizardOil,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var FullDebuffs = &proto.Debuffs{
	JudgementOfWisdom: true,
	Misery:            true,
}

var P1Gear = items.EquipmentSpecFromJsonString(`{
"items": [
	{
	  "id": 40516,
	  "enchant": 44877,
	  "gems": [
		41285,
		40025
	  ]
	},
	{
	  "id": 44661,
	  "gems": [
		40027
	  ]
	},
	{
	  "id": 40518,
	  "enchant": 44874,
	  "gems": [
		39998
	  ]
	},
	{
	  "id": 44005,
	  "enchant": 44472,
	  "gems": [
		40025
	  ]
	},
	{
	  "id": 40514,
	  "enchant": 44623,
	  "gems": [
		39998,
		40025
	  ]
	},
	{
	  "id": 40324,
	  "enchant": 44498,
	  "gems": [
		40025,
		0
	  ]
	},
	{
	  "id": 40302,
	  "enchant": 54999,
	  "gems": [
		0
	  ]
	},
	{
	  "id": 40327,
	  "gems": [
		39998
	  ]
	},
	{
	  "id": 40517,
	  "enchant": 41602,
	  "gems": [
		40049,
		40027
	  ]
	},
	{
	  "id": 40237,
	  "enchant": 60623,
	  "gems": [
		40025
	  ]
	},
	{
	  "id": 40399
	},
	{
	  "id": 48957
	},
	{
	  "id": 40255
	},
	{
	  "id": 39229
	},
	{
	  "id": 40395,
	  "enchant": 44487
	},
	{
	  "id": 40401,
	  "enchant": 60653
	},
	{
	  "id": 40708
	}
  ]
}`)
