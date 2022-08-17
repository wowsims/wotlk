package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) RotationActionCallback_BS_Frost(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false

	//fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())
	//if !(!dk.UnbreakableArmor.IsReady(sim) && dk.UnbreakableArmor.ReadyAt() > sim.CurrentTime+10*time.Second &&
	//	!dk.BloodTap.IsReady(sim) && dk.BloodTap.ReadyAt() > sim.CurrentTime+10*time.Second &&
	//	dk.CurrentRunicPower() >= 2.0*(fsCost-dk.fr.oblitRPRegen)) {
	casted = dk.BloodStrike.Cast(sim, target)
	//}

	s.Advance()
	return core.TernaryDuration(casted, -1, sim.CurrentTime)
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Obli(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
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
		waitTime = sim.CurrentTime
		s.Advance()
	}

	return core.TernaryDuration(casted, -1, waitTime)
}

// TODO: Improve this
func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_FS_KM(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false
	ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
	bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()
	if dk.KillingMachineAura.IsActive() && sim.CurrentTime+1500*time.Millisecond < core.MinDuration(ffExpiresAt, bpExpiresAt) {
		if dk.FrostStrike.CanCast(sim) {
			dk.FrostStrike.Cast(sim, target)
		}
	}

	s.Advance()
	return core.TernaryDuration(casted, -1, sim.CurrentTime)
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_FS_Dump_UntilBR(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false
	waitUntil := time.Duration(-1)

	abGcd := 1500 * time.Millisecond
	spGcd := dk.SpellGCD()
	br := dk.CurrentBloodRunes() + dk.CurrentDeathRunes()
	ddAt := dk.SpentBloodRuneReadyAt()
	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())

	tol := 0 * time.Millisecond
	if br == 0 {
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
			s.Advance()
			waitUntil = ddAt
		}
	} else if br == 1 {
		if sim.CurrentTime+abGcd <= ddAt+tol && dk.FrostStrike.CanCast(sim) && dk.KillingMachineAura.IsActive() {
			casted = dk.FrostStrike.Cast(sim, target)
		} else if sim.CurrentTime+spGcd <= ddAt+tol && dk.HowlingBlast.CanCast(sim) && dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
			casted = dk.HowlingBlast.Cast(sim, target)
		} else if sim.CurrentTime+abGcd <= ddAt+tol && dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() >= 100.0 {
			casted = dk.FrostStrike.Cast(sim, target)
		} else if sim.CurrentTime+spGcd <= ddAt+tol && dk.HowlingBlast.CanCast(sim) && dk.RimeAura.IsActive() {
			casted = dk.HowlingBlast.Cast(sim, target)
		} else if sim.CurrentTime+abGcd <= ddAt+tol && dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() >= 2.0*(fsCost-dk.fr.oblitRPRegen) {
			casted = dk.FrostStrike.Cast(sim, target)
		} else if sim.CurrentTime+spGcd <= ddAt+tol && dk.HornOfWinter.CanCast(sim) {
			casted = dk.HornOfWinter.Cast(sim, target)
		} else {
			s.Advance()
			waitUntil = sim.CurrentTime
		}
	} else if br == 2 {
		ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
		bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()
		if dk.KillingMachineAura.IsActive() && sim.CurrentTime+1500*time.Millisecond < core.MinDuration(ffExpiresAt, bpExpiresAt) &&
			dk.DeathRuneRevertAt() > sim.CurrentTime+1500*time.Millisecond {
			if dk.FrostStrike.CanCast(sim) {
				casted = dk.FrostStrike.Cast(sim, target)
			}
		}

		s.Advance()
		waitUntil = sim.CurrentTime
	}

	return core.TernaryDuration(casted, -1, waitUntil)
}
func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_FS_Dump_UntilUA(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false
	waitUntil := time.Duration(-1)

	//abGcd := 1500 * time.Millisecond
	//spGcd := dk.SpellGCD()
	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())

	if !dk.UnbreakableArmor.IsReady(sim) {
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
			s.Advance()
			waitUntil = dk.UnbreakableArmor.ReadyAt()
		}
	} else {
		s.Advance()
		waitUntil = sim.CurrentTime
	}

	return core.TernaryDuration(casted, -1, waitUntil)
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_FS_Dump(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false
	waitUntil := time.Duration(-1)

	abGcd := 1500 * time.Millisecond
	spGcd := dk.SpellGCD()
	fr := dk.NormalCurrentFrostRunes()
	uh := dk.NormalCurrentUnholyRunes()
	frAt := dk.NormalFrostRuneReadyAt(sim)
	uhAt := dk.NormalUnholyRuneReadyAt(sim)
	//frGp := dk.CurrentFrostRuneGrace(sim)
	//uhGp := dk.CurrentUnholyRuneGrace(sim)
	obAt := core.MaxDuration(frAt, uhAt)
	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())

	if fr == 2 && uh == 2 {
		casted := false
		waitTime := time.Duration(-1)

		ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
		bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()
		if sim.CurrentTime+1500*time.Millisecond < core.MinDuration(ffExpiresAt, bpExpiresAt) {
			if dk.Obliterate.CanCast(sim) {
				casted = dk.Obliterate.Cast(sim, target)
				if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
					dk.fr.oblitCount += 1
				}
			}
		} else {
			waitTime = sim.CurrentTime
			s.Advance()
		}

		return core.TernaryDuration(casted, -1, waitTime)
	} else if fr > 0 && uh > 0 {
		casted := false
		ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
		bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()
		if dk.fr.oblitCount == 1 && dk.FrostStrike.CanCast(sim) && dk.KillingMachineAura.IsActive() && sim.CurrentTime+1500*time.Millisecond < core.MinDuration(ffExpiresAt, bpExpiresAt) &&
			sim.CurrentTime+abGcd <= obAt+1000*time.Millisecond {
			dk.FrostStrike.Cast(sim, target)
		} else {
			if dk.Obliterate.CanCast(sim) {
				casted = dk.Obliterate.Cast(sim, target)
				advance := dk.LastOutcome.Matches(core.OutcomeLanded)
				if casted && advance {
					dk.fr.oblitCount += 1
				}

				if dk.fr.oblitCount == 2 {
					s.ConditionalAdvance(casted && advance)
					dk.fr.oblitCount = 0
				}
			} else {
				s.Advance()
			}
		}

		return core.TernaryDuration(casted, -1, sim.CurrentTime)
		/*
			compareToAb := (frGp-abGcd > -500*time.Millisecond) && (uhGp-abGcd > -500*time.Millisecond)
			compareToSp := (frGp-spGcd > -500*time.Millisecond) && (uhGp-spGcd > -500*time.Millisecond)
			if compareToAb && dk.FrostStrike.CanCast(sim) && dk.KillingMachineAura.IsActive() {
				casted = dk.FrostStrike.Cast(sim, target)
			} else if compareToSp && dk.HowlingBlast.CanCast(sim) && dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
				casted = dk.HowlingBlast.Cast(sim, target)
			} else if compareToAb && dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() >= 100.0 {
				casted = dk.FrostStrike.Cast(sim, target)
			} else if compareToSp && dk.HowlingBlast.CanCast(sim) && dk.RimeAura.IsActive() {
				casted = dk.HowlingBlast.Cast(sim, target)
			} else if compareToAb && dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() >= 2.0*(fsCost-dk.fr.oblitRPRegen) {
				casted = dk.FrostStrike.Cast(sim, target)
			} else if compareToSp && dk.HornOfWinter.CanCast(sim) {
				casted = dk.HornOfWinter.Cast(sim, target)
			} else {
				if sim.CurrentTime+sim.GetRemainingDuration() < obAt {
					casted = dk.FrostStrike.Cast(sim, target)
				} else {
					casted := false
					waitTime := time.Duration(-1)

					ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
					bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()
					if sim.CurrentTime+1500*time.Millisecond < core.MinDuration(ffExpiresAt, bpExpiresAt) {
						if dk.Obliterate.CanCast(sim) {
							casted = dk.Obliterate.Cast(sim, target)
							advance := dk.LastOutcome.Matches(core.OutcomeLanded)
							s.ConditionalAdvance(casted && advance)
						}
					} else {
						waitTime = sim.CurrentTime
						s.Advance()
					}

					return core.TernaryDuration(casted, -1, waitTime)
				}
			}
		*/
	} else {
		if sim.CurrentTime+abGcd <= obAt && dk.FrostStrike.CanCast(sim) && dk.KillingMachineAura.IsActive() {
			casted = dk.FrostStrike.Cast(sim, target)
		} else if sim.CurrentTime+spGcd <= obAt && dk.HowlingBlast.CanCast(sim) && dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
			casted = dk.HowlingBlast.Cast(sim, target)
		} else if sim.CurrentTime+abGcd <= obAt && dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() >= 100.0 {
			casted = dk.FrostStrike.Cast(sim, target)
		} else if sim.CurrentTime+spGcd <= obAt && dk.HowlingBlast.CanCast(sim) && dk.RimeAura.IsActive() {
			casted = dk.HowlingBlast.Cast(sim, target)
		} else if sim.CurrentTime+abGcd <= obAt && dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() >= 2.0*(fsCost-dk.fr.oblitRPRegen) {
			casted = dk.FrostStrike.Cast(sim, target)
		} else if sim.CurrentTime+spGcd <= obAt && dk.HornOfWinter.CanCast(sim) {
			casted = dk.HornOfWinter.Cast(sim, target)
		} else {
			waitUntil = obAt
		}
	}

	return core.TernaryDuration(casted, -1, waitUntil)
}

