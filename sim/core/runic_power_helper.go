package core

import (
	"fmt"
	"time"
)

type RuneCost uint16

func NewRuneCost(rp, blood, frost, unholy, death uint8) RuneCost {
	value := int16(0)
	if blood == 1 {
		value = 1
	} else if blood == 2 {
		value = 3
	}

	if frost == 1 {
		value += 1 << 2
	} else if frost == 2 {
		value += 3 << 2
	}

	if unholy == 1 {
		value += 1 << 4
	} else if unholy == 2 {
		value += 3 << 4
	}

	if death == 1 {
		value += 1 << 6
	} else if death == 2 {
		value += 3 << 6
	}

	value += int16(rp) << 8

	return RuneCost(value)
}

func (rc RuneCost) String() string {
	return fmt.Sprintf("RP: %d, Blood: %d, Frost: %d, Unholy: %d, Death: %d", rc.RunicPower(), rc.Blood(), rc.Frost(), rc.Unholy(), rc.Death())
}

// HasRune returns if this cost includes a rune portion.
//  If any bit is set in the rune bits it means that there is a rune cost.
func (rc RuneCost) HasRune() bool {
	const runebits = int16(0b11111111)
	return runebits&int16(rc) > 0
}

func (rc RuneCost) RunicPower() uint8 {
	const rpbits = uint16(0b1111111100000000)
	return uint8((uint16(rc) & rpbits) >> 8)
}

func (rc RuneCost) Blood() uint8 {
	runes := uint16(rc) & 0b11
	switch runes {
	case 0b00:
		return 0
	case 0b01:
		return 1
	case 0b11:
		return 2
	}
	return 0
}

func (rc RuneCost) Frost() uint8 {
	runes := uint16(rc) & 0b1100
	switch runes {
	case 0:
		return 0
	case 0b0100:
		return 1
	case 0b1100:
		return 2
	}
	return 0
}

func (rc RuneCost) Unholy() uint8 {
	runes := uint16(rc) & 0b110000
	switch runes {
	case 0:
		return 0
	case 0b010000:
		return 1
	case 0b110000:
		return 2
	}
	return 0
}

func (rc RuneCost) Death() uint8 {
	runes := uint16(rc) & 0b11000000
	switch runes {
	case 0:
		return 0
	case 0b01000000:
		return 1
	case 0b11000000:
		return 2
	}
	return 0
}

func RunesBothOfState(sim *Simulation, runes *[2]Rune, runeState RuneState) bool {
	return runes[0].state == runeState && runes[1].state == runeState
}

func RunesAtleastOneOfState(sim *Simulation, runes *[2]Rune, runeState RuneState) bool {
	return runes[0].state == runeState || runes[1].state == runeState
}

// TODO: Simplify this, its definitely possible
func (rp *runicPowerBar) LaunchBloodTapRegenPA(sim *Simulation, slot int32, spell *Spell) {
	r := &rp.bloodRunes[slot]

	pa := &PendingAction{
		NextActionAt: sim.CurrentTime + 20.0*time.Second,
		Priority:     ActionPriorityRegen,
	}

	pa.OnAction = func(sim *Simulation) {
		if !pa.cancelled {
			r.pas[1].Cancel(sim)
			r.pas[1] = nil
			if r.state == RuneState_Death {
				currRunes := rp.CurrentBloodRunes()
				rp.GainRuneMetrics(sim, rp.bloodRuneGainMetrics, "blood", currRunes, currRunes+1)
				rp.SetRuneToState(r, RuneState_Normal, RuneKind_Blood)

				currRunes = rp.CurrentDeathRunes()
				rp.SpendRuneMetrics(sim, spell.DeathRuneMetrics(), "death", currRunes, currRunes-1)
				if !rp.isACopy {
					rp.onBloodRuneGain(sim)
				}
			} else if r.state == RuneState_DeathSpent {

				if r.pas[0] == nil {
					panic("This should have a regen PA!")
				}

				rp.SetRuneToState(r, RuneState_Spent, RuneKind_Blood)
			}
		} else {
			r.pas[1] = nil
		}
	}

	r.pas[1] = pa
	//if !rp.isACopy {
	//sim.AddPendingAction(pa)
	//}
}

