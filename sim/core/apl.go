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

	// Action currently controlling this rotation (only used for certain actions, such as StrictSequence).
	controllingActions []APLActionImpl

	// Value that should evaluate to 'true' if the current channel is to be interrupted.
	// Will be nil when there is no active channel.
	interruptChannelIf APLValue

	// If true, can recast channel when interrupted.
	allowChannelRecastOnInterrupt bool

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

func (rot *APLRotation) ValidationWarning(message string, vals ...interface{}) {
	warning := fmt.Sprintf(message, vals...)
	rot.curWarnings = append(rot.curWarnings, warning)
}

// Invokes the fn function, and attributes all warnings generated during its invocation
// to the provided warningsList.
func (rot *APLRotation) doAndRecordWarnings(warningsList *[]string, isPrepull bool, fn func()) {
	rot.parsingPrepull = isPrepull
	fn()
	if warningsList != nil {
		*warningsList = append(*warningsList, rot.curWarnings...)
	}
	rot.curWarnings = nil
	rot.parsingPrepull = false
}

func (unit *Unit) newCustomRotation() *APLRotation {
	return unit.newAPLRotation(&proto.APLRotation{
		Type: proto.APLRotation_TypeAPL,
		PriorityList: []*proto.APLListItem{
			{
				Action: &proto.APLAction{
					Action: &proto.APLAction_CustomRotation{},
				},
			},
		},
	})
}

func (unit *Unit) newAPLRotation(config *proto.APLRotation) *APLRotation {
	if config == nil {
		return nil
	}

	rotation := &APLRotation{
		unit:                 unit,
		prepullWarnings:      make([][]string, len(config.PrepullActions)),
		priorityListWarnings: make([][]string, len(config.PriorityList)),
	}

	// Parse prepull actions
	for i, prepullItem := range config.PrepullActions {
		prepullIdx := i // Save to local variable for correct lambda capture behavior
		rotation.doAndRecordWarnings(&rotation.prepullWarnings[prepullIdx], true, func() {
			if !prepullItem.Hide {
				doAtVal := rotation.newAPLValue(prepullItem.DoAtValue)
				if doAtVal != nil {
					doAt := doAtVal.GetDuration(nil)
					if doAt > 0 {
						rotation.ValidationWarning("Invalid time for 'Do At', ignoring this Prepull Action")
					} else {
						action := rotation.newAPLAction(prepullItem.Action)
						if action != nil {
							rotation.prepullActions = append(rotation.prepullActions, action)
							unit.RegisterPrepullAction(doAt, func(sim *Simulation) {
								// Warnings for prepull cast failure are detected by running a fake prepull,
								// so this action.Execute needs to record warnings.
								rotation.doAndRecordWarnings(&rotation.prepullWarnings[prepullIdx], true, func() {
									action.Execute(sim)
								})
							})
						}
					}
				}
			}
		})
	}

	// Parse priority list
	var configIdxs []int
	for i, aplItem := range config.PriorityList {
		rotation.doAndRecordWarnings(&rotation.priorityListWarnings[i], false, func() {
			if !aplItem.Hide {
				action := rotation.newAPLAction(aplItem.Action)
				if action != nil {
					rotation.priorityList = append(rotation.priorityList, action)
					configIdxs = append(configIdxs, i)
				}
			}
		})
	}

	// Finalize
	for i, action := range rotation.prepullActions {
		rotation.doAndRecordWarnings(&rotation.prepullWarnings[i], true, func() {
			action.Finalize(rotation)
		})
	}
	for i, action := range rotation.priorityList {
		rotation.doAndRecordWarnings(&rotation.priorityListWarnings[i], false, func() {
			action.Finalize(rotation)
		})
	}

	// Remove MCDs that are referenced by APL actions, so that the Autocast Other Cooldowns
	// action does not include them.
	agent := unit.Env.GetAgentFromUnit(unit)
	if agent != nil {
		character := agent.GetCharacter()
		for _, action := range rotation.allAPLActions() {
			if castSpellAction, ok := action.impl.(*APLActionCastSpell); ok {
				character.removeInitialMajorCooldown(castSpellAction.spell.ActionID)
			}
		}
	}

	// If user has a Prepull potion set but does not use it in their APL settings, we enable it here.
	rotation.doAndRecordWarnings(nil, true, func() {
		prepotSpell := rotation.GetAPLSpell(ActionID{OtherID: proto.OtherAction_OtherActionPotion}.ToProto())
		if prepotSpell != nil {
			found := false
			for _, prepullAction := range rotation.allPrepullActions() {
				if castSpellAction, ok := prepullAction.impl.(*APLActionCastSpell); ok &&
					(castSpellAction.spell == prepotSpell || castSpellAction.spell.Flags.Matches(SpellFlagPotion)) {
					found = true
				}
			}
			if !found {
				unit.RegisterPrepullAction(-1*time.Second, func(sim *Simulation) {
					prepotSpell.Cast(sim, nil)
				})
			}
		}
	})

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
	if rot == nil || rot.priorityList == nil {
		return []*APLAction{}
	}

	return Flatten(MapSlice(rot.priorityList, func(action *APLAction) []*APLAction {
		// Check if action is nil before calling GetAllActions
		if action == nil {
			return []*APLAction{}
		}
		return action.GetAllActions()
	}))
}

