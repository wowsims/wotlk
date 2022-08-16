package protection

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
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
	rotationLoop:
		for _, spellNumber := range prot.RotationInput {
			switch spellNumber {
			case int32(proto.ProtectionPaladin_Rotation_JudgementOfWisdom):
				if prot.JudgementOfWisdom.IsReady(sim) {
					prot.JudgementOfWisdom.Cast(sim, target)
					break rotationLoop
				}
			case int32(proto.ProtectionPaladin_Rotation_HammerOfWrath):
				if isExecutePhase && prot.HammerOfWrath.IsReady(sim) {
					prot.HammerOfWrath.Cast(sim, target)
					break rotationLoop
				}
			case int32(proto.ProtectionPaladin_Rotation_Consecration):
				if prot.Consecration.IsReady(sim) {
					prot.Consecration.Cast(sim, target)
					break rotationLoop
				}
			case int32(proto.ProtectionPaladin_Rotation_HolyWrath):
				if prot.HolyWrath.IsReady(sim) {
					prot.HolyWrath.Cast(sim, target)
					break rotationLoop
				}
			case int32(proto.ProtectionPaladin_Rotation_Exorcism):
				if prot.Exorcism.IsReady(sim) {
					prot.Exorcism.Cast(sim, target)
					break rotationLoop
				}
			case int32(proto.ProtectionPaladin_Rotation_ShieldOfRighteousness):
				if prot.ShieldOfRighteousness.IsReady(sim) {
					prot.ShieldOfRighteousness.Cast(sim, target)
					break rotationLoop
				}
			case int32(proto.ProtectionPaladin_Rotation_AvengersShield):
				if prot.AvengersShield.IsReady(sim) {
					prot.AvengersShield.Cast(sim, target)
					break rotationLoop
				}
			case int32(proto.ProtectionPaladin_Rotation_HammerOfTheRighteous):
				if prot.HammerOfTheRighteous.IsReady(sim) {
					prot.HammerOfTheRighteous.Cast(sim, target)
					break rotationLoop
				}
			}
		}
	}

	// All possible next events
	events := []time.Duration{
		nextSwingAt,
		prot.GCD.ReadyAt(),
		prot.JudgementOfWisdom.CD.ReadyAt(),
		prot.HammerOfWrath.CD.ReadyAt(),
		prot.HolyWrath.CD.ReadyAt(),
		prot.Consecration.CD.ReadyAt(),
		prot.Exorcism.CD.ReadyAt(),
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
