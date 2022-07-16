package retribution

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (ret *RetributionPaladin) OnGCDReady(sim *core.Simulation) {
	ret.tryUseGCD(sim)
}

func (ret *RetributionPaladin) tryUseGCD(sim *core.Simulation) {
	if !ret.openerCompleted {
		ret.openingRotation(sim)
		return
	}
	ret.mainRotation(sim)

	if ret.GCD.IsReady(sim) {
		ret.DoNothing() // this means we had nothing to do and we are ok
	}
}

func (ret *RetributionPaladin) openingRotation(sim *core.Simulation) {
	ret.SealOfVengeance.Cast(sim, nil)
	ret.AutoAttacks.EnableAutoSwing(sim)
	ret.openerCompleted = true
}

func (ret *RetributionPaladin) mainRotation(sim *core.Simulation) {

	// Setup
	target := ret.CurrentTarget

	// gcdCD := ret.GCD.TimeToReady(sim)
	nextSwingAt := ret.AutoAttacks.NextAttackAt()

	// Needs 2pc t10 to be effective.
	if ret.GCD.IsReady(sim) {
		if nextSwingAt.Milliseconds() > 1500 {
			switch {
			case ret.JudgementOfWisdom.IsReady(sim):
				ret.JudgementOfWisdom.Cast(sim, target)
			case ret.CrusaderStrike.IsReady(sim):
				ret.CrusaderStrike.Cast(sim, target)
			case ret.DivineStorm.IsReady(sim):
				ret.DivineStorm.Cast(sim, target)
			case ret.Exorcism.IsReady(sim):
				ret.Exorcism.Cast(sim, target)
			case ret.Consecration.IsReady(sim):
				ret.Consecration.Cast(sim, target)
			}
		} else {
			switch {
			case ret.DivineStorm.IsReady(sim):
				ret.DivineStorm.Cast(sim, target)
			case ret.JudgementOfWisdom.IsReady(sim):
				ret.JudgementOfWisdom.Cast(sim, target)
			case ret.CrusaderStrike.IsReady(sim):
				ret.CrusaderStrike.Cast(sim, target)
			case ret.Exorcism.IsReady(sim):
				ret.Exorcism.Cast(sim, target)
			case ret.Consecration.IsReady(sim):
				ret.Consecration.Cast(sim, target)
			}
		}
	}

	// All possible next events
	events := []time.Duration{
		nextSwingAt,
		ret.GCD.ReadyAt(),
		// ret.JudgementOfWisdom.CD.ReadyAt(),
		// ret.CrusaderStrike.CD.ReadyAt(),
		// ret.DivineStorm.CD.ReadyAt(),
		// ret.Consecration.CD.ReadyAt(),
		// ret.Exorcism.CD.ReadyAt(),
	}

	ret.waitUntilNextEvent(sim, events)
}

// Helper function for finding the next event
func (ret *RetributionPaladin) waitUntilNextEvent(sim *core.Simulation, events []time.Duration) {
	// Find the minimum possible next event that is greater than the current time
	nextEventAt := time.Duration(math.MaxInt64) // any event will happen before forever.
	for _, elem := range events {
		if elem > sim.CurrentTime && elem < nextEventAt {
			nextEventAt = elem
		}
	}
	// If the next action is  the GCD, just return
	if nextEventAt == ret.GCD.ReadyAt() {
		return
	}

	ret.WaitUntil(sim, nextEventAt)
}
