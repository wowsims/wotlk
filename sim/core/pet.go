package core

import (
	"fmt"
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Extension of Agent interface, for Pets.
type PetAgent interface {
	Agent

	// The Pet controlled by this PetAgent.
	GetPet() *Pet
}

type PetStatInheritance func(ownerStats stats.Stats) stats.Stats

// Pet is an extension of Character, for any entity created by a player that can
// take actions on its own.
type Pet struct {
	Character

	Owner *Character

	// Calculates inherited stats based on owner stats or stat changes.
	statInheritance PetStatInheritance

	// No-op until finalized to prevent owner stats from affecting pet until we're ready.
	currentStatInheritance PetStatInheritance

	initialEnabled bool

	// Whether this pet is currently active. Pets which are active throughout a whole
	// encounter, like Hunter pets, are always enabled. Pets which are instead summoned,
	// such as Mage Water Elemental, begin as disabled and are enabled when summoned.
	enabled bool

	// Some pets expire after a certain duration. This is the pending action that disables
	// the pet on expiration.
	timeoutAction *PendingAction
}

func NewPet(name string, owner *Character, baseStats stats.Stats, statInheritance PetStatInheritance, enabledOnStart bool) Pet {
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
			},
			Name:       name,
			Party:      owner.Party,
			PartyIndex: owner.PartyIndex,
			baseStats:  baseStats,
		},
		Owner:           owner,
		statInheritance: statInheritance,
		initialEnabled:  enabledOnStart,
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

// Updates the stats for this pet in response to a stat change on the owner.
// addedStats is the amount of stats added to the owner (will be negative if the
// owner lost stats).
func (pet *Pet) addOwnerStats(sim *Simulation, addedStats stats.Stats) {
	inheritedChange := pet.currentStatInheritance(addedStats)
	pet.AddStatsDynamic(sim, inheritedChange)
}
func (pet *Pet) addOwnerStat(sim *Simulation, stat stats.Stat, addedAmount float64) {
	s := stats.Stats{}
	s[stat] = addedAmount
	pet.addOwnerStats(sim, s)
}

// This needs to be called after owner stats are finalized so we can inherit the
// final values.
func (pet *Pet) Finalize() {
	inheritedStats := pet.statInheritance(pet.Owner.GetStats())
	pet.AddStats(inheritedStats)
	pet.currentStatInheritance = pet.statInheritance
	pet.Character.Finalize(nil)
}

func (pet *Pet) reset(sim *Simulation, agent PetAgent) {
	pet.Character.reset(sim, agent)
	pet.enabled = false
	if pet.initialEnabled {
		pet.Enable(sim, agent)
	}
}
func (pet *Pet) advance(sim *Simulation, elapsedTime time.Duration) {
	pet.Character.advance(sim, elapsedTime)
}
func (pet *Pet) doneIteration(sim *Simulation) {
	pet.Character.doneIteration(sim)
}

func (pet *Pet) IsEnabled() bool {
	return pet.enabled
}

// petAgent should be the PetAgent which embeds this Pet.
func (pet *Pet) Enable(sim *Simulation, petAgent PetAgent) {
	if pet.enabled {
		return
	}

	pet.SetGCDTimer(sim, sim.CurrentTime)

	pet.enabled = true

	if sim.Log != nil {
		pet.Log(sim, "Pet summoned")
	}
}
func (pet *Pet) Disable(sim *Simulation) {
	if !pet.enabled {
		return
	}

	pet.CancelGCDTimer(sim)
	pet.AutoAttacks.CancelAutoSwing(sim)
	pet.enabled = false

	// If a pet is immediately re-summoned it might try to use GCD, so we need to
	// clear it.
	pet.Hardcast = Hardcast{}

	// Reset pet mana.
	pet.stats[stats.Mana] = pet.MaxMana()

	if pet.timeoutAction != nil {
		pet.timeoutAction.Cancel(sim)
		pet.timeoutAction = nil
	}

	if sim.Log != nil {
		pet.Log(sim, "Pet dismissed")
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
