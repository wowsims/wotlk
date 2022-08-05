package mage

import (
	"image/color"
	"log"
	"testing"
	"time"

	_ "github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func init() {
	RegisterMage()
}

func TestRand(t *testing.T) {
	mix := core.NewSplitMix(uint64(time.Now().UnixNano()))

	var dist []float64
	for i := 0; i < 100000; i++ {
		dist = append(dist, mix.NextFloat64())
	}
	hist(dist, "uniform", "Uniform Distribution")

}

func hist(dist []float64, name, title string) {
	n := len(dist)
	vals := make(plotter.Values, n)
	for i := 0; i < n; i++ {
		vals[i] = dist[i]
	}

	plt := plot.New()
	plt.Title.Text = title
	hist, err := plotter.NewHist(vals, 25) // 25 bins
	if err != nil {
		log.Println("Cannot plot:", err)
	}
	hist.FillColor = color.RGBA{R: 255, G: 127, B: 80, A: 255} // coral color
	plt.Add(hist)

	err = plt.Save(400, 200, name+".png")
	if err != nil {
		log.Panic(err)
	}
}

func TestArcane(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,

		Race: proto.Race_RaceTroll,

		GearSet: core.GearSetCombo{Label: "P1Arcane", GearSet: P1ArcaneGear},

		SpecOptions: core.SpecOptionsCombo{Label: "ArcaneRotation", SpecOptions: PlayerOptionsArcane},
		OtherSpecOptions: []core.SpecOptionsCombo{
			core.SpecOptionsCombo{Label: "AOE", SpecOptions: PlayerOptionsArcaneAOE},
		},

		RaidBuffs:   FullRaidBuffs,
		PartyBuffs:  FullArcanePartyBuffs,
		PlayerBuffs: FullIndividualBuffs,
		Consumes:    FullArcaneConsumes,
		Debuffs:     FullDebuffs,

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeSword,
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

func TestFire(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,

		Race: proto.Race_RaceTroll,

		GearSet: core.GearSetCombo{Label: "P1Fire", GearSet: P1FireGear},

		SpecOptions: core.SpecOptionsCombo{Label: "FireRotation", SpecOptions: PlayerOptionsFire},
		OtherSpecOptions: []core.SpecOptionsCombo{
			core.SpecOptionsCombo{Label: "AOE", SpecOptions: PlayerOptionsFireAOE},
		},

		RaidBuffs:   FullRaidBuffs,
		PartyBuffs:  FullFirePartyBuffs,
		PlayerBuffs: FullIndividualBuffs,
		Consumes:    FullFireConsumes,
		Debuffs:     FullDebuffs,

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeSword,
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

func TestFrost(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,

		Race: proto.Race_RaceTroll,

		GearSet: core.GearSetCombo{Label: "P1Frost", GearSet: P1FrostGear},

		SpecOptions: core.SpecOptionsCombo{Label: "FrostRotation", SpecOptions: PlayerOptionsFrost},
		OtherSpecOptions: []core.SpecOptionsCombo{
			core.SpecOptionsCombo{Label: "AOE", SpecOptions: PlayerOptionsFrostAOE},
		},

		RaidBuffs:   FullRaidBuffs,
		PartyBuffs:  FullFrostPartyBuffs,
		PlayerBuffs: FullIndividualBuffs,
		Consumes:    FullFrostConsumes,
		Debuffs:     FullDebuffs,

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeSword,
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
