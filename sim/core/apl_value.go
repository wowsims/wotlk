package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValue interface {
	// The type of value that will be returned.
	Type() proto.APLValueType

	// Gets the value, assuming it is a particular type. Usually only one of
	// these should be implemented in each class.
	GetBool(*Simulation) bool
	GetInt(*Simulation) int32
	GetFloat(*Simulation) float64
	GetDuration(*Simulation) time.Duration
	GetString(*Simulation) string
}

// Provides empty implementations for the GetX() value interface functions.
type defaultAPLValueImpl struct {
}

func (impl defaultAPLValueImpl) GetBool(sim *Simulation) bool {
	panic("Unimplemented GetBool")
}
func (impl defaultAPLValueImpl) GetInt(sim *Simulation) int32 {
	panic("Unimplemented GetInt")
}
func (impl defaultAPLValueImpl) GetFloat(sim *Simulation) float64 {
	panic("Unimplemented GetFloat")
}
func (impl defaultAPLValueImpl) GetDuration(sim *Simulation) time.Duration {
	panic("Unimplemented GetDuration")
}
func (impl defaultAPLValueImpl) GetString(sim *Simulation) string {
	panic("Unimplemented GetString")
}

func (rot *APLRotation) newAPLValue(config *proto.APLValue) APLValue {
	if config == nil {
		return nil
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

	// Encounter
	case *proto.APLValue_CurrentTime:
		return rot.newValueCurrentTime(config.GetCurrentTime())
	case *proto.APLValue_CurrentTimePercent:
		return rot.newValueCurrentTimePercent(config.GetCurrentTimePercent())
	case *proto.APLValue_RemainingTime:
		return rot.newValueRemainingTime(config.GetRemainingTime())
	case *proto.APLValue_RemainingTimePercent:
		return rot.newValueRemainingTimePercent(config.GetRemainingTimePercent())
	case *proto.APLValue_NumberTargets:
		return rot.newValueNumberTargets(config.GetNumberTargets())

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

	// GCD
	case *proto.APLValue_GcdIsReady:
		return rot.newValueGCDIsReady(config.GetGcdIsReady())
	case *proto.APLValue_GcdTimeToReady:
		return rot.newValueGCDTimeToReady(config.GetGcdTimeToReady())

	// Spells
	case *proto.APLValue_SpellCanCast:
		return rot.newValueSpellCanCast(config.GetSpellCanCast())
	case *proto.APLValue_SpellIsReady:
		return rot.newValueSpellIsReady(config.GetSpellIsReady())
	case *proto.APLValue_SpellTimeToReady:
		return rot.newValueSpellTimeToReady(config.GetSpellTimeToReady())

	// Auras
	case *proto.APLValue_AuraIsActive:
		return rot.newValueAuraIsActive(config.GetAuraIsActive())
	case *proto.APLValue_AuraRemainingTime:
		return rot.newValueAuraRemainingTime(config.GetAuraRemainingTime())
	case *proto.APLValue_AuraNumStacks:
		return rot.newValueAuraNumStacks(config.GetAuraNumStacks())

	// Dots
	case *proto.APLValue_DotIsActive:
		return rot.newValueDotIsActive(config.GetDotIsActive())
	case *proto.APLValue_DotRemainingTime:
		return rot.newValueDotRemainingTime(config.GetDotRemainingTime())

	default:
		return nil
	}
}
