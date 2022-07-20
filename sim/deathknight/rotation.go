package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
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

func (deathKnight *DeathKnight) shouldWaitForDnD(sim *core.Simulation, blood bool, frost bool, unholy bool) bool {
	return deathKnight.Rotation.UseDeathAndDecay && !(deathKnight.Talents.Morbidity == 0 || !(deathKnight.DeathAndDecay.CD.IsReady(sim) || deathKnight.DeathAndDecay.CD.TimeToReady(sim) < 4*time.Second) || ((!blood || deathKnight.CurrentBloodRunes() > 1) && (!frost || deathKnight.CurrentFrostRunes() > 1) && (!unholy || deathKnight.CurrentUnholyRunes() > 1)))
}

var recastedFF = false
var recastedBP = false

func (deathKnight *DeathKnight) shouldSpreadDisease(sim *core.Simulation) bool {
	return recastedFF && recastedBP && deathKnight.Env.GetNumTargets() > 1
}

func (deathKnight *DeathKnight) spreadDiseases(sim *core.Simulation, target *core.Unit) {
	deathKnight.Pestilence.Cast(sim, target)
	recastedFF = false
	recastedBP = false
}

func (deathKnight *DeathKnight) tryUseGCD(sim *core.Simulation) {
	//var spell *core.Spell
	var target = deathKnight.CurrentTarget

	if deathKnight.GCD.IsReady(sim) {
		// if sim.CurrentTime < time.Millisecond {
		// 	deathKnight.WaitUntil(sim, sim.CurrentTime+time.Millisecond)
		// 	return
		// }
		// UH DK rota
		if deathKnight.Talents.SummonGargoyle {
			if deathKnight.CanRaiseDead(sim) {
				deathKnight.RaiseDead.Cast(sim, target)
				return
			}
			diseaseRefreshDuration := time.Duration(deathKnight.Rotation.DiseaseRefreshDuration) * time.Second
			// Horn of Winter if you're the DK to refresh it and its not precasted/active
			if deathKnight.ShouldHornOfWinter(sim) {
				deathKnight.HornOfWinter.Cast(sim, target)
			} else if (!deathKnight.TargetHasDisease(FrostFeverAuraLabel, target) || deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim) < diseaseRefreshDuration) && deathKnight.CanIcyTouch(sim) {
				// Dont clip if theres half a second left to tick
				remainingDuration := deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim)
				if remainingDuration < time.Millisecond*500 && remainingDuration > 0 {
					deathKnight.WaitUntil(sim, sim.CurrentTime+remainingDuration+1)
				} else {
					deathKnight.IcyTouch.Cast(sim, target)
					recastedFF = true
				}
			} else if (!deathKnight.TargetHasDisease(BloodPlagueAuraLabel, target) || deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim) < diseaseRefreshDuration) && deathKnight.CanPlagueStrike(sim) {
				// Dont clip if theres half a second left to tick
				remainingDuration := deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim)
				if remainingDuration < time.Millisecond*500 && remainingDuration > 0 {
					deathKnight.WaitUntil(sim, sim.CurrentTime+remainingDuration+1)
				} else {
					deathKnight.PlagueStrike.Cast(sim, target)
					recastedBP = true
				}
			} else {
				if deathKnight.PresenceMatches(UnholyPresence) && (deathKnight.Rotation.ArmyOfTheDead != proto.DeathKnight_Rotation_AsMajorCd || !deathKnight.ArmyOfTheDead.CD.IsReady(sim)) && !deathKnight.SummonGargoyle.CD.IsReady(sim) && deathKnight.CanBloodPresence(sim) {
					// Swap to blood presence after gargoyle cast
					deathKnight.BloodPressence.Cast(sim, target)
					deathKnight.WaitUntil(sim, sim.CurrentTime+1)
				} else if deathKnight.Talents.Desolation > 0 && !deathKnight.DesolationAura.IsActive() && deathKnight.CanBloodStrike(sim) && !deathKnight.shouldWaitForDnD(sim, true, false, false) {
					// Desolation and Pestilence check
					if deathKnight.shouldSpreadDisease(sim) {
						deathKnight.spreadDiseases(sim, target)
					} else {
						deathKnight.BloodStrike.Cast(sim, target)
					}
				} else {
					if deathKnight.Rotation.UseDeathAndDecay {
						// Death and Decay Rotation
						if deathKnight.CanDeathAndDecay(sim) && deathKnight.AllDiseasesAreActive(target) {
							deathKnight.DeathAndDecay.Cast(sim, target)
						} else if deathKnight.CanGhoulFrenzy(sim) && deathKnight.Talents.MasterOfGhouls && (!deathKnight.Ghoul.GhoulFrenzyAura.IsActive() || deathKnight.Ghoul.GhoulFrenzyAura.RemainingDuration(sim) < 6*time.Second) && !deathKnight.shouldWaitForDnD(sim, false, false, true) {
							deathKnight.GhoulFrenzy.Cast(sim, target)
						} else if deathKnight.CanScourgeStrike(sim) && !deathKnight.shouldWaitForDnD(sim, false, true, true) {
							deathKnight.ScourgeStrike.Cast(sim, target)
						} else if !deathKnight.Talents.ScourgeStrike && deathKnight.CanIcyTouch(sim) && !deathKnight.shouldWaitForDnD(sim, false, true, false) {
							deathKnight.IcyTouch.Cast(sim, target)
						} else if !deathKnight.Talents.ScourgeStrike && deathKnight.CanPlagueStrike(sim) && !deathKnight.shouldWaitForDnD(sim, false, false, true) {
							deathKnight.PlagueStrike.Cast(sim, target)
						} else if deathKnight.CanBloodStrike(sim) && !deathKnight.shouldWaitForDnD(sim, true, false, false) {
							if deathKnight.shouldSpreadDisease(sim) {
								deathKnight.spreadDiseases(sim, target)
							} else if deathKnight.Env.GetNumTargets() > 2 {
								deathKnight.BloodBoil.Cast(sim, target)
							} else {
								deathKnight.BloodStrike.Cast(sim, target)
							}
						} else if deathKnight.CanDeathCoil(sim) && !deathKnight.SummonGargoyle.IsReady(sim) {
							deathKnight.DeathCoil.Cast(sim, target)
						} else if deathKnight.CanHornOfWinter(sim) {
							deathKnight.HornOfWinter.Cast(sim, target)
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
					} else {
						// Scourge Strike Rotation
						if deathKnight.CanGhoulFrenzy(sim) && deathKnight.Talents.MasterOfGhouls && (!deathKnight.Ghoul.GhoulFrenzyAura.IsActive() || deathKnight.Ghoul.GhoulFrenzyAura.RemainingDuration(sim) < 6*time.Second) {
							deathKnight.GhoulFrenzy.Cast(sim, target)
						} else if deathKnight.CanScourgeStrike(sim) {
							deathKnight.ScourgeStrike.Cast(sim, target)
						} else if deathKnight.CanBloodStrike(sim) {
							if deathKnight.shouldSpreadDisease(sim) {
								deathKnight.spreadDiseases(sim, target)
							} else if deathKnight.Env.GetNumTargets() > 2 {
								deathKnight.BloodBoil.Cast(sim, target)
							} else {
								deathKnight.BloodStrike.Cast(sim, target)
							}
						} else if deathKnight.CanDeathCoil(sim) && !deathKnight.SummonGargoyle.IsReady(sim) {
							deathKnight.DeathCoil.Cast(sim, target)
						} else if deathKnight.CanHornOfWinter(sim) {
							deathKnight.HornOfWinter.Cast(sim, target)
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

		// Start proper Frost rotation

		// Frost DK rota
		deathKnight.doDKRotation(sim, true)
	}
}
