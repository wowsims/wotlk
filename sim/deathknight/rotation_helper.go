package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

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
	RotationAction_SS
	RotationAction_DND
	RotationAction_GF
	RotationAction_DC
	RotationAction_Garg
	RotationAction_AOTD
	RotationAction_BP
	RotationAction_FP
	RotationAction_UP
)

type Sequence struct {
	idx        int
	numActions int
	actions    []RotationAction
}

type DoRotationEvent func(sim *core.Simulation, target *core.Unit)

type RotationHelper struct {
	opener   *Sequence
	onOpener bool

	CastSuccessful     bool
	justCastPestilence bool

	DoRotationEvent DoRotationEvent
}

func TernaryRotationAction(condition bool, t RotationAction, f RotationAction) RotationAction {
	if condition {
		return t
	} else {
		return f
	}
}

func (r *RotationHelper) DefineOpener(actions []RotationAction) {
	r.opener = &Sequence{
		idx:        0,
		numActions: len(actions),
		actions:    actions,
	}
}
