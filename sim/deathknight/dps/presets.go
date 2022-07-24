package dps

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var PlayerOptionsUnholy = &proto.Player_DeathKnight{
	DeathKnight: &proto.DeathKnight{
		Talents:  UnholyTalents,
		Options:  deathKnightOptions,
		Rotation: unholyRotation,
	},
}

var PlayerOptionsFrost = &proto.Player_DeathKnight{
	DeathKnight: &proto.DeathKnight{
		Talents:  FrostTalents,
		Options:  deathKnightOptions,
		Rotation: frostRotation,
	},
}

var UnholyTalents = &proto.DeathKnightTalents{
	ImprovedIcyTouch:  3,
	RunicPowerMastery: 2,
	BlackIce:          2,
	NervesOfColdSteel: 3,
	IcyTalons:         5,
	EndlessWinter:     2,

	ViciousStrikes:    2,
	Virulence:         3,
	Morbidity:         3,
	RavenousDead:      3,
	Outbreak:          3,
	Necrosis:          5,
	BloodCakedBlade:   3,
	NightOfTheDead:    2,
	Impurity:          5,
	Dirge:             2,
	MasterOfGhouls:    true,
	Desolation:        5,
	CryptFever:        3,
	BoneShield:        true,
	WanderingPlague:   3,
	EbonPlaguebringer: 3,
	ScourgeStrike:     true,
	RageOfRivendare:   5,
	SummonGargoyle:    true,
}

var FrostTalents = &proto.DeathKnightTalents{
	Butchery:       2,
	Subversion:     3,
	BladedArmor:    5,
	DarkConviction: 5,

	ImprovedIcyTouch:   3,
	RunicPowerMastery:  2,
	BlackIce:           5,
	NervesOfColdSteel:  3,
	IcyTalons:          5,
	Annihilation:       3,
	KillingMachine:     5,
	ChillOfTheGrave:    2,
	EndlessWinter:      2,
	GlacierRot:         3,
	ImprovedIcyTalons:  true,
	MercilessCombat:    2,
	Rime:               3,
	ThreatOfThassarian: 3,
	BloodOfTheNorth:    3,
	UnbreakableArmor:   true,
	FrostStrike:        true,
	GuileOfGorefiend:   3,
	TundraStalker:      5,
	HowlingBlast:       true,
}

var unholyRotation = &proto.DeathKnight_Rotation{
	UnholyPresenceOpener: true,
}

var frostRotation = &proto.DeathKnight_Rotation{
	UnholyPresenceOpener: false,
}

var deathKnightOptions = &proto.DeathKnight_Options{
	StartingRunicPower:  0,
	PetUptime:           1,
	PrecastGhoulFrenzy:  true,
	PrecastHornOfWinter: true,
}

var FullRaidBuffs = &proto.RaidBuffs{
	GiftOfTheWild:         proto.TristateEffect_TristateEffectImproved,
	SwiftRetribution:      true,
	StrengthOfEarthTotem:  proto.TristateEffect_TristateEffectImproved,
	IcyTalons:             true,
	AbominationsMight:     true,
	LeaderOfThePack:       proto.TristateEffect_TristateEffectImproved,
	SanctifiedRetribution: true,
	Bloodlust:             true,
	DevotionAura:          proto.TristateEffect_TristateEffectImproved,
}
var FullPartyBuffs = &proto.PartyBuffs{
	HeroicPresence: true,
}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings: true,
	BlessingOfMight: proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfEndlessRage,
	DefaultPotion: proto.Potions_PotionOfSpeed,
	PrepopPotion:  proto.Potions_PotionOfSpeed,
	Food:          proto.Food_FoodDragonfinFilet,
}

var FullDebuffs = &proto.Debuffs{
	BloodFrenzy:        true,
	FaerieFire:         proto.TristateEffect_TristateEffectImproved,
	JudgementOfWisdom:  true,
	Misery:             true,
	SunderArmor:        true,
	EbonPlaguebringer:  true,
	HeartOfTheCrusader: true,
}

var UnholyDwP1Gear = items.EquipmentSpecFromJsonString(`{"items": [
{
  "id": 44006,
  "enchant": 44879,
  "gems": [
	41400,
	22459
  ]
},
{
  "id": 44664,
  "gems": [
	39996
  ]
},
{
  "id": 40557,
  "enchant": 44871,
  "gems": [
	39996
  ]
},
{
  "id": 40403,
  "enchant": 44472
},
{
  "id": 40550,
  "enchant": 44623,
  "gems": [
	42142,
	40038
  ]
},
{
  "id": 40330,
  "enchant": 60616,
  "gems": [
	39996,
	0
  ]
},
{
  "id": 40347,
  "enchant": 54999,
  "gems": [
	39996,
	0
  ]
},
{
  "id": 40278,
  "gems": [
	42142,
	42142
  ]
},
{
  "id": 40294,
  "enchant": 38374
},
{
  "id": 40591,
  "enchant": 55016
},
{
  "id": 40717
},
{
  "id": 40075
},
{
  "id": 40431
},
{
  "id": 42987
},
{
  "id": 40189,
  "enchant": 53344
},
{
  "id": 40491,
  "enchant": 44495
},
{
  "id": 40207
}
]}`)

var FrostP1Gear = items.EquipmentSpecFromJsonString(`{ "items": [
{
	"id": 44006,
	"enchant": 44879,
	"gems": [
	41398,
	40022
	]
},
{
	"id": 44664,
	"gems": [
	39996
	]
},
{
	"id": 40557,
	"enchant": 44871,
	"gems": [
	39996
	]
},
{
	"id": 40403,
	"enchant": 55002
},
{
	"id": 40550,
	"enchant": 44623,
	"gems": [
	42142,
	39996
	]
},
{
	"id": 40330,
	"enchant": 60616,
	"gems": [
	39996,
	0
	]
},
{
	"id": 40552,
	"enchant": 54999,
	"gems": [
	39996,
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
	"id": 40556,
	"enchant": 38374,
	"gems": [
	42142,
	39996
	]
},
{
	"id": 40591,
	"enchant": 55016
},
{
	"id": 39401
},
{
	"id": 40075
},
{
	"id": 40684
},
{
	"id": 42987
},
{
	"id": 40189,
	"enchant": 53343
},
{
	"id": 40189,
	"enchant": 53344
},
{
	"id": 40207
}
]}`)
