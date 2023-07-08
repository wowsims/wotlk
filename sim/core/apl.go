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

	for i, action := range rotation.priorityList {
		action.impl.Finalize(rotation)

		rotation.priorityListWarnings[configIdxs[i]] = append(rotation.priorityListWarnings[configIdxs[i]], rotation.curWarnings...)
		rotation.curWarnings = nil

		// Remove MCDs that are referenced by APL actions.
		character := unit.Env.Raid.GetPlayerFromUnit(unit).GetCharacter()
		if castSpellAction, ok := action.impl.(*APLActionCastSpell); ok {
			character.removeInitialMajorCooldown(castSpellAction.spell.ActionID)
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

func APLRotationFromJsonString(jsonString string) *proto.APLRotation {
	apl := &proto.APLRotation{}
	data := []byte(jsonString)
	if err := protojson.Unmarshal(data, apl); err != nil {
		panic(err)
	}
	return apl
}
