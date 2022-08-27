package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) setupBloodRotations() {
	dk.RotationSequence.Clear().
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_HS).
		NewAction(dk.RotationActionCallback_DS).
		NewAction(dk.RotationActionCallback_HS).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_DRW).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_HS).
		NewAction(dk.RotationActionCallback_HS).
		NewAction(dk.RotationActionCallback_HS).
		NewAction(dk.RotationActionCallback_HS)

	dk.RotationSequence.NewAction(dk.RotationActionCallback_BloodRotation)
}

func (dk *DpsDeathknight) RotationActionCallback_BloodRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false

	if dk.DancingRuneWeapon.CanCast(sim) {
		casted = dk.DancingRuneWeapon.Cast(sim, target)
	}

	if !casted {
		if dk.blDiseaseCheck(sim, target, dk.HeartStrike, true, 1) {
			if dk.shShouldSpreadDisease(sim) {
				return dk.blSpreadDiseases(sim, target, s)
			} else {
				if dk.Talents.HeartStrike {
					casted = dk.HeartStrike.Cast(sim, target)
				} else {
					casted = dk.BloodStrike.Cast(sim, target)
				}
			}
		} else {
			dk.blRecastDiseasesSequence(sim)
			return sim.CurrentTime
		}
		if !casted {
			if dk.blDiseaseCheck(sim, target, dk.DeathStrike, true, 1) {
				casted = dk.DeathStrike.Cast(sim, target)
			} else {
				dk.blRecastDiseasesSequence(sim)
				return sim.CurrentTime
			}
			if !casted {
				if dk.blDeathCoilCheck(sim) {
					casted = dk.DeathCoil.Cast(sim, target)
				}
				if !casted && dk.HornOfWinter.CanCast(sim) {
					dk.HornOfWinter.Cast(sim, target)
				}
			}
		}
	}

	return -1
}

func (dk *DpsDeathknight) blRecastDiseasesSequence(sim *core.Simulation) {
	dk.RotationSequence.Clear().
		NewAction(dk.RotationActionBL_FF_ClipCheck).
		NewAction(dk.RotationActionBL_IT_Custom).
		NewAction(dk.RotationActionBL_BP_ClipCheck).
		NewAction(dk.RotationActionBL_PS_Custom).
		NewAction(dk.RotationAction_ResetToBloodMain)
}

func (dk *DpsDeathknight) RotationAction_ResetToBloodMain(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	dk.RotationSequence.Clear().
		NewAction(dk.RotationActionCallback_BloodRotation)

	return sim.CurrentTime
}

// Custom PS callback for tracking recasts for pestilence disease sync
func (dk *DpsDeathknight) RotationActionBL_PS_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := dk.PlagueStrike.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)

	dk.sr.recastedBP = casted && advance
	s.ConditionalAdvance(casted && advance)
	return -1
}

// Custom IT callback for tracking recasts for pestilence disease sync
func (dk *DpsDeathknight) RotationActionBL_IT_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := dk.IcyTouch.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)
	dk.sr.recastedFF = casted && advance
	s.ConditionalAdvance(casted && advance)
	return -1
}

func (dk *DpsDeathknight) RotationActionBL_FF_ClipCheck(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	dot := dk.FrostFeverDisease[target.Index]
	gracePeriod := dk.CurrentFrostRuneGrace(sim)
	return dk.RotationActionBL_DiseaseClipCheck(dot, gracePeriod, sim, target, s)
}

func (dk *DpsDeathknight) RotationActionBL_BP_ClipCheck(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	dot := dk.BloodPlagueDisease[target.Index]
	gracePeriod := dk.CurrentUnholyRuneGrace(sim)
	return dk.RotationActionBL_DiseaseClipCheck(dot, gracePeriod, sim, target, s)
}

// Check if we have enough rune grace period to delay the disease cast
// so we get more ticks without losing on rune cd
func (dk *DpsDeathknight) RotationActionBL_DiseaseClipCheck(dot *core.Dot, gracePeriod time.Duration, sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	// TODO: Play around with allowing rune cd to be wasted
	// for more disease ticks and see if its a worth option for the ui
	//runeCdWaste := 0 * time.Millisecond
	waitUntil := time.Duration(-1)
	if dot.TickCount < dot.NumberOfTicks-1 {
		nextTickAt := dot.ExpiresAt() - dot.TickLength*time.Duration((dot.NumberOfTicks-1)-dot.TickCount)
		if nextTickAt > sim.CurrentTime && (nextTickAt < sim.CurrentTime+gracePeriod || nextTickAt < sim.CurrentTime+400*time.Millisecond) {
			// Delay disease for next tick
			dk.LastOutcome = core.OutcomeMiss
			waitUntil = nextTickAt + 50*time.Millisecond
		} else {
			waitUntil = sim.CurrentTime
		}
	} else {
		waitUntil = sim.CurrentTime
	}

	s.Advance()
	return waitUntil
}
