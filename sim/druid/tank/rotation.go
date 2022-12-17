package tank

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (bear *FeralTankDruid) OnGCDReady(sim *core.Simulation) {
	bear.doRotation(sim)
}

func (bear *FeralTankDruid) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	bear.tryQueueMaul(sim)
}

func (bear *FeralTankDruid) doRotation(sim *core.Simulation) {
	if bear.GCD.IsReady(sim) {
		if bear.shouldSaveLacerateStacks(sim) && bear.CanLacerate(sim) {
			bear.Lacerate.Cast(sim, bear.CurrentTarget)
		} else if bear.shouldDemoRoar(sim) {
			bear.DemoralizingRoar.Cast(sim, bear.CurrentTarget)
		} else if bear.Berserk.IsReady(sim) {
			bear.Berserk.Cast(sim, nil)

			// Bundle Enrage + Barkskin with Berserk
			if bear.Enrage.IsReady(sim) {
				bear.Enrage.Cast(sim, nil)
			}
			if bear.Barkskin.IsReady(sim) {
				bear.Barkskin.Cast(sim, nil)
			}

			bear.UpdateMajorCooldowns()
		} else if bear.CanMangleBear(sim) {
			bear.MangleBear.Cast(sim, bear.CurrentTarget)
		} else if bear.shouldFaerieFire(sim) {
			bear.FaerieFire.Cast(sim, bear.CurrentTarget)
		} else if bear.shouldLacerate(sim) {
			bear.Lacerate.Cast(sim, bear.CurrentTarget)
		} else if bear.shouldSwipe(sim) {
			bear.SwipeBear.Cast(sim, bear.CurrentTarget)
		}
	}

	if bear.GCD.IsReady(sim) {
		nextAction := bear.FaerieFire.ReadyAt()

		if bear.MangleBear == nil {
			bear.WaitUntil(sim, nextAction)
		} else if !bear.MangleBear.IsReady(sim) {
			nextAction = core.MaxDuration(nextAction, sim.CurrentTime)
			nextMangle := bear.MangleBear.ReadyAt()

			if nextMangle < nextAction+time.Second {
				nextAction = nextMangle
			}

			if nextAction > sim.CurrentTime {
				bear.WaitUntil(sim, nextAction)
			}
		}
	}

	bear.tryQueueMaul(sim)
	bear.DoNothing() // means we intionally have no other action if all else fails.
}

func (bear *FeralTankDruid) shouldSaveLacerateStacks(sim *core.Simulation) bool {
	return bear.LacerateDot.GetStacks() == 5 &&
		bear.LacerateDot.RemainingDuration(sim) <= time.Millisecond*1500
}

func (bear *FeralTankDruid) shouldSwipe(sim *core.Simulation) bool {
	return bear.CanSwipeBear() &&
		((bear.MangleBear == nil) || (bear.MangleBear.ReadyAt() >= sim.CurrentTime+core.GCDDefault)) &&
		bear.CurrentRage()-bear.SwipeBear.DefaultCast.Cost >= bear.MaulRageThreshold
}

func (bear *FeralTankDruid) tryQueueMaul(sim *core.Simulation) {
	if bear.ShouldQueueMaul(sim) {
		bear.QueueMaul(sim)
	}
}

func (bear *FeralTankDruid) shouldDemoRoar(sim *core.Simulation) bool {
	return bear.ShouldDemoralizingRoar(sim, false, bear.Rotation.MaintainDemoralizingRoar)
}

func (bear *FeralTankDruid) shouldFaerieFire(sim *core.Simulation) bool {
	return bear.FaerieFire.IsReady(sim) && ((bear.MangleBear == nil) || (bear.MangleBear.ReadyAt() >= sim.CurrentTime+time.Second))
}

func (bear *FeralTankDruid) shouldLacerate(sim *core.Simulation) bool {
	return bear.CanLacerate(sim) && ((bear.MangleBear == nil) || (bear.MangleBear.ReadyAt() >= sim.CurrentTime+core.GCDDefault)) && ((bear.LacerateDot.GetStacks() < 5) || (bear.LacerateDot.RemainingDuration(sim) <= time.Duration(bear.Rotation.LacerateTime*float64(time.Second))))
}
