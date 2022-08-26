package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
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
	if dk.HornOfWinter.CanCast(sim) {
		dk.HornOfWinter.Cast(sim, target)
	}

	s.Advance()
	return -1
}

func (dk *Deathknight) RotationActionCallback_DRW(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	casted := dk.DancingRuneWeapon.Cast(sim, target)

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
	casted := dk.HeartStrike.Cast(sim, target)
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
			waitTime = sim.CurrentTime
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
			waitTime = sim.CurrentTime
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
			waitTime = sim.CurrentTime
		}
	}
	s.ConditionalAdvance(casted)
	return waitTime
}

func (dk *Deathknight) RotationActionCallback_RD(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	dk.RaiseDead.Cast(sim, target)

	s.Advance()
	return -1
}

func (dk *Deathknight) RotationActionCallback_Reset(sim *core.Simulation, target *core.Unit, s *Sequence) time.Duration {
	s.Reset()
	return -1
}

func (o *Sequence) DoAction(sim *core.Simulation, target *core.Unit, dk *Deathknight) time.Duration {
	action := o.actions[o.idx]
	return action(sim, target, o)
}

func (dk *Deathknight) Wait(sim *core.Simulation) {
	waitUntil := dk.AutoAttacks.MainhandSwingAt
	if dk.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
		waitUntil = core.MinDuration(waitUntil, dk.AutoAttacks.OffhandSwingAt)
	}

	anyRuneAt := dk.AnyRuneReadyAt(sim)
	if anyRuneAt != sim.CurrentTime {
		waitUntil = core.MinDuration(waitUntil, anyRuneAt)
	} else {
		waitUntil = core.MinDuration(waitUntil, dk.AnySpentRuneReadyAt())
	}

	if dk.ButcheryPA != nil {
		waitUntil = core.MinDuration(dk.ButcheryPA.NextActionAt, waitUntil)
	}
	waitUntil = core.MaxDuration(sim.CurrentTime, waitUntil)

	dk.WaitUntil(sim, waitUntil)
}

func (dk *Deathknight) DoRotation(sim *core.Simulation) {
	target := dk.CurrentTarget

	optWait := time.Duration(-1)
	if dk.RotationSequence.IsOngoing() {
		optWait = dk.RotationSequence.DoAction(sim, target, dk)
	}

	if dk.GCD.IsReady(sim) {
		for optWait == 0 && dk.GCD.IsReady(sim) {
			optWait = dk.RotationSequence.DoAction(sim, target, dk)
		}

		if optWait != -1 {
			dk.WaitUntil(sim, optWait)
		} else if dk.GCD.IsReady(sim) {
			dk.Wait(sim)
		}
	}
}
