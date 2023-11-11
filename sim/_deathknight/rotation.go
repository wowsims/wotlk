package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) OnAutoAttack(_ *core.Simulation, _ *core.Spell) {
}

func (dk *Deathknight) OnGCDReady(sim *core.Simulation) {
	dk.tryUseGCD(sim)
}

func (dk *Deathknight) tryUseGCD(sim *core.Simulation) {
	dk.DoRotation(sim)
}

func (dk *Deathknight) RotationActionCallback_IT(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.IcyTouch.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)

	s.ConditionalAdvance(casted && advance)
	return -1
}

func (dk *Deathknight) RotationActionCallback_PS(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.PlagueStrike.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)

	s.ConditionalAdvance(casted && advance)
	return -1
}

func (dk *Deathknight) RotationActionCallback_HW(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	dk.HornOfWinter.Cast(sim, target)

	s.Advance()
	return -1
}

func (dk *Deathknight) RotationActionCallback_DRW(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.DancingRuneWeapon.Cast(sim, target)

	s.ConditionalAdvance(casted)
	return -1
}

func (dk *Deathknight) RotationActionCallback_UF(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.UnholyFrenzy.Cast(sim, target)

	s.ConditionalAdvance(casted)
	return -1
}

func (dk *Deathknight) RotationActionCallback_DS(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.DeathStrike.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)

	s.ConditionalAdvance(casted && advance)
	return -1
}

func (dk *Deathknight) RotationActionCallback_HS(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := false
	if dk.Talents.HeartStrike {
		casted = dk.HeartStrike.Cast(sim, target)
	} else {
		casted = dk.BloodStrike.Cast(sim, target)
	}
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)

	s.ConditionalAdvance(casted && advance)
	return -1
}

func (dk *Deathknight) RotationActionCallback_BT(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.BloodTap.Cast(sim, target)
	s.ConditionalAdvance(casted)
	return sim.CurrentTime
}

func (dk *Deathknight) RotationActionCallback_ERW(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.EmpowerRuneWeapon.Cast(sim, target)
	s.ConditionalAdvance(casted)
	return sim.CurrentTime
}

func (dk *Deathknight) RotationActionCallback_Obli(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	if dk.Deathchill != nil && dk.Deathchill.IsReady(sim) {
		dk.Deathchill.Cast(sim, target)
	}
	casted := dk.Obliterate.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)
	s.ConditionalAdvance(casted && advance)
	return -1
}

func (dk *Deathknight) RotationActionCallback_FS(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	dk.FrostStrike.Cast(sim, target)

	s.Advance()
	return -1
}

func (dk *Deathknight) RotationActionCallback_HB(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	dk.HowlingBlast.Cast(sim, target)

	s.Advance()
	return -1
}

func (dk *Deathknight) RotationActionCallback_Pesti(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.Pestilence.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)
	s.ConditionalAdvance(casted && advance)
	return -1
}

func (dk *Deathknight) RotationActionCallback_BS(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.BloodStrike.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)
	s.ConditionalAdvance(casted && advance)
	return -1
}

func (dk *Deathknight) RotationActionCallback_BB(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.BloodBoil.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)
	s.ConditionalAdvance(casted && advance)
	return -1
}

func (dk *Deathknight) RotationActionCallback_SS(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.ScourgeStrike.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)
	s.ConditionalAdvance(casted && advance)
	return -1
}

func (dk *Deathknight) RotationActionCallback_DND(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.DeathAndDecay.Cast(sim, target)
	s.ConditionalAdvance(casted)
	return -1
}

func (dk *Deathknight) RotationActionCallback_GF(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.GhoulFrenzy.Cast(sim, target)
	s.ConditionalAdvance(casted)
	return -1
}

func (dk *Deathknight) RotationActionCallback_DC(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	dk.DeathCoil.Cast(sim, target)
	s.Advance()
	return -1
}

func (dk *Deathknight) RotationActionCallback_AOTD(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.ArmyOfTheDead.Cast(sim, target)
	s.ConditionalAdvance(casted)
	return -1
}

func (dk *Deathknight) RotationActionCallback_Garg(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.SummonGargoyle.Cast(sim, target)
	s.ConditionalAdvance(casted)
	return -1
}

