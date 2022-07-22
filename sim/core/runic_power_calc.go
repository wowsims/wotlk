package core

import (
	"time"
)

type CalcRune struct {
	state                    RuneState
	kind                     RuneKind
	pas                      [2]float64
	lastRegenTime            float64
	generatedByReapingOrBoTN bool
}

type CalcRunicPowerBar struct {
	maxRunicPower     float64
	currentRunicPower float64

	bloodRunes  [2]CalcRune
	frostRunes  [2]CalcRune
	unholyRunes [2]CalcRune
}

func GetCalcRune(r *Rune) CalcRune {
	cr := CalcRune{state: r.state,
		kind:                     r.kind,
		lastRegenTime:            float64(r.lastRegenTime),
		generatedByReapingOrBoTN: r.generatedByReapingOrBoTN,
	}

	pa0 := -1.0
	if r.pas[0] != nil {
		pa0 = float64(r.pas[0].NextActionAt)
	}

	pa1 := -1.0
	if r.pas[1] != nil {
		pa1 = float64(r.pas[1].NextActionAt)
	}

	cr.pas = [2]float64{pa0, pa1}

	return cr
}

func (rp *runicPowerBar) GetCalcRunicPowerBar() CalcRunicPowerBar {
	return CalcRunicPowerBar{
		maxRunicPower:     rp.maxRunicPower,
		currentRunicPower: rp.currentRunicPower,
		bloodRunes:        [2]CalcRune{GetCalcRune(&rp.bloodRunes[0]), GetCalcRune(&rp.bloodRunes[1])},
		frostRunes:        [2]CalcRune{GetCalcRune(&rp.frostRunes[0]), GetCalcRune(&rp.frostRunes[1])},
		unholyRunes:       [2]CalcRune{GetCalcRune(&rp.unholyRunes[0]), GetCalcRune(&rp.unholyRunes[1])},
	}
}

func (rp *CalcRunicPowerBar) addRunicPowerInterval(sim *Simulation, amount float64) {
	if amount < 0 {
		panic("Trying to add negative runic power!")
	}

	rp.currentRunicPower = MinFloat(rp.currentRunicPower+amount, rp.maxRunicPower)
}

func (rp *CalcRunicPowerBar) AddRunicPower(sim *Simulation, amount float64) {
	rp.addRunicPowerInterval(sim, amount)
}

func (rp *CalcRunicPowerBar) SpendRunicPower(sim *Simulation, amount float64) {
	if amount < 0 {
		panic("Trying to spend negative runic power!")
	}

	rp.currentRunicPower = rp.currentRunicPower - amount

}

func (rp *CalcRunicPowerBar) CurrentRunesOfType(rb *[2]CalcRune, runeState RuneState) int32 {
	return TernaryInt32(rb[0].state == runeState, 1, 0) + TernaryInt32(rb[1].state == runeState, 1, 0)
}

func (rp *CalcRunicPowerBar) DeathRuneRegenAt(r *CalcRune) float64 {
	readyAt := float64(NeverExpires)

	if r.state == RuneState_DeathSpent {
		if r.pas[1] >= 0 {
			readyAt = TernaryFloat64(r.pas[0] < r.pas[1], r.pas[0], readyAt)
		} else {
			readyAt = r.pas[0]
		}
	}

	return readyAt
}

func (rp *CalcRunicPowerBar) SpentDeathRuneReadyAt(sim *Simulation) float64 {
	readyAt := float64(NeverExpires)

	readyAt = MinFloat(readyAt, rp.DeathRuneRegenAt(&rp.bloodRunes[0]))
	readyAt = MinFloat(readyAt, rp.DeathRuneRegenAt(&rp.bloodRunes[1]))
	readyAt = MinFloat(readyAt, rp.DeathRuneRegenAt(&rp.frostRunes[0]))
	readyAt = MinFloat(readyAt, rp.DeathRuneRegenAt(&rp.frostRunes[1]))
	readyAt = MinFloat(readyAt, rp.DeathRuneRegenAt(&rp.unholyRunes[0]))
	readyAt = MinFloat(readyAt, rp.DeathRuneRegenAt(&rp.unholyRunes[1]))

	return readyAt
}
func (rp *CalcRunicPowerBar) DeathRuneReadyAt(sim *Simulation) float64 {
	readyAt := float64(NeverExpires)

	if rp.bloodRunes[0].state == RuneState_Death || rp.bloodRunes[1].state == RuneState_Death ||
		rp.frostRunes[0].state == RuneState_Death || rp.frostRunes[1].state == RuneState_Death ||
		rp.unholyRunes[0].state == RuneState_Death || rp.unholyRunes[1].state == RuneState_Death {
		readyAt = float64(sim.CurrentTime)
		return readyAt
	}

	return rp.SpentDeathRuneReadyAt(sim)
}

