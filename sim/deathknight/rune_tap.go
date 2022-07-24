package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (deathKnight *DeathKnight) registerRuneTapSpell() {
	if !deathKnight.Talents.RuneTap {
		return
	}

	actionID := core.ActionID{SpellID: 48982}
	cdTimer := deathKnight.NewTimer()
	cd := time.Minute * 1
	healthMetrics := deathKnight.NewHealthMetrics(actionID)

	healthGainMult := 0.0
	if deathKnight.Talents.ImprovedRuneTap == 1 {
		healthGainMult = 0.33
	} else if deathKnight.Talents.ImprovedRuneTap == 2 {
		healthGainMult = 0.66
	} else if deathKnight.Talents.ImprovedRuneTap == 3 {
		healthGainMult = 1.0
	}

	deathKnight.RuneTap = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dkSpellCost := deathKnight.DetermineOptimalCost(sim, 1, 0, 0)
			deathKnight.Spend(sim, spell, dkSpellCost)

			maxHealth := deathKnight.MaxHealth()
			deathKnight.GainHealth(sim, (1.0+healthGainMult)*(maxHealth*0.1), healthMetrics)
		},
	})
}

func (deathKnight *DeathKnight) CanRuneTap(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0, 1, 0, 0) && deathKnight.RuneTap.IsReady(sim)
}

func (deathKnight *DeathKnight) CastRuneTap(sim *core.Simulation, target *core.Unit) bool {
	if deathKnight.CanRuneTap(sim) {
		deathKnight.RuneTap.Cast(sim, target)
		return true
	}
	return false
}
