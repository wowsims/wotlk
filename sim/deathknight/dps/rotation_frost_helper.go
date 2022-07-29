package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

type FrostRotationAction uint8

const (
	FrostRotationAction_Obli FrostRotationAction = iota
	FrostRotationAction_BS
	FrostRotationAction_Pesti
)

type FrostRotation struct {
	idx        int
	numActions int
	actions    []FrostRotationAction
}

func (fr *FrostRotation) Reset(sim *core.Simulation) {
	fr.idx = 0
	fr.actions = []FrostRotationAction{
		FrostRotationAction_Obli,
		FrostRotationAction_Obli,
		FrostRotationAction_BS,
		FrostRotationAction_Pesti,
	}
	fr.numActions = len(fr.actions)
}

func (dk *DpsDeathknight) FrostDiseaseCheck(sim *core.Simulation, target *core.Unit, spell *core.Spell, costRunes bool, casts int) bool {
	ffRemaining := dk.FrostFeverDisease[target.Index].RemainingDuration(sim)
	bpRemaining := dk.BloodPlagueDisease[target.Index].RemainingDuration(sim)
	castGcd := dk.SpellGCD() * time.Duration(casts)

	if !dk.FrostFeverDisease[target.Index].IsActive() || ffRemaining < castGcd {
		// Refresh FF
		return false
	}
	if !dk.BloodPlagueDisease[target.Index].IsActive() || bpRemaining < castGcd {
		// Refresh BP
		return false
	}

	if dk.CanCast(sim, spell) && costRunes {
		ffExpiresAt := ffRemaining + sim.CurrentTime
		bpExpiresAt := bpRemaining + sim.CurrentTime

		crpb := dk.CopyRunicPowerBar()
		runeCostForSpell := dk.RuneAmountForSpell(spell)
		spellCost := crpb.DetermineOptimalCost(sim, runeCostForSpell.Blood, runeCostForSpell.Frost, runeCostForSpell.Unholy)

		crpb.Spend(sim, spell, spellCost)

		afterCastTime := sim.CurrentTime + castGcd
		currentBloodRunes := crpb.CurrentBloodRunes()
		nextBloodRuneAt := crpb.BloodRuneReadyAt(sim)

		// Check FF
		if dk.frCheckForDiseaseRecast(ffExpiresAt, afterCastTime, spellCost.Blood, currentBloodRunes, nextBloodRuneAt) {
			return false
		}

		// Check BP
		if dk.frCheckForDiseaseRecast(bpExpiresAt, afterCastTime, spellCost.Blood, currentBloodRunes, nextBloodRuneAt) {
			return false
		}
	}

	return true
}

func (dk *DpsDeathknight) frCheckForDiseaseRecast(expiresAt time.Duration, afterCastTime time.Duration,
	spellCost int, currentRunes int32, nextRuneAt time.Duration) bool {
	if spellCost > 0 && currentRunes == 0 {
		if expiresAt < nextRuneAt {
			return true
		}
	} else if afterCastTime > expiresAt {
		return true
	}
	return false
}
