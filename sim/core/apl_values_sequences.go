package core

import (
	"fmt"
	"time"

	"github.com/wowsims/classic/sim/core/proto"
)

type APLValueSequenceIsComplete struct {
	DefaultAPLValueImpl
	name     string
	sequence *APLActionSequence
}

func (rot *APLRotation) newValueSequenceIsComplete(config *proto.APLValueSequenceIsComplete) APLValue {
	if config.SequenceName == "" {
		rot.ValidationWarning("Sequence Is Complete() must provide a sequence name")
		return nil
	}
	return &APLValueSequenceIsComplete{
		name: config.SequenceName,
	}
}
func (value *APLValueSequenceIsComplete) Finalize(rot *APLRotation) {
	for _, otherAction := range rot.allAPLActions() {
		if sequence, ok := otherAction.impl.(*APLActionSequence); ok && sequence.name == value.name {
			value.sequence = sequence
			return
		}
	}
	rot.ValidationWarning("No sequence with name: '%s'", value.name)
}
func (value *APLValueSequenceIsComplete) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueSequenceIsComplete) GetBool(sim *Simulation) bool {
	return value.sequence.curIdx >= len(value.sequence.subactions)
}
func (value *APLValueSequenceIsComplete) String() string {
	return fmt.Sprintf("Sequence Is Complete(%s)", value.name)
}

type APLValueSequenceIsReady struct {
	DefaultAPLValueImpl
	name     string
	sequence *APLActionSequence
}

func (rot *APLRotation) newValueSequenceIsReady(config *proto.APLValueSequenceIsReady) APLValue {
	if config.SequenceName == "" {
		rot.ValidationWarning("Sequence Is Ready() must provide a sequence name")
		return nil
	}
	return &APLValueSequenceIsReady{
		name: config.SequenceName,
	}
}
func (value *APLValueSequenceIsReady) Finalize(rot *APLRotation) {
	for _, otherAction := range rot.allAPLActions() {
		if sequence, ok := otherAction.impl.(*APLActionSequence); ok && sequence.name == value.name {
			value.sequence = sequence
			return
		}
	}
	rot.ValidationWarning("No sequence with name: '%s'", value.name)
}
func (value *APLValueSequenceIsReady) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueSequenceIsReady) GetBool(sim *Simulation) bool {
	return value.sequence.IsReady(sim)
}
func (value *APLValueSequenceIsReady) String() string {
	return fmt.Sprintf("Sequence Is Ready(%s)", value.name)
}

type APLValueSequenceTimeToReady struct {
	DefaultAPLValueImpl
	name     string
	sequence *APLActionSequence
}

func (rot *APLRotation) newValueSequenceTimeToReady(config *proto.APLValueSequenceTimeToReady) APLValue {
	if config.SequenceName == "" {
		rot.ValidationWarning("Sequence Time To Ready() must provide a sequence name")
		return nil
	}
	return &APLValueSequenceTimeToReady{
		name: config.SequenceName,
	}
}
func (value *APLValueSequenceTimeToReady) Finalize(rot *APLRotation) {
	for _, otherAction := range rot.allAPLActions() {
		if sequence, ok := otherAction.impl.(*APLActionSequence); ok && sequence.name == value.name {
			value.sequence = sequence
			return
		}
	}
	rot.ValidationWarning("No sequence with name: '%s'", value.name)
}
func (value *APLValueSequenceTimeToReady) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueSequenceTimeToReady) GetDuration(sim *Simulation) time.Duration {
	if value.sequence.curIdx >= len(value.sequence.subactions) {
		return NeverExpires
	} else if subaction, ok := value.sequence.subactions[value.sequence.curIdx].impl.(*APLActionCastSpell); ok {
		return subaction.spell.TimeToReady(sim)
	} else if value.sequence.IsReady(sim) {
		return 0
	} else {
		return 3 * time.Second
	}
}
func (value *APLValueSequenceTimeToReady) String() string {
	return fmt.Sprintf("Sequence Time To Ready(%s)", value.name)
}
