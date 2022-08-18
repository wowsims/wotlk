package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

type UnholyRotation struct {
	dk *DpsDeathknight

	ffFirst bool
	syncFF  bool

	syncTimeFF time.Duration

	recastedFF bool
	recastedBP bool

	procTrackers []*ProcTracker
}

func (ur *UnholyRotation) Reset(sim *core.Simulation) {
	ur.syncFF = false

	ur.syncTimeFF = 0

	ur.recastedFF = false
	ur.recastedBP = false

	ur.resetProcTrackers()
}

func (dk *DpsDeathknight) getFirstDiseaseAction() deathknight.RotationAction {
	if dk.ur.ffFirst {
		return dk.RotationActionCallback_IT
	}
	return dk.RotationActionCallback_PS
}

func (dk *DpsDeathknight) getSecondDiseaseAction() deathknight.RotationAction {
	if dk.ur.ffFirst {
		return dk.RotationActionCallback_PS
	}
	return dk.RotationActionCallback_IT
}

func (dk *DpsDeathknight) getBloodRuneAction(isFirst bool) deathknight.RotationAction {
	if isFirst {
		if dk.Env.GetNumTargets() > 1 {
			return dk.RotationActionCallback_Pesti
		} else {
			return dk.RotationActionCallback_BS
		}
	} else {
		return dk.RotationActionCallback_BS
	}
}

func (dk *DpsDeathknight) desolationAuraCheck(sim *core.Simulation) bool {
	return !dk.DesolationAura.IsActive() || dk.DesolationAura.RemainingDuration(sim) < 10*time.Second ||
		dk.Rotation.BloodRuneFiller == proto.Deathknight_Rotation_BloodStrike
}

func (dk *DpsDeathknight) uhDiseaseCheck(sim *core.Simulation, target *core.Unit, spell *deathknight.RuneSpell, costRunes bool, casts int) bool {
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
	if spell.CanCast(sim) && costRunes {
		ffExpiresAt := ffRemaining + sim.CurrentTime
		bpExpiresAt := bpRemaining + sim.CurrentTime

		crpb := dk.CopyRunicPowerBar()
		spellCost := crpb.OptimalRuneCost(core.RuneCost(spell.DefaultCast.Cost))

		crpb.SpendRuneCost(sim, spell.Spell, spellCost)

		afterCastTime := sim.CurrentTime + castGcd
		currentFrostRunes := crpb.CurrentFrostRunes()
		currentUnholyRunes := crpb.CurrentUnholyRunes()
		nextFrostRuneAt := crpb.FrostRuneReadyAt(sim)
		nextUnholyRuneAt := crpb.UnholyRuneReadyAt(sim)

		// If FF is gonna drop while our runes are on CD
		if dk.uhRecastAvailableCheck(ffExpiresAt-dk.ur.syncTimeFF, afterCastTime, int(spellCost.Frost()), currentFrostRunes, nextFrostRuneAt) {
			return false
		}

		// If BP is gonna drop while our runes are on CD
		if dk.uhRecastAvailableCheck(bpExpiresAt, afterCastTime, int(spellCost.Unholy()), currentUnholyRunes, nextUnholyRuneAt) {
			return false
		}
	}

	return true
}

func (dk *DpsDeathknight) uhRecastAvailableCheck(expiresAt time.Duration, afterCastTime time.Duration,
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
		casted := dk.Pestilence.Cast(sim, target)
		landed := dk.LastOutcome.Matches(core.OutcomeLanded)

		// Reset flags on succesfull cast
		dk.ur.recastedFF = !(casted && landed)
		dk.ur.recastedBP = !(casted && landed)
		return casted
	} else {
		dk.recastDiseasesSequence(sim)
		return true
	}
}

// Simpler but somehow more effective for overall dps dnd check
func (dk *DpsDeathknight) uhShouldWaitForDnD(sim *core.Simulation, blood bool, frost bool, unholy bool) bool {
	return !(!(dk.DeathAndDecay.CD.IsReady(sim) || dk.DeathAndDecay.CD.TimeToReady(sim) <= 4*time.Second) || ((!blood || dk.CurrentBloodRunes() > 1) && (!frost || dk.CurrentFrostRunes() > 1) && (!unholy || dk.CurrentUnholyRunes() > 1)))
}

