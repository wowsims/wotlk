package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
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
