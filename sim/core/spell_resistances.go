package core

import (
	"fmt"
	"math"
	"strings"

	"github.com/wowsims/wotlk/sim/core/stats"
)

func (result *SpellResult) applyResistances(sim *Simulation, spell *Spell, isPeriodic bool, attackTable *AttackTable) {
	// TODO check why result.Outcome isn't updated with resists anymore
	resistanceMultiplier := spell.ResistanceMultiplier(sim, isPeriodic, attackTable)
	result.Damage *= resistanceMultiplier

	result.ResistanceMultiplier = resistanceMultiplier
	result.PreOutcomeDamage = result.Damage
}

// Modifies damage based on Armor or Magic resistances, depending on the damage type.
func (spell *Spell) ResistanceMultiplier(sim *Simulation, isPeriodic bool, attackTable *AttackTable) float64 {
	if spell.Flags.Matches(SpellFlagIgnoreResists) {
		return 1
	}

	if spell.SpellSchool.Matches(SpellSchoolPhysical) {
		// All physical dots (Bleeds) ignore armor.
		if isPeriodic && !spell.Flags.Matches(SpellFlagApplyArmorReduction) {
			return 1
		}

		// Physical resistance (armor).
		return attackTable.GetArmorDamageModifier(spell)
	}

	// Magical resistance.
	averageResist := attackTable.Defender.averageResist(spell.SpellSchool, attackTable.Attacker)
	if averageResist == 0 { // for equal or lower level mobs
		return 1
	}

	if spell.Flags.Matches(SpellFlagBinary) {
		if resistanceRoll := sim.RandomFloat("Binary Resist"); resistanceRoll < averageResist {
			return 0
		}
		return 1
	}

	thresholds := attackTable.Defender.partialResistRollThresholds(averageResist)

	switch resistanceRoll := sim.RandomFloat("Partial Resist"); {
	case resistanceRoll < thresholds[0].cumulativeChance:
		return thresholds[0].damageMultiplier()
	case resistanceRoll < thresholds[1].cumulativeChance:
		return thresholds[1].damageMultiplier()
	case resistanceRoll < thresholds[2].cumulativeChance:
		return thresholds[2].damageMultiplier()
	default:
		return thresholds[3].damageMultiplier()
	}
}

func (at *AttackTable) GetArmorDamageModifier(spell *Spell) float64 {
	armorConstant := float64(at.Attacker.Level)*467.5 - 22167.5
	defenderArmor := at.Defender.Armor()
	reducibleArmor := min((defenderArmor+armorConstant)/3, defenderArmor)
	armorPenRating := at.Attacker.stats[stats.ArmorPenetration] + spell.BonusArmorPenRating
	effectiveArmor := defenderArmor - reducibleArmor*at.Attacker.ArmorPenetrationPercentage(armorPenRating)
	return 1 - effectiveArmor/(effectiveArmor+armorConstant)
}

/*
 The following calculations are based on
 https://web.archive.org/web/20110207221537/http://elitistjerks.com/f15/t44675-resistance_mechanics_wotlk/
 This handles the mob vs. player case
  - average resist is calculated as AR = R / (R + C), C is 400 for level 80 mobs, assumed 510 for level 83 mobs
  - actual resist values come in multiples of 10%, with 3-4 values around the average resist
  - probability for a given resist value is P(x) = 0.5 - 2.5*|x - AR| (transformed for AR < 0.1 or AR > 0.9)
  - the resist cap is likely gone, since resists work like armor now
 https://web.archive.org/web/20110209210726/http://elitistjerks.com/f75/t38540-general_mage_discussion_information/p11/#post1171056
 This handles the player vs. mob partial resists case
  - average resist is still 2% percent per level vs. higher level mobs
  - otherwise it's modelled identical to the mob vs. player case
  - the resulting numbers have been verified in game (55% for 0%, 30% for 10%, 15% for 20% resists)
*/

func (unit *Unit) averageResist(school SpellSchool, attacker *Unit) float64 {
	resistance := unit.GetStat(school.ResistanceStat()) - attacker.stats[stats.SpellPenetration]
	if resistance <= 0 {
		return unit.levelBasedResist(attacker)
	}

	c := 5 * float64(attacker.Level)
	if attacker.Type == EnemyUnit && attacker.Level-unit.Level >= 3 {
		c = 510 // other values TBD, but not very useful in practice
	}

	return resistance/(c+resistance) + unit.levelBasedResist(attacker) // these may stack differently, but that's irrelevant in practice
}

func (unit *Unit) levelBasedResist(attacker *Unit) float64 {
	if unit.Type == EnemyUnit && unit.Level > attacker.Level {
		return 0.02 * float64(unit.Level-attacker.Level)
	}
	return 0
}

type Threshold struct {
	cumulativeChance float64
	bracket          int
}

func (x Threshold) damageMultiplier() float64 {
	return 1 - 0.1*float64(x.bracket)
}

type Thresholds [4]Threshold

func (x Thresholds) String() string {
	var sb strings.Builder
	var chance float64
	for _, t := range x {
		sb.WriteString(fmt.Sprintf("%.1f%% for %d%% ", (t.cumulativeChance-chance)*100, t.bracket*10))
		if t.cumulativeChance >= 1 {
			break
		}
		chance = t.cumulativeChance
	}
	return sb.String()
}

func (unit *Unit) partialResistRollThresholds(ar float64) Thresholds {
	if ar <= 0.1 { // always 0%, 10%, or 20%; this covers all player vs. mob cases, in practice
		return Thresholds{
			{cumulativeChance: 1 - 7.5*ar, bracket: 0},
			{cumulativeChance: 1 - 2.5*ar, bracket: 1},
			{cumulativeChance: 1, bracket: 2},
		}
	}

	if ar >= 0.9 { // always 80%, 90%, or 100%; only relevant for tests ;)
		return Thresholds{
			{cumulativeChance: 1 - 7.5*(1-ar), bracket: 10},
			{cumulativeChance: 1 - 2.5*(1-ar), bracket: 9},
			{cumulativeChance: 1, bracket: 8},
		}
	}

	p := func(x float64) float64 {
		return math.Max(0.5-2.5*math.Abs(x-ar), 0)
	}

	const eps = 1e-9 // imprecision guard (25-50-25 might become almost0-25-50-25-almost0)

	var thresholds Thresholds
	var cumulativeChance float64
	var index int
	for bracket := 0; bracket <= 10; bracket++ {
		if chance := p(float64(bracket) * 0.1); chance > eps {
			cumulativeChance += chance
			thresholds[index] = Threshold{cumulativeChance: cumulativeChance, bracket: bracket}
			index++
		}
	}

	if thresholds[index-1].cumulativeChance < 1 { // also guards against floating point imprecision
		thresholds[index-1].cumulativeChance = 1
	}

	return thresholds
}
