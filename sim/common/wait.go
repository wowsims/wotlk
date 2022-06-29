package common

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

type WaitReason byte

const (
	WaitReasonNone     WaitReason = iota // unknown why we waited
	WaitReasonOOM                        // no mana to cast
	WaitReasonRotation                   // waiting on rotation
	WaitReasonOptimal                    // waiting because its more optimal than casting.
)

type WaitAction struct {
	unit *core.Unit

	duration time.Duration
	reason   WaitReason
}

func (action WaitAction) GetActionID() core.ActionID {
	return core.ActionID{
		OtherID: proto.OtherAction_OtherActionWait,
	}
}

func (action WaitAction) GetName() string {
	return "Wait"
}

func (action WaitAction) GetTag() int32 {
	return 0
}

func (action WaitAction) GetUnit() *core.Unit {
	return action.unit
}

func (action WaitAction) GetDuration() time.Duration {
	return action.duration
}

func (action WaitAction) GetManaCost() float64 {
	return 0
}

func (action WaitAction) Cast(sim *core.Simulation) bool {
	switch action.reason {
	case WaitReasonNone:
		if sim.Log != nil {
			action.unit.Log(sim, "Idling for %s seconds, for no particular reason.", action.GetDuration())
		}
	case WaitReasonOOM:
		action.unit.Metrics.MarkOOM(action.unit, action.GetDuration())
		if sim.Log != nil {
			action.unit.Log(sim, "Not enough mana to cast, regenerating for %s.", action.GetDuration())
		}
	case WaitReasonRotation:
		if sim.Log != nil {
			action.unit.Log(sim, "Waiting for %s due to rotation / CDs.", action.GetDuration())
		}
	case WaitReasonOptimal:
		if sim.Log != nil {
			action.unit.Log(sim, "Waiting for %s because its more dps.", action.GetDuration())
		}
	}
	action.unit.SetGCDTimer(sim, sim.CurrentTime+action.GetDuration())

	return true
}

func NewWaitAction(sim *core.Simulation, unit *core.Unit, duration time.Duration, reason WaitReason) WaitAction {
	return WaitAction{
		unit:     unit,
		duration: duration,
		reason:   reason,
	}
}
