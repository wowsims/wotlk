package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type OnRuneSpend func(sim *Simulation)
type OnBloodRuneGain func(sim *Simulation)
type OnFrostRuneGain func(sim *Simulation)
type OnUnholyRuneGain func(sim *Simulation)
type OnDeathRuneGain func(sim *Simulation)
type OnRunicPowerGain func(sim *Simulation)

type RuneState uint8
type RuneKind uint8

const (
	RuneKind_Undef RuneKind = iota
	RuneKind_Blood
	RuneKind_Frost
	RuneKind_Unholy
	RuneKind_Death
)

const (
	RuneState_Normal RuneState = iota
	RuneState_Spent
	RuneState_Death
	RuneState_DeathSpent
)

type RuneAmount struct {
	Blood  int
	Frost  int
	Unholy int
	Death  int
}

type Rune struct {
	state         RuneState
	kind          RuneKind
	pas           [2]*PendingAction
	lastRegenTime time.Duration
	lastSpendTime time.Duration
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

	onRuneSpend      OnRuneSpend
	onBloodRuneGain  OnBloodRuneGain
	onFrostRuneGain  OnFrostRuneGain
	onUnholyRuneGain OnUnholyRuneGain
	onDeathRuneGain  OnDeathRuneGain
	onRunicPowerGain OnRunicPowerGain
	isACopy          bool
}

func (rp *runicPowerBar) CopyRunicPowerBar() runicPowerBar {
	crp := *rp
	crp.isACopy = true
	return crp
}

func ResetRunes(sim *Simulation, runes *[2]Rune, runeKind RuneKind) {
	runes[0].state = RuneState_Normal
	runes[0].kind = runeKind
	runes[0].lastRegenTime = -1
	runes[1].state = RuneState_Normal
	runes[1].kind = runeKind
	runes[1].lastRegenTime = -1

	if runes[0].pas[0] != nil {
		runes[0].pas[0].Cancel(sim)
		runes[0].pas[0] = nil
	}
	if runes[0].pas[1] != nil {
		runes[0].pas[1].Cancel(sim)
		runes[0].pas[1] = nil
	}
	if runes[1].pas[0] != nil {
		runes[1].pas[0].Cancel(sim)
		runes[1].pas[0] = nil
	}
	if runes[1].pas[1] != nil {
		runes[1].pas[1].Cancel(sim)
		runes[1].pas[1] = nil
	}
}

func (rp *runicPowerBar) reset(sim *Simulation) {
	if rp.unit == nil {
		return
	}

	ResetRunes(sim, &rp.bloodRunes, RuneKind_Blood)
	ResetRunes(sim, &rp.frostRunes, RuneKind_Frost)
	ResetRunes(sim, &rp.unholyRunes, RuneKind_Unholy)
}

func (unit *Unit) EnableRunicPowerBar(currentRunicPower float64, maxRunicPower float64,
	onRuneSpend OnRuneSpend,
	onBloodRuneGain OnBloodRuneGain,
	onFrostRuneGain OnFrostRuneGain,
	onUnholyRuneGain OnUnholyRuneGain,
	onDeathRuneGain OnDeathRuneGain,
	onRunicPowerGain OnRunicPowerGain) {
	unit.runicPowerBar = runicPowerBar{
		unit: unit,

		maxRunicPower:     maxRunicPower,
		currentRunicPower: currentRunicPower,

		bloodRunes:  [2]Rune{Rune{state: RuneState_Normal, kind: RuneKind_Blood, pas: [2]*PendingAction{nil, nil}, lastRegenTime: -1}, Rune{state: RuneState_Normal, kind: RuneKind_Blood, pas: [2]*PendingAction{nil, nil}, lastRegenTime: -1}},
		frostRunes:  [2]Rune{Rune{state: RuneState_Normal, kind: RuneKind_Frost, pas: [2]*PendingAction{nil, nil}, lastRegenTime: -1}, Rune{state: RuneState_Normal, kind: RuneKind_Frost, pas: [2]*PendingAction{nil, nil}, lastRegenTime: -1}},
		unholyRunes: [2]Rune{Rune{state: RuneState_Normal, kind: RuneKind_Unholy, pas: [2]*PendingAction{nil, nil}, lastRegenTime: -1}, Rune{state: RuneState_Normal, kind: RuneKind_Unholy, pas: [2]*PendingAction{nil, nil}, lastRegenTime: -1}},

		onRuneSpend:      onRuneSpend,
		onBloodRuneGain:  onBloodRuneGain,
		onFrostRuneGain:  onFrostRuneGain,
		onUnholyRuneGain: onUnholyRuneGain,
		onDeathRuneGain:  onDeathRuneGain,
		onRunicPowerGain: onRunicPowerGain,

		isACopy: false,
	}

	unit.bloodRuneGainMetrics = unit.NewBloodRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionBloodRuneGain, Tag: 1})
	unit.frostRuneGainMetrics = unit.NewFrostRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionFrostRuneGain, Tag: 1})
	unit.unholyRuneGainMetrics = unit.NewUnholyRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionUnholyRuneGain, Tag: 1})
	unit.deathRuneGainMetrics = unit.NewDeathRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionDeathRuneGain, Tag: 1})

	unit.runicPowerBar.unit = unit

}

