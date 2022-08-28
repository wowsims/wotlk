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

	if dk.Talents.Desolation > 0 && !dk.DesolationAura.IsActive() && dk.BloodStrike.CanCast(sim) {
		dk.BloodStrike.Cast(sim, dk.CurrentTarget)
	} else if !dk.FrostFeverDisease[dk.CurrentTarget.Index].IsActive() && dk.IcyTouch.CanCast(sim) {
		dk.IcyTouch.Cast(sim, dk.CurrentTarget)
	} else if !dk.BloodPlagueDisease[dk.CurrentTarget.Index].IsActive() && dk.PlagueStrike.CanCast(sim) {
		dk.PlagueStrike.Cast(sim, dk.CurrentTarget)
	} else if dk.Talents.Reaping > 0 && dk.BloodStrike.CanCast(sim) {
		dk.BloodStrike.Cast(sim, dk.CurrentTarget)
	} else if dk.ScourgeStrike.CanCast(sim) {
		dk.ScourgeStrike.Cast(sim, dk.CurrentTarget)
	} else if dk.BloodStrike.CanCast(sim) {
		dk.BloodStrike.Cast(sim, dk.CurrentTarget)
	} else if dk.DeathCoil.CanCast(sim) {
		dk.DeathCoil.Cast(sim, dk.CurrentTarget)
	} else {
		// This means we dont have the resources to do anything.
		dk.WaitUntil(sim, dk.RunicPowerBar.AnySpentRuneReadyAt())
		return 0
	}

	return dk.NextGCDAt()
}
