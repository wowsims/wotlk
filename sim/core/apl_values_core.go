package core

import (
	"fmt"
	"time"

	"github.com/wowsims/classic/sim/core/proto"
)

type APLValueDotIsActive struct {
	DefaultAPLValueImpl
	dot *Dot
}

func (rot *APLRotation) newValueDotIsActive(config *proto.APLValueDotIsActive) APLValue {
	dot := rot.GetAPLDot(rot.GetTargetUnit(config.TargetUnit), config.SpellId)
	if dot == nil {
		return nil
	}
	return &APLValueDotIsActive{
		dot: dot,
	}
}
func (value *APLValueDotIsActive) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueDotIsActive) GetBool(sim *Simulation) bool {
	return value.dot.IsActive()
}
func (value *APLValueDotIsActive) String() string {
	return fmt.Sprintf("Dot Is Active(%s)", value.dot.Spell.ActionID)
}

type APLValueDotRemainingTime struct {
	DefaultAPLValueImpl
	dot *Dot
}

func (rot *APLRotation) newValueDotRemainingTime(config *proto.APLValueDotRemainingTime) APLValue {
	dot := rot.GetAPLDot(rot.GetTargetUnit(config.TargetUnit), config.SpellId)
	if dot == nil {
		return nil
	}
	return &APLValueDotRemainingTime{
		dot: dot,
	}
}
func (value *APLValueDotRemainingTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueDotRemainingTime) GetDuration(sim *Simulation) time.Duration {
	return value.dot.RemainingDuration(sim)
}
func (value *APLValueDotRemainingTime) String() string {
	return fmt.Sprintf("Dot Remaining Time(%s)", value.dot.Spell.ActionID)
}
