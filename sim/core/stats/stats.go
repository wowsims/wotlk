package stats

import (
	"fmt"
	"strings"
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
)

type Stats [Len]float64

type Stat byte

// Use internal representation instead of proto.Stat so we can add functions
// and use 'byte' as the data type.
//
// This needs to stay synced with proto.Stat.
const (
	Strength Stat = iota
	Agility
	Stamina
	Intellect
	Spirit
	SpellPower
	HealingPower
	ArcaneSpellPower
	FireSpellPower
	FrostSpellPower
	HolySpellPower
	NatureSpellPower
	ShadowSpellPower
	MP5
	SpellHit
	SpellCrit
	SpellHaste
	SpellPenetration
	AttackPower
	MeleeHit
	MeleeCrit
	MeleeHaste
	ArmorPenetration
	Expertise
	Mana
	Energy
	Rage
	Armor
	RangedAttackPower
	Defense
	Block
	BlockValue
	Dodge
	Parry
	Resilience
	Health
	ArcaneResistance
	FireResistance
	FrostResistance
	NatureResistance
	ShadowResistance
	FeralAttackPower

	Len
)

func ProtoArrayToStatsList(protoStats []proto.Stat) []Stat {
	stats := make([]Stat, len(protoStats))
	for i, v := range protoStats {
		stats[i] = Stat(v)
	}
	return stats
}

func (s Stat) StatName() string {
	switch s {
	case Strength:
		return "Strength"
	case Agility:
		return "Agility"
	case Stamina:
		return "Stamina"
	case Intellect:
		return "Intellect"
	case Spirit:
		return "Spirit"
	case SpellCrit:
		return "SpellCrit"
	case SpellHit:
		return "SpellHit"
	case HealingPower:
		return "HealingPower"
	case SpellPower:
		return "SpellPower"
	case SpellHaste:
		return "SpellHaste"
	case MP5:
		return "MP5"
	case SpellPenetration:
		return "SpellPenetration"
	case FireSpellPower:
		return "FireSpellPower"
	case NatureSpellPower:
		return "NatureSpellPower"
	case FrostSpellPower:
		return "FrostSpellPower"
	case ShadowSpellPower:
		return "ShadowSpellPower"
	case HolySpellPower:
		return "HolySpellPower"
	case ArcaneSpellPower:
		return "ArcaneSpellPower"
	case AttackPower:
		return "AttackPower"
	case MeleeHit:
		return "MeleeHit"
	case MeleeHaste:
		return "MeleeHaste"
	case MeleeCrit:
		return "MeleeCrit"
	case Expertise:
		return "Expertise"
	case ArmorPenetration:
		return "ArmorPenetration"
	case Mana:
		return "Mana"
	case Energy:
		return "Energy"
	case Rage:
		return "Rage"
	case Armor:
		return "Armor"
	case RangedAttackPower:
		return "RangedAttackPower"
	case FeralAttackPower:
		return "FeralAttackPower"
	case Defense:
		return "Defense"
	case Block:
		return "Block"
	case BlockValue:
		return "BlockValue"
	case Dodge:
		return "Dodge"
	case Parry:
		return "Parry"
	case Resilience:
		return "Resilience"
	case Health:
		return "Health"
	case FireResistance:
		return "FireResistance"
	case NatureResistance:
		return "NatureResistance"
	case FrostResistance:
		return "FrostResistance"
	case ShadowResistance:
		return "ShadowResistance"
	case ArcaneResistance:
		return "ArcaneResistance"
	}

	return "none"
}

func FromFloatArray(values []float64) Stats {
	stats := Stats{}
	for i, v := range values {
		stats[i] = v
	}
	return stats
}

// Adds two Stats together, returning the new Stats.
func (stats Stats) Add(other Stats) Stats {
	newStats := Stats{}

	for i, thisStat := range stats {
		newStats[i] = thisStat + other[i]
	}

	return newStats
}

