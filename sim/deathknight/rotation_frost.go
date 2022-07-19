package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

type FrostRotationAction uint8

const (
	FrostRotationAction_Wait FrostRotationAction = iota
	FrostRotationAction_IT
	FrostRotationAction_PS
	FrostRotationAction_Obli
	FrostRotationAction_BS
	FrostRotationAction_BT
	FrostRotationAction_UA
	FrostRotationAction_Pesti
	FrostRotationAction_FS
	FrostRotationAction_HW
	FrostRotationAction_ERW
)

type FrostRotationSequence struct {
	building   bool
	idx        int
	numActions int
	actions    []FrostRotationAction
}

func (frs *FrostRotationSequence) beginBuildingSequence() {
	if frs.building {
		panic("Started building inside sequence!")
	}
	frs.idx = 0
	frs.building = true
}

func (frs *FrostRotationSequence) endBuildingSequence() {
	if !frs.building {
		panic("Ended building without a start!")
	}
	frs.idx = 0
	frs.building = false
}

type FrostRotation struct {
	numTargets int
	targets    []*core.Unit

	currSequence *FrostRotationSequence

	mainSequence  FrostRotationSequence
	constSequence FrostRotationSequence
}

func (deathKnight *DeathKnight) getIndexForTarget(t *core.Unit) int {
	fr := &deathKnight.FrostRotation
	idx := -1
	for i := 0; i < fr.numTargets; i++ {
		if t == fr.targets[i] {
			idx = i
			break
		}
	}
	if idx == -1 {
		panic("This cannot happen!")
	}
	return idx
}

func (frs *FrostRotationSequence) resetFrostRotationSequence() {
	frs.idx = 0
}

func (frs *FrostRotationSequence) addToFrostRotationSequence(action FrostRotationAction) {
	frs.actions[frs.idx] = action
	frs.idx += 1
	frs.numActions += 1
}

func (deathKnight *DeathKnight) setCurrentFrostRotationSequence(frs *FrostRotationSequence) {
	deathKnight.FrostRotation.currSequence = frs
}

func (deathKnight *DeathKnight) advanceFrostRotationSequence() {
	deathKnight.FrostRotation.currSequence.idx += 1
	if deathKnight.FrostRotation.currSequence.idx == deathKnight.FrostRotation.currSequence.numActions {
		if deathKnight.FrostRotation.currSequence == &deathKnight.FrostRotation.mainSequence {
			deathKnight.FrostRotation.currSequence.idx = 0
			deathKnight.FrostRotation.currSequence = &deathKnight.FrostRotation.constSequence
		} else {
			deathKnight.FrostRotation.currSequence.idx = 0
		}
	}
}

func (deathKnight *DeathKnight) setupFrostRotation() {
	fr := &deathKnight.FrostRotation
	fr.numTargets = int(deathKnight.Env.GetNumTargets())
	fr.targets = make([]*core.Unit, fr.numTargets)
	// TODO: make this nicer
	fr.mainSequence.actions = make([]FrostRotationAction, 64)
	fr.constSequence.actions = make([]FrostRotationAction, 64)
	for i := 0; i < fr.numTargets; i++ {
		fr.targets[i] = deathKnight.Env.GetTargetUnit(int32(i))
	}

	fr.mainSequence.beginBuildingSequence()
	if deathKnight.Rotation.RefreshHornOfWinter {
		fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_HW)
	}
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_IT)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_PS)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_BT)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_Pesti)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_Obli)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_FS)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_ERW)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_Obli)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_Obli)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_Obli)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_FS)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_FS)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_FS)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_Obli)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_Obli)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_BS)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_Pesti)
	fr.mainSequence.addToFrostRotationSequence(FrostRotationAction_FS)
	fr.mainSequence.endBuildingSequence()

	fr.constSequence.beginBuildingSequence()
	fr.constSequence.addToFrostRotationSequence(FrostRotationAction_Obli)
	fr.constSequence.addToFrostRotationSequence(FrostRotationAction_Obli)
	fr.constSequence.addToFrostRotationSequence(FrostRotationAction_Pesti)
	fr.constSequence.addToFrostRotationSequence(FrostRotationAction_BS)
	fr.constSequence.addToFrostRotationSequence(FrostRotationAction_FS)
	fr.constSequence.addToFrostRotationSequence(FrostRotationAction_FS)
	fr.constSequence.addToFrostRotationSequence(FrostRotationAction_Obli)
	fr.constSequence.addToFrostRotationSequence(FrostRotationAction_Obli)
	fr.constSequence.addToFrostRotationSequence(FrostRotationAction_FS)
	fr.constSequence.addToFrostRotationSequence(FrostRotationAction_Obli)
	fr.constSequence.addToFrostRotationSequence(FrostRotationAction_FS)
	fr.constSequence.addToFrostRotationSequence(FrostRotationAction_FS)
	fr.constSequence.endBuildingSequence()

	deathKnight.setCurrentFrostRotationSequence(&fr.mainSequence)
}

