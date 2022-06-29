package retribution

import (
	"math"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

func (ret *RetributionPaladin) OnGCDReady(sim *core.Simulation) {
	ret.tryUseGCD(sim)
}

func (ret *RetributionPaladin) OnManaTick(sim *core.Simulation) {
	if ret.FinishedWaitingForManaAndGCDReady(sim) {
		ret.tryUseGCD(sim)
	}
}

func (ret *RetributionPaladin) tryUseGCD(sim *core.Simulation) {
	if !ret.openerCompleted {
		ret.openingRotation(sim)
		return
	}
	ret.mainRotation(sim)
}

func (ret *RetributionPaladin) openingRotation(sim *core.Simulation) {
	target := ret.CurrentTarget

	// Cast selected judgement to keep on the boss
	if ret.JudgementOfWisdom.IsReady(sim) &&
		ret.judgement != proto.RetributionPaladin_Options_None {
		var judge *core.Spell
		switch ret.judgement {
		case proto.RetributionPaladin_Options_Wisdom:
			judge = ret.JudgementOfWisdom
		case proto.RetributionPaladin_Options_Crusader:
			judge = ret.JudgementOfTheCrusader
		}
		if judge != nil {
			judge.Cast(sim, target)
		}
	}

	// Cast Seal of Command
	if !ret.SealOfCommandAura.IsActive() {
		ret.SealOfCommand.Cast(sim, nil)
		return
	}

	// Cast Seal of Blood and enable attacks to twist
	if !ret.SealOfBloodAura.IsActive() {
		ret.SealOfBlood.Cast(sim, nil)
		ret.AutoAttacks.EnableAutoSwing(sim)
		ret.openerCompleted = true
	}
}

func (ret *RetributionPaladin) mainRotation(sim *core.Simulation) {
	// Need to check for SoC early
	socActive := ret.SealOfCommandAura.IsActive()

	// If mana is low, do the low mana rotation instead
	// Don't do the low mana rotation in the middle of a twist
	if ret.CurrentMana() <= 1000 && !socActive {
		ret.lowManaRotation(sim)
		return
	}

	// Setup
	target := ret.CurrentTarget

	gcdCD := ret.GCD.TimeToReady(sim)
	crusaderStrikeCD := ret.CrusaderStrike.TimeToReady(sim)
	nextCrusaderStrikeCD := ret.CrusaderStrike.CD.ReadyAt()
	judgementCD := ret.JudgementOfWisdom.TimeToReady(sim)

	sobActive := ret.SealOfBloodAura.IsActive()

	nextSwingAt := ret.AutoAttacks.NextAttackAt()
	timeTilNextSwing := nextSwingAt - sim.CurrentTime
	weaponSpeed := ret.AutoAttacks.MainhandSwingSpeed()

	spellGCD := ret.SpellGCD()

	inTwistWindow := (sim.CurrentTime >= nextSwingAt-twistWindow) && (sim.CurrentTime < ret.AutoAttacks.NextAttackAt())
	latestTwistStart := nextSwingAt - spellGCD
	possibleTwist := timeTilNextSwing > spellGCD+gcdCD
	willTwist := possibleTwist && (nextSwingAt+spellGCD <= nextCrusaderStrikeCD+ret.crusaderStrikeDelay)

	// Use Judgement if we will prep Seal of Command
	// TO-DO: Add more aggressive judgment logic
	// Should judge on crusader strike swings as well if we have enough time to refresh seal
	if judgementCD == 0 && sobActive && willTwist {
		ret.JudgementOfBlood.Cast(sim, target)
		sobActive = false
	}

	// Judgement can affect active seals and CDs
	nextJudgementCD := ret.JudgementOfWisdom.CD.ReadyAt()

	if gcdCD == 0 {
		if socActive && inTwistWindow {
			// If Seal of Command is Active, complete the twist
			ret.SealOfBlood.Cast(sim, nil)
		} else if crusaderStrikeCD == 0 && !willTwist &&
			(sobActive || spellGCD < timeTilNextSwing) {
			// Cast Crusader Strike if we won't swing naked and we aren't twisting
			ret.CrusaderStrike.Cast(sim, target)
		} else if willTwist && !socActive && (nextJudgementCD > latestTwistStart) {
			// Prep seal of command
			ret.SealOfCommand.Cast(sim, nil)
		} else if !sobActive && !socActive && !willTwist {
			// If no seal is active, cast Seal of Blood
			ret.SealOfBlood.Cast(sim, nil)
		} else if !willTwist && !socActive &&
			timeTilNextSwing+weaponSpeed > spellGCD*2 &&
			spellGCD < crusaderStrikeCD {
			// If there is literally nothing else to-do, cast fillers
			// Only if it won't clip crusader strike or seal twist
			ret.useFillers(sim, target)
		}
	}

	// All possible next events
	events := []time.Duration{
		nextSwingAt,
		nextSwingAt - twistWindow,
		ret.GCD.ReadyAt(),
		ret.JudgementOfWisdom.CD.ReadyAt(),
		ret.CrusaderStrike.CD.ReadyAt(),
	}

	ret.waitUntilNextEvent(sim, events)
}

//
func (ret *RetributionPaladin) useFillers(sim *core.Simulation, target *core.Unit) {

	// If the target is a demon and exorcism is up, cast exorcism
	// Only cast exorcism when above 40% mana
	if ret.Rotation.UseExorcism &&
		ret.CanExorcism(target) &&
		ret.Exorcism.IsReady(sim) &&
		ret.CurrentMana() > ret.MaxMana()*0.4 {

		ret.Exorcism.Cast(sim, target)
		return
	}

	// If we can't exorcise, try to consecrate
	// Only cast consecration when above 60% mana
	if ret.Rotation.ConsecrationRank != proto.RetributionPaladin_Rotation_None &&
		ret.Consecration.IsReady(sim) &&
		ret.CurrentMana() > ret.MaxMana()*0.6 {
		ret.Consecration.Cast(sim, target)
		return
	}
}

// Just roll seal of blood and cast crusader strike on CD to conserve mana
func (ret *RetributionPaladin) lowManaRotation(sim *core.Simulation) {
	target := ret.CurrentTarget

	sobExpiration := ret.SealOfBloodAura.ExpiresAt()
	nextSwingAt := ret.AutoAttacks.NextAttackAt()

	manaRegenAt := time.Duration(math.MaxInt64)
	// Roll seal of blood
	if sim.CurrentTime+time.Second >= sobExpiration {
		sobAndJudgementCost := ret.JudgementOfBlood.DefaultCast.Cost + ret.SealOfBlood.DefaultCast.Cost
		if ret.CanJudgementOfBlood(sim) && ret.CurrentMana() >= sobAndJudgementCost {
			ret.JudgementOfBlood.Cast(sim, target)
		}
		if ret.GCD.IsReady(sim) {
			if success := ret.SealOfBlood.Cast(sim, target); !success {
				// This should only happen in VERY BAD mana situations.
				manaRegenAt = sim.CurrentTime + ret.TimeUntilManaRegen(ret.SealOfBlood.CurCast.Cost)
			}
		}
	} else if ret.GCD.IsReady(sim) && ret.CrusaderStrike.CD.IsReady(sim) {
		spellGCD := ret.SpellGCD()
		sobAndCSCost := ret.CrusaderStrike.DefaultCast.Cost + ret.SealOfBlood.DefaultCast.Cost

		if !(spellGCD+sim.CurrentTime > nextSwingAt && sobExpiration < nextSwingAt) &&
			(ret.CurrentMana() >= sobAndCSCost) {
			// Crusader strike unless it will cause seal of blood to drop
			// Or we won't have enough mana to reseal
			ret.CrusaderStrike.Cast(sim, target)
		}
	}

	events := []time.Duration{
		ret.GCD.ReadyAt(),
		ret.CrusaderStrike.CD.ReadyAt(),
		manaRegenAt,
		sobExpiration - time.Second,
	}

	ret.waitUntilNextEvent(sim, events)
}

// Helper function for finding the next event
func (ret *RetributionPaladin) waitUntilNextEvent(sim *core.Simulation, events []time.Duration) {
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
		OnAction:     ret.mainRotation,
		NextActionAt: nextEventAt,
	}

	sim.AddPendingAction(pa)
}