func (rp *runicPowerBar) GainDeathRuneMetrics(sim *Simulation, spell *Spell, currRunes int32, newRunes int32) {
	if !rp.isACopy {
		metrics := rp.deathRuneGainMetrics
		metrics.AddEvent(1, float64(newRunes)-float64(currRunes))

		if sim.Log != nil {
			rp.unit.Log(sim, "Gained 1.000 death rune from %s (%d --> %d).", metrics.ActionID, currRunes, newRunes)
		}
	}
}

func (rp *runicPowerBar) SpendBloodRuneMetrics(sim *Simulation, spell *Spell, currRunes int32, newRunes int32) {
	if !rp.isACopy {
		metrics := spell.BloodRuneMetrics()

		metrics.AddEvent(-1, -1)

		if sim.Log != nil {
			rp.unit.Log(sim, "Spent 1.000 blood rune from %s (%d --> %d).", metrics.ActionID, currRunes, newRunes)
		}
	}
}

func (rp *runicPowerBar) CancelRuneRegenPA(sim *Simulation, r *Rune) {
	if r.pas[0] == nil {
		panic("Trying to cancel non-existant regen PA.")
	}

	if r.pas[0] != nil {
		r.pas[0].Cancel(sim)
	}
	r.pas[0] = nil
}

func (rp *runicPowerBar) CancelBloodTap(sim *Simulation) {
	runes := &rp.bloodRunes

	if runes[0].pas[1] != nil {
		runes[0].pas[1].OnAction(sim)
	} else if runes[1].pas[1] != nil {
		runes[1].pas[1].OnAction(sim)
	}
}

