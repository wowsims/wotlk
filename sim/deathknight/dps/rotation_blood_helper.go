package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

type BloodRotation struct {
	dk *DpsDeathknight

	drwSnapshot *core.SnapshotManager
	drwMaxDelay time.Duration

	bloodSpell *core.Spell

	activatingDrw bool
}

func (br *BloodRotation) Reset(sim *core.Simulation) {
	br.activatingDrw = false
	if br.drwSnapshot != nil {
		br.drwSnapshot.ResetProcTrackers()
	}
	br.drwMaxDelay = -1
}

func (br *BloodRotation) Initialize(dk *DpsDeathknight) {
}

func (dk *DpsDeathknight) blBloodRuneAction() deathknight.RotationAction {
	if dk.Env.GetNumTargets() > 1 {
		return dk.RotationActionCallback_Pesti
	} else {
		return dk.RotationActionBL_BS
	}
}

func (dk *DpsDeathknight) blDiseaseCheck(sim *core.Simulation, target *core.Unit, spell *core.Spell, costRunes bool, casts int) bool {
	// Early exit at end of fight
	if sim.GetRemainingDuration() < 10*time.Second {
		return true
	}

	ffRemaining := dk.FrostFeverSpell.Dot(target).RemainingDuration(sim)
	bpRemaining := dk.BloodPlagueSpell.Dot(target).RemainingDuration(sim)
	castGcd := core.GCDDefault * time.Duration(casts)

	// FF is not active or will drop before Gcd is ready after this cast
	if !dk.FrostFeverSpell.Dot(target).IsActive() || ffRemaining <= castGcd {
		return false
	}
	// BP is not active or will drop before Gcd is ready after this cast
	if !dk.BloodPlagueSpell.Dot(target).IsActive() || bpRemaining <= castGcd {
		return false
	}

	// If the ability we want to cast spends runes we check for possible disease drops
	// in the time we won't have runes to recast the disease
	if spell.CanCast(sim, nil) && costRunes {
		ffExpiresAt := ffRemaining + sim.CurrentTime
		bpExpiresAt := bpRemaining + sim.CurrentTime

		afterCastTime := sim.CurrentTime + castGcd
		if ffExpiresAt <= afterCastTime || bpExpiresAt <= afterCastTime {
			return false
		}

		crpb := dk.CopyRunicPowerBar()
		spellCost := crpb.OptimalRuneCost(core.RuneCost(spell.DefaultCast.Cost))

		crpb.SpendRuneCost(sim, spell, spellCost)

		if dk.sr.hasGod {
			currentBloodRunes := crpb.CurrentBloodRunes()
			nextBloodRuneAt := crpb.BloodRuneReadyAt(sim)

			// If FF is gonna drop while our runes are on CD
			if dk.shRecastAvailableCheck(ffExpiresAt, afterCastTime, int(spellCost.Blood()), int32(currentBloodRunes), nextBloodRuneAt) {
				return false
			}

			// If BP is gonna drop while our runes are on CD
			if dk.shRecastAvailableCheck(bpExpiresAt, afterCastTime, int(spellCost.Blood()), int32(currentBloodRunes), nextBloodRuneAt) {
				return false
			}
		} else {
			currentFrostRunes := crpb.CurrentFrostRunes()
			currentUnholyRunes := crpb.CurrentUnholyRunes()
			nextFrostRuneAt := crpb.FrostRuneReadyAt(sim)
			nextUnholyRuneAt := crpb.UnholyRuneReadyAt(sim)

			// If FF is gonna drop while our runes are on CD
			if dk.shRecastAvailableCheck(ffExpiresAt, afterCastTime, int(spellCost.Frost()), int32(currentFrostRunes), nextFrostRuneAt) {
				return false
			}

			// If BP is gonna drop while our runes are on CD
			if dk.shRecastAvailableCheck(bpExpiresAt, afterCastTime, int(spellCost.Unholy()), int32(currentUnholyRunes), nextUnholyRuneAt) {
				return false
			}
		}
	}

	return true
}

