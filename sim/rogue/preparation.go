package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (rogue *Rogue) registerPreparationCD() {
	if !rogue.Talents.Preparation {
		return
	}

	preparationSpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 14185},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second*8*60 - time.Second*90*rogue.Talents.FilthyTricks,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			// Reset Cooldown on Evasion, Sprint, Vanish (Overkill/Master of Subtlety), Cold Blood and Shadowstep
			// FIXME: rogue.ColdBlood.CD.Reset()
			rogue.Shadowstep.CD.Reset()
		},
	})
}
