package deathknight

import "github.com/wowsims/wotlk/sim/core"

type OpenerAction uint8

// Add your UH rotation Actions here and then on the DoNext function
const (
	OpenerAction_Skip OpenerAction = iota
	OpenerAction_IT
	OpenerAction_PS
	OpenerAction_Obli
	OpenerAction_BS
	OpenerAction_BT
	OpenerAction_UA
	OpenerAction_RD
	OpenerAction_Pesti
	OpenerAction_FS
	OpenerAction_HW
	OpenerAction_ERW
	OpenerAction_HB_Ghoul_RimeCheck
	OpenerAction_PrioMode
)

type OpenerID uint8

const (
	OpenerID_FrostSubBlood_Full OpenerID = iota
	OpenerID_FrostSubUnholy_Full
	OpenerID_Unholy_Full
	OpenerID_Count
	OpenerID_Unknown
)

type Opener struct {
	id         OpenerID
	idx        int
	numActions int
	actions    []OpenerAction
}

type RotationHelper struct {
	onOpener bool
	opener   *Opener
	openers  []Opener

	castSuccessful     bool
	justCastPestilence bool
}

func TernaryOpenerAction(condition bool, t OpenerAction, f OpenerAction) OpenerAction {
	if condition {
		return t
	} else {
		return f
	}
}

func (r *RotationHelper) DefineOpener(id OpenerID, actions []OpenerAction) {
	o := &r.openers[id]
	o.id = id
	o.idx = 0
	o.numActions = len(actions)
	o.actions = actions
}

func (deathKnight *DeathKnight) SetupRotation() {
	deathKnight.openers = make([]Opener, OpenerID_Count)

	// This defines the Sub Blood opener
	deathKnight.DefineOpener(OpenerID_FrostSubBlood_Full, []OpenerAction{
		OpenerAction_IT,
		OpenerAction_PS,
		OpenerAction_UA,
		OpenerAction_BT,
		OpenerAction_Obli,
		OpenerAction_FS,
		OpenerAction_Pesti,
		OpenerAction_ERW,
		OpenerAction_Obli,
		OpenerAction_Obli,
		OpenerAction_Obli,
		OpenerAction_FS,
		OpenerAction_HB_Ghoul_RimeCheck,
		OpenerAction_FS,
		OpenerAction_Obli,
		OpenerAction_Obli,
		OpenerAction_Pesti,
		OpenerAction_FS,
		OpenerAction_BS,
		OpenerAction_FS,
	})

	deathKnight.DefineOpener(OpenerID_FrostSubUnholy_Full, []OpenerAction{
		OpenerAction_IT,
		OpenerAction_PS,
		OpenerAction_BT,
		OpenerAction_Pesti,
		OpenerAction_UA,
		OpenerAction_Obli,
		OpenerAction_FS,
		OpenerAction_ERW,
		OpenerAction_Obli,
		OpenerAction_Obli,
		OpenerAction_Obli,
		OpenerAction_FS,
		OpenerAction_FS,
		OpenerAction_FS,
		OpenerAction_Obli,
		OpenerAction_Obli,
		OpenerAction_BS,
		OpenerAction_Pesti,
		OpenerAction_FS,
	})

	// To define the opener for Unholy for example (or any other UH tree you want
	// just define the enum accordingly) it goes as follows:
	deathKnight.DefineOpener(OpenerID_Unholy_Full, []OpenerAction{})

	// IMPORTANT
	openerId := OpenerID_Unknown
	// Also you need to update this to however you define spec
	if deathKnight.Talents.DarkConviction > 0 && deathKnight.Talents.HowlingBlast {
		openerId = OpenerID_FrostSubBlood_Full
	} else if deathKnight.Talents.BloodCakedBlade > 0 && deathKnight.Talents.HowlingBlast {
		openerId = OpenerID_FrostSubUnholy_Full
	} else if deathKnight.Talents.SummonGargoyle {
		openerId = OpenerID_Unholy_Full
	} else {
		panic("Unknown spec for rotation!")
	}

	deathKnight.opener = &deathKnight.openers[openerId]
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
		} else {

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
