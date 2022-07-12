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

type Rune struct {
	unit *Unit

	cd Cooldown
}

type DeathRune struct {
	unit *Unit

	cd Cooldown
}

type runicPowerBar struct {
	unit *Unit

	maxRunicPower     float64
	currentRunicPower float64

	bloodRunes  [2]Rune
	frostRunes  [2]Rune
	unholyRunes [2]Rune
	deathRunes  [6]DeathRune

	runeGainTrackers [4]int32

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

		bloodRunes: [2]Rune{
			Rune{
				unit: unit,
				cd:   Cooldown{unit.NewTimer(), 10.0 * time.Second},
			},
			Rune{
				unit: unit,
				cd:   Cooldown{unit.NewTimer(), 10.0 * time.Second},
			},
		},

		frostRunes: [2]Rune{
			Rune{
				unit: unit,
				cd:   Cooldown{unit.NewTimer(), 10.0 * time.Second},
			},
			Rune{
				unit: unit,
				cd:   Cooldown{unit.NewTimer(), 10.0 * time.Second},
			},
		},

		unholyRunes: [2]Rune{
			Rune{
				unit: unit,
				cd:   Cooldown{unit.NewTimer(), 10.0 * time.Second},
			},
			Rune{
				unit: unit,
				cd:   Cooldown{unit.NewTimer(), 10.0 * time.Second},
			},
		},

		deathRunes: [6]DeathRune{
			DeathRune{
				unit: unit,
				cd:   Cooldown{unit.NewTimer(), 20.0 * time.Second},
			},
			DeathRune{
				unit: unit,
				cd:   Cooldown{unit.NewTimer(), 20.0 * time.Second},
			},
			DeathRune{
				unit: unit,
				cd:   Cooldown{unit.NewTimer(), 20.0 * time.Second},
			},
			DeathRune{
				unit: unit,
				cd:   Cooldown{unit.NewTimer(), 20.0 * time.Second},
			},
			DeathRune{
				unit: unit,
				cd:   Cooldown{unit.NewTimer(), 20.0 * time.Second},
			},
			DeathRune{
				unit: unit,
				cd:   Cooldown{unit.NewTimer(), 20.0 * time.Second},
			},
		},

		runeGainTrackers: [4]int32{2, 2, 2, 0},

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

	unit.runicPowerBar.deathRunes[0].cd.Set(1<<63 - 1)
	unit.runicPowerBar.deathRunes[1].cd.Set(1<<63 - 1)
	unit.runicPowerBar.deathRunes[2].cd.Set(1<<63 - 1)
	unit.runicPowerBar.deathRunes[3].cd.Set(1<<63 - 1)
	unit.runicPowerBar.deathRunes[4].cd.Set(1<<63 - 1)
	unit.runicPowerBar.deathRunes[5].cd.Set(1<<63 - 1)
}

func (unit *Unit) HasRunicPower() bool {
	return unit.runicPowerBar.unit != nil
}

func (rp *runicPowerBar) CurrentRunicPower() float64 {
	return rp.currentRunicPower
}

func (r *Rune) IsReady(sim *Simulation) bool {
	return r.cd.IsReady(sim)
}

func (dr *DeathRune) IsReady(sim *Simulation) bool {
	return dr.cd.IsReady(sim)
}

func (rp *runicPowerBar) CurrentBloodRunes(sim *Simulation) int32 {
	total := int32(0)
	if rp.bloodRunes[1].IsReady(sim) {
		total += 1
	}
	if rp.bloodRunes[0].IsReady(sim) {
		total += 1
	}
	return total
}

func (rp *runicPowerBar) CurrentFrostRunes(sim *Simulation) int32 {
	total := int32(0)
	if rp.frostRunes[1].IsReady(sim) {
		total += 1
	}
	if rp.frostRunes[0].IsReady(sim) {
		total += 1
	}
	return total
}

func (rp *runicPowerBar) CurrentUnholyRunes(sim *Simulation) int32 {
	total := int32(0)
	if rp.unholyRunes[1].IsReady(sim) {
		total += 1
	}
	if rp.unholyRunes[0].IsReady(sim) {
		total += 1
	}
	return total
}

