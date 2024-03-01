package core

import (
	"fmt"
	"slices"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Extension of Agent interface, for Pets.
type PetAgent interface {
	Agent

	// The Pet controlled by this PetAgent.
	GetPet() *Pet
}

type OnPetEnable func(sim *Simulation)
type OnPetDisable func(sim *Simulation)

type PetStatInheritance func(ownerStats stats.Stats) stats.Stats
type PetMeleeSpeedInheritance func(amount float64)

// Pet is an extension of Character, for any entity created by a player that can
// take actions on its own.
type Pet struct {
	Character

	Owner *Character

	isGuardian     bool
	enabledOnStart bool

	OnPetEnable  OnPetEnable
	OnPetDisable OnPetDisable

	// Calculates inherited stats based on owner stats or stat changes.
	statInheritance        PetStatInheritance
	dynamicStatInheritance PetStatInheritance
	inheritedStats         stats.Stats

	// DK pets also inherit their owner's MeleeSpeed. This replace OwnerAttackSpeedChanged.
	dynamicMeleeSpeedInheritance PetMeleeSpeedInheritance

	isReset bool

	// Some pets expire after a certain duration. This is the pending action that disables
	// the pet on expiration.
	timeoutAction *PendingAction
}

func NewPet(name string, owner *Character, baseStats stats.Stats, statInheritance PetStatInheritance, enabledOnStart bool, isGuardian bool) Pet {
	pet := Pet{
		Character: Character{
			Unit: Unit{
				Type:        PetUnit,
				Index:       owner.Party.Raid.getNextPetIndex(),
				Label:       fmt.Sprintf("%s - %s", owner.Label, name),
				Level:       CharacterLevel,
				PseudoStats: stats.NewPseudoStats(),
				auraTracker: newAuraTracker(),
				Metrics:     NewUnitMetrics(),

				StatDependencyManager: stats.NewStatDependencyManager(),
			},
			Name:       name,
			Party:      owner.Party,
			PartyIndex: owner.PartyIndex,
			baseStats:  baseStats,
		},
		Owner:           owner,
		statInheritance: statInheritance,
		enabledOnStart:  enabledOnStart,
		isGuardian:      isGuardian,
	}
	pet.GCD = pet.NewTimer()

	pet.AddStats(baseStats)
	pet.addUniversalStatDependencies()
	pet.PseudoStats.InFrontOfTarget = owner.PseudoStats.InFrontOfTarget

	return pet
}

// Updates the stats for this pet in response to a stat change on the owner.
// addedStats is the amount of stats added to the owner (will be negative if the
// owner lost stats).
func (pet *Pet) addOwnerStats(sim *Simulation, addedStats stats.Stats) {
	inheritedChange := pet.dynamicStatInheritance(addedStats)

	pet.inheritedStats.AddInplace(&inheritedChange)
	pet.AddStatsDynamic(sim, inheritedChange)
}

func (pet *Pet) reset(sim *Simulation, agent PetAgent) {
	if pet.isReset {
		return
	}
	pet.isReset = true

	pet.Character.reset(sim, agent)

	pet.CancelGCDTimer(sim)
	pet.AutoAttacks.CancelAutoSwing(sim)

	pet.enabled = false
	if pet.enabledOnStart {
		pet.Enable(sim, agent)
	}
}
func (pet *Pet) doneIteration(sim *Simulation) {
	pet.Character.doneIteration(sim)
	pet.isReset = false
}

func (pet *Pet) IsGuardian() bool {
	return pet.isGuardian
}

// petAgent should be the PetAgent which embeds this Pet.
func (pet *Pet) Enable(sim *Simulation, petAgent PetAgent) {
	if pet.enabled {
		if sim.Log != nil {
			pet.Log(sim, "Pet already summoned")
		}
		return
	}

	// In case of Pre-pull guardian summoning we need to reset
	// TODO: Check if this has side effects
	if !pet.isReset {
		pet.reset(sim, petAgent)
	}

	pet.inheritedStats = pet.statInheritance(pet.Owner.GetStats())
	pet.AddStatsDynamic(sim, pet.inheritedStats)

	if !pet.isGuardian {
		pet.Owner.DynamicStatsPets = append(pet.Owner.DynamicStatsPets, pet)
		pet.dynamicStatInheritance = pet.statInheritance
	}

	//reset current mana after applying stats
	pet.manaBar.reset()

	// Call onEnable callbacks before enabling auto swing
	// to not have to reorder PAs multiple times
	pet.enabled = true

	if pet.OnPetEnable != nil {
		pet.OnPetEnable(sim)
	}

	pet.SetGCDTimer(sim, max(0, sim.CurrentTime))
	if sim.CurrentTime >= 0 {
		pet.AutoAttacks.EnableAutoSwing(sim)
	} else {
		sim.AddPendingAction(&PendingAction{
			NextActionAt: 0,
			OnAction:     pet.AutoAttacks.EnableAutoSwing,
		})
	}

	if sim.Log != nil {
		pet.Log(sim, "Pet stats: %s", pet.GetStats().FlatString())
		pet.Log(sim, "Pet inherited stats: %s", pet.ApplyStatDependencies(pet.inheritedStats).FlatString())
		pet.Log(sim, "Pet summoned")
	}

	sim.addTracker(&pet.auraTracker)

	if pet.HasFocusBar() {
		pet.focusBar.enable(sim)
	}
}

// Helper for enabling a pet that will expire after a certain duration.
func (pet *Pet) EnableWithTimeout(sim *Simulation, petAgent PetAgent, petDuration time.Duration) {
	pet.Enable(sim, petAgent)

	pet.timeoutAction = &PendingAction{
		NextActionAt: sim.CurrentTime + petDuration,
		OnAction: func(sim *Simulation) {
			pet.Disable(sim)
		},
	}
	sim.AddPendingAction(pet.timeoutAction)
}

// Enables and possibly updates how the pet inherits its owner's stats. DK use only.
func (pet *Pet) EnableDynamicStats(inheritance PetStatInheritance) {
	if !slices.Contains(pet.Owner.DynamicStatsPets, pet) {
		pet.Owner.DynamicStatsPets = append(pet.Owner.DynamicStatsPets, pet)
	}
	pet.dynamicStatInheritance = inheritance
}

// Enables and possibly updates how the pet inherits its owner's melee speed. DK use only.
func (pet *Pet) EnableDynamicMeleeSpeed(inheritance PetMeleeSpeedInheritance) {
	if !slices.Contains(pet.Owner.DynamicMeleeSpeedPets, pet) {
		pet.Owner.DynamicMeleeSpeedPets = append(pet.Owner.DynamicMeleeSpeedPets, pet)
	}
	pet.dynamicMeleeSpeedInheritance = inheritance
}

func (pet *Pet) Disable(sim *Simulation) {
	if !pet.enabled {
		if sim.Log != nil {
			pet.Log(sim, "No pet summoned")
		}
		return
	}

	// Remove inherited stats on dismiss if not permanent
	if pet.isGuardian || pet.timeoutAction != nil {
		pet.AddStatsDynamic(sim, pet.inheritedStats.Invert())
		pet.inheritedStats = stats.Stats{}
	}

	if pet.dynamicStatInheritance != nil {
		if idx := slices.Index(pet.Owner.DynamicStatsPets, pet); idx != -1 {
			pet.Owner.DynamicStatsPets = removeBySwappingToBack(pet.Owner.DynamicStatsPets, idx)
		}
		pet.dynamicStatInheritance = nil
	}

	if pet.dynamicMeleeSpeedInheritance != nil {
		if idx := slices.Index(pet.Owner.DynamicMeleeSpeedPets, pet); idx != -1 {
			pet.Owner.DynamicMeleeSpeedPets = removeBySwappingToBack(pet.Owner.DynamicMeleeSpeedPets, idx)
		}
		pet.dynamicMeleeSpeedInheritance = nil
	}

	pet.CancelGCDTimer(sim)
	pet.focusBar.disable(sim)
	pet.AutoAttacks.CancelAutoSwing(sim)
	pet.enabled = false

	// If a pet is immediately re-summoned it might try to use GCD, so we need to clear it.
	pet.Hardcast = Hardcast{}

	if pet.timeoutAction != nil {
		pet.timeoutAction.Cancel(sim)
		pet.timeoutAction = nil
	}

	if pet.OnPetDisable != nil {
		pet.OnPetDisable(sim)
	}

	pet.auraTracker.expireAll(sim)

	sim.removeTracker(&pet.auraTracker)

	if sim.Log != nil {
		pet.Log(sim, "Pet dismissed")
		pet.Log(sim, pet.GetStats().FlatString())
	}
}

// Default implementations for some Agent functions which most Pets don't need.
func (pet *Pet) GetCharacter() *Character {
	return &pet.Character
}
func (pet *Pet) AddRaidBuffs(_ *proto.RaidBuffs)   {}
func (pet *Pet) AddPartyBuffs(_ *proto.PartyBuffs) {}
func (pet *Pet) ApplyTalents()                     {}
func (pet *Pet) OnGCDReady(_ *Simulation)          {}
