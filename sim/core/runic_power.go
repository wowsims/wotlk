package core

import (
	"fmt"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type OnRune func(sim *Simulation)
type OnRunicPowerGain func(sim *Simulation)

type RuneKind uint8

const (
	RuneKind_Undef RuneKind = iota
	RuneKind_Blood
	RuneKind_Frost
	RuneKind_Unholy
	RuneKind_Death
)

type RuneMeta struct {
	lastRegenTime time.Duration // last time the rune regenerated.
	lastSpendTime time.Duration // last time the rune was spent.
	regenAt       time.Duration // time at which the rune will no longer be spent.
	revertAt      time.Duration // time at which rune will no longer be kind death.
	revertOnSpend bool
}

type RunicPowerBar struct {
	unit *Unit

	maxRunicPower     float64
	currentRunicPower float64

	// These flags are used to simplify pending action checks
	// |BFUDS| |BFUDS| |BFUDS| |BFUDS| |BFUDS| |BFUDS|
	runeStates int32
	runeMeta   [6]RuneMeta
	btslot     int8

	bloodRuneGainMetrics  *ResourceMetrics
	frostRuneGainMetrics  *ResourceMetrics
	unholyRuneGainMetrics *ResourceMetrics
	deathRuneGainMetrics  *ResourceMetrics

	onRuneSpend      OnRune
	onBloodRuneGain  OnRune
	onFrostRuneGain  OnRune
	onUnholyRuneGain OnRune
	onDeathRuneGain  OnRune
	onRunicPowerGain OnRunicPowerGain

	pa *PendingAction

	isACopy bool
	clone   *RunicPowerBar
}

func (rp *RunicPowerBar) DebugString() string {
	data := ""
	for i := int32(0); i < 6; i++ {
		data += fmt.Sprintf("Rune %d - D: %v S: %v\n\tRegenAt: %0.1f, RevertAt: %0.1f\n", i, rp.runeStates&isDeaths[i] != 0, rp.runeStates&isSpents[i] != 0, rp.runeMeta[i].regenAt.Seconds(), rp.runeMeta[i].revertAt.Seconds())
	}

	return data
}

// CopyRunicPowerBar will create a clone of the bar with the same
func (rp *RunicPowerBar) CopyRunicPowerBar() *RunicPowerBar {
	if rp.clone == nil {
		rp.clone = &RunicPowerBar{isACopy: true}
	}

	rp.clone.maxRunicPower = rp.maxRunicPower
	rp.clone.currentRunicPower = rp.currentRunicPower
	rp.clone.runeStates = rp.runeStates
	rp.clone.runeMeta = rp.runeMeta

	return rp.clone
}

func ResetRunes(runeMeta *RuneMeta) {
	runeMeta.regenAt = NeverExpires
	runeMeta.revertAt = NeverExpires
	runeMeta.lastRegenTime = -1
	runeMeta.lastSpendTime = -1
	runeMeta.revertOnSpend = false
}

func (rp *RunicPowerBar) reset(sim *Simulation) {
	if rp.unit == nil {
		return
	}
	if rp.pa != nil {
		rp.pa.Cancel(sim)
	}

	ResetRunes(&rp.runeMeta[0])
	ResetRunes(&rp.runeMeta[1])
	ResetRunes(&rp.runeMeta[2])
	ResetRunes(&rp.runeMeta[3])
	ResetRunes(&rp.runeMeta[4])
	ResetRunes(&rp.runeMeta[5])
	rp.runeStates = baseRuneState // unspent, no death
}

const baseRuneState = 0b100001000001000010000010000100

func (unit *Unit) EnableRunicPowerBar(currentRunicPower float64, maxRunicPower float64,
	onRuneSpend OnRune,
	onBloodRuneGain OnRune,
	onFrostRuneGain OnRune,
	onUnholyRuneGain OnRune,
	onDeathRuneGain OnRune,
	onRunicPowerGain OnRunicPowerGain) {
	unit.RunicPowerBar = RunicPowerBar{
		unit: unit,

		maxRunicPower:     maxRunicPower,
		currentRunicPower: currentRunicPower,

		runeStates: baseRuneState,

		onRuneSpend:      onRuneSpend,
		onBloodRuneGain:  onBloodRuneGain,
		onFrostRuneGain:  onFrostRuneGain,
		onUnholyRuneGain: onUnholyRuneGain,
		onDeathRuneGain:  onDeathRuneGain,
		onRunicPowerGain: onRunicPowerGain,
		isACopy:          false,
	}

	unit.bloodRuneGainMetrics = unit.NewBloodRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionBloodRuneGain, Tag: 1})
	unit.frostRuneGainMetrics = unit.NewFrostRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionFrostRuneGain, Tag: 1})
	unit.unholyRuneGainMetrics = unit.NewUnholyRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionUnholyRuneGain, Tag: 1})
	unit.deathRuneGainMetrics = unit.NewDeathRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionDeathRuneGain, Tag: 1})
}

