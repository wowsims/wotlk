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

	// Dots
	case *proto.APLValue_DotIsActive:
		return unit.newValueDotIsActive(config.GetDotIsActive())

	default:
		validationError("Unimplemented value type")
		return nil
	}
}
