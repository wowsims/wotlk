package rogue

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) OnEnergyGain(sim *core.Simulation) {
	if rogue.KillingSpreeAura.IsActive() {
		rogue.DoNothing()
		return
	}
	rogue.TryUseCooldowns(sim)
	if rogue.GCD.IsReady(sim) {
		rogue.rotation(sim)
	}
}

func (rogue *Rogue) OnGCDReady(sim *core.Simulation) {
	if rogue.KillingSpreeAura.IsActive() {
		rogue.DoNothing()
		return
	}
	rogue.rotation(sim)
}

func (rogue *Rogue) rotation(sim *core.Simulation) {
	numTargets := sim.GetNumTargets()
	var spell *core.Spell
	if numTargets > 1 && numTargets < 4 {
		spell = rogue.multiTargetChooseSpell(sim)
	} else if numTargets > 3 {
		spell = rogue.aoeChooseSpell(sim)
	} else {
		spell = rogue.singleTargetChooseSpell(sim)
	}
	if spell != nil {
		if spell.Cast(sim, rogue.CurrentTarget) {
			rogue.plan = None
		}
	}
	if rogue.GCD.IsReady(sim) {
		rogue.DoNothing()
	}
}

const (
	None Plan = iota
	Slice
	Expose
	Build
	Finish
	Fill
)

func (rogue *Rogue) updatePlan(sim *core.Simulation, refreshThreshold time.Duration) {
	sliceRemaining := RemainingAuraDuration(sim, rogue.SliceAndDiceAura)
	if sliceRemaining <= refreshThreshold {
		rogue.plan = Slice
		return
	}
	exposeRemaining := time.Second * math.MaxInt32
	if rogue.Rotation.MaintainExposeArmor {
		exposeRemaining = RemainingAuraDuration(sim, rogue.ExposeArmorAura)
	}
	if exposeRemaining <= refreshThreshold {
		rogue.plan = Expose
		return
	}
	if rogue.disabledMCDs != nil {
		rogue.EnableAllCooldowns(rogue.disabledMCDs)
		rogue.disabledMCDs = nil
	}
	timeUntilRefresh := core.MinDuration(sliceRemaining, exposeRemaining)
	cp := float64(rogue.ComboPoints())
	expectedCP := cp + rogue.expectedComboPoints(timeUntilRefresh, 0)
	desiredCP := rogue.desiredComboPointsForSlice(sim)
	desiredEnergy := rogue.SliceAndDice[1].DefaultCast.Cost
	if exposeRemaining < sliceRemaining {
		desiredCP = rogue.desiredComboPointsForExpose(sim)
		desiredEnergy = rogue.ExposeArmor[1].DefaultCast.Cost
	}
	energy := rogue.CurrentEnergy()
	expectedEnergy := energy + rogue.expectedEnergyGain(timeUntilRefresh)
	if expectedCP <= desiredCP {
		// We can spend energy but not CP unless it is on refreshing the expiring aura
		if (cp+rogue.BuilderComboPoints) <= desiredCP && (expectedEnergy-desiredEnergy) >= rogue.Builder.DefaultCast.Cost {
			rogue.plan = Build
		} else {
			// We don't have enough energy to build more
			rogue.plan = None
		}
		return
	} else if expectedEnergy <= desiredEnergy {
		// We can spend CP but not energy
		rogue.plan = None
		return
	} else {
		// We can freely spend energy or CP
		var remainingFinisherDuration time.Duration
		if rogue.Rotation.UseEnvenom {
			remainingFinisherDuration = RemainingAuraDuration(sim, rogue.EnvenomAura)
			desiredEnergy = rogue.Envenom[1].DefaultCast.Cost
		} else {
			remainingFinisherDuration = RemainingAuraDuration(sim, rogue.RuptureDot.Aura)
			desiredEnergy = rogue.Rupture[1].DefaultCast.Cost
		}
		durationToConsider := timeUntilRefresh
		if remainingFinisherDuration <= timeUntilRefresh {
			durationToConsider = remainingFinisherDuration
		}
		// Recompute expected cp and energy gains using finisher aura duration
		// How long can we delay the finisher (by filling or building cp)

		timeUntilCapCP := rogue.expectedTimeToComboPoints(5)
		timeUntilCapEnergy := rogue.expectedTimeToEnergy(100)

		// Cases where we needed to do something before the finisher aura
		if timeUntilCapCP < durationToConsider || timeUntilCapEnergy < durationToConsider {
			if cp < 5 {
				rogue.plan = Build
				return
			}
			// TODO: if there is enough energy to fill, cast the finisher, and refresh the slice/expose
			if energy >= 100 {
				rogue.plan = Fill
				return
			}
			rogue.plan = Finish
			return
			// Cases where we will cast the finisher before capping either
		} else {
			// Do we have time to cast another builder?
			// Can we cast another builder and still have the energy for our finisher
			if (energy + rogue.expectedEnergyGain(durationToConsider)) >= (rogue.Builder.DefaultCast.Cost + desiredEnergy) {
				rogue.plan = Build
				return
			}
			rogue.plan = Finish
			return
		}
	}
}

