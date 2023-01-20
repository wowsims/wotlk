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

func (rogue *Rogue) setSubtletyBuilder() {
	mhDagger := rogue.Equip[proto.ItemSlot_ItemSlotMainHand].WeaponType == proto.WeaponType_WeaponTypeDagger
	// Garrote
	if rogue.Garrote.Dot(rogue.CurrentTarget) != nil && rogue.ShadowDanceAura.IsActive() {
		rogue.Builder = rogue.Garrote
		rogue.BuilderPoints = 1
	} else
	// Ambush
	if rogue.ShadowDanceAura.IsActive() {
		rogue.Builder = rogue.Ambush
		rogue.BuilderPoints = 2
	} else
	// Backstab
	if mhDagger {
		rogue.Builder = rogue.Backstab
		rogue.BuilderPoints = 1
	} else
	// Hemorrhage
	{
		rogue.Builder = rogue.Hemorrhage
		rogue.BuilderPoints = 1
	}
	// Ghostly Strike
}

func (rogue *Rogue) setupSubtletyRotation(sim *core.Simulation) {
	rogue.subtletyPrios = make([]subtletyPrio, 0)

	// FIXME: Remove once added to UI
	rogue.Rotation.OpenWithGarrote = true
	rogue.Rotation.OpenWithPremeditation = true
	rogue.Rotation.OpenWithShadowstep = true

	if rogue.Rotation.OpenWithPremeditation {
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

	if rogue.Rotation.OpenWithShadowstep {
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
			if rogue.ComboPoints() > 0 && rogue.CurrentEnergy() > rogue.SliceAndDice[1].DefaultCast.Cost {
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
			return rogue.SliceAndDice[r.ComboPoints()].Cast(s, r.CurrentTarget)
		},
		rogue.SliceAndDice[1].DefaultCast.Cost,
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
				timeLeft := rogue.ExposeArmorAura.RemainingDuration(sim)
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
						if rogue.CurrentEnergy() >= rogue.ExposeArmor[1].DefaultCast.Cost {
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
				casted := r.ExposeArmor[r.ComboPoints()].Cast(sim, r.CurrentTarget)
				if casted {
					hasCastExpose = true
				}
				return casted
			},
			rogue.ExposeArmor[1].DefaultCast.Cost,
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
	if rogue.Rotation.SubtletyFinisherPriority == proto.Rogue_Rotation_Rupture {
		rogue.subtletyPrios = append(rogue.subtletyPrios, subtletyPrio{
			func(s *core.Simulation, r *Rogue) PriorityAction {
				if r.Shadowstep.IsReady(s) {
					// Can we cast Rupture now?
					if r.Rupture[0].Dot(r.CurrentTarget) == nil && rogue.ComboPoints() > 4 && rogue.CurrentEnergy() >= rogue.Rupture[1].DefaultCast.Cost+rogue.Shadowstep.DefaultCast.Cost {
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
			rogue.ShadowDance.DefaultCast.Cost,
		})
	}

	// Rupture
	if rogue.Rotation.SubtletyFinisherPriority == proto.Rogue_Rotation_Rupture {
		rogue.subtletyPrios = append(rogue.subtletyPrios, subtletyPrio{
			func(s *core.Simulation, r *Rogue) PriorityAction {
				if r.Rupture[0].Dot(r.CurrentTarget) != nil || s.GetRemainingDuration() < time.Second*22 {
					return Skip
				}
				if rogue.ComboPoints() > 4 && rogue.CurrentEnergy() >= rogue.Rupture[1].DefaultCast.Cost {
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
				return r.Rupture[r.ComboPoints()].Cast(s, r.CurrentTarget)
			},
			rogue.Rupture[1].DefaultCast.Cost,
		})
	}

	// Eviscerate
	rogue.subtletyPrios = append(rogue.subtletyPrios, subtletyPrio{
		func(s *core.Simulation, r *Rogue) PriorityAction {
			if rogue.Rotation.AllowCpUndercap {
				if r.ComboPoints() == 3 && r.CurrentEnergy() >= r.Eviscerate[1].DefaultCast.Cost {
					return Cast
				}
			}
			energyNeeded := r.Eviscerate[1].DefaultCast.Cost
			minimumCP := int32(4)
			if rogue.Rotation.AllowCpOvercap {
				if r.ComboPoints() == 4 && r.getExpectedComboPointPerSecond() >= 1 {
					return Wait
				}
				if r.ComboPoints() == 4 {
					return Build
				}
			}
			if r.ComboPoints() >= minimumCP && r.CurrentEnergy() >= energyNeeded {
				return Cast
			}
			if r.ComboPoints() < 4 && r.CurrentEnergy() >= r.Builder.DefaultCast.Cost+r.Eviscerate[1].DefaultCast.Cost {
				return Build
			}
			return Wait
		},
		func(s *core.Simulation, r *Rogue) bool {
			return rogue.Eviscerate[r.ComboPoints()].Cast(sim, rogue.CurrentTarget)
		},
		rogue.Eviscerate[1].DefaultCast.Cost,
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
				rogue.setSubtletyBuilder()
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