// Subtracts another Stats from this one, returning the new Stats.
func (stats Stats) Subtract(other Stats) Stats {
	newStats := Stats{}

	for k, v := range stats {
		newStats[k] = v - other[k]
	}

	return newStats
}

func (stats Stats) Multiply(multiplier float64) Stats {
	newStats := stats
	for k, v := range newStats {
		newStats[k] = v * multiplier
	}
	return newStats
}

// Multiplies two Stats together by multiplying the values of corresponding
// stats, like a dot product operation.
func (stats Stats) DotProduct(other Stats) Stats {
	newStats := Stats{}

	for k, v := range stats {
		newStats[k] = v * other[k]
	}

	return newStats
}

func (stats Stats) Equals(other Stats) bool {
	for i := range stats {
		if stats[i] != other[i] {
			return false
		}
	}

	return true
}

func (stats Stats) EqualsWithTolerance(other Stats, tolerance float64) bool {
	for i := range stats {
		if stats[i] < other[i]-tolerance || stats[i] > other[i]+tolerance {
			return false
		}
	}

	return true
}

func (stats Stats) String() string {
	var sb strings.Builder
	sb.WriteString("\n{\n")

	for statIdx, statValue := range stats {
		name := Stat(statIdx).StatName()
		if name == "none" || statValue == 0 {
			continue
		}

		fmt.Fprintf(&sb, "\t%s: %0.3f,\n", name, statValue)
	}

	sb.WriteString("\n}")
	return sb.String()
}

// Like String() but without the newlines.
func (stats Stats) FlatString() string {
	var sb strings.Builder
	sb.WriteString("{")

	for statIdx, statValue := range stats {
		name := Stat(statIdx).StatName()
		if name == "none" || statValue == 0 {
			continue
		}

		fmt.Fprintf(&sb, "%s: %0.3f,", name, statValue)
	}

	sb.WriteString("}")
	return sb.String()
}

func (stats Stats) ToFloatArray() []float64 {
	arr := make([]float64, len(stats))
	for i, v := range stats {
		arr[i] = v
	}
	return arr
}

// Given the current values for source and mod stats, should return the new
// value for the mod stat.
type StatModifier func(sourceValue float64, modValue float64) float64

// Represents a dependency between two stats, whereby the value of one stat
// modifies the value of the other.
//
// For example, many casters have a talent to increase their spell power by
// a percentage of their intellect.
type StatDependency struct {
	// The stat which will be used to control the amount of increase.
	SourceStat Stat

	// The stat which will be modified, depending on the value of SourceStat.
	ModifiedStat Stat

	// Applies the stat modification.
	Modifier StatModifier
}

type StatDependencyManager struct {
	// Stat dependencies for each stat.
	// First dimension is the modified stat. For each modified stat, stores a list of
	// dependencies for that stat.
	deps [Len][]StatDependency

	// Whether Finalize() has been called.
	finalized bool

	// Dependencies being managed, sorted so that their modifiers can be applied
	// in-order without any issues.
	sortedDeps []StatDependency
}

func (sdm *StatDependencyManager) AddStatDependency(dep StatDependency) {
	if sdm.finalized {
		panic("Stat dependencies may not be added once finalized!")
	}

	sdm.deps[dep.ModifiedStat] = append(sdm.deps[dep.ModifiedStat], dep)
}