func (unit *Unit) HasRunicPowerBar() bool {
	return unit.RunicPowerBar.unit != nil
}

func (rp *RunicPowerBar) BloodRuneGainMetrics() *ResourceMetrics {
	return rp.bloodRuneGainMetrics
}

func (rp *RunicPowerBar) FrostRuneGainMetrics() *ResourceMetrics {
	return rp.frostRuneGainMetrics
}

func (rp *RunicPowerBar) UnholyRuneGainMetrics() *ResourceMetrics {
	return rp.unholyRuneGainMetrics
}

func (rp *RunicPowerBar) DeathRuneGainMetrics() *ResourceMetrics {
	return rp.deathRuneGainMetrics
}

func (rp *RunicPowerBar) CurrentRunicPower() float64 {
	return rp.currentRunicPower
}

func (rp *RunicPowerBar) MaxRunicPower() float64 {
	return rp.maxRunicPower
}

func (rp *RunicPowerBar) PercentRunicPower() float64 {
	return rp.currentRunicPower / rp.maxRunicPower
}

func (rp *RunicPowerBar) addRunicPowerInterval(sim *Simulation, amount float64, metrics *ResourceMetrics) {
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

func (rp *RunicPowerBar) AddRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	rp.addRunicPowerInterval(sim, amount, metrics)
	if !rp.isACopy {
		rp.onRunicPowerGain(sim)
	}
}

func (rp *RunicPowerBar) spendRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
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

// DeathRuneRegenAt returns the time the given death rune will regen at.
//  If the rune is not death or not spent it returns NeverExpires
func (rp *RunicPowerBar) DeathRuneRegenAt(slot int32) time.Duration {
	// If not death or not spent, no regen time
	if isSpentDeath[slot]&rp.runeStates != isSpentDeath[slot] {
		return NeverExpires
	}
	return rp.runeMeta[slot].regenAt
}

// DeathRuneRevertAt returns the next time that a death rune will revert.
//  If there is no deathrune that needs to revert it returns `NeverExpires`.
func (rp *RunicPowerBar) DeathRuneRevertAt() time.Duration {
	readyAt := NeverExpires
	for i := int32(0); i < 6; i++ {
		if rp.runeStates&isDeaths[i] == isDeaths[i] {
			readyAt = MinDuration(readyAt, rp.runeMeta[i].revertAt)
		}
	}
	return readyAt
}

func (rp *RunicPowerBar) SpentDeathRuneReadyAt() time.Duration {
	readyAt := NeverExpires
	for i := int32(0); i < 6; i++ {
		readyAt = MinDuration(readyAt, rp.DeathRuneRegenAt(i))
	}
	return readyAt
}

func (rp *RunicPowerBar) CurrentRuneGrace(sim *Simulation, slot int32) time.Duration {
	if rp.runeMeta[slot].lastRegenTime < sim.CurrentTime {
		return time.Millisecond*2500 - MinDuration(2500*time.Millisecond, sim.CurrentTime-rp.runeMeta[slot].lastRegenTime)
	}
	return 0
}

const anyBloodSpent = 0b0000100001
const anyFrostSpent = 0b0000100001 << 10
const anyUnholySpent = 0b0000100001 << 20

func (rp *RunicPowerBar) CurrentBloodRuneGrace(sim *Simulation) time.Duration {
	return MaxDuration(rp.CurrentRuneGrace(sim, 0), rp.CurrentRuneGrace(sim, 1))
}

func (rp *RunicPowerBar) CurrentFrostRuneGrace(sim *Simulation) time.Duration {
	return MaxDuration(rp.CurrentRuneGrace(sim, 2), rp.CurrentRuneGrace(sim, 3))
}

func (rp *RunicPowerBar) CurrentUnholyRuneGrace(sim *Simulation) time.Duration {
	return MaxDuration(rp.CurrentRuneGrace(sim, 2), rp.CurrentRuneGrace(sim, 3))
}

func (rp *RunicPowerBar) NormalSpentBloodRuneReadyAt(sim *Simulation) time.Duration {
	readyAt := NeverExpires
	if rp.runeStates&isDeaths[0] == 0 && rp.runeStates&isSpents[0] != 0 {
		readyAt = rp.runeMeta[0].regenAt
	}
	if rp.runeStates&isDeaths[1] == 0 && rp.runeStates&isSpents[1] != 0 {
		readyAt = MinDuration(readyAt, rp.runeMeta[1].regenAt)
	}
	return readyAt
}

