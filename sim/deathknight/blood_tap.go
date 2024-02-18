package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) registerBloodTapSpell() {
	actionID := core.ActionID{SpellID: 45529}
	cdTimer := dk.NewTimer()
	cd := time.Minute * 1

	rpMetrics := dk.NewRunicPowerMetrics(actionID)
	dk.BloodTapAura = dk.RegisterAura(core.Aura{
		Label:    "Blood Tap",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.BloodTapConversion(sim)

			// Gain at the end, to take into account previous effects for callback
			amountOfRunicPower := 10.0
			dk.AddRunicPower(sim, amountOfRunicPower, rpMetrics)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.CancelBloodTap(sim)
		},
	})

	dk.BloodTap = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

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

	if !dk.Inputs.IsDps && dk.HasSetBonus(ItemSetScourgelordsPlate, 4) {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell: dk.BloodTap,
			Type:  core.CooldownTypeSurvival,
		})
	}
}
