package rogue

import (
	"github.com/wowsims/wotlk/sim/core/stats"
	"golang.org/x/exp/slices"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

type assassination_rotation struct {
	prios []prio

	builder *core.Spell
}

func (x *assassination_rotation) setup(sim *core.Simulation, rogue *Rogue) {
	x.builder = rogue.Mutilate

	x.prios = x.prios[:0]
	rogue.bleedCategory = rogue.CurrentTarget.GetExclusiveEffectCategory(core.BleedEffectCategory)

	// Garrote
	if rogue.Rotation.OpenWithGarrote {
		x.prios = append(x.prios, prio{
			func(sim *core.Simulation, rogue *Rogue) PriorityAction {
				if rogue.CurrentEnergy() > rogue.Garrote.DefaultCast.Cost {
					return Once
				}
				return Wait
			},
			func(sim *core.Simulation, rogue *Rogue) bool {
				return rogue.Garrote.Cast(sim, rogue.CurrentTarget)
			},
			rogue.Garrote.DefaultCast.Cost,
		})
	}

	// Slice And Dice
	x.prios = append(x.prios, prio{
		func(sim *core.Simulation, rogue *Rogue) PriorityAction {
			if rogue.SliceAndDiceAura.IsActive() {
				return Skip
			}
			if rogue.ComboPoints() > 0 && rogue.CurrentEnergy() > rogue.SliceAndDice.DefaultCast.Cost {
				return Cast
			}
			if rogue.ComboPoints() < 1 && rogue.CurrentEnergy() > x.builder.DefaultCast.Cost {
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
		x.prios = append(x.prios, prio{
			func(sim *core.Simulation, rogue *Rogue) PriorityAction {

				prioExpose := rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Once ||
					rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Maintain
				if prioExpose && !rogue.ExposeArmorAuras.Get(rogue.CurrentTarget).IsActive() {
					return Skip
				}

				if rogue.HungerForBloodAura.IsActive() {
					return Skip
				}

				if !x.targetHasBleed(sim, rogue) {
					return Skip
				}

				if x.targetHasBleed(sim, rogue) && rogue.CurrentEnergy() > rogue.HungerForBlood.DefaultCast.Cost {
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
		x.prios = append(x.prios, prio{
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
						if rogue.CurrentEnergy() >= x.builder.DefaultCast.Cost {
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
					cpGenerated := energyGained / x.builder.DefaultCast.Cost
					currentCp := float64(rogue.ComboPoints())
					if currentCp+cpGenerated > 5 {
						return Skip
					} else {
						if currentCp < 5 {
							if rogue.CurrentEnergy() >= x.builder.DefaultCast.Cost {
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
		x.prios = append(x.prios, prio{
			func(sim *core.Simulation, rogue *Rogue) PriorityAction {
				if x.targetHasBleed(sim, rogue) {
					return Skip
				}
				if rogue.HungerForBloodAura.IsActive() {
					return Skip
				}
				if rogue.ComboPoints() > 0 && rogue.CurrentEnergy() >= rogue.Rupture.DefaultCast.Cost {
					return Cast
				}
				if rogue.ComboPoints() < 1 && rogue.CurrentEnergy() >= x.builder.DefaultCast.Cost {
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
		x.prios = append(x.prios, prio{
			func(sim *core.Simulation, rogue *Rogue) PriorityAction {
				if rogue.HungerForBloodAura.IsActive() {
					return Skip
				}

				if !x.targetHasBleed(sim, rogue) {
					return Skip
				}

				if x.targetHasBleed(sim, rogue) && rogue.CurrentEnergy() >= rogue.HungerForBlood.DefaultCast.Cost {
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
	x.prios = append(x.prios, prio{
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
	x.prios = append(x.prios, prio{
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
				cost := rogue.Rupture.DefaultCast.Cost + x.builder.DefaultCast.Cost + rogue.Envenom.DefaultCast.Cost
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
				if cp < 4 && e >= x.builder.DefaultCast.Cost {
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
	x.prios = append(x.prios, prio{
		func(sim *core.Simulation, rogue *Rogue) PriorityAction {
			e, cp := rogue.CurrentEnergy(), rogue.ComboPoints()

			costEnv, costMut := rogue.Envenom.DefaultCast.Cost, x.builder.DefaultCast.Cost

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

			if e >= x.builder.DefaultCast.Cost {
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

func (x *assassination_rotation) run(sim *core.Simulation, rogue *Rogue) {
	for i, p := range x.prios {
		switch p.check(sim, rogue) {
		case Skip:
			continue
		case Build:
			if !x.builder.Cast(sim, rogue.CurrentTarget) {
				rogue.WaitForEnergy(sim, x.builder.DefaultCast.Cost)
				return
			}
		case Cast:
			if !p.cast(sim, rogue) {
				rogue.WaitForEnergy(sim, p.cost)
				return
			}
		case Once:
			if !p.cast(sim, rogue) {
				rogue.WaitForEnergy(sim, p.cost)
				return
			}
			x.prios = slices.Delete(x.prios, i, i+1)
		case Wait:
		}
		break
	}
	rogue.DoNothing()
}

func (x *assassination_rotation) targetHasBleed(_ *core.Simulation, rogue *Rogue) bool {
	return rogue.bleedCategory.AnyActive() || rogue.CurrentTarget.HasActiveAuraWithTag(RogueBleedTag)
}
