package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValueConst struct {
	valType proto.APLValueType

	intVal      int32
	floatVal    float64
	durationVal time.Duration
	stringVal   string
	boolVal     bool
}

func (rot *APLRotation) newValueConst(config *proto.APLValueConst) APLValue {
	result := &APLValueConst{
		valType:   proto.APLValueType_ValueTypeString,
		stringVal: config.Val,
		boolVal:   config.Val != "",
	}

	if durVal, err := time.ParseDuration(config.Val); err == nil {
		result.durationVal = durVal
		result.valType = proto.APLValueType_ValueTypeDuration
		return result
	}

	if intVal, err := strconv.Atoi(config.Val); err == nil {
		result.intVal = int32(intVal)
		result.floatVal = float64(result.intVal)
		result.durationVal = DurationFromSeconds(result.floatVal)
		result.valType = proto.APLValueType_ValueTypeInt
		return result
	}

	if len(config.Val) > 1 && config.Val[len(config.Val)-1] == '%' {
		if floatVal, err := strconv.ParseFloat(config.Val[0:len(config.Val)-1], 64); err == nil {
			result.floatVal = floatVal / 100.0
			result.durationVal = DurationFromSeconds(floatVal / 100.0)
			result.valType = proto.APLValueType_ValueTypeFloat
			return result
		}
	}

	if floatVal, err := strconv.ParseFloat(config.Val, 64); err == nil {
		result.floatVal = floatVal
		result.durationVal = DurationFromSeconds(floatVal)
		result.valType = proto.APLValueType_ValueTypeFloat
		return result
	}
	return result
}
func (value *APLValueConst) Type() proto.APLValueType {
	return value.valType
}
func (value *APLValueConst) GetBool(sim *Simulation) bool {
	return value.boolVal
}
func (value *APLValueConst) GetInt(sim *Simulation) int32 {
	return value.intVal
}
func (value *APLValueConst) GetFloat(sim *Simulation) float64 {
	return value.floatVal
}
func (value *APLValueConst) GetDuration(sim *Simulation) time.Duration {
	return value.durationVal
}
func (value *APLValueConst) GetString(sim *Simulation) string {
	return value.stringVal
}

type APLValueCoerced struct {
	valueType proto.APLValueType
	inner     APLValue
}

func (value *APLValueCoerced) Type() proto.APLValueType {
	return value.valueType
}
func (value *APLValueCoerced) GetBool(sim *Simulation) bool {
	switch value.inner.Type() {
	case proto.APLValueType_ValueTypeBool:
		return value.inner.GetBool(sim)
	case proto.APLValueType_ValueTypeInt:
		return value.inner.GetInt(sim) != 0
	case proto.APLValueType_ValueTypeFloat:
		return value.inner.GetFloat(sim) != 0
	case proto.APLValueType_ValueTypeDuration:
		return value.inner.GetDuration(sim) != 0
	case proto.APLValueType_ValueTypeString:
		return value.inner.GetString(sim) != ""
	}
	return false
}
func (value APLValueCoerced) GetInt(sim *Simulation) int32 {
	switch value.inner.Type() {
	case proto.APLValueType_ValueTypeBool:
		if value.inner.GetBool(sim) {
			return 1
		} else {
			return 0
		}
	case proto.APLValueType_ValueTypeInt:
		return value.inner.GetInt(sim)
	case proto.APLValueType_ValueTypeFloat:
		return int32(value.inner.GetFloat(sim))
	case proto.APLValueType_ValueTypeDuration:
		return int32(value.inner.GetDuration(sim).Seconds())
	case proto.APLValueType_ValueTypeString:
		panic("Cannot coerce string to int")
	}
	return 0
}
func (value APLValueCoerced) GetFloat(sim *Simulation) float64 {
	switch value.inner.Type() {
	case proto.APLValueType_ValueTypeBool:
		if value.inner.GetBool(sim) {
			return 1
		} else {
			return 0
		}
	case proto.APLValueType_ValueTypeInt:
		return float64(value.inner.GetInt(sim))
	case proto.APLValueType_ValueTypeFloat:
		return value.inner.GetFloat(sim)
	case proto.APLValueType_ValueTypeDuration:
		return value.inner.GetDuration(sim).Seconds()
	case proto.APLValueType_ValueTypeString:
		panic("Cannot coerce string to float")
	}
	return 0
}
func (value APLValueCoerced) GetDuration(sim *Simulation) time.Duration {
	switch value.inner.Type() {
	case proto.APLValueType_ValueTypeBool:
		panic("Cannot coerce bool to duration")
	case proto.APLValueType_ValueTypeInt:
		return time.Second * time.Duration(value.inner.GetInt(sim))
	case proto.APLValueType_ValueTypeFloat:
		return DurationFromSeconds(value.inner.GetFloat(sim))
	case proto.APLValueType_ValueTypeDuration:
		return value.inner.GetDuration(sim)
	case proto.APLValueType_ValueTypeString:
		panic("Cannot coerce string to duration")
	}
	return 0
}
func (value APLValueCoerced) GetString(sim *Simulation) string {
	switch value.inner.Type() {
	case proto.APLValueType_ValueTypeBool:
		panic("Cannot coerce bool to string")
	case proto.APLValueType_ValueTypeInt:
		return strconv.Itoa(int(value.inner.GetInt(sim)))
	case proto.APLValueType_ValueTypeFloat:
		return fmt.Sprintf("%.3f", value.inner.GetFloat(sim))
	case proto.APLValueType_ValueTypeDuration:
		return value.inner.GetDuration(sim).String()
	case proto.APLValueType_ValueTypeString:
		return value.inner.GetString(sim)
	}
	return ""
}

