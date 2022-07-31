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

func (ur *UnholyRotation) Reset(sim *core.Simulation) {
	ur.syncFF = false

	ur.syncTimeFF = 0

	ur.recastedFF = false
	ur.recastedBP = false
}

func (dk *DpsDeathknight) desolationAuraCheck(sim *core.Simulation) bool {
	return !dk.DesolationAura.IsActive() || dk.DesolationAura.RemainingDuration(sim) < 10*time.Second || dk.Env.GetNumTargets() == 1
}

func (dk *DpsDeathknight) uhDiseaseCheck(sim *core.Simulation, target *core.Unit, spell *core.Spell, costRunes bool, casts int) bool {
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

		crpb := dk.CopyRunicPowerBar()
		runeCostForSpell := dk.RuneAmountForSpell(spell)
		spellCost := crpb.DetermineOptimalCost(sim, runeCostForSpell.Blood, runeCostForSpell.Frost, runeCostForSpell.Unholy)

		crpb.Spend(sim, spell, spellCost)

		afterCastTime := sim.CurrentTime + castGcd
		currentFrostRunes := crpb.CurrentFrostRunes()
		currentUnholyRunes := crpb.CurrentUnholyRunes()
		nextFrostRuneAt := crpb.FrostRuneReadyAt(sim)
		nextUnholyRuneAt := crpb.UnholyRuneReadyAt(sim)

		// If FF is gonna drop while our runes are on CD
		if dk.uhCheckForDiseaseRecast(ffExpiresAt-dk.ur.syncTimeFF, afterCastTime, spellCost.Frost, currentFrostRunes, nextFrostRuneAt) {
			return false
		}

		// If BP is gonna drop while our runes are on CD
		if dk.uhCheckForDiseaseRecast(bpExpiresAt, afterCastTime, spellCost.Unholy, currentUnholyRunes, nextUnholyRuneAt) {
			return false
		}
	}

	return true
}

func (dk *DpsDeathknight) uhCheckForDiseaseRecast(expiresAt time.Duration, afterCastTime time.Duration,
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

func (dk *DpsDeathknight) uhShouldSpreadDisease(sim *core.Simulation) bool {
	return dk.ur.recastedFF && dk.ur.recastedBP && dk.Env.GetNumTargets() > 1
}

func (dk *DpsDeathknight) uhSpreadDiseases(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	if dk.uhDiseaseCheck(sim, target, dk.Pestilence, true, 1) {
		casted := dk.CastPestilence(sim, target)
		landed := dk.LastCastOutcome.Matches(core.OutcomeLanded)

		// Reset flags on succesfull cast
		dk.ur.recastedFF = !(casted && landed)
		dk.ur.recastedBP = !(casted && landed)
		return casted
	} else {
		dk.recastDiseasesSequence(sim)
		return true
	}
}

func (dk *DpsDeathknight) uhShouldWaitForDnD(sim *core.Simulation, blood bool, frost bool, unholy bool) bool {
	return !(!(dk.DeathAndDecay.CD.IsReady(sim) || dk.DeathAndDecay.CD.TimeToReady(sim) <= 4*time.Second) || ((!blood || dk.CurrentBloodRunes() > 1) && (!frost || dk.CurrentFrostRunes() > 1) && (!unholy || dk.CurrentUnholyRunes() > 1)))
}

func (dk *DpsDeathknight) uhGhoulFrenzyCheck(sim *core.Simulation, target *core.Unit) bool {
	// If no Ghoul Frenzy Aura or duration less then 10 seconds we try recasting
	if !dk.GhoulFrenzyAura.IsActive() || dk.GhoulFrenzyAura.RemainingDuration(sim) < 10*time.Second {
		if dk.CanBloodTap(sim) && dk.GhoulFrenzy.IsReady(sim) && dk.AllBloodRunesSpent() && dk.AllUnholySpent() && dk.SummonGargoyle.CD.TimeToReady(sim) > time.Second*60 {
			// Use Ghoul Frenzy with a Blood Tap and Blood rune if all blood runes are on CD and Garg wont come off cd in less then a minute.
			// The gargoyle check is there because you should BT -> UP -> Garg (Not in the sim yet)
			if dk.uhDiseaseCheck(sim, target, dk.GhoulFrenzy, true, 1) {
				dk.ghoulFrenzySequence(sim, true)
				return true
			} else {
				dk.recastDiseasesSequence(sim)
				return true
			}
		} else if !dk.Rotation.BtGhoulFrenzy && dk.CanGhoulFrenzy(sim) && dk.CanIcyTouch(sim) {
			// Use Ghoul Frenzy with an Unholy Rune and sync the frost rune with Icy Touch
			if dk.uhDiseaseCheck(sim, target, dk.GhoulFrenzy, true, 5) && dk.uhDiseaseCheck(sim, target, dk.IcyTouch, true, 5) {
				dk.ghoulFrenzySequence(sim, false)
				return true
			} else {
				dk.recastDiseasesSequence(sim)
				return true
			}
		}
	}
	return false
}

func (dk *DpsDeathknight) uhGargoyleCheck(sim *core.Simulation, target *core.Unit) bool {
	if dk.uhGargoyleCanCast(sim) {
		if !dk.PresenceMatches(deathknight.UnholyPresence) {
			dk.CastBloodTap(sim, dk.CurrentTarget)
			dk.CastUnholyPresence(sim, dk.CurrentTarget)
		}

		if dk.CastSummonGargoyle(sim, target) {
			return true
		}
	}

	// Go back to Blood Presence after gargoyle cast
	if dk.PresenceMatches(deathknight.UnholyPresence) && !dk.CanSummonGargoyle(sim) {
		if dk.BloodTapAura.IsActive() {
			dk.BloodTapAura.Deactivate(sim)
		}
		if dk.CastBloodPresence(sim, target) {
			dk.WaitUntil(sim, sim.CurrentTime)
			return true
		}
	}
	return false
}

func (dk *DpsDeathknight) GargoyleProcCheck(sim *core.Simulation) bool {
	return false
}

func (dk *DpsDeathknight) uhGargoyleCanCast(sim *core.Simulation) bool {
	if dk.Opener.IsOngoing() {
		return false
	}
	if !dk.SummonGargoyle.IsReady(sim) {
		return false
	}
	if dk.CurrentRunicPower() < dk.SummonGargoyle.DefaultCast.Cost {
		return false
	}
	if !dk.PresenceMatches(deathknight.UnholyPresence) && !dk.CanBloodTap(sim) {
		return false
	}
	if dk.GargoyleProcCheck(sim) {
		return false
	}

	return true
}
