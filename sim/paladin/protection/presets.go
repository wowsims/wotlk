package protection

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var defaultProtTalents = &proto.PaladinTalents{
	DivineStrength:                5,
	Anticipation:                  5,
	DivineSacrifice:               true,
	ImprovedRighteousFury:         3,
	Toughness:                     5,
	DivineGuardian:                2,
	BlessingOfSanctuary:           true,
	Reckoning:                     3,
	SacredDuty:                    2,
	OneHandedWeaponSpecialization: 3,
	SpiritualAttunement:           1,
	HolyShield:                    true,
	ArdentDefender:                3,
	Redoubt:                       3,
	CombatExpertise:               3,
	TouchedByTheLight:             3,
	AvengersShield:                true,
	GuardedByTheLight:             2,
	ShieldOfTheTemplar:            3,
	JudgementsOfTheJust:           2,
	HammerOfTheRighteous:          true,
	Deflection:                    5,
	Benediction:                   1,
	ImprovedJudgements:            1,
	HeartOfTheCrusader:            3,
	Vindication:                   2,
	SealOfCommand:                 true,
	PursuitOfJustice:              2,
	Crusade:                       3,
}

var defaultProtRotation = &proto.ProtectionPaladin_Rotation{}

var defaultProtOptions = &proto.ProtectionPaladin_Options{
	Judgement: proto.PaladinJudgement_JudgementOfWisdom,
	Seal:      proto.PaladinSeal_Vengeance,
	Aura:      proto.PaladinAura_RetributionAura,
}

var DefaultOptions = &proto.Player_ProtectionPaladin{
	ProtectionPaladin: &proto.ProtectionPaladin{
		Talents:  defaultProtTalents,
		Options:  defaultProtOptions,
		Rotation: defaultProtRotation,
	},
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfStoneblood,
	Food:            proto.Food_FoodDragonfinFilet,
	DefaultPotion:   proto.Potions_IndestructiblePotion,
	PrepopPotion:    proto.Potions_IndestructiblePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40581,
		"enchant": 3818,
		"gems": [
			41396,
			36767
		]
	},
	{
		"id": 40387
	},
	{
		"id": 40584,
		"enchant": 3852,
		"gems": [
			49110
		]
	},
	{
		"id": 40410,
		"enchant": 3605
	},
	{
		"id": 40579,
		"enchant": 3832,
		"gems": [
			36767,
			40022
		]
	},
	{
		"id": 39764,
		"enchant": 3850,
		"gems": [
			0
		]
	},
	{
		"id": 40580,
		"enchant": 3860,
		"gems": [
			40008,
			0
		]
	},
	{
		"id": 39759,
		"enchant": 3601,
		"gems": [
			40008,
			40008
		]
	},
	{
		"id": 40589,
		"enchant": 3822
	},
	{
		"id": 39717,
		"enchant": 3606,
		"gems": [
			40089
		]
	},
	{
		"id": 40718
	},
	{
		"id": 40107
	},
	{
		"id": 44063,
		"gems": [
			36767,
			40089
		]
	},
	{
		"id": 37220
	},
	{
		"id": 40345,
		"enchant": 3788
	},
	{
		"id": 40400,
		"enchant": 3849
	},
	{
		"id": 40707
	}
]}`)
