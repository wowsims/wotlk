package deathknight

import "github.com/wowsims/wotlk/sim/core"

// RuneSpell is a wrapper around a normal core.Spell that allows for management of spending
// runes and runic power. Specifically this also allows for "refunding" of missed refundable spells.

type RuneSpellCanCast func(sim *core.Simulation) bool
type RuneSpellOnCast func(sim *core.Simulation)

type RuneType int8

const (
	RuneTypeNone RuneType = 1 << iota
	RuneTypeBlood
	RuneTypeFrost
	RuneTypeUnholy
)

type RuneSpell struct {
	*core.Spell
	Refundable         bool
	DeathConvertChance float64
	ConvertType        RuneType
	dk                 *Deathknight

	canCast RuneSpellCanCast
}

func (rs *RuneSpell) OnResult(sim *core.Simulation, result *core.SpellResult) {
	cost := core.RuneCost(rs.Spell.CurCast.Cost) // cost was already optimized in RuneSpell.Cast
	if cost == 0 {
		return // it was free this time. we don't care
	}
	if result.Landed() {
		slot1, slot2, slot3 := rs.Spell.Unit.SpendRuneCost(sim, rs.Spell, cost)
		if !sim.Proc(rs.DeathConvertChance, "Blood of The North / Reaping / DRM") {
			return
		}

		for _, slot := range [3]int8{slot1, slot2, slot3} {
			if (rs.ConvertType&RuneTypeBlood != 0 && (slot == 0 || slot == 1)) ||
				(rs.ConvertType&RuneTypeFrost != 0 && (slot == 2 || slot == 3)) ||
				rs.ConvertType&RuneTypeUnholy != 0 && (slot == 4 || slot == 5) {

				// If the slot to be converted is already blood-tapped, then we convert the other blood rune
				if (slot == 0 || slot == 1) && rs.dk.IsBloodTappedRune(slot) && rs.ConvertType&RuneTypeBlood != 0 {
					otherRune := (slot + 1) % 2
					rs.dk.ConvertToDeath(sim, otherRune, core.NeverExpires)
				} else {
					rs.dk.ConvertToDeath(sim, slot, core.NeverExpires)
				}
			}
		}
	}
	// misses just don't get spent as a way to avoid having to cancel regeneration PAs
}

func (rs *RuneSpell) DoCost(sim *core.Simulation) {
	cost := core.RuneCost(rs.Spell.CurCast.Cost)
	// Spend now if there is no way to refund the spell
	if !rs.Refundable {
		rs.Spell.Unit.SpendRuneCost(sim, rs.Spell, cost)
	}
}

func (rs *RuneSpell) castInternal(sim *core.Simulation, target *core.Unit) bool {
	result := rs.Spell.Cast(sim, target)
	if !result {
		return result
	}

	rs.DoCost(sim)

	return result
}

func (rs *RuneSpell) CanCast(sim *core.Simulation) bool {
	if rs == nil {
		return false
	} else if rs.canCast == nil {
		return true
	} else {
		return rs.canCast(sim)
	}
}

func (rs *RuneSpell) Cast(sim *core.Simulation, target *core.Unit) bool {
	if rs.CanCast(sim) {
		return rs.castInternal(sim, target)
	}
	return false
}

// RegisterSpell will connect the underlying spell to the given RuneSpell.
//
//	If no RuneSpell is provided, it will be constructed here.
func (dk *Deathknight) RegisterSpell(rs *RuneSpell, spellConfig core.SpellConfig, canCast func(sim *core.Simulation) bool) *RuneSpell {
	if rs == nil {
		rs = &RuneSpell{}
	}
	rs.dk = dk
	rs.canCast = canCast
	rs.Spell = dk.Character.RegisterSpell(spellConfig)
	return rs
}
