package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) setupUnholyRotations() {

	dk.DefineOpener(deathknight.RotationID_UnholySsUnholyPresence_Full, []deathknight.RotationAction{
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

	dk.DefineOpener(deathknight.RotationID_UnholySsArmyUnholyPresence_Full, []deathknight.RotationAction{
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

	dk.DefineOpener(deathknight.RotationID_UnholySsBloodPresence_Full, []deathknight.RotationAction{
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

	dk.DefineOpener(deathknight.RotationID_UnholySsArmyBloodPresence_Full, []deathknight.RotationAction{
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

	dk.DefineOpener(deathknight.RotationID_UnholyDnd_Full, []deathknight.RotationAction{
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

func (dk *DpsDeathknight) UnholyDiseaseCheckWrapper(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
	success := false

	if !dk.TargetHasDisease(deathknight.FrostFeverAuraLabel, target) || dk.FrostFeverDisease[target.Index].RemainingDuration(sim) < spell.CurCast.GCD {
		success = dk.CastIcyTouch(sim, target)
	} else if !dk.TargetHasDisease(deathknight.BloodPlagueAuraLabel, target) || dk.BloodPlagueDisease[target.Index].RemainingDuration(sim) < spell.CurCast.GCD {
		success = dk.CastPlagueStrike(sim, target)
	} else {
		if dk.CanCast(sim, spell) {
			ffExpiresIn := dk.FrostFeverDisease[target.Index].RemainingDuration(sim)
			bpExpiresIn := dk.BloodPlagueDisease[target.Index].RemainingDuration(sim)
			ffExpiresAt := ffExpiresIn + sim.CurrentTime
			bpExpiresAt := bpExpiresIn + sim.CurrentTime
			if spell.CurCast.GCD > ffExpiresIn || spell.CurCast.GCD > bpExpiresIn {
				return success
			}

			crpb := dk.CopyRunicPowerBar()
			runeCostForSpell := dk.RuneAmountForSpell(spell)
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
					if dk.CanCast(sim, spell) {
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

func (dk *DpsDeathknight) shouldWaitForDnD(sim *core.Simulation, blood bool, frost bool, unholy bool) bool {
	return dk.Rotation.UseDeathAndDecay && !(dk.Talents.Morbidity == 0 || !(dk.DeathAndDecay.CD.IsReady(sim) || dk.DeathAndDecay.CD.TimeToReady(sim) < 4*time.Second) || ((!blood || dk.CurrentBloodRunes() > 1) && (!frost || dk.CurrentFrostRunes() > 1) && (!unholy || dk.CurrentUnholyRunes() > 1)))
}

var recastedFF = false
var recastedBP = false

func (dk *DpsDeathknight) shouldSpreadDisease(sim *core.Simulation) bool {
	return recastedFF && recastedBP && dk.Env.GetNumTargets() > 1
}

func (dk *DpsDeathknight) spreadDiseases(sim *core.Simulation, target *core.Unit) {
	dk.Pestilence.Cast(sim, target)
	recastedFF = false
	recastedBP = false
}

func (dk *DpsDeathknight) doUnholyRotation(sim *core.Simulation, target *core.Unit) bool {
	casted := &dk.CastSuccessful
	// I suggest adding the a wrapper around each spell you cast like this:
	// dk.YourWrapper(sim, target, dk.FrostStrike) that returns a bool for when you casted
	// since the waiting code relies on knowing if you actually casted

	diseaseRefreshDuration := time.Duration(dk.Rotation.DiseaseRefreshDuration) * time.Second
	// Horn of Winter if you're the DK to refresh it and its not precasted/active
	if dk.ShouldHornOfWinter(sim) {
		dk.HornOfWinter.Cast(sim, target)
		*casted = true
	} else if (!dk.TargetHasDisease(deathknight.FrostFeverAuraLabel, target) || dk.FrostFeverDisease[target.Index].RemainingDuration(sim) < diseaseRefreshDuration) && dk.CanIcyTouch(sim) {
		// Dont clip if theres half a second left to tick
		remainingDuration := dk.FrostFeverDisease[target.Index].RemainingDuration(sim)
		if remainingDuration < time.Millisecond*500 && remainingDuration > 0 {
			dk.WaitUntil(sim, sim.CurrentTime+remainingDuration+1)
		} else {
			dk.IcyTouch.Cast(sim, target)
			*casted = true
			recastedFF = true
		}
	} else if (!dk.TargetHasDisease(deathknight.BloodPlagueAuraLabel, target) || dk.BloodPlagueDisease[target.Index].RemainingDuration(sim) < diseaseRefreshDuration) && dk.CanPlagueStrike(sim) {
		// Dont clip if theres half a second left to tick
		remainingDuration := dk.BloodPlagueDisease[target.Index].RemainingDuration(sim)
		if remainingDuration < time.Millisecond*500 && remainingDuration > 0 {
			dk.WaitUntil(sim, sim.CurrentTime+remainingDuration+1)
		} else {
			dk.PlagueStrike.Cast(sim, target)
			*casted = true
			recastedBP = true
		}
	} else {
		if dk.Talents.Desolation > 0 && !dk.DesolationAura.IsActive() && dk.CanBloodStrike(sim) && !dk.shouldWaitForDnD(sim, true, false, false) {
			// Desolation and Pestilence check
			if dk.shouldSpreadDisease(sim) {
				dk.spreadDiseases(sim, target)
				*casted = true
			} else {
				dk.BloodStrike.Cast(sim, target)
				*casted = true
			}
		} else {
			if dk.Rotation.UseDeathAndDecay {
				// Death and Decay Rotation
				if dk.CanDeathAndDecay(sim) && dk.AllDiseasesAreActive(target) {
					dk.DeathAndDecay.Cast(sim, target)
					*casted = true
				} else if dk.CanGhoulFrenzy(sim) && dk.Talents.MasterOfGhouls && (!dk.Ghoul.GhoulFrenzyAura.IsActive() || dk.Ghoul.GhoulFrenzyAura.RemainingDuration(sim) < 6*time.Second) && !dk.shouldWaitForDnD(sim, false, false, true) {
					dk.GhoulFrenzy.Cast(sim, target)
					*casted = true
				} else if dk.CanScourgeStrike(sim) && !dk.shouldWaitForDnD(sim, false, true, true) {
					dk.ScourgeStrike.Cast(sim, target)
					*casted = true
				} else if !dk.Talents.ScourgeStrike && dk.CanIcyTouch(sim) && !dk.shouldWaitForDnD(sim, false, true, false) {
					dk.IcyTouch.Cast(sim, target)
					*casted = true
				} else if !dk.Talents.ScourgeStrike && dk.CanPlagueStrike(sim) && !dk.shouldWaitForDnD(sim, false, false, true) {
					dk.PlagueStrike.Cast(sim, target)
					*casted = true
				} else if dk.CanBloodStrike(sim) && !dk.shouldWaitForDnD(sim, true, false, false) {
					if dk.shouldSpreadDisease(sim) {
						dk.spreadDiseases(sim, target)
						*casted = true
					} else if dk.Env.GetNumTargets() > 2 {
						dk.BloodBoil.Cast(sim, target)
						*casted = true
					} else {
						dk.BloodStrike.Cast(sim, target)
						*casted = true
					}
				} else if dk.CanDeathCoil(sim) && !dk.SummonGargoyle.IsReady(sim) {
					dk.DeathCoil.Cast(sim, target)
					*casted = true
				} else if dk.CanHornOfWinter(sim) {
					dk.HornOfWinter.Cast(sim, target)
					*casted = true
				} else {
					// Probably want to make this just return *casted as casted should be false in this case, the wait time will be handled after the return
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
				if dk.CanGhoulFrenzy(sim) && dk.Talents.MasterOfGhouls && (!dk.Ghoul.GhoulFrenzyAura.IsActive() || dk.Ghoul.GhoulFrenzyAura.RemainingDuration(sim) < 6*time.Second) {
					dk.GhoulFrenzy.Cast(sim, target)
					*casted = true
				} else if dk.CanScourgeStrike(sim) {
					dk.ScourgeStrike.Cast(sim, target)
					*casted = true
				} else if dk.CanBloodStrike(sim) {
					if dk.shouldSpreadDisease(sim) {
						dk.spreadDiseases(sim, target)
						*casted = true
					} else if dk.Env.GetNumTargets() > 2 {
						dk.BloodBoil.Cast(sim, target)
						*casted = true
					} else {
						dk.BloodStrike.Cast(sim, target)
						*casted = true
					}
				} else if dk.CanDeathCoil(sim) && !dk.SummonGargoyle.IsReady(sim) {
					dk.DeathCoil.Cast(sim, target)
					*casted = true
				} else if dk.CanHornOfWinter(sim) {
					dk.HornOfWinter.Cast(sim, target)
					*casted = true
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
	return true
}
