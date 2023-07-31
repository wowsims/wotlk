package core

import (
	"fmt"
	"strings"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLActionSequence struct {
	defaultAPLActionImpl
	unit    *Unit
	name    string
	actions []*APLAction
	curIdx  int
}

func (rot *APLRotation) newActionSequence(config *proto.APLActionSequence) APLActionImpl {
	actions := MapSlice(config.Actions, func(action *proto.APLAction) *APLAction {
		return rot.newAPLAction(action)
	})
	actions = FilterSlice(actions, func(action *APLAction) bool { return action != nil })

	return &APLActionSequence{
		unit:    rot.unit,
		name:    config.Name,
		actions: actions,
	}
}
func (action *APLActionSequence) GetInnerActions() []*APLAction {
	return Flatten(MapSlice(action.actions, func(action *APLAction) []*APLAction { return action.GetAllActions() }))
}
func (action *APLActionSequence) Finalize(rot *APLRotation) {
	for _, subaction := range action.actions {
		subaction.impl.Finalize(rot)
	}
}
func (action *APLActionSequence) Reset(*Simulation) {
	action.curIdx = 0
}
func (action *APLActionSequence) IsReady(sim *Simulation) bool {
	return action.curIdx < len(action.actions) && action.actions[action.curIdx].IsReady(sim)
}
func (action *APLActionSequence) Execute(sim *Simulation) {
	action.actions[action.curIdx].Execute(sim)
	action.curIdx++
}
func (action *APLActionSequence) String() string {
	return "Sequence(" + strings.Join(MapSlice(action.actions, func(subaction *APLAction) string { return fmt.Sprintf("(%s)", subaction) }), "+") + ")"
}

type APLActionResetSequence struct {
	defaultAPLActionImpl
	unit     *Unit
	name     string
	sequence *APLActionSequence
}

func (rot *APLRotation) newActionResetSequence(config *proto.APLActionResetSequence) APLActionImpl {
	if config.SequenceName == "" {
		rot.validationWarning("Reset Sequence must provide a sequence name")
		return nil
	}
	return &APLActionResetSequence{
		unit: rot.unit,
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
	rot.validationWarning("No sequence with name: '%s'", action.name)
}
func (action *APLActionResetSequence) IsReady(sim *Simulation) bool {
	return action.sequence != nil
}
func (action *APLActionResetSequence) Execute(sim *Simulation) {
	action.sequence.curIdx = 0
}
func (action *APLActionResetSequence) String() string {
	return fmt.Sprintf("Reset Sequence(name = '%s')", action.name)
}

type APLActionStrictSequence struct {
	defaultAPLActionImpl
	unit    *Unit
	actions []*APLAction
	curIdx  int
}

func (rot *APLRotation) newActionStrictSequence(config *proto.APLActionStrictSequence) APLActionImpl {
	actions := MapSlice(config.Actions, func(action *proto.APLAction) *APLAction {
		return rot.newAPLAction(action)
	})
	actions = FilterSlice(actions, func(action *APLAction) bool { return action != nil })

	return &APLActionStrictSequence{
		unit:    rot.unit,
		actions: actions,
	}
}
func (action *APLActionStrictSequence) GetInnerActions() []*APLAction {
	return Flatten(MapSlice(action.actions, func(action *APLAction) []*APLAction { return action.GetAllActions() }))
}
func (action *APLActionStrictSequence) Finalize(rot *APLRotation) {
	for _, subaction := range action.actions {
		subaction.impl.Finalize(rot)
	}
}
func (action *APLActionStrictSequence) Reset(*Simulation) {
	action.curIdx = 0
}
func (action *APLActionStrictSequence) IsReady(sim *Simulation) bool {
	for _, subaction := range action.actions {
		if !subaction.IsReady(sim) {
			return false
		}
	}
	return true
}
func (action *APLActionStrictSequence) Execute(sim *Simulation) {
	action.unit.Rotation.strictSequence = action
	if !action.actions[action.curIdx].IsReady(sim) {
		action.curIdx = 0
		action.unit.Rotation.strictSequence = nil
		return
	}

	action.actions[action.curIdx].Execute(sim)
	action.curIdx++

	if action.curIdx == len(action.actions) {
		action.curIdx = 0
		action.unit.Rotation.strictSequence = nil
	}
}
func (action *APLActionStrictSequence) String() string {
	return "Strict Sequence(" + strings.Join(MapSlice(action.actions, func(subaction *APLAction) string { return fmt.Sprintf("(%s)", subaction) }), "+") + ")"
}
