package core

import (
	"fmt"
	"time"
)

// Note that this is only used when the hardcast and GCD actions
func (unit *Unit) newHardcastAction(sim *Simulation) {
	if unit.hardcastAction != nil && !unit.hardcastAction.consumed {
		unit.hardcastAction.Cancel(sim)
		unit.hardcastAction = nil
	}

	if unit.hardcastAction == nil {
		pa := &PendingAction{
			NextActionAt: unit.Hardcast.Expires,
			OnAction: func(sim *Simulation) {
				// Don't need to do anything, the Advance() call will take care of the hardcast.
			},
		}
		unit.hardcastAction = pa
	} else {
		unit.hardcastAction.cancelled = false
		unit.hardcastAction.NextActionAt = unit.Hardcast.Expires
	}

	sim.AddPendingAction(unit.hardcastAction)
}

func (unit *Unit) NextGCDAt() time.Duration {
	return unit.gcdAction.NextActionAt
}

func (unit *Unit) SetGCDTimer(sim *Simulation, gcdReadyAt time.Duration) {
	unit.GCD.Set(gcdReadyAt)

	if unit.gcdAction.consumed {
		unit.gcdAction.cancelled = false
		unit.gcdAction.NextActionAt = gcdReadyAt
	} else {
		unit.gcdAction.Cancel(sim)
		oldAction := unit.gcdAction.OnAction
		unit.gcdAction = &PendingAction{
			NextActionAt: gcdReadyAt,
			Priority:     ActionPriorityGCD,
			OnAction:     oldAction,
		}
	}
	sim.AddPendingAction(unit.gcdAction)
}

// Call this to stop the GCD loop for a unit.
// This is mostly used for pets that get summoned / expire.
func (unit *Unit) CancelGCDTimer(sim *Simulation) {
	unit.gcdAction.Cancel(sim)
}

func (unit *Unit) IsWaiting() bool {
	return unit.waitStartTime != 0
}
func (unit *Unit) IsWaitingForMana() bool {
	return unit.waitingForMana != 0
}

func (unit *Unit) IsWaitingForEnergy() bool {
	return unit.waitingForEnergy != 0
}

// Assumes that IsWaitingForMana() == true
func (unit *Unit) DoneWaitingForMana(sim *Simulation) bool {
	if unit.CurrentMana() >= unit.waitingForMana {
		unit.Metrics.MarkOOM(unit, sim.CurrentTime-unit.waitStartTime)
		unit.waitStartTime = 0
		unit.waitingForMana = 0
		return true
	}
	return false
}

func (unit *Unit) DoneWaitingForEnergy(sim *Simulation) bool {
	if unit.CurrentEnergy() >= unit.waitingForEnergy {
		unit.waitStartTime = 0
		unit.waitingForEnergy = 0
		return true
	}
	return false
}

// Returns true if the unit was waiting for mana but is now finished AND
// the GCD is also ready.
func (unit *Unit) FinishedWaitingForManaAndGCDReady(sim *Simulation) bool {
	if !unit.IsWaitingForMana() || !unit.DoneWaitingForMana(sim) {
		return false
	}

	return unit.GCD.IsReady(sim)
}

func (unit *Unit) WaitUntil(sim *Simulation, readyTime time.Duration) {
	if readyTime < sim.CurrentTime {
		panic(unit.Label + ": cannot wait negative time")
	}
	unit.waitStartTime = sim.CurrentTime
	unit.SetGCDTimer(sim, readyTime)
	if sim.Log != nil {
		unit.Log(sim, "Pausing GCD for %s due to rotation / CDs.", readyTime-sim.CurrentTime)
	}
}

func (unit *Unit) HardcastWaitUntil(sim *Simulation, readyTime time.Duration, onComplete CastFunc) {
	if unit.Hardcast.Expires >= sim.CurrentTime {
		fmt.Printf("Sim current time: %0.2f\n", sim.CurrentTime.Seconds())
		panic(fmt.Sprintf("Hardcast already in use, will finish at: %0.2f", unit.Hardcast.Expires.Seconds()))
	}

	unit.Hardcast.Expires = readyTime
	unit.Hardcast.OnComplete = onComplete
	unit.newHardcastAction(sim)
}

func (unit *Unit) WaitForMana(sim *Simulation, desiredMana float64) {
	if !unit.IsWaitingForMana() {
		unit.waitStartTime = sim.CurrentTime
	}
	unit.waitingForMana = desiredMana
	if sim.Log != nil {
		unit.Log(sim, "Not enough mana to cast, pausing GCD until mana >= %0.01f.", desiredMana)
	}
}

func (unit *Unit) WaitForEnergy(sim *Simulation, desiredEnergy float64) {
	if !unit.IsWaitingForEnergy() {
		unit.waitStartTime = sim.CurrentTime
	}
	unit.waitingForEnergy = desiredEnergy
	if sim.Log != nil {
		unit.Log(sim, "Not enough energy to cast, pausing GCD until energy >= %0.01f.", desiredEnergy)
	}
}

func (unit *Unit) doneIterationGCD(simDuration time.Duration) {
	if unit.IsWaitingForMana() {
		unit.Metrics.MarkOOM(unit, simDuration-unit.waitStartTime)
		unit.waitStartTime = 0
		unit.waitingForMana = 0
	} else if unit.IsWaitingForEnergy() {
		unit.waitStartTime = 0
		unit.waitingForEnergy = 0
	} else if unit.IsWaiting() {
		unit.waitStartTime = 0
	}
}
