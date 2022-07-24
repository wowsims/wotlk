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

type RotationID uint8

const (
	RotationID_Default RotationID = iota
	RotationID_FrostSubBlood_Full
	RotationID_FrostSubUnholy_Full

	RotationID_UnholySsUnholyPresence_Full
	RotationID_UnholySsArmyUnholyPresence_Full
	RotationID_UnholySsBloodPresence_Full
	RotationID_UnholySsArmyBloodPresence_Full
	RotationID_UnholyDnd_Full
	RotationID_Count
	RotationID_Unknown
)

type Sequence struct {
	id         RotationID
	idx        int
	numActions int
	actions    []RotationAction
}

type SetupRotationEvent func() RotationID
type DoRotationEvent func(sim *core.Simulation, target *core.Unit)

type RotationHelper struct {
	opener   *Sequence
	openers  []Sequence
	onOpener bool

	sequence *Sequence

	CastSuccessful     bool
	justCastPestilence bool

	SetupRotationEvent SetupRotationEvent
	DoRotationEvent    DoRotationEvent
}

func (dk *Deathknight) GetRotationId() RotationID {
	return dk.opener.id
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

func (dk *Deathknight) SetupRotation() {
	dk.openers = make([]Sequence, RotationID_Count)

	rotationId := RotationID_Unknown
	if dk.SetupRotationEvent != nil {
		rotationId = dk.SetupRotationEvent()
	} else {
		panic("Missing SetupRotationEvent. Please assign during spec creation")
	}

	dk.opener = &dk.openers[rotationId]
	dk.onOpener = true
}
