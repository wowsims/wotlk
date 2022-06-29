package tank

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
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
		} else if bear.Rotation.MaintainFaerieFire && bear.ShouldFaerieFire(sim) {
			bear.FaerieFire.Cast(sim, bear.CurrentTarget)
		} else if bear.shouldDemoRoar(sim) {
			bear.DemoralizingRoar.Cast(sim, bear.CurrentTarget)
		} else if bear.Rotation.Swipe == proto.FeralTankDruid_Rotation_SwipeSpam {
			if bear.CanSwipe() {
				bear.Swipe.Cast(sim, bear.CurrentTarget)
			}
		} else if bear.CanMangleBear(sim) {
			bear.Mangle.Cast(sim, bear.CurrentTarget)
		} else if bear.shouldSwipe(sim) {
			bear.Swipe.Cast(sim, bear.CurrentTarget)
		} else if bear.CanLacerate(sim) {
			bear.Lacerate.Cast(sim, bear.CurrentTarget)
		}
	}

	if bear.GCD.IsReady(sim) && !bear.Mangle.IsReady(sim) && bear.Rotation.Swipe != proto.FeralTankDruid_Rotation_SwipeSpam {
		bear.WaitUntil(sim, bear.Mangle.ReadyAt())
	}

	bear.tryQueueMaul(sim)
}

func (bear *FeralTankDruid) shouldSaveLacerateStacks(sim *core.Simulation) bool {
	return bear.LacerateDot.GetStacks() == 5 &&
		bear.LacerateDot.RemainingDuration(sim) <= time.Millisecond*1500
}

func (bear *FeralTankDruid) shouldSwipe(sim *core.Simulation) bool {
	ap := bear.GetStat(stats.AttackPower) + bear.PseudoStats.MobTypeAttackPower + bear.CurrentTarget.PseudoStats.BonusMeleeAttackPower

	return bear.Rotation.Swipe == proto.FeralTankDruid_Rotation_SwipeWithEnoughAP &&
		bear.CanSwipe() &&
		bear.LacerateDot.GetStacks() == 5 &&
		bear.LacerateDot.RemainingDuration(sim) > time.Millisecond*3000 &&
		ap >= float64(bear.Rotation.SwipeApThreshold)
}

func (bear *FeralTankDruid) tryQueueMaul(sim *core.Simulation) {
	if bear.ShouldQueueMaul(sim) {
		bear.QueueMaul(sim)
	}
}

func (bear *FeralTankDruid) shouldDemoRoar(sim *core.Simulation) bool {
	return bear.ShouldDemoralizingRoar(sim, false, bear.Rotation.MaintainDemoralizingRoar)
}
