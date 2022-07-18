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

var recastedFF = false
var recastedBP = false

func (deathKnight *DeathKnight) shouldSpreadDisease() bool {
	return recastedFF && recastedBP && deathKnight.Env.GetNumTargets() > 1
}

func (deathKnight *DeathKnight) spreadDisease(sim *core.Simulation, target *core.Unit) {
	deathKnight.Pestilence.Cast(sim, target)
	recastedFF = false
	recastedBP = false
}

func (deathKnight *DeathKnight) tryUseGCD(sim *core.Simulation) {
	//var spell *core.Spell
	var target = deathKnight.CurrentTarget

	if deathKnight.GCD.IsReady(sim) {
		// UH DK rota
		if deathKnight.Talents.SummonGargoyle {
			if (!deathKnight.TargetHasDisease(FrostFeverAuraLabel, target) || deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim) < 6*time.Second) && deathKnight.CanIcyTouch(sim) {
				deathKnight.IcyTouch.Cast(sim, target)
				recastedFF = true
			} else if (!deathKnight.TargetHasDisease(BloodPlagueAuraLabel, target) || deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim) < 6*time.Second) && deathKnight.CanPlagueStrike(sim) {
				deathKnight.PlagueStrike.Cast(sim, target)
				recastedBP = true
			} else {
				if deathKnight.PresenceMatches(UnholyPresence) && !deathKnight.SummonGargoyle.CD.IsReady(sim) && deathKnight.CanBloodPresence(sim) {
					// Swap to blood after gargoyle
					deathKnight.BloodPressence.Cast(sim, target)
					deathKnight.WaitUntil(sim, sim.CurrentTime+1)
				} else if deathKnight.Talents.Desolation > 0 && !deathKnight.DesolationAura.IsActive() && deathKnight.CanBloodStrike(sim) && !deathKnight.shouldWaitForDnD(sim, true, false, false) {
					// Desolation check
					if deathKnight.shouldSpreadDisease() {
						deathKnight.spreadDisease(sim, target)
					} else {
						deathKnight.BloodStrike.Cast(sim, target)
					}
				} else {
					if deathKnight.Rotation.UseDeathAndDecay {
						// DW Rota
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
							if deathKnight.shouldSpreadDisease() {
								deathKnight.spreadDisease(sim, target)
							} else if deathKnight.Env.GetNumTargets() > 2 {
								deathKnight.BloodBoil.Cast(sim, target)
							} else {
								deathKnight.BloodStrike.Cast(sim, target)
							}
						} else if deathKnight.CanDeathCoil(sim) && !deathKnight.SummonGargoyle.IsReady(sim) {
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
					} else {
						// No DnD Rota
						if deathKnight.CanGhoulFrenzy(sim) && deathKnight.Talents.MasterOfGhouls && (!deathKnight.Ghoul.GhoulFrenzyAura.IsActive() || deathKnight.Ghoul.GhoulFrenzyAura.RemainingDuration(sim) < 6*time.Second) {
							deathKnight.GhoulFrenzy.Cast(sim, target)
						} else if deathKnight.CanScourgeStrike(sim) {
							deathKnight.ScourgeStrike.Cast(sim, target)
						} else if deathKnight.CanBloodStrike(sim) {
							if deathKnight.shouldSpreadDisease() {
								deathKnight.spreadDisease(sim, target)
							} else if deathKnight.Env.GetNumTargets() > 2 {
								deathKnight.BloodBoil.Cast(sim, target)
							} else {
								deathKnight.BloodStrike.Cast(sim, target)
							}
						} else if deathKnight.CanDeathCoil(sim) && !deathKnight.SummonGargoyle.IsReady(sim) {
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
		}

		// Frost DK rota
		if deathKnight.Talents.HowlingBlast {
			if deathKnight.ShouldHornOfWinter(sim) {
				deathKnight.HornOfWinter.Cast(sim, target)
			} else if (!deathKnight.TargetHasDisease(FrostFeverAuraLabel, target) || deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim) < 6*time.Second) && deathKnight.CanIcyTouch(sim) {
				deathKnight.IcyTouch.Cast(sim, target)
			} else if (!deathKnight.TargetHasDisease(BloodPlagueAuraLabel, target) || deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim) < 6*time.Second) && deathKnight.CanPlagueStrike(sim) {
				deathKnight.PlagueStrike.Cast(sim, target)
			} else {
				if deathKnight.CanBloodTap(sim) && deathKnight.AllDiseasesAreActive(target) {
					deathKnight.BloodTap.Cast(sim, target)
				} else if deathKnight.CanUnbreakableArmor(sim) && deathKnight.AllDiseasesAreActive(target) {
					deathKnight.UnbreakableArmor.Cast(sim, target)
				} else if deathKnight.CanObliterate(sim) && deathKnight.AllDiseasesAreActive(target) {
					deathKnight.Obliterate.Cast(sim, target)
				} else if deathKnight.CanHowlingBlast(sim) && deathKnight.AllDiseasesAreActive(target) {
					deathKnight.HowlingBlast.Cast(sim, target)
				} else if deathKnight.CanFrostStrike(sim) && deathKnight.AllDiseasesAreActive(target) {
					deathKnight.FrostStrike.Cast(sim, target)
				} else if deathKnight.CanBloodStrike(sim) && deathKnight.AllDiseasesAreActive(target) {
					deathKnight.BloodStrike.Cast(sim, target)
				} else if deathKnight.CanIcyTouch(sim) {
					deathKnight.IcyTouch.Cast(sim, target)
				} else if deathKnight.CanPlagueStrike(sim) {
					deathKnight.PlagueStrike.Cast(sim, target)
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
