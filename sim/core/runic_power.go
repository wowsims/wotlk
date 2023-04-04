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

func (rp *RunicPowerBar) Print() {
	fmt.Print(rp.DebugString())
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
	rp.clone.runeCD = rp.runeCD
	rp.clone.runeStates = rp.runeStates
	rp.clone.runeMeta = rp.runeMeta

	return rp.clone
}

func ResetRunes(runeMeta *RuneMeta) {
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

	ResetRunes(&rp.runeMeta[0])
	ResetRunes(&rp.runeMeta[1])
	ResetRunes(&rp.runeMeta[2])
	ResetRunes(&rp.runeMeta[3])
	ResetRunes(&rp.runeMeta[4])
	ResetRunes(&rp.runeMeta[5])
	rp.runeStates = baseRuneState // unspent, no death
}

const baseRuneState = int16(0)

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
//
//	If the rune is not death or not spent it returns NeverExpires
func (rp *RunicPowerBar) DeathRuneRegenAt(slot int32) time.Duration {
	// If not death or not spent, no regen time
	if isSpentDeath[slot]&rp.runeStates != isSpentDeath[slot] {
		return NeverExpires
	}
	return rp.runeMeta[slot].regenAt
}

// DeathRuneRevertAt returns the next time that a death rune will revert.
//
//	If there is no deathrune that needs to revert it returns `NeverExpires`.
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

func (rp *RunicPowerBar) RuneGraceRemaining(sim *Simulation, slot int32) time.Duration {
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

const anyBloodSpent = 0b0101
const anyFrostSpent = 0b0101 << 4
const anyUnholySpent = 0b0101 << 8

func (rp *RunicPowerBar) CurrentBloodRuneGrace(sim *Simulation) time.Duration {
	return MaxDuration(rp.RuneGraceRemaining(sim, 0), rp.RuneGraceRemaining(sim, 1))
}

func (rp *RunicPowerBar) CurrentFrostRuneGrace(sim *Simulation) time.Duration {
	return MaxDuration(rp.RuneGraceRemaining(sim, 2), rp.RuneGraceRemaining(sim, 3))
}

func (rp *RunicPowerBar) CurrentUnholyRuneGrace(sim *Simulation) time.Duration {
	return MaxDuration(rp.RuneGraceRemaining(sim, 4), rp.RuneGraceRemaining(sim, 5))
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
//
//	It will be NeverExpires if there is no rune pending regeneration.
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
			rp.SpendRuneMetrics(sim, metrics, 1)
			rp.GainRuneMetrics(sim, rp.deathRuneGainMetrics, 1)
			rp.onDeathRuneGain(sim)
		}
	}
}

func (rp *RunicPowerBar) LeftBloodRuneReady() bool {
	const unspentBlood1 = isSpent
	if rp.runeStates&unspentBlood1 == 0 {
		return true
	} else {
		return false
	}
}

