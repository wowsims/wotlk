package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
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
	baseCost := float64(core.NewRuneCost(0, 1, 0, 0, 0))
	dk.RuneTap = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:     actionID,
		Flags:        core.SpellFlagNoOnCastComplete,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				// TODO: Does not invoke GCD?
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			maxHealth := dk.MaxHealth()
			dk.GainHealth(sim, (1.0+healthGainMult)*(maxHealth*0.1), healthMetrics)
		},
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 0, 1, 0, 0) && dk.RuneTap.IsReady(sim)
	}, nil)

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell:    dk.RuneTap.Spell,
			Priority: core.CooldownPriorityLow, // Use low prio so other actives get used first.
			Type:     core.CooldownTypeSurvival,
		})
	}
}
