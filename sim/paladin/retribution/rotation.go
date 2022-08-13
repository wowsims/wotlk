package retribution

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (ret *RetributionPaladin) OnGCDReady(sim *core.Simulation) {
	ret.SelectedRotation(sim)

	if ret.GCD.IsReady(sim) {
		ret.DoNothing() // this means we had nothing to do and we are ok
	}
}

func (ret *RetributionPaladin) customRotation(sim *core.Simulation) {
	// Setup
	target := ret.CurrentTarget

	nextSwingAt := ret.AutoAttacks.NextAttackAt()
	isExecutePhase := sim.IsExecutePhase20()

	if ret.GCD.IsReady(sim) {
	rotationLoop:
		for _, spellNumber := range ret.RotationInput {
			switch spellNumber {
			case int32(proto.RetributionPaladin_Rotation_JudgementOfWisdom):
				if ret.JudgementOfWisdom.IsReady(sim) {
					ret.JudgementOfWisdom.Cast(sim, target)
					break rotationLoop
				}
			case int32(proto.RetributionPaladin_Rotation_DivineStorm):
				if ret.DivineStorm.IsReady(sim) {
					ret.DivineStorm.Cast(sim, target)
					break rotationLoop
				}
			case int32(proto.RetributionPaladin_Rotation_HammerOfWrath):
				if isExecutePhase && ret.HammerOfWrath.IsReady(sim) {
					ret.HammerOfWrath.Cast(sim, target)
					break rotationLoop
				}
			case int32(proto.RetributionPaladin_Rotation_Consecration):
				if ret.Consecration.IsReady(sim) {
					ret.Consecration.Cast(sim, target)
					break rotationLoop
				}
			case int32(proto.RetributionPaladin_Rotation_HolyWrath):
				if ret.HolyWrath.IsReady(sim) {
					ret.HolyWrath.Cast(sim, target)
					break rotationLoop
				}
			case int32(proto.RetributionPaladin_Rotation_CrusaderStrike):
				if ret.CrusaderStrike.IsReady(sim) {
					ret.CrusaderStrike.Cast(sim, target)
					break rotationLoop
				}
			case int32(proto.RetributionPaladin_Rotation_Exorcism):
				if ret.Exorcism.IsReady(sim) && ret.ArtOfWarInstantCast.IsActive() {
					ret.Exorcism.Cast(sim, target)
					break rotationLoop
				}
			}
		}
	}

	// All possible next events
	events := []time.Duration{
		nextSwingAt,
		ret.GCD.ReadyAt(),
		ret.JudgementOfWisdom.CD.ReadyAt(),
		ret.DivineStorm.CD.ReadyAt(),
		ret.HammerOfWrath.CD.ReadyAt(),
		ret.HolyWrath.CD.ReadyAt(),
		ret.CrusaderStrike.CD.ReadyAt(),
		ret.Consecration.CD.ReadyAt(),
		ret.Exorcism.CD.ReadyAt(),
	}

	ret.waitUntilNextEvent(sim, events, ret.customRotation)

}

func (ret *RetributionPaladin) castSequenceRotation(sim *core.Simulation) {
	if len(ret.RotationInput) == 0 {
		return
	}

	// Setup
	target := ret.CurrentTarget
	isExecutePhase := sim.IsExecutePhase20()

	nextReadyAt := sim.CurrentTime
	if ret.GCD.IsReady(sim) {
		switch ret.RotationInput[ret.CastSequenceIndex] {
		case int32(proto.RetributionPaladin_Rotation_JudgementOfWisdom):
			if ret.JudgementOfWisdom.IsReady(sim) {
				ret.JudgementOfWisdom.Cast(sim, target)
				ret.CastSequenceIndex = (ret.CastSequenceIndex + 1) % int32(len(ret.RotationInput))
			} else {
				nextReadyAt = ret.JudgementOfWisdom.ReadyAt()
			}
		case int32(proto.RetributionPaladin_Rotation_DivineStorm):
			if ret.DivineStorm.IsReady(sim) {
				ret.DivineStorm.Cast(sim, target)
				ret.CastSequenceIndex = (ret.CastSequenceIndex + 1) % int32(len(ret.RotationInput))
			} else {
				nextReadyAt = ret.DivineStorm.ReadyAt()
			}
		case int32(proto.RetributionPaladin_Rotation_HammerOfWrath):
			if isExecutePhase && ret.HammerOfWrath.IsReady(sim) {
				ret.HammerOfWrath.Cast(sim, target)
				ret.CastSequenceIndex = (ret.CastSequenceIndex + 1) % int32(len(ret.RotationInput))
			} else {
				nextReadyAt = ret.HammerOfWrath.ReadyAt()
			}
		case int32(proto.RetributionPaladin_Rotation_Consecration):
			if ret.Consecration.IsReady(sim) {
				ret.Consecration.Cast(sim, target)
				ret.CastSequenceIndex = (ret.CastSequenceIndex + 1) % int32(len(ret.RotationInput))
			} else {
				nextReadyAt = ret.Consecration.ReadyAt()
			}
		case int32(proto.RetributionPaladin_Rotation_HolyWrath):
			if ret.HolyWrath.IsReady(sim) {
				ret.HolyWrath.Cast(sim, target)
				ret.CastSequenceIndex = (ret.CastSequenceIndex + 1) % int32(len(ret.RotationInput))
			} else {
				nextReadyAt = ret.HolyWrath.ReadyAt()
			}
		case int32(proto.RetributionPaladin_Rotation_CrusaderStrike):
			if ret.CrusaderStrike.IsReady(sim) {
				ret.CrusaderStrike.Cast(sim, target)
				ret.CastSequenceIndex = (ret.CastSequenceIndex + 1) % int32(len(ret.RotationInput))
			} else {
				nextReadyAt = ret.CrusaderStrike.ReadyAt()
			}
		case int32(proto.RetributionPaladin_Rotation_Exorcism):
			if ret.Exorcism.IsReady(sim) {
				ret.Exorcism.Cast(sim, target)
				ret.CastSequenceIndex = (ret.CastSequenceIndex + 1) % int32(len(ret.RotationInput))
			} else {
				nextReadyAt = ret.Exorcism.ReadyAt()
			}
		}
	}

	events := []time.Duration{
		ret.GCD.ReadyAt(),
		nextReadyAt,
	}

	ret.waitUntilNextEvent(sim, events, ret.castSequenceRotation)
}

