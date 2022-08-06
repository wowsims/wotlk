package rogue

import (
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

func (rogue *Rogue) OnGCDReady(sim *core.Simulation) {
	rogue.doRotation(sim)
}

func (rogue *Rogue) doAssassinationRotation(sim *core.Simulation) {
	switch rogue.plan {
	case PlanNone:
		rogue.doAssassinationNone(sim)
	case PlanSliceASAP:
		rogue.doPlanSliceASAP(sim)
	case PlanOpener:
		rogue.doPlanSliceASAP(sim)
	}
}

func (rogue *Rogue) doRotation(sim *core.Simulation) {
	if rogue.KillingSpreeAura.IsActive() {
		rogue.DoNothing()
		return
	}
	if rogue.Rotation.UseEnvenom {
		rogue.doAssassinationRotation(sim)
	} else {
		switch rogue.plan {
		case PlanNone:
			rogue.doPlanNone(sim)
		case PlanSliceASAP:
			rogue.doPlanSliceASAP(sim)
		case PlanMaximalSlice:
			rogue.doPlanMaximalSlice(sim)
		case PlanExposeArmor:
			rogue.doPlanExposeArmor(sim)
		case PlanFillBeforeEA:
			rogue.doPlanFillBeforeEA(sim)
		case PlanFillBeforeSND:
			rogue.doPlanFillBeforeSND(sim)
		case PlanOpener:
			rogue.doPlanOpener(sim)
		}
	}

	// If rogue decided to not use GCD, mark exsplicitly this is ok.
	if rogue.GCD.IsReady(sim) {
		rogue.DoNothing()
	}
}

// Opening rotation.
func (rogue *Rogue) doPlanOpener(sim *core.Simulation) {
	// Can add other opener logic here if we want.
	rogue.plan = PlanSliceASAP
	rogue.doPlanSliceASAP(sim)
}

// Cast SND as the next finisher, using no more builders than necessary.
func (rogue *Rogue) doPlanSliceASAP(sim *core.Simulation) {
	if rogue.doneSND {
		rogue.plan = PlanNone
		rogue.doPlanNone(sim)
		return
	}

	energy := rogue.CurrentEnergy()
	comboPoints := rogue.ComboPoints()
	target := rogue.CurrentTarget
	sndTimeRemaining := rogue.SliceAndDiceAura.RemainingDuration(sim)

	if comboPoints > 0 {
		if energy >= SliceAndDiceEnergyCost || rogue.DeathmantleProcAura.IsActive() {
			if rogue.canPoolEnergy(sim, energy) && sndTimeRemaining > time.Second*2 {
				return
			}
			rogue.SliceAndDice[comboPoints].Cast(sim, nil)
			if rogue.disabledMCDs != nil {
				rogue.EnableAllCooldowns(rogue.disabledMCDs)
				rogue.disabledMCDs = nil
			}
			rogue.plan = PlanNone
		}
		return
	} else {
		if energy >= rogue.Builder.DefaultCast.Cost {
			rogue.castBuilder(sim, target)
		}
	}
}

// Get the biggest Slice we can, but still leaving time for EA if necessary.
func (rogue *Rogue) doPlanMaximalSlice(sim *core.Simulation) {
	if rogue.doneSND {
		rogue.plan = PlanNone
		rogue.doPlanNone(sim)
		return
	}

	energy := rogue.CurrentEnergy()
	comboPoints := rogue.ComboPoints()
	target := rogue.CurrentTarget
	sndTimeRemaining := rogue.SliceAndDiceAura.RemainingDuration(sim)

	remainingSimDuration := sim.GetRemainingDuration()
	if rogue.sliceAndDiceDurations[comboPoints] >= remainingSimDuration {
		if energy >= SliceAndDiceEnergyCost || rogue.DeathmantleProcAura.IsActive() {
			if rogue.canPoolEnergy(sim, energy) && sndTimeRemaining > time.Second*2 {
				return
			}
			rogue.SliceAndDice[comboPoints].Cast(sim, nil)
			rogue.plan = PlanNone
		}
		return
	}

	if sndTimeRemaining <= time.Second && comboPoints > 0 {
		if energy >= SliceAndDiceEnergyCost || rogue.DeathmantleProcAura.IsActive() {
			rogue.SliceAndDice[comboPoints].Cast(sim, nil)
			rogue.plan = PlanNone
		}
		return
	}

	if rogue.MaintainingExpose(target) {
		eaTimeRemaining := rogue.ExposeArmorAura.RemainingDuration(sim)
		if rogue.eaBuildTime+buildTimeBuffer > eaTimeRemaining {
			// Cast our slice and start prepping for EA.
			if comboPoints == 0 {
				rogue.plan = PlanExposeArmor
				rogue.doPlanExposeArmor(sim)
				return
			}
			if energy >= SliceAndDiceEnergyCost || rogue.DeathmantleProcAura.IsActive() {
				if rogue.canPoolEnergy(sim, energy) && sndTimeRemaining > time.Second*2 {
					return
				}
				rogue.SliceAndDice[comboPoints].Cast(sim, nil)
				rogue.plan = PlanExposeArmor
				return
			}
		} else {
			if comboPoints == 5 {
				if energy >= SliceAndDiceEnergyCost || rogue.DeathmantleProcAura.IsActive() {
					if rogue.canPoolEnergy(sim, energy) && sndTimeRemaining > time.Second*2 {
						return
					}
					rogue.SliceAndDice[comboPoints].Cast(sim, nil)
					rogue.plan = PlanFillBeforeEA
					return
				}
			} else if energy >= rogue.Builder.DefaultCast.Cost {
				rogue.castBuilder(sim, target)
			}
		}
	} else {
		if comboPoints == 5 {
			if energy >= SliceAndDiceEnergyCost || rogue.DeathmantleProcAura.IsActive() {
				if rogue.canPoolEnergy(sim, energy) && sndTimeRemaining > time.Second*2 {
					return
				}
				rogue.SliceAndDice[comboPoints].Cast(sim, nil)
				rogue.plan = PlanFillBeforeSND
				return
			}
		} else if energy >= rogue.Builder.DefaultCast.Cost {
			rogue.castBuilder(sim, target)
		}
	}
}

// Build towards and cast a 5 pt Expose Armor.
func (rogue *Rogue) doPlanExposeArmor(sim *core.Simulation) {
	if rogue.doneEA {
		rogue.plan = PlanNone
		rogue.doPlanNone(sim)
		return
	}

	energy := rogue.CurrentEnergy()
	comboPoints := rogue.ComboPoints()
	target := rogue.CurrentTarget

	if comboPoints == 5 {
		if energy >= rogue.ExposeArmor.DefaultCast.Cost || rogue.DeathmantleProcAura.IsActive() {
			eaTimeRemaining := rogue.ExposeArmorAura.RemainingDuration(sim)
			if rogue.canPoolEnergy(sim, energy) && eaTimeRemaining > time.Second*2 {
				return
			}
			rogue.ExposeArmor.Cast(sim, target)
			rogue.plan = PlanNone
		}
		return
	} else {
		if energy >= rogue.Builder.DefaultCast.Cost {
			rogue.castBuilder(sim, target)
		}
	}
}

// Optional dps finisher followed by EA.
func (rogue *Rogue) doPlanFillBeforeEA(sim *core.Simulation) {
	energy := rogue.CurrentEnergy()
	comboPoints := rogue.ComboPoints()
	target := rogue.CurrentTarget
	eaTimeRemaining := rogue.ExposeArmorAura.RemainingDuration(sim)

	if rogue.eaBuildTime+buildTimeBuffer > eaTimeRemaining {
		// Cast our finisher and start prepping for EA.
		if comboPoints == 0 {
			rogue.plan = PlanExposeArmor
			rogue.doPlanExposeArmor(sim)
			return
		} else {
			if comboPoints < rogue.Rotation.MinComboPointsForDamageFinisher {
				rogue.plan = PlanExposeArmor
				return
			}
			if rogue.tryUseDamageFinisher(sim, energy, comboPoints) {
				rogue.plan = PlanExposeArmor
				return
			}
		}
	} else {
		if comboPoints == 5 {
			rogue.tryUseDamageFinisher(sim, energy, comboPoints)
		} else if energy >= rogue.Builder.DefaultCast.Cost {
			rogue.castBuilder(sim, target)
		}
	}
}

// Optional dps finisher followed by SND.
func (rogue *Rogue) doPlanFillBeforeSND(sim *core.Simulation) {
	energy := rogue.CurrentEnergy()
	comboPoints := rogue.ComboPoints()
	target := rogue.CurrentTarget
	sndTimeRemaining := rogue.SliceAndDiceAura.RemainingDuration(sim)

	if !rogue.doneSND && rogue.eaBuildTime+buildTimeBuffer > sndTimeRemaining {
		// Cast our finisher and start prepping for SND.
		if comboPoints == 0 {
			rogue.plan = PlanMaximalSlice
			rogue.doPlanMaximalSlice(sim)
			return
		} else {
			if comboPoints < rogue.Rotation.MinComboPointsForDamageFinisher {
				rogue.plan = PlanMaximalSlice
				return
			}
			if rogue.tryUseDamageFinisher(sim, energy, comboPoints) {
				rogue.plan = PlanMaximalSlice
				return
			}
		}
	} else {
		if comboPoints == 5 || (comboPoints > 0 && sim.GetRemainingDuration() < time.Second*2) {
			rogue.tryUseDamageFinisher(sim, energy, comboPoints)
		} else if energy >= rogue.Builder.DefaultCast.Cost {
			rogue.castBuilder(sim, target)
		}
	}
}

func (rogue *Rogue) doAssassinationNone(sim *core.Simulation) {
	energy := rogue.CurrentEnergy()
	if energy < 15 {
		return
	}
	comboPoints := rogue.ComboPoints()
	sndTimeRemaining := rogue.SliceAndDiceAura.RemainingDuration(sim)
	if sndTimeRemaining <= 0 {
		rogue.plan = PlanSliceASAP
		return
	}
	hungerTimeRemaining := rogue.HungerForBloodAura.RemainingDuration(sim)
	if hungerTimeRemaining < time.Second*2 {
		rogue.HungerForBlood.Cast(sim, nil)
		rogue.plan = PlanNone
		return
	}
	envenomTimeRemaining := rogue.GetAura("Envenom").RemainingDuration(sim)
	if envenomTimeRemaining <= time.Second*1 && comboPoints >= rogue.Rotation.MinComboPointsForDamageFinisher {
		if energy >= rogue.Envenom[comboPoints].DefaultCast.Cost {
			rogue.Envenom[comboPoints].Cast(sim, rogue.CurrentTarget)
			rogue.plan = PlanNone
			return
		}
	}
	if comboPoints <= 3 {
		if energy >= rogue.Builder.BaseCost {
			rogue.Builder.Cast(sim, rogue.CurrentTarget)
			rogue.plan = PlanNone
		}
	}
}

func (rogue *Rogue) doPlanNone(sim *core.Simulation) {
	energy := rogue.CurrentEnergy()
	if energy < 25 {
		// No ability costs < 25 energy so just wait.
		return
	}

	comboPoints := rogue.ComboPoints()
	target := rogue.CurrentTarget

	if comboPoints == 0 {
		// No option other than using a builder.
		if energy >= rogue.Builder.DefaultCast.Cost {
			rogue.castBuilder(sim, target)
		}
		return
	}

	sndTimeRemaining := rogue.SliceAndDiceAura.RemainingDuration(sim)

	if !rogue.MaintainingExpose(target) {
		if rogue.doneSND || sndTimeRemaining > rogue.eaBuildTime+buildTimeBuffer {
			rogue.plan = PlanFillBeforeSND
			rogue.doPlanFillBeforeSND(sim)
		} else {
			rogue.plan = PlanMaximalSlice
			rogue.doPlanMaximalSlice(sim)
		}
		return
	}

	eaTimeRemaining := rogue.ExposeArmorAura.RemainingDuration(sim)
	energyForEANext := rogue.Builder.DefaultCast.Cost*float64(5-comboPoints) + rogue.ExposeArmor.DefaultCast.Cost
	eaNextBuildTime := core.MaxDuration(0, time.Duration(((energyForEANext-energy)/rogue.energyPerSecondAvg)*float64(time.Second)))
	spareTime := core.MaxDuration(0, eaTimeRemaining-eaNextBuildTime)
	if spareTime <= buildTimeBuffer {
		rogue.plan = PlanExposeArmor
		rogue.doPlanExposeArmor(sim)
		return
	}

	if sndTimeRemaining == 0 {
		rogue.plan = PlanSliceASAP
		rogue.doPlanSliceASAP(sim)
		return
	}

	if sndTimeRemaining > eaTimeRemaining {
		rogue.plan = PlanFillBeforeEA
		rogue.doPlanFillBeforeEA(sim)
		return
	}

	if rogue.doneSND {
		rogue.plan = PlanFillBeforeSND
		rogue.doPlanFillBeforeSND(sim)
		return
	}

	rogue.plan = PlanMaximalSlice
	rogue.doPlanMaximalSlice(sim)
}

func (rogue *Rogue) canPoolEnergy(sim *core.Simulation, energy float64) bool {
	return sim.GetRemainingDuration() >= time.Second*6 && energy <= 85 && ((rogue.AdrenalineRushAura == nil || !rogue.AdrenalineRushAura.IsActive()) || energy <= 70)
}

func (rogue *Rogue) castBuilder(sim *core.Simulation, target *core.Unit) {
	if rogue.Rotation.UseShiv && rogue.DeadlyPoisonDot.IsActive() && rogue.DeadlyPoisonDot.RemainingDuration(sim) < time.Second*2 && rogue.CurrentEnergy() >= rogue.Shiv.DefaultCast.Cost {
		rogue.Shiv.Cast(sim, target)
	} else {
		rogue.Builder.Cast(sim, target)
	}
}

func (rogue *Rogue) tryUseDamageFinisher(sim *core.Simulation, energy float64, comboPoints int32) bool {
	newRuptureDuration := core.MinDuration(rogue.RuptureDuration(comboPoints), sim.GetRemainingDuration())
	if rogue.RuptureDot.IsActive() {
		newRuptureDuration -= core.MinDuration(rogue.RuptureDot.RemainingDuration(sim), sim.GetRemainingDuration())
	}
	if rogue.Rotation.UseRupture &&
		newRuptureDuration >= time.Second*10 &&
		(sim.GetNumTargets() == 1 || (rogue.BladeFlurryAura == nil || !rogue.BladeFlurryAura.IsActive())) {
		if energy >= RuptureEnergyCost || rogue.DeathmantleProcAura.IsActive() {
			rogue.Rupture[comboPoints].Cast(sim, rogue.CurrentTarget)
		}
		return true
	}
	if rogue.Rotation.UseEnvenom &&
		energy >= rogue.Envenom[comboPoints].DefaultCast.Cost {
		rogue.Envenom[comboPoints].Cast(sim, rogue.CurrentTarget)
		return true
	}
	if energy >= rogue.Eviscerate[comboPoints].DefaultCast.Cost || rogue.DeathmantleProcAura.IsActive() {
		rogue.Eviscerate[comboPoints].Cast(sim, rogue.CurrentTarget)
		return true
	}
	return false
}
