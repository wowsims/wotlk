package retribution

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (ret *RetributionPaladin) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	if ret.SealOfVengeanceAura.IsActive() && core.MinInt32(ret.MaxSoVTargets, ret.Env.GetNumTargets()) > 1 {
		minVengeanceDotDuration := time.Second * 15
		var minVengeanceDotDurationTarget *core.Unit
		minVengeanceDotStacks := int32(5)
		var minVengeanceDotStacksTarget *core.Unit
		for i := int32(0); i < core.MinInt32(ret.MaxSoVTargets, ret.Env.GetNumTargets()); i++ {
			target := ret.Env.GetTargetUnit(i)
			dot := ret.SovDotSpell.Dot(target)
			remainingDuration := dot.RemainingDuration(sim)
			stackCount := dot.GetStacks()

			if remainingDuration < minVengeanceDotDuration && remainingDuration > 0 {
				minVengeanceDotDuration = remainingDuration
				minVengeanceDotDurationTarget = target
			}

			if stackCount < minVengeanceDotStacks {
				minVengeanceDotStacks = stackCount
				minVengeanceDotStacksTarget = target
			}
		}

		if minVengeanceDotDuration < ret.WeaponFromMainHand(0).SwingDuration*2 {
			ret.CurrentTarget = minVengeanceDotDurationTarget
		} else if ret.SovDotSpell.Dot(ret.CurrentTarget).GetStacks() == 5 && minVengeanceDotStacks < 5 {
			ret.CurrentTarget = minVengeanceDotStacksTarget
		} else {
			ret.CurrentTarget = ret.Env.Encounter.TargetUnits[0]
		}
	}
}

func (ret *RetributionPaladin) OnGCDReady(sim *core.Simulation) {
	ret.SelectedRotation(sim)

	if ret.GCD.IsReady(sim) {
		ret.DoNothing() // this means we had nothing to do and we are ok
	}
}

func (ret *RetributionPaladin) customRotation(sim *core.Simulation) {
	// Setup
	target := ret.Env.Encounter.TargetUnits[0]

	nextSwingAt := ret.AutoAttacks.NextAttackAt()
	isExecutePhase := sim.IsExecutePhase20()

	if ret.HandOfReckoning != nil && ret.HandOfReckoning.IsReady(sim) {
		ret.HandOfReckoning.Cast(sim, ret.CurrentTarget)
	}

	if ret.GCD.IsReady(sim) {
	rotationLoop:
		for _, spell := range ret.RotationInput {
			if spell == ret.HammerOfWrath && !isExecutePhase {
				continue
			}

			if spell == ret.HammerOfWrath && isExecutePhase && ret.HoldLastAvengingWrathUntilExecution {
				if ret.AvengingWrath.IsReady(sim) {
					success := ret.AvengingWrath.Cast(sim, target)
					if !success {
						ret.WaitForMana(sim, ret.AvengingWrath.CurCast.Cost)
					}
				}
			}

			if spell == ret.Consecration && !ret.checkConsecrationClipping(sim) {
				// This is a skip, so we take the opposite of the clip check.
				continue
			}

			if spell == ret.Exorcism && !ret.ArtOfWarInstantCast.IsActive() {
				continue
			}

			if spell == ret.DivinePlea && ret.CurrentMana() > (ret.MaxMana()*ret.DivinePleaPercentage) {
				continue
			}

			if spell.IsReady(sim) {
				success := spell.Cast(sim, target)
				if !success {
					ret.WaitForMana(sim, spell.CurCast.Cost)
				}
				break rotationLoop
			}
		}
	}

	// All possible next events
	events := []time.Duration{
		nextSwingAt,
		ret.GCD.ReadyAt(),
		ret.SelectedJudgement.CD.ReadyAt(),
		ret.DivineStorm.CD.ReadyAt(),
		ret.HammerOfWrath.CD.ReadyAt(),
		ret.HolyWrath.CD.ReadyAt(),
		ret.CrusaderStrike.CD.ReadyAt(),
		ret.Consecration.CD.ReadyAt(),
		ret.Exorcism.CD.ReadyAt(),
		ret.DivinePlea.CD.ReadyAt(),
	}

	if ret.HandOfReckoning != nil {
		events = append(events, ret.HandOfReckoning.CD.ReadyAt())
	}

	ret.waitUntilNextEvent(sim, events, ret.customRotation)

}