func (rp *RunicPowerBar) NormalSpentFrostRuneReadyAt(sim *Simulation) time.Duration {
	readyAt := NeverExpires
	if rp.runeStates&isDeaths[2] == 0 && rp.runeStates&isSpents[2] != 0 {
		readyAt = rp.runeMeta[2].regenAt
	}
	if rp.runeStates&isDeaths[3] == 0 && rp.runeStates&isSpents[3] != 0 {
		readyAt = MinDuration(readyAt, rp.runeMeta[3].regenAt)
	}
	return readyAt
}

func (rp *RunicPowerBar) NormalFrostRuneReadyAt(sim *Simulation) time.Duration {
	readyAt := NeverExpires
	if rp.runeStates&isDeaths[2] == 0 && rp.runeStates&isSpents[2] != 0 {
		readyAt = rp.runeMeta[2].regenAt
	}
	if rp.runeStates&isDeaths[3] == 0 && rp.runeStates&isSpents[3] != 0 {
		readyAt = MinDuration(readyAt, rp.runeMeta[3].regenAt)
	}
	if (rp.runeStates&isDeaths[2] == 0 && rp.runeStates&isSpents[2] == 0) || (rp.runeStates&isDeaths[3] == 0 && rp.runeStates&isSpents[3] == 0) {
		readyAt = sim.CurrentTime
	}
	return readyAt
}

func (rp *RunicPowerBar) NormalSpentUnholyRuneReadyAt(sim *Simulation) time.Duration {
	readyAt := NeverExpires
	if rp.runeStates&isDeaths[4] == 0 && rp.runeStates&isSpents[4] != 0 {
		readyAt = rp.runeMeta[4].regenAt
	}
	if rp.runeStates&isDeaths[5] == 0 && rp.runeStates&isSpents[5] != 0 {
		readyAt = MinDuration(readyAt, rp.runeMeta[5].regenAt)
	}
	return readyAt
}

func (rp *RunicPowerBar) NormalUnholyRuneReadyAt(sim *Simulation) time.Duration {
	readyAt := NeverExpires
	if rp.runeStates&isDeaths[4] == 0 && rp.runeStates&isSpents[4] != 0 {
		readyAt = rp.runeMeta[4].regenAt
	}
	if rp.runeStates&isDeaths[5] == 0 && rp.runeStates&isSpents[5] != 0 {
		readyAt = MinDuration(readyAt, rp.runeMeta[5].regenAt)
	}
	if (rp.runeStates&isDeaths[4] == 0 && rp.runeStates&isSpents[4] == 0) || (rp.runeStates&isDeaths[5] == 0 && rp.runeStates&isSpents[5] == 0) {
		readyAt = sim.CurrentTime
	}
	return readyAt
}

func (rp *RunicPowerBar) SpentBloodRuneReadyAt() time.Duration {
	return MinDuration(rp.runeMeta[0].regenAt, rp.runeMeta[1].regenAt)
}

func (rp *RunicPowerBar) SpentFrostRuneReadyAt() time.Duration {
	return MinDuration(rp.runeMeta[2].regenAt, rp.runeMeta[3].regenAt)
}

func (rp *RunicPowerBar) SpentUnholyRuneReadyAt() time.Duration {
	return MinDuration(rp.runeMeta[4].regenAt, rp.runeMeta[5].regenAt)
}

func (rp *RunicPowerBar) BloodRuneReadyAt(sim *Simulation) time.Duration {
	if rp.runeStates&anyBloodSpent != anyBloodSpent { // if any are not spent
		return sim.CurrentTime
	}
	return MinDuration(rp.runeMeta[0].regenAt, rp.runeMeta[1].regenAt)
}

func (rp *RunicPowerBar) FrostRuneReadyAt(sim *Simulation) time.Duration {
	if rp.runeStates&anyFrostSpent != anyFrostSpent { // if any are not spent
		return sim.CurrentTime
	}
	return MinDuration(rp.runeMeta[2].regenAt, rp.runeMeta[3].regenAt)
}

func (rp *RunicPowerBar) UnholyRuneReadyAt(sim *Simulation) time.Duration {
	if rp.runeStates&anyUnholySpent != anyUnholySpent { // if any are not spent
		return sim.CurrentTime
	}
	return MinDuration(rp.runeMeta[4].regenAt, rp.runeMeta[5].regenAt)
}

