package elemental

import (
	"log"
	"os"
	"testing"

	_ "github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterElementalShaman()
}

func GetAplRotation(dir string, file string) core.RotationCombo {
	filePath := dir + "/" + file + ".apl.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to load apl json file: %s, %s", filePath, err)
	}

	return core.RotationCombo{Label: file, Rotation: core.APLRotationFromJsonString(string(data))}
}

func TestElemental(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassShaman,
		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceOrc},

		GearSet:     core.GearSetCombo{Label: "P1", GearSet: P1Gear},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Adaptive", SpecOptions: PlayerOptionsAdaptive},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "EleFireElemental", SpecOptions: PlayerOptionsAdaptiveFireElemental},
		},
		Rotation: GetAplRotation("../../../ui/elemental_shaman/apls", "default"),

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeFist,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeShield,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeMail,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeTotem,
			},
		},

		EPReferenceStat: proto.Stat_StatSpellPower,
		StatsToWeigh: []proto.Stat{
			proto.Stat_StatIntellect,
			proto.Stat_StatSpellPower,
			proto.Stat_StatSpellHit,
			proto.Stat_StatSpellCrit,
		},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:          proto.Race_RaceOrc,
				Class:         proto.Class_ClassShaman,
				Equipment:     P1Gear,
				TalentsString: StandardTalents,
				Glyphs:        StandardGlyphs,
				Consumes:      FullConsumes,
				Spec:          PlayerOptionsAdaptive,
				Buffs:         core.FullIndividualBuffs,
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

var StandardTalents = "0532001523212351322301351-005052031"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.ShamanMajorGlyph_GlyphOfLava),
	Major2: int32(proto.ShamanMajorGlyph_GlyphOfTotemOfWrath),
	Major3: int32(proto.ShamanMajorGlyph_GlyphOfLightningBolt),
}

var NoTotems = &proto.ShamanTotems{}
var BasicTotems = &proto.ShamanTotems{
	Earth: proto.EarthTotem_TremorTotem,
	Air:   proto.AirTotem_WrathOfAirTotem,
	Water: proto.WaterTotem_ManaSpringTotem,
	Fire:  proto.FireTotem_TotemOfWrath,
}

var FireElementalBasicTotems = &proto.ShamanTotems{
	Earth:            proto.EarthTotem_TremorTotem,
	Air:              proto.AirTotem_WrathOfAirTotem,
	Water:            proto.WaterTotem_ManaSpringTotem,
	Fire:             proto.FireTotem_TotemOfWrath,
	UseFireElemental: true,
}

var PlayerOptionsAdaptive = &proto.Player_ElementalShaman{
	ElementalShaman: &proto.ElementalShaman{
		Options: &proto.ElementalShaman_Options{
			Shield:    proto.ShamanShield_WaterShield,
			Bloodlust: true,
			Totems:    BasicTotems,
		},
		Rotation: &proto.ElementalShaman_Rotation{},
	},
}

var PlayerOptionsAdaptiveFireElemental = &proto.Player_ElementalShaman{
	ElementalShaman: &proto.ElementalShaman{
		Options: &proto.ElementalShaman_Options{
			Shield:    proto.ShamanShield_WaterShield,
			Bloodlust: true,
			Totems:    FireElementalBasicTotems,
		},
		Rotation: &proto.ElementalShaman_Rotation{},
	},
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfBlindingLight,
	Food:            proto.Food_FoodBlackenedBasilisk,
	DefaultPotion:   proto.Potions_SuperManaPotion,
	PrepopPotion:    proto.Potions_DestructionPotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var P1Gear = core.EquipmentSpecFromJsonString(`{"items": [
	{"id":40516,"enchant":3820,"gems":[41285,40027]},
	{"id":44661,"gems":[39998]},
	{"id":40286,"enchant":3810},
	{"id":44005,"enchant":3722,"gems":[40027]},
	{"id":40514,"enchant":3832,"gems":[42144,42144]},
	{"id":40324,"enchant":2332,"gems":[42144,0]},
	{"id":40302,"enchant":3246,"gems":[0]},
	{"id":40301,"gems":[40014]},
	{"id":40560,"enchant":3721},
	{"id":40519,"enchant":3826},
	{"id":37694},
	{"id":40399},
	{"id":40432},
	{"id":40255},
	{"id":40395,"enchant":3834},
	{"id":40401,"enchant":1128},
	{"id":40267}
]}`)
