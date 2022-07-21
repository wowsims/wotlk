package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

type OpenerAction uint8

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
	OpenerAction_HB_Ghoul_FS_RimeCheck
	OpenerAction_PrioMode
)

type OpenerID uint8

const (
	OpenerID_FrostSubBlood_Full OpenerID = iota
	OpenerID_FrostSubUnholy_Full
	OpenerID_Unholy_Full
	OpenerID_Count
)

type Opener struct {
	id         OpenerID
	idx        int
	numActions int
	actions    []OpenerAction
}

type DKRotation struct {
	onOpener bool
	opener   *Opener
	openers  []Opener

	canAdvance         bool
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

func (r *DKRotation) DefineOpener(id OpenerID, actions []OpenerAction) {
	o := &r.openers[id]
	o.id = id
	o.idx = 0
	o.numActions = len(actions)
	o.actions = actions
}

func (deathKnight *DeathKnight) SetupRotation() {
	r := &deathKnight.DKRotation
	r.openers = make([]Opener, OpenerID_Count)

	r.DefineOpener(OpenerID_FrostSubBlood_Full, []OpenerAction{
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
		OpenerAction_HB_Ghoul_FS_RimeCheck,
		OpenerAction_FS,
		OpenerAction_Obli,
		OpenerAction_Obli,
		OpenerAction_Pesti,
		OpenerAction_FS,
		OpenerAction_BS,
		OpenerAction_FS,
	})

	openerId := OpenerID_FrostSubBlood_Full
	if deathKnight.Talents.BloodCakedBlade > 0 {
		openerId = OpenerID_FrostSubUnholy_Full
	} else if deathKnight.Talents.SummonGargoyle {
		openerId = OpenerID_Unholy_Full
	}

	r.opener = &r.openers[openerId]
	r.onOpener = true
}

func (deathKnight *DeathKnight) DoRotation(sim *core.Simulation) {
	if !deathKnight.Talents.HowlingBlast {
		return
	}

	target := deathKnight.CurrentTarget

	if !deathKnight.Rotation.GetWipFrostRotation() {
		if deathKnight.ShouldHornOfWinter(sim) {
			deathKnight.HornOfWinter.Cast(sim, target)
		} else if (!deathKnight.TargetHasDisease(FrostFeverAuraLabel, target) || deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim) < 6*time.Second) && deathKnight.CanIcyTouch(sim) {
			deathKnight.IcyTouch.Cast(sim, target)
			recastedFF = true
		} else if (!deathKnight.TargetHasDisease(BloodPlagueAuraLabel, target) || deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim) < 6*time.Second) && deathKnight.CanPlagueStrike(sim) {
			deathKnight.PlagueStrike.Cast(sim, target)
			recastedBP = true
		} else {
			if deathKnight.CanBloodTap(sim) && deathKnight.AllDiseasesAreActive(target) {
				deathKnight.BloodTap.Cast(sim, target)
				deathKnight.WaitUntil(sim, sim.CurrentTime+1)
			} else if deathKnight.CanUnbreakableArmor(sim) && deathKnight.AllDiseasesAreActive(target) {
				deathKnight.UnbreakableArmor.Cast(sim, target)
				deathKnight.WaitUntil(sim, sim.CurrentTime+1)
			} else if deathKnight.CanPestilence(sim) && deathKnight.shouldSpreadDisease(sim) {
				deathKnight.spreadDiseases(sim, target)
			} else if deathKnight.CanObliterate(sim) && deathKnight.AllDiseasesAreActive(target) {
				deathKnight.Obliterate.Cast(sim, target)
			} else if deathKnight.CanHowlingBlast(sim) && deathKnight.AllDiseasesAreActive(target) {
				deathKnight.HowlingBlast.Cast(sim, target)
			} else if deathKnight.CanFrostStrike(sim) && deathKnight.AllDiseasesAreActive(target) {
				deathKnight.FrostStrike.Cast(sim, target)
			} else if deathKnight.CanBloodStrike(sim) && deathKnight.AllDiseasesAreActive(target) {
				deathKnight.BloodStrike.Cast(sim, target)
			} else if deathKnight.CanIcyTouch(sim) {
				deathKnight.IcyTouch.Cast(sim, target)
			} else if deathKnight.CanPlagueStrike(sim) {
				deathKnight.PlagueStrike.Cast(sim, target)
			} else if deathKnight.CanHornOfWinter(sim) {
				deathKnight.HornOfWinter.Cast(sim, target)
			} else {
				if deathKnight.GCD.IsReady(sim) && !deathKnight.IsWaiting() {
					// This means we did absolutely nothing.
					// Wait until our next auto attack to decide again.
					waitUntil := deathKnight.AutoAttacks.MainhandSwingAt
					if deathKnight.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
						waitUntil = core.MinDuration(waitUntil, deathKnight.AutoAttacks.OffhandSwingAt)
					}
					waitUntil = core.MinDuration(time.Duration(0.1*float64(waitUntil-sim.CurrentTime)+float64(waitUntil)), deathKnight.AnyRuneReadyAt(sim))
					deathKnight.WaitUntil(sim, waitUntil)
				}
			}
		}

	} else {
		opener := deathKnight.DKRotation.opener

		if !opener.DoNext(sim, deathKnight) {
			if deathKnight.GCD.IsReady(sim) && !deathKnight.IsWaiting() {
				waitUntil := deathKnight.AutoAttacks.MainhandSwingAt
				if deathKnight.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
					waitUntil = core.MinDuration(waitUntil, deathKnight.AutoAttacks.OffhandSwingAt)
				}
				waitUntil = core.MinDuration(waitUntil, deathKnight.AnyRuneReadyAt(sim))
				deathKnight.WaitUntil(sim, waitUntil)
			} else { // No resources
				waitUntil := deathKnight.AnySpentRuneReadyAt(sim)
				deathKnight.WaitUntil(sim, waitUntil)
			}
		}
	}
}

