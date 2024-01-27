package core

import (
	"fmt"

	"github.com/wowsims/sod/sim/core/proto"
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

type APLActionItemSwap struct {
	defaultAPLActionImpl
	character *Character
	swapSet   proto.APLActionItemSwap_SwapSet
}

func (rot *APLRotation) newActionItemSwap(config *proto.APLActionItemSwap) APLActionImpl {
	if config.SwapSet == proto.APLActionItemSwap_Unknown {
		rot.ValidationWarning("Unknown item swap set")
		return nil
	}

	character := rot.unit.Env.Raid.GetPlayerFromUnit(rot.unit).GetCharacter()
	if !character.ItemSwap.IsEnabled() {
		if config.SwapSet != proto.APLActionItemSwap_Main {
			rot.ValidationWarning("No swap set configured in Settings.")
		}
		return nil
	}

	return &APLActionItemSwap{
		character: character,
		swapSet:   config.SwapSet,
	}
}
func (action *APLActionItemSwap) IsReady(sim *Simulation) bool {
	return (action.swapSet == proto.APLActionItemSwap_Main) == action.character.ItemSwap.IsSwapped()
}
func (action *APLActionItemSwap) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.character.Log(sim, "Item Swap to set %s", action.swapSet)
	}

	if action.swapSet == proto.APLActionItemSwap_Main {
		action.character.ItemSwap.reset(sim)
	} else {
		action.character.ItemSwap.SwapItems(sim, action.character.ItemSwap.slots, true)
	}
}
func (action *APLActionItemSwap) String() string {
	return fmt.Sprintf("Item Swap(%s)", action.swapSet)
}

type APLActionMove struct {
	defaultAPLActionImpl
	unit      *Unit
	moveRange APLValue
}

func (rot *APLRotation) newActionMove(config *proto.APLActionMove) APLActionImpl {
	return &APLActionMove{
		unit:      rot.unit,
		moveRange: rot.newAPLValue(config.RangeFromTarget),
	}
}
func (action *APLActionMove) IsReady(sim *Simulation) bool {
	return !action.unit.Moving && action.moveRange.GetFloat(sim) != action.unit.DistanceFromTarget && action.unit.Hardcast.Expires < sim.CurrentTime
}
func (action *APLActionMove) Execute(sim *Simulation) {
	moveRange := action.moveRange.GetFloat(sim)
	if sim.Log != nil {
		action.unit.Log(sim, "Moving to %s", moveRange)
	}

	action.unit.MoveTo(moveRange, sim)
}
func (action *APLActionMove) String() string {
	return fmt.Sprintf("Move(%s)", action.moveRange)
}
