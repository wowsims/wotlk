package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rot *APLRotation) aplGetSpell(spellId *proto.ActionID) *Spell {
	spell := rot.unit.GetSpell(ProtoToActionID(spellId))
	if spell == nil {
		rot.validationWarning("%s does not know spell %s", rot.unit.Label, ProtoToActionID(spellId))
	}
	return spell
}

func (rot *APLRotation) aplGetDot(spellId *proto.ActionID) *Dot {
	spell := rot.aplGetSpell(spellId)

	if spell == nil {
		return nil
	} else if spell.AOEDot() != nil {
		return spell.AOEDot()
	} else {
		return spell.CurDot()
	}
}

func (rot *APLRotation) aplGetMultidotSpell(spellId *proto.ActionID) *Spell {
	spell := rot.aplGetSpell(spellId)
	if spell == nil {
		return nil
	} else if spell.CurDot() == nil {
		rot.validationWarning("Spell %s does not have an associated DoT", ProtoToActionID(spellId))
		return nil
	}
	return spell
}

type APLValueDotIsActive struct {
	defaultAPLValueImpl
	dot *Dot
}

func (rot *APLRotation) newValueDotIsActive(config *proto.APLValueDotIsActive) APLValue {
	dot := rot.aplGetDot(config.SpellId)
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

func (rot *APLRotation) newValueDotRemainingTime(config *proto.APLValueDotRemainingTime) APLValue {
	dot := rot.aplGetDot(config.SpellId)
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
