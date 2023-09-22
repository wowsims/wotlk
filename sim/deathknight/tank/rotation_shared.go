package tank

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *TankDeathknight) DoDiseaseChecks(sim *core.Simulation, target *core.Unit, _ *deathknight.Sequence) bool {
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

func (dk *TankDeathknight) DoFrostCast(sim *core.Simulation, target *core.Unit, _ *deathknight.Sequence) bool {
	if dk.Talents.FrostStrike && dk.FrostStrike.CanCast(sim, target) {
		dk.FrostStrike.Cast(sim, target)
		return true
	}

	if dk.Talents.HowlingBlast && dk.FreezingFogAura.IsActive() && dk.HowlingBlast.CanCast(sim, target) {
		dk.HowlingBlast.Cast(sim, target)
		return true
	}

	return false
}

func (dk *TankDeathknight) DoBloodCast(sim *core.Simulation, target *core.Unit, _ *deathknight.Sequence) bool {
	t := sim.CurrentTime
	recast := 3 * time.Second // 2 GCDs for miss
	ff := dk.FrostFeverSpell.Dot(target).ExpiresAt() - t
	bp := dk.BloodPlagueSpell.Dot(target).ExpiresAt() - t
	b := dk.CurrentBloodRunes()

	if b >= 1 {
		if dk.NormalSpentBloodRuneReadyAt(sim)-t < ff-recast && dk.NormalSpentBloodRuneReadyAt(sim)-t < bp-recast {
			dk.BloodSpell.Cast(sim, target)
			return true
		}
	}

	return false
}
