package enhancement

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var BasicRaidBuffs = &proto.RaidBuffs{}
var BasicPartyBuffs = &proto.PartyBuffs{}
var BasicIndividualBuffs = &proto.IndividualBuffs{}

var StandardTalents = &proto.ShamanTalents{
	Concussion:           5,
	CallOfFlame:          3,
	ElementalDevastation: 3,
	ElementalFocus:       true,
	ElementalFury:        5,
	ImprovedFireNova:     2,

	EnhancingTotems:         3,
	AncestralKnowledge:      4,
	ThunderingStrikes:       5,
	ElementalWeapons:        3,
	ShamanisticFocus:        true, //1/2 imp stormstrike might be better, yet to be determined
	Flurry:                  5,
	ImprovedWindfuryTotem:   2,
	SpiritWeapons:           true,
	MentalDexterity:         3,
	UnleashedRage:           3,
	WeaponMastery:           3,
	DualWieldSpecialization: 3,
	DualWield:               true,
	Stormstrike:             true,
	StaticShock:             3,
	LavaLash:                true,
	MentalQuickness:         3,
	ShamanisticRage:         true,
	MaelstromWeapon:         5,
	FeralSpirit:             true,
}

var DefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.ShamanMajorGlyph_GlyphOfStormstrike),
	Major2: int32(proto.ShamanMajorGlyph_GlyphOfFlametongueWeapon),
	Major3: int32(proto.ShamanMajorGlyph_GlyphOfFeralSpirit),

	//minor glyphs are all qol passive stuff
}

var PlayerOptionsBasic = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Talents:  StandardTalents,
		Options:  enhShamOptions,
		Rotation: enhShamRotation,
	},
}

var enhShamRotation = &proto.EnhancementShaman_Rotation{
	Totems: &proto.ShamanTotems{
		Earth: proto.EarthTotem_StrengthOfEarthTotem,
		Air:   proto.AirTotem_WindfuryTotem,
		Water: proto.WaterTotem_ManaSpringTotem,
		Fire:  proto.FireTotem_MagmaTotem,
	},
	WeaveLatency: 200,
}

var enhShamOptions = &proto.EnhancementShaman_Options{
	Shield:    proto.ShamanShield_LightningShield,
	Bloodlust: true,
	SyncType:  proto.ShamanSyncType_SyncMainhandOffhandSwings,
	ImbueMh:   proto.ShamanImbue_FlametongueWeaponDownrank, //phase 1 (wraith strike) only
	ImbueOh:   proto.ShamanImbue_FlametongueWeapon,
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance:     true,
	GiftOfTheWild:        proto.TristateEffect_TristateEffectImproved,
	FerociousInspiration: true,
	BattleShout:          proto.TristateEffect_TristateEffectImproved,
	LeaderOfThePack:      proto.TristateEffect_TristateEffectImproved,
	TrueshotAura:         true,
}
var FullPartyBuffs = &proto.PartyBuffs{}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings:  true,
	BlessingOfWisdom: proto.TristateEffect_TristateEffectImproved,
	BlessingOfMight:  proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	DefaultConjured: proto.Conjured_ConjuredFlameCap,
}

var FullDebuffs = &proto.Debuffs{
	BloodFrenzy:       true,
	SunderArmor:       true,
	FaerieFire:        proto.TristateEffect_TristateEffectImproved,
	JudgementOfWisdom: true,
	Misery:            true,
}

var Phase1Gear = items.EquipmentSpecFromJsonString(`{"items": [
        {
          "id": 40543,
          "enchant": 44879,
          "gems": [
            41398,
            40014
          ]
        },
        {
          "id": 44661,
          "gems": [
            40014
          ]
        },
        {
          "id": 40524,
          "enchant": 44871,
          "gems": [
            40014
          ]
        },
        {
          "id": 40403,
          "enchant": 55002
        },
        {
          "id": 40523,
		  "enchant": 44489,
          "gems": [
            40003,
            40014
          ]
        },
        {
          "id": 40282,
          "enchant": 44484,
          "gems": [
            42702,
            0
          ]
        },
        {
          "id": 40520,
          "enchant": 54999,
          "gems": [
            42154,
            0
          ]
        },
        {
          "id": 40275,
          "gems": [
            42156
          ]
        },
        {
          "id": 40522,
          "enchant": 38374,
          "gems": [
            39999,
            42156
          ]
        },
        {
          "id": 40367,
          "enchant": 55016,
          "gems": [
            40058
          ]
        },
        {
          "id": 40474
        },
        {
          "id": 40074
        },
        {
          "id": 40684
        },
        {
          "id": 37390
        },
        {
          "id": 39763,
          "enchant": 44492
        },
        {
          "id": 39468,
          "enchant": 44492
        },
        {
          "id": 40322
        }
      ]}`)
