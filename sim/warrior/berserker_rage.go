package warrior

import (
	"time"

	"github.com/wowsims/classic/sod/sim/core"
)

func (warrior *Warrior) registerBerserkerRageSpell() {
	if warrior.Level < 32 {
		return
	}

	actionID := core.ActionID{SpellID: 18499}

	warrior.BerserkerRageAura = warrior.RegisterAura(core.Aura{
		Label:    "Berserker Rage",
		ActionID: actionID,
		Duration: time.Second * 10,
	})

	warrior.BerserkerRage = warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 30,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			warrior.BerserkerRageAura.Activate(sim)
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool { return warrior.StanceMatches(BerserkerStance) },
	})
}