func (deathKnight *DeathKnight) nextFrostRotationSequenceAction() FrostRotationAction {
	return deathKnight.FrostRotation.currSequence.actions[deathKnight.FrostRotation.currSequence.idx]
}

func (deathKnight *DeathKnight) doFrostRotation(sim *core.Simulation) {
	if !deathKnight.Talents.HowlingBlast {
		return
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
	} else {
		nextFRAction := deathKnight.nextFrostRotationSequenceAction()
		casted := false
		// TODO: Check for hits on main disease appliers && have a prio bracket for when we're "lost"
		switch nextFRAction {
		case FrostRotationAction_Wait:
			// TODO:
		case FrostRotationAction_IT:
			if deathKnight.CanIcyTouch(sim) {
				deathKnight.IcyTouch.Cast(sim, target)
				casted = true
			}
		case FrostRotationAction_PS:
			if deathKnight.CanPlagueStrike(sim) {
				deathKnight.PlagueStrike.Cast(sim, target)
				casted = true
			}
		case FrostRotationAction_Obli:
			if deathKnight.CanObliterate(sim) {
				deathKnight.Obliterate.Cast(sim, target)
				casted = true
			}
		case FrostRotationAction_BS:
			if deathKnight.CanBloodStrike(sim) {
				deathKnight.BloodStrike.Cast(sim, target)
				casted = true
			}
		case FrostRotationAction_BT:
			if deathKnight.CanBloodTap(sim) {
				deathKnight.BloodTap.Cast(sim, target)
				casted = true
			}
		case FrostRotationAction_UA:
			if deathKnight.CanUnbreakableArmor(sim) {
				deathKnight.UnbreakableArmor.Cast(sim, target)
				casted = true
			}
		case FrostRotationAction_Pesti:
			if deathKnight.CanPestilence(sim) {
				deathKnight.Pestilence.Cast(sim, target)
				casted = true
			}
		case FrostRotationAction_FS:
			if deathKnight.CanFrostStrike(sim) {
				deathKnight.FrostStrike.Cast(sim, target)
				casted = true
			}
		case FrostRotationAction_HW:
			if deathKnight.CanHornOfWinter(sim) {
				deathKnight.HornOfWinter.Cast(sim, target)
				casted = true
			}
		case FrostRotationAction_ERW:
			if deathKnight.CanEmpowerRuneWeapon(sim) {
				deathKnight.EmpowerRuneWeapon.Cast(sim, target)
				casted = true
			}
		}

		if !casted {
			if deathKnight.GCD.IsReady(sim) && !deathKnight.IsWaiting() {
				nextSwing := deathKnight.AutoAttacks.MainhandSwingAt
				if deathKnight.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
					nextSwing = core.MinDuration(nextSwing, deathKnight.AutoAttacks.OffhandSwingAt)
				}
				deathKnight.WaitUntil(sim, nextSwing)
			}
		} else {
			deathKnight.advanceFrostRotationSequence()
		}
	}
}

func (deathKnight *DeathKnight) resetFrostRotation(sim *core.Simulation) {
	deathKnight.FrostRotation.mainSequence.resetFrostRotationSequence()
	deathKnight.FrostRotation.constSequence.resetFrostRotationSequence()
}
