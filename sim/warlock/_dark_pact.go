package warlock

import (
	"math"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) registerDarkPactSpell() {
	if !warlock.Talents.DarkPact {
		return
	}

	actionID := core.ActionID{SpellID: 59092}
	baseRestore := 1200.0
	manaMetrics := warlock.NewManaMetrics(actionID)
	petManaMetrics := warlock.Pet.NewManaMetrics(actionID)

	warlock.DarkPact = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		FlatThreatBonus:  80,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			maxDrain := baseRestore + 0.96*warlock.GetStat(stats.SpellPower)
			actualDrain := math.Min(maxDrain, warlock.Pet.CurrentMana())

			warlock.Pet.SpendMana(sim, actualDrain, petManaMetrics)
			warlock.AddMana(sim, actualDrain, manaMetrics)
		},
	})
}
