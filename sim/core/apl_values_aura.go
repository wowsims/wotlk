package core

import (
	"fmt"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValueAuraIsActive struct {
	DefaultAPLValueImpl
	aura AuraReference
}

func (rot *APLRotation) newValueAuraIsActive(config *proto.APLValueAuraIsActive) APLValue {
	aura := rot.GetAPLAura(rot.GetSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
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
	return value.aura.Get().IsActive()
}
func (value *APLValueAuraIsActive) String() string {
	return fmt.Sprintf("Aura Active(%s)", value.aura.String())
}

type APLValueAuraIsActiveWithReactionTime struct {
	DefaultAPLValueImpl
	aura         AuraReference
	reactionTime time.Duration
}

func (rot *APLRotation) newValueAuraIsActiveWithReactionTime(config *proto.APLValueAuraIsActiveWithReactionTime) APLValue {
	aura := rot.GetAPLAura(rot.GetSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLValueAuraIsActiveWithReactionTime{
		aura:         aura,
		reactionTime: rot.unit.ReactionTime,
	}
}
func (value *APLValueAuraIsActiveWithReactionTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAuraIsActiveWithReactionTime) GetBool(sim *Simulation) bool {
	aura := value.aura.Get()
	return aura.IsActive() && aura.TimeActive(sim) >= value.reactionTime
}
func (value *APLValueAuraIsActiveWithReactionTime) String() string {
	return fmt.Sprintf("Aura Active With Reaction Time(%s)", value.aura.String())
}

type APLValueAuraRemainingTime struct {
	DefaultAPLValueImpl
	aura AuraReference
}

func (rot *APLRotation) newValueAuraRemainingTime(config *proto.APLValueAuraRemainingTime) APLValue {
	aura := rot.GetAPLAura(rot.GetSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
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
	return value.aura.Get().RemainingDuration(sim)
}
func (value *APLValueAuraRemainingTime) String() string {
	return fmt.Sprintf("Aura Remaining Time(%s)", value.aura.String())
}

type APLValueAuraNumStacks struct {
	DefaultAPLValueImpl
	aura AuraReference
}

func (rot *APLRotation) newValueAuraNumStacks(config *proto.APLValueAuraNumStacks) APLValue {
	aura := rot.GetAPLAura(rot.GetSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	if aura.Get().MaxStacks == 0 {
		rot.ValidationWarning("%s is not a stackable aura", ProtoToActionID(config.AuraId))
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
	return value.aura.Get().GetStacks()
}
func (value *APLValueAuraNumStacks) String() string {
	return fmt.Sprintf("Aura Num Stacks(%s)", value.aura.String())
}

type APLValueAuraInternalCooldown struct {
	DefaultAPLValueImpl
	aura AuraReference
}

func (rot *APLRotation) newValueAuraInternalCooldown(config *proto.APLValueAuraInternalCooldown) APLValue {
	aura := rot.GetAPLICDAura(rot.GetSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
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
	return value.aura.Get().Icd.TimeToReady(sim)
}
func (value *APLValueAuraInternalCooldown) String() string {
	return fmt.Sprintf("Aura Remaining ICD(%s)", value.aura.String())
}

type APLValueAuraICDIsReadyWithReactionTime struct {
	DefaultAPLValueImpl
	aura         AuraReference
	reactionTime time.Duration
}

func (rot *APLRotation) newValueAuraICDIsReadyWithReactionTime(config *proto.APLValueAuraICDIsReadyWithReactionTime) APLValue {
	aura := rot.GetAPLICDAura(rot.GetSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLValueAuraICDIsReadyWithReactionTime{
		aura:         aura,
		reactionTime: rot.unit.ReactionTime,
	}
}
func (value *APLValueAuraICDIsReadyWithReactionTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAuraICDIsReadyWithReactionTime) GetBool(sim *Simulation) bool {
	aura := value.aura.Get()
	return aura.Icd.IsReady(sim) || (aura.IsActive() && aura.TimeActive(sim) < value.reactionTime)
}
func (value *APLValueAuraICDIsReadyWithReactionTime) String() string {
	return fmt.Sprintf("Aura ICD Is Ready with Reaction Time(%s)", value.aura.String())
}

type APLValueAuraShouldRefresh struct {
	DefaultAPLValueImpl
	aura       AuraReference
	maxOverlap APLValue
}

func (rot *APLRotation) newValueAuraShouldRefresh(config *proto.APLValueAuraShouldRefresh) APLValue {
	aura := rot.GetAPLAura(rot.GetTargetUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}

	maxOverlap := rot.coerceTo(rot.newAPLValue(config.MaxOverlap), proto.APLValueType_ValueTypeDuration)
	if maxOverlap == nil {
		maxOverlap = rot.newValueConst(&proto.APLValueConst{Val: "0ms"})
	}

	return &APLValueAuraShouldRefresh{
		aura:       aura,
		maxOverlap: maxOverlap,
	}
}
func (value *APLValueAuraShouldRefresh) GetInnerValues() []APLValue {
	return []APLValue{value.maxOverlap}
}
func (value *APLValueAuraShouldRefresh) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAuraShouldRefresh) GetBool(sim *Simulation) bool {
	return value.aura.Get().ShouldRefreshExclusiveEffects(sim, value.maxOverlap.GetDuration(sim))
}
func (value *APLValueAuraShouldRefresh) String() string {
	return fmt.Sprintf("Should Refresh Aura(%s)", value.aura.String())
}
