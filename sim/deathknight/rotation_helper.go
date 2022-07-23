package deathknight

import "github.com/wowsims/wotlk/sim/core/proto"

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

	deathKnight.setupFrostRotations()
	deathKnight.setupUnholyRotations()

	// IMPORTANT
	rotationId := RotationID_Unknown
	// Also you need to update this to however you define spec
	if deathKnight.Talents.DarkConviction > 0 && deathKnight.Talents.HowlingBlast {
		rotationId = RotationID_FrostSubBlood_Full
	} else if deathKnight.Talents.BloodCakedBlade > 0 && deathKnight.Talents.HowlingBlast {
		rotationId = RotationID_FrostSubUnholy_Full
	} else if deathKnight.Talents.HowlingBlast {
		rotationId = RotationID_FrostSubBlood_Full
	} else if deathKnight.Talents.SummonGargoyle {
		if deathKnight.Rotation.UseDeathAndDecay {
			rotationId = RotationID_UnholyDnd_Full
		} else {
			if deathKnight.Rotation.ArmyOfTheDead == proto.DeathKnight_Rotation_AsMajorCd {
				if deathKnight.Rotation.UnholyPresenceOpener {
					rotationId = RotationID_UnholySsArmyUnholyPresence_Full
				} else {
					rotationId = RotationID_UnholySsArmyBloodPresence_Full
				}
			} else {
				if deathKnight.Rotation.UnholyPresenceOpener {
					rotationId = RotationID_UnholySsUnholyPresence_Full
				} else {
					rotationId = RotationID_UnholySsBloodPresence_Full
				}
			}
		}
	} else {
		rotationId = RotationID_Default
	}

	deathKnight.opener = &deathKnight.openers[rotationId]
	deathKnight.onOpener = true
}
