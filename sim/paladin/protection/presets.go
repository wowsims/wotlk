package protection

import (
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
)

var defaultProtTalents = &proto.PaladinTalents{
	Redoubt:                       5,
	Precision:                     3,
	Toughness:                     5,
	BlessingOfKings:               true,
	ImprovedRighteousFury:         3,
	Anticipation:                  5,
	BlessingOfSanctuary:           true,
	Reckoning:                     4,
	SacredDuty:                    2,
	OneHandedWeaponSpecialization: 5,
	HolyShield:                    true,
	ImprovedHolyShield:            2,
	CombatExpertise:               5,
	AvengersShield:                true,

	Benediction:       5,
	ImprovedJudgement: 2,
	Deflection:        5,
	PursuitOfJustice:  3,
	Crusade:           3,
}

var defaultProtRotation = &proto.ProtectionPaladin_Rotation{
	ConsecrationRank:  6,
	UseExorcism:       false,
	MaintainJudgement: proto.PaladinJudgement_JudgementOfWisdom,
}

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
}
var FullPartyBuffs = &proto.PartyBuffs{
	MoonkinAura:     proto.TristateEffect_TristateEffectRegular,
	TotemOfWrath:    1,
	WrathOfAirTotem: proto.TristateEffect_TristateEffectImproved,
	ManaSpringTotem: proto.TristateEffect_TristateEffectRegular,
}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings: true,
	//BlessingOfSanctuary: true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
	BlessingOfMight:  proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	Flask:              proto.Flask_FlaskOfBlindingLight,
	Food:               proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:      proto.Potions_SuperManaPotion,
	NumStartingPotions: 1,
	DefaultConjured:    proto.Conjured_ConjuredDarkRune,
	MainHandImbue:      proto.WeaponImbue_WeaponImbueSuperiorWizardOil,
}

var FullDebuffs = &proto.Debuffs{
	Misery:                      true,
	CurseOfElements:             proto.TristateEffect_TristateEffectImproved,
	IsbUptime:                   1,
	BloodFrenzy:                 true,
	ExposeArmor:                 proto.TristateEffect_TristateEffectImproved,
	FaerieFire:                  proto.TristateEffect_TristateEffectImproved,
	CurseOfRecklessness:         true,
	HuntersMark:                 proto.TristateEffect_TristateEffectImproved,
	ExposeWeaknessUptime:        1,
	ExposeWeaknessHunterAgility: 800,
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
]}`)
var Phase4Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 32521,
		"enchant": 29191,
		"gems": [
			25896,
			32196
		]
	},
	{
		"id": 32362
	},
	{
		"id": 30998,
		"enchant": 28911,
		"gems": [
			32200,
			32196
		]
	},
	{
		"id": 33593,
		"enchant": 33148
	},
	{
		"id": 30991,
		"enchant": 27957,
		"gems": [
			32196,
			32196,
			32221
		]
	},
	{
		"id": 32232,
		"enchant": 22534
	},
	{
		"id": 30985,
		"enchant": 33153,
		"gems": [
			32196
		]
	},
	{
		"id": 32342,
		"gems": [
			32200,
			32200
		]
	},
	{
		"id": 30995,
		"enchant": 24274,
		"gems": [
			32200
		]
	},
	{
		"id": 32245,
		"enchant": 35297,
		"gems": [
			32200,
			32200
		]
	},
	{
		"id": 32261,
		"enchant": 22536
	},
	{
		"id": 29172,
		"enchant": 22536
	},
	{
		"id": 31858
	},
	{
		"id": 33829
	},
	{
		"id": 30910,
		"enchant": 22555
	},
	{
		"id": 32375,
		"enchant": 28282
	},
	{
		"id": 33504
	}
]}`)
