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

	// If we have spent all our runes and ERW is ready, lets use it!
	if dk.EmpowerRuneWeapon.IsReady(sim) && dk.AllRunesSpent() {
		// Also a good time to swap presence.
		if dk.Presence == deathknight.UnholyPresence && dk.BloodTap.IsReady(sim) {
			dk.BloodTap.Cast(sim, dk.CurrentTarget)
			dk.BloodPresence.Cast(sim, dk.CurrentTarget)
		}
		dk.EmpowerRuneWeapon.Cast(sim, dk.CurrentTarget)
	}

	// TODO: should we make this the default or somehow configurable?
	useBTForGF := true
	// If we need GF and we can't cast it right now, but BT is ready, lets use it!
	canGFWithBT := !dk.GhoulFrenzyAura.IsActive() && dk.BloodTap.CanCast(sim)

	if dk.Talents.Desolation > 0 && !dk.DesolationAura.IsActive() && dk.BloodStrike.CanCast(sim) {
		dk.BloodStrike.Cast(sim, dk.CurrentTarget)
	} else if dk.FrostFeverDisease[dk.CurrentTarget.Index].RemainingDuration(sim) < time.Second*4 && dk.IcyTouch.CanCast(sim) {
		dk.IcyTouch.Cast(sim, dk.CurrentTarget)
	} else if dk.BloodPlagueDisease[dk.CurrentTarget.Index].RemainingDuration(sim) < time.Second*4 && dk.PlagueStrike.CanCast(sim) {
		dk.PlagueStrike.Cast(sim, dk.CurrentTarget)
	} else if (!dk.GhoulFrenzyAura.IsActive() && dk.GhoulFrenzy.CanCast(sim)) || (useBTForGF && canGFWithBT) {
		if !dk.GhoulFrenzy.CanCast(sim) && dk.BloodTap.CanCast(sim) {
			dk.BloodTap.Cast(sim, dk.CurrentTarget)
		}
		dk.GhoulFrenzy.Cast(sim, dk.CurrentTarget)
	} else if dk.SummonGargoyle.CanCast(sim) {
		dk.SummonGargoyle.Cast(sim, dk.CurrentTarget)
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

func (dk *DpsDeathknight) RotationActionCallback_AutoDW(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {

	if !dk.GCD.IsReady(sim) {
		return dk.NextGCDAt()
	}

	// If we don't have any runes, might as well pop ERW.
	if dk.EmpowerRuneWeapon.IsReady(sim) && dk.AllRunesSpent() {
		dk.EmpowerRuneWeapon.Cast(sim, dk.CurrentTarget)
		if dk.Presence == deathknight.UnholyPresence {
			// good time to swap stances too
			dk.BloodPresence.Cast(sim, dk.CurrentTarget)
		}
	}

	useBTForGF := true
	// If we need GF and we can't cast it right now, but BT is ready, lets use it!
	canGFWithBT := !dk.GhoulFrenzyAura.IsActive() && dk.BloodTap.CanCast(sim)

	if dk.Talents.Desolation > 0 && dk.DesolationAura.RemainingDuration(sim) < time.Second && dk.BloodStrike.CanCast(sim) {
		dk.BloodStrike.Cast(sim, dk.CurrentTarget)
	} else if dk.FrostFeverDisease[dk.CurrentTarget.Index].RemainingDuration(sim) < time.Second*4 && dk.IcyTouch.CanCast(sim) {
		dk.IcyTouch.Cast(sim, dk.CurrentTarget)
	} else if dk.BloodPlagueDisease[dk.CurrentTarget.Index].RemainingDuration(sim) < time.Second*4 && dk.PlagueStrike.CanCast(sim) {
		dk.PlagueStrike.Cast(sim, dk.CurrentTarget)
	} else if (!dk.GhoulFrenzyAura.IsActive() && dk.GhoulFrenzy.CanCast(sim)) || (useBTForGF && canGFWithBT) {
		if !dk.GhoulFrenzy.CanCast(sim) && dk.BloodTap.CanCast(sim) {
			dk.BloodTap.Cast(sim, dk.CurrentTarget)
		}
		dk.GhoulFrenzy.Cast(sim, dk.CurrentTarget)
	} else if !dk.DeathAndDecayDot.IsActive() && dk.DeathAndDecay.CanCast(sim) {
		dk.DeathAndDecay.Cast(sim, dk.CurrentTarget)
	} else if dk.SummonGargoyle.CanCast(sim) {
		dk.SummonGargoyle.Cast(sim, dk.CurrentTarget)
	} else {
		dndAt := dk.DeathAndDecay.ReadyAt()

		// Prio DC if we have a lot of extra RP
		if dk.DeathCoil.CanCast(sim) && dk.CurrentRunicPower() > 80 {
			dk.DeathCoil.Cast(sim, dk.CurrentTarget)
		} else if dndAt > sim.CurrentTime {
			numDeath := dk.CurrentDeathRunes()
			numBlood := dk.CurrentBloodRunes()
			numFrost := dk.CurrentFrostRunes()
			numUnholy := dk.CurrentUnholyRunes()

			nbr := dk.SpentBloodRuneReadyAt()
			nfr := dk.SpentFrostRuneReadyAt()
			nur := dk.SpentUnholyRuneReadyAt()

			if numDeath+numBlood >= 2 || (numBlood == 1 && nbr < dndAt) {
				dk.BloodBoil.Cast(sim, dk.CurrentTarget)
			} else {
				// TODO: check which one has the shorter duration for the DOT.
				if numDeath+numFrost >= 2 || (numFrost == 1 && nfr < dndAt) {
					dk.IcyTouch.Cast(sim, dk.CurrentTarget)
				} else if numDeath+numUnholy >= 2 || (numUnholy == 1 && nur < dndAt) {
					dk.PlagueStrike.Cast(sim, dk.CurrentTarget)
				}
			}
		}

		// If we didn't do anything else spend some RP or make some more
		if dk.GCD.IsReady(sim) {
			if dk.DeathCoil.CanCast(sim) {
				dk.DeathCoil.Cast(sim, dk.CurrentTarget)
			} else if dk.HornOfWinter.CanCast(sim) && (dk.CurrentRunicPower() < 80 || !dk.HornOfWinterAura.IsActive()) {
				dk.HornOfWinter.Cast(sim, dk.CurrentTarget)
			} else {
				waitUntil := dk.RunicPowerBar.AnySpentRuneReadyAt()
				if waitUntil == sim.CurrentTime {
					waitUntil = dk.AutoAttacks.NextAttackAt()
				}
				dk.WaitUntil(sim, waitUntil)
				return 0
			}
		}
	}

	return dk.NextGCDAt()
}
