package core

import (
	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValueCurrentMana struct {
	defaultAPLValueImpl
	unit *Unit
}

func (unit *Unit) newValueCurrentMana(config *proto.APLValueCurrentMana) APLValue {
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

func (unit *Unit) newValueCurrentManaPercent(config *proto.APLValueCurrentManaPercent) APLValue {
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

func (unit *Unit) newValueCurrentRage(config *proto.APLValueCurrentRage) APLValue {
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

func (unit *Unit) newValueCurrentEnergy(config *proto.APLValueCurrentEnergy) APLValue {
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

func (unit *Unit) newValueCurrentComboPoints(config *proto.APLValueCurrentComboPoints) APLValue {
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