func (unit *Unit) HasRunicPowerBar() bool {
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

func (rp *runicPowerBar) MaxRunicPower() float64 {
	return rp.maxRunicPower
}

func (rp *runicPowerBar) PercentRunicPower() float64 {
	return rp.currentRunicPower / rp.maxRunicPower
}

func (rp *runicPowerBar) addRunicPowerInterval(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to add negative runic power!")
	}

	newRunicPower := MinFloat(rp.currentRunicPower+amount, rp.maxRunicPower)

	if !rp.isACopy {
		metrics.AddEvent(amount, newRunicPower-rp.currentRunicPower)

		if sim.Log != nil {
			rp.unit.Log(sim, "Gained %0.3f runic power from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, rp.currentRunicPower, newRunicPower)
		}
	}

	rp.currentRunicPower = newRunicPower
}

func (rp *runicPowerBar) AddRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	rp.addRunicPowerInterval(sim, amount, metrics)
	rp.onRunicPowerGain(sim)
}

func (rp *runicPowerBar) spendRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to spend negative runic power!")
	}

	newRunicPower := rp.currentRunicPower - amount

	if !rp.isACopy {
		metrics.AddEvent(-amount, -amount)

		if sim.Log != nil {
			rp.unit.Log(sim, "Spent %0.3f runic power from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, rp.currentRunicPower, newRunicPower)
		}
	}

	rp.currentRunicPower = newRunicPower

}

func (rp *runicPowerBar) CurrentRunesOfType(rb *[2]Rune, runeState RuneState) int32 {
	return TernaryInt32(rb[0].state == runeState, 1, 0) + TernaryInt32(rb[1].state == runeState, 1, 0)
}

func (rp *runicPowerBar) DeathRuneRegenAt(r *Rune) time.Duration {
	readyAt := time.Duration(NeverExpires)

	if r.state == RuneState_DeathSpent {
		if r.pas[1] != nil {
			readyAt = time.Duration(TernaryDuration(r.pas[0].NextActionAt < r.pas[1].NextActionAt, r.pas[0].NextActionAt, readyAt))
		} else {
			readyAt = r.pas[0].NextActionAt
		}
	}

	return readyAt
}

func (rp *runicPowerBar) SpentDeathRuneReadyAt(sim *Simulation) time.Duration {
	readyAt := time.Duration(NeverExpires)

	readyAt = MinDuration(readyAt, rp.DeathRuneRegenAt(&rp.bloodRunes[0]))
	readyAt = MinDuration(readyAt, rp.DeathRuneRegenAt(&rp.bloodRunes[1]))
	readyAt = MinDuration(readyAt, rp.DeathRuneRegenAt(&rp.frostRunes[0]))
	readyAt = MinDuration(readyAt, rp.DeathRuneRegenAt(&rp.frostRunes[1]))
	readyAt = MinDuration(readyAt, rp.DeathRuneRegenAt(&rp.unholyRunes[0]))
	readyAt = MinDuration(readyAt, rp.DeathRuneRegenAt(&rp.unholyRunes[1]))

	return readyAt
}

