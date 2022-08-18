package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerGhoulFrenzySpell() {
	if !dk.Talents.GhoulFrenzy {
		return
	}
	baseCost := float64(core.NewRuneCost(10, 0, 0, 1, 0))
	dk.GhoulFrenzy = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 63560},
		SpellSchool:  core.SpellSchoolShadow,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: baseCost,
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
		},
	}, func(sim *core.Simulation) bool {
		return dk.Talents.GhoulFrenzy && dk.Ghoul.IsEnabled() && dk.CastCostPossible(sim, 0.0, 0, 0, 1) && dk.GhoulFrenzy.IsReady(sim)
	}, nil)

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
			dk.Ghoul.MultiplyMeleeSpeed(sim, 1.25)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.Ghoul.MultiplyMeleeSpeed(sim, 1/1.25)
		},
	})
}
