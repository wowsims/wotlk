package rogue

import (
	"log"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

type subtlety_rotation struct {
	prios []prio
}

func (x *subtlety_rotation) setup(sim *core.Simulation, rogue *Rogue) {
	x.setSubtletyBuilder(sim, rogue)

	x.prios = x.prios[:0]

	if rogue.Rotation.OpenWithPremeditation && rogue.Talents.Premeditation {
		hasCastPremeditation := false
		x.prios = append(x.prios, prio{
			func(s *core.Simulation, r *Rogue) PriorityAction {
				if hasCastPremeditation {
					return Skip
				}
				if rogue.Premeditation.IsReady(s) {
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
			rogue.Premeditation.DefaultCast.Cost,
		})
	}

	if rogue.Rotation.OpenWithShadowstep && rogue.Talents.Shadowstep {
		hasCastShadowstep := false
		x.prios = append(x.prios, prio{
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
		x.prios = append(x.prios, prio{
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
			if rogue.ComboPoints() < 1 && rogue.CurrentEnergy() > rogue.Builder.DefaultCast.Cost && rogue.getExpectedComboPointsPerSecond() >= 0.7 {
				return Wait
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
						if rogue.CurrentEnergy() >= rogue.Builder.DefaultCast.Cost && rogue.getExpectedComboPointsPerSecond() < 1 {
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
			if rogue.ComboPoints() < 5 && rogue.CurrentEnergy() >= rogue.Builder.DefaultCast.Cost+rogue.Rupture.DefaultCast.Cost {
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
				if rogue.ComboPoints() < 5 && rogue.CurrentEnergy() >= rogue.Builder.DefaultCast.Cost+rogue.Envenom.DefaultCast.Cost {
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
			if rogue.ComboPoints() < 5 && rogue.CurrentEnergy() >= rogue.Builder.DefaultCast.Cost+rogue.Eviscerate.DefaultCast.Cost {
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
	if !rogue.GCD.IsReady(sim) {
		log.Panic("gcd not ready")
	}

	prioIndex := 0
	for prioIndex < len(x.prios) {
		prio := x.prios[prioIndex]
		switch prio.check(sim, rogue) {
		case Skip:
			prioIndex += 1
		case Build:
			if rogue.ComboPoints() == 4 && rogue.CurrentEnergy() <= rogue.maxEnergy-10 {
				// just wait for HaT proc - if it happens, a finisher will follow and often cost effectively 0 energy,
				// so we add another GCD worth of energy headroom
				rogue.DoNothing()
				return
			}

			if rogue.GCD.IsReady(sim) {
				x.setSubtletyBuilder(sim, rogue)
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

func (x *subtlety_rotation) setSubtletyBuilder(sim *core.Simulation, rogue *Rogue) {
	// Garrote
	if !rogue.Garrote.CurDot().Aura.IsActive() && rogue.ShadowDanceAura.IsActive() && !rogue.PseudoStats.InFrontOfTarget {
		rogue.Builder = rogue.Garrote
		rogue.BuilderPoints = 1
		return
	}
	// Ambush
	if rogue.ShadowDanceAura.IsActive() && !rogue.PseudoStats.InFrontOfTarget && rogue.HasDagger(core.MainHand) {
		rogue.Builder = rogue.Ambush
		rogue.BuilderPoints = 2
		return
	}
	// Backstab
	if !rogue.Rotation.HemoWithDagger && !rogue.PseudoStats.InFrontOfTarget && rogue.HasDagger(core.MainHand) {
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
		return
	}

	// Sinister Strike
	rogue.Builder = rogue.SinisterStrike
	rogue.BuilderPoints = 1
}
