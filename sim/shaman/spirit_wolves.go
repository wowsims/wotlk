package shaman

import (
	"math"
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type SpiritWolf struct {
	core.Pet

	shamanOwner *Shaman
}

type SpiritWolves struct {
	SpiritWolf1 *SpiritWolf
	SpiritWolf2 *SpiritWolf
}

func (SpiritWolves *SpiritWolves) EnableWithTimeout(sim *core.Simulation) {
	SpiritWolves.SpiritWolf1.EnableWithTimeout(sim, SpiritWolves.SpiritWolf1, time.Second*45)
	SpiritWolves.SpiritWolf2.EnableWithTimeout(sim, SpiritWolves.SpiritWolf2, time.Second*45)
}

func (SpiritWolves *SpiritWolves) CancelGCDTimer(sim *core.Simulation) {
	SpiritWolves.SpiritWolf1.CancelGCDTimer(sim)
	SpiritWolves.SpiritWolf2.CancelGCDTimer(sim)
}

var spiritWolfBaseStats = stats.Stats{
	stats.Stamina:   361,
	stats.Spirit:    109,
	stats.Intellect: 65,
	stats.Armor:     9616,

	stats.Agility:     113,
	stats.Strength:    331,
	stats.AttackPower: -20,

	// Add 1.8% because pets aren't affected by that component of crit suppression.
	stats.MeleeCrit: (1.1515 + 1.8) * core.CritRatingPerCritChance,
}

func (shaman *Shaman) NewSpiritWolf(index int) *SpiritWolf {
	spiritWolf := &SpiritWolf{
		Pet:         core.NewPet("Spirit Wolf "+strconv.Itoa(index), &shaman.Character, spiritWolfBaseStats, shaman.makeStatInheritance(), false, false),
		shamanOwner: shaman,
	}

	spiritWolf.EnableAutoAttacks(spiritWolf, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  246,
			BaseDamageMax:  372,
			SwingSpeed:     1.5,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	})

	spiritWolf.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	spiritWolf.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance/83.3)
	core.ApplyPetConsumeEffects(&spiritWolf.Character, shaman.Consumes)

	shaman.AddPet(spiritWolf)

	return spiritWolf
}

const PetExpertiseScale = 3.25

func (shaman *Shaman) makeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		ownerHitChance := ownerStats[stats.MeleeHit] / core.MeleeHitRatingPerHitChance
		hitRatingFromOwner := math.Floor(ownerHitChance) * core.MeleeHitRatingPerHitChance

		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.3,
			stats.Armor:       ownerStats[stats.Armor] * 0.35,
			stats.AttackPower: ownerStats[stats.AttackPower] * (core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFeralSpirit), 0.61, 0.31)),

			stats.MeleeHit:  hitRatingFromOwner,
			stats.Expertise: math.Floor(math.Floor(ownerHitChance)*PetExpertiseScale) * core.ExpertisePerQuarterPercentReduction,
		}
	}
}

func (spiritWolf *SpiritWolf) Initialize() {
	// Nothing
}

func (spiritWolf *SpiritWolf) ExecuteCustomRotation(_ *core.Simulation) {
}

func (spiritWolf *SpiritWolf) Reset(sim *core.Simulation) {
	spiritWolf.Disable(sim)
	if sim.Log != nil {
		spiritWolf.Log(sim, "Base Stats: %s", spiritWolfBaseStats)
		inheritedStats := spiritWolf.shamanOwner.makeStatInheritance()(spiritWolf.shamanOwner.GetStats())
		spiritWolf.Log(sim, "Inherited Stats: %s", inheritedStats)
		spiritWolf.Log(sim, "Total Stats: %s", spiritWolf.GetStats())
	}
}

func (spiritWolf *SpiritWolf) GetPet() *core.Pet {
	return &spiritWolf.Pet
}
