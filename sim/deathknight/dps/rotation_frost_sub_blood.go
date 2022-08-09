package dps

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

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
	fr := dk.CurrentFrostRunes()
	ur := dk.CurrentUnholyRunes()
	dr := dk.CurrentDeathRunes()
	frAt := dk.NormalSpentFrostRuneReadyAt(sim)
	uhAt := dk.NormalSpentUnholyRuneReadyAt(sim)
	obAt := core.MaxDuration(frAt, uhAt)

	if dk.NextCast == dk.Pestilence {
		casted := dk.CastPestilence(sim, target)
		if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
			dk.NextCast = nil
			dk.fr.oblitCount = 0
		}
		return casted
	}

	if dk.LastCast == dk.Obliterate {
		if dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
			if dk.CanHowlingBlast(sim) && dk.LastOutcome.Matches(core.OutcomeLanded) {
				return dk.CastHowlingBlast(sim, target)
			}
		} else if dk.KillingMachineAura.IsActive() {
			if dk.CanFrostStrike(sim) && dk.LastOutcome.Matches(core.OutcomeLanded) {
				return dk.CastFrostStrike(sim, target)
			}
		}
	}

	if dk.fr.oblitCount == 3 && dk.CanBloodTap(sim) {
		if dk.CastBloodTap(sim, target) {
			casted := dk.CastPestilence(sim, target)
			if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
				dk.fr.oblitCount = 0
			} else {
				dk.NextCast = dk.Pestilence
			}
			return casted
		}
	}

	if dk.ShouldHornOfWinter(sim) {
		return dk.CastHornOfWinter(sim, target)
	} else if !ff {
		return dk.CastIcyTouch(sim, target)
	} else if !bp {
		return dk.CastPlagueStrike(sim, target)
	} else if dk.CanObliterate(sim) && fr > 0 && ur > 0 {
		casted := false
		if dk.fr.oblitCount < 2 {
			casted = dk.CastObliterate(sim, target)
			if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
				dk.fr.oblitCount += 1
			}
		} else if dk.fr.oblitCount == 2 && dk.BloodTap.IsReady(sim) {
			casted = dk.CastObliterate(sim, target)
			if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
				dk.fr.oblitCount += 1
			}
		}
		return casted
	} else if dk.CanObliterate(sim) && dk.BloodTap.IsReady(sim) && fr == 0 && ur == 0 && dr == 2 {
		casted := false
		if dk.fr.oblitCount < 2 {
			casted = dk.CastObliterate(sim, target)
			if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
				dk.fr.oblitCount += 1
			}
		} else if dk.fr.oblitCount == 2 && dk.BloodTap.IsReady(sim) {
			casted = dk.CastObliterate(sim, target)
			if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
				dk.fr.oblitCount += 1
			}
		}
		return casted
	} else if dk.CanBloodStrike(sim) && dk.fr.canBloodStrike && dk.fr.oblitCount >= 2 {
		casted := false
		if dk.CanUnbreakableArmor(sim) {
			casted = dk.CastUnbreakableArmor(sim, target)
			dk.castAllMajorCooldowns(sim)
			if casted {
				dk.fr.canBloodStrike = false
				casted = dk.CastPestilence(sim, target)
				if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
					dk.fr.oblitCount = 0
					dk.fr.canBloodStrike = true
				}
			} else {
				dk.WaitUntil(sim, sim.CurrentTime)
			}
		} else {
			casted = dk.CastBloodStrike(sim, target)
			if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
				dk.fr.canBloodStrike = false
			}
		}
		return casted
	} else if dk.CanPestilence(sim) && dk.fr.oblitCount >= 2 {
		casted := dk.CastPestilence(sim, target)
		if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
			dk.fr.oblitCount = 0
			dk.fr.canBloodStrike = true
		}
		return casted
	} else if sim.CurrentTime+gcd < obAt {
		if dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
			return dk.CastHowlingBlast(sim, target)
		} else if dk.KillingMachineAura.IsActive() {
			return dk.CastFrostStrike(sim, target)
		} else if dk.CurrentRunicPower() > 100.0 {
			return dk.CastFrostStrike(sim, target)
		} else if dk.RimeAura.IsActive() {
			return dk.CastHowlingBlast(sim, target)
		} else if dk.CurrentRunicPower()-float64(core.RuneCost(dk.FrostStrike.CurCast.Cost).RunicPower()) > 14.0 {
			return dk.CastFrostStrike(sim, target)
		} else {
			return dk.CastHornOfWinter(sim, target)
		}
	} else {
		return false
	}
}

