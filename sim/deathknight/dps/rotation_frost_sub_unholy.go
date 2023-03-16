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
		NewAction(dk.RotationActionCallback_UA_Frost).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_HB).
		NewAction(dk.RotationActionCallback_Frost_Pesti_ERW).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_HB).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationAction_CancelBT).
		NewAction(dk.RotationActionCallback_RD).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_HB).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_HB).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_HB).
		NewAction(dk.RotationActionCallback_Pesti).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence1)
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_FS_HB(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if !dk.canCastAbilityBeforeDiseasesExpire(sim, target) {
		s.Advance()
		return sim.CurrentTime
	}
	return dk.RotationActionCallback_Frost_FS_HB(sim, target, s)
}

func (dk *DpsDeathknight) canCastAbilityBeforeDiseasesExpire(sim *core.Simulation, target *core.Unit) bool {
	ffExpiresAt := dk.FrostFeverSpell.Dot(target).ExpiresAt()
	bpExpiresAt := dk.BloodPlagueSpell.Dot(target).ExpiresAt()
	return sim.CurrentTime+1500*time.Millisecond < core.MinDuration(ffExpiresAt, bpExpiresAt)
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_Obli(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false
	advance := true
	waitTime := time.Duration(-1)

	if dk.canCastAbilityBeforeDiseasesExpire(sim, target) {
		if dk.Obliterate.CanCast(sim, nil) {
			if dk.Deathchill != nil && dk.Deathchill.IsReady(sim) {
				dk.Deathchill.Cast(sim, target)
			}
			casted = dk.Obliterate.Cast(sim, target)
			advance = dk.LastOutcome.Matches(core.OutcomeLanded)
		}

		s.ConditionalAdvance(casted && advance)
	} else {
		if dk.Obliterate.CanCast(sim, nil) {
			if dk.Deathchill != nil && dk.Deathchill.IsReady(sim) {
				dk.Deathchill.Cast(sim, target)
			}
			casted = dk.Obliterate.Cast(sim, target)
			advance = dk.LastOutcome.Matches(core.OutcomeLanded)

			if casted && advance {
				return dk.RotationActionCallback_FrostSubUnholy_RecoverFromPestiMiss(sim, target, s)
			}
		}
	}

	return core.TernaryDuration(casted, -1, waitTime)
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_FS_KM(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := dk.RotationActionCallback_LastSecondsCast(sim, target)

	if !casted {
		if !dk.canCastAbilityBeforeDiseasesExpire(sim, target) {
			s.Advance()
			return sim.CurrentTime
		}

		spell := dk.RegularPrioPickSpell(sim, target, core.NeverExpires)
		if spell != nil {
			casted = spell.Cast(sim, target)
		}

		s.Advance()
	}

	return core.TernaryDuration(casted, -1, sim.CurrentTime)
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_Dump_Until_Deaths(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	/*
		We need to have the first death up before we UA + BT + Oblit, since if only the 2nd
		death rune is up, UA then Blood Tap will convert and refresh the first and the second will have a 10s CD from UA
	*/
	if dk.LeftBloodRuneReady() {
		s.Advance()
		return sim.CurrentTime
	}

	spell := dk.RegularPrioPickSpell(sim, target, dk.DeathRuneRegenAt(1))

	if spell != nil {
		spell.Cast(sim, target)
	}

	return -1
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_UA_Check1(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.UnbreakableArmor.CanCast(sim, nil) && dk.BloodTap.CanCast(sim, nil) {
		s.Clear().
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Dump_Until_Deaths).
			NewAction(dk.RotationActionCallback_UA_Frost).
			NewAction(dk.RotationActionCallback_BT).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence2)
	} else {
		s.Advance()
	}
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_UA_Check2(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.UnbreakableArmor.CanCast(sim, nil) && dk.BloodTap.CanCast(sim, nil) {
		s.Clear().
			NewAction(dk.RotationActionCallback_UA_Frost).
			NewAction(dk.RotationActionCallback_BT).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence2)
	} else {
		s.Advance()
	}
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_UA_Check3(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if (dk.UnbreakableArmor.TimeToReady(sim) < 2500*time.Millisecond+sim.CurrentTime) && (dk.BloodTap.TimeToReady(sim) < 2500*time.Millisecond+sim.CurrentTime) {
		s.Clear().
			NewAction(dk.RotationActionCallback_BT).
			NewAction(dk.RotationActionCallback_UA_Frost).
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence1)
	} else {
		s.Advance()
	}
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_Sequence1(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	s.Clear().
		NewAction(dk.RotationActionCallback_EndOfFightCheck).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_Dump).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_UA_Check1).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Sequence2)
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_Pesti(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	desiredGrace := 10 * time.Millisecond
	currentGrace := dk.RuneGraceAt(0, sim.CurrentTime)
	ffExpiresAt := dk.FrostFeverSpell.Dot(target).ExpiresAt()
	bpExpiresAt := dk.BloodPlagueSpell.Dot(target).ExpiresAt()
	diseaseExpiresAt := core.MinDuration(ffExpiresAt, bpExpiresAt)
	waitUntil := sim.CurrentTime + (desiredGrace - currentGrace)

	if diseaseExpiresAt <= sim.CurrentTime {
		return dk.RotationActionCallback_FrostSubUnholy_RecoverFromPestiMiss(sim, target, s)
	}
	if currentGrace < desiredGrace && diseaseExpiresAt > waitUntil {
		return waitUntil
	}
	casted := dk.Pestilence.Cast(sim, target)
	s.ConditionalAdvance(casted && dk.LastOutcome.Matches(core.OutcomeLanded))
	return -1
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_Sequence2(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	s.Clear().
		NewAction(dk.RotationAction_CancelBT).
		NewAction(dk.RotationActionCallback_EndOfFightCheck).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_Dump).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_KM).
		NewAction(dk.RotationActionCallback_EndOfFightCheck).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Pesti).
		//NewAction(dk.RotationActionCallback_FrostSubUnholy_UA_Check3).
		NewAction(dk.RotationActionCallback_BS).
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
			NewAction(dk.RotationActionCallback_Frost_FS_HB).
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
			NewAction(dk.RotationActionCallback_Frost_FS_HB).
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
		delayAmount := time.Second
		spell := dk.RegularPrioPickSpell(sim, target, obAt+delayAmount)
		if spell != nil {
			casted = spell.Cast(sim, target)
		} else {
			waitUntil = obAt
			s.Advance()
		}
	}

	return core.TernaryDuration(casted, -1, waitUntil)
}
