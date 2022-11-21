package dps

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var FrostDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.DeathknightMajorGlyph_GlyphOfFrostStrike),
	Major2: int32(proto.DeathknightMajorGlyph_GlyphOfObliterate),
	Major3: int32(proto.DeathknightMajorGlyph_GlyphOfDisease),
	// No interesting minor glyphs.
}

var UnholyDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.DeathknightMajorGlyph_GlyphOfTheGhoul),
	Major2: int32(proto.DeathknightMajorGlyph_GlyphOfDarkDeath),
	Major3: int32(proto.DeathknightMajorGlyph_GlyphOfDeathAndDecay),
	// No interesting minor glyphs.
}

var PlayerOptionsUnholy = &proto.Player_Deathknight{
	Deathknight: &proto.Deathknight{
		Talents:  UnholyTalents,
		Options:  deathKnightOptions,
		Rotation: unholyRotation,
	},
}

var PlayerOptionsFrost = &proto.Player_Deathknight{
	Deathknight: &proto.Deathknight{
		Talents:  FrostTalents,
		Options:  deathKnightOptions,
		Rotation: frostRotation,
	},
}

var UnholyTalents = &proto.DeathknightTalents{
	ImprovedIcyTouch:  3,
	RunicPowerMastery: 2,
	BlackIce:          3,
	NervesOfColdSteel: 3,
	IcyTalons:         5,
	EndlessWinter:     2,

	ViciousStrikes:    2,
	Virulence:         3,
	Morbidity:         3,
	RavenousDead:      3,
	Outbreak:          0,
	Necrosis:          5,
	BloodCakedBlade:   3,
	NightOfTheDead:    2,
	Impurity:          5,
	Dirge:             2,
	MasterOfGhouls:    true,
	Desolation:        5,
	GhoulFrenzy:       true,
	CryptFever:        3,
	BoneShield:        true,
	WanderingPlague:   3,
	EbonPlaguebringer: 3,
	ScourgeStrike:     true,
	RageOfRivendare:   5,
	SummonGargoyle:    true,
}

var FrostTalents = &proto.DeathknightTalents{
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

var unholyRotation = &proto.Deathknight_Rotation{
	UseDeathAndDecay:     true,
	StartingPresence:     proto.Deathknight_Rotation_Unholy,
	BlPresence:           proto.Deathknight_Rotation_Blood,
	UseEmpowerRuneWeapon: true,
	UseGargoyle:          true,
	BtGhoulFrenzy:        true,
	HoldErwArmy:          false,
	BloodRuneFiller:      proto.Deathknight_Rotation_BloodBoil,
	ArmyOfTheDead:        proto.Deathknight_Rotation_AsMajorCd,
	BloodTap:             proto.Deathknight_Rotation_GhoulFrenzy,
}

var frostRotation = &proto.Deathknight_Rotation{}

var deathKnightOptions = &proto.Deathknight_Options{
	StartingRunicPower:  0,
	PetUptime:           1,
	PrecastGhoulFrenzy:  false,
	PrecastHornOfWinter: true,
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfEndlessRage,
	DefaultPotion: proto.Potions_PotionOfSpeed,
	PrepopPotion:  proto.Potions_PotionOfSpeed,
	Food:          proto.Food_FoodDragonfinFilet,
}

var UnholyDwP1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 44006,
		"enchant": 3817,
		"gems": [
			41398,
			42702
		]
	},
	{
		"id": 39421
	},
	{
		"id": 40557,
		"enchant": 3808,
		"gems": [
			39996
		]
	},
	{
		"id": 40403,
		"enchant": 3831
	},
	{
		"id": 40550,
		"enchant": 3832,
		"gems": [
			42142,
			39996
		]
	},
	{
		"id": 40330,
		"enchant": 3845,
		"gems": [
			39996,
			0
		]
	},
	{
		"id": 40347,
		"enchant": 3604,
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
		"enchant": 3823
	},
	{
		"id": 39706,
		"enchant": 3606,
		"gems": [
			39996
		]
	},
	{
		"id": 39401
	},
	{
		"id": 40075
	},
	{
		"id": 37390
	},
	{
		"id": 42987
	},
	{
		"id": 40402,
		"enchant": 3368
	},
	{
		"id": 40491,
		"enchant": 3368
	},
	{
		"id": 42620
	}
]}`)

var FrostP1Gear = core.EquipmentSpecFromJsonString(`{ "items": [
	{
		"id": 44006,
		"enchant": 3817,
		"gems": [
			41398,
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
		"id": 40557,
		"enchant": 3808,
		"gems": [
			39996
		]
	},
	{
		"id": 40403,
		"enchant": 3831
	},
	{
		"id": 40550,
		"enchant": 3832,
		"gems": [
			42142,
			39996
		]
	},
	{
		"id": 40330,
		"enchant": 3845,
		"gems": [
			39996,
			0
		]
	},
	{
		"id": 40552,
		"enchant": 3604,
		"gems": [
			39996,
			0
		]
	},
	{
		"id": 40278,
		"gems": [
			39996,
			42142
		]
	},
	{
		"id": 40556,
		"enchant": 3823,
		"gems": [
			42142,
			39996
		]
	},
	{
		"id": 40591,
		"enchant": 3606
	},
	{
		"id": 39401
	},
	{
		"id": 40075
	},
	{
		"id": 40256
	},
	{
		"id": 42987
	},
	{
		"id": 40189,
		"enchant": 3370
	},
	{
		"id": 40189,
		"enchant": 3368
	},
	{
		"id": 40207
	}
]}`)
