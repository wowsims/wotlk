package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValueAutoTimeToNext struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueAutoTimeToNext(_ *proto.APLValueAutoTimeToNext) APLValue {
	return &APLValueAutoTimeToNext{
		unit: rot.unit,
	}
}
func (value *APLValueAutoTimeToNext) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueAutoTimeToNext) GetDuration(sim *Simulation) time.Duration {
	return max(0, value.unit.AutoAttacks.NextAttackAt()-sim.CurrentTime)
}
func (value *APLValueAutoTimeToNext) String() string {
	return "Auto Time To Next"
}