// AnySpentRuneReadyAt returns the next time that a rune will regenerate.
//  It will be NeverExpires if there is no rune pending regeneration.
func (rp *RunicPowerBar) AnySpentRuneReadyAt() time.Duration {
	return MinDuration(MinDuration(rp.SpentBloodRuneReadyAt(), rp.SpentFrostRuneReadyAt()), rp.SpentUnholyRuneReadyAt())
}

func (rp *RunicPowerBar) AnyRuneReadyAt(sim *Simulation) time.Duration {
	return MinDuration(MinDuration(rp.BloodRuneReadyAt(sim), rp.FrostRuneReadyAt(sim)), rp.UnholyRuneReadyAt(sim))
}

// ConvertFromDeath reverts the rune to its original type.
func (rp *RunicPowerBar) ConvertFromDeath(sim *Simulation, slot int8) {
	rp.runeStates = ^isDeaths[slot] & rp.runeStates
	rp.runeMeta[slot].revertAt = NeverExpires
	rp.runeMeta[slot].revertOnSpend = false

	if !rp.isACopy && rp.runeStates&isSpents[slot] == 0 {
		metrics := rp.bloodRuneGainMetrics
		onGain := rp.onBloodRuneGain
		if slot == 2 || slot == 3 {
			metrics = rp.frostRuneGainMetrics
			onGain = rp.onFrostRuneGain
		} else if slot == 4 || slot == 5 {
			metrics = rp.unholyRuneGainMetrics
			onGain = rp.onUnholyRuneGain
		}
		rp.SpendRuneMetrics(sim, rp.deathRuneGainMetrics, 1)
		rp.GainRuneMetrics(sim, metrics, 1)
		onGain(sim)
	}
}

// ConvertToDeath converts the given slot to death and sets up the revertion conditions
// ConvertToDeath converts the given slot to death and sets up the revertion conditions
func (rp *RunicPowerBar) ConvertToDeath(sim *Simulation, slot int8, revertOnSpend bool, revertAt time.Duration) {
	if slot == -1 {
		return
	}
	rp.runeStates = rp.runeStates | isDeaths[slot]

	// revertOnSpend == true overrides anything
	rp.runeMeta[slot].revertOnSpend = rp.runeMeta[slot].revertOnSpend || revertOnSpend

	if rp.runeMeta[slot].revertOnSpend {
		rp.runeMeta[slot].revertAt = NeverExpires
	} else {
		if rp.runeMeta[slot].revertAt != NeverExpires {
			rp.runeMeta[slot].revertAt = MaxDuration(rp.runeMeta[slot].revertAt, revertAt)
		} else {
			rp.runeMeta[slot].revertAt = revertAt
		}
	}

	// Note we gained
	if !rp.isACopy {
		metrics := rp.bloodRuneGainMetrics
		if slot == 2 || slot == 3 {
			metrics = rp.frostRuneGainMetrics
		} else if slot == 4 || slot == 5 {
			metrics = rp.unholyRuneGainMetrics
		}
		if rp.runeStates&isSpents[slot] == 0 {
			// Only lose/gain if it wasn't spent (which it should be at this point)
			rp.SpendRuneMetrics(sim, metrics, 1)
			rp.GainRuneMetrics(sim, rp.deathRuneGainMetrics, 1)
			rp.onDeathRuneGain(sim)
		}
	}
}

func (rp *RunicPowerBar) LeftBloodRuneReady() bool {
	const unspentBlood1 = isDeath | isSpent
	if rp.runeStates&unspentBlood1 == 0 {
		return true
	} else {
		return false
	}
}

func (rp *RunicPowerBar) RightBloodRuneReady() bool {
	const unspentBlood1 = isDeath | isSpent
	const unspentBlood2 = unspentBlood1 << 5
	if rp.runeStates&unspentBlood2 == 0 {
		return true
	} else {
		return false
	}
}

func (rp *RunicPowerBar) RuneIsDeath(slot int8) bool {
	return (rp.runeStates & isDeaths[slot]) != 0
}

func (rp *RunicPowerBar) CurrentBloodRunes() int8 {
	const unspentBlood1 = isDeath | isSpent
	const unspentBlood2 = unspentBlood1 << 5

	var count int8
	if rp.runeStates&unspentBlood1 == 0 {
		count++
	}
	if rp.runeStates&unspentBlood2 == 0 {
		count++
	}

	return count
}

