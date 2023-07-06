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

func (unit *Unit) newAPLValue(config *proto.APLValue) APLValue {
	if config == nil {
		return nil
	}

	switch config.Value.(type) {
	// Operators
	case *proto.APLValue_Const:
		return unit.newValueConst(config.GetConst())
	case *proto.APLValue_And:
		return unit.newValueAnd(config.GetAnd())
	case *proto.APLValue_Or:
		return unit.newValueOr(config.GetOr())
	case *proto.APLValue_Not:
		return unit.newValueNot(config.GetNot())
	case *proto.APLValue_Cmp:
		return unit.newValueCompare(config.GetCmp())

	// Encounter
	case *proto.APLValue_CurrentTime:
		return unit.newValueCurrentTime(config.GetCurrentTime())
	case *proto.APLValue_CurrentTimePercent:
		return unit.newValueCurrentTimePercent(config.GetCurrentTimePercent())
	case *proto.APLValue_RemainingTime:
		return unit.newValueRemainingTime(config.GetRemainingTime())
	case *proto.APLValue_RemainingTimePercent:
		return unit.newValueRemainingTimePercent(config.GetRemainingTimePercent())

	// Resources
	case *proto.APLValue_CurrentMana:
		return unit.newValueCurrentMana(config.GetCurrentMana())
	case *proto.APLValue_CurrentManaPercent:
		return unit.newValueCurrentManaPercent(config.GetCurrentManaPercent())
	case *proto.APLValue_CurrentRage:
		return unit.newValueCurrentRage(config.GetCurrentRage())
	case *proto.APLValue_CurrentEnergy:
		return unit.newValueCurrentEnergy(config.GetCurrentEnergy())
	case *proto.APLValue_CurrentComboPoints:
		return unit.newValueCurrentComboPoints(config.GetCurrentComboPoints())

	// GCD
	case *proto.APLValue_GcdIsReady:
		return unit.newValueGCDIsReady(config.GetGcdIsReady())
	case *proto.APLValue_GcdTimeToReady:
		return unit.newValueGCDTimeToReady(config.GetGcdTimeToReady())
	
	// Spells
	case *proto.APLValue_SpellCanCast:
		return unit.newValueSpellCanCast(config.GetSpellCanCast())
	case *proto.APLValue_SpellIsReady:
		return unit.newValueSpellIsReady(config.GetSpellIsReady())
	case *proto.APLValue_SpellTimeToReady:
		return unit.newValueSpellTimeToReady(config.GetSpellTimeToReady())

	// Auras
	case *proto.APLValue_AuraIsActive:
		return unit.newValueAuraIsActive(config.GetAuraIsActive())
	case *proto.APLValue_AuraRemainingTime:
		return unit.newValueAuraRemainingTime(config.GetAuraRemainingTime())
	case *proto.APLValue_AuraNumStacks:
		return unit.newValueAuraNumStacks(config.GetAuraNumStacks())

	// Dots
	case *proto.APLValue_DotIsActive:
		return unit.newValueDotIsActive(config.GetDotIsActive())
	case *proto.APLValue_DotRemainingTime:
		return unit.newValueDotRemainingTime(config.GetDotRemainingTime())

	default:
		validationError("Unimplemented value type")
		return nil
	}
}
