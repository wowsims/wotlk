package deathknight

import (
	"container/list"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

type DKRotationAction uint8

const (
	DKRotationAction_Skip DKRotationAction = iota
	DKRotationAction_IT
	DKRotationAction_PS
	DKRotationAction_Obli
	DKRotationAction_BS
	DKRotationAction_BT
	DKRotationAction_UA
	DKRotationAction_Pesti
	DKRotationAction_FS
	DKRotationAction_HW
	DKRotationAction_ERW
)

type DKRotationSequence struct {
	idx        int
	numActions int
	repeatable bool
	actions    []DKRotationAction
}

type DKRotation struct {
	numTargets int
	targets    []*core.Unit

	currSequence     *list.Element
	mainSequenceList *list.List
}

func (deathKnight *DeathKnight) getIndexForTarget(t *core.Unit) int {
	r := &deathKnight.DKRotation
	idx := -1
	for i := 0; i < r.numTargets; i++ {
		if t == r.targets[i] {
			idx = i
			break
		}
	}
	if idx == -1 {
		panic("This cannot happen!")
	}
	return idx
}

func (rs *DKRotationSequence) resetSequence() {
	rs.idx = 0
}

func TernaryRotationAction(condition bool, t DKRotationAction, f DKRotationAction) DKRotationAction {
	if condition {
		return t
	} else {
		return f
	}
}

func initSequence(repeatable bool, actions []DKRotationAction) *DKRotationSequence {
	var seq DKRotationSequence
	seq.idx = 0
	seq.numActions = len(actions)
	seq.actions = actions
	seq.repeatable = repeatable
	return &seq
}

func (deathKnight *DeathKnight) setupTargets() {
	r := &deathKnight.DKRotation
	r.numTargets = int(deathKnight.Env.GetNumTargets())
	r.targets = make([]*core.Unit, r.numTargets)
	for i := 0; i < r.numTargets; i++ {
		r.targets[i] = deathKnight.Env.GetTargetUnit(int32(i))
	}
}

func (deathKnight *DeathKnight) setupDKRotation() {
	r := &deathKnight.DKRotation
	deathKnight.setupTargets()

	mainSequence := initSequence(false, []DKRotationAction{
		TernaryRotationAction(deathKnight.Options.RefreshHornOfWinter, DKRotationAction_HW, DKRotationAction_Skip),
		DKRotationAction_IT,
		DKRotationAction_PS,
		DKRotationAction_BT,
		DKRotationAction_Pesti,
		DKRotationAction_Obli,
		DKRotationAction_FS,
		DKRotationAction_ERW,
		DKRotationAction_Obli,
		DKRotationAction_Obli,
		DKRotationAction_Obli,
		DKRotationAction_FS,
		DKRotationAction_FS,
		DKRotationAction_FS,
		DKRotationAction_Obli,
		DKRotationAction_Obli,
		DKRotationAction_BS,
		DKRotationAction_Pesti,
		DKRotationAction_FS,
	})

	constSequence := initSequence(true, []DKRotationAction{
		DKRotationAction_Obli,
		DKRotationAction_Obli,
		DKRotationAction_Pesti,
		DKRotationAction_BS,
		DKRotationAction_FS,
		DKRotationAction_FS,
		DKRotationAction_Obli,
		DKRotationAction_Obli,
		DKRotationAction_FS,
		DKRotationAction_Obli,
		DKRotationAction_FS,
		DKRotationAction_FS,
	})

	r.mainSequenceList = list.New()
	r.mainSequenceList.PushBack(mainSequence)
	r.mainSequenceList.PushBack(constSequence)

	r.currSequence = r.mainSequenceList.Front()
}

func (deathKnight *DeathKnight) nextDKRotationSequenceAction() DKRotationAction {
	seq := deathKnight.DKRotation.currSequence.Value.(*DKRotationSequence)
	return seq.actions[seq.idx]
}

func (deathKnight *DeathKnight) advanceDKRotationSequenceAction() bool {
	seq := deathKnight.DKRotation.currSequence.Value.(*DKRotationSequence)
	if seq.idx+1 >= seq.numActions {
		return true
	} else {
		seq.idx += 1
		return false
	}
}

func (deathKnight *DeathKnight) advanceDKRotation() {
	if deathKnight.advanceDKRotationSequenceAction() {
		r := &deathKnight.DKRotation
		seq := r.currSequence.Value.(*DKRotationSequence)
		if !seq.repeatable {
			r.currSequence = r.currSequence.Next()
		} else {
			seq.resetSequence()
		}
	}
}

func (deathKnight *DeathKnight) doDKRotation(sim *core.Simulation, allowAdvance bool) bool {
	if !deathKnight.Talents.HowlingBlast {
		return false
	}

	target := deathKnight.CurrentTarget

	const USE_BAD_ROTA = true

	if USE_BAD_ROTA {
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
			} else if deathKnight.CanUnbreakableArmor(sim) && deathKnight.AllDiseasesAreActive(target) {
				deathKnight.UnbreakableArmor.Cast(sim, target)
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
					nextSwing := deathKnight.AutoAttacks.MainhandSwingAt
					if deathKnight.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
						nextSwing = core.MinDuration(nextSwing, deathKnight.AutoAttacks.OffhandSwingAt)
					}
					deathKnight.WaitUntil(sim, nextSwing)
				}
			}
		}

		return false
	} else {
		nextFRAction := deathKnight.nextDKRotationSequenceAction()
		casted := false
		// TODO: Check for hits on main disease appliers && have a prio bracket for when we're "lost"
		switch nextFRAction {
		case DKRotationAction_Skip:
			deathKnight.advanceDKRotation()
			casted = deathKnight.doDKRotation(sim, false)
		case DKRotationAction_IT:
			if deathKnight.CanIcyTouch(sim) {
				deathKnight.IcyTouch.Cast(sim, target)
				casted = true
			}
		case DKRotationAction_PS:
			if deathKnight.CanPlagueStrike(sim) {
				deathKnight.PlagueStrike.Cast(sim, target)
				casted = true
			}
		case DKRotationAction_Obli:
			if deathKnight.CanObliterate(sim) {
				deathKnight.Obliterate.Cast(sim, target)
				casted = true
			}
		case DKRotationAction_BS:
			if deathKnight.CanBloodStrike(sim) {
				deathKnight.BloodStrike.Cast(sim, target)
				casted = true
			}
		case DKRotationAction_BT:
			if deathKnight.CanBloodTap(sim) {
				deathKnight.BloodTap.Cast(sim, target)
				casted = true
			}
		case DKRotationAction_UA:
			if deathKnight.CanUnbreakableArmor(sim) {
				deathKnight.UnbreakableArmor.Cast(sim, target)
				casted = true
			}
		case DKRotationAction_Pesti:
			if deathKnight.CanPestilence(sim) {
				deathKnight.Pestilence.Cast(sim, target)
				casted = true
			}
		case DKRotationAction_FS:
			if deathKnight.CanFrostStrike(sim) {
				deathKnight.FrostStrike.Cast(sim, target)
				casted = true
			}
		case DKRotationAction_HW:
			if deathKnight.CanHornOfWinter(sim) {
				deathKnight.HornOfWinter.Cast(sim, target)
				casted = true
			}
		case DKRotationAction_ERW:
			if deathKnight.CanEmpowerRuneWeapon(sim) {
				deathKnight.EmpowerRuneWeapon.Cast(sim, target)
				casted = true
			}
		}

		if !casted {
			// TODO: Prio stuff.
			//if deathKnight.RaiseDead.IsReady(sim) {
			//	deathKnight.RaiseDead.Cast(sim, target)
			//}

			if deathKnight.GCD.IsReady(sim) && !deathKnight.IsWaiting() {
				// Wait until 1/10th of the swing
				nextSwing := deathKnight.AutoAttacks.MainhandSwingAt
				if deathKnight.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
					nextSwing = core.MinDuration(nextSwing, deathKnight.AutoAttacks.OffhandSwingAt)
				}
				deathKnight.WaitUntil(sim, time.Duration(0.1*float64(nextSwing-sim.CurrentTime))+sim.CurrentTime)
			}
		} else {
			if !(deathKnight.LastCastOutcome == core.OutcomeMiss &&
				(nextFRAction == DKRotationAction_IT ||
					nextFRAction == DKRotationAction_PS ||
					nextFRAction == DKRotationAction_Pesti)) && allowAdvance {
				deathKnight.advanceDKRotation()
			}
		}

		return casted
	}
}

func (deathKnight *DeathKnight) resetDKRotation(sim *core.Simulation) {
	r := &deathKnight.DKRotation

	for e := r.mainSequenceList.Front(); e != nil; e = e.Next() {
		seq := e.Value.(*DKRotationSequence)
		seq.resetSequence()
	}

	r.currSequence = r.mainSequenceList.Front()
}
