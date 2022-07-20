package deathknight

import (
	"container/list"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

type DKRotationAction uint8

const (
	DKRotationAction_Skip DKRotationAction = iota
	DKRotationAction_EnableDiseaseCheck
	DKRotationAction_DisableDiseaseCheck
	DKRotationAction_ReapplyDiseases
	DKRotationAction_IT
	DKRotationAction_PS
	DKRotationAction_Obli
	DKRotationAction_BS
	DKRotationAction_BT
	DKRotationAction_UA
	DKRotationAction_RD
	DKRotationAction_Pesti
	DKRotationAction_FS
	DKRotationAction_HW
	DKRotationAction_ERW
	DKRotationAction_HB_Ghoul_FS_RimeCheck
	DKRotationAction_PrioMode
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

	currSequence      *list.Element
	mainSequenceList  *list.List
	bloodSequenceList *list.List

	lastFFApplication   time.Duration
	lastBPApplication   time.Duration
	diseaseCheckEnabled bool
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
		DKRotationAction_IT,
		DKRotationAction_PS,
		DKRotationAction_EnableDiseaseCheck,
		DKRotationAction_UA,
		DKRotationAction_BT,
		DKRotationAction_Obli,
		DKRotationAction_FS,
		DKRotationAction_Pesti,
		DKRotationAction_ERW,
		DKRotationAction_Obli,
		DKRotationAction_Obli,
		DKRotationAction_Obli,
		DKRotationAction_FS,
		DKRotationAction_HB_Ghoul_FS_RimeCheck,
		DKRotationAction_FS,
		DKRotationAction_Obli,
		DKRotationAction_Obli,
		DKRotationAction_Pesti,
		DKRotationAction_FS,
		DKRotationAction_BS,
		DKRotationAction_FS,
	})

	constSequence := initSequence(true, []DKRotationAction{
		DKRotationAction_Obli,
		DKRotationAction_FS,
		DKRotationAction_Obli,
		DKRotationAction_FS,
		DKRotationAction_BS,
		DKRotationAction_FS,
		DKRotationAction_Pesti,
		DKRotationAction_FS,
		DKRotationAction_PrioMode,
	})

	r.mainSequenceList = list.New()
	r.mainSequenceList.PushBack(mainSequence)
	r.mainSequenceList.PushBack(constSequence)

	r.currSequence = r.mainSequenceList.Front()

	r.diseaseCheckEnabled = false
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

func (seq *DKRotationSequence) doNothing() {
	return
}

