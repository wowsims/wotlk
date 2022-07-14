package deathknight

import (
	//"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	//"github.com/wowsims/wotlk/sim/core/proto"
	//"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	if deathKnight.GCD.IsReady(sim) {
		deathKnight.tryUseGCD(sim)
	}
}

func (deathKnight *DeathKnight) OnGCDReady(sim *core.Simulation) {
	deathKnight.tryUseGCD(sim)
}

const (
	DKRotation_Wait uint8 = iota
	DKRotation_IT
	DKRotation_PS
	DKRotation_Obli
	DKRotation_BS
	DKRotation_BT
	DKRotation_UA
	DKRotation_Pesti
	DKRotation_FS
)

func (deathKnight *DeathKnight) tryUseGCD(sim *core.Simulation) {
	//var spell *core.Spell
	var target = deathKnight.CurrentTarget

	if deathKnight.GCD.IsReady(sim) {

		// Disease has highest prio so first check them
		if !deathKnight.FrostFeverDisease.IsActive() && deathKnight.CanIcyTouch(sim) {
			deathKnight.IcyTouch.Cast(sim, target)
			return
		} else if !deathKnight.BloodPlagueDisease.IsActive() && deathKnight.CanPlagueStrike(sim) {
			deathKnight.PlagueStrike.Cast(sim, target)
			return
		}

		// Desolation check
		if deathKnight.Talents.Desolation > 0 {
			if (!deathKnight.DesolationAura.IsActive() || deathKnight.DesolationAura.RemainingDuration(sim) < 6*time.Second) && deathKnight.CanBloodStrike(sim) {
				deathKnight.BloodStrike.Cast(sim, target)
				return
			}
		}

		// Unholy checks
		if deathKnight.Talents.ScourgeStrike {
			if deathKnight.CanScourgeStrike(sim) && deathKnight.FrostFeverDisease.IsActive() && deathKnight.BloodPlagueDisease.IsActive() {
				deathKnight.ScourgeStrike.Cast(sim, target)
			} else if deathKnight.CanBloodStrike(sim) {
				deathKnight.BloodStrike.Cast(sim, target)
			} else {
				nextCD := deathKnight.IcyTouch.ReadyAt()

				if nextCD > sim.CurrentTime {
					deathKnight.WaitUntil(sim, nextCD)
				}
			}
		} else if deathKnight.Talents.HowlingBlast {
			if deathKnight.CanObliterate(sim) && deathKnight.FrostFeverDisease.IsActive() && deathKnight.BloodPlagueDisease.IsActive() {
				deathKnight.Obliterate.Cast(sim, target)
			} else if deathKnight.CanHowlingBlast(sim) && deathKnight.FrostFeverDisease.IsActive() && deathKnight.BloodPlagueDisease.IsActive() {
				deathKnight.HowlingBlast.Cast(sim, target)
			} else if deathKnight.CanFrostStrike(sim) && deathKnight.FrostFeverDisease.IsActive() && deathKnight.BloodPlagueDisease.IsActive() {
				deathKnight.FrostStrike.Cast(sim, target)
			} else if deathKnight.CanBloodStrike(sim) && deathKnight.FrostFeverDisease.IsActive() && deathKnight.BloodPlagueDisease.IsActive() {
				deathKnight.BloodStrike.Cast(sim, target)
			} else if deathKnight.CanIcyTouch(sim) {
				deathKnight.IcyTouch.Cast(sim, target)
			} else if deathKnight.CanPlagueStrike(sim) {
				deathKnight.PlagueStrike.Cast(sim, target)
			} else if deathKnight.CanBloodTap(sim) && deathKnight.FrostFeverDisease.IsActive() && deathKnight.BloodPlagueDisease.IsActive() {
				deathKnight.BloodTap.Cast(sim, target)
			} else {
				nextCD := deathKnight.IcyTouch.ReadyAt()

				if nextCD > sim.CurrentTime {
					deathKnight.WaitUntil(sim, nextCD)
				}
			}
		}
	}
}
