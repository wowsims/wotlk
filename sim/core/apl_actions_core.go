package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLActionCastSpell struct {
	spell *Spell
}

func (unit *Unit) newActionCastSpell(config *proto.APLActionCastSpell) APLActionImpl {
	spell := unit.GetSpell(ProtoToActionID(config.SpellId))
	if spell == nil {
		validationWarning("No spell found for id: %s", ProtoToActionID(config.SpellId).String())
		return nil
	}
	return &APLActionCastSpell{
		spell: spell,
	}
}
func (action *APLActionCastSpell) IsAvailable(sim *Simulation) bool {
	return action.spell.CanCast(sim, action.spell.Unit.CurrentTarget)
}
func (action *APLActionCastSpell) Execute(sim *Simulation) {
	action.spell.Cast(sim, action.spell.Unit.CurrentTarget)
}

type APLActionWait struct {
	unit     *Unit
	duration time.Duration
}

func (unit *Unit) newActionWait(config *proto.APLActionWait) APLActionImpl {
	return &APLActionWait{
		unit:     unit,
		duration: DurationFromProto(config.Duration),
	}
}
func (action *APLActionWait) IsAvailable(sim *Simulation) bool {
	return true
}
func (action *APLActionWait) Execute(sim *Simulation) {
	action.unit.WaitUntil(sim, sim.CurrentTime+action.duration)
}