func (deathKnight *DeathKnight) doDKRotation(sim *core.Simulation, advance bool) bool {
	if !deathKnight.Talents.HowlingBlast {
		return false
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

		return false
	} else {
		seq := deathKnight.DKRotation.currSequence.Value.(*DKRotationSequence)
		seq.doNothing()
		nextAction := deathKnight.nextDKRotationSequenceAction()
		casted := false
		skip := false

		if nextAction == DKRotationAction_EnableDiseaseCheck {
			if deathKnight.DKRotation.diseaseCheckEnabled {
				panic("Enable disease check inside another enable disease check.")
			}
			deathKnight.DKRotation.diseaseCheckEnabled = true
			deathKnight.advanceDKRotation()
			skip = true

		} else if nextAction == DKRotationAction_DisableDiseaseCheck {
			if !deathKnight.DKRotation.diseaseCheckEnabled {
				panic("Disable disease check while not enabled.")
			}
			deathKnight.DKRotation.diseaseCheckEnabled = false
			deathKnight.advanceDKRotation()
			skip = true
		}

		if deathKnight.DKRotation.diseaseCheckEnabled && (!deathKnight.TargetHasDisease(FrostFeverAuraLabel, target) || !deathKnight.TargetHasDisease(BloodPlagueAuraLabel, target)) {
			nextAction = DKRotationAction_ReapplyDiseases
		}

		if deathKnight.DKRotation.diseaseCheckEnabled && nextAction != DKRotationAction_ReapplyDiseases {
			if !deathKnight.TargetHasDisease(FrostFeverAuraLabel, target) || deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim) < 2*time.Second ||
				!deathKnight.TargetHasDisease(BloodPlagueAuraLabel, target) || deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim) < 2*time.Second {
				if deathKnight.CanPestilence(sim) {
					deathKnight.Pestilence.Cast(sim, target)
					deathKnight.DKRotation.lastFFApplication = sim.CurrentTime
					deathKnight.DKRotation.lastBPApplication = sim.CurrentTime
					skip = true
				} else {

				}
			}
		}

		if !skip {
			// TODO: Check for hits on main disease appliers && have a prio bracket for when we're "lost"
			switch nextAction {
			case DKRotationAction_Skip:
				deathKnight.advanceDKRotation()
				casted = deathKnight.doDKRotation(sim, false)
			case DKRotationAction_IT:
				if deathKnight.CanIcyTouch(sim) {
					deathKnight.IcyTouch.Cast(sim, target)
					if deathKnight.LastCastOutcome != core.OutcomeMiss {
						deathKnight.DKRotation.lastFFApplication = sim.CurrentTime
					}
					casted = true
				}
			case DKRotationAction_PS:
				if deathKnight.CanPlagueStrike(sim) {
					deathKnight.PlagueStrike.Cast(sim, target)
					if deathKnight.LastCastOutcome != core.OutcomeMiss {
						deathKnight.DKRotation.lastBPApplication = sim.CurrentTime
					}
					casted = true
				}
			case DKRotationAction_ReapplyDiseases:
				if !deathKnight.TargetHasDisease(FrostFeverAuraLabel, target) {
					if deathKnight.CanIcyTouch(sim) {
						deathKnight.IcyTouch.Cast(sim, target)
						if deathKnight.LastCastOutcome != core.OutcomeMiss {
							deathKnight.DKRotation.lastFFApplication = sim.CurrentTime
						}
						casted = true
						advance = false
					}
				} else if !deathKnight.TargetHasDisease(BloodPlagueAuraLabel, target) {
					if deathKnight.CanPlagueStrike(sim) {
						deathKnight.PlagueStrike.Cast(sim, target)
						if deathKnight.LastCastOutcome != core.OutcomeMiss {
							deathKnight.DKRotation.lastBPApplication = sim.CurrentTime
						}
						casted = true
						advance = false
					}
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
					deathKnight.WaitUntil(sim, sim.CurrentTime+1)
					casted = true
				}
			case DKRotationAction_UA:
				if deathKnight.CanUnbreakableArmor(sim) {
					deathKnight.UnbreakableArmor.Cast(sim, target)
					deathKnight.WaitUntil(sim, sim.CurrentTime+1)
					casted = true
				}
			case DKRotationAction_Pesti:
				if deathKnight.CanPestilence(sim) && (sim.CurrentTime-deathKnight.DKRotation.lastFFApplication > 3*time.Second || sim.CurrentTime-deathKnight.DKRotation.lastBPApplication > 3*time.Second) {
					deathKnight.Pestilence.Cast(sim, target)
					casted = true
				} else {
					deathKnight.advanceDKRotation()
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
					deathKnight.WaitUntil(sim, sim.CurrentTime+1)
					casted = true
				}
			case DKRotationAction_RD:
				if deathKnight.CanRaiseDead(sim) {
					deathKnight.RaiseDead.Cast(sim, target)
					casted = true
				}
			case DKRotationAction_HB_Ghoul_FS_RimeCheck:
				if deathKnight.RimeAura.IsActive() {
					if deathKnight.CanHowlingBlast(sim) {
						deathKnight.HowlingBlast.Cast(sim, target)
						casted = true
					}
				} else {
					if deathKnight.CanRaiseDead(sim) {
						deathKnight.RaiseDead.Cast(sim, target)
						casted = true
					} else if deathKnight.CanFrostStrike(sim) {
						deathKnight.FrostStrike.Cast(sim, target)
						casted = true
					}
				}
			case DKRotationAction_PrioMode:

				casted = true
			}

			if !casted {
				// TODO: Prio stuff.
				if deathKnight.CanObliterate(sim) {
					deathKnight.Obliterate.Cast(sim, target)
				} else if deathKnight.KillingMachineAura.IsActive() {
					if deathKnight.CastCostPossible(sim, 0, 0, 1, 1) && deathKnight.RimeAura.IsActive() {
						if deathKnight.CurrentRunicPower() < 110 {
							if deathKnight.CanHowlingBlast(sim) {
								deathKnight.HowlingBlast.Cast(sim, target)
							}
						} else {
							if deathKnight.CanFrostStrike(sim) {
								deathKnight.FrostStrike.Cast(sim, target)
							}
						}
					} else {
						if deathKnight.CanFrostStrike(sim) {
							deathKnight.FrostStrike.Cast(sim, target)
						}
					}
				} else if deathKnight.CanFrostStrike(sim) {
					deathKnight.FrostStrike.Cast(sim, target)
				} else if deathKnight.CanHornOfWinter(sim) {
					deathKnight.HornOfWinter.Cast(sim, target)

				}

				if deathKnight.GCD.IsReady(sim) && !deathKnight.IsWaiting() {
					// Wait until 1/10th of the swing
					waitUntil := deathKnight.AutoAttacks.MainhandSwingAt
					if deathKnight.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
						waitUntil = core.MinDuration(waitUntil, deathKnight.AutoAttacks.OffhandSwingAt)
					}
					waitUntil = core.MinDuration(time.Duration(0.1*float64(waitUntil)), deathKnight.AnyRuneReadyAt(sim))
					deathKnight.WaitUntil(sim, waitUntil)
				}
			} else {
				if !(deathKnight.LastCastOutcome == core.OutcomeMiss &&
					(nextAction == DKRotationAction_IT ||
						nextAction == DKRotationAction_PS ||
						nextAction == DKRotationAction_Pesti)) && advance {
					deathKnight.advanceDKRotation()
				}
			}
		} else {
			if deathKnight.GCD.IsReady(sim) && !deathKnight.IsWaiting() {
				// Wait until 1/10th of the swing
				waitUntil := deathKnight.AutoAttacks.MainhandSwingAt
				if deathKnight.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
					waitUntil = core.MinDuration(waitUntil, deathKnight.AutoAttacks.OffhandSwingAt)
				}
				waitUntil = core.MinDuration(time.Duration(0.1*float64(waitUntil-sim.CurrentTime)+float64(waitUntil)), deathKnight.AnyRuneReadyAt(sim))
				deathKnight.WaitUntil(sim, waitUntil)
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

	r.diseaseCheckEnabled = false
}
