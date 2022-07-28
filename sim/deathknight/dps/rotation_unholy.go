package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) RotationActionCallback_UnholyPrioRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false
	// I suggest adding the a wrapper around each spell you cast like this:
	// dk.YourWrapper(sim, target, dk.FrostStrike) that returns a bool for when you casted
	// since the waiting code relies on knowing if you actually casted

	diseaseRefreshDuration := time.Duration(dk.Rotation.DiseaseRefreshDuration) * time.Second
	// Horn of Winter if you're the DK to refresh it and its not precasted/active
	if dk.ShouldHornOfWinter(sim) {
		dk.HornOfWinter.Cast(sim, target)
		casted = true
	} else if (!dk.TargetHasDisease(deathknight.FrostFeverAuraLabel, target) || dk.FrostFeverDisease[target.Index].RemainingDuration(sim) < diseaseRefreshDuration) && dk.CanIcyTouch(sim) {
		// Dont clip if theres half a second left to tick
		remainingDuration := dk.FrostFeverDisease[target.Index].RemainingDuration(sim)
		if remainingDuration < time.Millisecond*500 && remainingDuration > 0 {
			dk.WaitUntil(sim, sim.CurrentTime+remainingDuration+1)
		} else {
			dk.IcyTouch.Cast(sim, target)
			casted = true
			recastedFF = true
		}
	} else if (!dk.TargetHasDisease(deathknight.BloodPlagueAuraLabel, target) || dk.BloodPlagueDisease[target.Index].RemainingDuration(sim) < diseaseRefreshDuration) && dk.CanPlagueStrike(sim) {
		// Dont clip if theres half a second left to tick
		remainingDuration := dk.BloodPlagueDisease[target.Index].RemainingDuration(sim)
		if remainingDuration < time.Millisecond*500 && remainingDuration > 0 {
			dk.WaitUntil(sim, sim.CurrentTime+remainingDuration+1)
		} else {
			dk.PlagueStrike.Cast(sim, target)
			casted = true
			recastedBP = true
		}
	} else {
		if dk.Talents.Desolation > 0 && !dk.DesolationAura.IsActive() && dk.CanBloodStrike(sim) && !dk.shouldWaitForDnD(sim, true, false, false) {
			// Desolation and Pestilence check
			if dk.shouldSpreadDisease(sim) {
				dk.spreadDiseases(sim, target)
				casted = true
			} else {
				dk.BloodStrike.Cast(sim, target)
				casted = true
			}
		} else {
			if dk.Rotation.UseDeathAndDecay {
				// Death and Decay Rotation
				if dk.CanDeathAndDecay(sim) && dk.AllDiseasesAreActive(target) {
					dk.DeathAndDecay.Cast(sim, target)
					casted = true
				} else if dk.CanGhoulFrenzy(sim) && (!dk.Ghoul.GhoulFrenzyAura.IsActive() || dk.Ghoul.GhoulFrenzyAura.RemainingDuration(sim) < 6*time.Second) && !dk.shouldWaitForDnD(sim, false, false, true) {
					dk.GhoulFrenzy.Cast(sim, target)
					casted = true
				} else if dk.CanScourgeStrike(sim) && !dk.shouldWaitForDnD(sim, false, true, true) {
					dk.ScourgeStrike.Cast(sim, target)
					casted = true
				} else if !dk.Talents.ScourgeStrike && dk.CanIcyTouch(sim) && !dk.shouldWaitForDnD(sim, false, true, false) {
					dk.IcyTouch.Cast(sim, target)
					casted = true
				} else if !dk.Talents.ScourgeStrike && dk.CanPlagueStrike(sim) && !dk.shouldWaitForDnD(sim, false, false, true) {
					dk.PlagueStrike.Cast(sim, target)
					casted = true
				} else if dk.CanBloodStrike(sim) && !dk.shouldWaitForDnD(sim, true, false, false) {
					if dk.shouldSpreadDisease(sim) {
						dk.spreadDiseases(sim, target)
						casted = true
					} else if dk.Env.GetNumTargets() > 2 {
						dk.BloodBoil.Cast(sim, target)
						casted = true
					} else {
						dk.BloodStrike.Cast(sim, target)
						casted = true
					}
				} else if dk.CanDeathCoil(sim) && !dk.SummonGargoyle.IsReady(sim) {
					dk.DeathCoil.Cast(sim, target)
					casted = true
				} else if dk.CanHornOfWinter(sim) {
					dk.HornOfWinter.Cast(sim, target)
					casted = true
				} else {
					// Probably want to make this just return casted as casted should be false in this case, the wait time will be handled after the return
					if dk.GCD.IsReady(sim) && !dk.IsWaiting() {
						// This means we did absolutely nothing.
						// Wait until our next auto attack to decide again.
						nextSwing := dk.AutoAttacks.MainhandSwingAt
						if dk.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
							nextSwing = core.MinDuration(nextSwing, dk.AutoAttacks.OffhandSwingAt)
						}
						dk.WaitUntil(sim, nextSwing)
					}
				}
			} else {
				// Scourge Strike Rotation
				if dk.CanGhoulFrenzy(sim) && (!dk.Ghoul.GhoulFrenzyAura.IsActive() || dk.Ghoul.GhoulFrenzyAura.RemainingDuration(sim) < 6*time.Second) {
					dk.GhoulFrenzy.Cast(sim, target)
					casted = true
				} else if dk.CanScourgeStrike(sim) {
					dk.ScourgeStrike.Cast(sim, target)
					casted = true
				} else if dk.CanBloodStrike(sim) {
					if dk.shouldSpreadDisease(sim) {
						dk.spreadDiseases(sim, target)
						casted = true
					} else if dk.Env.GetNumTargets() > 2 {
						dk.BloodBoil.Cast(sim, target)
						casted = true
					} else {
						dk.BloodStrike.Cast(sim, target)
						casted = true
					}
				} else if dk.CanDeathCoil(sim) && !dk.SummonGargoyle.IsReady(sim) {
					dk.DeathCoil.Cast(sim, target)
					casted = true
				} else if dk.CanHornOfWinter(sim) {
					dk.HornOfWinter.Cast(sim, target)
					casted = true
				} else {
					if dk.GCD.IsReady(sim) && !dk.IsWaiting() {
						// This means we did absolutely nothing.
						// Wait until our next auto attack to decide again.
						nextSwing := dk.AutoAttacks.MainhandSwingAt
						if dk.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
							nextSwing = core.MinDuration(nextSwing, dk.AutoAttacks.OffhandSwingAt)
						}
						dk.WaitUntil(sim, nextSwing)
					}
				}
			}
		}
	}

	return casted
}

