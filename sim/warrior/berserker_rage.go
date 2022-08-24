package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (warrior *Warrior) registerBerserkerRageSpell() {
	actionID := core.ActionID{SpellID: 18499}
	rageBonus := 10 * float64(warrior.Talents.ImprovedBerserkerRage)
	rageMetrics := warrior.NewRageMetrics(actionID)
	cooldownDur := time.Second * 30
	if warrior.Talents.IntensifyRage == 1 {
		cooldownDur = time.Duration(float64(cooldownDur) * 0.89)
	} else if warrior.Talents.IntensifyRage == 2 {
		cooldownDur = time.Duration(float64(cooldownDur) * 0.78)
	} else if warrior.Talents.IntensifyRage == 3 {
		cooldownDur = time.Duration(float64(cooldownDur) * 0.67)
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
