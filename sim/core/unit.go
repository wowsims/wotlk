package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type UnitType int
type SpellRegisteredHandler func(spell *Spell)

const (
	PlayerUnit UnitType = iota
	EnemyUnit
	PetUnit
)

type PowerBarType int

const (
	ManaBar PowerBarType = iota
	EnergyBar
	RageBar
	RunicPower
)

type DynamicDamageTakenModifier func(sim *Simulation, spell *Spell, result *SpellResult)

// Unit is an abstraction of a Character/Boss/Pet/etc, containing functionality
// shared by all of them.
type Unit struct {
	Type UnitType

	// Index of this unit with its group.
	//  For Players, this is the 0-indexed raid index (0-24).
	//  For Enemies, this is its enemy index.
	//  For Pets, this is the same as the owner's index.
	Index int32

	// Unique index of this unit among all units in the environment.
	// This is used as the index for attack tables.
	UnitIndex int32

	// Unique label for logging.
	Label string

	Level int32 // Level of Unit, e.g. Bosses are 83.

	MobType proto.MobType

	// Amount of time it takes for the human agent to react to in-game events.
	// Used by certain APL values and actions.
	ReactionTime time.Duration

	// Amount of time following a post-GCD channel tick, to when the next action can be performed.
	ChannelClipDelay time.Duration

	// How far this unit is from its target(s). Measured in yards, this is used
	// for calculating spell travel time for certain spells.
	DistanceFromTarget float64

	// How many casts on average a Valkyr will get off during its lifetime.
	NibelungAverageCasts float64

	// Environment in which this Unit exists. This will be nil until after the
	// construction phase.
	Env *Environment

	// Whether this unit is able to perform actions.
	enabled bool

	// Stats this Unit will have at the very start of each Sim iteration.
	// Includes all equipment / buffs / permanent effects but not temporary
	// effects from items / abilities.
	initialStats            stats.Stats
	initialStatsWithoutDeps stats.Stats

	initialPseudoStats stats.PseudoStats

	// Cast speed without any temporary effects.
	initialCastSpeed float64

	// Melee swing speed without any temporary effects.
	initialMeleeSwingSpeed float64

	// Ranged swing speed without any temporary effects.
	initialRangedSwingSpeed float64

	// Provides aura tracking behavior.
	auraTracker

	// Current stats, including temporary effects but not dependencies.
	statsWithoutDeps stats.Stats

	// Current stats, including temporary effects and dependencies.
	stats stats.Stats

	// Provides stat dependency management behavior.
	stats.StatDependencyManager

	PseudoStats stats.PseudoStats

	currentPowerBar PowerBarType
	healthBar
	manaBar
	rageBar
	energyBar
	focusBar
	runicPowerBar

	// All spells that can be cast by this unit.
	Spellbook                 []*Spell
	spellRegistrationHandlers []SpellRegisteredHandler

	// Pets owned by this Unit.
	PetAgents []PetAgent

	DynamicStatsPets      []*Pet
	DynamicMeleeSpeedPets []*Pet

	// AutoAttacks is the manager for auto attack swings.
	// Must be enabled to use, with "EnableAutoAttacks()".
	AutoAttacks AutoAttacks

	Rotation *APLRotation

	// Statistics describing the results of the sim.
	Metrics UnitMetrics

	cdTimers []*Timer

	AttackTables                []*AttackTable
	DynamicDamageTakenModifiers []DynamicDamageTakenModifier

	GCD *Timer

	// Used for applying the effect of a hardcast spell when casting finishes.
	//  For channeled spells, only Expires is set.
	// No more than one cast may be active at any given time.
	Hardcast Hardcast

	// GCD-related PendingActions.
	gcdAction      *PendingAction
	hardcastAction *PendingAction

	// Cached mana return values per tick.
	manaTickWhileCasting    float64
	manaTickWhileNotCasting float64

	CastSpeed float64

	CurrentTarget *Unit
	defaultTarget *Unit

	// The currently-channeled DOT spell, otherwise nil.
	ChanneledDot *Dot
}

// Units can be disabled for several reasons:
//  1. Downtime for temporary pets (e.g. Water Elemental)
//  2. Enemy units in various phases (not yet implemented)
//  3. Dead units (not yet implemented)
func (unit *Unit) IsEnabled() bool {
	return unit.enabled
}

func (unit *Unit) IsActive() bool {
	return unit.IsEnabled() && unit.CurrentHealthPercent() > 0
}

func (unit *Unit) IsOpponent(other *Unit) bool {
	return (unit.Type == EnemyUnit) != (other.Type == EnemyUnit)
}

func (unit *Unit) GetOpponents() []*Unit {
	if unit.Type == EnemyUnit {
		return unit.Env.Raid.AllUnits
	} else {
		return unit.Env.Encounter.TargetUnits
	}
}

