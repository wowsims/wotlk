package common

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

type CustomSpellCondition func(*core.Simulation) bool

type CustomSpell struct {
	Spell     *core.Spell
	Condition CustomSpellCondition
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
		if customSpell.Spell != nil {
			cr.spells = append(cr.spells, customSpell)
		}
	}

	if len(cr.spells) == 0 {
		return nil
	} else {
		return cr
	}
}

func (cr *CustomRotation) ChooseSpell(sim *core.Simulation) *core.Spell {
	for _, customSpell := range cr.spells {
		if customSpell.Condition(sim) {
			return customSpell.Spell
		}
	}
	return nil
}

func (cr *CustomRotation) Cast(sim *core.Simulation) {
	spell := cr.ChooseSpell(sim)

	if spell == nil {
		cr.character.WaitUntil(sim, sim.CurrentTime+time.Millisecond*500)
	}

	success := spell.Cast(sim, cr.character.CurrentTarget)
	if !success {
		if cr.character.HasManaBar() {
			cr.character.WaitForMana(sim, spell.CurCast.Cost)
		}
	}
}
