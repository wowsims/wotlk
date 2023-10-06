package core

import (
	"fmt"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValueCurrentRuneCount struct {
	DefaultAPLValueImpl
	unit     *Unit
	runeType proto.APLValueRuneType
}

func (rot *APLRotation) newValueCurrentRuneCount(config *proto.APLValueCurrentRuneCount) APLValue {
	unit := rot.unit
	if !unit.HasRunicPowerBar() {
		rot.ValidationWarning("%s does not use Runes", unit.Label)
		return nil
	}
	return &APLValueCurrentRuneCount{
		unit:     unit,
		runeType: config.RuneType,
	}
}
func (value *APLValueCurrentRuneCount) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueCurrentRuneCount) GetInt(_ *Simulation) int32 {
	switch value.runeType {
	case proto.APLValueRuneType_RuneBlood:
		return int32(value.unit.CurrentBloodOrDeathRunes())
	case proto.APLValueRuneType_RuneFrost:
		return int32(value.unit.CurrentFrostOrDeathRunes())
	case proto.APLValueRuneType_RuneUnholy:
		return int32(value.unit.CurrentUnholyOrDeathRunes())
	case proto.APLValueRuneType_RuneDeath:
		return int32(value.unit.CurrentDeathRunes())
	}
	return 0
}
func (value *APLValueCurrentRuneCount) String() string {
	return fmt.Sprintf("Current Rune Count(%s)", value.runeType)
}

type APLValueCurrentNonDeathRuneCount struct {
	DefaultAPLValueImpl
	unit     *Unit
	runeType proto.APLValueRuneType
}

func (rot *APLRotation) newValueCurrentNonDeathRuneCount(config *proto.APLValueCurrentNonDeathRuneCount) APLValue {
	unit := rot.unit
	if !unit.HasRunicPowerBar() {
		rot.ValidationWarning("%s does not use Runes", unit.Label)
		return nil
	}
	return &APLValueCurrentNonDeathRuneCount{
		unit:     unit,
		runeType: config.RuneType,
	}
}
func (value *APLValueCurrentNonDeathRuneCount) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueCurrentNonDeathRuneCount) GetInt(_ *Simulation) int32 {
	switch value.runeType {
	case proto.APLValueRuneType_RuneBlood:
		return int32(value.unit.CurrentBloodRunes())
	case proto.APLValueRuneType_RuneFrost:
		return int32(value.unit.CurrentFrostRunes())
	case proto.APLValueRuneType_RuneUnholy:
		return int32(value.unit.CurrentUnholyRunes())
	}
	return 0
}
func (value *APLValueCurrentNonDeathRuneCount) String() string {
	return fmt.Sprintf("Current Non-Death Rune Count(%s)", value.runeType)
}

type APLValueCurrentRuneActive struct {
	DefaultAPLValueImpl
	unit     *Unit
	runeSlot int8
}

func (rot *APLRotation) newValueCurrentRuneActive(config *proto.APLValueCurrentRuneActive) APLValue {
	unit := rot.unit
	if !unit.HasRunicPowerBar() {
		rot.ValidationWarning("%s does not use Runes", unit.Label)
		return nil
	}
	return &APLValueCurrentRuneActive{
		unit:     unit,
		runeSlot: int8(config.RuneSlot) - 1, // 0 is Unknown
	}
}
func (value *APLValueCurrentRuneActive) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueCurrentRuneActive) GetBool(_ *Simulation) bool {
	return value.unit.RuneIsActive(value.runeSlot)
}
func (value *APLValueCurrentRuneActive) String() string {
	return fmt.Sprintf("Current Rune Active(%d)", value.runeSlot)
}

type APLValueCurrentRuneDeath struct {
	DefaultAPLValueImpl
	unit     *Unit
	runeSlot int8
}

func (rot *APLRotation) newValueCurrentRuneDeath(config *proto.APLValueCurrentRuneDeath) APLValue {
	unit := rot.unit
	if !unit.HasRunicPowerBar() {
		rot.ValidationWarning("%s does not use Runes", unit.Label)
		return nil
	}
	return &APLValueCurrentRuneDeath{
		unit:     unit,
		runeSlot: int8(config.RuneSlot) - 1, // 0 is Unknown
	}
}
func (value *APLValueCurrentRuneDeath) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueCurrentRuneDeath) GetBool(_ *Simulation) bool {
	return value.unit.RuneIsDeath(int8(value.runeSlot))
}
func (value *APLValueCurrentRuneDeath) String() string {
	return fmt.Sprintf("Current Rune Death(%d)", value.runeSlot)
}

type APLValueRuneCooldown struct {
	DefaultAPLValueImpl
	unit     *Unit
	runeType proto.APLValueRuneType
}

