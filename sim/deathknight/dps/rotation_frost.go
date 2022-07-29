package dps

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) RotationActionCallback_HB_Ghoul_RimeCheck(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false
	if dk.RimeAura.IsActive() {
		casted = dk.CastHowlingBlast(sim, target)
	} else {
		casted = dk.CastRaiseDead(sim, target)
	}

	s.ConditionalAdvance(true)
	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_FrostPrioRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	fr := &dk.fr

	casted := false
	numActions := fr.numActions
	nextAction := fr.actions[fr.idx]
	advance := true

	switch nextAction {
	case FrostRotationAction_Obli:
		if dk.FrostDiseaseCheck(sim, target, dk.Obliterate, true, core.TernaryInt(fr.idx == 0, 2, 1)) {
			casted = dk.CastObliterate(sim, target)
		} else {
			casted = dk.CastPestilence(sim, target)
		}
	case FrostRotationAction_BS:
		casted = dk.CastBloodStrike(sim, target)
		advance = casted && dk.LastCastOutcome.Matches(core.OutcomeLanded)
	case FrostRotationAction_Pesti:
		casted = dk.CastPestilence(sim, target)
		advance = casted && dk.LastCastOutcome.Matches(core.OutcomeLanded)
	}

	if fr.idx+1 < numActions {
		if advance {
			fr.idx += 1
		}
	} else {
		fr.Reset(sim)
	}

	return casted
}

func (dk *DpsDeathknight) setupFrostSubBloodOpener() {
	dk.Opener.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_UA).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_Pesti).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_HB_Ghoul_RimeCheck).
		NewAction(dk.RotationActionCallback_FS).
		//NewAction(dk.RotationActionCallback_FS).
		//NewAction(dk.RotationActionCallback_HW).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Pesti).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_FS)

	dk.Main.
		NewAction(dk.RotationActionCallback_FrostPrioRotation)
}

func (dk *DpsDeathknight) setupFrostSubUnholyOpener() {
	dk.Opener.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_Pesti).
		NewAction(dk.RotationActionCallback_UA).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_Pesti).
		NewAction(dk.RotationActionCallback_FS)

	dk.Main.
		NewAction(dk.RotationActionCallback_FrostPrioRotation)
}

func (dk *DpsDeathknight) FrostDiseaseCheckWrapper(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
	success := false

	if !dk.FrostFeverDisease[target.Index].IsActive() {
		success = dk.CastIcyTouch(sim, target)
	} else if !dk.BloodPlagueDisease[target.Index].IsActive() {
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
