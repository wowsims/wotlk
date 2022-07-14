package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
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
	RuneState_DeathSpent
	RuneState_Death
)

type DKRuneCost struct {
	blood  int
	frost  int
	unholy int
	death  int
}

type Rune struct {
	state RuneState
	pas   [3]*PendingAction
}

type runicPowerBar struct {
	unit *Unit

	maxRunicPower     float64
	currentRunicPower float64

	bloodRunes  [2]Rune
	frostRunes  [2]Rune
	unholyRunes [2]Rune

	bloodRuneGainMetrics  *ResourceMetrics
	frostRuneGainMetrics  *ResourceMetrics
	unholyRuneGainMetrics *ResourceMetrics
	deathRuneGainMetrics  *ResourceMetrics

	onBloodRuneGain  OnBloodRuneGain
	onFrostRuneGain  OnFrostRuneGain
	onUnholyRuneGain OnUnholyRuneGain
	onDeathRuneGain  OnDeathRuneGain
	onRunicPowerGain OnRunicPowerGain
}

func ResetRune(sim *Simulation, runes *[2]Rune) {
	runes[0].state = RuneState_Normal
	runes[1].state = RuneState_Normal

	if runes[0].pas[0] != nil {
		runes[0].pas[0].Cancel(sim)
	}
	runes[0].pas[0] = nil

	if runes[0].pas[1] != nil {
		runes[0].pas[1].Cancel(sim)
	}
	runes[0].pas[1] = nil

	if runes[0].pas[2] != nil {
		runes[0].pas[2].Cancel(sim)
	}
	runes[0].pas[2] = nil

	if runes[1].pas[0] != nil {
		runes[1].pas[0].Cancel(sim)
	}
	runes[1].pas[0] = nil

	if runes[1].pas[1] != nil {
		runes[1].pas[1].Cancel(sim)
	}
	runes[1].pas[1] = nil

	if runes[1].pas[2] != nil {
		runes[1].pas[2].Cancel(sim)
	}
	runes[1].pas[2] = nil
}

func (rp *runicPowerBar) ResetRunicPowerBar(sim *Simulation) {
	if rp.unit == nil {
		return
	}

	ResetRune(sim, &rp.bloodRunes)
	ResetRune(sim, &rp.frostRunes)
	ResetRune(sim, &rp.unholyRunes)
}

func (unit *Unit) EnableRunicPowerBar(currentRunicPower float64, maxRunicPower float64,
	onBloodRuneGain OnBloodRuneGain,
	onFrostRuneGain OnFrostRuneGain,
	onUnholyRuneGain OnUnholyRuneGain,
	onDeathRuneGain OnDeathRuneGain,
	onRunicPowerGain OnRunicPowerGain) {
	unit.runicPowerBar = runicPowerBar{
		unit: unit,

		maxRunicPower:     maxRunicPower,
		currentRunicPower: currentRunicPower,

		bloodRunes:  [2]Rune{Rune{state: RuneState_Normal, pas: [3]*PendingAction{nil, nil, nil}}, Rune{state: RuneState_Normal, pas: [3]*PendingAction{nil, nil, nil}}},
		frostRunes:  [2]Rune{Rune{state: RuneState_Normal, pas: [3]*PendingAction{nil, nil, nil}}, Rune{state: RuneState_Normal, pas: [3]*PendingAction{nil, nil, nil}}},
		unholyRunes: [2]Rune{Rune{state: RuneState_Normal, pas: [3]*PendingAction{nil, nil, nil}}, Rune{state: RuneState_Normal, pas: [3]*PendingAction{nil, nil, nil}}},

		onBloodRuneGain:  onBloodRuneGain,
		onFrostRuneGain:  onFrostRuneGain,
		onUnholyRuneGain: onUnholyRuneGain,
		onDeathRuneGain:  onDeathRuneGain,
		onRunicPowerGain: onRunicPowerGain,
	}

	unit.bloodRuneGainMetrics = unit.NewBloodRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionBloodRuneGain, Tag: 1})
	unit.frostRuneGainMetrics = unit.NewFrostRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionFrostRuneGain, Tag: 1})
	unit.unholyRuneGainMetrics = unit.NewUnholyRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionUnholyRuneGain, Tag: 1})
	unit.deathRuneGainMetrics = unit.NewDeathRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionDeathRuneGain, Tag: 1})
}

