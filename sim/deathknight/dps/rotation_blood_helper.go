package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

type BloodRotation struct {
	dk *DpsDeathknight

	drwSnapshot *core.SnapshotManager

	activatingDrw bool
}

func (br *BloodRotation) Reset(sim *core.Simulation) {
	br.activatingDrw = false
	br.drwSnapshot.ResetProcTrackers()
}

func (dk *DpsDeathknight) blDiseaseCheck(sim *core.Simulation, target *core.Unit, spell *deathknight.RuneSpell, costRunes bool, casts int) bool {
	return dk.shDiseaseCheck(sim, target, spell, costRunes, casts, 0)
}

func (dk *DpsDeathknight) blSpreadDiseases(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.blDiseaseCheck(sim, target, dk.Pestilence, true, 1) {
		casted := dk.Pestilence.Cast(sim, target)
		landed := dk.LastOutcome.Matches(core.OutcomeLanded)

		// Reset flags on succesfull cast
		dk.sr.recastedFF = !(casted && landed)
		dk.sr.recastedBP = !(casted && landed)
		return -1
	} else {
		dk.blRecastDiseasesSequence(sim)
		return sim.CurrentTime
	}
}

// Save up Runic Power for DRW - Allow casts above 100 RP when DRW is ready or above 85 (for death strike glyph) when not
func (dk *DpsDeathknight) blDeathCoilCheck(sim *core.Simulation) bool {
	canCastDrw := dk.DancingRuneWeapon != nil && dk.DancingRuneWeapon.IsReady(sim) || dk.DancingRuneWeapon.CD.TimeToReady(sim) < 5*time.Second
	currentRP := dk.CurrentRunicPower()
	return (!canCastDrw && currentRP >= 65) || (canCastDrw && dk.CurrentRunicPower() >= 90)
}

// Combined checks for casting gargoyle sequence & going back to blood presence after
func (dk *DpsDeathknight) blDrwCheck(sim *core.Simulation, target *core.Unit, castTime time.Duration) bool {
	if dk.blDrwCanCast(sim, castTime) {

		// Unholy Presence
		// if !dk.PresenceMatches(deathknight.UnholyPresence) {
		// 	if dk.CurrentUnholyRunes() == 0 {
		// 		if dk.BloodTap.IsReady(sim) {
		// 			dk.BloodTap.Cast(sim, dk.CurrentTarget)
		// 		} else {
		// 			return false
		// 		}
		// 	}
		// 	dk.UnholyPresence.Cast(sim, dk.CurrentTarget)
		// }

		dk.br.activatingDrw = true
		dk.br.drwSnapshot.ActivateMajorCooldowns(sim)
		dk.br.activatingDrw = false

		if dk.DancingRuneWeapon.Cast(sim, target) {
			dk.br.drwSnapshot.ResetProcTrackers()
			return true
		}
	}

	// Go back to Blood Presence after Drw
	// if !dk.DancingRuneWeapon.IsReady(sim) && dk.PresenceMatches(deathknight.UnholyPresence) {
	// 	return dk.BloodPresence.Cast(sim, target)
	// }

	return false
}

func (dk *DpsDeathknight) blDrwCanCast(sim *core.Simulation, castTime time.Duration) bool {
	if !dk.DancingRuneWeapon.IsReady(sim) {
		return false
	}
	if !dk.CastCostPossible(sim, 60.0, 0, 0, 0) {
		return false
	}
	// if dk.CurrentDeathRunes() < 2 {
	// 	return false
	// }
	// if !dk.PresenceMatches(deathknight.UnholyPresence) && (!dk.BloodTap.CanCast(sim) && dk.CurrentUnholyRunes() == 0) {
	// 	return false
	// }
	if !dk.br.drwSnapshot.CanSnapShot(sim, castTime) {
		return false
	}

	return true
}
