package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) setupUnholySsUnholyPresenceOpener() {
	dk.DefineOpener([]deathknight.RotationAction{
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
}

func (dk *DpsDeathknight) setupUnholySsArmyUnholyPresenceOpener() {
	dk.DefineOpener([]deathknight.RotationAction{
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
}

func (dk *DpsDeathknight) setupUnholySsBloodPresenceOpener() {
	dk.DefineOpener([]deathknight.RotationAction{
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
}

func (dk *DpsDeathknight) setupUnholySsArmyBloodPresenceOpener() {
	dk.DefineOpener([]deathknight.RotationAction{
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
}

func (dk *DpsDeathknight) setupUnholyDndUnholyPresenceOpener() {
	dk.DefineOpener([]deathknight.RotationAction{
		deathknight.RotationAction_IT,
		deathknight.RotationAction_PS,
		deathknight.RotationAction_BS,
		deathknight.RotationAction_DND,
		deathknight.RotationAction_Garg,
		deathknight.RotationAction_ERW,
		deathknight.RotationAction_BP,
		deathknight.RotationAction_SS,
		deathknight.RotationAction_SS,
		deathknight.RotationAction_BS,
	})
}

func (dk *DpsDeathknight) setupUnholyDndBloodPresenceOpener() {
	dk.DefineOpener([]deathknight.RotationAction{
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
		deathknight.RotationAction_SS,
		deathknight.RotationAction_BS,
	})
}

var syncBp = false

func (dk *DpsDeathknight) UnholyDiseaseCheckWrapper(sim *core.Simulation, target *core.Unit, spell *core.Spell, costRunes bool) bool {
	success := false

	dropTimeAllowed := time.Millisecond * -100
	ffRemaining := dk.FrostFeverDisease[target.Index].RemainingDuration(sim) + dropTimeAllowed
	bpRemaining := dk.BloodPlagueDisease[target.Index].RemainingDuration(sim) + dropTimeAllowed

	if !dk.TargetHasDisease(deathknight.FrostFeverAuraLabel, target) || ffRemaining < spell.CurCast.GCD {
		// Refresh FF
		success = dk.CastIcyTouch(sim, target)
	} else if syncBp || !dk.TargetHasDisease(deathknight.BloodPlagueAuraLabel, target) || bpRemaining < spell.CurCast.GCD {
		// Refresh BP
		syncBp = false
		success = dk.CastPlagueStrike(sim, target)
	} else if dk.CanCast(sim, spell) {
		ffExpiresAt := ffRemaining + sim.CurrentTime
		bpExpiresAt := bpRemaining + sim.CurrentTime

		crpb := dk.CopyRunicPowerBar()
		runeCostForSpell := dk.RuneAmountForSpell(spell)
		spellCost := crpb.DetermineOptimalCost(sim, runeCostForSpell.Blood, runeCostForSpell.Frost, runeCostForSpell.Unholy)

		crpb.Spend(sim, spell, spellCost)

		// Check FF
		if spellCost.Frost > 0 && crpb.CurrentFrostRunes() == 0 && crpb.CurrentDeathRunes() == 0 {
			nextFrostRuneAt := crpb.FrostRuneReadyAt(sim)
			nextDeathRuneAt := crpb.DeathRuneReadyAt(sim)

			// Can cast FF with frost before expire
			ff1 := ffExpiresAt > nextFrostRuneAt && spell.CurCast.GCD < ffRemaining
			// Can cast FF with death before expire
			ff2 := ffExpiresAt > nextDeathRuneAt && spell.CurCast.GCD < ffRemaining

			if !ff1 && !ff2 {
				// Refresh FF
				success = dk.CastIcyTouch(sim, target)
				syncBp = success
				return success
			}
		}

		// Check BP
		if spellCost.Unholy > 0 && crpb.CurrentUnholyRunes() == 0 && crpb.CurrentDeathRunes() == 0 {
			nextUnholyRuneAt := crpb.UnholyRuneReadyAt(sim)
			nextDeathRuneAt := crpb.DeathRuneReadyAt(sim)

			// Can cast BP with unholy before expire
			bp1 := bpExpiresAt > nextUnholyRuneAt && spell.CurCast.GCD < bpRemaining
			// Can cast BP with death before expire
			bp2 := bpExpiresAt > nextDeathRuneAt && spell.CurCast.GCD < bpRemaining

			if !bp1 && !bp2 {
				// Refresh BP
				success = dk.CastPlagueStrike(sim, target)
				return success
			}
		}

		// We have runes left for disease after this cast
		spell.Cast(sim, target)
		success = true
	}

	return success
}

func (dk *DpsDeathknight) doUnholySsRotation(sim *core.Simulation, target *core.Unit) {
	casted := &dk.CastSuccessful

	if dk.ShouldHornOfWinter(sim) {
		*casted = dk.CastHornOfWinter(sim, target)
	} else {
		*casted = dk.UnholyDiseaseCheckWrapper(sim, target, dk.ScourgeStrike, true)
		if !*casted {
			*casted = dk.UnholyDiseaseCheckWrapper(sim, target, dk.BloodStrike, true)
			if !*casted {
				*casted = dk.UnholyDiseaseCheckWrapper(sim, target, dk.DeathCoil, false)
				if !*casted {
					*casted = dk.UnholyDiseaseCheckWrapper(sim, target, dk.HornOfWinter, false)
				}
			}
		}
	}
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

func (dk *DpsDeathknight) doUnholyRotation(sim *core.Simulation, target *core.Unit) {
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
				} else if dk.CanGhoulFrenzy(sim) && (!dk.Ghoul.GhoulFrenzyAura.IsActive() || dk.Ghoul.GhoulFrenzyAura.RemainingDuration(sim) < 6*time.Second) && !dk.shouldWaitForDnD(sim, false, false, true) {
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
				if dk.CanGhoulFrenzy(sim) && (!dk.Ghoul.GhoulFrenzyAura.IsActive() || dk.Ghoul.GhoulFrenzyAura.RemainingDuration(sim) < 6*time.Second) {
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
}
