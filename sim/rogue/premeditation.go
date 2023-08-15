package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (rogue *Rogue) registerPremeditation() {
	if !rogue.Talents.Premeditation {
		return
	}

	comboMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 14183})

	rogue.Premeditation = rogue.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 14183},
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: 0,
				GCD:  0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 20,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.IsStealthed()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.AddComboPoints(sim, 2, comboMetrics)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    rogue.Premeditation,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityLow,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return rogue.ComboPoints() <= 2 && rogue.ShadowDanceAura.IsActive()
		},
	})
}
