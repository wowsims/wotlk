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
	BonusArmor
	RunicPower
	BloodRune
	FrostRune
	UnholyRune
	DeathRune
	// DO NOT add new stats here without discussing it first; new stats come with
	// a performance penalty.

	Len
)

var PseudoStatsLen = len(proto.PseudoStat_name)
var UnitStatsLen = int(Len) + PseudoStatsLen

type SchoolIndex byte

const (
	SchoolIndexNone     SchoolIndex = 0
	SchoolIndexPhysical SchoolIndex = iota
	SchoolIndexArcane
	SchoolIndexFire
	SchoolIndexFrost
	SchoolIndexHoly
	SchoolIndexNature
	SchoolIndexShadow

	SchoolLen
)

func NewSchoolFloatArray() [SchoolLen]float64 {
	return [SchoolLen]float64{
		1, 1, 1, 1, 1, 1, 1, 1,
	}
}

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
	case SpellPower:
		return "SpellPower"
	case SpellHaste:
		return "SpellHaste"
	case MP5:
		return "MP5"
	case SpellPenetration:
		return "SpellPenetration"
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
	case BonusArmor:
		return "BonusArmor"
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
	var stats Stats
	copy(stats[:], values)
	return stats
}

// Adds two Stats together, returning the new Stats.
func (stats Stats) Add(other Stats) Stats {
	for k := range stats {
		stats[k] += other[k]
	}
	return stats
}

// Adds another to Stats to this, in-place. For performance, only.
func (stats *Stats) AddInplace(other *Stats) {
	for k := range stats {
		stats[k] += other[k]
	}
}

// Subtracts another Stats from this one, returning the new Stats.
func (stats Stats) Subtract(other Stats) Stats {
	for k := range stats {
		stats[k] -= other[k]
	}
	return stats
}

func (stats Stats) Invert() Stats {
	for k, v := range stats {
		stats[k] = -v
	}
	return stats
}

func (stats Stats) Multiply(multiplier float64) Stats {
	for k := range stats {
		stats[k] *= multiplier
	}
	return stats
}

// Multiplies two Stats together by multiplying the values of corresponding
// stats, like a dot product operation.
func (stats Stats) DotProduct(other Stats) Stats {
	for k := range stats {
		stats[k] *= other[k]
	}
	return stats
}

func (stats Stats) Equals(other Stats) bool {
	return stats == other
}

func (stats Stats) EqualsWithTolerance(other Stats, tolerance float64) bool {
	for k, v := range stats {
		if v < other[k]-tolerance || v > other[k]+tolerance {
			return false
		}
	}
	return true
}

func (stats Stats) String() string {
	var sb strings.Builder
	sb.WriteString("\n{\n")

	for statIdx, statValue := range stats {
		if statValue == 0 {
			continue
		}
		if name := Stat(statIdx).StatName(); name != "none" {
			_, _ = fmt.Fprintf(&sb, "\t%s: %0.3f,\n", name, statValue)
		}
	}

	sb.WriteString("\n}")
	return sb.String()
}

// Like String() but without the newlines.
func (stats Stats) FlatString() string {
	var sb strings.Builder
	sb.WriteString("{")

	for statIdx, statValue := range stats {
		if statValue == 0 {
			continue
		}
		if name := Stat(statIdx).StatName(); name != "none" {
			_, _ = fmt.Fprintf(&sb, "\"%s\": %0.3f,", name, statValue)
		}
	}

	sb.WriteString("}")
	return sb.String()
}

func (stats Stats) ToFloatArray() []float64 {
	return stats[:]
}

