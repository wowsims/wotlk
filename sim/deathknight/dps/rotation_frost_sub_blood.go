package dps

import (
	"math"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) setupFrostSubBloodERWOpener() {
	dk.Opener.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_UA).
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
		NewAction(dk.RotationActionCallback_FS)

	dk.Main.
		NewAction(dk.RotationActionCallback_FrostSubBlood_PrioRotation)
}

func (dk *DpsDeathknight) setupFrostSubBloodNoERWOpener() {
	dk.Opener.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_UA).
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

	return casted
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

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBlood_PrioRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	gcd := dk.SpellGCD()
	ffActive := dk.FrostFeverDisease[target.Index].IsActive()
	bpActive := dk.BloodPlagueDisease[target.Index].IsActive()
	ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
	bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()
	fbExpireAt := core.MinDuration(ffExpiresAt, bpExpiresAt)
	frostGracePeriod := dk.CurrentFrostRuneGrace(sim)
	unholyGracePeriod := dk.CurrentUnholyRuneGrace(sim)
	currGCDCastsTillExpire := math.Floor(float64(fbExpireAt-sim.CurrentTime) / float64(gcd))
	currBloodRunes := dk.CurrentBloodRunes()
	currFrostRunes := dk.CurrentFrostRunes()
	currUnholyRunes := dk.CurrentUnholyRunes()
	currDeathRunes := dk.CurrentDeathRunes()
	spentBloodRuneAt := dk.SpentBloodRuneReadyAt(sim)
	spentFrostRuneAt := dk.NormalSpentFrostRuneReadyAt(sim)
	spentUnholyRuneAt := dk.NormalSpentUnholyRuneReadyAt(sim)
	oblitRunesAt := core.MaxDuration(spentFrostRuneAt, spentUnholyRuneAt)

	if dk.ShouldHornOfWinter(sim) {
		casted = dk.CastHornOfWinter(sim, target)
	} else if dk.NextCast == dk.BloodStrike {
		casted = dk.CastBloodStrike(sim, target)
		if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
			dk.NextCast = dk.Pestilence
		}
	} else if dk.NextCast == dk.Pestilence {
		if (dk.LastCast != dk.FrostStrike) || (dk.LastCast != dk.HowlingBlast) || (dk.LastCast != dk.HornOfWinter) {
			if sim.CurrentTime+2*gcd < fbExpireAt {
				casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
			}
		}

		if !casted {
			casted = dk.CastPestilence(sim, target)
			if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
				dk.NextCast = nil
			}
		}
	} else {
		if ffActive {
			if bpActive {
				if currGCDCastsTillExpire > 2 {
					if dk.LastCast == dk.Obliterate {
						if dk.KillingMachineAura.IsActive() && sim.CurrentTime+gcd < spentFrostRuneAt && sim.CurrentTime+gcd < spentUnholyRuneAt &&
							frostGracePeriod-gcd > 0 && unholyGracePeriod-gcd > 0 {
							casted = dk.CastFrostStrike(sim, target)
						} else if currFrostRunes > 0 && currUnholyRunes > 0 {
							casted = dk.CastObliterate(sim, target)
						}
					} else if currFrostRunes > 0 && currUnholyRunes > 0 {
						casted = dk.CastObliterate(sim, target)
					}

					if !casted {
						// TODO: improve this, it breaks some runes
						if currBloodRunes > 0 || currDeathRunes > 0 {
							refreshingOverlapsOblitRunes := sim.CurrentTime+2*gcd > oblitRunesAt
							canRefreshAfterOblits := fbExpireAt > oblitRunesAt+2*gcd
							if refreshingOverlapsOblitRunes && canRefreshAfterOblits {
								if currBloodRunes+currDeathRunes == 2 {
									casted = dk.CastBloodStrike(sim, target)
									if casted && !dk.LastOutcome.Matches(core.OutcomeLanded) {
										dk.NextCast = dk.BloodStrike
									} else {
										dk.NextCast = dk.Pestilence
									}
								} else {
									casted = dk.CastObliterate(sim, target)
								}
							} else {
								if dk.KillingMachineAura.IsActive() {
									casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
								} else {
									if dk.CanUnbreakableArmor(sim) {
										if sim.CurrentTime+gcd > spentBloodRuneAt-gcd {
											casted = dk.CastUnbreakableArmor(sim, target)
											if casted {
												casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
											}
											dk.NextCast = dk.Pestilence
										}
									} else {
										if sim.CurrentTime+gcd > spentBloodRuneAt-gcd {
											casted = dk.CastBloodStrike(sim, target)
											if casted && !dk.LastOutcome.Matches(core.OutcomeLanded) {
												dk.NextCast = dk.BloodStrike
											} else {
												dk.NextCast = dk.Pestilence
											}
										} else {
											casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
										}
									}
								}
							}
						} else if sim.CurrentTime+gcd < spentFrostRuneAt && sim.CurrentTime+gcd < spentUnholyRuneAt {
							casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
						}
					}
				} else if currGCDCastsTillExpire > 1 {
					if sim.CurrentTime+gcd > spentBloodRuneAt-gcd {
						if dk.CanUnbreakableArmor(sim) {
							casted = dk.CastUnbreakableArmor(sim, target)
							if casted {
								casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
							} else {
								dk.NextCast = dk.Pestilence
							}
						} else {
							casted = dk.CastBloodStrike(sim, target)
							if casted && !dk.LastOutcome.Matches(core.OutcomeLanded) {
								dk.NextCast = dk.BloodStrike
							}
						}
					}

					if !casted {
						casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
					}
				} else {
					casted = dk.CastPestilence(sim, target)
					if !casted {
						casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
					}
				}
			} else {
				casted = dk.CastPlagueStrike(sim, target)
				if !casted {
					casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
				}
			}
		} else {
			casted = dk.CastIcyTouch(sim, target)
			if !casted {
				casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
			}
		}
	}
	return casted
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
		casted = false
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
