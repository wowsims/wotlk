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
		NewAction(dk.RotationActionCallback_Pesti).
		NewAction(dk.RotationActionCallback_ERW).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_HB_Ghoul_RimeCheck).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_Pesti).
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
		NewAction(dk.RotationActionCallback_Pesti).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_RD).
		NewAction(dk.RotationActionCallback_FS_HB_Advance).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FS_HB_Advance).
		NewAction(dk.RotationActionCallback_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_Pesti).
		NewAction(dk.RotationActionCallback_FS_HB_Advance).
		NewAction(dk.RotationActionCallback_FS_HB_Advance)

	dk.Main.
		NewAction(dk.RotationActionCallback_FrostSubBloodPrioRotation)
}

func (dk *DpsDeathknight) setupFrostSubUnholyOpener() {
	dk.Opener.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_Pesti).
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
		NewAction(dk.RotationActionCallback_Pesti).
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
		if !dk.LastCastOutcome.Matches(core.OutcomeLanded) {
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

func (dk *DpsDeathknight) RotationActionCallback_FrostSubBloodPrioRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	fr := &dk.fr

	gcd := dk.SpellGCD()
	ffActive := dk.FrostFeverDisease[target.Index].IsActive()
	bpActive := dk.BloodPlagueDisease[target.Index].IsActive()
	ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
	bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()
	frostGracePeriod := dk.CurrentFrostRuneGrace(sim)
	unholyGracePeriod := dk.CurrentUnholyRuneGrace(sim)
	currGCDCastsTillExpire := math.Floor(float64(core.MinDuration(ffExpiresAt, bpExpiresAt)-sim.CurrentTime) / float64(gcd))
	currFrostRunes := dk.CurrentFrostRunes()
	currUnholyRunes := dk.CurrentUnholyRunes()
	spentFrostRuneAt := dk.SpentFrostRuneReadyAt(sim)
	spentUnholyRuneAt := dk.SpentUnholyRuneReadyAt(sim)

	if dk.ShouldHornOfWinter(sim) {
		casted = dk.CastHornOfWinter(sim, target)
	} else if fr.nextSpell == dk.Pestilence {
		casted = dk.FrostRotationCast(sim, target, dk.Pestilence)
		if dk.LastCastOutcome.Matches(core.OutcomeLanded) {
			fr.nextSpell = nil
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
					} else if dk.CurrentFrostRunes() > 0 && dk.CurrentUnholyRunes() > 0 {
						casted = dk.FrostRotationCast(sim, target, dk.Obliterate)
					}

					if !casted {
						if dk.CurrentBloodRunes() > 0 || dk.CurrentDeathRunes() > 0 {
							if dk.KillingMachineAura.IsActive() {
								casted = dk.RotationActionCallback_FS_HB(sim, target, s)
							} else {
								if dk.CanUnbreakableArmor(sim) {
									casted = dk.FrostRotationCast(sim, target, dk.UnbreakableArmor)
									if casted {
										casted = dk.RotationActionCallback_FS_HB(sim, target, s)
									}
								} else {
									casted = dk.FrostRotationCast(sim, target, dk.BloodStrike)
								}
								fr.nextSpell = dk.Pestilence
							}
						} else if sim.CurrentTime+gcd < dk.SpentFrostRuneReadyAt(sim) && sim.CurrentTime+gcd < dk.SpentUnholyRuneReadyAt(sim) {
							casted = dk.RotationActionCallback_FS_HB(sim, target, s)
						}
					}
				} else if currGCDCastsTillExpire > 1 {
					if dk.CanUnbreakableArmor(sim) {
						casted = dk.FrostRotationCast(sim, target, dk.UnbreakableArmor)
						if casted {
							casted = dk.RotationActionCallback_FS_HB(sim, target, s)
						} else {
							fr.nextSpell = dk.Pestilence
						}
					} else {
						casted = dk.FrostRotationCast(sim, target, dk.BloodStrike)
					}

					if !casted {
						casted = dk.RotationActionCallback_FS_HB(sim, target, s)
					}
				} else {
					casted = dk.FrostRotationCast(sim, target, dk.Pestilence)
					if !casted {
						casted = dk.RotationActionCallback_FS_HB(sim, target, s)
					}
				}
			} else {
				casted = dk.FrostRotationCast(sim, target, dk.PlagueStrike)
				if !casted {
					casted = dk.RotationActionCallback_FS_HB(sim, target, s)
				}
			}
		} else {
			casted = dk.FrostRotationCast(sim, target, dk.IcyTouch)
			if !casted {
				casted = dk.RotationActionCallback_FS_HB(sim, target, s)
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

	s.ConditionalAdvance(true)
	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_FS_HB(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false
	if dk.KillingMachineAura.IsActive() && !dk.RimeAura.IsActive() {
		casted = dk.FrostRotationCast(sim, target, dk.FrostStrike)
	} else if dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
		if dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() < 110 {
			casted = dk.FrostRotationCast(sim, target, dk.HowlingBlast)
		} else if dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() > 110 {
			casted = dk.FrostRotationCast(sim, target, dk.HowlingBlast)
		} else if !dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() > 110 {
			casted = dk.FrostRotationCast(sim, target, dk.FrostStrike)
		} else if !dk.CastCostPossible(sim, 0, 0, 1, 1) && dk.CurrentRunicPower() < 110 {
			casted = dk.FrostRotationCast(sim, target, dk.FrostStrike)
		}
	} else if !dk.KillingMachineAura.IsActive() && dk.RimeAura.IsActive() {
		if dk.CurrentRunicPower() < 110 {
			casted = dk.FrostRotationCast(sim, target, dk.HowlingBlast)
		} else {
			casted = dk.FrostRotationCast(sim, target, dk.FrostStrike)
		}
	} else {
		casted = dk.FrostRotationCast(sim, target, dk.FrostStrike)
		if !casted {
			casted = dk.FrostRotationCast(sim, target, dk.HornOfWinter)
		}
	}
	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_FS_HB_Advance(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := dk.RotationActionCallback_FS_HB(sim, target, s)
	s.Advance()
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