func (rp *runicPowerBar) CurrentDeathRunes(sim *Simulation) int32 {
	total := int32(0)
	if rp.deathRunes[5].IsReady(sim) {
		total += 1
	}
	if rp.deathRunes[4].IsReady(sim) {
		total += 1
	}
	if rp.deathRunes[3].IsReady(sim) {
		total += 1
	}
	if rp.deathRunes[2].IsReady(sim) {
		total += 1
	}
	if rp.deathRunes[1].IsReady(sim) {
		total += 1
	}
	if rp.deathRunes[0].IsReady(sim) {
		total += 1
	}
	return total
}

func (rp *runicPowerBar) GainRuneMetrics(sim *Simulation, metrics *ResourceMetrics, runeName string, newRunes int32, prevRunes int32) {
	metrics.AddEvent(1, float64(newRunes-prevRunes))

	if sim.Log != nil {
		rp.unit.Log(sim, "Gained %s Rune (%d --> %d).", runeName, prevRunes, newRunes)
	}
}

func (rp *runicPowerBar) CheckRuneGainTrackers(sim *Simulation) {
	newRunes := rp.CurrentBloodRunes(sim)
	prevRunes := rp.runeGainTrackers[0]
	if newRunes > prevRunes {
		rp.onBloodRuneGain(sim)
		rp.GainRuneMetrics(sim, rp.bloodRuneGainMetrics, "Blood", newRunes, prevRunes)
	}

	newRunes = rp.CurrentFrostRunes(sim)
	prevRunes = rp.runeGainTrackers[1]
	if newRunes > prevRunes {
		rp.onFrostRuneGain(sim)
		rp.GainRuneMetrics(sim, rp.frostRuneGainMetrics, "Frost", newRunes, prevRunes)
	}

	newRunes = rp.CurrentUnholyRunes(sim)
	prevRunes = rp.runeGainTrackers[2]
	if newRunes > prevRunes {
		rp.onUnholyRuneGain(sim)
		rp.GainRuneMetrics(sim, rp.unholyRuneGainMetrics, "Unholy", newRunes, prevRunes)
	}

	newRunes = rp.CurrentDeathRunes(sim)
	prevRunes = rp.runeGainTrackers[3]
	if newRunes > prevRunes {
		rp.onDeathRuneGain(sim)
		rp.GainRuneMetrics(sim, rp.deathRuneGainMetrics, "Death", newRunes, prevRunes)
	}
}

func (rp *runicPowerBar) UpdateRuneGainTrackers(sim *Simulation) {
	rp.runeGainTrackers[0] = rp.CurrentBloodRunes(sim)
	rp.runeGainTrackers[1] = rp.CurrentFrostRunes(sim)
	rp.runeGainTrackers[2] = rp.CurrentUnholyRunes(sim)
	rp.runeGainTrackers[3] = rp.CurrentDeathRunes(sim)
}

func (rp *runicPowerBar) BloodRuneReadyAt(sim *Simulation) time.Duration {
	return MinDuration(rp.bloodRunes[0].cd.ReadyAt(), rp.bloodRunes[1].cd.ReadyAt())
}

func (rp *runicPowerBar) FrostRuneReadyAt(sim *Simulation) time.Duration {
	return MinDuration(rp.frostRunes[0].cd.ReadyAt(), rp.frostRunes[1].cd.ReadyAt())
}

