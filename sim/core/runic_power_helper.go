package core

func RunesBothOfState(sim *Simulation, runes *[2]Rune, runeState RuneState) bool {
	return runes[0].state == runeState && runes[1].state == runeState
}

func RunesAtleastOneOfState(sim *Simulation, runes *[2]Rune, runeState RuneState) bool {
	return runes[0].state == runeState || runes[1].state == runeState
}

func (rp *runicPowerBar) CorrectBloodTapConversion(sim *Simulation, bloodGainMetrics *ResourceMetrics, deathGainMetrics *ResourceMetrics, spell *Spell) {
	runes := &rp.bloodRunes

	startingBloodRunes := rp.CurrentBloodRunes()
	startingDeathRunes := rp.CurrentDeathRunes()

	//currRunes := rp.CurrentDeathRunes()
	//rp.GenerateRuneMetrics(sim, metrics, "Death", currRunes, currRunes+1)

	slot := int32(0)
	// Point 1
	if RunesBothOfState(sim, runes, RuneState_Normal) { // Point 1.1
		// Both are active, we convert leftmost into death rune
		slot = 0
		SetRuneAtSlotToState(runes, slot, RuneState_Death)
		rp.LaunchBloodTapRegenPA(sim, slot)
	} else if RunesBothOfState(sim, runes, RuneState_Spent) { // Point 1.2
		// Both are spent, we convert leftmost into death rune
		slot = 0
		SetRuneAtSlotToState(runes, slot, RuneState_Death)
		if runes[slot].pas[0] != nil {
			runes[slot].pas[0].Cancel(sim)
			runes[slot].pas[0] = nil
		}
		rp.LaunchBloodTapRegenPA(sim, slot)
	} else if RunesAtleastOneOfState(sim, runes, RuneState_Normal) && RunesAtleastOneOfState(sim, runes, RuneState_Spent) { // Point 2
		// One is active one is spent, we convert the active one into a death rune and the spent one remains spent
		slot = TernaryInt32(runes[0].state == RuneState_Normal, 0, 1)
		SetRuneAtSlotToState(runes, slot, RuneState_Death)
		rp.LaunchBloodTapRegenPA(sim, slot)
	} else if !RunesAtleastOneOfState(sim, runes, RuneState_Normal) && !RunesAtleastOneOfState(sim, runes, RuneState_Spent) { // Point 3
		// We have 2 death runes (spent or active)
		if RunesBothOfState(sim, runes, RuneState_DeathSpent) {
			// Both death runes are spent
			slot = 0
			SetRuneAtSlotToState(runes, slot, RuneState_Death)
			if runes[slot].pas[0] != nil {
				runes[slot].pas[0].Cancel(sim)
				runes[slot].pas[0] = nil
			}
			rp.LaunchBloodTapRegenPA(sim, slot)
		} else if RunesBothOfState(sim, runes, RuneState_Death) {
			// Both death runes are active
			// Reset CD of highest CD one?
		} else {
			// Only one death rune is spent
			slot = TernaryInt32(runes[0].state == RuneState_DeathSpent, 0, 1)
			SetRuneAtSlotToState(runes, slot, RuneState_Death)
			if runes[slot].pas[0] != nil {
				runes[slot].pas[0].Cancel(sim)
				runes[slot].pas[0] = nil
			}
			rp.LaunchBloodTapRegenPA(sim, slot)
		}
	} else if RunesAtleastOneOfState(sim, runes, RuneState_Spent) && RunesAtleastOneOfState(sim, runes, RuneState_DeathSpent) { // Point 5
		// One spent blood rune and one spent death rune, we convert the spent blood rune to a spent death rune and activate the other spent death rune
		slot = TernaryInt32(runes[0].state == RuneState_Spent, 0, 1)
		SetRuneAtSlotToState(runes, slot, RuneState_DeathSpent)
		if runes[slot].pas[0] != nil {
			runes[slot].pas[0].Cancel(sim)
			runes[slot].pas[0] = nil
		}
		rp.LaunchBloodTapRegenPA(sim, slot)

		slot = TernaryInt32(runes[0].state == RuneState_DeathSpent, 0, 1)
		SetRuneAtSlotToState(runes, slot, RuneState_Death)
		if runes[slot].pas[0] != nil {
			runes[slot].pas[0].Cancel(sim)
			runes[slot].pas[0] = nil
		}
		rp.LaunchBloodTapRegenPA(sim, slot)
	} else if RunesAtleastOneOfState(sim, runes, RuneState_Normal) && RunesAtleastOneOfState(sim, runes, RuneState_Death) { // Point 4
		// One active blood rune && one active death rune, we convert the blood rune into a death rune
		slot = TernaryInt32(runes[0].state == RuneState_Normal, 0, 1)
		SetRuneAtSlotToState(runes, slot, RuneState_Death)
	} else if RunesAtleastOneOfState(sim, runes, RuneState_Normal) && RunesAtleastOneOfState(sim, runes, RuneState_DeathSpent) ||
		RunesAtleastOneOfState(sim, runes, RuneState_Spent) && RunesAtleastOneOfState(sim, runes, RuneState_Death) {
		// We have one blood rune and one death rune where one is active and the other spent (no mattter which)
		if RunesAtleastOneOfState(sim, runes, RuneState_Normal) && RunesAtleastOneOfState(sim, runes, RuneState_DeathSpent) {
			// We have an active blood rune and a spent death rune, we convert the blood rune into a death rune and activate the spent death rune
			slot = TernaryInt32(runes[0].state == RuneState_Normal, 0, 1)
			SetRuneAtSlotToState(runes, slot, RuneState_Death)
			slot = TernaryInt32(runes[0].state == RuneState_DeathSpent, 0, 1)
			SetRuneAtSlotToState(runes, slot, RuneState_Death)
			if runes[slot].pas[0] != nil {
				runes[slot].pas[0].Cancel(sim)
				runes[slot].pas[0] = nil
			}
			rp.LaunchBloodTapRegenPA(sim, slot)
		} else {
			// We have an active death rune and a spent blood rune, we convert the blood rune into a death rune
			slot = TernaryInt32(runes[0].state == RuneState_Spent, 0, 1)
			SetRuneAtSlotToState(runes, slot, RuneState_Death)
			if runes[slot].pas[0] != nil {
				runes[slot].pas[0].Cancel(sim)
				runes[slot].pas[0] = nil
			}
			rp.LaunchBloodTapRegenPA(sim, slot)
		}
	}

	currBloodRunes := rp.CurrentBloodRunes()
	currDeathRunes := rp.CurrentDeathRunes()
	totalChangeInBloodRunes := currBloodRunes - startingBloodRunes
	totalChangeInDeathRunes := currDeathRunes - startingDeathRunes

	if totalChangeInBloodRunes > 0 {
		rp.GenerateRuneMetrics(sim, bloodGainMetrics, "Blood", startingBloodRunes, currBloodRunes)
	} else if totalChangeInBloodRunes < 0 {
		rp.SpendRuneMetrics(sim, spell.BloodRuneMetrics(), "Blood", startingBloodRunes, currBloodRunes)
	}

	if totalChangeInDeathRunes > 0 {
		rp.GenerateRuneMetrics(sim, deathGainMetrics, "Death", startingDeathRunes, currDeathRunes)
	} else if totalChangeInDeathRunes < 0 {
		rp.SpendRuneMetrics(sim, spell.DeathRuneMetrics(), "Death", startingDeathRunes, currDeathRunes)
	}
}
