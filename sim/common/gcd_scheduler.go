package common

// Helper module for planning GCD-bound actions in advance.

import (
	"fmt"
	"time"

	"github.com/wowsims/tbc/sim/core"
)

const Unresolved = time.Duration(-1)

// Returns whether the cast was successful.
type AbilityCaster func(sim *core.Simulation) bool

type ScheduledAbility struct {
	// When to cast this ability. Might not cast at this time if there are conflicts.
	DesiredCastAt time.Duration

	// Limits the search window for conflict resolution.
	MinCastAt time.Duration
	MaxCastAt time.Duration

	// Override the default conflict resolution behavior of searching after the
	// desired cast time first. Instead, search before the desired cast time first.
	PrioritizeEarlierForConflicts bool

	// How much GCD time will be used by this ability.
	Duration time.Duration

	// How to cast this ability.
	TryCast AbilityCaster

	// When the ability will be cast.
	castAt time.Duration

	// When the ability cast will be completed.
	doneAt time.Duration
}

type GCDScheduler struct {
	// Scheduled abilities, sorted from soonest to latest CastAt time.
	schedule []ScheduledAbility

	nextAbilityIndex int

	managedMCDIDs []core.ActionID
	managedMCDs   []*core.MajorCooldown
}

// Returns the actual time at which the ability will be cast.
func (gs *GCDScheduler) Schedule(newAbility ScheduledAbility) time.Duration {
	newAbility.castAt = newAbility.DesiredCastAt
	newAbility.doneAt = newAbility.DesiredCastAt + newAbility.Duration

	oldLen := len(gs.schedule)
	if oldLen == 0 {
		gs.schedule = append(gs.schedule, newAbility)
		return newAbility.castAt
	}

	// Find the index at which this ability should be inserted, ignoring priority for now.
	var desiredIndex = 0
	for _, scheduledAbility := range gs.schedule {
		if newAbility.castAt < scheduledAbility.castAt {
			break
		}
		desiredIndex++
	}

	// If the insert was at the end with no overlap, can just append.
	if desiredIndex == oldLen && gs.schedule[oldLen-1].doneAt <= newAbility.castAt {
		gs.schedule = append(gs.schedule, newAbility)
		return newAbility.castAt
	}

	conflictBefore := desiredIndex > 0 && gs.schedule[desiredIndex-1].doneAt > newAbility.castAt
	conflictAfter := desiredIndex < oldLen && gs.schedule[desiredIndex].castAt < newAbility.doneAt
	if !conflictBefore && !conflictAfter {
		gs.schedule = append(gs.schedule, newAbility)
		copy(gs.schedule[desiredIndex+1:], gs.schedule[desiredIndex:])
		gs.schedule[desiredIndex] = newAbility
		return newAbility.castAt
	}

	// If we're here, we have a conflict.
	var castAt time.Duration
	if newAbility.PrioritizeEarlierForConflicts {
		castAt = gs.scheduleBefore(newAbility, desiredIndex, conflictAfter)
		if castAt == Unresolved {
			castAt = gs.scheduleAfter(newAbility, desiredIndex, conflictBefore)
		}
	} else {
		castAt = gs.scheduleAfter(newAbility, desiredIndex, conflictBefore)
		if castAt == Unresolved {
			castAt = gs.scheduleBefore(newAbility, desiredIndex, conflictAfter)
		}
	}
	return castAt
}

func (gs *GCDScheduler) scheduleBefore(newAbility ScheduledAbility, desiredIndex int, conflictAfter bool) time.Duration {
	curCastAt := newAbility.castAt
	if conflictAfter {
		curCastAt = gs.schedule[desiredIndex].castAt - newAbility.Duration
	}

	curIndex := desiredIndex
	for curIndex >= 0 && curCastAt >= newAbility.MinCastAt {
		conflictBefore := curIndex > 0 && gs.schedule[curIndex-1].doneAt > curCastAt
		if conflictBefore {
			curCastAt = gs.schedule[curIndex-1].castAt - newAbility.Duration
			curIndex--
		} else {
			newAbility.castAt = curCastAt
			newAbility.doneAt = curCastAt + newAbility.Duration

			gs.schedule = append(gs.schedule, newAbility)
			copy(gs.schedule[curIndex+1:], gs.schedule[curIndex:])
			gs.schedule[curIndex] = newAbility

			return newAbility.castAt
			break
		}
	}

	return Unresolved
}

func (gs *GCDScheduler) scheduleAfter(newAbility ScheduledAbility, desiredIndex int, conflictBefore bool) time.Duration {
	curCastAt := newAbility.castAt
	if conflictBefore {
		curCastAt = gs.schedule[desiredIndex-1].doneAt
	}

	curIndex := desiredIndex
	oldLen := len(gs.schedule)
	for curIndex <= oldLen && curCastAt <= newAbility.MaxCastAt {
		conflictAfter := curIndex < oldLen && gs.schedule[curIndex].castAt < curCastAt+newAbility.Duration
		if conflictAfter {
			curCastAt = gs.schedule[curIndex].doneAt
			curIndex++
		} else {
			newAbility.castAt = curCastAt
			newAbility.doneAt = curCastAt + newAbility.Duration

			gs.schedule = append(gs.schedule, newAbility)
			copy(gs.schedule[curIndex+1:], gs.schedule[curIndex:])
			gs.schedule[curIndex] = newAbility

			return newAbility.castAt
			break
		}
	}

	return Unresolved
}

