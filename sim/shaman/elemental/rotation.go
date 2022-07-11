package elemental

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (eleShaman *ElementalShaman) GetPresimOptions(_ proto.Player) *core.PresimOptions {
	return eleShaman.rotation.GetPresimOptions()
}

func (eleShaman *ElementalShaman) OnGCDReady(sim *core.Simulation) {
	eleShaman.tryUseGCD(sim)
}

func (eleShaman *ElementalShaman) tryUseGCD(sim *core.Simulation) {
	if eleShaman.TryDropTotems(sim) {
		return
	}

	eleShaman.rotation.DoAction(eleShaman, sim)
	//actionSuccessful := newAction.Cast(sim)
	//if actionSuccessful {
	//	eleShaman.rotation.OnActionAccepted(eleShaman, sim, newAction)
	//} else {
	//	// Only way for a shaman spell to fail is due to mana cost.
	//	// Wait until we have enough mana to cast.
	//	eleShaman.WaitForMana(sim, newAction.GetManaCost())
	//}
}

// Picks which attacks / abilities the Shaman does.
type Rotation interface {
	GetPresimOptions() *core.PresimOptions

	// Returns the action this rotation would like to take next.
	DoAction(*ElementalShaman, *core.Simulation)

	// Returns this rotation to its initial state. Called before each Sim iteration.
	Reset(*ElementalShaman, *core.Simulation)
}

// ################################################################
//                             ADAPTIVE
// ################################################################
type AdaptiveRotation struct {
	// manaTracker common.ManaSpendingRateTracker

	// hasClearcasting bool
	// baseRotation    Rotation // The rotation used most of the time
	// surplusRotation Rotation // The rotation used when we have extra mana

	LB *core.Spell
}

func (rotation *AdaptiveRotation) DoAction(eleShaman *ElementalShaman, sim *core.Simulation) {
	target := sim.GetTargetUnit(0)

	if eleShaman.CurrentManaPercent() < 0.9 && eleShaman.Thunderstorm.IsReady(sim) {
		eleShaman.Thunderstorm.Cast(sim, target)
		return
	}

	// TODO: If lvb CD < FlameShockDot.Duration then we need to cast FS
	fsUp := eleShaman.FlameShockDot.IsActive()
	if !fsUp && eleShaman.FlameShock.IsReady(sim) {
		if !eleShaman.FlameShock.Cast(sim, target) {
			eleShaman.WaitForMana(sim, eleShaman.FlameShock.CurCast.Cost)
		}
		return
	} else if fsUp && eleShaman.LavaBurst.IsReady(sim) {
		if !eleShaman.LavaBurst.Cast(sim, target) {
			eleShaman.WaitForMana(sim, eleShaman.LavaBurst.CurCast.Cost)
		}
		return
	}

	if !eleShaman.LightningBolt.Cast(sim, target) {
		if sim.Log != nil {
			eleShaman.Log(sim, "Failed to cast LB, cost: %0.1f, current mana: %0.1f\n")
		}
		eleShaman.WaitForMana(sim, eleShaman.LightningBolt.CurCast.Cost)
	}

	// rotation.manaTracker.Update(sim, eleShaman.GetCharacter())
}

func (rotation *AdaptiveRotation) Reset(eleShaman *ElementalShaman, sim *core.Simulation) {
	// rotation.manaTracker.Reset()
	// rotation.baseRotation.Reset(eleShaman, sim)
	// rotation.surplusRotation.Reset(eleShaman, sim)
}

func (rotation *AdaptiveRotation) GetPresimOptions() *core.PresimOptions {
	return &core.PresimOptions{
		SetPresimPlayerOptions: func(player *proto.Player) {
			// player.Spec.(*proto.Player_ElementalShaman).ElementalShaman.Rotation.Type = proto.ElementalShaman_Rotation_CLOnClearcast
		},

		OnPresimResult: func(presimResult proto.UnitMetrics, iterations int32, duration time.Duration) bool {
			return true
		},
	}
}

func NewAdaptiveRotation(talents *proto.ShamanTalents) *AdaptiveRotation {
	return &AdaptiveRotation{
		// hasClearcasting: talents.ElementalFocus,
		// manaTracker:     common.NewManaSpendingRateTracker(),
	}
}

// A single action that an Agent can take.
type AgentAction interface {
	GetActionID() core.ActionID

	// TODO: Maybe change this to 'ResourceCost'
	// Amount of mana required to perform the action.
	GetManaCost() float64

	// Do the action. Returns whether the action was successful. An unsuccessful
	// action indicates that the prerequisites, like resource cost, were not met.
	Cast(sim *core.Simulation) bool
}