func (unit *Unit) LogLabel() string {
	return "[" + unit.Label + "]"
}

func (unit *Unit) Log(sim *Simulation, message string, vals ...interface{}) {
	sim.Log(unit.LogLabel()+" "+message, vals...)
}

func (unit *Unit) GetInitialStat(stat stats.Stat) float64 {
	return unit.initialStats[stat]
}
func (unit *Unit) GetStats() stats.Stats {
	return unit.stats
}
func (unit *Unit) GetStat(stat stats.Stat) float64 {
	return unit.stats[stat]
}

func (unit *Unit) AddStats(stat stats.Stats) {
	if unit.Env != nil && unit.Env.IsFinalized() {
		panic("Already finalized, use AddStatsDynamic instead!")
	}
	unit.stats = unit.stats.Add(stat)
}
func (unit *Unit) AddStat(stat stats.Stat, amount float64) {
	if unit.Env != nil && unit.Env.IsFinalized() {
		panic("Already finalized, use AddStatDynamic instead!")
	}
	unit.stats[stat] += amount
}

func (unit *Unit) AddDynamicDamageTakenModifier(ddtm DynamicDamageTakenModifier) {
	if unit.Env != nil && unit.Env.IsFinalized() {
		panic("Already finalized, cannot add dynamic damage taken modifier!")
	}
	unit.DynamicDamageTakenModifiers = append(unit.DynamicDamageTakenModifiers, ddtm)
}

func (unit *Unit) AddStatsDynamic(sim *Simulation, bonus stats.Stats) {
	if unit.Env == nil {
		panic("Environment not constructed.")
	} else if !unit.Env.IsFinalized() && !unit.Env.MeasuringStats {
		panic("Not finalized, use AddStats instead!")
	}

	unit.statsWithoutDeps.AddInplace(&bonus)

	bonus = unit.ApplyStatDependencies(bonus)

	if sim.Log != nil {
		unit.Log(sim, "Dynamic stat change: %s", bonus.FlatString())
	}

	unit.stats.AddInplace(&bonus)
	unit.processDynamicBonus(sim, bonus)
}

func (unit *Unit) AddStatDynamic(sim *Simulation, stat stats.Stat, amount float64) {
	bonus := stats.Stats{}
	bonus[stat] = amount
	unit.AddStatsDynamic(sim, bonus)
}

func (unit *Unit) processDynamicBonus(sim *Simulation, bonus stats.Stats) {
	if bonus[stats.MP5] != 0 || bonus[stats.Intellect] != 0 || bonus[stats.Spirit] != 0 {
		unit.UpdateManaRegenRates()
	}
	if bonus[stats.MeleeHaste] != 0 {
		unit.AutoAttacks.UpdateSwingTimers(sim)
	}
	if bonus[stats.SpellHaste] != 0 {
		unit.updateCastSpeed()
	}

	for _, pet := range unit.DynamicStatsPets {
		pet.addOwnerStats(sim, bonus)
	}
}

func (unit *Unit) EnableDynamicStatDep(sim *Simulation, dep *stats.StatDependency) {
	if unit.StatDependencyManager.EnableDynamicStatDep(dep) {
		oldStats := unit.stats
		unit.stats = unit.ApplyStatDependencies(unit.statsWithoutDeps)
		unit.processDynamicBonus(sim, unit.stats.Subtract(oldStats))

		if sim.Log != nil {
			unit.Log(sim, "Dynamic dep enabled (%s): %s", dep.String(), unit.stats.Subtract(oldStats).FlatString())
		}
	}
}
func (unit *Unit) DisableDynamicStatDep(sim *Simulation, dep *stats.StatDependency) {
	if unit.StatDependencyManager.DisableDynamicStatDep(dep) {
		oldStats := unit.stats
		unit.stats = unit.ApplyStatDependencies(unit.statsWithoutDeps)
		unit.processDynamicBonus(sim, unit.stats.Subtract(oldStats))

		if sim.Log != nil {
			unit.Log(sim, "Dynamic dep disabled (%s): %s", dep.String(), unit.stats.Subtract(oldStats).FlatString())
		}
	}
}

// Returns whether the indicates stat is currently modified by a temporary bonus.
func (unit *Unit) HasTemporaryBonusForStat(stat stats.Stat) bool {
	return unit.initialStats[stat] != unit.stats[stat]
}

// Returns if spell casting has any temporary increases active.
func (unit *Unit) HasTemporarySpellCastSpeedIncrease() bool {
	return unit.CastSpeed != unit.initialCastSpeed
}

// Returns if melee swings have any temporary increases active.
func (unit *Unit) HasTemporaryMeleeSwingSpeedIncrease() bool {
	return unit.SwingSpeed() != unit.initialMeleeSwingSpeed
}

