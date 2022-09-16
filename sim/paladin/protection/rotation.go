package protection

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	//"github.com/wowsims/wotlk/sim/core/proto"
)

func (prot *ProtectionPaladin) OnGCDReady(sim *core.Simulation) {
	prot.SelectedRotation(sim)

	if prot.GCD.IsReady(sim) {
		prot.DoNothing()
	}
}

func (prot *ProtectionPaladin) nextCDAt(sim *core.Simulation) time.Duration {
	nextCDAt := core.MinDuration(prot.HolyShield.ReadyAt(), prot.JudgementOfWisdom.ReadyAt())
	nextCDAt = core.MinDuration(nextCDAt, prot.Consecration.ReadyAt())
	return nextCDAt
}

func (prot *ProtectionPaladin) customRotation(sim *core.Simulation) {
	// Setup
	target := prot.CurrentTarget

	nextSwingAt := prot.AutoAttacks.NextAttackAt()
	isExecutePhase := sim.IsExecutePhase20()

	if prot.GCD.IsReady(sim) {
	//rotationLoop:
	
		if isExecutePhase && prot.HammerOfWrath.IsReady(sim) {
			// Always cast HoW if ready
			prot.HammerOfWrath.Cast(sim, target)
		} else
		if prot.HammerOfTheRighteous.IsReady(sim) {
			// Always cast HotR if ready
			prot.HammerOfTheRighteous.Cast(sim, target)
		} else
		if prot.ShieldOfRighteousness.IsReady(sim) && (prot.HammerOfTheRighteous.TimeToReady(sim) < time.Millisecond * 3100) {
			// Cast ShoR if ready but only if you've spent a global since HotR
			prot.ShieldOfRighteousness.Cast(sim, target)
		} else
		if prot.HasGlyphAS && prot.AvengersShield.IsReady(sim) {
			// AS prio if glyphed.
			prot.AvengersShield.Cast(sim, target)
		} else
		if prot.HolyShield.IsReady(sim) {
			// Top priority 9 is Holy Shield
			prot.HolyShield.Cast(sim, target)
		} else
		if prot.JudgementOfWisdom.IsReady(sim) {
			// Lower prio 9, should be last prio for DPS but higher prio for Libram
			prot.JudgementOfWisdom.Cast(sim, target)
		} else
		if prot.Consecration.IsReady(sim) {
			// Lower priority 9, Judgement is better damage but triggers Libram of Obstruction
			prot.Consecration.Cast(sim, target)
		}	
		// Do not ever cast Exorcism
		// Do not ever cast Holy Wrath
		// Do not cast AS unglyphed ... yet ... dropping Judgements TBD optional behavior
	

	}

	// All possible next events
	events := []time.Duration{
		nextSwingAt,
		prot.GCD.ReadyAt(),
		prot.JudgementOfWisdom.ReadyAt(),
		prot.HammerOfWrath.ReadyAt(),
		prot.Consecration.ReadyAt(),
		prot.HolyWrath.ReadyAt(),
		prot.Exorcism.ReadyAt(),
		prot.ShieldOfRighteousness.ReadyAt(),
		prot.AvengersShield.ReadyAt(),
		prot.HammerOfTheRighteous.ReadyAt(),
		prot.HolyShield.ReadyAt(),
	}

	prot.waitUntilNextEvent(sim, events, prot.customRotation)

}

// Helper function for finding the next event
func (prot *ProtectionPaladin) waitUntilNextEvent(sim *core.Simulation, events []time.Duration, rotationCallback func(*core.Simulation)) {
	// Find the minimum possible next event that is greater than the current time
	nextEventAt := time.Duration(math.MaxInt64) // any event will happen before forever.
	for _, elem := range events {
		if elem > sim.CurrentTime && elem < nextEventAt {
			nextEventAt = elem
		}
	}
	// If the next action is  the GCD, just return
	if nextEventAt == prot.GCD.ReadyAt() {
		return
	}

	// Otherwise add a pending action for the next time
	pa := &core.PendingAction{
		Priority:     core.ActionPriorityLow,
		OnAction:     rotationCallback,
		NextActionAt: nextEventAt,
	}

	sim.AddPendingAction(pa)
}
