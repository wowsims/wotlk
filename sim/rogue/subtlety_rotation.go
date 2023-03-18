package rogue

import (
	"golang.org/x/exp/slices"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

type subtlety_rotation struct {
	prios []prio

	builder *core.Spell
}

func (x *subtlety_rotation) setup(sim *core.Simulation, rogue *Rogue) {
	x.setSubtletyBuilder(sim, rogue)

	x.prios = x.prios[:0]

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
			if rogue.ComboPoints() < 1 && rogue.CurrentEnergy() > x.builder.DefaultCast.Cost && rogue.getExpectedComboPointsPerSecond() >= 0.7 {
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
						if rogue.CurrentEnergy() >= x.builder.DefaultCast.Cost && rogue.getExpectedComboPointsPerSecond() < 1 {
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
					comboGained := rogue.getExpectedComboPointsPerSecond() * timeLeft.Seconds()
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
			if rogue.allMCDsDisabled {
				for _, mcd := range rogue.GetMajorCooldowns() {
					mcd.Enable()
				}
				rogue.allMCDsDisabled = false
			}
			return Skip
		},
		func(_ *core.Simulation, _ *Rogue) bool {
			return false
		},
		0,
	})

	//Shadowstep
	if rogue.Talents.Shadowstep {
		x.prios = append(x.prios, prio{
			func(sim *core.Simulation, rogue *Rogue) PriorityAction {
				if rogue.Shadowstep.IsReady(sim) {
					// Can we cast Rupture now?
					if !rogue.Rupture.CurDot().IsActive() && rogue.ComboPoints() > 4 && rogue.CurrentEnergy() >= rogue.Rupture.DefaultCast.Cost+rogue.Shadowstep.DefaultCast.Cost {
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
			if sim.GetRemainingDuration().Seconds() < rogue.getExpectedSecondsPerComboPoint() && rogue.ComboPoints() >= 2 && rogue.CurrentEnergy() >= rogue.Eviscerate.DefaultCast.Cost {
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

func (x *subtlety_rotation) run(sim *core.Simulation, rogue *Rogue) {
	for i, p := range x.prios {
		switch p.check(sim, rogue) {
		case Skip:
			continue
		case Build:
			if rogue.ComboPoints() == 4 && rogue.CurrentEnergy() <= rogue.maxEnergy-10 {
				// just wait for HaT proc - if it happens, a finisher will follow and often cost effectively 0 energy,
				// so we add another GCD worth of energy headroom
				break
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
		case Wait:
		}
		break
	}
	rogue.DoNothing()
}

func (x *subtlety_rotation) setSubtletyBuilder(sim *core.Simulation, rogue *Rogue) {
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
	// Ghostly Strike
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
