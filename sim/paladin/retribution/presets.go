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
	UseDivinePlea:        true,
	DivinePleaPercentage: 0.75,
	HolyWrathThreshold:   4,
}

var defaultRetOptions = &proto.RetributionPaladin_Options{
	Judgement:            proto.PaladinJudgement_JudgementOfWisdom,
	Seal:                 proto.PaladinSeal_Vengeance,
	Aura:                 proto.PaladinAura_RetributionAura,
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
	JudgementOfLight:  true,
	Misery:            true,
	CurseOfElements:   true,
	BloodFrenzy:       true,
	SunderArmor:       true,
	FaerieFire:        proto.TristateEffect_TristateEffectImproved,
	CurseOfWeakness:   proto.TristateEffect_TristateEffectImproved,
}

var Phase1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40576,
		"enchant": 44879,
		"gems": [
			41398,
			40037
		]
	},
	{
		"id": 44664,
		"gems": [
			42142
		]
	},
	{
		"id": 40578,
		"enchant": 44871,
		"gems": [
			49110
		]
	},
	{
		"id": 40403,
		"enchant": 55002
	},
	{
		"id": 40574,
		"enchant": 44489,
		"gems": [
			42142,
			39996
		]
	},
	{
		"id": 40186,
		"enchant": 44484,
		"gems": [
			0
		]
	},
	{
		"id": 40541,
		"enchant": 54999,
		"gems": [
			0
		]
	},
	{
		"id": 40205,
		"gems": [
			39996
		]
	},
	{
		"id": 40577,
		"enchant": 38374,
		"gems": [
			42142,
			40038
		]
	},
	{
		"id": 39701,
		"enchant": 55016
	},
	{
		"id": 40075
	},
	{
		"id": 40474
	},
	{
		"id": 42987
	},
	{
		"id": 40431
	},
	{
		"id": 40384,
		"enchant": 44492
	},
	{},
	{
		"id": 42852
	}
]}`)
