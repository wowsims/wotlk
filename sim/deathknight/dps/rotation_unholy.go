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
			dk.recastedFF = true
		}
	} else if (!dk.TargetHasDisease(deathknight.BloodPlagueAuraLabel, target) || dk.BloodPlagueDisease[target.Index].RemainingDuration(sim) < diseaseRefreshDuration) && dk.CanPlagueStrike(sim) {
		// Dont clip if theres half a second left to tick
		remainingDuration := dk.BloodPlagueDisease[target.Index].RemainingDuration(sim)
		if remainingDuration < time.Millisecond*500 && remainingDuration > 0 {
			dk.WaitUntil(sim, sim.CurrentTime+remainingDuration+1)
		} else {
			dk.PlagueStrike.Cast(sim, target)
			casted = true
			dk.recastedBP = true
		}
	} else {
		if dk.Talents.Desolation > 0 && !dk.DesolationAura.IsActive() && dk.CanBloodStrike(sim) && !dk.shouldWaitForDnD(sim, true, false, false) {
			// Desolation and Pestilence check
			if dk.shouldSpreadDisease(sim) {
				dk.spreadDiseases(sim, target, s)
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
						dk.spreadDiseases(sim, target, s)
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
						dk.spreadDiseases(sim, target, s)
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

func (dk *DpsDeathknight) getBloodRuneAction(isFirst bool) deathknight.RotationAction {
	if isFirst {
		if dk.Env.GetNumTargets() > 1 && dk.Env.Encounter.Duration <= time.Second*30 {
			return dk.RotationActionCallback_Pesti
		} else {
			return dk.RotationActionCallback_BS
		}
	} else {
		if dk.Env.GetNumTargets() > 1 && dk.Env.Encounter.Duration > time.Second*30 {
			return dk.RotationActionCallback_Pesti
		} else {
			return dk.RotationActionCallback_BS
		}
	}
}

func (dk *DpsDeathknight) setupUnholySsOpener() {
	dk.Opener.
		NewAction(dk.getFirstDiseaseAction()).
		NewAction(dk.getSecondDiseaseAction()).
		NewAction(dk.getBloodRuneAction(true)).
		NewAction(dk.RotationActionCallback_SS).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_UP).
		NewAction(dk.RotationActionCallback_Garg).
		NewAction(dk.RotationAction_CancelBT).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_BP).
		NewAction(dk.RotationActionCallback_SS).
		NewAction(dk.RotationActionCallback_SS).
		NewAction(dk.getBloodRuneAction(false))

	dk.Main.NewAction(dk.RotationActionCallback_UnholySsRotation)
}

func (dk *DpsDeathknight) setupUnholySsArmyOpener() {
	dk.Opener.
		NewAction(dk.getFirstDiseaseAction()).
		NewAction(dk.getSecondDiseaseAction()).
		NewAction(dk.getBloodRuneAction(true)).
		NewAction(dk.RotationActionCallback_SS).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_UP).
		NewAction(dk.RotationActionCallback_Garg).
		NewAction(dk.RotationAction_CancelBT).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_AOTD).
		NewAction(dk.RotationActionCallback_BP).
		NewAction(dk.RotationActionCallback_SS)

	dk.Main.NewAction(dk.RotationActionCallback_UnholySsRotation)
}

func (dk *DpsDeathknight) setupUnholyDndOpener() {
	dk.Opener.
		NewAction(dk.getFirstDiseaseAction()).
		NewAction(dk.getSecondDiseaseAction()).
		NewAction(dk.getBloodRuneAction(true)).
		NewAction(dk.RotationActionCallback_DND).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_UP).
		NewAction(dk.RotationActionCallback_Garg).
		NewAction(dk.RotationAction_CancelBT).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_BP).
		NewAction(dk.RotationActionCallback_SS).
		NewAction(dk.RotationActionCallback_SS).
		NewAction(dk.getBloodRuneAction(false))

	dk.Main.NewAction(dk.RotationActionCallback_UnholyPrioRotation)
}

func (dk *DpsDeathknight) RotationAction_CancelBT(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.BloodTapAura.Deactivate(sim)
	dk.WaitUntil(sim, sim.CurrentTime)
	s.ConditionalAdvance(true)
	return true
}

func (dk *DpsDeathknight) RotationAction_ResetToMain(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.Main.ResetSequence().
		NewAction(dk.RotationActionCallback_UnholySsRotation)

	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

func (dk *DpsDeathknight) RotationAction_IT_SetSync(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	ffRemaining := dk.FrostFeverDisease[target.Index].RemainingDuration(sim)
	casted := dk.CastIcyTouch(sim, target)
	advance := dk.LastCastOutcome.Matches(core.OutcomeLanded)
	if casted && advance {
		dk.syncTimeFF = dk.FrostFeverDisease[target.Index].Duration - ffRemaining
	}

	s.ConditionalAdvance(casted && advance)
	return casted
}

func (dk *DpsDeathknight) ghoulFrenzySequence() {
	ffFirst := dk.Inputs.FirstDisease == proto.Deathknight_Rotation_FrostFever
	if ffFirst {
		dk.Main.ResetSequence().
			NewAction(dk.RotationAction_IT_SetSync).
			NewAction(dk.RotationActionCallback_GF).
			NewAction(dk.RotationAction_ResetToMain)
	} else {
		dk.Main.ResetSequence().
			NewAction(dk.RotationActionCallback_GF).
			NewAction(dk.RotationAction_IT_SetSync).
			NewAction(dk.RotationAction_ResetToMain)
	}
}

func (dk *DpsDeathknight) RotationActionCallback_UnholySsRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	if dk.ShouldHornOfWinter(sim) {
		casted = dk.CastHornOfWinter(sim, target)
	} else {
		// Ghoul Frenzy usage with FF sync
		if dk.Talents.GhoulFrenzy && (!dk.GhoulFrenzyAura.IsActive() || dk.GhoulFrenzyAura.RemainingDuration(sim) < 10*time.Second) {
			ffRemaining := dk.FrostFeverDisease[target.Index].RemainingDuration(sim)
			bpRemaining := dk.BloodPlagueDisease[target.Index].RemainingDuration(sim)
			minRemaining := core.MinDuration(ffRemaining, bpRemaining)

			if minRemaining > dk.SpellGCD()*2 {
				casted = dk.UnholyDiseaseCheckWrapper(sim, target, dk.GhoulFrenzy, true)
				if casted && dk.lastCastSpell == dk.GhoulFrenzy {
					dk.ghoulFrenzySequence()
					dk.WaitUntil(sim, sim.CurrentTime)
				}
			}
		}
		if !casted {
			casted = dk.UnholyDiseaseCheckWrapper(sim, target, dk.ScourgeStrike, true)
			if !casted {
				if dk.shouldSpreadDisease(sim) {
					casted = dk.spreadDiseases(sim, target, s)
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
	}

	return casted
}
