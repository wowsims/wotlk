package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

type subtletyPrio struct {
	check GetAction
	cast  DoAction
	cost  float64
}

func (rogue *Rogue) setSubtletyBuilder(sim *core.Simulation) {
	mhDagger := rogue.Equip[proto.ItemSlot_ItemSlotMainHand].WeaponType == proto.WeaponType_WeaponTypeDagger
	// Garrote
	if !rogue.Garrote.Dot(rogue.CurrentTarget).Aura.IsActive() && rogue.ShadowDanceAura.IsActive() && !rogue.PseudoStats.InFrontOfTarget {
		rogue.Builder = rogue.Garrote
		rogue.BuilderPoints = 1
		return
	}
	// Ambush
	if rogue.ShadowDanceAura.IsActive() && mhDagger && !rogue.PseudoStats.InFrontOfTarget {
		rogue.Builder = rogue.Ambush
		rogue.BuilderPoints = 2
		return
	}
	// Backstab
	if mhDagger && !rogue.Rotation.HemoWithDagger && !rogue.PseudoStats.InFrontOfTarget {
		rogue.Builder = rogue.Backstab
		rogue.BuilderPoints = 1
		return
	}
	// Ghostly Strike
	if rogue.Talents.GhostlyStrike && rogue.Rotation.UseGhostlyStrike && rogue.GhostlyStrike.IsReady(sim) {
		rogue.Builder = rogue.GhostlyStrike
		rogue.BuilderPoints = 1
		return
	}
	// Hemorrhage
	if rogue.Talents.Hemorrhage {
		rogue.Builder = rogue.Hemorrhage
		rogue.BuilderPoints = 1
	} else
	// Sinister Strike
	{
		rogue.Builder = rogue.SinisterStrike
		rogue.BuilderPoints = 1
	}
}

