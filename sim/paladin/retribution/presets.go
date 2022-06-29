package retribution

import (
	"github.com/wowsims/tbc/sim/core/items"
	"github.com/wowsims/tbc/sim/core/proto"
)

var defaultRetTalents = &proto.PaladinTalents{
	Benediction:                   5,
	ImprovedSealOfTheCrusader:     3,
	ImprovedJudgement:             2,
	Conviction:                    5,
	SealOfCommand:                 true,
	Crusade:                       3,
	SanctityAura:                  true,
	TwoHandedWeaponSpecialization: 3,
	ImprovedSanctityAura:          2,
	Vengeance:                     5,
	SanctifiedJudgement:           2,
	SanctifiedSeals:               3,
	Fanaticism:                    5,
	CrusaderStrike:                true,
	Precision:                     3,
	DivineStrength:                5,
}

var defaultRetRotation = &proto.RetributionPaladin_Rotation{
	ConsecrationRank: proto.RetributionPaladin_Rotation_None,
	UseExorcism:      false,
}

var defaultRetOptions = &proto.RetributionPaladin_Options{
	Judgement:             proto.RetributionPaladin_Options_Crusader,
	Aura:                  proto.PaladinAura_SanctityAura,
	CrusaderStrikeDelayMs: 1700,
	HasteLeewayMs:         100,
	DamageTakenPerSecond:  0,
}

var DefaultOptions = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Talents:  defaultRetTalents,
		Options:  defaultRetOptions,
		Rotation: defaultRetRotation,
	},
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance: true,
	GiftOfTheWild:    proto.TristateEffect_TristateEffectImproved,
	DivineSpirit:     proto.TristateEffect_TristateEffectImproved,
}

var FullPartyBuffs = &proto.PartyBuffs{
	Bloodlust:            1,
	Drums:                proto.Drums_DrumsOfBattle,
	BraidedEterniumChain: true,
	ManaSpringTotem:      proto.TristateEffect_TristateEffectRegular,
	StrengthOfEarthTotem: proto.StrengthOfEarthType_EnhancingTotems,
	WindfuryTotemRank:    5,
	BattleShout:          proto.TristateEffect_TristateEffectImproved,
	WindfuryTotemIwt:     2,
}

var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:     true,
	BlessingOfMight:     proto.TristateEffect_TristateEffectImproved,
	BlessingOfSalvation: true,
	UnleashedRage:       true,
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfRelentlessAssault,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
	Food:            proto.Food_FoodRoastedClefthoof,
	SuperSapper:     true,
}

var FullDebuffs = &proto.Debuffs{
	JudgementOfWisdom:           true,
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

var Phase4Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 32235,
		"enchant": 29192,
		"gems": [
			32409,
			32193
		]
	},
	{
		"id": 30022
	},
	{
		"id": 30055,
		"enchant": 28888,
		"gems": [
			32193
		]
	},
	{
		"id": 33590,
		"enchant": 34004
	},
	{
		"id": 30905,
		"enchant": 24003,
		"gems": [
			32211,
			32193,
			32217
		]
	},
	{
		"id": 32574,
		"enchant": 27899
	},
	{
		"id": 29947,
		"enchant": 33995
	},
	{
		"id": 30106,
		"gems": [
			32193,
			32211
		]
	},
	{
		"id": 30900,
		"enchant": 29535,
		"gems": [
			32193,
			32193,
			32193
		]
	},
	{
		"id": 32366,
		"enchant": 22544,
		"gems": [
			32193,
			32217
		]
	},
	{
		"id": 32526
	},
	{
		"id": 30834
	},
	{
		"id": 33831
	},
	{
		"id": 28830
	},
	{
		"id": 32332,
		"enchant": 22559
	},
	{
		"id": 27484
	}
]}`)
