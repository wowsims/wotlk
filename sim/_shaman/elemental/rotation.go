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
	if eleShaman.IsUsingAPL {
		return
	}

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
	fnmm      float64
	clmm      float64
	lvbFSWait time.Duration
}

func (rotation *AdaptiveRotation) DoAction(eleShaman *ElementalShaman, sim *core.Simulation) {
	target := eleShaman.CurrentTarget

	shouldTS := false
	cmp := eleShaman.CurrentManaPercent()
	percent := sim.GetRemainingDurationPercent() - 0.1
	if eleShaman.Env.GetNumTargets() > 1 {
		percent = 0.9 // single target we need less mana.
	}
	if cmp < percent {
		shouldTS = true
	}
	if shouldTS && eleShaman.Thunderstorm.IsReady(sim) {
		eleShaman.Thunderstorm.Cast(sim, target)
		return
	}

	fsTime := eleShaman.FlameShock.CurDot().RemainingDuration(sim)
	lvTime := eleShaman.LavaBurst.CD.TimeToReady(sim)
	lvCastTime := eleShaman.ApplyCastSpeed(eleShaman.LavaBurst.DefaultCast.CastTime)
	if fsTime <= 0 && eleShaman.FlameShock.IsReady(sim) {
		if !eleShaman.FlameShock.Cast(sim, target) {
			eleShaman.WaitForMana(sim, eleShaman.FlameShock.CurCast.Cost)
		}
		return
	} else if fsTime > lvCastTime {
		if lvTime <= 0 {
			if !eleShaman.LavaBurst.Cast(sim, target) {
				eleShaman.WaitForMana(sim, eleShaman.LavaBurst.CurCast.Cost)
			}
			return
		} else if lvTime <= rotation.lvbFSWait && fsTime > lvCastTime+lvTime {
			// If we have enough time to wait lvbFSWait and still have FS up, we should just wait to cast LvB.
			eleShaman.WaitUntil(sim, sim.CurrentTime+lvTime)
			return
		}
	}

	fsCD := eleShaman.FlameShock.CD.TimeToReady(sim)
	if fsCD > fsTime {
		fsTime = fsCD
	}
	// If FS will be needed and is ready in < lvbFSWait time, delay.
	if fsTime <= rotation.lvbFSWait {
		eleShaman.WaitUntil(sim, sim.CurrentTime+fsTime)
		return
	}

	if cmp > rotation.clmm && eleShaman.ChainLightning.IsReady(sim) {
		clTime := eleShaman.ApplyCastSpeed(eleShaman.ChainLightning.DefaultCast.CastTime)
		// Only CL if cast time > 1 second than CL or there is more than 1 target.
		if clTime > core.GCDMin || eleShaman.Env.GetNumTargets() > 1 {
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

func (rotation *AdaptiveRotation) Reset(_ *ElementalShaman, sim *core.Simulation) {
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

func NewAdaptiveRotation(options *proto.ElementalShaman_Rotation) *AdaptiveRotation {
	if options.LvbFsWaitMs == 0 {
		options.LvbFsWaitMs = 175
	}
	return &AdaptiveRotation{
		lvbFSWait: time.Duration(options.LvbFsWaitMs) * time.Millisecond,
	}
}

// ################################################################
//
//	MANUAL
//
// ################################################################
type ManualRotation struct {
	options *proto.ElementalShaman_Rotation
}

func (rotation *ManualRotation) DoAction(eleShaman *ElementalShaman, sim *core.Simulation) {
	target := eleShaman.CurrentTarget
	lvbFSWait := time.Duration(rotation.options.LvbFsWaitMs) * time.Millisecond
	shouldTS := false
	cmp := eleShaman.CurrentManaPercent()

	// TODO: expose these percents to let user tweak
	percent := sim.GetRemainingDurationPercent() - 0.1
	if eleShaman.Env.GetNumTargets() > 1 {
		percent = 0.9
	}
	if cmp < percent {
		shouldTS = true
	}
	if shouldTS && rotation.options.UseThunderstorm && eleShaman.Thunderstorm.IsReady(sim) {
		eleShaman.Thunderstorm.Cast(sim, target)
		return
	}

	fsRemain := eleShaman.FlameShock.CurDot().RemainingDuration(sim)
	needFS := fsRemain <= 0
	// Only overwrite if lvb is ready right now.
	if !needFS && rotation.options.OverwriteFlameshock && eleShaman.LavaBurst.CD.TimeToReady(sim) <= core.GCDDefault {
		lvbTime := max(eleShaman.ApplyCastSpeed(eleShaman.LavaBurst.DefaultCast.CastTime), core.GCDMin)
		if fsRemain < lvbTime {
			needFS = true
		}
	}

	allowLvB := true
	if rotation.options.AlwaysCritLvb {
		lvbTime := max(eleShaman.ApplyCastSpeed(eleShaman.LavaBurst.DefaultCast.CastTime), core.GCDMin)
		if fsRemain <= lvbTime {
			allowLvB = false
		}
	}

	shouldCL := rotation.options.UseChainLightning && cmp > (rotation.options.ClMinManaPer/100) && eleShaman.ChainLightning.IsReady(sim)
	clTime := eleShaman.ApplyCastSpeed(eleShaman.ChainLightning.DefaultCast.CastTime)

	// Never cast CL if single target and cast time <= 1 second
	// This is effecively a waste of haste (that is probably temporary.)
	if clTime <= time.Second && eleShaman.Env.GetNumTargets() == 1 {
		shouldCL = false
	}
	lvbCD := eleShaman.LavaBurst.CD.TimeToReady(sim)
	if shouldCL && rotation.options.UseClOnlyGap {
		shouldCL = false
		clCast := max(eleShaman.ApplyCastSpeed(eleShaman.ChainLightning.DefaultCast.CastTime), core.GCDMin)
		// If LvB CD < CL cast time, we should use CL to pass the time until then.
		// Or if FS is about to expire and we didn't cast LvB.
		if fsRemain <= clCast || (lvbCD <= clCast) {
			shouldCL = true
		}
	}

	fsCD := eleShaman.FlameShock.CD.TimeToReady(sim)
	if fsCD > fsRemain {
		fsRemain = fsCD
	}

	// If FS will be needed and is ready in < lvbFSWait time, delay.
	if fsRemain <= lvbFSWait && fsRemain > 0 {
		eleShaman.WaitUntil(sim, sim.CurrentTime+fsRemain)
		return
	}

	// If LvB will be ready in < lvbFSWait time, delay
	if lvbCD <= lvbFSWait && lvbCD > 0 {
		eleShaman.WaitUntil(sim, sim.CurrentTime+lvbCD)
		return
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

func (rotation *ManualRotation) Reset(_ *ElementalShaman, _ *core.Simulation) {
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

func NewManualRotation(options *proto.ElementalShaman_Rotation) *ManualRotation {
	return &ManualRotation{
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
