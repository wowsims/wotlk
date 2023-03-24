package rogue

import (
	"golang.org/x/exp/slices"
	"log"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

type rotation_combat struct {
	prios []prio
}

func (x *rotation_combat) setup(_ *core.Simulation, rogue *Rogue) {
	x.prios = x.prios[:0]

	ssCost := rogue.SinisterStrike.DefaultCast.Cost
	sndCost := rogue.SliceAndDice.DefaultCast.Cost
	rupCost := rogue.Rupture.DefaultCast.Cost
	evisCost := rogue.Eviscerate.DefaultCast.Cost

	baseEps := 10 * rogue.EnergyTickMultiplier
	maxPool := rogue.maxEnergy - 3*float64(rogue.Talents.CombatPotency)

	// estimate of energy per second while nothing is cast
	energyPerSecond := func() float64 {
		if rogue.Talents.CombatPotency == 0 {
			return 10 * rogue.EnergyTickMultiplier
		}

		attackTable := rogue.AttackTables[rogue.CurrentTarget.UnitIndex]
		spell := rogue.AutoAttacks.OHAuto

		landChance := 1.0
		if miss := attackTable.BaseMissChance + 0.19 - spell.PhysicalHitChance(rogue.CurrentTarget); miss > 0 {
			landChance -= miss
		}
		if dodge := attackTable.BaseDodgeChance - spell.ExpertisePercentage() - spell.Unit.PseudoStats.DodgeReduction; dodge > 0 {
			landChance -= dodge
		}
		landsPerSecond := landChance * (1 / rogue.AutoAttacks.OffhandSwingSpeed().Seconds())
		return 10*rogue.EnergyTickMultiplier + landsPerSecond*0.2*float64(rogue.Talents.CombatPotency)*3
	}

	// Garrote
	if rogue.Rotation.OpenWithGarrote && !rogue.PseudoStats.InFrontOfTarget {
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
			cp, e := rogue.ComboPoints(), rogue.CurrentEnergy()

			if rogue.SliceAndDiceAura.IsActive() {
				if cp == 5 { // pool for snd if pooling for rupture fails
					rupDur := rogue.Rupture.CurDot().RemainingDuration(sim)
					if e+rupDur.Seconds()*energyPerSecond() > maxPool {
						sndDur := rogue.SliceAndDiceAura.RemainingDuration(sim)
						if e+sndDur.Seconds()*energyPerSecond() <= maxPool {
							return Wait
						}
					}
					return Skip
				}

				if cp >= 1 { // don't build if it reduces uptime
					sndDur := rogue.SliceAndDiceAura.RemainingDuration(sim)
					if e+sndDur.Seconds()*energyPerSecond() < sndCost+ssCost || sndDur < time.Second {
						return Wait
					}
				}
				return Skip
			}

			// end of fight - heuristically, 2s of snd beat a 3 CP eviscerate for DPE, and 3s are close to a 5 CP one.
			if cp >= 3 && sim.GetRemainingDuration() < time.Duration(2000+600*cp)*time.Millisecond {
				return Skip
			}

			if cp >= 1 && e >= sndCost {
				return Cast
			}
			if cp < 1 && e >= ssCost {
				return Build
			}
			return Wait
		},
		func(sim *core.Simulation, rogue *Rogue) bool {
			return rogue.SliceAndDice.Cast(sim, rogue.CurrentTarget)
		},
		sndCost,
	})

	// Expose armor - update this as well
	if rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Once || rogue.Rotation.ExposeArmorFrequency == proto.Rogue_Rotation_Maintain {
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
						if rogue.CurrentEnergy() >= ssCost {
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
					energyGained := energyPerSecond() * timeLeft.Seconds()
					cpGenerated := energyGained / ssCost
					currentCp := float64(rogue.ComboPoints())
					if currentCp+cpGenerated > 5 {
						return Skip
					} else {
						if currentCp < 5 {
							if rogue.CurrentEnergy() >= ssCost {
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

	// Enable CDs
	x.prios = append(x.prios, prio{
		func(sim *core.Simulation, rogue *Rogue) PriorityAction {
			for _, mcd := range rogue.GetMajorCooldowns() {
				mcd.Enable()
			}
			return Once
		},
		func(s *core.Simulation, r *Rogue) bool {
			return true
		},
		0,
	})

	const ruptureMinDuration = time.Second * 10 // heuristically, 4-5 rupture ticks are better DPE than eviscerate

	// seconds a 5 cp rupture can be delayed to match a 4 cp rupture's dps. for rup4to5 and rup3to4, this delay is < 2s,
	// which also means that clipping 3 or 4 cp ruptures is usually a dps loss
	rup4to5 := rogue.RuptureDuration(4).Seconds() * (1 - rogue.RuptureDamage(4)/rogue.RuptureDamage(5))
	rup3to4 := rogue.RuptureDuration(3).Seconds() * (1 - rogue.RuptureDamage(3)/rogue.RuptureDamage(4))

	// Rupture
	x.prios = append(x.prios, prio{
		func(sim *core.Simulation, rogue *Rogue) PriorityAction {
			cp, e := rogue.ComboPoints(), rogue.CurrentEnergy()

			if sim.GetRemainingDuration() < ruptureMinDuration {
				return Skip
			}

			rupDot := rogue.Rupture.CurDot()
			if !rupDot.IsActive() {
				if cp == 5 && e >= rupCost {
					return Cast
				}
				if cp == 4 && e+rup4to5*energyPerSecond() < ssCost+rupCost {
					return Cast
				}
				if cp == 3 && e+rup3to4*energyPerSecond() < ssCost+rupCost {
					return Cast
				}
				if e >= ssCost {
					return Build
				}
				return Wait
			}

			// there's ample time to rebuild, simply skip
			dur := rupDot.RemainingDuration(sim).Seconds()
			if e+dur*baseEps > maxPool {
				return Skip
			}

			if cp == 5 {
				if e+dur*energyPerSecond() > maxPool {
					return Skip // can't pool any longer, maybe we can fit in Eviscerate
				}
				return Wait
			}
			if cp == 4 && e+(dur+rup4to5)*energyPerSecond() < ssCost+rupCost {
				return Wait
			}
			if cp == 3 && e+(dur+rup3to4)*energyPerSecond() < ssCost+rupCost {
				return Wait
			}
			if e >= ssCost {
				return Build
			}
			return Wait
		},
		func(sim *core.Simulation, rogue *Rogue) bool {
			return rogue.Rupture.Cast(sim, rogue.CurrentTarget)
		},
		rupCost,
	})

	ssPerCp := 1.0
	if rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfSinisterStrike) {
		attackTable := rogue.AttackTables[rogue.CurrentTarget.UnitIndex]
		crit := rogue.SinisterStrike.PhysicalCritChance(rogue.CurrentTarget, attackTable)
		ssPerCp = 1 / (1 + crit*0.5)
	}

	// Eviscerate
	x.prios = append(x.prios, prio{
		func(sim *core.Simulation, rogue *Rogue) PriorityAction {
			e, cp := rogue.CurrentEnergy(), rogue.ComboPoints()

			if dur := sim.GetRemainingDuration(); dur <= ruptureMinDuration {
				// end of fight handling - build towards a 3+ cp eviscerate, or just sinister strike
				switch cp {
				case 5:
					if e >= evisCost {
						return Cast
					}
					return Wait
				default:
					if e+dur.Seconds()*energyPerSecond() >= ssCost+evisCost {
						return Build
					}
					if cp >= 3 && e >= evisCost {
						return Cast
					}
					if cp < 3 && e >= ssCost {
						return Build
					}
				}
				return Wait
			}

			// we only get here if there's ample time left on rupture, or rupture pooling failed: in these cases, we
			// can try to fill in a 5 cp eviscerate, if it's not too disruptive. lower cp eviscerates aren't worth it,
			// since sinister spam isn't all that much worse
			if cp <= 4 {
				return Build
			}

			rupDot := rogue.Rupture.CurDot()

			ruthCP := 0.2 * float64(rogue.Talents.Ruthlessness)
			cost := evisCost + (4-ruthCP)*ssCost*ssPerCp + rupCost

			rupDur := rupDot.RemainingDuration(sim)
			sndDur := rogue.SliceAndDiceAura.RemainingDuration(sim)
			if sndDur < rupDur {
				cost += sndCost + (1-ruthCP)*ssCost*ssPerCp
			}

			if avail := e + rupDur.Seconds()*energyPerSecond(); avail >= cost {
				return Cast
			}
			return Build
		},
		func(sim *core.Simulation, rogue *Rogue) bool {
			return rogue.Eviscerate.Cast(sim, rogue.CurrentTarget)
		},
		evisCost,
	})
}

func (x *rotation_combat) run(sim *core.Simulation, rogue *Rogue) {
	if rogue.KillingSpreeAura.IsActive() {
		rogue.DoNothing()
		return
	}

	for i := 0; i < len(x.prios); i++ {
		switch p := x.prios[i]; p.check(sim, rogue) {
		case Skip:
			continue
		case Build:
			if !rogue.SinisterStrike.Cast(sim, rogue.CurrentTarget) {
				rogue.WaitForEnergy(sim, rogue.SinisterStrike.DefaultCast.Cost)
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
			i--
		case Wait:
			rogue.DoNothing()
			return
		}

		if !rogue.GCD.IsReady(sim) {
			return
		}
	}
	log.Panic("skipped all prios")
}
