package core

import (
	"fmt"
	"time"

	"github.com/wowsims/classic/sim/core/proto"
)

type APLValueCurrentTime struct {
	DefaultAPLValueImpl
}

func (rot *APLRotation) newValueCurrentTime(config *proto.APLValueCurrentTime) APLValue {
	return &APLValueCurrentTime{}
}
func (value *APLValueCurrentTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueCurrentTime) GetDuration(sim *Simulation) time.Duration {
	return sim.CurrentTime
}
func (value *APLValueCurrentTime) String() string {
	return "Current Time"
}

type APLValueCurrentTimePercent struct {
	DefaultAPLValueImpl
}

func (rot *APLRotation) newValueCurrentTimePercent(config *proto.APLValueCurrentTimePercent) APLValue {
	return &APLValueCurrentTimePercent{}
}
func (value *APLValueCurrentTimePercent) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentTimePercent) GetFloat(sim *Simulation) float64 {
	return sim.CurrentTime.Seconds() / sim.Duration.Seconds()
}
func (value *APLValueCurrentTimePercent) String() string {
	return fmt.Sprintf("Current Time %%")
}

type APLValueRemainingTime struct {
	DefaultAPLValueImpl
}

func (rot *APLRotation) newValueRemainingTime(config *proto.APLValueRemainingTime) APLValue {
	return &APLValueRemainingTime{}
}
func (value *APLValueRemainingTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueRemainingTime) GetDuration(sim *Simulation) time.Duration {
	return sim.GetRemainingDuration()
}
func (value *APLValueRemainingTime) String() string {
	return "Remaining Time"
}

type APLValueRemainingTimePercent struct {
	DefaultAPLValueImpl
}

func (rot *APLRotation) newValueRemainingTimePercent(config *proto.APLValueRemainingTimePercent) APLValue {
	return &APLValueRemainingTimePercent{}
}
func (value *APLValueRemainingTimePercent) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueRemainingTimePercent) GetFloat(sim *Simulation) float64 {
	return sim.GetRemainingDurationPercent()
}
func (value *APLValueRemainingTimePercent) String() string {
	return fmt.Sprintf("Remaining Time %%")
}

type APLValueNumberTargets struct {
	DefaultAPLValueImpl
}

func (rot *APLRotation) newValueNumberTargets(config *proto.APLValueNumberTargets) APLValue {
	return &APLValueNumberTargets{}
}
func (value *APLValueNumberTargets) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueNumberTargets) GetInt(sim *Simulation) int32 {
	return sim.GetNumTargets()
}
func (value *APLValueNumberTargets) String() string {
	return "Num Targets"
}

type APLValueIsExecutePhase struct {
	DefaultAPLValueImpl
	threshold proto.APLValueIsExecutePhase_ExecutePhaseThreshold
}

func (rot *APLRotation) newValueIsExecutePhase(config *proto.APLValueIsExecutePhase) APLValue {
	if config.Threshold == proto.APLValueIsExecutePhase_Unknown {
		return nil
	}
	return &APLValueIsExecutePhase{
		threshold: config.Threshold,
	}
}
func (value *APLValueIsExecutePhase) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueIsExecutePhase) GetBool(sim *Simulation) bool {
	if value.threshold == proto.APLValueIsExecutePhase_E20 {
		return sim.IsExecutePhase20()
	} else if value.threshold == proto.APLValueIsExecutePhase_E25 {
		return sim.IsExecutePhase25()
	} else if value.threshold == proto.APLValueIsExecutePhase_E35 {
		return sim.IsExecutePhase35()
	} else {
		panic("Should never reach here")
	}
}
func (value *APLValueIsExecutePhase) String() string {
	return "Is Execute Phase"
}