// Returns if ranged swings have any temporary increases active.
func (unit *Unit) HasTemporaryRangedSwingSpeedIncrease() bool {
	return unit.RangedSwingSpeed() != unit.initialRangedSwingSpeed
}

func (unit *Unit) InitialCastSpeed() float64 {
	return unit.initialCastSpeed
}

func (unit *Unit) SpellGCD() time.Duration {
	return max(GCDMin, unit.ApplyCastSpeed(GCDDefault))
}

func (unit *Unit) updateCastSpeed() {
	unit.CastSpeed = 1 / (unit.PseudoStats.CastSpeedMultiplier * (1 + (unit.stats[stats.SpellHaste] / (HasteRatingPerHastePercent * 100))))
}
func (unit *Unit) MultiplyCastSpeed(amount float64) {
	unit.PseudoStats.CastSpeedMultiplier *= amount
	unit.updateCastSpeed()
}

func (unit *Unit) ApplyCastSpeed(dur time.Duration) time.Duration {
	return time.Duration(float64(dur) * unit.CastSpeed)
}
func (unit *Unit) ApplyCastSpeedForSpell(dur time.Duration, spell *Spell) time.Duration {
	return time.Duration(float64(dur) * unit.CastSpeed * spell.CastTimeMultiplier)
}

func (unit *Unit) SwingSpeed() float64 {
	return unit.PseudoStats.MeleeSpeedMultiplier * (1 + (unit.stats[stats.MeleeHaste] / (unit.PseudoStats.MeleeHasteRatingPerHastePercent * 100)))
}

func (unit *Unit) Armor() float64 {
	return unit.PseudoStats.ArmorMultiplier * unit.stats[stats.Armor]
}

func (unit *Unit) BlockValue() float64 {
	return unit.PseudoStats.BlockValueMultiplier * unit.stats[stats.BlockValue]
}

func (unit *Unit) ArmorPenetrationPercentage(armorPenRating float64) float64 {
	return max(min(armorPenRating/ArmorPenPerPercentArmor, 100.0)*0.01, 0.0)
}

func (unit *Unit) RangedSwingSpeed() float64 {
	return unit.PseudoStats.RangedSpeedMultiplier * (1 + (unit.stats[stats.MeleeHaste] / (HasteRatingPerHastePercent * 100)))
}

// MultiplyMeleeSpeed will alter the attack speed multiplier and change swing speed of all autoattack swings in progress.
func (unit *Unit) MultiplyMeleeSpeed(sim *Simulation, amount float64) {
	unit.PseudoStats.MeleeSpeedMultiplier *= amount

	for _, pet := range unit.DynamicMeleeSpeedPets {
		pet.dynamicMeleeSpeedInheritance(amount)
	}
	unit.AutoAttacks.UpdateSwingTimers(sim)
}

func (unit *Unit) MultiplyRangedSpeed(sim *Simulation, amount float64) {
	unit.PseudoStats.RangedSpeedMultiplier *= amount
	unit.AutoAttacks.UpdateSwingTimers(sim)
}

// Helper for when both MultiplyMeleeSpeed and MultiplyRangedSpeed are needed.
func (unit *Unit) MultiplyAttackSpeed(sim *Simulation, amount float64) {
	unit.PseudoStats.MeleeSpeedMultiplier *= amount
	unit.PseudoStats.RangedSpeedMultiplier *= amount

	for _, pet := range unit.DynamicMeleeSpeedPets {
		pet.dynamicMeleeSpeedInheritance(amount)
	}
	unit.AutoAttacks.UpdateSwingTimers(sim)
}

func (unit *Unit) AddBonusRangedHitRating(amount float64) {
	unit.OnSpellRegistered(func(spell *Spell) {
		if spell.ProcMask.Matches(ProcMaskRanged) {
			spell.BonusHitRating += amount
		}
	})
}
func (unit *Unit) AddBonusRangedCritRating(amount float64) {
	unit.OnSpellRegistered(func(spell *Spell) {
		if spell.ProcMask.Matches(ProcMaskRanged) {
			spell.BonusCritRating += amount
		}
	})
}

func (unit *Unit) SetCurrentPowerBar(bar PowerBarType) {
	unit.currentPowerBar = bar
}

func (unit *Unit) GetCurrentPowerBar() PowerBarType {
	return unit.currentPowerBar
}

