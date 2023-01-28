package tank

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *TankDeathknight) DoDefensiveCds(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	if dk.CurrentHealthPercent() < dk.HpPercentForDefensives {
		// Roll out Defensives in Prio Order
		if dk.Talents.RuneTap {
			if !dk.RuneTap.CanCast(sim, target) && dk.BloodTap.CanCast(sim, nil) {
				dk.BloodTap.Cast(sim, nil)
			}
			if dk.RuneTap.CanCast(sim, target) {
				dk.RuneTap.Cast(sim, target)
				return true
			}
		}

		if dk.DeathPact.CanCast(sim, target) {
			dk.DeathPact.Cast(sim, target)
			return true
		}

		if dk.AntiMagicShell.CanCast(sim, target) {
			dk.AntiMagicShell.Cast(sim, target)
			return true
		}

		if dk.Talents.VampiricBlood {
			if !dk.VampiricBlood.CanCast(sim, target) && dk.BloodTap.CanCast(sim, nil) {
				dk.BloodTap.Cast(sim, nil)
			}
			if dk.VampiricBlood.CanCast(sim, target) {
				dk.VampiricBlood.Cast(sim, target)
				return true
			}
		}

		if dk.Talents.UnbreakableArmor {
			if !dk.UnbreakableArmor.CanCast(sim, target) && dk.BloodTap.CanCast(sim, nil) {
				dk.BloodTap.Cast(sim, nil)
			}
			if dk.UnbreakableArmor.CanCast(sim, target) {
				dk.UnbreakableArmor.Cast(sim, target)
				return true
			}
		}

		if dk.IceboundFortitude.CanCast(sim, target) {
			dk.IceboundFortitude.Cast(sim, target)
			return true
		}
	}

	if dk.CurrentHealthPercent() < 0.5 && dk.RaiseDead.IsReady(sim) {
		dk.RaiseDead.Cast(sim, nil)
		return true
	}

	return false
}

func (dk *TankDeathknight) DoDiseaseChecks(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	t := sim.CurrentTime
	recast := 3 * time.Second // 2 GCDs for miss
	ff := dk.FrostFeverSpell.Dot(target).ExpiresAt() - t
	bp := dk.BloodPlagueSpell.Dot(target).ExpiresAt() - t

	if ff <= 0 && dk.IcyTouch.CanCast(sim, target) {
		dk.IcyTouch.Cast(sim, target)
		return true
	}

	if bp <= 0 && dk.PlagueStrike.CanCast(sim, target) {
		dk.PlagueStrike.Cast(sim, target)
		return true
	}

	if ff <= recast || bp <= recast && dk.Pestilence.CanCast(sim, target) {
		dk.Pestilence.Cast(sim, target)
		return true
	}

	return false
}
