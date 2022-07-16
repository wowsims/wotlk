package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type GhoulPet struct {
	core.Pet
	focusBar

	dkOwner *DeathKnight

	ClawAbility PetAbility

	uptimePercent float64
}

func (deathKnight *DeathKnight) NewGhoulPet(permanent bool) *GhoulPet {
	ghoulPet := &GhoulPet{
		Pet: core.NewPet(
			"Ghoul",
			&deathKnight.Character,
			ghoulPetBaseStats,
			deathKnight.makeStatInheritance(),
			permanent,
		),
		dkOwner: deathKnight,
	}

	// NightOfTheDead
	ghoulPet.PseudoStats.DamageTakenMultiplier *= (1.0 - float64(deathKnight.Talents.NightOfTheDead)*0.45)

	ghoulPet.EnableFocusBar(func(sim *core.Simulation) {
		if ghoulPet.GCD.IsReady(sim) {
			ghoulPet.OnGCDReady(sim)
		}
	})

	ghoulPet.EnableAutoAttacks(ghoulPet, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  42,
			BaseDamageMax:  68,
			SwingSpeed:     2,
			SwingDuration:  time.Second * 2,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	})

	ghoulPet.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Strength,
		ModifiedStat: stats.AttackPower,
		Modifier: func(strength float64, attackPower float64) float64 {
			return attackPower + strength
		},
	})
	ghoulPet.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.MeleeCrit,
		Modifier: func(agility float64, meleeCrit float64) float64 {
			return meleeCrit + (agility/83.3)*core.CritRatingPerCritChance
		},
	})

	core.ApplyPetConsumeEffects(&ghoulPet.Character, deathKnight.Consumes)

	deathKnight.AddPet(ghoulPet)

	return ghoulPet
}

func (ghoul *GhoulPet) GetPet() *core.Pet {
	return &ghoul.Pet
}

func (ghoulPet *GhoulPet) Initialize() {
	ghoulPet.ClawAbility = ghoulPet.NewPetAbility(Claw)
}

func (ghoulPet *GhoulPet) Reset(sim *core.Simulation) {
	if ghoulPet.IsEnabled() {
		ghoulPet.focusBar.reset(sim)
	} else {
		ghoulPet.AutoAttacks.CancelAutoSwing(sim)
	}

	if sim.Log != nil {
		ghoulPet.Log(sim, "Total Pet stats: %s", ghoulPet.GetStats())
		inheritedStats := ghoulPet.dkOwner.makeStatInheritance()(ghoulPet.dkOwner.GetStats())
		ghoulPet.Log(sim, "Inherited Pet stats: %s", inheritedStats)
	}

	if ghoulPet.dkOwner.Talents.MasterOfGhouls {
		ghoulPet.uptimePercent = core.MinFloat(1, core.MaxFloat(0, ghoulPet.dkOwner.Options.PetUptime))
	} else {
		ghoulPet.uptimePercent = 1.0
	}
}

func (ghoulPet *GhoulPet) OnGCDReady(sim *core.Simulation) {
	// Apply uptime for permanent ghoul
	if ghoulPet.dkOwner.Talents.MasterOfGhouls {
		percentRemaining := sim.GetRemainingDurationPercent()
		if percentRemaining < 1.0-ghoulPet.uptimePercent { // once fight is % completed, disable pet.
			ghoulPet.Disable(sim)
			ghoulPet.focusBar.Cancel(sim)
			return
		}
	}

	target := ghoulPet.CurrentTarget

	if !ghoulPet.ClawAbility.TryCast(sim, target, ghoulPet) {
		ghoulPet.DoNothing()
	}
}

var ghoulPetBaseStats = stats.Stats{
	stats.Agility:     856,
	stats.Strength:    331,
	stats.AttackPower: 836,

	// Add 1.8% because pets aren't affected by that component of crit suppression.
	stats.MeleeCrit: (1.1515 + 1.8) * core.CritRatingPerCritChance,
}

func (deathKnight *DeathKnight) makeStatInheritance() core.PetStatInheritance {
	ravenousDead := 1.0 + 0.2*float64(deathKnight.Talents.RavenousDead)
	glyphBonus := 0.0
	if deathKnight.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfTheGhoul) {
		glyphBonus = 0.4
	}

	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.Stamina:  ownerStats[stats.Stamina] * (glyphBonus + 0.7*ravenousDead),
			stats.Strength: ownerStats[stats.Strength] * (glyphBonus + 0.7*ravenousDead),

			stats.MeleeHit:   ownerStats[stats.MeleeHit],
			stats.SpellHit:   ownerStats[stats.SpellHit],
			stats.Expertise:  ownerStats[stats.Expertise],
			stats.MeleeHaste: ownerStats[stats.MeleeHaste],
			stats.SpellHaste: ownerStats[stats.MeleeHaste],
		}
	}
}
