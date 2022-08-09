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
	rogue.energyPerSecondAvg = rogue.energyPerSecondCalculator()
	if rogue.GCD.IsReady(sim) {
		rogue.rotation(sim)
	}
}

func (rogue *Rogue) OnGCDReady(sim *core.Simulation) {
	rogue.rotation(sim)
}

func (rogue *Rogue) rotation(sim *core.Simulation) {
	var spell *core.Spell
	if sim.GetNumTargets() > 1 {
		spell = rogue.multiTargetChooseSpell(sim)
	} else {
		spell = rogue.singleTargetChooseSpell(sim)
	}
	if spell != nil {
		spell.Cast(sim, rogue.CurrentTarget)
	}
	if rogue.GCD.IsReady(sim) {
		rogue.DoNothing()
	}
}

func (rogue *Rogue) chooseSpell(sim *core.Simulation, refreshThreshold time.Duration, filler proto.Rogue_Rotation_Filler) *core.Spell {
	sliceRemaining := rogue.SliceAndDiceAura.RemainingDuration(sim)
	exposeRemaining := time.Second * math.MaxInt32
	envenomRemaining := time.Second * math.MaxInt32
	ruptureRemaining := rogue.RuptureDot.RemainingDuration(sim)
	if rogue.Rotation.MaintainExposeArmor {
		exposeRemaining = rogue.ExposeArmorAura.RemainingDuration(sim)
	}
	if rogue.Rotation.UseEnvenom {
		envenomRemaining = rogue.EnvenomAura.RemainingDuration(sim)
	}
	cp := rogue.ComboPoints()
	if sliceRemaining <= refreshThreshold {
		if cp > 0 {
			return rogue.SliceAndDice[cp]
		} else {
			return rogue.Builder
		}
	}
	if exposeRemaining <= refreshThreshold {
		if cp > 0 {
			return rogue.ExposeArmor[cp]
		} else {
			return rogue.Builder
		}
	}
	if rogue.disabledMCDs != nil {
		rogue.EnableAllCooldowns(rogue.disabledMCDs)
		rogue.disabledMCDs = nil
	}
	timeUntilMaintenance := core.MinDuration(sliceRemaining, exposeRemaining)
	cpGained := rogue.expectedComboPoints(timeUntilMaintenance)
	if (cp+cpGained) < 5 && rogue.Talents.CutToTheChase < 1 {
		return rogue.Builder
	}
	if rogue.Talents.HungerForBlood && rogue.HungerForBloodAura.RemainingDuration(sim) <= refreshThreshold {
		// TODO : needs to have a bleed
		return rogue.HungerForBlood
	}
	if envenomRemaining <= refreshThreshold {
		if cp > 2 {
			return rogue.Envenom[cp]
		}
	}
	if rogue.Rotation.UseRupture && ruptureRemaining <= refreshThreshold {
		if cp > 2 {
			return rogue.Rupture[cp]
		}
	}
	timeUntilMaintenance = core.MinDuration(timeUntilMaintenance, core.MinDuration(envenomRemaining, ruptureRemaining))
	extraEnergy := rogue.CurrentEnergy() + rogue.expectedEnergyGain(timeUntilMaintenance) - rogue.Builder.DefaultCast.Cost
	switch filler {
	case proto.Rogue_Rotation_FanOfKnives:
		// Will we run out of energy
		if extraEnergy >= rogue.FanOfKnives.DefaultCast.Cost {
			return rogue.FanOfKnives
		}
	case proto.Rogue_Rotation_Eviscerate:
		// Will we run out of energy or combo points
		if cp > 0 && extraEnergy >= rogue.Eviscerate[cp].DefaultCast.Cost {
			return rogue.Eviscerate[cp]
		}
	}
	if cp <= int32(5-rogue.BuilderComboPoints) {
		return rogue.Builder
	}
	return nil
}

func (rogue *Rogue) multiTargetChooseSpell(sim *core.Simulation) *core.Spell {
	return rogue.chooseSpell(sim, time.Second*2, proto.Rogue_Rotation_FanOfKnives)
}

func (rogue *Rogue) singleTargetChooseSpell(sim *core.Simulation) *core.Spell {
	return rogue.chooseSpell(sim, 0, rogue.Rotation.Filler)
}

func (rogue *Rogue) expectedEnergyGain(duration time.Duration) float64 {
	return rogue.energyPerSecondAvg * duration.Seconds()
}

func (rogue *Rogue) expectedComboPoints(duration time.Duration) int32 {
	builderCasts := core.MinFloat((rogue.CurrentEnergy()+rogue.expectedEnergyGain(duration))/rogue.Builder.DefaultCast.Cost, duration.Seconds())
	return int32(math.Floor(builderCasts) * rogue.BuilderComboPoints)
}