func (ret *RetributionPaladin) castSequenceRotation(sim *core.Simulation) {
	if len(ret.RotationInput) == 0 {
		return
	}

	// Setup
	target := ret.Env.Encounter.TargetUnits[0]
	isExecutePhase := sim.IsExecutePhase20()

	nextReadyAt := sim.CurrentTime

	if hc := ret.Hardcast; ret.HandOfReckoning != nil && ret.HandOfReckoning.IsReady(sim) && !(hc.Expires > sim.CurrentTime) {
		ret.HandOfReckoning.Cast(sim, ret.CurrentTarget)
	}

	if ret.GCD.IsReady(sim) {
		if ret.UseDivinePlea && ret.DivinePlea.IsReady(sim) && ret.CurrentMana() < (ret.MaxMana()*ret.DivinePleaPercentage) {
			ret.DivinePlea.Cast(sim, nil)
		} else {
			currentSpell := ret.RotationInput[ret.CastSequenceIndex]

			if currentSpell == ret.HammerOfWrath && !isExecutePhase {
				return
			}

			if currentSpell.IsReady(sim) {
				success := currentSpell.Cast(sim, target)
				if success {
					ret.CastSequenceIndex = (ret.CastSequenceIndex + 1) % int32(len(ret.RotationInput))
				} else {
					ret.WaitForMana(sim, currentSpell.CurCast.Cost)
				}
			} else {
				nextReadyAt = currentSpell.ReadyAt()
			}
		}
	}

	events := []time.Duration{
		ret.GCD.ReadyAt(),
		nextReadyAt,
	}

	if ret.HandOfReckoning != nil {
		events = append(events, ret.HandOfReckoning.CD.ReadyAt())
	}

	ret.waitUntilNextEvent(sim, events, ret.castSequenceRotation)
}

