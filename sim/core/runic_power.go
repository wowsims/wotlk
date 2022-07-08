package core

type OnRunicPowerGain func(sim *Simulation)

type RunicPowerBar struct {
	unit *Unit

	maxRunicPower     float64
	currentRunicPower float64

	onRunicPowerGain OnRunicPowerGain
}

func (unit *Unit) EnableRunicPowerBar(maxRunicPower float64,
	onRunicPowerGain OnRunicPowerGain) {
	unit.runicPowerBar = RunicPowerBar{
		unit: unit,

		maxRunicPower:     maxRunicPower,
		currentRunicPower: maxRunicPower,

		onRunicPowerGain: onRunicPowerGain,
	}
}

func (unit *Unit) HasRunicPower() bool {
	return unit.runicPowerBar.unit != nil
}

func (rp *RunicPowerBar) CurrentRunicPower() float64 {
	return rp.currentRunicPower
}

func (rp *RunicPowerBar) addRunicPowerInterval(sim *Simulation, amount float64, metrics *ResourceMetrics) {
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

func (rp *RunicPowerBar) AddRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	rp.addRunicPowerInterval(sim, amount, metrics)
	rp.onRunicPowerGain(sim)
}

func (rp *RunicPowerBar) SpendRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
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
