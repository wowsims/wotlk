package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) setupUnholyRotations() {
	if dk.Rotation.BloodTap == proto.Deathknight_Rotation_GhoulFrenzy && !dk.Talents.GhoulFrenzy {
		dk.Rotation.BloodTap = proto.Deathknight_Rotation_IcyTouch
	}

	dk.setupGargoyleCooldowns()

	dk.RotationSequence.Clear().
		NewAction(dk.getFirstDiseaseAction()).
		NewAction(dk.getSecondDiseaseAction()).
		NewAction(dk.getBloodRuneAction(true))

	if dk.Rotation.UseDeathAndDecay || !dk.Talents.ScourgeStrike {
		if dk.Rotation.DeathAndDecayPrio == proto.Deathknight_Rotation_MaxRuneDowntime {
			dk.RotationSequence.
				NewAction(dk.RotationActionCallback_DND).
				NewAction(dk.RotationActionCallback_UnholyDndRotation)
		} else {
			dk.dndExperimentalOpener()
		}
	} else {
		dk.RotationSequence.NewAction(dk.RotationActionCallback_UnholySsRotation)
	}
}

func (dk *DpsDeathknight) RotationActionCallback_UnholyDndRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false

	if dk.uhGargoyleCheck(sim, target, 100*time.Millisecond) {
		dk.uhAfterGargoyleSequence(sim)
		return sim.CurrentTime
	}

	if dk.Talents.GhoulFrenzy && !dk.uhShouldWaitForDnD(sim, false, true, true) {
		if dk.uhGhoulFrenzyCheck(sim, target) {
			return sim.CurrentTime
		}
	}

	if dk.uhEmpoweredRuneWeapon(sim, target) {
		return sim.CurrentTime
	}

	if dk.uhBloodTap(sim, target) {
		return sim.CurrentTime
	}

	// What follows is a simple APL where every cast is checked against current diseses
	// And if the cast would leave the DK with not enough runes to cast disease before falloff
	// the cast is canceled and a disease recast is queued. Priority is as follows:
	// Death and Decay -> Scourge Strike -> Blood Strike (or Pesti/BB on Aoe) -> Death Coil -> Horn of Winter
	if !casted {
		if dk.uhDiseaseCheck(sim, target, dk.DeathAndDecay, true, 1) {
			if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
				dk.uhAfterGargoyleSequence(sim)
				return sim.CurrentTime
			}
			casted = dk.DeathAndDecay.Cast(sim, target)
		} else {
			if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()*2+250*time.Millisecond) {
				dk.uhAfterGargoyleSequence(sim)
				return sim.CurrentTime
			}
			dk.uhRecastDiseasesSequence(sim)
			return sim.CurrentTime
		}
		if !casted {
			if dk.uhDiseaseCheck(sim, target, dk.ScourgeStrike, true, 1) {
				if !dk.uhShouldWaitForDnD(sim, false, true, true) {
					if dk.Talents.ScourgeStrike && dk.ScourgeStrike.IsReady(sim) {
						if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
							dk.uhAfterGargoyleSequence(sim)
							return sim.CurrentTime
						}
						casted = dk.ScourgeStrike.Cast(sim, target)
					} else if dk.IcyTouch.CanCast(sim) && dk.PlagueStrike.CanCast(sim) {
						if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()*2+50*time.Millisecond) {
							dk.uhAfterGargoyleSequence(sim)
							return sim.CurrentTime
						}
						dk.uhRecastDiseasesSequence(sim)
						return sim.CurrentTime
					}
				}
			} else {
				dk.uhRecastDiseasesSequence(sim)
				return sim.CurrentTime
			}
			if !casted {
				if dk.shShouldSpreadDisease(sim) {
					if !dk.uhShouldWaitForDnD(sim, true, false, false) {
						if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
							dk.uhAfterGargoyleSequence(sim)
							return sim.CurrentTime
						}
						casted = dk.uhSpreadDiseases(sim, target, s)
					}
				} else {
					if !dk.uhShouldWaitForDnD(sim, true, false, false) {
						if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
							dk.uhAfterGargoyleSequence(sim)
							return sim.CurrentTime
						}
						if dk.desolationAuraCheck(sim) {
							casted = dk.BloodStrike.Cast(sim, target)
						} else {
							casted = dk.BloodBoil.Cast(sim, target)
						}
					}
				}
				if !casted {
					if dk.uhDeathCoilCheck(sim) {
						if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
							dk.uhAfterGargoyleSequence(sim)
							return sim.CurrentTime
						}
						casted = dk.DeathCoil.Cast(sim, target)
					}
					if !casted {
						if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
							dk.uhAfterGargoyleSequence(sim)
							return sim.CurrentTime
						}

						if dk.HornOfWinter.CanCast(sim) {
							casted = dk.HornOfWinter.Cast(sim, target)
						}
					}
				}
			}
		}
	}

	// Gargoyle cast needs to be checked more often then default rotation on gcd/resource gain checks
	if dk.SummonGargoyle.IsReady(sim) && dk.GCD.IsReady(sim) {
		return sim.CurrentTime + 100*time.Millisecond
	}

	return -1
}