func (ret *RetributionPaladin) mainRotation(sim *core.Simulation) {

	// Setup
	target := ret.CurrentTarget

	nextSwingAt := ret.AutoAttacks.NextAttackAt()
	isExecutePhase := sim.IsExecutePhase20()

	nextPrimaryAbility := core.MinDuration(ret.CrusaderStrike.CD.ReadyAt(), ret.DivineStorm.CD.ReadyAt())
	nextPrimaryAbility = core.MinDuration(nextPrimaryAbility, ret.JudgementOfWisdom.CD.ReadyAt())
	nextPrimaryAbilityDelta := nextPrimaryAbility - sim.CurrentTime

	if ret.GCD.IsReady(sim) {
		switch {
		case ret.JudgementOfWisdom.IsReady(sim):
			ret.JudgementOfWisdom.Cast(sim, target)
		case ret.HasLightswornBattlegear2Pc && ret.DivineStorm.IsReady(sim):
			ret.DivineStorm.Cast(sim, target)
		case ret.Env.GetNumTargets() == 1 && isExecutePhase && ret.HammerOfWrath.IsReady(sim):
			ret.HammerOfWrath.Cast(sim, target)
		case ret.Env.GetNumTargets() > 1 && ret.Consecration.IsReady(sim):
			ret.Consecration.Cast(sim, target)
		case ret.DemonAndUndeadTargetCount >= ret.HolyWrathThreshold && ret.HolyWrath.IsReady(sim):
			ret.HolyWrath.Cast(sim, target)
		case ret.UseDivinePlea && ret.CurrentMana() < (ret.MaxMana()*ret.DivinePleaPercentage) && ret.DivinePlea.IsReady(sim):
			ret.DivinePlea.Cast(sim, &ret.Unit)
		case ret.CrusaderStrike.IsReady(sim):
			ret.CrusaderStrike.Cast(sim, target)
		case ret.DivineStorm.IsReady(sim):
			ret.DivineStorm.Cast(sim, target)
		case (target.MobType == proto.MobType_MobTypeDemon || target.MobType == proto.MobType_MobTypeUndead) &&
			nextPrimaryAbilityDelta.Milliseconds() > int64(ret.ExoSlack) && ret.Exorcism.IsReady(sim) && ret.ArtOfWarInstantCast.IsActive():
			ret.Exorcism.Cast(sim, target)
		case nextPrimaryAbilityDelta.Milliseconds() > int64(ret.ConsSlack) && ret.Consecration.IsReady(sim):
			ret.Consecration.Cast(sim, target)
		case nextPrimaryAbilityDelta.Milliseconds() > int64(ret.ExoSlack) && ret.Exorcism.IsReady(sim) && ret.ArtOfWarInstantCast.IsActive():
			ret.Exorcism.Cast(sim, target)
		case ret.DemonAndUndeadTargetCount >= 1 && ret.HolyWrath.IsReady(sim):
			ret.HolyWrath.Cast(sim, target)
		}
	}

	// All possible next events
	events := []time.Duration{
		nextSwingAt,
		ret.GCD.ReadyAt(),
		ret.JudgementOfWisdom.CD.ReadyAt(),
		ret.DivineStorm.CD.ReadyAt(),
		ret.HammerOfWrath.CD.ReadyAt(),
		ret.HolyWrath.CD.ReadyAt(),
		ret.CrusaderStrike.CD.ReadyAt(),
		ret.Consecration.CD.ReadyAt(),
		ret.Exorcism.CD.ReadyAt(),
		ret.DivinePlea.CD.ReadyAt(),
	}

	ret.waitUntilNextEvent(sim, events, ret.mainRotation)
}

// Helper function for finding the next event
func (ret *RetributionPaladin) waitUntilNextEvent(sim *core.Simulation, events []time.Duration, rotationCallback func(*core.Simulation)) {
	// Find the minimum possible next event that is greater than the current time
	nextEventAt := time.Duration(math.MaxInt64) // any event will happen before forever.
	for _, elem := range events {
		if elem > sim.CurrentTime && elem < nextEventAt {
			nextEventAt = elem
		}
	}
	// If the next action is  the GCD, just return
	if nextEventAt == ret.GCD.ReadyAt() {
		return
	}

	// Otherwise add a pending action for the next time
	pa := &core.PendingAction{
		Priority:     core.ActionPriorityLow,
		OnAction:     rotationCallback,
		NextActionAt: nextEventAt,
	}

	sim.AddPendingAction(pa)
}
