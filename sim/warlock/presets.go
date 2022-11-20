package warlock

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var defaultDestroTalents = &proto.WarlockTalents{
	// destro
	ImprovedShadowBolt: 5,
	Bane:               5,
	Ruin:               5,
	Intensity:          2,
	DestructiveReach:   2,
	Devastation:        true,
	Aftermath:          2,
	ImprovedImmolate:   3,
	Emberstorm:         5,
	Conflagrate:        true,
	Backlash:           3,
	Shadowburn:         true,
	ShadowAndFlame:     5,
	Backdraft:          3,
	EmpoweredImp:       3,
	FireAndBrimstone:   5,
	ChaosBolt:          true,
	Shadowfury:         true,
	Pyroclasm:          3,
	DemonicPower:       2,
	Cataclysm:          3,
	SoulLeech:          3,
	ImprovedSoulLeech:  2,
	// demo
	FelSynergy:  2,
	ImprovedImp: 3,
}

var defaultDestroRotation = &proto.Warlock_Rotation{
	Type:         proto.Warlock_Rotation_Destruction,
	PrimarySpell: proto.Warlock_Rotation_Incinerate,
	SecondaryDot: proto.Warlock_Rotation_Immolate,
	SpecSpell:    proto.Warlock_Rotation_ChaosBolt,
	Curse:        proto.Warlock_Rotation_Doom,
	Corruption:   false,
}

var defaultDestroOptions = &proto.Warlock_Options{
	Armor:       proto.Warlock_Options_FelArmor,
	Summon:      proto.Warlock_Options_Imp,
	WeaponImbue: proto.Warlock_Options_GrandFirestone,
}

var DefaultDestroWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Talents:  defaultDestroTalents,
		Options:  defaultDestroOptions,
		Rotation: defaultDestroRotation,
	},
}

// ---------------------------------------
var DefaultAfflictionWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Talents:  defaultAfflictionTalents,
		Options:  defaultAfflictionOptions,
		Rotation: defaultAfflictionRotation,
	},
}

var defaultAfflictionTalents = &proto.WarlockTalents{
	// Affliction
	ImprovedCurseOfAgony:  2,
	Suppression:           3,
	ImprovedCorruption:    5,
	SoulSiphon:            2,
	FelConcentration:      3,
	Nightfall:             2,
	EmpoweredCorruption:   3,
	ShadowEmbrace:         5,
	SiphonLife:            true,
	ImprovedFelhunter:     2,
	ShadowMastery:         5,
	Eradication:           3,
	Contagion:             5,
	DeathsEmbrace:         3,
	UnstableAffliction:    true,
	Pandemic:              true,
	EverlastingAffliction: 5,
	Haunt:                 true,
	// Destro
	ImprovedShadowBolt: 5,
	Bane:               5,
	Ruin:               5,
	Intensity:          1,
}

var defaultAfflictionOptions = &proto.Warlock_Options{
	Armor:       proto.Warlock_Options_FelArmor,
	Summon:      proto.Warlock_Options_Felhunter,
	WeaponImbue: proto.Warlock_Options_GrandSpellstone,
}

var defaultAfflictionRotation = &proto.Warlock_Rotation{
	Type:         proto.Warlock_Rotation_Affliction,
	PrimarySpell: proto.Warlock_Rotation_ShadowBolt,
	SecondaryDot: proto.Warlock_Rotation_UnstableAffliction,
	SpecSpell:    proto.Warlock_Rotation_Haunt,
	Curse:        proto.Warlock_Rotation_Agony,
	Corruption:   true,
	DetonateSeed: true,
}

var defaultAfflictionGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarlockMajorGlyph_GlyphOfQuickDecay),
	Major2: int32(proto.WarlockMajorGlyph_GlyphOfLifeTap),
	Major3: int32(proto.WarlockMajorGlyph_GlyphOfHaunt),
}

// ---------------------------------------
var DefaultDemonologyWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Talents:  defaultDemonologyTalents,
		Options:  defaultDemonologyOptions,
		Rotation: defaultDemonologyRotation,
	},
}

var defaultDemonologyTalents = &proto.WarlockTalents{
	// Demonology
	ImprovedHealthstone: 2,
	DemonicEmbrace:      3,
	FelSynergy:          2,
	DemonicBrutality:    3,
	FelVitality:         3,
	SoulLink:            true,
	DemonicAegis:        3,
	UnholyPower:         5,
	ManaFeed:            true,
	MasterConjuror:      2,
	MasterDemonologist:  5,
	MoltenCore:          3,
	DemonicEmpowerment:  true,
	DemonicKnowledge:    3,
	DemonicTactics:      5,
	Decimation:          2,
	SummonFelguard:      true,
	Nemesis:             3,
	DemonicPact:         5,
	Metamorphosis:       true,
	// Destro
	ImprovedShadowBolt: 5,
	Bane:               5,
	Ruin:               5,
	Intensity:          2,
}

var defaultDemonologyOptions = &proto.Warlock_Options{
	Armor:       proto.Warlock_Options_FelArmor,
	Summon:      proto.Warlock_Options_Felguard,
	WeaponImbue: proto.Warlock_Options_GrandSpellstone,
}

var defaultDemonologyRotation = &proto.Warlock_Rotation{
	Type:         proto.Warlock_Rotation_Demonology,
	PrimarySpell: proto.Warlock_Rotation_ShadowBolt,
	SecondaryDot: proto.Warlock_Rotation_Immolate,
	Curse:        proto.Warlock_Rotation_Doom,
	Corruption:   true,
	DetonateSeed: true,
}

// ---------------------------------------------------------

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	DefaultPotion: proto.Potions_PotionOfWildMagic,
	PrepopPotion:  proto.Potions_PotionOfWildMagic,
	Food:          proto.Food_FoodFishFeast,
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40421,
		"enchant": 3820,
		"gems": [
			41285,
			40051
		]
	},
	{
		"id": 44661,
		"gems": [
			40026
		]
	},
	{
		"id": 40424,
		"enchant": 3810,
		"gems": [
			39998
		]
	},
	{
		"id": 44005,
		"enchant": 3722,
		"gems": [
			40026
		]
	},
	{
		"id": 40423,
		"enchant": 3832,
		"gems": [
			39998,
			40051
		]
	},
	{
		"id": 44008,
		"enchant": 2332,
		"gems": [
			39998,
			0
		]
	},
	{
		"id": 40420,
		"enchant": 3604,
		"gems": [
			39998,
			0
		]
	},
	{
		"id": 40561,
		"gems": [
			39998
		]
	},
	{
		"id": 40560,
		"enchant": 3719
	},
	{
		"id": 40558,
		"enchant": 3606
	},
	{
		"id": 40399
	},
	{
		"id": 40719
	},
	{
		"id": 40432
	},
	{
		"id": 40255
	},
	{
		"id": 40396,
		"enchant": 3834
	},
	{
		"id": 39766
	},
	{
		"id": 39712
	}
]}`)
