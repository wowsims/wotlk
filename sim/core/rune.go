package core

type OnRuneGain func(sim *Simulation)

type runeBar struct {
	unit *Unit

	name string

	maxRunes     float64
	currentRunes float64

	onRuneGain OnRuneGain
}

func (unit *Unit) EnableRuneBars(onBloodRuneGain OnRuneGain,
	onFrostRuneGain OnRuneGain,
	onUnholyRuneGain OnRuneGain,
	onDeathRuneGain OnRuneGain) {

	unit.runicPowerBar.bloodRunesBar = runeBar{
		unit: unit,

		name: "Blood",

		maxRunes:     2,
		currentRunes: 2,

		onRuneGain: onBloodRuneGain,
	}

	unit.runicPowerBar.frostRunesBar = runeBar{
		unit: unit,

		name: "Frost",

		maxRunes:     2,
		currentRunes: 2,

		onRuneGain: onFrostRuneGain,
	}

	unit.runicPowerBar.unholyRunesBar = runeBar{
		unit: unit,

		name: "Unholy",

		maxRunes:     2,
		currentRunes: 2,

		onRuneGain: onUnholyRuneGain,
	}

	unit.runicPowerBar.deathRunesBar = runeBar{
		unit: unit,

		name: "Death",

		maxRunes:     3,
		currentRunes: 0,

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

func (rb *runeBar) SpendRune(sim *Simulation, metrics *ResourceMetrics) {
	if rb.currentRunes <= 0 {
		panic("Trying to spend negative runic power!")
	}

	newRunes := rb.currentRunes - 1
	metrics.AddEvent(-1, -1)

	if sim.Log != nil {
		rb.unit.Log(sim, "Spent one %s Rune from %s (%d --> %d).", rb.name, metrics.ActionID, int32(rb.currentRunes), int32(newRunes))
	}

	rb.currentRunes = newRunes
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
