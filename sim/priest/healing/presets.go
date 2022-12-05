package healing

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var DiscTalents = &proto.PriestTalents{
	TwinDisciplines:            5,
	ImprovedInnerFire:          3,
	ImprovedPowerWordFortitude: 2,
	Meditation:                 3,
	InnerFocus:                 true,
	ImprovedPowerWordShield:    3,
	MentalAgility:              3,
	MentalStrength:             5,
	SoulWarding:                true,
	FocusedPower:               2,
	Enlightenment:              3,
	PowerInfusion:              true,
	ImprovedFlashHeal:          3,
	RenewedHope:                1,
	Rapture:                    3,
	Aspiration:                 2,
	DivineAegis:                3,
	PainSuppression:            true,
	Grace:                      2,
	BorrowedTime:               5,
	Penance:                    true,

	HealingFocus:       2,
	ImprovedRenew:      3,
	HolySpecialization: 5,
	SpellWarding:       1,
	DesperatePrayer:    true,
	Inspiration:        3,
	ImprovedHealing:    3,
}

var HolyTalents = &proto.PriestTalents{
	TwinDisciplines:            5,
	ImprovedInnerFire:          3,
	ImprovedPowerWordFortitude: 2,
	Meditation:                 3,
	InnerFocus:                 true,
	ImprovedPowerWordShield:    1,
	MentalAgility:              3,

	HealingFocus:       2,
	ImprovedRenew:      3,
	HolySpecialization: 4,
	DivineFury:         5,
	DesperatePrayer:    true,
	Inspiration:        3,
	HolyReach:          2,
	HealingPrayers:     2,
	SpiritOfRedemption: true,
	SpiritualGuidance:  5,
	SurgeOfLight:       2,
	SpiritualHealing:   5,
	HolyConcentration:  3,
	EmpoweredHealing:   4,
	Serendipity:        3,
	EmpoweredRenew:     1,
	CircleOfHealing:    true,
	DivineProvidence:   5,
	GuardianSpirit:     true,
}

var DiscGlyphs = &proto.Glyphs{
	Major1: int32(proto.PriestMajorGlyph_GlyphOfPowerWordShield),
	Major2: int32(proto.PriestMajorGlyph_GlyphOfFlashHeal),
	Major3: int32(proto.PriestMajorGlyph_GlyphOfPenance),
	// No interesting minor glyphs.
}

var HolyGlyphs = &proto.Glyphs{
	Major1: int32(proto.PriestMajorGlyph_GlyphOfPrayerOfHealing),
	Major2: int32(proto.PriestMajorGlyph_GlyphOfRenew),
	Major3: int32(proto.PriestMajorGlyph_GlyphOfCircleOfHealing),
	// No interesting minor glyphs.
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	Food:          proto.Food_FoodFishFeast,
	DefaultPotion: proto.Potions_RunicManaInjector,
	PrepopPotion:  proto.Potions_PotionOfWildMagic,
}

var PlayerOptionsDisc = &proto.Player_HealingPriest{
	HealingPriest: &proto.HealingPriest{
		Talents: DiscTalents,
		Options: &proto.HealingPriest_Options{
			UseInnerFire:   true,
			UseShadowfiend: true,
			RaptureChance:  0.8,
		},
		Rotation: &proto.HealingPriest_Rotation{},
	},
}

var PlayerOptionsHoly = &proto.Player_HealingPriest{
	HealingPriest: &proto.HealingPriest{
		Talents: HolyTalents,
		Options: &proto.HealingPriest_Options{
			UseInnerFire:   true,
			UseShadowfiend: true,
		},
		Rotation: &proto.HealingPriest_Rotation{},
	},
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{
		"id": 40456,
		"enchant": 3819,
		"gems": [
			41401,
			39998
		]
	},
	{
		"id": 44657,
		"gems": [
			40047
		]
	},
	{
		"id": 40450,
		"enchant": 3809,
		"gems": [
			42144
		]
	},
	{
		"id": 40724,
		"enchant": 3859
	},
	{
		"id": 40194,
		"enchant": 3832,
		"gems": [
			42144
		]
	},
	{
		"id": 40741,
		"enchant": 2332,
		"gems": [
			0
		]
	},
	{
		"id": 40445,
		"enchant": 3246,
		"gems": [
			42144,
			0
		]
	},
	{
		"id": 40271,
		"enchant": 3601,
		"gems": [
			40027,
			39998
		]
	},
	{
		"id": 40398,
		"enchant": 3719,
		"gems": [
			39998,
			39998
		]
	},
	{
		"id": 40236,
		"enchant": 3606
	},
	{
		"id": 40108
	},
	{
		"id": 40433
	},
	{
		"id": 37835
	},
	{
		"id": 40258
	},
	{
		"id": 40395,
		"enchant": 3834
	},
	{
		"id": 40350
	},
	{
		"id": 40245
	}
]}`)