func (dk *DpsDeathknight) RotationActionCallback_UnholySsRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false

	if dk.uhGargoyleCheck(sim, target, 100*time.Millisecond) {
		dk.uhAfterGargoyleSequence(sim)
		return sim.CurrentTime
	}

	if dk.Talents.GhoulFrenzy {
		if dk.uhGhoulFrenzyCheck(sim, target) {
			return sim.CurrentTime
		}
	}

	if dk.uhEmpoweredRuneWeapon(sim, target) {
		return sim.CurrentTime
	}

	if dk.uhBloodTap(sim, target) {
		return sim.CurrentTime
	}

	// What follows is a simple APL where every cast is checked against current diseses
	// And if the cast would leave the DK with not enough runes to cast disease before falloff
	// the cast is canceled and a disease recast is queued. Priority is as follows:
	// Scourge Strike -> Blood Strike (or Pesti/BB on Aoe) -> Death Coil -> Horn of Winter
	if !casted {
		if dk.uhDiseaseCheck(sim, target, dk.ScourgeStrike, true, 1) {
			if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
				dk.uhAfterGargoyleSequence(sim)
				return sim.CurrentTime
			}
			casted = dk.ScourgeStrike.Cast(sim, target)
		} else {
			if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()*2+50*time.Millisecond) {
				dk.uhAfterGargoyleSequence(sim)
				return sim.CurrentTime
			}
			dk.uhRecastDiseasesSequence(sim)
			return sim.CurrentTime
		}
		if !casted {
			if dk.shShouldSpreadDisease(sim) {
				if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
					dk.uhAfterGargoyleSequence(sim)
					return sim.CurrentTime
				}
				casted = dk.uhSpreadDiseases(sim, target, s)
			} else {
				if dk.uhDiseaseCheck(sim, target, dk.BloodStrike, true, 1) {
					if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
						dk.uhAfterGargoyleSequence(sim)
						return sim.CurrentTime
					}
					if dk.desolationAuraCheck(sim) {
						casted = dk.BloodStrike.Cast(sim, target)
					} else {
						casted = dk.BloodBoil.Cast(sim, target)
					}
				} else {
					if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()*2+50*time.Millisecond) {
						dk.uhAfterGargoyleSequence(sim)
						return sim.CurrentTime
					}
					dk.uhRecastDiseasesSequence(sim)
					return sim.CurrentTime
				}
			}
			if !casted {
				if dk.uhDeathCoilCheck(sim) {
					if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
						dk.uhAfterGargoyleSequence(sim)
						return sim.CurrentTime
					}
					casted = dk.DeathCoil.Cast(sim, target)
				}
				if !casted {
					if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
						dk.uhAfterGargoyleSequence(sim)
						return sim.CurrentTime
					}

					if dk.HornOfWinter.CanCast(sim) {
						casted = dk.HornOfWinter.Cast(sim, target)
					}
				}
			}
		}
	}

	// Gargoyle cast needs to be checked more often then default rotation on gcd/resource gain checks
	if dk.SummonGargoyle.IsReady(sim) && dk.GCD.IsReady(sim) {
		return sim.CurrentTime + 100*time.Millisecond
	}

	return -1
}

