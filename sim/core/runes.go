package core

type OnRunicPowerGain func(sim *Simulation)
type OnBloodRuneGain func(sim *Simulation)
type OnFrostRuneGain func(sim *Simulation)
type OnUnholyRuneGain func(sim *Simulation)
type OnDeathRuneGain func(sim *Simulation)

type runeSystem struct {
	unit *Unit

	maxRunes    int32
	bloodRunes  int32
	frostRunes  int32
	unholyRunes int32

	maxDeathRunes int32
	deathRunes    int32

	maxRunicPower     float64
	currentRunicPower float64

	onRunicPowerGain OnRunicPowerGain
	onBloodRuneGain  OnBloodRuneGain
	onFrostRuneGain  OnFrostRuneGain
	onUnholyRuneGain OnUnholyRuneGain
	onDeathRuneGain  OnDeathRuneGain
}

func (unit *Unit) EnableRuneSystem(maxRunicPower float64,
	onRunicPowerGain OnRunicPowerGain,
	onBloodRuneGain OnBloodRuneGain,
	onFrostRuneGain OnFrostRuneGain,
	onUnholyRuneGain OnUnholyRuneGain,
	onDeathRuneGain OnDeathRuneGain) {
	unit.runeSystem = runeSystem{
		unit: unit,

		maxRunes:    2,
		bloodRunes:  2,
		frostRunes:  2,
		unholyRunes: 2,

		maxDeathRunes: 3,
		deathRunes:    0,

		maxRunicPower:     maxRunicPower,
		currentRunicPower: maxRunicPower,

		onRunicPowerGain: onRunicPowerGain,
		onBloodRuneGain:  onBloodRuneGain,
		onFrostRuneGain:  onFrostRuneGain,
		onUnholyRuneGain: onUnholyRuneGain,
		onDeathRuneGain:  onDeathRuneGain,
	}
}

func (unit *Unit) HasRuneSystem() bool {
	return unit.runeSystem.unit != nil
}

func (rs *runeSystem) CurrentRunicPower() float64 {
	return rs.currentRunicPower
}

func (rs *runeSystem) CurrentBloodRunes() int32 {
	return rs.bloodRunes
}

func (rs *runeSystem) CurrentFrostRunes() int32 {
	return rs.frostRunes
}

func (rs *runeSystem) CurrentUnholyRunes() int32 {
	return rs.unholyRunes
}

func (rs *runeSystem) CurrentDeathRunes() int32 {
	return rs.deathRunes
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

func (rs *runeSystem) SpendBloodRune(sim *Simulation, metrics *ResourceMetrics) {
	if rs.bloodRunes <= 0 {
		panic("Trying to spend blood runes but we are out!")
	}

	newBloodRunes := rs.bloodRunes - 1
	metrics.AddEvent(-1, -1)

	if sim.Log != nil {
		rs.unit.Log(sim, "Spent 1 blood rune from %s (%3d --> %3d).", metrics.ActionID, rs.bloodRunes, newBloodRunes)
	}

	rs.bloodRunes = newBloodRunes
}

func (rs *runeSystem) SpendFrostRune(sim *Simulation, metrics *ResourceMetrics) {
	if rs.frostRunes <= 0 {
		panic("Trying to spend frost runes but we are out!")
	}

	newFrostRunes := rs.frostRunes - 1
	metrics.AddEvent(-1, -1)

	if sim.Log != nil {
		rs.unit.Log(sim, "Spent 1 frost rune from %s (%3d --> %3d).", metrics.ActionID, rs.frostRunes, newFrostRunes)
	}

	rs.frostRunes = newFrostRunes
}

func (rs *runeSystem) SpendUnholyRune(sim *Simulation, metrics *ResourceMetrics) {
	if rs.unholyRunes <= 0 {
		panic("Trying to spend unholy runes but we are out!")
	}

	newUnholyRunes := rs.unholyRunes - 1
	metrics.AddEvent(-1, -1)

	if sim.Log != nil {
		rs.unit.Log(sim, "Spent 1 unholy rune from %s (%3d --> %3d).", metrics.ActionID, rs.unholyRunes, newUnholyRunes)
	}

	rs.unholyRunes = newUnholyRunes
}
