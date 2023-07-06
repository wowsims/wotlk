package core

import (
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
func (action *APLActionCastSpell) GetInnerActions() []*APLAction { return nil }
func (action *APLActionCastSpell) Finalize()                     {}
func (action *APLActionCastSpell) Reset(*Simulation)             {}
func (action *APLActionCastSpell) IsReady(sim *Simulation) bool {
	return action.spell.CanCast(sim, action.spell.Unit.CurrentTarget)
}
func (action *APLActionCastSpell) Execute(sim *Simulation) {
	action.spell.Cast(sim, action.spell.Unit.CurrentTarget)
}

type APLActionAutocastOtherCooldowns struct {
	character *Character

	nextReadyMCD *MajorCooldown
}

func (unit *Unit) newActionAutocastOtherCooldowns(config *proto.APLActionAutocastOtherCooldowns) APLActionImpl {
	return &APLActionAutocastOtherCooldowns{
		character: unit.Env.Raid.GetPlayerFromUnit(unit).GetCharacter(),
	}
}
func (action *APLActionAutocastOtherCooldowns) GetInnerActions() []*APLAction { return nil }
func (action *APLActionAutocastOtherCooldowns) Finalize()                     {}
func (action *APLActionAutocastOtherCooldowns) Reset(*Simulation) {
	action.nextReadyMCD = nil
}
func (action *APLActionAutocastOtherCooldowns) IsReady(sim *Simulation) bool {
	action.nextReadyMCD = action.character.getFirstReadyMCD(sim)
	return action.nextReadyMCD != nil
}
func (action *APLActionAutocastOtherCooldowns) Execute(sim *Simulation) {
	action.nextReadyMCD.tryActivateHelper(sim, action.character)
	action.character.UpdateMajorCooldowns()
}

type APLActionWait struct {
	unit     *Unit
	duration APLValue
}

func (unit *Unit) newActionWait(config *proto.APLActionWait) APLActionImpl {
	return &APLActionWait{
		unit:     unit,
		duration: unit.coerceTo(unit.newAPLValue(config.Duration), proto.APLValueType_ValueTypeDuration),
	}
}
func (action *APLActionWait) GetInnerActions() []*APLAction { return nil }
func (action *APLActionWait) Finalize()                     {}
func (action *APLActionWait) Reset(*Simulation)             {}
func (action *APLActionWait) IsReady(sim *Simulation) bool {
	return action.duration != nil
}
func (action *APLActionWait) Execute(sim *Simulation) {
	action.unit.WaitUntil(sim, sim.CurrentTime+action.duration.GetDuration(sim))
}