func (rogue *Rogue) setupSubtletyRotation(sim *core.Simulation) {
	rogue.subtletyPrios = make([]subtletyPrio, 0)

	if rogue.Rotation.OpenWithPremeditation && rogue.Talents.Premeditation {
		hasCastPremeditation := false
		rogue.subtletyPrios = append(rogue.subtletyPrios, subtletyPrio{
			func(s *core.Simulation, r *Rogue) PriorityAction {
				if hasCastPremeditation {
					return Skip
				}
				if rogue.Preparation.IsReady(s) {
					return Cast
				}
				return Wait
			},
			func(s *core.Simulation, r *Rogue) bool {
				casted := r.Premeditation.Cast(s, r.CurrentTarget)
				if casted {
					hasCastPremeditation = true
				}
				return casted
			},
			rogue.Preparation.DefaultCast.Cost,
		})
	}

	if rogue.Rotation.OpenWithShadowstep && rogue.Talents.Shadowstep {
		hasCastShadowstep := false
		rogue.subtletyPrios = append(rogue.subtletyPrios, subtletyPrio{
			func(s *core.Simulation, r *Rogue) PriorityAction {
				if hasCastShadowstep {
					return Skip
				}
				if rogue.CurrentEnergy() > rogue.Shadowstep.DefaultCast.Cost {
					return Cast
				}
				return Wait
			},
			func(s *core.Simulation, r *Rogue) bool {
				casted := rogue.Shadowstep.Cast(sim, rogue.CurrentTarget)
				if casted {
					hasCastShadowstep = true
				}
				return casted
			},
			rogue.Shadowstep.DefaultCast.Cost,
		})
	}

	// Garrote
	if rogue.Rotation.OpenWithGarrote {
		hasCastGarrote := false
		rogue.subtletyPrios = append(rogue.subtletyPrios, subtletyPrio{
			func(s *core.Simulation, r *Rogue) PriorityAction {
				if hasCastGarrote {
					return Skip
				}
				if rogue.CurrentEnergy() > rogue.Garrote.DefaultCast.Cost {
					return Cast
				}
				return Wait
			},
			func(s *core.Simulation, r *Rogue) bool {
				casted := rogue.Garrote.Cast(sim, rogue.CurrentTarget)
				if casted {
					hasCastGarrote = true
				}
				return casted
			},
			rogue.Garrote.DefaultCast.Cost,
		})
	}

	// Slice and Dice
	rogue.subtletyPrios = append(rogue.subtletyPrios, subtletyPrio{
		func(s *core.Simulation, r *Rogue) PriorityAction {
			if rogue.SliceAndDiceAura.IsActive() {
				return Skip
			}
			if rogue.ComboPoints() > 0 && rogue.CurrentEnergy() > rogue.SliceAndDice.DefaultCast.Cost {
				return Cast
			}
			if rogue.ComboPoints() < 1 && rogue.CurrentEnergy() > rogue.Builder.DefaultCast.Cost && rogue.getExpectedComboPointPerSecond() >= 0.7 {
				return Wait
			}
			if rogue.ComboPoints() < 1 && rogue.CurrentEnergy() > rogue.Builder.DefaultCast.Cost {
				return Build
			}
			return Wait
		},
		func(s *core.Simulation, r *Rogue) bool {
			return rogue.SliceAndDice.Cast(s, r.CurrentTarget)
		},
		rogue.SliceAndDice.DefaultCast.Cost,
	})

	// Expose armor
	if rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Once ||
		rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Maintain {
		hasCastExpose := false
		rogue.subtletyPrios = append(rogue.subtletyPrios, subtletyPrio{
			func(s *core.Simulation, r *Rogue) PriorityAction {
				if hasCastExpose && rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Once {
					return Skip
				}
				timeLeft := rogue.ExposeArmorAuras.Get(rogue.CurrentTarget).RemainingDuration(sim)
				minPoints := core.MaxInt32(1, core.MinInt32(rogue.Rotation.MinimumComboPointsExposeArmor, 5))
				if rogue.Rotation.ExposeArmorFrequency != proto.Rogue_Rotation_Once {
					minPoints = 1
				}
				if timeLeft <= 0 {
					if rogue.ComboPoints() < minPoints {
						if rogue.CurrentEnergy() >= rogue.Builder.DefaultCast.Cost && rogue.getExpectedComboPointPerSecond() < 1 {
							return Build
						} else {
							return Wait
						}
					} else {
						if rogue.CurrentEnergy() >= rogue.ExposeArmor.DefaultCast.Cost {
							return Cast
						} else {
							return Wait
						}
					}
				} else {
					energyGained := rogue.getExpectedEnergyPerSecond() * timeLeft.Seconds()
					comboGained := rogue.getExpectedComboPointPerSecond() * timeLeft.Seconds()
					cpGenerated := energyGained/rogue.Builder.DefaultCast.Cost + comboGained
					currentCP := float64(rogue.ComboPoints())
					if currentCP+cpGenerated > 5 {
						return Skip
					} else {
						if currentCP < 5 {
							if rogue.CurrentEnergy() >= rogue.Builder.DefaultCast.Cost {
								return Build
							}
						}
						return Wait
					}
				}
			},
			func(s *core.Simulation, r *Rogue) bool {
				casted := r.ExposeArmor.Cast(sim, r.CurrentTarget)
				if casted {
					hasCastExpose = true
				}
				return casted
			},
			rogue.ExposeArmor.DefaultCast.Cost,
		})
	}

	// Enable CDS
	rogue.subtletyPrios = append(rogue.subtletyPrios, subtletyPrio{
		func(s *core.Simulation, r *Rogue) PriorityAction {
			if r.allMCDsDisabled {
				for _, mcd := range r.GetMajorCooldowns() {
					mcd.Enable()
				}
				r.allMCDsDisabled = false
			}
			return Skip
		},
		func(s *core.Simulation, r *Rogue) bool {
			return false
		},
		0,
	})

	//Shadowstep
	if rogue.Talents.Shadowstep {
		rogue.subtletyPrios = append(rogue.subtletyPrios, subtletyPrio{
			func(s *core.Simulation, r *Rogue) PriorityAction {
				if r.Shadowstep.IsReady(s) {
					// Can we cast Rupture now?
					if !r.Rupture.CurDot().IsActive() && rogue.ComboPoints() > 4 && rogue.CurrentEnergy() >= rogue.Rupture.DefaultCast.Cost+rogue.Shadowstep.DefaultCast.Cost {
						return Cast
					} else {
						return Skip
					}
				}
				return Skip
			},
			func(s *core.Simulation, r *Rogue) bool {
				return r.Shadowstep.Cast(s, r.CurrentTarget)
			},
			rogue.Shadowstep.DefaultCast.Cost,
		})
	}

	// Rupture
	rogue.subtletyPrios = append(rogue.subtletyPrios, subtletyPrio{
		func(s *core.Simulation, r *Rogue) PriorityAction {
			if r.Rupture.CurDot().IsActive() || s.GetRemainingDuration() < time.Second*22 {
				return Skip
			}
			if rogue.ComboPoints() > 4 && rogue.CurrentEnergy() >= rogue.Rupture.DefaultCast.Cost {
				return Cast
			}
			if rogue.ComboPoints() < 5 && rogue.CurrentEnergy()+rogue.getExpectedEnergyPerSecond() >= rogue.maxEnergy {
				return Build
			}
			if rogue.ComboPoints() < 5 && rogue.getExpectedComboPointPerSecond() >= 0.7 {
				return Wait
			}
			if rogue.ComboPoints() < 5 && rogue.CurrentEnergy() >= rogue.Builder.DefaultCast.Cost {
				return Build
			}
			return Wait
		},
		func(s *core.Simulation, r *Rogue) bool {
			return r.Rupture.Cast(s, r.CurrentTarget)
		},
		rogue.Rupture.DefaultCast.Cost,
	})

	//Envenom
	if rogue.Rotation.SubtletyFinisherPriority == proto.Rogue_Rotation_SubtletyEnvenom {
		rogue.subtletyPrios = append(rogue.subtletyPrios, subtletyPrio{
			func(s *core.Simulation, r *Rogue) PriorityAction {
				minimumCP := int32(5)
				if !r.DeadlyPoison.Dot(r.CurrentTarget).Aura.IsActive() {
					return Skip
				}
				if r.EnvenomAura.IsActive() {
					return Skip
				}
				if rogue.ComboPoints() >= minimumCP && rogue.CurrentEnergy() >= rogue.Envenom.DefaultCast.Cost {
					return Cast
				}
				if rogue.ComboPoints() < minimumCP && rogue.CurrentEnergy()+rogue.getExpectedEnergyPerSecond() >= rogue.maxEnergy {
					return Build
				}
				if rogue.ComboPoints() < minimumCP && rogue.getExpectedComboPointPerSecond() >= 0.7 {
					return Wait
				}
				if rogue.ComboPoints() < minimumCP && rogue.CurrentEnergy() >= rogue.Builder.DefaultCast.Cost {
					return Build
				}
				return Wait
			},
			func(s *core.Simulation, r *Rogue) bool {
				return rogue.Envenom.Cast(sim, rogue.CurrentTarget)
			},
			rogue.Envenom.DefaultCast.Cost,
		})
	}

	// Eviscerate
	rogue.subtletyPrios = append(rogue.subtletyPrios, subtletyPrio{
		func(s *core.Simulation, r *Rogue) PriorityAction {
			energyNeeded := r.Eviscerate.DefaultCast.Cost
			minimumCP := int32(5)
			if r.ComboPoints() >= minimumCP && r.CurrentEnergy() >= energyNeeded {
				return Cast
			}
			if r.ComboPoints() < minimumCP && r.CurrentEnergy() >= r.Builder.DefaultCast.Cost+r.Eviscerate.DefaultCast.Cost {
				return Build
			}
			return Wait
		},
		func(s *core.Simulation, r *Rogue) bool {
			return rogue.Eviscerate.Cast(sim, rogue.CurrentTarget)
		},
		rogue.Eviscerate.DefaultCast.Cost,
	})
}

func (rogue *Rogue) doSubtletyRotation(sim *core.Simulation) {
	prioIndex := 0
	for prioIndex < len(rogue.subtletyPrios) {
		prio := rogue.subtletyPrios[prioIndex]
		switch prio.check(sim, rogue) {
		case Skip:
			prioIndex += 1
		case Build:
			if rogue.GCD.IsReady(sim) {
				rogue.setSubtletyBuilder(sim)
				if !rogue.Builder.Cast(sim, rogue.CurrentTarget) {
					rogue.WaitForEnergy(sim, rogue.Builder.DefaultCast.Cost)
					return
				}
			}
			rogue.DoNothing()
			return
		case Cast:
			if rogue.GCD.IsReady(sim) {
				if !prio.cast(sim, rogue) {
					rogue.WaitForEnergy(sim, prio.cost)
					return
				}
			}
			rogue.DoNothing()
			return
		case Wait:
			rogue.DoNothing()
			return
		}
	}
	rogue.DoNothing()
}

func (rogue *Rogue) OnCanActSubtlety(sim *core.Simulation) {
	if rogue.GCD.IsReady(sim) {
		rogue.doSubtletyRotation(sim)
	}
}