func (rp *runicPowerBar) DeathRuneReadyAt(sim *Simulation) time.Duration {
	readyAt := time.Duration(NeverExpires)

	if rp.bloodRunes[0].state == RuneState_Death || rp.bloodRunes[1].state == RuneState_Death ||
		rp.frostRunes[0].state == RuneState_Death || rp.frostRunes[1].state == RuneState_Death ||
		rp.unholyRunes[0].state == RuneState_Death || rp.unholyRunes[1].state == RuneState_Death {
		readyAt = sim.CurrentTime
		return readyAt
	}

	return rp.SpentDeathRuneReadyAt(sim)
}

func (rp *runicPowerBar) CurrentRuneGrace(sim *Simulation, runes *[2]Rune) time.Duration {
	if runes[0].pas[0] == nil {
		return time.Millisecond*2500 - MinDuration(2500*time.Millisecond, sim.CurrentTime-runes[0].lastRegenTime)
	} else if runes[1].pas[0] == nil {
		return time.Millisecond*2500 - MinDuration(2500*time.Millisecond, sim.CurrentTime-runes[1].lastRegenTime)
	}
	return 0
}

func (rp *runicPowerBar) IsLeftBloodRuneNormal() bool {
	return (rp.bloodRunes[0].state == RuneState_Normal)
}

func (rp *runicPowerBar) CurrentBloodRuneGrace(sim *Simulation) time.Duration {
	return rp.CurrentRuneGrace(sim, &rp.bloodRunes)
}

func (rp *runicPowerBar) CurrentFrostRuneGrace(sim *Simulation) time.Duration {
	return rp.CurrentRuneGrace(sim, &rp.frostRunes)
}

func (rp *runicPowerBar) CurrentUnholyRuneGrace(sim *Simulation) time.Duration {
	return rp.CurrentRuneGrace(sim, &rp.unholyRunes)
}

func (rp *runicPowerBar) SpentRuneReadyAt(sim *Simulation, runes *[2]Rune) time.Duration {
	readyAt := rp.SpentDeathRuneReadyAt(sim)

	if runes[0].pas[0] != nil {
		readyAt = MinDuration(readyAt, runes[0].pas[0].NextActionAt)
	}

	if runes[1].pas[0] != nil {
		readyAt = MinDuration(readyAt, runes[1].pas[0].NextActionAt)
	}

	return readyAt
}

func (rp *runicPowerBar) RuneReadyAt(sim *Simulation, runes *[2]Rune) time.Duration {
	if runes[0].state == RuneState_Normal || runes[0].state == RuneState_Death ||
		runes[1].state == RuneState_Normal || runes[1].state == RuneState_Death {
		return sim.CurrentTime
	}

	return rp.SpentRuneReadyAt(sim, runes)
}

func (rp *runicPowerBar) SpentBloodRuneReadyAt(sim *Simulation) time.Duration {
	return rp.SpentRuneReadyAt(sim, &rp.bloodRunes)
}

func (rp *runicPowerBar) SpentFrostRuneReadyAt(sim *Simulation) time.Duration {
	return rp.SpentRuneReadyAt(sim, &rp.frostRunes)
}

func (rp *runicPowerBar) SpentUnholyRuneReadyAt(sim *Simulation) time.Duration {
	return rp.SpentRuneReadyAt(sim, &rp.unholyRunes)
}

func (rp *runicPowerBar) BloodRuneReadyAt(sim *Simulation) time.Duration {
	return rp.RuneReadyAt(sim, &rp.bloodRunes)
}

func (rp *runicPowerBar) FrostRuneReadyAt(sim *Simulation) time.Duration {
	return rp.RuneReadyAt(sim, &rp.frostRunes)
}

func (rp *runicPowerBar) UnholyRuneReadyAt(sim *Simulation) time.Duration {
	return rp.RuneReadyAt(sim, &rp.unholyRunes)
}

func (rp *runicPowerBar) AnySpentRuneReadyAt(sim *Simulation) time.Duration {
	return MinDuration(MinDuration(rp.SpentRuneReadyAt(sim, &rp.bloodRunes), rp.SpentRuneReadyAt(sim, &rp.frostRunes)), rp.SpentRuneReadyAt(sim, &rp.unholyRunes))
}

func (rp *runicPowerBar) AnyRuneReadyAt(sim *Simulation) time.Duration {
	return MinDuration(MinDuration(rp.RuneReadyAt(sim, &rp.bloodRunes), rp.RuneReadyAt(sim, &rp.frostRunes)), rp.RuneReadyAt(sim, &rp.unholyRunes))
}

