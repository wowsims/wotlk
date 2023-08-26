package core

import (
	"fmt"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValueAutoTimeToNext struct {
	DefaultAPLValueImpl
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
func (value *APLValueAutoTimeToNext) String() string {
	return fmt.Sprintf("Auto Time To Next")
}
