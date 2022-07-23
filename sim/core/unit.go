package core

import (
	"errors"
	"fmt"
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

// Unit is an abstraction of a Character/Boss/Pet/etc, containing functionality
// shared by all of them.
type Unit struct {
	Type UnitType

	// Index of this unit with its group.
	//  For Players, this is the 0-indexed raid index (0-24).
	//  For Enemies, this is its enemy index.
	//  For Pets, this is the same as the owner's index.
	Index int32

	// Index of this unit as it appears in attack/defense tables.
	// This is different from Index because there can be gaps in the raid.
	TableIndex int32

	// Unique label for logging.
	Label string

	Level int32 // Level of Unit, e.g. Bosses are 83.

	MobType proto.MobType

	// Environment in which this Unit exists. This will be nil until after the
	// construction phase.
	Env *Environment

	// Stats this Unit will have at the very start of each Sim iteration.
	// Includes all equipment / buffs / permanent effects but not temporary
	// effects from items / abilities.
	initialStats stats.Stats

	initialPseudoStats stats.PseudoStats

	// Cast speed without any temporary effects.
	initialCastSpeed float64

	// Melee swing speed without any temporary effects.
	initialMeleeSwingSpeed float64

	// Ranged swing speed without any temporary effects.
	initialRangedSwingSpeed float64

	// Provides aura tracking behavior.
	auraTracker

	// Current stats, including temporary effects.
	stats stats.Stats

	// Provides stat dependency management behavior.
	statBonuses [stats.Len]stats.Bonuses

	PseudoStats stats.PseudoStats

	healthBar
	manaBar
	rageBar
	energyBar
	runicPowerBar

	// All spells that can be cast by this unit.
	Spellbook []*Spell

	// AutoAttacks is the manager for auto attack swings.
	// Must be enabled to use, with "EnableAutoAttacks()".
	AutoAttacks AutoAttacks

	// Statistics describing the results of the sim.
	Metrics UnitMetrics

	cdTimers []*Timer

	AttackTables  []*AttackTable
	DefenseTables []*AttackTable

	GCD       *Timer
	doNothing bool // flags that this character chose to do nothing.

	// Used for applying the effects of hardcast / channeled spells at a later time.
	// By definition there can be only 1 hardcast spell being cast at any moment.
	Hardcast Hardcast

	// GCD-related PendingActions.
	gcdAction      *PendingAction
	hardcastAction *PendingAction

	// Fields related to waiting for certain events to happen.
	waitingForMana float64
	waitStartTime  time.Duration

	// Cached mana return values per tick.
	manaTickWhileCasting    float64
	manaTickWhileNotCasting float64

	CastSpeed float64

	CurrentTarget *Unit
}

// DoNothing will explicitly declare that the character is intentionally doing nothing.
//  If the GCD is not used during OnGCDReady and this flag is set, OnGCDReady will not be called again
//  until it is used in some other way (like from an auto attack or resource regeneration).
func (char *Character) DoNothing() {
	char.doNothing = true
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

func (unit *Unit) AddStatsDynamic(sim *Simulation, stat stats.Stats) {
	if unit.Env == nil || !unit.Env.IsFinalized() {
		panic("Not finalized, use AddStats instead!")
	}

	stat[stats.Mana] = 0 // TODO: Mana needs special treatment

	for k, v := range stat {
		unit.AddStatDynamic(sim, stats.Stat(k), v)
	}
}
func (unit *Unit) AddStatDynamic(sim *Simulation, stat stats.Stat, amount float64) {
	if unit.Env == nil || !unit.Env.IsFinalized() {
		panic("Not finalized, use AddStats instead!")
	}

	added := amount * unit.statBonuses[stat].Multiplier

	if stat == stats.MeleeHaste {
		unit.AddMeleeHaste(sim, added)
	} else {
		unit.stats[stat] += added

		if stat == stats.MP5 || stat == stats.Intellect || stat == stats.Spirit {
			unit.UpdateManaRegenRates()
		} else if stat == stats.SpellHaste {
			unit.updateCastSpeed()
		} else if stat == stats.Armor {
			unit.updateArmor()
		} else if stat == stats.ArmorPenetration {
			unit.updateArmorPen()
		} else if stat == stats.SpellPenetration {
			unit.updateSpellPen()
		} else if stat == stats.ArcaneResistance {
			unit.updateResistances()
		} else if stat == stats.FireResistance {
			unit.updateResistances()
		} else if stat == stats.FrostResistance {
			unit.updateResistances()
		} else if stat == stats.NatureResistance {
			unit.updateResistances()
		} else if stat == stats.ShadowResistance {
			unit.updateResistances()
		}
	}

	// Now apply stat dependencies
	for k, v := range unit.statBonuses[stat].Deps {
		if v == 1 {
			continue
		}
		unit.AddStatDynamic(sim, k, (v-1)*added) // this should handle descending
	}
}

// applyStatDependencies will apply all stat dependencies.
func (unit *Unit) applyStatDependencies(ss stats.Stats) stats.Stats {
	news := stats.Stats{}

	var addstat func(s stats.Stat, v float64)

	addstat = func(s stats.Stat, v float64) {
		if unit.statBonuses[s].Multiplier == 0 {
			unit.statBonuses[s].Multiplier = 1
		}
		added := v * unit.statBonuses[s].Multiplier
		news[s] += added
		for k, v := range unit.statBonuses[s].Deps {
			if v == 1 {
				continue
			}
			addstat(k, (v-1)*added)
		}
	}

	for s, v := range ss {
		if v == 0 {
			continue
		}
		addstat(stats.Stat(s), v)
	}

	return news
}

// AddStatDependency will add source stat * ratio to the modified stat.
func (unit *Unit) AddStatDependency(source, modified stats.Stat, multiplier float64) {
	if unit.Env != nil && unit.Env.IsFinalized() {
		panic("Already finalized, can't add more dependencies!")
	}
	if source == modified {
		if unit.statBonuses[source].Multiplier == 0 {
			unit.statBonuses[source].Multiplier = multiplier
		} else {
			unit.statBonuses[source].Multiplier *= multiplier
		}
		return
	}
	if unit.statBonuses[source].Deps == nil {
		unit.statBonuses[source].Deps = map[stats.Stat]float64{
			modified: multiplier,
		}
	} else if unit.statBonuses[source].Deps[modified] == 0 {
		unit.statBonuses[source].Deps[modified] = multiplier
	} else {
		unit.statBonuses[source].Deps[modified] *= multiplier
	}
}

// AddStatDependencyDynamic will dynamically adjust stats based on the change to the dependency.
func (unit *Unit) AddStatDependencyDynamic(sim *Simulation, source, modified stats.Stat, multiplier float64) {
	if unit.Env == nil || !unit.Env.IsFinalized() {
		panic("Not finalized, use AddStatDependency instead!")
	}
	if source == modified {
		oldMultiplier := 1.0
		if unit.statBonuses[source].Multiplier == 0 {
			unit.statBonuses[source].Multiplier = multiplier
		} else {
			oldMultiplier = unit.statBonuses[source].Multiplier
			unit.statBonuses[source].Multiplier *= multiplier
		}
		// Now modify the stat itself
		stat := unit.stats[source]
		bonus := ((stat * unit.statBonuses[source].Multiplier / oldMultiplier) - stat) / unit.statBonuses[source].Multiplier
		unit.AddStatDynamic(sim, source, bonus)
		return
	}

	oldMultiplier := 1.0
	if unit.statBonuses[source].Deps == nil {
		unit.statBonuses[source].Deps = map[stats.Stat]float64{
			modified: multiplier,
		}
	} else if unit.statBonuses[source].Deps[modified] == 0 {
		unit.statBonuses[source].Deps[modified] = multiplier
	} else {
		oldMultiplier = unit.statBonuses[source].Deps[modified]
		unit.statBonuses[source].Deps[modified] *= multiplier
	}

	stat := unit.stats[source]
	bonus := ((stat * unit.statBonuses[source].Deps[modified] / oldMultiplier) - stat) / unit.statBonuses[source].Deps[modified]
	// Now apply the newly gained stats
	unit.AddStatDynamic(sim, modified, bonus)
}

// finalizeStatDeps will descend the tree of each stat's depedencies and verify
// there are no circular dependencies
func (unit *Unit) finalizeStatDeps() {
	seen := map[stats.Stat]struct{}{}

	var walk func(m map[stats.Stat]float64) error

	walk = func(m map[stats.Stat]float64) error {
		for k := range m {
			if _, ok := seen[k]; ok {
				return errors.New("circular dependency in stats: " + k.StatName())
			}
			seen[k] = struct{}{}
			err := walk(unit.statBonuses[k].Deps)
			if err != nil {
				return fmt.Errorf("%w from: %s", err, k.StatName())
			}
			delete(seen, k)
		}
		return nil
	}

	for s := range unit.stats {
		if unit.statBonuses[s].Multiplier == 0 {
			unit.statBonuses[s].Multiplier = 1
		}
		seen[stats.Stat(s)] = struct{}{}
		if err := walk(unit.statBonuses[s].Deps); err != nil {
			panic(err)
		}
		delete(seen, stats.Stat(s))
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

func (unit *Unit) SwingSpeed() float64 {
	return unit.PseudoStats.MeleeSpeedMultiplier * (1 + (unit.stats[stats.MeleeHaste] / (unit.PseudoStats.MeleeHasteRatingPerHastePercent * 100)))
}

func (unit *Unit) Armor() float64 {
	return unit.PseudoStats.ArmorMultiplier * unit.stats[stats.Armor]
}

func (unit *Unit) ArmorPenetrationPercentage() float64 {
	return MaxFloat(MinFloat(unit.stats[stats.ArmorPenetration]/ArmorPenPerPercentArmor, 100.0)*0.01, 0.0)
}

func (unit *Unit) RangedSwingSpeed() float64 {
	return unit.PseudoStats.RangedSpeedMultiplier * (1 + (unit.stats[stats.MeleeHaste] / (HasteRatingPerHastePercent * 100)))
}

func (unit *Unit) AddMeleeHaste(sim *Simulation, amount float64) {
	if amount > 0 {
		mod := 1 + (amount / (unit.PseudoStats.MeleeHasteRatingPerHastePercent * 100))
		unit.AutoAttacks.ModifySwingTime(sim, mod)
	} else {
		mod := 1 / (1 + (-amount / (unit.PseudoStats.MeleeHasteRatingPerHastePercent * 100)))
		unit.AutoAttacks.ModifySwingTime(sim, mod)
	}
	unit.stats[stats.MeleeHaste] += amount

	// Could add melee haste to pets too, but not aware of any pets that scale with
	// owner's melee haste.
}

// MultiplyMeleeSpeed will alter the attack speed multiplier and change swing speed of all autoattack swings in progress.
func (unit *Unit) MultiplyMeleeSpeed(sim *Simulation, amount float64) {
	unit.PseudoStats.MeleeSpeedMultiplier *= amount
	unit.AutoAttacks.ModifySwingTime(sim, amount)
}

func (unit *Unit) MultiplyRangedSpeed(sim *Simulation, amount float64) {
	unit.PseudoStats.RangedSpeedMultiplier *= amount
}

// Helper for when both MultiplyMeleeSpeed and MultiplyRangedSpeed are needed.
func (unit *Unit) MultiplyAttackSpeed(sim *Simulation, amount float64) {
	unit.PseudoStats.MeleeSpeedMultiplier *= amount
	unit.PseudoStats.RangedSpeedMultiplier *= amount
	unit.AutoAttacks.ModifySwingTime(sim, amount)
}

func (unit *Unit) finalize() {
	if unit.Env.IsFinalized() {
		panic("Unit already finalized!")
	}

	// Make sure we dont accidentally set initial stats instead of stats.
	if !unit.initialStats.Equals(stats.Stats{}) {
		panic("Initial stats may not be set before finalized: " + unit.initialStats.String())
	}

	unit.updateCastSpeed()

	// All stats added up to this point are part of the 'initial' stats.
	unit.initialStats = unit.stats
	unit.initialPseudoStats = unit.PseudoStats
	unit.initialCastSpeed = unit.CastSpeed
	unit.initialMeleeSwingSpeed = unit.SwingSpeed()
	unit.initialRangedSwingSpeed = unit.RangedSwingSpeed()

	unit.applyParryHaste()

	for _, spell := range unit.Spellbook {
		spell.finalize()
	}
}

func (unit *Unit) init(sim *Simulation) {
	unit.auraTracker.init(sim)
}

func (unit *Unit) reset(sim *Simulation, agent Agent) {
	unit.Metrics.reset()
	unit.stats = unit.initialStats
	unit.PseudoStats = unit.initialPseudoStats
	unit.auraTracker.reset(sim)
	for _, spell := range unit.Spellbook {
		spell.reset(sim)
	}

	unit.healthBar.reset(sim)
	unit.UpdateManaRegenRates()

	unit.energyBar.reset(sim)
	unit.rageBar.reset(sim)
	unit.runicPowerBar.reset(sim)

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
