package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/stats"
)

type TickEffects func(*Simulation, *Dot) func()

type Dot struct {
	Spell *Spell

	// Embed Aura so we can use IsActive/Refresh/etc directly.
	*Aura

	NumberOfTicks int           // number of ticks over the whole duration
	TickLength    time.Duration // time between each tick

	// If true, tick length will be shortened based on casting speed.
	AffectedByCastSpeed bool

	TickEffects        TickEffects
	snapshotEffect     *SpellEffect
	snapshotMultiplier float64

	tickFn     func()
	tickAction *PendingAction
	tickPeriod time.Duration

	// Number of ticks since last call to Apply().
	TickCount int

	lastTickTime time.Duration
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
	oldCrit := dot.snapshotEffect.BonusSpellCritRating

	// by not calling 'dot.TakeSnapshot' we dont recalculate dmg multiplier
	dot.tickFn = dot.TickEffects(sim, dot)

	dot.snapshotEffect.BonusSpellCritRating = oldCrit

	dot.RecomputeAuraDuration()
	dot.Aura.Refresh(sim)
}

// Reapply will reset the current DOT by being reapplied.
// This will re-snapshot.
func (dot *Dot) Reapply(sim *Simulation) {
	dot.RecomputeAuraDuration() // calculate aura duration
	dot.TakeSnapshot(sim)       // snapshots dmg / sp / crit
	dot.Aura.Refresh(sim)       // resets aura with new duration
}

// func (dot *Dot) Refresh(sim *Simulation, restartTimer bool) {
// 	dot.Aura.Refresh(sim)
// 	if restartTimer {
// 		dot.tickAction.NextActionAt = sim.CurrentTime + dot.tickPeriod
// 	}
// }

func (dot *Dot) Apply(sim *Simulation) {
	dot.Cancel(sim)
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
func (dot *Dot) TakeSnapshot(sim *Simulation) {
	dot.snapshotMultiplier = dot.snapshotEffect.snapshotAttackModifiers(dot.Spell)
	dot.tickFn = dot.TickEffects(sim, dot)
}

// Forces an instant tick. Does not reset the tick timer or aura duration,
// the tick is simply an extra tick.
func (dot *Dot) TickOnce() {
	dot.tickFn()
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

	oldOnGain := dot.Aura.OnGain
	oldOnExpire := dot.Aura.OnExpire
	dot.Aura.OnGain = func(aura *Aura, sim *Simulation) {
		dot.TakeSnapshot(sim)

		periodicOptions := basePeriodicOptions
		periodicOptions.Period = dot.tickPeriod
		dot.tickAction = NewPeriodicAction(sim, periodicOptions)
		sim.AddPendingAction(dot.tickAction)

		if oldOnGain != nil {
			oldOnGain(aura, sim)
		}
	}
	dot.Aura.OnExpire = func(aura *Aura, sim *Simulation) {
		if dot.tickAction != nil {
			dot.tickAction.Cancel(sim)
			dot.tickAction = nil
		}

		if oldOnExpire != nil {
			oldOnExpire(aura, sim)
		}
	}
	dot.snapshotEffect = &SpellEffect{}

	return dot
}

func TickFuncSnapshot(target *Unit, baseEffect SpellEffect) TickEffects {
	return func(sim *Simulation, dot *Dot) func() {
		*dot.snapshotEffect = baseEffect
		dot.snapshotEffect.DamageMultiplier *= dot.snapshotMultiplier
		dot.snapshotEffect.Target = target

		baseDamage := dot.snapshotEffect.calculateBaseDamage(sim, dot.Spell)
		dot.snapshotEffect.BonusSpellCritRating = dot.snapshotEffect.BonusSpellCritRating +
			dot.Spell.Unit.GetStat(stats.SpellCrit) + dot.Spell.Unit.PseudoStats.BonusSpellCritRating +
			target.PseudoStats.BonusCritRatingTaken + target.PseudoStats.BonusSpellCritRatingTaken
		dot.snapshotEffect.DamageMultiplier = 1
		dot.snapshotEffect.BaseDamage = BaseDamageConfigFlat(baseDamage)

		effectsFunc := ApplyEffectFuncDirectDamage(*dot.snapshotEffect)
		return func() {
			effectsFunc(sim, target, dot.Spell)
		}
	}
}

func TickFuncAOESnapshot(env *Environment, baseEffect SpellEffect) TickEffects {
	return func(sim *Simulation, dot *Dot) func() {
		target := dot.Spell.Unit.CurrentTarget
		*dot.snapshotEffect = baseEffect
		dot.snapshotEffect.Target = target
		baseDamage := dot.snapshotEffect.calculateBaseDamage(sim, dot.Spell)
		dot.snapshotEffect.DamageMultiplier = 1
		dot.snapshotEffect.BaseDamage = BaseDamageConfigFlat(baseDamage)

		effectsFunc := ApplyEffectFuncAOEDamage(env, *dot.snapshotEffect)
		return func() {
			effectsFunc(sim, target, dot.Spell)
		}
	}
}
func TickFuncAOESnapshotCapped(env *Environment, aoeCap float64, baseEffect SpellEffect) TickEffects {
	return func(sim *Simulation, dot *Dot) func() {
		target := dot.Spell.Unit.CurrentTarget
		*dot.snapshotEffect = baseEffect
		dot.snapshotEffect.Target = target
		baseDamage := dot.snapshotEffect.calculateBaseDamage(sim, dot.Spell)
		dot.snapshotEffect.DamageMultiplier = 1
		dot.snapshotEffect.BaseDamage = BaseDamageConfigFlat(baseDamage)

		effectsFunc := ApplyEffectFuncAOEDamageCapped(env, aoeCap, *dot.snapshotEffect)
		return func() {
			effectsFunc(sim, target, dot.Spell)
		}
	}
}

func TickFuncApplyEffects(effectsFunc ApplySpellEffects) TickEffects {
	return func(sim *Simulation, dot *Dot) func() {
		return func() {
			effectsFunc(sim, dot.Spell.Unit.CurrentTarget, dot.Spell)
		}
	}
}
