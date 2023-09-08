package rogue

import (
	"golang.org/x/exp/slices"
	"log"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

type rotation_generic struct {
	prios []prio

	builder *core.Spell
}

func (x *rotation_generic) setup(_ *core.Simulation, rogue *Rogue) {
	x.prios = x.prios[:0]

	x.builder = rogue.SinisterStrike
	if rogue.HasDagger(core.MainHand) && !rogue.PseudoStats.InFrontOfTarget {
		x.builder = rogue.Backstab
	}
	if rogue.CanMutilate() {
		x.builder = rogue.Mutilate
	}
	if rogue.Talents.Hemorrhage {
		x.builder = rogue.Hemorrhage
	}

	bldCost := x.builder.DefaultCast.Cost
	sndCost := rogue.SliceAndDice.DefaultCast.Cost
	rupCost := rogue.Rupture.DefaultCast.Cost

	baseEps := 10 * rogue.EnergyTickMultiplier
	maxPool := rogue.maxEnergy - 3*float64(rogue.Talents.CombatPotency) - 2*float64(rogue.Talents.FocusedAttacks)/3.0

	ruthCp := 0.2 * float64(rogue.Talents.Ruthlessness)
	rsPerCp := float64(rogue.Talents.RelentlessStrikes)

	// estimate of energy per second while nothing is cast
	energyPerSecond := func() float64 {
		var eps float64
		if rogue.Talents.CombatPotency > 0 {
			spell := rogue.AutoAttacks.OHAuto
			at := rogue.AttackTables[rogue.CurrentTarget.UnitIndex]

			landChance := 1.0
			if miss := at.BaseMissChance + 0.19 - spell.PhysicalHitChance(at); miss > 0 {
				landChance -= miss
			}
			if dodge := at.BaseDodgeChance - spell.ExpertisePercentage() - spell.Unit.PseudoStats.DodgeReduction; dodge > 0 {
				landChance -= dodge
			}

			landsPerSecond := landChance / rogue.AutoAttacks.OffhandSwingSpeed().Seconds()

			eps += landsPerSecond * 0.2 * float64(rogue.Talents.CombatPotency) * 3
		}
		if rogue.Talents.FocusedAttacks > 0 {
			getCritChance := func(spell *core.Spell) float64 {
				at := rogue.AttackTables[rogue.CurrentTarget.UnitIndex]

				critCap := 1.0 - at.BaseGlanceChance
				if miss := at.BaseMissChance + 0.19 - spell.PhysicalHitChance(at); miss > 0 {
					critCap -= miss
				}
				if dodge := at.BaseDodgeChance - spell.ExpertisePercentage() - rogue.PseudoStats.DodgeReduction; dodge > 0 {
					critCap -= dodge
				}

				critChance := spell.PhysicalCritChance(at)
				if critChance > critCap {
					critChance = critCap
				}
				return critChance
			}

			critsPerSecond := getCritChance(rogue.AutoAttacks.MHAuto)/rogue.AutoAttacks.MainhandSwingSpeed().Seconds() +
				getCritChance(rogue.AutoAttacks.OHAuto)/rogue.AutoAttacks.OffhandSwingSpeed().Seconds()
			procChance := []float64{0, 0.33, 0.66, 1}[rogue.Talents.FocusedAttacks]

			eps += critsPerSecond * procChance * 2
		}
		return 10*rogue.EnergyTickMultiplier + eps
	}

	// Glyph of Backstab support
	var bonusDuration float64
	rupRemaining := func(sim *core.Simulation) time.Duration {
		if dot := rogue.Rupture.CurDot(); dot.IsActive() {
			return dot.RemainingDuration(sim)
		}
		return 0
	}

	if x.builder == rogue.Backstab && rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfBackstab) {
		bonusDuration = 6
		rupRemaining = func(sim *core.Simulation) time.Duration {
			if dot := rogue.Rupture.CurDot(); dot.IsActive() {
				dur := dot.RemainingDuration(sim)
				dur += dot.TickLength * time.Duration(dot.MaxStacks+3-dot.NumberOfTicks)
				return dur
			}
			return 0
		}
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

			if sndDur := rogue.SliceAndDiceAura.RemainingDuration(sim); sndDur > 0 {
				if cp == 5 { // pool for snd if pooling for rupture fails
					rupDur := rupRemaining(sim)
					if e+rupDur.Seconds()*energyPerSecond() > maxPool {
						if e+sndDur.Seconds()*energyPerSecond() <= maxPool {
							return Wait
						}
					}
					return Skip
				}

				if cp >= 1 { // don't build if it reduces uptime
					if e+sndDur.Seconds()*energyPerSecond() < sndCost+bldCost || sndDur < time.Second {
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
			if cp < 1 && e >= bldCost {
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
						if rogue.CurrentEnergy() >= bldCost {
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
					cpGenerated := energyGained / bldCost
					currentCp := float64(rogue.ComboPoints())
					if currentCp+cpGenerated > 5 {
						return Skip
					} else {
						if currentCp < 5 {
							if rogue.CurrentEnergy() >= bldCost {
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
	rup4to5 := (rogue.RuptureDuration(4).Seconds() + bonusDuration) * (1 - rogue.RuptureDamage(4)/rogue.RuptureDamage(5))
	rup3to4 := (rogue.RuptureDuration(3).Seconds() + bonusDuration) * (1 - rogue.RuptureDamage(3)/rogue.RuptureDamage(4))

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
				if cp == 4 && e+rup4to5*energyPerSecond() < bldCost+rupCost {
					return Cast
				}
				if cp == 3 && e+rup3to4*energyPerSecond() < bldCost+rupCost {
					return Cast
				}
				if e >= bldCost {
					return Build
				}
				return Wait
			}

			// there's ample time to rebuild, simply skip
			dur := rupRemaining(sim).Seconds()
			if e+dur*baseEps > maxPool {
				return Skip
			}

			if cp == 5 {
				if e+dur*energyPerSecond() > maxPool {
					return Skip // can't pool any longer, maybe we can fit in Eviscerate
				}
				return Wait
			}
			if cp == 4 && e+(dur+rup4to5)*energyPerSecond() < bldCost+rupCost {
				return Wait
			}
			if cp == 3 && e+(dur+rup3to4)*energyPerSecond() < bldCost+rupCost {
				return Wait
			}
			if e >= bldCost {
				return Build
			}
			return Wait
		},
		func(sim *core.Simulation, rogue *Rogue) bool {
			return rogue.Rupture.Cast(sim, rogue.CurrentTarget)
		},
		rupCost,
	})

	bldPerCp := 1.0
	if x.builder == rogue.SinisterStrike {
		attackTable := rogue.AttackTables[rogue.CurrentTarget.UnitIndex]
		crit := rogue.SinisterStrike.PhysicalCritChance(attackTable)
		var extraChance float64
		if rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfSinisterStrike) {
			extraChance = 0.5
		}
		bldPerCp = 1 / (1 + crit*(extraChance+0.2*float64(rogue.Talents.SealFate)))
	}
	if x.builder == rogue.Backstab {
		attackTable := rogue.AttackTables[rogue.CurrentTarget.UnitIndex]
		crit := rogue.Backstab.PhysicalCritChance(attackTable)
		bldPerCp = 1 / (1 + crit*(0.2*float64(rogue.Talents.SealFate)))
	}
	if x.builder == rogue.Hemorrhage {
		attackTable := rogue.AttackTables[rogue.CurrentTarget.UnitIndex]
		crit := rogue.Hemorrhage.PhysicalCritChance(attackTable)
		bldPerCp = 1 / (1 + crit*(0.2*float64(rogue.Talents.SealFate)))
	}
	if x.builder == rogue.Mutilate {
		attackTable := rogue.AttackTables[rogue.CurrentTarget.UnitIndex]
		critMH := rogue.MutilateMH.PhysicalCritChance(attackTable)
		critOH := rogue.MutilateOH.PhysicalCritChance(attackTable)
		crit := 1 - (1-critMH)*(1-critOH)
		bldPerCp = 1 / (2 + crit*(0.2*float64(rogue.Talents.SealFate)))
	}

	// direct damage finisher (Eviscerate/Envenom)
	finisher := rogue.Eviscerate
	if rogue.Talents.MasterPoisoner > 0 {
		finisher = rogue.Envenom
	}
	finisherCost := finisher.DefaultCast.Cost
	x.prios = append(x.prios, prio{
		func(sim *core.Simulation, rogue *Rogue) PriorityAction {
			e, cp := rogue.CurrentEnergy(), rogue.ComboPoints()

			if dur := sim.GetRemainingDuration(); dur <= ruptureMinDuration {
				// end of fight handling - build towards a 3+ cp finisher, or just spam the builder
				switch cp {
				case 5:
					if e >= finisherCost {
						return Cast
					}
					return Wait
				default:
					if e+dur.Seconds()*energyPerSecond() >= bldCost+finisherCost {
						return Build
					}
					if cp >= 3 && e >= finisherCost {
						return Cast
					}
					if cp < 3 && e >= bldCost {
						return Build
					}
				}
				return Wait
			}

			// we only get here if there's ample time left on rupture, or rupture pooling failed: in these cases, we
			// can try to fill in a 5 cp finisher, if it's not too disruptive. lower cp finishers aren't worth it,
			// since builder spam isn't all that much worse
			if cp <= 4 {
				return Build
			}

			cost := finisherCost - 5*rsPerCp + (4-ruthCp)*bldCost*bldPerCp + rupCost

			rupDur := rupRemaining(sim)
			sndDur := rogue.SliceAndDiceAura.RemainingDuration(sim)
			if sndDur < rupDur {
				cost += sndCost - 1*rsPerCp + (1-ruthCp)*bldCost*bldPerCp
			}

			if avail := e + rupDur.Seconds()*energyPerSecond(); avail >= cost {
				return Cast
			}

			// we'd lose a CP here, so we just wait...
			if e <= maxPool {
				return Wait
			}

			// ... and if that doesn't work, allow to clip snd
			if sndDur < rogue.sliceAndDiceDurations[2]-rogue.sliceAndDiceDurations[1] {
				rogue.SliceAndDice.Cast(sim, rogue.CurrentTarget)
				return Wait
			}

			return Build
		},
		func(sim *core.Simulation, rogue *Rogue) bool {
			return finisher.Cast(sim, rogue.CurrentTarget)
		},
		finisherCost,
	})
}

func (x *rotation_generic) run(sim *core.Simulation, rogue *Rogue) {
	for i := 0; i < len(x.prios); i++ {
		switch p := x.prios[i]; p.check(sim, rogue) {
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
