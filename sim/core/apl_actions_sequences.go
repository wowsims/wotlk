package core

import (
	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLActionSequence struct {
	unit    *Unit
	name    string
	actions []*APLAction
	curIdx  int
}

func (unit *Unit) newActionSequence(config *proto.APLActionSequence) APLActionImpl {
	actions := MapSlice(config.Actions, func(action *proto.APLAction) *APLAction {
		return unit.newAPLAction(action)
	})
	actions = FilterSlice(actions, func(action *APLAction) bool { return action != nil })

	return &APLActionSequence{
		unit:    unit,
		name:    config.Name,
		actions: actions,
	}
}
func (action *APLActionSequence) GetInnerActions() []*APLAction {
	return Flatten(MapSlice(action.actions, func(action *APLAction) []*APLAction { return action.GetAllActions() }))
}
func (action *APLActionSequence) Finalize() {}
func (action *APLActionSequence) Reset(*Simulation) {
	action.curIdx = 0
}
func (action *APLActionSequence) IsAvailable(sim *Simulation) bool {
	return action.curIdx < len(action.actions) && action.actions[action.curIdx].IsAvailable(sim)
}
func (action *APLActionSequence) Execute(sim *Simulation) {
	action.actions[action.curIdx].Execute(sim)
	action.curIdx++
}

type APLActionResetSequence struct {
	unit     *Unit
	name     string
	sequence *APLActionSequence
}

func (unit *Unit) newActionResetSequence(config *proto.APLActionResetSequence) APLActionImpl {
	if config.SequenceName == "" {
		validationError("Reset Sequence must provide a sequence name: %s", config.SequenceName)
		return nil
	}
	return &APLActionResetSequence{
		unit: unit,
		name: config.SequenceName,
	}
}
func (action *APLActionResetSequence) GetInnerActions() []*APLAction { return nil }
func (action *APLActionResetSequence) Finalize() {
	for _, otherAction := range action.unit.allAPLActions() {
		if sequence, ok := otherAction.impl.(*APLActionSequence); ok && sequence.name == action.name {
			action.sequence = sequence
			return
		}
	}
	validationWarning("No sequence with name: %s", action.name)
}
func (action *APLActionResetSequence) Reset(*Simulation) {}
func (action *APLActionResetSequence) IsAvailable(sim *Simulation) bool {
	return true
}
func (action *APLActionResetSequence) Execute(sim *Simulation) {
	action.sequence.curIdx = 0
}

type APLActionStrictSequence struct {
	unit    *Unit
	actions []*APLAction
	curIdx  int
}

func (unit *Unit) newActionStrictSequence(config *proto.APLActionStrictSequence) APLActionImpl {
	actions := MapSlice(config.Actions, func(action *proto.APLAction) *APLAction {
		return unit.newAPLAction(action)
	})
	actions = FilterSlice(actions, func(action *APLAction) bool { return action != nil })

	return &APLActionStrictSequence{
		unit:    unit,
		actions: actions,
	}
}
func (action *APLActionStrictSequence) GetInnerActions() []*APLAction {
	return Flatten(MapSlice(action.actions, func(action *APLAction) []*APLAction { return action.GetAllActions() }))
}
func (action *APLActionStrictSequence) Finalize() {}
func (action *APLActionStrictSequence) Reset(*Simulation) {
	action.curIdx = 0
}
func (action *APLActionStrictSequence) IsAvailable(sim *Simulation) bool {
	return action.actions[action.curIdx].IsAvailable(sim)
}
func (action *APLActionStrictSequence) Execute(sim *Simulation) {
	action.unit.Rotation.strictSequence = action
	if !action.IsAvailable(sim) {
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
