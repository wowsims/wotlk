package core

import (
	"fmt"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLActionCastSpell struct {
	defaultAPLActionImpl
	spell  *Spell
	target UnitReference
}

func (rot *APLRotation) newActionCastSpell(config *proto.APLActionCastSpell) APLActionImpl {
	spell := rot.aplGetSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	target := rot.getTargetUnit(config.Target)
	if target.Get() == nil {
		return nil
	}
	return &APLActionCastSpell{
		spell:  spell,
		target: target,
	}
}
func (action *APLActionCastSpell) IsReady(sim *Simulation) bool {
	return action.spell.CanCast(sim, action.target.Get()) && (!action.spell.Flags.Matches(SpellFlagMCD) || action.spell.Unit.GCD.IsReady(sim))
}
func (action *APLActionCastSpell) Execute(sim *Simulation) {
	action.spell.Cast(sim, action.target.Get())
}
func (action *APLActionCastSpell) String() string {
	return fmt.Sprintf("Cast Spell(%s)", action.spell.ActionID)
}

type APLActionMultidot struct {
	defaultAPLActionImpl
	spell      *Spell
	maxDots    int32
	maxOverlap APLValue

	nextTarget *Unit
}

func (rot *APLRotation) newActionMultidot(config *proto.APLActionMultidot) APLActionImpl {
	unit := rot.unit

	spell := rot.aplGetMultidotSpell(config.SpellId)
	if spell == nil {
		return nil
	}

	maxOverlap := rot.coerceTo(rot.newAPLValue(config.MaxOverlap), proto.APLValueType_ValueTypeDuration)
	if maxOverlap == nil {
		maxOverlap = rot.newValueConst(&proto.APLValueConst{Val: "0ms"})
	}

	maxDots := config.MaxDots
	numTargets := unit.Env.GetNumTargets()
	if spell.Flags.Matches(SpellFlagHelpful) {
		numTargets = int32(len(unit.Env.Raid.AllPlayerUnits))
	}
	if numTargets < maxDots {
		rot.validationWarning("Encounter only has %d targets. Using that for Max Dots instead of %d", numTargets, maxDots)
		maxDots = numTargets
	}

	return &APLActionMultidot{
		spell:      spell,
		maxDots:    maxDots,
		maxOverlap: maxOverlap,
	}
}
func (action *APLActionMultidot) Reset(*Simulation) {
	action.nextTarget = nil
}
func (action *APLActionMultidot) IsReady(sim *Simulation) bool {
	maxOverlap := action.maxOverlap.GetDuration(sim)

	if action.spell.Flags.Matches(SpellFlagHelpful) {
		for i := int32(0); i < action.maxDots; i++ {
			target := sim.Raid.AllPlayerUnits[i]
			dot := action.spell.Dot(target)
			if (!dot.IsActive() || dot.RemainingDuration(sim) < maxOverlap) && action.spell.CanCast(sim, target) {
				action.nextTarget = target
				return true
			}
		}
	} else {
		for i := int32(0); i < action.maxDots; i++ {
			target := sim.Encounter.TargetUnits[i]
			dot := action.spell.Dot(target)
			if (!dot.IsActive() || dot.RemainingDuration(sim) < maxOverlap) && action.spell.CanCast(sim, target) {
				action.nextTarget = target
				return true
			}
		}
	}
	return false
}
func (action *APLActionMultidot) Execute(sim *Simulation) {
	action.spell.Cast(sim, action.nextTarget)
}
func (action *APLActionMultidot) String() string {
	return fmt.Sprintf("Multidot(%s)", action.spell.ActionID)
}

type APLActionMultishield struct {
	defaultAPLActionImpl
	spell      *Spell
	maxShields int32
	maxOverlap APLValue

	nextTarget *Unit
}

func (rot *APLRotation) newActionMultishield(config *proto.APLActionMultishield) APLActionImpl {
	unit := rot.unit

	spell := rot.aplGetMultishieldSpell(config.SpellId)
	if spell == nil {
		return nil
	}

	maxOverlap := rot.coerceTo(rot.newAPLValue(config.MaxOverlap), proto.APLValueType_ValueTypeDuration)
	if maxOverlap == nil {
		maxOverlap = rot.newValueConst(&proto.APLValueConst{Val: "0ms"})
	}

	maxShields := config.MaxShields
	numTargets := int32(len(unit.Env.Raid.AllPlayerUnits))
	if numTargets < maxShields {
		rot.validationWarning("Encounter only has %d targets. Using that for Max Shields instead of %d", numTargets, maxShields)
		maxShields = numTargets
	}

	return &APLActionMultishield{
		spell:      spell,
		maxShields: maxShields,
		maxOverlap: maxOverlap,
	}
}
func (action *APLActionMultishield) Reset(*Simulation) {
	action.nextTarget = nil
}
func (action *APLActionMultishield) IsReady(sim *Simulation) bool {
	maxOverlap := action.maxOverlap.GetDuration(sim)

	for i := int32(0); i < action.maxShields; i++ {
		target := sim.Raid.AllPlayerUnits[i]
		shield := action.spell.Shield(target)
		if (!shield.IsActive() || shield.RemainingDuration(sim) < maxOverlap) && action.spell.CanCast(sim, target) {
			action.nextTarget = target
			return true
		}
	}
	return false
}
func (action *APLActionMultishield) Execute(sim *Simulation) {
	action.spell.Cast(sim, action.nextTarget)
}
func (action *APLActionMultishield) String() string {
	return fmt.Sprintf("Multishield(%s)", action.spell.ActionID)
}

type APLActionAutocastOtherCooldowns struct {
	defaultAPLActionImpl
	character *Character

	nextReadyMCD *MajorCooldown
}

func (rot *APLRotation) newActionAutocastOtherCooldowns(config *proto.APLActionAutocastOtherCooldowns) APLActionImpl {
	unit := rot.unit
	return &APLActionAutocastOtherCooldowns{
		character: unit.Env.Raid.GetPlayerFromUnit(unit).GetCharacter(),
	}
}
func (action *APLActionAutocastOtherCooldowns) Reset(*Simulation) {
	action.nextReadyMCD = nil
}
func (action *APLActionAutocastOtherCooldowns) IsReady(sim *Simulation) bool {
	action.nextReadyMCD = action.character.getFirstReadyMCD(sim)

	// Explicitly check for GCD because MCDs are usually desired to be cast immediately
	// before the next spell, rather than immediately after the previous spell. This is
	// true even for MCDs which do not require the GCD.
	return action.nextReadyMCD != nil && action.character.GCD.IsReady(sim)
}
func (action *APLActionAutocastOtherCooldowns) Execute(sim *Simulation) {
	action.nextReadyMCD.tryActivateHelper(sim, action.character)
	action.character.UpdateMajorCooldowns()
}
func (action *APLActionAutocastOtherCooldowns) String() string {
	return fmt.Sprintf("Autocast Other Cooldowns")
}

type APLActionWait struct {
	defaultAPLActionImpl
	unit     *Unit
	duration APLValue
}

func (rot *APLRotation) newActionWait(config *proto.APLActionWait) APLActionImpl {
	unit := rot.unit
	return &APLActionWait{
		unit:     unit,
		duration: rot.coerceTo(rot.newAPLValue(config.Duration), proto.APLValueType_ValueTypeDuration),
	}
}
func (action *APLActionWait) IsReady(sim *Simulation) bool {
	return action.duration != nil
}
func (action *APLActionWait) Execute(sim *Simulation) {
	action.unit.WaitUntil(sim, sim.CurrentTime+action.duration.GetDuration(sim))
}
func (action *APLActionWait) String() string {
	return fmt.Sprintf("Wait(%s)", action.duration)
}