func (rp *CalcRunicPowerBar) SpentRuneReadyAt(sim *Simulation, runes *[2]CalcRune) float64 {
	readyAt := rp.SpentDeathRuneReadyAt(sim)

	if runes[0].pas[0] >= 0 {
		readyAt = MinFloat(readyAt, runes[0].pas[0])
	}

	if runes[1].pas[0] >= 0 {
		readyAt = MinFloat(readyAt, runes[1].pas[0])
	}

	return readyAt
}

func (rp *CalcRunicPowerBar) RuneReadyAt(sim *Simulation, runes *[2]CalcRune) float64 {
	readyAt := rp.DeathRuneReadyAt(sim)

	if runes[0].state == RuneState_Normal || runes[0].state == RuneState_Death ||
		runes[1].state == RuneState_Normal || runes[1].state == RuneState_Death {
		readyAt = float64(sim.CurrentTime)
		return readyAt
	}

	return rp.SpentRuneReadyAt(sim, runes)
}

func (rp *CalcRunicPowerBar) SpentBloodRuneReadyAt(sim *Simulation) float64 {
	return rp.SpentRuneReadyAt(sim, &rp.bloodRunes)
}

func (rp *CalcRunicPowerBar) SpentFrostRuneReadyAt(sim *Simulation) float64 {
	return rp.SpentRuneReadyAt(sim, &rp.frostRunes)
}

func (rp *CalcRunicPowerBar) SpentUnholyRuneReadyAt(sim *Simulation) float64 {
	return rp.SpentRuneReadyAt(sim, &rp.unholyRunes)
}

func (rp *CalcRunicPowerBar) BloodRuneReadyAt(sim *Simulation) float64 {
	return rp.RuneReadyAt(sim, &rp.bloodRunes)
}

func (rp *CalcRunicPowerBar) FrostRuneReadyAt(sim *Simulation) float64 {
	return rp.RuneReadyAt(sim, &rp.frostRunes)
}

func (rp *CalcRunicPowerBar) UnholyRuneReadyAt(sim *Simulation) float64 {
	return rp.RuneReadyAt(sim, &rp.unholyRunes)
}

func (rp *CalcRunicPowerBar) AnySpentRuneReadyAt(sim *Simulation) float64 {
	return MinFloat(MinFloat(rp.SpentRuneReadyAt(sim, &rp.bloodRunes), rp.SpentRuneReadyAt(sim, &rp.frostRunes)), rp.SpentRuneReadyAt(sim, &rp.unholyRunes))
}

func (rp *CalcRunicPowerBar) AnyRuneReadyAt(sim *Simulation) float64 {
	return MinFloat(MinFloat(rp.RuneReadyAt(sim, &rp.bloodRunes), rp.RuneReadyAt(sim, &rp.frostRunes)), rp.RuneReadyAt(sim, &rp.unholyRunes))
}

func (rp *CalcRunicPowerBar) CurrentBloodRunes() int32 {
	return rp.CurrentRunesOfType(&rp.bloodRunes, RuneState_Normal)
}

func (rp *CalcRunicPowerBar) CurrentFrostRunes() int32 {
	return rp.CurrentRunesOfType(&rp.frostRunes, RuneState_Normal)
}

func (rp *CalcRunicPowerBar) CurrentUnholyRunes() int32 {
	return rp.CurrentRunesOfType(&rp.unholyRunes, RuneState_Normal)
}

func (rp *CalcRunicPowerBar) CurrentDeathRunes() int32 {
	return rp.CurrentRunesOfType(&rp.bloodRunes, RuneState_Death) + rp.CurrentRunesOfType(&rp.frostRunes, RuneState_Death) + rp.CurrentRunesOfType(&rp.unholyRunes, RuneState_Death)
}

func (rp *CalcRunicPowerBar) CurrentRunicPower() float64 {
	return rp.currentRunicPower
}

