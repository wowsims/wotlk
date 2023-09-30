package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

// Time between energy ticks.
const EnergyTickDuration = time.Millisecond * 100
const EnergyPerTick = 1.0

// OnEnergyGain is called any time energy is increased.
type OnEnergyGain func(sim *Simulation)

type energyBar struct {
	unit *Unit

	maxEnergy     float64
	currentEnergy float64

	comboPoints int32

	onEnergyGain OnEnergyGain
	tickAction   *PendingAction

	// Multiplies energy regen from ticks.
	EnergyTickMultiplier float64

	regenMetrics        *ResourceMetrics
	EnergyRefundMetrics *ResourceMetrics
}

func (unit *Unit) EnableEnergyBar(maxEnergy float64, onEnergyGain OnEnergyGain) {
	unit.SetCurrentPowerBar(EnergyBar)

	unit.energyBar = energyBar{
		unit:      unit,
		maxEnergy: MaxFloat(100, maxEnergy),
		onEnergyGain: func(sim *Simulation) {
			if sim.CurrentTime < 0 {
				return
			}

			if !sim.Options.Interactive && (!unit.IsWaitingForEnergy() || unit.DoneWaitingForEnergy(sim)) {
				if unit.IsUsingAPL {
					unit.Rotation.DoNextAction(sim)
				} else {
					onEnergyGain(sim)
				}
			}
		},
		EnergyTickMultiplier: 1,
		regenMetrics:         unit.NewEnergyMetrics(ActionID{OtherID: proto.OtherAction_OtherActionEnergyRegen}),
		EnergyRefundMetrics:  unit.NewEnergyMetrics(ActionID{OtherID: proto.OtherAction_OtherActionRefund}),
	}
}

func (unit *Unit) HasEnergyBar() bool {
	return unit.energyBar.unit != nil
}

func (eb *energyBar) CurrentEnergy() float64 {
	return eb.currentEnergy
}

func (eb *energyBar) NextEnergyTickAt() time.Duration {
	return eb.tickAction.NextActionAt
}

func (eb *energyBar) addEnergyInternal(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to add negative energy!")
	}

	newEnergy := MinFloat(eb.currentEnergy+amount, eb.maxEnergy)
	metrics.AddEvent(amount, newEnergy-eb.currentEnergy)

	if sim.Log != nil {
		eb.unit.Log(sim, "Gained %0.3f energy from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, eb.currentEnergy, newEnergy)
	}

	eb.currentEnergy = newEnergy
}
func (eb *energyBar) AddEnergy(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	eb.addEnergyInternal(sim, amount, metrics)
	eb.onEnergyGain(sim)
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

	eb.addEnergyInternal(sim, partialTickAmount, eb.regenMetrics)
	eb.onEnergyGain(sim)

	eb.newTickAction(sim, false, sim.CurrentTime)
}

func (eb *energyBar) AddComboPoints(sim *Simulation, pointsToAdd int32, metrics *ResourceMetrics) {
	newComboPoints := MinInt32(eb.comboPoints+pointsToAdd, 5)
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

func (eb *energyBar) newTickAction(sim *Simulation, randomTickTime bool, startAt time.Duration) {
	if eb.tickAction != nil {
		eb.tickAction.Cancel(sim)
	}

	nextTickDuration := EnergyTickDuration
	if randomTickTime {
		nextTickDuration = time.Duration(sim.RandomFloat("Energy Tick") * float64(EnergyTickDuration))
	}

	pa := &PendingAction{
		NextActionAt: startAt + nextTickDuration,
		Priority:     ActionPriorityRegen,
	}
	pa.OnAction = func(sim *Simulation) {
		eb.addEnergyInternal(sim, EnergyPerTick*eb.EnergyTickMultiplier, eb.regenMetrics)
		eb.onEnergyGain(sim)

		pa.NextActionAt = sim.CurrentTime + EnergyTickDuration
		sim.AddPendingAction(pa)
	}
	eb.tickAction = pa
	sim.AddPendingAction(pa)
}

func (eb *energyBar) reset(sim *Simulation) {
	if eb.unit == nil {
		return
	}

	eb.currentEnergy = eb.maxEnergy
	eb.comboPoints = 0
	eb.newTickAction(sim, true, sim.Environment.PrepullStartTime())
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

func (ec *EnergyCost) MeetsRequirement(spell *Spell) bool {
	spell.CurCast.Cost = spell.ApplyCostModifiers(spell.CurCast.Cost)
	return spell.Unit.CurrentEnergy() >= spell.CurCast.Cost
}
func (ec *EnergyCost) LogCostFailure(sim *Simulation, spell *Spell) {
	spell.Unit.Log(sim,
		"Failed casting %s, not enough energy. (Current Energy = %0.03f, Energy Cost = %0.03f)",
		spell.ActionID, spell.Unit.CurrentEnergy(), spell.CurCast.Cost)
}
func (ec *EnergyCost) SpendCost(sim *Simulation, spell *Spell) {
	if spell.CurCast.Cost > 0 {
		spell.Unit.SpendEnergy(sim, spell.CurCast.Cost, ec.ResourceMetrics)
	}
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
