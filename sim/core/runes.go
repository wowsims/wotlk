package core

type OnRunicPowerGain func(sim *Simulation)

type rune struct {
	exists        bool
	is_death_rune bool
}

type runeSystem struct {
	unit *Unit

	bloodRunes  int32
	frostRunes  int32
	unholyRunes int32
	deathRunes  int32

	maxRunicPower     float64
	currentRunicPower float64

	onRunicPowerGain OnRunicPowerGain
}

func (unit *Unit) EnableRuneSystem(maxRunicPower float64, onRunicPowerGain OnRunicPowerGain) {
	unit.runeSystem = runeSystem{
		unit: unit,

		bloodRunes:  2,
		frostRunes:  2,
		unholyRunes: 2,

		maxRunicPower:     100,
		currentRunicPower: 0,

		onRunicPowerGain: onRunicPowerGain,
	}
}

func (unit *Unit) HasRuneSystem() bool {
	return unit.runeSystem.unit != nil
}

func (rs *runeSystem) CurrentRunicPower() float64 {
	return rs.currentRunicPower
}

func (rs *runeSystem) addRunicPowerInterval(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to add negative runic power!")
	}

	newRunicPower := MinFloat(rs.currentRunicPower+amount, rs.maxRunicPower)
	metrics.AddEvent(amount, newRunicPower-rs.currentRunicPower)

	if sim.Log != nil {
		rs.unit.Log(sim, "Gained %0.3f runic power from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, rs.currentRunicPower, newRunicPower)
	}

	rs.currentRunicPower = newRunicPower
}

func (rs *runeSystem) AddRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	rs.addRunicPowerInterval(sim, amount, metrics)
	rs.onRunicPowerGain(sim)
}

func (rs *runeSystem) SpendRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to spend negative runic power!")
	}

	newRunicPower := rs.currentRunicPower - amount
	metrics.AddEvent(-amount, -amount)

	if sim.Log != nil {
		rs.unit.Log(sim, "Spent %0.3f runic power from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, rs.currentRunicPower, newRunicPower)
	}

	rs.currentRunicPower = newRunicPower
}
