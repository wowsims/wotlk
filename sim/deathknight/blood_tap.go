package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) registerBloodTapSpell() {
	actionID := core.ActionID{SpellID: 45529}
	cdTimer := dk.NewTimer()
	cd := time.Minute * 1

	dk.BloodTapAura = dk.RegisterAura(core.Aura{
		Label:    "Blood Tap",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.CorrectBloodTapConversion(sim,
				dk.BloodRuneGainMetrics(),
				dk.DeathRuneGainMetrics(),
				dk.BloodTap)

			// Gain at the end, to take into account previous effects for callback
			amountOfRunicPower := 10.0
			dk.AddRunicPower(sim, amountOfRunicPower, dk.BloodTap.RunicPowerMetrics())
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.CancelBloodTap(sim)
		},
	})

	dk.BloodTap = dk.RegisterSpell(core.SpellConfig{
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
			dk.BloodTapAura.Activate(sim)
		},
	})
}
