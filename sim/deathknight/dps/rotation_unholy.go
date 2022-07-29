package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) RotationActionCallback_UnholyPrioRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	ur := &dk.ur
	// I suggest adding the a wrapper around each spell you cast like this:
	// dk.YourWrapper(sim, target, dk.FrostStrike) that returns a bool for when you casted
	// since the waiting code relies on knowing if you actually casted

	diseaseRefreshDuration := time.Duration(dk.Rotation.DiseaseRefreshDuration) * time.Second
	// Horn of Winter if you're the DK to refresh it and its not precasted/active
	if dk.ShouldHornOfWinter(sim) {
		dk.HornOfWinter.Cast(sim, target)
		casted = true
	} else if (!dk.FrostFeverDisease[target.Index].IsActive() || dk.FrostFeverDisease[target.Index].RemainingDuration(sim) < diseaseRefreshDuration) && dk.CanIcyTouch(sim) {
		// Dont clip if theres half a second left to tick
		remainingDuration := dk.FrostFeverDisease[target.Index].RemainingDuration(sim)
		if remainingDuration < time.Millisecond*500 && remainingDuration > 0 {
			dk.WaitUntil(sim, sim.CurrentTime+remainingDuration+1)
		} else {
			dk.IcyTouch.Cast(sim, target)
			casted = true
			ur.recastedFF = true
		}
	} else if (!dk.BloodPlagueDisease[target.Index].IsActive() || dk.BloodPlagueDisease[target.Index].RemainingDuration(sim) < diseaseRefreshDuration) && dk.CanPlagueStrike(sim) {
		// Dont clip if theres half a second left to tick
		remainingDuration := dk.BloodPlagueDisease[target.Index].RemainingDuration(sim)
		if remainingDuration < time.Millisecond*500 && remainingDuration > 0 {
			dk.WaitUntil(sim, sim.CurrentTime+remainingDuration+1)
		} else {
			dk.PlagueStrike.Cast(sim, target)
			casted = true
			ur.recastedBP = true
		}
	} else {
		if dk.Talents.Desolation > 0 && !dk.DesolationAura.IsActive() && dk.CanBloodStrike(sim) && !dk.uhShouldWaitForDnD(sim, true, false, false) {
			// Desolation and Pestilence check
			if dk.uhShouldSpreadDisease(sim) {
				dk.uhSpreadDiseases(sim, target, s)
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
				} else if dk.CanGhoulFrenzy(sim) && (!dk.Ghoul.GhoulFrenzyAura.IsActive() || dk.Ghoul.GhoulFrenzyAura.RemainingDuration(sim) < 6*time.Second) && !dk.uhShouldWaitForDnD(sim, false, false, true) {
					dk.GhoulFrenzy.Cast(sim, target)
					casted = true
				} else if dk.CanScourgeStrike(sim) && !dk.uhShouldWaitForDnD(sim, false, true, true) {
					dk.ScourgeStrike.Cast(sim, target)
					casted = true
				} else if !dk.Talents.ScourgeStrike && dk.CanIcyTouch(sim) && !dk.uhShouldWaitForDnD(sim, false, true, false) {
					dk.IcyTouch.Cast(sim, target)
					casted = true
				} else if !dk.Talents.ScourgeStrike && dk.CanPlagueStrike(sim) && !dk.uhShouldWaitForDnD(sim, false, false, true) {
					dk.PlagueStrike.Cast(sim, target)
					casted = true
				} else if dk.CanBloodStrike(sim) && !dk.uhShouldWaitForDnD(sim, true, false, false) {
					if dk.uhShouldSpreadDisease(sim) {
						dk.uhSpreadDiseases(sim, target, s)
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
					if dk.uhShouldSpreadDisease(sim) {
						dk.uhSpreadDiseases(sim, target, s)
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
	if dk.ur.ffFirst {
		return dk.RotationActionCallback_IT
	}
	return dk.RotationActionCallback_PS
}

func (dk *DpsDeathknight) getSecondDiseaseAction() deathknight.RotationAction {
	if dk.ur.ffFirst {
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
	s.Advance()
	return true
}

func (dk *DpsDeathknight) RotationAction_ResetToMain(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.Main.Clear().
		NewAction(dk.RotationActionCallback_UnholySsRotation)

	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

func (dk *DpsDeathknight) RotationAction_PS_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := dk.RotationActionCallback_PS(sim, target, s)
	advance := dk.LastCastOutcome.Matches(core.OutcomeLanded)
	dk.ur.recastedBP = casted && advance
	return casted
}

func (dk *DpsDeathknight) RotationAction_IT_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := dk.RotationActionCallback_IT(sim, target, s)
	advance := dk.LastCastOutcome.Matches(core.OutcomeLanded)
	if casted && advance {
		dk.ur.recastedFF = true
		dk.ur.syncTimeFF = 0
	}
	return casted
}

func (dk *DpsDeathknight) RotationAction_IT_SetSync(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	ffRemaining := dk.FrostFeverDisease[target.Index].RemainingDuration(sim)
	casted := dk.RotationActionCallback_IT(sim, target, s)
	advance := dk.LastCastOutcome.Matches(core.OutcomeLanded)
	if casted && advance {
		dk.ur.syncTimeFF = dk.FrostFeverDisease[target.Index].Duration - ffRemaining
	}

	return casted
}

func (dk *DpsDeathknight) RotationAction_FF_ClipCheck(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dot := dk.FrostFeverDisease[target.Index]
	gracePeriod := dk.CurrentFrostRuneGrace(sim)
	return dk.RotationAction_DiseaseClipCheck(dot, gracePeriod, sim, target, s)
}

func (dk *DpsDeathknight) RotationAction_BP_ClipCheck(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dot := dk.BloodPlagueDisease[target.Index]
	gracePeriod := dk.CurrentUnholyRuneGrace(sim)
	return dk.RotationAction_DiseaseClipCheck(dot, gracePeriod, sim, target, s)
}

func (dk *DpsDeathknight) RotationAction_DiseaseClipCheck(dot *core.Dot, gracePeriod time.Duration, sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	//runeCdWaste := 0 * time.Millisecond
	if dot.TickCount < dot.NumberOfTicks-1 {
		nextTickAt := dot.ExpiresAt() - dot.TickLength*time.Duration((dot.NumberOfTicks-1)-dot.TickCount)
		if nextTickAt > sim.CurrentTime && (nextTickAt < sim.CurrentTime+gracePeriod || nextTickAt < sim.CurrentTime+400*time.Millisecond) {
			// Delay disease for next tick
			dk.LastCastOutcome = core.OutcomeMiss
			dk.WaitUntil(sim, nextTickAt+50*time.Millisecond)
		} else {
			dk.WaitUntil(sim, sim.CurrentTime)
		}
	} else {
		dk.WaitUntil(sim, sim.CurrentTime)
	}

	s.Advance()
	return true
}

func (dk *DpsDeathknight) ghoulFrenzySequence(sim *core.Simulation, bloodTap bool) {
	if bloodTap {
		dk.Main.Clear().
			NewAction(dk.RotationActionCallback_BT).
			NewAction(dk.RotationActionCallback_GF).
			NewAction(dk.RotationAction_CancelBT)
	} else {
		if dk.ffFirst {
			dk.Main.Clear().
				NewAction(dk.RotationAction_IT_SetSync).
				NewAction(dk.RotationActionCallback_GF)
		} else {
			dk.Main.Clear().
				NewAction(dk.RotationActionCallback_GF).
				NewAction(dk.RotationAction_IT_SetSync)
		}
	}
	dk.Main.NewAction(dk.RotationAction_ResetToMain)
	dk.WaitUntil(sim, sim.CurrentTime)
}

func (dk *DpsDeathknight) recastDiseasesSequence(sim *core.Simulation) {
	dk.Main.Clear()

	if dk.ur.ffFirst {
		dk.Main.
			NewAction(dk.RotationAction_FF_ClipCheck).
			NewAction(dk.RotationAction_IT_Custom).
			NewAction(dk.RotationAction_BP_ClipCheck).
			NewAction(dk.RotationAction_PS_Custom).
			NewAction(dk.RotationAction_ResetToMain)
	} else {
		dk.Main.
			NewAction(dk.RotationAction_BP_ClipCheck).
			NewAction(dk.RotationAction_PS_Custom).
			NewAction(dk.RotationAction_FF_ClipCheck).
			NewAction(dk.RotationAction_IT_Custom).
			NewAction(dk.RotationAction_ResetToMain)
	}

	dk.WaitUntil(sim, sim.CurrentTime)
}

func (dk *DpsDeathknight) RotationActionCallback_UnholySsRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	if dk.Talents.GhoulFrenzy {
		if dk.Rotation.BtGhoulFrenzy && !dk.Rotation.UseDeathAndDecay {
			// Use Ghoul Frenzy only with a Blood Tap and Blood rune.
			// That means 50% uptime on GF at maximum but more Scourge Strikes
			if dk.CanBloodTap(sim) && dk.GhoulFrenzy.IsReady(sim) && dk.AllBloodRunesSpent() {
				if dk.UnholyDiseaseCheckWrapper(sim, target, dk.GhoulFrenzy, true, 1) {
					dk.ghoulFrenzySequence(sim, true)
					return true
				} else {
					dk.recastDiseasesSequence(sim)
					return true
				}
			}
		} else {
			// Use Ghoul Frenzy with an Unholy Rune and sync the frost rune with Icy Touch
			// That means 100% uptime on GF at maximum but less Scourge Strikes
			if dk.CanGhoulFrenzy(sim) && dk.CanIcyTouch(sim) &&
				(!dk.GhoulFrenzyAura.IsActive() || dk.GhoulFrenzyAura.RemainingDuration(sim) < 10*time.Second) {
				if dk.UnholyDiseaseCheckWrapper(sim, target, dk.GhoulFrenzy, true, 5) && dk.UnholyDiseaseCheckWrapper(sim, target, dk.IcyTouch, true, 5) {
					dk.ghoulFrenzySequence(sim, false)
					return true
				} else {
					dk.recastDiseasesSequence(sim)
					return true
				}
			}
		}
	}

	if !casted {
		if dk.UnholyDiseaseCheckWrapper(sim, target, dk.ScourgeStrike, true, 1) {
			casted = dk.CastScourgeStrike(sim, target)
		} else {
			dk.recastDiseasesSequence(sim)
			return true
		}
		if !casted {
			if dk.uhShouldSpreadDisease(sim) {
				casted = dk.uhSpreadDiseases(sim, target, s)
			} else {
				if dk.UnholyDiseaseCheckWrapper(sim, target, dk.BloodStrike, true, 1) {
					casted = dk.CastBloodStrike(sim, target)
				} else {
					dk.recastDiseasesSequence(sim)
					return true
				}
			}
			if !casted {
				casted = dk.CastDeathCoil(sim, target)
				if !casted {
					casted = dk.CastHornOfWinter(sim, target)
				}
			}
		}
	}

	return casted
}
