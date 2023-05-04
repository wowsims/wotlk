package core

import (
	"fmt"
	"time"
)

// A cast corresponds to any action which causes the in-game castbar to be
// shown, and activates the GCD. Note that a cast can also be instant, i.e.
// the effects are applied immediately even though the GCD is still activated.

// Callback for when a cast is finished, i.e. when the in-game castbar reaches full.
type OnCastComplete func(aura *Aura, sim *Simulation, spell *Spell)

type Hardcast struct {
	Expires    time.Duration
	ActionID   ActionID
	OnComplete func(*Simulation, *Unit)
	Target     *Unit
}

// Input for constructing the CastSpell function for a spell.
type CastConfig struct {
	// Default cast values with all static effects applied.
	DefaultCast Cast

	// Dynamic modifications for each cast.
	ModifyCast func(*Simulation, *Spell, *Cast)

	// Ignores haste when calculating the GCD and cast time for this cast.
	IgnoreHaste bool

	CD       Cooldown
	SharedCD Cooldown

	// Callbacks for providing additional custom behavior.
	OnCastComplete func(*Simulation, *Spell)
}

type Cast struct {
	// Amount of resource that will be consumed by this cast.
	Cost float64

	// The length of time the GCD will be on CD as a result of this cast.
	GCD time.Duration

	// The amount of time between the call to spell.Cast() and when the spell
	// effects are invoked.
	CastTime time.Duration

	// Additional GCD delay after the cast completes.
	ChannelTime time.Duration

	// Additional GCD delay after the cast ends. Never affected by cast speed.
	// This is typically used for latency.
	AfterCastDelay time.Duration
}

func (cast Cast) EffectiveTime() time.Duration {
	gcd := cast.GCD
	if cast.GCD != 0 {
		// TODO: isn't this wrong for spells like shadowfury, that have a reduced GCD?
		gcd = MaxDuration(GCDMin, gcd)
	}
	fullCastTime := cast.CastTime + cast.ChannelTime + cast.AfterCastDelay
	return MaxDuration(gcd, fullCastTime)
}

var emptyCast Cast

type CastFunc func(*Simulation, *Unit)
type CastSuccessFunc func(*Simulation, *Unit) bool

func (spell *Spell) makeCastFunc(config CastConfig, onCastComplete CastFunc) CastSuccessFunc {
	return spell.wrapCastFuncInit(config,
		spell.wrapCastFuncExtraCond(config,
			spell.wrapCastFuncCDsReady(config,
				spell.wrapCastFuncResources(config,
					spell.wrapCastFuncHaste(config,
						spell.wrapCastFuncGCD(config,
							spell.wrapCastFuncCooldown(config,
								spell.wrapCastFuncSharedCooldown(config,
									spell.makeCastFuncWait(config, onCastComplete)))))))))
}

func (spell *Spell) ApplyCostModifiers(cost float64) float64 {
	cost -= spell.Unit.PseudoStats.CostReduction
	cost = MaxFloat(0, cost*spell.Unit.PseudoStats.CostMultiplier)
	return MaxFloat(0, cost*spell.CostMultiplier)
}

func (spell *Spell) wrapCastFuncInit(config CastConfig, onCastComplete CastSuccessFunc) CastSuccessFunc {
	if spell.DefaultCast == emptyCast {
		return onCastComplete
	}

	if config.ModifyCast == nil {
		return func(sim *Simulation, target *Unit) bool {
			spell.CurCast = spell.DefaultCast
			return onCastComplete(sim, target)
		}
	} else {
		modifyCast := config.ModifyCast
		return func(sim *Simulation, target *Unit) bool {
			spell.CurCast = spell.DefaultCast
			cost := spell.CurCast.Cost
			modifyCast(sim, spell, &spell.CurCast)
			if cost != spell.CurCast.Cost {
				// Costs need to be modified using the unit and spell multipliers, so that
				// their affects are also visible in the spell.CanCast() function, which
				// does not invoke ModifyCast.
				panic("May not modify cost in ModifyCast!")
			}
			return onCastComplete(sim, target)
		}
	}
}

func (spell *Spell) wrapCastFuncExtraCond(config CastConfig, onCastComplete CastSuccessFunc) CastSuccessFunc {
	if spell.ExtraCastCondition == nil {
		return onCastComplete
	} else {
		return func(sim *Simulation, target *Unit) bool {
			if spell.ExtraCastCondition(sim, target) {
				return onCastComplete(sim, target)
			} else {
				if sim.Log != nil {
					sim.Log("Failed cast because of extra condition")
				}
				return false
			}
		}
	}
}