func (rp *runicPowerBar) CurrentBloodRunes() int32 {
	return rp.CurrentRunesOfType(&rp.bloodRunes, RuneState_Normal)
}

func (rp *runicPowerBar) CurrentFrostRunes() int32 {
	return rp.CurrentRunesOfType(&rp.frostRunes, RuneState_Normal)
}

func (rp *runicPowerBar) CurrentUnholyRunes() int32 {
	return rp.CurrentRunesOfType(&rp.unholyRunes, RuneState_Normal)
}

func (rp *runicPowerBar) CurrentDeathRunes() int32 {
	return rp.CurrentRunesOfType(&rp.bloodRunes, RuneState_Death) + rp.CurrentRunesOfType(&rp.frostRunes, RuneState_Death) + rp.CurrentRunesOfType(&rp.unholyRunes, RuneState_Death)
}

func (rp *runicPowerBar) AllRunesOfTypeSpent(runes *[2]Rune) bool {
	return (runes[0].state == RuneState_Spent || runes[0].state == RuneState_DeathSpent) &&
		(runes[1].state == RuneState_Spent || runes[1].state == RuneState_DeathSpent)
}

func (rp *runicPowerBar) AllBloodRunesSpent() bool {
	return rp.AllRunesOfTypeSpent(&rp.bloodRunes)
}

func (rp *runicPowerBar) AllFrostSpent() bool {
	return rp.AllRunesOfTypeSpent(&rp.frostRunes)
}

