package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) RotationActionCallback_Auto(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {

	if !dk.GCD.IsReady(sim) {
		return dk.NextGCDAt()
	}

	// 1. If you have desolation - Maintain Desolation
	// 2. Maintain Frost Fever
	// 3. Maintain Blood Plague
	// 4. If you have reaping - Spend Blood Runes (Blood Strike)
	// 5. Use extra runes on Scourge Strike
	// 6. If you don't have reaping, spend extra blood runes (Blood Strike)
	// 7. Use Runic Power on Death Coil

	// If we have spent all our runes and ERW is ready, lets use it!
	if dk.EmpowerRuneWeapon.IsReady(sim) && dk.AllRunesSpent() {
		if dk.Presence == deathknight.UnholyPresence && dk.BloodTap.IsReady(sim) {
			dk.BloodTap.Cast(sim, dk.CurrentTarget)
			dk.BloodPresence.Cast(sim, dk.CurrentTarget)
		}
		dk.EmpowerRuneWeapon.Cast(sim, dk.CurrentTarget)
	}

	if dk.Talents.Desolation > 0 && !dk.DesolationAura.IsActive() && dk.BloodStrike.CanCast(sim) {
		dk.BloodStrike.Cast(sim, dk.CurrentTarget)
	} else if dk.FrostFeverDisease[dk.CurrentTarget.Index].RemainingDuration(sim) < time.Second*4 && dk.IcyTouch.CanCast(sim) {
		dk.IcyTouch.Cast(sim, dk.CurrentTarget)
	} else if dk.BloodPlagueDisease[dk.CurrentTarget.Index].RemainingDuration(sim) < time.Second*4 && dk.PlagueStrike.CanCast(sim) {
		dk.PlagueStrike.Cast(sim, dk.CurrentTarget)
	} else if !dk.GhoulFrenzyAura.IsActive() && dk.GhoulFrenzy.CanCast(sim) {
		dk.GhoulFrenzy.Cast(sim, dk.CurrentTarget)
	} else if dk.SummonGargoyle.CanCast(sim) {
		dk.SummonGargoyle.Cast(sim, dk.CurrentTarget)
		return dk.NextGCDAt()
	} else if dk.Talents.Reaping > 0 && dk.BloodStrike.CanCast(sim) {
		dk.BloodStrike.Cast(sim, dk.CurrentTarget)
	} else if dk.ScourgeStrike.CanCast(sim) {
		dk.ScourgeStrike.Cast(sim, dk.CurrentTarget)
	} else if dk.BloodStrike.CanCast(sim) {
		dk.BloodStrike.Cast(sim, dk.CurrentTarget)
	} else if dk.DeathCoil.CanCast(sim) {
		dk.DeathCoil.Cast(sim, dk.CurrentTarget)
	} else {
		if dk.HornOfWinter.CanCast(sim) {
			dk.HornOfWinter.Cast(sim, dk.CurrentTarget)
		} else {
			// This means we dont have the resources to do anything.
			dk.WaitUntil(sim, dk.RunicPowerBar.AnySpentRuneReadyAt())
			return 0
		}
	}

	return dk.NextGCDAt()
}
