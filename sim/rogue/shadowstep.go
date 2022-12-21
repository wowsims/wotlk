package rogue

import (
	"time"

	"wowsims.com/wotlk/sim/core"
)

func (rogue *Rogue) registerShadowstepCD() {
	if !rogue.Talents.Shadowstep {
		return
	}

	actionID := core.ActionID{SpellID: 36554}
	target := rogue.CurrentTarget

	shadowstepAura := target.GetOrRegisterAura(core.Aura{
		Label:		"Shadowstep",
		ActionID:	actionID,
		Duration:	time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Damage of your next ability is increased by 20% and the threat caused is reduced by 50%.
			rogue.PseudoStats.DamageDealtMultiplier *= 0.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Remove
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// if spell is no ability
			// return
			// Remove aura on hit
			// 20% more damage by ability
		}
	})

	shadowstepSpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.Energy,
		baseCost := rogue.costModifier(10 - 5 * float64(rogue.Talents.FilthyTricks))
		BaseCost:	  baseCost
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: BaseCost,
				GCD:  time.Second,
			},
			IgnoreHaste: truek
			CD: core.Cooldown{
				Timer:		rogue.NewTimer(),
				Duration:	time.Second * (30 - 5 * float64(rogue.Talents.FilthyTricks))
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.shadowstepAura.Activate(sim)
		},
	})

	
	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell: rogue.shadowstepSpell,
		Type:  core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return rogue.CurrentEnergy > 35
		}
	})
}