func (rp *CalcRunicPowerBar) DetermineOptimalCost(sim *Simulation, bloodAmount int, frostAmount int, unholyAmount int) DKRuneCost {
	totalBloodRunes := int(rp.CurrentBloodRunes())
	startingBloodRunes := totalBloodRunes
	totalFrostRunes := int(rp.CurrentFrostRunes())
	startingFrostRunes := totalFrostRunes
	totalUnholyRunes := int(rp.CurrentUnholyRunes())
	startingUnholyRunes := totalUnholyRunes
	totalDeathRunes := int(rp.CurrentDeathRunes())
	startingDeathRunes := totalDeathRunes

	if int(rp.CurrentBloodRunes()) >= bloodAmount {
		totalBloodRunes -= bloodAmount
	} else {
		if totalDeathRunes > 0 {
			totalDeathRunes -= bloodAmount
		}
	}

	if int(rp.CurrentFrostRunes()) >= frostAmount {
		totalFrostRunes -= frostAmount
	} else {
		if totalDeathRunes > 0 {
			totalDeathRunes -= frostAmount
		}
	}

	if int(rp.CurrentUnholyRunes()) >= unholyAmount {
		totalUnholyRunes -= unholyAmount
	} else {
		if totalDeathRunes > 0 {
			totalDeathRunes -= unholyAmount
		}
	}

	spellCost := DKRuneCost{
		Blood:  startingBloodRunes - totalBloodRunes,
		Frost:  startingFrostRunes - totalFrostRunes,
		Unholy: startingUnholyRunes - totalUnholyRunes,
		Death:  startingDeathRunes - totalDeathRunes,
	}

	return spellCost
}

func (rp *CalcRunicPowerBar) CastCostPossibleFor(sim *Simulation, currentRunes *DKRuneCost, bloodAmount int, frostAmount int, unholyAmount int) bool {
	totalDeathRunes := currentRunes.Death

	if currentRunes.Blood < bloodAmount {
		if totalDeathRunes > 0 {
			totalDeathRunes -= 1
		} else {
			return false
		}
	}

	if currentRunes.Frost < frostAmount {
		if totalDeathRunes > 0 {
			totalDeathRunes -= 1
		} else {
			return false
		}
	}

	if currentRunes.Unholy < unholyAmount {
		if totalDeathRunes > 0 {
			totalDeathRunes -= 1
		} else {
			return false
		}
	}

	return true
}

func (rp *CalcRunicPowerBar) CastCostPossible(sim *Simulation, runicPowerAmount float64, bloodAmount int32, frostAmount int32, unholyAmount int32) bool {
	totalDeathRunes := rp.CurrentDeathRunes()

	if rp.CurrentRunicPower() < runicPowerAmount {
		return false
	}

	if rp.CurrentBloodRunes() < bloodAmount {
		if totalDeathRunes > 0 {
			totalDeathRunes -= 1
		} else {
			return false
		}
	}

	if rp.CurrentFrostRunes() < frostAmount {
		if totalDeathRunes > 0 {
			totalDeathRunes -= 1
		} else {
			return false
		}
	}

	if rp.CurrentUnholyRunes() < unholyAmount {
		if totalDeathRunes > 0 {
			totalDeathRunes -= 1
		} else {
			return false
		}
	}

	return true
}

func (rp *CalcRunicPowerBar) DetermineRunesAfterCost(sim *Simulation, bloodAmount int, frostAmount int, unholyAmount int) DKRuneCost {
	spellCost := DKRuneCost{-1, -1, -1, -1}
	if rp.CastCostPossible(sim, 0, int32(bloodAmount), int32(frostAmount), int32(unholyAmount)) {
		totalBloodRunes := int(rp.CurrentBloodRunes())
		totalFrostRunes := int(rp.CurrentFrostRunes())
		totalUnholyRunes := int(rp.CurrentUnholyRunes())
		totalDeathRunes := int(rp.CurrentDeathRunes())

		if int(rp.CurrentBloodRunes()) >= bloodAmount {
			totalBloodRunes -= bloodAmount
		} else {
			totalDeathRunes -= bloodAmount
		}

		if int(rp.CurrentFrostRunes()) >= frostAmount {
			totalFrostRunes -= frostAmount
		} else {
			totalDeathRunes -= frostAmount
		}

		if int(rp.CurrentUnholyRunes()) >= unholyAmount {
			totalUnholyRunes -= unholyAmount
		} else {
			totalDeathRunes -= unholyAmount
		}

		spellCost = DKRuneCost{
			Blood:  totalBloodRunes,
			Frost:  totalFrostRunes,
			Unholy: totalUnholyRunes,
			Death:  totalDeathRunes,
		}
	}

	return spellCost
}

