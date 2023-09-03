package core

import (
	"fmt"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValueCurrentHealth struct {
	DefaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentHealth(config *proto.APLValueCurrentHealth) APLValue {
	unit := rot.GetSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasHealthBar() {
		rot.ValidationWarning("%s does not use Health", unit.Get().Label)
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
func (value *APLValueCurrentHealth) String() string {
	return "Current Health"
}

type APLValueCurrentHealthPercent struct {
	DefaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentHealthPercent(config *proto.APLValueCurrentHealthPercent) APLValue {
	unit := rot.GetSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasHealthBar() {
		rot.ValidationWarning("%s does not use Health", unit.Get().Label)
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
func (value *APLValueCurrentHealthPercent) String() string {
	return fmt.Sprintf("Current Health %%")
}

type APLValueCurrentMana struct {
	DefaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentMana(config *proto.APLValueCurrentMana) APLValue {
	unit := rot.GetSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasManaBar() {
		rot.ValidationWarning("%s does not use Mana", unit.Get().Label)
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
func (value *APLValueCurrentMana) String() string {
	return "Current Mana"
}

type APLValueCurrentManaPercent struct {
	DefaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentManaPercent(config *proto.APLValueCurrentManaPercent) APLValue {
	unit := rot.GetSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasManaBar() {
		rot.ValidationWarning("%s does not use Mana", unit.Get().Label)
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
func (value *APLValueCurrentManaPercent) String() string {
	return fmt.Sprintf("Current Mana %%")
}

type APLValueCurrentRage struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentRage(config *proto.APLValueCurrentRage) APLValue {
	unit := rot.unit
	if !unit.HasRageBar() {
		rot.ValidationWarning("%s does not use Rage", unit.Label)
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
func (value *APLValueCurrentRage) String() string {
	return "Current Rage"
}

type APLValueCurrentEnergy struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentEnergy(config *proto.APLValueCurrentEnergy) APLValue {
	unit := rot.unit
	if !unit.HasEnergyBar() {
		rot.ValidationWarning("%s does not use Energy", unit.Label)
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
func (value *APLValueCurrentEnergy) String() string {
	return "Current Energy"
}

type APLValueCurrentComboPoints struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentComboPoints(config *proto.APLValueCurrentComboPoints) APLValue {
	unit := rot.unit
	if !unit.HasEnergyBar() {
		rot.ValidationWarning("%s does not use Combo Points", unit.Label)
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
func (value *APLValueCurrentComboPoints) String() string {
	return "Current Combo Points"
}

type APLValueCurrentRunicPower struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentRunicPower(config *proto.APLValueCurrentRunicPower) APLValue {
	unit := rot.unit
	if !unit.HasRunicPowerBar() {
		rot.ValidationWarning("%s does not use Runic Power", unit.Label)
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
func (value *APLValueCurrentRunicPower) String() string {
	return "Current Runic Power"
}