func (deathKnight *DeathKnight) DiseaseCheckWrapper(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
	success := false

	if !deathKnight.TargetHasDisease(FrostFeverAuraLabel, target) {
		success = deathKnight.CastIcyTouch(sim, target)
	} else if !deathKnight.TargetHasDisease(BloodPlagueAuraLabel, target) {
		success = deathKnight.CastPlagueStrike(sim, target)
	} else if deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim) < 4*time.Second ||
		deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim) < 4*time.Second {
		success = deathKnight.CastPestilence(sim, target)
		if deathKnight.LastCastOutcome == core.OutcomeMiss {
			// Deal with pestilence miss
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

func (o *Opener) DoNext(sim *core.Simulation, deathKnight *DeathKnight) bool {
	target := deathKnight.CurrentTarget
	casted := &deathKnight.DKRotation.castSuccessful
	advance := true
	*casted = false

	if o.idx < o.numActions {
		action := o.actions[o.idx]

		switch action {
		case OpenerAction_IT:
			*casted = deathKnight.CastIcyTouch(sim, target)
			advance = deathKnight.LastCastOutcome != core.OutcomeMiss
		case OpenerAction_PS:
			*casted = deathKnight.CastPlagueStrike(sim, target)
			advance = deathKnight.LastCastOutcome != core.OutcomeMiss
		case OpenerAction_UA:
			*casted = deathKnight.CastUnbreakableArmor(sim, target)
			deathKnight.WaitUntil(sim, sim.CurrentTime)
		case OpenerAction_BT:
			*casted = deathKnight.CastBloodTap(sim, target)
			deathKnight.WaitUntil(sim, sim.CurrentTime)
		case OpenerAction_Obli:
			*casted = deathKnight.CastObliterate(sim, target)
		case OpenerAction_FS:
			*casted = deathKnight.CastFrostStrike(sim, target)
		case OpenerAction_Pesti:
			*casted = deathKnight.CastPestilence(sim, target)
			if deathKnight.LastCastOutcome == core.OutcomeMiss {
				advance = false
			}
		case OpenerAction_ERW:
			*casted = deathKnight.CastEmpowerRuneWeapon(sim, target)
			deathKnight.WaitUntil(sim, sim.CurrentTime)
		case OpenerAction_HB_Ghoul_FS_RimeCheck:
			if deathKnight.RimeAura.IsActive() {
				*casted = deathKnight.CastHowlingBlast(sim, target)
			} else {
				*casted = deathKnight.CastRaiseDead(sim, target)
			}
		case OpenerAction_BS:
			*casted = deathKnight.CastBloodStrike(sim, target)
		}

		if *casted && advance {
			o.idx += 1
		}

	} else {
		deathKnight.DKRotation.onOpener = false

		if deathKnight.ShouldHornOfWinter(sim) {
			*casted = deathKnight.CastHornOfWinter(sim, target)
		} else {
			*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.Obliterate)
			if !*casted {
				if deathKnight.KillingMachineAura.IsActive() && !deathKnight.RimeAura.IsActive() {
					*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.FrostStrike)
				} else if deathKnight.KillingMachineAura.IsActive() && deathKnight.RimeAura.IsActive() {
					if deathKnight.CastCostPossible(sim, 0, 0, 1, 1) && deathKnight.CurrentRunicPower() < 110 {
						*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.HowlingBlast)
					} else if deathKnight.CastCostPossible(sim, 0, 0, 1, 1) && deathKnight.CurrentRunicPower() > 110 {
						*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.HowlingBlast)
					} else if !deathKnight.CastCostPossible(sim, 0, 0, 1, 1) && deathKnight.CurrentRunicPower() > 110 {
						*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.FrostStrike)
					} else if !deathKnight.CastCostPossible(sim, 0, 0, 1, 1) && deathKnight.CurrentRunicPower() < 110 {
						*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.FrostStrike)
					}
				} else if !deathKnight.KillingMachineAura.IsActive() && deathKnight.RimeAura.IsActive() {
					if deathKnight.CurrentRunicPower() < 110 {
						*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.HowlingBlast)
					} else {
						*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.FrostStrike)
					}
				} else {
					*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.FrostStrike)
					if !*casted {
						*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.HornOfWinter)
					}
				}
			}
		}
	}

	return *casted
}

func (deathKnight *DeathKnight) resetDKRotation(sim *core.Simulation) {
	deathKnight.DKRotation.opener.idx = 0
}
