package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type EnvironmentState int

const (
	Created EnvironmentState = iota
	Constructed
	Initialized
	Finalized
)

// Callback for doing something after finalization.
type PostFinalizeEffect func()

type Environment struct {
	State EnvironmentState

	Raid      *Raid
	Encounter Encounter
	AllUnits  []*Unit

	BaseDuration      time.Duration // base duration
	DurationVariation time.Duration // variation per duration

	// Effects to invoke when the Env is finalized.
	postFinalizeEffects []PostFinalizeEffect
}

func NewEnvironment(raidProto proto.Raid, encounterProto proto.Encounter) (*Environment, *proto.RaidStats) {
	env := &Environment{
		State: Created,
	}

	env.construct(raidProto, encounterProto)
	raidStats := env.initialize(raidProto, encounterProto)
	env.finalize(raidProto, encounterProto, raidStats)

	return env, raidStats
}

// The construction phase.
func (env *Environment) construct(raidProto proto.Raid, encounterProto proto.Encounter) {
	env.Encounter = NewEncounter(encounterProto)
	env.BaseDuration = env.Encounter.Duration
	env.DurationVariation = env.Encounter.DurationVariation
	env.Raid = NewRaid(raidProto)

	env.Raid.updatePlayersAndPets()

	env.AllUnits = append(env.Encounter.TargetUnits, env.Raid.AllUnits...)

	for unitIndex, unit := range env.AllUnits {
		unit.Env = env
		unit.UnitIndex = int32(unitIndex)
	}

	for _, unit := range env.Raid.AllUnits {
		unit.CurrentTarget = &env.Encounter.Targets[0].Unit
	}

	// Apply extra debuffs from raid.
	if raidProto.Debuffs != nil && len(env.Encounter.Targets) > 0 {
		applyDebuffEffects(&env.Encounter.Targets[0].Unit, *raidProto.Debuffs)
	}

	// Assign target or target using Tanks field.
	for _, target := range env.Encounter.Targets {
		if target.Index < int32(len(encounterProto.Targets)) {
			targetProto := encounterProto.Targets[target.Index]
			if targetProto.TankIndex >= 0 && targetProto.TankIndex < int32(len(raidProto.Tanks)) {
				raidTargetProto := raidProto.Tanks[targetProto.TankIndex]
				if raidTargetProto != nil {
					raidTarget := env.Raid.GetPlayerFromRaidTarget(*raidTargetProto)
					if raidTarget != nil {
						target.CurrentTarget = &raidTarget.GetCharacter().Unit
					}
				}
			}
		}
	}

	env.State = Constructed
}

// The initialization phase.
func (env *Environment) initialize(raidProto proto.Raid, encounterProto proto.Encounter) *proto.RaidStats {
	for _, target := range env.Encounter.Targets {
		if target.Index < int32(len(encounterProto.Targets)) {
			target.initialize(encounterProto.Targets[target.Index])
		} else {
			target.initialize(nil)
		}
	}

	for _, party := range env.Raid.Parties {
		for _, playerOrPet := range party.PlayersAndPets {
			playerOrPet.GetCharacter().initialize(playerOrPet)
		}
	}

	raidStats := env.Raid.applyCharacterEffects(raidProto)

	for _, party := range env.Raid.Parties {
		for _, playerOrPet := range party.PlayersAndPets {
			playerOrPet.Initialize()
		}
	}

	env.State = Initialized
	return raidStats
}

// The finalization phase.
func (env *Environment) finalize(raidProto proto.Raid, encounterProto proto.Encounter, raidStats *proto.RaidStats) {
	for _, target := range env.Encounter.Targets {
		target.finalize()
	}

	for partyIdx, party := range env.Raid.Parties {
		for _, player := range party.Players {
			character := player.GetCharacter()
			character.Finalize(raidStats.Parties[partyIdx].Players[character.PartyIndex])

			for _, petAgent := range character.Pets {
				petAgent.GetPet().Finalize()
			}
		}
	}

	for _, finalizeEffect := range env.postFinalizeEffects {
		finalizeEffect()
	}
	env.postFinalizeEffects = nil

	env.setupAttackTables()

	env.State = Finalized
}

func (env *Environment) setupAttackTables() {
	raidUnits := env.Raid.AllUnits
	if len(raidUnits) == 0 {
		return
	}

	for _, unit := range env.AllUnits {
		unit.AttackTables = make([]*AttackTable, len(env.AllUnits))
		unit.DefenseTables = make([]*AttackTable, len(env.AllUnits))
	}

	for i := 0; i < len(env.AllUnits); i++ {
		for j := 0; j < len(env.AllUnits); j++ {
			attacker := env.AllUnits[i]
			defender := env.AllUnits[j]
			attackTable := NewAttackTable(attacker, defender)
			attacker.AttackTables[j] = attackTable
			defender.DefenseTables[i] = attackTable
		}
	}
}

func (env *Environment) IsFinalized() bool {
	return env.State >= Finalized
}

// The maximum possible duration for any iteration.
func (env *Environment) GetMaxDuration() time.Duration {
	return env.BaseDuration + env.DurationVariation
}

func (env *Environment) GetNumTargets() int32 {
	return int32(len(env.Encounter.Targets))
}

func (env *Environment) GetTarget(index int32) *Target {
	return env.Encounter.Targets[index]
}
func (env *Environment) GetTargetUnit(index int32) *Unit {
	return &env.Encounter.Targets[index].Unit
}
func (env *Environment) NextTarget(target *Unit) *Target {
	return env.Encounter.Targets[target.Index].NextTarget()
}
func (env *Environment) NextTargetUnit(target *Unit) *Unit {
	return &env.NextTarget(target).Unit
}

// Registers a callback to this Character which will be invoked after all Units
// are finalized.
func (env *Environment) RegisterPostFinalizeEffect(postFinalizeEffect PostFinalizeEffect) {
	if env.IsFinalized() {
		panic("Finalize effects may not be added once finalized!")
	}

	env.postFinalizeEffects = append(env.postFinalizeEffects, postFinalizeEffect)
}
