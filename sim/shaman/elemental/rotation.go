package elemental

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// func (eleShaman *ElementalShaman) GetPresimOptions(_ proto.Player) *core.PresimOptions {
// 	return eleShaman.rotation.GetPresimOptions()
// }

func (eleShaman *ElementalShaman) OnGCDReady(sim *core.Simulation) {
	eleShaman.tryUseGCD(sim)
}

func (eleShaman *ElementalShaman) tryUseGCD(sim *core.Simulation) {
	if eleShaman.TryDropTotems(sim) {
		return
	}

	eleShaman.rotation.DoAction(eleShaman, sim)
}

// Picks which attacks / abilities the Shaman does.
type Rotation interface {
	// GetPresimOptions() *core.PresimOptions

	// Returns the action this rotation would like to take next.
	DoAction(*ElementalShaman, *core.Simulation)

	// Returns this rotation to its initial state. Called before each Sim iteration.
	Reset(*ElementalShaman, *core.Simulation)
}

// ################################################################
//                             ADAPTIVE
// ################################################################
type AdaptiveRotation struct {
}

func (rotation *AdaptiveRotation) DoAction(eleShaman *ElementalShaman, sim *core.Simulation) {
	target := sim.GetTargetUnit(0)

	if eleShaman.CurrentManaPercent() < 0.9 && eleShaman.Thunderstorm.IsReady(sim) {
		eleShaman.Thunderstorm.Cast(sim, target)
		return
	}

	fsExpires := eleShaman.FlameShockDot.ExpiresAt()
	if fsExpires < sim.CurrentTime+eleShaman.LavaBurst.CurCast.CastTime && eleShaman.FlameShock.IsReady(sim) {
		if !eleShaman.FlameShock.Cast(sim, target) {
			eleShaman.WaitForMana(sim, eleShaman.FlameShock.CurCast.Cost)
		}
		return
	} else if eleShaman.LavaBurst.IsReady(sim) { // probably better to cast LvB if we can't shock than cast LB
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
}

func (rotation *AdaptiveRotation) Reset(eleShaman *ElementalShaman, sim *core.Simulation) {
}

// func (rotation *AdaptiveRotation) GetPresimOptions() *core.PresimOptions {
// 	return &core.PresimOptions{
// 		SetPresimPlayerOptions: func(player *proto.Player) {
// 		},

// 		OnPresimResult: func(presimResult proto.UnitMetrics, iterations int32, duration time.Duration) bool {
// 			return true
// 		},
// 	}
// }

func NewAdaptiveRotation(talents *proto.ShamanTalents) *AdaptiveRotation {
	return &AdaptiveRotation{}
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
