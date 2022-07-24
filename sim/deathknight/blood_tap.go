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

			amountOfRunicPower := 10.0
			dk.AddRunicPower(sim, amountOfRunicPower, dk.BloodTap.RunicPowerMetrics())
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
			dk.BloodTapAura.Prioritize()
		},
	})
}

func (dk *Deathknight) CanBloodTap(sim *core.Simulation) bool {
	return dk.BloodTap.IsReady(sim) && dk.BloodTap.CD.IsReady(sim)
}

func (dk *Deathknight) CastBloodTap(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanBloodTap(sim) {
		dk.BloodTap.Cast(sim, target)
		return true
	}
	return false
}
