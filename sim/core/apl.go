package core

import (
	"fmt"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type APLRotation struct {
	unit         *Unit
	priorityList []*APLAction

	// Current strict sequence
	strictSequence *APLActionStrictSequence
}

func (unit *Unit) newAPLRotation(config *proto.APLRotation) *APLRotation {
	if config == nil || !config.Enabled {
		return nil
	}

	priorityList := MapSlice(config.PriorityList, func(aplItem *proto.APLListItem) *APLAction {
		if aplItem.Hide {
			return nil
		} else {
			return unit.newAPLAction(aplItem.Action)
		}
	})
	priorityList = FilterSlice(priorityList, func(action *APLAction) bool { return action != nil })

	rotation := &APLRotation{
		unit:         unit,
		priorityList: priorityList,
	}

	for _, action := range rotation.allAPLActions() {
		action.impl.Finalize()

		// Remove MCDs that are referenced by APL actions.
		character := unit.Env.Raid.GetPlayerFromUnit(unit).GetCharacter()
		if castSpellAction, ok := action.impl.(*APLActionCastSpell); ok {
			character.removeInitialMajorCooldown(castSpellAction.spell.ActionID)
		}
	}

	return rotation
}

// Returns all action objects as an unstructured list. Used for easily finding specific actions.
func (rot *APLRotation) allAPLActions() []*APLAction {
	return Flatten(MapSlice(rot.priorityList, func(action *APLAction) []*APLAction { return action.GetAllActions() }))
}
func (unit *Unit) allAPLActions() []*APLAction {
	return unit.Rotation.allAPLActions()
}

func (rot *APLRotation) reset(sim *Simulation) {
	rot.strictSequence = nil
	for _, action := range rot.allAPLActions() {
		action.impl.Reset(sim)
	}
}

// We intentionally try to mimic the behavior of simc APL to avoid confusion
// and leverage the community's existing familiarity.
// https://github.com/simulationcraft/simc/wiki/ActionLists
func (apl *APLRotation) DoNextAction(sim *Simulation) {
	if apl.strictSequence == nil {
		for _, action := range apl.priorityList {
			if action.IsReady(sim) {
				action.Execute(sim)
				if apl.unit.GCD.IsReady(sim) {
					apl.unit.WaitUntil(sim, sim.CurrentTime)
				}
				return
			}
		}
	} else {
		apl.strictSequence.Execute(sim)
	}

	if sim.Log != nil {
		apl.unit.Log(sim, "No available actions!")
	}
	if apl.unit.GCD.IsReady(sim) {
		apl.unit.WaitUntil(sim, sim.CurrentTime+time.Millisecond*500)
	} else {
		apl.unit.DoNothing()
	}
}

func validationError(message string, vals ...interface{}) {
	panic("Validation Error: " + fmt.Sprintf(message, vals...))
}

// For validation issues that we can manage internally. Will probably make this a test-only panic later.
func validationWarning(message string, vals ...interface{}) {
	panic("Validation Warning: " + fmt.Sprintf(message, vals...))
}

func APLRotationFromJsonString(jsonString string) *proto.APLRotation {
	apl := &proto.APLRotation{}
	data := []byte(jsonString)
	if err := protojson.Unmarshal(data, apl); err != nil {
		panic(err)
	}
	return apl
}
