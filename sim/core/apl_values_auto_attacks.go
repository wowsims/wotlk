package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValueAutoTimeToNext struct {
	defaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueAutoTimeToNext(config *proto.APLValueAutoTimeToNext) APLValue {
	return &APLValueAutoTimeToNext{
		unit: rot.unit,
	}
}
func (value *APLValueAutoTimeToNext) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueAutoTimeToNext) GetDuration(sim *Simulation) time.Duration {
	return MaxDuration(0, value.unit.AutoAttacks.NextAttackAt()-sim.CurrentTime)
}
