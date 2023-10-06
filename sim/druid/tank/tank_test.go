package tank

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterFeralTankDruid()
}

func TestFeralTank(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDruid,
		Race:  proto.Race_RaceTauren,

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsDefault},
		Rotation:    core.GetAplRotation("../../../ui/feral_tank_druid/apls", "default"),

		IsTank:          true,
		InFrontOfTarget: true,

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeLeather,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeIdol,
			},
		},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:      proto.Race_RaceTauren,
				Class:     proto.Class_ClassDruid,
				Equipment: P1Gear,
				Consumes:  FullConsumes,
				Spec:      PlayerOptionsDefault,
				Buffs:     core.FullIndividualBuffs,

				InFrontOfTarget: true,
			},
			core.FullPartyBuffs,
			core.FullRaidBuffs,
			core.FullDebuffs),
		Encounter: &proto.Encounter{
			Duration: 300,
			Targets: []*proto.Target{
				core.NewDefaultTarget(),
			},
		},
		SimOptions: core.AverageDefaultSimTestOptions,
	}

	core.RaidBenchmark(b, rsr)
}

var StandardTalents = "-503232132322010353120300313511-20350001"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.DruidMajorGlyph_GlyphOfMaul),
	Major2: int32(proto.DruidMajorGlyph_GlyphOfSurvivalInstincts),
	Major3: int32(proto.DruidMajorGlyph_GlyphOfFrenziedRegeneration),
}

var PlayerOptionsDefault = &proto.Player_FeralTankDruid{
	FeralTankDruid: &proto.FeralTankDruid{
		Options: &proto.FeralTankDruid_Options{
			InnervateTarget: &proto.UnitReference{}, // no Innervate
			StartingRage:    20,
		},
		Rotation: &proto.FeralTankDruid_Rotation{},
	},
}

var FullConsumes = &proto.Consumes{
	BattleElixir:    proto.BattleElixir_GurusElixir,
	GuardianElixir:  proto.GuardianElixir_GiftOfArthas,
	Food:            proto.Food_FoodBlackenedDragonfin,
	DefaultPotion:   proto.Potions_IndestructiblePotion,
	DefaultConjured: proto.Conjured_ConjuredHealthstone,
	ThermalSapper:   true,
	FillerExplosive: proto.Explosive_ExplosiveSaroniteBomb,
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{"id":40329,"enchant":67839,"gems":[41339,40008]},
	{"id":40387},
	{"id":40494,"enchant":44957,"gems":[40008]},
	{"id":40252,"enchant":3294},
	{"id":40471,"enchant":3832,"gems":[42702,40088]},
	{"id":40186,"enchant":3850,"gems":[40008,0]},
	{"id":40472,"enchant":63770,"gems":[40008,0]},
	{"id":43591,"gems":[40008,40008,40008]},
	{"id":44011,"enchant":38373,"gems":[40008,40008]},
	{"id":40243,"enchant":55016,"gems":[40008]},
	{"id":40370},
	{"id":37784},
	{"id":44253},
	{"id":37220},
	{"id":40280,"enchant":2673},
	{},
	{"id":38365}
]}`)
