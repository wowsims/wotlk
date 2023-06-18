package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

var OverkillActionID = core.ActionID{SpellID: 58426}

func (rogue *Rogue) registerOverkillCD() {
	if !rogue.Talents.Overkill {
		return
	}
	rogue.OverkillAura = rogue.RegisterAura(core.Aura{
		Label:    "Overkill",
		ActionID: OverkillActionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.ApplyEnergyTickMultiplier(0.3)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.ApplyEnergyTickMultiplier(-0.3)
		},
	})
	rogue.Overkill = rogue.RegisterSpell(core.SpellConfig{
		ActionID: OverkillActionID,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * time.Duration(180-30*rogue.Talents.Elusiveness),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.OverkillAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell: rogue.Overkill,
		Type:  core.CooldownTypeDPS,

		ShouldActivate: func(sim *core.Simulation, c *core.Character) bool {
			return !rogue.OverkillAura.IsActive() && rogue.CurrentEnergy() < 50
		},
	})

}
