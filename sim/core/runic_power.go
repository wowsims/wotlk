package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type OnRune func(sim *Simulation)
type OnRunicPowerGain func(sim *Simulation)

type RuneMeta struct {
	lastRegenTime time.Duration // last time the rune regenerated.
	lastSpendTime time.Duration // last time the rune was spent.
	regenAt       time.Duration // time at which the rune will no longer be spent.
	revertAt      time.Duration // time at which rune will no longer be kind death.
}

type RunicPowerBar struct {
	unit *Unit

	maxRunicPower     float64
	currentRunicPower float64
	runeCD            time.Duration

	// These flags are used to simplify pending action checks
	// |DS|DS|DS|DS|DS|DS|
	runeStates int16
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

// Constants for finding runes
// |DS|DS|DS|DS|DS|DS|
const (
	baseRuneState = 0 // unspent, no death

	allDeath = 0b101010101010
	allSpent = 0b010101010101

	anyBloodSpent  = 0b0101 << 0
	anyFrostSpent  = 0b0101 << 4
	anyUnholySpent = 0b0101 << 8
)

var (
	isDeaths     = [6]int16{0b10 << 0, 0b10 << 2, 0b10 << 4, 0b10 << 6, 0b10 << 8, 0b10 << 10}
	isSpents     = [6]int16{0b01 << 0, 0b01 << 2, 0b01 << 4, 0b01 << 6, 0b01 << 8, 0b01 << 10}
	isSpentDeath = [6]int16{0b11 << 0, 0b11 << 2, 0b11 << 4, 0b11 << 6, 0b11 << 8, 0b11 << 10}
)

func (rp *RunicPowerBar) DebugString() string {
	ss := make([]string, len(rp.runeMeta))
	for i := range rp.runeMeta {
		ss[i] += fmt.Sprintf("Rune %d - D: %v S: %v\n\tRegenAt: %0.1f, RevertAt: %0.1f", i, rp.runeStates&isDeaths[i] != 0, rp.runeStates&isSpents[i] != 0, rp.runeMeta[i].regenAt.Seconds(), rp.runeMeta[i].revertAt.Seconds())
	}
	return strings.Join(ss, "\n")
}

// CopyRunicPowerBar will create a clone of the bar with the same rune state
func (rp *RunicPowerBar) CopyRunicPowerBar() *RunicPowerBar {
	if rp.clone == nil {
		rp.clone = &RunicPowerBar{isACopy: true}
	}

	rp.clone.maxRunicPower = rp.maxRunicPower
	rp.clone.currentRunicPower = rp.currentRunicPower
	rp.clone.runeCD = rp.runeCD
	rp.clone.runeStates = rp.runeStates
	rp.clone.runeMeta = rp.runeMeta

	return rp.clone
}

func resetRunes(runeMeta *RuneMeta) {
	runeMeta.regenAt = NeverExpires
	runeMeta.revertAt = NeverExpires
	runeMeta.lastRegenTime = -1
	runeMeta.lastSpendTime = -1
}

func (rp *RunicPowerBar) reset(sim *Simulation) {
	if rp.unit == nil {
		return
	}
	if rp.pa != nil {
		rp.pa.Cancel(sim)
	}

	resetRunes(&rp.runeMeta[0])
	resetRunes(&rp.runeMeta[1])
	resetRunes(&rp.runeMeta[2])
	resetRunes(&rp.runeMeta[3])
	resetRunes(&rp.runeMeta[4])
	resetRunes(&rp.runeMeta[5])
	rp.runeStates = baseRuneState
}

func (unit *Unit) EnableRunicPowerBar(currentRunicPower float64, maxRunicPower float64, runeCD time.Duration,
	onRuneSpend OnRune,
	onBloodRuneGain OnRune,
	onFrostRuneGain OnRune,
	onUnholyRuneGain OnRune,
	onDeathRuneGain OnRune,
	onRunicPowerGain OnRunicPowerGain) {
	unit.SetCurrentPowerBar(RunicPower)
	unit.RunicPowerBar = RunicPowerBar{
		unit: unit,

		maxRunicPower:     maxRunicPower,
		currentRunicPower: currentRunicPower,
		runeCD:            runeCD,

		runeStates: baseRuneState,

		onRuneSpend:      onRuneSpend,
		onBloodRuneGain:  onBloodRuneGain,
		onFrostRuneGain:  onFrostRuneGain,
		onUnholyRuneGain: onUnholyRuneGain,
		onDeathRuneGain:  onDeathRuneGain,
		onRunicPowerGain: onRunicPowerGain,
		isACopy:          false,
		btslot:           -1,
	}

	unit.bloodRuneGainMetrics = unit.NewBloodRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionBloodRuneGain, Tag: 1})
	unit.frostRuneGainMetrics = unit.NewFrostRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionFrostRuneGain, Tag: 1})
	unit.unholyRuneGainMetrics = unit.NewUnholyRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionUnholyRuneGain, Tag: 1})
	unit.deathRuneGainMetrics = unit.NewDeathRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionDeathRuneGain, Tag: 1})
}

