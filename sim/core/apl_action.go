package core

import (
	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLAction struct {
	condition APLValue
	impl      APLActionImpl
}

func (action *APLAction) IsAvailable(sim *Simulation) bool {
	return (action.condition == nil || action.condition.GetBool(sim)) && action.impl.IsAvailable(sim)
}

func (action *APLAction) Execute(sim *Simulation) {
	action.impl.Execute(sim)
}

type APLActionImpl interface {
	// Whether this action is available to be used right now.
	IsAvailable(*Simulation) bool

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
	case *proto.APLAction_CastSpell:
		return unit.newActionCastSpell(config.GetCastSpell())
	case *proto.APLAction_Wait:
		return unit.newActionWait(config.GetWait())
	default:
		validationError("Unimplemented action type")
		return nil
	}
}
