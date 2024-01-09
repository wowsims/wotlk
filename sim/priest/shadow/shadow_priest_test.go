package shadow

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get caster sets included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterShadowPriest()
}

func TestShadow(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassPriest,
		Race:       proto.Race_RaceUndead,
		OtherRaces: []proto.Race{proto.Race_RaceNightElf, proto.Race_RaceDraenei},

		GearSet:  core.GetGearSet("../../../ui/shadow_priest/gear_sets", "p1"),
		Talents:  DefaultTalents,
		Glyphs:   DefaultGlyphs,
		Consumes: FullConsumes,

		SpecOptions: core.SpecOptionsCombo{Label: "Ideal", SpecOptions: PlayerOptionsIdeal},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "Basic", SpecOptions: PlayerOptionsBasic},
			{Label: "Clipping", SpecOptions: PlayerOptionsClipping},
		},

		Rotation: core.GetAplRotation("../../../ui/shadow_priest/apls", "default"),

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeCloth,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeWand,
			},
		},
	}))
}

var DefaultTalents = "05032031--325023051223010323151301351"
var DefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.PriestMajorGlyph_GlyphOfShadow),
	Major2: int32(proto.PriestMajorGlyph_GlyphOfMindFlay),
	Major3: int32(proto.PriestMajorGlyph_GlyphOfDispersion),
	// No dps increasing minor glyphs.
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfPureDeath,
	Food:            proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:   proto.Potions_SuperManaPotion,
	PrepopPotion:    proto.Potions_PotionOfWildMagic,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var PlayerOptionsBasic = &proto.Player_ShadowPriest{
	ShadowPriest: &proto.ShadowPriest{
		Options: &proto.ShadowPriest_Options{
			Armor: proto.ShadowPriest_Options_InnerFire,
		},
	},
}
var PlayerOptionsClipping = &proto.Player_ShadowPriest{
	ShadowPriest: &proto.ShadowPriest{
		Options: &proto.ShadowPriest_Options{
			Armor: proto.ShadowPriest_Options_InnerFire,
		},
	},
}
var PlayerOptionsIdeal = &proto.Player_ShadowPriest{
	ShadowPriest: &proto.ShadowPriest{
		Options: &proto.ShadowPriest_Options{
			Armor: proto.ShadowPriest_Options_InnerFire,
		},
	},
}
