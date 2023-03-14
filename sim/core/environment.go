package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"golang.org/x/exp/slices"
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

// Callback for doing something on prepull.
type PrepullAction struct {
	DoAt   time.Duration
	Action func(*Simulation)
}

type Environment struct {
	State EnvironmentState

	// Whether stats are currently being measured. Used to disable some validation
	// checks which are otherwise helpful.
	MeasuringStats bool

	Raid      *Raid
	Encounter Encounter
	AllUnits  []*Unit

	BaseDuration      time.Duration // base duration
	DurationVariation time.Duration // variation per duration

	// Effects to invoke when the Env is finalized.
	postFinalizeEffects []PostFinalizeEffect

	prepullActions []PrepullAction
}

func NewEnvironment(raidProto *proto.Raid, encounterProto *proto.Encounter) (*Environment, *proto.RaidStats) {
	env := &Environment{
		State: Created,
	}

	env.construct(raidProto, encounterProto)
	raidStats := env.initialize(raidProto, encounterProto)
	env.finalize(raidProto, encounterProto, raidStats)

	return env, raidStats
}

// The construction phase.
func (env *Environment) construct(raidProto *proto.Raid, encounterProto *proto.Encounter) {
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
		unit.CurrentTarget = env.Encounter.TargetUnits[0]
	}

	// Apply extra debuffs from raid.
	if raidProto.Debuffs != nil && len(env.Encounter.TargetUnits) > 0 {
		for targetIdx, targetUnit := range env.Encounter.TargetUnits {
			applyDebuffEffects(targetUnit, targetIdx, raidProto.Debuffs, raidProto)
		}
	}

	// Assign target or target using Tanks field.
	for _, target := range env.Encounter.Targets {
		if target.Index < int32(len(encounterProto.Targets)) {
			targetProto := encounterProto.Targets[target.Index]
			if targetProto.TankIndex >= 0 && targetProto.TankIndex < int32(len(raidProto.Tanks)) {
				raidTargetProto := raidProto.Tanks[targetProto.TankIndex]
				if raidTargetProto != nil {
					raidTarget := env.Raid.GetPlayerFromRaidTarget(raidTargetProto)
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
func (env *Environment) initialize(raidProto *proto.Raid, encounterProto *proto.Encounter) *proto.RaidStats {
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
func (env *Environment) finalize(raidProto *proto.Raid, _ *proto.Encounter, raidStats *proto.RaidStats) {
	for _, target := range env.Encounter.Targets {
		target.finalize()
	}

	for _, party := range env.Raid.Parties {
		for _, player := range party.Players {
			character := player.GetCharacter()
			character.Finalize()
			for _, petAgent := range character.Pets {
				petAgent.GetPet().Finalize()
			}
		}
	}

	for partyIdx, party := range env.Raid.Parties {
		partyProto := raidProto.Parties[partyIdx]
		for playerIdx, player := range party.Players {
			if playerIdx >= len(partyProto.Players) {
				// This happens for target dummies.
				continue
			}
			playerProto := partyProto.Players[playerIdx]
			char := player.GetCharacter()
			char.Rotation = char.newAPLRotation(playerProto.Rotation)
		}
	}

	for _, finalizeEffect := range env.postFinalizeEffects {
		finalizeEffect()
	}
	env.postFinalizeEffects = nil

	slices.SortStableFunc(env.prepullActions, func(a1, a2 PrepullAction) bool {
		return a1.DoAt < a2.DoAt
	})

	env.setupAttackTables()

	for partyIdx, party := range env.Raid.Parties {
		for _, player := range party.Players {
			character := player.GetCharacter()
			character.FillPlayerStats(raidStats.Parties[partyIdx].Players[character.PartyIndex])
		}
	}

	env.State = Finalized
}

func (env *Environment) setupAttackTables() {
	raidUnits := env.Raid.AllUnits
	if len(raidUnits) == 0 {
		return
	}

	for _, attacker := range env.AllUnits {
		attacker.AttackTables = make([]*AttackTable, len(env.AllUnits))
		for idx, defender := range env.AllUnits {
			attacker.AttackTables[idx] = NewAttackTable(attacker, defender)
		}
	}
}

func (env *Environment) IsFinalized() bool {
	return env.State >= Finalized
}

func (env *Environment) reset(sim *Simulation) {
	// Reset primary targets damage taken for tracking health fights.
	env.Encounter.DamageTaken = 0

	// Targets need to be reset before the raid, so that players can check for
	// the presence of permanent target auras in their Reset handlers.
	for _, target := range env.Encounter.Targets {
		target.Reset(sim)
	}

	env.Raid.reset(sim)
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

// Registers a callback to this Unit which will be invoked on the prepull at the specified
// negative time.
func (unit *Unit) RegisterPrepullAction(doAt time.Duration, action func(*Simulation)) {
	env := unit.Env
	if env.IsFinalized() {
		panic("Prepull actions may not be added once finalized!")
	}
	if doAt > 0 {
		panic("Prepull DoAt must not be positive!")
	}

	env.prepullActions = append(env.prepullActions, PrepullAction{
		DoAt:   doAt,
		Action: action,
	})
}
