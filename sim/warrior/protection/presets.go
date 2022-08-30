package protection

import (
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var StandardTalents = &proto.WarriorTalents{
	ImprovedHeroicStrike: 2,
	Deflection:           5,
	TacticalMastery:      3,
	DeepWounds:           3,
	Impale:               2,

	ArmoredToTheTeeth: 3,
	Cruelty:           2,

	ImprovedThunderClap:           3,
	Anticipation:                  5,
	ShieldSpecialization:          5,
	Incite:                        3,
	ImprovedRevenge:               2,
	LastStand:                     true,
	ShieldMastery:                 2,
	Toughness:                     5,
	ConcussionBlow:                true,
	GagOrder:                      2,
	OneHandedWeaponSpecialization: 5,
	Vigilance:                     true,
	ImprovedDefensiveStance:       2,
	Vitality:                      3,
	Warbringer:                    true,
	Devastate:                     true,
	CriticalBlock:                 3,
	SwordAndBoard:                 3,
	DamageShield:                  2,
	Shockwave:                     true,
}

var PlayerOptionsBasic = &proto.Player_ProtectionWarrior{
	ProtectionWarrior: &proto.ProtectionWarrior{
		Talents:  StandardTalents,
		Options:  warriorOptions,
		Rotation: warriorRotation,
	},
}

var DefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarriorMajorGlyph_GlyphOfBlocking),
	Major2: int32(proto.WarriorMajorGlyph_GlyphOfDevastate),
	Major3: int32(proto.WarriorMajorGlyph_GlyphOfVigilance),
	// No interesting minor glyphs.
}

var warriorRotation = &proto.ProtectionWarrior_Rotation{
	HsRageThreshold: 30,
}

var warriorOptions = &proto.ProtectionWarrior_Options{
	Shout:                proto.WarriorShout_WarriorShoutCommanding,
	PrecastShout:         false,
	PrecastShoutT2:       false,
	PrecastShoutSapphire: false,

	StartingRage: 0,
}

var FullRaidBuffs = &proto.RaidBuffs{
	ArcaneBrilliance:     true,
	GiftOfTheWild:        proto.TristateEffect_TristateEffectImproved,
	Thorns:               proto.TristateEffect_TristateEffectImproved,
	Bloodlust:            true,
	BattleShout:          proto.TristateEffect_TristateEffectImproved,
	LeaderOfThePack:      proto.TristateEffect_TristateEffectImproved,
	StrengthOfEarthTotem: proto.TristateEffect_TristateEffectImproved,
	WindfuryTotem:        proto.TristateEffect_TristateEffectImproved,
	UnleashedRage:        true,
}
var FullPartyBuffs = &proto.PartyBuffs{}
var FullIndividualBuffs = &proto.IndividualBuffs{
	BlessingOfKings: true,
	BlessingOfMight: proto.TristateEffect_TristateEffectImproved,
}

var FullConsumes = &proto.Consumes{
	BattleElixir:   proto.BattleElixir_ElixirOfMastery,
	GuardianElixir: proto.GuardianElixir_GiftOfArthas,
}

var FullDebuffs = &proto.Debuffs{
	BloodFrenzy:   true,
	FaerieFire:    proto.TristateEffect_TristateEffectImproved,
	Misery:        true,
	ShadowEmbrace: true,
	ScorpidSting:  true,
}

var P1Gear = items.EquipmentSpecFromJsonString(`{"items": [
	{
	  "id": 40546,
	  "enchant": 44878,
	  "gems": [
		41380,
		40034
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
	  "id": 40544,
	  "enchant": 44489,
	  "gems": [
		40008,
		40008
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
	  "id": 40545,
	  "enchant": 63770,
	  "gems": [
		49110,
		0
	  ]
	},
	{
	  "id": 39759,
	  "enchant": 54793,
	  "gems": [
		40008,
		36767
	  ]
	},
	{
	  "id": 40589,
	  "enchant": 38373
	},
	{
	  "id": 40297,
	  "enchant": 44491
	},
	{
	  "id": 40370
	},
	{
	  "id": 40718
	},
	{
	  "id": 40257
	},
	{
	  "id": 44063,
	  "gems": [
		36767,
		40089
	  ]
	},
	{
	  "id": 40402,
	  "enchant": 22559
	},
	{
	  "id": 40400,
	  "enchant": 44936
	},
	{
	  "id": 41168,
	  "gems": [
		36767
	  ]
	}
  ]}`)
