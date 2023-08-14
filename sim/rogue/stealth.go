package rogue

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (rogue *Rogue) registerStealthAura() {
	// These spells do NOT break stealth when cast
	nonbreakingSpells := []*core.Spell{rogue.TricksOfTheTrade, rogue.SliceAndDice}

	if rogue.Talents.Premeditation {
		nonbreakingSpells = append(nonbreakingSpells, rogue.Premeditation)
	}

	rogue.StealthAura = rogue.RegisterAura(core.Aura{
		Label:    "Stealth",
		ActionID: core.ActionID{SpellID: 1787},
		Duration: core.NeverExpires,
		/* OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Do not break stealth on certain spell casts
			for _, nobreak := range nonbreakingSpells {
				if nobreak.ActionID == spell.ActionID {
					return
				}
			}

			aura.Deactivate(sim)
		}, */
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Pause auto attacks
			rogue.AutoAttacks.CancelAutoSwing(sim)
			// Stealth triggered auras
			if rogue.Talents.Overkill {
				rogue.OverkillAura.Activate(sim)
			}
			if rogue.Talents.MasterOfSubtlety > 0 {
				rogue.MasterOfSubtletyAura.Activate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.AutoAttacks.EnableAutoSwing(sim)
			if rogue.Talents.Overkill {
				rogue.OverkillAura.Deactivate(sim)
				rogue.OverkillAura.Activate(sim)
			}
			if rogue.Talents.MasterOfSubtlety > 0 {
				rogue.MasterOfSubtletyAura.Deactivate(sim)
				rogue.MasterOfSubtletyAura.Activate(sim)
			}
		},
		// Stealth breaks on damage taken (if not absorbed)
		// This may be desirable later, but not applicable currently
	}) 
}
