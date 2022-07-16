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

func (deathKnight *DeathKnight) shouldWaitForDnD(sim *core.Simulation, blood bool, frost bool, unholy bool) bool {
	return !(deathKnight.Talents.Morbidity == 0 || !(deathKnight.DeathAndDecay.CD.IsReady(sim) || deathKnight.DeathAndDecay.CD.TimeToReady(sim) < 6*time.Second) || ((!blood || deathKnight.CurrentBloodRunes() > 1) && (!frost || deathKnight.CurrentFrostRunes() > 1) && (!unholy || deathKnight.CurrentUnholyRunes() > 1)))
}

func (deathKnight *DeathKnight) tryUseGCD(sim *core.Simulation) {
	//var spell *core.Spell
	var target = deathKnight.CurrentTarget

	if deathKnight.GCD.IsReady(sim) {
		// UH DK rota
		if deathKnight.Talents.SummonGargoyle {
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
					if deathKnight.Talents.Morbidity > 0 && deathKnight.CanDeathAndDecay(sim) && deathKnight.AllDiseasesAreActive() {
						deathKnight.DeathAndDecay.Cast(sim, target)
					} else if deathKnight.CanGhoulFrenzy(sim) && deathKnight.Talents.MasterOfGhouls && (!deathKnight.Ghoul.GhoulFrenzyAura.IsActive() || deathKnight.Ghoul.GhoulFrenzyAura.RemainingDuration(sim) < 6*time.Second) && !deathKnight.shouldWaitForDnD(sim, false, false, true) {
						deathKnight.GhoulFrenzy.Cast(sim, target)
					} else if deathKnight.CanScourgeStrike(sim) && (deathKnight.Talents.Morbidity == 0 || !deathKnight.shouldWaitForDnD(sim, false, true, true)) {
						deathKnight.ScourgeStrike.Cast(sim, target)
					} else if !deathKnight.Talents.ScourgeStrike && deathKnight.CanIcyTouch(sim) && !deathKnight.shouldWaitForDnD(sim, false, true, false) {
						deathKnight.IcyTouch.Cast(sim, target)
					} else if !deathKnight.Talents.ScourgeStrike && deathKnight.CanPlagueStrike(sim) && !deathKnight.shouldWaitForDnD(sim, false, false, true) {
						deathKnight.PlagueStrike.Cast(sim, target)
					} else if deathKnight.CanBloodStrike(sim) && !deathKnight.shouldWaitForDnD(sim, true, false, false) {
						deathKnight.BloodStrike.Cast(sim, target)
					} else if deathKnight.CanDeathCoil(sim) {
						deathKnight.DeathCoil.Cast(sim, target)
					} else {
						if deathKnight.GCD.IsReady(sim) && !deathKnight.IsWaiting() {
							// This means we did absolutely nothing.
							// Wait until our next auto attack to decide again.
							nextSwing := deathKnight.AutoAttacks.MainhandSwingAt
							if deathKnight.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
								nextSwing = core.MinDuration(nextSwing, deathKnight.AutoAttacks.OffhandSwingAt)
							}
							deathKnight.WaitUntil(sim, nextSwing)
						}
					}
				}
			}
		}

		// Frost DK rota
		if deathKnight.Talents.HowlingBlast {
			if (!deathKnight.FrostFeverDisease.IsActive() || deathKnight.FrostFeverDisease.RemainingDuration(sim) < 6*time.Second) && deathKnight.CanIcyTouch(sim) {
				deathKnight.IcyTouch.Cast(sim, target)
			} else if (!deathKnight.BloodPlagueDisease.IsActive() || deathKnight.BloodPlagueDisease.RemainingDuration(sim) < 6*time.Second) && deathKnight.CanPlagueStrike(sim) {
				deathKnight.PlagueStrike.Cast(sim, target)
			} else {
				if deathKnight.CanObliterate(sim) && deathKnight.FrostFeverDisease.IsActive() && deathKnight.BloodPlagueDisease.IsActive() {
					deathKnight.Obliterate.Cast(sim, target)
				} else if deathKnight.CanBloodTap(sim) && deathKnight.FrostFeverDisease.IsActive() && deathKnight.BloodPlagueDisease.IsActive() {
					deathKnight.BloodTap.Cast(sim, target)
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
				} else {
					if deathKnight.GCD.IsReady(sim) && !deathKnight.IsWaiting() {
						// This means we did absolutely nothing.
						// Wait until our next auto attack to decide again.
						nextSwing := deathKnight.AutoAttacks.MainhandSwingAt
						if deathKnight.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
							nextSwing = core.MinDuration(nextSwing, deathKnight.AutoAttacks.OffhandSwingAt)
						}
						deathKnight.WaitUntil(sim, nextSwing)
					}
				}
			}
		}
	}
}