func (rp *CalcRunicPowerBar) SetRuneToState(r *CalcRune, runeState RuneState, runeKind RuneKind) {
	if (r.state == RuneState_Spent || r.state == RuneState_Normal) && (runeState == RuneState_Death || runeState == RuneState_DeathSpent) {
		r.kind = RuneKind_Death
	} else if (r.state == RuneState_DeathSpent || r.state == RuneState_Death) && (runeState != RuneState_Death && runeState != RuneState_DeathSpent) {
		if runeKind == RuneKind_Undef {
			panic("You have to set a rune kind here!")
		}
		r.kind = runeKind
	}

	r.state = runeState
}

func (rp *CalcRunicPowerBar) SetRuneAtIdxSlotToState(runeBarIdx int32, slot int32, runeState RuneState, runeKind RuneKind) {
	rb := &rp.bloodRunes
	if runeBarIdx == 1 {
		rb = &rp.frostRunes
	} else if runeBarIdx == 2 {
		rb = &rp.unholyRunes
	}

	// TODO: safeguard this?
	if (rb[slot].state == RuneState_Spent || rb[slot].state == RuneState_Normal) && (runeState == RuneState_Death || runeState == RuneState_DeathSpent) {
		rb[slot].kind = RuneKind_Death
	} else if (rb[slot].state == RuneState_DeathSpent || rb[slot].state == RuneState_Death) && (runeState != RuneState_Death && runeState != RuneState_DeathSpent) {
		if runeKind == RuneKind_Undef {
			panic("You have to set a rune kind here!")
		}
		rb[slot].kind = runeKind
	}

	rb[slot].state = runeState
}

func (rp *CalcRunicPowerBar) SetRuneAtSlotToState(rb *[2]CalcRune, slot int32, runeState RuneState, runeKind RuneKind) {
	// TODO: safeguard this?
	if (rb[slot].state == RuneState_Spent || rb[slot].state == RuneState_Normal) && (runeState == RuneState_Death || runeState == RuneState_DeathSpent) {
		rb[slot].kind = RuneKind_Death
	} else if (rb[slot].state == RuneState_DeathSpent || rb[slot].state == RuneState_Death) && (runeState != RuneState_Death && runeState != RuneState_DeathSpent) {
		if runeKind == RuneKind_Undef {
			panic("You have to set a rune kind here!")
		}
		rb[slot].kind = runeKind
	}

	rb[slot].state = runeState
}

func (rp *CalcRunicPowerBar) RegenRuneAndCancelPAs(sim *Simulation, r *CalcRune) {
	if r.state == RuneState_Spent {
		r.state = RuneState_Normal

		if r.pas[0] >= 0 {
			r.lastRegenTime = float64(sim.CurrentTime)
			r.pas[0] = -1
		}

		r.generatedByReapingOrBoTN = false
	} else if r.state == RuneState_DeathSpent {
		r.state = RuneState_Death

		if r.pas[0] >= 0 {
			r.lastRegenTime = float64(sim.CurrentTime)
			r.pas[0] = -1
		}

		r.generatedByReapingOrBoTN = false
	}
}

func (rp *CalcRunicPowerBar) RegenAllRunes(sim *Simulation) {
	rp.RegenRuneAndCancelPAs(sim, &rp.bloodRunes[0])
	rp.RegenRuneAndCancelPAs(sim, &rp.bloodRunes[1])
	rp.RegenRuneAndCancelPAs(sim, &rp.frostRunes[0])
	rp.RegenRuneAndCancelPAs(sim, &rp.frostRunes[1])
	rp.RegenRuneAndCancelPAs(sim, &rp.unholyRunes[0])
	rp.RegenRuneAndCancelPAs(sim, &rp.unholyRunes[1])
}

func (rp *CalcRunicPowerBar) GenerateRune(sim *Simulation, r *CalcRune) {
	if r.state == RuneState_Spent {
		if r.kind == RuneKind_Death {
			panic("Rune has wrong type for state.")
		}
		r.state = RuneState_Normal
		r.lastRegenTime = float64(sim.CurrentTime)
	} else if r.state == RuneState_DeathSpent {
		if r.kind != RuneKind_Death {
			panic("Rune has wrong type for state.")
		}
		r.state = RuneState_Death
		r.lastRegenTime = float64(sim.CurrentTime)
	}
}

