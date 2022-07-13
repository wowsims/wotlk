package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (deathKnight *DeathKnight) registerBloodTapSpell() {
	actionID := core.ActionID{SpellID: 45529}
	cdTimer := deathKnight.NewTimer()
	cd := time.Minute * 1

	deathKnight.BloodTap = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			deathKnight.GenerateDeathRuneFromBloodRune(sim, deathKnight.DeathRuneGainMetrics(), spell)
		},
	})
}

func (deathKnight *DeathKnight) CanBloodTap(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 1, 0, 0) && deathKnight.BloodTap.IsReady(sim)
}
