package core

import (
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
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

	Level int32 // Level of Unit, e.g. Bosses are lvl 73.

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

	PseudoStats stats.PseudoStats

	// TODO: Put these inside a 'manaBar' object.
	manaCastingMetrics    *ResourceMetrics
	manaNotCastingMetrics *ResourceMetrics
	JowManaMetrics        *ResourceMetrics
	VtManaMetrics         *ResourceMetrics

	healthBar
	rageBar
	energyBar

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

	GCD *Timer

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

	if stat[stats.MeleeHaste] != 0 {
		unit.AddMeleeHaste(sim, stat[stats.MeleeHaste])
		stat[stats.MeleeHaste] = 0
	}

	unit.stats = unit.stats.Add(stat)

	if stat[stats.MP5] != 0 || stat[stats.Intellect] != 0 || stat[stats.Spirit] != 0 {
		unit.UpdateManaRegenRates()
	}
	if stat[stats.SpellHaste] != 0 {
		unit.updateCastSpeed()
	}
	if stat[stats.Armor] != 0 {
		unit.updateArmor()
	}
	if stat[stats.ArmorPenetration] != 0 {
		unit.updateArmorPen()
	}
	if stat[stats.SpellPenetration] != 0 {
		unit.updateSpellPen()
	}
	if stat[stats.ArcaneResistance] != 0 || stat[stats.FireResistance] != 0 || stat[stats.FrostResistance] != 0 || stat[stats.NatureResistance] != 0 || stat[stats.ShadowResistance] != 0 {
		unit.updateResistances()
	}
}
func (unit *Unit) AddStatDynamic(sim *Simulation, stat stats.Stat, amount float64) {
	if unit.Env == nil || !unit.Env.IsFinalized() {
		panic("Not finalized, use AddStats instead!")
	}

	if stat == stats.MeleeHaste {
		unit.AddMeleeHaste(sim, amount)
		return
	}

	unit.stats[stat] += amount

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
	return unit.PseudoStats.MeleeSpeedMultiplier * (1 + (unit.stats[stats.MeleeHaste] / (HasteRatingPerHastePercent * 100)))
}

func (unit *Unit) RangedSwingSpeed() float64 {
	return unit.PseudoStats.RangedSpeedMultiplier * (1 + (unit.stats[stats.MeleeHaste] / (HasteRatingPerHastePercent * 100)))
}

func (unit *Unit) AddMeleeHaste(sim *Simulation, amount float64) {
	if amount > 0 {
		mod := 1 + (amount / (HasteRatingPerHastePercent * 100))
		unit.AutoAttacks.ModifySwingTime(sim, mod)
	} else {
		mod := 1 / (1 + (-amount / (HasteRatingPerHastePercent * 100)))
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
