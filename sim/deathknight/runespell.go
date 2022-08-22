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

	CanCast RuneSpellCanCast
	onCast  RuneSpellOnCast
}

func (rs *RuneSpell) OnOutcome(sim *core.Simulation, outcome core.HitOutcome) {
	cost := core.RuneCost(rs.Spell.CurCast.Cost) // cost was already optimized in RuneSpell.Cast
	if cost == 0 {
		return // it was free this time. we dont care
	}
	if outcome.Matches(core.OutcomeLanded) {
		slot1, slot2, slot3 := rs.Spell.Unit.SpendRuneCost(sim, rs.Spell, cost)
		if rs.DeathConvertChance == 0 {
			return
		}
		if rs.DeathConvertChance < 1 && sim.RandomFloat("Blood of The North / Reaping / DRM") > rs.DeathConvertChance {
			return
		}
		if (rs.ConvertType&RuneTypeBlood != 0 && (slot1 == 0 || slot1 == 1)) ||
			(rs.ConvertType&RuneTypeFrost != 0 && (slot1 == 2 || slot1 == 3)) ||
			rs.ConvertType&RuneTypeUnholy != 0 && (slot1 == 4 || slot1 == 5) {
			rs.dk.ConvertToDeath(sim, slot1, true, core.NeverExpires)
		}
		if (rs.ConvertType&RuneTypeBlood != 0 && (slot2 == 0 || slot2 == 1)) ||
			(rs.ConvertType&RuneTypeFrost != 0 && (slot2 == 2 || slot2 == 3)) ||
			rs.ConvertType&RuneTypeUnholy != 0 && (slot2 == 4 || slot2 == 5) {
			rs.dk.ConvertToDeath(sim, slot2, true, core.NeverExpires)
		}
		if (rs.ConvertType&RuneTypeBlood != 0 && (slot3 == 0 || slot3 == 1)) ||
			(rs.ConvertType&RuneTypeFrost != 0 && (slot3 == 2 || slot3 == 3)) ||
			rs.ConvertType&RuneTypeUnholy != 0 && (slot3 == 4 || slot3 == 5) {
			rs.dk.ConvertToDeath(sim, slot3, true, core.NeverExpires)
		}
	}
	// misses just dont get spent as a way to avoid having to cancel regeneration PAs
}

func (rs *RuneSpell) DoCost(sim *core.Simulation) {
	cost := core.RuneCost(rs.Spell.CurCast.Cost)
	// Spend now if there is no way to refund the spell
	if !cost.HasRune() || !rs.Refundable {
		rs.Spell.Unit.SpendRuneCost(sim, rs.Spell, cost)
	}
}

func (rs *RuneSpell) castInternal(sim *core.Simulation, target *core.Unit) bool {
	result := rs.Spell.Cast(sim, target)
	if !result {
		return result
	}

	rs.dk.LastCast = rs

	rs.DoCost(sim)

	if rs.onCast != nil {
		rs.onCast(sim)
	}

	return result
}

func (rs *RuneSpell) Cast(sim *core.Simulation, target *core.Unit) bool {
	if rs.CanCast == nil {
		return rs.castInternal(sim, target)
	} else {
		if rs.CanCast(sim) {
			return rs.castInternal(sim, target)
		}
		return false
	}
}

// RegisterSpell will connect the underlying spell to the given RuneSpell.
//  If no RuneSpell is provided, it will be constructed here.
func (dk *Deathknight) RegisterSpell(rs *RuneSpell, spellConfig core.SpellConfig, canCast func(sim *core.Simulation) bool, onCast func(sim *core.Simulation)) *RuneSpell {
	if rs == nil {
		rs = &RuneSpell{}
	}
	rs.dk = dk
	rs.CanCast = canCast
	rs.onCast = onCast
	rs.Spell = dk.Character.RegisterSpell(spellConfig)
	return rs
}

// withRuneRefund is a wrapper around spell effects that on a miss provides a rune refund.
func (dk *Deathknight) withRuneRefund(rs *RuneSpell, baseEffect core.SpellEffect, isAOE bool) core.ApplySpellEffects {
	var baseEffects []core.SpellEffect
	if isAOE && dk.Env.GetNumTargets() > 1 {
		numTargets := dk.Env.GetNumTargets()
		baseEffects = make([]core.SpellEffect, numTargets)
		for i := range baseEffects {
			baseEffects[i] = baseEffect
			baseEffects[i].Target = dk.Env.GetTargetUnit(int32(i))
		}
	} else {
		baseEffects = []core.SpellEffect{baseEffect}
	}

	rs.Refundable = true
	return core.ApplyEffectFuncWithOutcome(baseEffects, rs.OnOutcome)
}