func (unit *Unit) HasRunicPowerBar() bool {
	return unit.RunicPowerBar.unit != nil
}

func (rp *RunicPowerBar) SetRuneCd(runeCd time.Duration) {
	rp.runeCD = runeCd
}

func (rp *RunicPowerBar) CurrentRunicPower() float64 {
	return rp.currentRunicPower
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
// If the rune is not death or not spent it returns NeverExpires
func (rp *RunicPowerBar) DeathRuneRegenAt(slot int32) time.Duration {
	// If not death or not spent, no regen time
	if isSpentDeath[slot]&rp.runeStates != isSpentDeath[slot] {
		return NeverExpires
	}
	return rp.runeMeta[slot].regenAt
}

// DeathRuneRevertAt returns the next time that a death rune will revert.
// If there is no death rune that needs to revert it returns `NeverExpires`.
func (rp *RunicPowerBar) DeathRuneRevertAt() time.Duration {
	readyAt := NeverExpires
	for i := int32(0); i < 6; i++ {
		if rp.runeStates&isDeaths[i] == isDeaths[i] {
			readyAt = MinDuration(readyAt, rp.runeMeta[i].revertAt)
		}
	}
	return readyAt
}

func (rp *RunicPowerBar) RuneGraceRemaining(sim *Simulation, slot int8) time.Duration {
	lastRegenTime := rp.runeMeta[slot].lastRegenTime

	// pre-pull casts should not get rune-grace
	if sim.CurrentTime <= 0 || lastRegenTime <= 0 {
		return 0
	}

	if lastRegenTime < sim.CurrentTime {
		return time.Millisecond*2500 - MinDuration(2500*time.Millisecond, sim.CurrentTime-lastRegenTime)
	}
	return 0
}

func (rp *RunicPowerBar) CurrentRuneGrace(sim *Simulation, slot int8) time.Duration {
	lastRegenTime := rp.runeMeta[slot].lastRegenTime

	// pre-pull casts should not get rune-grace
	if sim.CurrentTime <= 0 || lastRegenTime <= 0 {
		return 0
	}

	if lastRegenTime < sim.CurrentTime {
		return MinDuration(2500*time.Millisecond, sim.CurrentTime-lastRegenTime)
	}
	return 0
}

func (rp *RunicPowerBar) CurrentBloodRuneGrace(sim *Simulation) time.Duration {
	return MaxDuration(rp.CurrentRuneGrace(sim, 0), rp.CurrentRuneGrace(sim, 1))
}

func (rp *RunicPowerBar) CurrentFrostRuneGrace(sim *Simulation) time.Duration {
	return MaxDuration(rp.CurrentRuneGrace(sim, 2), rp.CurrentRuneGrace(sim, 3))
}

func (rp *RunicPowerBar) CurrentUnholyRuneGrace(sim *Simulation) time.Duration {
	return MaxDuration(rp.CurrentRuneGrace(sim, 4), rp.CurrentRuneGrace(sim, 5))
}

func (rp *RunicPowerBar) FrostRuneGraceRemaining(sim *Simulation) time.Duration {
	return MaxDuration(rp.RuneGraceRemaining(sim, 2), rp.RuneGraceRemaining(sim, 3))
}

func (rp *RunicPowerBar) UnholyRuneGraceRemaining(sim *Simulation) time.Duration {
	return MaxDuration(rp.RuneGraceRemaining(sim, 4), rp.RuneGraceRemaining(sim, 5))
}

func (rp *RunicPowerBar) normalSpentRuneReadyAt(slot int8) time.Duration {
	readyAt := NeverExpires
	if t := rp.runeMeta[slot].regenAt; t < readyAt && rp.runeStates&isSpentDeath[slot] == isSpents[slot] {
		readyAt = t
	}
	if t := rp.runeMeta[slot+1].regenAt; t < readyAt && rp.runeStates&isSpentDeath[slot+1] == isSpents[slot+1] {
		readyAt = t
	}
	return readyAt
}

// NormalSpentBloodRuneReadyAt returns the earliest time a spent non-death blood rune is ready.
func (rp *RunicPowerBar) NormalSpentBloodRuneReadyAt(_ *Simulation) time.Duration {
	return rp.normalSpentRuneReadyAt(0)
}

func (rp *RunicPowerBar) normalRuneReadyAt(sim *Simulation, slot int8) time.Duration {
	if rp.runeStates&isSpentDeath[slot] == 0 || rp.runeStates&isSpentDeath[slot+1] == 0 {
		return sim.CurrentTime
	}
	return rp.normalSpentRuneReadyAt(slot)
}

// NormalFrostRuneReadyAt returns the earliest time a non-death frost rune is ready.
func (rp *RunicPowerBar) NormalFrostRuneReadyAt(sim *Simulation) time.Duration {
	return rp.normalRuneReadyAt(sim, 2)
}

func (rp *RunicPowerBar) NormalUnholyRuneReadyAt(sim *Simulation) time.Duration {
	return rp.normalRuneReadyAt(sim, 4)
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

func (rp *RunicPowerBar) BloodDeathRuneBothReadyAt() time.Duration {
	if rp.runeStates&isDeaths[0] != 0 && rp.runeStates&isDeaths[1] != 0 {
		if MaxDuration(rp.runeMeta[0].regenAt, rp.runeMeta[1].regenAt) > 150000000*time.Minute {
			return MinDuration(rp.runeMeta[0].regenAt, rp.runeMeta[1].regenAt)
		} else {
			return MaxDuration(rp.runeMeta[0].regenAt, rp.runeMeta[1].regenAt)
		}
	} else {
		return -1
	}
}

func (rp *RunicPowerBar) RuneReadyAt(sim *Simulation, slot int8) time.Duration {
	if rp.runeStates&isSpents[slot] != isSpents[slot] {
		return sim.CurrentTime
	}
	return rp.runeMeta[slot].regenAt
}

func (rp *RunicPowerBar) SpendRuneReadyAt(slot int8, spendAt time.Duration) time.Duration {
	runeGraceDuration := rp.RuneGraceAt(slot, spendAt)
	return spendAt + (rp.runeCD - runeGraceDuration)
}

// BloodRuneReadyAt returns the earliest time a (possibly death-converted) blood rune is ready.
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

func (rp *RunicPowerBar) NextRuneTypeReadyAt(sim *Simulation, left int8, right int8) time.Duration {
	if rp.runeStates&isSpents[left] != isSpents[left] && rp.runeStates&isSpents[right] != isSpents[right] {
		// Both are ready so return current time
		return sim.CurrentTime
	} else if rp.runeStates&isSpents[left] == isSpents[left] || rp.runeStates&isSpents[right] == isSpents[right] {
		// One is spent so return the time it will regen at
		return MinDuration(rp.runeMeta[left].regenAt, rp.runeMeta[right].regenAt)
	}
	// Both are spent so return the last one to regen at
	return MaxDuration(rp.runeMeta[left].regenAt, rp.runeMeta[right].regenAt)
}

// TODO Can possibly replaced by either BloodRuneReadyAt() or SpentBloodRuneReadyAt() variants, depending on semantics.
func (rp *RunicPowerBar) NextBloodRuneReadyAt(sim *Simulation) time.Duration {
	return rp.NextRuneTypeReadyAt(sim, 0, 1)
}

func (rp *RunicPowerBar) NextFrostRuneReadyAt(sim *Simulation) time.Duration {
	return rp.NextRuneTypeReadyAt(sim, 2, 3)
}

func (rp *RunicPowerBar) NextUnholyRuneReadyAt(sim *Simulation) time.Duration {
	return rp.NextRuneTypeReadyAt(sim, 4, 5)
}

// AnySpentRuneReadyAt returns the next time that a rune will regenerate.
// It will be NeverExpires if there is no rune pending regeneration.
func (rp *RunicPowerBar) AnySpentRuneReadyAt() time.Duration {
	return MinDuration(MinDuration(rp.SpentBloodRuneReadyAt(), rp.SpentFrostRuneReadyAt()), rp.SpentUnholyRuneReadyAt())
}

func (rp *RunicPowerBar) AnyRuneReadyAt(sim *Simulation) time.Duration {
	return MinDuration(MinDuration(rp.BloodRuneReadyAt(sim), rp.FrostRuneReadyAt(sim)), rp.UnholyRuneReadyAt(sim))
}

// ConvertFromDeath reverts the rune to its original type.
func (rp *RunicPowerBar) ConvertFromDeath(sim *Simulation, slot int8) {
	rp.runeStates ^= isDeaths[slot]
	rp.runeMeta[slot].revertAt = NeverExpires

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
		rp.spendRuneMetrics(sim, rp.deathRuneGainMetrics, 1)
		rp.gainRuneMetrics(sim, metrics, 1)
		onGain(sim)
	}
}

// ConvertToDeath converts the given slot to death and sets up the reversion conditions
func (rp *RunicPowerBar) ConvertToDeath(sim *Simulation, slot int8, revertAt time.Duration) {
	if slot == -1 {
		return
	}
	rp.runeStates |= isDeaths[slot]

	if rp.btslot != slot {
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
			rp.spendRuneMetrics(sim, metrics, 1)
			rp.gainRuneMetrics(sim, rp.deathRuneGainMetrics, 1)
			rp.onDeathRuneGain(sim)
		}
	}
}