func (rp *runicPowerBar) UnholyRuneReadyAt(sim *Simulation) time.Duration {
	return MinDuration(rp.unholyRunes[0].cd.ReadyAt(), rp.unholyRunes[1].cd.ReadyAt())
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

func (rp *runicPowerBar) CastCostPossible(sim *Simulation, runicPowerAmount float64, bloodAmount int32, frostAmount int32, unholyAmount int32, deathAmount int32) bool {
	return (rp.currentRunicPower > runicPowerAmount) &&
		(rp.CurrentBloodRunes(sim) >= bloodAmount) &&
		(rp.CurrentFrostRunes(sim) >= frostAmount) &&
		(rp.CurrentUnholyRunes(sim) >= unholyAmount) &&
		(rp.CurrentDeathRunes(sim) >= deathAmount)
}

// func (rb *runeBar) SpendRune(sim *Simulation, metrics *ResourceMetrics) {
// 	availableRune := rb.GetAvailableRune(sim)

// 	if availableRune.available {
// 		newRunes := rb.currentRunes - 1
// 		metrics.AddEvent(-1, -1)

// 		if sim.Log != nil {
// 			rb.unit.Log(sim, "Spent %s Rune(%d) from %s (%d --> %d).", rb.RuneTypeName(), availableRune.slot, metrics.ActionID, int32(rb.currentRunes), int32(newRunes))
// 		}

// 		rb.currentRunes = newRunes
// 		rb.cooldowns[availableRune.slot].Use(sim)
// 	}
// }

func (r *Rune) Spend(sim *Simulation) {
	r.cd.Use(sim)
}

func (dr *DeathRune) Spend(sim *Simulation) {
	dr.cd.Use(sim)
}

func (rp *runicPowerBar) SpendRuneMetrics(sim *Simulation, metrics *ResourceMetrics, runeName string, slot int32, currentRunes int32, newRunes int32) {
	metrics.AddEvent(-1, -1)
	if sim.Log != nil {
		rp.unit.Log(sim, "Spent %s Rune(slot %d) from %s (%d --> %d).", runeName, slot, metrics.ActionID, currentRunes, newRunes)
	}
}

func (rp *runicPowerBar) SpendBloodRune(sim *Simulation, metrics *ResourceMetrics) {
	currentBloodRunes := rp.CurrentBloodRunes(sim)

	slot := int32(-1)
	if rp.bloodRunes[0].IsReady(sim) {
		slot = 0
	} else if rp.bloodRunes[1].IsReady(sim) {
		slot = 1
	} else {
		panic("Trying to spend blood rune but we have none!")
	}

	rp.bloodRunes[slot].Spend(sim)
	rp.SpendRuneMetrics(sim, metrics, "Blood", slot, currentBloodRunes, rp.CurrentBloodRunes(sim))
}

func (rp *runicPowerBar) SpendFrostRune(sim *Simulation, metrics *ResourceMetrics) {
	currentFrostRunes := rp.CurrentFrostRunes(sim)

	slot := int32(-1)
	if rp.frostRunes[0].IsReady(sim) {
		slot = 0
	} else if rp.frostRunes[1].IsReady(sim) {
		slot = 1
	} else {
		panic("Trying to spend frost rune but we have none!")
	}

	rp.frostRunes[slot].Spend(sim)
	rp.SpendRuneMetrics(sim, metrics, "Frost", slot, currentFrostRunes, rp.CurrentFrostRunes(sim))
}

func (rp *runicPowerBar) SpendUnholyRune(sim *Simulation, metrics *ResourceMetrics) {
	currentUnholyRunes := rp.CurrentUnholyRunes(sim)

	slot := int32(-1)
	if rp.unholyRunes[0].IsReady(sim) {
		slot = 0
	} else if rp.unholyRunes[1].IsReady(sim) {
		slot = 1
	} else {
		panic("Trying to spend unholy rune but we have none!")
	}

	rp.unholyRunes[slot].Spend(sim)
	rp.SpendRuneMetrics(sim, metrics, "Unholy", slot, currentUnholyRunes, rp.CurrentUnholyRunes(sim))
}

func (rp *runicPowerBar) SpendDeathRune(sim *Simulation, metrics *ResourceMetrics) {
	currentDeathRunes := rp.CurrentDeathRunes(sim)

	slot := int32(-1)
	if rp.deathRunes[0].IsReady(sim) {
		slot = 0
	} else if rp.deathRunes[1].IsReady(sim) {
		slot = 1
	} else if rp.deathRunes[2].IsReady(sim) {
		slot = 2
	} else if rp.deathRunes[3].IsReady(sim) {
		slot = 3
	} else if rp.deathRunes[4].IsReady(sim) {
		slot = 4
	} else if rp.deathRunes[5].IsReady(sim) {
		slot = 5
	} else {
		panic("Trying to spend death rune but we have none!")
	}

	rp.deathRunes[slot].Spend(sim)
	rp.SpendRuneMetrics(sim, metrics, "Death", slot, currentDeathRunes, rp.CurrentDeathRunes(sim))
}
