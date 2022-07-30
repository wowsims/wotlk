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

	if dk.Rotation.DeathAndDecayPrio == proto.Deathknight_Rotation_MaxRuneDowntime {
		dk.Main.NewAction(dk.RotationActionCallback_UnholyDndRotation)
	} else {
		dk.dndStartSequence()
	}
}

func (dk *DpsDeathknight) dndStartSequence() {
	dk.Main.Clear().NewAction(dk.getFirstDiseaseAction()).
		NewAction(dk.getSecondDiseaseAction()).
		NewAction(dk.getBloodRuneAction(true)).
		NewAction(dk.RotationAction_Dnd_Custom).
		NewAction(dk.RotationAction_UnholyDndRotationGhoulFrenzyCheck)
}

// Custom Dnd callback with delay
func (dk *DpsDeathknight) RotationAction_Dnd_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := dk.CastDeathAndDecay(sim, target)
	if !casted {
		if !dk.DeathAndDecay.CD.IsReady(sim) {
			dk.WaitUntil(sim, dk.DeathAndDecay.ReadyAt())
			return true
		}
	} else {
		s.Advance()
	}
	return casted
}

func (dk *DpsDeathknight) RotationAction_UnholyDndRotationGhoulFrenzyCheck(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.Main.Clear()

	if dk.Talents.GhoulFrenzy && (!dk.GhoulFrenzyAura.IsActive() || dk.GhoulFrenzyAura.RemainingDuration(sim) < time.Second*10) {
		if dk.ur.ffFirst {
			dk.Main.NewAction(dk.RotationActionCallback_IT).
				NewAction(dk.RotationActionCallback_GF)
		} else {
			dk.Main.NewAction(dk.RotationActionCallback_GF).
				NewAction(dk.RotationActionCallback_IT)
		}
	} else {
		if dk.Talents.ScourgeStrike {
			dk.Main.NewAction(dk.RotationActionCallback_SS)
		} else {
			dk.Main.NewAction(dk.getFirstDiseaseAction()).
				NewAction(dk.getSecondDiseaseAction())
		}
	}

	if dk.desolationAuraCheck(sim) {
		dk.Main.NewAction(dk.RotationActionCallback_BS)
	} else {
		dk.Main.NewAction(dk.RotationActionCallback_BB)
	}
	dk.Main.NewAction(dk.RotationAction_UnholyDndRotationEnd)

	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

func (dk *DpsDeathknight) RotationAction_UnholyDndRotationEnd(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.dndStartSequence()
	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

func (dk *DpsDeathknight) RotationActionCallback_UnholyDndRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	if dk.Talents.GhoulFrenzy && !dk.uhShouldWaitForDnD(sim, false, true, true) {
		// If no Ghoul Frenzy Aura or duration less then 10 seconds we try recasting
		if !dk.GhoulFrenzyAura.IsActive() || dk.GhoulFrenzyAura.RemainingDuration(sim) < 10*time.Second {
			if dk.CanBloodTap(sim) && dk.GhoulFrenzy.IsReady(sim) && dk.AllBloodRunesSpent() && dk.AllUnholySpent() && dk.SummonGargoyle.CD.TimeToReady(sim) > time.Second*60 {
				// Use Ghoul Frenzy with a Blood Tap and Blood rune if all blood runes are on CD and Garg wont come off cd in less then a minute.
				// The gargoyle check is there because you should BT -> UP -> Garg (Not in the sim yet)
				if dk.uhDiseaseCheck(sim, target, dk.GhoulFrenzy, true, 1) {
					dk.ghoulFrenzySequence(sim, true)
					return true
				} else {
					dk.recastDiseasesSequence(sim)
					return true
				}
			} else if !dk.Rotation.BtGhoulFrenzy && dk.CanGhoulFrenzy(sim) && dk.CanIcyTouch(sim) {
				// Use Ghoul Frenzy with an Unholy Rune and sync the frost rune with Icy Touch
				if dk.uhDiseaseCheck(sim, target, dk.GhoulFrenzy, true, 5) && dk.uhDiseaseCheck(sim, target, dk.IcyTouch, true, 5) {
					dk.ghoulFrenzySequence(sim, false)
					return true
				} else {
					dk.recastDiseasesSequence(sim)
					return true
				}
			}
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
					casted = dk.CastScourgeStrike(sim, target)
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
					casted = dk.CastDeathCoil(sim, target)
					if !casted {
						casted = dk.CastHornOfWinter(sim, target)
					}
				}
			}
		}
	}

	return casted
}

