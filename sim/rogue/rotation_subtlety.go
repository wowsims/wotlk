package rogue

import (
	"golang.org/x/exp/slices"
	"log"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

type rotation_subtlety struct {
	prios []prio

	builder *core.Spell
}

func (x *rotation_subtlety) setup(sim *core.Simulation, rogue *Rogue) {
	x.setSubtletyBuilder(sim, rogue)

	x.prios = x.prios[:0]

	secondsPerComboPoint := func() float64 {
		honorAmongThievesChance := []float64{0, 0.33, 0.66, 1.0}[rogue.Talents.HonorAmongThieves]
		return 1 + 1/(float64(rogue.Options.HonorOfThievesCritRate+100)/100*honorAmongThievesChance)
	}

	comboPointsPerSecond := func() float64 {
		return 1 / secondsPerComboPoint()
	}

	energyPerSecond := func() float64 {
		return 10 * rogue.EnergyTickMultiplier
	}

	if rogue.Rotation.OpenWithPremeditation && rogue.Talents.Premeditation {
		x.prios = append(x.prios, prio{
			func(s *core.Simulation, r *Rogue) PriorityAction {
				if rogue.Premeditation.IsReady(s) {
					return Once
				}
				return Wait
			},
			func(s *core.Simulation, r *Rogue) bool {
				return r.Premeditation.Cast(s, r.CurrentTarget)
			},
			rogue.Premeditation.DefaultCast.Cost,
		})
	}

	if rogue.Rotation.OpenWithShadowstep && rogue.Talents.Shadowstep {
		x.prios = append(x.prios, prio{
			func(s *core.Simulation, r *Rogue) PriorityAction {
				if rogue.CurrentEnergy() > rogue.Shadowstep.DefaultCast.Cost {
					return Once
				}
				return Wait
			},
			func(s *core.Simulation, r *Rogue) bool {
				return rogue.Shadowstep.Cast(sim, rogue.CurrentTarget)
			},
			rogue.Shadowstep.DefaultCast.Cost,
		})
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

	// Slice and Dice
	x.prios = append(x.prios, prio{
		func(sim *core.Simulation, rogue *Rogue) PriorityAction {
			if rogue.SliceAndDiceAura.IsActive() {
				return Skip
			}
			// end of combat handling - prefer Eviscerate over a mostly wasted SnD
			if rogue.ComboPoints() >= 2 && rogue.sliceAndDiceDurations[rogue.ComboPoints()] >= 2*sim.GetRemainingDuration() {
				return Skip
			}
			if rogue.ComboPoints() >= 1 && rogue.CurrentEnergy() > rogue.SliceAndDice.DefaultCast.Cost {
				return Cast
			}
			if rogue.ComboPoints() < 1 && rogue.CurrentEnergy() > x.builder.DefaultCast.Cost && comboPointsPerSecond() >= 0.7 {
				return Wait
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

	// Expose armor
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
						if rogue.CurrentEnergy() >= x.builder.DefaultCast.Cost && comboPointsPerSecond() < 1 {
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
					comboGained := comboPointsPerSecond() * timeLeft.Seconds()
					cpGenerated := energyGained/x.builder.DefaultCast.Cost + comboGained
					currentCP := float64(rogue.ComboPoints())
					if currentCP+cpGenerated > 5 {
						return Skip
					} else {
						if currentCP < 5 {
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

	// Enable CDS
	x.prios = append(x.prios, prio{
		func(sim *core.Simulation, rogue *Rogue) PriorityAction {
			for _, mcd := range rogue.GetMajorCooldowns() {
				mcd.Enable()
			}
			return Once
		},
		func(_ *core.Simulation, _ *Rogue) bool {
			return true
		},
		0,
	})

	//Shadowstep
	if rogue.Talents.Shadowstep {
		x.prios = append(x.prios, prio{
			func(sim *core.Simulation, rogue *Rogue) PriorityAction {
				if rogue.Shadowstep.IsReady(sim) {
					// Can we cast Rupture now?
					if !rogue.Rupture.CurDot().IsActive() && rogue.ComboPoints() >= 5 && rogue.CurrentEnergy() >= rogue.Rupture.DefaultCast.Cost+rogue.Shadowstep.DefaultCast.Cost {
						return Cast
					} else {
						return Skip
					}
				}
				return Skip
			},
			func(sim *core.Simulation, rogue *Rogue) bool {
				return rogue.Shadowstep.Cast(sim, rogue.CurrentTarget)
			},
			rogue.Shadowstep.DefaultCast.Cost,
		})
	}

	const ruptureMinDuration = time.Second * 8 // heuristically, 3-4 Rupture ticks are better DPE than Eviscerate or Envenom

	// Rupture
	x.prios = append(x.prios, prio{
		func(sim *core.Simulation, rogue *Rogue) PriorityAction {
			if rogue.Rupture.CurDot().IsActive() || sim.GetRemainingDuration() < ruptureMinDuration {
				return Skip
			}
			if rogue.ComboPoints() >= 5 && rogue.CurrentEnergy() >= rogue.Rupture.DefaultCast.Cost {
				return Cast
			}
			// don't explicitly wait here, to shorten downtime
			if rogue.ComboPoints() < 5 && rogue.CurrentEnergy() >= x.builder.DefaultCast.Cost+rogue.Rupture.DefaultCast.Cost {
				return Build
			}
			return Wait
		},
		func(sim *core.Simulation, rogue *Rogue) bool {
			return rogue.Rupture.Cast(sim, rogue.CurrentTarget)
		},
		rogue.Rupture.DefaultCast.Cost,
	})

	//Envenom
	if rogue.Rotation.SubtletyFinisherPriority == proto.Rogue_Rotation_SubtletyEnvenom {
		x.prios = append(x.prios, prio{
			func(sim *core.Simulation, rogue *Rogue) PriorityAction {
				if !rogue.DeadlyPoison.CurDot().Aura.IsActive() {
					return Skip
				}
				if rogue.EnvenomAura.IsActive() {
					return Skip
				}
				if rogue.ComboPoints() >= 5 && rogue.CurrentEnergy() >= rogue.Envenom.DefaultCast.Cost {
					return Cast
				}
				if rogue.ComboPoints() < 5 && rogue.CurrentEnergy() >= x.builder.DefaultCast.Cost+rogue.Envenom.DefaultCast.Cost {
					return Build
				}
				return Wait
			},
			func(sim *core.Simulation, rogue *Rogue) bool {
				return rogue.Envenom.Cast(sim, rogue.CurrentTarget)
			},
			rogue.Envenom.DefaultCast.Cost,
		})
	}

	// Eviscerate
	x.prios = append(x.prios, prio{
		func(sim *core.Simulation, rogue *Rogue) PriorityAction {
			// end of combat handling - prefer Eviscerate over Builder, heuristically
			if sim.GetRemainingDuration().Seconds() < secondsPerComboPoint() && rogue.ComboPoints() >= 2 && rogue.CurrentEnergy() >= rogue.Eviscerate.DefaultCast.Cost {
				return Cast
			}
			if rogue.ComboPoints() >= 5 && rogue.CurrentEnergy() >= rogue.Eviscerate.DefaultCast.Cost {
				return Cast
			}
			if rogue.ComboPoints() < 5 && rogue.CurrentEnergy() >= x.builder.DefaultCast.Cost+rogue.Eviscerate.DefaultCast.Cost {
				return Build
			}
			return Wait
		},
		func(sim *core.Simulation, rogue *Rogue) bool {
			return rogue.Eviscerate.Cast(sim, rogue.CurrentTarget)
		},
		rogue.Eviscerate.DefaultCast.Cost,
	})
}

func (x *rotation_subtlety) run(sim *core.Simulation, rogue *Rogue) {
	for i := 0; i < len(x.prios); i++ {
		switch p := x.prios[i]; p.check(sim, rogue) {
		case Skip:
			continue
		case Build:
			if rogue.ComboPoints() == 4 && rogue.CurrentEnergy() <= rogue.maxEnergy-10 {
				// just wait for HaT proc - if it happens, a finisher will follow and often cost effectively 0 energy,
				// so we add another GCD worth of energy headroom
				rogue.DoNothing()
				return
			}

			x.setSubtletyBuilder(sim, rogue)
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

func (x *rotation_subtlety) setSubtletyBuilder(sim *core.Simulation, rogue *Rogue) {
	// Garrote
	if !rogue.Garrote.CurDot().Aura.IsActive() && rogue.ShadowDanceAura.IsActive() && !rogue.PseudoStats.InFrontOfTarget {
		x.builder = rogue.Garrote
		return
	}
	// Ambush
	if rogue.ShadowDanceAura.IsActive() && !rogue.PseudoStats.InFrontOfTarget && rogue.HasDagger(core.MainHand) {
		x.builder = rogue.Ambush
		return
	}
	// Backstab
	if !rogue.Rotation.HemoWithDagger && !rogue.PseudoStats.InFrontOfTarget && rogue.HasDagger(core.MainHand) {
		x.builder = rogue.Backstab
		return
	}
	// Ghostly Strike -- should only be considered when glyphed
	if rogue.Talents.GhostlyStrike && rogue.Rotation.UseGhostlyStrike && rogue.GhostlyStrike.IsReady(sim) {
		x.builder = rogue.GhostlyStrike
		return
	}
	// Hemorrhage
	if rogue.Talents.Hemorrhage {
		x.builder = rogue.Hemorrhage
		return
	}

	// Sinister Strike
	x.builder = rogue.SinisterStrike
}
