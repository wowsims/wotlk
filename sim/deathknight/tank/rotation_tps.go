package tank

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *TankDeathknight) TankRA_Tps(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if !dk.GCD.IsReady(sim) {
		return dk.NextGCDAt()
	}

	t := sim.CurrentTime
	ff := dk.FrostFeverSpell.Dot(target).ExpiresAt() - t
	bp := dk.BloodPlagueSpell.Dot(target).ExpiresAt() - t

	if ff <= 0 && dk.IcyTouch.CanCast(sim, target) {
		dk.IcyTouch.Cast(sim, target)
		return -1
	}

	if bp <= 0 && dk.PlagueStrike.CanCast(sim, target) {
		dk.PlagueStrike.Cast(sim, target)
		return -1
	}

	if ff <= 2*time.Second || bp <= 2*time.Second && dk.Pestilence.CanCast(sim, target) {
		dk.Pestilence.Cast(sim, target)
		return -1
	}

	if dk.switchIT && dk.IcyTouch.CanCast(sim, target) {
		dk.IcyTouch.Cast(sim, target)

		if dk.DeathRunesInFU() == 0 {
			dk.switchIT = false
		}

		return -1
	}

	if !dk.switchIT && dk.FuSpell.CanCast(sim, target) {
		dk.FuSpell.Cast(sim, target)

		if dk.DeathRunesInFU() == 4 {
			dk.switchIT = true
		}

		return -1
	}

	if dk.Rotation.BloodTapPrio == proto.TankDeathknight_Rotation_Offensive {
		if dk.BloodTap.CanCast(sim, target) {
			dk.BloodTap.Cast(sim, target)
			dk.IcyTouch.Cast(sim, target)
			dk.CancelBloodTap(sim)
			return -1
		}
	}

	if dk.DoFrostCast(sim, target, s) {
		return -1
	}

	if dk.DoBloodCast(sim, target, s) {
		return -1
	}

	return -1
}
