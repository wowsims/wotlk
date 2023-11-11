package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

var OverkillActionID = core.ActionID{SpellID: 58426}

func (rogue *Rogue) registerOverkill() {
	if !rogue.Talents.Overkill {
		return
	}

	effectDuration := time.Second * 20
	if rogue.StealthAura.IsActive() {
		effectDuration = core.NeverExpires
	}

	rogue.OverkillAura = rogue.RegisterAura(core.Aura{
		Label:    "Overkill",
		ActionID: OverkillActionID,
		Duration: effectDuration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.ApplyEnergyTickMultiplier(0.3)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.ApplyEnergyTickMultiplier(-0.3)
		},
	})
}
