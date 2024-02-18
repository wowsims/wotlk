package core

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

// Time between energy ticks.
const EnergyTickDuration = time.Millisecond * 100
const EnergyPerTick = 1.0

type energyBar struct {
	unit *Unit

	maxEnergy     float64
	currentEnergy float64

	comboPoints int32

	// List of energy levels that might affect APL decisions. E.g:
	// [10, 15, 20, 30, 60, 85]
	energyDecisionThresholds []int

	// Slice with len == maxEnergy+1 with each index corresponding to an amount of energy. Looks like this:
	// [0, 0, 0, 0, 1, 1, 1, 2, 2, 2, 2, 2, 3, 3, ...]
	// Increments by 1 at each value of energyDecisionThresholds.
	cumulativeEnergyDecisionThresholds []int

	nextEnergyTick time.Duration

	// Multiplies energy regen from ticks.
	EnergyTickMultiplier float64

	regenMetrics        *ResourceMetrics
	EnergyRefundMetrics *ResourceMetrics
}

func (unit *Unit) EnableEnergyBar(maxEnergy float64) {
	unit.SetCurrentPowerBar(EnergyBar)

	unit.energyBar = energyBar{
		unit:                 unit,
		maxEnergy:            max(100, maxEnergy),
		EnergyTickMultiplier: 1,
		regenMetrics:         unit.NewEnergyMetrics(ActionID{OtherID: proto.OtherAction_OtherActionEnergyRegen}),
		EnergyRefundMetrics:  unit.NewEnergyMetrics(ActionID{OtherID: proto.OtherAction_OtherActionRefund}),
	}
}

// Computes the energy thresholds.
func (eb *energyBar) setupEnergyThresholds() {
	if eb.unit == nil {
		return
	}
	var energyThresholds []int

	// Energy thresholds from spell costs.
	for _, action := range eb.unit.Rotation.allAPLActions() {
		for _, spell := range action.GetAllSpells() {
			if _, ok := spell.Cost.(*EnergyCost); ok {
				energyThresholds = append(energyThresholds, int(math.Ceil(spell.DefaultCast.Cost)))
			}
		}
	}

	// Energy thresholds from conditional comparisons.
	for _, action := range eb.unit.Rotation.allAPLActions() {
		for _, value := range action.GetAllAPLValues() {
			if cmpValue, ok := value.(*APLValueCompare); ok {
				_, lhsIsEnergy := cmpValue.lhs.(*APLValueCurrentEnergy)
				_, rhsIsEnergy := cmpValue.rhs.(*APLValueCurrentEnergy)
				if !lhsIsEnergy && !rhsIsEnergy {
					continue
				}

				lhsConstVal := getConstAPLFloatValue(cmpValue.lhs)
				rhsConstVal := getConstAPLFloatValue(cmpValue.rhs)

				if lhsIsEnergy && rhsConstVal != -1 {
					energyThresholds = append(energyThresholds, int(math.Ceil(rhsConstVal)))
				} else if rhsIsEnergy && lhsConstVal != -1 {
					energyThresholds = append(energyThresholds, int(math.Ceil(lhsConstVal)))
				}
			}
		}
	}

	slices.SortStableFunc(energyThresholds, func(t1, t2 int) int {
		return t1 - t2
	})

	// Add each unique value to the final thresholds list.
	curVal := 0
	for _, threshold := range energyThresholds {
		if threshold > curVal {
			eb.energyDecisionThresholds = append(eb.energyDecisionThresholds, threshold)
			curVal = threshold
		}
	}

	curEnergy := 0
	cumulativeVal := 0
	eb.cumulativeEnergyDecisionThresholds = make([]int, int(eb.maxEnergy)+1)
	for _, threshold := range eb.energyDecisionThresholds {
		for curEnergy < threshold {
			eb.cumulativeEnergyDecisionThresholds[curEnergy] = cumulativeVal
			curEnergy++
		}
		cumulativeVal++
	}
	for curEnergy < len(eb.cumulativeEnergyDecisionThresholds) {
		eb.cumulativeEnergyDecisionThresholds[curEnergy] = cumulativeVal
		curEnergy++
	}
}

