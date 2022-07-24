package dps

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (deathKnight *DpsDeathknight) setupFrostRotations() {

	// This defines the Sub Blood opener
	deathKnight.DefineOpener(deathknight.RotationID_FrostSubBlood_Full, []deathknight.RotationAction{
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

	// This defines the Sub Unholy opener
	deathKnight.DefineOpener(deathknight.RotationID_FrostSubUnholy_Full, []deathknight.RotationAction{
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

func (deathKnight *DpsDeathknight) FrostDiseaseCheckWrapper(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
	success := false

	if !deathKnight.TargetHasDisease(deathknight.FrostFeverAuraLabel, target) {
		success = deathKnight.CastIcyTouch(sim, target)
	} else if !deathKnight.TargetHasDisease(deathknight.BloodPlagueAuraLabel, target) {
		success = deathKnight.CastPlagueStrike(sim, target)
	} else if deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim) < spell.CurCast.GCD ||
		deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim) < spell.CurCast.GCD {
		success = deathKnight.CastPestilence(sim, target)
		if deathKnight.LastCastOutcome == core.OutcomeMiss {
			// Deal with pestilence miss
			// TODO:

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

func (deathKnight *DpsDeathknight) doFrostRotation(sim *core.Simulation, target *core.Unit) {
	casted := &deathKnight.CastSuccessful

	if deathKnight.ShouldHornOfWinter(sim) {
		*casted = deathKnight.CastHornOfWinter(sim, target)
	} else {
		*casted = deathKnight.FrostDiseaseCheckWrapper(sim, target, deathKnight.Obliterate)
		if !*casted {
			if deathKnight.KillingMachineAura.IsActive() && !deathKnight.RimeAura.IsActive() {
				*casted = deathKnight.FrostDiseaseCheckWrapper(sim, target, deathKnight.FrostStrike)
			} else if deathKnight.KillingMachineAura.IsActive() && deathKnight.RimeAura.IsActive() {
				if deathKnight.CastCostPossible(sim, 0, 0, 1, 1) && deathKnight.CurrentRunicPower() < 110 {
					*casted = deathKnight.FrostDiseaseCheckWrapper(sim, target, deathKnight.HowlingBlast)
				} else if deathKnight.CastCostPossible(sim, 0, 0, 1, 1) && deathKnight.CurrentRunicPower() > 110 {
					*casted = deathKnight.FrostDiseaseCheckWrapper(sim, target, deathKnight.HowlingBlast)
				} else if !deathKnight.CastCostPossible(sim, 0, 0, 1, 1) && deathKnight.CurrentRunicPower() > 110 {
					*casted = deathKnight.FrostDiseaseCheckWrapper(sim, target, deathKnight.FrostStrike)
				} else if !deathKnight.CastCostPossible(sim, 0, 0, 1, 1) && deathKnight.CurrentRunicPower() < 110 {
					*casted = deathKnight.FrostDiseaseCheckWrapper(sim, target, deathKnight.FrostStrike)
				}
			} else if !deathKnight.KillingMachineAura.IsActive() && deathKnight.RimeAura.IsActive() {
				if deathKnight.CurrentRunicPower() < 110 {
					*casted = deathKnight.FrostDiseaseCheckWrapper(sim, target, deathKnight.HowlingBlast)
				} else {
					*casted = deathKnight.FrostDiseaseCheckWrapper(sim, target, deathKnight.FrostStrike)
				}
			} else {
				*casted = deathKnight.FrostDiseaseCheckWrapper(sim, target, deathKnight.FrostStrike)
				if !*casted {
					*casted = deathKnight.FrostDiseaseCheckWrapper(sim, target, deathKnight.HornOfWinter)
				}
			}
		}
	}
}