func (rp *CalcRunicPowerBar) SpendRuneFromType(rb *[2]CalcRune, runeState RuneState) int32 {
	spendState := RuneState_Spent
	if runeState == RuneState_Death {
		spendState = RuneState_DeathSpent
	}

	slot := int32(-1)
	if rb[0].state == runeState {
		rb[0].state = spendState
		slot = 0
	} else if rb[1].state == runeState {
		rb[1].state = spendState
		slot = 1
	} else {
		panic("Trying to spend rune that does not exist!")
	}
	return slot
}

func (rp *CalcRunicPowerBar) LaunchRuneRegenPA(sim *Simulation, r *CalcRune) {
	runeGracePeriod := 0.0
	if r.lastRegenTime != -1 {
		runeGracePeriod = MinFloat(2.5, (float64(sim.CurrentTime)-r.lastRegenTime)/float64(1*time.Second))
	}
	r.pas[0] = float64(sim.CurrentTime) + 10.0 - runeGracePeriod
}

func (rp *CalcRunicPowerBar) SpendBloodRune(sim *Simulation) int32 {
	currRunes := rp.CurrentBloodRunes()
	if currRunes <= 0 {
		panic("Trying to spend blood runes that don't exist!")
	}

	spendSlot := rp.SpendRuneFromType(&rp.bloodRunes, RuneState_Normal)

	r := &rp.bloodRunes[spendSlot]
	rp.LaunchRuneRegenPA(sim, r)

	return spendSlot
}

func (rp *CalcRunicPowerBar) SpendFrostRune(sim *Simulation) int32 {
	currRunes := rp.CurrentFrostRunes()
	if currRunes <= 0 {
		panic("Trying to spend frost runes that don't exist!")
	}

	spendSlot := rp.SpendRuneFromType(&rp.frostRunes, RuneState_Normal)

	r := &rp.frostRunes[spendSlot]
	rp.LaunchRuneRegenPA(sim, r)

	return spendSlot
}

func (rp *CalcRunicPowerBar) SpendUnholyRune(sim *Simulation) int32 {
	currRunes := rp.CurrentUnholyRunes()
	if currRunes <= 0 {
		panic("Trying to spend unholy runes that don't exist!")
	}

	spendSlot := rp.SpendRuneFromType(&rp.unholyRunes, RuneState_Normal)

	r := &rp.unholyRunes[spendSlot]
	rp.LaunchRuneRegenPA(sim, r)

	return spendSlot
}

func (rp *CalcRunicPowerBar) SpendDeathRune(sim *Simulation) {
	currRunes := rp.CurrentDeathRunes()
	if currRunes <= 0 {
		panic("Trying to spend death runes that don't exist!")
	}

	runeTypeIdx := 0
	spendSlot := rp.SpendRuneFromType(&rp.bloodRunes, RuneState_Death)
	if spendSlot < 0 {
		runeTypeIdx += 1
		spendSlot = rp.SpendRuneFromType(&rp.frostRunes, RuneState_Death)
		if spendSlot < 0 {
			runeTypeIdx += 1
			spendSlot = rp.SpendRuneFromType(&rp.unholyRunes, RuneState_Death)
		}
	}

	r := &rp.bloodRunes[spendSlot]
	if runeTypeIdx == 1 {
		r = &rp.frostRunes[spendSlot]
	} else if runeTypeIdx == 2 {
		r = &rp.unholyRunes[spendSlot]
	}

	if r.generatedByReapingOrBoTN {
		rp.SetRuneToState(r, RuneState_Spent, RuneKind_Blood)
		r.generatedByReapingOrBoTN = false
	}

	rp.LaunchRuneRegenPA(sim, r)
}

func (rp *CalcRunicPowerBar) Spend(sim *Simulation, cost DKRuneCost) {
	for i := 0; i < cost.Blood; i++ {
		rp.SpendBloodRune(sim)
	}
	for i := 0; i < cost.Frost; i++ {
		rp.SpendFrostRune(sim)
	}
	for i := 0; i < cost.Unholy; i++ {
		rp.SpendUnholyRune(sim)
	}
	for i := 0; i < cost.Death; i++ {
		rp.SpendDeathRune(sim)
	}
}
