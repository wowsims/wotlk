package core

import (
	"fmt"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type APLRotation struct {
	unit           *Unit
	prepullActions []*APLAction
	priorityList   []*APLAction

	// Current strict sequence
	strictSequence *APLAction

	// Used inside of actions/value to determine whether they will occur during the prepull or regular rotation.
	parsingPrepull bool

	// Used to avoid recursive APL loops.
	inLoop bool

	// Validation warnings that occur during proto parsing.
	// We return these back to the user for display in the UI.
	curWarnings          []string
	prepullWarnings      [][]string
	priorityListWarnings [][]string
}

func (rot *APLRotation) validationWarning(message string, vals ...interface{}) {
	warning := fmt.Sprintf(message, vals...)
	rot.curWarnings = append(rot.curWarnings, warning)
}

func (unit *Unit) newAPLRotation(config *proto.APLRotation) *APLRotation {
	if config == nil || !config.Enabled {
		return nil
	}

	rotation := &APLRotation{
		unit: unit,
	}

	// Parse prepull actions
	rotation.parsingPrepull = true
	for _, prepullItem := range config.PrepullActions {
		if !prepullItem.Hide {
			doAtVal := rotation.newAPLValue(prepullItem.DoAtValue)
			if doAtVal != nil {
				doAt := doAtVal.GetDuration(nil)
				if doAt > 0 {
					rotation.validationWarning("Invalid time for 'Do At', ignoring this Prepull Action")
				} else {
					action := rotation.newAPLAction(prepullItem.Action)
					if action != nil {
						rotation.prepullActions = append(rotation.prepullActions, action)
						unit.RegisterPrepullAction(doAt, func(sim *Simulation) {
							action.Execute(sim)
						})
					}
				}
			}
		}

		rotation.prepullWarnings = append(rotation.prepullWarnings, rotation.curWarnings)
		rotation.curWarnings = nil
	}
	rotation.parsingPrepull = false

	// Parse priority list
	var configIdxs []int
	for i, aplItem := range config.PriorityList {
		if !aplItem.Hide {
			action := rotation.newAPLAction(aplItem.Action)
			if action != nil {
				rotation.priorityList = append(rotation.priorityList, action)
				configIdxs = append(configIdxs, i)
			}
		}

		rotation.priorityListWarnings = append(rotation.priorityListWarnings, rotation.curWarnings)
		rotation.curWarnings = nil
	}

	// Finalize
	for _, action := range rotation.prepullActions {
		action.impl.Finalize(rotation)
		rotation.curWarnings = nil
	}
	for i, action := range rotation.priorityList {
		action.impl.Finalize(rotation)

		rotation.priorityListWarnings[configIdxs[i]] = append(rotation.priorityListWarnings[configIdxs[i]], rotation.curWarnings...)
		rotation.curWarnings = nil
	}

	// Remove MCDs that are referenced by APL actions, so that the Autocast Other Cooldowns
	// action does not include them.
	character := unit.Env.Raid.GetPlayerFromUnit(unit).GetCharacter()
	for _, action := range rotation.allAPLActions() {
		if castSpellAction, ok := action.impl.(*APLActionCastSpell); ok {
			character.removeInitialMajorCooldown(castSpellAction.spell.ActionID)
		}
	}

	// If user has a Prepull potion set but does not use it in their APL settings, we enable it here.
	rotation.parsingPrepull = true
	prepotSpell := rotation.aplGetSpell(ActionID{OtherID: proto.OtherAction_OtherActionPotion}.ToProto())
	rotation.parsingPrepull = false
	if prepotSpell != nil {
		found := false
		for _, prepullAction := range rotation.allPrepullActions() {
			if castSpellAction, ok := prepullAction.impl.(*APLActionCastSpell); ok && castSpellAction.spell == prepotSpell {
				found = true
			}
		}
		if !found {
			unit.RegisterPrepullAction(-1*time.Second, func(sim *Simulation) {
				prepotSpell.Cast(sim, nil)
			})
		}
	}

	return rotation
}
func (rot *APLRotation) getStats() *proto.APLStats {
	return &proto.APLStats{
		PrepullActions: MapSlice(rot.prepullWarnings, func(warnings []string) *proto.APLActionStats { return &proto.APLActionStats{Warnings: warnings} }),
		PriorityList:   MapSlice(rot.priorityListWarnings, func(warnings []string) *proto.APLActionStats { return &proto.APLActionStats{Warnings: warnings} }),
	}
}

// Returns all action objects as an unstructured list. Used for easily finding specific actions.
func (rot *APLRotation) allAPLActions() []*APLAction {
	return Flatten(MapSlice(rot.priorityList, func(action *APLAction) []*APLAction { return action.GetAllActions() }))
}

// Returns all action objects from the prepull as an unstructured list. Used for easily finding specific actions.
func (rot *APLRotation) allPrepullActions() []*APLAction {
	return Flatten(MapSlice(rot.prepullActions, func(action *APLAction) []*APLAction { return action.GetAllActions() }))
}

func (rot *APLRotation) reset(sim *Simulation) {
	rot.strictSequence = nil
	rot.inLoop = false
	for _, action := range rot.allAPLActions() {
		action.impl.Reset(sim)
	}
}

// We intentionally try to mimic the behavior of simc APL to avoid confusion
// and leverage the community's existing familiarity.
// https://github.com/simulationcraft/simc/wiki/ActionLists
func (apl *APLRotation) DoNextAction(sim *Simulation) {
	if apl.inLoop {
		return
	}

	i := 0
	apl.inLoop = true
	for nextAction := apl.getNextAction(sim); nextAction != nil; i, nextAction = i+1, apl.getNextAction(sim) {
		if i > 1000 {
			panic(fmt.Sprintf("[USER_ERROR] Infinite loop detected, current action:\n%s", nextAction))
		}

		nextAction.Execute(sim)
	}
	apl.inLoop = false

	if sim.Log != nil && i == 0 {
		apl.unit.Log(sim, "No available actions!")
	}

	if apl.unit.GCD.IsReady(sim) {
		apl.unit.WaitUntil(sim, sim.CurrentTime+time.Millisecond*500)
	} else {
		apl.unit.DoNothing()
	}
}

func (apl *APLRotation) getNextAction(sim *Simulation) *APLAction {
	if apl.strictSequence != nil {
		ss := apl.strictSequence.impl.(*APLActionStrictSequence)
		if ss.actions[ss.curIdx].IsReady(sim) {
			return apl.strictSequence
		} else if apl.unit.GCD.IsReady(sim) {
			// If the GCD is ready when the next subaction isn't, it means the sequence is bad
			// so reset and exit the sequence.
			ss.curIdx = 0
			apl.strictSequence = nil
		} else {
			// Return nil to wait for the GCD to become ready.
			return nil
		}
	}

	for _, action := range apl.priorityList {
		if action.IsReady(sim) {
			if _, ok := action.impl.(*APLActionStrictSequence); ok {
				apl.strictSequence = action
			}
			return action
		}
	}

	return nil
}

func APLRotationFromJsonString(jsonString string) *proto.APLRotation {
	apl := &proto.APLRotation{}
	data := []byte(jsonString)
	if err := protojson.Unmarshal(data, apl); err != nil {
		panic(err)
	}
	return apl
}
