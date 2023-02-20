package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

// Time between focus ticks.
const MaxFocus = 100.0
const tickDuration = time.Second * 1
const BaseFocusPerTick = 5.0

// OnFocusGain is called any time focus is increased.
type OnFocusGain func(sim *Simulation)

type focusBar struct {
	unit *Unit

	focusPerTick float64

	currentFocus float64

	onFocusGain OnFocusGain
	tickAction  *PendingAction
}

func (unit *Unit) EnableFocusBar(regenMultiplier float64, onFocusGain OnFocusGain) {
	unit.focusBar = focusBar{
		unit:         unit,
		focusPerTick: BaseFocusPerTick * regenMultiplier,
		onFocusGain:  onFocusGain,
	}
}

func (unit *Unit) HasFocusBar() bool {
	return unit.focusBar.unit != nil
}

func (fb *focusBar) CurrentFocus() float64 {
	return fb.currentFocus
}

func (fb *focusBar) AddFocus(sim *Simulation, amount float64, actionID ActionID) {
	if amount < 0 {
		panic("Trying to add negative focus!")
	}

	newFocus := MinFloat(fb.currentFocus+amount, MaxFocus)

	if sim.Log != nil {
		fb.unit.Log(sim, "Gained %0.3f focus from %s (%0.3f --> %0.3f).", amount, actionID, fb.currentFocus, newFocus)
	}

	fb.currentFocus = newFocus

	if fb.onFocusGain != nil {
		fb.onFocusGain(sim)
	}
}

func (fb *focusBar) SpendFocus(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to spend negative focus!")
	}

	newFocus := fb.currentFocus - amount
	metrics.AddEvent(-amount, -amount)

	if sim.Log != nil {
		fb.unit.Log(sim, "Spent %0.3f focus from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, fb.currentFocus, newFocus)
	}

	fb.currentFocus = newFocus
}

func (fb *focusBar) reset(sim *Simulation) {
	if fb.unit == nil {
		return
	}

	fb.currentFocus = MaxFocus

	pa := &PendingAction{
		Priority:     ActionPriorityRegen,
		NextActionAt: tickDuration,
	}
	pa.OnAction = func(sim *Simulation) {
		fb.AddFocus(sim, fb.focusPerTick, ActionID{OtherID: proto.OtherAction_OtherActionFocusRegen})

		pa.NextActionAt = sim.CurrentTime + tickDuration
		sim.AddPendingAction(pa)
	}
	fb.tickAction = pa
	sim.AddPendingAction(pa)
}

func (fb *focusBar) Cancel(sim *Simulation) {
	if fb.tickAction != nil {
		fb.tickAction.Cancel(sim)
		fb.tickAction = nil
	}
}

type FocusCostOptions struct {
	Cost float64
}
type FocusCost struct {
	Refund          float64
	ResourceMetrics *ResourceMetrics
}

func newFocusCost(spell *Spell, options FocusCostOptions) *FocusCost {
	spell.DefaultCast.Cost = options.Cost

	return &FocusCost{
		ResourceMetrics: spell.Unit.NewFocusMetrics(spell.ActionID),
	}
}

func (fc *FocusCost) MeetsRequirement(spell *Spell) bool {
	spell.CurCast.Cost = MaxFloat(0, spell.CurCast.Cost*spell.Unit.PseudoStats.CostMultiplier)
	return spell.Unit.CurrentFocus() >= spell.CurCast.Cost
}
func (fc *FocusCost) LogCostFailure(sim *Simulation, spell *Spell) {
	spell.Unit.Log(sim,
		"Failed casting %s, not enough focus. (Current Focus = %0.03f, Focus Cost = %0.03f)",
		spell.ActionID, spell.Unit.CurrentFocus(), spell.CurCast.Cost)
}
func (fc *FocusCost) SpendCost(sim *Simulation, spell *Spell) {
	spell.Unit.SpendFocus(sim, spell.CurCast.Cost, fc.ResourceMetrics)
}
func (fc *FocusCost) IssueRefund(sim *Simulation, spell *Spell) {
}

func (spell *Spell) FocusMetrics() *ResourceMetrics {
	return spell.Cost.(*FocusCost).ResourceMetrics
}
