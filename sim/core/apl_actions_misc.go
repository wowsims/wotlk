package core

import (
	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLActionChangeTarget struct {
	unit      *Unit
	newTarget UnitReference
}

func (rot *APLRotation) newActionChangeTarget(config *proto.APLActionChangeTarget) APLActionImpl {
	newTarget := rot.getSourceUnit(config.NewTarget)
	if newTarget.Get() == nil {
		return nil
	}
	return &APLActionChangeTarget{
		newTarget: newTarget,
	}
}
func (action *APLActionChangeTarget) GetInnerActions() []*APLAction { return nil }
func (action *APLActionChangeTarget) Finalize(*APLRotation)         {}
func (action *APLActionChangeTarget) Reset(*Simulation)             {}
func (action *APLActionChangeTarget) IsReady(sim *Simulation) bool {
	return action.unit.CurrentTarget != action.newTarget.Get()
}
func (action *APLActionChangeTarget) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.unit.Log(sim, "Changing target to %s", action.newTarget.Get().Label)
	}
	action.unit.CurrentTarget = action.newTarget.Get()
}

type APLActionCancelAura struct {
	aura *Aura
}

func (rot *APLRotation) newActionCancelAura(config *proto.APLActionCancelAura) APLActionImpl {
	aura := rot.aplGetAura(&proto.UnitReference{Type: proto.UnitReference_Self}, config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLActionCancelAura{
		aura: aura.Get(),
	}
}
func (action *APLActionCancelAura) GetInnerActions() []*APLAction { return nil }
func (action *APLActionCancelAura) Finalize(*APLRotation)         {}
func (action *APLActionCancelAura) Reset(*Simulation)             {}
func (action *APLActionCancelAura) IsReady(sim *Simulation) bool {
	return action.aura.IsActive()
}
func (action *APLActionCancelAura) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.aura.Unit.Log(sim, "Cancelling aura %s", action.aura.ActionID)
	}
	action.aura.Deactivate(sim)
}

type APLActionTriggerICD struct {
	aura *Aura
}

func (rot *APLRotation) newActionTriggerICD(config *proto.APLActionTriggerICD) APLActionImpl {
	aura := rot.aplGetICDAura(&proto.UnitReference{Type: proto.UnitReference_Self}, config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLActionTriggerICD{
		aura: aura.Get(),
	}
}
func (action *APLActionTriggerICD) GetInnerActions() []*APLAction { return nil }
func (action *APLActionTriggerICD) Finalize(*APLRotation)         {}
func (action *APLActionTriggerICD) Reset(*Simulation)             {}
func (action *APLActionTriggerICD) IsReady(sim *Simulation) bool {
	return action.aura.IsActive()
}
func (action *APLActionTriggerICD) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.aura.Unit.Log(sim, "Triggering ICD %s", action.aura.ActionID)
	}
	action.aura.Icd.Use(sim)
}