func (ret *RetributionPaladin) mainRotation(sim *core.Simulation) {

	// Setup
	target := ret.Env.Encounter.TargetUnits[0]

	nextSwingAt := ret.AutoAttacks.NextAttackAt()
	isExecutePhase := sim.IsExecutePhase20()

	nextPrimaryAbility := core.MinDuration(ret.CrusaderStrike.CD.ReadyAt(), ret.DivineStorm.CD.ReadyAt())
	nextPrimaryAbility = core.MinDuration(nextPrimaryAbility, ret.SelectedJudgement.CD.ReadyAt())
	nextPrimaryAbilityDelta := nextPrimaryAbility - sim.CurrentTime

	if ret.HandOfReckoning != nil && ret.HandOfReckoning.IsReady(sim) {
		ret.HandOfReckoning.Cast(sim, ret.CurrentTarget)
	}

	if ret.GCD.IsReady(sim) {
		switch {
		case isExecutePhase && ret.HammerOfWrath.IsReady(sim) && ret.HoldLastAvengingWrathUntilExecution:
			if ret.AvengingWrath.IsReady(sim) {
				success := ret.AvengingWrath.Cast(sim, target)
				if !success {
					ret.WaitForMana(sim, ret.AvengingWrath.CurCast.Cost)
				}
			}

			success := ret.HammerOfWrath.Cast(sim, target)
			if !success {
				ret.WaitForMana(sim, ret.HammerOfWrath.CurCast.Cost)
			}
		case ret.SelectedJudgement.IsReady(sim):
			success := ret.SelectedJudgement.Cast(sim, target)
			if !success {
				ret.WaitForMana(sim, ret.SelectedJudgement.CurCast.Cost)
			}
		case ret.HasLightswornBattlegear2Pc && ret.DivineStorm.IsReady(sim):
			success := ret.DivineStorm.Cast(sim, target)
			if !success {
				ret.WaitForMana(sim, ret.DivineStorm.CurCast.Cost)
			}
		case ret.Env.GetNumTargets() > 1 && ret.Consecration.IsReady(sim):
			success := ret.Consecration.Cast(sim, target)
			if !success {
				ret.WaitForMana(sim, ret.Consecration.CurCast.Cost)
			}
		case isExecutePhase && ret.HammerOfWrath.IsReady(sim):
			success := ret.HammerOfWrath.Cast(sim, target)
			if !success {
				ret.WaitForMana(sim, ret.HammerOfWrath.CurCast.Cost)
			}
		case ret.DemonAndUndeadTargetCount >= ret.HolyWrathThreshold && ret.HolyWrath.IsReady(sim):
			success := ret.HolyWrath.Cast(sim, target)
			if !success {
				ret.WaitForMana(sim, ret.HolyWrath.CurCast.Cost)
			}
		case ret.UseDivinePlea && ret.CurrentMana() < (ret.MaxMana()*ret.DivinePleaPercentage) && ret.DivinePlea.IsReady(sim):
			ret.DivinePlea.Cast(sim, nil)
		case ret.CrusaderStrike.IsReady(sim):
			success := ret.CrusaderStrike.Cast(sim, target)
			if !success {
				ret.WaitForMana(sim, ret.CrusaderStrike.CurCast.Cost)
			}
		case ret.DivineStorm.IsReady(sim):
			success := ret.DivineStorm.Cast(sim, target)
			if !success {
				ret.WaitForMana(sim, ret.DivineStorm.CurCast.Cost)
			}
		case (target.MobType == proto.MobType_MobTypeDemon || target.MobType == proto.MobType_MobTypeUndead) &&
			nextPrimaryAbilityDelta.Milliseconds() > int64(ret.ExoSlack) && ret.Exorcism.IsReady(sim) && ret.ArtOfWarInstantCast.IsActive():
			success := ret.Exorcism.Cast(sim, target)
			if !success {
				ret.WaitForMana(sim, ret.Exorcism.CurCast.Cost)
			}
		case nextPrimaryAbilityDelta.Milliseconds() > int64(ret.ConsSlack) && ret.Consecration.IsReady(sim) && ret.checkConsecrationClipping(sim):
			success := ret.Consecration.Cast(sim, target)
			if !success {
				ret.WaitForMana(sim, ret.Consecration.CurCast.Cost)
			}
		case nextPrimaryAbilityDelta.Milliseconds() > int64(ret.ExoSlack) && ret.Exorcism.IsReady(sim) && ret.ArtOfWarInstantCast.IsActive():
			success := ret.Exorcism.Cast(sim, target)
			if !success {
				ret.WaitForMana(sim, ret.Exorcism.CurCast.Cost)
			}
		case ret.DemonAndUndeadTargetCount >= 1 && ret.HolyWrath.IsReady(sim):
			success := ret.HolyWrath.Cast(sim, target)
			if !success {
				ret.WaitForMana(sim, ret.HolyWrath.CurCast.Cost)
			}
		}
	}

	// All possible next events
	events := []time.Duration{
		nextSwingAt,
		ret.GCD.ReadyAt(),
		ret.SelectedJudgement.CD.ReadyAt(),
		ret.DivineStorm.CD.ReadyAt(),
		ret.HammerOfWrath.CD.ReadyAt(),
		ret.HolyWrath.CD.ReadyAt(),
		ret.CrusaderStrike.CD.ReadyAt(),
		ret.Consecration.CD.ReadyAt(),
		ret.Exorcism.CD.ReadyAt(),
		ret.DivinePlea.CD.ReadyAt(),
	}

	if ret.HandOfReckoning != nil {
		events = append(events, ret.HandOfReckoning.CD.ReadyAt())
	}

	ret.waitUntilNextEvent(sim, events, ret.mainRotation)
}

func (ret *RetributionPaladin) checkConsecrationClipping(sim *core.Simulation) bool {
	if ret.AvoidClippingConsecration {
		// TODO: sim.Duration is the wrong value to check
		return sim.CurrentTime+ret.Consecration.AOEDot().TickLength*4 <= sim.Duration
	} else {
		// If we're not configured to check, always return success.
		return true
	}
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
