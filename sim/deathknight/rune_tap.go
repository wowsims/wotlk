package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) registerRuneTapSpell() {
	if !dk.Talents.RuneTap {
		return
	}

	actionID := core.ActionID{SpellID: 48982}
	cdTimer := dk.NewTimer()
	cd := time.Minute * 1
	healthMetrics := dk.NewHealthMetrics(actionID)

	healthGainMult := 0.0
	if dk.Talents.ImprovedRuneTap == 1 {
		healthGainMult = 0.33
	} else if dk.Talents.ImprovedRuneTap == 2 {
		healthGainMult = 0.66
	} else if dk.Talents.ImprovedRuneTap == 3 {
		healthGainMult = 1.0
	}

	dk.RuneTap = dk.RegisterSpell(core.SpellConfig{
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
			maxHealth := dk.MaxHealth()
			dk.GainHealth(sim, (1.0+healthGainMult)*(maxHealth*0.1), healthMetrics)

			dkSpellCost := dk.DetermineOptimalCost(sim, 1, 0, 0)
			dk.Spend(sim, spell, dkSpellCost)
		},
	})
}

func (dk *Deathknight) CanRuneTap(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0, 1, 0, 0) && dk.RuneTap.IsReady(sim)
}

func (dk *Deathknight) CastRuneTap(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanRuneTap(sim) {
		dk.RuneTap.Cast(sim, target)
		return true
	}
	return false
}
