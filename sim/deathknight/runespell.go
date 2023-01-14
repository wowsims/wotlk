package deathknight

import "github.com/wowsims/wotlk/sim/core"

// RuneSpell is a wrapper around a normal core.Spell that allows for management of spending
// runes and runic power. Specifically this also allows for "refunding" of missed refundable spells.

type RuneSpell struct {
	*core.Spell
	Refundable bool
	dk         *Deathknight
}

func (rs *RuneSpell) DoCost(sim *core.Simulation) {
	// Spend now if there is no way to refund the spell
	if !rs.Refundable {
		cost := core.RuneCost(rs.Spell.CurCast.Cost)
		rs.Spell.Unit.SpendRuneCost(sim, rs.Spell, cost)
	}
}

func (rs *RuneSpell) CanCast(sim *core.Simulation) bool {
	if rs == nil {
		return false
	} else {
		return rs.Spell.CanCast(sim, nil)
	}
}

func (rs *RuneSpell) Cast(sim *core.Simulation, target *core.Unit) bool {
	if !rs.CanCast(sim) {
		return false
	}
	result := rs.Spell.Cast(sim, target)
	if !result {
		return result
	}

	rs.DoCost(sim)

	return result
}

// RegisterSpell will connect the underlying spell to the given RuneSpell.
//
//	If no RuneSpell is provided, it will be constructed here.
func (dk *Deathknight) RegisterSpell(rs *RuneSpell, spellConfig core.SpellConfig) *RuneSpell {
	if rs == nil {
		rs = &RuneSpell{}
	}
	rs.dk = dk
	rs.Spell = dk.Character.RegisterSpell(spellConfig)
	return rs
}
