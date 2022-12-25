package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (rogue *Rogue) registerPreparationCD() {
	if !rogue.Talents.Preparation {
		return
	}

	var affectedSpells = []*core.Spell{rogue.ColdBlood, rogue.Shadowstep, rogue.MasterOfSubtlety, rogue.Overkill}

	rogue.Preparation = rogue.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 14185},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: 0,
				GCD:  0,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second*8*60 - time.Second*time.Duration(90*rogue.Talents.FilthyTricks),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			// Reset Cooldown on Evasion, Sprint, Vanish (Overkill/Master of Subtlety), Cold Blood and Shadowstep
			for _, affectedSpell := range affectedSpells {
				if affectedSpell != nil {
					affectedSpell.CD.Reset()
				}
			}
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    rogue.Preparation,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return !rogue.MasterOfSubtlety.CD.IsReady(sim)
		},
	})
}
