package dps

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) RotationActionCallback_UA_Frost(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := dk.CastUnbreakableArmor(sim, target)
	dk.WaitUntil(sim, sim.CurrentTime)
	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_BT_Frost(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := dk.CastBloodTap(sim, target)
	dk.WaitUntil(sim, sim.CurrentTime)
	return casted
}

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

func (dk *DpsDeathknight) RotationActionCallback_HB_FS(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false
	if dk.KillingMachineAura.IsActive() && !dk.RimeAura.IsActive() {
		casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.FrostStrike)
	} else if dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
		if dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() < 110 {
			casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.HowlingBlast)
		} else if dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() > 110 {
			casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.HowlingBlast)
		} else if !dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() > 110 {
			casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.FrostStrike)
		} else if !dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() < 110 {
			casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.FrostStrike)
		}
	} else if !dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
		if dk.CurrentRunicPower() < 110 {
			casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.HowlingBlast)
		} else {
			casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.FrostStrike)
		}
	} else {
		casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.FrostStrike)
		if !casted {
			casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.HornOfWinter)
		}
	}
	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_FrostPrioRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	fr := &dk.fr

	casted := false

	if fr.nextSpell == dk.UnbreakableArmor {
		casted = dk.RotationActionCallback_UA_Frost(sim, target, s)
		if casted {
			fr.nextSpell = dk.Pestilence
		}
	} else if fr.nextSpell == dk.Pestilence {
		casted = dk.CastPestilence(sim, target)
		if casted && dk.LastCastOutcome.Matches(core.OutcomeLanded) {
			fr.nextSpell = nil
		}
	} else if dk.ShouldHornOfWinter(sim) {
		casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.HornOfWinter)
	} else {
		if dk.KillingMachineAura.IsActive() && !dk.RimeAura.IsActive() && fr.lastSpell == dk.Obliterate {
			casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.FrostStrike)
		} else {
			casted = dk.FrostDiseaseCheckWrapper(sim, target, dk.Obliterate)
			if !casted {
				casted = dk.RotationActionCallback_HB_FS(sim, target, s)
			}
		}
	}

	if !casted && dk.IsLeftBloodRuneNormal() && dk.CanBloodTap(sim) {
		casted = dk.RotationActionCallback_BT_Frost(sim, target, s)
		fr.nextSpell = dk.UnbreakableArmor
	} else if !casted && (dk.CurrentBloodRunes() > 0 || dk.CurrentDeathRunes() > 0) {
		casted = dk.CastBloodStrike(sim, target)
		fr.nextSpell = dk.Pestilence
	}

	return casted
}

// func (dk *DpsDeathknight) RotationActionCallback_FrostPrioRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
// 	fr := &dk.fr

// 	casted := false
// 	numActions := fr.numActions
// 	nextAction := fr.actions[fr.idx]
// 	advance := true

// 	switch nextAction {
// 	case FrostRotationAction_Obli:
// 		if dk.FrostDiseaseCheck(sim, target, dk.Obliterate, true, core.TernaryInt(fr.idx == 0, 2, 1)) {
// 			casted = dk.CastObliterate(sim, target)
// 		} else {
// 			if !dk.FrostFeverDisease[target.Index].IsActive() {
// 				casted = dk.CastIcyTouch(sim, target)
// 			} else if !dk.BloodPlagueDisease[target.Index].IsActive() {
// 				casted = dk.CastPlagueStrike(sim, target)
// 			} else {
// 				casted = dk.CastPestilence(sim, target)
// 			}
// 		}
// 	case FrostRotationAction_BS:
// 		casted = dk.CastBloodStrike(sim, target)
// 		advance = casted && dk.LastCastOutcome.Matches(core.OutcomeLanded)
// 	case FrostRotationAction_Pesti:
// 		casted = dk.CastPestilence(sim, target)
// 		advance = casted && dk.LastCastOutcome.Matches(core.OutcomeLanded)
// 	}

// 	if fr.idx+1 < numActions {
// 		if advance {
// 			fr.idx += 1
// 		}
// 	} else {
// 		fr.Reset(sim)
// 	}

// 	return casted
// }

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
		NewAction(dk.RotationActionCallback_FS).
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

func (dk *DpsDeathknight) FrostDiseaseCheckWrapper(sim *core.Simulation, target *core.Unit, spell *deathknight.RuneSpell) bool {
	fr := &dk.fr

	gcd := dk.SpellGCD()

	success := false

	if !dk.FrostFeverDisease[target.Index].IsActive() {
		success = dk.CastIcyTouch(sim, target)
		fr.SetLastSpell(success, dk.IcyTouch)
	} else if !dk.BloodPlagueDisease[target.Index].IsActive() {
		success = dk.CastPlagueStrike(sim, target)
		fr.SetLastSpell(success, dk.PlagueStrike)
	} else if dk.FrostFeverDisease[target.Index].RemainingDuration(sim) < gcd ||
		dk.BloodPlagueDisease[target.Index].RemainingDuration(sim) < gcd {
		success = dk.CastPestilence(sim, target)
		fr.SetLastSpell(success, dk.Pestilence)
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
			if gcd > ffExpiresIn || gcd > bpExpiresIn {
				success = false
				return success
			}

			crpb := dk.CopyRunicPowerBar()
			spellCost := crpb.OptimalRuneCost(core.RuneCost(spell.DefaultCast.Cost))

			// Add whichever non-frost specific checks you want here, I guess you'll need them.

			if !(dk.RimeAura.IsActive() && spell == dk.HowlingBlast) {
				crpb.SpendRuneCost(sim, spell.Spell, spellCost)
			}

			if crpb.CurrentBloodRunes() == 0 {
				nextBloodRuneAt := float64(crpb.BloodRuneReadyAt(sim))

				ff1 := float64(ffExpiresAt) > nextBloodRuneAt
				ff2 := sim.CurrentTime+gcd < ffExpiresAt

				bp1 := float64(bpExpiresAt) > nextBloodRuneAt
				bp2 := sim.CurrentTime+gcd < bpExpiresAt

				ff := ff1 && ff2
				bp := bp1 && bp2

				if ff || bp {
					spell.Cast(sim, target)
					success = true
				} else {
					return success
				}
			} else {
				spell.Cast(sim, target)
				success = true
			}
		}

		fr.SetLastSpell(success, spell)
	}

	return success
}
