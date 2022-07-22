package deathknight

import "github.com/wowsims/wotlk/sim/core"

func (deathKnight *DeathKnight) setupFrostRotations() {

	// This defines the Sub Blood opener
	deathKnight.DefineOpener(RotationID_FrostSubBlood_Full, []RotationAction{
		RotationAction_IT,
		RotationAction_PS,
		RotationAction_UA,
		RotationAction_BT,
		RotationAction_Obli,
		RotationAction_FS,
		RotationAction_Pesti,
		RotationAction_ERW,
		RotationAction_Obli,
		RotationAction_Obli,
		RotationAction_Obli,
		RotationAction_FS,
		RotationAction_HB_Ghoul_RimeCheck,
		RotationAction_FS,
		RotationAction_Obli,
		RotationAction_Obli,
		RotationAction_Pesti,
		RotationAction_FS,
		RotationAction_BS,
		RotationAction_FS,
	})

	// This defines the Sub Unholy opener
	deathKnight.DefineOpener(RotationID_FrostSubUnholy_Full, []RotationAction{
		RotationAction_IT,
		RotationAction_PS,
		RotationAction_BT,
		RotationAction_Pesti,
		RotationAction_UA,
		RotationAction_Obli,
		RotationAction_FS,
		RotationAction_ERW,
		RotationAction_Obli,
		RotationAction_Obli,
		RotationAction_Obli,
		RotationAction_FS,
		RotationAction_FS,
		RotationAction_FS,
		RotationAction_Obli,
		RotationAction_Obli,
		RotationAction_BS,
		RotationAction_Pesti,
		RotationAction_FS,
	})
}

func (deathKnight *DeathKnight) FrostDiseaseCheckWrapper(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
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

			crpb := deathKnight.GetCalcRunicPowerBar()
			spellCost := DetermineOptimalCostForSpell(&crpb, sim, deathKnight, spell)

			// Add whichever non-frost specific checks you want here, I guess you'll need them.

			if !(deathKnight.RimeAura.IsActive() && spell == deathKnight.HowlingBlast) {
				crpb.Spend(sim, spellCost)
			}

			if crpb.CurrentBloodRunes() == 0 && crpb.CurrentDeathRunes() == 0 {
				nextBloodRuneAt := crpb.BloodRuneReadyAt(sim)
				nextDeathRuneAt := crpb.DeathRuneReadyAt(sim)

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

func (deathKnight *DeathKnight) doFrostRotation(sim *core.Simulation, target *core.Unit) {
	casted := &deathKnight.castSuccessful

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
