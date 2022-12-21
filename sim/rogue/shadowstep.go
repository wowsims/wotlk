package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (rogue *Rogue) registerShadowstepCD() {
	if !rogue.Talents.Shadowstep {
		return
	}

	actionID := core.ActionID{SpellID: 36554}
	target := rogue.CurrentTarget
	baseCost := rogue.costModifier(10 - 5*float64(rogue.Talents.FilthyTricks))

	shadowstepAura := target.GetOrRegisterAura(core.Aura{
		Label:    "Shadowstep",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Damage of your next ability is increased by 20% and the threat caused is reduced by 50%.
			// 20% damage, but only ability casts
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// if spell is no ability
			// return
			// Remove aura on hit
			// 20% more damage by ability
		},
	})

	shadowstepSpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.Energy,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: BaseCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * (30 - 5*float64(rogue.Talents.FilthyTricks)),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.shadowstepAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    rogue.shadowstepSpell,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return rogue.CurrentEnergy > 35
		},
	})
}
