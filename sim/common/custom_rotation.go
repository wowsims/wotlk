package common

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

type CustomSpellCondition func(*core.Simulation) bool

type CustomSpell struct {
	Spell     *core.Spell
	RuneSpell *deathknight.RuneSpell
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
		if customSpell.Spell != nil || customSpell.RuneSpell != nil {
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
			if customSpell.Spell != nil || customSpell.RuneSpell != nil {
				return &customSpell
			}
		}
	}
	return nil
}

func (cr *CustomRotation) Cast(sim *core.Simulation) {
	spell := cr.ChooseSpell(sim)

	if spell == nil {
		return
	}

	success := false
	if spell.Spell != nil {
		success = spell.Spell.Cast(sim, cr.character.CurrentTarget)
	} else if spell.RuneSpell != nil {
		success = spell.RuneSpell.Cast(sim, cr.character.CurrentTarget)
	}
	if !success {
		if cr.character.HasManaBar() && spell.Spell != nil && spell.RuneSpell == nil {
			cr.character.WaitForMana(sim, spell.Spell.CurCast.Cost)
		}
	}
}
