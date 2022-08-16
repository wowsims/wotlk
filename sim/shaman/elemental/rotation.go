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
	target := eleShaman.CurrentTarget

	shouldTS := false
	cmp := eleShaman.CurrentManaPercent()
	percent := 0.75
	if len(eleShaman.Env.Encounter.Targets) > 1 {
		percent = 0.9 // single target we need less mana.
	}
	if cmp < percent {
		shouldTS = true
	}
	if shouldTS && eleShaman.Thunderstorm.IsReady(sim) {
		eleShaman.Thunderstorm.Cast(sim, target)
		return
	}

	if eleShaman.FlameShockDot.RemainingDuration(sim) <= 0 && eleShaman.FlameShock.IsReady(sim) {
		if !eleShaman.FlameShock.Cast(sim, target) {
			eleShaman.WaitForMana(sim, eleShaman.FlameShock.CurCast.Cost)
		}
		return
	} else if eleShaman.FlameShockDot.RemainingDuration(sim) > eleShaman.ApplyCastSpeed(eleShaman.LavaBurst.DefaultCast.CastTime) && eleShaman.LavaBurst.IsReady(sim) {
		if !eleShaman.LavaBurst.Cast(sim, target) {
			eleShaman.WaitForMana(sim, eleShaman.LavaBurst.CurCast.Cost)
		}
		return
	} else if len(eleShaman.Env.Encounter.Targets) > 1 && cmp > 0.33 && eleShaman.ChainLightning.IsReady(sim) {
		if !eleShaman.ChainLightning.Cast(sim, target) {
			eleShaman.WaitForMana(sim, eleShaman.ChainLightning.CurCast.Cost)
		}
		return
	} else if len(eleShaman.Env.Encounter.Targets) > 1 && cmp > 0.66 && eleShaman.FireNova.IsReady(sim) {
		if !eleShaman.FireNova.Cast(sim, target) {
			eleShaman.WaitForMana(sim, eleShaman.FireNova.CurCast.Cost)
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

// ################################################################
//                             MANUAL
// ################################################################
type ManualRotation struct {
	// talents
	options *proto.ElementalShaman_Rotation
}

func (rotation *ManualRotation) DoAction(eleShaman *ElementalShaman, sim *core.Simulation) {
	target := eleShaman.CurrentTarget

	shouldTS := false
	cmp := eleShaman.CurrentManaPercent()
	percent := 0.75
	if len(eleShaman.Env.Encounter.Targets) > 1 {
		percent = 0.9
	}
	if cmp < percent {
		shouldTS = true
	}
	if shouldTS && rotation.options.UseThunderstorm && eleShaman.Thunderstorm.IsReady(sim) {
		eleShaman.Thunderstorm.Cast(sim, target)
		return
	}

	fsRemain := eleShaman.FlameShockDot.RemainingDuration(sim)
	needFS := fsRemain <= 0
	// Only overwrite if lvb is ready right now.
	if !needFS && rotation.options.OverwriteFlameshock && eleShaman.LavaBurst.CD.TimeToReady(sim) <= 0 {
		lvbTime := core.MaxDuration(eleShaman.ApplyCastSpeed(eleShaman.LavaBurst.DefaultCast.CastTime), core.GCDMin)
		if lvbTime < fsRemain {
			needFS = true
		}
	}

	allowLvB := true
	if rotation.options.AlwaysCritLvb {
		lvbTime := core.MaxDuration(eleShaman.ApplyCastSpeed(eleShaman.LavaBurst.DefaultCast.CastTime), core.GCDMin)
		if fsRemain <= lvbTime {
			allowLvB = false
		}
	}

	shouldCL := rotation.options.UseChainLightning && cmp > (rotation.options.ClMinManaPer/100) && eleShaman.ChainLightning.IsReady(sim)
	if shouldCL && rotation.options.UseClOnlyGap {
		// If LvB CD < CL cast time, we should use CL in its place.
		lvbCD := eleShaman.LavaBurst.CD.TimeToReady(sim)
		clCast := core.MaxDuration(eleShaman.ApplyCastSpeed(eleShaman.ChainLightning.DefaultCast.CastTime), core.GCDMin)
		if lvbCD <= clCast {
			// need to check that FS has enough time to make casting CL worth.
			if rotation.options.AlwaysCritLvb {
				lvbTime := core.MaxDuration(eleShaman.ApplyCastSpeed(eleShaman.LavaBurst.DefaultCast.CastTime), core.GCDMin)
				if fsRemain <= core.MaxDuration(clCast, lvbCD)+lvbTime {
					shouldCL = false
				}
			}
		} else {
			shouldCL = false
		}
	}

	if needFS && eleShaman.FlameShock.IsReady(sim) {
		if !eleShaman.FlameShock.Cast(sim, target) {
			eleShaman.WaitForMana(sim, eleShaman.FlameShock.CurCast.Cost)
		}
		return
	} else if allowLvB && eleShaman.LavaBurst.IsReady(sim) {
		if !eleShaman.LavaBurst.Cast(sim, target) {
			eleShaman.WaitForMana(sim, eleShaman.LavaBurst.CurCast.Cost)
		}
		return
	} else if shouldCL {
		if !eleShaman.ChainLightning.Cast(sim, target) {
			eleShaman.WaitForMana(sim, eleShaman.ChainLightning.CurCast.Cost)
		}
		return
	} else if rotation.options.UseFireNova && cmp > (rotation.options.FnMinManaPer/100) && eleShaman.FireNova.IsReady(sim) {
		if !eleShaman.FireNova.Cast(sim, target) {
			eleShaman.WaitForMana(sim, eleShaman.FireNova.CurCast.Cost)
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

func (rotation *ManualRotation) Reset(eleShaman *ElementalShaman, sim *core.Simulation) {
}

// func (rotation *ManualRotation) GetPresimOptions() *core.PresimOptions {
// 	return &core.PresimOptions{
// 		SetPresimPlayerOptions: func(player *proto.Player) {
// 		},

// 		OnPresimResult: func(presimResult proto.UnitMetrics, iterations int32, duration time.Duration) bool {
// 			return true
// 		},
// 	}
// }

func NewManualRotation(talents *proto.ShamanTalents, options *proto.ElementalShaman_Rotation) *ManualRotation {
	return &ManualRotation{
		// talents: talents,
		options: options,
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
