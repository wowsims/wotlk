package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) setupUnholyRotations() {
	if dk.Rotation.BloodTap == proto.Deathknight_Rotation_GhoulFrenzy && !dk.Talents.GhoulFrenzy {
		dk.Rotation.BloodTap = proto.Deathknight_Rotation_IcyTouch
	}

	dk.Inputs.FuStrike = deathknight.FuStrike_ScourgeStrike
	if dk.Talents.Annihilation > 0 {
		dk.Inputs.FuStrike = deathknight.FuStrike_Obliterate
	}

	if dk.Talents.SummonGargoyle && dk.Rotation.UseGargoyle {
		dk.setupGargoyleCooldowns()
	}

	dk.RotationSequence.Clear().
		NewAction(dk.RotationActionCallback_MindFreezeFiller).
		NewAction(dk.getFirstDiseaseAction()).
		NewAction(dk.getSecondDiseaseAction()).
		NewAction(dk.getBloodRuneAction(true))

	if dk.Talents.ScourgeStrike {
		if dk.ur.sigil == Sigil_Virulence {
			dk.RotationSequence.
				NewAction(dk.RotationActionCallback_SS).
				NewAction(dk.RotationActionCallback_BS)
		} else if dk.ur.sigil == Sigil_HangedMan {
			dk.RotationSequence.
				NewAction(dk.RotationActionCallback_SS).
				NewAction(dk.RotationActionCallback_BS).
				NewAction(dk.RotationActionCallback_DC).
				NewAction(dk.RotationActionCallback_ERW).
				NewAction(dk.RotationActionCallback_SS).
				NewAction(dk.RotationActionCallback_BS).
				NewAction(dk.RotationActionCallback_SS).
				NewAction(dk.RotationActionCallback_BS)
		}
	}

	if dk.Rotation.UseDeathAndDecay || (!dk.Talents.ScourgeStrike && dk.Talents.Annihilation == 0) {
		if !dk.Talents.ScourgeStrike || dk.ur.sigil == Sigil_Other {
			dk.RotationSequence.
				NewAction(dk.RotationActionCallback_DND)
		}
		dk.RotationSequence.
			NewAction(dk.RotationActionCallback_UnholyDndRotation)
	} else {
		dk.RotationSequence.NewAction(dk.RotationActionCallback_UnholySsRotation)
	}
}

func (dk *DpsDeathknight) setupWeaponSwap() {
	if !dk.ItemSwap.IsEnabled() {
		return
	}

	if mh := dk.ItemSwap.GetItem(proto.ItemSlot_ItemSlotMainHand); mh != nil {
		if mh.Enchant.EffectID == 3790 {
			dk.ur.mhSwap = WeaponSwap_BlackMagic
		} else if mh.Enchant.EffectID == 3789 {
			dk.ur.mhSwap = WeaponSwap_Berserking
		} else if mh.Enchant.EffectID == 3368 {
			dk.ur.mhSwap = WeaponSwap_FallenCrusader
		}
	}

	if oh := dk.ItemSwap.GetItem(proto.ItemSlot_ItemSlotOffHand); oh != nil {
		if oh.Enchant.EffectID == 3790 {
			dk.ur.ohSwap = WeaponSwap_BlackMagic
		} else if oh.Enchant.EffectID == 3789 {
			dk.ur.ohSwap = WeaponSwap_Berserking
		} else if oh.Enchant.EffectID == 3368 {
			dk.ur.ohSwap = WeaponSwap_FallenCrusader
		}
	}

	if dk.ur.mhSwap != WeaponSwap_None || dk.ur.ohSwap != WeaponSwap_None {
		core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
			Label: "Weapon Swap Check",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if dk.GCD.TimeToReady(sim) >= time.Second && dk.Hardcast.Expires < sim.CurrentTime {
					if sim.Log != nil {
						sim.Log("Swap Check: %0.2f", dk.GCD.TimeToReady(sim).Seconds())
					}
					dk.weaponSwapCheck(sim)
				}
			},
		}))
	}
}

