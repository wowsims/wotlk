package dps

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) setupFrostRotations(rotationId RotationID) {
	if rotationId == RotationID_FrostSubBlood_Full {
		// This defines the Sub Blood opener
		dk.DefineOpener([]deathknight.RotationAction{
			deathknight.RotationAction_IT,
			deathknight.RotationAction_PS,
			deathknight.RotationAction_UA,
			deathknight.RotationAction_BT,
			deathknight.RotationAction_Obli,
			deathknight.RotationAction_FS,
			deathknight.RotationAction_Pesti,
			deathknight.RotationAction_ERW,
			deathknight.RotationAction_Obli,
			deathknight.RotationAction_Obli,
			deathknight.RotationAction_Obli,
			deathknight.RotationAction_FS,
			deathknight.RotationAction_HB_Ghoul_RimeCheck,
			deathknight.RotationAction_FS,
			deathknight.RotationAction_Obli,
			deathknight.RotationAction_Obli,
			deathknight.RotationAction_Pesti,
			deathknight.RotationAction_FS,
			deathknight.RotationAction_BS,
			deathknight.RotationAction_FS,
		})
	} else if rotationId == RotationID_FrostSubUnholy_Full {
		// This defines the Sub Unholy opener
		dk.DefineOpener([]deathknight.RotationAction{
			deathknight.RotationAction_IT,
			deathknight.RotationAction_PS,
			deathknight.RotationAction_BT,
			deathknight.RotationAction_Pesti,
			deathknight.RotationAction_UA,
			deathknight.RotationAction_Obli,
			deathknight.RotationAction_FS,
			deathknight.RotationAction_ERW,
			deathknight.RotationAction_Obli,
			deathknight.RotationAction_Obli,
			deathknight.RotationAction_Obli,
			deathknight.RotationAction_FS,
			deathknight.RotationAction_FS,
			deathknight.RotationAction_FS,
			deathknight.RotationAction_Obli,
			deathknight.RotationAction_Obli,
			deathknight.RotationAction_BS,
			deathknight.RotationAction_Pesti,
			deathknight.RotationAction_FS,
		})
	}

}

func (dk *DpsDeathknight) FrostDiseaseCheckWrapper(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
	success := false

	if !dk.TargetHasDisease(deathknight.FrostFeverAuraLabel, target) {
		success = dk.CastIcyTouch(sim, target)
	} else if !dk.TargetHasDisease(deathknight.BloodPlagueAuraLabel, target) {
		success = dk.CastPlagueStrike(sim, target)
	} else if dk.FrostFeverDisease[target.Index].RemainingDuration(sim) < spell.CurCast.GCD ||
		dk.BloodPlagueDisease[target.Index].RemainingDuration(sim) < spell.CurCast.GCD {
		success = dk.CastPestilence(sim, target)
		if dk.LastCastOutcome == core.OutcomeMiss {
			// Deal with pestilence miss
			// TODO:

		}
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

			// Add whichever non-frost specific checks you want here, I guess you'll need them.

			if !(dk.RimeAura.IsActive() && spell == dk.HowlingBlast) {
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

func (dk *DpsDeathknight) doFrostRotation(sim *core.Simulation, target *core.Unit) {
	casted := &dk.CastSuccessful

	if dk.ShouldHornOfWinter(sim) {
		*casted = dk.CastHornOfWinter(sim, target)
	} else {
		*casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.Obliterate)
		if !*casted {
			if dk.KillingMachineAura.IsActive() && !dk.RimeAura.IsActive() {
				*casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.FrostStrike)
			} else if dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
				if dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() < 110 {
					*casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.HowlingBlast)
				} else if dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() > 110 {
					*casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.HowlingBlast)
				} else if !dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() > 110 {
					*casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.FrostStrike)
				} else if !dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() < 110 {
					*casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.FrostStrike)
				}
			} else if !dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
				if dk.CurrentRunicPower() < 110 {
					*casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.HowlingBlast)
				} else {
					*casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.FrostStrike)
				}
			} else {
				*casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.FrostStrike)
				if !*casted {
					*casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.HornOfWinter)
				}
			}
		}
	}
}