type PseudoStats struct {
	///////////////////////////////////////////////////
	// Effects that apply when this unit is the attacker.
	///////////////////////////////////////////////////

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

	BonusMHDps     float64
	BonusOHDps     float64
	BonusRangedDps float64

	DisableDWMissPenalty bool    // Used by Heroic Strike and Cleave
	IncreasedMissChance  float64 // Insect Swarm and Scorpid Sting
	DodgeReduction       float64 // Used by Warrior talent 'Weapon Mastery' and SWP boss auras.

	MobTypeAttackPower float64 // Bonus AP against mobs of the current type.
	MobTypeSpellPower  float64 // Bonus SP against mobs of the current type.

	ThreatMultiplier float64 // Modulates the threat generated. Affected by things like salv.

	DamageDealtMultiplier       float64            // All damage
	SchoolDamageDealtMultiplier [SchoolLen]float64 // For specific spell schools (arcane, fire, shadow, etc).

	// Treat melee haste as a pseudostat so that shamans, death knights, paladins, and druids can get the correct scaling
	MeleeHasteRatingPerHastePercent float64

	// Important when unit is attacker or target
	BlockValueMultiplier float64

	// Only used for NPCs, governs variance in enemy auto-attack damage
	DamageSpread float64

	///////////////////////////////////////////////////
	// Effects that apply when this unit is the target.
	///////////////////////////////////////////////////

	CanBlock bool
	CanParry bool
	Stunned  bool // prevents blocks, dodges, and parries

	ParryHaste bool

	// Avoidance % not affected by Diminishing Returns
	BaseDodge float64
	BaseParry float64
	//BaseMiss is not needed, this is always 5%

	ReducedCritTakenChance float64 // Reduces chance to be crit.

	BonusRangedAttackPowerTaken float64 // Hunters mark
	BonusSpellCritRatingTaken   float64 // Imp Shadow Bolt / Imp Scorch / Winter's Chill debuff
	BonusCritRatingTaken        float64 // Totem of Wrath / Master Poisoner / Heart of the Crusader
	BonusMeleeHitRatingTaken    float64 // Formerly Imp FF and SW Radiance;
	BonusSpellHitRatingTaken    float64 // Imp FF

	BonusPhysicalDamageTaken float64 // Hemo, Gift of Arthas, etc
	BonusHealingTaken        float64 // Talisman of Troll Divinity

	DamageTakenMultiplier       float64            // All damage
	SchoolDamageTakenMultiplier [SchoolLen]float64 // For specific spell schools (arcane, fire, shadow, etc.)

	DiseaseDamageTakenMultiplier          float64
	PeriodicPhysicalDamageTakenMultiplier float64

	ArmorMultiplier float64 // Major/minor/special multiplicative armor modifiers

	ReducedPhysicalHitTakenChance float64
	ReducedArcaneHitTakenChance   float64
	ReducedFireHitTakenChance     float64
	ReducedFrostHitTakenChance    float64
	ReducedNatureHitTakenChance   float64
	ReducedShadowHitTakenChance   float64

	HealingTakenMultiplier float64
}

func NewPseudoStats() PseudoStats {
	return PseudoStats{
		CostMultiplier: 1,

		CastSpeedMultiplier:   1,
		MeleeSpeedMultiplier:  1,
		RangedSpeedMultiplier: 1,
		SpiritRegenMultiplier: 1,

		ThreatMultiplier: 1,

		DamageDealtMultiplier:       1,
		SchoolDamageDealtMultiplier: NewSchoolFloatArray(),

		MeleeHasteRatingPerHastePercent: 32.79,

		BlockValueMultiplier: 1,

		DamageSpread: 0.3333,

		// Target effects.
		DamageTakenMultiplier:       1,
		SchoolDamageTakenMultiplier: NewSchoolFloatArray(),

		DiseaseDamageTakenMultiplier:          1,
		PeriodicPhysicalDamageTakenMultiplier: 1,

		ArmorMultiplier: 1,

		HealingTakenMultiplier: 1,
	}
}

type UnitStat int

func (s UnitStat) IsStat() bool                                 { return int(s) < int(Len) }
func (s UnitStat) IsPseudoStat() bool                           { return !s.IsStat() }
func (s UnitStat) EqualsStat(other Stat) bool                   { return int(s) == int(other) }
func (s UnitStat) EqualsPseudoStat(other proto.PseudoStat) bool { return int(s) == int(other) }
func (s UnitStat) StatIdx() int {
	if !s.IsStat() {
		panic("Is a pseudo stat")
	}
	return int(s)
}
func (s UnitStat) PseudoStatIdx() int {
	if s.IsStat() {
		panic("Is a regular stat")
	}
	return int(s) - int(Len)
}
func (s UnitStat) AddToStatsProto(p *proto.UnitStats, value float64) {
	if s.IsStat() {
		p.Stats[s.StatIdx()] += value
	} else {
		p.PseudoStats[s.PseudoStatIdx()] += value
	}
}

func UnitStatFromIdx(s int) UnitStat                     { return UnitStat(s) }
func UnitStatFromStat(s Stat) UnitStat                   { return UnitStat(s) }
func UnitStatFromPseudoStat(s proto.PseudoStat) UnitStat { return UnitStat(int(s) + int(Len)) }
