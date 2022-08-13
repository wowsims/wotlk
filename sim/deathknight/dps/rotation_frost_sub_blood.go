package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_FS_KM(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.KillingMachineAura.IsActive() {
		if dk.FrostStrike.CanCast(sim) {
			dk.FrostStrike.Cast(sim, target)
		}
	}

	s.Advance()
	return -1
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_FS_Dump_UntilUA(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	waitUntil := time.Duration(-1)

	//abGcd := 1500 * time.Millisecond
	//spGcd := dk.SpellGCD()
	uaAt := dk.UnbreakableArmor.ReadyAt()
	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())

	if !dk.UnbreakableArmor.IsReady(sim) {
		if dk.FrostStrike.CanCast(sim) && dk.KillingMachineAura.IsActive() {
			dk.FrostStrike.Cast(sim, target)
		} else if dk.HowlingBlast.CanCast(sim) && dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
			dk.HowlingBlast.Cast(sim, target)
		} else if dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() > 100.0 {
			dk.FrostStrike.Cast(sim, target)
		} else if dk.HowlingBlast.CanCast(sim) && dk.RimeAura.IsActive() {
			dk.HowlingBlast.Cast(sim, target)
		} else if dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() > 2.0*(fsCost-dk.fr.oblitRPRegen) {
			dk.FrostStrike.Cast(sim, target)
		} else if dk.HornOfWinter.CanCast(sim) {
			dk.HornOfWinter.Cast(sim, target)
		} else {
			waitUntil = uaAt
			s.Advance()
		}
	} else {
		s.Advance()
	}

	return waitUntil
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_FS_Dump(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	waitUntil := time.Duration(-1)

	abGcd := 1500 * time.Millisecond
	spGcd := dk.SpellGCD()
	frAt := dk.FrostRuneReadyAt(sim)
	uhAt := dk.UnholyRuneReadyAt(sim)
	obAt := core.MaxDuration(frAt, uhAt)
	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())

	if sim.CurrentTime+abGcd < obAt && dk.FrostStrike.CanCast(sim) && dk.KillingMachineAura.IsActive() {
		dk.FrostStrike.Cast(sim, target)
	} else if sim.CurrentTime+spGcd < obAt && dk.HowlingBlast.CanCast(sim) && dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
		dk.HowlingBlast.Cast(sim, target)
	} else if sim.CurrentTime+abGcd < obAt && dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() > 100.0 {
		dk.FrostStrike.Cast(sim, target)
	} else if sim.CurrentTime+spGcd < obAt && dk.HowlingBlast.CanCast(sim) && dk.RimeAura.IsActive() {
		dk.HowlingBlast.Cast(sim, target)
	} else if sim.CurrentTime+abGcd < obAt && dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() > 2.0*(fsCost-dk.fr.oblitRPRegen) {
		dk.FrostStrike.Cast(sim, target)
	} else if sim.CurrentTime+spGcd < obAt && dk.HornOfWinter.CanCast(sim) {
		dk.HornOfWinter.Cast(sim, target)
	} else {
		waitUntil = obAt
		s.Advance()
	}

	return waitUntil
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
			NewAction(dk.RotationActionCallback_UA_Frost).
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_KM).
			NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_Dump).
			NewAction(dk.RotationActionCallback_FrostSubBlood_SequenceRotation)
	} else if uaState == FrostSubBloodUAState_Soon {
		s.Clear().
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_Dump_UntilUA).
			NewAction(dk.RotationActionCallback_UA_Frost).
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_KM).
			NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_Dump).
			NewAction(dk.RotationActionCallback_FrostSubBlood_SequenceRotation)
	} else if uaState == FrostSubBloodUAState_CD {
		s.Clear().NewAction(dk.RotationActionCallback_FrostSubBlood_Obli_Check)
	}

	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Obli_Check(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.BloodTap.IsReady(sim) && dk.CurrentDeathBloodRunes() == 2 {
		s.Clear().
			NewAction(dk.RotationActionCallback_Obli).
			NewAction(dk.RotationActionCallback_BT).
			NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_Dump).
			NewAction(dk.RotationActionCallback_FrostSubBlood_SequenceRotation)
	} else {
		s.Clear().
			NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_KM).
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FrostSubBlood_FS_Dump).
			NewAction(dk.RotationActionCallback_FrostSubBlood_SequenceRotation)
	}
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_SequenceRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	s.Clear().
		NewAction(dk.RotationActionCallback_FrostSubBlood_FS_Dump).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FrostSubBlood_FS_KM).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FrostSubBlood_FS_KM).
		NewAction(dk.RotationActionCallback_FrostSubBlood_UA_Check)
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_PrioRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	//Priority List:
	// HoW if refreshing
	// FF & BP up
	// Obliterate
	// Pesti -> BS // BS -> Pesti
	// HB if KM & Rime
	// FS if KM
	// FS if RP > 100
	// HB if Rime
	// FS
	// HW

	gcd := dk.SpellGCD()
	ff := dk.FrostFeverDisease[target.Index].IsActive()
	bp := dk.BloodPlagueDisease[target.Index].IsActive()
	fbAt := core.MinDuration(dk.FrostFeverDisease[target.Index].ExpiresAt(), dk.BloodPlagueDisease[target.Index].ExpiresAt())
	fr := dk.CurrentFrostRunes()
	ur := dk.CurrentUnholyRunes()
	dr := dk.CurrentDeathRunes()
	frAt := dk.NormalSpentFrostRuneReadyAt(sim)
	uhAt := dk.NormalSpentUnholyRuneReadyAt(sim)
	obAt := core.MaxDuration(frAt, uhAt)
	fsCost := float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower())

	if dk.fr.oblitCount == 3 && dk.BloodTap.CanCast(sim) {
		if dk.BloodTap.Cast(sim, target) {
			casted := dk.Pestilence.Cast(sim, target)
			if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
				dk.fr.oblitCount = 0
			} else {
				dk.NextCast = dk.Pestilence
			}
			return casted
		}
	}

	if ff && bp && sim.CurrentTime+gcd > fbAt {
		casted := dk.Pestilence.Cast(sim, target)
		if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
			dk.NextCast = nil
			dk.fr.oblitCount = 0
		}
		return casted
	}

	if dk.NextCast == dk.Pestilence && ff && bp {
		casted := dk.Pestilence.Cast(sim, target)
		if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
			dk.NextCast = nil
			dk.fr.oblitCount = 0
		}
		return casted
	}

	if dk.LastCast == dk.Obliterate && !(dk.Obliterate.CanCast(sim) && dk.BloodTap.IsReady(sim) && fr == 0 && ur == 0 && dr == 2) {
		if dk.KillingMachineAura.IsActive() {
			if dk.FrostStrike.CanCast(sim) && dk.LastOutcome.Matches(core.OutcomeLanded) {
				return dk.FrostStrike.Cast(sim, target)
			}
		}
	}

	if dk.ShouldHornOfWinter(sim) {
		return dk.HornOfWinter.Cast(sim, target)
	} else if !ff {
		dk.fr.oblitCount = 0
		dk.RotationActionCallback_FrostSubBlood_RecoverFromPestiMiss(sim, target, s)
		return false
	} else if !bp {
		dk.fr.oblitCount = 0
		dk.RotationActionCallback_FrostSubBlood_RecoverFromPestiMiss(sim, target, s)
		return false
	} else if dk.Obliterate.CanCast(sim) && fr > 0 && ur > 0 {
		casted := false
		if dk.fr.oblitCount < 2 {
			casted = dk.Obliterate.Cast(sim, target)
			if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
				dk.fr.oblitCount += 1
			}
		} else if dk.fr.oblitCount == 2 && dk.BloodTap.IsReady(sim) {
			casted = dk.Obliterate.Cast(sim, target)
			if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
				dk.fr.oblitCount += 1
			}
		}
		return casted
	} else if dk.Obliterate.CanCast(sim) && dk.BloodTap.CanCast(sim) && fr == 0 && ur == 0 && dr == 2 {
		casted := false
		if dk.fr.oblitCount < 2 {
			casted = dk.Obliterate.Cast(sim, target)
			if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
				dk.fr.oblitCount += 1
			}
		} else if dk.fr.oblitCount == 2 && dk.BloodTap.IsReady(sim) {
			casted = dk.Obliterate.Cast(sim, target)
			if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
				dk.fr.oblitCount += 1
			}
		}
		return casted
	} else if dk.Pestilence.CanCast(sim) && dk.fr.shouldPesti && dk.fr.oblitCount >= 2 {
		casted := dk.Pestilence.Cast(sim, target)
		if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
			dk.fr.shouldPesti = false
		} else {
			dk.NextCast = dk.Pestilence
		}
		return casted
	} else if dk.BloodStrike.CanCast(sim) && dk.fr.oblitCount >= 2 {
		casted := false
		if dk.UnbreakableArmor.CanCast(sim) && fr == 0 {
			casted = dk.UnbreakableArmor.Cast(sim, target)
			dk.castAllMajorCooldowns(sim)
			dk.WaitUntil(sim, sim.CurrentTime)
			dk.fr.shouldPesti = true
			dk.fr.oblitCount = 0
		} else {
			if dk.KillingMachineAura.IsActive() && dk.CurrentRunicPower() > dk.MaxRunicPower()-fsCost {
				casted = dk.FrostStrike.Cast(sim, target)
				if casted {
					dk.fr.shouldPesti = true
					dk.fr.oblitCount = 0
				}
			} else {
				casted = dk.BloodStrike.Cast(sim, target)
				if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
					dk.fr.shouldPesti = true
					dk.fr.oblitCount = 0
				} else {
					dk.fr.shouldPesti = true
					dk.fr.oblitCount = 0
				}
			}
		}
		return casted
	} else if sim.CurrentTime+gcd < obAt {
		if dk.FrostStrike.CanCast(sim) && dk.KillingMachineAura.IsActive() {
			return dk.FrostStrike.Cast(sim, target)
		} else if dk.HowlingBlast.CanCast(sim) && dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
			return dk.HowlingBlast.Cast(sim, target)
		} else if dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() > 100.0 {
			return dk.FrostStrike.Cast(sim, target)
		} else if dk.HowlingBlast.CanCast(sim) && dk.RimeAura.IsActive() {
			return dk.HowlingBlast.Cast(sim, target)
		} else if dk.FrostStrike.CanCast(sim) && dk.CurrentRunicPower() > 2.0*(fsCost-dk.fr.oblitRPRegen) {
			return dk.FrostStrike.Cast(sim, target)
		} else {
			return dk.HornOfWinter.Cast(sim, target)
		}
	} else {
		return false
	}
}