func (unit *Unit) finalize() {
	if unit.Env.IsFinalized() {
		panic("Unit already finalized!")
	}

	// Make sure we don't accidentally set initial stats instead of stats.
	if !unit.initialStats.Equals(stats.Stats{}) {
		panic("Initial stats may not be set before finalized: " + unit.initialStats.String())
	}

	unit.defaultTarget = unit.CurrentTarget
	unit.applyParryHaste()
	unit.updateCastSpeed()

	// All stats added up to this point are part of the 'initial' stats.
	unit.initialStatsWithoutDeps = unit.stats
	unit.initialPseudoStats = unit.PseudoStats
	unit.initialCastSpeed = unit.CastSpeed
	unit.initialMeleeSwingSpeed = unit.SwingSpeed()
	unit.initialRangedSwingSpeed = unit.RangedSwingSpeed()

	unit.StatDependencyManager.FinalizeStatDeps()
	unit.initialStats = unit.ApplyStatDependencies(unit.initialStatsWithoutDeps)
	unit.statsWithoutDeps = unit.initialStatsWithoutDeps
	unit.stats = unit.initialStats

	unit.AutoAttacks.finalize()

	for _, spell := range unit.Spellbook {
		spell.finalize()
	}
}

func (unit *Unit) reset(sim *Simulation, _ Agent) {
	unit.enabled = true
	unit.resetCDs(sim)
	unit.Hardcast.Expires = startingCDTime
	unit.ChanneledDot = nil
	unit.Metrics.reset()
	unit.ResetStatDeps()
	unit.statsWithoutDeps = unit.initialStatsWithoutDeps
	unit.stats = unit.initialStats
	unit.PseudoStats = unit.initialPseudoStats
	unit.auraTracker.reset(sim)
	// Spellbook needs to be reset AFTER auras.
	for _, spell := range unit.Spellbook {
		spell.reset(sim)
	}

	unit.manaBar.reset()
	unit.focusBar.reset(sim)
	unit.healthBar.reset(sim)
	unit.UpdateManaRegenRates()

	unit.energyBar.reset(sim)
	unit.rageBar.reset(sim)
	unit.runicPowerBar.reset(sim)

	unit.AutoAttacks.reset(sim)

	if unit.Rotation != nil {
		unit.Rotation.reset(sim)
	}

	unit.DynamicStatsPets = unit.DynamicStatsPets[:0]
	unit.DynamicMeleeSpeedPets = unit.DynamicMeleeSpeedPets[:0]

	if unit.Type != PetUnit {
		sim.addTracker(&unit.auraTracker)
	}
}

func (unit *Unit) startPull(sim *Simulation) {
	unit.AutoAttacks.startPull(sim)

	if unit.Type == PlayerUnit {
		unit.SetGCDTimer(sim, max(0, unit.GCD.ReadyAt()))
	}
}

func (unit *Unit) doneIteration(sim *Simulation) {
	unit.Hardcast = Hardcast{}

	unit.manaBar.doneIteration(sim)
	unit.rageBar.doneIteration()

	unit.auraTracker.doneIteration(sim)
	for _, spell := range unit.Spellbook {
		spell.doneIteration()
	}
}

func (unit *Unit) GetSpellsMatchingSchool(school SpellSchool) []*Spell {
	var spells []*Spell
	for _, spell := range unit.Spellbook {
		if spell.SpellSchool.Matches(school) {
			spells = append(spells, spell)
		}
	}
	return spells
}

func (unit *Unit) GetUnit(ref *proto.UnitReference) *Unit {
	return unit.Env.GetUnit(ref, unit)
}

func (unit *Unit) GetMetadata() *proto.UnitMetadata {
	metadata := &proto.UnitMetadata{
		Name: unit.Label,
	}

	metadata.Spells = MapSlice(unit.Spellbook, func(spell *Spell) *proto.SpellStats {
		return &proto.SpellStats{
			Id: spell.ActionID.ToProto(),

			IsCastable:      spell.Flags.Matches(SpellFlagAPL),
			IsChanneled:     spell.Flags.Matches(SpellFlagChanneled),
			IsMajorCooldown: spell.Flags.Matches(SpellFlagMCD),
			HasDot:          spell.dots != nil || spell.aoeDot != nil,
			HasShield:       spell.shields != nil || spell.selfShield != nil,
			PrepullOnly:     spell.Flags.Matches(SpellFlagPrepullOnly),
			EncounterOnly:   spell.Flags.Matches(SpellFlagEncounterOnly),
			HasCastTime:     spell.DefaultCast.CastTime > 0,
		}
	})

	aplAuras := FilterSlice(unit.auras, func(aura *Aura) bool {
		return !aura.ActionID.IsEmptyAction()
	})
	metadata.Auras = MapSlice(aplAuras, func(aura *Aura) *proto.AuraStats {
		return &proto.AuraStats{
			Id:                 aura.ActionID.ToProto(),
			MaxStacks:          aura.MaxStacks,
			HasIcd:             aura.Icd != nil,
			HasExclusiveEffect: len(aura.ExclusiveEffects) > 0,
		}
	})

	return metadata
}

func (unit *Unit) ExecuteCustomRotation(sim *Simulation) {
	panic("Unimplemented ExecuteCustomRotation")
}