func (dk *DpsDeathknight) desolationAuraCheck(sim *core.Simulation) bool {
	return !dk.DesolationAura.IsActive() || dk.DesolationAura.RemainingDuration(sim) < 10*time.Second || dk.Env.GetNumTargets() == 1
}

func (dk *DpsDeathknight) RotationActionCallback_UnholySsRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	if dk.Talents.GhoulFrenzy {
		// If no Ghoul Frenzy Aura or duration less then 10 seconds we try recasting
		if !dk.GhoulFrenzyAura.IsActive() || dk.GhoulFrenzyAura.RemainingDuration(sim) < 10*time.Second {
			if dk.CanBloodTap(sim) && dk.GhoulFrenzy.IsReady(sim) && dk.AllBloodRunesSpent() && dk.AllUnholySpent() && dk.SummonGargoyle.CD.TimeToReady(sim) > time.Second*60 {
				// Use Ghoul Frenzy with a Blood Tap and Blood rune if all blood runes are on CD and Garg wont come off cd in less then a minute.
				// The gargoyle check is there because you should BT -> UP -> Garg (Not in the sim yet)
				if dk.uhDiseaseCheck(sim, target, dk.GhoulFrenzy, true, 1) {
					dk.ghoulFrenzySequence(sim, true)
					return true
				} else {
					dk.recastDiseasesSequence(sim)
					return true
				}
			} else if !dk.Rotation.BtGhoulFrenzy && dk.CanGhoulFrenzy(sim) && dk.CanIcyTouch(sim) {
				// Use Ghoul Frenzy with an Unholy Rune and sync the frost rune with Icy Touch
				if dk.uhDiseaseCheck(sim, target, dk.GhoulFrenzy, true, 5) && dk.uhDiseaseCheck(sim, target, dk.IcyTouch, true, 5) {
					dk.ghoulFrenzySequence(sim, false)
					return true
				} else {
					dk.recastDiseasesSequence(sim)
					return true
				}
			}
		}
	}

	// What follows is a simple APL where every cast is checked against current diseses
	// And if the cast would leave the DK with not enough runes to cast disease before falloff
	// the cast is canceled and a disease recast is queued. Priority is as follows:
	// Scourge Strike -> Blood Strike (or Pesti on Aoe) -> Death Coil -> Horn of Winter
	if !casted {
		if dk.uhDiseaseCheck(sim, target, dk.ScourgeStrike, true, 1) {
			casted = dk.CastScourgeStrike(sim, target)
		} else {
			dk.recastDiseasesSequence(sim)
			return true
		}
		if !casted {
			if dk.uhShouldSpreadDisease(sim) {
				casted = dk.uhSpreadDiseases(sim, target, s)
			} else {
				if dk.uhDiseaseCheck(sim, target, dk.BloodStrike, true, 1) {
					if dk.Env.GetNumTargets() > 1 && dk.DesolationAura.IsActive() && dk.DesolationAura.RemainingDuration(sim) > time.Second*10 {
						casted = dk.CastBloodBoil(sim, target)
					} else {
						casted = dk.CastBloodStrike(sim, target)
					}
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

	if dk.Rotation.UseDeathAndDecay {
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

	if dk.Rotation.UseDeathAndDecay {
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
	advance := dk.LastCastOutcome.Matches(core.OutcomeLanded)
	dk.ur.recastedBP = casted && advance
	return casted
}

// Custom IT callback for tracking recasts for pestilence disease sync
func (dk *DpsDeathknight) RotationAction_IT_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := dk.RotationActionCallback_IT(sim, target, s)
	advance := dk.LastCastOutcome.Matches(core.OutcomeLanded)
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
