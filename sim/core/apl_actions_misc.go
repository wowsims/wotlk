package core

import (
	"fmt"

	"github.com/wowsims/classic/sod/sim/core/proto"
)

type APLActionChangeTarget struct {
	defaultAPLActionImpl
	unit      *Unit
	newTarget UnitReference
}

func (rot *APLRotation) newActionChangeTarget(config *proto.APLActionChangeTarget) APLActionImpl {
	newTarget := rot.GetSourceUnit(config.NewTarget)
	if newTarget.Get() == nil {
		return nil
	}
	return &APLActionChangeTarget{
		newTarget: newTarget,
	}
}
func (action *APLActionChangeTarget) IsReady(sim *Simulation) bool {
	return action.unit.CurrentTarget != action.newTarget.Get()
}
func (action *APLActionChangeTarget) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.unit.Log(sim, "Changing target to %s", action.newTarget.Get().Label)
	}
	action.unit.CurrentTarget = action.newTarget.Get()
}
func (action *APLActionChangeTarget) String() string {
	return fmt.Sprintf("Change Target(%s)", action.newTarget.Get().Label)
}

type APLActionCancelAura struct {
	defaultAPLActionImpl
	aura *Aura
}

type APLActionActivateAura struct {
	defaultAPLActionImpl
	aura *Aura
}

func (rot *APLRotation) newActionCancelAura(config *proto.APLActionCancelAura) APLActionImpl {
	aura := rot.GetAPLAura(rot.GetSourceUnit(&proto.UnitReference{Type: proto.UnitReference_Self}), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLActionCancelAura{
		aura: aura.Get(),
	}
}

func (rot *APLRotation) newActionActivateAura(config *proto.APLActionActivateAura) APLActionImpl {
	aura := rot.GetAPLAura(rot.GetSourceUnit(&proto.UnitReference{Type: proto.UnitReference_Self}), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLActionActivateAura{
		aura: aura.Get(),
	}
}

func (action *APLActionCancelAura) IsReady(sim *Simulation) bool {
	return action.aura.IsActive()
}
func (action *APLActionCancelAura) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.aura.Unit.Log(sim, "Cancelling aura %s", action.aura.ActionID)
	}
	action.aura.Deactivate(sim)
}
func (action *APLActionCancelAura) String() string {
	return fmt.Sprintf("Cancel Aura(%s)", action.aura.ActionID)
}

func (action *APLActionActivateAura) IsReady(sim *Simulation) bool {
	return true
}

func (action *APLActionActivateAura) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.aura.Unit.Log(sim, "Activating aura %s", action.aura.ActionID)
	}
	action.aura.Activate(sim)
}

func (action *APLActionActivateAura) String() string {
	return fmt.Sprintf("Activate Aura(%s)", action.aura.ActionID)
}

type APLActionTriggerICD struct {
	defaultAPLActionImpl
	aura *Aura
}

func (rot *APLRotation) newActionTriggerICD(config *proto.APLActionTriggerICD) APLActionImpl {
	aura := rot.GetAPLICDAura(rot.GetSourceUnit(&proto.UnitReference{Type: proto.UnitReference_Self}), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLActionTriggerICD{
		aura: aura.Get(),
	}
}
func (action *APLActionTriggerICD) IsReady(sim *Simulation) bool {
	return action.aura.IsActive()
}
func (action *APLActionTriggerICD) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.aura.Unit.Log(sim, "Triggering ICD %s", action.aura.ActionID)
	}
	action.aura.Icd.Use(sim)
}
func (action *APLActionTriggerICD) String() string {
	return fmt.Sprintf("Trigger ICD(%s)", action.aura.ActionID)
}
