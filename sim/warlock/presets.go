package warlock

import (
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
)

var defaultDestroTalents = &proto.WarlockTalents{
	// destro
	ImprovedShadowBolt: 5,
	Bane:               5,
	Devastation:        5,
	Shadowburn:         true,
	DestructiveReach:   2,
	ImprovedImmolate:   5,
	Ruin:               true,
	Emberstorm:         5,
	Backlash:           3,
	Conflagrate:        true,
	ShadowAndFlame:     5,
	//demo
	DemonicEmbrace:     5,
	ImprovedVoidwalker: 1,
	FelIntellect:       3,
	// FelDomination: true,
	// MasterSummoner: 2,
	FelStamina:       3,
	DemonicAegis:     3,
	DemonicSacrifice: true,
}

var defaultDestroRotation = &proto.Warlock_Rotation{
	PrimarySpell: proto.Warlock_Rotation_Shadowbolt,
	Immolate:     true,
}

var defaultDestroOptions = &proto.Warlock_Options{
	Armor:           proto.Warlock_Options_FelArmor,
	Summon:          proto.Warlock_Options_Succubus,
	SacrificeSummon: true,
}

var DefaultDestroWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Talents:  defaultDestroTalents,
		Options:  defaultDestroOptions,
		Rotation: defaultDestroRotation,
	},
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance: true,
	GiftOfTheWild:    proto.TristateEffect_TristateEffectImproved,
	DivineSpirit:     proto.TristateEffect_TristateEffectImproved,
}

var FullPartyBuffs = &proto.PartyBuffs{
	Bloodlust:       1,
	Drums:           proto.Drums_DrumsOfBattle,
	ManaSpringTotem: proto.TristateEffect_TristateEffectRegular,
	WrathOfAirTotem: proto.TristateEffect_TristateEffectRegular,
	TotemOfWrath:    1,
}

var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:     true,
	BlessingOfSalvation: true,
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfPureDeath,
	DefaultPotion:   proto.Potions_DestructionPotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
	Food:            proto.Food_FoodBlackenedBasilisk,
}

var FullDebuffs = &proto.Debuffs{
	JudgementOfWisdom:           true,
	Misery:                      true,
	BloodFrenzy:                 true,
	ExposeArmor:                 proto.TristateEffect_TristateEffectImproved,
	FaerieFire:                  proto.TristateEffect_TristateEffectImproved,
	CurseOfRecklessness:         true,
	HuntersMark:                 proto.TristateEffect_TristateEffectImproved,
	ExposeWeaknessUptime:        1,
	ExposeWeaknessHunterAgility: 800,
}

var Phase4Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 31051,
		"enchant": 29191,
		"gems": [
		34220,
		32218
		]
	},
	{
		"id": 33281
	},
	{
		"id": 31054,
		"enchant": 28886,
		"gems": [
		32215,
		32218
		]
	},
	{
		"id": 32524,
		"enchant": 33150
	},
	{
		"id": 30107,
		"enchant": 24003,
		"gems": [
		32196,
		32196,
		32196
		]
	},
	{
		"id": 32586,
		"enchant": 22534
	},
	{
		"id": 31050,
		"enchant": 28272,
		"gems": [
		32196
		]
	},
	{
		"id": 30888,
		"gems": [
		32196,
		32196
		]
	},
	{
		"id": 31053,
		"enchant": 24274,
		"gems": [
		32196
		]
	},
	{
		"id": 32239,
		"enchant": 35297,
		"gems": [
		32218,
		32215
		]
	},
	{
		"id": 32527,
		"enchant": 22536
	},
	{
		"id": 33497,
		"enchant": 22536
	},
	{
		"id": 32483
	},
	{
		"id": 33829
	},
	{
		"id": 32374,
		"enchant": 22561
	},
	{},
	{
		"id": 33192,
		"gems": [
		32215
		]
	}
]}`)
