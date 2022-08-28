package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerMarkOfBloodSpell() {
	if !dk.Talents.MarkOfBlood {
		return
	}

	actionID := core.ActionID{SpellID: 49005}
	cdTimer := dk.NewTimer()
	cd := time.Minute * 3

	baseCost := float64(core.NewRuneCost(10, 0, 1, 0, 0))
	var markOfBloodAura *core.Aura = nil
	dk.MarkOfBlood = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:     actionID,
		Flags:        core.SpellFlagNoOnCastComplete,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: baseCost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
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
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 0, 1, 0, 0) && dk.MarkOfBlood.IsReady(sim)
	}, nil)

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell:    dk.MarkOfBlood.Spell,
			Priority: core.CooldownPriorityLow, // Use low prio so other actives get used first.
			Type:     core.CooldownTypeSurvival,
		})
	}
}
