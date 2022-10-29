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

//nolint:unused
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

	// Forced CD remaining on HotR/ShoR to cast the other. Can't be exactly 3sec or lusted consecration GCDs will desync us.
	gapSlack := time.Millisecond * 4000

	if prot.GCD.IsReady(sim) {

		if !prot.Rotation.UseCustomPrio {

			// Standard rotation. Enforce 6sec CDs to have 1 GCD between, filling with 9sec abilities.
			if isExecutePhase && prot.HammerOfWrath.IsReady(sim) {
				// Always cast HoW if ready
				prot.HammerOfWrath.Cast(sim, target)
			} else if prot.Rotation.HammerFirst && prot.HammerOfTheRighteous.IsReady(sim) {
				// Always cast HotR if ready
				prot.HammerOfTheRighteous.Cast(sim, target)
			} else if prot.Rotation.HammerFirst &&
				prot.ShieldOfRighteousness.IsReady(sim) &&
				(prot.HammerOfTheRighteous.TimeToReady(sim) < gapSlack) {
				// Cast ShoR if ready but only if you've spent a global since HotR
				prot.ShieldOfRighteousness.Cast(sim, target)
			} else if !prot.Rotation.HammerFirst && prot.ShieldOfRighteousness.IsReady(sim) {
				// Always cast ShoR if ready
				prot.ShieldOfRighteousness.Cast(sim, target)
			} else if !prot.Rotation.HammerFirst &&
				prot.HammerOfTheRighteous.IsReady(sim) &&
				(prot.ShieldOfRighteousness.TimeToReady(sim) < gapSlack) {
				// Cast HotR if ready but only if you've spent a global since ShoR
				prot.HammerOfTheRighteous.Cast(sim, target)
			} else if prot.HolyShield.IsReady(sim) {
				// Top priority 9 is Holy Shield
				prot.HolyShield.Cast(sim, target)
			} else if prot.HasGlyphAS && prot.AvengersShield.IsReady(sim) {
				// AS prio if glyphed. This will push out Cons/Judge which may not be good, but assumed desired based on the glyph choice
				prot.AvengersShield.Cast(sim, target)
			} else if prot.Consecration.IsReady(sim) {
				prot.Consecration.Cast(sim, target)
			} else if prot.JudgementOfWisdom.IsReady(sim) {
				prot.JudgementOfWisdom.Cast(sim, target)
			}
			// Do not ever cast Exorcism, Holy Wrath, or unglyphed AS
			// TODO: Possible to dynamically affect Judgement<>AS priority based on Libram bonus at SBV softcap?

		} else {

			// Custom rotation
		rotationLoop:
			for _, spellNumber := range prot.RotationInput {
				// In priority order, fire the first spell which is Ready
				// Still enforce Hammer/Shield being separated by a GCD
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
					if prot.ShieldOfRighteousness.IsReady(sim) && (prot.HammerOfTheRighteous.TimeToReady(sim) < gapSlack) {
						prot.ShieldOfRighteousness.Cast(sim, target)
						break rotationLoop
					}
				case int32(proto.ProtectionPaladin_Rotation_AvengersShield):
					if prot.AvengersShield.IsReady(sim) {
						prot.AvengersShield.Cast(sim, target)
						break rotationLoop
					}
				case int32(proto.ProtectionPaladin_Rotation_HammerOfTheRighteous):
					if prot.HammerOfTheRighteous.IsReady(sim) && (prot.ShieldOfRighteousness.TimeToReady(sim) < gapSlack) {
						prot.HammerOfTheRighteous.Cast(sim, target)
						break rotationLoop
					}
				case int32(proto.ProtectionPaladin_Rotation_HolyShield):
					if prot.HolyShield.IsReady(sim) {
						prot.HolyShield.Cast(sim, target)
						break rotationLoop
					}
				}
			}

		}

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