// Wraps a value so that it is converted into a Boolean.
func (rot *APLRotation) coerceTo(value APLValue, newType proto.APLValueType) APLValue {
	if value == nil {
		return nil
	} else if value.Type() == newType {
		return value
	} else if constVal, ok := value.(*APLValueConst); ok {
		// For the special case of APLValueConst, we can skip the wrapper and
		// simply make a copy with a different type.
		newVal := &APLValueConst{}
		*newVal = *constVal
		newVal.valType = newType
		return newVal
	} else {
		return &APLValueCoerced{
			valueType: newType,
			inner:     value,
		}
	}
}

// Types that come later in the list are higher 'priority'.
var aplValueTypeOrder = []proto.APLValueType{
	proto.APLValueType_ValueTypeInt,
	proto.APLValueType_ValueTypeFloat,
	proto.APLValueType_ValueTypeDuration,
	proto.APLValueType_ValueTypeString,
	proto.APLValueType_ValueTypeBool,
}

// Coerces 2 values into the same type, returning the two new values.
func (rot *APLRotation) coerceToSameType(value1 APLValue, value2 APLValue) (APLValue, APLValue) {
	var coercionType proto.APLValueType
	for _, listType := range aplValueTypeOrder {
		if value1.Type() == listType || value2.Type() == listType {
			coercionType = listType
		}
	}
	return rot.coerceTo(value1, coercionType), rot.coerceTo(value2, coercionType)
}

type APLValueCompare struct {
	defaultAPLValueImpl
	op  proto.APLValueCompare_ComparisonOperator
	lhs APLValue
	rhs APLValue
}

