package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

//type RotationAction func(sim *core.Simulation, target *core.Unit)

type RotationAction func(sim *core.Simulation, target *core.Unit, s *Sequence) bool

// Add your UH rotation Actions here and then on the DoNext function

func (s *Sequence) NewAction(action RotationAction) *Sequence {
	s.actions = append(s.actions, action)
	s.numActions += 1
	return s
}

/*
const (
	RotationAction_Skip RotationActionCallback = Func_RotationAction_Skip
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
	RotationAction_SS
	RotationAction_DND
	RotationAction_GF
	RotationAction_DC
	RotationAction_Garg
	RotationAction_AOTD
	RotationAction_BP
	RotationAction_FP
	RotationAction_UP
	RotationAction_RedoSequence
	RotationAction_FS_IF_KM
)
*/

type Sequence struct {
	idx        int
	numActions int
	actions    []RotationAction
}

func (o *Sequence) IsOngoing() bool {
	return o.idx < o.numActions
}

func (o *Sequence) Reset() {
	o.idx = 0
}

func (o *Sequence) Advance(condition bool) {
	o.idx += 1
}

func (o *Sequence) ConditionalAdvance(condition bool) {
	if condition {
		o.idx += 1
	}
}

func (o *Sequence) GetNextAction() RotationAction {
	if o.idx+1 < o.numActions {
		return o.actions[o.idx+1]
	} else {
		return nil
	}
}

type RotationHelper struct {
	Opener *Sequence
	Main   *Sequence
}

func TernaryRotationAction(condition bool, t RotationAction, f RotationAction) RotationAction {
	if condition {
		return t
	} else {
		return f
	}
}