func (dk *DpsDeathknight) RotationActionCallback_UnholyDndRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.uhGargoyleCheck(sim, target, 100*time.Millisecond) {
		dk.uhAfterGargoyleSequence(sim)
		return sim.CurrentTime
	}

	if dk.CurrentRunicPower() > 100 && dk.GCD.IsReady(sim) && dk.DeathCoil.Cast(sim, target) {
		return -1
	}

	if dk.Talents.GhoulFrenzy && !dk.uhShouldWaitForDnD(sim, false, true, true) {
		if dk.uhGhoulFrenzyCheck(sim, target) {
			return sim.CurrentTime
		}
	}

	if dk.uhBloodTap(sim, target) {
		return sim.CurrentTime
	}

	if dk.uhEmpoweredRuneWeapon(sim, target) {
		return sim.CurrentTime
	}

	if dk.uhMindFreeze(sim, target) {
		return sim.CurrentTime
	}

	cast := false
	prioSs, prioBs := dk.bonusProcRotationChecks(sim)

	if dk.uhDiseaseCheck(sim, target, dk.DeathAndDecay, true, 1) {
		if prioSs {
			if dk.uhGargoyleCheck(sim, target, core.GCDDefault*2+50*time.Millisecond) {
				dk.uhAfterGargoyleSequence(sim)
				return sim.CurrentTime
			}
			cast = dk.ScourgeStrike.Cast(sim, target)
			if cast {
				dk.RotationSequence.Clear().
					NewAction(dk.RotationActionCallback_BS).
					NewAction(dk.RotationActionCallback_UnholyDndRotation)
			}
		} else {
			if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
				dk.uhAfterGargoyleSequence(sim)
				return sim.CurrentTime
			}
			cast = dk.DeathAndDecay.Cast(sim, target)
		}
	} else {
		if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+core.GCDDefault+50*time.Millisecond) {
			dk.uhAfterGargoyleSequence(sim)
			return sim.CurrentTime
		}
		dk.uhRecastDiseasesSequence(sim)
		return sim.CurrentTime
	}

	if !cast {
		if dk.uhDiseaseCheck(sim, target, dk.ScourgeStrike, true, 1) {
			if !dk.uhShouldWaitForDnD(sim, false, true, true) {
				if dk.Talents.ScourgeStrike && dk.ScourgeStrike.IsReady(sim) {
					if dk.uhGargoyleCheck(sim, target, core.GCDDefault+50*time.Millisecond) {
						dk.uhAfterGargoyleSequence(sim)
						return sim.CurrentTime
					}
					cast = dk.ScourgeStrike.Cast(sim, target)
				} else if dk.IcyTouch.CanCast(sim, nil) && dk.PlagueStrike.CanCast(sim, nil) {
					if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+core.GCDDefault+50*time.Millisecond) {
						dk.uhAfterGargoyleSequence(sim)
						return sim.CurrentTime
					}
					dk.uhRecastDiseasesSequence(sim)
					return sim.CurrentTime
				}
			}
		} else {
			dk.uhRecastDiseasesSequence(sim)
			return sim.CurrentTime
		}
		if !cast {
			if dk.shShouldSpreadDisease(sim) {
				if !dk.uhShouldWaitForDnD(sim, true, false, false) {
					if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
						dk.uhAfterGargoyleSequence(sim)
						return sim.CurrentTime
					}
					cast = dk.uhSpreadDiseases(sim, target, s)
				}
			} else {
				if !dk.uhShouldWaitForDnD(sim, true, false, false) {
					if dk.uhGargoyleCheck(sim, target, core.GCDDefault+50*time.Millisecond) {
						dk.uhAfterGargoyleSequence(sim)
						return sim.CurrentTime
					}
					if dk.desolationAuraCheck(sim) || prioBs {
						cast = dk.BloodStrike.Cast(sim, target)
					} else {
						cast = dk.BloodBoil.Cast(sim, target)
					}
				}
			}
			if !cast {
				if dk.uhDeathCoilCheck(sim) {
					if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
						dk.uhAfterGargoyleSequence(sim)
						return sim.CurrentTime
					}
					cast = dk.DeathCoil.Cast(sim, target)
				}
				if !cast {
					if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
						dk.uhAfterGargoyleSequence(sim)
						return sim.CurrentTime
					}

					if dk.HornOfWinter.CanCast(sim, nil) {
						dk.HornOfWinter.Cast(sim, target)
					}
				}
			}
		}
	}

	// Gargoyle cast needs to be checked more often then default rotation on gcd/resource gain checks
	if dk.Rotation.UseGargoyle && dk.SummonGargoyle.IsReady(sim) && dk.GCD.IsReady(sim) {
		return sim.CurrentTime + 100*time.Millisecond
	}

	return -1
}