func (rogue *Rogue) chooseSpell(sim *core.Simulation, refreshThreshold time.Duration, filler proto.Rogue_Rotation_Filler) *core.Spell {

	if rogue.plan == None {
		rogue.updatePlan(sim, refreshThreshold)
	}

	cp := rogue.ComboPoints()
	switch rogue.plan {
	case None:
		return nil
	case Slice:
		if cp > 0 {
			return rogue.SliceAndDice[cp]
		} else {
			return rogue.Builder
		}
	case Expose:
		if cp > 0 {
			return rogue.ExposeArmor[cp]
		} else {
			return rogue.Builder
		}
	case Build:
		return rogue.Builder
	case Finish:
		var finisher *core.Spell
		if cp < 3 {
			finisher = rogue.Builder
		} else if rogue.Rotation.UseEnvenom {
			finisher = rogue.Envenom[cp]
		} else {
			finisher = rogue.Rupture[cp]
		}
		return finisher
	case Fill:
		filler := rogue.Builder
		if rogue.Rotation.Filler == proto.Rogue_Rotation_FanOfKnives {
			filler = rogue.FanOfKnives
		} else if cp > 1 && rogue.Rotation.Filler == proto.Rogue_Rotation_Eviscerate {
			filler = rogue.Eviscerate[cp]
		}
		return filler
	}
	return nil
}

func RemainingAuraDuration(sim *core.Simulation, aura *core.Aura) time.Duration {
	remainingDuration := aura.RemainingDuration(sim)
	if remainingDuration < sim.GetRemainingDuration() {
		return remainingDuration
	} else {
		return time.Second * math.MaxInt32
	}
}

func (rogue *Rogue) aoeChooseSpell(sim *core.Simulation) *core.Spell {
	sliceRemaining := RemainingAuraDuration(sim, rogue.SliceAndDiceAura)
	cp := rogue.ComboPoints()
	if sliceRemaining <= 0 {
		if cp <= 0 {
			return rogue.Builder
		}
		return rogue.SliceAndDice[cp]
	}
	if rogue.disabledMCDs != nil {
		rogue.EnableAllCooldowns(rogue.disabledMCDs)
		rogue.disabledMCDs = nil
	}
	return rogue.FanOfKnives
}

func (rogue *Rogue) desiredComboPointsForSlice(sim *core.Simulation) float64 {
	timeRemaining := sim.GetRemainingDuration()
	for i := range rogue.sliceAndDiceDurations {
		if rogue.sliceAndDiceDurations[i] > timeRemaining {
			return float64(i)
		}
	}
	return 5.0
}

func (rogue *Rogue) desiredComboPointsForExpose(sim *core.Simulation) float64 {
	timeRemaining := sim.GetRemainingDuration()
	for i := range rogue.exposeArmorDurations {
		if rogue.exposeArmorDurations[i] > timeRemaining {
			return float64(i)
		}
	}
	return 5.0
}

func (rogue *Rogue) multiTargetChooseSpell(sim *core.Simulation) *core.Spell {
	return rogue.chooseSpell(sim, time.Second*2, proto.Rogue_Rotation_FanOfKnives)
}

func (rogue *Rogue) singleTargetChooseSpell(sim *core.Simulation) *core.Spell {
	return rogue.chooseSpell(sim, 0, rogue.Rotation.Filler)
}

func (rogue *Rogue) expectedTimeToEnergy(energy float64) time.Duration {
	energyNeeded := core.MaxFloat(energy-rogue.CurrentEnergy(), 0)
	return time.Duration(energyNeeded/rogue.PredictedEnergyPerSecond) * time.Second
}

func (rogue *Rogue) expectedTimeToComboPoints(points float64) time.Duration {
	pointsNeeded := core.MaxFloat(0, points-float64(rogue.ComboPoints()))
	castsNeeded := pointsNeeded / rogue.BuilderComboPoints
	energyNeeded := castsNeeded * rogue.Builder.DefaultCast.Cost
	return rogue.expectedTimeToEnergy(energyNeeded)
}

func (rogue *Rogue) expectedEnergyGain(duration time.Duration) float64 {
	return rogue.PredictedEnergyPerSecond * duration.Seconds()
}

func (rogue *Rogue) expectedComboPoints(duration time.Duration, energyUsed float64) float64 {
	builderCasts := core.MinFloat((rogue.CurrentEnergy()+rogue.expectedEnergyGain(duration))/rogue.Builder.DefaultCast.Cost, duration.Seconds())
	return math.Floor(builderCasts) * rogue.BuilderComboPoints
}