func (rp *RunicPowerBar) LeftBloodRuneReady() bool {
	return rp.runeStates&isSpents[0] == 0
}

func (rp *RunicPowerBar) RightBloodRuneReady() bool {
	return rp.runeStates&isSpents[1] == 0
}

func (rp *RunicPowerBar) RuneIsActive(slot int8) bool {
	return rp.runeStates&isSpents[slot] == 0
}

func (rp *RunicPowerBar) RuneIsDeath(slot int8) bool {
	return rp.runeStates&isDeaths[slot] != 0
}

// rune state to count of non-death, non-spent runes (0b00)
var rs2c = []int8{
	0b0000: 2, 0b0001: 1, 0b0010: 1, 0b0011: 1, 0b0100: 1, 0b0101: 0, 0b0110: 0, 0b0111: 0,
	0b1000: 1, 0b1001: 0, 0b1010: 0, 0b1011: 0, 0b1100: 1, 0b1101: 0, 0b1110: 0, 0b1111: 0,
}

func (rp *RunicPowerBar) CurrentBloodRunes() int8 {
	return rs2c[rp.runeStates&0b1111]
}

func (rp *RunicPowerBar) CurrentFrostRunes() int8 {
	return rs2c[(rp.runeStates>>4)&0b1111]
}

func (rp *RunicPowerBar) CurrentUnholyRunes() int8 {
	return rs2c[(rp.runeStates>>8)&0b1111]
}