func (dk *DpsDeathknight) uhAfterGargoyleSequence(sim *core.Simulation) {
	if dk.Rotation.UseEmpowerRuneWeapon && dk.EmpowerRuneWeapon.IsReady(sim) {
		dk.RotationSequence.Clear()

		if dk.BloodTapAura.IsActive() {
			dk.RotationSequence.NewAction(dk.RotationActionUH_CancelBT)
		}

		didErw := false
		if dk.Rotation.ArmyOfTheDead != proto.Deathknight_Rotation_DoNotUse && dk.ArmyOfTheDead.IsReady(sim) {
			// If not enough runes for aotd cast ERW
			if dk.CurrentBloodRunes() < 1 || dk.CurrentFrostRunes() < 1 || dk.CurrentUnholyRunes() < 1 {
				dk.RotationSequence.NewAction(dk.RotationActionCallback_ERW)
				didErw = true
			}
			dk.RotationSequence.NewAction(dk.RotationActionCallback_AOTD)
		} else {
			// If no runes soon cast ERW
			if dk.CurrentBloodRunes() < 1 && dk.CurrentFrostRunes() < 1 && dk.CurrentUnholyRunes() < 1 && dk.AnyRuneReadyAt(sim)-sim.CurrentTime > 2*time.Second {
				dk.RotationSequence.NewAction(dk.RotationActionCallback_ERW)
				didErw = true
			}
		}

		if !dk.PresenceMatches(deathknight.BloodPresence) {
			if didErw || dk.CurrentBloodRunes() > 0 {
				dk.RotationSequence.NewAction(dk.RotationActionCallback_BP)
			} else if !didErw && !dk.Rotation.BtGhoulFrenzy && dk.BloodTap.IsReady(sim) {
				dk.RotationSequence.
					NewAction(dk.RotationActionCallback_BT).
					NewAction(dk.RotationActionCallback_BP).
					NewAction(dk.RotationActionUH_CancelBT)
			}
		}

		if dk.Rotation.UseDeathAndDecay || !dk.Talents.ScourgeStrike {
			dk.RotationSequence.NewAction(dk.RotationActionUH_ResetToDndMain)
		} else {
			dk.RotationSequence.NewAction(dk.RotationActionUH_ResetToSsMain)
		}
	}
}

func (dk *DpsDeathknight) uhGhoulFrenzySequence(sim *core.Simulation, bloodTap bool) {
	if bloodTap {
		dk.RotationSequence.Clear().
			NewAction(dk.RotationActionCallback_BT).
			NewAction(dk.RotationActionCallback_GF).
			NewAction(dk.RotationActionUH_CancelBT)
	} else {
		if dk.ur.ffFirst {
			dk.RotationSequence.Clear().
				NewAction(dk.RotationActionUH_IT_SetSync).
				NewAction(dk.RotationActionCallback_GF)
		} else {
			dk.RotationSequence.Clear().
				NewAction(dk.RotationActionCallback_GF).
				NewAction(dk.RotationActionUH_IT_SetSync)
		}
	}

	if dk.Rotation.UseDeathAndDecay || !dk.Talents.ScourgeStrike {
		dk.RotationSequence.NewAction(dk.RotationActionUH_ResetToDndMain)
	} else {
		dk.RotationSequence.NewAction(dk.RotationActionUH_ResetToSsMain)
	}
}

func (dk *DpsDeathknight) uhRecastDiseasesSequence(sim *core.Simulation) {
	dk.RotationSequence.Clear()

	// If we have glyph of Disease and both dots active try to refresh with pesti
	didPesti := false
	if dk.ur.hasGod {
		if dk.FrostFeverDisease[dk.CurrentTarget.Index].IsActive() && dk.BloodPlagueDisease[dk.CurrentTarget.Index].IsActive() {
			didPesti = true
			dk.RotationSequence.NewAction(dk.RotationActionCallback_Pesti_Custom)
		}
	}

	// If we did not pesti queue normal dot refresh
	if !didPesti {
		if dk.ur.ffFirst {
			dk.RotationSequence.
				NewAction(dk.RotationActionUH_FF_ClipCheck).
				NewAction(dk.RotationActionUH_IT_Custom).
				NewAction(dk.RotationActionUH_BP_ClipCheck).
				NewAction(dk.RotationActionUH_PS_Custom)
		} else {
			dk.RotationSequence.
				NewAction(dk.RotationActionUH_BP_ClipCheck).
				NewAction(dk.RotationActionUH_PS_Custom).
				NewAction(dk.RotationActionUH_FF_ClipCheck).
				NewAction(dk.RotationActionUH_IT_Custom)
		}
	}

	if dk.Rotation.UseDeathAndDecay || !dk.Talents.ScourgeStrike {
		dk.RotationSequence.NewAction(dk.RotationActionUH_ResetToDndMain)
	} else {
		dk.RotationSequence.NewAction(dk.RotationActionUH_ResetToSsMain)
	}
}

