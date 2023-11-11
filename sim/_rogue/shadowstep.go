package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (rogue *Rogue) registerShadowstepCD() {
	if !rogue.Talents.Shadowstep {
		return
	}

	actionID := core.ActionID{SpellID: 36554}
	baseCost := 10 - 5*float64(rogue.Talents.FilthyTricks)
	var affectedSpells []*core.Spell

	rogue.ShadowstepAura = rogue.RegisterAura(core.Aura{
		Label:    "Shadowstep",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range rogue.Spellbook {
				if spell.Flags.Matches(SpellFlagBuilder|SpellFlagFinisher) && spell.DamageMultiplier > 0 {
					affectedSpells = append(affectedSpells, spell)
				}
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Damage of your next ability is increased by 20% and the threat caused is reduced by 50%.
			for _, spell := range affectedSpells {
				spell.DamageMultiplier *= 1.2
				spell.ThreatMultiplier *= 0.5
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplier *= 1 / 1.2
				spell.ThreatMultiplier *= 1 / 0.5
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			for _, affectedSpell := range affectedSpells {
				if spell == affectedSpell {
					aura.Deactivate(sim)
				}
			}
		},
	})

	rogue.Shadowstep = rogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,
		EnergyCost: core.EnergyCostOptions{
			Cost:   baseCost,
			Refund: 0,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * time.Duration(30-5*rogue.Talents.FilthyTricks),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.ShadowstepAura.Activate(sim)
		},
	})
}
