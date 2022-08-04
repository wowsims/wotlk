package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

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
		if dk.Env.GetNumTargets() > 1 {
			return dk.RotationActionCallback_Pesti
		} else {
			return dk.RotationActionCallback_BS
		}
	} else {
		return dk.RotationActionCallback_BS
	}
}

func (dk *DpsDeathknight) setupUnholyOpener() {
	dk.setupGargoyleCooldowns()

	dk.Opener.
		NewAction(dk.getFirstDiseaseAction()).
		NewAction(dk.getSecondDiseaseAction()).
		NewAction(dk.getBloodRuneAction(true))

	if dk.Rotation.UseDeathAndDecay || !dk.Talents.ScourgeStrike {
		if dk.Rotation.DeathAndDecayPrio == proto.Deathknight_Rotation_MaxRuneDowntime {
			dk.Main.Clear().NewAction(dk.RotationActionCallback_UnholyDndRotation)
		} else {
			dk.dndStartOpener()
		}
	} else {
		dk.Main.NewAction(dk.RotationActionCallback_UnholySsRotation)
	}
}

func (dk *DpsDeathknight) afterGargoyleOpener(sim *core.Simulation) {
	if dk.Rotation.UseEmpowerRuneWeapon && dk.EmpowerRuneWeapon.IsReady(sim) {
		dk.Main.Clear()

		if dk.BloodTapAura.IsActive() {
			dk.Main.NewAction(dk.RotationAction_CancelBT)
		}

		if dk.Rotation.ArmyOfTheDead != proto.Deathknight_Rotation_DoNotUse && dk.ArmyOfTheDead.IsReady(sim) {
			// If not enough runes for aotd cast ERW
			if dk.CurrentBloodRunes() < 1 || dk.CurrentFrostRunes() < 1 || dk.CurrentUnholyRunes() < 1 {
				dk.Main.NewAction(dk.RotationActionCallback_ERW)
			}
			dk.Main.NewAction(dk.RotationActionCallback_AOTD)
		} else {
			// If no runes cast ERW TODO: Figure out when to do it after
			if dk.CurrentBloodRunes() < 1 && dk.CurrentFrostRunes() < 1 && dk.CurrentUnholyRunes() < 1 {
				dk.Main.NewAction(dk.RotationActionCallback_ERW)
			}
		}

		dk.Main.NewAction(dk.RotationActionCallback_BP)

		if dk.Rotation.UseDeathAndDecay || !dk.Talents.ScourgeStrike {
			dk.Main.NewAction(dk.RotationAction_ResetToDndMain)
		} else {
			dk.Main.NewAction(dk.RotationAction_ResetToSsMain)
		}
	}
}

