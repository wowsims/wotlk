package deathknight

import "github.com/wowsims/wotlk/sim/core"

type RotationAction uint8

// Add your UH rotation Actions here and then on the DoNext function
const (
	RotationAction_Skip RotationAction = iota
	RotationAction_IT
	RotationAction_PS
	RotationAction_Obli
	RotationAction_BS
	RotationAction_BT
	RotationAction_UA
	RotationAction_RD
	RotationAction_Pesti
	RotationAction_FS
	RotationAction_HW
	RotationAction_ERW
	RotationAction_HB_Ghoul_RimeCheck
	RotationAction_PrioMode
)

type RotationID uint8

const (
	RotationID_Default RotationID = iota
	RotationID_FrostSubBlood_Full
	RotationID_FrostSubUnholy_Full
	RotationID_Unholy_Full
	RotationID_Count
	RotationID_Unknown
)

type Sequence struct {
	id         RotationID
	idx        int
	numActions int
	actions    []RotationAction
}

type RotationHelper struct {
	onOpener bool
	opener   *Sequence
	openers  []Sequence

	sequence *Sequence

	castSuccessful     bool
	justCastPestilence bool
}

func TernaryRotationAction(condition bool, t RotationAction, f RotationAction) RotationAction {
	if condition {
		return t
	} else {
		return f
	}
}

func (r *RotationHelper) DefineOpener(id RotationID, actions []RotationAction) {
	o := &r.openers[id]
	o.id = id
	o.idx = 0
	o.numActions = len(actions)
	o.actions = actions
}

func (r *RotationHelper) PushSequence(actions []RotationAction) {
	seq := &Sequence{}
	seq.id = RotationID_Unknown
	seq.idx = 0
	seq.numActions = len(actions)
	seq.actions = actions
	r.sequence = seq
}

func (deathKnight *DeathKnight) SetupRotation() {
	deathKnight.openers = make([]Sequence, RotationID_Count)

	// This defines the Sub Blood opener
	deathKnight.DefineOpener(RotationID_Default, []RotationAction{})

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

	// To define the opener for Unholy for example (or any other UH tree you want
	// just define the enum accordingly) it goes as follows:
	deathKnight.DefineOpener(RotationID_Unholy_Full, []RotationAction{})

	// IMPORTANT
	rotationId := RotationID_Unknown
	// Also you need to update this to however you define spec
	if deathKnight.Talents.DarkConviction > 0 && deathKnight.Talents.HowlingBlast {
		rotationId = RotationID_FrostSubBlood_Full
	} else if deathKnight.Talents.BloodCakedBlade > 0 && deathKnight.Talents.HowlingBlast {
		rotationId = RotationID_FrostSubUnholy_Full
	} else if deathKnight.Talents.SummonGargoyle {
		rotationId = RotationID_Unholy_Full
	} else {
		rotationId = RotationID_Default
	}

	deathKnight.opener = &deathKnight.openers[rotationId]
	deathKnight.onOpener = true
}

func (deathKnight *DeathKnight) DiseaseCheckWrapper(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
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
			if deathKnight.opener.id == RotationID_FrostSubUnholy_Full && deathKnight.onOpener {
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
