package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (deathKnight *Deathknight) registerGhoulFrenzySpell() {
	if !deathKnight.Talents.GhoulFrenzy {
		return
	}

	deathKnight.GhoulFrenzy = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 63560},
		SpellSchool: core.SpellSchoolShadow,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = deathKnight.getModifiedGCD()
			},
			CD: core.Cooldown{
				Timer:    deathKnight.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			deathKnight.GhoulFrenzyAura.Activate(sim)
			deathKnight.Ghoul.GhoulFrenzyAura.Activate(sim)

			dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 0, 1)
			deathKnight.Spend(sim, spell, dkSpellCost)

			amountOfRunicPower := 10.0
			deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
		},
	})

	deathKnight.GhoulFrenzyAura = deathKnight.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 63560},
		Label:    "Ghoul Frenzy",
		Duration: time.Second * 30.0,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			if deathKnight.Options.PrecastGhoulFrenzy {
				deathKnight.GhoulFrenzyAura.Activate(sim)
				deathKnight.GhoulFrenzyAura.UpdateExpires(sim.CurrentTime + time.Second*20)
			}
		},
	})

	deathKnight.Ghoul.GhoulFrenzyAura = deathKnight.Ghoul.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 63560},
		Label:    "Ghoul Frenzy",
		Duration: time.Second * 30.0,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			if deathKnight.Options.PrecastGhoulFrenzy {
				deathKnight.Ghoul.GhoulFrenzyAura.Activate(sim)
				deathKnight.Ghoul.GhoulFrenzyAura.UpdateExpires(sim.CurrentTime + time.Second*20)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.MeleeSpeedMultiplier *= 1.25
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.MeleeSpeedMultiplier /= 1.25
		},
	})
}

func (deathKnight *Deathknight) CanGhoulFrenzy(sim *core.Simulation) bool {
	return deathKnight.Talents.GhoulFrenzy && deathKnight.Ghoul.IsEnabled() && deathKnight.CastCostPossible(sim, 0.0, 0, 0, 1) && deathKnight.GhoulFrenzy.IsReady(sim)
}

func (deathKnight *Deathknight) CastGhoulFrenzy(sim *core.Simulation, target *core.Unit) bool {
	if deathKnight.CanGhoulFrenzy(sim) {
		deathKnight.GhoulFrenzy.Cast(sim, target)
		return true
	}
	return false
}
