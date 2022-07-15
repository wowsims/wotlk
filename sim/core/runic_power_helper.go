package core

import "time"

func RunesBothOfState(sim *Simulation, runes *[2]Rune, runeState RuneState) bool {
	return runes[0].state == runeState && runes[1].state == runeState
}

func RunesAtleastOneOfState(sim *Simulation, runes *[2]Rune, runeState RuneState) bool {
	return runes[0].state == runeState || runes[1].state == runeState
}

func (rp *runicPowerBar) FlagBloodRuneSlotAsBoTN(slot int32) {
	rp.bloodRunes[slot].botnFlag = true
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
			if r.state == RuneState_Death {
				r.pas[1] = nil

				currRunes := rp.CurrentBloodRunes()
				rp.GainRuneMetrics(sim, rp.bloodRuneGainMetrics, "blood", currRunes, currRunes+1)
				rp.SetRuneToState(r, RuneState_Normal, RuneKind_Blood)
				rp.onBloodRuneGain(sim)

				currRunes = rp.CurrentDeathRunes()
				rp.SpendRuneMetrics(sim, spell.DeathRuneMetrics(), "death", currRunes, currRunes-1)
			} else if r.state == RuneState_DeathSpent {
				r.pas[1] = nil

				currPA := r.pas[0]
				if currPA == nil {
					panic("This should have a regen PA!")
				}

				timeTillPA := currPA.NextActionAt
				currPA.Cancel(sim)
				rp.SetRuneToState(r, RuneState_Spent, RuneKind_Blood)

				pa := &PendingAction{
					NextActionAt: timeTillPA,
					Priority:     ActionPriorityRegen,
				}

				pa.OnAction = func(sim *Simulation) {
					if !pa.cancelled {
						r.pas[0] = nil
						currRunes := rp.CurrentBloodRunes()
						rp.GainRuneMetrics(sim, rp.bloodRuneGainMetrics, "blood", currRunes, currRunes+1)
						rp.SetRuneToState(r, RuneState_Normal, RuneKind_Blood)
						rp.onBloodRuneGain(sim)
					} else {
						r.pas[0] = nil
					}
				}

				r.pas[0] = pa
				sim.AddPendingAction(pa)
			}
		} else {
			r.pas[1] = nil
		}
	}

	r.pas[1] = pa
	sim.AddPendingAction(pa)
}

func (rp *runicPowerBar) GainDeathRuneMetrics(sim *Simulation, spell *Spell, currRunes int32, newRunes int32) {
	metrics := rp.deathRuneGainMetrics
	metrics.AddEvent(1, float64(newRunes)-float64(currRunes))

	if sim.Log != nil {
		rp.unit.Log(sim, "Gained 1.000 death rune from %s (%d --> %d).", metrics.ActionID, currRunes, newRunes)
	}
}

