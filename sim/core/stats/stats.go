package stats

import (
	"fmt"
	"strings"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
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
	RunicPower
	BloodRune
	FrostRune
	UnholyRune
	DeathRune

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
	case RunicPower:
		return "Runic Power"
	case BloodRune:
		return "Blood Rune"
	case FrostRune:
		return "Frost Rune"
	case UnholyRune:
		return "Unholy Rune"
	case DeathRune:
		return "Death Rune"
	}

	return "none"
}

func FromFloatArray(values []float64) Stats {
	stats := Stats{}
	copy(stats[:], values)
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
		fmt.Fprintf(&sb, "\"%s\": %0.3f,", name, statValue)
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

	BonusRangedHitRating      float64 // Hit rating for ranged only.
	BonusMeleeCritRating      float64 // Crit rating for melee only (not ranged).
	BonusRangedCritRating     float64 // Crit rating for ranged only.
	BonusFireCritRating       float64 // Crit rating for fire spells only.
	BonusShadowCritRating     float64 // Crit rating for shadow spells only. Warlock stuff. You wouldn't understand.
	BonusMHCritRating         float64 // Talents, e.g. Rogue Dagger specialization
	BonusOHCritRating         float64 // Talents, e.g. Rogue Dagger specialization
	BonusMeleeSpellCritRating float64 // Crit rating for melee special attacks, used for Warrior Recklessness
	BonusMHArmorPenRating     float64 // Talents, e.g. Rogue Mace specialization
	BonusOHArmorPenRating     float64 // Talents, e.g. Rogue Mace specialization
	BonusSpellCritRating      float64 // Crit rating bonus to spells

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
	DiseaseDamageDealtMultiplier  float64

	// Modifiers for spells with the SpellFlagAgentReserved1 flag set.
	BonusCritRatingAgentReserved1       float64
	AgentReserved1DamageDealtMultiplier float64

	// Treat melee haste as a pseudostat so that shamans, death knights, paladins, and druids can get the correct scaling
	MeleeHasteRatingPerHastePercent float64

	HealingDealtMultiplier float64

	///////////////////////////////////////////////////
	// Effects that apply when this unit is the target.
	///////////////////////////////////////////////////

	CanBlock bool
	CanParry bool

	ParryHaste bool

	ReducedCritTakenChance float64 // Reduces chance to be crit.

	BonusMeleeAttackPowerTaken  float64 // Imp Hunters mark, EW
	BonusRangedAttackPowerTaken float64 // Hunters mark, EW
	BonusSpellCritRatingTaken   float64 // Imp Shadow Bolt / Imp Scorch / Winter's Chill debuff
	BonusCritRatingTaken        float64 // Totem of Wrath / Master Poisoner / Heart of the Crusader
	BonusMeleeHitRatingTaken    float64 //
	BonusSpellHitRatingTaken    float64 // Imp FF

	BonusDamageTaken         float64 // Blessing of Sanctuary
	BonusPhysicalDamageTaken float64 // Hemo, Gift of Arthas, etc
	BonusHolyDamageTaken     float64 // Judgement of the Crusader

	DamageTakenMultiplier float64 // All damage

	ArmorMultiplier float64 // Major/minor/special multipicative armor modifiers

	PhysicalDamageTakenMultiplier float64
	ArcaneDamageTakenMultiplier   float64
	FireDamageTakenMultiplier     float64
	FrostDamageTakenMultiplier    float64
	HolyDamageTakenMultiplier     float64
	NatureDamageTakenMultiplier   float64
	ShadowDamageTakenMultiplier   float64
	DiseaseDamageTakenMultiplier  float64

	ReducedPhysicalHitTakenChance float64
	ReducedArcaneHitTakenChance   float64
	ReducedFireHitTakenChance     float64
	ReducedFrostHitTakenChance    float64
	ReducedNatureHitTakenChance   float64
	ReducedShadowHitTakenChance   float64

	PeriodicPhysicalDamageTakenMultiplier float64
	PeriodicShadowDamageTakenMultiplier   float64

	HealingTakenMultiplier float64
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

		PhysicalDamageDealtMultiplier:       1,
		ArcaneDamageDealtMultiplier:         1,
		FireDamageDealtMultiplier:           1,
		FrostDamageDealtMultiplier:          1,
		HolyDamageDealtMultiplier:           1,
		NatureDamageDealtMultiplier:         1,
		ShadowDamageDealtMultiplier:         1,
		DiseaseDamageDealtMultiplier:        1,
		AgentReserved1DamageDealtMultiplier: 1,

		MeleeHasteRatingPerHastePercent: 32.79,

		HealingDealtMultiplier: 1,

		// Target effects.
		DamageTakenMultiplier: 1,

		ArmorMultiplier: 1,

		PhysicalDamageTakenMultiplier: 1,
		ArcaneDamageTakenMultiplier:   1,
		FireDamageTakenMultiplier:     1,
		FrostDamageTakenMultiplier:    1,
		HolyDamageTakenMultiplier:     1,
		NatureDamageTakenMultiplier:   1,
		ShadowDamageTakenMultiplier:   1,
		DiseaseDamageTakenMultiplier:  1,

		PeriodicPhysicalDamageTakenMultiplier: 1,
		PeriodicShadowDamageTakenMultiplier:   1,

		HealingTakenMultiplier: 1,
	}
}
