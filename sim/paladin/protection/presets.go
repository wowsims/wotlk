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
    GiftOfTheWild: proto.TristateEffect_TristateEffectImproved,
    PowerWordFortitude: proto.TristateEffect_TristateEffectImproved,
    StrengthOfEarthTotem: proto.TristateEffect_TristateEffectImproved,
    ArcaneBrilliance: true,
    UnleashedRage: true,
    LeaderOfThePack: proto.TristateEffect_TristateEffectRegular,
    IcyTalons: true,
    TotemOfWrath: true,
    DemonicPact: 500,
    SwiftRetribution: true,
    MoonkinAura: proto.TristateEffect_TristateEffectRegular,
    SanctifiedRetribution: true,
    ManaSpringTotem: proto.TristateEffect_TristateEffectRegular,
    Bloodlust: true,
    Thorns: proto.TristateEffect_TristateEffectImproved,
    DevotionAura: proto.TristateEffect_TristateEffectImproved,
    ShadowProtection: true,
}
var FullPartyBuffs = &proto.PartyBuffs{}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:     true,
	BlessingOfSanctuary: true,
	BlessingOfWisdom:    proto.TristateEffect_TristateEffectImproved,
	BlessingOfMight:     proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfStoneblood,
	Food:            proto.Food_FoodDragonfinFilet,
	DefaultPotion:   proto.Potions_IndestructiblePotion,
	PrepopPotion:	 proto.Potions_IndestructiblePotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var FullDebuffs = &proto.Debuffs{
    JudgementOfWisdom: true,
    JudgementOfLight: true,
    Misery: true,
    FaerieFire: proto.TristateEffect_TristateEffectImproved,
    EbonPlaguebringer: true,
    TotemOfWrath: true,
    ShadowMastery: true,
    BloodFrenzy: true,
    Mangle: true,
    ExposeArmor: true,
    SunderArmor: true,
    Vindication: true,
    ThunderClap: proto.TristateEffect_TristateEffectImproved,
    InsectSwarm: true,
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
        {
          "id": 42549,
          "enchant": 44878,
          "gems": [
            41396,
            40089
          ]
        },
        {
          "id": 43282
        },
        {
          "id": 37635,
          "enchant": 44957,
          "gems": [
            40089
          ]
        },
        {
          "id": 44188,
          "enchant": 55002
        },
        {
          "id": 30991,
          "enchant": 47766,
          "gems": [
            40039,
            40039,
            40089
          ]
        },
        {
          "id": 37682,
          "enchant": 44944,
          "gems": [
            0
          ]
        },
        {
          "id": 44183,
          "enchant": 63770,
          "gems": [
            0
          ]
        },
        {
          "id": 37379,
          "enchant": 54793,
          "gems": [
            40022,
            40008
          ]
        },
        {
          "id": 37292,
          "enchant": 38373,
          "gems": [
            40089
          ]
        },
        {
          "id": 44243,
          "enchant": 44528
        },
        {
          "id": 37186,
          "enchant": 59636
        },
        {
          "id": 29297,
          "enchant": 59636
        },
        {
          "id": 40767
        },
        {
          "id": 37220
        },
        {
          "id": 37179,
          "enchant": 22559
        },
        {
          "id": 43085,
          "enchant": 44936
        },
        {
          "id": 40707
        }
      ]}`)
