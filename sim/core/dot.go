package core

import (
	"time"
)

type OnSnapshot func(sim *Simulation, target *Unit, dot *Dot, isRollover bool)
type OnTick func(sim *Simulation, target *Unit, dot *Dot)

type Dot struct {
	Spell *Spell

	// Embed Aura so we can use IsActive/Refresh/etc directly.
	*Aura

	NumberOfTicks int32         // number of ticks over the whole duration
	TickLength    time.Duration // time between each tick

	// If true, tick length will be shortened based on casting speed.
	AffectedByCastSpeed bool

	OnSnapshot OnSnapshot
	OnTick     OnTick

	SnapshotBaseDamage         float64
	SnapshotCritChance         float64
	SnapshotAttackerMultiplier float64

	tickAction *PendingAction
	tickPeriod time.Duration

	// Number of ticks since last call to Apply().
	TickCount int32

	lastTickTime time.Duration
}

// TickPeriod is how fast the snapshotted dot ticks.
func (dot *Dot) TickPeriod() time.Duration {
	return dot.tickPeriod
}

func (dot *Dot) TimeUntilNextTick(sim *Simulation) time.Duration {
	return dot.lastTickTime + dot.tickPeriod - sim.CurrentTime
}

// Roll over = gets carried over with everlasting refresh and doesn't get applied if triggered when the spell is already up.
// - Example: critical strike rating, internal % damage modifiers: buffs or debuffs on player
// Nevermelting Ice, Shadow Mastery (ISB), Trick of the Trades, Deaths Embrace, Thadius Polarity, Hera Spores, Crit on weapons from swapping

// Snapshot = calculation happens at refresh and application (stays up even if buff falls of, until new refresh or application)
// - Example: Spell power, Haste rating
// Blood Fury, Lightweave Embroid, Eradication, Bloodlust

// Dynamic = realtime update
// - Example: external % damage modifier debuffs on target
// Haunt, Curse of Shadow, Shadow Embrace

// Rollover is used to reset the duration of a dot from an external spell (not casting the dot itself)
// This keeps the snapshotted crit and %dmg modifiers.
// However sp and haste are recalculated.
func (dot *Dot) Rollover(sim *Simulation) {
	dot.TakeSnapshot(sim, true)

	dot.RecomputeAuraDuration() // recalculate haste
	dot.Aura.Refresh(sim)       // update aura's duration

	oldNextTick := dot.tickAction.NextActionAt
	dot.tickAction.Cancel(sim) // remove old PA ticker

	// recreate with new period, resetting the next tick.
	periodicOptions := dot.basePeriodicOptions()
	periodicOptions.Period = dot.tickPeriod
	dot.tickAction = NewPeriodicAction(sim, periodicOptions)
	dot.tickAction.NextActionAt = oldNextTick
	sim.AddPendingAction(dot.tickAction)
}

func (dot *Dot) Apply(sim *Simulation) {
	dot.Cancel(sim)
	dot.TickCount = 0
	dot.RecomputeAuraDuration()
	dot.Aura.Activate(sim)
}

// ApplyOrReset is used for rolling dots that reset the tick timer on reapplication.
// This is more efficient than Apply(), and works around tickAction.CleanUp() wrongly generating
// an extra ticks if (re-)application and tick happen at the same time.
func (dot *Dot) ApplyOrReset(sim *Simulation) {
	if !dot.IsActive() {
		dot.Apply(sim)
		return
	}

	dot.TakeSnapshot(sim, true)

	dot.RecomputeAuraDuration() // recalculate haste
	dot.Aura.Refresh(sim)       // update aura's duration

	dot.TickCount = 0

	dot.tickAction.NextActionAt = -1 // prevent tickAction.CleanUp() from adding an extra tick
	dot.tickAction.Cancel(sim)       // remove old PA ticker

	// recreate with new period, resetting the next tick.
	periodicOptions := dot.basePeriodicOptions()
	periodicOptions.Period = dot.tickPeriod
	dot.tickAction = NewPeriodicAction(sim, periodicOptions)
	sim.AddPendingAction(dot.tickAction)
}

