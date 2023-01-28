package tank

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *TankDeathknight) TankRA_Tps(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if !dk.GCD.IsReady(sim) {
		return dk.NextGCDAt()
	}

	t := sim.CurrentTime
	ff := dk.FrostFeverSpell.Dot(target).ExpiresAt() - t
	bp := dk.BloodPlagueSpell.Dot(target).ExpiresAt() - t
	b, _, _ := dk.NormalCurrentRunes()

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

	if !dk.switchIT && dk.DeathStrike.CanCast(sim, target) {
		dk.DeathStrike.Cast(sim, target)

		if dk.DeathRunesInFU() == 4 {
			dk.switchIT = true
		}

		return -1
	}

	if dk.BloodTap.CanCast(sim, target) {
		dk.BloodTap.Cast(sim, target)
		dk.IcyTouch.Cast(sim, target)
		dk.CancelBloodTap(sim)
		return -1
	}

	if dk.Talents.FrostStrike && dk.CurrentRunicPower() > 60 && dk.FrostStrike.CanCast(sim, target) {
		dk.FrostStrike.Cast(sim, target)
		return -1
	}

	if dk.Talents.HowlingBlast && dk.RimeAura.IsActive() && dk.HowlingBlast.CanCast(sim, target) {
		dk.HowlingBlast.Cast(sim, target)
		return -1
	}

	if b >= 1 && dk.NormalSpentBloodRuneReadyAt(sim)-t < ff-2*time.Second && dk.NormalSpentBloodRuneReadyAt(sim)-t < bp-2*time.Second && dk.BloodSpell.CanCast(sim, nil) {
		dk.BloodSpell.Cast(sim, target)
		return -1
	}

	return -1
}
