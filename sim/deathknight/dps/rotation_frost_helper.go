package dps

import (
	"github.com/wowsims/wotlk/sim/core"
)

type FrostRotation struct {
	lastSpell *core.Spell
	nextSpell *core.Spell

	firstBloodStrike bool
}

func (fr *FrostRotation) Reset(sim *core.Simulation) {
	fr.nextSpell = nil
	fr.lastSpell = nil
	fr.firstBloodStrike = true
}

func (dk *DpsDeathknight) FrostRotationCast(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
	fr := &dk.fr
	canCast := dk.CanCast(sim, spell)
	if canCast {
		fr.lastSpell = spell
		return spell.Cast(sim, target)
	}
	return false
}
