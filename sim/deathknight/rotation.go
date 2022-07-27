package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	if !dk.Opener.IsOngoing() {
		if dk.GCD.IsReady(sim) {
			dk.tryUseGCD(sim)
		}
	}
}

func (dk *Deathknight) OnGCDReady(sim *core.Simulation) {
	dk.tryUseGCD(sim)
}

func (dk *Deathknight) tryUseGCD(sim *core.Simulation) {
	dk.DoRotation(sim)
}

func RotationActionCallback_IT(sim *core.Simulation, target *core.Unit, dk *Deathknight, s *Sequence) bool {
	casted := dk.CastIcyTouch(sim, target)
	advance := dk.LastCastOutcome != core.OutcomeMiss

	s.ConditionalAdvance(casted && advance)
	return casted
}

func RotationActionCallback_PS(sim *core.Simulation, target *core.Unit, dk *Deathknight, s *Sequence) bool {
	casted := dk.CastPlagueStrike(sim, target)
	advance := dk.LastCastOutcome != core.OutcomeMiss

	s.ConditionalAdvance(casted && advance)
	return casted
}

func RotationActionCallback_UA(sim *core.Simulation, target *core.Unit, dk *Deathknight, s *Sequence) bool {
	casted := dk.CastUnbreakableArmor(sim, target)
	dk.WaitUntil(sim, sim.CurrentTime)

	s.ConditionalAdvance(casted)
	return casted
}

func RotationActionCallback_BT(sim *core.Simulation, target *core.Unit, dk *Deathknight, s *Sequence) bool {
	casted := dk.CastBloodTap(sim, target)
	dk.WaitUntil(sim, sim.CurrentTime)

	s.ConditionalAdvance(casted)
	return casted
}

func RotationActionCallback_ERW(sim *core.Simulation, target *core.Unit, dk *Deathknight, s *Sequence) bool {
	casted := dk.CastEmpowerRuneWeapon(sim, target)
	dk.WaitUntil(sim, sim.CurrentTime)

	s.ConditionalAdvance(casted)
	return casted
}

func RotationActionCallback_Obli(sim *core.Simulation, target *core.Unit, dk *Deathknight, s *Sequence) bool {
	casted := dk.CastObliterate(sim, target)

	s.ConditionalAdvance(true)
	return casted
}

func RotationActionCallback_FS(sim *core.Simulation, target *core.Unit, dk *Deathknight, s *Sequence) bool {
	casted := dk.CastFrostStrike(sim, target)

	s.ConditionalAdvance(true)
	return casted
}

func RotationActionCallback_Pesti(sim *core.Simulation, target *core.Unit, dk *Deathknight, s *Sequence) bool {
	casted := dk.CastPestilence(sim, target)
	advance := dk.LastCastOutcome != core.OutcomeMiss

	s.ConditionalAdvance(advance)
	return casted
}

func RotationActionCallback_BS(sim *core.Simulation, target *core.Unit, dk *Deathknight, s *Sequence) bool {
	casted := dk.CastBloodStrike(sim, target)

	s.ConditionalAdvance(true)
	return casted
}