// rune state to count of death, non-spent runes (0b10)
var rs2d = []int8{
	0b0000: 0, 0b0001: 0, 0b0010: 1, 0b0011: 0, 0b0100: 0, 0b0101: 0, 0b0110: 1, 0b0111: 0,
	0b1000: 1, 0b1001: 1, 0b1010: 2, 0b1011: 1, 0b1100: 0, 0b1101: 0, 0b1110: 1, 0b1111: 0,
}

func (rp *RunicPowerBar) CurrentDeathRunes() int8 {
	return rs2d[rp.runeStates&0b1111] + rs2d[(rp.runeStates>>4)&0b1111] + rs2d[(rp.runeStates>>8)&0b1111]
}

func (rp *RunicPowerBar) DeathRunesInFU() int8 {
	return rs2d[(rp.runeStates>>4)&0b1111] + rs2d[(rp.runeStates>>8)&0b1111]
}

func (rp *RunicPowerBar) AllRunesSpent() bool {
	return rp.runeStates&allSpent == allSpent
}

func (rp *RunicPowerBar) OptimalRuneCost(cost RuneCost) RuneCost {
	var b, f, u, d int8

	if b = cost.Blood(); b > 0 {
		if cb := rp.CurrentBloodRunes(); cb < b {
			d += b - cb
			b = cb
		}
	}

	if f = cost.Frost(); f > 0 {
		if cf := rp.CurrentFrostRunes(); cf < f {
			d += f - cf
			f = cf
		}
	}

	if u = cost.Unholy(); u > 0 {
		if cu := rp.CurrentUnholyRunes(); cu < u {
			d += u - cu
			u = cu
		}
	}

	if d == 0 {
		return cost
	}

	d += cost.Death()

	if cd := rp.CurrentDeathRunes(); cd >= d {
		return NewRuneCost(cost.RunicPower(), b, f, u, d)
	}

	return 0
}

