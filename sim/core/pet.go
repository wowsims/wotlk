package core

import (
	"fmt"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Extension of Agent interface, for Pets.
type PetAgent interface {
	Agent

	// The Pet controlled by this PetAgent.
	GetPet() *Pet
	OwnerAttackSpeedChanged(sim *Simulation)
}

type OnPetEnable func(sim *Simulation)
type OnPetDisable func(sim *Simulation)
type PetStatInheritance func(ownerStats stats.Stats) stats.Stats

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
	statInheritance PetStatInheritance

	// No-op until finalized to prevent owner stats from affecting pet until we're ready.
	currentStatInheritance PetStatInheritance
	inheritedStats         stats.Stats

	// Whether this pet is currently active. Pets which are active throughout a whole
	// encounter, like Hunter pets, are always enabled. Pets which are instead summoned,
	// such as Mage Water Elemental, begin as disabled and are enabled when summoned.
	enabled bool
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
	pet.currentStatInheritance = func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{}
	}

	pet.AddStats(baseStats)
	pet.addUniversalStatDependencies()
	pet.PseudoStats.InFrontOfTarget = owner.PseudoStats.InFrontOfTarget

	return pet
}

// Add a default base if pets dont need this
func (pet *Pet) OwnerAttackSpeedChanged(sim *Simulation) {}

// Updates the stats for this pet in response to a stat change on the owner.
// addedStats is the amount of stats added to the owner (will be negative if the
// owner lost stats).
func (pet *Pet) addOwnerStats(sim *Simulation, addedStats stats.Stats) {
	// Temporary pets dont update stats after summon
	if pet.isGuardian {
		return
	}
	inheritedChange := pet.currentStatInheritance(addedStats)
	pet.inheritedStats = pet.inheritedStats.Add(inheritedChange)
	pet.AddStatsDynamic(sim, inheritedChange)
}

func (pet *Pet) Finalize() {
	pet.Character.Finalize(nil)
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
func (pet *Pet) advance(sim *Simulation, elapsedTime time.Duration) {
	pet.Character.advance(sim, elapsedTime)
}
func (pet *Pet) doneIteration(sim *Simulation) {
	pet.Character.doneIteration(sim)
	pet.isReset = false
}

func (pet *Pet) IsEnabled() bool {
	return pet.enabled
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
	pet.currentStatInheritance = pet.statInheritance

	pet.SetGCDTimer(sim, sim.CurrentTime)
	pet.AutoAttacks.EnableAutoSwing(sim)

	pet.enabled = true

	if pet.OnPetEnable != nil {
		pet.OnPetEnable(sim)
	}

	if sim.Log != nil {
		pet.Log(sim, "Pet stats: %s", pet.GetStats())
		pet.Log(sim, "Pet inherited stats: %s", pet.ApplyStatDependencies(pet.inheritedStats))
		pet.Log(sim, "Pet summoned")
	}
}
func (pet *Pet) Disable(sim *Simulation) {
	if !pet.enabled {
		if sim.Log != nil {
			pet.Log(sim, "No pet summoned")
		}
		return
	}

	// Remove inherited stats on dismiss if not permanent
	if pet.isGuardian {
		pet.AddStatsDynamic(sim, pet.inheritedStats.Multiply(-1))
		pet.inheritedStats = stats.Stats{}
		pet.currentStatInheritance = func(ownerStats stats.Stats) stats.Stats {
			return stats.Stats{}
		}
	}

	pet.CancelGCDTimer(sim)
	pet.AutoAttacks.CancelAutoSwing(sim)
	pet.enabled = false
	pet.DoNothing() // mark it is as doing nothing now.

	// If a pet is immediately re-summoned it might try to use GCD, so we need to
	// clear it.
	pet.Hardcast = Hardcast{}

	// TODO When MaxMana includes bonus mana will move this to the enable,
	//so bonus mana from the owner inheritence applies to the pet currentMana
	pet.manaBar.reset()

	if pet.timeoutAction != nil {
		pet.timeoutAction.Cancel(sim)
		pet.timeoutAction = nil
	}

	if pet.OnPetDisable != nil {
		pet.OnPetDisable(sim)
	}

	if sim.Log != nil {
		pet.Log(sim, "Pet dismissed")

		if sim.Log != nil {
			pet.Log(sim, pet.GetStats().String())
		}
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

// Default implementations for some Agent functions which most Pets don't need.
func (pet *Pet) GetCharacter() *Character {
	return &pet.Character
}
func (pet *Pet) AddRaidBuffs(raidBuffs *proto.RaidBuffs)    {}
func (pet *Pet) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {}
func (pet *Pet) ApplyTalents()                              {}
