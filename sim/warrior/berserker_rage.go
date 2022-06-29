package warrior

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

func (warrior *Warrior) registerBerserkerRageSpell() {
	actionID := core.ActionID{SpellID: 18499}
	rageBonus := 5 * float64(warrior.Talents.ImprovedBerserkerRage)
	rageMetrics := warrior.NewRageMetrics(actionID)

	warrior.BerserkerRage = warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warrior.AddRage(sim, rageBonus, rageMetrics)
		},
	})
}

func (warrior *Warrior) ShouldBerserkerRage(sim *core.Simulation) bool {
	return warrior.Talents.ImprovedBerserkerRage > 0 && warrior.CurrentRage() < 80 && warrior.BerserkerRage.IsReady(sim)
}