func (unit *Unit) HasEnergyBar() bool {
	return unit.energyBar.unit != nil
}

func (eb *energyBar) CurrentEnergy() float64 {
	return eb.currentEnergy
}

func (eb *energyBar) NextEnergyTickAt() time.Duration {
	return eb.nextEnergyTick
}

func (eb *energyBar) onEnergyGain(sim *Simulation, crossedThreshold bool) {
	if sim.CurrentTime < 0 {
		return
	}

	if !sim.Options.Interactive && crossedThreshold {
		eb.unit.Rotation.DoNextAction(sim)
	}
}

func (eb *energyBar) addEnergyInternal(sim *Simulation, amount float64, metrics *ResourceMetrics) bool {
	if amount < 0 {
		panic("Trying to add negative energy!")
	}

	newEnergy := min(eb.currentEnergy+amount, eb.maxEnergy)
	metrics.AddEvent(amount, newEnergy-eb.currentEnergy)

	if sim.Log != nil {
		eb.unit.Log(sim, "Gained %0.3f energy from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, eb.currentEnergy, newEnergy)
	}

	crossedThreshold := eb.cumulativeEnergyDecisionThresholds == nil || eb.cumulativeEnergyDecisionThresholds[int(eb.currentEnergy)] != eb.cumulativeEnergyDecisionThresholds[int(newEnergy)]
	eb.currentEnergy = newEnergy

	return crossedThreshold
}
func (eb *energyBar) AddEnergy(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	crossedThreshold := eb.addEnergyInternal(sim, amount, metrics)
	eb.onEnergyGain(sim, crossedThreshold)
}

func (eb *energyBar) SpendEnergy(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to spend negative energy!")
	}

	newEnergy := eb.currentEnergy - amount
	metrics.AddEvent(-amount, -amount)

	if sim.Log != nil {
		eb.unit.Log(sim, "Spent %0.3f energy from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, eb.currentEnergy, newEnergy)
	}

	eb.currentEnergy = newEnergy
}

func (eb *energyBar) ComboPoints() int32 {
	return eb.comboPoints
}

// Gives an immediate partial energy tick and restarts the tick timer.
func (eb *energyBar) ResetEnergyTick(sim *Simulation) {
	timeSinceLastTick := sim.CurrentTime - (eb.NextEnergyTickAt() - EnergyTickDuration)
	partialTickAmount := (EnergyPerTick * eb.EnergyTickMultiplier) * (float64(timeSinceLastTick) / float64(EnergyTickDuration))

	crossedThreshold := eb.addEnergyInternal(sim, partialTickAmount, eb.regenMetrics)
	eb.onEnergyGain(sim, crossedThreshold)

	eb.nextEnergyTick = sim.CurrentTime + EnergyTickDuration
	sim.RescheduleTask(eb.nextEnergyTick)
}

func (eb *energyBar) AddComboPoints(sim *Simulation, pointsToAdd int32, metrics *ResourceMetrics) {
	newComboPoints := min(eb.comboPoints+pointsToAdd, 5)
	metrics.AddEvent(float64(pointsToAdd), float64(newComboPoints-eb.comboPoints))

	if sim.Log != nil {
		eb.unit.Log(sim, "Gained %d combo points from %s (%d --> %d)", pointsToAdd, metrics.ActionID, eb.comboPoints, newComboPoints)
	}

	eb.comboPoints = newComboPoints
}

func (eb *energyBar) SpendComboPoints(sim *Simulation, metrics *ResourceMetrics) {
	if sim.Log != nil {
		eb.unit.Log(sim, "Spent %d combo points from %s (%d --> %d).", eb.comboPoints, metrics.ActionID, eb.comboPoints, 0)
	}
	metrics.AddEvent(float64(-eb.comboPoints), float64(-eb.comboPoints))
	eb.comboPoints = 0
}

func (eb *energyBar) RunTask(sim *Simulation) time.Duration {
	if sim.CurrentTime < eb.nextEnergyTick {
		return eb.nextEnergyTick
	}

	crossedThreshold := eb.addEnergyInternal(sim, EnergyPerTick*eb.EnergyTickMultiplier, eb.regenMetrics)
	eb.onEnergyGain(sim, crossedThreshold)

	eb.nextEnergyTick = sim.CurrentTime + EnergyTickDuration
	return eb.nextEnergyTick
}

func (eb *energyBar) reset(sim *Simulation) {
	if eb.unit == nil {
		return
	}

	eb.currentEnergy = eb.maxEnergy
	eb.comboPoints = 0

	if eb.unit.Type != PetUnit {
		eb.enable(sim, sim.Environment.PrepullStartTime())
	}
}

func (eb *energyBar) enable(sim *Simulation, startAt time.Duration) {
	sim.AddTask(eb)
	eb.nextEnergyTick = startAt + time.Duration(sim.RandomFloat("Energy Tick")*float64(EnergyTickDuration))
	sim.RescheduleTask(eb.nextEnergyTick)

	if eb.cumulativeEnergyDecisionThresholds != nil && sim.Log != nil {
		eb.unit.Log(sim, "[DEBUG] APL Energy decision thresholds: %v", eb.energyDecisionThresholds)
	}
}

func (eb *energyBar) disable(sim *Simulation) {
	eb.nextEnergyTick = NeverExpires
	sim.RemoveTask(eb)
}

type EnergyCostOptions struct {
	Cost float64

	Refund        float64
	RefundMetrics *ResourceMetrics // Optional, will default to unit.EnergyRefundMetrics if not supplied.
}
type EnergyCost struct {
	Refund            float64
	RefundMetrics     *ResourceMetrics
	ResourceMetrics   *ResourceMetrics
	ComboPointMetrics *ResourceMetrics
}

func newEnergyCost(spell *Spell, options EnergyCostOptions) *EnergyCost {
	spell.DefaultCast.Cost = options.Cost
	if options.Refund > 0 && options.RefundMetrics == nil {
		options.RefundMetrics = spell.Unit.EnergyRefundMetrics
	}

	return &EnergyCost{
		Refund:            options.Refund,
		RefundMetrics:     options.RefundMetrics,
		ResourceMetrics:   spell.Unit.NewEnergyMetrics(spell.ActionID),
		ComboPointMetrics: spell.Unit.NewComboPointMetrics(spell.ActionID),
	}
}

func (ec *EnergyCost) MeetsRequirement(_ *Simulation, spell *Spell) bool {
	spell.CurCast.Cost = spell.ApplyCostModifiers(spell.CurCast.Cost)
	return spell.Unit.CurrentEnergy() >= spell.CurCast.Cost
}
func (ec *EnergyCost) CostFailureReason(_ *Simulation, spell *Spell) string {
	return fmt.Sprintf("not enough energy (Current Energy = %0.03f, Energy Cost = %0.03f)", spell.Unit.CurrentEnergy(), spell.CurCast.Cost)
}
func (ec *EnergyCost) SpendCost(sim *Simulation, spell *Spell) {
	spell.Unit.SpendEnergy(sim, spell.CurCast.Cost, ec.ResourceMetrics)
}
func (ec *EnergyCost) IssueRefund(sim *Simulation, spell *Spell) {
	if ec.Refund > 0 {
		spell.Unit.AddEnergy(sim, ec.Refund*spell.CurCast.Cost, ec.RefundMetrics)
	}
}

func (spell *Spell) EnergyMetrics() *ResourceMetrics {
	return spell.Cost.(*EnergyCost).ComboPointMetrics
}

func (spell *Spell) ComboPointMetrics() *ResourceMetrics {
	return spell.Cost.(*EnergyCost).ComboPointMetrics
}
