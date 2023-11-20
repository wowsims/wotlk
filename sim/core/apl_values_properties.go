package core

import (
	"time"

	"github.com/wowsims/classic/sim/core/proto"
)

type APLValueChannelClipDelay struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueChannelClipDelay(config *proto.APLValueChannelClipDelay) APLValue {
	return &APLValueChannelClipDelay{
		unit: rot.unit,
	}
}
func (value *APLValueChannelClipDelay) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueChannelClipDelay) GetDuration(sim *Simulation) time.Duration {
	return value.unit.ChannelClipDelay
}
func (value *APLValueChannelClipDelay) String() string {
	return "Channel Clip Delay()"
}

type APLValueFrontOfTarget struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueFrontOfTarget(config *proto.APLValueFrontOfTarget) APLValue {
	return &APLValueFrontOfTarget{
		unit: rot.unit,
	}
}
func (value *APLValueFrontOfTarget) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueFrontOfTarget) GetBool(sim *Simulation) bool {
	return value.unit.PseudoStats.InFrontOfTarget
}
func (value *APLValueFrontOfTarget) String() string {
	return "Front of Target()"
}
