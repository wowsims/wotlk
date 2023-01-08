package core

import (
	"time"
)

type DelayedActionOptions struct {
	// When the action should be performed.
	DoAt time.Duration

	Priority ActionPriority

	OnAction func(*Simulation)
	CleanUp  func(*Simulation)
}

func NewDelayedAction(sim *Simulation, options DelayedActionOptions) *PendingAction {
	pa := &PendingAction{
		NextActionAt: options.DoAt,
		Priority:     options.Priority,
	}

	pa.OnAction = func(sim *Simulation) {
		options.OnAction(sim)
	}
	pa.CleanUp = func(sim *Simulation) {
		if options.CleanUp != nil {
			options.CleanUp(sim)
		}
	}

	return pa
}

// Convenience for immediately creating and starting a delayed action.
func StartDelayedAction(sim *Simulation, options DelayedActionOptions) *PendingAction {
	pa := NewDelayedAction(sim, options)
	sim.AddPendingAction(pa)
	return pa
}

type PeriodicActionOptions struct {
	// How often the action should be performed.
	Period time.Duration

	// Number of times to perform the action before stopping.
	// 0 indicates a permanent periodic action.
	NumTicks int

	// Whether the first tick should happen immediately. If false, first tick will
	// wait for Period.
	TickImmediately bool

	Priority ActionPriority

	OnAction func(*Simulation)
	CleanUp  func(*Simulation)
}

func NewPeriodicAction(sim *Simulation, options PeriodicActionOptions) *PendingAction {
	pa := &PendingAction{
		NextActionAt: TernaryDuration(options.TickImmediately, sim.CurrentTime, sim.CurrentTime+options.Period),
		Priority:     options.Priority,
	}

	tickIndex := 0

	pa.OnAction = func(sim *Simulation) {
		options.OnAction(sim)
		tickIndex++

		if options.NumTicks == 0 || tickIndex < options.NumTicks {
			// Refresh action.
			pa.NextActionAt = sim.CurrentTime + options.Period
			sim.AddPendingAction(pa)
		} else {
			pa.Cancel(sim)
		}
	}
	pa.CleanUp = func(sim *Simulation) {
		if options.CleanUp != nil {
			options.CleanUp(sim)
		}
	}

	return pa
}

// Convenience for immediately creating and starting a periodic action.
func StartPeriodicAction(sim *Simulation, options PeriodicActionOptions) *PendingAction {
	pa := NewPeriodicAction(sim, options)
	sim.AddPendingAction(pa)
	return pa
}
