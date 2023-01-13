package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) registerMarkOfBloodSpell() {
	if !dk.Talents.MarkOfBlood {
		return
	}

	actionID := core.ActionID{SpellID: 49005}
	cdTimer := dk.NewTimer()
	cd := time.Minute * 3

	var markOfBloodAura *core.Aura = nil
	dk.MarkOfBlood = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.GetModifiedGCD()
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if markOfBloodAura == nil {
				markOfBloodAura = core.MarkOfBloodAura(target)
			}

			markOfBloodAura.Activate(sim)
		},
	})

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell:    dk.MarkOfBlood.Spell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeSurvival,
			CanActivate: func(sim *core.Simulation, character *core.Character) bool {
				return dk.MarkOfBlood.CanCast(sim)
			},
		})
	}
}