// Populates sortedDeps. Panics if there are any dependency cycles.
// TODO: Figure out if we need to separate additive / multiplicative dependencies.
func (sdm *StatDependencyManager) Sort() {
	sdm.sortedDeps = []StatDependency{}

	// Set of stats we're done processing.
	processedStats := map[Stat]struct{}{}

	for len(processedStats) < int(Len) {
		numNewlyProcessed := 0
		for i := 0; i < int(Len); i++ {
			stat := Stat(i)

			if _, alreadyProcessed := processedStats[stat]; alreadyProcessed {
				continue
			}

			// If all deps for this stat have been processed or are the same stat, we can process it.
			allDepsProcessed := true
			for _, dep := range sdm.deps[stat] {
				_, depAlreadyProcessed := processedStats[dep.SourceStat]

				if !depAlreadyProcessed && dep.SourceStat != stat {
					allDepsProcessed = false
				}
			}
			if !allDepsProcessed {
				continue
			}

			// Process this stat by adding its deps to sortedDeps.

			// Add deps from other stats first.
			for _, dep := range sdm.deps[stat] {
				if dep.SourceStat != stat {
					sdm.sortedDeps = append(sdm.sortedDeps, dep)
				}
			}

			// Now add deps from the same stat.
			for _, dep := range sdm.deps[stat] {
				if dep.SourceStat == stat {
					sdm.sortedDeps = append(sdm.sortedDeps, dep)
				}
			}

			// Mark this stat as processed.
			processedStats[stat] = struct{}{}
			numNewlyProcessed++
		}

		// If we couldn't process any new stats but there are still stats left,
		// there must be a circular dependency.
		if numNewlyProcessed == 0 {
			panic("Circular stat dependency detected")
		}
	}
}

func (sdm *StatDependencyManager) Finalize() {
	if sdm.finalized {
		return
	}
	sdm.finalized = true

	sdm.Sort()
}

// Applies all stat dependencies and returns the new Stats.
func (sdm *StatDependencyManager) ApplyStatDependencies(stats Stats) Stats {
	newStats := stats
	for _, dep := range sdm.sortedDeps {
		newStats[dep.ModifiedStat] = dep.Modifier(newStats[dep.SourceStat], newStats[dep.ModifiedStat])
	}

	return newStats
}
func (sdm *StatDependencyManager) SortAndApplyStatDependencies(stats Stats) Stats {
	sdm.Sort()
	return sdm.ApplyStatDependencies(stats)
}

