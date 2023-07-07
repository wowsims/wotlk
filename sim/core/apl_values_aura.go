package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

func (unit *Unit) aplGetAura(auraId *proto.ActionID) *Aura {
	return unit.GetAuraByID(ProtoToActionID(auraId))
}

type APLValueAuraIsActive struct {
	defaultAPLValueImpl
	aura *Aura
}

func (unit *Unit) newValueAuraIsActive(config *proto.APLValueAuraIsActive) APLValue {
	aura := unit.aplGetAura(config.AuraId)
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

func (unit *Unit) newValueAuraRemainingTime(config *proto.APLValueAuraRemainingTime) APLValue {
	aura := unit.aplGetAura(config.AuraId)
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

func (unit *Unit) newValueAuraNumStacks(config *proto.APLValueAuraNumStacks) APLValue {
	aura := unit.aplGetAura(config.AuraId)
	if aura == nil {
		return nil
	}
	if aura.MaxStacks == 0 {
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
