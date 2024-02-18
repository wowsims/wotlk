package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValue interface {
	// Returns all inner APLValues.
	GetInnerValues() []APLValue

	// The type of value that will be returned.
	Type() proto.APLValueType

	// Gets the value, assuming it is a particular type. Usually only one of
	// these should be implemented in each class.
	GetBool(*Simulation) bool
	GetInt(*Simulation) int32
	GetFloat(*Simulation) float64
	GetDuration(*Simulation) time.Duration
	GetString(*Simulation) string

	// Performs optional post-processing.
	Finalize(*APLRotation)

	// Pretty-print string for debugging.
	String() string
}

// Provides empty implementations for the GetX() value interface functions.
type DefaultAPLValueImpl struct {
}

func (impl DefaultAPLValueImpl) GetInnerValues() []APLValue { return nil }
func (impl DefaultAPLValueImpl) Finalize(*APLRotation)      {}

func (impl DefaultAPLValueImpl) GetBool(sim *Simulation) bool {
	panic("Unimplemented GetBool")
}
func (impl DefaultAPLValueImpl) GetInt(sim *Simulation) int32 {
	panic("Unimplemented GetInt")
}
func (impl DefaultAPLValueImpl) GetFloat(sim *Simulation) float64 {
	panic("Unimplemented GetFloat")
}
func (impl DefaultAPLValueImpl) GetDuration(sim *Simulation) time.Duration {
	panic("Unimplemented GetDuration")
}
func (impl DefaultAPLValueImpl) GetString(sim *Simulation) string {
	panic("Unimplemented GetString")
}