func (dk *DpsDeathknight) RotationActionCallback_Pesti_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	// If we have both dots active try to refresh with pesti and move to normal rotation
	if dk.FrostFeverDisease[dk.CurrentTarget.Index].IsActive() && dk.BloodPlagueDisease[dk.CurrentTarget.Index].IsActive() {
		dk.Pestilence.Cast(sim, target)
		s.Advance()

		return -1
	} else {
		// If a disease has dropped do normal reapply
		dk.RotationSequence.Clear()

		if dk.ur.ffFirst {
			dk.RotationSequence.
				NewAction(dk.RotationActionUH_FF_ClipCheck).
				NewAction(dk.RotationActionUH_IT_Custom).
				NewAction(dk.RotationActionUH_BP_ClipCheck).
				NewAction(dk.RotationActionUH_PS_Custom)
		} else {
			dk.RotationSequence.
				NewAction(dk.RotationActionUH_BP_ClipCheck).
				NewAction(dk.RotationActionUH_PS_Custom).
				NewAction(dk.RotationActionUH_FF_ClipCheck).
				NewAction(dk.RotationActionUH_IT_Custom)
		}
		return sim.CurrentTime
	}
}

func (dk *DpsDeathknight) RotationActionUH_CancelBT(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	dk.BloodTapAura.Deactivate(sim)
	s.Advance()
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionUH_ResetToSsMain(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	dk.RotationSequence.Clear().
		NewAction(dk.RotationActionCallback_UnholySsRotation)
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionUH_ResetToDndMain(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	dk.RotationSequence.Clear().
		NewAction(dk.RotationActionCallback_UnholyDndRotation)
	return sim.CurrentTime
}

// Custom PS callback for tracking recasts for pestilence disease sync
func (dk *DpsDeathknight) RotationActionUH_PS_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
		dk.uhAfterGargoyleSequence(sim)
		return sim.CurrentTime
	}
	casted := dk.PlagueStrike.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)

	dk.sr.recastedBP = casted && advance
	s.ConditionalAdvance(casted && advance)
	return -1
}

// Custom IT callback for tracking recasts for pestilence disease sync
func (dk *DpsDeathknight) RotationActionUH_IT_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
		dk.uhAfterGargoyleSequence(sim)
		return sim.CurrentTime
	}
	casted := dk.IcyTouch.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)
	if casted && advance {
		dk.sr.recastedFF = true
		dk.ur.syncTimeFF = 0
	}
	s.ConditionalAdvance(casted && advance)
	return -1
}

// Custom IT callback for ghoul frenzy frost rune sync
func (dk *DpsDeathknight) RotationActionUH_IT_SetSync(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	ffRemaining := dk.FrostFeverDisease[target.Index].RemainingDuration(sim)
	dk.RotationActionCallback_IT(sim, target, s)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)
	if !dk.GCD.IsReady(sim) && advance {
		dk.ur.syncTimeFF = dk.FrostFeverDisease[target.Index].Duration - ffRemaining
	}

	return -1
}

func (dk *DpsDeathknight) RotationActionUH_FF_ClipCheck(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	dot := dk.FrostFeverDisease[target.Index]
	gracePeriod := dk.CurrentFrostRuneGrace(sim)
	return dk.RotationActionUH_DiseaseClipCheck(dot, gracePeriod, sim, target, s)
}

func (dk *DpsDeathknight) RotationActionUH_BP_ClipCheck(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	dot := dk.BloodPlagueDisease[target.Index]
	gracePeriod := dk.CurrentUnholyRuneGrace(sim)
	return dk.RotationActionUH_DiseaseClipCheck(dot, gracePeriod, sim, target, s)
}

// Check if we have enough rune grace period to delay the disease cast
// so we get more ticks without losing on rune cd
func (dk *DpsDeathknight) RotationActionUH_DiseaseClipCheck(dot *core.Dot, gracePeriod time.Duration, sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	// TODO: Play around with allowing rune cd to be wasted
	// for more disease ticks and see if its a worth option for the ui
	//runeCdWaste := 0 * time.Millisecond
	waitTime := sim.CurrentTime
	if dot.TickCount < dot.NumberOfTicks-1 {
		nextTickAt := dot.ExpiresAt() - dot.TickLength*time.Duration((dot.NumberOfTicks-1)-dot.TickCount)
		if nextTickAt > sim.CurrentTime && (nextTickAt < sim.CurrentTime+gracePeriod || nextTickAt < sim.CurrentTime+400*time.Millisecond) {
			// Delay disease for next tick
			dk.LastOutcome = core.OutcomeMiss

			if dk.uhGargoyleCheck(sim, target, nextTickAt-sim.CurrentTime+50*time.Millisecond) {
				dk.uhAfterGargoyleSequence(sim)
				return waitTime
			}

			waitTime = nextTickAt + 50*time.Millisecond
		} else {
			waitTime = sim.CurrentTime
		}
	} else {
		waitTime = sim.CurrentTime
	}

	s.Advance()
	return waitTime
}