func (dk *DpsDeathknight) getFirstDiseaseAction() deathknight.RotationAction {
	if dk.Inputs.FirstDisease == proto.Deathknight_Rotation_FrostFever {
		return dk.RotationActionCallback_IT
	}
	return dk.RotationActionCallback_PS
}

func (dk *DpsDeathknight) getSecondDiseaseAction() deathknight.RotationAction {
	if dk.Inputs.FirstDisease == proto.Deathknight_Rotation_FrostFever {
		return dk.RotationActionCallback_PS
	}
	return dk.RotationActionCallback_IT
}

func (dk *DpsDeathknight) setupUnholySsOpener() {
	dk.Opener.NewAction(dk.getFirstDiseaseAction())
	dk.Opener.NewAction(dk.getSecondDiseaseAction())
	dk.Opener.NewAction(dk.RotationActionCallback_BS)
	dk.Opener.NewAction(dk.RotationActionCallback_SS)
	dk.Opener.NewAction(dk.RotationActionCallback_BT)
	dk.Opener.NewAction(dk.RotationActionCallback_UP)
	dk.Opener.NewAction(dk.RotationActionCallback_Garg)
	dk.Opener.NewAction(dk.RotationActionCallback_ERW)
	dk.Opener.NewAction(dk.RotationActionCallback_BP)
	dk.Opener.NewAction(dk.RotationActionCallback_SS)
	dk.Opener.NewAction(dk.RotationActionCallback_SS)
	dk.Opener.NewAction(dk.RotationActionCallback_BS)

	dk.Main.NewAction(dk.RotationActionCallback_UnholySsRotation)
}

func (dk *DpsDeathknight) setupUnholySsArmyOpener() {
	dk.Opener.NewAction(dk.getFirstDiseaseAction())
	dk.Opener.NewAction(dk.getSecondDiseaseAction())
	dk.Opener.NewAction(dk.RotationActionCallback_BS)
	dk.Opener.NewAction(dk.RotationActionCallback_SS)
	dk.Opener.NewAction(dk.RotationActionCallback_BT)
	dk.Opener.NewAction(dk.RotationActionCallback_UP)
	dk.Opener.NewAction(dk.RotationActionCallback_Garg)
	dk.Opener.NewAction(dk.RotationActionCallback_ERW)
	dk.Opener.NewAction(dk.RotationActionCallback_AOTD)
	dk.Opener.NewAction(dk.RotationActionCallback_BP)
	dk.Opener.NewAction(dk.RotationActionCallback_SS)

	dk.Main.NewAction(dk.RotationActionCallback_UnholySsRotation)
}

