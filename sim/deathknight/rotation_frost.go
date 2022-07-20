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
	numTargets int
	targets    []*core.Unit

	opener  *Opener
	openers []Opener
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

func (deathKnight *DeathKnight) setupTargets() {
	r := &deathKnight.DKRotation
	r.numTargets = int(deathKnight.Env.GetNumTargets())
	r.targets = make([]*core.Unit, r.numTargets)
	for i := 0; i < r.numTargets; i++ {
		r.targets[i] = deathKnight.Env.GetTargetUnit(int32(i))
	}
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
	deathKnight.setupTargets()
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
			}
		}

	}
}

func (o *Opener) DoNext(sim *core.Simulation, deathKnight *DeathKnight) bool {
	target := deathKnight.CurrentTarget

	action := o.actions[o.idx]

	castSuccessful := false
	switch action {
	case OpenerAction_IT:
		castSuccessful = deathKnight.CastIcyTouch(sim, target)
	case OpenerAction_PS:
		castSuccessful = deathKnight.CastPlagueStrike(sim, target)
	}

	return castSuccessful
}

func (deathKnight *DeathKnight) resetDKRotation(sim *core.Simulation) {
	deathKnight.DKRotation.opener.idx = 0
}
