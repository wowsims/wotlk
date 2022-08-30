package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type UnitType int

const (
	PlayerUnit UnitType = iota
	EnemyUnit
	PetUnit
)

type DynamicDamageTakenModifier func(sim *Simulation, spellEffect *SpellEffect)

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

	// How far this unit is from its target(s). Measured in yards, this is used
	// for calculating spell travel time for certain spells.
	DistanceFromTarget float64

	// Environment in which this Unit exists. This will be nil until after the
	// construction phase.
	Env *Environment

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

	healthBar
	manaBar
	rageBar
	energyBar
	RunicPowerBar

	// All spells that can be cast by this unit.
	Spellbook []*Spell

	// Pets owned by this Unit.
	Pets []PetAgent

	// AutoAttacks is the manager for auto attack swings.
	// Must be enabled to use, with "EnableAutoAttacks()".
	AutoAttacks AutoAttacks

	// Statistics describing the results of the sim.
	Metrics UnitMetrics

	cdTimers []*Timer

	AttackTables                []*AttackTable
	DefenseTables               []*AttackTable
	DynamicDamageTakenModifiers []DynamicDamageTakenModifier

	GCD       *Timer
	doNothing bool // flags that this character chose to do nothing.

	// Used for applying the effects of hardcast / channeled spells at a later time.
	// By definition there can be only 1 hardcast spell being cast at any moment.
	Hardcast Hardcast

	// GCD-related PendingActions.
	gcdAction      *PendingAction
	hardcastAction *PendingAction

	// Fields related to waiting for certain events to happen.
	waitingForEnergy float64
	waitingForMana   float64
	waitStartTime    time.Duration

	// Cached mana return values per tick.
	manaTickWhileCasting    float64
	manaTickWhileNotCasting float64

	CastSpeed float64

	CurrentTarget *Unit
}

// DoNothing will explicitly declare that the character is intentionally doing nothing.
//
//	If the GCD is not used during OnGCDReady and this flag is set, OnGCDReady will not be called again
//	until it is used in some other way (like from an auto attack or resource regeneration).
func (char *Character) DoNothing() {
	char.doNothing = true
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
	if unit.Env == nil || !unit.Env.IsFinalized() {
		panic("Not finalized, use AddStats instead!")
	}

	unit.statsWithoutDeps = unit.statsWithoutDeps.Add(bonus)

	bonus = unit.ApplyStatDependencies(bonus)

	if sim.Log != nil {
		unit.Log(sim, "Dynamic stat change: %s", bonus.FlatString())
	}

	unit.stats = unit.stats.Add(bonus)
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
		unit.AutoAttacks.UpdateSwingTime(sim)
	}
	if bonus[stats.SpellHaste] != 0 {
		unit.updateCastSpeed()
	}
	if bonus[stats.Armor] != 0 {
		unit.updateArmor()
	}
	if bonus[stats.ArmorPenetration] != 0 {
		unit.updateArmorPen()
	}
	if bonus[stats.SpellPenetration] != 0 {
		unit.updateSpellPen()
	}
	if bonus[stats.ArcaneResistance] != 0 {
		unit.updateResistances()
	}
	if bonus[stats.FireResistance] != 0 {
		unit.updateResistances()
	}
	if bonus[stats.FrostResistance] != 0 {
		unit.updateResistances()
	}
	if bonus[stats.NatureResistance] != 0 {
		unit.updateResistances()
	}
	if bonus[stats.ShadowResistance] != 0 {
		unit.updateResistances()
	}

	if len(unit.Pets) > 0 {
		for _, petAgent := range unit.Pets {
			petAgent.GetPet().addOwnerStats(sim, bonus)
		}
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
			unit.Log(sim, "Dynamic dep enabled (%s): %s", dep.String(), unit.stats.Subtract(oldStats).FlatString())
		}
	}
}

func (unit *Unit) updateArmor() {
	for _, table := range unit.DefenseTables {
		if table != nil {
			table.UpdateArmorDamageReduction()
		}
	}
}
func (unit *Unit) updateArmorPen() {
	for _, table := range unit.AttackTables {
		if table != nil {
			table.UpdateArmorDamageReduction()
		}
	}
}
func (unit *Unit) updateResistances() {
	for _, table := range unit.DefenseTables {
		if table != nil {
			table.UpdatePartialResists()
		}
	}
}
func (unit *Unit) updateSpellPen() {
	for _, table := range unit.AttackTables {
		if table != nil {
			table.UpdatePartialResists()
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
	return MaxDuration(GCDMin, unit.ApplyCastSpeed(GCDDefault))
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

func (unit *Unit) ArmorPenetrationPercentage(armorPenRating float64) float64 {
	return MaxFloat(MinFloat(armorPenRating/ArmorPenPerPercentArmor, 100.0)*0.01, 0.0)
}

func (unit *Unit) RangedSwingSpeed() float64 {
	return unit.PseudoStats.RangedSpeedMultiplier * (1 + (unit.stats[stats.MeleeHaste] / (HasteRatingPerHastePercent * 100)))
}

// MultiplyMeleeSpeed will alter the attack speed multiplier and change swing speed of all autoattack swings in progress.
func (unit *Unit) MultiplyMeleeSpeed(sim *Simulation, amount float64) {
	unit.PseudoStats.MeleeSpeedMultiplier *= amount
	unit.AutoAttacks.UpdateSwingTime(sim)
}

func (unit *Unit) MultiplyRangedSpeed(sim *Simulation, amount float64) {
	unit.PseudoStats.RangedSpeedMultiplier *= amount
}

// Helper for when both MultiplyMeleeSpeed and MultiplyRangedSpeed are needed.
func (unit *Unit) MultiplyAttackSpeed(sim *Simulation, amount float64) {
	unit.PseudoStats.MeleeSpeedMultiplier *= amount
	unit.PseudoStats.RangedSpeedMultiplier *= amount
	unit.AutoAttacks.UpdateSwingTime(sim)
}

func (unit *Unit) finalize() {
	if unit.Env.IsFinalized() {
		panic("Unit already finalized!")
	}

	// Make sure we dont accidentally set initial stats instead of stats.
	if !unit.initialStats.Equals(stats.Stats{}) {
		panic("Initial stats may not be set before finalized: " + unit.initialStats.String())
	}

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

	for _, spell := range unit.Spellbook {
		spell.finalize()
	}
}

func (unit *Unit) init(sim *Simulation) {
	unit.auraTracker.init(sim)
}

func (unit *Unit) reset(sim *Simulation, agent Agent) {
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
	unit.healthBar.reset(sim)
	unit.UpdateManaRegenRates()

	unit.energyBar.reset(sim)
	unit.rageBar.reset(sim)
	unit.RunicPowerBar.reset(sim)

	unit.AutoAttacks.reset(sim)
}

// Advance moves time forward counting down auras, CDs, mana regen, etc
func (unit *Unit) advance(sim *Simulation, elapsedTime time.Duration) {
	unit.auraTracker.advance(sim)

	if unit.Hardcast.Expires != 0 && unit.Hardcast.Expires <= sim.CurrentTime {
		unit.Hardcast.Expires = 0
		unit.Hardcast.OnExpire(sim)
	}
}

func (unit *Unit) doneIteration(sim *Simulation) {
	unit.Hardcast = Hardcast{}
	unit.doneIterationGCD(sim.CurrentTime)

	unit.doneIterationMana()
	unit.rageBar.doneIteration()

	unit.auraTracker.doneIteration(sim)
	for _, spell := range unit.Spellbook {
		spell.doneIteration()
	}
	unit.resetCDs(sim)
}
