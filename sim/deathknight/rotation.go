package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	if !dk.Opener.IsOngoing() && !dk.Inputs.IsDps {
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

func (dk *Deathknight) RotationActionCallback_IT(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.IcyTouch.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)

	s.ConditionalAdvance(casted && advance)
	return casted
}

func (dk *Deathknight) RotationActionCallback_PS(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.PlagueStrike.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)

	s.ConditionalAdvance(casted && advance)
	return casted
}

func (dk *Deathknight) RotationActionCallback_HW(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.HornOfWinter.Cast(sim, target)

	s.ConditionalAdvance(true)
	return casted
}

func (dk *Deathknight) RotationActionCallback_UA(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.UnbreakableArmor.Cast(sim, target)
	dk.WaitUntil(sim, sim.CurrentTime)
	s.ConditionalAdvance(casted)
	return casted
}

func (dk *Deathknight) RotationActionCallback_BT(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.BloodTap.Cast(sim, target)
	dk.WaitUntil(sim, sim.CurrentTime)

	s.ConditionalAdvance(casted)
	return casted
}

func (dk *Deathknight) RotationActionCallback_ERW(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.EmpowerRuneWeapon.Cast(sim, target)
	dk.WaitUntil(sim, sim.CurrentTime)

	s.ConditionalAdvance(casted)
	return casted
}

func (dk *Deathknight) RotationActionCallback_Obli(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.Obliterate.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)

	s.ConditionalAdvance(casted && advance)
	return casted
}

func (dk *Deathknight) RotationActionCallback_FS(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.FrostStrike.Cast(sim, target)

	s.Advance()
	return casted
}

func (dk *Deathknight) RotationActionCallback_Pesti(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.Pestilence.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)

	s.ConditionalAdvance(casted && advance)
	return casted
}

func (dk *Deathknight) RotationActionCallback_BS(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.BloodStrike.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)

	s.ConditionalAdvance(casted && advance)
	return casted
}

func (dk *Deathknight) RotationActionCallback_BB(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.BloodBoil.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)

	s.ConditionalAdvance(casted && advance)
	return casted
}

func (dk *Deathknight) RotationActionCallback_SS(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.ScourgeStrike.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)

	s.ConditionalAdvance(casted && advance)
	return casted
}

func (dk *Deathknight) RotationActionCallback_DND(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.DeathAndDecay.Cast(sim, target)

	s.ConditionalAdvance(casted)
	return casted
}

func (dk *Deathknight) RotationActionCallback_GF(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.GhoulFrenzy.Cast(sim, target)

	s.ConditionalAdvance(casted)
	return casted
}

func (dk *Deathknight) RotationActionCallback_DC(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.DeathCoil.Cast(sim, target)

	s.Advance()
	return casted
}

func (dk *Deathknight) RotationActionCallback_AOTD(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.ArmyOfTheDead.Cast(sim, target)

	s.ConditionalAdvance(casted)
	return casted
}

func (dk *Deathknight) RotationActionCallback_Garg(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.SummonGargoyle.Cast(sim, target)

	s.ConditionalAdvance(casted)
	return casted
}

func (dk *Deathknight) RotationActionCallback_BP(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.BloodPresence.Cast(sim, target)
	if !casted && !dk.BloodPresence.IsReady(sim) {
		dk.WaitUntil(sim, dk.BloodPresence.CD.ReadyAt())
	} else {
		dk.WaitUntil(sim, sim.CurrentTime)
	}
	s.ConditionalAdvance(casted)
	return casted
}

func (dk *Deathknight) RotationActionCallback_FP(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.FrostPresence.Cast(sim, target)
	if !casted && !dk.FrostPresence.IsReady(sim) {
		dk.WaitUntil(sim, dk.FrostPresence.CD.ReadyAt())
	} else {
		dk.WaitUntil(sim, sim.CurrentTime)
	}
	s.ConditionalAdvance(casted)
	return casted
}

func (dk *Deathknight) RotationActionCallback_UP(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.UnholyPresence.Cast(sim, target)
	if !casted && !dk.UnholyPresence.IsReady(sim) {
		dk.WaitUntil(sim, dk.UnholyPresence.CD.ReadyAt())
	} else {
		dk.WaitUntil(sim, sim.CurrentTime)
	}
	s.ConditionalAdvance(casted)
	return casted
}

func (dk *Deathknight) RotationActionCallback_RD(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	casted := dk.RaiseDead.Cast(sim, target)

	s.ConditionalAdvance(true)
	return casted
}

func (dk *Deathknight) RotationActionCallback_Reset(sim *core.Simulation, target *core.Unit, s *Sequence) bool {
	s.Reset()
	return false
}

func (o *Sequence) DoAction(sim *core.Simulation, target *core.Unit, dk *Deathknight) bool {
	action := o.actions[o.idx]
	return action(sim, target, o)
}

func (dk *Deathknight) Wait(sim *core.Simulation) {
	waitUntil := dk.AutoAttacks.MainhandSwingAt
	if dk.AutoAttacks.OffhandSwingAt > sim.CurrentTime {
		waitUntil = core.MinDuration(waitUntil, dk.AutoAttacks.OffhandSwingAt)
	}
	waitUntil = core.MinDuration(waitUntil, dk.AnySpentRuneReadyAt())
	if dk.ButcheryPA != nil {
		waitUntil = core.MinDuration(dk.ButcheryPA.NextActionAt, waitUntil)
	}
	waitUntil = core.MaxDuration(sim.CurrentTime, waitUntil)
	dk.WaitUntil(sim, waitUntil)
}

func (dk *Deathknight) DoRotation(sim *core.Simulation) {
	target := dk.CurrentTarget

	casted := false
	if dk.Opener.IsOngoing() {
		casted = dk.Opener.DoAction(sim, target, dk)
	} else if dk.Main.IsOngoing() {
		casted = dk.Main.DoAction(sim, target, dk)
	}

	if !casted || (dk.GCD.IsReady(sim) && !dk.IsWaiting()) {
		dk.Wait(sim)
	}
}