func (rp *RunicPowerBar) SpendRuneCost(sim *Simulation, spell *Spell, cost RuneCost) (int8, int8, int8) {
	if !cost.HasRune() {
		if rpc := cost.RunicPower(); rpc > 0 {
			rp.spendRunicPower(sim, float64(cost.RunicPower()), spell.RunicPowerMetrics())
		}
		return -1, -1, -1
	}

	slots := [3]int8{-1, -1, -1}
	idx := 0
	for i := int8(0); i < cost.Blood(); i++ {
		slots[idx] = rp.spendRune(sim, 0, spell.BloodRuneMetrics())
		idx++
	}
	for i := int8(0); i < cost.Frost(); i++ {
		slots[idx] = rp.spendRune(sim, 2, spell.FrostRuneMetrics())
		idx++
	}
	for i := int8(0); i < cost.Unholy(); i++ {
		slots[idx] = rp.spendRune(sim, 4, spell.UnholyRuneMetrics())
		idx++
	}
	for i := int8(0); i < cost.Death(); i++ {
		slots[idx] = rp.spendDeathRune(sim, spell.DeathRuneMetrics())
		idx++
	}

	if rpc := cost.RunicPower(); rpc > 0 {
		rp.AddRunicPower(sim, float64(rpc), spell.RunicPowerMetrics())
	}
	return slots[0], slots[1], slots[2]
}

