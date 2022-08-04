package dps

import (
	"math"
	"time"

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
		NewAction(dk.RotationActionCallback_Pesti_Opener).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_RD).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Pesti_Opener).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_FS)

	dk.Main.
		NewAction(dk.RotationActionCallback_FrostSubBloodPrioRotation)
}

func (dk *DpsDeathknight) setupFrostSubBloodNoERWOpener() {
	dk.Opener.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_UA).
		NewAction(dk.RotationActionCallback_Pesti_Opener).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_RD).
		NewAction(dk.RotationActionCallback_FS_Star_SubBlood_Advance).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FS_Star_SubBlood_Advance).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_Pesti_Opener).
		NewAction(dk.RotationActionCallback_FS_Star_SubBlood_Advance).
		NewAction(dk.RotationActionCallback_FS_Star_SubBlood_Advance)

	dk.Main.
		NewAction(dk.RotationActionCallback_FrostSubBloodPrioRotation)
}

func (dk *DpsDeathknight) setupFrostSubUnholyERWOpener() {
	dk.Opener.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_Pesti_Opener).
		NewAction(dk.RotationActionCallback_UA).
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_FS_Star).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_Pesti_Opener).
		NewAction(dk.RotationActionCallback_FS)

	dk.Main.
		NewAction(dk.RotationActionCallback_FrostSubUnholyStep1)
}

func (dk *DpsDeathknight) setupFrostSubUnholyNoERWOpener() {
	dk.Opener.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_UA).
		NewAction(dk.RotationActionCallback_Pesti_Opener).
		NewAction(dk.RotationActionCallback_FS_Star).
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_Pesti_Opener).
		NewAction(dk.RotationActionCallback_FS)

	dk.Main.
		NewAction(dk.RotationActionCallback_FrostSubUnholyStep1)
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholyStep1(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.Main.Clear().
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_Pesti_SubUnholy).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_FS_Star).
		NewAction(dk.RotationActionCallback_FS_Star).
		NewAction(dk.RotationActionCallback_FrostSubUnholyStep2)
	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholyStep2(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.Main.Clear().
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_UA_SubUnholy).
		NewAction(dk.RotationActionCallback_BT_SubUnholy).
		NewAction(dk.RotationActionCallback_Obli_SubUnholy).
		NewAction(dk.RotationActionCallback_FS_Star).
		NewAction(dk.RotationActionCallback_FS_Star).
		NewAction(dk.RotationActionCallback_FrostSubUnholyStep1)
	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholyRecoverFromPestiMiss_BS(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.Main.Clear().
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FS_Star).
		NewAction(dk.RotationActionCallback_FrostSubUnholyStep1)
	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholyRecoverFromPestiMiss(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.Main.Clear().
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FS_Star).
		NewAction(dk.RotationActionCallback_FrostSubUnholyStep1)
	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBloodRecoverFromPestiMiss(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	fr := &dk.fr

	if fr.lastSpell == dk.BloodStrike {
		dk.Opener.Clear().
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FS).
			NewAction(dk.RotationActionCallback_IT).
			NewAction(dk.RotationActionCallback_PS).
			NewAction(dk.RotationActionCallback_Obli).
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_Pesti_Opener).
			NewAction(dk.RotationActionCallback_FrostSubBloodPrioRotation)
	} else {
		dk.Opener.Clear().
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FS).
			NewAction(dk.RotationActionCallback_IT).
			NewAction(dk.RotationActionCallback_PS).
			NewAction(dk.RotationActionCallback_Obli).
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_Pesti_Opener).
			NewAction(dk.RotationActionCallback_FrostSubBloodPrioRotation)
	}

	//dk.WaitUntil(sim, sim.CurrentTime)
	return false
}