// Like Apply(), but does not reset the tick timer.
func (dot *Dot) ApplyOrRefresh(sim *Simulation) {
	dot.TickCount = 0
	dot.RecomputeAuraDuration()
	dot.Aura.Activate(sim)
}

func (dot *Dot) Cancel(sim *Simulation) {
	if dot.Aura.IsActive() {
		dot.Aura.Deactivate(sim)
	}
}

// Call this after manually changing NumberOfTicks or TickLength.
func (dot *Dot) RecomputeAuraDuration() {
	if dot.AffectedByCastSpeed {
		dot.tickPeriod = dot.Spell.Unit.ApplyCastSpeed(dot.TickLength)
		dot.Aura.Duration = dot.tickPeriod * time.Duration(dot.NumberOfTicks)
	} else {
		dot.tickPeriod = dot.TickLength
		dot.Aura.Duration = dot.tickPeriod * time.Duration(dot.NumberOfTicks)
	}
}

// Takes a new snapshot of this Dot's effects.
//
// In most cases this will be called automatically, and should only be called
// to force a new snapshot to be taken.
//
//	doRollover will apply previously snapshotted crit/%dmg instead of recalculating.
func (dot *Dot) TakeSnapshot(sim *Simulation, doRollover bool) {
	if dot.OnSnapshot != nil {
		dot.OnSnapshot(sim, dot.Unit, dot, doRollover)
	}
}

// Forces an instant tick. Does not reset the tick timer or aura duration,
// the tick is simply an extra tick.
func (dot *Dot) TickOnce(sim *Simulation) {
	dot.lastTickTime = sim.CurrentTime
	dot.OnTick(sim, dot.Unit, dot)
}

func (dot *Dot) basePeriodicOptions() PeriodicActionOptions {
	return PeriodicActionOptions{
		//Priority: ActionPriorityDOT,
		OnAction: func(sim *Simulation) {
			if dot.lastTickTime != sim.CurrentTime {
				dot.TickCount++
				dot.TickOnce(sim)
			}
		},
		CleanUp: func(sim *Simulation) {
			// In certain cases, the last tick and the dot aura expiration can happen in
			// different orders, so we might need to apply the last tick.
			if dot.tickAction != nil && dot.tickAction.NextActionAt == sim.CurrentTime {
				if dot.lastTickTime != sim.CurrentTime {
					dot.TickCount++
					dot.TickOnce(sim)
				}
			}
		},
	}
}

func NewDot(config Dot) *Dot {
	dot := &Dot{}
	*dot = config

	dot.tickPeriod = dot.TickLength
	dot.Aura.Duration = dot.TickLength * time.Duration(dot.NumberOfTicks)

	dot.Aura.ApplyOnGain(func(aura *Aura, sim *Simulation) {
		dot.lastTickTime = sim.CurrentTime
		dot.TakeSnapshot(sim, false)

		periodicOptions := dot.basePeriodicOptions()
		periodicOptions.Period = dot.tickPeriod
		dot.tickAction = NewPeriodicAction(sim, periodicOptions)
		sim.AddPendingAction(dot.tickAction)
	})
	dot.Aura.ApplyOnExpire(func(aura *Aura, sim *Simulation) {
		if dot.tickAction != nil {
			dot.tickAction.Cancel(sim)
			dot.tickAction = nil
		}
	})

	return dot
}

// Creates HoTs for all allied units.
func NewAllyHotArray(caster *Unit, config Dot, auraConfig Aura) []*Dot {
	hots := make([]*Dot, len(caster.Env.AllUnits))
	for _, target := range caster.Env.AllUnits {
		if !caster.IsOpponent(target) {
			config.Aura = target.GetOrRegisterAura(auraConfig)
			hots[target.UnitIndex] = NewDot(config)
		}
	}
	return hots
}