func (dk *DpsDeathknight) uhGhoulFrenzyCheck(sim *core.Simulation, target *core.Unit) bool {
	// If no Ghoul Frenzy Aura or duration less then 10 seconds we try recasting
	if !dk.GhoulFrenzyAura.IsActive() || dk.GhoulFrenzyAura.RemainingDuration(sim) < 10*time.Second {
		if (dk.Rotation.BloodTap == proto.Deathknight_Rotation_GhoulFrenzy || dk.Rotation.BtGhoulFrenzy) && dk.BloodTap.CanCast(sim) && dk.GhoulFrenzy.IsReady(sim) && dk.AllBloodRunesSpent() && dk.AllUnholySpent() && dk.SummonGargoyle.CD.TimeToReady(sim) > time.Second*55 {
			// Use Ghoul Frenzy with a Blood Tap and Blood rune if all blood runes are on CD and Garg wont come off cd in less then a minute.
			// The gargoyle check is there because you should BT -> UP -> Garg (Not in the sim yet)
			if dk.uhDiseaseCheck(sim, target, dk.GhoulFrenzy, true, 1) {
				dk.ghoulFrenzySequence(sim, true)
				return true
			} else {
				dk.recastDiseasesSequence(sim)
				return true
			}
		} else if !dk.Rotation.BtGhoulFrenzy && dk.GhoulFrenzy.CanCast(sim) && dk.IcyTouch.CanCast(sim) {
			if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()*2+50*time.Millisecond) {
				dk.afterGargoyleSequence(sim)
				return true
			}
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

func (dk *DpsDeathknight) uhBloodTap(sim *core.Simulation, target *core.Unit) bool {
	if !dk.GCD.IsReady(sim) || dk.SummonGargoyle.IsReady(sim) {
		return false
	}

	if dk.Rotation.BloodTap != proto.Deathknight_Rotation_GhoulFrenzy && dk.BloodTap.IsReady(sim) && dk.CurrentBloodRunes() == 0 {
		switch dk.Rotation.BloodTap {
		case proto.Deathknight_Rotation_IcyTouch:
			if dk.CurrentFrostRunes() == 0 {
				dk.BloodTap.Cast(sim, dk.CurrentTarget)
				dk.IcyTouch.Cast(sim, target)
				return true
			}
		case proto.Deathknight_Rotation_BloodStrikeBT:
			dk.BloodTap.Cast(sim, dk.CurrentTarget)
			dk.BloodStrike.Cast(sim, target)
			return true
		case proto.Deathknight_Rotation_BloodBoilBT:
			dk.BloodTap.Cast(sim, dk.CurrentTarget)
			dk.BloodBoil.Cast(sim, target)
			return true
		}
	}

	return false
}

func (dk *DpsDeathknight) uhEmpoweredRuneWeapon(sim *core.Simulation, target *core.Unit) bool {
	if !dk.Rotation.UseEmpowerRuneWeapon || dk.SummonGargoyle.IsReady(sim) {
		return false
	}

	if dk.EmpowerRuneWeapon.IsReady(sim) && dk.CurrentBloodRunes() == 0 && dk.CurrentFrostRunes() == 0 && dk.CurrentUnholyRunes() == 0 {
		dk.EmpowerRuneWeapon.Cast(sim, target)
		return true
	}
	return false
}

// Save up Runic Power for Summon Gargoyle - Allow casts above 100 rp or garg CD > 5 sec
func (dk *DpsDeathknight) uhDeathCoilCheck(sim *core.Simulation) bool {
	return !(dk.SummonGargoyle.IsReady(sim) || dk.SummonGargoyle.CD.TimeToReady(sim) < 5*time.Second) || dk.CurrentRunicPower() >= 100
}

// Combined checks for casting gargoyle sequence & going back to blood presence after
func (dk *DpsDeathknight) uhGargoyleCheck(sim *core.Simulation, target *core.Unit, castTime time.Duration) bool {
	if dk.uhGargoyleCanCast(sim, castTime) {
		if !dk.PresenceMatches(deathknight.UnholyPresence) {
			dk.BloodTap.Cast(sim, dk.CurrentTarget)
			dk.UnholyPresence.Cast(sim, dk.CurrentTarget)
		}

		if dk.SummonGargoyle.Cast(sim, target) {
			dk.ur.resetProcTrackers()
			return true
		}
	}

	// Go back to Blood Presence after gargoyle cast
	if dk.PresenceMatches(deathknight.UnholyPresence) && !dk.SummonGargoyle.IsReady(sim) {
		if dk.BloodTapAura.IsActive() {
			dk.BloodTapAura.Deactivate(sim)
		}
		if dk.BloodPresence.Cast(sim, target) {
			return true
		}
	}
	return false
}

func (dk *DpsDeathknight) uhGargoyleCanCast(sim *core.Simulation, castTime time.Duration) bool {
	if !dk.SummonGargoyle.IsReady(sim) {
		return false
	}
	if !dk.CastCostPossible(sim, 60.0, 0, 0, 0) {
		return false
	}
	if !dk.PresenceMatches(deathknight.UnholyPresence) && !dk.BloodTap.CanCast(sim) {
		return false
	}
	if dk.GargoyleProcCheck(sim, castTime) {
		return false
	}

	return true
}

func (dk *DpsDeathknight) setupGargoyleCooldowns() {
	// hyperspeed accelerators
	dk.gargoyleCooldownSync(core.ActionID{SpellID: 54758}, false)

	// berserking (troll)
	dk.gargoyleCooldownSync(core.ActionID{SpellID: 26297}, false)

	// potion of speed
	dk.gargoyleCooldownSync(core.ActionID{ItemID: 40211}, true)
}

func (dk *DpsDeathknight) gargoyleCooldownSync(actionID core.ActionID, isPotion bool) {
	if dk.Character.HasMajorCooldown(actionID) {
		majorCd := dk.Character.GetMajorCooldown(actionID)
		majorCd.ShouldActivate = func(sim *core.Simulation, character *core.Character) bool {
			return dk.SummonGargoyle.CD.IsReady(sim) || (dk.SummonGargoyle.CD.TimeToReady(sim) > majorCd.Spell.CD.Duration && !isPotion) || dk.SummonGargoyle.CD.ReadyAt() > dk.Env.Encounter.Duration
		}
	}
}

func logMessage(sim *core.Simulation, message string) {
	if sim.Log != nil {
		sim.Log(message)
	}
}

func (dk *DpsDeathknight) GargoyleProcCheck(sim *core.Simulation, castTime time.Duration) bool {
	for _, procTracker := range dk.ur.procTrackers {
		if !procTracker.didActivate && procTracker.aura.IsActive() {
			procTracker.didActivate = true
			procTracker.expiresAt = procTracker.aura.ExpiresAt()
		}

		// A proc is about to drop
		if procTracker.didActivate && procTracker.expiresAt <= sim.CurrentTime+castTime {
			logMessage(sim, "Proc dropping "+procTracker.aura.Label)
			return false
		}
	}

	for _, procTracker := range dk.ur.procTrackers {
		if !procTracker.didActivate {
			// logMessage(sim, "Waiting on procs..")
			return true
		}
	}

	return false
}
