package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValueSpellCanCast struct {
	defaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellCanCast(config *proto.APLValueSpellCanCast) APLValue {
	spell := rot.aplGetSpell(config.SpellId)
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

func (rot *APLRotation) newValueSpellIsReady(config *proto.APLValueSpellIsReady) APLValue {
	spell := rot.aplGetSpell(config.SpellId)
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

func (rot *APLRotation) newValueSpellTimeToReady(config *proto.APLValueSpellTimeToReady) APLValue {
	spell := rot.aplGetSpell(config.SpellId)
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

type APLValueSpellCastTime struct {
	defaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellCastTime(config *proto.APLValueSpellCastTime) APLValue {
	spell := rot.aplGetSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellCastTime{
		spell: spell,
	}
}
func (value *APLValueSpellCastTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueSpellCastTime) GetDuration(sim *Simulation) time.Duration {
	return value.spell.Unit.ApplyCastSpeedForSpell(value.spell.DefaultCast.CastTime, value.spell)
}

type APLValueSpellChannelTime struct {
	defaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellChannelTime(config *proto.APLValueSpellChannelTime) APLValue {
	spell := rot.aplGetSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellChannelTime{
		spell: spell,
	}
}
func (value *APLValueSpellChannelTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueSpellChannelTime) GetDuration(sim *Simulation) time.Duration {
	return value.spell.Unit.ApplyCastSpeedForSpell(value.spell.DefaultCast.ChannelTime, value.spell)
}

type APLValueSpellTravelTime struct {
	defaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellTravelTime(config *proto.APLValueSpellTravelTime) APLValue {
	spell := rot.aplGetSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellTravelTime{
		spell: spell,
	}
}
func (value *APLValueSpellTravelTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueSpellTravelTime) GetDuration(sim *Simulation) time.Duration {
	return time.Duration(float64(time.Second) * value.spell.Unit.DistanceFromTarget / value.spell.MissileSpeed)
}