func (dk *DpsDeathknight) blSpreadDiseases(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.blDiseaseCheck(sim, target, dk.Pestilence, true, 1) {
		casted := dk.Pestilence.Cast(sim, target)
		landed := dk.LastOutcome.Matches(core.OutcomeLanded)

		// Reset flags on succesfull cast
		dk.sr.recastedFF = !(casted && landed)
		dk.sr.recastedBP = !(casted && landed)
		return -1
	} else {
		dk.blRecastDiseasesSequence(sim)
		return sim.CurrentTime
	}
}

// Save up Runic Power for DRW - Allow casts above 100 RP when DRW is ready or above 85 (for death strike glyph) when not
func (dk *DpsDeathknight) blDeathCoilCheck(sim *core.Simulation) bool {
	canCastDrw := dk.Talents.DancingRuneWeapon && dk.DancingRuneWeapon != nil && (dk.DancingRuneWeapon.IsReady(sim) || dk.DancingRuneWeapon.CD.TimeToReady(sim) < 5*time.Second)
	currentRP := dk.CurrentRunicPower()
	return (!canCastDrw && currentRP >= 65) || (canCastDrw && dk.CurrentRunicPower() >= 100)
}

func (dk *DpsDeathknight) blBloodTapCheck(sim *core.Simulation, target *core.Unit) bool {
	if dk.CurrentBloodRunes() > 0 {
		return false
	}

	if (!dk.Talents.DancingRuneWeapon || dk.RuneWeapon.IsEnabled()) && dk.BloodTap.IsReady(sim) {
		return dk.BloodTap.Cast(sim, target)
	}

	return false
}

// Combined checks for casting gargoyle sequence & going back to blood presence after
func (dk *DpsDeathknight) blDrwCheck(sim *core.Simulation, target *core.Unit, castTime time.Duration) bool {
	if dk.blDrwCanCast(sim, castTime) {

		dk.br.activatingDrw = true
		dk.br.drwSnapshot.ActivateMajorCooldowns(sim)
		dk.br.activatingDrw = false

		if dk.DancingRuneWeapon.Cast(sim, target) {
			dk.br.drwSnapshot.ResetProcTrackers()
			dk.br.drwMaxDelay = -1
		}
		return true
	}

	return false
}

func (dk *DpsDeathknight) blDrwCanCast(sim *core.Simulation, castTime time.Duration) bool {
	if !dk.Talents.DancingRuneWeapon {
		return false
	}
	if !dk.Rotation.UseDancingRuneWeapon {
		return false
	}
	if !dk.DancingRuneWeapon.IsReady(sim) {
		return false
	}
	if !dk.CastCostPossible(sim, 60.0, 0, 0, 0) {
		return false
	}
	// Setup max delay possible
	if dk.br.drwMaxDelay == -1 {
		drwCd := dk.DancingRuneWeapon.CD.Duration
		timeLeft := sim.GetRemainingDuration()
		for timeLeft > drwCd {
			timeLeft = timeLeft - (drwCd + 2*time.Second)
		}
		dk.br.drwMaxDelay = timeLeft - 2*time.Second
	}
	// Cast it if holding will result in less total DRWs for the encounter
	if sim.CurrentTime > dk.br.drwMaxDelay {
		return true
	}
	// Cast it if holding will take from its duration
	if sim.GetRemainingDuration() < 20*time.Second {
		return true
	}
	// Make sure we can instantly put diseases up with the rune weapon
	if !dk.sr.hasGod && (dk.CurrentFrostRunes() < 1 || dk.CurrentUnholyRunes() < 1) {
		return false
	}
	if dk.sr.hasGod && dk.CurrentBloodRunes() < 1 {
		return false
	}
	if !dk.br.drwSnapshot.CanSnapShot(sim, castTime) {
		return false
	}

	return true
}