type FrostSubBloodUAState uint8

const (
	FrostSubBloodUAState_Now FrostSubBloodUAState = iota
	FrostSubBloodUAState_Soon
	FrostSubBloodUAState_CD
)

func (dk *DpsDeathknight) frCheckForUATime(sim *core.Simulation) FrostSubBloodUAState {
	if dk.UnbreakableArmor.IsReady(sim) {
		return FrostSubBloodUAState_Now
	} else {
		if dk.UnbreakableArmor.ReadyAt() < sim.CurrentTime+dk.SpellGCD()*2 {
			return FrostSubBloodUAState_Soon
		} else {
			return FrostSubBloodUAState_CD
		}
	}
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_UA_Check(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	uaState := dk.frCheckForUATime(sim)
	if uaState == FrostSubBloodUAState_Now {
		s.Clear().
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_Dump_UntilBR).
			NewAction(dk.RotationActionCallback_UA_Frost).
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_KM).
			NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
			NewAction(dk.RotationActionCallback_FrostSubBlood_SequenceRotation)
	} else if uaState == FrostSubBloodUAState_Soon {
		s.Clear().
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_Dump_UntilBR).
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_Dump_UntilUA).
			NewAction(dk.RotationActionCallback_UA_Frost).
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_KM).
			NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
			NewAction(dk.RotationActionCallback_FrostSubBlood_SequenceRotation)
	} else if uaState == FrostSubBloodUAState_CD {
		s.Clear().NewAction(dk.RotationActionCallback_FrostSubBlood_Obli_Check)
	}

	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Obli_Check(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.BloodTap.IsReady(sim) && dk.RuneIsDeath(0) && dk.RuneIsDeath(1) {
		s.Clear().
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_Dump_UntilBR).
			NewAction(dk.RotationActionCallback_FrostSubBlood_Obli).
			NewAction(dk.RotationActionCallback_BT).
			NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
			NewAction(dk.RotationActionCallback_FrostSubBlood_SequenceRotation)
	} else {
		s.Clear().
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_Dump_UntilBR).
			NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_KM).
			NewAction(dk.RotationActionCallback_BS_Frost).
			NewAction(dk.RotationActionCallback_FrostSubBlood_SequenceRotation)
	}
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_SequenceRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	s.Clear().
		NewAction(dk.RotationActionCallback_FrostSubBlood_FS_Dump).
		NewAction(dk.RotationActionCallback_FrostSubBlood_UA_Check)
	return sim.CurrentTime
}