type PseudoStats struct {
	///////////////////////////////////////////////////
	// Effects that apply when this unit is the attacker.
	///////////////////////////////////////////////////

	NoCost         bool    // If set, spells cost no mana/energy/rage.
	CostMultiplier float64 // Multiplies spell cost.
	CostReduction  float64 // Reduces spell cost.

	CastSpeedMultiplier   float64
	MeleeSpeedMultiplier  float64
	RangedSpeedMultiplier float64

	FiveSecondRuleRefreshTime time.Duration // last time a spell was cast
	SpiritRegenRateCasting    float64       // percentage of spirit regen allowed during casting

	// Both of these are currently only used for innervate.
	ForceFullSpiritRegen  bool    // If set, automatically uses full spirit regen regardless of FSR refresh time.
	SpiritRegenMultiplier float64 // Multiplier on spirit portion of mana regen.

	// If true, allows block/parry.
	InFrontOfTarget bool

	// "Apply Aura: Mod Damage Done (Physical)", applies to abilities with EffectSpellCoefficient > 0.
	//  This includes almost all "(Normalized) Weapon Damage", but also some "School Damage (Physical)" abilities.
	BonusDamage float64 // Comes from '+X Weapon Damage' effects

	BonusRangedHitRating  float64 // Hit rating for ranged only.
	BonusMeleeCritRating  float64 // Crit rating for melee only (not ranged).
	BonusRangedCritRating float64 // Crit rating for ranged only.
	BonusFireCritRating   float64 // Crit rating for fire spells only (Combustion).
	BonusMHCritRating     float64 // Talents, e.g. Rogue Dagger specialization
	BonusOHCritRating     float64 // Talents, e.g. Rogue Dagger specialization

	DisableDWMissPenalty bool    // Used by Heroic Strike and Cleave
	IncreasedMissChance  float64 // Insect Swarm and Scorpid Sting
	DodgeReduction       float64 // Used by Warrior talent 'Weapon Mastery' and SWP boss auras.

	MobTypeAttackPower float64 // Bonus AP against mobs of the current type.
	MobTypeSpellPower  float64 // Bonus SP against mobs of the current type.

	// For Human and Orc weapon racials
	BonusMHExpertiseRating float64
	BonusOHExpertiseRating float64

	ThreatMultiplier          float64 // Modulates the threat generated. Affected by things like salv.
	HolySpellThreatMultiplier float64 // Righteous Fury

	DamageDealtMultiplier       float64 // All damage
	RangedDamageDealtMultiplier float64

	PhysicalDamageDealtMultiplier float64
	ArcaneDamageDealtMultiplier   float64
	FireDamageDealtMultiplier     float64
	FrostDamageDealtMultiplier    float64
	HolyDamageDealtMultiplier     float64
	NatureDamageDealtMultiplier   float64
	ShadowDamageDealtMultiplier   float64

	// Modifiers for spells with the SpellFlagAgentReserved1 flag set.
	BonusCritRatingAgentReserved1       float64
	AgentReserved1DamageDealtMultiplier float64

	///////////////////////////////////////////////////
	// Effects that apply when this unit is the target.
	///////////////////////////////////////////////////

	CanBlock bool
	CanParry bool
	CanCrush bool

	ParryHaste bool

	ReducedCritTakenChance float64 // Reduces chance to be crit.

	BonusMeleeAttackPower  float64 // Imp Hunters mark, EW
	BonusRangedAttackPower float64 // Hunters mark, EW
	BonusCritRating        float64 // Imp Judgement of the Crusader
	BonusFrostCritRating   float64 // Winter's Chill
	BonusMeleeHitRating    float64 // Imp FF

	BonusDamageTaken         float64 // Blessing of Sanctuary
	BonusPhysicalDamageTaken float64 // Hemo, Gift of Arthas, etc
	BonusHolyDamageTaken     float64 // Judgement of the Crusader

	DamageTakenMultiplier float64 // All damage

	PhysicalDamageTakenMultiplier float64
	ArcaneDamageTakenMultiplier   float64
	FireDamageTakenMultiplier     float64
	FrostDamageTakenMultiplier    float64
	HolyDamageTakenMultiplier     float64
	NatureDamageTakenMultiplier   float64
	ShadowDamageTakenMultiplier   float64

	PeriodicPhysicalDamageTakenMultiplier float64
}

func NewPseudoStats() PseudoStats {
	return PseudoStats{
		CostMultiplier: 1,

		CastSpeedMultiplier:  1,
		MeleeSpeedMultiplier: 1,
		//RangedSpeedMultiplier: 1, // Leave at 0 so we can use this to ignore ranged stuff for non-hunters.
		SpiritRegenMultiplier: 1,

		ThreatMultiplier:          1,
		HolySpellThreatMultiplier: 1,

		DamageDealtMultiplier:       1,
		RangedDamageDealtMultiplier: 1,

		PhysicalDamageDealtMultiplier: 1,
		ArcaneDamageDealtMultiplier:   1,
		FireDamageDealtMultiplier:     1,
		FrostDamageDealtMultiplier:    1,
		HolyDamageDealtMultiplier:     1,
		NatureDamageDealtMultiplier:   1,
		ShadowDamageDealtMultiplier:   1,

		AgentReserved1DamageDealtMultiplier: 1,

		// Target effects.
		DamageTakenMultiplier: 1,

		PhysicalDamageTakenMultiplier: 1,
		ArcaneDamageTakenMultiplier:   1,
		FireDamageTakenMultiplier:     1,
		FrostDamageTakenMultiplier:    1,
		HolyDamageTakenMultiplier:     1,
		NatureDamageTakenMultiplier:   1,
		ShadowDamageTakenMultiplier:   1,

		PeriodicPhysicalDamageTakenMultiplier: 1,
	}
}
