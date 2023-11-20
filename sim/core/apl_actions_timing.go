package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/wowsims/classic/sim/core/proto"
)

type APLActionWait struct {
	defaultAPLActionImpl
	unit     *Unit
	duration APLValue

	curWaitTime time.Duration
}

func (rot *APLRotation) newActionWait(config *proto.APLActionWait) APLActionImpl {
	unit := rot.unit
	durationVal := rot.coerceTo(rot.newAPLValue(config.Duration), proto.APLValueType_ValueTypeDuration)
	if durationVal == nil {
		return nil
	}

	return &APLActionWait{
		unit:     unit,
		duration: durationVal,
	}
}
func (action *APLActionWait) GetAPLValues() []APLValue {
	return []APLValue{action.duration}
}
func (action *APLActionWait) IsReady(sim *Simulation) bool {
	return action.duration.GetDuration(sim) > 0
}

func (action *APLActionWait) Execute(sim *Simulation) {
	action.unit.Rotation.pushControllingAction(action)
	action.curWaitTime = sim.CurrentTime + action.duration.GetDuration(sim)

	pa := &PendingAction{
		Priority:     ActionPriorityLow,
		OnAction:     action.unit.gcdAction.OnAction,
		NextActionAt: action.curWaitTime,
	}
	sim.AddPendingAction(pa)
}

func (action *APLActionWait) GetNextAction(sim *Simulation) *APLAction {
	if sim.CurrentTime >= action.curWaitTime {
		action.unit.Rotation.popControllingAction(action)
		return action.unit.Rotation.getNextAction(sim)
	} else {
		return nil
	}
}

func (action *APLActionWait) String() string {
	return fmt.Sprintf("Wait(%s)", action.duration)
}

type APLActionWaitUntil struct {
	defaultAPLActionImpl
	unit      *Unit
	condition APLValue
}

func (rot *APLRotation) newActionWaitUntil(config *proto.APLActionWaitUntil) APLActionImpl {
	unit := rot.unit
	conditionVal := rot.coerceTo(rot.newAPLValue(config.Condition), proto.APLValueType_ValueTypeBool)
	if conditionVal == nil {
		return nil
	}

	return &APLActionWaitUntil{
		unit:      unit,
		condition: conditionVal,
	}
}
func (action *APLActionWaitUntil) GetAPLValues() []APLValue {
	return []APLValue{action.condition}
}
func (action *APLActionWaitUntil) IsReady(sim *Simulation) bool {
	return !action.condition.GetBool(sim)
}

func (action *APLActionWaitUntil) Execute(sim *Simulation) {
	action.unit.Rotation.pushControllingAction(action)
}

func (action *APLActionWaitUntil) GetNextAction(sim *Simulation) *APLAction {
	if action.condition.GetBool(sim) {
		action.unit.Rotation.popControllingAction(action)
		return action.unit.Rotation.getNextAction(sim)
	} else {
		return nil
	}
}

func (action *APLActionWaitUntil) String() string {
	return fmt.Sprintf("WaitUntil(%s)", action.condition)
}

type APLActionSchedule struct {
	defaultAPLActionImpl
	innerAction *APLAction

	timings       []time.Duration
	nextTimingIdx int
}

func (rot *APLRotation) newActionSchedule(config *proto.APLActionSchedule) APLActionImpl {
	innerAction := rot.newAPLAction(config.InnerAction)
	if innerAction == nil {
		return nil
	}

	timingStrs := strings.Split(config.Schedule, ",")
	if len(timingStrs) == 0 {
		return nil
	}

	timings := make([]time.Duration, len(timingStrs))
	valid := true
	for i, timingStr := range timingStrs {
		if durVal, err := time.ParseDuration(strings.TrimSpace(timingStr)); err == nil {
			timings[i] = durVal
		} else {
			rot.ValidationWarning("Invalid duration value '%s'", strings.TrimSpace(timingStr))
			valid = false
		}
	}
	if !valid {
		return nil
	}

	return &APLActionSchedule{
		innerAction: innerAction,
		timings:     timings,
	}
}
func (action *APLActionSchedule) Reset(*Simulation) {
	action.nextTimingIdx = 0
}
func (action *APLActionSchedule) GetInnerActions() []*APLAction {
	return []*APLAction{action.innerAction}
}
func (action *APLActionSchedule) IsReady(sim *Simulation) bool {
	return action.nextTimingIdx < len(action.timings) &&
		sim.CurrentTime >= action.timings[action.nextTimingIdx] &&
		action.innerAction.IsReady(sim)
}

func (action *APLActionSchedule) Execute(sim *Simulation) {
	action.nextTimingIdx++
	action.innerAction.Execute(sim)
}

func (action *APLActionSchedule) String() string {
	return fmt.Sprintf("Schedule(%s, %s)", action.timings, action.innerAction)
}