func (dk *DpsDeathknight) RotationActionCallback_UnholySsRotation(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.uhGargoyleCheck(sim, target, 100*time.Millisecond) {
		dk.uhAfterGargoyleSequence(sim)
		return sim.CurrentTime
	}

	if dk.Talents.GhoulFrenzy {
		if dk.uhGhoulFrenzyCheck(sim, target) {
			return sim.CurrentTime
		}
	}

	if dk.uhBloodTap(sim, target) {
		return sim.CurrentTime
	}

	if dk.uhEmpoweredRuneWeapon(sim, target) {
		return sim.CurrentTime
	}

	if dk.uhMindFreeze(sim, target) {
		return sim.CurrentTime
	}

	casted := false
	fuStrike := dk.ScourgeStrike
	if dk.Inputs.FuStrike == deathknight.FuStrike_Obliterate {
		fuStrike = dk.Obliterate
	}
	if dk.uhDiseaseCheck(sim, target, fuStrike, true, 1) {
		if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
			dk.uhAfterGargoyleSequence(sim)
			return sim.CurrentTime
		}
		casted = fuStrike.Cast(sim, target)
	} else {
		if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()*2+50*time.Millisecond) {
			dk.uhAfterGargoyleSequence(sim)
			return sim.CurrentTime
		}
		dk.uhRecastDiseasesSequence(sim)
		return sim.CurrentTime
	}
	if !casted {
		if dk.shShouldSpreadDisease(sim) {
			if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
				dk.uhAfterGargoyleSequence(sim)
				return sim.CurrentTime
			}
			casted = dk.uhSpreadDiseases(sim, target, s)
		} else {
			if dk.uhDiseaseCheck(sim, target, dk.BloodStrike, true, 1) {
				if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
					dk.uhAfterGargoyleSequence(sim)
					return sim.CurrentTime
				}
				if dk.desolationAuraCheck(sim) {
					casted = dk.BloodStrike.Cast(sim, target)
				} else {
					casted = dk.BloodBoil.Cast(sim, target)
				}
			} else {
				if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()*2+50*time.Millisecond) {
					dk.uhAfterGargoyleSequence(sim)
					return sim.CurrentTime
				}
				dk.uhRecastDiseasesSequence(sim)
				return sim.CurrentTime
			}
		}
		if !casted {
			if dk.uhDeathCoilCheck(sim) {
				if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
					dk.uhAfterGargoyleSequence(sim)
					return sim.CurrentTime
				}
				casted = dk.DeathCoil.Cast(sim, target)
			}
			if !casted {
				if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
					dk.uhAfterGargoyleSequence(sim)
					return sim.CurrentTime
				}

				if dk.HornOfWinter.CanCast(sim, nil) {
					dk.HornOfWinter.Cast(sim, target)
				}
			}
		}
	}

	// Gargoyle cast needs to be checked more often then default rotation on gcd/resource gain checks
	if dk.Rotation.UseGargoyle && dk.SummonGargoyle.IsReady(sim) && dk.GCD.IsReady(sim) {
		return sim.CurrentTime + 100*time.Millisecond
	}

	return -1
}

