package dps

import (
	"github.com/wowsims/wotlk/sim/core"
)

type FrostRotation struct {
	lastSpell *core.Spell
	nextSpell *core.Spell
}

func (fr *FrostRotation) Reset(sim *core.Simulation) {
	fr.nextSpell = nil
	fr.lastSpell = nil
}

func (dk *DpsDeathknight) FrostRotationCast(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
	fr := &dk.fr
	if dk.CanCast(sim, spell) {
		spell.Cast(sim, target)
		fr.lastSpell = spell
		return true
	}
	return false
}
