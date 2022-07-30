package enhancement

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

//////////////////////////////////////
// Priority Rotation - Configurable //
//////////////////////////////////////
type PriorityRotation struct {
	BaseRotation
}

func (rotation *PriorityRotation) DoAction(enh *EnhancementShaman, sim *core.Simulation) {
	target := sim.GetTargetUnit(0)
	for i := 1; i <= 9; i++ {
		if rotation.shouldCast(i, enh, sim, target) {
			if !rotation.cast(i, enh, sim, target) {
				enh.WaitForMana(sim, rotation.getCurrentCastCost(i, enh))
			}
			return
		}
	}

	enh.DoNothing()
	return
}

func (rotation *PriorityRotation) shouldCast(spell int, enh *EnhancementShaman, sim *core.Simulation, target *core.Unit) bool {
	switch spell {
	case 1:
		return rotation.shouldCastStormstrikeNoDebuff(enh, sim, target)
	case 2:
		return rotation.shouldCastLightningBoltInstant(enh, sim, target)
	case 3:
		return rotation.shouldCastStormstrike(enh, sim, target)
	case 4:
		return rotation.shouldCastFlameShock(enh, sim, target)
	case 5:
		return rotation.shouldCastLightningBoltWeave(enh, sim, target)
	case 6:
		return rotation.shouldCastEarthShock(enh, sim, target)
	case 7:
		return rotation.shouldCastLightningShield(enh, sim, target)
	case 8:
		return rotation.shouldCastFireNova(enh, sim, target)
	case 9:
		return rotation.shouldCastLavaLash(enh, sim, target)
	default:
		return false
	}
}

func (rotation *PriorityRotation) cast(spell int, enh *EnhancementShaman, sim *core.Simulation, target *core.Unit) bool {
	switch spell {
	case 1:
		return enh.Stormstrike.Cast(sim, target)
	case 2:
		return enh.LightningBolt.Cast(sim, target)
	case 3:
		return enh.Stormstrike.Cast(sim, target)
	case 4:
		return enh.FlameShock.Cast(sim, target)
	case 5:
		return enh.LightningBolt.Cast(sim, target)
	case 6:
		return enh.EarthShock.Cast(sim, target)
	case 7:
		return enh.LightningShield.Cast(sim, target)
	case 8:
		return enh.FireNova.Cast(sim, target)
	case 9:
		return enh.LavaLash.Cast(sim, target)
	default:
		return false
	}
}

func (rotation *PriorityRotation) getCurrentCastCost(spell int, enh *EnhancementShaman) float64 {
	switch spell {
	case 1:
		return enh.Stormstrike.CurCast.Cost
	case 2:
		return enh.LightningBolt.CurCast.Cost
	case 3:
		return enh.Stormstrike.CurCast.Cost
	case 4:
		return enh.FlameShock.CurCast.Cost
	case 5:
		return enh.LightningBolt.CurCast.Cost
	case 6:
		return enh.EarthShock.CurCast.Cost
	case 7:
		return enh.LightningShield.CurCast.Cost
	case 8:
		return enh.FireNova.CurCast.Cost
	case 9:
		return enh.LavaLash.CurCast.Cost
	default:
		return 0
	}
}

func NewPriorityRotation(talents *proto.ShamanTalents) *PriorityRotation {
	pr := new(PriorityRotation)
	return pr
}