func (o *Sequence) DoAction(sim *core.Simulation, target *core.Unit, dk *Deathknight) bool {
	action := o.actions[o.idx]
	return action(sim, target, dk, o)

	/*
			minClickLatency := time.Millisecond * 0

			switch action {
			case RotationAction_IT:
				casted = dk.CastIcyTouch(sim, target)
				// Add this line if you care about recasting a spell in the opener in
				// case it missed
				advance = dk.LastCastOutcome != core.OutcomeMiss
			case RotationAction_PS:
				casted = dk.CastPlagueStrike(sim, target)
				advance = dk.LastCastOutcome.Matches(core.OutcomeHit | core.OutcomeCrit)
			case RotationAction_UA:
				casted = dk.CastUnbreakableArmor(sim, target)
				// Add this line if your spell does not incur a GCD or you will hang!
				dk.WaitUntil(sim, sim.CurrentTime+minClickLatency)
			case RotationAction_BT:
				casted = dk.CastBloodTap(sim, target)
				dk.WaitUntil(sim, sim.CurrentTime+minClickLatency)
			case RotationAction_Obli:
				casted = dk.CastObliterate(sim, target)
			case RotationAction_FS:
				casted = dk.CastFrostStrike(sim, target)
				forceAdvance = true
			case RotationAction_Pesti:
				casted = dk.CastPestilence(sim, target)
				if dk.LastCastOutcome == core.OutcomeMiss {
					advance = false
				}
			case RotationAction_ERW:
				casted = dk.CastEmpowerRuneWeapon(sim, target)
				dk.WaitUntil(sim, sim.CurrentTime+minClickLatency)
			case RotationAction_HB_Ghoul_RimeCheck:
				// You can do custom actions, this is deciding whether to HB or raise dead
				if dk.RimeAura.IsActive() {
					casted = dk.CastHowlingBlast(sim, target)
				} else {
					casted = dk.CastRaiseDead(sim, target)
				}
			case RotationAction_BS:
				casted = dk.CastBloodStrike(sim, target)
				advance = dk.LastCastOutcome != core.OutcomeMiss
			case RotationAction_SS:
				casted = dk.CastScourgeStrike(sim, target)
				advance = dk.LastCastOutcome.Matches(core.OutcomeHit | core.OutcomeCrit)
			case RotationAction_DND:
				casted = dk.CastDeathAndDecay(sim, target)
			case RotationAction_GF:
				casted = dk.CastGhoulFrenzy(sim, target)
			case RotationAction_DC:
				casted = dk.CastDeathCoil(sim, target)
			case RotationAction_Garg:
				casted = dk.CastSummonGargoyle(sim, target)
			case RotationAction_AOTD:
				casted = dk.CastArmyOfTheDead(sim, target)
			case RotationAction_BP:
				casted = dk.CastBloodPresence(sim, target)
				if !casted {
					dk.WaitUntil(sim, dk.BloodPresence.CD.ReadyAt())
				} else {
					dk.WaitUntil(sim, sim.CurrentTime+minClickLatency)
				}
			case RotationAction_FP:
				casted = dk.CastFrostPresence(sim, target)
				if !casted {
					dk.WaitUntil(sim, dk.FrostPresence.CD.ReadyAt())
				} else {
					dk.WaitUntil(sim, sim.CurrentTime+minClickLatency)
				}
			case RotationAction_UP:
				casted = dk.CastUnholyPresence(sim, target)
				if !casted {
					dk.WaitUntil(sim, dk.UnholyPresence.CD.ReadyAt())
				} else {
					dk.WaitUntil(sim, sim.CurrentTime+minClickLatency)
				}
			case RotationAction_RedoSequence:
				o.Reset()
			case RotationAction_FS_IF_KM:
				if dk.KillingMachineAura.IsActive() && !dk.RimeAura.IsActive() {
					casted = dk.CastFrostStrike(sim, target)
				} else if dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
					if dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() < 110 {
						casted = dk.CastHowlingBlast(sim, target)
					} else if dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() > 110 {
						casted = dk.CastHowlingBlast(sim, target)
					} else if !dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() > 110 {
						casted = dk.CastFrostStrike(sim, target)
					} else if !dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() < 110 {
						casted = dk.CastFrostStrike(sim, target)
					}
				} else if !dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
					if dk.CurrentRunicPower() < 110 {
						casted = dk.CastHowlingBlast(sim, target)
					} else {
						casted = dk.CastFrostStrike(sim, target)
					}
				} else {
					casted = dk.CastFrostStrike(sim, target)
					if !casted {
						casted = dk.CastHornOfWinter(sim, target)
					}
				}
				forceAdvance = true
			}

			// Advances the opener
			if (casted && advance) || forceAdvance {
				o.idx += 1
			}

		return casted
	*/
}

func (dk *Deathknight) Wait(sim *core.Simulation) {
	waitUntil := dk.AutoAttacks.MainhandSwingAt
	if dk.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
		waitUntil = core.MinDuration(waitUntil, dk.AutoAttacks.OffhandSwingAt)
	}
	waitUntil = core.MinDuration(waitUntil, dk.AnyRuneReadyAt(sim))
	dk.WaitUntil(sim, waitUntil)
}

func (dk *Deathknight) WaitForResources(sim *core.Simulation) {
	waitUntil := dk.AnySpentRuneReadyAt(sim)
	dk.WaitUntil(sim, waitUntil)
}

func (dk *Deathknight) DoRotation(sim *core.Simulation) {
	target := dk.CurrentTarget

	if dk.Opener.IsOngoing() {
		dk.Opener.DoAction(sim, target, dk)
	} else {
		if dk.Main.IsOngoing() {
			dk.Main.DoAction(sim, target, dk)
		} else {
			if dk.GCD.IsReady(sim) && !dk.IsWaiting() {
				dk.Wait(sim)
			} else { // No resources
				dk.WaitForResources(sim)
			}
		}
	}
}

func (dk *Deathknight) ResetRotation(sim *core.Simulation) {
	dk.Opener.Reset()
	dk.Main.Reset()
}
