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

	GhoulFrenzyAura *core.Aura

	ClawAbility PetAbility

	uptimePercent float64
}

func (deathKnight *DeathKnight) NewArmyGhoulPet(index int) *GhoulPet {
	ghoulPet := &GhoulPet{
		Pet: core.NewPet(
			"Army of the Dead", //+strconv.Itoa(index),
			&deathKnight.Character,
			ghoulPetBaseStats,
			deathKnight.armyGhoulStatInheritance(),
			false,
		),
		dkOwner: deathKnight,
	}

	ghoulPet.PseudoStats.DamageTakenMultiplier *= 0.1

	deathKnight.SetupGhoul(ghoulPet)

	return ghoulPet
}

func (deathKnight *DeathKnight) NewGhoulPet(permanent bool) *GhoulPet {
	ghoulPet := &GhoulPet{
		Pet: core.NewPet(
			"Ghoul",
			&deathKnight.Character,
			ghoulPetBaseStats,
			deathKnight.ghoulStatInheritance(),
			permanent,
		),
		dkOwner: deathKnight,
	}

	// NightOfTheDead
	ghoulPet.PseudoStats.DamageTakenMultiplier *= (1.0 - float64(deathKnight.Talents.NightOfTheDead)*0.45)

	deathKnight.SetupGhoul(ghoulPet)

	return ghoulPet
}

func (deathKnight *DeathKnight) SetupGhoul(ghoulPet *GhoulPet) {
	ghoulPet.Pet.OnPetEnable = ghoulPet.enable
	ghoulPet.Pet.OnPetDisable = ghoulPet.disable

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

	ghoulPet.AddStatDependency(stats.Strength, stats.AttackPower, 1.0+1)
	ghoulPet.AddStatDependency(stats.Agility, stats.MeleeCrit, 1.0+(core.CritRatingPerCritChance/83.3))

	core.ApplyPetConsumeEffects(&ghoulPet.Character, deathKnight.Consumes)

	deathKnight.AddPet(ghoulPet)
}

func (ghoulPet *GhoulPet) IsPetGhoul() bool {
	return ghoulPet.dkOwner.Talents.MasterOfGhouls && ghoulPet == ghoulPet.dkOwner.Ghoul
}

func (ghoul *GhoulPet) GetPet() *core.Pet {
	return &ghoul.Pet
}

func (ghoulPet *GhoulPet) Initialize() {
	ghoulPet.ClawAbility = ghoulPet.NewPetAbility(Claw)
}

func (ghoulPet *GhoulPet) Reset(sim *core.Simulation) {
	if ghoulPet.IsPetGhoul() {
		ghoulPet.uptimePercent = core.MinFloat(1, core.MaxFloat(0, ghoulPet.dkOwner.Options.PetUptime))
	} else {
		ghoulPet.uptimePercent = 1.0
	}

	ghoulPet.AutoAttacks.CancelAutoSwing(sim)
}

func (ghoulPet *GhoulPet) OnGCDReady(sim *core.Simulation) {
	// Apply uptime for permanent pet ghoul
	if ghoulPet.IsPetGhoul() {
		percentRemaining := sim.GetRemainingDurationPercent()
		if percentRemaining < 1.0-ghoulPet.uptimePercent { // once fight is % completed, disable pet.
			ghoulPet.Pet.Disable(sim)
			return
		}
	}

	target := ghoulPet.CurrentTarget

	if !ghoulPet.ClawAbility.TryCast(sim, target, ghoulPet) {
		ghoulPet.DoNothing()
	}
}

func (ghoulPet *GhoulPet) enable(sim *core.Simulation) {
	ghoulPet.focusBar.Enable(sim)

	// Snapshot extra % speed modifiers from dk owner
	if !ghoulPet.PermanentPet {
		ghoulPet.PseudoStats.MeleeSpeedMultiplier = ghoulPet.dkOwner.PseudoStats.MeleeSpeedMultiplier
	}
}

func (ghoulPet *GhoulPet) disable(sim *core.Simulation) {
	ghoulPet.focusBar.Disable(sim)

	// Clear snapshot speed
	if !ghoulPet.PermanentPet {
		ghoulPet.PseudoStats.MeleeSpeedMultiplier = 1
	}
}

var ghoulPetBaseStats = stats.Stats{
	stats.Agility:     856,
	stats.Strength:    331,
	stats.AttackPower: 836,

	// Add 1.8% because pets aren't affected by that component of crit suppression.
	stats.MeleeCrit: (1.1515 + 1.8) * core.CritRatingPerCritChance,
}

func (deathKnight *DeathKnight) ghoulStatInheritance() core.PetStatInheritance {
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

func (deathKnight *DeathKnight) armyGhoulStatInheritance() core.PetStatInheritance {
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
