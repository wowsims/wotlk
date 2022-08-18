package core

import (
	"fmt"
	"time"

	"github.com/wowsims/wotlk/sim/core/stats"
)

// A cast corresponds to any action which causes the in-game castbar to be
// shown, and activates the GCD. Note that a cast can also be instant, i.e.
// the effects are applied immediately even though the GCD is still activated.

// Callback for when a cast is finished, i.e. when the in-game castbar reaches full.
type OnCastComplete func(aura *Aura, sim *Simulation, spell *Spell)

type Hardcast struct {
	Expires    time.Duration
	OnComplete func(*Simulation, *Unit)
	Target     *Unit
}

func (hc *Hardcast) OnExpire(sim *Simulation) {
	hc.OnComplete(sim, hc.Target)
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
	AfterCast      func(*Simulation, *Spell)
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

type CastFunc func(*Simulation, *Unit)
type CastSuccessFunc func(*Simulation, *Unit) bool

func (spell *Spell) makeCastFunc(config CastConfig, onCastComplete CastFunc) CastSuccessFunc {
	return spell.wrapCastFuncInit(config,
		spell.wrapCastFuncResources(config,
			spell.wrapCastFuncHaste(config,
				spell.wrapCastFuncGCD(config,
					spell.wrapCastFuncCooldown(config,
						spell.wrapCastFuncSharedCooldown(config,
							spell.makeCastFuncWait(config, onCastComplete)))))))
}

func (spell *Spell) ApplyCostModifiers(cost float64) float64 {
	if spell.Unit.PseudoStats.NoCost {
		return 0
	} else {
		cost -= spell.BaseCost * (1 - spell.Unit.PseudoStats.CostMultiplier)
		cost -= spell.Unit.PseudoStats.CostReduction
		return MaxFloat(0, cost*spell.CostMultiplier)
	}
}

func (spell *Spell) wrapCastFuncInit(config CastConfig, onCastComplete CastSuccessFunc) CastSuccessFunc {
	empty := Cast{}
	if config.DefaultCast == empty {
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
			modifyCast(sim, spell, &spell.CurCast)
			return onCastComplete(sim, target)
		}
	}
}

func (spell *Spell) wrapCastFuncResources(config CastConfig, onCastComplete CastFunc) CastSuccessFunc {
	if spell.ResourceType == 0 || config.DefaultCast.Cost == 0 {
		if spell.ResourceType != 0 {
			panic("ResourceType set for spell " + spell.ActionID.String() + " but no cost")
		}
		if config.DefaultCast.Cost != 0 {
			panic("Cost set for spell " + spell.ActionID.String() + " but no ResourceType")
		}
		return func(sim *Simulation, target *Unit) bool {
			onCastComplete(sim, target)
			return true
		}
	}

	switch spell.ResourceType {
	case stats.Mana:
		return func(sim *Simulation, target *Unit) bool {
			spell.CurCast.Cost = spell.ApplyCostModifiers(spell.CurCast.Cost)
			if spell.Unit.CurrentMana() < spell.CurCast.Cost {
				if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
					spell.Unit.Log(sim, "Failed casting %s, not enough mana. (Current Mana = %0.03f, Mana Cost = %0.03f)",
						spell.ActionID, spell.Unit.CurrentMana(), spell.CurCast.Cost)
				}
				return false
			}

			// Mana is subtracted at the end of the cast.
			onCastComplete(sim, target)
			return true
		}
	case stats.Rage:
		return func(sim *Simulation, target *Unit) bool {
			spell.CurCast.Cost = spell.ApplyCostModifiers(spell.CurCast.Cost)
			if spell.Unit.CurrentRage() < spell.CurCast.Cost {
				return false
			}
			spell.Unit.SpendRage(sim, spell.CurCast.Cost, spell.ResourceMetrics)
			onCastComplete(sim, target)
			return true
		}
	case stats.Energy:
		return func(sim *Simulation, target *Unit) bool {
			spell.CurCast.Cost = spell.ApplyCostModifiers(spell.CurCast.Cost)
			if spell.Unit.CurrentEnergy() < spell.CurCast.Cost {
				return false
			}
			spell.Unit.SpendEnergy(sim, spell.CurCast.Cost, spell.ResourceMetrics)
			onCastComplete(sim, target)
			return true
		}
	case stats.RunicPower:
		return func(sim *Simulation, target *Unit) bool {
			// Rune spending is currently handled in DK codebase.
			// This verifies that the user has the resources but does not spend.
			if spell.CurCast.Cost != 0 {
				cost := RuneCost(spell.CurCast.Cost)
				if !cost.HasRune() {
					if float64(cost.RunicPower()) > spell.Unit.CurrentRunicPower() {
						return false
					}
				} else {
					// Given cost might not be what is actually paid.
					//  Calculate what combination of runes can actually pay for this spell.
					optCost := spell.Unit.OptimalRuneCost(cost)
					if optCost == 0 { // no combo of runes to fulfill cost
						return false
					}
					spell.CurCast.Cost = float64(optCost) // assign chosen runes to the cost
				}
			}
			onCastComplete(sim, target)
			return true
		}
	}

	panic("Invalid resource type")
}