func (rot *APLRotation) newValueCompare(config *proto.APLValueCompare) APLValue {
	lhs, rhs := rot.coerceToSameType(rot.newAPLValue(config.Lhs), rot.newAPLValue(config.Rhs))
	if lhs.Type() == proto.APLValueType_ValueTypeBool && !(config.Op == proto.APLValueCompare_OpEq || config.Op == proto.APLValueCompare_OpNe) {
		rot.validationWarning("Bool types only allow Equals and NotEquals comparisons!")
		return nil
	}
	return &APLValueCompare{
		op:  config.Op,
		lhs: lhs,
		rhs: rhs,
	}
}
func (value *APLValueCompare) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueCompare) GetBool(sim *Simulation) bool {
	switch value.lhs.Type() {
	case proto.APLValueType_ValueTypeBool:
		switch value.op {
		case proto.APLValueCompare_OpEq:
			return value.lhs.GetBool(sim) == value.rhs.GetBool(sim)
		case proto.APLValueCompare_OpNe:
			return value.lhs.GetBool(sim) != value.rhs.GetBool(sim)
		}
	case proto.APLValueType_ValueTypeInt:
		switch value.op {
		case proto.APLValueCompare_OpEq:
			return value.lhs.GetInt(sim) == value.rhs.GetInt(sim)
		case proto.APLValueCompare_OpNe:
			return value.lhs.GetInt(sim) != value.rhs.GetInt(sim)
		case proto.APLValueCompare_OpLt:
			return value.lhs.GetInt(sim) < value.rhs.GetInt(sim)
		case proto.APLValueCompare_OpLe:
			return value.lhs.GetInt(sim) <= value.rhs.GetInt(sim)
		case proto.APLValueCompare_OpGt:
			return value.lhs.GetInt(sim) > value.rhs.GetInt(sim)
		case proto.APLValueCompare_OpGe:
			return value.lhs.GetInt(sim) >= value.rhs.GetInt(sim)
		}
	case proto.APLValueType_ValueTypeFloat:
		switch value.op {
		case proto.APLValueCompare_OpEq:
			return value.lhs.GetFloat(sim) == value.rhs.GetFloat(sim)
		case proto.APLValueCompare_OpNe:
			return value.lhs.GetFloat(sim) != value.rhs.GetFloat(sim)
		case proto.APLValueCompare_OpLt:
			return value.lhs.GetFloat(sim) < value.rhs.GetFloat(sim)
		case proto.APLValueCompare_OpLe:
			return value.lhs.GetFloat(sim) <= value.rhs.GetFloat(sim)
		case proto.APLValueCompare_OpGt:
			return value.lhs.GetFloat(sim) > value.rhs.GetFloat(sim)
		case proto.APLValueCompare_OpGe:
			return value.lhs.GetFloat(sim) >= value.rhs.GetFloat(sim)
		}
	case proto.APLValueType_ValueTypeDuration:
		switch value.op {
		case proto.APLValueCompare_OpEq:
			return value.lhs.GetDuration(sim) == value.rhs.GetDuration(sim)
		case proto.APLValueCompare_OpNe:
			return value.lhs.GetDuration(sim) != value.rhs.GetDuration(sim)
		case proto.APLValueCompare_OpLt:
			return value.lhs.GetDuration(sim) < value.rhs.GetDuration(sim)
		case proto.APLValueCompare_OpLe:
			return value.lhs.GetDuration(sim) <= value.rhs.GetDuration(sim)
		case proto.APLValueCompare_OpGt:
			return value.lhs.GetDuration(sim) > value.rhs.GetDuration(sim)
		case proto.APLValueCompare_OpGe:
			return value.lhs.GetDuration(sim) >= value.rhs.GetDuration(sim)
		}
	case proto.APLValueType_ValueTypeString:
		switch value.op {
		case proto.APLValueCompare_OpEq:
			return value.lhs.GetString(sim) == value.rhs.GetString(sim)
		case proto.APLValueCompare_OpNe:
			return value.lhs.GetString(sim) != value.rhs.GetString(sim)
		case proto.APLValueCompare_OpLt:
			return value.lhs.GetString(sim) < value.rhs.GetString(sim)
		case proto.APLValueCompare_OpLe:
			return value.lhs.GetString(sim) <= value.rhs.GetString(sim)
		case proto.APLValueCompare_OpGt:
			return value.lhs.GetString(sim) > value.rhs.GetString(sim)
		case proto.APLValueCompare_OpGe:
			return value.lhs.GetString(sim) >= value.rhs.GetString(sim)
		}
	}
	return false
}

type APLValueAnd struct {
	defaultAPLValueImpl
	vals []APLValue
}

func (rot *APLRotation) newValueAnd(config *proto.APLValueAnd) APLValue {
	vals := MapSlice(config.Vals, func(val *proto.APLValue) APLValue {
		return rot.coerceTo(rot.newAPLValue(val), proto.APLValueType_ValueTypeBool)
	})
	vals = FilterSlice(vals, func(val APLValue) bool { return val != nil })
	if len(vals) == 0 {
		return nil
	}
	return &APLValueAnd{
		vals: vals,
	}
}
func (value *APLValueAnd) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAnd) GetBool(sim *Simulation) bool {
	for _, val := range value.vals {
		if !val.GetBool(sim) {
			return false
		}
	}
	return true
}

type APLValueOr struct {
	defaultAPLValueImpl
	vals []APLValue
}

func (rot *APLRotation) newValueOr(config *proto.APLValueOr) APLValue {
	vals := MapSlice(config.Vals, func(val *proto.APLValue) APLValue {
		return rot.coerceTo(rot.newAPLValue(val), proto.APLValueType_ValueTypeBool)
	})
	vals = FilterSlice(vals, func(val APLValue) bool { return val != nil })
	if len(vals) == 0 {
		return nil
	}
	return &APLValueOr{
		vals: vals,
	}
}
func (value *APLValueOr) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueOr) GetBool(sim *Simulation) bool {
	for _, val := range value.vals {
		if val.GetBool(sim) {
			return true
		}
	}
	return false
}

type APLValueNot struct {
	defaultAPLValueImpl
	val APLValue
}

func (rot *APLRotation) newValueNot(config *proto.APLValueNot) APLValue {
	val := rot.coerceTo(rot.newAPLValue(config.Val), proto.APLValueType_ValueTypeBool)
	if val == nil {
		return nil
	}
	return &APLValueNot{
		val: val,
	}
}
func (value *APLValueNot) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueNot) GetBool(sim *Simulation) bool {
	return !value.val.GetBool(sim)
}