func (dk *DpsDeathknight) setupUnholyDndOpener() {
	dk.Opener.NewAction(dk.getFirstDiseaseAction())
	dk.Opener.NewAction(dk.getSecondDiseaseAction())
	dk.Opener.NewAction(dk.RotationActionCallback_BS)
	dk.Opener.NewAction(dk.RotationActionCallback_DND)
	dk.Opener.NewAction(dk.RotationActionCallback_BT)
	dk.Opener.NewAction(dk.RotationActionCallback_UP)
	dk.Opener.NewAction(dk.RotationActionCallback_Garg)
	dk.Opener.NewAction(dk.RotationActionCallback_ERW)
	dk.Opener.NewAction(dk.RotationActionCallback_BP)
	dk.Opener.NewAction(dk.RotationActionCallback_SS)
	dk.Opener.NewAction(dk.RotationActionCallback_SS)
	dk.Opener.NewAction(dk.RotationActionCallback_BS)

	dk.Main.NewAction(dk.RotationActionCallback_UnholyPrioRotation)
}

var syncDisease = false
var recastedFF = false
var recastedBP = false

func (dk *DpsDeathknight) UnholyDiseaseCheckWrapper(sim *core.Simulation, target *core.Unit, spell *core.Spell, costRunes bool) bool {
	success := false

	ffFirst := dk.Inputs.FirstDisease == proto.Deathknight_Rotation_FrostFever

	dropTimeAllowed := time.Millisecond * -100
	ffRemaining := dk.FrostFeverDisease[target.Index].RemainingDuration(sim) + dropTimeAllowed
	bpRemaining := dk.BloodPlagueDisease[target.Index].RemainingDuration(sim) + dropTimeAllowed
	castGcd := core.MinDuration(core.GCDMin, dk.ApplyCastSpeed(spell.CurCast.GCD))
	gracePeriodFrost := dk.CurrentFrostRuneGrace(sim)
	gracePeriodUnholy := dk.CurrentUnholyRuneGrace(sim)

	if ffFirst {
		if !dk.TargetHasDisease(deathknight.FrostFeverAuraLabel, target) || ffRemaining < castGcd {
			// Refresh FF
			success = dk.CastIcyTouch(sim, target)
			recastedFF = success
		} else if syncDisease || !dk.TargetHasDisease(deathknight.BloodPlagueAuraLabel, target) || bpRemaining < castGcd {
			// Refresh BP
			if syncDisease {
				dk.LastCastOutcome = core.OutcomeMiss
				success = dk.castClipDisease(false, gracePeriodUnholy, sim, dk.CanPlagueStrike(sim), dk.PlagueStrike, dk.BloodPlagueDisease[target.Index], target)
			} else {
				success = dk.CastPlagueStrike(sim, target)
			}
			recastedBP = success && dk.LastCastOutcome.Matches(core.OutcomeHit|core.OutcomeCrit)
			syncDisease = !recastedBP
		}
	} else {
		if !dk.TargetHasDisease(deathknight.BloodPlagueAuraLabel, target) || bpRemaining < castGcd {
			// Refresh BP
			success = dk.CastPlagueStrike(sim, target)
			recastedBP = success
		} else if syncDisease || !dk.TargetHasDisease(deathknight.FrostFeverAuraLabel, target) || ffRemaining < castGcd {
			// Refresh FF
			if syncDisease {
				dk.LastCastOutcome = core.OutcomeMiss
				success = dk.castClipDisease(false, gracePeriodFrost, sim, dk.CanIcyTouch(sim), dk.IcyTouch, dk.FrostFeverDisease[target.Index], target)
			} else {
				success = dk.CastIcyTouch(sim, target)
			}
			recastedFF = success && dk.LastCastOutcome.Matches(core.OutcomeHit|core.OutcomeCrit)
			syncDisease = !recastedFF
		}
	}

	if !success && dk.CanCast(sim, spell) {
		ffExpiresAt := ffRemaining + sim.CurrentTime
		bpExpiresAt := bpRemaining + sim.CurrentTime

		crpb := dk.CopyRunicPowerBar()
		runeCostForSpell := dk.RuneAmountForSpell(spell)
		spellCost := crpb.DetermineOptimalCost(sim, runeCostForSpell.Blood, runeCostForSpell.Frost, runeCostForSpell.Unholy)

		crpb.Spend(sim, spell, spellCost)

		afterCastTime := sim.CurrentTime + castGcd
		currentFrostRunes := crpb.CurrentFrostRunes()
		currentUnholyRunes := crpb.CurrentUnholyRunes()
		nextFrostRuneAt := crpb.FrostRuneReadyAt(sim)
		nextUnholyRuneAt := crpb.UnholyRuneReadyAt(sim)

		if ffFirst {
			// Check FF
			if dk.checkForDiseaseRecast(ffExpiresAt, afterCastTime, spellCost.Frost, currentFrostRunes, nextFrostRuneAt) {
				success = dk.castClipDisease(true, gracePeriodFrost, sim, dk.CanIcyTouch(sim), dk.IcyTouch, dk.FrostFeverDisease[target.Index], target)
				recastedFF = success
				return success
			}

			// Check BP
			if dk.checkForDiseaseRecast(bpExpiresAt, afterCastTime, spellCost.Unholy, currentUnholyRunes, nextUnholyRuneAt) {
				success = dk.castClipDisease(false, gracePeriodUnholy, sim, dk.CanPlagueStrike(sim), dk.PlagueStrike, dk.BloodPlagueDisease[target.Index], target)
				recastedBP = success
				return success
			}
		} else {
			// Check BP
			if dk.checkForDiseaseRecast(bpExpiresAt, afterCastTime, spellCost.Unholy, currentUnholyRunes, nextUnholyRuneAt) {
				success = dk.castClipDisease(true, gracePeriodUnholy, sim, dk.CanPlagueStrike(sim), dk.PlagueStrike, dk.BloodPlagueDisease[target.Index], target)
				recastedBP = success
				return success
			}

			// Check FF
			if dk.checkForDiseaseRecast(ffExpiresAt, afterCastTime, spellCost.Frost, currentFrostRunes, nextFrostRuneAt) {
				success = dk.castClipDisease(false, gracePeriodFrost, sim, dk.CanIcyTouch(sim), dk.IcyTouch, dk.FrostFeverDisease[target.Index], target)
				recastedFF = success
				return success
			}
		}

		// We have runes left for disease after this cast
		spell.Cast(sim, target)
		success = true
	}

	return success
}

