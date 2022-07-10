package core

type OnRunicPowerGain func(sim *Simulation)

type runicPowerBar struct {
	unit *Unit

	maxRunicPower     float64
	currentRunicPower float64

	bloodRunesBar  runeBar
	frostRunesBar  runeBar
	unholyRunesBar runeBar
	deathRunesBar  runeBar

	onRunicPowerGain OnRunicPowerGain
}

func (unit *Unit) EnableRunicPowerBar(maxRunicPower float64,
	onRunicPowerGain OnRunicPowerGain) {
	unit.runicPowerBar = runicPowerBar{
		unit: unit,

		maxRunicPower:     maxRunicPower,
		currentRunicPower: maxRunicPower,

		onRunicPowerGain: onRunicPowerGain,
	}
}

func (unit *Unit) HasRunicPower() bool {
	return unit.runicPowerBar.unit != nil
}

func (rp *runicPowerBar) CurrentRunicPower() float64 {
	return rp.currentRunicPower
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

func (rp *runicPowerBar) CurrentBloodRunes() float64 {
	rb := &rp.bloodRunesBar
	return rb.CurrentRunes()
}

func (rp *runicPowerBar) CurrentFrostRunes() float64 {
	rb := &rp.frostRunesBar
	return rb.CurrentRunes()
}

func (rp *runicPowerBar) CurrentUnholyRunes() float64 {
	rb := &rp.unholyRunesBar
	return rb.CurrentRunes()
}

func (rp *runicPowerBar) CurrentDeathRunes() float64 {
	rb := &rp.deathRunesBar
	return rb.CurrentRunes()
}

func (rp *runicPowerBar) SpendBloodRune(sim *Simulation, metrics *ResourceMetrics) {
	rb := &rp.bloodRunesBar
	rb.SpendRune(sim, metrics)
}

func (rp *runicPowerBar) SpendFrostRune(sim *Simulation, metrics *ResourceMetrics) {
	rb := &rp.frostRunesBar
	rb.SpendRune(sim, metrics)
}

func (rp *runicPowerBar) SpendUnholyRune(sim *Simulation, metrics *ResourceMetrics) {
	rb := &rp.unholyRunesBar
	rb.SpendRune(sim, metrics)
}

func (rp *runicPowerBar) SpendDeathRune(sim *Simulation, metrics *ResourceMetrics) {
	rb := &rp.deathRunesBar
	rb.SpendRune(sim, metrics)
}
