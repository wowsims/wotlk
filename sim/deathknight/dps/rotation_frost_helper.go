package dps

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

type FrostRotation struct {
	lastSpell *deathknight.RuneSpell
	nextSpell *deathknight.RuneSpell
}

func (fr *FrostRotation) Reset(sim *core.Simulation) {
	fr.nextSpell = nil
	fr.lastSpell = nil
}

func (dk *DpsDeathknight) FrostRotationCast(sim *core.Simulation, target *core.Unit, spell *deathknight.RuneSpell) bool {
	fr := &dk.fr
	if dk.CanCast(sim, spell) {
		spell.Cast(sim, target)
		fr.lastSpell = spell
		return true
	}
	return false
}
