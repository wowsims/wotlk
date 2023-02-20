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

	return &APLRotation{
		unit:         unit,
		priorityList: priorityList,
	}
}

// We intentionally try to mimic the behavior of simc APL to avoid confusion
// and leverage the community's existing familiarity.
// https://github.com/simulationcraft/simc/wiki/ActionLists
func (apl *APLRotation) DoNextAction(sim *Simulation) {
	for _, action := range apl.priorityList {
		if action.IsAvailable(sim) {
			action.Execute(sim)
			return
		}
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
