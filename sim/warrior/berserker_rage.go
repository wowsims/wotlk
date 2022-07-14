package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (warrior *Warrior) registerBerserkerRageSpell() {
	actionID := core.ActionID{SpellID: 18499}
	rageBonus := 5 * float64(warrior.Talents.ImprovedBerserkerRage)
	rageMetrics := warrior.NewRageMetrics(actionID)
	cooldownDur := time.Second * 30
	if warrior.Talents.IntensifyRage == 1 {
		cooldownDur *= (100 - 11) / 100
	} else if warrior.Talents.IntensifyRage == 2 {
		cooldownDur *= (100 - 22) / 100
	} else if warrior.Talents.IntensifyRage == 3 {
		cooldownDur *= (100 - 33) / 100
	}
	warrior.BerserkerRage = warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownDur,
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