func (rp *RunicPowerBar) CurrentFrostRunes() int8 {
	const unspentFrost1 = (isDeath | isSpent) << 10
	const unspentFrost2 = unspentFrost1 << 5

	var count int8
	if rp.runeStates&unspentFrost1 == 0 {
		count++
	}
	if rp.runeStates&unspentFrost2 == 0 {
		count++
	}

	return count
}

func (rp *RunicPowerBar) CurrentUnholyRunes() int8 {
	const unspentUnholy1 = (isDeath | isSpent) << 20
	const unspentUnholy2 = unspentUnholy1 << 5

	var count int8
	if rp.runeStates&unspentUnholy1 == 0 {
		count++
	}
	if rp.runeStates&unspentUnholy2 == 0 {
		count++
	}

	return count
}

func (rp *RunicPowerBar) CurrentDeathRunes() int8 {
	var count int8
	for i := range rp.runeMeta {
		if rp.runeStates&isDeaths[i] != 0 && rp.runeStates&isSpents[i] == 0 {
			count++
		}
	}
	return count
}

func (rp *RunicPowerBar) NormalCurrentBloodRunes() int32 {
	const unspentBlood1 = isSpent
	const unspentBlood2 = unspentBlood1 << 5

	var count int32
	if rp.runeStates&unspentBlood1 == 0 {
		count++
	}
	if rp.runeStates&unspentBlood2 == 0 {
		count++
	}

	return count
}

func (rp *RunicPowerBar) NormalCurrentFrostRunes() int32 {
	const unspentFrost1 = (isSpent) << 10
	const unspentFrost2 = unspentFrost1 << 5

	var count int32
	if rp.runeStates&unspentFrost1 == 0 {
		count++
	}
	if rp.runeStates&unspentFrost2 == 0 {
		count++
	}

	return count
}

func (rp *RunicPowerBar) NormalCurrentUnholyRunes() int32 {
	const unspentUnholy1 = (isSpent) << 20
	const unspentUnholy2 = unspentUnholy1 << 5

	var count int32
	if rp.runeStates&unspentUnholy1 == 0 {
		count++
	}
	if rp.runeStates&unspentUnholy2 == 0 {
		count++
	}

	return count
}

func (rp *RunicPowerBar) AllBloodRunesSpent() bool {
	const checkBloodSpent = isSpent & (isSpent << 5)
	return rp.runeStates&checkBloodSpent == checkBloodSpent
}

func (rp *RunicPowerBar) AllFrostSpent() bool {
	const checkFrostSpent = (isSpent << 10) & (isSpent << 15)
	return rp.runeStates&checkFrostSpent == checkFrostSpent
}

func (rp *RunicPowerBar) AllUnholySpent() bool {
	const checkUnholySpent = (isSpent << 20) & (isSpent << 25)
	return rp.runeStates&checkUnholySpent == checkUnholySpent
}