func (spell *Spell) wrapCastFuncCDsReady(config CastConfig, onCastComplete CastSuccessFunc) CastSuccessFunc {
	if spell.Unit.PseudoStats.GracefulCastCDFailures {
		return func(sim *Simulation, target *Unit) bool {
			if spell.IsReady(sim) {
				return onCastComplete(sim, target)
			} else {
				if sim.Log != nil {
					sim.Log("Failed cast because of CDs")
				}
				return false
			}
		}
	} else {
		return onCastComplete
	}
}

func (spell *Spell) wrapCastFuncResources(config CastConfig, onCastComplete CastFunc) CastSuccessFunc {
	if spell.Cost == nil {
		return func(sim *Simulation, target *Unit) bool {
			onCastComplete(sim, target)
			return true
		}
	}

	return func(sim *Simulation, target *Unit) bool {
		if !spell.Cost.MeetsRequirement(spell) {
			if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
				spell.Cost.LogCostFailure(sim, spell)
			}
			return false
		}
		onCastComplete(sim, target)
		return true
	}
}

func (spell *Spell) wrapCastFuncHaste(config CastConfig, onCastComplete CastFunc) CastFunc {
	if config.IgnoreHaste || (spell.DefaultCast.GCD == 0 && spell.DefaultCast.CastTime == 0 && spell.DefaultCast.ChannelTime == 0) {
		return onCastComplete
	}

	return func(sim *Simulation, target *Unit) {
		spell.CurCast.GCD = spell.Unit.ApplyCastSpeed(spell.CurCast.GCD)
		spell.CurCast.CastTime = spell.Unit.ApplyCastSpeedForSpell(spell.CurCast.CastTime, spell)
		spell.CurCast.ChannelTime = spell.Unit.ApplyCastSpeedForSpell(spell.CurCast.ChannelTime, spell)

		onCastComplete(sim, target)
	}
}

func (spell *Spell) wrapCastFuncGCD(config CastConfig, onCastComplete CastFunc) CastFunc {
	if spell.DefaultCast == emptyCast { // spells that are not actually cast (e.g. auto attacks, procs)
		return onCastComplete
	}

	if spell.DefaultCast.GCD == 0 { // mostly cooldowns (e.g. nature's swiftness, presence of mind)
		return func(sim *Simulation, target *Unit) {
			if hc := spell.Unit.Hardcast; hc.Expires > sim.CurrentTime {
				panic(fmt.Sprintf("Trying to cast %s but casting/channeling %v for %s, curTime = %s", spell.ActionID, hc.ActionID, hc.Expires-sim.CurrentTime, sim.CurrentTime))
			}
			onCastComplete(sim, target)
		}
	}

	return func(sim *Simulation, target *Unit) {
		// By panicking if spell is on CD, we force each sim to properly check for their own CDs.
		if spell.CurCast.GCD != 0 && !spell.Unit.GCD.IsReady(sim) {
			panic(fmt.Sprintf("Trying to cast %s but GCD on cooldown for %s, curTime = %s", spell.ActionID, spell.Unit.GCD.TimeToReady(sim), sim.CurrentTime))
		}

		if hc := spell.Unit.Hardcast; hc.Expires > sim.CurrentTime {
			panic(fmt.Sprintf("Trying to cast %s but casting/channeling %v for %s, curTime = %s", spell.ActionID, hc.ActionID, hc.Expires-sim.CurrentTime, sim.CurrentTime))
		}

		effectiveTime := spell.CurCast.EffectiveTime()
		if effectiveTime != 0 {
			spell.SpellMetrics[target.UnitIndex].TotalCastTime += effectiveTime
			spell.Unit.SetGCDTimer(sim, sim.CurrentTime+effectiveTime)
		}

		onCastComplete(sim, target)
	}
}

