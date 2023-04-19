package core

import (
	"time"
)

type ActionPriority int32

const (
	ActionPriorityLow ActionPriority = -1
	ActionPriorityGCD ActionPriority = 0

	// Higher than GCD because regen can cause GCD actions (if we were waiting
	// for mana).
	ActionPriorityRegen ActionPriority = 1

	// Autos can cause regen (JoW, rage, energy procs, etc) so they should be
	// higher prio so that we never go backwards in the priority order.
	ActionPriorityAuto ActionPriority = 2

	// DOTs need to be higher than anything else so that dots can properly expire before we take other actions.
	ActionPriorityDOT ActionPriority = 3
)

type PendingAction struct {
	NextActionAt time.Duration
	Priority     ActionPriority

	OnAction func(*Simulation)
	CleanUp  func(*Simulation)

	cancelled bool
	consumed  bool

	repetitions    int
	maxRepetitions int
	period         time.Duration
}

func (pa *PendingAction) RunOnAction(sim *Simulation) {
	if pa.OnAction != nil {
		pa.OnAction(sim)
	}

	if pa.period > 0 {
		pa.repetitions--
		if pa.repetitions > 0 || pa.maxRepetitions == 0 {
			pa.NextActionAt = sim.CurrentTime + pa.period
			sim.AddPendingAction(pa)
		} else {
			pa.Cancel(sim)
		}
	}
}

func (pa *PendingAction) Cancel(sim *Simulation) {
	if pa.cancelled {
		return
	}

	if pa.CleanUp != nil {
		pa.CleanUp(sim)
		pa.CleanUp = nil
	}

	pa.cancelled = true
}
