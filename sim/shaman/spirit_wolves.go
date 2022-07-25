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
	SpiritWolves.SpiritWolf2.EnableWithTimeout(sim, SpiritWolves.SpiritWolf1, time.Second*45)
}

func (SpiritWolves *SpiritWolves) CancelGCDTimer(sim *core.Simulation) {
	SpiritWolves.SpiritWolf1.CancelGCDTimer(sim)
	SpiritWolves.SpiritWolf2.CancelGCDTimer(sim)
}

// Source: https://web.archive.org/web/20201120214816/https://github.com/dalaranwow/dalaran-wow/issues/4670
var spiritWolfBaseStats = stats.Stats{
	stats.Agility:     113,
	stats.Strength:    331,
	stats.AttackPower: 836, // uncertain number. 1 STR = 2AP, but surely AGI has a contribution?. This is copy-value from dk/ghoul

	// Add 1.8% because pets aren't affected by that component of crit suppression.
	stats.MeleeCrit: (1.1515 + 1.8) * core.CritRatingPerCritChance,
}

func (shaman *Shaman) NewSpiritWolf(index int) *SpiritWolf {
	spiritWolf := &SpiritWolf{
		Pet: core.NewPet(
			"Spirit Wolf "+strconv.Itoa(index),
			&shaman.Character,
			spiritWolfBaseStats,
			shaman.makeStatInheritance(),
			false,
		),
		shamanOwner: shaman,
	}

	spiritWolf.EnableAutoAttacks(spiritWolf, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  218,
			BaseDamageMax:  311,
			SwingSpeed:     1.5,
			SwingDuration:  time.Millisecond * 1500,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	})

	spiritWolf.AddStatDependency(stats.Strength, stats.AttackPower, 1.0+2)
	spiritWolf.AddStatDependency(stats.Agility, stats.MeleeCrit, 1.0+(core.CritRatingPerCritChance/83.3))
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
			stats.AttackPower: ownerStats[stats.AttackPower] * (0.3 + core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFeralSpirit), 0.3, 0)),
			// stats.SpellPower:  ownerStats[stats.AttackPower] * 0.128, FIXME: Is this relevant even tough technically true

			stats.MeleeHit: hitRatingFromOwner,
			// stats.SpellHit:  hitRatingFromOwner * 2, FIXME: Is it relevant?
			stats.Expertise: math.Floor((math.Floor(ownerHitChance) * PetExpertiseScale)) * core.ExpertisePerQuarterPercentReduction,
		}
	}
}

func (spiritWolf *SpiritWolf) Initialize() {}

func (spiritWolf *SpiritWolf) OnGCDReady(sim *core.Simulation) {
	spiritWolf.DoNothing()
}

func (spiritWolf *SpiritWolf) Reset(sim *core.Simulation) {
	spiritWolf.AutoAttacks.CancelAutoSwing(sim)
	if sim.Log != nil {
		spiritWolf.Log(sim, "Total Spirit Wolf Stats: %s", spiritWolf.GetStats())
		inheritedStats := spiritWolf.shamanOwner.makeStatInheritance()(spiritWolf.shamanOwner.GetStats())
		spiritWolf.Log(sim, "Inherited Pet stats: %s", inheritedStats)
	}

	// spiritWolf.uptimePercent = core.MinFloat(1, core.MaxFloat(0, spiritWolf.shamanOwner.Options.PetUptime))
}

func (spiritWolf *SpiritWolf) GetPet() *core.Pet {
	return &spiritWolf.Pet
}
