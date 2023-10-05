package feral

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (cat *FeralDruid) doAoeRotation(sim *core.Simulation) (bool, time.Duration) {
	rotation := &cat.Rotation

	curEnergy := cat.CurrentEnergy()
	curCp := cat.ComboPoints()
	isClearcast := cat.ClearcastingAura.IsActive()
	simTimeRemain := sim.GetRemainingDuration()
	latencySecs := cat.latency.Seconds()
	shiftCost := cat.CatForm.DefaultCast.Cost

	waitForTf := cat.Talents.Berserk && (cat.TigersFury.ReadyAt() <= cat.BerserkAura.Duration) && (cat.TigersFury.ReadyAt()+time.Second < simTimeRemain-cat.BerserkAura.Duration)
	berserkNow := cat.Berserk.IsReady(sim) && !waitForTf && !isClearcast

	useBuilder := curCp == 0 && (!cat.SavageRoarAura.IsActive() || cat.SavageRoarAura.RemainingDuration(sim) <= time.Second)

	mangleNow := useBuilder && rotation.AoeMangleBuilder
	rakeNow := useBuilder && !rotation.AoeMangleBuilder

	ffThresh := 87.0
	if cat.BerserkAura.IsActive() {
		ffThresh = rotation.BerserkFfThresh
	}
	ffNow := cat.FaerieFire.CanCast(sim, cat.CurrentTarget) && !isClearcast && curEnergy < ffThresh

	if ffNow {
		simTimeSecs := sim.GetRemainingDuration().Seconds()
		maxSwipesWithoutFF := (int)((curEnergy + simTimeSecs*10) / cat.SwipeCat.DefaultCast.Cost)
		numSwipesWithoutFF := min(maxSwipesWithoutFF, int(simTimeSecs)+1)
		numSwipesWithFF := min(maxSwipesWithoutFF+1, int(simTimeSecs))
		ffNow = numSwipesWithFF > numSwipesWithoutFF
	}

	roarNow := curCp >= 1 && (!cat.SavageRoarAura.IsActive() || cat.clipRoar(sim))

	nextFfEnergy := curEnergy + float64((cat.FaerieFire.TimeToReady(sim)+cat.latency)/core.EnergyTickDuration)
	waitForFf := (cat.FaerieFire.TimeToReady(sim) < time.Second-rotation.MaxFfDelay) && (nextFfEnergy < ffThresh) && !isClearcast

	furorCap := min(20.0*float64(cat.Talents.Furor), 85)
	flowershiftEnergy := min(furorCap, 75) - 10*cat.SpellGCD().Seconds() - 20*latencySecs

	flowerEnd := time.Duration(float64(sim.CurrentTime) + float64(cat.SpellGCD()) + (2.5+2*latencySecs)*float64(time.Second))
	flowerFfDelay := flowerEnd - cat.FaerieFire.ReadyAt()
	flowershiftNow := rotation.FlowerWeave && (curEnergy <= flowershiftEnergy) && !isClearcast && !cat.BerserkAura.IsActive() && !cat.tfExpectedBefore(sim, flowerEnd) && flowerFfDelay < rotation.MaxFfDelay

	if flowershiftNow {
		// if we cant cast and get back then abandon flowershift
		if cat.CurrentMana() <= shiftCost+cat.GiftOfTheWild.DefaultCast.Cost {
			flowershiftNow = false
			cat.Metrics.MarkOOM(sim)
		}
	}

	if flowershiftNow {
		energyToDump := curEnergy + ((flowerEnd - sim.CurrentTime).Seconds() * 10)
		flowershiftNow = flowerEnd+time.Duration(math.Floor(energyToDump/42)*float64(time.Second)) < sim.CurrentTime+simTimeRemain
	}

	pendingPool := PoolingActions{}

	if cat.SavageRoarAura.IsActive() {
		roarCost := core.Ternary(cat.berserkExpectedAt(sim, cat.SavageRoarAura.ExpiresAt()), cat.SavageRoar.DefaultCast.Cost*0.5, cat.SavageRoar.DefaultCast.Cost)
		pendingPool.addAction(cat.SavageRoarAura.ExpiresAt(), roarCost)

		if curCp == 0 && cat.SavageRoarAura.RemainingDuration(sim) > time.Second {
			expireTime := cat.SavageRoarAura.ExpiresAt() - time.Second
			if cat.FaerieFire.TimeToReady(sim) > expireTime-sim.CurrentTime {
				builderCost := core.Ternary(rotation.AoeMangleBuilder, cat.MangleCat.DefaultCast.Cost, cat.Rake.DefaultCast.Cost)
				builderCost = core.Ternary(cat.berserkExpectedAt(sim, expireTime), builderCost*0.5, builderCost)
				pendingPool.addAction(expireTime, builderCost)
			}
		}
	}

	pendingPool.sort()

	floatingEnergy := pendingPool.calcFloatingEnergy(cat, sim)
	excessE := curEnergy - floatingEnergy

	timeToNextAction := time.Duration(0)

	if !cat.CatFormAura.IsActive() && rotation.FlowerWeave {
		// If the previous GotW cast was unsuccessful and we still have
		// leeway available, then try again. Otherwise, shift back into Cat
		// Form.
		if flowershiftNow {
			cat.flowerCast(sim)
		} else {
			cat.readyToShift = true
		}
	} else {
		if ffNow {
			cat.FaerieFire.Cast(sim, cat.CurrentTarget)
			return false, 0
		} else if berserkNow {
			cat.Berserk.Cast(sim, nil)
			cat.UpdateMajorCooldowns()
			return false, 0
		} else if roarNow {
			if cat.SavageRoar.CanCast(sim, cat.CurrentTarget) {
				cat.SavageRoar.Cast(sim, nil)
				return false, 0
			}
			timeToNextAction = time.Duration((cat.CurrentSavageRoarCost() - curEnergy) * float64(core.EnergyTickDuration))
		} else if mangleNow && !waitForFf {
			if cat.MangleCat.CanCast(sim, cat.CurrentTarget) {
				cat.MangleCat.Cast(sim, cat.CurrentTarget)
				return false, 0
			}
			timeToNextAction = time.Duration((cat.CurrentMangleCatCost() - curEnergy) * float64(core.EnergyTickDuration))
		} else if rakeNow && !waitForFf {
			if cat.Rake.CanCast(sim, cat.CurrentTarget) {
				cat.Rake.Cast(sim, cat.CurrentTarget)
				return false, 0
			}
			timeToNextAction = time.Duration((cat.CurrentRakeCost() - curEnergy) * float64(core.EnergyTickDuration))
		} else if flowershiftNow && curEnergy < 42 {
			cat.readyToGift = true
		} else {
			if excessE > cat.CurrentSwipeCatCost() || isClearcast {
				cat.SwipeCat.Cast(sim, cat.CurrentTarget)
				return false, 0
			}
			timeToNextAction = time.Duration((cat.CurrentSwipeCatCost() - excessE) * float64(core.EnergyTickDuration))
		}
	}

	// Model in latency when waiting on Energy for our next action
	nextAction := sim.CurrentTime + timeToNextAction
	paValid, rt := pendingPool.nextRefreshTime()
	if paValid {
		nextAction = min(nextAction, rt)
	}

	return true, nextAction
}
