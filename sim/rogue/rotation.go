package rogue

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

const buildTimeBuffer = time.Second * 0

const (
	PlanNone = iota
	PlanOpener
	PlanExposeArmor
	PlanSliceASAP
	PlanFillBeforeEA
	PlanFillBeforeSND
	PlanMaximalSlice
)

func (rogue *Rogue) OnEnergyGain(sim *core.Simulation) {
	rogue.TryUseCooldowns(sim)
	if rogue.GCD.IsReady(sim) {
		rogue.rotation(sim)
	}
}

func (rogue *Rogue) OnGCDReady(sim *core.Simulation) {
	rogue.rotation(sim)
}

func (rogue *Rogue) rotation(sim *core.Simulation) {
	if rogue.KillingSpreeAura.IsActive() {
		rogue.DoNothing()
		return
	}
	var spell *core.Spell
	if sim.GetNumTargets() > 1 {
		spell = rogue.multiTargetChooseSpell(sim)
	} else {
		spell = rogue.singleTargetChooseSpell(sim)
	}
	spell.Cast(sim, rogue.CurrentTarget)
	if rogue.GCD.IsReady(sim) {
		rogue.DoNothing()
	}
}

func (rogue *Rogue) multiTargetChooseSpell(sim *core.Simulation) *core.Spell {
	return rogue.FanOfKnifes
}

func (rogue *Rogue) singleTargetChooseSpell(sim *core.Simulation) *core.Spell {
	sliceRemaining := rogue.SliceAndDiceAura.RemainingDuration(sim)
	eaRemaining := time.Second * math.MaxInt32
	refreshThreshold := time.Second * 0
	if rogue.Rotation.MaintainExposeArmor {
		eaRemaining = rogue.ExposeArmorAura.RemainingDuration(sim)
	}
	cp := rogue.ComboPoints()
	if sliceRemaining <= refreshThreshold {
		if cp > 0 {
			return rogue.SliceAndDice[cp]
		} else {
			return rogue.Builder
		}
	}
	if rogue.Rotation.MaintainExposeArmor && rogue.ExposeArmorAura.RemainingDuration(sim) <= refreshThreshold {
		if cp > 0 {
			return rogue.ExposeArmor
		} else {
			return rogue.Builder
		}
	}
	if rogue.disabledMCDs != nil {
		rogue.EnableAllCooldowns(rogue.disabledMCDs)
		rogue.disabledMCDs = nil
	}
	cpGained := rogue.expectedComboPoints(core.MinDuration(sliceRemaining, eaRemaining))
	if (cp+cpGained) < 5 && rogue.Talents.CutToTheChase < 1 {
		return rogue.Builder
	}
	if rogue.Talents.HungerForBlood && rogue.HungerForBloodAura.RemainingDuration(sim) <= refreshThreshold {
		// TODO : needs to have a bleed
		return rogue.HungerForBlood
	}
	if rogue.Rotation.UseEnvenom && rogue.EnvenomAura.RemainingDuration(sim) <= refreshThreshold {
		if cp > 2 {
			return rogue.Envenom[cp]
		}
	}
	if rogue.Rotation.UseRupture && rogue.RuptureDot.RemainingDuration(sim) <= refreshThreshold {
		if cp > 0 {
			return rogue.Rupture[cp]
		}
	}
	return rogue.Builder
}

func (rogue *Rogue) expectedComboPoints(duration time.Duration) int32 {
	expectedEnergyGain := rogue.energyPerSecondAvg * duration.Seconds()
	builderCasts := core.MinFloat((rogue.CurrentEnergy()+expectedEnergyGain)/rogue.Builder.DefaultCast.Cost, duration.Seconds())
	return int32(math.Floor(builderCasts) * rogue.BuilderComboPoints)
}