func (dk *DpsDeathknight) RotationActionCallback_UnholyDndRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	if dk.uhGargoyleCheck(sim, target) {
		dk.afterGargoyleOpener(sim)
		return true
	}

	if dk.Talents.GhoulFrenzy && !dk.uhShouldWaitForDnD(sim, false, true, true) {
		if dk.uhGhoulFrenzyCheck(sim, target) {
			return true
		}
	}

	// What follows is a simple APL where every cast is checked against current diseses
	// And if the cast would leave the DK with not enough runes to cast disease before falloff
	// the cast is canceled and a disease recast is queued. Priority is as follows:
	// Death and Decay -> Scourge Strike -> Blood Strike (or Pesti/BB on Aoe) -> Death Coil -> Horn of Winter
	if !casted {
		if dk.uhDiseaseCheck(sim, target, dk.DeathAndDecay, true, 1) {
			casted = dk.CastDeathAndDecay(sim, target)
		} else {
			dk.recastDiseasesSequence(sim)
			return true
		}
		if !casted {
			if dk.uhDiseaseCheck(sim, target, dk.ScourgeStrike, true, 1) {
				if !dk.uhShouldWaitForDnD(sim, false, true, true) {
					if dk.Talents.ScourgeStrike && dk.ScourgeStrike.IsReady(sim) {
						casted = dk.ScourgeStrike.Cast(sim, target)
					} else if dk.CanIcyTouch(sim) && dk.CanPlagueStrike(sim) {
						dk.recastDiseasesSequence(sim)
						return true
					}
				}
			} else {
				dk.recastDiseasesSequence(sim)
				return true
			}
			if !casted {
				if dk.uhShouldSpreadDisease(sim) {
					if !dk.uhShouldWaitForDnD(sim, true, false, false) {
						casted = dk.uhSpreadDiseases(sim, target, s)
					}
				} else {
					if !dk.uhShouldWaitForDnD(sim, true, false, false) {
						if dk.desolationAuraCheck(sim) {
							casted = dk.CastBloodStrike(sim, target)
						} else {
							casted = dk.CastBloodBoil(sim, target)
						}
					}
				}
				if !casted {
					if dk.uhDeathCoilCheck(sim) {
						casted = dk.CastDeathCoil(sim, target)
					}
					if !casted {
						casted = dk.CastHornOfWinter(sim, target)
					}
				}
			}
		}
	}

	// Gargoyle cast needs to be checked more often then default rotation on gcd/resource gain checks
	if dk.SummonGargoyle.IsReady(sim) && dk.GCD.IsReady(sim) {
		dk.WaitUntil(sim, sim.CurrentTime+100*time.Millisecond)
		return true
	}

	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_UnholySsRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	if dk.uhGargoyleCheck(sim, target) {
		dk.afterGargoyleOpener(sim)
		return true
	}

	if dk.Talents.GhoulFrenzy {
		if dk.uhGhoulFrenzyCheck(sim, target) {
			return true
		}
	}

	// What follows is a simple APL where every cast is checked against current diseses
	// And if the cast would leave the DK with not enough runes to cast disease before falloff
	// the cast is canceled and a disease recast is queued. Priority is as follows:
	// Scourge Strike -> Blood Strike (or Pesti/BB on Aoe) -> Death Coil -> Horn of Winter
	if !casted {
		if dk.uhDiseaseCheck(sim, target, dk.ScourgeStrike, true, 1) {
			casted = dk.ScourgeStrike.Cast(sim, target)
		} else {
			dk.recastDiseasesSequence(sim)
			return true
		}
		if !casted {
			if dk.uhShouldSpreadDisease(sim) {
				casted = dk.uhSpreadDiseases(sim, target, s)
			} else {
				if dk.uhDiseaseCheck(sim, target, dk.BloodStrike, true, 1) {
					if dk.desolationAuraCheck(sim) {
						casted = dk.CastBloodStrike(sim, target)
					} else {
						casted = dk.CastBloodBoil(sim, target)
					}
				} else {
					dk.recastDiseasesSequence(sim)
					return true
				}
			}
			if !casted {
				if dk.uhDeathCoilCheck(sim) {
					casted = dk.CastDeathCoil(sim, target)
				}
				if !casted {
					casted = dk.CastHornOfWinter(sim, target)
				}
			}
		}
	}

	// Gargoyle cast needs to be checked more often then default rotation on gcd/resource gain checks
	if dk.SummonGargoyle.IsReady(sim) && dk.GCD.IsReady(sim) {
		dk.WaitUntil(sim, sim.CurrentTime+100*time.Millisecond)
		return true
	}

	return casted
}

func (dk *DpsDeathknight) ghoulFrenzySequence(sim *core.Simulation, bloodTap bool) {
	if bloodTap {
		dk.Main.Clear().
			NewAction(dk.RotationActionCallback_BT).
			NewAction(dk.RotationActionCallback_GF).
			NewAction(dk.RotationAction_CancelBT)
	} else {
		if dk.ur.ffFirst {
			dk.Main.Clear().
				NewAction(dk.RotationAction_IT_SetSync).
				NewAction(dk.RotationActionCallback_GF)
		} else {
			dk.Main.Clear().
				NewAction(dk.RotationActionCallback_GF).
				NewAction(dk.RotationAction_IT_SetSync)
		}
	}

	if dk.Rotation.UseDeathAndDecay || !dk.Talents.ScourgeStrike {
		dk.Main.NewAction(dk.RotationAction_ResetToDndMain)
	} else {
		dk.Main.NewAction(dk.RotationAction_ResetToSsMain)
	}
	dk.WaitUntil(sim, sim.CurrentTime)
}

