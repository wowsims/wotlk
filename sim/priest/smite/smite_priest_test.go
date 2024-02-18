package smite

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common" // imported to get caster sets included.
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterSmitePriest()
}

func TestSmite(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassPriest,
		Race:  proto.Race_RaceUndead,

		GearSet:     core.GetGearSet("../../../ui/smite_priest/gear_sets", "p1"),
		Talents:     DefaultTalents,
		Glyphs:      DefaultGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
		Rotation:    core.GetAplRotation("../../../ui/smite_priest/apls", "default"),

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

var DefaultTalents = "05332031013005023310001-005551002020152-00502"
var DefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.PriestMajorGlyph_GlyphOfSmite),
	Major2: int32(proto.PriestMajorGlyph_GlyphOfHolyNova),
	Major3: int32(proto.PriestMajorGlyph_GlyphOfShadowWordDeath),
	// No interesting minor glyphs.
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	Food:          proto.Food_FoodFishFeast,
	DefaultPotion: proto.Potions_RunicManaInjector,
	PrepopPotion:  proto.Potions_PotionOfWildMagic,
}

var PlayerOptionsBasic = &proto.Player_SmitePriest{
	SmitePriest: &proto.SmitePriest{
		Options: &proto.SmitePriest_Options{
			UseInnerFire:   true,
			UseShadowfiend: true,
		},
	},
}
