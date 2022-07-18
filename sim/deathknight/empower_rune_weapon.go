package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (deathKnight *DeathKnight) registerEmpowerRuneWeaponSpell() {
	actionID := core.ActionID{SpellID: 47568}
	cdTimer := deathKnight.NewTimer()
	cd := time.Minute * 5

	deathKnight.EmpowerRuneWeapon = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			deathKnight.RegenAllRunes(sim)

			amountOfRunicPower := 25.0
			deathKnight.AddRunicPower(sim, amountOfRunicPower, deathKnight.UnbreakableArmor.RunicPowerMetrics())
		},
	})
}

func (deathKnight *DeathKnight) CanEmpowerRuneWeapon(sim *core.Simulation) bool {
	return deathKnight.EmpowerRuneWeapon.IsReady(sim) && deathKnight.EmpowerRuneWeapon.CD.IsReady(sim)
}
