package retribution

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var defaultRetTalents = &proto.PaladinTalents{
	SealsOfThePure:                5,
	DivineIntellect:               5,
	AuraMastery:                   true,
	DivineStrength:                5,
	Benediction:                   5,
	ImprovedJudgements:            2,
	HeartOfTheCrusader:            3,
	ImprovedBlessingOfMight:       2,
	Conviction:                    5,
	SealOfCommand:                 true,
	PursuitOfJustice:              2,
	SanctityOfBattle:              3,
	Crusade:                       3,
	TwoHandedWeaponSpecialization: 3,
	SanctifiedRetribution:         true,
	Vengeance:                     3,
	TheArtOfWar:                   2,
	Repentance:                    true,
	JudgementsOfTheWise:           3,
	Fanaticism:                    3,
	SanctifiedWrath:               2,
	SwiftRetribution:              3,
	CrusaderStrike:                true,
	SanctifiedLight:               3,
	RighteousVengeance:            3,
	DivineStorm:                   true,
}

var defaultRetGlyphs = &proto.Glyphs{
	Major1: int32(proto.PaladinMajorGlyph_GlyphOfSealOfVengeance),
	Major2: int32(proto.PaladinMajorGlyph_GlyphOfJudgement),
	Major3: int32(proto.PaladinMajorGlyph_GlyphOfConsecration),
	Minor1: int32(proto.PaladinMinorGlyph_GlyphOfSenseUndead),
	Minor2: int32(proto.PaladinMinorGlyph_GlyphOfLayOnHands),
	Minor3: int32(proto.PaladinMinorGlyph_GlyphOfBlessingOfKings),
}

var defaultRetRotation = &proto.RetributionPaladin_Rotation{
	ConsSlack:            500,
	ExoSlack:             500,
	DivinePleaPercentage: 0.75,
}

var defaultRetOptions = &proto.RetributionPaladin_Options{
	Judgement:            proto.PaladinJudgement_JudgementOfWisdom,
	Seal:                 proto.PaladinSeal_Vengeance,
	Aura:                 proto.PaladinAura_RetributionAura,
	UseDivinePlea:        true,
	DamageTakenPerSecond: 0,
}

var DefaultOptions = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Talents:  defaultRetTalents,
		Options:  defaultRetOptions,
		Rotation: defaultRetRotation,
	},
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance:     true,
	GiftOfTheWild:        proto.TristateEffect_TristateEffectImproved,
	DivineSpirit:         true,
	Bloodlust:            true,
	ManaSpringTotem:      proto.TristateEffect_TristateEffectRegular,
	StrengthOfEarthTotem: proto.TristateEffect_TristateEffectImproved,
	WindfuryTotem:        proto.TristateEffect_TristateEffectImproved,
	BattleShout:          proto.TristateEffect_TristateEffectImproved,
	UnleashedRage:        true,
}

var FullPartyBuffs = &proto.PartyBuffs{
	BraidedEterniumChain: true,
}

var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings: true,
	BlessingOfMight: proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfRelentlessAssault,
	DefaultPotion:   proto.Potions_HastePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
	Food:            proto.Food_FoodRoastedClefthoof,
	ThermalSapper:   true,
}

var FullDebuffs = &proto.Debuffs{
	JudgementOfWisdom: true,
	Misery:            true,
	CurseOfElements:   true,
	BloodFrenzy:       true,
	SunderArmor:       true,
	FaerieFire:        proto.TristateEffect_TristateEffectImproved,
	CurseOfWeakness:   proto.TristateEffect_TristateEffectImproved,
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
