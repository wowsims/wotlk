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
	Finalize(*APLRotation)

	// Invoked before each sim iteration.
	Reset(*Simulation)

	// Whether this action is available to be used right now.
	IsReady(*Simulation) bool

	// Performs the action.
	Execute(*Simulation)
}

func (rot *APLRotation) newAPLAction(config *proto.APLAction) *APLAction {
	if config == nil {
		return nil
	}

	action := &APLAction{
		condition: rot.coerceTo(rot.newAPLValue(config.Condition), proto.APLValueType_ValueTypeBool),
		impl:      rot.newAPLActionImpl(config),
	}

	if action.impl == nil {
		return nil
	} else {
		return action
	}
}

func (rot *APLRotation) newAPLActionImpl(config *proto.APLAction) APLActionImpl {
	if config == nil {
		return nil
	}

	switch config.Action.(type) {
	case *proto.APLAction_Sequence:
		return rot.newActionSequence(config.GetSequence())
	case *proto.APLAction_ResetSequence:
		return rot.newActionResetSequence(config.GetResetSequence())
	case *proto.APLAction_StrictSequence:
		return rot.newActionStrictSequence(config.GetStrictSequence())
	case *proto.APLAction_CastSpell:
		return rot.newActionCastSpell(config.GetCastSpell())
	case *proto.APLAction_Multidot:
		return rot.newActionMultidot(config.GetMultidot())
	case *proto.APLAction_AutocastOtherCooldowns:
		return rot.newActionAutocastOtherCooldowns(config.GetAutocastOtherCooldowns())
	case *proto.APLAction_Wait:
		return rot.newActionWait(config.GetWait())
	default:
		return nil
	}
}