func (dk *DpsDeathknight) setupFrostSubBloodERWOpener() {
	dk.setupUnbreakableArmorCooldowns()

	dk.RotationSequence.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_UA_Frost).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Obli).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Obli).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_RD).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Obli).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Obli).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FrostSubBlood_SequenceRotation)
}

func (dk *DpsDeathknight) setupFrostSubBloodNoERWOpener() {
	dk.setupUnbreakableArmorCooldowns()

	dk.RotationSequence.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Obli).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_UA_Frost).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_RD).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Opener_FS_Star).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Obli).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Opener_FS_Star).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Opener_FS_Star).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Opener_FS_Star).
		NewAction(dk.RotationActionCallback_FrostSubBlood_SequenceRotation)
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_RecoverFromPestiMiss(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.LastCast == dk.BloodStrike {
		s.Clear().
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FS).
			NewAction(dk.RotationActionCallback_IT).
			NewAction(dk.RotationActionCallback_PS).
			NewAction(dk.RotationActionCallback_Obli).
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
			NewAction(dk.RotationActionCallback_FrostSubBlood_SequenceRotation)
	} else {
		s.Clear().
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FS).
			NewAction(dk.RotationActionCallback_IT).
			NewAction(dk.RotationActionCallback_PS).
			NewAction(dk.RotationActionCallback_Obli).
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
			NewAction(dk.RotationActionCallback_FrostSubBlood_SequenceRotation)
	}

	dk.NextCast = nil
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Sequence_Pesti(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
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
		return dk.RotationActionCallback_FrostSubBlood_RecoverFromPestiMiss(sim, target, s)
	} else {
		casted = dk.Pestilence.Cast(sim, target)
		advance := dk.LastOutcome.Matches(core.OutcomeLanded)
		if !casted || (casted && !dk.LastOutcome.Matches(core.OutcomeLanded)) {

			if sim.CurrentTime+dk.SpellGCD() > ffExpiresAt || sim.CurrentTime+dk.SpellGCD() > bpExpiresAt {
				return dk.RotationActionCallback_FrostSubBlood_RecoverFromPestiMiss(sim, target, s)
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

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Main_Pesti(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false

	ffActive := dk.FrostFeverDisease[target.Index].IsActive()
	bpActive := dk.BloodPlagueDisease[target.Index].IsActive()

	if !ffActive || !bpActive {
		return dk.RotationActionCallback_FrostSubBlood_RecoverFromPestiMiss(sim, target, s)
	} else {
		casted = dk.Pestilence.Cast(sim, target)
		if !casted || (casted && !dk.LastOutcome.Matches(core.OutcomeLanded)) {
			ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
			bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()

			if sim.CurrentTime+dk.SpellGCD() > ffExpiresAt || sim.CurrentTime+dk.SpellGCD() > bpExpiresAt {
				return dk.RotationActionCallback_FrostSubBlood_RecoverFromPestiMiss(sim, target, s)
			} else {
				return -1
			}
		} else {
			return -1
		}
	}
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Main_FS_Star(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false
	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())
	if dk.FrostStrike.CanCast(sim) && (dk.PercentRunicPower() >= 0.95 || (dk.KillingMachineAura.IsActive() && dk.CurrentRunicPower() >= 2.0*(fsCost-dk.fr.oblitRPRegen))) {
		casted = dk.FrostStrike.Cast(sim, target)
	} else if dk.HowlingBlast.CanCast(sim) && dk.RimeAura.IsActive() {
		casted = dk.HowlingBlast.Cast(sim, target)
	} else if dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() >= 2.0*(fsCost-dk.fr.oblitRPRegen) {
		casted = dk.FrostStrike.Cast(sim, target)
		if !casted {
			casted = dk.HornOfWinter.Cast(sim, target)
		}
	} else {
		casted = dk.HornOfWinter.Cast(sim, target)
	}

	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Opener_FS_Star(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	casted := false
	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())
	if dk.FrostStrike.CanCast(sim) && (dk.PercentRunicPower() >= 0.95 || (dk.KillingMachineAura.IsActive() && dk.CurrentRunicPower() >= 2.0*(fsCost-dk.fr.oblitRPRegen))) {
		dk.FrostStrike.Cast(sim, target)
		s.Advance()
	} else if dk.HowlingBlast.CanCast(sim) && dk.RimeAura.IsActive() {
		casted = dk.HowlingBlast.Cast(sim, target)
		s.ConditionalAdvance(casted)
	} else if dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() >= 2.0*(fsCost-dk.fr.oblitRPRegen) {
		casted = dk.FrostStrike.Cast(sim, target)
		if !casted {
			dk.HornOfWinter.Cast(sim, target)
		}
		s.Advance()
	} else {
		s.Advance()
	}

	return -1
}
