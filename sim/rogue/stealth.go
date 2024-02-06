package rogue

import (
	"github.com/wowsims/sod/sim/core"
)

func (rogue *Rogue) registerStealthAura() {
	// TODO: Add Stealth spell for use with prepull in APL
	rogue.StealthAura = rogue.RegisterAura(core.Aura{
		Label:    "Stealth",
		ActionID: core.ActionID{SpellID: 1787},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Stealth triggered auras
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
		},
		// Stealth breaks on damage taken (if not absorbed)
		// This may be desirable later, but not applicable currently
	})
}
