package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) setupFrostSubUnholyERWOpener() {
	dk.setupUnbreakableArmorCooldowns()

	dk.RotationSequence.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence_Pesti).
		NewAction(dk.RotationActionCallback_UA_Frost).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_RD).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence_Pesti).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence1)
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_Obli(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false
	advance := true
	waitTime := time.Duration(-1)

	ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
	bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()
	if sim.CurrentTime+1500*time.Millisecond < core.MinDuration(ffExpiresAt, bpExpiresAt) {
		if dk.Obliterate.CanCast(sim) {
			casted = dk.Obliterate.Cast(sim, target)
			advance = dk.LastOutcome.Matches(core.OutcomeLanded)
		}

		s.ConditionalAdvance(casted && advance)
	} else {
		if dk.Obliterate.CanCast(sim) {
			casted = dk.Obliterate.Cast(sim, target)
			advance = dk.LastOutcome.Matches(core.OutcomeLanded)

			if casted && advance {
				return dk.RotationActionCallback_FrostSubUnholy_RecoverFromPestiMiss(sim, target, s)
			}
		}
	}

	return core.TernaryDuration(casted, -1, waitTime)
}

// TODO: Improve this
func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_FS_KM(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := dk.RotationActionCallback_LastSecondsCast(sim, target)

	if !casted {
		spell := dk.RegularPrioPickSpell(sim, core.NeverExpires)
		if spell != nil {
			casted = spell.Cast(sim, target)
		}

		s.Advance()
	}

	return core.TernaryDuration(casted, -1, sim.CurrentTime)
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_UA_Check1(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.UnbreakableArmor.CanCast(sim) && dk.BloodTap.CanCast(sim) {
		s.Clear().
			NewAction(dk.RotationActionCallback_UA_Frost).
			NewAction(dk.RotationActionCallback_BT).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence1)
	} else {
		s.Advance()
	}
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_UA_Check2(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.UnbreakableArmor.CanCast(sim) && dk.BloodTap.CanCast(sim) {
		s.Clear().
			NewAction(dk.RotationActionCallback_UA_Frost).
			NewAction(dk.RotationActionCallback_BT).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence1)
	} else {
		s.Advance()
	}
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_UA_Check3(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.UnbreakableArmor.CanCast(sim) && dk.BloodTap.CanCast(sim) {
		s.Clear().
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_BT).
			NewAction(dk.RotationActionCallback_UA_Frost).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence_Pesti).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence2)
	} else {
		s.Advance()
	}
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_Sequence1(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	s.Clear().
		NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_Dump).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_UA_Check3).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence_Pesti).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence2)
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_Sequence2(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	s.Clear().
		NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_Dump).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_KM).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_UA_Check1).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_UA_Check2).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence1)
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_RecoverFromPestiMiss(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.LastCast == dk.BloodStrike {
		s.Clear().
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FS).
			NewAction(dk.RotationActionCallback_IT).
			NewAction(dk.RotationActionCallback_PS).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
			NewAction(dk.RotationActionCallback_FS).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_KM).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence1)
	} else {
		s.Clear().
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FS).
			NewAction(dk.RotationActionCallback_IT).
			NewAction(dk.RotationActionCallback_PS).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
			NewAction(dk.RotationActionCallback_FS).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_KM).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence1)
	}

	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_FS_Dump(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false
	waitUntil := time.Duration(-1)

	fr := dk.NormalCurrentFrostRunes()
	uh := dk.NormalCurrentUnholyRunes()

	if fr > 0 && uh > 0 {
		s.Advance()
		return sim.CurrentTime
	}

	casted = dk.RotationActionCallback_LastSecondsCast(sim, target)
	if !casted {
		frAt := dk.NormalFrostRuneReadyAt(sim)
		uhAt := dk.NormalUnholyRuneReadyAt(sim)
		obAt := core.MaxDuration(frAt, uhAt)

		spell := dk.RegularPrioPickSpell(sim, obAt)
		if spell != nil {
			casted = spell.Cast(sim, target)
		} else {
			waitUntil = obAt
			s.Advance()
		}
	}

	return core.TernaryDuration(casted, -1, waitUntil)
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_Sequence_Pesti(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false
	waitUntil := time.Duration(-1)

	ffActive := dk.FrostFeverDisease[target.Index].IsActive()
	bpActive := dk.BloodPlagueDisease[target.Index].IsActive()
	ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
	bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()
	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())

	if core.MinDuration(ffExpiresAt, bpExpiresAt) > sim.CurrentTime+sim.GetRemainingDuration() {
		if dk.FrostStrike.CanCast(sim) && dk.KillingMachineAura.IsActive() {
			casted = dk.FrostStrike.Cast(sim, target)
		} else if dk.HowlingBlast.CanCast(sim) && dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
			casted = dk.HowlingBlast.Cast(sim, target)
		} else if dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() >= 100.0 {
			casted = dk.FrostStrike.Cast(sim, target)
		} else if dk.HowlingBlast.CanCast(sim) && dk.RimeAura.IsActive() {
			casted = dk.HowlingBlast.Cast(sim, target)
		} else if dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() >= 2.0*(fsCost-dk.fr.oblitRPRegen) {
			casted = dk.FrostStrike.Cast(sim, target)
		} else if dk.HornOfWinter.CanCast(sim) {
			casted = dk.HornOfWinter.Cast(sim, target)
		} else {
			waitUntil = -1
		}

		return core.TernaryDuration(casted, -1, waitUntil)
	}

	if !ffActive || !bpActive {
		return dk.RotationActionCallback_FrostSubUnholy_RecoverFromPestiMiss(sim, target, s)
	} else {
		casted = dk.Pestilence.Cast(sim, target)
		advance := dk.LastOutcome.Matches(core.OutcomeLanded)
		if !casted || (casted && !dk.LastOutcome.Matches(core.OutcomeLanded)) {

			if sim.CurrentTime+dk.SpellGCD() > ffExpiresAt || sim.CurrentTime+dk.SpellGCD() > bpExpiresAt {
				return dk.RotationActionCallback_FrostSubUnholy_RecoverFromPestiMiss(sim, target, s)
			} else {
				s.ConditionalAdvance(casted && advance)
				return core.TernaryDuration(casted, -1, waitUntil)
			}
		} else {
			s.ConditionalAdvance(casted && advance)
			return core.TernaryDuration(casted, -1, waitUntil)
		}
	}
}
