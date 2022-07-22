package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (deathKnight *DeathKnight) setupUnholyRotations() {

	// This defines the Sub Blood opener
	deathKnight.DefineOpener(RotationID_Unholy_Full, []RotationAction{})
}

func (deathKnight *DeathKnight) UnholyDiseaseCheckWrapper(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
	success := false

	if !deathKnight.TargetHasDisease(FrostFeverAuraLabel, target) {
		success = deathKnight.CastIcyTouch(sim, target)
	} else if !deathKnight.TargetHasDisease(BloodPlagueAuraLabel, target) {
		success = deathKnight.CastPlagueStrike(sim, target)
	} else if deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim) < spell.CurCast.GCD ||
		deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim) < spell.CurCast.GCD {
		success = deathKnight.CastPestilence(sim, target)
		if deathKnight.LastCastOutcome == core.OutcomeMiss {
			// Deal with pestilence miss
			// TODO:
			if deathKnight.opener.id == RotationID_FrostSubUnholy_Full {
				deathKnight.PushSequence([]RotationAction{
					RotationAction_BS,
					RotationAction_FS,
					RotationAction_IT,
					RotationAction_PS,
					RotationAction_Obli,
					RotationAction_Obli,
					RotationAction_FS,
					RotationAction_FS,
				})
			}
		}
	} else {
		if deathKnight.CanCast(sim, spell) {
			ffExpiresIn := deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim)
			bpExpiresIn := deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim)
			ffExpiresAt := ffExpiresIn + sim.CurrentTime
			bpExpiresAt := bpExpiresIn + sim.CurrentTime
			if spell.CurCast.GCD > ffExpiresIn || spell.CurCast.GCD > bpExpiresIn {
				return success
			}

			crpb := deathKnight.CopyRunicPowerBar()
			runeCostForSpell := deathKnight.RuneAmountForSpell(spell)
			spellCost := crpb.DetermineOptimalCost(sim, runeCostForSpell.Blood, runeCostForSpell.Frost, runeCostForSpell.Unholy)

			// Add whichever non-frost specific checks you want here, I guess you'll need them.

			if !(deathKnight.RimeAura.IsActive() && spell == deathKnight.HowlingBlast) {
				crpb.Spend(sim, spell, spellCost)
			}

			if crpb.CurrentBloodRunes() == 0 && crpb.CurrentDeathRunes() == 0 {
				nextBloodRuneAt := float64(crpb.BloodRuneReadyAt(sim))
				nextDeathRuneAt := float64(crpb.DeathRuneReadyAt(sim))

				ff1 := (float64(ffExpiresAt) > nextBloodRuneAt) && (float64(ffExpiresAt)-nextBloodRuneAt < float64(spell.CurCast.GCD))
				ff2 := (float64(ffExpiresAt) > nextDeathRuneAt) && (float64(ffExpiresAt)-nextDeathRuneAt < float64(spell.CurCast.GCD))
				bp1 := (float64(bpExpiresAt) > nextBloodRuneAt) && (float64(bpExpiresAt)-nextBloodRuneAt < float64(spell.CurCast.GCD))
				bp2 := (float64(bpExpiresAt) > nextDeathRuneAt) && (float64(bpExpiresAt)-nextDeathRuneAt < float64(spell.CurCast.GCD))

				if (ff1 || ff2) && (bp1 || bp2) {
					if deathKnight.CanCast(sim, spell) {
						spell.Cast(sim, target)
						success = true
					}
				} else {
					return success
				}
			} else {
				spell.Cast(sim, target)
				success = true
			}
		}
	}

	return success
}

