package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerDeathPactSpell() {
	actionID := core.ActionID{SpellID: 48743}
	cdTimer := dk.NewTimer()
	cd := time.Minute * 2

	hpMetrics := dk.NewHealthMetrics(actionID)

	rs := &RuneSpell{}
	baseCost := float64(core.NewRuneCost(40.0, 0, 0, 0, 0))
	dk.DeathPact = dk.RegisterSpell(rs, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

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
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healthGain := 0.4 * dk.Ghoul.MaxHealth()
			dk.GainHealth(sim, healthGain, hpMetrics)
			dk.Ghoul.Pet.Disable(sim)

			rs.DoCost(sim)
		},
	}, func(sim *core.Simulation) bool {
		return dk.Ghoul.Pet.IsEnabled() && dk.DeathPact.IsReady(sim)
	}, nil)

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell:    dk.DeathPact.Spell,
			Type:     core.CooldownTypeDPS,
			Priority: core.CooldownPriorityLow,
			ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
				return dk.DeathPact.CanCast(sim) && dk.CurrentHealthPercent() <= 0.75
			},
		})
	}
}
