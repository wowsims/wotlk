package common

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

const manaBuffer = 500.0
const manaTrackingWindowSeconds = 60
const manaTrackingWindow = time.Second * manaTrackingWindowSeconds

// 3 * (# of seconds) should be plenty of slots
const manaSnapshotsBufferSize = manaTrackingWindowSeconds

// Tracks how fast mana is being spent. This is used by some specs to decide
// whether to use more mana-efficient or higher-dps spells.
type ManaSpendingRateTracker struct {
	// Circular array buffer for recent mana snapshots, within a time window
	manaSnapshots      [manaSnapshotsBufferSize]manaSnapshot
	numSnapshots       int32
	firstSnapshotIndex int32

	manaSpentDuringWindow  float64
	manaGainedDuringWindow float64

	previousManaSpent  float64
	previousManaGained float64
	previousCastSpeed  float64
}

type manaSnapshot struct {
	time       time.Duration // time this snapshot was taken
	manaSpent  float64       // total amount of mana spent up to this time
	manaGained float64       // total amount of mana gained, minus bonus mana (pots/runes/innervates).

	manaSpentDelta  float64
	manaGainedDelta float64
}

func NewManaSpendingRateTracker() ManaSpendingRateTracker {
	return ManaSpendingRateTracker{}
}

// This needs to be called on sim reset.
func (tracker *ManaSpendingRateTracker) Reset() {
	tracker.manaSnapshots = [manaSnapshotsBufferSize]manaSnapshot{}
	tracker.firstSnapshotIndex = 0
	tracker.numSnapshots = 0
	tracker.manaSpentDuringWindow = 0
	tracker.manaGainedDuringWindow = 0
	tracker.previousManaSpent = 0
	tracker.previousManaGained = 0
	tracker.previousCastSpeed = 1
}

func (tracker *ManaSpendingRateTracker) getOldestSnapshot() manaSnapshot {
	return tracker.manaSnapshots[tracker.firstSnapshotIndex]
}

func (tracker *ManaSpendingRateTracker) purgeExpiredSnapshots(sim *core.Simulation) {
	expirationCutoff := sim.CurrentTime - manaTrackingWindow

	curIndex := tracker.firstSnapshotIndex
	for tracker.numSnapshots > 0 && tracker.manaSnapshots[curIndex].time < expirationCutoff {
		tracker.manaSpentDuringWindow -= tracker.manaSnapshots[curIndex].manaSpentDelta
		tracker.manaGainedDuringWindow -= tracker.manaSnapshots[curIndex].manaGainedDelta
		curIndex = (curIndex + 1) % manaSnapshotsBufferSize
		tracker.numSnapshots--
	}
	tracker.firstSnapshotIndex = curIndex
}

// This needs to be called at regular intervals to update the tracker's data.
func (tracker *ManaSpendingRateTracker) Update(sim *core.Simulation, character *core.Character) {
	if tracker.numSnapshots >= manaSnapshotsBufferSize {
		panic("Mana tracker snapshot buffer is full")
	}

	// Scale down mana spent/gained so we don't get bad estimates from lust/drums/etc.
	manaDeltaCoefficient := character.InitialCastSpeed() / tracker.previousCastSpeed
	manaSpent := character.Metrics.ManaSpent
	manaGained := character.Metrics.ManaGained - character.Metrics.BonusManaGained

	snapshot := manaSnapshot{
		time:            sim.CurrentTime,
		manaSpent:       manaSpent,
		manaGained:      manaGained,
		manaSpentDelta:  (manaSpent - tracker.previousManaSpent) * manaDeltaCoefficient,
		manaGainedDelta: (manaGained - tracker.previousManaGained) * manaDeltaCoefficient / character.PseudoStats.SpiritRegenMultiplier,
	}
	//if sim.Log != nil {
	//	character.Log(sim, "Init speed: %0.02f, prev cast speed: %0.02f, Mana gained: %0.02f, Mana gained delta: %0.02f", character.InitialCastSpeed(), tracker.previousCastSpeed, snapshot.manaGained, snapshot.manaGainedDelta)
	//}

	nextIndex := (tracker.firstSnapshotIndex + tracker.numSnapshots) % manaSnapshotsBufferSize
	tracker.previousCastSpeed = 1 / character.CastSpeed
	tracker.previousManaSpent = snapshot.manaSpent
	tracker.previousManaGained = snapshot.manaGained
	tracker.manaSpentDuringWindow += snapshot.manaSpentDelta
	tracker.manaGainedDuringWindow += snapshot.manaGainedDelta
	tracker.manaSnapshots[nextIndex] = snapshot
	tracker.numSnapshots++

	tracker.purgeExpiredSnapshots(sim)
}

func (tracker *ManaSpendingRateTracker) ManaSpentPerSecond(sim *core.Simulation, character *core.Character) float64 {
	tracker.purgeExpiredSnapshots(sim)
	oldestSnapshot := tracker.getOldestSnapshot()

	manaSpent := tracker.manaSpentDuringWindow - tracker.manaGainedDuringWindow
	timeDelta := sim.CurrentTime - oldestSnapshot.time
	if timeDelta == 0 {
		return 0
	}

	return manaSpent / timeDelta.Seconds()
}

// The amount of mana we will need to spend over the remaining sim duration
// at the current rate of mana spending.
func (tracker *ManaSpendingRateTracker) ProjectedManaCost(sim *core.Simulation, character *core.Character) float64 {
	manaSpentPerSecond := tracker.ManaSpentPerSecond(sim, character)

	projectedManaCost := manaSpentPerSecond * sim.GetRemainingDuration().Seconds()

	//if sim.Log != nil {
	//	remainingManaPool := character.CurrentMana() + character.ExpectedBonusMana - manaBuffer
	//	character.Log(sim, "Mana spent: %0.02f, Mana gained: %0.02f, BonusManaGained: %0.02f, Projected: %0.02f, total: %0.02f (%0.02f + %0.02f)", character.Metrics.ManaSpent, character.Metrics.ManaGained, character.Metrics.BonusManaGained, projectedManaCost, remainingManaPool, character.CurrentMana(), character.ExpectedBonusMana)
	//}

	return projectedManaCost
}

func (tracker *ManaSpendingRateTracker) ProjectedRemainingMana(sim *core.Simulation, character *core.Character) float64 {
	return character.CurrentMana() + character.ExpectedBonusMana - manaBuffer
}

func (tracker *ManaSpendingRateTracker) ProjectedManaSurplus(sim *core.Simulation, character *core.Character) bool {
	// If we've gone OOM at least once, stop using surplus rotations.
	// Spending time not casting while OOM will throw off the mana spend / gain rates so this is necessary.
	if character.Metrics.WentOOM {
		return false
	}

	return tracker.ProjectedManaCost(sim, character) < tracker.ProjectedRemainingMana(sim, character)
}
