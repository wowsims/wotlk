package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValueBool struct {
	defaultAPLValueImpl
	val bool
}

func (unit *Unit) newValueBool(config *proto.APLValueBool) APLValue {
	return &APLValueBool{
		val: config.Val,
	}
}
func (value *APLValueBool) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueBool) GetBool(sim *Simulation) bool {
	return value.val
}

type APLValueInt struct {
	defaultAPLValueImpl
	val int32
}

func (unit *Unit) newValueInt(config *proto.APLValueInt) APLValue {
	return &APLValueInt{
		val: config.Val,
	}
}
func (value *APLValueInt) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueInt) GetInt(sim *Simulation) int32 {
	return value.val
}

type APLValueFloat struct {
	defaultAPLValueImpl
	val float64
}

func (unit *Unit) newValueFloat(config *proto.APLValueFloat) APLValue {
	return &APLValueFloat{
		val: config.Val,
	}
}
func (value *APLValueFloat) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueFloat) GetFloat(sim *Simulation) float64 {
	return value.val
}

type APLValueDuration struct {
	defaultAPLValueImpl
	val time.Duration
}

func (unit *Unit) newValueDuration(config *proto.APLValueDuration) APLValue {
	return &APLValueDuration{
		val: DurationFromProto(config.Val),
	}
}
func (value *APLValueDuration) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueDuration) GetDuration(sim *Simulation) time.Duration {
	return value.val
}

type APLValueString struct {
	defaultAPLValueImpl
	val string
}

func (unit *Unit) newValueString(config *proto.APLValueString) APLValue {
	return &APLValueString{
		val: config.Val,
	}
}
func (value *APLValueString) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeString
}
func (value *APLValueString) GetString(sim *Simulation) string {
	return value.val
}