func (rp *RunicPowerBar) CastCostPossible(sim *Simulation, runicPowerAmount float64, bloodAmount int8, frostAmount int8, unholyAmount int8) bool {
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

func (rp *RunicPowerBar) OptimalRuneCost(cost RuneCost) RuneCost {
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

func (rp *RunicPowerBar) SpendRuneCost(sim *Simulation, spell *Spell, cost RuneCost) (int8, int8, int8) {
	r := [3]int8{-1, -1, -1}
	idx := 0

	for i := uint8(0); i < cost.Blood(); i++ {
		r[idx] = rp.SpendBloodRune(sim, spell.BloodRuneMetrics())
		idx++
	}
	for i := uint8(0); i < cost.Frost(); i++ {
		r[idx] = rp.SpendFrostRune(sim, spell.FrostRuneMetrics())
		idx++
	}
	for i := uint8(0); i < cost.Unholy(); i++ {
		r[idx] = rp.SpendUnholyRune(sim, spell.UnholyRuneMetrics())
		idx++
	}
	for i := uint8(0); i < cost.Death(); i++ {
		r[idx] = rp.SpendDeathRune(sim, spell.DeathRuneMetrics())
		idx++
	}
	rpc := cost.RunicPower()
	hasRune := cost.HasRune()
	if rpc <= 0 {
		return r[0], r[1], r[2]
	}
	if !hasRune {
		rp.spendRunicPower(sim, float64(rpc), spell.RunicPowerMetrics())
	} else {
		rp.AddRunicPower(sim, float64(rpc), spell.RunicPowerMetrics())
	}
	return r[0], r[1], r[2]
}

// GainRuneMetrics should be called after gaining the rune
func (rp *RunicPowerBar) GainRuneMetrics(sim *Simulation, metrics *ResourceMetrics, gainAmount int8) {
	if !rp.isACopy {
		metrics.AddEvent(float64(gainAmount), float64(gainAmount))

		if sim.Log != nil {

			var name string
			var currRunes int8

			switch metrics.Type {
			case proto.ResourceType_ResourceTypeDeathRune:
				name = "death"
				currRunes = rp.CurrentDeathRunes()
			case proto.ResourceType_ResourceTypeBloodRune:
				name = "blood"
				currRunes = rp.CurrentBloodRunes()
			case proto.ResourceType_ResourceTypeFrostRune:
				name = "frost"
				currRunes = rp.CurrentFrostRunes()
			case proto.ResourceType_ResourceTypeUnholyRune:
				name = "unholy"
				currRunes = rp.CurrentUnholyRunes()
			default:
				panic("invalid metrics for rune gaining")
			}

			rp.unit.Log(sim, "Gained %0.3f %s rune from %s (%d --> %d).", float64(gainAmount), name, metrics.ActionID, currRunes-gainAmount, currRunes)
		}
	}
}

// SpendRuneMetrics should be called after spending the rune
func (rp *RunicPowerBar) SpendRuneMetrics(sim *Simulation, metrics *ResourceMetrics, spendAmount int8) {
	if !rp.isACopy {
		metrics.AddEvent(-float64(spendAmount), -float64(spendAmount))

		if sim.Log != nil {
			var name string
			var currRunes int8

			switch metrics.Type {
			case proto.ResourceType_ResourceTypeDeathRune:
				name = "death"
				currRunes = rp.CurrentDeathRunes()
			case proto.ResourceType_ResourceTypeBloodRune:
				name = "blood"
				currRunes = rp.CurrentBloodRunes()
			case proto.ResourceType_ResourceTypeFrostRune:
				name = "frost"
				currRunes = rp.CurrentFrostRunes()
			case proto.ResourceType_ResourceTypeUnholyRune:
				name = "unholy"
				currRunes = rp.CurrentUnholyRunes()
			default:
				panic("invalid metrics for rune spending")
			}

			rp.unit.Log(sim, "Spent 1.000 %s rune from %s (%d --> %d).", name, metrics.ActionID, currRunes+spendAmount, currRunes)
		}
	}
}

func (rp *RunicPowerBar) BloodRuneSpentAt(dur time.Duration) int32 {
	if rp.runeMeta[0].lastSpendTime == dur {
		return 0
	}
	if rp.runeMeta[1].lastSpendTime == dur {
		return 1
	}
	return -1
}

func (rp *RunicPowerBar) FrostRuneSpentAt(dur time.Duration) int32 {
	if rp.runeMeta[2].lastSpendTime == dur {
		return 2
	}
	if rp.runeMeta[3].lastSpendTime == dur {
		return 3
	}
	return -1
}

func (rp *RunicPowerBar) UnholyRuneSpentAt(dur time.Duration) int32 {
	if rp.runeMeta[4].lastSpendTime == dur {
		return 4
	}
	if rp.runeMeta[5].lastSpendTime == dur {
		return 5
	}
	return -1
}

func (rp *RunicPowerBar) RegenRune(regenAt time.Duration, slot int32) {
	checkSpent := isSpents[slot]
	if checkSpent&rp.runeStates > 0 {
		rp.runeStates = ^checkSpent & rp.runeStates // unset spent flag for this rune.
		rp.runeMeta[slot].lastRegenTime = regenAt
		rp.runeMeta[slot].regenAt = NeverExpires
	}
}

func (rp *RunicPowerBar) RegenAllRunes(sim *Simulation) {
	startBlood := rp.CurrentBloodRunes()
	startFrost := rp.CurrentFrostRunes()
	startUnholy := rp.CurrentUnholyRunes()
	startDeath := rp.CurrentDeathRunes()

	rp.RegenRune(sim.CurrentTime, 0)
	rp.RegenRune(sim.CurrentTime, 1)
	rp.RegenRune(sim.CurrentTime, 2)
	rp.RegenRune(sim.CurrentTime, 3)
	rp.RegenRune(sim.CurrentTime, 4)
	rp.RegenRune(sim.CurrentTime, 5)

	if !rp.isACopy {
		if rp.CurrentBloodRunes()-startBlood > 0 {
			rp.GainRuneMetrics(sim, rp.bloodRuneGainMetrics, rp.CurrentBloodRunes()-startBlood)
			rp.onBloodRuneGain(sim)
		}

		if rp.CurrentFrostRunes()-startFrost > 0 {
			rp.GainRuneMetrics(sim, rp.frostRuneGainMetrics, rp.CurrentFrostRunes()-startFrost)
			rp.onFrostRuneGain(sim)
		}

		if rp.CurrentUnholyRunes()-startUnholy > 0 {
			rp.GainRuneMetrics(sim, rp.unholyRuneGainMetrics, rp.CurrentUnholyRunes()-startUnholy)
			rp.onUnholyRuneGain(sim)
		}

		if rp.CurrentDeathRunes()-startDeath > 0 {
			rp.GainRuneMetrics(sim, rp.deathRuneGainMetrics, rp.CurrentDeathRunes()-startDeath)
			rp.onDeathRuneGain(sim)
		}
	}
}

func (rp *RunicPowerBar) SpendRuneFromKind(sim *Simulation, rkind RuneKind) int8 {
	var rb int8
	if rkind == RuneKind_Frost {
		rb = 2
	} else if rkind == RuneKind_Unholy {
		rb = 4
	} else if rkind == RuneKind_Death {
		panic("use 'ReadyRuneByKind' to find death rune")
	}
	spent1 := isSpents[rb]
	spent2 := isSpents[rb+1]

	spendable1 := spent1 | isDeaths[rb] // verify rune is not spent and not death
	spendable2 := spent2 | isDeaths[rb+1]

	slot := int8(-1)
	// Figure out which rune is spendable (not death and not spent)
	// Then mark the spend bit for that rune.
	if rp.runeStates&spendable1 == 0 {
		rp.runeStates = rp.runeStates | spent1
		slot = rb
	} else if rp.runeStates&spendable2 == 0 {
		rp.runeStates = rp.runeStates | spent2
		slot = rb + 1
	} else {
		panic("Trying to spend rune that does not exist!")
	}

	rp.runeMeta[slot].lastSpendTime = sim.CurrentTime

	if rp.onRuneSpend != nil {
		rp.onRuneSpend(sim)
	}

	return slot
}

func (rp *RunicPowerBar) LaunchRuneRegen(sim *Simulation, slot int8) {
	var runeGracePeriod time.Duration
	if rp.runeMeta[slot].lastRegenTime != -1 {
		runeGracePeriod = MinDuration(time.Millisecond*2500, sim.CurrentTime-rp.runeMeta[slot].lastRegenTime)
	}
	rp.runeMeta[slot].regenAt = sim.CurrentTime + (time.Second*10 - runeGracePeriod)

	rp.launchPA(sim, rp.runeMeta[slot].regenAt)
}

func (rp *RunicPowerBar) launchPA(sim *Simulation, at time.Duration) {
	if rp.isACopy {
		return
	}
	if rp.pa != nil {
		// If this new regen is before currently scheduled one, we must cancel old regen and start a new one.
		if rp.pa.NextActionAt > at {
			rp.pa.Cancel(sim)
			rp.pa = nil
		} else {
			return
		}
	}
	pa := &PendingAction{
		NextActionAt: at,
		Priority:     ActionPriorityRegen,
	}
	pa.OnAction = func(sim *Simulation) {
		if !pa.cancelled {
			// regenerate and revert
			rp.Advance(sim, sim.CurrentTime)

			// Check when we need next check
			pa.NextActionAt = MinDuration(rp.AnySpentRuneReadyAt(), rp.DeathRuneRevertAt())
			if pa.NextActionAt < NeverExpires {
				sim.AddPendingAction(pa)
			}
		}
	}
	rp.pa = pa
	sim.AddPendingAction(pa)

}

// Constants for finding runes

//                |BFUDSBFUDSBFUDSBFUDSBFUDSBFUDS|
const checkDeath = 0b000100001000010000100001000010
const checkSpent = 0b000010000100001000010000100001

const isDeath = int32(0b00010)
const isSpent = int32(0b00001)

var isDeaths = [6]int32{
	isDeath,
	isDeath << 5,
	isDeath << 10,
	isDeath << 15,
	isDeath << 20,
	isDeath << 25,
}

var isSpents = [6]int32{
	isSpent,
	isSpent << 5,
	isSpent << 10,
	isSpent << 15,
	isSpent << 20,
	isSpent << 25,
}

var isSpentDeath = [6]int32{
	(isDeath | isSpent),
	(isDeath | isSpent) << 5,
	(isDeath | isSpent) << 10,
	(isDeath | isSpent) << 15,
	(isDeath | isSpent) << 20,
	(isDeath | isSpent) << 25,
}

func (rp *RunicPowerBar) Advance(sim *Simulation, newTime time.Duration) {
	if rp.runeStates&checkDeath > 0 {
		for i := int8(0); i < 6; i++ {
			if rp.runeMeta[i].revertAt <= newTime {
				if rp.btslot == i {
					rp.btslot = -1 // this was the BT slot.
				}
				rp.ConvertFromDeath(sim, i)
			}
		}
	}

	if rp.runeStates&checkSpent > 0 {
		rp.findAndRegen(sim, newTime)
	}
}

func (rp *RunicPowerBar) TryRegenRune(sim *Simulation, newTime time.Duration, slot int32) {
	if rp.runeMeta[slot].regenAt > newTime {
		return
	}
	if rp.runeStates&isSpents[slot] == 0 {
		return
	}

	metrics := rp.bloodRuneGainMetrics
	onGain := rp.onBloodRuneGain
	if rp.runeStates&(isDeaths[slot]) > 0 {
		metrics = rp.deathRuneGainMetrics
		onGain = rp.onDeathRuneGain
	} else if slot == 2 || slot == 3 {
		metrics = rp.frostRuneGainMetrics
		onGain = rp.onFrostRuneGain
	} else if slot == 4 || slot == 5 {
		metrics = rp.unholyRuneGainMetrics
		onGain = rp.onUnholyRuneGain
	}
	rp.RegenRune(newTime, slot)
	if !rp.isACopy {
		rp.GainRuneMetrics(sim, metrics, 1)
		onGain(sim)
	}
}

func (rp *RunicPowerBar) findAndRegen(sim *Simulation, newTime time.Duration) {
	rp.TryRegenRune(sim, newTime, 0)
	rp.TryRegenRune(sim, newTime, 1)
	rp.TryRegenRune(sim, newTime, 2)
	rp.TryRegenRune(sim, newTime, 3)
	rp.TryRegenRune(sim, newTime, 4)
	rp.TryRegenRune(sim, newTime, 5)
}

func (rp *RunicPowerBar) SpendBloodRune(sim *Simulation, metrics *ResourceMetrics) int8 {
	currRunes := rp.CurrentBloodRunes()
	if currRunes <= 0 {
		panic("Trying to spend blood runes that don't exist!")
	}

	spendSlot := rp.SpendRuneFromKind(sim, RuneKind_Blood)
	rp.SpendRuneMetrics(sim, metrics, 1)
	rp.LaunchRuneRegen(sim, spendSlot)
	return spendSlot
}

func (rp *RunicPowerBar) SpendFrostRune(sim *Simulation, metrics *ResourceMetrics) int8 {
	currRunes := rp.CurrentFrostRunes()
	if currRunes <= 0 {
		panic("Trying to spend frost runes that don't exist!")
	}

	spendSlot := rp.SpendRuneFromKind(sim, RuneKind_Frost)
	rp.SpendRuneMetrics(sim, metrics, 1)
	rp.LaunchRuneRegen(sim, spendSlot)
	return spendSlot
}

func (rp *RunicPowerBar) SpendUnholyRune(sim *Simulation, metrics *ResourceMetrics) int8 {
	currRunes := rp.CurrentUnholyRunes()
	if currRunes <= 0 {
		panic("Trying to spend unholy runes that don't exist!")
	}

	spendSlot := rp.SpendRuneFromKind(sim, RuneKind_Unholy)
	rp.SpendRuneMetrics(sim, metrics, 1)
	rp.LaunchRuneRegen(sim, spendSlot)
	return spendSlot
}

// ReadyDeathRune returns the slot of first available death rune.
//  Returns -1 if there are no ready death runes
func (rp *RunicPowerBar) ReadyDeathRune() int8 {
	for i := int8(0); i < 6; i++ {
		if rp.runeStates&isDeaths[i] != 0 && rp.runeStates&isSpents[i] == 0 {
			return i
		}
	}
	return -1
}

func (rp *RunicPowerBar) SpendDeathRune(sim *Simulation, metrics *ResourceMetrics) int8 {
	if rp.runeStates&checkDeath == 0 {
		panic("Trying to spend death runes that don't exist!")
	}

	slot := rp.ReadyDeathRune()
	if rp.runeMeta[slot].revertOnSpend {
		// disable revert at
		rp.runeMeta[slot].revertOnSpend = false
		rp.runeMeta[slot].revertAt = NeverExpires
		// clear death bit to revert.
		rp.runeStates = ^isDeaths[slot] & rp.runeStates
	}

	// mark spent bit to spend
	rp.runeStates = rp.runeStates | isSpents[slot]
	rp.runeMeta[slot].lastSpendTime = sim.CurrentTime

	rp.SpendRuneMetrics(sim, metrics, 1)
	rp.LaunchRuneRegen(sim, slot)
	return slot
}
