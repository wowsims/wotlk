package common

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// Custom condition for an action.
type CustomCondition func(*core.Simulation) bool

// Custom action based on a condition. Returns a bool and the CurCast cost.
type CustomAction func(*core.Simulation, *core.Unit) (bool, float64)

type CustomSpell struct {
	Action    CustomAction
	Condition CustomCondition
}

type CustomRotation struct {
	character *core.Character
	spells    []CustomSpell
}

func NewCustomRotation(crProto *proto.CustomRotation, character *core.Character, spellsMap map[int32]CustomSpell) *CustomRotation {
	if crProto == nil || len(crProto.Spells) == 0 {
		return nil
	}

	cr := &CustomRotation{
		character: character,
	}
	for _, customSpellProto := range crProto.Spells {
		customSpell := spellsMap[customSpellProto.Spell]
		if customSpell.Action != nil {
			cr.spells = append(cr.spells, customSpell)
		}
	}

	if len(cr.spells) == 0 {
		return nil
	} else {
		return cr
	}
}

func (cr *CustomRotation) ChooseSpell(sim *core.Simulation) *CustomSpell {
	for _, customSpell := range cr.spells {
		if customSpell.Condition(sim) {
			if customSpell.Action != nil {
				return &customSpell
			}
		}
	}
	return nil
}

func (cr *CustomRotation) Cast(sim *core.Simulation) bool {
	spell := cr.ChooseSpell(sim)

	if spell == nil {
		return false
	}

	success := false
	cost := 0.0
	if spell.Action != nil {
		success, cost = spell.Action(sim, cr.character.CurrentTarget)
	}
	if !success {
		if cr.character.HasManaBar() && spell.Action != nil {
			cr.character.WaitForMana(sim, cost)
		}
	}

	return true
}