func (dk *DpsDeathknight) checkForDiseaseRecast(expiresAt time.Duration, afterCastTime time.Duration,
	spellCost int, currentRunes int32, nextRuneAt time.Duration) bool {
	if spellCost > 0 && currentRunes == 0 {
		if expiresAt < nextRuneAt {
			return true
		}
	} else if afterCastTime > expiresAt {
		return true
	}
	return false
}

func (dk *DpsDeathknight) castClipDisease(mainDisease bool, gracePeriod time.Duration, sim *core.Simulation, canCast bool, spell *core.Spell, dot *core.Dot, target *core.Unit) bool {
	if canCast {
		// Dont drop disease due to %dmg modifiers
		if dot.TickCount < dot.NumberOfTicks-1 {
			nextTickAt := dot.ExpiresAt() - dot.TickLength*time.Duration((dot.NumberOfTicks-1)-dot.TickCount)
			if nextTickAt > sim.CurrentTime && (nextTickAt < sim.CurrentTime+gracePeriod || nextTickAt < sim.CurrentTime+400*time.Millisecond) {
				// Delay disease for next tick
				dk.WaitUntil(sim, nextTickAt+50*time.Millisecond)
				return true
			}
		}

		spell.Cast(sim, target)
		success := dk.LastCastOutcome.Matches(core.OutcomeCrit | core.OutcomeHit)
		if mainDisease {
			syncDisease = success
		}
		return true
	}
	return false
}

func (dk *DpsDeathknight) shouldSpreadDisease(sim *core.Simulation) bool {
	return recastedFF && recastedBP && dk.Env.GetNumTargets() > 1
}

func (dk *DpsDeathknight) spreadDiseases(sim *core.Simulation, target *core.Unit) bool {
	casted := dk.UnholyDiseaseCheckWrapper(sim, target, dk.Pestilence, true)

	// Reset flags on succesfull cast
	recastedFF = !casted
	recastedBP = !casted
	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_UnholySsRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	if dk.ShouldHornOfWinter(sim) {
		casted = dk.CastHornOfWinter(sim, target)
	} else {
		casted = dk.UnholyDiseaseCheckWrapper(sim, target, dk.ScourgeStrike, true)
		if !casted {
			if dk.shouldSpreadDisease(sim) {
				casted = dk.spreadDiseases(sim, target)
			} else {
				casted = dk.UnholyDiseaseCheckWrapper(sim, target, dk.BloodStrike, true)
			}
			if !casted {
				casted = dk.UnholyDiseaseCheckWrapper(sim, target, dk.DeathCoil, false)
				if !casted {
					casted = dk.UnholyDiseaseCheckWrapper(sim, target, dk.HornOfWinter, false)
				}
			}
		}
	}

	return casted
}

func (dk *DpsDeathknight) shouldWaitForDnD(sim *core.Simulation, blood bool, frost bool, unholy bool) bool {
	return dk.Rotation.UseDeathAndDecay && !(dk.Talents.Morbidity == 0 || !(dk.DeathAndDecay.CD.IsReady(sim) || dk.DeathAndDecay.CD.TimeToReady(sim) < 4*time.Second) || ((!blood || dk.CurrentBloodRunes() > 1) && (!frost || dk.CurrentFrostRunes() > 1) && (!unholy || dk.CurrentUnholyRunes() > 1)))
}
