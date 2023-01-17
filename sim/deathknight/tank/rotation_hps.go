package tank

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *TankDeathknight) TankRA_Hps(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if !dk.GCD.IsReady(sim) {
		return dk.NextGCDAt()
	}

	t := sim.CurrentTime
	ff := dk.FrostFeverSpell.Dot(target).ExpiresAt() - t
	bp := dk.BloodPlagueSpell.Dot(target).ExpiresAt() - t
	fd := dk.CurrentFrostRunes() + dk.CurrentDeathRunes()
	ud := dk.CurrentUnholyRunes() + dk.CurrentDeathRunes()
	b, _, _ := dk.NormalCurrentRunes()

	if ff <= 0 && dk.IcyTouch.CanCast(sim, nil) {
		dk.IcyTouch.Cast(sim, target)
		return -1
	}

	if bp <= 0 && dk.PlagueStrike.CanCast(sim, nil) {
		dk.PlagueStrike.Cast(sim, target)
		return -1
	}

	if ff <= 2*time.Second || bp <= 2*time.Second && dk.Pestilence.CanCast(sim, nil) {
		dk.Pestilence.Cast(sim, target)
		return -1
	}

	if fd > 0 && ud > 0 && dk.DeathStrike.CanCast(sim, nil) && dk.CurrentHealthPercent() < 1.0 {
		dk.DeathStrike.Cast(sim, target)
		return -1
	}

	if dk.BloodTap.CanCast(sim, nil) {
		dk.BloodTap.Cast(sim, target)
		dk.IcyTouch.Cast(sim, target)
		dk.CancelBloodTap(sim)
		return -1
	}

	if b >= 1 && dk.NormalSpentBloodRuneReadyAt(sim)-t < ff-2*time.Second && dk.NormalSpentBloodRuneReadyAt(sim)-t < bp-2*time.Second && dk.BloodSpell.CanCast(sim, nil) {
		dk.BloodSpell.Cast(sim, target)
		return -1
	}

	return -1
}
