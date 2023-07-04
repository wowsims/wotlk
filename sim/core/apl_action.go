package core

import (
	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLAction struct {
	condition APLValue
	impl      APLActionImpl
}

func (action *APLAction) IsReady(sim *Simulation) bool {
	return (action.condition == nil || action.condition.GetBool(sim)) && action.impl.IsReady(sim)
}

func (action *APLAction) Execute(sim *Simulation) {
	action.impl.Execute(sim)
}

// Returns this Action, along with all inner Actions.
func (action *APLAction) GetAllActions() []*APLAction {
	actions := action.impl.GetInnerActions()
	actions = append(actions, action)
	return actions
}

type APLActionImpl interface {
	// Returns all inner Actions.
	GetInnerActions() []*APLAction

	// Performs optional post-processing.
	Finalize()

	// Invoked before each sim iteration.
	Reset(*Simulation)

	// Whether this action is available to be used right now.
	IsReady(*Simulation) bool

	// Performs the action.
	Execute(*Simulation)
}

func (unit *Unit) newAPLAction(config *proto.APLAction) *APLAction {
	if config == nil {
		return nil
	}

	return &APLAction{
		condition: unit.coerceTo(unit.newAPLValue(config.Condition), proto.APLValueType_ValueTypeBool),
		impl:      unit.newAPLActionImpl(config),
	}
}

func (unit *Unit) newAPLActionImpl(config *proto.APLAction) APLActionImpl {
	if config == nil {
		return nil
	}

	switch config.Action.(type) {
	case *proto.APLAction_Sequence:
		return unit.newActionSequence(config.GetSequence())
	case *proto.APLAction_ResetSequence:
		return unit.newActionResetSequence(config.GetResetSequence())
	case *proto.APLAction_StrictSequence:
		return unit.newActionStrictSequence(config.GetStrictSequence())
	case *proto.APLAction_CastSpell:
		return unit.newActionCastSpell(config.GetCastSpell())
	case *proto.APLAction_AutocastOtherCooldowns:
		return unit.newActionAutocastOtherCooldowns(config.GetAutocastOtherCooldowns())
	case *proto.APLAction_Wait:
		return unit.newActionWait(config.GetWait())
	default:
		validationError("Unimplemented action type")
		return nil
	}
}
