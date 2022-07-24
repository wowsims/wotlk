package enhancement

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var BasicRaidBuffs = &proto.RaidBuffs{}
var BasicPartyBuffs = &proto.PartyBuffs{}
var BasicIndividualBuffs = &proto.IndividualBuffs{}

var StandardTalents = &proto.ShamanTalents{
	Concussion:           5,
	CallOfFlame:          3,
	ElementalDevastation: 3,
	ElementalFocus:       true,
	ElementalFury:        5,
	ImprovedFireNova:     2,

	EnhancingTotems:         3,
	AncestralKnowledge:      2,
	ThunderingStrikes:       5,
	ImprovedShields:         3,
	ElementalWeapons:        3,
	ShamanisticFocus:        true, //1/2 imp stormstrike might be better, yet to be determined
	Flurry:                  5,
	ImprovedWindfuryTotem:   2,
	SpiritWeapons:           true,
	MentalDexterity:         3,
	UnleashedRage:           3,
	WeaponMastery:           3,
	DualWieldSpecialization: 3,
	DualWield:               true,
	Stormstrike:             true,
	StaticShock:             3,
	LavaLash:                true,
	MentalQuickness:         3,
	ShamanisticRage:         true,
	MaelstromWeapon:         5,
	FeralSpirit:             true,
}

var PlayerOptionsBasic = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Talents:  StandardTalents,
		Options:  enhShamOptions,
		Rotation: enhShamRotation,
	},
}

var enhShamRotation = &proto.EnhancementShaman_Rotation{
	Totems: &proto.ShamanTotems{
		Earth: proto.EarthTotem_StrengthOfEarthTotem,
		Air:   proto.AirTotem_WindfuryTotem,
		Water: proto.WaterTotem_ManaSpringTotem,
		Fire:  proto.FireTotem_MagmaTotem,
	},
}

var enhShamOptions = &proto.EnhancementShaman_Options{
	Shield:             proto.ShamanShield_LightningShield,
	Bloodlust:          true,
	DelayOffhandSwings: true, //might not be default anymore, depending on what imbues we run. still useful regardless
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance:     true,
	GiftOfTheWild:        proto.TristateEffect_TristateEffectImproved,
	FerociousInspiration: true,
	BattleShout:          proto.TristateEffect_TristateEffectImproved,
	LeaderOfThePack:      proto.TristateEffect_TristateEffectImproved,
	TrueshotAura:         true,
}
var FullPartyBuffs = &proto.PartyBuffs{}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:  true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
	BlessingOfMight:  proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	MainHandImbue:   proto.WeaponImbue_WeaponImbueShamanWindfury,
	OffHandImbue:    proto.WeaponImbue_WeaponImbueShamanFlametongue,
	DefaultConjured: proto.Conjured_ConjuredFlameCap,
}

var FullDebuffs = &proto.Debuffs{
	BloodFrenzy:       true,
	SunderArmor:       true,
	FaerieFire:        proto.TristateEffect_TristateEffectImproved,
	JudgementOfWisdom: true,
	Misery:            true,
}

var Phase2Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 30190,
		"enchant": 29192,
		"gems": [
			32409,
			24058
		]
	},
	{
		"id": 30017
	},
	{
		"id": 30055,
		"enchant": 28888,
		"gems": [
			24027
		]
	},
	{
		"id": 29994,
		"enchant": 34004
	},
	{
		"id": 30185,
		"enchant": 24003,
		"gems": [
			24027,
			24054,
			24058
		]
	},
	{
		"id": 30091,
		"enchant": 27899,
		"gems": [
			24027
		]
	},
	{
		"id": 30189,
		"enchant": 33995
	},
	{
		"id": 30106,
		"gems": [
			24027,
			24054
		]
	},
	{
		"id": 30192,
		"enchant": 29535,
		"gems": [
			24027
		]
	},
	{
		"id": 30039,
		"enchant": 28279
	},
	{
		"id": 29997
	},
	{
		"id": 30052
	},
	{
		"id": 28830
	},
	{
		"id": 29383
	},
	{
		"id": 32944,
		"enchant": 22559
	},
	{
		"id": 29996,
		"enchant": 22559
	},
	{
		"id": 27815
	}
]}`)
