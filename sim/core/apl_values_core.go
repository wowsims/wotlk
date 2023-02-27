package core

import (
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

type APLValueDotIsActive struct {
	defaultAPLValueImpl
	dot *Dot
}

func (unit *Unit) newValueDotIsActive(config *proto.APLValueDotIsActive) APLValue {
	return &APLValueDotIsActive{
		dot: unit.aplGetDot(config.SpellId),
	}
}
func (value *APLValueDotIsActive) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueDotIsActive) GetBool(sim *Simulation) bool {
	return value.dot.IsActive()
}