func (dk *Deathknight) RotationActionCallback_BP(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.BloodPresence.Cast(sim, target)

	waitTime := time.Duration(-1)
	if !casted && !dk.BloodPresence.IsReady(sim) {
		if dk.BloodPresence.CD.ReadyAt() != sim.CurrentTime {
			waitTime = dk.BloodPresence.CD.ReadyAt()
		}
	}

	s.ConditionalAdvance(casted)
	return waitTime
}

func (dk *Deathknight) RotationActionCallback_FP(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.FrostPresence.Cast(sim, target)

	waitTime := time.Duration(-1)
	if !casted && !dk.FrostPresence.IsReady(sim) {
		if dk.FrostPresence.CD.ReadyAt() != sim.CurrentTime {
			waitTime = dk.FrostPresence.CD.ReadyAt()
		}
	}
	s.ConditionalAdvance(casted)
	return waitTime
}

func (dk *Deathknight) RotationActionCallback_UP(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.UnholyPresence.Cast(sim, target)

	waitTime := time.Duration(-1)
	if !casted && !dk.UnholyPresence.IsReady(sim) {
		if dk.UnholyPresence.CD.ReadyAt() != sim.CurrentTime {
			waitTime = dk.UnholyPresence.CD.ReadyAt()
		}
	}
	s.ConditionalAdvance(casted)
	return waitTime
}

func (dk *Deathknight) RotationActionCallback_RD(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	if !dk.Talents.MasterOfGhouls {
		dk.RaiseDead.Cast(sim, target)
	}

	s.Advance()
	return -1
}

func (s *Sequence) DoAction(sim *core.Simulation, target *core.Unit, _ *Deathknight) time.Duration {
	action := s.actions[s.idx]
	return action(sim, target, s)
}

func (dk *Deathknight) Wait(sim *core.Simulation) {
	waitUntil := sim.CurrentTime + time.Millisecond*200

	anyRuneAt := dk.AnyRuneReadyAt(sim)
	if anyRuneAt != sim.CurrentTime {
		waitUntil = min(waitUntil, anyRuneAt)
	} else {
		waitUntil = min(waitUntil, dk.AnySpentRuneReadyAt())
	}
	if sim.Log != nil {
		dk.Log(sim, "DK Wait: %s, any at: %s, any spent at: %s", waitUntil, anyRuneAt, dk.AnySpentRuneReadyAt())
	}

	if dk.ButcheryPA != nil {
		waitUntil = min(dk.ButcheryPA.NextActionAt, waitUntil)
	}
	waitUntil = max(sim.CurrentTime, waitUntil)

	if !dk.Inputs.IsDps {
		target := dk.CurrentTarget
		if dk.IsMainTank() {
			if targetSwingAt := target.AutoAttacks.NextAttackAt(); targetSwingAt > sim.CurrentTime {
				waitUntil = min(waitUntil, targetSwingAt)
			}
		}
	}

	dk.WaitUntil(sim, waitUntil)
}

func (dk *Deathknight) IsMainTank() bool {
	return dk.CurrentTarget.CurrentTarget == &dk.Unit
}

func (dk *Deathknight) DoRotation(sim *core.Simulation) {
	if dk.IsUsingAPL {
		return
	}

	target := dk.CurrentTarget

	optWait := time.Duration(-1)
	if dk.RotationSequence.IsOngoing() {
		if sim.Log != nil {
			dk.Log(sim, "DoSequenceAction")
		}
		optWait = dk.RotationSequence.DoAction(sim, target, dk)
	}

	if dk.GCD.IsReady(sim) {
		if sim.Log != nil {
			dk.Log(sim, "DoGCD")
		}
		for optWait == 0 && dk.GCD.IsReady(sim) {
			if sim.Log != nil {
				dk.Log(sim, "DoAction")
			}
			optWait = dk.RotationSequence.DoAction(sim, target, dk)
		}

		if optWait != -1 {
			if optWait < sim.CurrentTime {
				dk.Wait(sim)
			} else {
				dk.WaitUntil(sim, optWait)
			}
		} else if dk.GCD.IsReady(sim) {
			dk.Wait(sim)
		}
	}
}