func (dk *DpsDeathknight) uhAfterGargoyleSequence(sim *core.Simulation) {
	if dk.Rotation.UseEmpowerRuneWeapon && dk.EmpowerRuneWeapon.IsReady(sim) {
		dk.RotationSequence.Clear()

		if dk.BloodTapAura.IsActive() {
			dk.RotationSequence.NewAction(dk.RotationAction_CancelBT)
		}

		didErw := false
		if dk.Inputs.ArmyOfTheDeadType != proto.Deathknight_Rotation_DoNotUse && dk.ArmyOfTheDead.IsReady(sim) {
			// If not enough runes for aotd cast ERW
			if dk.CurrentBloodRunes() < 1 || dk.CurrentFrostRunes() < 1 || dk.CurrentUnholyRunes() < 1 {
				dk.RotationSequence.NewAction(dk.RotationActionCallback_ERW)
				didErw = true
			}
			dk.RotationSequence.
				NewAction(dk.RotationActionCallback_Haste_Snapshot).
				NewAction(dk.RotationActionCallback_AOTD)
		}

		if dk.Rotation.BlPresence == proto.Deathknight_Rotation_Blood && !dk.PresenceMatches(deathknight.BloodPresence) && (dk.Rotation.PreNerfedGargoyle || dk.Rotation.GargoylePresence == proto.Deathknight_Rotation_Blood) {
			if didErw || dk.CurrentBloodRunes() > 0 {
				dk.RotationSequence.NewAction(dk.RotationActionCallback_BP)
			} else if !didErw && !dk.Rotation.BtGhoulFrenzy && dk.BloodTap.IsReady(sim) {
				dk.RotationSequence.
					NewAction(dk.RotationActionCallback_BT).
					NewAction(dk.RotationActionCallback_BP).
					NewAction(dk.RotationAction_CancelBT)
			}
		}

		if dk.Rotation.UseDeathAndDecay || (!dk.Talents.ScourgeStrike && dk.Talents.Annihilation == 0) {
			dk.RotationSequence.NewAction(dk.RotationActionUH_ResetToDndMain)
		} else {
			dk.RotationSequence.NewAction(dk.RotationActionUH_ResetToSsMain)
		}
	} else if dk.Inputs.ArmyOfTheDeadType == proto.Deathknight_Rotation_AsMajorCd && dk.ArmyOfTheDead.IsReady(sim) {
		dk.RotationSequence.Clear()
		dk.RotationSequence.
			NewAction(dk.RotationActionCallback_Haste_Snapshot).
			NewAction(dk.RotationActionCallback_AOTD)

		if dk.Rotation.UseDeathAndDecay || (!dk.Talents.ScourgeStrike && dk.Talents.Annihilation == 0) {
			dk.RotationSequence.NewAction(dk.RotationActionUH_ResetToDndMain)
		} else {
			dk.RotationSequence.NewAction(dk.RotationActionUH_ResetToSsMain)
		}
	}
}

func (dk *DpsDeathknight) RotationActionCallback_Haste_Snapshot(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	dk.ur.gargoyleSnapshot.ActivateMajorCooldowns(sim)
	s.Advance()
	return sim.CurrentTime
}