func (rp *runicPowerBar) CorrectBloodTapConversion(sim *Simulation, bloodGainMetrics *ResourceMetrics, deathGainMetrics *ResourceMetrics, spell *Spell) {
	runes := &rp.bloodRunes

	currBloodRunes := rp.CurrentBloodRunes()
	currDeathRunes := rp.CurrentDeathRunes()

	slot := int32(0)
	// Point 1
	if RunesBothOfState(sim, runes, RuneState_Normal) { // Point 1.1
		// Both are active, we convert leftmost into death rune
		slot = 0
		rp.SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
		rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+1)
		rp.SpendBloodRuneMetrics(sim, spell, currBloodRunes, currBloodRunes-1)
		rp.LaunchBloodTapRegenPA(sim, slot, spell)
	} else if RunesBothOfState(sim, runes, RuneState_Spent) { // Point 1.2
		// Both are spent, we convert leftmost into death rune
		slot = 0
		rp.SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
		rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+1)
		rp.CancelRuneRegenPA(sim, &runes[slot])
		rp.LaunchBloodTapRegenPA(sim, slot, spell)
	} else if RunesAtleastOneOfState(sim, runes, RuneState_Normal) && RunesAtleastOneOfState(sim, runes, RuneState_Spent) { // Point 2
		// One is active one is spent, we convert the active one into a death rune and the spent one remains spent
		slot = TernaryInt32(runes[0].state == RuneState_Normal, 0, 1)
		rp.SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
		rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+1)
		rp.SpendBloodRuneMetrics(sim, spell, currBloodRunes, currBloodRunes-1)
		rp.LaunchBloodTapRegenPA(sim, slot, spell)
	} else if !RunesAtleastOneOfState(sim, runes, RuneState_Normal) && !RunesAtleastOneOfState(sim, runes, RuneState_Spent) { // Point 3
		// We have 2 death runes (spent or active)
		if RunesBothOfState(sim, runes, RuneState_DeathSpent) {
			// Both death runes are spent
			slot = 0
			rp.SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
			rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+1)
			rp.CancelRuneRegenPA(sim, &runes[slot])
			rp.LaunchBloodTapRegenPA(sim, slot, spell)
		} else if RunesBothOfState(sim, runes, RuneState_Death) {
			// Both death runes are active
			// Reset CD of highest CD one?
		} else {
			// Only one death rune is spent
			slot = TernaryInt32(runes[0].state == RuneState_DeathSpent, 0, 1)
			rp.SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
			rp.CancelRuneRegenPA(sim, &runes[slot])
			rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+1)
			rp.LaunchBloodTapRegenPA(sim, slot, spell)
		}
	} else if RunesAtleastOneOfState(sim, runes, RuneState_Spent) && RunesAtleastOneOfState(sim, runes, RuneState_DeathSpent) { // Point 5
		// One spent blood rune and one spent death rune, we convert the spent blood rune to a spent death rune and activate the other spent death rune
		slot = TernaryInt32(runes[0].state == RuneState_Spent, 0, 1)
		rp.SetRuneAtSlotToState(runes, slot, RuneState_DeathSpent, RuneKind_Death)
		rp.LaunchBloodTapRegenPA(sim, slot, spell)

		slot = TernaryInt32(runes[0].state == RuneState_DeathSpent, 0, 1)
		rp.SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
		rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+2)
		rp.CancelRuneRegenPA(sim, &runes[slot])
	} else if RunesAtleastOneOfState(sim, runes, RuneState_Normal) && RunesAtleastOneOfState(sim, runes, RuneState_Death) { // Point 4
		// One active blood rune && one active death rune, we convert the blood rune into a death rune
		slot = TernaryInt32(runes[0].state == RuneState_Normal, 0, 1)
		rp.SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
		rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+1)
		rp.SpendBloodRuneMetrics(sim, spell, currBloodRunes, currBloodRunes-1)
	} else if RunesAtleastOneOfState(sim, runes, RuneState_Normal) && RunesAtleastOneOfState(sim, runes, RuneState_DeathSpent) ||
		RunesAtleastOneOfState(sim, runes, RuneState_Spent) && RunesAtleastOneOfState(sim, runes, RuneState_Death) {
		// We have one blood rune and one death rune where one is active and the other spent (no mattter which)
		if RunesAtleastOneOfState(sim, runes, RuneState_Normal) && RunesAtleastOneOfState(sim, runes, RuneState_DeathSpent) {
			// We have an active blood rune and a spent death rune, we convert the blood rune into a death rune and activate the spent death rune
			slot = TernaryInt32(runes[0].state == RuneState_Normal, 0, 1)
			rp.SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
			rp.SpendBloodRuneMetrics(sim, spell, currBloodRunes, currBloodRunes-1)
			rp.LaunchBloodTapRegenPA(sim, slot, spell)

			slot = TernaryInt32(runes[0].state == RuneState_DeathSpent, 0, 1)
			rp.SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
			rp.CancelRuneRegenPA(sim, &runes[slot])
			rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+2)
		} else {
			// We have an active death rune and a spent blood rune, we convert the blood rune into a death rune
			slot = TernaryInt32(runes[0].state == RuneState_Spent, 0, 1)
			rp.SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
			rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+1)
			rp.CancelRuneRegenPA(sim, &runes[slot])
			rp.LaunchBloodTapRegenPA(sim, slot, spell)
		}
	}

	if !rp.isACopy {
		rp.onDeathRuneGain(sim)
	}
}

// so in english
// 1. try to convert active blood rune -> death rune
// 2. if no active blood, convert inactive blood rune -> death rune, and then convert one inactive death rune -> active

// psuedocode
// rune = findFirstActiveBlood()
// if !rune {
//   rune = findFirstInactiveBlood()
// }
// // possible we already have 2 death runes
// if rune {
//   rune.death = true
// }

// deathrune = findFirstInactiveDeath()
// if deathrune {
//   deathrune.active = true
// }
