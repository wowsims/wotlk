package core

import (
	"time"
)

type OnBloodRuneGain func(sim *Simulation)
type OnFrostRuneGain func(sim *Simulation)
type OnUnholyRuneGain func(sim *Simulation)
type OnDeathRuneGain func(sim *Simulation)
type OnRunicPowerGain func(sim *Simulation)

type RuneState uint8

const (
	RuneState_Spent RuneState = iota
	RuneState_Normal
	RuneState_Death
)

type Rune struct {
	state RuneState
	pas   [2]*PendingAction
}

type runicPowerBar struct {
	unit *Unit

	maxRunicPower     float64
	currentRunicPower float64

	bloodRunes  [2]Rune
	frostRunes  [2]Rune
	unholyRunes [2]Rune

	onBloodRuneGain  OnBloodRuneGain
	onFrostRuneGain  OnFrostRuneGain
	onUnholyRuneGain OnUnholyRuneGain
	onDeathRuneGain  OnDeathRuneGain
	onRunicPowerGain OnRunicPowerGain
}

func (unit *Unit) EnableRunicPowerBar(maxRunicPower float64,
	onBloodRuneGain OnBloodRuneGain,
	onFrostRuneGain OnFrostRuneGain,
	onUnholyRuneGain OnUnholyRuneGain,
	onDeathRuneGain OnDeathRuneGain,
	onRunicPowerGain OnRunicPowerGain) {
	unit.runicPowerBar = runicPowerBar{
		unit: unit,

		maxRunicPower:     maxRunicPower,
		currentRunicPower: maxRunicPower,

		bloodRunes:  [2]Rune{Rune{state: RuneState_Normal, pas: [2]*PendingAction{nil, nil}}, Rune{state: RuneState_Normal, pas: [2]*PendingAction{nil, nil}}},
		frostRunes:  [2]Rune{Rune{state: RuneState_Normal, pas: [2]*PendingAction{nil, nil}}, Rune{state: RuneState_Normal, pas: [2]*PendingAction{nil, nil}}},
		unholyRunes: [2]Rune{Rune{state: RuneState_Normal, pas: [2]*PendingAction{nil, nil}}, Rune{state: RuneState_Normal, pas: [2]*PendingAction{nil, nil}}},

		onBloodRuneGain:  onBloodRuneGain,
		onFrostRuneGain:  onFrostRuneGain,
		onUnholyRuneGain: onUnholyRuneGain,
		onDeathRuneGain:  onDeathRuneGain,
		onRunicPowerGain: onRunicPowerGain,
	}
}

func (unit *Unit) HasRunicPower() bool {
	return unit.runicPowerBar.unit != nil
}

func (rp *runicPowerBar) CurrentRunicPower() float64 {
	return rp.currentRunicPower
}

func (rp *runicPowerBar) addRunicPowerInterval(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to add negative runic power!")
	}

	newRunicPower := MinFloat(rp.currentRunicPower+amount, rp.maxRunicPower)
	metrics.AddEvent(amount, newRunicPower-rp.currentRunicPower)

	if sim.Log != nil {
		rp.unit.Log(sim, "Gained %0.3f runic power from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, rp.currentRunicPower, newRunicPower)
	}

	rp.currentRunicPower = newRunicPower
}

func (rp *runicPowerBar) AddRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	rp.addRunicPowerInterval(sim, amount, metrics)
	rp.onRunicPowerGain(sim)
}

func (rp *runicPowerBar) SpendRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to spend negative runic power!")
	}

	newRunicPower := rp.currentRunicPower - amount
	metrics.AddEvent(-amount, -amount)

	if sim.Log != nil {
		rp.unit.Log(sim, "Spent %0.3f runic power from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, rp.currentRunicPower, newRunicPower)
	}

	rp.currentRunicPower = newRunicPower

}

func CurrentRunesOfType(rb *[2]Rune, runeState RuneState) int32 {
	return TernaryInt32(rb[0].state == runeState, 1, 0) + TernaryInt32(rb[1].state == runeState, 1, 0)
}

func (rp *runicPowerBar) CurrentBloodRunes() int32 {
	return CurrentRunesOfType(&rp.bloodRunes, RuneState_Normal)
}

func (rp *runicPowerBar) CurrentFrostRunes() int32 {
	return CurrentRunesOfType(&rp.frostRunes, RuneState_Normal)
}

func (rp *runicPowerBar) CurrentUnholyRunes() int32 {
	return CurrentRunesOfType(&rp.unholyRunes, RuneState_Normal)
}

func (rp *runicPowerBar) CurrentDeathRunes() int32 {
	return CurrentRunesOfType(&rp.bloodRunes, RuneState_Death) + CurrentRunesOfType(&rp.frostRunes, RuneState_Death) + CurrentRunesOfType(&rp.unholyRunes, RuneState_Death)
}

func (rp *runicPowerBar) CastCostPossible(sim *Simulation, runicPowerAmount float64, bloodAmount int32, frostAmount int32, unholyAmount int32, deathAmount int32) bool {
	return (rp.currentRunicPower > runicPowerAmount) &&
		(rp.CurrentBloodRunes() >= bloodAmount) &&
		(rp.CurrentFrostRunes() >= frostAmount) &&
		(rp.CurrentUnholyRunes() >= unholyAmount) &&
		(rp.CurrentDeathRunes() >= deathAmount)
}

func (rp *runicPowerBar) GenerateRuneMetrics(sim *Simulation, metrics *ResourceMetrics, name string, currRunes int32, newRunes int32) {
	metrics.AddEvent(1, 1)

	if sim.Log != nil {
		rp.unit.Log(sim, "Generated %s Rune from %s (%d --> %d).", name, metrics.ActionID, currRunes, newRunes)
	}
}

