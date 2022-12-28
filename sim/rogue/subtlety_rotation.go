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

func (rogue *Rogue) setupSubtletyRotation(sim *core.Simulation) {
	rogue.subtletyPrios = make([]subtletyPrio, 0)

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
						if rogue.CurrentEnergy() >= rogue.Builder.DefaultCast.Cost {
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
					cpGenerated := energyGained / rogue.Builder.DefaultCast.Cost
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

	// Rupture
	if rogue.Rotation.SubtletyFinisherPriority == proto.Rogue_Rotation_Rupture {
		rogue.subtletyPrios = append(rogue.subtletyPrios, subtletyPrio{
			func(s *core.Simulation, r *Rogue) PriorityAction {
				if r.ruptureDot.IsActive() || s.GetRemainingDuration() < time.Second*18 {
					return Skip
				}
				if rogue.ComboPoints() > 3 && rogue.CurrentEnergy() >= rogue.Rupture[1].DefaultCast.Cost {
					return Cast
				}
				if rogue.ComboPoints() < 4 && rogue.CurrentEnergy() >= rogue.Builder.DefaultCast.Cost {
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
			energyNeeded := core.MinFloat(r.maxEnergy, float64(rogue.Rotation.EnvenomEnergyThreshold))
			// Don't pool when fight is about to end
			if s.GetRemainingDuration() <= time.Second*4 {
				energyNeeded = r.Eviscerate[1].DefaultCast.Cost
			}
			energyNeeded = core.MaxFloat(r.Eviscerate[1].DefaultCast.Cost, energyNeeded)
			minimumCP := int32(4)
			if rogue.Rotation.AllowCpOvercap {
				if r.ComboPoints() == 4 {
					return Build
				}
			}
			if r.ComboPoints() >= minimumCP && r.CurrentEnergy() >= energyNeeded {
				return Cast
			}
			if r.ComboPoints() < 4 && r.CurrentEnergy() >= r.Builder.DefaultCast.Cost {
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
	if !rogue.HonorAmongThievesDot.IsActive() {
		rogue.HonorAmongThieves.Cast(sim, rogue.CurrentTarget)
	}
	if rogue.KillingSpreeAura.IsActive() {
		rogue.DoNothing()
		return
	}
	rogue.TryUseCooldowns(sim)
	if rogue.GCD.IsReady(sim) {
		rogue.doSubtletyRotation(sim)
	}
}
