package core

import (
	"time"
)

// Stored value is the time at which the cooldown will be available again.
type Timer time.Duration

type Cooldown struct {
	*Timer

	// Default amount of time after activation before this CD can be used again.
	// Note that some CDs won't use this, e.g. the GCD.
	Duration time.Duration
}

func (unit *Unit) NewTimer() *Timer {
	if len(unit.cdTimers) > 100 {
		panic("Over 100 timers! There is probably one being registered every iteration.")
	}

	newTimer := new(Timer)
	unit.cdTimers = append(unit.cdTimers, newTimer)
	return newTimer
}

func (unit *Unit) resetCDs(sim *Simulation) {
	for _, timer := range unit.cdTimers {
		timer.Reset()
	}
}

func (timer *Timer) ReadyAt() time.Duration {
	return time.Duration(*timer)
}

func (timer *Timer) Set(t time.Duration) {
	*timer = Timer(t)
}

func (timer *Timer) Reset() {
	*timer = 0
}

func (timer *Timer) TimeToReady(sim *Simulation) time.Duration {
	return MaxDuration(0, time.Duration(*timer)-sim.CurrentTime)
}

func (timer *Timer) IsReady(sim *Simulation) bool {
	return time.Duration(*timer) <= sim.CurrentTime
}

// Puts this CD on cooldown, using the default duration.
func (cd *Cooldown) Use(sim *Simulation) {
	*cd.Timer = Timer(sim.CurrentTime + cd.Duration)
}

func BothTimersReadyAt(t1 *Timer, t2 *Timer) time.Duration {
	readyAt := time.Duration(0)
	if t1 != nil {
		readyAt = t1.ReadyAt()
	}
	if t2 != nil {
		readyAt = MaxDuration(readyAt, t2.ReadyAt())
	}
	return readyAt
}

func BothTimersReady(t1 *Timer, t2 *Timer, sim *Simulation) bool {
	return (t1 == nil || t1.IsReady(sim)) && (t2 == nil || t2.IsReady(sim))
}

func MaxTimeToReady(t1 *Timer, t2 *Timer, sim *Simulation) time.Duration {
	remaining := time.Duration(0)
	if t1 != nil {
		remaining = t1.TimeToReady(sim)
	}
	if t2 != nil {
		remaining = MaxDuration(remaining, t2.TimeToReady(sim))
	}
	return remaining
}

// Helper for shared timers that are not always needed, so it is only
// allocated if necessary.
func (unit *Unit) GetOrInitTimer(timer **Timer) *Timer {
	if *timer == nil {
		*timer = unit.NewTimer()
	}
	return *timer
}