func (dk *DpsDeathknight) uhGhoulFrenzySequence(sim *core.Simulation, bloodTap bool) {
	if bloodTap {
		dk.RotationSequence.Clear().
			NewAction(dk.RotationActionCallback_BT).
			NewAction(dk.RotationActionCallback_GF).
			NewAction(dk.RotationAction_CancelBT)
	} else {
		if dk.sr.ffFirst {
			dk.RotationSequence.Clear().
				NewAction(dk.RotationActionUH_IT_SetSync).
				NewAction(dk.RotationActionCallback_GF)
		} else {
			dk.RotationSequence.Clear().
				NewAction(dk.RotationActionCallback_GF).
				NewAction(dk.RotationActionUH_IT_SetSync)
		}
	}

	if dk.Rotation.UseDeathAndDecay || (!dk.Talents.ScourgeStrike && dk.Talents.Annihilation == 0) {
		dk.RotationSequence.NewAction(dk.RotationActionUH_ResetToDndMain)
	} else {
		dk.RotationSequence.NewAction(dk.RotationActionUH_ResetToSsMain)
	}
}

func (dk *DpsDeathknight) uhRecastDiseasesSequence(sim *core.Simulation) {
	dk.RotationSequence.Clear()

	// If we have glyph of Disease and both dots active try to refresh with pesti
	didPesti := false
	if dk.sr.hasGod {
		if dk.FrostFeverSpell.Dot(dk.CurrentTarget).IsActive() && dk.BloodPlagueSpell.Dot(dk.CurrentTarget).IsActive() {
			didPesti = true
			dk.RotationSequence.NewAction(dk.RotationActionCallback_Pesti_Custom)
		}
	}

	// If we did not pesti queue normal dot refresh
	if !didPesti {
		if dk.sr.ffFirst {
			dk.RotationSequence.
				NewAction(dk.RotationActionUH_FF_ClipCheck).
				NewAction(dk.RotationActionUH_IT_Custom).
				NewAction(dk.RotationActionUH_BP_ClipCheck).
				NewAction(dk.RotationActionUH_PS_Custom)
		} else {
			dk.RotationSequence.
				NewAction(dk.RotationActionUH_BP_ClipCheck).
				NewAction(dk.RotationActionUH_PS_Custom).
				NewAction(dk.RotationActionUH_FF_ClipCheck).
				NewAction(dk.RotationActionUH_IT_Custom)
		}
	}

	if dk.Rotation.UseDeathAndDecay || (!dk.Talents.ScourgeStrike && dk.Talents.Annihilation == 0) {
		dk.RotationSequence.NewAction(dk.RotationActionUH_ResetToDndMain)
	} else {
		dk.RotationSequence.NewAction(dk.RotationActionUH_ResetToSsMain)
	}
}

func (dk *DpsDeathknight) RotationActionCallback_MindFreezeFiller(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	// Use Mind Freeze off GCD to proc extra effects in the opener
	if dk.Talents.EndlessWinter == 2 && dk.SummonGargoyle.IsReady(sim) {
		if dk.MindFreezeSpell.IsReady(sim) {
			dk.MindFreezeSpell.Cast(sim, target)
		}
	}
	s.Advance()
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionCallback_Pesti_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	// If we have both dots active try to refresh with pesti and move to normal rotation
	if dk.FrostFeverSpell.Dot(dk.CurrentTarget).IsActive() && dk.BloodPlagueSpell.Dot(dk.CurrentTarget).IsActive() {
		dk.Pestilence.Cast(sim, target)
		s.Advance()

		return -1
	} else {
		// If a disease has dropped do normal reapply
		dk.RotationSequence.Clear()

		if dk.sr.ffFirst {
			dk.RotationSequence.
				NewAction(dk.RotationActionUH_FF_ClipCheck).
				NewAction(dk.RotationActionUH_IT_Custom).
				NewAction(dk.RotationActionUH_BP_ClipCheck).
				NewAction(dk.RotationActionUH_PS_Custom)
		} else {
			dk.RotationSequence.
				NewAction(dk.RotationActionUH_BP_ClipCheck).
				NewAction(dk.RotationActionUH_PS_Custom).
				NewAction(dk.RotationActionUH_FF_ClipCheck).
				NewAction(dk.RotationActionUH_IT_Custom)
		}
		return sim.CurrentTime
	}
}

