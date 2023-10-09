package warlock

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterWarlock()
}

func TestAffliction(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassWarlock,
		Race:  proto.Race_RaceOrc,

		GearSet:     core.GearSetCombo{Label: "P3", GearSet: P3Gear_affliction},
		Talents:     AfflictionTalents,
		Glyphs:      AfflictionGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Affliction Warlock", SpecOptions: DefaultAfflictionWarlock},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "AffItemSwap", SpecOptions: afflictionItemSwap},
		},

		ItemFilter: ItemFilter,
	}))
}

func TestDemonology(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassWarlock,
		Race:  proto.Race_RaceOrc,

		GearSet:     core.GearSetCombo{Label: "P2", GearSet: P3Gear_demo},
		Talents:     DemonologyTalents,
		Glyphs:      DemonologyGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Demonology Warlock", SpecOptions: DefaultDemonologyWarlock},
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../ui/warlock/apls", "demo"),
		},

		ItemFilter: ItemFilter,
	}))
}

func TestDestruction(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassWarlock,
		Race:  proto.Race_RaceOrc,

		GearSet:     core.GearSetCombo{Label: "P2", GearSet: P3Gear_destro},
		Talents:     DestructionTalents,
		Glyphs:      DestructionGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Destruction Warlock", SpecOptions: DefaultDestroWarlock},
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../ui/warlock/apls", "destro"),
		},
		ItemFilter: ItemFilter,
	}))
}

var ItemFilter = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeDagger,
	},
	HandTypes: []proto.HandType{
		proto.HandType_HandTypeOffHand,
	},
	ArmorType: proto.ArmorType_ArmorTypeCloth,
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeWand,
	},
}

var AfflictionTalents = "2350002030023510253500331151--550000051"
var DemonologyTalents = "-203203301035012530135201351-550000052"
var DestructionTalents = "-03310030003-05203205210331051335230351"
var AfflictionGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarlockMajorGlyph_GlyphOfQuickDecay),
	Major2: int32(proto.WarlockMajorGlyph_GlyphOfLifeTap),
	Major3: int32(proto.WarlockMajorGlyph_GlyphOfHaunt),
}
var DemonologyGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarlockMajorGlyph_GlyphOfQuickDecay),
	Major2: int32(proto.WarlockMajorGlyph_GlyphOfLifeTap),
	Major3: int32(proto.WarlockMajorGlyph_GlyphOfFelguard),
}
var DestructionGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarlockMajorGlyph_GlyphOfConflagrate),
	Major2: int32(proto.WarlockMajorGlyph_GlyphOfLifeTap),
	Major3: int32(proto.WarlockMajorGlyph_GlyphOfIncinerate),
}

var defaultDestroRotation = &proto.Warlock_Rotation{
	Type:         proto.Warlock_Rotation_Destruction,
	PrimarySpell: proto.Warlock_Rotation_Incinerate,
	SecondaryDot: proto.Warlock_Rotation_Immolate,
	SpecSpell:    proto.Warlock_Rotation_ChaosBolt,
	Curse:        proto.Warlock_Rotation_Doom,
	Corruption:   false,
	DetonateSeed: true,
}

var defaultDestroOptions = &proto.Warlock_Options{
	Armor:       proto.Warlock_Options_FelArmor,
	Summon:      proto.Warlock_Options_Imp,
	WeaponImbue: proto.Warlock_Options_GrandFirestone,
}

var DefaultDestroWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Options:  defaultDestroOptions,
		Rotation: defaultDestroRotation,
	},
}

// ---------------------------------------
var DefaultAfflictionWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Options:  defaultAfflictionOptions,
		Rotation: defaultAfflictionRotation,
	},
}

var afflictionItemSwap = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Options:  defaultAfflictionOptions,
		Rotation: afflictionItemSwapRotation,
	},
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