func (dk *DpsDeathknight) recastDiseasesSequence(sim *core.Simulation) {
	dk.Main.Clear()

	if dk.ur.ffFirst {
		dk.Main.
			NewAction(dk.RotationAction_FF_ClipCheck).
			NewAction(dk.RotationAction_IT_Custom).
			NewAction(dk.RotationAction_BP_ClipCheck).
			NewAction(dk.RotationAction_PS_Custom)
	} else {
		dk.Main.
			NewAction(dk.RotationAction_BP_ClipCheck).
			NewAction(dk.RotationAction_PS_Custom).
			NewAction(dk.RotationAction_FF_ClipCheck).
			NewAction(dk.RotationAction_IT_Custom)
	}

	if dk.Rotation.UseDeathAndDecay || !dk.Talents.ScourgeStrike {
		dk.Main.NewAction(dk.RotationAction_ResetToDndMain)
	} else {
		dk.Main.NewAction(dk.RotationAction_ResetToSsMain)
	}
	dk.WaitUntil(sim, sim.CurrentTime)
}

func (dk *DpsDeathknight) RotationAction_CancelBT(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.BloodTapAura.Deactivate(sim)
	dk.WaitUntil(sim, sim.CurrentTime)
	s.Advance()
	return true
}

func (dk *DpsDeathknight) RotationAction_ResetToSsMain(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.Main.Clear().
		NewAction(dk.RotationActionCallback_UnholySsRotation)

	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

func (dk *DpsDeathknight) RotationAction_ResetToDndMain(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.Main.Clear().
		NewAction(dk.RotationActionCallback_UnholyDndRotation)

	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

// Custom PS callback for tracking recasts for pestilence disease sync
func (dk *DpsDeathknight) RotationAction_PS_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := dk.RotationActionCallback_PS(sim, target, s)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)
	dk.ur.recastedBP = casted && advance
	return casted
}

// Custom IT callback for tracking recasts for pestilence disease sync
func (dk *DpsDeathknight) RotationAction_IT_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := dk.RotationActionCallback_IT(sim, target, s)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)
	if casted && advance {
		dk.ur.recastedFF = true
		dk.ur.syncTimeFF = 0
	}
	return casted
}

// Custom IT callback for ghoul frenzy frost rune sync
func (dk *DpsDeathknight) RotationAction_IT_SetSync(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	ffRemaining := dk.FrostFeverDisease[target.Index].RemainingDuration(sim)
	casted := dk.RotationActionCallback_IT(sim, target, s)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)
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

// Check if we have enough rune grace period to delay the disease cast
// so we get more ticks without losing on rune cd
func (dk *DpsDeathknight) RotationAction_DiseaseClipCheck(dot *core.Dot, gracePeriod time.Duration, sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	// TODO: Play around with allowing rune cd to be wasted
	// for more disease ticks and see if its a worth option for the ui
	//runeCdWaste := 0 * time.Millisecond
	if dot.TickCount < dot.NumberOfTicks-1 {
		nextTickAt := dot.ExpiresAt() - dot.TickLength*time.Duration((dot.NumberOfTicks-1)-dot.TickCount)
		if nextTickAt > sim.CurrentTime && (nextTickAt < sim.CurrentTime+gracePeriod || nextTickAt < sim.CurrentTime+400*time.Millisecond) {
			// Delay disease for next tick
			dk.LastOutcome = core.OutcomeMiss
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