func (rp *runicPowerBar) SpendBloodRuneMetrics(sim *Simulation, spell *Spell, currRunes int32, newRunes int32) {
	metrics := spell.BloodRuneMetrics()

	metrics.AddEvent(-1, -1)

	if sim.Log != nil {
		rp.unit.Log(sim, "Spent 1.000 blood rune from %s (%d --> %d).", metrics.ActionID, currRunes, newRunes)
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
		SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
		rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+1)
		rp.SpendBloodRuneMetrics(sim, spell, currBloodRunes, currBloodRunes-1)
		rp.LaunchBloodTapRegenPA(sim, slot, spell)
	} else if RunesBothOfState(sim, runes, RuneState_Spent) { // Point 1.2
		// Both are spent, we convert leftmost into death rune
		slot = 0
		SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
		rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+1)
		runes[slot].pas[0].Cancel(sim)
		runes[slot].pas[0] = nil
		rp.LaunchBloodTapRegenPA(sim, slot, spell)
	} else if RunesAtleastOneOfState(sim, runes, RuneState_Normal) && RunesAtleastOneOfState(sim, runes, RuneState_Spent) { // Point 2
		// One is active one is spent, we convert the active one into a death rune and the spent one remains spent
		slot = TernaryInt32(runes[0].state == RuneState_Normal, 0, 1)
		SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
		rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+1)
		rp.SpendBloodRuneMetrics(sim, spell, currBloodRunes, currBloodRunes-1)
		rp.LaunchBloodTapRegenPA(sim, slot, spell)
	} else if !RunesAtleastOneOfState(sim, runes, RuneState_Normal) && !RunesAtleastOneOfState(sim, runes, RuneState_Spent) { // Point 3
		// We have 2 death runes (spent or active)
		if RunesBothOfState(sim, runes, RuneState_DeathSpent) {
			// Both death runes are spent
			slot = 0
			SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
			rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+1)
			runes[slot].pas[0].Cancel(sim)
			runes[slot].pas[0] = nil
			rp.LaunchBloodTapRegenPA(sim, slot, spell)
		} else if RunesBothOfState(sim, runes, RuneState_Death) {
			// Both death runes are active
			// Reset CD of highest CD one?
		} else {
			// Only one death rune is spent
			slot = TernaryInt32(runes[0].state == RuneState_DeathSpent, 0, 1)
			SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
			rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+1)
			runes[slot].pas[0].Cancel(sim)
			runes[slot].pas[0] = nil
			rp.LaunchBloodTapRegenPA(sim, slot, spell)
		}
	} else if RunesAtleastOneOfState(sim, runes, RuneState_Spent) && RunesAtleastOneOfState(sim, runes, RuneState_DeathSpent) { // Point 5
		// One spent blood rune and one spent death rune, we convert the spent blood rune to a spent death rune and activate the other spent death rune
		slot = TernaryInt32(runes[0].state == RuneState_Spent, 0, 1)
		SetRuneAtSlotToState(runes, slot, RuneState_DeathSpent, RuneKind_Death)
		runes[slot].pas[0].Cancel(sim)
		runes[slot].pas[0] = nil
		rp.LaunchBloodTapRegenPA(sim, slot, spell)

		slot = TernaryInt32(runes[0].state == RuneState_DeathSpent, 0, 1)
		SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
		rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+2)
		runes[slot].pas[0].Cancel(sim)
		runes[slot].pas[0] = nil
		rp.LaunchBloodTapRegenPA(sim, slot, spell)
	} else if RunesAtleastOneOfState(sim, runes, RuneState_Normal) && RunesAtleastOneOfState(sim, runes, RuneState_Death) { // Point 4
		// One active blood rune && one active death rune, we convert the blood rune into a death rune
		slot = TernaryInt32(runes[0].state == RuneState_Normal, 0, 1)
		SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
		rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+1)
		rp.SpendBloodRuneMetrics(sim, spell, currBloodRunes, currBloodRunes-1)
	} else if RunesAtleastOneOfState(sim, runes, RuneState_Normal) && RunesAtleastOneOfState(sim, runes, RuneState_DeathSpent) ||
		RunesAtleastOneOfState(sim, runes, RuneState_Spent) && RunesAtleastOneOfState(sim, runes, RuneState_Death) {
		// We have one blood rune and one death rune where one is active and the other spent (no mattter which)
		if RunesAtleastOneOfState(sim, runes, RuneState_Normal) && RunesAtleastOneOfState(sim, runes, RuneState_DeathSpent) {
			// We have an active blood rune and a spent death rune, we convert the blood rune into a death rune and activate the spent death rune
			slot = TernaryInt32(runes[0].state == RuneState_Normal, 0, 1)
			SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
			rp.SpendBloodRuneMetrics(sim, spell, currBloodRunes, currBloodRunes-1)
			rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+2)
			slot = TernaryInt32(runes[0].state == RuneState_DeathSpent, 0, 1)
			SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
			runes[slot].pas[0].Cancel(sim)
			runes[slot].pas[0] = nil
			rp.LaunchBloodTapRegenPA(sim, slot, spell)
		} else {
			// We have an active death rune and a spent blood rune, we convert the blood rune into a death rune
			slot = TernaryInt32(runes[0].state == RuneState_Spent, 0, 1)
			SetRuneAtSlotToState(runes, slot, RuneState_Death, RuneKind_Death)
			rp.GainDeathRuneMetrics(sim, spell, currDeathRunes, currDeathRunes+1)
			runes[slot].pas[0].Cancel(sim)
			runes[slot].pas[0] = nil
			rp.LaunchBloodTapRegenPA(sim, slot, spell)
		}
	}
}