func (rp *runicPowerBar) AllUnholySpent() bool {
	return rp.AllRunesOfTypeSpent(&rp.unholyRunes)
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

func (rp *runicPowerBar) OptimalRuneCost(cost RuneCost) RuneCost {
	bh := uint8(rp.CurrentBloodRunes())
	fh := uint8(rp.CurrentFrostRunes())
	uh := uint8(rp.CurrentUnholyRunes())
	dh := uint8(rp.CurrentDeathRunes())
	current := NewRuneCost(cost.RunicPower(), bh, fh, uh, dh)
	if current&cost == cost {
		return cost // if we match the cost then we dont need deathrunes
	}

	neededDeath := cost.Death() // just in case death was passed in as a cost.

	newCost := NewRuneCost(cost.RunicPower(), 0, 0, 0, 0)

	if c := cost.Blood(); bh < c {
		neededDeath += c - bh
	} else if c == 1 {
		newCost = newCost | 0b01
	} else if c == 2 {
		newCost = newCost | 0b11
	}

	if c := cost.Frost(); fh < c {
		neededDeath += c - fh
	} else if c == 1 {
		newCost = newCost | 0b0100
	} else if c == 2 {
		newCost = newCost | 0b1100
	}

	if c := cost.Unholy(); uh < c {
		neededDeath += c - uh
	} else if c == 1 {
		newCost = newCost | 0b010000
	} else if c == 2 {
		newCost = newCost | 0b110000
	}

	if neededDeath > dh {
		return 0 // can't cast
	} else if neededDeath == 1 {
		newCost = newCost | 0b01000000
	} else if neededDeath == 2 {
		newCost = newCost | 0b11000000
	}

	return newCost
}

func (rc *RuneAmount) IsValid() bool {
	return rc.Blood >= 0 && rc.Frost >= 0 && rc.Unholy >= 0 && rc.Death >= 0
}

func (rp *runicPowerBar) SpendRuneCost(sim *Simulation, spell *Spell, cost RuneCost) {
	for i := uint8(0); i < cost.Blood(); i++ {
		rp.SpendBloodRune(sim, spell.BloodRuneMetrics())
	}
	for i := uint8(0); i < cost.Frost(); i++ {
		rp.SpendFrostRune(sim, spell.FrostRuneMetrics())
	}
	for i := uint8(0); i < cost.Unholy(); i++ {
		rp.SpendUnholyRune(sim, spell.UnholyRuneMetrics())
	}
	for i := uint8(0); i < cost.Death(); i++ {
		rp.SpendDeathRune(sim, spell.DeathRuneMetrics())
	}
	rpc := cost.RunicPower()
	hasRune := cost.HasRune()
	if rpc <= 0 {
		return
	}
	if !hasRune {
		rp.spendRunicPower(sim, float64(rpc), spell.RunicPowerMetrics())
	} else {
		rp.AddRunicPower(sim, float64(rpc), spell.RunicPowerMetrics())
	}
}

func (rp *runicPowerBar) GainRuneMetrics(sim *Simulation, metrics *ResourceMetrics, name string, currRunes int32, newRunes int32) {
	if !rp.isACopy {
		metrics.AddEvent(1, float64(newRunes)-float64(currRunes))

		if sim.Log != nil {
			rp.unit.Log(sim, "Gained 1.000 %s rune from %s (%d --> %d).", name, metrics.ActionID, currRunes, newRunes)
		}
	}
}

func (rp *runicPowerBar) SpendRuneMetrics(sim *Simulation, metrics *ResourceMetrics, name string, currRunes int32, newRunes int32) {
	if !rp.isACopy {
		metrics.AddEvent(-1, -1)

		if sim.Log != nil {
			rp.unit.Log(sim, "Spent 1.000 %s rune from %s (%d --> %d).", name, metrics.ActionID, currRunes, newRunes)
		}
	}
}

func (rp *runicPowerBar) SetRuneToState(r *Rune, runeState RuneState, runeKind RuneKind) {
	if (r.state == RuneState_Spent || r.state == RuneState_Normal) && (runeState == RuneState_Death || runeState == RuneState_DeathSpent) {
		r.kind = RuneKind_Death
	} else if (r.state == RuneState_DeathSpent || r.state == RuneState_Death) && (runeState != RuneState_Death && runeState != RuneState_DeathSpent) {
		r.kind = runeKind
	}
	r.state = runeState
}

// LastSpentRune gives the slot of the last rune of given type to have been spent.
func (rp *runicPowerBar) LastSpentRuneofType(kind RuneKind) int32 {
	rb := &rp.bloodRunes
	if kind == RuneKind_Frost {
		rb = &rp.frostRunes
	} else if kind == RuneKind_Unholy {
		rb = &rp.unholyRunes
	} else if kind == RuneKind_Death {
		panic("havent implemented finding last spent death rune.")
	}

	// if rune 1 was most recently spent And its the right kind
	//  or if its the only correct kind
	if rb[0].lastSpendTime < rb[1].lastSpendTime || rb[0].kind != kind {
		if rb[1].kind == kind && rb[1].state == RuneState_Spent {
			return 1
		}
	}

	// In this case if rune 0 was the right kind, its the only option left.
	if rb[0].kind == kind && rb[0].state == RuneState_Spent {
		return 0
	}

	// This means no runes of the given kind can be found... what do?
	return -1
}

func (rp *runicPowerBar) SetRuneAtIdxSlotToState(runeBarIdx int32, slot int32, runeState RuneState, runeKind RuneKind) {
	rb := &rp.bloodRunes
	if runeBarIdx == 1 {
		rb = &rp.frostRunes
	} else if runeBarIdx == 2 {
		rb = &rp.unholyRunes
	}

	if (rb[slot].state == RuneState_Spent || rb[slot].state == RuneState_Normal) && (runeState == RuneState_Death || runeState == RuneState_DeathSpent) {
		rb[slot].kind = RuneKind_Death
	} else if (rb[slot].state == RuneState_DeathSpent || rb[slot].state == RuneState_Death) && (runeState != RuneState_Death && runeState != RuneState_DeathSpent) {
		rb[slot].kind = runeKind
	}
	rb[slot].state = runeState
}

func (rp *runicPowerBar) SetRuneAtSlotToState(rb *[2]Rune, slot int32, runeState RuneState, runeKind RuneKind) {
	if (rb[slot].state == RuneState_Spent || rb[slot].state == RuneState_Normal) && (runeState == RuneState_Death || runeState == RuneState_DeathSpent) {
		rb[slot].kind = RuneKind_Death
	} else if (rb[slot].state == RuneState_DeathSpent || rb[slot].state == RuneState_Death) && (runeState != RuneState_Death && runeState != RuneState_DeathSpent) {
		if runeKind == RuneKind_Undef {
			panic("You have to set a rune kind here!")
		}
		rb[slot].kind = runeKind
	}
	rb[slot].state = runeState
}

func (rp *runicPowerBar) RegenRuneAndCancelPAs(sim *Simulation, r *Rune) {
	if r.state == RuneState_Spent {
		r.state = RuneState_Normal

		if r.pas[0] != nil {
			r.lastRegenTime = sim.CurrentTime
			r.pas[0].Cancel(sim)
			r.pas[0] = nil
		}
	} else if r.state == RuneState_DeathSpent {
		r.state = RuneState_Death

		if r.pas[0] != nil {
			r.lastRegenTime = sim.CurrentTime
			r.pas[0].Cancel(sim)
			r.pas[0] = nil
		}
	}
}

func (rp *runicPowerBar) RegenAllRunes(sim *Simulation) {
	startBlood := rp.CurrentBloodRunes()
	startFrost := rp.CurrentFrostRunes()
	startUnholy := rp.CurrentUnholyRunes()
	startDeath := rp.CurrentDeathRunes()

	rp.RegenRuneAndCancelPAs(sim, &rp.bloodRunes[0])
	rp.RegenRuneAndCancelPAs(sim, &rp.bloodRunes[1])
	rp.RegenRuneAndCancelPAs(sim, &rp.frostRunes[0])
	rp.RegenRuneAndCancelPAs(sim, &rp.frostRunes[1])
	rp.RegenRuneAndCancelPAs(sim, &rp.unholyRunes[0])
	rp.RegenRuneAndCancelPAs(sim, &rp.unholyRunes[1])

	if !rp.isACopy {
		if rp.CurrentBloodRunes()-startBlood > 0 {
			rp.GainRuneMetrics(sim, rp.bloodRuneGainMetrics, "blood", startBlood, rp.CurrentBloodRunes())
			rp.onBloodRuneGain(sim)
		}

		if rp.CurrentFrostRunes()-startFrost > 0 {
			rp.GainRuneMetrics(sim, rp.frostRuneGainMetrics, "frost", startFrost, rp.CurrentFrostRunes())
			rp.onFrostRuneGain(sim)
		}

		if rp.CurrentUnholyRunes()-startUnholy > 0 {
			rp.GainRuneMetrics(sim, rp.unholyRuneGainMetrics, "unholy", startUnholy, rp.CurrentUnholyRunes())
			rp.onUnholyRuneGain(sim)
		}

		if rp.CurrentDeathRunes()-startDeath > 0 {
			rp.GainRuneMetrics(sim, rp.deathRuneGainMetrics, "death", startDeath, rp.CurrentDeathRunes())
			rp.onDeathRuneGain(sim)
		}
	}
}

func (rp *runicPowerBar) GenerateRune(sim *Simulation, r *Rune) {
	if r.state == RuneState_Spent {
		if r.kind == RuneKind_Death {
			panic("Rune has wrong type for state.")
		}
		r.state = RuneState_Normal
		r.lastRegenTime = sim.CurrentTime
	} else if r.state == RuneState_DeathSpent {
		if r.kind != RuneKind_Death {
			panic("Rune has wrong type for state.")
		}
		r.state = RuneState_Death
		r.lastRegenTime = sim.CurrentTime
	}
}

func (rp *runicPowerBar) SpendRuneFromType(sim *Simulation, rb *[2]Rune, runeState RuneState) int32 {
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

	rb[slot].lastSpendTime = sim.CurrentTime

	if rp.onRuneSpend != nil {
		rp.onRuneSpend(sim)
	}

	return slot
}

func (rp *runicPowerBar) LaunchRuneRegenPA(sim *Simulation, r *Rune) {
	runeGracePeriod := 0.0
	if r.lastRegenTime != -1 {
		runeGracePeriod = MinFloat(2.5, float64(sim.CurrentTime-r.lastRegenTime)/float64(1*time.Second))
	}
	pa := &PendingAction{
		NextActionAt: sim.CurrentTime + time.Second*time.Duration(10.0-runeGracePeriod),
		Priority:     ActionPriorityRegen,
	}
	pa.OnAction = func(sim *Simulation) {
		if !pa.cancelled {
			r.pas[0] = nil

			currRunes := int32(-1)
			switch r.kind {
			case RuneKind_Blood:
				currRunes = rp.CurrentBloodRunes()
			case RuneKind_Frost:
				currRunes = rp.CurrentFrostRunes()
			case RuneKind_Unholy:
				currRunes = rp.CurrentUnholyRunes()
			case RuneKind_Death:
				currRunes = rp.CurrentDeathRunes()
			}

			rp.GenerateRune(sim, r)

			switch r.kind {
			case RuneKind_Blood:
				rp.GainRuneMetrics(sim, rp.bloodRuneGainMetrics, "blood", currRunes, currRunes+1)
				if !rp.isACopy {
					rp.onBloodRuneGain(sim)
				}
			case RuneKind_Frost:
				rp.GainRuneMetrics(sim, rp.frostRuneGainMetrics, "frost", currRunes, currRunes+1)
				if !rp.isACopy {
					rp.onFrostRuneGain(sim)
				}
			case RuneKind_Unholy:
				rp.GainRuneMetrics(sim, rp.unholyRuneGainMetrics, "unholy", currRunes, currRunes+1)
				if !rp.isACopy {
					rp.onUnholyRuneGain(sim)
				}
			case RuneKind_Death:
				rp.GainRuneMetrics(sim, rp.deathRuneGainMetrics, "death", currRunes, currRunes+1)
				if !rp.isACopy {
					rp.onDeathRuneGain(sim)
				}
			}
		} else {
			r.pas[0] = nil
		}
	}

	r.pas[0] = pa
	if !rp.isACopy {
		sim.AddPendingAction(pa)
	}
}

func (rp *runicPowerBar) SpendBloodRune(sim *Simulation, metrics *ResourceMetrics) int32 {
	currRunes := rp.CurrentBloodRunes()
	if currRunes <= 0 {
		panic("Trying to spend blood runes that don't exist!")
	}

	rp.SpendRuneMetrics(sim, metrics, "blood", currRunes, currRunes-1)
	spendSlot := rp.SpendRuneFromType(sim, &rp.bloodRunes, RuneState_Normal)

	r := &rp.bloodRunes[spendSlot]
	rp.LaunchRuneRegenPA(sim, r)

	return spendSlot
}

func (rp *runicPowerBar) SpendFrostRune(sim *Simulation, metrics *ResourceMetrics) int32 {
	currRunes := rp.CurrentFrostRunes()
	if currRunes <= 0 {
		panic("Trying to spend frost runes that don't exist!")
	}

	rp.SpendRuneMetrics(sim, metrics, "frost", currRunes, currRunes-1)
	spendSlot := rp.SpendRuneFromType(sim, &rp.frostRunes, RuneState_Normal)

	r := &rp.frostRunes[spendSlot]
	rp.LaunchRuneRegenPA(sim, r)

	return spendSlot
}

func (rp *runicPowerBar) SpendUnholyRune(sim *Simulation, metrics *ResourceMetrics) int32 {
	currRunes := rp.CurrentUnholyRunes()
	if currRunes <= 0 {
		panic("Trying to spend unholy runes that don't exist!")
	}

	rp.SpendRuneMetrics(sim, metrics, "unholy", currRunes, currRunes-1)
	spendSlot := rp.SpendRuneFromType(sim, &rp.unholyRunes, RuneState_Normal)

	r := &rp.unholyRunes[spendSlot]
	rp.LaunchRuneRegenPA(sim, r)

	return spendSlot
}

func (rp *runicPowerBar) SpendDeathRune(sim *Simulation, metrics *ResourceMetrics) {
	currRunes := rp.CurrentDeathRunes()
	if currRunes <= 0 {
		panic("Trying to spend death runes that don't exist!")
	}

	rp.SpendRuneMetrics(sim, metrics, "death", currRunes, currRunes-1)

	runeTypeIdx := 0
	spendSlot := rp.SpendRuneFromType(sim, &rp.bloodRunes, RuneState_Death)
	if spendSlot < 0 {
		runeTypeIdx += 1
		spendSlot = rp.SpendRuneFromType(sim, &rp.frostRunes, RuneState_Death)
		if spendSlot < 0 {
			runeTypeIdx += 1
			spendSlot = rp.SpendRuneFromType(sim, &rp.unholyRunes, RuneState_Death)
		}
	}

	r := &rp.bloodRunes[spendSlot]
	if runeTypeIdx == 1 {
		r = &rp.frostRunes[spendSlot]
	} else if runeTypeIdx == 2 {
		r = &rp.unholyRunes[spendSlot]
	}

	rp.LaunchRuneRegenPA(sim, r)
}
