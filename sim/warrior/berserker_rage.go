package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (warrior *Warrior) registerBerserkerRageSpell() {
	actionID := core.ActionID{SpellID: 18499}
	rageBonus := 10 * float64(warrior.Talents.ImprovedBerserkerRage)
	rageMetrics := warrior.NewRageMetrics(actionID)

	warrior.BerserkerRage = warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: warrior.intensifyRageCooldown(time.Second * 30),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warrior.AddRage(sim, rageBonus, rageMetrics)
		},
	})
}