var afflictionItemSwapRotation = &proto.Warlock_Rotation{
	Type:             proto.Warlock_Rotation_Affliction,
	PrimarySpell:     proto.Warlock_Rotation_ShadowBolt,
	SecondaryDot:     proto.Warlock_Rotation_UnstableAffliction,
	SpecSpell:        proto.Warlock_Rotation_Haunt,
	Curse:            proto.Warlock_Rotation_Agony,
	Corruption:       true,
	DetonateSeed:     true,
	EnableWeaponSwap: true,
	WeaponSwap: &proto.ItemSwap{
		MhItem: &proto.ItemSpec{
			Id:      45457,
			Enchant: 3790,
			Gems:    []int32{40013, 40013},
		},
	},
}

// ---------------------------------------
var DefaultDemonologyWarlock = &proto.Player_Warlock{
	Warlock: &proto.Warlock{
		Options:  defaultDemonologyOptions,
		Rotation: defaultDemonologyRotation,
	},
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

var P3Gear_affliction = core.EquipmentSpecFromJsonString(`{"items":[
	{"id":47796,"enchant":3820,"gems":[41285,40133]},
	{"id":47468,"gems":[40155]},
	{"id":47793,"enchant":3810,"gems":[40155]},
	{"id":47551,"enchant":3722,"gems":[40113]},
	{"id":47462,"enchant":1144,"gems":[40133,40155,40113]},
	{"id":47485,"enchant":2332,"gems":[40113,0]},
	{"id":47797,"enchant":3604,"gems":[40113,0]},
	{"id":47419,"enchant":3599,"gems":[40133,40113,40113]},
	{"id":47795,"enchant":3872,"gems":[40113,40153]},
	{"id":47454,"enchant":3606,"gems":[40133,40113]},
	{"id":45495,"gems":[40113]},
	{"id":47489,"gems":[40155]},
	{"id":45518},{"id":45466},
	{"id":47422,"enchant":3834,"gems":[40113]},
	{"id":48032,"gems":[40155]},
	{"id":45294,"gems":[40155]}
]}`)

var P3Gear_demo = core.EquipmentSpecFromJsonString(`{"items":[
	{"id":47796,"enchant":3820,"gems":[41285,40133]},
	{"id":45133,"gems":[40153]},
	{"id":47793,"enchant":3810,"gems":[40113]},
	{"id":47554,"enchant":3722,"gems":[40113]},
	{"id":47794,"enchant":1144,"gems":[40113,40133]},
	{"id":47485,"enchant":2332,"gems":[40133,0]},
	{"id":47788,"enchant":3604,"gems":[40113,0]},
	{"id":47419,"enchant":3599,"gems":[40133,40113,40113]},
	{"id":47435,"enchant":3872,"gems":[40113,40133,40133]},
	{"id":47454,"enchant":3606,"gems":[40133,40113]},
	{"id":45495,"gems":[40133]},
	{"id":47489,"gems":[40113]},
	{"id":45518},
	{"id":40255},
	{"id":47422,"enchant":3834,"gems":[40133]},
	{"id":47470},
	{"id":45294,"gems":[40113]}
]}`)

var P3Gear_destro = core.EquipmentSpecFromJsonString(`{"items":[
	{"id":47796,"enchant":3820,"gems":[41285,40133]},
	{"id":47468,"gems":[40153]},
	{"id":47793,"enchant":3810,"gems":[40155]},
	{"id":47551,"enchant":3722,"gems":[40113]},
	{"id":47794,"enchant":1144,"gems":[40113,40133]},
	{"id":47467,"enchant":2332,"gems":[40153,0]},
	{"id":47788,"enchant":3604,"gems":[40113,0]},
	{"id":47419,"enchant":3599,"gems":[40133,40113,40113]},
	{"id":47435,"enchant":3872,"gems":[40113,40133,40133]},
	{"id":47454,"enchant":3606,"gems":[40133,40113]},
	{"id":45495,"gems":[40133]},
	{"id":47489,"gems":[40155]},
	{"id":45518},
	{"id":47477},
	{"id":47422,"enchant":3834,"gems":[40133]},
	{"id":47437},
	{"id":45294,"gems":[40113]}
]}`)
