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

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
        {
          "id": 40581,
          "enchant": 44878,
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
          "enchant": 44957,
          "gems": [
            49110
          ]
        },
        {
          "id": 40410,
          "enchant": 55002
        },
        {
          "id": 40579,
          "enchant": 44489,
          "gems": [
            36767,
            40022
          ]
        },
        {
          "id": 39764,
          "enchant": 44944,
          "gems": [
            0
          ]
        },
        {
          "id": 40580,
          "enchant": 63770,
          "gems": [
            40008,
            0
          ]
        },
        {
          "id": 39759,
          "enchant": 54793,
          "gems": [
            40008,
            40008
          ]
        },
        {
          "id": 40589,
          "enchant": 38373
        },
        {
          "id": 39717,
          "enchant": 55016,
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
          "enchant": 44496
        },
        {
          "id": 40400,
          "enchant": 44936
        },
        {
          "id": 40707
        }
      ]}`)
