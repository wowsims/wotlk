package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (deathKnight *DpsDeathknight) setupUnholyRotations() {

	deathKnight.DefineOpener(deathknight.RotationID_UnholySsUnholyPresence_Full, []deathknight.RotationAction{
		deathknight.RotationAction_IT,
		deathknight.RotationAction_PS,
		deathknight.RotationAction_BS,
		deathknight.RotationAction_SS,
		deathknight.RotationAction_Garg,
		deathknight.RotationAction_ERW,
		deathknight.RotationAction_BP,
		deathknight.RotationAction_SS,
		deathknight.RotationAction_SS,
		deathknight.RotationAction_BS,
	})

	deathKnight.DefineOpener(deathknight.RotationID_UnholySsArmyUnholyPresence_Full, []deathknight.RotationAction{
		deathknight.RotationAction_IT,
		deathknight.RotationAction_PS,
		deathknight.RotationAction_BS,
		deathknight.RotationAction_SS,
		deathknight.RotationAction_Garg,
		deathknight.RotationAction_ERW,
		deathknight.RotationAction_AOTD,
		deathknight.RotationAction_BP,
		deathknight.RotationAction_SS,
	})

	deathKnight.DefineOpener(deathknight.RotationID_UnholySsBloodPresence_Full, []deathknight.RotationAction{
		deathknight.RotationAction_IT,
		deathknight.RotationAction_PS,
		deathknight.RotationAction_BS,
		deathknight.RotationAction_SS,
		deathknight.RotationAction_BT,
		deathknight.RotationAction_UP,
		deathknight.RotationAction_Garg,
		deathknight.RotationAction_ERW,
		deathknight.RotationAction_BP,
		deathknight.RotationAction_SS,
		deathknight.RotationAction_SS,
		deathknight.RotationAction_BS,
	})

	deathKnight.DefineOpener(deathknight.RotationID_UnholySsArmyBloodPresence_Full, []deathknight.RotationAction{
		deathknight.RotationAction_IT,
		deathknight.RotationAction_PS,
		deathknight.RotationAction_BS,
		deathknight.RotationAction_SS,
		deathknight.RotationAction_BT,
		deathknight.RotationAction_UP,
		deathknight.RotationAction_Garg,
		deathknight.RotationAction_ERW,
		deathknight.RotationAction_AOTD,
		deathknight.RotationAction_BP,
		deathknight.RotationAction_SS,
	})

	deathKnight.DefineOpener(deathknight.RotationID_UnholyDnd_Full, []deathknight.RotationAction{
		deathknight.RotationAction_IT,
		deathknight.RotationAction_PS,
		deathknight.RotationAction_BS,
		deathknight.RotationAction_DND,
		deathknight.RotationAction_BT,
		deathknight.RotationAction_UP,
		deathknight.RotationAction_Garg,
		deathknight.RotationAction_ERW,
		deathknight.RotationAction_BP,
		deathknight.RotationAction_SS,
		deathknight.RotationAction_BS,
		deathknight.RotationAction_SS,
	})
}

func (deathKnight *DpsDeathknight) UnholyDiseaseCheckWrapper(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
	success := false

	if !deathKnight.TargetHasDisease(deathknight.FrostFeverAuraLabel, target) || deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim) < spell.CurCast.GCD {
		success = deathKnight.CastIcyTouch(sim, target)
	} else if !deathKnight.TargetHasDisease(deathknight.BloodPlagueAuraLabel, target) || deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim) < spell.CurCast.GCD {
		success = deathKnight.CastPlagueStrike(sim, target)
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

			crpb.Spend(sim, spell, spellCost)

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

func (deathKnight *DpsDeathknight) shouldWaitForDnD(sim *core.Simulation, blood bool, frost bool, unholy bool) bool {
	return deathKnight.Rotation.UseDeathAndDecay && !(deathKnight.Talents.Morbidity == 0 || !(deathKnight.DeathAndDecay.CD.IsReady(sim) || deathKnight.DeathAndDecay.CD.TimeToReady(sim) < 4*time.Second) || ((!blood || deathKnight.CurrentBloodRunes() > 1) && (!frost || deathKnight.CurrentFrostRunes() > 1) && (!unholy || deathKnight.CurrentUnholyRunes() > 1)))
}

var recastedFF = false
var recastedBP = false

func (deathKnight *DpsDeathknight) shouldSpreadDisease(sim *core.Simulation) bool {
	return recastedFF && recastedBP && deathKnight.Env.GetNumTargets() > 1
}

func (deathKnight *DpsDeathknight) spreadDiseases(sim *core.Simulation, target *core.Unit) {
	deathKnight.Pestilence.Cast(sim, target)
	recastedFF = false
	recastedBP = false
}

func (deathKnight *DpsDeathknight) doUnholyRotation(sim *core.Simulation, target *core.Unit) bool {
	casted := &deathKnight.CastSuccessful
	// I suggest adding the a wrapper around each spell you cast like this:
	// deathKnight.YourWrapper(sim, target, deathKnight.FrostStrike) that returns a bool for when you casted
	// since the waiting code relies on knowing if you actually casted

	diseaseRefreshDuration := time.Duration(deathKnight.Rotation.DiseaseRefreshDuration) * time.Second
	// Horn of Winter if you're the DK to refresh it and its not precasted/active
	if deathKnight.ShouldHornOfWinter(sim) {
		deathKnight.HornOfWinter.Cast(sim, target)
		*casted = true
	} else if (!deathKnight.TargetHasDisease(deathknight.FrostFeverAuraLabel, target) || deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim) < diseaseRefreshDuration) && deathKnight.CanIcyTouch(sim) {
		// Dont clip if theres half a second left to tick
		remainingDuration := deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim)
		if remainingDuration < time.Millisecond*500 && remainingDuration > 0 {
			deathKnight.WaitUntil(sim, sim.CurrentTime+remainingDuration+1)
		} else {
			deathKnight.IcyTouch.Cast(sim, target)
			*casted = true
			recastedFF = true
		}
	} else if (!deathKnight.TargetHasDisease(deathknight.BloodPlagueAuraLabel, target) || deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim) < diseaseRefreshDuration) && deathKnight.CanPlagueStrike(sim) {
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
		if deathKnight.Talents.Desolation > 0 && !deathKnight.DesolationAura.IsActive() && deathKnight.CanBloodStrike(sim) && !deathKnight.shouldWaitForDnD(sim, true, false, false) {
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
