package tank

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *TankDeathknight) TankRA_Hps(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if !dk.GCD.IsReady(sim) {
		return dk.NextGCDAt()
	}

	if dk.DoDefensiveCds(sim, target, s) {
		return -1
	}

	if dk.DoDiseaseChecks(sim, target, s) {
		return -1
	}

	t := sim.CurrentTime
	recast := 3 * time.Second // 2 GCDs for miss
	ff := dk.FrostFeverSpell.Dot(target).ExpiresAt() - t
	bp := dk.BloodPlagueSpell.Dot(target).ExpiresAt() - t
	fd := dk.CurrentFrostRunes() + dk.CurrentDeathRunes()
	ud := dk.CurrentUnholyRunes() + dk.CurrentDeathRunes()
	b, _, _ := dk.NormalCurrentRunes()

	if fd > 0 && ud > 0 && dk.DeathStrike.CanCast(sim, target) && dk.CurrentHealthPercent() < 0.75 {
		dk.DeathStrike.Cast(sim, target)
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

	if dk.Talents.FrostStrike && dk.CurrentRunicPower() > 60 && dk.FrostStrike.CanCast(sim, target) {
		dk.FrostStrike.Cast(sim, target)
		return -1
	}

	if b >= 1 {
		if dk.NormalSpentBloodRuneReadyAt(sim)-t < ff-recast && dk.NormalSpentBloodRuneReadyAt(sim)-t < bp-recast {
			dk.BloodSpell.Cast(sim, target)
		}
	}

	return -1
}
