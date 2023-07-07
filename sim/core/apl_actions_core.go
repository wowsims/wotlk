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
func (action *APLActionCastSpell) GetInnerActions() []*APLAction { return nil }
func (action *APLActionCastSpell) Finalize()                     {}
func (action *APLActionCastSpell) Reset(*Simulation)             {}
func (action *APLActionCastSpell) IsReady(sim *Simulation) bool {
	return action.spell.CanCast(sim, action.spell.Unit.CurrentTarget)
}
func (action *APLActionCastSpell) Execute(sim *Simulation) {
	action.spell.Cast(sim, action.spell.Unit.CurrentTarget)
}

type APLActionMultidot struct {
	spell *Spell
	maxDots int32
	refreshWindow time.Duration

	nextTarget *Unit
}

func (unit *Unit) newActionMultidot(config *proto.APLActionMultidot) APLActionImpl {
	spell := unit.aplGetMultidotSpell(config.SpellId)

	refreshWindow := time.Duration(0)
	canRollover := false
	if canRollover {
		refreshWindow = time.Second*3
	}
	return &APLActionMultidot{
		spell: spell,
		maxDots: MinInt32(config.MaxDots, unit.Env.GetNumTargets()),
		refreshWindow: refreshWindow,
	}
}
func (action *APLActionMultidot) GetInnerActions() []*APLAction { return nil }
func (action *APLActionMultidot) Finalize()                     {}
func (action *APLActionMultidot) Reset(*Simulation) {
	action.nextTarget = nil
}
func (action *APLActionMultidot) IsReady(sim *Simulation) bool {
	for i := 0; i < action.maxDots; i++ {
		target := sim.Encounter.GetTarget(i)
		dot := action.spell.Dot(target)
		shouldPreserveSnapshot := dot.
		if dot.IsActive() && dot.RemainingDuration(sim) < time.Second*3 {
			action.nextTarget = target
			return true
		}
	}
	return false
}
func (action *APLActionMultidot) Execute(sim *Simulation) {
	action.spell.Cast(sim, action.nextTarget)
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
