package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) registerEmpowerRuneWeaponSpell() {
	actionID := core.ActionID{SpellID: 47568}
	cdTimer := dk.NewTimer()
	cd := time.Minute * 5

	dk.EmpowerRuneWeapon = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.RegenAllRunes(sim)

			amountOfRunicPower := 25.0
			dk.AddRunicPower(sim, amountOfRunicPower, dk.EmpowerRuneWeapon.RunicPowerMetrics())
		},
	}, func(sim *core.Simulation) bool {
		return dk.EmpowerRuneWeapon.IsReady(sim)
	}, func(sim *core.Simulation) {
		dk.UpdateMajorCooldowns()
	})
}
