package core

import (
	"fmt"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValueBossSpellIsCasting struct {
	DefaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueBossSpellIsCasting(config *proto.APLValueBossSpellIsCasting) APLValue {
	spell := rot.GetTargetAPLSpell(config.SpellId, rot.GetTargetUnit(config.TargetUnit))
	if spell == nil {
		return nil
	}
	return &APLValueBossSpellIsCasting{
		spell: spell,
	}
}
func (value *APLValueBossSpellIsCasting) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueBossSpellIsCasting) GetBool(sim *Simulation) bool {
	return value.spell.Unit.Hardcast.ActionID == value.spell.ActionID && value.spell.Unit.Hardcast.Expires > sim.CurrentTime
}
func (value *APLValueBossSpellIsCasting) String() string {
	return fmt.Sprintf("Boss is Casting(%s)", value.spell.ActionID)
}

type APLValueBossSpellTimeToReady struct {
	DefaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueBossSpellTimeToReady(config *proto.APLValueBossSpellTimeToReady) APLValue {
	spell := rot.GetTargetAPLSpell(config.SpellId, rot.GetTargetUnit(config.TargetUnit))
	if spell == nil {
		return nil
	}
	return &APLValueBossSpellTimeToReady{
		spell: spell,
	}
}
func (value *APLValueBossSpellTimeToReady) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueBossSpellTimeToReady) GetDuration(sim *Simulation) time.Duration {
	return value.spell.TimeToReady(sim)
}
func (value *APLValueBossSpellTimeToReady) String() string {
	return fmt.Sprintf("Boss Spell Time to Ready(%s)", value.spell.ActionID)
}
