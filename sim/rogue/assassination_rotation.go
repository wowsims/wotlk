package rogue

import (
	"github.com/wowsims/wotlk/sim/core/stats"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

type PriorityAction int32

const (
	Skip PriorityAction = iota
	Build
	Cast
	Wait
)

type GetAction func(*core.Simulation, *Rogue) PriorityAction
type DoAction func(*core.Simulation, *Rogue) bool

type assassinationPrio struct {
	check GetAction
	cast  DoAction
	cost  float64
}

func (rogue *Rogue) targetHasBleed(_ *core.Simulation) bool {
	return rogue.bleedCategory.AnyActive() || rogue.CurrentTarget.HasActiveAuraWithTag(RogueBleedTag)
}

func (rogue *Rogue) setupAssassinationRotation(sim *core.Simulation) {
	rogue.assassinationPrios = rogue.assassinationPrios[:0]
	rogue.bleedCategory = rogue.CurrentTarget.GetExclusiveEffectCategory(core.BleedEffectCategory)

	// Garrote
	if rogue.Rotation.OpenWithGarrote {
		hasCastGarrote := false
		rogue.assassinationPrios = append(rogue.assassinationPrios, assassinationPrio{
			func(sim *core.Simulation, rogue *Rogue) PriorityAction {
				if hasCastGarrote {
					return Skip
				}
				if rogue.CurrentEnergy() > rogue.Garrote.DefaultCast.Cost {
					return Cast
				}
				return Wait
			},
			func(sim *core.Simulation, rogue *Rogue) bool {
				casted := rogue.Garrote.Cast(sim, rogue.CurrentTarget)
				if casted {
					hasCastGarrote = true
				}
				return casted
			},
			rogue.Garrote.DefaultCast.Cost,
		})
	}

	// Slice And Dice
	rogue.assassinationPrios = append(rogue.assassinationPrios, assassinationPrio{
		func(sim *core.Simulation, rogue *Rogue) PriorityAction {
			if rogue.SliceAndDiceAura.IsActive() {
				return Skip
			}
			if rogue.ComboPoints() > 0 && rogue.CurrentEnergy() > rogue.SliceAndDice.DefaultCast.Cost {
				return Cast
			}
			if rogue.ComboPoints() < 1 && rogue.CurrentEnergy() > rogue.Builder.DefaultCast.Cost {
				return Build
			}
			return Wait
		},
		func(sim *core.Simulation, rogue *Rogue) bool {
			return rogue.SliceAndDice.Cast(sim, rogue.CurrentTarget)
		},
		rogue.SliceAndDice.DefaultCast.Cost,
	})

	// Hunger while planning
	if rogue.Talents.HungerForBlood {
		rogue.assassinationPrios = append(rogue.assassinationPrios, assassinationPrio{
			func(sim *core.Simulation, rogue *Rogue) PriorityAction {

				prioExpose := rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Once ||
					rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Maintain
				if prioExpose && !rogue.ExposeArmorAuras.Get(rogue.CurrentTarget).IsActive() {
					return Skip
				}

				if rogue.HungerForBloodAura.IsActive() {
					return Skip
				}

				if !rogue.targetHasBleed(sim) {
					return Skip
				}

				if rogue.targetHasBleed(sim) && rogue.CurrentEnergy() > rogue.HungerForBlood.DefaultCast.Cost {
					return Cast
				}
				return Wait
			},
			func(sim *core.Simulation, rogue *Rogue) bool {
				return rogue.HungerForBlood.Cast(sim, rogue.CurrentTarget)
			},
			rogue.HungerForBlood.DefaultCast.Cost,
		})
	}

	// Expose armor
	if rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Once ||
		rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Maintain {
		hasCastExpose := false
		rogue.assassinationPrios = append(rogue.assassinationPrios, assassinationPrio{
			func(sim *core.Simulation, rogue *Rogue) PriorityAction {
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
						if rogue.CurrentEnergy() >= rogue.Builder.DefaultCast.Cost {
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
					cpGenerated := energyGained / rogue.Builder.DefaultCast.Cost
					currentCp := float64(rogue.ComboPoints())
					if currentCp+cpGenerated > 5 {
						return Skip
					} else {
						if currentCp < 5 {
							if rogue.CurrentEnergy() >= rogue.Builder.DefaultCast.Cost {
								return Build
							}
						}
						return Wait
					}
				}
			},
			func(sim *core.Simulation, rogue *Rogue) bool {
				casted := rogue.ExposeArmor.Cast(sim, rogue.CurrentTarget)
				if casted {
					hasCastExpose = true
				}
				return casted
			},
			rogue.ExposeArmor.DefaultCast.Cost,
		})
	}

	// Rupture for Bleed
	if rogue.Rotation.RuptureForBleed {
		rogue.assassinationPrios = append(rogue.assassinationPrios, assassinationPrio{
			func(sim *core.Simulation, rogue *Rogue) PriorityAction {
				if rogue.targetHasBleed(sim) {
					return Skip
				}
				if rogue.HungerForBloodAura.IsActive() {
					return Skip
				}
				if rogue.ComboPoints() > 0 && rogue.CurrentEnergy() >= rogue.Rupture.DefaultCast.Cost {
					return Cast
				}
				if rogue.ComboPoints() < 1 && rogue.CurrentEnergy() >= rogue.Builder.DefaultCast.Cost {
					return Build
				}
				return Wait
			},
			func(sim *core.Simulation, rogue *Rogue) bool {
				return rogue.Rupture.Cast(sim, rogue.CurrentTarget)
			},
			rogue.Rupture.DefaultCast.Cost,
		})
	}

	// Hunger for Blood
	if rogue.Talents.HungerForBlood {
		rogue.assassinationPrios = append(rogue.assassinationPrios, assassinationPrio{
			func(sim *core.Simulation, rogue *Rogue) PriorityAction {
				if rogue.HungerForBloodAura.IsActive() {
					return Skip
				}

				if !rogue.targetHasBleed(sim) {
					return Skip
				}

				if rogue.targetHasBleed(sim) && rogue.CurrentEnergy() >= rogue.HungerForBlood.DefaultCast.Cost {
					return Cast
				}
				return Wait
			},
			func(s *core.Simulation, r *Rogue) bool {
				return rogue.HungerForBlood.Cast(sim, rogue.CurrentTarget)
			},
			rogue.HungerForBlood.DefaultCast.Cost,
		})
	}

	// Enable CDs
	rogue.assassinationPrios = append(rogue.assassinationPrios, assassinationPrio{
		func(sim *core.Simulation, rogue *Rogue) PriorityAction {
			if rogue.allMCDsDisabled {
				for _, mcd := range rogue.GetMajorCooldowns() {
					if mcd.Spell != rogue.ColdBlood {
						mcd.Enable()
					}
				}
				rogue.allMCDsDisabled = false
			}
			return Skip
		},
		func(s *core.Simulation, r *Rogue) bool {
			return false
		},
		0,
	})

	// estimate of energy per second while nothing is cast
	energyPerSecond := func() float64 {
		if rogue.Talents.FocusedAttacks == 0 {
			return 10 * rogue.EnergyTickMultiplier
		}

		procChance := []float64{0, 0.33, 0.66, 1}[rogue.Talents.FocusedAttacks]
		critSuppression := rogue.AttackTables[rogue.CurrentTarget.UnitIndex].CritSuppression
		effectiveCrit := rogue.GetStat(stats.MeleeCrit)/(core.CritRatingPerCritChance*100) - critSuppression
		critsPerSecond := effectiveCrit * procChance * (1/rogue.AutoAttacks.MainhandSwingSpeed().Seconds() + 1/rogue.AutoAttacks.OffhandSwingSpeed().Seconds())
		return 10*rogue.EnergyTickMultiplier + critsPerSecond*2
	}

	// Rupture
	rogue.assassinationPrios = append(rogue.assassinationPrios, assassinationPrio{
		func(sim *core.Simulation, rogue *Rogue) PriorityAction {
			cp, e := rogue.ComboPoints(), rogue.CurrentEnergy()

			if rogue.Rotation.AssassinationFinisherPriority == proto.Rogue_Rotation_EnvenomRupture {
				if rogue.Rupture.CurDot().IsActive() || sim.GetRemainingDuration() < rogue.RuptureDuration(4) {
					return Skip
				}
				if !rogue.EnvenomAura.IsActive() || cp < 4 || rogue.Talents.Ruthlessness < 3 {
					return Skip
				}

				// use Rupture if you can re-cast Envenom with minimal delay, hoping for a Ruthlessness proc ;)
				avail := e + rogue.EnvenomAura.RemainingDuration(sim).Seconds()*energyPerSecond()
				cost := rogue.Rupture.DefaultCast.Cost + rogue.Builder.DefaultCast.Cost + rogue.Envenom.DefaultCast.Cost
				if avail >= cost {
					return Cast
				}
				return Skip

			} else {
				if rogue.Rupture.CurDot().IsActive() || sim.GetRemainingDuration() < time.Second*18 {
					return Skip
				}
				if cp >= 4 && e >= rogue.Rupture.DefaultCast.Cost {
					return Cast
				}
				if cp < 4 && e >= rogue.Builder.DefaultCast.Cost {
					return Build
				}
				return Wait
			}
		},
		func(sim *core.Simulation, rogue *Rogue) bool {
			return rogue.Rupture.Cast(sim, rogue.CurrentTarget)
		},
		rogue.Rupture.DefaultCast.Cost,
	})

	// Envenom
	rogue.assassinationPrios = append(rogue.assassinationPrios, assassinationPrio{
		func(sim *core.Simulation, rogue *Rogue) PriorityAction {
			e, cp := rogue.CurrentEnergy(), rogue.ComboPoints()

			costEnv, costMut := rogue.Envenom.DefaultCast.Cost, rogue.Builder.DefaultCast.Cost

			// end of combat handling - possibly use low CP Envenoms instead of doing nothing
			if dur := sim.GetRemainingDuration(); dur <= 10*time.Second {
				avail := e + dur.Seconds()*energyPerSecond()

				if cp == 3 && avail < costMut+costEnv && e >= costEnv {
					return Cast
				}

				if cp >= 1 && avail < costMut && e >= costEnv {
					return Cast
				}
			}

			if cp >= 4 {
				eps := energyPerSecond()

				if rogue.EnvenomAura.IsActive() {
					// don't clip Envenom, unless you'd energy cap
					if e < rogue.maxEnergy-eps && sim.GetRemainingDuration() >= rogue.EnvenomDuration(5) {
						return Wait
					}
					return Cast
				}

				// pool, so two Mutilate casts fit into the next uptime; this is a very minor DPS gain, and primarily for lower gear levels
				cost := costEnv + costMut + costMut
				if cp == 5 && rogue.Talents.RelentlessStrikes == 5 {
					cost -= 25
				}
				avail := e + rogue.EnvenomDuration(cp).Seconds()*eps
				if avail < cost {
					return Wait
				}
				return Cast
			}

			if e >= rogue.Builder.DefaultCast.Cost {
				return Build
			}
			return Wait
		},
		func(sim *core.Simulation, rogue *Rogue) bool {
			if rogue.ColdBlood.IsReady(sim) && rogue.ComboPoints() == 5 {
				rogue.ColdBlood.Cast(sim, rogue.CurrentTarget)
			}
			return rogue.Envenom.Cast(sim, rogue.CurrentTarget)
		},
		rogue.Envenom.DefaultCast.Cost,
	})
}

func (rogue *Rogue) doAssassinationRotation(sim *core.Simulation) {
	prioIndex := 0
	for prioIndex < len(rogue.assassinationPrios) {
		prio := rogue.assassinationPrios[prioIndex]
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

func (rogue *Rogue) OnCanAct(sim *core.Simulation) {
	if rogue.KillingSpreeAura.IsActive() {
		rogue.DoNothing()
		return
	}
	rogue.TryUseCooldowns(sim)
	if rogue.GCD.IsReady(sim) {
		rogue.doAssassinationRotation(sim)
	}
}
