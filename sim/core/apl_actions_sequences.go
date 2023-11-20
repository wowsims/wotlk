package core

import (
	"fmt"
	"strings"

	"github.com/wowsims/classic/sim/core/proto"
)

type APLActionSequence struct {
	defaultAPLActionImpl
	unit       *Unit
	name       string
	subactions []*APLAction
	curIdx     int
}

func (rot *APLRotation) newActionSequence(config *proto.APLActionSequence) APLActionImpl {
	subactions := MapSlice(config.Actions, func(action *proto.APLAction) *APLAction {
		return rot.newAPLAction(action)
	})
	subactions = FilterSlice(subactions, func(action *APLAction) bool { return action != nil })
	if len(subactions) == 0 {
		return nil
	}

	return &APLActionSequence{
		unit:       rot.unit,
		name:       config.Name,
		subactions: subactions,
	}
}
func (action *APLActionSequence) GetInnerActions() []*APLAction {
	return Flatten(MapSlice(action.subactions, func(action *APLAction) []*APLAction { return action.GetAllActions() }))
}
func (action *APLActionSequence) Finalize(rot *APLRotation) {
	for _, subaction := range action.subactions {
		subaction.impl.Finalize(rot)
	}
}
func (action *APLActionSequence) Reset(*Simulation) {
	action.curIdx = 0
}
func (action *APLActionSequence) IsReady(sim *Simulation) bool {
	return action.curIdx < len(action.subactions) && action.subactions[action.curIdx].IsReady(sim)
}
func (action *APLActionSequence) Execute(sim *Simulation) {
	action.subactions[action.curIdx].Execute(sim)
	action.curIdx++
}
func (action *APLActionSequence) String() string {
	return "Sequence(" + strings.Join(MapSlice(action.subactions, func(subaction *APLAction) string { return fmt.Sprintf("(%s)", subaction) }), "+") + ")"
}

type APLActionResetSequence struct {
	defaultAPLActionImpl
	name     string
	sequence *APLActionSequence
}

func (rot *APLRotation) newActionResetSequence(config *proto.APLActionResetSequence) APLActionImpl {
	if config.SequenceName == "" {
		rot.ValidationWarning("Reset Sequence must provide a sequence name")
		return nil
	}
	return &APLActionResetSequence{
		name: config.SequenceName,
	}
}
func (action *APLActionResetSequence) Finalize(rot *APLRotation) {
	for _, otherAction := range rot.allAPLActions() {
		if sequence, ok := otherAction.impl.(*APLActionSequence); ok && sequence.name == action.name {
			action.sequence = sequence
			return
		}
	}
	rot.ValidationWarning("No sequence with name: '%s'", action.name)
}
func (action *APLActionResetSequence) IsReady(sim *Simulation) bool {
	return action.sequence != nil && action.sequence.curIdx != 0
}
func (action *APLActionResetSequence) Execute(sim *Simulation) {
	action.sequence.curIdx = 0
}
func (action *APLActionResetSequence) String() string {
	return fmt.Sprintf("Reset Sequence(name = '%s')", action.name)
}

type APLActionStrictSequence struct {
	defaultAPLActionImpl
	unit       *Unit
	subactions []*APLAction
	curIdx     int

	subactionSpells []*Spell
}

func (rot *APLRotation) newActionStrictSequence(config *proto.APLActionStrictSequence) APLActionImpl {
	subactions := MapSlice(config.Actions, func(action *proto.APLAction) *APLAction {
		return rot.newAPLAction(action)
	})
	subactions = FilterSlice(subactions, func(action *APLAction) bool { return action != nil })
	if len(subactions) == 0 {
		return nil
	}

	return &APLActionStrictSequence{
		unit:       rot.unit,
		subactions: subactions,
	}
}
func (action *APLActionStrictSequence) GetInnerActions() []*APLAction {
	return Flatten(MapSlice(action.subactions, func(action *APLAction) []*APLAction { return action.GetAllActions() }))
}
func (action *APLActionStrictSequence) Finalize(rot *APLRotation) {
	for _, subaction := range action.subactions {
		subaction.impl.Finalize(rot)
		action.subactionSpells = append(action.subactionSpells, subaction.GetAllSpells()...)
	}
}
func (action *APLActionStrictSequence) Reset(*Simulation) {
	action.curIdx = 0
}
func (action *APLActionStrictSequence) IsReady(sim *Simulation) bool {
	if !action.unit.GCD.IsReady(sim) {
		return false
	}
	if !action.subactions[0].IsReady(sim) {
		return false
	}
	for _, spell := range action.subactionSpells {
		if !spell.IsReady(sim) {
			return false
		}
	}
	return true
}
func (action *APLActionStrictSequence) Execute(sim *Simulation) {
	action.unit.Rotation.pushControllingAction(action)
}
func (action *APLActionStrictSequence) GetNextAction(sim *Simulation) *APLAction {
	if action.subactions[action.curIdx].IsReady(sim) {
		nextAction := action.subactions[action.curIdx]

		action.curIdx++
		if action.curIdx == len(action.subactions) {
			action.curIdx = 0
			action.unit.Rotation.popControllingAction(action)
		}

		return nextAction
	} else if action.unit.GCD.IsReady(sim) {
		// If the GCD is ready when the next subaction isn't, it means the sequence is bad
		// so reset and exit the sequence.
		action.curIdx = 0
		action.unit.Rotation.popControllingAction(action)
		return action.unit.Rotation.getNextAction(sim)
	} else {
		// Return nil to wait for the GCD to become ready.
		return nil
	}
}
func (action *APLActionStrictSequence) String() string {
	return "Strict Sequence(" + strings.Join(MapSlice(action.subactions, func(subaction *APLAction) string { return fmt.Sprintf("(%s)", subaction) }), "+") + ")"
}