func (unit *Unit) HasRunicPower() bool {
	return unit.runicPowerBar.unit != nil
}

func (rp *runicPowerBar) BloodRuneGainMetrics() *ResourceMetrics {
	return rp.bloodRuneGainMetrics
}

func (rp *runicPowerBar) FrostRuneGainMetrics() *ResourceMetrics {
	return rp.frostRuneGainMetrics
}

func (rp *runicPowerBar) UnholyRuneGainMetrics() *ResourceMetrics {
	return rp.unholyRuneGainMetrics
}

func (rp *runicPowerBar) DeathRuneGainMetrics() *ResourceMetrics {
	return rp.deathRuneGainMetrics
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

func RuneReadyAt(sim *Simulation, runes *[2]Rune) time.Duration {
	readyAt := NeverExpires
	pa := runes[0].pas[0]
	if pa != nil {
		readyAt = MinDuration(readyAt, pa.NextActionAt)
	} else {
		return sim.CurrentTime
	}

	pa = runes[1].pas[0]
	if pa != nil {
		readyAt = MinDuration(readyAt, pa.NextActionAt)
	} else {
		return sim.CurrentTime
	}

	return readyAt
}

func (rp *runicPowerBar) BloodRuneReadyAt(sim *Simulation) time.Duration {
	return RuneReadyAt(sim, &rp.bloodRunes)
}

func (rp *runicPowerBar) FrostRuneReadyAt(sim *Simulation) time.Duration {
	return RuneReadyAt(sim, &rp.frostRunes)
}

func (rp *runicPowerBar) UnholyRuneReadyAt(sim *Simulation) time.Duration {
	return RuneReadyAt(sim, &rp.unholyRunes)
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

func (rp *runicPowerBar) CastCostPossible(sim *Simulation, runicPowerAmount float64, bloodAmount int32, frostAmount int32, unholyAmount int32) bool {
	totalDeathRunes := rp.CurrentDeathRunes()

	if rp.CurrentRunicPower() < runicPowerAmount {
		return false
	}

	if rp.CurrentBloodRunes() < bloodAmount {
		if totalDeathRunes > 0 {
			totalDeathRunes -= 1
		} else {
			return false
		}
	}

	if rp.CurrentFrostRunes() < frostAmount {
		if totalDeathRunes > 0 {
			totalDeathRunes -= 1
		} else {
			return false
		}
	}

	if rp.CurrentUnholyRunes() < unholyAmount {
		if totalDeathRunes > 0 {
			totalDeathRunes -= 1
		} else {
			return false
		}
	}

	return true
}

func (rp *runicPowerBar) DetermineOptimalCost(sim *Simulation, bloodAmount int, frostAmount int, unholyAmount int) DKRuneCost {
	totalBloodRunes := int(rp.CurrentBloodRunes())
	startingBloodRunes := totalBloodRunes
	totalFrostRunes := int(rp.CurrentFrostRunes())
	startingFrostRunes := totalFrostRunes
	totalUnholyRunes := int(rp.CurrentUnholyRunes())
	startingUnholyRunes := totalUnholyRunes
	totalDeathRunes := int(rp.CurrentDeathRunes())
	startingDeathRunes := totalDeathRunes

	if int(rp.CurrentBloodRunes()) >= bloodAmount {
		totalBloodRunes -= bloodAmount
	} else {
		totalDeathRunes -= bloodAmount
	}

	if int(rp.CurrentFrostRunes()) >= frostAmount {
		totalFrostRunes -= frostAmount
	} else {
		totalDeathRunes -= frostAmount
	}

	if int(rp.CurrentUnholyRunes()) >= unholyAmount {
		totalUnholyRunes -= unholyAmount
	} else {
		totalDeathRunes -= unholyAmount
	}

	spellCost := DKRuneCost{
		blood:  startingBloodRunes - totalBloodRunes,
		frost:  startingFrostRunes - totalFrostRunes,
		unholy: startingUnholyRunes - totalUnholyRunes,
		death:  startingDeathRunes - totalDeathRunes,
	}

	return spellCost
}

func (rp *runicPowerBar) Spend(sim *Simulation, spell *Spell, cost DKRuneCost) {
	for i := 0; i < cost.blood; i++ {
		rp.SpendBloodRune(sim, spell.BloodRuneMetrics())
	}
	for i := 0; i < cost.frost; i++ {
		rp.SpendFrostRune(sim, spell.FrostRuneMetrics())
	}
	for i := 0; i < cost.unholy; i++ {
		rp.SpendUnholyRune(sim, spell.UnholyRuneMetrics())
	}
	for i := 0; i < cost.death; i++ {
		rp.SpendDeathRune(sim, spell.DeathRuneMetrics())
	}
}

func (rp *runicPowerBar) GenerateRuneMetrics(sim *Simulation, metrics *ResourceMetrics, name string, currRunes int32, newRunes int32) {
	metrics.AddEvent(1, float64(newRunes)-float64(currRunes))

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

func SetRuneAtSlotToState(rb *[2]Rune, slot int32, runeState RuneState) {
	// TODO: safeguard this?
	rb[slot].state = runeState
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
	spendState := RuneState_Spent
	if runeState == RuneState_Death {
		spendState = RuneState_DeathSpent
	}

	slot := int32(-1)
	if rb[0].state == runeState {
		rb[0].state = spendState
		slot = 0
	} else if rb[1].state == runeState {
		rb[1].state = spendState
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

func (rp *runicPowerBar) GenerateDeathRuneAtSlot(sim *Simulation, metrics *ResourceMetrics, runes *[2]Rune, slot int32) {
	currRunes := rp.CurrentDeathRunes()
	rp.GenerateRuneMetrics(sim, metrics, "Death", currRunes, currRunes+1)
	SetRuneAtSlotToState(runes, slot, RuneState_Death)
}

func (rp *runicPowerBar) SpendBloodRune(sim *Simulation, metrics *ResourceMetrics) int32 {
	currRunes := rp.CurrentBloodRunes()
	if currRunes <= 0 {
		panic("Trying to spend blood runes that don't exist!")
	}

	rp.SpendRuneMetrics(sim, metrics, "Blood", currRunes, currRunes-1)
	spendSlot := SpendRuneFromType(&rp.bloodRunes, RuneState_Normal)

	pa := &PendingAction{
		NextActionAt: sim.CurrentTime + 10*time.Second,
		Priority:     ActionPriorityRegen,
	}

	pa.OnAction = func(sim *Simulation) {
		if !pa.cancelled {
			rp.GenerateBloodRune(sim, rp.bloodRuneGainMetrics)
			rp.onBloodRuneGain(sim)
			rp.bloodRunes[spendSlot].pas[0] = nil
		}
	}

	rp.bloodRunes[spendSlot].pas[0] = pa
	sim.AddPendingAction(pa)

	return spendSlot
}

func (rp *runicPowerBar) SpendFrostRune(sim *Simulation, metrics *ResourceMetrics) int32 {
	currRunes := rp.CurrentFrostRunes()
	if currRunes <= 0 {
		panic("Trying to spend frost runes that don't exist!")
	}

	rp.SpendRuneMetrics(sim, metrics, "Frost", currRunes, currRunes-1)
	spendSlot := SpendRuneFromType(&rp.frostRunes, RuneState_Normal)

	pa := &PendingAction{
		NextActionAt: sim.CurrentTime + 10*time.Second,
		Priority:     ActionPriorityRegen,
	}

	pa.OnAction = func(sim *Simulation) {
		if !pa.cancelled {
			rp.GenerateFrostRune(sim, rp.frostRuneGainMetrics)
			rp.onFrostRuneGain(sim)
		}
		rp.frostRunes[spendSlot].pas[0] = nil
	}

	rp.frostRunes[spendSlot].pas[0] = pa
	sim.AddPendingAction(pa)

	return spendSlot
}

func (rp *runicPowerBar) SpendUnholyRune(sim *Simulation, metrics *ResourceMetrics) int32 {
	currRunes := rp.CurrentUnholyRunes()
	if currRunes <= 0 {
		panic("Trying to spend unholy runes that don't exist!")
	}

	rp.SpendRuneMetrics(sim, metrics, "Unholy", currRunes, currRunes-1)
	spendSlot := SpendRuneFromType(&rp.unholyRunes, RuneState_Normal)

	pa := &PendingAction{
		NextActionAt: sim.CurrentTime + 10*time.Second,
		Priority:     ActionPriorityRegen,
	}

	pa.OnAction = func(sim *Simulation) {
		if !pa.cancelled {
			rp.GenerateUnholyRune(sim, rp.unholyRuneGainMetrics)
			rp.onUnholyRuneGain(sim)
		}
		rp.unholyRunes[spendSlot].pas[0] = nil
	}

	rp.unholyRunes[spendSlot].pas[0] = pa
	sim.AddPendingAction(pa)

	return spendSlot
}

func (rp *runicPowerBar) SpendDeathRune(sim *Simulation, metrics *ResourceMetrics) {
	currRunes := rp.CurrentDeathRunes()
	if currRunes <= 0 {
		panic("Trying to spend death runes that don't exist!")
	}

	rp.SpendRuneMetrics(sim, metrics, "Death", currRunes, currRunes-1)

	runeTypeIdx := 0
	spendSlot := SpendRuneFromType(&rp.bloodRunes, RuneState_Death)
	if spendSlot < 0 {
		runeTypeIdx += 1
		spendSlot = SpendRuneFromType(&rp.frostRunes, RuneState_Death)
		if spendSlot < 0 {
			runeTypeIdx += 1
			SpendRuneFromType(&rp.unholyRunes, RuneState_Death)
		}
	}

	pa := &PendingAction{
		NextActionAt: sim.CurrentTime + 10*time.Second,
		Priority:     ActionPriorityRegen,
	}

	pa.OnAction = func(sim *Simulation) {
		if runeTypeIdx == 0 {
			if !pa.cancelled {
				rp.GenerateDeathRuneAtSlot(sim, rp.deathRuneGainMetrics, &rp.bloodRunes, spendSlot)
				rp.onDeathRuneGain(sim)
			}
			rp.bloodRunes[spendSlot].pas[1] = nil
		} else if runeTypeIdx == 1 {
			if !pa.cancelled {
				rp.GenerateDeathRuneAtSlot(sim, rp.deathRuneGainMetrics, &rp.frostRunes, spendSlot)
				rp.onDeathRuneGain(sim)
			}
			rp.frostRunes[spendSlot].pas[1] = nil
		} else if runeTypeIdx == 2 {
			if !pa.cancelled {
				rp.GenerateDeathRuneAtSlot(sim, rp.deathRuneGainMetrics, &rp.unholyRunes, spendSlot)
				rp.onDeathRuneGain(sim)
			}
			rp.unholyRunes[spendSlot].pas[1] = nil
		}
	}

	if runeTypeIdx == 0 {
		rp.bloodRunes[spendSlot].pas[1] = pa
	} else if runeTypeIdx == 1 {
		rp.frostRunes[spendSlot].pas[1] = pa
	} else if runeTypeIdx == 2 {
		rp.unholyRunes[spendSlot].pas[1] = pa
	}
	sim.AddPendingAction(pa)
}
