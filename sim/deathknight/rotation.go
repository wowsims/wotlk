package deathknight

import (
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
		// Disease check
		if (!deathKnight.FrostFeverDisease.IsActive() || deathKnight.FrostFeverDisease.RemainingDuration(sim) < 6*time.Second) && deathKnight.CanIcyTouch(sim) {
			deathKnight.IcyTouch.Cast(sim, target)
		} else if (!deathKnight.BloodPlagueDisease.IsActive() || deathKnight.BloodPlagueDisease.RemainingDuration(sim) < 6*time.Second) && deathKnight.CanPlagueStrike(sim) {
			deathKnight.PlagueStrike.Cast(sim, target)
		} else {
			// Desolation check
			if deathKnight.Talents.Desolation > 0 && !deathKnight.DesolationAura.IsActive() {
				if deathKnight.CanBloodStrike(sim) {
					deathKnight.BloodStrike.Cast(sim, target)
				}
			} else {
				// Unholy checks
				if deathKnight.Talents.ScourgeStrike {
					if deathKnight.CanDeathAndDecay(sim) && deathKnight.AllDiseasesAreActive() {
						deathKnight.DeathAndDecay.Cast(sim, target)
					} else if deathKnight.CanScourgeStrike(sim) && !(deathKnight.DeathAndDecay.CD.IsReady(sim) || deathKnight.DeathAndDecay.CD.TimeToReady(sim) < 6*time.Second) {
						deathKnight.ScourgeStrike.Cast(sim, target)
					} else if deathKnight.CanBloodStrike(sim) && !(deathKnight.DeathAndDecay.CD.IsReady(sim) || deathKnight.DeathAndDecay.CD.TimeToReady(sim) < 3*time.Second) {
						deathKnight.BloodStrike.Cast(sim, target)
					} else if deathKnight.CanDeathCoil(sim) {
						deathKnight.DeathCoil.Cast(sim, target)
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
	}
}
