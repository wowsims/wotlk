package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

type UnholyRotation struct {
	lastCastSpell *core.Spell

	ffFirst bool
	syncFF  bool

	syncTimeFF time.Duration

	recastedFF bool
	recastedBP bool
}

func (ur *UnholyRotation) ResetUnholyRotation(sim *core.Simulation) {
	ur.syncFF = false

	ur.syncTimeFF = 0

	ur.recastedFF = false
	ur.recastedBP = false
}

func (dk *DpsDeathknight) shouldWaitForDnD(sim *core.Simulation, blood bool, frost bool, unholy bool) bool {
	return dk.Rotation.UseDeathAndDecay && !(dk.Talents.Morbidity == 0 || !(dk.DeathAndDecay.CD.IsReady(sim) || dk.DeathAndDecay.CD.TimeToReady(sim) < 4*time.Second) || ((!blood || dk.CurrentBloodRunes() > 1) && (!frost || dk.CurrentFrostRunes() > 1) && (!unholy || dk.CurrentUnholyRunes() > 1)))
}

func (dk *DpsDeathknight) UnholyDiseaseCheckWrapper(sim *core.Simulation, target *core.Unit, spell *core.Spell, costRunes bool, casts int) bool {
	ffRemaining := dk.FrostFeverDisease[target.Index].RemainingDuration(sim)
	bpRemaining := dk.BloodPlagueDisease[target.Index].RemainingDuration(sim)
	castGcd := dk.SpellGCD() * time.Duration(casts)

	if !dk.TargetHasDisease(deathknight.FrostFeverAuraLabel, target) || ffRemaining < castGcd {
		// Refresh FF
		return false
	}
	if !dk.TargetHasDisease(deathknight.BloodPlagueAuraLabel, target) || bpRemaining < castGcd {
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
		currentFrostRunes := crpb.CurrentFrostRunes()
		currentUnholyRunes := crpb.CurrentUnholyRunes()
		nextFrostRuneAt := crpb.FrostRuneReadyAt(sim)
		nextUnholyRuneAt := crpb.UnholyRuneReadyAt(sim)

		// Check FF
		if dk.checkForDiseaseRecast(ffExpiresAt-dk.syncTimeFF, afterCastTime, spellCost.Frost, currentFrostRunes, nextFrostRuneAt) {
			return false
		}

		// Check BP
		if dk.checkForDiseaseRecast(bpExpiresAt, afterCastTime, spellCost.Unholy, currentUnholyRunes, nextUnholyRuneAt) {
			return false
		}
	}

	return true
}

func (dk *DpsDeathknight) checkForDiseaseRecast(expiresAt time.Duration, afterCastTime time.Duration,
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

func (dk *DpsDeathknight) castClipDisease(mainDisease bool, gracePeriod time.Duration, sim *core.Simulation, canCast bool, spell *core.Spell, dot *core.Dot, target *core.Unit) bool {
	if canCast {
		// Dont drop disease due to %dmg modifiers
		if dot.TickCount < dot.NumberOfTicks-1 {
			nextTickAt := dot.ExpiresAt() - dot.TickLength*time.Duration((dot.NumberOfTicks-1)-dot.TickCount)
			if nextTickAt > sim.CurrentTime && (nextTickAt < sim.CurrentTime+gracePeriod || nextTickAt < sim.CurrentTime+400*time.Millisecond) {
				// Delay disease for next tick
				dk.LastCastOutcome = core.OutcomeMiss
				dk.WaitUntil(sim, nextTickAt+50*time.Millisecond)
				return true
			}
		}

		spell.Cast(sim, target)
		dk.lastCastSpell = spell
		success := dk.LastCastOutcome.Matches(core.OutcomeLanded)
		if success && spell == dk.IcyTouch {
			dk.syncTimeFF = 0
		}
		return true
	}
	return false
}

func (dk *DpsDeathknight) shouldSpreadDisease(sim *core.Simulation) bool {
	return dk.recastedFF && dk.recastedBP && dk.Env.GetNumTargets() > 1
}

func (dk *DpsDeathknight) spreadDiseases(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := dk.UnholyDiseaseCheckWrapper(sim, target, dk.Pestilence, true, 1)
	if casted {
		dk.CastPestilence(sim, target)
		landed := dk.LastCastOutcome.Matches(core.OutcomeLanded)

		// Reset flags on succesfull cast
		dk.recastedFF = !(casted && landed)
		dk.recastedBP = !(casted && landed)
		return casted
	} else {
		dk.recastDiseasesSequence(sim)
		return true
	}
}
