package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

type FrostRotation struct {
	lastSpell *deathknight.RuneSpell
	nextSpell *deathknight.RuneSpell
}

func (fr *FrostRotation) Reset(sim *core.Simulation) {
	fr.nextSpell = nil
	fr.lastSpell = nil
}

func (fr *FrostRotation) SetLastSpell(condition bool, spell *deathknight.RuneSpell) {
	if condition {
		fr.lastSpell = spell
	}
}

func (dk *DpsDeathknight) FrostDiseaseCheck(sim *core.Simulation, target *core.Unit, spell *deathknight.RuneSpell, costRunes bool, casts int) bool {
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
		spellCost := crpb.OptimalRuneCost(core.RuneCost(spell.DefaultCast.Cost))

		crpb.SpendRuneCost(sim, spell.Spell, spellCost)

		afterCastTime := sim.CurrentTime + castGcd
		currentBloodRunes := crpb.CurrentBloodRunes()
		nextBloodRuneAt := crpb.NormalBloodRuneReadyAt(sim)

		// Check FF
		if dk.frCheckForDiseaseRecast(ffExpiresAt, afterCastTime, int(spellCost.Blood()), currentBloodRunes, nextBloodRuneAt) {
			return false
		}

		// Check BP
		if dk.frCheckForDiseaseRecast(bpExpiresAt, afterCastTime, int(spellCost.Blood()), currentBloodRunes, nextBloodRuneAt) {
			return false
		}
	} else {
		return false
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
