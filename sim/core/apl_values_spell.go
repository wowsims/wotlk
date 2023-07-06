package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

func (unit *Unit) aplGetSpell(spellId *proto.ActionID) *Spell {
	return unit.GetSpell(ProtoToActionID(spellId))
}

type APLValueSpellCanCast struct {
	defaultAPLValueImpl
	spell *Spell
}

func (unit *Unit) newValueSpellCanCast(config *proto.APLValueSpellCanCast) APLValue {
	spell := unit.aplGetSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellCanCast{
		spell: spell,
	}
}
func (value *APLValueSpellCanCast) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueSpellCanCast) GetBool(sim *Simulation) bool {
	return value.spell.CanCast(sim, value.spell.Unit.CurrentTarget)
}

type APLValueSpellIsReady struct {
	defaultAPLValueImpl
	spell *Spell
}

func (unit *Unit) newValueSpellIsReady(config *proto.APLValueSpellIsReady) APLValue {
	spell := unit.aplGetSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellIsReady{
		spell: spell,
	}
}
func (value *APLValueSpellIsReady) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueSpellIsReady) GetBool(sim *Simulation) bool {
	return value.spell.IsReady(sim)
}

type APLValueSpellTimeToReady struct {
	defaultAPLValueImpl
	spell *Spell
}

func (unit *Unit) newValueSpellTimeToReady(config *proto.APLValueSpellTimeToReady) APLValue {
	spell := unit.aplGetSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellTimeToReady{
		spell: spell,
	}
}
func (value *APLValueSpellTimeToReady) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueSpellTimeToReady) GetDuration(sim *Simulation) time.Duration {
	return value.spell.TimeToReady(sim)
}
