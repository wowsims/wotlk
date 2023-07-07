package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

func (unit *Unit) aplGetDot(spellId *proto.ActionID) *Dot {
	spell := unit.GetSpell(ProtoToActionID(spellId))
	if spell == nil {
		return nil
	}

	if spell.AOEDot() != nil {
		return spell.AOEDot()
	} else {
		return spell.CurDot()
	}
}

func (unit *Unit) aplGetMultidotSpell(spellId *proto.ActionID) *Spell {
	spell := unit.GetSpell(ProtoToActionID(spellId))
	if spell == nil || spell.CurDot() == nil {
		return nil
	}
	return spell
}

type APLValueDotIsActive struct {
	defaultAPLValueImpl
	dot *Dot
}

func (unit *Unit) newValueDotIsActive(config *proto.APLValueDotIsActive) APLValue {
	dot := unit.aplGetDot(config.SpellId)
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

type APLValueDotRemainingTime struct {
	defaultAPLValueImpl
	dot *Dot
}

func (unit *Unit) newValueDotRemainingTime(config *proto.APLValueDotRemainingTime) APLValue {
	dot := unit.aplGetDot(config.SpellId)
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
