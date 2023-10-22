package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValueConst struct {
	DefaultAPLValueImpl
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

	if strings.ToLower(config.Val) == "true" {
		result.boolVal = true
		result.valType = proto.APLValueType_ValueTypeBool
		return result
	} else if strings.ToLower(config.Val) == "false" {
		result.boolVal = false
		result.valType = proto.APLValueType_ValueTypeBool
		return result
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
func (value *APLValueConst) GetBool(_ *Simulation) bool {
	return value.boolVal
}
func (value *APLValueConst) GetInt(_ *Simulation) int32 {
	return value.intVal
}
func (value *APLValueConst) GetFloat(_ *Simulation) float64 {
	return value.floatVal
}
func (value *APLValueConst) GetDuration(_ *Simulation) time.Duration {
	return value.durationVal
}
func (value *APLValueConst) GetString(_ *Simulation) string {
	return value.stringVal
}
func (value *APLValueConst) String() string {
	return value.stringVal
}

type APLValueCoerced struct {
	DefaultAPLValueImpl
	valueType proto.APLValueType
	inner     APLValue
}

func (value *APLValueCoerced) GetInnerValues() []APLValue {
	return []APLValue{value.inner}
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
func (value *APLValueCoerced) String() string {
	return value.inner.String()
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

func higherOrderType(type1 proto.APLValueType, type2 proto.APLValueType) proto.APLValueType {
	for _, listType := range aplValueTypeOrder {
		if listType == type1 {
			return type2
		} else if listType == type2 {
			return type1
		}
	}
	return type1
}
func highestOrderTypeList(values []APLValue) proto.APLValueType {
	coercionType := aplValueTypeOrder[0]
	for _, val := range values {
		if val != nil {
			coercionType = higherOrderType(coercionType, val.Type())
		}
	}
	return coercionType
}
func (rot *APLRotation) coerceAllToSameType(values []APLValue) []APLValue {
	coercionType := highestOrderTypeList(values)
	return MapSlice(values, func(val APLValue) APLValue { return rot.coerceTo(val, coercionType) })
}

// Coerces 2 values into the same type, returning the two new values.
func (rot *APLRotation) coerceToSameType(value1 APLValue, value2 APLValue) (APLValue, APLValue) {
	coerced := rot.coerceAllToSameType([]APLValue{value1, value2})
	return coerced[0], coerced[1]
}

// Utility function which returns the constant float value of a Const or Coerced(Const) APL value.
// Returns -1 if the value is not a constant, or does not have a float value.
func getConstAPLFloatValue(value APLValue) float64 {
	if constValue, isConst := value.(*APLValueConst); isConst {
		return constValue.GetFloat(nil)
	} else if coercedValue, isCoerced := value.(*APLValueCoerced); isCoerced {
		if _, innerIsConst := coercedValue.inner.(*APLValueConst); innerIsConst {
			return coercedValue.GetFloat(nil)
		}
	}
	return -1
}

type APLValueCompare struct {
	DefaultAPLValueImpl
	op  proto.APLValueCompare_ComparisonOperator
	lhs APLValue
	rhs APLValue
}

func (rot *APLRotation) newValueCompare(config *proto.APLValueCompare) APLValue {
	lhs, rhs := rot.coerceToSameType(rot.newAPLValue(config.Lhs), rot.newAPLValue(config.Rhs))
	if lhs == nil || rhs == nil {
		return nil
	}

	if lhs.Type() == proto.APLValueType_ValueTypeBool && !(config.Op == proto.APLValueCompare_OpEq || config.Op == proto.APLValueCompare_OpNe) {
		rot.ValidationWarning("Bool types only allow Equals and NotEquals comparisons!")
		return nil
	}
	return &APLValueCompare{
		op:  config.Op,
		lhs: lhs,
		rhs: rhs,
	}
}
func (value *APLValueCompare) GetInnerValues() []APLValue {
	return []APLValue{value.lhs, value.rhs}
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
func (value *APLValueCompare) String() string {
	return fmt.Sprintf("%s %s %s", value.lhs, value.op, value.rhs)
}

type APLValueMath struct {
	DefaultAPLValueImpl
	op  proto.APLValueMath_MathOperator
	lhs APLValue
	rhs APLValue
}

func (rot *APLRotation) newValueMath(config *proto.APLValueMath) APLValue {
	lhs, rhs := rot.newAPLValue(config.Lhs), rot.newAPLValue(config.Rhs)
	if config.Op == proto.APLValueMath_OpAdd || config.Op == proto.APLValueMath_OpSub {
		lhs, rhs = rot.coerceToSameType(lhs, rhs)
	}
	if lhs == nil || rhs == nil {
		return nil
	}

	if lhs.Type() == proto.APLValueType_ValueTypeBool || rhs.Type() == proto.APLValueType_ValueTypeBool {
		rot.ValidationWarning("Bool types not allowed in Math Operations!")
		return nil
	}

	if lhs.Type() == proto.APLValueType_ValueTypeString || rhs.Type() == proto.APLValueType_ValueTypeString {
		rot.ValidationWarning("String types not allowed in Math Operations!")
		return nil
	}

	return &APLValueMath{
		op:  config.Op,
		lhs: lhs,
		rhs: rhs,
	}
}
func (value *APLValueMath) GetInnerValues() []APLValue {
	return []APLValue{value.lhs, value.rhs}
}
func (value *APLValueMath) Type() proto.APLValueType {
	return value.lhs.Type()
}
func (value *APLValueMath) GetInt(sim *Simulation) int32 {
	switch value.op {
	case proto.APLValueMath_OpAdd:
		return value.lhs.GetInt(sim) + value.rhs.GetInt(sim)
	case proto.APLValueMath_OpSub:
		return value.lhs.GetInt(sim) - value.rhs.GetInt(sim)
	case proto.APLValueMath_OpMul:
		return value.lhs.GetInt(sim) * value.rhs.GetInt(sim)
	case proto.APLValueMath_OpDiv:
		return value.lhs.GetInt(sim) / value.rhs.GetInt(sim)
	}
	return 0
}
func (value *APLValueMath) GetFloat(sim *Simulation) float64 {
	switch value.op {
	case proto.APLValueMath_OpAdd:
		return value.lhs.GetFloat(sim) + value.rhs.GetFloat(sim)
	case proto.APLValueMath_OpSub:
		return value.lhs.GetFloat(sim) - value.rhs.GetFloat(sim)
	case proto.APLValueMath_OpMul:
		return value.lhs.GetFloat(sim) * value.rhs.GetFloat(sim)
	case proto.APLValueMath_OpDiv:
		return value.lhs.GetFloat(sim) / value.rhs.GetFloat(sim)
	}
	return 0
}
func (value *APLValueMath) GetDuration(sim *Simulation) time.Duration {
	switch value.op {
	case proto.APLValueMath_OpAdd:
		return value.lhs.GetDuration(sim) + value.rhs.GetDuration(sim)
	case proto.APLValueMath_OpSub:
		return value.lhs.GetDuration(sim) - value.rhs.GetDuration(sim)
	case proto.APLValueMath_OpMul:
		left := value.lhs.GetDuration(sim)
		right := value.rhs.GetDuration(sim)

		switch value.lhs.Type() {
		case proto.APLValueType_ValueTypeInt:
			left = time.Duration(value.lhs.GetInt(sim))
		case proto.APLValueType_ValueTypeFloat:
			left = time.Duration(value.lhs.GetFloat(sim))
		}

		switch value.rhs.Type() {
		case proto.APLValueType_ValueTypeInt:
			right = time.Duration(value.rhs.GetInt(sim))
		case proto.APLValueType_ValueTypeFloat:
			right = time.Duration(value.rhs.GetFloat(sim))
		}
		return left * right
	case proto.APLValueMath_OpDiv:
		divider := value.rhs.GetDuration(sim)
		if value.rhs.Type() == proto.APLValueType_ValueTypeFloat {
			divider = time.Duration(value.rhs.GetFloat(sim))
		} else if value.rhs.Type() == proto.APLValueType_ValueTypeInt {
			divider = time.Duration(value.rhs.GetInt(sim))
		}
		return value.lhs.GetDuration(sim) / divider
	}
	return 0
}
func (value *APLValueMath) String() string {
	return fmt.Sprintf("Math(%s %s %s)", value.lhs, value.op, value.rhs)
}

type APLValueMax struct {
	DefaultAPLValueImpl
	vals []APLValue
}

func (rot *APLRotation) newValueMax(config *proto.APLValueMax) APLValue {
	vals := MapSlice(config.Vals, func(val *proto.APLValue) APLValue {
		return rot.newAPLValue(val)
	})
	vals = rot.coerceAllToSameType(vals)
	vals = FilterSlice(vals, func(val APLValue) bool { return val != nil })
	if len(vals) == 0 {
		return nil
	} else if len(vals) == 1 {
		return vals[0]
	}
	return &APLValueMax{
		vals: vals,
	}
}
func (value *APLValueMax) GetInnerValues() []APLValue {
	return value.vals
}
func (value *APLValueMax) Type() proto.APLValueType {
	return value.vals[0].Type()
}
func (value *APLValueMax) GetInt(sim *Simulation) int32 {
	result := value.vals[0].GetInt(sim)
	for i := 1; i < len(value.vals); i++ {
		result = max(result, value.vals[i].GetInt(sim))
	}
	return result
}
func (value *APLValueMax) GetFloat(sim *Simulation) float64 {
	result := value.vals[0].GetFloat(sim)
	for i := 1; i < len(value.vals); i++ {
		result = max(result, value.vals[i].GetFloat(sim))
	}
	return result
}
func (value *APLValueMax) GetDuration(sim *Simulation) time.Duration {
	result := value.vals[0].GetDuration(sim)
	for i := 1; i < len(value.vals); i++ {
		result = max(result, value.vals[i].GetDuration(sim))
	}
	return result
}
func (value *APLValueMax) String() string {
	return fmt.Sprintf("Max(%s)", strings.Join(MapSlice(value.vals, func(subvalue APLValue) string { return fmt.Sprintf("(%s)", subvalue) }), ", "))
}

type APLValueMin struct {
	DefaultAPLValueImpl
	vals []APLValue
}

func (rot *APLRotation) newValueMin(config *proto.APLValueMin) APLValue {
	vals := MapSlice(config.Vals, func(val *proto.APLValue) APLValue {
		return rot.newAPLValue(val)
	})
	vals = rot.coerceAllToSameType(vals)
	vals = FilterSlice(vals, func(val APLValue) bool { return val != nil })
	if len(vals) == 0 {
		return nil
	} else if len(vals) == 1 {
		return vals[0]
	}
	return &APLValueMin{
		vals: vals,
	}
}
func (value *APLValueMin) GetInnerValues() []APLValue {
	return value.vals
}
func (value *APLValueMin) Type() proto.APLValueType {
	return value.vals[0].Type()
}
func (value *APLValueMin) GetInt(sim *Simulation) int32 {
	result := value.vals[0].GetInt(sim)
	for _, v := range value.vals[1:] {
		result = min(result, v.GetInt(sim))
	}
	return result
}
func (value *APLValueMin) GetFloat(sim *Simulation) float64 {
	result := value.vals[0].GetFloat(sim)
	for _, v := range value.vals[1:] {
		result = min(result, v.GetFloat(sim))
	}
	return result
}
func (value *APLValueMin) GetDuration(sim *Simulation) time.Duration {
	result := value.vals[0].GetDuration(sim)
	for _, v := range value.vals[1:] {
		result = min(result, v.GetDuration(sim))
	}
	return result
}
func (value *APLValueMin) String() string {
	return fmt.Sprintf("Min(%s)", strings.Join(MapSlice(value.vals, func(subvalue APLValue) string { return fmt.Sprintf("(%s)", subvalue) }), ", "))
}

type APLValueAnd struct {
	DefaultAPLValueImpl
	vals []APLValue
}

func (rot *APLRotation) newValueAnd(config *proto.APLValueAnd) APLValue {
	vals := MapSlice(config.Vals, func(val *proto.APLValue) APLValue {
		return rot.coerceTo(rot.newAPLValue(val), proto.APLValueType_ValueTypeBool)
	})
	vals = FilterSlice(vals, func(val APLValue) bool { return val != nil })
	if len(vals) == 0 {
		return nil
	} else if len(vals) == 1 {
		return vals[0]
	}
	return &APLValueAnd{
		vals: vals,
	}
}
func (value *APLValueAnd) GetInnerValues() []APLValue {
	return value.vals
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
func (value *APLValueAnd) String() string {
	return strings.Join(MapSlice(value.vals, func(subvalue APLValue) string { return fmt.Sprintf("(%s)", subvalue) }), " AND ")
}

type APLValueOr struct {
	DefaultAPLValueImpl
	vals []APLValue
}

func (rot *APLRotation) newValueOr(config *proto.APLValueOr) APLValue {
	vals := MapSlice(config.Vals, func(val *proto.APLValue) APLValue {
		return rot.coerceTo(rot.newAPLValue(val), proto.APLValueType_ValueTypeBool)
	})
	vals = FilterSlice(vals, func(val APLValue) bool { return val != nil })
	if len(vals) == 0 {
		return nil
	} else if len(vals) == 1 {
		return vals[0]
	}
	return &APLValueOr{
		vals: vals,
	}
}
func (value *APLValueOr) GetInnerValues() []APLValue {
	return value.vals
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
func (value *APLValueOr) String() string {
	return strings.Join(MapSlice(value.vals, func(subvalue APLValue) string { return fmt.Sprintf("(%s)", subvalue) }), " OR ")
}

type APLValueNot struct {
	DefaultAPLValueImpl
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
func (value *APLValueNot) GetInnerValues() []APLValue {
	return []APLValue{value.val}
}
func (value *APLValueNot) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueNot) GetBool(sim *Simulation) bool {
	return !value.val.GetBool(sim)
}
func (value *APLValueNot) String() string {
	return fmt.Sprintf("Not(%s)", value.val)
}
