package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rot *APLRotation) aplGetSourceUnit(ref *proto.UnitReference) *Unit {
	if ref == nil {
		return rot.unit
	}

	unit := rot.unit.GetUnit(ref)
	if unit == nil {
		rot.validationWarning("No unit found matching reference: %s", ref)
	}
	return unit
}

func (rot *APLRotation) aplGetAura(sourceRef *proto.UnitReference, auraId *proto.ActionID) *Aura {
	sourceUnit := rot.aplGetSourceUnit(sourceRef)
	if sourceUnit == nil {
		return nil
	}

	aura := sourceUnit.GetAuraByID(ProtoToActionID(auraId))
	if aura == nil {
		rot.validationWarning("No aura found for: %s", ProtoToActionID(auraId))
	}
	return aura
}

func (rot *APLRotation) aplGetProcAura(sourceRef *proto.UnitReference, auraId *proto.ActionID) *Aura {
	sourceUnit := rot.aplGetSourceUnit(sourceRef)
	if sourceUnit == nil {
		return nil
	}

	aura := sourceUnit.GetIcdAuraByID(ProtoToActionID(auraId))
	if aura == nil {
		rot.validationWarning("No aura found for: %s", ProtoToActionID(auraId))
	}
	return aura
}

type APLValueAuraIsActive struct {
	defaultAPLValueImpl
	aura *Aura
}

func (rot *APLRotation) newValueAuraIsActive(config *proto.APLValueAuraIsActive) APLValue {
	aura := rot.aplGetAura(config.SourceUnit, config.AuraId)
	if aura == nil {
		return nil
	}
	return &APLValueAuraIsActive{
		aura: aura,
	}
}
func (value *APLValueAuraIsActive) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAuraIsActive) GetBool(sim *Simulation) bool {
	return value.aura.IsActive()
}

type APLValueAuraRemainingTime struct {
	defaultAPLValueImpl
	aura *Aura
}

func (rot *APLRotation) newValueAuraRemainingTime(config *proto.APLValueAuraRemainingTime) APLValue {
	aura := rot.aplGetAura(config.SourceUnit, config.AuraId)
	if aura == nil {
		return nil
	}
	return &APLValueAuraRemainingTime{
		aura: aura,
	}
}
func (value *APLValueAuraRemainingTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueAuraRemainingTime) GetDuration(sim *Simulation) time.Duration {
	return value.aura.RemainingDuration(sim)
}

type APLValueAuraNumStacks struct {
	defaultAPLValueImpl
	aura *Aura
}

func (rot *APLRotation) newValueAuraNumStacks(config *proto.APLValueAuraNumStacks) APLValue {
	aura := rot.aplGetAura(config.SourceUnit, config.AuraId)
	if aura == nil {
		return nil
	}
	if aura.MaxStacks == 0 {
		rot.validationWarning("%s is not a stackable aura", ProtoToActionID(config.AuraId))
		return nil
	}
	return &APLValueAuraNumStacks{
		aura: aura,
	}
}
func (value *APLValueAuraNumStacks) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueAuraNumStacks) GetInt(sim *Simulation) int32 {
	return value.aura.GetStacks()
}

type APLValueAuraInternalCooldown struct {
	defaultAPLValueImpl
	aura *Aura
}

func (rot *APLRotation) newValueAuraInternalCooldown(config *proto.APLValueAuraInternalCooldown) APLValue {
	aura := rot.aplGetProcAura(config.SourceUnit, config.AuraId)
	if aura == nil {
		return nil
	}
	return &APLValueAuraInternalCooldown{
		aura: aura,
	}
}
func (value *APLValueAuraInternalCooldown) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueAuraInternalCooldown) GetDuration(sim *Simulation) time.Duration {
	return value.aura.Icd.TimeToReady(sim)
}