func (rp *runicPowerBar) SpendRuneMetrics(sim *Simulation, metrics *ResourceMetrics, name string, currRunes int32, newRunes int32) {
	metrics.AddEvent(-1, -1)

	if sim.Log != nil {
		rp.unit.Log(sim, "Spent %s Rune from %s (%d --> %d).", name, metrics.ActionID, currRunes, newRunes)
	}
}

func GenerateRuneOfType(rb *[2]Rune, runeState RuneState) int32 {
	slot := int32(-1)
	if rb[0].state == RuneState_Spent {
		rb[0].state = runeState
		slot = 0
	} else if rb[1].state == RuneState_Spent {
		rb[1].state = runeState
		slot = 1
	}
	return slot
}

func SpendRuneFromType(rb *[2]Rune, runeState RuneState) int32 {
	slot := int32(-1)
	if rb[0].state == runeState {
		rb[0].state = RuneState_Spent
		slot = 0
	} else if rb[1].state == runeState {
		rb[1].state = RuneState_Spent
		slot = 1
	} else {
		panic("Trying to spend rune that does not exist!")
	}
	return slot
}

func (rp *runicPowerBar) GenerateBloodRune(sim *Simulation, metrics *ResourceMetrics) {
	currRunes := rp.CurrentBloodRunes()
	rp.GenerateRuneMetrics(sim, metrics, "Blood", currRunes, currRunes+1)
	GenerateRuneOfType(&rp.bloodRunes, RuneState_Normal)
}

func (rp *runicPowerBar) GenerateFrostRune(sim *Simulation, metrics *ResourceMetrics) {
	currRunes := rp.CurrentFrostRunes()
	rp.GenerateRuneMetrics(sim, metrics, "Frost", currRunes, currRunes+1)
	GenerateRuneOfType(&rp.frostRunes, RuneState_Normal)
}

func (rp *runicPowerBar) GenerateUnholyRune(sim *Simulation, metrics *ResourceMetrics) {
	currRunes := rp.CurrentUnholyRunes()
	rp.GenerateRuneMetrics(sim, metrics, "Unholy", currRunes, currRunes+1)
	GenerateRuneOfType(&rp.unholyRunes, RuneState_Normal)
}

func (rp *runicPowerBar) SpendBloodRune(sim *Simulation, metrics *ResourceMetrics) {
	currRunes := rp.CurrentBloodRunes()
	rp.SpendRuneMetrics(sim, metrics, "Blood", currRunes, currRunes-1)
	spendSlot := SpendRuneFromType(&rp.bloodRunes, RuneState_Normal)

	pa := &rp.bloodRunes[spendSlot].pas[0]
	if *pa == nil {
		*pa = &PendingAction{}
	}

	(*pa).NextActionAt = sim.CurrentTime + 10*time.Second
	(*pa).Priority = ActionPriorityRegen
	(*pa).OnAction = func(sim *Simulation) {
		rp.GenerateBloodRune(sim, metrics)
		rp.onBloodRuneGain(sim)
	}

	sim.AddPendingAction(*pa)
}

func (rp *runicPowerBar) SpendFrostRune(sim *Simulation, metrics *ResourceMetrics) {
	currRunes := rp.CurrentFrostRunes()
	rp.SpendRuneMetrics(sim, metrics, "Frost", currRunes, currRunes-1)
	spendSlot := SpendRuneFromType(&rp.frostRunes, RuneState_Normal)

	pa := &rp.frostRunes[spendSlot].pas[0]
	if *pa == nil {
		*pa = &PendingAction{}
	}

	(*pa).NextActionAt = sim.CurrentTime + 10*time.Second
	(*pa).Priority = ActionPriorityRegen
	(*pa).OnAction = func(sim *Simulation) {
		rp.GenerateFrostRune(sim, metrics)
		rp.onFrostRuneGain(sim)
	}

	sim.AddPendingAction(*pa)
}

func (rp *runicPowerBar) SpendUnholyRune(sim *Simulation, metrics *ResourceMetrics) {
	currRunes := rp.CurrentUnholyRunes()
	rp.SpendRuneMetrics(sim, metrics, "Unholy", currRunes, currRunes-1)
	spendSlot := SpendRuneFromType(&rp.unholyRunes, RuneState_Normal)

	pa := &rp.unholyRunes[spendSlot].pas[0]
	if *pa == nil {
		*pa = &PendingAction{}
	}

	(*pa).NextActionAt = sim.CurrentTime + 10*time.Second
	(*pa).Priority = ActionPriorityRegen
	(*pa).OnAction = func(sim *Simulation) {
		rp.GenerateUnholyRune(sim, metrics)
		rp.onUnholyRuneGain(sim)
	}

	sim.AddPendingAction(*pa)
}

// TODO: Implement this
func (rp *runicPowerBar) SpendDeathRune(sim *Simulation, metrics *ResourceMetrics) {
	currRunes := rp.CurrentUnholyRunes()
	rp.SpendRuneMetrics(sim, metrics, "Unholy", currRunes, currRunes-1)

	//pa := &PendingAction{
	//	NextActionAt: sim.CurrentTime + 30*time.Second,
	//	Priority:     ActionPriorityRegen,
	//}
	//
	//pa.OnAction = func(sim *Simulation) {
	//	if !pa.cancelled {
	//		rp.GenerateBloodRune(sim, metrics)
	//	}
	//}
	//
	//sim.AddPendingAction(pa)
}
