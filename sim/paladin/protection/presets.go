package protection

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var defaultProtTalents = &proto.PaladinTalents{
	// Redoubt:                       5,
	// Precision:                     3,
	// Toughness:                     5,
	// BlessingOfKings:               true,
	// ImprovedRighteousFury:         3,
	// Anticipation:                  5,
	// BlessingOfSanctuary:           true,
	// Reckoning:                     4,
	// SacredDuty:                    2,
	// OneHandedWeaponSpecialization: 5,
	// HolyShield:                    true,
	// ImprovedHolyShield:            2,
	// CombatExpertise:               5,
	// AvengersShield:                true,

	// Benediction:       5,
	// ImprovedJudgement: 2,
	// Deflection:        5,
	// PursuitOfJustice:  3,
	// Crusade:           3,
}

var defaultProtRotation = &proto.ProtectionPaladin_Rotation{}

var defaultProtOptions = &proto.ProtectionPaladin_Options{
	Aura: proto.PaladinAura_RetributionAura,
}

var DefaultOptions = &proto.Player_ProtectionPaladin{
	ProtectionPaladin: &proto.ProtectionPaladin{
		Talents:  defaultProtTalents,
		Options:  defaultProtOptions,
		Rotation: defaultProtRotation,
	},
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance:   true,
	GiftOfTheWild:      proto.TristateEffect_TristateEffectImproved,
	PowerWordFortitude: proto.TristateEffect_TristateEffectRegular,
	MoonkinAura:        proto.TristateEffect_TristateEffectRegular,
	TotemOfWrath:       true,
	WrathOfAirTotem:    true,
	ManaSpringTotem:    proto.TristateEffect_TristateEffectRegular,
}
var FullPartyBuffs = &proto.PartyBuffs{}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:     true,
	BlessingOfSanctuary: true,
	BlessingOfWisdom:    proto.TristateEffect_TristateEffectImproved,
	BlessingOfMight:     proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfBlindingLight,
	Food:            proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:   proto.Potions_SuperManaPotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var FullDebuffs = &proto.Debuffs{
	Misery:          true,
	CurseOfElements: true,
	BloodFrenzy:     true,
	SunderArmor:     true,
	FaerieFire:      proto.TristateEffect_TristateEffectImproved,
	CurseOfWeakness: proto.TristateEffect_TristateEffectImproved,
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
]}`)
var Phase4Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 34401,
		"enchant": 29191,
		"gems": [
			35501,
			32200
		]
	},
	{
		"id": 34178
	},
	{
		"id": 30998,
		"enchant": 23549,
		"gems": [
			32200,
			32196
		]
	},
	{
		"id": 34190,
		"enchant": 35756
	},
	{
		"id": 34945,
		"enchant": 27957,
		"gems": [
			32223
		]
	},
	{
		"id": 34433,
		"enchant": 22533,
		"gems": [
			32200,
			0
		]
	},
	{
		"id": 30985,
		"enchant": 33153,
		"gems": [
			32215,
			0
		]
	},
	{
		"id": 34488,
		"gems": [
			32200,
			0
		]
	},
	{
		"id": 34382,
		"enchant": 24274,
		"gems": [
			32200,
			32200,
			32215
		]
	},
	{
		"id": 34947,
		"gems": [
			32215
		]
	},
	{
		"id": 34889,
		"enchant": 22536
	},
	{
		"id": 34889,
		"enchant": 22536
	},
	{
		"id": 33829
	},
	{
		"id": 35014,
		"enchant": 22555
	},
	{
		"id": 34185,
		"enchant": 28282,
		"gems": [
			32215
		]
	},
	{
		"id": 33504
	}
]}`)
