package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

type FrostRotationAction uint8

const (
	FrostRotationAction_Wait FrostRotationAction = iota
	FrostRotationAction_IT
	FrostRotationAction_PS
	FrostRotationAction_Obli
	FrostRotationAction_BS
	FrostRotationAction_BT
	FrostRotationAction_UA
	FrostRotationAction_Pesti
	FrostRotationAction_FS
	FrostRotationAction_HW
)

type FrostRotation struct {
	actionIdx    int
	mainSequence []FrostRotationAction
}

func (deathKnight *DeathKnight) setupFrostRotation() {
	fr := &deathKnight.FrostRotation
	fr.actionIdx = 0

	//if deathKnight.Options.RefreshHornOfWinter {
	//	fr.mainSequence.append(FrostRotationAction_HW)
	//}
	//fr.mainSequence.append(FrostRotationAction_IT)
	//fr.mainSequence.append(FrostRotationAction_PS)
}

func (deathKnight *DeathKnight) doFrostRotation(sim *core.Simulation) {
	if !deathKnight.Talents.HowlingBlast {
		return
	}

	target := deathKnight.CurrentTarget

	if deathKnight.ShouldHornOfWinter(sim) {
		deathKnight.HornOfWinter.Cast(sim, target)
	} else if (!deathKnight.TargetHasDisease(FrostFeverAuraLabel, target) || deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim) < 6*time.Second) && deathKnight.CanIcyTouch(sim) {
		deathKnight.IcyTouch.Cast(sim, target)
		recastedFF = true
	} else if (!deathKnight.TargetHasDisease(BloodPlagueAuraLabel, target) || deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim) < 6*time.Second) && deathKnight.CanPlagueStrike(sim) {
		deathKnight.PlagueStrike.Cast(sim, target)
		recastedBP = true
	} else {
		if deathKnight.CanBloodTap(sim) && deathKnight.AllDiseasesAreActive(target) {
			deathKnight.BloodTap.Cast(sim, target)
		} else if deathKnight.CanUnbreakableArmor(sim) && deathKnight.AllDiseasesAreActive(target) {
			deathKnight.UnbreakableArmor.Cast(sim, target)
		} else if deathKnight.CanPestilence(sim) && deathKnight.shouldSpreadDisease(sim) {
			deathKnight.spreadDiseases(sim, target)
		} else if deathKnight.CanObliterate(sim) && deathKnight.AllDiseasesAreActive(target) {
			deathKnight.Obliterate.Cast(sim, target)
		} else if deathKnight.CanHowlingBlast(sim) && deathKnight.AllDiseasesAreActive(target) {
			deathKnight.HowlingBlast.Cast(sim, target)
		} else if deathKnight.CanFrostStrike(sim) && deathKnight.AllDiseasesAreActive(target) {
			deathKnight.FrostStrike.Cast(sim, target)
		} else if deathKnight.CanBloodStrike(sim) && deathKnight.AllDiseasesAreActive(target) {
			deathKnight.BloodStrike.Cast(sim, target)
		} else if deathKnight.CanIcyTouch(sim) {
			deathKnight.IcyTouch.Cast(sim, target)
		} else if deathKnight.CanPlagueStrike(sim) {
			deathKnight.PlagueStrike.Cast(sim, target)
		} else if deathKnight.CanHornOfWinter(sim) {
			deathKnight.HornOfWinter.Cast(sim, target)
		} else {
			if deathKnight.GCD.IsReady(sim) && !deathKnight.IsWaiting() {
				// This means we did absolutely nothing.
				// Wait until our next auto attack to decide again.
				nextSwing := deathKnight.AutoAttacks.MainhandSwingAt
				if deathKnight.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
					nextSwing = core.MinDuration(nextSwing, deathKnight.AutoAttacks.OffhandSwingAt)
				}
				deathKnight.WaitUntil(sim, nextSwing)
			}
		}
	}
}

func (deathKnight *DeathKnight) resetFrostRotation(sim *core.Simulation) {
	deathKnight.FrostRotation.actionIdx = 0
}
