package core

import (
	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValueCurrentMana struct {
	defaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentMana(config *proto.APLValueCurrentMana) APLValue {
	unit := rot.unit
	if !unit.HasManaBar() {
		return nil
	}
	return &APLValueCurrentMana{
		unit: unit,
	}
}
func (value *APLValueCurrentMana) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentMana) GetFloat(sim *Simulation) float64 {
	return value.unit.CurrentMana()
}

type APLValueCurrentManaPercent struct {
	defaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentManaPercent(config *proto.APLValueCurrentManaPercent) APLValue {
	unit := rot.unit
	if !unit.HasManaBar() {
		return nil
	}
	return &APLValueCurrentManaPercent{
		unit: unit,
	}
}
func (value *APLValueCurrentManaPercent) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentManaPercent) GetFloat(sim *Simulation) float64 {
	return value.unit.CurrentManaPercent()
}

type APLValueCurrentRage struct {
	defaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentRage(config *proto.APLValueCurrentRage) APLValue {
	unit := rot.unit
	if !unit.HasRageBar() {
		return nil
	}
	return &APLValueCurrentRage{
		unit: unit,
	}
}
func (value *APLValueCurrentRage) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentRage) GetFloat(sim *Simulation) float64 {
	return value.unit.CurrentRage()
}

type APLValueCurrentEnergy struct {
	defaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentEnergy(config *proto.APLValueCurrentEnergy) APLValue {
	unit := rot.unit
	if !unit.HasEnergyBar() {
		return nil
	}
	return &APLValueCurrentEnergy{
		unit: unit,
	}
}
func (value *APLValueCurrentEnergy) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentEnergy) GetFloat(sim *Simulation) float64 {
	return value.unit.CurrentEnergy()
}

type APLValueCurrentComboPoints struct {
	defaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentComboPoints(config *proto.APLValueCurrentComboPoints) APLValue {
	unit := rot.unit
	if !unit.HasEnergyBar() {
		return nil
	}
	return &APLValueCurrentComboPoints{
		unit: unit,
	}
}
func (value *APLValueCurrentComboPoints) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueCurrentComboPoints) GetInt(sim *Simulation) int32 {
	return value.unit.ComboPoints()
}

type APLValueCurrentRunicPower struct {
	defaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentRunicPower(config *proto.APLValueCurrentRunicPower) APLValue {
	unit := rot.unit
	if !unit.HasRunicPowerBar() {
		return nil
	}
	return &APLValueCurrentRunicPower{
		unit: unit,
	}
}
func (value *APLValueCurrentRunicPower) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueCurrentRunicPower) GetInt(sim *Simulation) int32 {
	return int32(value.unit.CurrentRunicPower())
}
