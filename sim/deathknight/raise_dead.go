package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (deathKnight *DeathKnight) registerRaiseDeadCD() {
	// If talented as permanent pet skip this spell
	if deathKnight.Talents.MasterOfGhouls {
		return
	}

	raiseDeadAura := deathKnight.RegisterAura(core.Aura{
		Label:    "Raise Dead",
		ActionID: core.ActionID{SpellID: 46584},
		Duration: time.Minute * 1,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			deathKnight.Ghoul.Enable(sim, deathKnight.Ghoul)
			deathKnight.Ghoul.focusBar.reset(sim)
			deathKnight.Ghoul.AutoAttacks.EnableAutoSwing(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			deathKnight.Ghoul.Disable(sim)
			deathKnight.Ghoul.focusBar.Cancel(sim)
		},
	})

	deathKnight.RaiseDead = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 46584},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = deathKnight.getModifiedGCD()
			},
			CD: core.Cooldown{
				Timer:    deathKnight.NewTimer(),
				Duration: time.Minute*3 - time.Second*45*time.Duration(deathKnight.Talents.NightOfTheDead),
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			raiseDeadAura.Activate(sim)
		},
	})
}

func (deathKnight *DeathKnight) CanRaiseDead(sim *core.Simulation) bool {
	return deathKnight.RaiseDead.IsReady(sim)
}