// Returns all action objects from the prepull as an unstructured list. Used for easily finding specific actions.
func (rot *APLRotation) allPrepullActions() []*APLAction {
	return Flatten(MapSlice(rot.prepullActions, func(action *APLAction) []*APLAction { return action.GetAllActions() }))
}

func (rot *APLRotation) reset(sim *Simulation) {
	rot.controllingActions = nil
	rot.inLoop = false
	rot.interruptChannelIf = nil
	rot.allowChannelRecastOnInterrupt = false
	for _, action := range rot.allAPLActions() {
		action.impl.Reset(sim)
	}
}

// We intentionally try to mimic the behavior of simc APL to avoid confusion
// and leverage the community's existing familiarity.
// https://github.com/simulationcraft/simc/wiki/ActionLists
func (apl *APLRotation) DoNextAction(sim *Simulation) {
	if sim.CurrentTime < 0 {
		return
	}

	if apl.inLoop {
		return
	}

	if apl.unit.ChanneledDot != nil {
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

	gcdReady := apl.unit.GCD.IsReady(sim)
	if gcdReady {
		apl.unit.WaitUntil(sim, sim.CurrentTime+time.Millisecond*50)
	}
}

func (apl *APLRotation) getNextAction(sim *Simulation) *APLAction {
	if len(apl.controllingActions) != 0 {
		return apl.controllingActions[len(apl.controllingActions)-1].GetNextAction(sim)
	}

	for _, action := range apl.priorityList {
		if action.IsReady(sim) {
			return action
		}
	}

	return nil
}

func (apl *APLRotation) pushControllingAction(ca APLActionImpl) {
	apl.controllingActions = append(apl.controllingActions, ca)
}

func (apl *APLRotation) popControllingAction(ca APLActionImpl) {
	if len(apl.controllingActions) == 0 || apl.controllingActions[len(apl.controllingActions)-1] != ca {
		panic("Wrong APL controllingAction in pop()")
	}
	apl.controllingActions = apl.controllingActions[:len(apl.controllingActions)-1]
}

func (apl *APLRotation) shouldInterruptChannel(sim *Simulation) bool {
	channeledDot := apl.unit.ChanneledDot

	if channeledDot.MaxTicksRemaining() == 0 {
		// Channel has ended, but apl.unit.ChanneledDot hasn't been cleared yet meaning the aura is still active.
		return false
	}

	if apl.interruptChannelIf == nil || !apl.interruptChannelIf.GetBool(sim) {
		// Continue the channel.
		return false
	}

	// Allow next action to interrupt the channel, but if the action is the same action then it still needs to continue.
	nextAction := apl.getNextAction(sim)
	if nextAction == nil {
		return false
	}

	if channelAction, ok := nextAction.impl.(*APLActionChannelSpell); ok && channelAction.spell == channeledDot.Spell {
		// Newly selected action is channeling the same spell, so continue the channel unless recast is allowed.
		return apl.allowChannelRecastOnInterrupt
	}

	return true
}

func APLRotationFromJsonString(jsonString string) *proto.APLRotation {
	apl := &proto.APLRotation{}
	data := []byte(jsonString)
	if err := protojson.Unmarshal(data, apl); err != nil {
		panic(err)
	}
	return apl
}
