package dps

/*
func (dk *DpsDeathknight) setupFrostSubUnholyERWOpener() {
	dk.setupUnbreakableArmorCooldowns()

	dk.Opener.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Pesti).
		NewAction(dk.RotationActionCallback_UA_Frost).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_Star).
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
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Pesti).
		NewAction(dk.RotationActionCallback_FS)

	dk.Main.
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Step1)
}

func (dk *DpsDeathknight) setupFrostSubUnholyNoERWOpener() {
	dk.setupUnbreakableArmorCooldowns()

	dk.Opener.
		NewAction(dk.RotationActionCallback_IT).
		NewAction(dk.RotationActionCallback_PS).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_BT).
		NewAction(dk.RotationActionCallback_UA_Frost).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Pesti).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_Star).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FS).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
		NewAction(dk.RotationActionCallback_BS).
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Pesti).
		NewAction(dk.RotationActionCallback_FS)

	dk.Main.
		NewAction(dk.RotationActionCallback_FrostSubUnholy_Step1)
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_Step1(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	gcd := dk.SpellGCD()
	ffActive := dk.FrostFeverDisease[target.Index].IsActive()
	bpActive := dk.BloodPlagueDisease[target.Index].IsActive()
	ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
	bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()
	fbExpireAt := core.MinDuration(ffExpiresAt, bpExpiresAt)
	bloodGracePeriod := dk.CurrentBloodRuneGrace(sim)
	frostGracePeriod := dk.CurrentFrostRuneGrace(sim)
	unholyGracePeriod := dk.CurrentUnholyRuneGrace(sim)
	currBloodRunes := dk.CurrentBloodRunes()
	currFrostRunes := dk.CurrentFrostRunes()
	currUnholyRunes := dk.CurrentUnholyRunes()
	currDeathRunes := dk.CurrentDeathRunes()
	spentBloodRuneAt := dk.BloodRuneReadyAt(sim)
	spentFrostRuneAt := dk.NormalSpentFrostRuneReadyAt(sim)
	spentUnholyRuneAt := dk.NormalSpentUnholyRuneReadyAt(sim)
	oblitRunesAt := core.MaxDuration(spentFrostRuneAt, spentUnholyRuneAt)

	if dk.ShouldHornOfWinter(sim) {
		casted = dk.HornOfWinter.Cast(sim, target)
	} else if dk.NextCast == dk.BloodStrike {
		if (dk.LastCast != dk.FrostStrike) && (dk.LastCast != dk.HowlingBlast) && (dk.LastCast != dk.HornOfWinter) && (dk.LastCast != dk.UnbreakableArmor) {
			if sim.CurrentTime+2*gcd < fbExpireAt && bloodGracePeriod-gcd > 0 {
				casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
			}
		}

		if !casted {
			casted = dk.BloodStrike.Cast(sim, target)
			if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
				dk.NextCast = nil
				dk.Main.Clear().NewAction(dk.RotationActionCallback_FrostSubUnholy_Step2)
				return casted
			}
		}
	} else if dk.NextCast == dk.Pestilence {
		casted = dk.Pestilence.Cast(sim, target)
		if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
			dk.NextCast = dk.BloodStrike
		}
	} else {
		if ffActive {
			if bpActive {
				if dk.LastCast == dk.Obliterate {
					if dk.KillingMachineAura.IsActive() && sim.CurrentTime+gcd < spentFrostRuneAt && sim.CurrentTime+gcd < spentUnholyRuneAt &&
						frostGracePeriod-gcd > 0 && unholyGracePeriod-gcd > 0 {
						casted = dk.FrostStrike.Cast(sim, target)
					}
				}

				if !casted {
					if currFrostRunes > 0 && currUnholyRunes > 0 {
						casted = dk.Obliterate.Cast(sim, target)
						if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
							dk.fr.oblitCount += 1
						}
					} else if currDeathRunes == 2 && dk.fr.uaCycle && dk.BloodTap.CanCast(sim) {
						casted = dk.Obliterate.Cast(sim, target)
						if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
							dk.fr.oblitCount += 1
							dk.fr.uaCycle = false
							dk.fr.delayUACycle = false
							dk.NextCast = dk.BloodTap
						}
					} else if currDeathRunes == 2 && dk.fr.uaCycle && !dk.BloodTap.CanCast(sim) && dk.BloodTap.CD.ReadyAt() < fbExpireAt {
						casted = dk.Obliterate.Cast(sim, target)
						if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
							dk.fr.oblitCount += 1
							dk.fr.uaCycle = false
							dk.fr.delayUACycle = false
							dk.NextCast = dk.BloodTap
						}
					} else if currDeathRunes == 2 && dk.fr.uaCycle && !dk.BloodTap.CanCast(sim) && dk.BloodTap.CD.ReadyAt() >= fbExpireAt {
						dk.fr.delayUACycle = true
					}
				}

				if casted && (dk.LastCast == dk.Obliterate) && !dk.LastOutcome.Matches(core.OutcomeLanded) {
					dk.NextCast = dk.Obliterate
					return casted
				}

				if !casted {
					// TODO: improve this, it breaks some runes
					if currBloodRunes > 0 || currDeathRunes > 0 {
						refreshingOverlapsOblitRunes := sim.CurrentTime+2*gcd > oblitRunesAt
						canRefreshAfterOblits := fbExpireAt > oblitRunesAt+2*gcd
						if refreshingOverlapsOblitRunes && canRefreshAfterOblits {
							if currBloodRunes+currDeathRunes == 2 {
								casted = dk.Pestilence.Cast(sim, target)
								if casted && !dk.LastOutcome.Matches(core.OutcomeLanded) {
									dk.NextCast = dk.Pestilence
								} else {
									dk.NextCast = dk.BloodStrike
									dk.fr.oblitCount = 0
								}
							} else {
								casted = dk.Obliterate.Cast(sim, target)
								dk.fr.oblitCount += core.TernaryInt32(casted, 1, 0)
							}
						} else {
							if dk.KillingMachineAura.IsActive() {
								if fbExpireAt > sim.CurrentTime+gcd {
									casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
								}
							} else {
								if sim.CurrentTime+gcd > spentBloodRuneAt-gcd && dk.fr.oblitCount >= 2 && dk.NextCast != dk.BloodStrike && !dk.fr.uaCycle {
									casted = dk.Pestilence.Cast(sim, target)
									if casted && !dk.LastOutcome.Matches(core.OutcomeLanded) {
										dk.NextCast = dk.Pestilence
									} else {
										dk.NextCast = dk.BloodStrike
										dk.fr.oblitCount = 0
									}
								}
							}
						}
					} else if sim.CurrentTime+gcd < spentFrostRuneAt && sim.CurrentTime+gcd < spentUnholyRuneAt {
						casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
					}
				}

				if !casted {
					if (dk.fr.oblitCount >= 2 && dk.NextCast != dk.BloodStrike && !dk.fr.uaCycle) || dk.fr.delayUACycle {
						if (sim.CurrentTime+gcd > spentBloodRuneAt-gcd) || (currBloodRunes+currDeathRunes >= 1) {
							casted = dk.Pestilence.Cast(sim, target)
							if casted && !dk.LastOutcome.Matches(core.OutcomeLanded) {
								dk.NextCast = dk.Pestilence
							} else if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
								dk.fr.oblitCount = 0
								dk.NextCast = dk.BloodStrike
							}
						}
					}
				}

				if !casted {
					if oblitRunesAt+0*time.Millisecond > sim.CurrentTime+gcd {
						casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
					}
				}
			} else {
				casted = dk.PlagueStrike.Cast(sim, target)
				if !casted {
					casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
				}
			}
		} else {
			casted = dk.IcyTouch.Cast(sim, target)
			if !casted {
				casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
			}
		}
	}
	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_Step2(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	gcd := dk.SpellGCD()
	ffActive := dk.FrostFeverDisease[target.Index].IsActive()
	bpActive := dk.BloodPlagueDisease[target.Index].IsActive()
	//ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
	//bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()
	//fbExpireAt := core.MinDuration(ffExpiresAt, bpExpiresAt)
	//bloodGracePeriod := dk.CurrentBloodRuneGrace(sim)
	frostGracePeriod := dk.CurrentFrostRuneGrace(sim)
	unholyGracePeriod := dk.CurrentUnholyRuneGrace(sim)
	//currBloodRunes := dk.CurrentBloodRunes()
	currFrostRunes := dk.CurrentFrostRunes()
	currUnholyRunes := dk.CurrentUnholyRunes()
	currDeathRunes := dk.CurrentDeathRunes()
	bloodRuneAt := dk.BloodRuneReadyAt(sim)
	spentFrostRuneAt := dk.NormalSpentFrostRuneReadyAt(sim)
	spentUnholyRuneAt := dk.NormalSpentUnholyRuneReadyAt(sim)
	oblitRunesAt := core.MaxDuration(spentFrostRuneAt, spentUnholyRuneAt)

	if dk.ShouldHornOfWinter(sim) {
		casted = dk.HornOfWinter.Cast(sim, target)
	} else if dk.NextCast == dk.BloodTap {
		casted = dk.BloodTap.Cast(sim, target)
		if casted {
			casted = dk.UnbreakableArmor.Cast(sim, target)
			dk.castAllMajorCooldowns(sim)
			dk.WaitUntil(sim, sim.CurrentTime)
		}

		dk.fr.oblitCount = 0
		dk.NextCast = nil
		dk.Main.Clear().NewAction(dk.RotationActionCallback_FrostSubUnholy_Step1)
	} else if dk.NextCast == dk.Obliterate {
		if dk.fr.oblitCount == 2 {
			casted = dk.Obliterate.Cast(sim, target)
			if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
				dk.NextCast = nil
				dk.fr.oblitCount = 0
				dk.Main.Clear().NewAction(dk.RotationActionCallback_FrostSubUnholy_Step1)
				return casted
			}
		} else {
			casted = dk.Obliterate.Cast(sim, target)
			if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
				dk.NextCast = nil
				dk.fr.oblitCount += 1
			}
		}
	} else {
		if ffActive {
			if bpActive {
				skipOblitCheck := false
				if dk.fr.oblitCount == 2 && dk.BloodTap.CanCast(sim) && dk.UnbreakableArmor.CanCast(sim) && currDeathRunes < 2 && bloodRuneAt > sim.CurrentTime+gcd {
					casted = dk.FrostStrike.Cast(sim, target)
					dk.NextCast = dk.BloodTap
				} else if dk.fr.oblitCount == 2 && dk.BloodTap.CanCast(sim) && dk.UnbreakableArmor.CanCast(sim) && currDeathRunes == 2 {
					casted = dk.Obliterate.Cast(sim, target)
					if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
						dk.fr.oblitCount += 1
						dk.NextCast = dk.BloodTap
					} else {
						skipOblitCheck = true
					}
				} else if dk.LastCast == dk.Obliterate {
					if dk.KillingMachineAura.IsActive() && sim.CurrentTime+gcd < spentFrostRuneAt && sim.CurrentTime+gcd < spentUnholyRuneAt &&
						frostGracePeriod-gcd > 0 && unholyGracePeriod-gcd > 0 {
						casted = dk.FrostStrike.Cast(sim, target)
					}
				}

				if !casted {
					if currFrostRunes > 0 && currUnholyRunes > 0 {
						casted = dk.Obliterate.Cast(sim, target)
						if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
							dk.fr.oblitCount += 1
						}
					} else if currDeathRunes == 2 {
						casted = dk.Obliterate.Cast(sim, target)
						if casted && dk.LastOutcome.Matches(core.OutcomeLanded) {
							dk.fr.oblitCount = 0
							dk.Main.Clear().NewAction(dk.RotationActionCallback_FrostSubUnholy_Step1)
							return casted
						}
					}
				}

				if casted && (dk.LastCast == dk.Obliterate) && !skipOblitCheck && !dk.LastOutcome.Matches(core.OutcomeLanded) {
					dk.NextCast = dk.Obliterate
					return casted
				}

				if !casted {
					if oblitRunesAt+0*time.Millisecond > sim.CurrentTime+gcd {
						casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
					}
				}
			} else {
				casted = dk.PlagueStrike.Cast(sim, target)
				if !casted {
					casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
				}
			}
		} else {
			casted = dk.IcyTouch.Cast(sim, target)
			if !casted {
				casted = dk.RotationActionCallback_FrostSubBlood_Main_FS_Star(sim, target, s)
			}
		}
	}
	return casted
	//s.Clear().
	//	NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
	//	NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
	//	NewAction(dk.RotationActionCallback_FS).
	//	NewAction(dk.RotationActionCallback_FrostSubUnholy_UA_BT).
	//	NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
	//	NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_Star).
	//	NewAction(dk.RotationActionCallback_FrostSubUnholy_FS_Star).
	//	NewAction(dk.RotationActionCallback_FrostSubUnholy_Step1)
	//return false
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_RecoverFromPestiMiss(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	if dk.LastCast == dk.BloodStrike {
		s.Clear().
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FS).
			NewAction(dk.RotationActionCallback_IT).
			NewAction(dk.RotationActionCallback_PS).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Pesti).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Step1)
	} else {
		s.Clear().
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FS).
			NewAction(dk.RotationActionCallback_IT).
			NewAction(dk.RotationActionCallback_PS).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Obli).
			NewAction(dk.RotationActionCallback_BS).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Pesti).
			NewAction(dk.RotationActionCallback_FrostSubUnholy_Step1)
	}

	return false
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_Obli(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false
	ffActive := dk.FrostFeverDisease[target.Index].IsActive()
	bpActive := dk.BloodPlagueDisease[target.Index].IsActive()

	if ffActive && bpActive {
		ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
		bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()

		if dk.Obliterate.CanCast(sim) {
			runeCost := dk.OptimalRuneCost(core.RuneCost(dk.Obliterate.DefaultCast.Cost))

			if runeCost.Death() == 0 {
				casted = dk.Obliterate.Cast(sim, target)
				s.ConditionalAdvance(casted)
			} else {
				if sim.CurrentTime+(10*time.Second-dk.CurrentBloodRuneGrace(sim)) < core.MinDuration(ffExpiresAt, bpExpiresAt) {
					casted = dk.Obliterate.Cast(sim, target)
					s.ConditionalAdvance(casted)
				} else {
					s.Advance()
				}

			}
		}
	} else {
		s.Advance()
	}

	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_Pesti(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false

	ffActive := dk.FrostFeverDisease[target.Index].IsActive()
	bpActive := dk.BloodPlagueDisease[target.Index].IsActive()

	if !ffActive || !bpActive {
		dk.RotationActionCallback_FrostSubUnholy_RecoverFromPestiMiss(sim, target, s)
		return casted
	} else {
		casted = dk.Pestilence.Cast(sim, target)
		advance := dk.LastOutcome.Matches(core.OutcomeLanded)
		if !casted || (casted && !dk.LastOutcome.Matches(core.OutcomeLanded)) {
			ffExpiresAt := dk.FrostFeverDisease[target.Index].ExpiresAt()
			bpExpiresAt := dk.BloodPlagueDisease[target.Index].ExpiresAt()

			if sim.CurrentTime+dk.SpellGCD() > ffExpiresAt || sim.CurrentTime+dk.SpellGCD() > bpExpiresAt {
				dk.RotationActionCallback_FrostSubUnholy_RecoverFromPestiMiss(sim, target, s)
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

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_FS_Star(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	casted := false
	if dk.PercentRunicPower() >= 0.95 || (dk.KillingMachineAura.IsActive() && dk.CurrentRunicPower() >= 32.0) {
		casted = dk.FrostStrike.Cast(sim, target)
	} else if dk.RimeAura.IsActive() {
		casted = dk.HowlingBlast.Cast(sim, target)
	} else if dk.CurrentRunicPower() >= 32.0 {
		casted = dk.FrostStrike.Cast(sim, target)
		if !casted {
			casted = dk.HornOfWinter.Cast(sim, target)
		}
	} else {
		casted = false
	}

	s.Advance()
	return casted
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_UA_BT(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	if dk.UnbreakableArmor.CanCast(sim) && dk.BloodTap.CanCast(sim) {
		casted := dk.UnbreakableArmor.Cast(sim, target)
		casted = casted && dk.BloodTap.Cast(sim, target)
		dk.WaitUntil(sim, sim.CurrentTime)
		s.ConditionalAdvance(casted)
		return casted
	}

	s.Advance()
	dk.WaitUntil(sim, sim.CurrentTime)
	return true
}

func (dk *DpsDeathknight) RotationActionCallback_FrostSubUnholy_CancelBT(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	dk.BloodTapAura.Deactivate(sim)
	dk.WaitUntil(sim, sim.CurrentTime)
	s.Advance()
	return true
}
*/