func (spell *Spell) wrapCastFuncCooldown(config CastConfig, onCastComplete CastFunc) CastFunc {
	if config.CD.Timer == nil {
		return onCastComplete
	}

	if config.CD.Duration == 0 {
		panic("Cooldown specified but no duration!")
	}

	return func(sim *Simulation, target *Unit) {
		// By panicking if spell is on CD, we force each sim to properly check for their own CDs.
		if !spell.CD.IsReady(sim) {
			panic(fmt.Sprintf("Trying to cast %s but is still on cooldown for %s, curTime = %s", spell.ActionID, spell.CD.TimeToReady(sim), sim.CurrentTime))
		}

		spell.CD.Set(sim.CurrentTime + spell.CurCast.CastTime + spell.CD.Duration)

		onCastComplete(sim, target)
	}
}

func (spell *Spell) wrapCastFuncSharedCooldown(config CastConfig, onCastComplete CastFunc) CastFunc {
	if config.SharedCD.Timer == nil {
		return onCastComplete
	}

	if config.SharedCD.Duration == 0 {
		panic("SharedCooldown specified but no duration!")
	}

	return func(sim *Simulation, target *Unit) {
		// By panicking if spell is on CD, we force each sim to properly check for their own CDs.
		if !spell.SharedCD.IsReady(sim) {
			panic(fmt.Sprintf("Trying to cast %s but is still on shared cooldown for %s, curTime = %s", spell.ActionID, spell.SharedCD.TimeToReady(sim), sim.CurrentTime))
		}

		spell.SharedCD.Set(sim.CurrentTime + spell.CurCast.CastTime + spell.SharedCD.Duration)

		onCastComplete(sim, target)
	}
}

func (spell *Spell) makeCastFuncWait(config CastConfig, onCastComplete CastFunc) CastFunc {
	if !spell.Flags.Matches(SpellFlagNoOnCastComplete) {
		oldOnCastComplete1 := onCastComplete
		configOnCastComplete := config.OnCastComplete
		onCastComplete = func(sim *Simulation, target *Unit) {
			oldOnCastComplete1(sim, target)
			if configOnCastComplete != nil {
				configOnCastComplete(sim, spell)
			}
			spell.Unit.OnCastComplete(sim, spell)
		}
	}

	if spell.Cost != nil {
		oldOnCastComplete2 := onCastComplete
		onCastComplete = func(sim *Simulation, target *Unit) {
			spell.Cost.SpendCost(sim, spell)
			oldOnCastComplete2(sim, target)
		}
	}

	if spell.DefaultCast.ChannelTime > 0 {
		return func(sim *Simulation, target *Unit) {
			spell.Unit.Hardcast = Hardcast{Expires: sim.CurrentTime + spell.CurCast.ChannelTime, ActionID: spell.ActionID}
			if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
				spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s, Effective Time = %s)",
					spell.ActionID, MaxFloat(0, spell.CurCast.Cost), spell.CurCast.CastTime, spell.CurCast.EffectiveTime())
				spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
			}
			onCastComplete(sim, target)
		}
	}

	if spell.DefaultCast.CastTime == 0 {
		if spell.Flags.Matches(SpellFlagNoLogs) {
			return onCastComplete
		} else {
			return func(sim *Simulation, target *Unit) {
				if sim.Log != nil {
					spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s, Effective Time = %s)",
						spell.ActionID, MaxFloat(0, spell.CurCast.Cost), spell.CurCast.CastTime, spell.CurCast.EffectiveTime())
					spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
				}
				onCastComplete(sim, target)
			}
		}
	} else {
		if !spell.Flags.Matches(SpellFlagNoLogs) {
			oldOnCastComplete3 := onCastComplete
			onCastComplete = func(sim *Simulation, target *Unit) {
				if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
					spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
				}
				oldOnCastComplete3(sim, target)
			}
		}

		return func(sim *Simulation, target *Unit) {
			if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
				spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s, Effective Time = %s)",
					spell.ActionID, MaxFloat(0, spell.CurCast.Cost), spell.CurCast.CastTime, spell.CurCast.EffectiveTime())
			}

			// For instant-cast spells we can skip creating an aura.
			if spell.CurCast.CastTime == 0 {
				onCastComplete(sim, target)
			} else {
				spell.Unit.Hardcast = Hardcast{
					Expires:    sim.CurrentTime + spell.CurCast.CastTime,
					ActionID:   spell.ActionID,
					OnComplete: onCastComplete,
					Target:     target,
				}

				// If hardcast and GCD happen at the same time then we don't need a separate action.
				if spell.Unit.Hardcast.Expires != spell.Unit.NextGCDAt() {
					spell.Unit.newHardcastAction(sim)
				}
			}
		}
	}
}
