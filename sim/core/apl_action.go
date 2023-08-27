package core

import (
	"fmt"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLAction struct {
	condition APLValue
	impl      APLActionImpl
}

func (action *APLAction) Finalize(rot *APLRotation) {
	action.impl.Finalize(rot)
	for _, value := range action.GetAllAPLValues() {
		value.Finalize(rot)
	}
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

// Returns all APLValues used by this action and all of its inner Actions.
func (action *APLAction) GetAllAPLValues() []APLValue {
	var values []APLValue
	for _, a := range action.GetAllActions() {
		values = append(values, a.impl.GetAPLValues()...)
		if a.condition != nil {
			values = append(values, a.condition)
			values = append(values, a.condition.GetInnerValues()...)
		}
	}
	return FilterSlice(values, func(val APLValue) bool { return val != nil })
}

func (action *APLAction) String() string {
	if action.condition == nil {
		return fmt.Sprintf("ACTION = %s", action.impl)
	} else {
		return fmt.Sprintf("IF = %s\nACTION = %s", action.condition, action.impl)
	}
}

type APLActionImpl interface {
	// Returns all inner APL Actions.
	GetInnerActions() []*APLAction

	// Returns all APLValues used by this Action (but not by inner Actions).
	GetAPLValues() []APLValue

	// Performs optional post-processing.
	Finalize(*APLRotation)

	// Invoked before each sim iteration.
	Reset(*Simulation)

	// Whether this action is available to be used right now.
	IsReady(*Simulation) bool

	// Performs the action.
	Execute(*Simulation)

	// Pretty-print string for debugging.
	String() string
}

// Provides empty implementations for the action impl interface functions.
type defaultAPLActionImpl struct {
}

func (impl defaultAPLActionImpl) GetInnerActions() []*APLAction { return nil }
func (impl defaultAPLActionImpl) GetAPLValues() []APLValue      { return nil }
func (impl defaultAPLActionImpl) Finalize(*APLRotation)         {}
func (impl defaultAPLActionImpl) Reset(*Simulation)             {}

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

	customAction := rot.unit.Env.Raid.GetPlayerFromUnit(rot.unit).NewAPLAction(rot, config)
	if customAction != nil {
		return customAction
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
	case *proto.APLAction_Multishield:
		return rot.newActionMultishield(config.GetMultishield())
	case *proto.APLAction_AutocastOtherCooldowns:
		return rot.newActionAutocastOtherCooldowns(config.GetAutocastOtherCooldowns())
	case *proto.APLAction_ChangeTarget:
		return rot.newActionChangeTarget(config.GetChangeTarget())
	case *proto.APLAction_ActivateAura:
		return rot.newActionActivateAura(config.GetActivateAura())
	case *proto.APLAction_CancelAura:
		return rot.newActionCancelAura(config.GetCancelAura())
	case *proto.APLAction_TriggerIcd:
		return rot.newActionTriggerICD(config.GetTriggerIcd())
	case *proto.APLAction_Wait:
		return rot.newActionWait(config.GetWait())
	default:
		return nil
	}
}

// Default implementation of Agent.NewAPLAction so each spec doesn't need this boilerplate.
func (unit *Unit) NewAPLAction(rot *APLRotation, config *proto.APLAction) APLActionImpl {
	return nil
}
