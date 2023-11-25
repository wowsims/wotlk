package warrior

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (warrior *Warrior) registerBloodrageCD() {
	actionID := core.ActionID{SpellID: 2687}
	rageMetrics := warrior.NewRageMetrics(actionID)

	instantRage := 10.0 + []float64{2, 5}[warrior.Talents.ImprovedBloodrage]
	ragePerSec := 1.0

	warrior.BloodrageAura = warrior.RegisterAura(core.Aura{
		Label:    "Bloodrage",
		ActionID: actionID,
		Duration: time.Second * 10,
	})

	warrior.Bloodrage = warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warrior.BloodrageAura.Activate(sim)
			warrior.AddRage(sim, instantRage, rageMetrics)

			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 10,
				Period:   time.Second * 1,
				OnAction: func(sim *core.Simulation) {
					warrior.AddRage(sim, ragePerSec, rageMetrics)
				},
			})
		},
	})
}