func (deathKnight *DeathKnight) doUnholyRotation(sim *core.Simulation, target *core.Unit) bool {
	casted := &deathKnight.castSuccessful
	// I suggest adding the a wrapper around each spell you cast like this:
	// deathKnight.YourWrapper(sim, target, deathKnight.FrostStrike) that returns a bool for when you casted
	// since the waiting code relies on knowing if you actually casted

	if deathKnight.CanRaiseDead(sim) {
		deathKnight.RaiseDead.Cast(sim, target)
		*casted = true
		return *casted
	}
	diseaseRefreshDuration := time.Duration(deathKnight.Rotation.DiseaseRefreshDuration) * time.Second
	// Horn of Winter if you're the DK to refresh it and its not precasted/active
	if deathKnight.ShouldHornOfWinter(sim) {
		deathKnight.HornOfWinter.Cast(sim, target)
		*casted = true
	} else if (!deathKnight.TargetHasDisease(FrostFeverAuraLabel, target) || deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim) < diseaseRefreshDuration) && deathKnight.CanIcyTouch(sim) {
		// Dont clip if theres half a second left to tick
		remainingDuration := deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim)
		if remainingDuration < time.Millisecond*500 && remainingDuration > 0 {
			deathKnight.WaitUntil(sim, sim.CurrentTime+remainingDuration+1)
		} else {
			deathKnight.IcyTouch.Cast(sim, target)
			*casted = true
			recastedFF = true
		}
	} else if (!deathKnight.TargetHasDisease(BloodPlagueAuraLabel, target) || deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim) < diseaseRefreshDuration) && deathKnight.CanPlagueStrike(sim) {
		// Dont clip if theres half a second left to tick
		remainingDuration := deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim)
		if remainingDuration < time.Millisecond*500 && remainingDuration > 0 {
			deathKnight.WaitUntil(sim, sim.CurrentTime+remainingDuration+1)
		} else {
			deathKnight.PlagueStrike.Cast(sim, target)
			*casted = true
			recastedBP = true
		}
	} else {
		if deathKnight.PresenceMatches(UnholyPresence) && (deathKnight.Rotation.ArmyOfTheDead != proto.DeathKnight_Rotation_AsMajorCd || !deathKnight.ArmyOfTheDead.CD.IsReady(sim)) && !deathKnight.SummonGargoyle.CD.IsReady(sim) && deathKnight.CanBloodPresence(sim) {
			// Swap to blood presence after gargoyle cast
			deathKnight.BloodPressence.Cast(sim, target)
			deathKnight.WaitUntil(sim, sim.CurrentTime+1)
		} else if deathKnight.Talents.Desolation > 0 && !deathKnight.DesolationAura.IsActive() && deathKnight.CanBloodStrike(sim) && !deathKnight.shouldWaitForDnD(sim, true, false, false) {
			// Desolation and Pestilence check
			if deathKnight.shouldSpreadDisease(sim) {
				deathKnight.spreadDiseases(sim, target)
				*casted = true
			} else {
				deathKnight.BloodStrike.Cast(sim, target)
				*casted = true
			}
		} else {
			if deathKnight.Rotation.UseDeathAndDecay {
				// Death and Decay Rotation
				if deathKnight.CanDeathAndDecay(sim) && deathKnight.AllDiseasesAreActive(target) {
					deathKnight.DeathAndDecay.Cast(sim, target)
					*casted = true
				} else if deathKnight.CanGhoulFrenzy(sim) && deathKnight.Talents.MasterOfGhouls && (!deathKnight.Ghoul.GhoulFrenzyAura.IsActive() || deathKnight.Ghoul.GhoulFrenzyAura.RemainingDuration(sim) < 6*time.Second) && !deathKnight.shouldWaitForDnD(sim, false, false, true) {
					deathKnight.GhoulFrenzy.Cast(sim, target)
					*casted = true
				} else if deathKnight.CanScourgeStrike(sim) && !deathKnight.shouldWaitForDnD(sim, false, true, true) {
					deathKnight.ScourgeStrike.Cast(sim, target)
					*casted = true
				} else if !deathKnight.Talents.ScourgeStrike && deathKnight.CanIcyTouch(sim) && !deathKnight.shouldWaitForDnD(sim, false, true, false) {
					deathKnight.IcyTouch.Cast(sim, target)
					*casted = true
				} else if !deathKnight.Talents.ScourgeStrike && deathKnight.CanPlagueStrike(sim) && !deathKnight.shouldWaitForDnD(sim, false, false, true) {
					deathKnight.PlagueStrike.Cast(sim, target)
					*casted = true
				} else if deathKnight.CanBloodStrike(sim) && !deathKnight.shouldWaitForDnD(sim, true, false, false) {
					if deathKnight.shouldSpreadDisease(sim) {
						deathKnight.spreadDiseases(sim, target)
						*casted = true
					} else if deathKnight.Env.GetNumTargets() > 2 {
						deathKnight.BloodBoil.Cast(sim, target)
						*casted = true
					} else {
						deathKnight.BloodStrike.Cast(sim, target)
						*casted = true
					}
				} else if deathKnight.CanDeathCoil(sim) && !deathKnight.SummonGargoyle.IsReady(sim) {
					deathKnight.DeathCoil.Cast(sim, target)
					*casted = true
				} else if deathKnight.CanHornOfWinter(sim) {
					deathKnight.HornOfWinter.Cast(sim, target)
					*casted = true
				} else {
					// Probably want to make this just return *casted as casted should be false in this case, the wait time will be handled after the return
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
			} else {
				// Scourge Strike Rotation
				if deathKnight.CanGhoulFrenzy(sim) && deathKnight.Talents.MasterOfGhouls && (!deathKnight.Ghoul.GhoulFrenzyAura.IsActive() || deathKnight.Ghoul.GhoulFrenzyAura.RemainingDuration(sim) < 6*time.Second) {
					deathKnight.GhoulFrenzy.Cast(sim, target)
					*casted = true
				} else if deathKnight.CanScourgeStrike(sim) {
					deathKnight.ScourgeStrike.Cast(sim, target)
					*casted = true
				} else if deathKnight.CanBloodStrike(sim) {
					if deathKnight.shouldSpreadDisease(sim) {
						deathKnight.spreadDiseases(sim, target)
						*casted = true
					} else if deathKnight.Env.GetNumTargets() > 2 {
						deathKnight.BloodBoil.Cast(sim, target)
						*casted = true
					} else {
						deathKnight.BloodStrike.Cast(sim, target)
						*casted = true
					}
				} else if deathKnight.CanDeathCoil(sim) && !deathKnight.SummonGargoyle.IsReady(sim) {
					deathKnight.DeathCoil.Cast(sim, target)
					*casted = true
				} else if deathKnight.CanHornOfWinter(sim) {
					deathKnight.HornOfWinter.Cast(sim, target)
					*casted = true
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
	}
	return true
}
