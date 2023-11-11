package tank

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (bear *FeralTankDruid) OnGCDReady(sim *core.Simulation) {
	bear.doRotation(sim)
}

func (bear *FeralTankDruid) OnAutoAttack(sim *core.Simulation, _ *core.Spell) {
	bear.tryQueueMaul(sim)
}

func (bear *FeralTankDruid) doRotation(sim *core.Simulation) {
	if bear.GCD.IsReady(sim) {
		if bear.shouldSaveLacerateStacks(sim) && bear.Lacerate.CanCast(sim, bear.CurrentTarget) {
			bear.Lacerate.Cast(sim, bear.CurrentTarget)
		} else if bear.shouldDemoRoar(sim) {
			bear.DemoralizingRoar.Cast(sim, bear.CurrentTarget)
		} else if bear.Berserk.IsReady(sim) {
			bear.Berserk.Cast(sim, nil)
		} else if bear.MangleBear.CanCast(sim, bear.CurrentTarget) {
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
			nextAction = max(nextAction, sim.CurrentTime)
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
	lacerateDot := bear.Lacerate.CurDot()
	return lacerateDot.GetStacks() == 5 &&
		lacerateDot.RemainingDuration(sim) <= time.Millisecond*1500
}

func (bear *FeralTankDruid) shouldSwipe(sim *core.Simulation) bool {
	return bear.SwipeBear.CanCast(sim, bear.CurrentTarget) &&
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
	lacerateDot := bear.Lacerate.CurDot()
	return bear.Lacerate.CanCast(sim, bear.CurrentTarget) && ((bear.MangleBear == nil) || (bear.MangleBear.ReadyAt() >= sim.CurrentTime+core.GCDDefault)) && ((lacerateDot.GetStacks() < 5) || (lacerateDot.RemainingDuration(sim) <= time.Duration(bear.Rotation.LacerateTime*float64(time.Second))))
}