func (dk *DpsDeathknight) RotationActionUH_ResetToSsMain(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	dk.RotationSequence.Clear().
		NewAction(dk.RotationActionCallback_UnholySsRotation)
	return sim.CurrentTime
}

func (dk *DpsDeathknight) RotationActionUH_ResetToDndMain(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	dk.RotationSequence.Clear().
		NewAction(dk.RotationActionCallback_UnholyDndRotation)
	return sim.CurrentTime
}

// Custom PS callback for tracking recasts for pestilence disease sync
func (dk *DpsDeathknight) RotationActionUH_PS_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
		dk.uhAfterGargoyleSequence(sim)
		return sim.CurrentTime
	}
	casted := dk.PlagueStrike.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)

	dk.sr.recastedBP = casted && advance
	s.ConditionalAdvance(casted && advance)
	return -1
}

// Custom IT callback for tracking recasts for pestilence disease sync
func (dk *DpsDeathknight) RotationActionUH_IT_Custom(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()+50*time.Millisecond) {
		dk.uhAfterGargoyleSequence(sim)
		return sim.CurrentTime
	}
	casted := dk.IcyTouch.Cast(sim, target)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)
	if casted && advance {
		dk.sr.recastedFF = true
		dk.ur.syncTimeFF = 0
	}
	s.ConditionalAdvance(casted && advance)
	return -1
}

// Custom IT callback for ghoul frenzy frost rune sync
func (dk *DpsDeathknight) RotationActionUH_IT_SetSync(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	ffRemaining := dk.FrostFeverSpell.Dot(target).RemainingDuration(sim)
	dk.RotationActionCallback_IT(sim, target, s)
	advance := dk.LastOutcome.Matches(core.OutcomeLanded)
	if !dk.GCD.IsReady(sim) && advance {
		dk.ur.syncTimeFF = dk.FrostFeverSpell.Dot(target).Duration - ffRemaining
	}

	return -1
}

func (dk *DpsDeathknight) RotationActionUH_FF_ClipCheck(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	dot := dk.FrostFeverSpell.Dot(target)
	gracePeriod := dk.CurrentFrostRuneGrace(sim)
	return dk.RotationActionUH_DiseaseClipCheck(dot, gracePeriod, sim, target, s)
}

func (dk *DpsDeathknight) RotationActionUH_BP_ClipCheck(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	dot := dk.BloodPlagueSpell.Dot(target)
	gracePeriod := dk.CurrentUnholyRuneGrace(sim)
	return dk.RotationActionUH_DiseaseClipCheck(dot, gracePeriod, sim, target, s)
}

// Check if we have enough rune grace period to delay the disease cast
// so we get more ticks without losing on rune cd
func (dk *DpsDeathknight) RotationActionUH_DiseaseClipCheck(dot *core.Dot, gracePeriod time.Duration, sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	// TODO: Play around with allowing rune cd to be wasted
	// for more disease ticks and see if its a worth option for the ui
	//runeCdWaste := 0 * time.Millisecond
	waitTime := sim.CurrentTime
	if dot.TickCount < dot.NumberOfTicks-1 {
		nextTickAt := dot.ExpiresAt() - dot.TickLength*time.Duration((dot.NumberOfTicks-1)-dot.TickCount)
		if nextTickAt > sim.CurrentTime && (nextTickAt < sim.CurrentTime+gracePeriod || nextTickAt < sim.CurrentTime+400*time.Millisecond) {
			// Delay disease for next tick
			dk.LastOutcome = core.OutcomeMiss

			if dk.uhGargoyleCheck(sim, target, nextTickAt-sim.CurrentTime+50*time.Millisecond) {
				dk.uhAfterGargoyleSequence(sim)
				return waitTime
			}

			waitTime = nextTickAt + 50*time.Millisecond
		} else {
			waitTime = sim.CurrentTime
		}
	} else {
		waitTime = sim.CurrentTime
	}

	s.Advance()
	return waitTime
}