func (dk *DpsDeathknight) setupFrostSubBloodERWOpener() {
	dk.setupUnbreakableArmorCooldowns()

	dk.RotationSequence.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_UA_Frost).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_RD).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
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
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_UA_Frost).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_RD).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Opener_FS_Star).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Opener_FS_Star).
		NewAction(dk.RotationActionCallback_Obli).
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

	ffActive := dk.FrostFeverDisease[target.Index].IsActive()
	bpActive := dk.BloodPlagueDisease[target.Index].IsActive()

	if !ffActive || !bpActive {
		return dk.RotationActionCallback_FrostSubBlood_RecoverFromPestiMiss(sim, target, s)
	} else {
		casted = dk.Pestilence.Cast(sim, target)
		advance := dk.LastOutcome.Matches(core.OutcomeLanded)
		if !casted || (casted && !dk.LastOutcome.Matches(core.OutcomeLanded)) {
			ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
			bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()

			if sim.CurrentTime+dk.SpellGCD() > ffExpiresAt || sim.CurrentTime+dk.SpellGCD() > bpExpiresAt {
				return dk.RotationActionCallback_FrostSubBlood_RecoverFromPestiMiss(sim, target, s)
			} else {
				s.ConditionalAdvance(casted && advance)
				return -1
			}
		} else {
			s.ConditionalAdvance(casted && advance)
			return -1
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
