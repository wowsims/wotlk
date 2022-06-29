package core

import (
	"time"
)

type TickEffects func(*Simulation, *Spell) func()

type Dot struct {
	Spell *Spell

	// Embed Aura so we can use IsActive/Refresh/etc directly.
	*Aura

	NumberOfTicks int           // number of ticks over the whole duration
	TickLength    time.Duration // time between each tick

	// If true, tick length will be shortened based on casting speed.
	AffectedByCastSpeed bool

	TickEffects TickEffects

	tickFn     func()
	tickAction *PendingAction
	tickPeriod time.Duration

	// Number of ticks since last call to Apply().
	TickCount int

	lastTickTime time.Duration
}

func (dot *Dot) Apply(sim *Simulation) {
	dot.Cancel(sim)

	dot.TickCount = 0
	if dot.AffectedByCastSpeed {
		dot.tickPeriod = dot.Spell.Unit.ApplyCastSpeed(dot.TickLength)
		dot.Aura.Duration = dot.tickPeriod * time.Duration(dot.NumberOfTicks)
	}
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
func (dot *Dot) TakeSnapshot(sim *Simulation) {
	dot.tickFn = dot.TickEffects(sim, dot.Spell)
}

func NewDot(config Dot) *Dot {
	dot := &Dot{}
	*dot = config

	basePeriodicOptions := PeriodicActionOptions{
		OnAction: func(sim *Simulation) {
			if dot.lastTickTime != sim.CurrentTime {
				dot.lastTickTime = sim.CurrentTime
				dot.TickCount++
				dot.tickFn()
			}
		},
		CleanUp: func(sim *Simulation) {
			// In certain cases, the last tick and the dot aura expiration can happen in
			// different orders, so we might need to apply the last tick.
			if dot.tickAction.NextActionAt == sim.CurrentTime {
				if dot.lastTickTime != sim.CurrentTime {
					dot.lastTickTime = sim.CurrentTime
					dot.TickCount++
					dot.tickFn()
				}
			}
		},
	}

	dot.tickPeriod = dot.TickLength
	dot.Aura.Duration = dot.TickLength * time.Duration(dot.NumberOfTicks)

	dot.Aura.OnGain = func(aura *Aura, sim *Simulation) {
		dot.TakeSnapshot(sim)

		periodicOptions := basePeriodicOptions
		periodicOptions.Period = dot.tickPeriod
		dot.tickAction = NewPeriodicAction(sim, periodicOptions)
		sim.AddPendingAction(dot.tickAction)
	}
	dot.Aura.OnExpire = func(aura *Aura, sim *Simulation) {
		if dot.tickAction != nil {
			dot.tickAction.Cancel(sim)
			dot.tickAction = nil
		}
	}

	return dot
}

func TickFuncSnapshot(target *Unit, baseEffect SpellEffect) TickEffects {
	snapshotEffect := &SpellEffect{}
	return func(sim *Simulation, spell *Spell) func() {
		*snapshotEffect = baseEffect
		snapshotEffect.Target = target
		baseDamage := snapshotEffect.calculateBaseDamage(sim, spell) * snapshotEffect.DamageMultiplier
		snapshotEffect.DamageMultiplier = 1
		snapshotEffect.BaseDamage = BaseDamageConfigFlat(baseDamage)

		effectsFunc := ApplyEffectFuncDirectDamage(*snapshotEffect)
		return func() {
			effectsFunc(sim, target, spell)
		}
	}
}
func TickFuncAOESnapshot(env *Environment, baseEffect SpellEffect) TickEffects {
	snapshotEffect := &SpellEffect{}
	return func(sim *Simulation, spell *Spell) func() {
		target := spell.Unit.CurrentTarget
		*snapshotEffect = baseEffect
		snapshotEffect.Target = target
		baseDamage := snapshotEffect.calculateBaseDamage(sim, spell) * snapshotEffect.DamageMultiplier
		snapshotEffect.DamageMultiplier = 1
		snapshotEffect.BaseDamage = BaseDamageConfigFlat(baseDamage)

		effectsFunc := ApplyEffectFuncAOEDamage(env, *snapshotEffect)
		return func() {
			effectsFunc(sim, target, spell)
		}
	}
}
func TickFuncAOESnapshotCapped(env *Environment, aoeCap float64, baseEffect SpellEffect) TickEffects {
	snapshotEffect := &SpellEffect{}
	return func(sim *Simulation, spell *Spell) func() {
		target := spell.Unit.CurrentTarget
		*snapshotEffect = baseEffect
		snapshotEffect.Target = target
		baseDamage := snapshotEffect.calculateBaseDamage(sim, spell) * snapshotEffect.DamageMultiplier
		snapshotEffect.DamageMultiplier = 1
		snapshotEffect.BaseDamage = BaseDamageConfigFlat(baseDamage)

		effectsFunc := ApplyEffectFuncAOEDamageCapped(env, aoeCap, *snapshotEffect)
		return func() {
			effectsFunc(sim, target, spell)
		}
	}
}

func TickFuncApplyEffects(effectsFunc ApplySpellEffects) TickEffects {
	return func(sim *Simulation, spell *Spell) func() {
		return func() {
			effectsFunc(sim, spell.Unit.CurrentTarget, spell)
		}
	}
}
