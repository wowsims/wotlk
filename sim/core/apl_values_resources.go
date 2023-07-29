package core

import (
	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValueCurrentHealth struct {
	defaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentHealth(config *proto.APLValueCurrentHealth) APLValue {
	unit := rot.getSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasHealthBar() {
		rot.validationWarning("%s does not use Health", unit.Get().Label)
		return nil
	}
	return &APLValueCurrentHealth{
		unit: unit,
	}
}
func (value *APLValueCurrentHealth) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentHealth) GetFloat(sim *Simulation) float64 {
	return value.unit.Get().CurrentHealth()
}

type APLValueCurrentHealthPercent struct {
	defaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentHealthPercent(config *proto.APLValueCurrentHealthPercent) APLValue {
	unit := rot.getSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasHealthBar() {
		rot.validationWarning("%s does not use Health", unit.Get().Label)
		return nil
	}
	return &APLValueCurrentHealthPercent{
		unit: unit,
	}
}
func (value *APLValueCurrentHealthPercent) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentHealthPercent) GetFloat(sim *Simulation) float64 {
	return value.unit.Get().CurrentHealthPercent()
}

type APLValueCurrentMana struct {
	defaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentMana(config *proto.APLValueCurrentMana) APLValue {
	unit := rot.getSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasManaBar() {
		rot.validationWarning("%s does not use Mana", unit.Get().Label)
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
	return value.unit.Get().CurrentMana()
}

type APLValueCurrentManaPercent struct {
	defaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentManaPercent(config *proto.APLValueCurrentManaPercent) APLValue {
	unit := rot.getSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasManaBar() {
		rot.validationWarning("%s does not use Mana", unit.Get().Label)
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
	return value.unit.Get().CurrentManaPercent()
}

type APLValueCurrentRage struct {
	defaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentRage(config *proto.APLValueCurrentRage) APLValue {
	unit := rot.unit
	if !unit.HasRageBar() {
		rot.validationWarning("%s does not use Rage", unit.Label)
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
		rot.validationWarning("%s does not use Energy", unit.Label)
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
		rot.validationWarning("%s does not use Combo Points", unit.Label)
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
		rot.validationWarning("%s does not use Runic Power", unit.Label)
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