func (spell *Spell) wrapCastFuncHaste(config CastConfig, onCastComplete CastFunc) CastFunc {
	if config.IgnoreHaste || (config.DefaultCast.GCD == 0 && config.DefaultCast.CastTime == 0 && config.DefaultCast.ChannelTime == 0) {
		return onCastComplete
	}

	return func(sim *Simulation, target *Unit) {
		spell.CurCast.GCD = spell.Unit.ApplyCastSpeed(spell.CurCast.GCD)
		spell.CurCast.CastTime = spell.Unit.ApplyCastSpeed(spell.CurCast.CastTime)
		spell.CurCast.ChannelTime = spell.Unit.ApplyCastSpeed(spell.CurCast.ChannelTime)

		onCastComplete(sim, target)
	}
}

func (spell *Spell) wrapCastFuncGCD(config CastConfig, onCastComplete CastFunc) CastFunc {
	if config.DefaultCast.GCD == 0 {
		return onCastComplete
	}

	return func(sim *Simulation, target *Unit) {
		// By panicking if spell is on CD, we force each sim to properly check for their own CDs.
		if spell.CurCast.GCD != 0 && !spell.Unit.GCD.IsReady(sim) {
			panic(fmt.Sprintf("Trying to cast %s but GCD on cooldown for %s", spell.ActionID, spell.Unit.GCD.TimeToReady(sim)))
		}

		gcd := spell.CurCast.GCD
		if spell.CurCast.GCD != 0 {
			gcd = MaxDuration(GCDMin, gcd)
		}

		fullCastTime := spell.CurCast.CastTime + spell.CurCast.ChannelTime + spell.CurCast.AfterCastDelay

		if gcd != 0 || fullCastTime != 0 {
			spell.Unit.SetGCDTimer(sim, sim.CurrentTime+MaxDuration(gcd, fullCastTime))
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
			panic(fmt.Sprintf("Trying to cast %s but is still on cooldown for %s", spell.ActionID, spell.CD.TimeToReady(sim)))
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
			panic(fmt.Sprintf("Trying to cast %s but is still on shared cooldown for %s", spell.ActionID, spell.SharedCD.TimeToReady(sim)))
		}

		spell.SharedCD.Set(sim.CurrentTime + spell.CurCast.CastTime + spell.SharedCD.Duration)

		onCastComplete(sim, target)
	}
}

func (spell *Spell) makeCastFuncWait(config CastConfig, onCastComplete CastFunc) CastFunc {
	if !spell.Flags.Matches(SpellFlagNoOnCastComplete) {
		configOnCastComplete := config.OnCastComplete
		configAfterCast := config.AfterCast
		oldOnCastComplete1 := onCastComplete
		onCastComplete = func(sim *Simulation, target *Unit) {
			spell.Unit.OnCastComplete(sim, spell)
			if configOnCastComplete != nil {
				configOnCastComplete(sim, spell)
			}
			oldOnCastComplete1(sim, target)
			if configAfterCast != nil {
				configAfterCast(sim, spell)
			}
		}
	}

	if spell.ResourceType == stats.Mana && config.DefaultCast.Cost != 0 {
		oldOnCastComplete2 := onCastComplete
		onCastComplete = func(sim *Simulation, target *Unit) {
			if spell.CurCast.Cost > 0 {
				spell.Unit.SpendMana(sim, spell.CurCast.Cost, spell.ResourceMetrics)
				spell.Unit.PseudoStats.FiveSecondRuleRefreshTime = sim.CurrentTime + time.Second*5
			}
			oldOnCastComplete2(sim, target)
		}
	}

	if config.DefaultCast.CastTime == 0 {
		if spell.Flags.Matches(SpellFlagNoLogs) {
			return onCastComplete
		} else {
			return func(sim *Simulation, target *Unit) {
				if sim.Log != nil {
					// Hunter fake cast has no ID.
					if !spell.ActionID.IsEmptyAction() {
						spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s)",
							spell.ActionID, MaxFloat(0, spell.CurCast.Cost), spell.CurCast.CastTime)
						spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
					}
				}
				onCastComplete(sim, target)
			}
		}
	} else {
		if !spell.Flags.Matches(SpellFlagNoLogs) {
			oldOnCastComplete3 := onCastComplete
			onCastComplete = func(sim *Simulation, target *Unit) {
				if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
					// Hunter fake cast has no ID.
					if !spell.ActionID.SameAction(ActionID{}) {
						spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
					}
				}
				oldOnCastComplete3(sim, target)
			}
		}

		return func(sim *Simulation, target *Unit) {
			if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
				spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s)",
					spell.ActionID, MaxFloat(0, spell.CurCast.Cost), spell.CurCast.CastTime)
			}

			// For instant-cast spells we can skip creating an aura.
			if spell.CurCast.CastTime == 0 {
				onCastComplete(sim, target)
			} else {
				spell.Unit.Hardcast.Expires = sim.CurrentTime + spell.CurCast.CastTime
				spell.Unit.Hardcast.OnComplete = onCastComplete
				spell.Unit.Hardcast.Target = target

				// If hardcast and GCD happen at the same time then we don't need a separate action.
				if spell.Unit.Hardcast.Expires != spell.Unit.NextGCDAt() {
					spell.Unit.newHardcastAction(sim)
				}

				if spell.Unit.AutoAttacks.IsEnabled() {
					// Delay autoattacks until the cast is complete.
					spell.Unit.AutoAttacks.DelayMeleeUntil(sim, spell.Unit.Hardcast.Expires)
				}
			}
		}
	}
}