func (rot *APLRotation) newAPLValue(config *proto.APLValue) APLValue {
	if config == nil {
		return nil
	}

	customValue := rot.unit.Env.GetAgentFromUnit(rot.unit).NewAPLValue(rot, config)
	if customValue != nil {
		return customValue
	}

	switch config.Value.(type) {
	// Operators
	case *proto.APLValue_Const:
		return rot.newValueConst(config.GetConst())
	case *proto.APLValue_And:
		return rot.newValueAnd(config.GetAnd())
	case *proto.APLValue_Or:
		return rot.newValueOr(config.GetOr())
	case *proto.APLValue_Not:
		return rot.newValueNot(config.GetNot())
	case *proto.APLValue_Cmp:
		return rot.newValueCompare(config.GetCmp())
	case *proto.APLValue_Math:
		return rot.newValueMath(config.GetMath())
	case *proto.APLValue_Max:
		return rot.newValueMax(config.GetMax())
	case *proto.APLValue_Min:
		return rot.newValueMin(config.GetMin())

	// Encounter
	case *proto.APLValue_CurrentTime:
		return rot.newValueCurrentTime(config.GetCurrentTime())
	case *proto.APLValue_CurrentTimePercent:
		return rot.newValueCurrentTimePercent(config.GetCurrentTimePercent())
	case *proto.APLValue_RemainingTime:
		return rot.newValueRemainingTime(config.GetRemainingTime())
	case *proto.APLValue_RemainingTimePercent:
		return rot.newValueRemainingTimePercent(config.GetRemainingTimePercent())
	case *proto.APLValue_IsExecutePhase:
		return rot.newValueIsExecutePhase(config.GetIsExecutePhase())
	case *proto.APLValue_NumberTargets:
		return rot.newValueNumberTargets(config.GetNumberTargets())

	// Boss
	case *proto.APLValue_BossSpellIsCasting:
		return rot.newValueBossSpellIsCasting(config.GetBossSpellIsCasting())
	case *proto.APLValue_BossSpellTimeToReady:
		return rot.newValueBossSpellTimeToReady(config.GetBossSpellTimeToReady())

	// Resources
	case *proto.APLValue_CurrentHealth:
		return rot.newValueCurrentHealth(config.GetCurrentHealth())
	case *proto.APLValue_CurrentHealthPercent:
		return rot.newValueCurrentHealthPercent(config.GetCurrentHealthPercent())
	case *proto.APLValue_CurrentMana:
		return rot.newValueCurrentMana(config.GetCurrentMana())
	case *proto.APLValue_CurrentManaPercent:
		return rot.newValueCurrentManaPercent(config.GetCurrentManaPercent())
	case *proto.APLValue_CurrentRage:
		return rot.newValueCurrentRage(config.GetCurrentRage())
	case *proto.APLValue_CurrentEnergy:
		return rot.newValueCurrentEnergy(config.GetCurrentEnergy())
	case *proto.APLValue_CurrentComboPoints:
		return rot.newValueCurrentComboPoints(config.GetCurrentComboPoints())
	case *proto.APLValue_CurrentRunicPower:
		return rot.newValueCurrentRunicPower(config.GetCurrentRunicPower())

	// Resources Runes
	case *proto.APLValue_CurrentRuneCount:
		return rot.newValueCurrentRuneCount(config.GetCurrentRuneCount())
	case *proto.APLValue_CurrentNonDeathRuneCount:
		return rot.newValueCurrentNonDeathRuneCount(config.GetCurrentNonDeathRuneCount())
	case *proto.APLValue_CurrentRuneActive:
		return rot.newValueCurrentRuneActive(config.GetCurrentRuneActive())
	case *proto.APLValue_CurrentRuneDeath:
		return rot.newValueCurrentRuneDeath(config.GetCurrentRuneDeath())
	case *proto.APLValue_RuneCooldown:
		return rot.newValueRuneCooldown(config.GetRuneCooldown())
	case *proto.APLValue_NextRuneCooldown:
		return rot.newValueNextRuneCooldown(config.GetNextRuneCooldown())
	case *proto.APLValue_RuneSlotCooldown:
		return rot.newValueRuneSlotCooldown(config.GetRuneSlotCooldown())
	case *proto.APLValue_RuneGrace:
		return rot.newValueRuneGrace(config.GetRuneGrace())
	case *proto.APLValue_RuneSlotGrace:
		return rot.newValueRuneSlotGrace(config.GetRuneSlotGrace())

	// GCD
	case *proto.APLValue_GcdIsReady:
		return rot.newValueGCDIsReady(config.GetGcdIsReady())
	case *proto.APLValue_GcdTimeToReady:
		return rot.newValueGCDTimeToReady(config.GetGcdTimeToReady())

	// Auto attacks
	case *proto.APLValue_AutoTimeToNext:
		return rot.newValueAutoTimeToNext(config.GetAutoTimeToNext())

	// Spells
	case *proto.APLValue_SpellCanCast:
		return rot.newValueSpellCanCast(config.GetSpellCanCast())
	case *proto.APLValue_SpellIsReady:
		return rot.newValueSpellIsReady(config.GetSpellIsReady())
	case *proto.APLValue_SpellTimeToReady:
		return rot.newValueSpellTimeToReady(config.GetSpellTimeToReady())
	case *proto.APLValue_SpellCastTime:
		return rot.newValueSpellCastTime(config.GetSpellCastTime())
	case *proto.APLValue_SpellTravelTime:
		return rot.newValueSpellTravelTime(config.GetSpellTravelTime())
	case *proto.APLValue_SpellCpm:
		return rot.newValueSpellCPM(config.GetSpellCpm())
	case *proto.APLValue_SpellIsChanneling:
		return rot.newValueSpellIsChanneling(config.GetSpellIsChanneling())
	case *proto.APLValue_SpellChanneledTicks:
		return rot.newValueSpellChanneledTicks(config.GetSpellChanneledTicks())

	// Auras
	case *proto.APLValue_AuraIsActive:
		return rot.newValueAuraIsActive(config.GetAuraIsActive())
	case *proto.APLValue_AuraIsActiveWithReactionTime:
		return rot.newValueAuraIsActiveWithReactionTime(config.GetAuraIsActiveWithReactionTime())
	case *proto.APLValue_AuraRemainingTime:
		return rot.newValueAuraRemainingTime(config.GetAuraRemainingTime())
	case *proto.APLValue_AuraNumStacks:
		return rot.newValueAuraNumStacks(config.GetAuraNumStacks())
	case *proto.APLValue_AuraInternalCooldown:
		return rot.newValueAuraInternalCooldown(config.GetAuraInternalCooldown())
	case *proto.APLValue_AuraIcdIsReadyWithReactionTime:
		return rot.newValueAuraICDIsReadyWithReactionTime(config.GetAuraIcdIsReadyWithReactionTime())
	case *proto.APLValue_AuraShouldRefresh:
		return rot.newValueAuraShouldRefresh(config.GetAuraShouldRefresh())

	// Dots
	case *proto.APLValue_DotIsActive:
		return rot.newValueDotIsActive(config.GetDotIsActive())
	case *proto.APLValue_DotRemainingTime:
		return rot.newValueDotRemainingTime(config.GetDotRemainingTime())

	// Sequences
	case *proto.APLValue_SequenceIsComplete:
		return rot.newValueSequenceIsComplete(config.GetSequenceIsComplete())
	case *proto.APLValue_SequenceIsReady:
		return rot.newValueSequenceIsReady(config.GetSequenceIsReady())
	case *proto.APLValue_SequenceTimeToReady:
		return rot.newValueSequenceTimeToReady(config.GetSequenceTimeToReady())

	// Properties
	case *proto.APLValue_ChannelClipDelay:
		return rot.newValueChannelClipDelay(config.GetChannelClipDelay())

	default:
		return nil
	}
}

// Default implementation of Agent.NewAPLValue so each spec doesn't need this boilerplate.
func (unit *Unit) NewAPLValue(rot *APLRotation, config *proto.APLValue) APLValue {
	return nil
}