func (rp *RunicPowerBar) RightBloodRuneReady() bool {
	const unspentBlood1 = isSpent
	const unspentBlood2 = unspentBlood1 << 2
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
	const unspentBlood2 = unspentBlood1 << 2

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
	const unspentFrost1 = (isDeath | isSpent) << 4
	const unspentFrost2 = unspentFrost1 << 2

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
	const unspentUnholy1 = (isDeath | isSpent) << 8
	const unspentUnholy2 = unspentUnholy1 << 2

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

func (rp *RunicPowerBar) DeathRunesInFU() int8 {
	var count int8
	for i := 2; i < len(rp.runeMeta); i++ {
		if rp.runeStates&isDeaths[i] != 0 {
			count++
		}
	}
	return count
}

func (rp *RunicPowerBar) NormalCurrentBloodRunes() int32 {
	const unspentBlood1 = isSpent
	const unspentBlood2 = unspentBlood1 << 2

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
	const unspentFrost1 = (isSpent) << 4
	const unspentFrost2 = unspentFrost1 << 2

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
	const unspentUnholy1 = (isSpent) << 8
	const unspentUnholy2 = unspentUnholy1 << 2

	var count int32
	if rp.runeStates&unspentUnholy1 == 0 {
		count++
	}
	if rp.runeStates&unspentUnholy2 == 0 {
		count++
	}

	return count
}

func (rp *RunicPowerBar) NormalCurrentRunes() (int32, int32, int32) {
	return rp.NormalCurrentBloodRunes(), rp.NormalCurrentFrostRunes(), rp.NormalCurrentUnholyRunes()
}
func (rp *RunicPowerBar) AllRunesSpent() bool {
	const allSpent = isSpent | (isSpent << 2) | (isSpent << 4) | (isSpent << 6) | (isSpent << 8) | (isSpent << 10)
	return rp.runeStates&allSpent == allSpent
}

func (rp *RunicPowerBar) AllBloodRunesSpent() bool {
	const checkBloodSpent = isSpent | (isSpent << 2)
	return rp.runeStates&checkBloodSpent == checkBloodSpent
}

func (rp *RunicPowerBar) AllFrostSpent() bool {
	const checkFrostSpent = (isSpent << 4) | (isSpent << 6)
	return rp.runeStates&checkFrostSpent == checkFrostSpent
}

func (rp *RunicPowerBar) AllUnholySpent() bool {
	const checkUnholySpent = (isSpent << 8) | (isSpent << 10)
	return rp.runeStates&checkUnholySpent == checkUnholySpent
}

func (rp *RunicPowerBar) CastCostPossible(sim *Simulation, runicPowerAmount float64, bloodAmount int8, frostAmount int8, unholyAmount int8) bool {
	if rp.CurrentRunicPower() < runicPowerAmount {
		return false
	}

	var deficit int8
	if d := bloodAmount - rp.CurrentBloodRunes(); d > 0 {
		deficit += d
	}
	if d := frostAmount - rp.CurrentFrostRunes(); d > 0 {
		deficit += d
	}
	if d := unholyAmount - rp.CurrentUnholyRunes(); d > 0 {
		deficit += d
	}
	return deficit <= rp.CurrentDeathRunes()
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
		newCost |= 0b01
	} else if c == 2 {
		newCost |= 0b11
	}

	if c := cost.Frost(); fh < c {
		neededDeath += c - fh
	} else if c == 1 {
		newCost |= 0b0100
	} else if c == 2 {
		newCost |= 0b1100
	}

	if c := cost.Unholy(); uh < c {
		neededDeath += c - uh
	} else if c == 1 {
		newCost |= 0b010000
	} else if c == 2 {
		newCost |= 0b110000
	}

	if neededDeath > dh {
		return 0 // can't cast
	} else if neededDeath == 1 {
		newCost |= 0b01000000
	} else if neededDeath == 2 {
		newCost |= 0b11000000
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

func (rp *RunicPowerBar) RegenRune(sim *Simulation, regenAt time.Duration, slot int8) {
	checkSpent := isSpents[slot]
	if checkSpent&rp.runeStates > 0 {
		rp.runeStates = ^checkSpent & rp.runeStates // unset spent flag for this rune.
		rp.runeMeta[slot].lastRegenTime = regenAt
		rp.runeMeta[slot].regenAt = NeverExpires

		if !rp.isACopy {
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

			rp.GainRuneMetrics(sim, metrics, 1)
			onGain(sim)
		}
	}
}

func (rp *RunicPowerBar) RegenAllRunes(sim *Simulation) {
	rp.RegenRune(sim, sim.CurrentTime, 0)
	rp.RegenRune(sim, sim.CurrentTime, 1)
	rp.RegenRune(sim, sim.CurrentTime, 2)
	rp.RegenRune(sim, sim.CurrentTime, 3)
	rp.RegenRune(sim, sim.CurrentTime, 4)
	rp.RegenRune(sim, sim.CurrentTime, 5)
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
		rp.runeStates |= spent1
		slot = rb
	} else if rp.runeStates&spendable2 == 0 {
		rp.runeStates |= spent2
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

func (rp *RunicPowerBar) LaunchRuneRegen(sim *Simulation, slot int8) {
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

// Constants for finding runes

// |DS|DS|DS|DS|DS|DS|
const checkDeath = int16(0b101010101010)
const checkSpent = int16(0b010101010101)

const isDeath = int16(0b10)
const isSpent = int16(0b01)

var isDeaths = [6]int16{
	isDeath,
	isDeath << 2,
	isDeath << 4,
	isDeath << 6,
	isDeath << 8,
	isDeath << 10,
}

var isSpents = [6]int16{
	isSpent,
	isSpent << 2,
	isSpent << 4,
	isSpent << 6,
	isSpent << 8,
	isSpent << 10,
}

var isSpentDeath = [6]int16{
	(isDeath | isSpent),
	(isDeath | isSpent) << 2,
	(isDeath | isSpent) << 4,
	(isDeath | isSpent) << 6,
	(isDeath | isSpent) << 8,
	(isDeath | isSpent) << 10,
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

func (rp *RunicPowerBar) TryRegenRune(sim *Simulation, newTime time.Duration, slot int8) {
	if rp.runeMeta[slot].regenAt > newTime {
		return
	}
	if rp.runeStates&isSpents[slot] == 0 {
		return
	}

	rp.RegenRune(sim, newTime, slot)
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
//
//	Returns -1 if there are no ready death runes
func (rp *RunicPowerBar) ReadyDeathRune() int8 {
	// Death runes are spent in the order Unholy -> Frost -> Blood in-game...
	for runeType := int8(2); runeType >= 0; runeType-- {
		for i := runeType * 2; i < (runeType+1)*2; i++ {
			if rp.runeStates&isDeaths[i] != 0 && rp.runeStates&isSpents[i] == 0 {
				return i
			}
		}
	}
	return -1
}

func (rp *RunicPowerBar) IsBloodTappedRune(slot int8) bool {
	return slot == rp.btslot
}

func (rp *RunicPowerBar) SpendDeathRune(sim *Simulation, metrics *ResourceMetrics) int8 {
	if rp.runeStates&checkDeath == 0 {
		panic("Trying to spend death runes that don't exist!")
	}

	slot := rp.ReadyDeathRune()
	if rp.btslot != slot {
		// disable revert at
		rp.runeMeta[slot].revertAt = NeverExpires
		// clear death bit to revert.
		rp.runeStates = ^isDeaths[slot] & rp.runeStates
	}

	// mark spent bit to spend
	rp.runeStates |= isSpents[slot]
	rp.runeMeta[slot].lastSpendTime = sim.CurrentTime

	rp.SpendRuneMetrics(sim, metrics, 1)
	rp.LaunchRuneRegen(sim, slot)
	return slot
}

type RuneConvertType int8

const (
	RuneConvertTypeNone RuneConvertType = 1 << iota
	RuneConvertTypeBlood
	RuneConvertTypeFrost
	RuneConvertTypeUnholy
)

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
	baseCost := float64(NewRuneCost(uint8(options.RunicPowerCost), uint8(options.BloodRuneCost), uint8(options.FrostRuneCost), uint8(options.UnholyRuneCost), 0))
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
	//rp := &spell.Unit.RunicPowerBar
	spell.CurCast.Cost *= spell.CostMultiplier
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
		cost := RuneCost(spell.CurCast.Cost)
		spell.Unit.SpendRuneCost(sim, spell, cost)
	}
	if rc.RunicPowerGain > 0 && spell.CurCast.Cost > 0 {
		spell.Unit.AddRunicPower(sim, rc.RunicPowerGain, spell.RunicPowerMetrics())
	}
}
func (rc *RuneCostImpl) SpendRefundableCost(sim *Simulation, spell *Spell, result *SpellResult) {
	cost := RuneCost(spell.CurCast.Cost) // cost was already optimized in RuneSpell.Cast
	if cost == 0 {
		return // it was free this time. we don't care
	}
	if result.Landed() {
		spell.Unit.SpendRuneCost(sim, spell, cost)
	}
}
func (spell *Spell) SpendRefundableCost(sim *Simulation, result *SpellResult) {
	spell.Cost.(*RuneCostImpl).SpendRefundableCost(sim, spell, result)
}
func (rc *RuneCostImpl) SpendRefundableCostAndConvertBloodRune(sim *Simulation, spell *Spell, result *SpellResult, convertChance float64) {
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

	for _, slot := range [2]int8{slot1, slot2} {
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
	spell.Cost.(*RuneCostImpl).SpendRefundableCostAndConvertBloodRune(sim, spell, result, convertChance)
}
func (rc *RuneCostImpl) SpendRefundableCostAndConvertFrostOrUnholyRune(sim *Simulation, spell *Spell, result *SpellResult, convertChance float64) {
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

	for _, slot := range [3]int8{slot1, slot2, slot3} {
		if slot == 2 || slot == 3 || slot == 4 || slot == 5 {
			spell.Unit.ConvertToDeath(sim, slot, NeverExpires)
		}
	}
}
func (spell *Spell) SpendRefundableCostAndConvertFrostOrUnholyRune(sim *Simulation, result *SpellResult, convertChance float64) {
	spell.Cost.(*RuneCostImpl).SpendRefundableCostAndConvertFrostOrUnholyRune(sim, spell, result, convertChance)
}
func (rc *RuneCostImpl) IssueRefund(sim *Simulation, spell *Spell) {
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