func (dk *DpsDeathknight) RotationActionCallback_Obli_SubUnholy(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	if dk.CanObliterate(sim) {
		runeCost := dk.OptimalRuneCost(core.RuneCost(dk.Obliterate.DefaultCast.Cost))

		if runeCost.Death() == 0 {
			casted = dk.CastObliterate(sim, target)
			s.ConditionalAdvance(casted)
		} else {
			ffActive := dk.FrostFeverDisease[target.Index].IsActive()
			bpActive := dk.BloodPlagueDisease[target.Index].IsActive()
			ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
			bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()

			if ffActive && bpActive {
				if sim.CurrentTime+(10*time.Second-dk.CurrentBloodRuneGrace(sim)) < core.MinDuration(ffExpiresAt, bpExpiresAt) {
					casted = dk.CastObliterate(sim, target)
					s.ConditionalAdvance(casted)
				} else {
					s.Advance()
				}
			} else {
				s.Advance()
			}
		}
	}

	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_Pesti_SubUnholy(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	ffActive := dk.FrostFeverDisease[target.Index].IsActive()
	bpActive := dk.BloodPlagueDisease[target.Index].IsActive()

	if !ffActive || !bpActive {
		return dk.RotationActionCallback_FrostSubUnholyRecoverFromPestiMiss(sim, target, s)
	} else {
		casted := dk.CastPestilence(sim, target)
		if !casted || (casted && !dk.LastCastOutcome.Matches(core.OutcomeLanded)) {
			ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
			bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()

			if sim.CurrentTime+dk.SpellGCD() > ffExpiresAt || sim.CurrentTime+dk.SpellGCD() > bpExpiresAt {
				return dk.RotationActionCallback_FrostSubUnholyRecoverFromPestiMiss(sim, target, s)
			} else {
				s.ConditionalAdvance(false)
				return casted
			}
		} else {
			s.ConditionalAdvance(casted)
			return casted
		}
	}
}

func (dk *DpsDeathknight) RotationActionCallback_Pesti_Opener(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	fr := &dk.fr

	casted := false
	advance := true

	ffActive := dk.FrostFeverDisease[target.Index].IsActive()
	bpActive := dk.BloodPlagueDisease[target.Index].IsActive()

	if ffActive && bpActive {
		casted = dk.CastPestilence(sim, target)
		advance = dk.LastCastOutcome.Matches(core.OutcomeLanded)
		if casted {
			fr.lastSpell = dk.Pestilence
		}
	} else {
		return dk.RotationActionCallback_FrostSubBloodRecoverFromPestiMiss(sim, target, s)
		// TODO: I need to deal with this in the opener
	}

	s.ConditionalAdvance(casted && advance)
	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBloodPrioRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	fr := &dk.fr

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
	} else if fr.nextSpell == dk.BloodStrike {
		casted = dk.FrostRotationCast(sim, target, dk.BloodStrike)
		if casted && dk.LastCastOutcome.Matches(core.OutcomeLanded) {
			fr.nextSpell = dk.Pestilence
		}
	} else if fr.nextSpell == dk.Pestilence {
		if (fr.lastSpell != dk.FrostStrike) || (fr.lastSpell != dk.HowlingBlast) || (fr.lastSpell != dk.HornOfWinter) {
			if sim.CurrentTime+2*gcd < fbExpireAt {
				casted = dk.RotationActionCallback_FS_Star_SubBlood(sim, target, s)
			}
		}

		if !casted {
			casted = dk.FrostRotationCast(sim, target, dk.Pestilence)
			if casted && dk.LastCastOutcome.Matches(core.OutcomeLanded) {
				fr.nextSpell = nil
			}
		}
	} else {
		if ffActive {
			if bpActive {
				if currGCDCastsTillExpire > 2 {
					if fr.lastSpell == dk.Obliterate {
						if dk.KillingMachineAura.IsActive() && sim.CurrentTime+gcd < spentFrostRuneAt && sim.CurrentTime+gcd < spentUnholyRuneAt &&
							frostGracePeriod-gcd > 0 && unholyGracePeriod-gcd > 0 {
							casted = dk.FrostRotationCast(sim, target, dk.FrostStrike)
						} else if currFrostRunes > 0 && currUnholyRunes > 0 {
							casted = dk.FrostRotationCast(sim, target, dk.Obliterate)
						}
					} else if currFrostRunes > 0 && currUnholyRunes > 0 {
						casted = dk.FrostRotationCast(sim, target, dk.Obliterate)
					}

					if !casted {
						// TODO: improve this, it breaks some runes
						if currBloodRunes > 0 || currDeathRunes > 0 {
							refreshingOverlapsOblitRunes := sim.CurrentTime+2*gcd > oblitRunesAt
							canRefreshAfterOblits := fbExpireAt > oblitRunesAt+2*gcd
							if refreshingOverlapsOblitRunes && canRefreshAfterOblits {
								if currBloodRunes+currDeathRunes == 2 {
									casted = dk.FrostRotationCast(sim, target, dk.BloodStrike)
									if casted && !dk.LastCastOutcome.Matches(core.OutcomeLanded) {
										fr.nextSpell = dk.BloodStrike
									} else {
										fr.nextSpell = dk.Pestilence
									}
								} else {
									casted = dk.FrostRotationCast(sim, target, dk.Obliterate)
								}
							} else {
								if dk.KillingMachineAura.IsActive() {
									casted = dk.RotationActionCallback_FS_Star_SubBlood(sim, target, s)
								} else {
									if dk.CanUnbreakableArmor(sim) {
										if sim.CurrentTime+gcd > spentBloodRuneAt-gcd {
											casted = dk.FrostRotationCast(sim, target, dk.UnbreakableArmor)
											if casted {
												casted = dk.RotationActionCallback_FS_Star_SubBlood(sim, target, s)
											}
											fr.nextSpell = dk.Pestilence
										}
									} else {
										if sim.CurrentTime+gcd > spentBloodRuneAt-gcd {
											casted = dk.FrostRotationCast(sim, target, dk.BloodStrike)
											if casted && !dk.LastCastOutcome.Matches(core.OutcomeLanded) {
												fr.nextSpell = dk.BloodStrike
											} else {
												fr.nextSpell = dk.Pestilence
											}
										} else {
											casted = dk.RotationActionCallback_FS_Star_SubBlood(sim, target, s)
										}
									}
								}
							}
						} else if sim.CurrentTime+gcd < spentFrostRuneAt && sim.CurrentTime+gcd < spentUnholyRuneAt {
							casted = dk.RotationActionCallback_FS_Star_SubBlood(sim, target, s)
						}
					}
				} else if currGCDCastsTillExpire > 1 {
					if sim.CurrentTime+gcd > spentBloodRuneAt-gcd {
						if dk.CanUnbreakableArmor(sim) {
							casted = dk.FrostRotationCast(sim, target, dk.UnbreakableArmor)
							if casted {
								casted = dk.RotationActionCallback_FS_Star_SubBlood(sim, target, s)
							} else {
								fr.nextSpell = dk.Pestilence
							}
						} else {
							casted = dk.FrostRotationCast(sim, target, dk.BloodStrike)
							if casted && !dk.LastCastOutcome.Matches(core.OutcomeLanded) {
								fr.nextSpell = dk.BloodStrike
							}
						}
					}

					if !casted {
						casted = dk.RotationActionCallback_FS_Star_SubBlood(sim, target, s)
					}
				} else {
					casted = dk.FrostRotationCast(sim, target, dk.Pestilence)
					if !casted {
						casted = dk.RotationActionCallback_FS_Star_SubBlood(sim, target, s)
					}
				}
			} else {
				casted = dk.FrostRotationCast(sim, target, dk.PlagueStrike)
				if !casted {
					casted = dk.RotationActionCallback_FS_Star_SubBlood(sim, target, s)
				}
			}
		} else {
			casted = dk.FrostRotationCast(sim, target, dk.IcyTouch)
			if !casted {
				casted = dk.RotationActionCallback_FS_Star_SubBlood(sim, target, s)
			}
		}
	}
	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_UA_SubUnholy(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	if dk.CanUnbreakableArmor(sim) {
		casted := dk.UnbreakableArmor.Cast(sim, target)
		dk.WaitUntil(sim, sim.CurrentTime)
		s.ConditionalAdvance(casted)
		return casted
	}

	s.Advance()
	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

func (dk *DpsDeathknight) RotationActionCallback_BT_SubUnholy(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	if dk.CanBloodTap(sim) {
		casted := dk.BloodTap.Cast(sim, target)
		dk.WaitUntil(sim, sim.CurrentTime)
		s.ConditionalAdvance(casted)
		return casted
	}

	s.Advance()
	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

func (dk *DpsDeathknight) RotationActionCallback_UA_Frost(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := dk.CastUnbreakableArmor(sim, target)
	if casted {
		dk.fr.lastSpell = dk.UnbreakableArmor
	}
	dk.WaitUntil(sim, sim.CurrentTime)
	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_BT_Frost(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := dk.CastBloodTap(sim, target)
	if casted {
		dk.fr.lastSpell = dk.BloodTap
	}
	dk.WaitUntil(sim, sim.CurrentTime)
	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_HB_Ghoul_RimeCheck(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false
	if dk.RimeAura.IsActive() {
		casted = dk.CastHowlingBlast(sim, target)
		if casted {
			dk.fr.lastSpell = dk.HowlingBlast
		}
	} else {
		casted = dk.CastRaiseDead(sim, target)
		if casted {
			dk.fr.lastSpell = dk.RaiseDead
		}
	}

	s.ConditionalAdvance(casted)
	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_FS_Star(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
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

	s.Advance()
	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_FS_Star_SubBlood(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false
	if dk.PercentRunicPower() >= 0.95 || (dk.KillingMachineAura.IsActive() && dk.CurrentRunicPower() >= 32.0) {
		casted = dk.FrostRotationCast(sim, target, dk.FrostStrike)
	} else if dk.RimeAura.IsActive() {
		casted = dk.FrostRotationCast(sim, target, dk.HowlingBlast)
	} else if dk.CurrentRunicPower() >= 32.0 {
		casted = dk.FrostRotationCast(sim, target, dk.FrostStrike)
		if !casted {
			casted = dk.FrostRotationCast(sim, target, dk.HornOfWinter)
		}
	} else {
		casted = false
	}

	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_FS_Star_SubBlood_Advance(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false
	if dk.PercentRunicPower() >= 0.95 || (dk.KillingMachineAura.IsActive() && dk.CurrentRunicPower() >= 32.0) {
		casted = dk.FrostRotationCast(sim, target, dk.FrostStrike)
		s.Advance()
	} else if dk.RimeAura.IsActive() {
		casted = dk.FrostRotationCast(sim, target, dk.HowlingBlast)
		s.ConditionalAdvance(casted)
	} else if dk.CurrentRunicPower() >= 32.0 {
		casted = dk.FrostRotationCast(sim, target, dk.FrostStrike)
		if !casted {
			casted = dk.FrostRotationCast(sim, target, dk.HornOfWinter)
		}
		s.Advance()
	} else {
		casted = false
		s.Advance()
	}

	return casted
}