func (rot *APLRotation) newValueRuneCooldown(config *proto.APLValueRuneCooldown) APLValue {
	unit := rot.unit
	if !unit.HasRunicPowerBar() {
		rot.ValidationWarning("%s does not use Runes", unit.Label)
		return nil
	}
	return &APLValueRuneCooldown{
		unit:     unit,
		runeType: config.RuneType,
	}
}
func (value *APLValueRuneCooldown) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueRuneCooldown) GetDuration(sim *Simulation) time.Duration {
	returnValue := time.Duration(0)
	switch value.runeType {
	case proto.APLValueRuneType_RuneBlood:
		returnValue = value.unit.BloodRuneReadyAt(sim) - sim.CurrentTime
	case proto.APLValueRuneType_RuneFrost:
		returnValue = value.unit.FrostRuneReadyAt(sim) - sim.CurrentTime
	case proto.APLValueRuneType_RuneUnholy:
		returnValue = value.unit.UnholyRuneReadyAt(sim) - sim.CurrentTime
	}
	return max(0, returnValue)
}
func (value *APLValueRuneCooldown) String() string {
	return fmt.Sprintf("Rune Cooldown(%s)", value.runeType)
}

type APLValueNextRuneCooldown struct {
	DefaultAPLValueImpl
	unit     *Unit
	runeType proto.APLValueRuneType
}

func (rot *APLRotation) newValueNextRuneCooldown(config *proto.APLValueNextRuneCooldown) APLValue {
	unit := rot.unit
	if !unit.HasRunicPowerBar() {
		rot.ValidationWarning("%s does not use Runes", unit.Label)
		return nil
	}
	return &APLValueNextRuneCooldown{
		unit:     unit,
		runeType: config.RuneType,
	}
}
func (value *APLValueNextRuneCooldown) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueNextRuneCooldown) GetDuration(sim *Simulation) time.Duration {
	returnValue := time.Duration(0)
	switch value.runeType {
	case proto.APLValueRuneType_RuneBlood:
		returnValue = value.unit.NextBloodRuneReadyAt(sim) - sim.CurrentTime
	case proto.APLValueRuneType_RuneFrost:
		returnValue = value.unit.NextFrostRuneReadyAt(sim) - sim.CurrentTime
	case proto.APLValueRuneType_RuneUnholy:
		returnValue = value.unit.NextUnholyRuneReadyAt(sim) - sim.CurrentTime
	}
	return max(0, returnValue)
}
func (value *APLValueNextRuneCooldown) String() string {
	return fmt.Sprintf("Next Rune Cooldown(%s)", value.runeType)
}

type APLValueRuneSlotCooldown struct {
	DefaultAPLValueImpl
	unit     *Unit
	runeSlot int8
}

func (rot *APLRotation) newValueRuneSlotCooldown(config *proto.APLValueRuneSlotCooldown) APLValue {
	unit := rot.unit
	if !unit.HasRunicPowerBar() {
		rot.ValidationWarning("%s does not use Runes", unit.Label)
		return nil
	}
	return &APLValueRuneSlotCooldown{
		unit:     unit,
		runeSlot: int8(config.RuneSlot) - 1, // 0 is Unknown
	}
}
func (value *APLValueRuneSlotCooldown) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueRuneSlotCooldown) GetDuration(sim *Simulation) time.Duration {
	return max(0, value.unit.RuneReadyAt(sim, value.runeSlot)-sim.CurrentTime)
}
func (value *APLValueRuneSlotCooldown) String() string {
	return fmt.Sprintf("Rune Slot Cooldown(%d)", value.runeSlot)
}

type APLValueRuneGrace struct {
	DefaultAPLValueImpl
	unit     *Unit
	runeType proto.APLValueRuneType
}

func (rot *APLRotation) newValueRuneGrace(config *proto.APLValueRuneGrace) APLValue {
	unit := rot.unit
	if !unit.HasRunicPowerBar() {
		rot.ValidationWarning("%s does not use Runes", unit.Label)
		return nil
	}
	return &APLValueRuneGrace{
		unit:     unit,
		runeType: config.RuneType,
	}
}
func (value *APLValueRuneGrace) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueRuneGrace) GetDuration(sim *Simulation) time.Duration {
	switch value.runeType {
	case proto.APLValueRuneType_RuneBlood:
		return value.unit.CurrentBloodRuneGrace(sim)
	case proto.APLValueRuneType_RuneFrost:
		return value.unit.CurrentFrostRuneGrace(sim)
	case proto.APLValueRuneType_RuneUnholy:
		return value.unit.CurrentUnholyRuneGrace(sim)
	}
	return 0
}
func (value *APLValueRuneGrace) String() string {
	return fmt.Sprintf("Rune Grace(%s)", value.runeType)
}

type APLValueRuneSlotGrace struct {
	DefaultAPLValueImpl
	unit     *Unit
	runeSlot int8
}

func (rot *APLRotation) newValueRuneSlotGrace(config *proto.APLValueRuneSlotGrace) APLValue {
	unit := rot.unit
	if !unit.HasRunicPowerBar() {
		rot.ValidationWarning("%s does not use Runes", unit.Label)
		return nil
	}
	return &APLValueRuneSlotGrace{
		unit:     unit,
		runeSlot: int8(config.RuneSlot) - 1, // 0 is Unknown
	}
}
func (value *APLValueRuneSlotGrace) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueRuneSlotGrace) GetDuration(sim *Simulation) time.Duration {
	return value.unit.CurrentRuneGrace(sim, value.runeSlot)
}
func (value *APLValueRuneSlotGrace) String() string {
	return fmt.Sprintf("Rune Slot Grace(%d)", value.runeSlot)
}
