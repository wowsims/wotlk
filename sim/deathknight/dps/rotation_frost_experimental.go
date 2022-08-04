package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) frRecastAvailableCheck(expiresAt time.Duration, afterCastTime time.Duration, spellCost int, currentRunes int32, nextRuneAt time.Duration) bool {
	if spellCost > 0 && currentRunes == 0 {
		if expiresAt < nextRuneAt {
			return true
		}
	} else if afterCastTime > expiresAt {
		return true
	}
	return false
}

func (dk *DpsDeathknight) frDiseaseCheck(sim *core.Simulation, target *core.Unit, spell *deathknight.RuneSpell, costRunes bool, casts int) bool {
	ffRemaining := dk.FrostFeverDisease[target.Index].RemainingDuration(sim)
	bpRemaining := dk.BloodPlagueDisease[target.Index].RemainingDuration(sim)
	castGcd := dk.SpellGCD() * time.Duration(casts)

	// FF is not active or will drop before Gcd is ready after this cast
	if !dk.FrostFeverDisease[target.Index].IsActive() || ffRemaining < castGcd {
		return false
	}
	// BP is not active or will drop before Gcd is ready after this cast
	if !dk.BloodPlagueDisease[target.Index].IsActive() || bpRemaining < castGcd {
		return false
	}

	// If the ability we want to cast spends runes we check for possible disease drops
	// in the time we won't have runes to recast the disease
	if dk.CanCast(sim, spell) && costRunes {
		ffExpiresAt := ffRemaining + sim.CurrentTime
		bpExpiresAt := bpRemaining + sim.CurrentTime
		diseasesExpireAt := core.MinDuration(ffExpiresAt, bpExpiresAt)

		crpb := dk.CopyRunicPowerBar()
		spellCost := crpb.OptimalRuneCost(core.RuneCost(spell.DefaultCast.Cost))

		crpb.SpendRuneCost(sim, spell.Spell, spellCost)

		afterCastTime := sim.CurrentTime + castGcd
		currentBloodRunes := crpb.CurrentBloodRunes() + crpb.CurrentDeathRunes()
		nextBloodRuneAt := crpb.BloodRuneReadyAt(sim)

		// If FF is gonna drop while our runes are on CD
		if dk.frRecastAvailableCheck(diseasesExpireAt, afterCastTime, int(spellCost.Blood()), currentBloodRunes, nextBloodRuneAt) {
			return false
		}
	}

	return true
}
