package elemental

import (
	"time"

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
//
//	ADAPTIVE
//
// ################################################################
type AdaptiveRotation struct {
	fnmm float64
	clmm float64
}

func (rotation *AdaptiveRotation) DoAction(eleShaman *ElementalShaman, sim *core.Simulation) {
	target := eleShaman.CurrentTarget

	shouldTS := false
	cmp := eleShaman.CurrentManaPercent()
	percent := 0.55
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
	}

	if cmp > rotation.clmm && eleShaman.ChainLightning.IsReady(sim) {
		lbTime := eleShaman.ApplyCastSpeed(eleShaman.LightningBolt.DefaultCast.CastTime)
		// Only CL if LB is slower than CL or there is more than 1 target.
		if lbTime > time.Second || len(eleShaman.Env.Encounter.Targets) > 1 {
			if !eleShaman.ChainLightning.Cast(sim, target) {
				eleShaman.WaitForMana(sim, eleShaman.ChainLightning.CurCast.Cost)
			}
			return
		}
	} else if cmp > rotation.fnmm && eleShaman.FireNova.IsReady(sim) {
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
	rotation.fnmm = 1.0
	rotation.clmm = 1.0
	if len(sim.Encounter.Targets) > 4 {
		// 5+ targets FN is better
		rotation.fnmm = 0.33
		// Allow CL as long as you have decent mana (leaving most mana for FN)
		rotation.clmm = 0.5
	} else if len(sim.Encounter.Targets) == 4 {
		// 4 targets, enable both similar prio, prob looking at real AoE now (short fight)
		rotation.clmm = 0.33
		rotation.fnmm = 0.33
	} else if len(sim.Encounter.Targets) == 3 {
		// 3 targets, enable both, but prio CL (more efficient)
		//  Still trying to be very mana efficient as 3 targets
		//  is still often a "boss fight" and could be long.
		rotation.clmm = 0.33
		rotation.fnmm = 0.66
	} else if len(sim.Encounter.Targets) == 2 {
		// enable CL with 2
		rotation.clmm = 0.33
	}
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
//
//	MANUAL
//
// ################################################################
type ManualRotation struct {
	// talents
	options *proto.ElementalShaman_Rotation
}

func (rotation *ManualRotation) DoAction(eleShaman *ElementalShaman, sim *core.Simulation) {
	target := eleShaman.CurrentTarget

	shouldTS := false
	cmp := eleShaman.CurrentManaPercent()

	// TODO: expose these percents to let user tweak
	percent := 0.55
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
		if fsRemain < lvbTime {
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
	lbTime := eleShaman.ApplyCastSpeed(eleShaman.LightningBolt.DefaultCast.CastTime)

	// Never cast CL if single target and LB cast time == CL cast time.
	if lbTime <= time.Second && len(eleShaman.Env.Encounter.Targets) == 1 {
		shouldCL = false // never CL if your LB is just as fast.
	}
	if shouldCL && rotation.options.UseClOnlyGap {
		shouldCL = false
		lvbCD := eleShaman.LavaBurst.CD.TimeToReady(sim)
		clCast := core.MaxDuration(eleShaman.ApplyCastSpeed(eleShaman.ChainLightning.DefaultCast.CastTime), core.GCDMin)
		// If LvB CD < CL cast time, we should use CL to pass the time until then.
		// Or if FS is about to expire and we didn't cast LvB.
		if fsRemain <= clCast || (lvbCD <= clCast) {
			shouldCL = true
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
