package core

import (
	"time"
)

type OnRuneGain func(sim *Simulation)

const (
	RuneType_Blood int32 = iota
	RuneType_Frost
	RuneType_Unholy
	RuneType_Death
)

const (
	RuneTypeBloodName  string = "Blood"
	RuneTypeFrostName  string = "Frost"
	RuneTypeUnholyName string = "Unholy"
	RuneTypeDeathName  string = "Death"
	RuneTypeUndefName  string = "Undefined"
)

func (rb *runeBar) RuneTypeName() string {
	switch rb.runeType {
	case RuneType_Blood:
		return RuneTypeBloodName
	case RuneType_Frost:
		return RuneTypeFrostName
	case RuneType_Unholy:
		return RuneTypeUnholyName
	case RuneType_Death:
		return RuneTypeDeathName
	}
	// Should never get this!
	return RuneTypeUndefName
}

type AvailableRune struct {
	slot      int32
	available bool
}

type runeBar struct {
	unit *Unit

	runeType int32

	maxRunes     float64
	currentRunes float64

	cooldowns [3]Cooldown

	onRuneGain OnRuneGain
}

func (unit *Unit) EnableRuneBars(onBloodRuneGain OnRuneGain,
	onFrostRuneGain OnRuneGain,
	onUnholyRuneGain OnRuneGain,
	onDeathRuneGain OnRuneGain) {

	unit.runicPowerBar.bloodRunesBar = runeBar{
		unit: unit,

		runeType: RuneType_Blood,

		maxRunes:     2,
		currentRunes: 2,

		cooldowns: [3]Cooldown{
			Cooldown{unit.NewTimer(), 10 * time.Second},
			Cooldown{unit.NewTimer(), 10 * time.Second},
			Cooldown{nil, 0},
		},

		onRuneGain: onBloodRuneGain,
	}

	unit.runicPowerBar.frostRunesBar = runeBar{
		unit: unit,

		runeType: RuneType_Frost,

		maxRunes:     2,
		currentRunes: 2,

		cooldowns: [3]Cooldown{
			Cooldown{unit.NewTimer(), 10 * time.Second},
			Cooldown{unit.NewTimer(), 10 * time.Second},
			Cooldown{nil, 0},
		},

		onRuneGain: onFrostRuneGain,
	}

	unit.runicPowerBar.unholyRunesBar = runeBar{
		unit: unit,

		runeType: RuneType_Unholy,

		maxRunes:     2,
		currentRunes: 2,

		cooldowns: [3]Cooldown{
			Cooldown{unit.NewTimer(), 10 * time.Second},
			Cooldown{unit.NewTimer(), 10 * time.Second},
			Cooldown{nil, 0},
		},

		onRuneGain: onUnholyRuneGain,
	}

	unit.runicPowerBar.deathRunesBar = runeBar{
		unit: unit,

		runeType: RuneType_Death,

		maxRunes:     3,
		currentRunes: 0,

		cooldowns: [3]Cooldown{
			Cooldown{unit.NewTimer(), 10 * time.Second},
			Cooldown{unit.NewTimer(), 10 * time.Second},
			Cooldown{unit.NewTimer(), 10 * time.Second},
		},

		onRuneGain: onDeathRuneGain,
	}
}

func (unit *Unit) HasBloodRuneBar() bool {
	return unit.runicPowerBar.bloodRunesBar.unit != nil
}
func (unit *Unit) HasFrostRuneBar() bool {
	return unit.runicPowerBar.frostRunesBar.unit != nil
}
func (unit *Unit) HasUnholyRuneBar() bool {
	return unit.runicPowerBar.unholyRunesBar.unit != nil
}
func (unit *Unit) HasDeathRuneBar() bool {
	return unit.runicPowerBar.deathRunesBar.unit != nil
}

func (rb *runeBar) CurrentRunes() float64 {
	return rb.currentRunes
}

func (rb *runeBar) AnyAvailableRune(sim *Simulation) bool {
	available := false
	if rb.cooldowns[1].IsReady(sim) {
		available = true
	} else if rb.cooldowns[0].IsReady(sim) {
		available = true
	}

	return available
}

func (rb *runeBar) GetAvailableRune(sim *Simulation) AvailableRune {
	available := false
	slot := -1
	if rb.cooldowns[1].IsReady(sim) {
		available = true
		slot = 1
	} else if rb.cooldowns[0].IsReady(sim) {
		available = true
		slot = 0
	}

	return AvailableRune{int32(slot), available}
}

func (rb *runeBar) SpendRune(sim *Simulation, metrics *ResourceMetrics) {
	availableRune := rb.GetAvailableRune(sim)

	if availableRune.available {
		newRunes := rb.currentRunes - 1
		metrics.AddEvent(-1, -1)

		if sim.Log != nil {
			rb.unit.Log(sim, "Spent %s Rune(%d) from %s (%d --> %d).", rb.RuneTypeName(), availableRune.slot, metrics.ActionID, int32(rb.currentRunes), int32(newRunes))
		}

		rb.currentRunes = newRunes
		rb.cooldowns[availableRune.slot].Use(sim)
	}
}

//func (rp *runicPowerBar) addRunicPowerInterval(sim *Simulation, amount float64, metrics *ResourceMetrics) {
//	if amount < 0 {
//		panic("Trying to add negative runic power!")
//	}
//
//	newRunicPower := MinFloat(rp.currentRunicPower+amount, rp.maxRunicPower)
//	metrics.AddEvent(amount, newRunicPower-rp.currentRunicPower)
//
//	if sim.Log != nil {
//		rp.unit.Log(sim, "Gained %0.3f runic power from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, rp.currentRunicPower, newRunicPower)
//	}
//
//	rp.currentRunicPower = newRunicPower
//}
//
//func (rp *runicPowerBar) AddRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
//	rp.addRunicPowerInterval(sim, amount, metrics)
//	rp.onRunicPowerGain(sim)
//}
//
