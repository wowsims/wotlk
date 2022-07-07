package core

type OnRunicPowerGain func(sim *Simulation)

type RunicPower struct {
	unit *Unit

	maxRunicPower     float64
	currentRunicPower float64

	onRunicPowerGain OnRunicPowerGain
}

func (unit *Unit) EnableRunicPowerBar(maxRunicPower float64,
	onRunicPowerGain OnRunicPowerGain) {
	unit.runeSystem = runeSystem{
		unit: unit,

		maxRunicPower:     maxRunicPower,
		currentRunicPower: maxRunicPower,

		onRunicPowerGain: onRunicPowerGain,
	}
}

func (unit *Unit) HasRunicPower() bool {
	return unit.runicPower.unit != nil
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
