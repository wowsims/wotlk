package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) registerGhoulFrenzySpell() {
	if !dk.Talents.GhoulFrenzy {
		return
	}

	dk.GhoulFrenzy = dk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 63560},
		SpellSchool: core.SpellSchoolShadow,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dk.GhoulFrenzyAura.Activate(sim)
			dk.Ghoul.GhoulFrenzyAura.Activate(sim)

			dkSpellCost := dk.DetermineCost(sim, core.DKCastEnum_U)
			dk.Spend(sim, spell, dkSpellCost)

			amountOfRunicPower := 10.0
			dk.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
		},
	})

	dk.GhoulFrenzyAura = dk.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 63560},
		Label:    "Ghoul Frenzy",
		Duration: time.Second * 30.0,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			if dk.Inputs.PrecastGhoulFrenzy {
				dk.GhoulFrenzyAura.Activate(sim)
				dk.GhoulFrenzyAura.UpdateExpires(sim.CurrentTime + time.Second*20)
			}
		},
	})

	dk.Ghoul.GhoulFrenzyAura = dk.Ghoul.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 63560},
		Label:    "Ghoul Frenzy",
		Duration: time.Second * 30.0,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			if dk.Inputs.PrecastGhoulFrenzy {
				dk.Ghoul.GhoulFrenzyAura.Activate(sim)
				dk.Ghoul.GhoulFrenzyAura.UpdateExpires(sim.CurrentTime + time.Second*20)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyMeleeSpeed(sim, 1.25)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyMeleeSpeed(sim, 1/1.25)
		},
	})
}

func (dk *Deathknight) CanGhoulFrenzy(sim *core.Simulation) bool {
	return dk.Talents.GhoulFrenzy && dk.Ghoul.IsEnabled() && dk.CastCostPossible(sim, 0.0, 0, 0, 1) && dk.GhoulFrenzy.IsReady(sim)
}

func (dk *Deathknight) CastGhoulFrenzy(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanGhoulFrenzy(sim) {
		dk.GhoulFrenzy.Cast(sim, target)
		return true
	}
	return false
}
