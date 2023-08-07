package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

var OverkillActionID = core.ActionID{SpellID: 58426}

// TODO: Infinite length while Stealth aura active
func (rogue *Rogue) registerOverkill() {
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
}
