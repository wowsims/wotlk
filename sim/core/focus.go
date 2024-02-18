package core

import (
	"fmt"
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

	nextFocusTick time.Duration

	regenMetrics  *ResourceMetrics
	refundMetrics *ResourceMetrics
}

func (unit *Unit) EnableFocusBar(regenMultiplier float64, onFocusGain OnFocusGain) {
	unit.focusBar = focusBar{
		unit:          unit,
		focusPerTick:  BaseFocusPerTick * regenMultiplier,
		onFocusGain:   onFocusGain,
		regenMetrics:  unit.NewEnergyMetrics(ActionID{OtherID: proto.OtherAction_OtherActionFocusRegen}),
		refundMetrics: unit.NewEnergyMetrics(ActionID{OtherID: proto.OtherAction_OtherActionRefund}),
	}
}

func (unit *Unit) HasFocusBar() bool {
	return unit.focusBar.unit != nil
}

func (fb *focusBar) CurrentFocus() float64 {
	return fb.currentFocus
}

func (fb *focusBar) AddFocus(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to add negative focus!")
	}

	newFocus := min(fb.currentFocus+amount, MaxFocus)
	metrics.AddEvent(amount, newFocus-fb.currentFocus)

	if sim.Log != nil {
		fb.unit.Log(sim, "Gained %0.3f focus from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, fb.currentFocus, newFocus)
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

	if fb.unit.Type != PetUnit {
		fb.enable(sim)
	}
}

func (fb *focusBar) enable(sim *Simulation) {
	sim.AddTask(fb)
	fb.nextFocusTick = sim.CurrentTime + tickDuration
	sim.RescheduleTask(fb.nextFocusTick)
}

func (fb *focusBar) disable(sim *Simulation) {
	fb.nextFocusTick = NeverExpires
	sim.RemoveTask(fb)
}

func (fb *focusBar) RunTask(sim *Simulation) time.Duration {
	if sim.CurrentTime < fb.nextFocusTick {
		return fb.nextFocusTick
	}

	fb.AddFocus(sim, fb.focusPerTick, fb.regenMetrics)

	fb.nextFocusTick = sim.CurrentTime + tickDuration
	return fb.nextFocusTick
}

type FocusCostOptions struct {
	Cost float64

	Refund        float64
	RefundMetrics *ResourceMetrics // Optional, will default to unit.FocusRefundMetrics if not supplied
}
type FocusCost struct {
	Refund          float64
	RefundMetrics   *ResourceMetrics
	ResourceMetrics *ResourceMetrics
}

func newFocusCost(spell *Spell, options FocusCostOptions) *FocusCost {
	spell.DefaultCast.Cost = options.Cost
	if options.Refund > 0 && options.RefundMetrics == nil {
		options.RefundMetrics = spell.Unit.refundMetrics
	}
	return &FocusCost{
		Refund:          options.Refund,
		RefundMetrics:   options.RefundMetrics,
		ResourceMetrics: spell.Unit.NewFocusMetrics(spell.ActionID),
	}
}

func (fc *FocusCost) MeetsRequirement(_ *Simulation, spell *Spell) bool {
	spell.CurCast.Cost = max(0, spell.CurCast.Cost*spell.Unit.PseudoStats.CostMultiplier)
	return spell.Unit.CurrentFocus() >= spell.CurCast.Cost
}
func (fc *FocusCost) CostFailureReason(_ *Simulation, spell *Spell) string {
	return fmt.Sprintf("not enough focus (Current Focus = %0.03f, Focus Cost = %0.03f)", spell.Unit.CurrentFocus(), spell.CurCast.Cost)
}
func (fc *FocusCost) SpendCost(sim *Simulation, spell *Spell) {
	spell.Unit.SpendFocus(sim, spell.CurCast.Cost, fc.ResourceMetrics)
}
func (fc *FocusCost) IssueRefund(sim *Simulation, spell *Spell) {
	if fc.Refund > 0 {
		spell.Unit.AddFocus(sim, fc.Refund*spell.CurCast.Cost, fc.RefundMetrics)
	}
}