// gainRuneMetrics should be called after gaining the rune
func (rp *RunicPowerBar) gainRuneMetrics(sim *Simulation, metrics *ResourceMetrics, gainAmount int8) {
	if rp.isACopy {
		return
	}

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

// spendRuneMetrics should be called after spending the rune
func (rp *RunicPowerBar) spendRuneMetrics(sim *Simulation, metrics *ResourceMetrics, spendAmount int8) {
	if rp.isACopy {
		return
	}

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

func (rp *RunicPowerBar) regenRune(sim *Simulation, regenAt time.Duration, slot int8) {
	if rp.runeStates&isSpents[slot] == 0 {
		return
	}

	rp.runeStates ^= isSpents[slot] // unset spent flag for this rune.
	rp.runeMeta[slot].lastRegenTime = regenAt
	rp.runeMeta[slot].regenAt = NeverExpires

	if !rp.isACopy {
		metrics := rp.bloodRuneGainMetrics
		onGain := rp.onBloodRuneGain
		if rp.runeStates&isDeaths[slot] > 0 {
			metrics = rp.deathRuneGainMetrics
			onGain = rp.onDeathRuneGain
		} else if slot == 2 || slot == 3 {
			metrics = rp.frostRuneGainMetrics
			onGain = rp.onFrostRuneGain
		} else if slot == 4 || slot == 5 {
			metrics = rp.unholyRuneGainMetrics
			onGain = rp.onUnholyRuneGain
		}

		rp.gainRuneMetrics(sim, metrics, 1)
		onGain(sim)
	}
}

func (rp *RunicPowerBar) RegenAllRunes(sim *Simulation) {
	rp.regenRune(sim, sim.CurrentTime, 0)
	rp.regenRune(sim, sim.CurrentTime, 1)
	rp.regenRune(sim, sim.CurrentTime, 2)
	rp.regenRune(sim, sim.CurrentTime, 3)
	rp.regenRune(sim, sim.CurrentTime, 4)
	rp.regenRune(sim, sim.CurrentTime, 5)
}

func (rp *RunicPowerBar) RuneGraceAt(slot int8, at time.Duration) (runeGraceDuration time.Duration) {
	lastRegenTime := rp.runeMeta[slot].lastRegenTime
	// pre-pull casts should not get rune-grace
	if at <= 0 || lastRegenTime <= 0 {
		return 0
	}
	if lastRegenTime != -1 {
		runeGraceDuration = MinDuration(time.Millisecond*2500, at-lastRegenTime)
	}
	return runeGraceDuration
}

func (rp *RunicPowerBar) launchRuneRegen(sim *Simulation, slot int8) {
	runeGracePeriod := rp.RuneGraceAt(slot, sim.CurrentTime)
	rp.runeMeta[slot].regenAt = sim.CurrentTime + (rp.runeCD - runeGracePeriod)

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

func (rp *RunicPowerBar) Advance(sim *Simulation, newTime time.Duration) {
	if rp.runeStates&allDeath > 0 {
		for i := int8(0); i < int8(len(rp.runeMeta)); i++ {
			if rp.runeMeta[i].revertAt <= newTime {
				if rp.btslot == i {
					rp.btslot = -1 // this was the BT slot.
				}
				rp.ConvertFromDeath(sim, i)
			}
		}
	}

	if rp.runeStates&allSpent > 0 {
		rp.findAndRegen(sim, newTime)
	}
}

func (rp *RunicPowerBar) tryRegenRune(sim *Simulation, newTime time.Duration, slot int8) {
	if rp.runeMeta[slot].regenAt > newTime {
		return
	}
	if rp.runeStates&isSpents[slot] == 0 {
		return
	}
	rp.regenRune(sim, newTime, slot)
}

func (rp *RunicPowerBar) findAndRegen(sim *Simulation, newTime time.Duration) {
	rp.tryRegenRune(sim, newTime, 0)
	rp.tryRegenRune(sim, newTime, 1)
	rp.tryRegenRune(sim, newTime, 2)
	rp.tryRegenRune(sim, newTime, 3)
	rp.tryRegenRune(sim, newTime, 4)
	rp.tryRegenRune(sim, newTime, 5)
}

func (rp *RunicPowerBar) spendRune(sim *Simulation, firstSlot int8, metrics *ResourceMetrics) int8 {
	slot := rp.findReadyRune(firstSlot)
	rp.runeStates |= isSpents[slot]

	rp.runeMeta[slot].lastSpendTime = sim.CurrentTime

	if rp.onRuneSpend != nil {
		rp.onRuneSpend(sim)
	}

	rp.spendRuneMetrics(sim, metrics, 1)
	rp.launchRuneRegen(sim, slot)
	return slot
}

func (rp *RunicPowerBar) findReadyRune(slot int8) int8 {
	if rp.runeStates&isSpentDeath[slot] == 0 {
		return slot
	}
	if rp.runeStates&isSpentDeath[slot+1] == 0 {
		return slot + 1
	}
	panic(fmt.Sprintf("findReadyRune(%d) - no slot found (runeStates = %12b)", slot, rp.runeStates))
}

func (rp *RunicPowerBar) spendDeathRune(sim *Simulation, metrics *ResourceMetrics) int8 {
	slot := rp.findReadyDeathRune()
	if rp.btslot != slot {
		rp.runeMeta[slot].revertAt = NeverExpires // disable revert at
		rp.runeStates ^= isDeaths[slot]           // clear death bit to revert.
	}

	// mark spent bit to spend
	rp.runeStates |= isSpents[slot]
	rp.runeMeta[slot].lastSpendTime = sim.CurrentTime

	rp.spendRuneMetrics(sim, metrics, 1)
	rp.launchRuneRegen(sim, slot)
	return slot
}

// findReadyDeathRune returns the slot of first available death rune.
func (rp *RunicPowerBar) findReadyDeathRune() int8 {
	for _, slot := range []int8{4, 5, 2, 3, 0, 1} { // Death runes are spent in the order Unholy -> Frost -> Blood in-game...
		if rp.runeStates&isSpentDeath[slot] == isDeaths[slot] {
			return slot
		}
	}
	panic(fmt.Sprintf("findReadyDeathRune() - no slot found (runeStates = %12b)", rp.runeStates))
}

func (rp *RunicPowerBar) IsBloodTappedRune(slot int8) bool {
	return slot == rp.btslot
}

type RuneCostOptions struct {
	BloodRuneCost  int8
	FrostRuneCost  int8
	UnholyRuneCost int8
	RunicPowerCost float64
	RunicPowerGain float64
	Refundable     bool
}

type RuneCostImpl struct {
	BloodRuneCost  int8
	FrostRuneCost  int8
	UnholyRuneCost int8
	RunicPowerCost float64
	RunicPowerGain float64
	Refundable     bool

	runicPowerMetrics *ResourceMetrics
	bloodRuneMetrics  *ResourceMetrics
	frostRuneMetrics  *ResourceMetrics
	unholyRuneMetrics *ResourceMetrics
	deathRuneMetrics  *ResourceMetrics
}

func newRuneCost(spell *Spell, options RuneCostOptions) *RuneCostImpl {
	baseCost := float64(NewRuneCost(int8(options.RunicPowerCost), options.BloodRuneCost, options.FrostRuneCost, options.UnholyRuneCost, 0))
	spell.DefaultCast.Cost = baseCost
	spell.CurCast.Cost = baseCost

	return &RuneCostImpl{
		BloodRuneCost:  options.BloodRuneCost,
		FrostRuneCost:  options.FrostRuneCost,
		UnholyRuneCost: options.UnholyRuneCost,
		RunicPowerCost: options.RunicPowerCost,
		RunicPowerGain: options.RunicPowerGain,
		Refundable:     options.Refundable,

		runicPowerMetrics: Ternary(options.RunicPowerCost > 0 || options.RunicPowerGain > 0, spell.Unit.NewRunicPowerMetrics(spell.ActionID), nil),
		bloodRuneMetrics:  Ternary(options.BloodRuneCost > 0, spell.Unit.NewBloodRuneMetrics(spell.ActionID), nil),
		frostRuneMetrics:  Ternary(options.FrostRuneCost > 0, spell.Unit.NewFrostRuneMetrics(spell.ActionID), nil),
		unholyRuneMetrics: Ternary(options.UnholyRuneCost > 0, spell.Unit.NewUnholyRuneMetrics(spell.ActionID), nil),
		deathRuneMetrics:  spell.Unit.NewDeathRuneMetrics(spell.ActionID),
	}
}

func (rc *RuneCostImpl) MeetsRequirement(spell *Spell) bool {
	spell.CurCast.Cost *= spell.CostMultiplier // TODO this looks fishy - multiplying and rune costs don't go well together

	cost := RuneCost(spell.CurCast.Cost)
	if cost == 0 {
		return true
	}

	if !cost.HasRune() {
		if float64(cost.RunicPower()) > spell.Unit.CurrentRunicPower() {
			return false
		}
	}

	optCost := spell.Unit.OptimalRuneCost(cost)
	if optCost == 0 { // no combo of runes to fulfill cost
		return false
	}
	spell.CurCast.Cost = float64(optCost) // assign chosen runes to the cost
	return true
}

func (rc *RuneCostImpl) LogCostFailure(sim *Simulation, spell *Spell) {
	spell.Unit.Log(sim, "Failed casting %s, not enough RP or runes.", spell.ActionID)
}

func (rc *RuneCostImpl) SpendCost(sim *Simulation, spell *Spell) {
	// Spend now if there is no way to refund the spell
	if !rc.Refundable {
		spell.Unit.SpendRuneCost(sim, spell, RuneCost(spell.CurCast.Cost))
	}
	if rc.RunicPowerGain > 0 && spell.CurCast.Cost > 0 {
		spell.Unit.AddRunicPower(sim, rc.RunicPowerGain, spell.RunicPowerMetrics())
	}
}

func (rc *RuneCostImpl) spendRefundableCost(sim *Simulation, spell *Spell, result *SpellResult) {
	cost := RuneCost(spell.CurCast.Cost) // cost was already optimized in RuneSpell.Cast
	if cost == 0 {
		return // it was free this time. we don't care
	}
	if result.Landed() {
		spell.Unit.SpendRuneCost(sim, spell, cost)
	}
}

func (spell *Spell) SpendRefundableCost(sim *Simulation, result *SpellResult) {
	spell.Cost.(*RuneCostImpl).spendRefundableCost(sim, spell, result)
}

func (rc *RuneCostImpl) spendRefundableCostAndConvertBloodRune(sim *Simulation, spell *Spell, result *SpellResult, convertChance float64) {
	cost := RuneCost(spell.CurCast.Cost) // cost was already optimized in RuneSpell.Cast
	if cost == 0 {
		return // it was free this time. we don't care
	}
	if !result.Landed() {
		// misses just don't get spent as a way to avoid having to cancel regeneration PAs
		return
	}
	slot1, slot2, _ := spell.Unit.SpendRuneCost(sim, spell, cost)
	if !sim.Proc(convertChance, "Blood of The North / Reaping / DRM") {
		return
	}

	for _, slot := range []int8{slot1, slot2} {
		if slot == 0 || slot == 1 {
			// If the slot to be converted is already blood-tapped, then we convert the other blood rune
			if spell.Unit.IsBloodTappedRune(slot) {
				otherRune := (slot + 1) % 2
				spell.Unit.ConvertToDeath(sim, otherRune, NeverExpires)
			} else {
				spell.Unit.ConvertToDeath(sim, slot, NeverExpires)
			}
		}
	}
}

func (spell *Spell) SpendRefundableCostAndConvertBloodRune(sim *Simulation, result *SpellResult, convertChance float64) {
	spell.Cost.(*RuneCostImpl).spendRefundableCostAndConvertBloodRune(sim, spell, result, convertChance)
}

func (rc *RuneCostImpl) spendRefundableCostAndConvertFrostOrUnholyRune(sim *Simulation, spell *Spell, result *SpellResult, convertChance float64) {
	cost := RuneCost(spell.CurCast.Cost) // cost was already optimized in RuneSpell.Cast
	if cost == 0 {
		return // it was free this time. we don't care
	}
	if !result.Landed() {
		// misses just don't get spent as a way to avoid having to cancel regeneration PAs
		return
	}
	slot1, slot2, slot3 := spell.Unit.SpendRuneCost(sim, spell, cost)
	if !sim.Proc(convertChance, "Blood of The North / Reaping / DRM") {
		return
	}

	for _, slot := range []int8{slot1, slot2, slot3} {
		if slot == 2 || slot == 3 || slot == 4 || slot == 5 {
			spell.Unit.ConvertToDeath(sim, slot, NeverExpires)
		}
	}
}

func (spell *Spell) SpendRefundableCostAndConvertFrostOrUnholyRune(sim *Simulation, result *SpellResult, convertChance float64) {
	spell.Cost.(*RuneCostImpl).spendRefundableCostAndConvertFrostOrUnholyRune(sim, spell, result, convertChance)
}

func (rc *RuneCostImpl) IssueRefund(_ *Simulation, _ *Spell) {
	// Instead of issuing refunds we just don't charge the cost of spells which
	// miss; this is better for perf since we'd have to cancel the regen actions.
}

func (spell *Spell) RunicPowerMetrics() *ResourceMetrics {
	return spell.Cost.(*RuneCostImpl).runicPowerMetrics
}

func (spell *Spell) BloodRuneMetrics() *ResourceMetrics {
	return spell.Cost.(*RuneCostImpl).bloodRuneMetrics
}

func (spell *Spell) FrostRuneMetrics() *ResourceMetrics {
	return spell.Cost.(*RuneCostImpl).frostRuneMetrics
}

func (spell *Spell) UnholyRuneMetrics() *ResourceMetrics {
	return spell.Cost.(*RuneCostImpl).unholyRuneMetrics
}

func (spell *Spell) DeathRuneMetrics() *ResourceMetrics {
	return spell.Cost.(*RuneCostImpl).deathRuneMetrics
}
