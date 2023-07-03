package common

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

type CustomRotationType byte

const (
	Basic CustomRotationType = iota
	CPM
)

// Custom condition for an action.
type CustomCondition func(*core.Simulation) bool

// Custom action based on a condition. Returns a bool and the CurCast cost.
type CustomAction func(*core.Simulation, *core.Unit) (bool, float64)

type CustomSpell struct {
	Spell      *core.Spell // Might be nil if this is not a spell action.
	Action     CustomAction
	Condition  CustomCondition
	DesiredCPM float64

	casts int // Number of casts thus far in the current iteration.
}

type CustomRotation struct {
	rotationType CustomRotationType
	character    *core.Character
	spells       []CustomSpell
}

func NewCustomRotation(crProto *proto.CustomRotation, character *core.Character, spellsMap map[int32]CustomSpell) *CustomRotation {
	if crProto == nil || len(crProto.Spells) == 0 {
		return nil
	}

	cr := &CustomRotation{
		rotationType: Basic,
		character:    character,
	}

	for _, customSpellProto := range crProto.Spells {
		customSpell := spellsMap[customSpellProto.Spell]
		customSpell.DesiredCPM = customSpellProto.CastsPerMinute
		if customSpell.DesiredCPM > 0 {
			cr.rotationType = CPM
		}
		if customSpell.Action == nil && customSpell.Spell != nil {
			spell := customSpell.Spell
			customSpell.Action = func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				success := spell.Cast(sim, target)
				return success, spell.CurCast.Cost
			}
		}
		if customSpell.Condition == nil {
			spell := customSpell.Spell
			customSpell.Condition = func(sim *core.Simulation) bool {
				return spell.CanCast(sim, character.CurrentTarget)
			}
		}
		if customSpell.Action != nil {
			cr.spells = append(cr.spells, customSpell)
		}
	}

	if len(cr.spells) == 0 {
		return nil
	} else {
		cr.character.RegisterResetEffect(func(sim *core.Simulation) {
			cr.reset(sim)
		})
		return cr
	}
}

func (cr *CustomRotation) reset(sim *core.Simulation) {
	for i, _ := range cr.spells {
		cr.spells[i].casts = 0
	}
}

func (cr *CustomRotation) ChooseSpell(sim *core.Simulation) *CustomSpell {
	if cr.rotationType == Basic {
		return cr.chooseSpellBasic(sim)
	} else {
		return cr.chooseSpellCPM(sim)
	}
}

func (cr *CustomRotation) chooseSpellBasic(sim *core.Simulation) *CustomSpell {
	for _, customSpell := range cr.spells {
		if customSpell.Condition(sim) {
			return &customSpell
		}
	}
	return nil
}

func (cr *CustomRotation) chooseSpellCPM(sim *core.Simulation) *CustomSpell {
	for i, _ := range cr.spells {
		customSpell := &cr.spells[i]
		if customSpell.CPM(sim) <= customSpell.DesiredCPM && customSpell.Condition(sim) {
			return customSpell
		}
	}
	return nil
}

func (cs *CustomSpell) CPM(sim *core.Simulation) float64 {
	if sim.CurrentTime == 0 {
		return 0
	}
	return float64(cs.casts) / (float64(sim.CurrentTime) / float64(time.Minute))
}

func (cr *CustomRotation) Cast(sim *core.Simulation) bool {
	if cr == nil {
		panic("Custom Rotation is empty")
	}

	spell := cr.ChooseSpell(sim)

	if spell == nil {
		cr.character.WaitUntil(sim, sim.CurrentTime+time.Millisecond*100)
		return false
	}

	success, cost := spell.Action(sim, cr.character.CurrentTarget)
	if success {
		spell.casts++
	} else {
		if cr.character.HasManaBar() {
			cr.character.WaitForMana(sim, cost)
		}
	}

	return true
}
