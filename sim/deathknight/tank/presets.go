package tank

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var PlayerOptionsBloodTank = &proto.Player_TankDeathknight{
	TankDeathknight: &proto.TankDeathknight{
		Talents:  BloodTankTalents,
		Options:  deathKnightOptions,
		Rotation: bloodTankRotation,
	},
}

var BloodTankTalents = &proto.DeathknightTalents{
	BladeBarrier:         5,
	BladedArmor:          5,
	ScentOfBlood:         1,
	RuneTap:              true,
	DarkConviction:       5,
	DeathRuneMastery:     3,
	ImprovedRuneTap:      3,
	SpellDeflection:      3,
	BloodyStrikes:        3,
	VeteranOfTheThirdWar: 3,
	BloodyVengeance:      2,
	AbominationsMight:    2,
	Hysteria:             true,
	ImprovedDeathStrike:  2,
	VampiricBlood:        true,
	WillOfTheNecropolis:  3,
	ImprovedIcyTouch:     3,
	Toughness:            5,
	BlackIce:             5,
	IcyTalons:            5,
	Lichborne:            true,
	EndlessWinter:        2,
	FrigidDreadplate:     3,
	GlacierRot:           2,
	ImprovedIcyTalons:    true,
	Anticipation:         2,
}

var bloodTankRotation = &proto.TankDeathknight_Rotation{}

var deathKnightOptions = &proto.TankDeathknight_Options{
	StartingRunicPower: 0,
}

var FullRaidBuffs = &proto.RaidBuffs{
	GiftOfTheWild:         proto.TristateEffect_TristateEffectImproved,
	PowerWordFortitude:    proto.TristateEffect_TristateEffectImproved,
	AbominationsMight:     true,
	SwiftRetribution:      true,
	Bloodlust:             true,
	StrengthOfEarthTotem:  proto.TristateEffect_TristateEffectImproved,
	LeaderOfThePack:       proto.TristateEffect_TristateEffectImproved,
	SanctifiedRetribution: true,
	DevotionAura:          proto.TristateEffect_TristateEffectImproved,
	RetributionAura:       true,
	IcyTalons:             true,
}
var FullPartyBuffs = &proto.PartyBuffs{
	HeroicPresence: true,
}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:     true,
	BlessingOfMight:     proto.TristateEffect_TristateEffectImproved,
	BlessingOfSanctuary: true,
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfStoneblood,
	DefaultPotion: proto.Potions_IndestructiblePotion,
	PrepopPotion:  proto.Potions_IndestructiblePotion,
	Food:          proto.Food_FoodDragonfinFilet,
}

var FullDebuffs = &proto.Debuffs{
	SunderArmor:        true,
	Mangle:             true,
	DemoralizingShout:  proto.TristateEffect_TristateEffectImproved,
	JudgementOfLight:   true,
	FaerieFire:         proto.TristateEffect_TristateEffectRegular,
	Misery:             true,
	FrostFever:         proto.TristateEffect_TristateEffectImproved,
	BloodFrenzy:        true,
	EbonPlaguebringer:  true,
	HeartOfTheCrusader: true,
}

var Glyphs = &proto.Glyphs{
	Major1: int32(proto.DeathknightMajorGlyph_GlyphOfDarkCommand),
	Major2: int32(proto.DeathknightMajorGlyph_GlyphOfObliterate),
	Major3: int32(proto.DeathknightMajorGlyph_GlyphOfVampiricBlood),
}

var BloodP1Gear = items.EquipmentSpecFromJsonString(`{"items": [
    {
      "id": 40565,
      "enchant": 67839,
      "gems": [
        41380,
        36767
      ]
    },
    {
      "id": 40387
    },
    {
      "id": 39704,
      "enchant": 44957,
      "gems": [
        40008
      ]
    },
    {
      "id": 40252,
      "enchant": 55002
    },
    {
      "id": 40559,
      "gems": [
        40008,
        40022
      ]
    },
    {
      "id": 40306,
      "enchant": 44944,
      "gems": [
        40008,
        0
      ]
    },
    {
      "id": 40563,
      "enchant": 63770,
      "gems": [
        40008,
        0
      ]
    },
    {
      "id": 39759,
      "gems": [
        40008,
        40008
      ]
    },
    {
      "id": 40567,
      "enchant": 38373,
      "gems": [
        40008,
        40008
      ]
    },
    {
      "id": 40297,
      "enchant": 44491
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
        36767
      ]
    },
    {
      "id": 42341,
      "gems": [
        40008,
        40008
      ]
    },
    {
      "id": 40406,
      "enchant": 62158
    },
    {},
    {
      "id": 40207
    }
  ]}`)