// Schedules a group of abilities that must be cast back-to-back.
// Most settings are taken from the first ability.
func (gs *GCDScheduler) ScheduleGroup(newAbilities []ScheduledAbility) time.Duration {
	if len(newAbilities) == 0 {
		panic("Empty ability group!")
	}

	totalDuration := time.Duration(0)
	for _, newAbility := range newAbilities {
		totalDuration += newAbility.Duration
	}

	// Schedule a 'group ability' which is just a fake ability for reserving the time slots.
	groupAbility := ScheduledAbility{
		DesiredCastAt:                 newAbilities[0].DesiredCastAt,
		MinCastAt:                     newAbilities[0].MinCastAt,
		MaxCastAt:                     newAbilities[0].MaxCastAt,
		PrioritizeEarlierForConflicts: newAbilities[0].PrioritizeEarlierForConflicts,
		Duration:                      totalDuration,
	}
	groupCastAt := gs.Schedule(groupAbility)

	if groupCastAt == Unresolved {
		return Unresolved
	}

	// Update internals for the individual abilities, now that we know when they'll be cast.
	nextCastAt := groupCastAt
	for i, _ := range newAbilities {
		newAbilities[i].castAt = nextCastAt
		newAbilities[i].doneAt = nextCastAt + newAbilities[i].Duration
		nextCastAt = newAbilities[i].doneAt
	}

	// Replace the group ability with the individual abilities.
	for i, ability := range gs.schedule {
		if ability.castAt == groupCastAt {
			temp := make([]ScheduledAbility, len(gs.schedule)+len(newAbilities)-1)
			for j := 0; j < i; j++ {
				temp[j] = gs.schedule[j]
			}
			for j := 0; j < len(newAbilities); j++ {
				temp[i+j] = newAbilities[j]
			}
			for j := i + 1; j < len(gs.schedule); j++ {
				temp[j+len(newAbilities)-1] = gs.schedule[j]
			}
			gs.schedule = temp
			break
		}
	}

	return groupCastAt
}

// Takes ownership of a MCD, adding it to the schedule and removing it from the
// character's managed cooldowns.
func (gs *GCDScheduler) ScheduleMCD(character *core.Character, mcdID core.ActionID) {
	mcd := character.GetInitialMajorCooldown(mcdID)
	if mcd.Spell == nil {
		panic("No spell for MCD: " + mcdID.String())
	}
	if !mcd.Spell.ActionID.SameAction(mcdID) {
		return
	}

	mcdIdx := len(gs.managedMCDIDs)
	gs.managedMCDIDs = append(gs.managedMCDIDs, mcdID)
	gs.managedMCDs = append(gs.managedMCDs, nil)

	mcdAction := ScheduledAbility{
		Duration: core.GCDDefault,
		TryCast: func(sim *core.Simulation) bool {
			success := gs.managedMCDs[mcdIdx].TryActivate(sim, character)
			if success {
				character.UpdateMajorCooldowns()
			} else {
				character.EnableMajorCooldown(gs.managedMCDIDs[mcdIdx])
				gs.managedMCDs[mcdIdx].Spell.DefaultCast.GCD = 0
			}
			return success
		},
	}
	timings := mcd.GetTimings()
	curTime := time.Duration(0)
	maxDuration := character.Env.GetMaxDuration()
	i := 0
	if len(timings) > 0 {
		curTime = core.MaxDuration(curTime, timings[0])
	}
	for curTime <= maxDuration {
		ability := mcdAction
		ability.DesiredCastAt = curTime
		ability.MinCastAt = curTime
		ability.MaxCastAt = curTime + time.Second*30

		actualCastAt := gs.Schedule(ability)

		curTime = actualCastAt + mcd.Spell.CD.Duration
		i++
		if len(timings) > i {
			curTime = core.MaxDuration(curTime, timings[i])
		}
	}
}

func (gs *GCDScheduler) Reset(sim *core.Simulation, character *core.Character) {
	gs.nextAbilityIndex = 0

	for i, mcdID := range gs.managedMCDIDs {
		gs.managedMCDs[i] = character.GetMajorCooldown(mcdID)
		character.DisableMajorCooldown(mcdID)
	}
}

// Returns whether the cast was a success.
func (gs *GCDScheduler) DoNextAbility(sim *core.Simulation, character *core.Character) bool {
	if gs.nextAbilityIndex >= len(gs.schedule) {
		// It's possible for this function to get called near the end of the iteration,
		// after the final scheduled ability.
		return true
	}

	expectedCastAt := gs.schedule[gs.nextAbilityIndex].castAt
	if expectedCastAt > sim.CurrentTime {
		character.SetGCDTimer(sim, expectedCastAt)
		return true
	} else if expectedCastAt < sim.CurrentTime {
		panic(fmt.Sprintf("Missed scheduled cast! Expected %s but is now %s", expectedCastAt, sim.CurrentTime))
	}

	success := gs.schedule[gs.nextAbilityIndex].TryCast(sim)
	gs.nextAbilityIndex++

	if gs.nextAbilityIndex < len(gs.schedule) {
		nextCastAt := gs.schedule[gs.nextAbilityIndex].castAt
		if nextCastAt > character.NextGCDAt() {
			character.SetGCDTimer(sim, nextCastAt)
		}
	}

	return success
}
