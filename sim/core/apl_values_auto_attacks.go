package core

import (
	"time"

	"github.com/wowsims/sod/sim/core/proto"
)

type APLValueAutoTimeToNext struct {
	DefaultAPLValueImpl
	unit     *Unit
	autoType proto.APLValueAutoTimeToNext_AutoType
}

func (rot *APLRotation) newValueAutoTimeToNext(config *proto.APLValueAutoTimeToNext) APLValue {
	return &APLValueAutoTimeToNext{
		unit:     rot.unit,
		autoType: config.AutoType,
	}
}
func (value *APLValueAutoTimeToNext) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueAutoTimeToNext) GetDuration(sim *Simulation) time.Duration {
	switch value.autoType {
	case proto.APLValueAutoTimeToNext_Melee:
		return max(0, value.unit.AutoAttacks.NextAttackAt()-sim.CurrentTime)
	case proto.APLValueAutoTimeToNext_MainHand:
		return max(0, value.unit.AutoAttacks.MainhandSwingAt()-sim.CurrentTime)
	case proto.APLValueAutoTimeToNext_OffHand:
		return max(0, value.unit.AutoAttacks.OffhandSwingAt()-sim.CurrentTime)
	case proto.APLValueAutoTimeToNext_Ranged:
		return max(0, value.unit.AutoAttacks.NextRangedAttackAt()-sim.CurrentTime)
	}
	// defaults to Any
	return max(0, value.unit.AutoAttacks.NextAnyAttackAt()-sim.CurrentTime)
}
func (value *APLValueAutoTimeToNext) String() string {
	return "Auto Time To Next"
}