func (dk *DpsDeathknight) setupFrostSubBloodERWOpener() {
	dk.setupUnbreakableArmorCooldowns()

	dk.Opener.
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
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_HW).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_FS)

	dk.Main.
		NewAction(dk.RotationActionCallback_FrostSubBlood_PrioRotation)
}

func (dk *DpsDeathknight) setupFrostSubBloodNoERWOpener() {
	dk.setupUnbreakableArmorCooldowns()

	dk.Opener.
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
		NewAction(dk.RotationActionCallback_FrostSubBlood_Opener_FS_Star)

	dk.Main.
		NewAction(dk.RotationActionCallback_FrostSubBlood_PrioRotation)
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_RecoverFromPestiMiss(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	if dk.LastCast == dk.BloodStrike {
		s.Clear().
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FS).
			NewAction(dk.RotationActionCallback_IT).
			NewAction(dk.RotationActionCallback_PS).
			NewAction(dk.RotationActionCallback_Obli).
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FrostSubBlood_Sequence_Pesti).
			NewAction(dk.RotationActionCallback_FrostSubBlood_PrioRotation)
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
			NewAction(dk.RotationActionCallback_FrostSubBlood_PrioRotation)
	}

	return false
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Sequence_Pesti(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	ffActive := dk.FrostFeverDisease[target.Index].IsActive()
	bpActive := dk.BloodPlagueDisease[target.Index].IsActive()

	if !ffActive || !bpActive {
		dk.RotationActionCallback_FrostSubBlood_RecoverFromPestiMiss(sim, target, s)
		return casted
	} else {
		casted = dk.CastPestilence(sim, target)
		advance := dk.LastOutcome.Matches(core.OutcomeLanded)
		if !casted || (casted && !dk.LastOutcome.Matches(core.OutcomeLanded)) {
			ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
			bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()

			if sim.CurrentTime+dk.SpellGCD() > ffExpiresAt || sim.CurrentTime+dk.SpellGCD() > bpExpiresAt {
				dk.RotationActionCallback_FrostSubBlood_RecoverFromPestiMiss(sim, target, s)
				return casted
			} else {
				s.ConditionalAdvance(casted && advance)
				return casted
			}
		} else {
			s.ConditionalAdvance(casted && advance)
			return casted
		}
	}
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Main_Pesti(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	ffActive := dk.FrostFeverDisease[target.Index].IsActive()
	bpActive := dk.BloodPlagueDisease[target.Index].IsActive()

	if !ffActive || !bpActive {
		dk.RotationActionCallback_FrostSubBlood_RecoverFromPestiMiss(sim, target, s)
		return casted
	} else {
		casted = dk.CastPestilence(sim, target)
		if !casted || (casted && !dk.LastOutcome.Matches(core.OutcomeLanded)) {
			ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
			bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()

			if sim.CurrentTime+dk.SpellGCD() > ffExpiresAt || sim.CurrentTime+dk.SpellGCD() > bpExpiresAt {
				dk.RotationActionCallback_FrostSubBlood_RecoverFromPestiMiss(sim, target, s)
				return casted
			} else {

				return casted
			}
		} else {
			return casted
		}
	}
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Main_FS_Star(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false
	if dk.PercentRunicPower() >= 0.95 || (dk.KillingMachineAura.IsActive() && dk.CurrentRunicPower() >= 32.0) {
		casted = dk.CastFrostStrike(sim, target)
	} else if dk.RimeAura.IsActive() {
		casted = dk.CastHowlingBlast(sim, target)
	} else if dk.CurrentRunicPower() >= 32.0 {
		casted = dk.CastFrostStrike(sim, target)
		if !casted {
			casted = dk.CastHornOfWinter(sim, target)
		}
	} else {
		casted = dk.CastHornOfWinter(sim, target)
	}

	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_Opener_FS_Star(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false
	if dk.PercentRunicPower() >= 0.95 || (dk.KillingMachineAura.IsActive() && dk.CurrentRunicPower() >= 32.0) {
		casted = dk.CastFrostStrike(sim, target)
		s.Advance()
	} else if dk.RimeAura.IsActive() {
		casted = dk.CastHowlingBlast(sim, target)
		s.ConditionalAdvance(casted)
	} else if dk.CurrentRunicPower() >= 32.0 {
		casted = dk.CastFrostStrike(sim, target)
		if !casted {
			casted = dk.CastHornOfWinter(sim, target)
		}
		s.Advance()
	} else {
		casted = false
		s.Advance()
	}

	return casted
}
