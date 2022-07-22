package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	//"github.com/wowsims/wotlk/sim/core/proto"
	//"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	if !deathKnight.onOpener {
		if deathKnight.GCD.IsReady(sim) {
			deathKnight.tryUseGCD(sim)
		}
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
	if deathKnight.GCD.IsReady(sim) {
		deathKnight.DoRotation(sim)
	}
}

func (o *Opener) DoNext(sim *core.Simulation, deathKnight *DeathKnight) bool {
	target := deathKnight.CurrentTarget
	casted := &deathKnight.castSuccessful
	advance := true
	*casted = false

	if o.idx < o.numActions {
		action := o.actions[o.idx]

		switch action {
		case OpenerAction_IT:
			*casted = deathKnight.CastIcyTouch(sim, target)
			// Add this line if you care about recasting a spell in the opener in
			// case it missed
			advance = deathKnight.LastCastOutcome != core.OutcomeMiss
		case OpenerAction_PS:
			*casted = deathKnight.CastPlagueStrike(sim, target)
			advance = deathKnight.LastCastOutcome != core.OutcomeMiss
		case OpenerAction_UA:
			*casted = deathKnight.CastUnbreakableArmor(sim, target)
			// Add this line if your spell does not incur a GCD or you will hang!
			deathKnight.WaitUntil(sim, sim.CurrentTime)
		case OpenerAction_BT:
			*casted = deathKnight.CastBloodTap(sim, target)
			deathKnight.WaitUntil(sim, sim.CurrentTime)
		case OpenerAction_Obli:
			*casted = deathKnight.CastObliterate(sim, target)
		case OpenerAction_FS:
			*casted = deathKnight.CastFrostStrike(sim, target)
		case OpenerAction_Pesti:
			*casted = deathKnight.CastPestilence(sim, target)
			if deathKnight.LastCastOutcome == core.OutcomeMiss {
				advance = false
			}
		case OpenerAction_ERW:
			*casted = deathKnight.CastEmpowerRuneWeapon(sim, target)
			deathKnight.WaitUntil(sim, sim.CurrentTime)
		case OpenerAction_HB_Ghoul_RimeCheck:
			// You can do custom actions, this is deciding whether to HB or raise dead
			if deathKnight.RimeAura.IsActive() {
				*casted = deathKnight.CastHowlingBlast(sim, target)
			} else {
				*casted = deathKnight.CastRaiseDead(sim, target)
			}
		case OpenerAction_BS:
			*casted = deathKnight.CastBloodStrike(sim, target)
		}

		// Advances the opener
		if *casted && advance {
			o.idx += 1
		}
	} else {
		deathKnight.onOpener = false

		if deathKnight.opener.id == OpenerID_FrostSubBlood_Full || deathKnight.opener.id == OpenerID_FrostSubUnholy_Full {
			if deathKnight.ShouldHornOfWinter(sim) {
				*casted = deathKnight.CastHornOfWinter(sim, target)
			} else {
				*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.Obliterate)
				if !*casted {
					if deathKnight.KillingMachineAura.IsActive() && !deathKnight.RimeAura.IsActive() {
						*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.FrostStrike)
					} else if deathKnight.KillingMachineAura.IsActive() && deathKnight.RimeAura.IsActive() {
						if deathKnight.CastCostPossible(sim, 0, 0, 1, 1) && deathKnight.CurrentRunicPower() < 110 {
							*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.HowlingBlast)
						} else if deathKnight.CastCostPossible(sim, 0, 0, 1, 1) && deathKnight.CurrentRunicPower() > 110 {
							*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.HowlingBlast)
						} else if !deathKnight.CastCostPossible(sim, 0, 0, 1, 1) && deathKnight.CurrentRunicPower() > 110 {
							*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.FrostStrike)
						} else if !deathKnight.CastCostPossible(sim, 0, 0, 1, 1) && deathKnight.CurrentRunicPower() < 110 {
							*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.FrostStrike)
						}
					} else if !deathKnight.KillingMachineAura.IsActive() && deathKnight.RimeAura.IsActive() {
						if deathKnight.CurrentRunicPower() < 110 {
							*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.HowlingBlast)
						} else {
							*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.FrostStrike)
						}
					} else {
						*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.FrostStrike)
						if !*casted {
							*casted = deathKnight.DiseaseCheckWrapper(sim, target, deathKnight.HornOfWinter)
						}
					}
				}
			}
		} else if deathKnight.opener.id == OpenerID_Unholy_Full {
			// I suggest adding the a wrapper around each spell you cast like this:
			// deathKnight.YourWrapper(sim, target, deathKnight.FrostStrike) that returns a bool for when you casted
			// since the waiting code relies on knowing if you actually casted

			if deathKnight.CanRaiseDead(sim) {
				deathKnight.RaiseDead.Cast(sim, target)
				*casted = true
				return *casted
			}
			diseaseRefreshDuration := time.Duration(deathKnight.Rotation.DiseaseRefreshDuration) * time.Second
			// Horn of Winter if you're the DK to refresh it and its not precasted/active
			if deathKnight.ShouldHornOfWinter(sim) {
				deathKnight.HornOfWinter.Cast(sim, target)
				*casted = true
			} else if (!deathKnight.TargetHasDisease(FrostFeverAuraLabel, target) || deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim) < diseaseRefreshDuration) && deathKnight.CanIcyTouch(sim) {
				// Dont clip if theres half a second left to tick
				remainingDuration := deathKnight.FrostFeverDisease[target.Index].RemainingDuration(sim)
				if remainingDuration < time.Millisecond*500 && remainingDuration > 0 {
					deathKnight.WaitUntil(sim, sim.CurrentTime+remainingDuration+1)
				} else {
					deathKnight.IcyTouch.Cast(sim, target)
					*casted = true
					recastedFF = true
				}
			} else if (!deathKnight.TargetHasDisease(BloodPlagueAuraLabel, target) || deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim) < diseaseRefreshDuration) && deathKnight.CanPlagueStrike(sim) {
				// Dont clip if theres half a second left to tick
				remainingDuration := deathKnight.BloodPlagueDisease[target.Index].RemainingDuration(sim)
				if remainingDuration < time.Millisecond*500 && remainingDuration > 0 {
					deathKnight.WaitUntil(sim, sim.CurrentTime+remainingDuration+1)
				} else {
					deathKnight.PlagueStrike.Cast(sim, target)
					*casted = true
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
						*casted = true
					} else {
						deathKnight.BloodStrike.Cast(sim, target)
						*casted = true
					}
				} else {
					if deathKnight.Rotation.UseDeathAndDecay {
						// Death and Decay Rotation
						if deathKnight.CanDeathAndDecay(sim) && deathKnight.AllDiseasesAreActive(target) {
							deathKnight.DeathAndDecay.Cast(sim, target)
							*casted = true
						} else if deathKnight.CanGhoulFrenzy(sim) && deathKnight.Talents.MasterOfGhouls && (!deathKnight.Ghoul.GhoulFrenzyAura.IsActive() || deathKnight.Ghoul.GhoulFrenzyAura.RemainingDuration(sim) < 6*time.Second) && !deathKnight.shouldWaitForDnD(sim, false, false, true) {
							deathKnight.GhoulFrenzy.Cast(sim, target)
							*casted = true
						} else if deathKnight.CanScourgeStrike(sim) && !deathKnight.shouldWaitForDnD(sim, false, true, true) {
							deathKnight.ScourgeStrike.Cast(sim, target)
							*casted = true
						} else if !deathKnight.Talents.ScourgeStrike && deathKnight.CanIcyTouch(sim) && !deathKnight.shouldWaitForDnD(sim, false, true, false) {
							deathKnight.IcyTouch.Cast(sim, target)
							*casted = true
						} else if !deathKnight.Talents.ScourgeStrike && deathKnight.CanPlagueStrike(sim) && !deathKnight.shouldWaitForDnD(sim, false, false, true) {
							deathKnight.PlagueStrike.Cast(sim, target)
							*casted = true
						} else if deathKnight.CanBloodStrike(sim) && !deathKnight.shouldWaitForDnD(sim, true, false, false) {
							if deathKnight.shouldSpreadDisease(sim) {
								deathKnight.spreadDiseases(sim, target)
								*casted = true
							} else if deathKnight.Env.GetNumTargets() > 2 {
								deathKnight.BloodBoil.Cast(sim, target)
								*casted = true
							} else {
								deathKnight.BloodStrike.Cast(sim, target)
								*casted = true
							}
						} else if deathKnight.CanDeathCoil(sim) && !deathKnight.SummonGargoyle.IsReady(sim) {
							deathKnight.DeathCoil.Cast(sim, target)
							*casted = true
						} else if deathKnight.CanHornOfWinter(sim) {
							deathKnight.HornOfWinter.Cast(sim, target)
							*casted = true
						} else {
							// Probably want to make this just return *casted as casted should be false in this case, the wait time will be handled after the return
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
							*casted = true
						} else if deathKnight.CanScourgeStrike(sim) {
							deathKnight.ScourgeStrike.Cast(sim, target)
							*casted = true
						} else if deathKnight.CanBloodStrike(sim) {
							if deathKnight.shouldSpreadDisease(sim) {
								deathKnight.spreadDiseases(sim, target)
								*casted = true
							} else if deathKnight.Env.GetNumTargets() > 2 {
								deathKnight.BloodBoil.Cast(sim, target)
								*casted = true
							} else {
								deathKnight.BloodStrike.Cast(sim, target)
								*casted = true
							}
						} else if deathKnight.CanDeathCoil(sim) && !deathKnight.SummonGargoyle.IsReady(sim) {
							deathKnight.DeathCoil.Cast(sim, target)
							*casted = true
						} else if deathKnight.CanHornOfWinter(sim) {
							deathKnight.HornOfWinter.Cast(sim, target)
							*casted = true
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
		// Other prio lists for other specs here just else if {...
	}

	return *casted
}

func (deathKnight *DeathKnight) DoRotation(sim *core.Simulation) {
	opener := deathKnight.opener
	if !opener.DoNext(sim, deathKnight) {
		if deathKnight.GCD.IsReady(sim) && !deathKnight.IsWaiting() {
			waitUntil := deathKnight.AutoAttacks.MainhandSwingAt
			if deathKnight.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
				waitUntil = core.MinDuration(waitUntil, deathKnight.AutoAttacks.OffhandSwingAt)
			}
			waitUntil = core.MinDuration(waitUntil, deathKnight.AnyRuneReadyAt(sim))
			deathKnight.WaitUntil(sim, waitUntil)
		} else { // No resources
			waitUntil := deathKnight.AnySpentRuneReadyAt(sim)
			deathKnight.WaitUntil(sim, waitUntil)
		}
	}
}

func (deathKnight *DeathKnight) ResetRotation(sim *core.Simulation) {
	deathKnight.opener.idx = 0
}
