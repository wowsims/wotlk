package core

import (
	"github.com/wowsims/wotlk/sim/core/stats"
)

// This function should do 3 things:
//  1. Set the Outcome of the hit effect.
//  2. Update spell outcome metrics.
//  3. Modify the damage if necessary.
type OutcomeApplier func(sim *Simulation, result *SpellResult, attackTable *AttackTable)

func (spell *Spell) OutcomeAlwaysHit(_ *Simulation, result *SpellResult, _ *AttackTable) {
	result.Outcome = OutcomeHit
	spell.SpellMetrics[result.Target.UnitIndex].Hits++
}
func (spell *Spell) OutcomeAlwaysMiss(_ *Simulation, result *SpellResult, _ *AttackTable) {
	result.Outcome = OutcomeMiss
	result.Damage = 0
	spell.SpellMetrics[result.Target.UnitIndex].Misses++
}

// A tick always hits, but we don't count them as hits in the metrics.
func (dot *Dot) OutcomeTick(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	result.Outcome = OutcomeHit
}

func (dot *Dot) OutcomeTickCounted(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	result.Outcome = OutcomeHit
	dot.Spell.SpellMetrics[result.Target.UnitIndex].Hits++
}

func (dot *Dot) OutcomeTickPhysicalCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if dot.Spell.PhysicalCritCheck(sim, result.Target, attackTable) {
		result.Outcome = OutcomeCrit
		result.Damage *= dot.Spell.CritMultiplier
	} else {
		result.Outcome = OutcomeHit
	}
}

func (dot *Dot) OutcomeSnapshotCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if dot.Spell.CritMultiplier == 0 {
		panic("Spell " + dot.Spell.ActionID.String() + " missing CritMultiplier")
	}
	if sim.RandomFloat("Snapshot Crit Roll") < dot.SnapshotCritChance {
		result.Outcome = OutcomeCrit
		result.Damage *= dot.Spell.CritMultiplier
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Crits++
	} else {
		result.Outcome = OutcomeHit
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Hits++
	}
}

func (dot *Dot) OutcomeMagicHitAndSnapshotCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if dot.Spell.CritMultiplier == 0 {
		panic("Spell " + dot.Spell.ActionID.String() + " missing CritMultiplier")
	}
	if dot.Spell.MagicHitCheck(sim, attackTable) {
		if sim.RandomFloat("Snapshot Crit Roll") < dot.SnapshotCritChance {
			result.Outcome = OutcomeCrit
			result.Damage *= dot.Spell.CritMultiplier
			dot.Spell.SpellMetrics[result.Target.UnitIndex].Crits++
		} else {
			result.Outcome = OutcomeHit
			dot.Spell.SpellMetrics[result.Target.UnitIndex].Hits++
		}
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Misses++
	}
}

func (spell *Spell) OutcomeMagicHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	if spell.MagicHitCheck(sim, attackTable) {
		if spell.MagicCritCheck(sim, result.Target) {
			result.Outcome = OutcomeCrit
			result.Damage *= spell.CritMultiplier
			spell.SpellMetrics[result.Target.UnitIndex].Crits++
		} else {
			result.Outcome = OutcomeHit
			spell.SpellMetrics[result.Target.UnitIndex].Hits++
		}
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
		spell.SpellMetrics[result.Target.UnitIndex].Misses++
	}
}

func (spell *Spell) OutcomeMagicCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	if spell.MagicCritCheck(sim, result.Target) {
		result.Outcome = OutcomeCrit
		result.Damage *= spell.CritMultiplier
		spell.SpellMetrics[result.Target.UnitIndex].Crits++
	} else {
		result.Outcome = OutcomeHit
		spell.SpellMetrics[result.Target.UnitIndex].Hits++
	}
}

func (spell *Spell) OutcomeHealingCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	if spell.HealingCritCheck(sim) {
		result.Outcome = OutcomeCrit
		result.Damage *= spell.CritMultiplier
		spell.SpellMetrics[result.Target.UnitIndex].Crits++
	} else {
		result.Outcome = OutcomeHit
		spell.SpellMetrics[result.Target.UnitIndex].Hits++
	}
}

func (spell *Spell) OutcomeCritFixedChance(critChance float64) OutcomeApplier {
	return func(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
		if spell.CritMultiplier == 0 {
			panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
		}
		if spell.fixedCritCheck(sim, critChance) {
			result.Outcome = OutcomeCrit
			result.Damage *= spell.CritMultiplier
			spell.SpellMetrics[result.Target.UnitIndex].Crits++
		} else {
			result.Outcome = OutcomeHit
			spell.SpellMetrics[result.Target.UnitIndex].Hits++
		}
	}
}

func (spell *Spell) OutcomeMagicCritFixedChance(critChance float64) OutcomeApplier {
	return func(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
		if spell.CritMultiplier == 0 {
			panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
		}
		if spell.MagicHitCheck(sim, attackTable) {
			if spell.fixedCritCheck(sim, critChance) {
				result.Outcome = OutcomeCrit
				result.Damage *= spell.CritMultiplier
				spell.SpellMetrics[result.Target.UnitIndex].Crits++
			} else {
				result.Outcome = OutcomeHit
				spell.SpellMetrics[result.Target.UnitIndex].Hits++
			}
		} else {
			result.Outcome = OutcomeMiss
			result.Damage = 0
			spell.SpellMetrics[result.Target.UnitIndex].Misses++
		}
	}
}

func (spell *Spell) OutcomeTickMagicHit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if spell.MagicHitCheck(sim, attackTable) {
		result.Outcome = OutcomeHit
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
	}
}
func (spell *Spell) OutcomeMagicHit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if spell.MagicHitCheck(sim, attackTable) {
		result.Outcome = OutcomeHit
		spell.SpellMetrics[result.Target.UnitIndex].Hits++
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
		spell.SpellMetrics[result.Target.UnitIndex].Misses++
	}
}

func (spell *Spell) OutcomeMeleeWhite(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMiss(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableGlance(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableBlock(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableCrit(spell, attackTable, roll, &chance) {
			result.applyAttackTableHit(spell)
		}
	} else {
		if !result.applyAttackTableMiss(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableGlance(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableCrit(spell, attackTable, roll, &chance) {
			result.applyAttackTableHit(spell)
		}
	}
}

func (spell *Spell) OutcomeMeleeSpecialHit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance) {
			result.applyAttackTableHit(spell)
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) {
			result.applyAttackTableHit(spell)
		}
	}
}

func (spell *Spell) OutcomeMeleeSpecialHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance) {
			if result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
				result.applyAttackTableBlock(spell, attackTable, roll, &chance)
			} else {
				if !result.applyAttackTableBlock(spell, attackTable, roll, &chance) {
					result.applyAttackTableHit(spell)
				}
			}
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
			result.applyAttackTableHit(spell)
		}
	}
}

// Like OutcomeMeleeSpecialHitAndCrit, but blocks prevent crits (all weapon damage based attacks).
func (spell *Spell) OutcomeMeleeWeaponSpecialHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if spell.Unit.PseudoStats.InFrontOfTarget {
		roll := sim.RandomFloat("White Hit Table")
		chance := 0.0

		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableBlock(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
			result.applyAttackTableHit(spell)
		}
	} else {
		spell.OutcomeMeleeSpecialHitAndCrit(sim, result, attackTable)
	}
}

func (spell *Spell) OutcomeMeleeWeaponSpecialNoCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableBlock(spell, attackTable, roll, &chance) {
			result.applyAttackTableHit(spell)
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) {
			result.applyAttackTableHit(spell)
		}
	}
}

func (spell *Spell) OutcomeMeleeSpecialNoBlockDodgeParry(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
		!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
		result.applyAttackTableHit(spell)
	}
}

func (spell *Spell) OutcomeMeleeSpecialNoBlockDodgeParryNoCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) {
		result.applyAttackTableHit(spell)
	}
}

func (spell *Spell) OutcomeMeleeSpecialCritOnly(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if !result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
		result.applyAttackTableHit(spell)
	}
}

func (spell *Spell) OutcomeRangedHit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) {
		result.applyAttackTableHit(spell)
	}
}

func (spell *Spell) OutcomeRangedHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if spell.Unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) {
			if result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
				result.applyAttackTableBlock(spell, attackTable, roll, &chance)
			} else {
				if !result.applyAttackTableBlock(spell, attackTable, roll, &chance) {
					result.applyAttackTableHit(spell)
				}
			}
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
			result.applyAttackTableHit(spell)
		}
	}
}
func (dot *Dot) OutcomeRangedHitAndCritSnapshot(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if dot.Spell.Unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(dot.Spell, attackTable, roll, &chance) {
			if result.applyAttackTableCritSeparateRollSnapshot(sim, dot) {
				result.applyAttackTableBlock(dot.Spell, attackTable, roll, &chance)
			} else {
				if !result.applyAttackTableBlock(dot.Spell, attackTable, roll, &chance) {
					result.applyAttackTableHit(dot.Spell)
				}
			}
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(dot.Spell, attackTable, roll, &chance) &&
			!result.applyAttackTableCritSeparateRollSnapshot(sim, dot) {
			result.applyAttackTableHit(dot.Spell)
		}
	}
}

func (spell *Spell) OutcomeRangedCritOnly(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	// Block already checks for this, but we can skip the RNG roll which is expensive.
	if spell.Unit.PseudoStats.InFrontOfTarget {
		roll := sim.RandomFloat("White Hit Table")
		chance := 0.0

		if result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
			result.applyAttackTableBlock(spell, attackTable, roll, &chance)
		} else {
			if !result.applyAttackTableBlock(spell, attackTable, roll, &chance) {
				result.applyAttackTableHit(spell)
			}
		}
	} else {
		if !result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
			result.applyAttackTableHit(spell)
		}
	}
}

func (spell *Spell) OutcomeEnemyMeleeWhite(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	roll := sim.RandomFloat("Enemy White Hit Table")
	chance := 0.0

	if !result.applyEnemyAttackTableMiss(spell, attackTable, roll, &chance) &&
		!result.applyEnemyAttackTableDodge(spell, attackTable, roll, &chance) &&
		!result.applyEnemyAttackTableParry(spell, attackTable, roll, &chance) &&
		!result.applyEnemyAttackTableBlock(spell, attackTable, roll, &chance) &&
		!result.applyEnemyAttackTableCrit(spell, attackTable, roll, &chance) {
		result.applyAttackTableHit(spell)
	}
}

func (spell *Spell) fixedCritCheck(sim *Simulation, critChance float64) bool {
	return sim.RandomFloat("Fixed Crit Roll") < critChance
}

func (result *SpellResult) applyAttackTableMiss(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	missChance := attackTable.BaseMissChance - spell.PhysicalHitChance(result.Target)
	if spell.Unit.AutoAttacks.IsDualWielding && !spell.Unit.PseudoStats.DisableDWMissPenalty {
		missChance += 0.19
	}
	*chance = MaxFloat(0, missChance)

	if roll < *chance {
		result.Outcome = OutcomeMiss
		spell.SpellMetrics[result.Target.UnitIndex].Misses++
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableMissNoDWPenalty(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	missChance := attackTable.BaseMissChance - spell.PhysicalHitChance(result.Target)
	*chance = MaxFloat(0, missChance)

	if roll < *chance {
		result.Outcome = OutcomeMiss
		spell.SpellMetrics[result.Target.UnitIndex].Misses++
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableBlock(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += attackTable.BaseBlockChance

	if roll < *chance {
		result.Outcome |= OutcomeBlock
		spell.SpellMetrics[result.Target.UnitIndex].Blocks++
		result.Damage = MaxFloat(0, result.Damage-result.Target.BlockValue())
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableDodge(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	if spell.Flags.Matches(SpellFlagCannotBeDodged) {
		return false
	}

	*chance += MaxFloat(0, attackTable.BaseDodgeChance-spell.ExpertisePercentage()-spell.Unit.PseudoStats.DodgeReduction)

	if roll < *chance {
		result.Outcome = OutcomeDodge
		spell.SpellMetrics[result.Target.UnitIndex].Dodges++
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableParry(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += MaxFloat(0, attackTable.BaseParryChance-spell.ExpertisePercentage())

	if roll < *chance {
		result.Outcome = OutcomeParry
		spell.SpellMetrics[result.Target.UnitIndex].Parries++
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableGlance(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += attackTable.BaseGlanceChance

	if roll < *chance {
		result.Outcome = OutcomeGlance
		spell.SpellMetrics[result.Target.UnitIndex].Glances++
		// TODO glancing blow damage reduction is actually a range ([65%, 85%] vs. +3, [80%, 90%] vs. +2, [91%, 99%] vs. +1 and +0)
		result.Damage *= attackTable.GlanceMultiplier
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableCrit(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	*chance += spell.PhysicalCritChance(result.Target, attackTable)

	if roll < *chance {
		result.Outcome = OutcomeCrit
		spell.SpellMetrics[result.Target.UnitIndex].Crits++
		result.Damage *= spell.CritMultiplier
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableCritSeparateRoll(sim *Simulation, spell *Spell, attackTable *AttackTable) bool {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	if spell.PhysicalCritCheck(sim, result.Target, attackTable) {
		result.Outcome = OutcomeCrit
		spell.SpellMetrics[result.Target.UnitIndex].Crits++
		result.Damage *= spell.CritMultiplier
		return true
	}
	return false
}
func (result *SpellResult) applyAttackTableCritSeparateRollSnapshot(sim *Simulation, dot *Dot) bool {
	if dot.Spell.CritMultiplier == 0 {
		panic("Spell " + dot.Spell.ActionID.String() + " missing CritMultiplier")
	}
	if sim.RandomFloat("Physical Crit Roll") < dot.SnapshotCritChance {
		result.Outcome = OutcomeCrit
		result.Damage *= dot.Spell.CritMultiplier
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Crits++
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableHit(spell *Spell) {
	result.Outcome = OutcomeHit
	spell.SpellMetrics[result.Target.UnitIndex].Hits++
}

func (result *SpellResult) applyEnemyAttackTableMiss(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	missChance := attackTable.BaseMissChance + spell.Unit.PseudoStats.IncreasedMissChance + result.Target.GetDiminishedMissChance() + result.Target.PseudoStats.ReducedPhysicalHitTakenChance
	if spell.Unit.AutoAttacks.IsDualWielding && !spell.Unit.PseudoStats.DisableDWMissPenalty {
		missChance += 0.19
	}
	*chance = MaxFloat(0, missChance)

	if roll < *chance {
		result.Outcome = OutcomeMiss
		spell.SpellMetrics[result.Target.UnitIndex].Misses++
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyEnemyAttackTableBlock(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	if !result.Target.PseudoStats.CanBlock {
		return false
	}

	blockChance := attackTable.BaseBlockChance +
		result.Target.stats[stats.Block]/BlockRatingPerBlockChance/100 +
		result.Target.stats[stats.Defense]*DefenseRatingToChanceReduction
	*chance += MaxFloat(0, blockChance)

	if roll < *chance {
		result.Outcome |= OutcomeBlock
		spell.SpellMetrics[result.Target.UnitIndex].Blocks++
		result.Damage = MaxFloat(0, result.Damage-result.Target.BlockValue())
		return true
	}
	return false
}

func (result *SpellResult) applyEnemyAttackTableDodge(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	dodgeChance := attackTable.BaseDodgeChance +
		result.Target.PseudoStats.BaseDodge +
		result.Target.GetDiminishedDodgeChance() -
		spell.Unit.PseudoStats.DodgeReduction
	*chance += MaxFloat(0, dodgeChance)

	if roll < *chance {
		result.Outcome = OutcomeDodge
		spell.SpellMetrics[result.Target.UnitIndex].Dodges++
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyEnemyAttackTableParry(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	if !result.Target.PseudoStats.CanParry {
		return false
	}

	parryChance := attackTable.BaseParryChance +
		result.Target.PseudoStats.BaseParry +
		result.Target.GetDiminishedParryChance()
	*chance += MaxFloat(0, parryChance)

	if roll < *chance {
		result.Outcome = OutcomeParry
		spell.SpellMetrics[result.Target.UnitIndex].Parries++
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyEnemyAttackTableCrit(spell *Spell, _ *AttackTable, roll float64, chance *float64) bool {

	critRating := spell.Unit.stats[stats.MeleeCrit] + spell.BonusCritRating
	critChance := critRating / (CritRatingPerCritChance * 100)
	critChance -= result.Target.stats[stats.Defense] * DefenseRatingToChanceReduction
	critChance -= result.Target.stats[stats.Resilience] / ResilienceRatingPerCritReductionChance / 100
	critChance -= result.Target.PseudoStats.ReducedCritTakenChance
	*chance += MaxFloat(0, critChance)

	if roll < *chance {
		result.Outcome = OutcomeCrit
		spell.SpellMetrics[result.Target.UnitIndex].Crits++
		// Assume PvE enemies do not use damage reduction multiplier component in WotLK
		//resilCritMultiplier := 1 - result.Target.stats[stats.Resilience]/ResilienceRatingPerCritDamageReductionPercent/100
		result.Damage *= 2
		return true
	}
	return false
}

func (spell *Spell) OutcomeExpectedTick(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	// result.Damage *= 1
}
func (spell *Spell) OutcomeExpectedMagicAlwaysHit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	// result.Damage *= 1
}
func (spell *Spell) OutcomeExpectedMagicHit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	averageMultiplier := 1.0
	averageMultiplier -= spell.SpellChanceToMiss(attackTable)

	result.Damage *= averageMultiplier
}

func (spell *Spell) OutcomeExpectedMagicCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}

	averageMultiplier := 1.0
	averageMultiplier += spell.SpellCritChance(result.Target) * (spell.CritMultiplier - 1)

	result.Damage *= averageMultiplier
}

func (spell *Spell) OutcomeExpectedMagicHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}

	averageMultiplier := 1.0
	averageMultiplier -= spell.SpellChanceToMiss(attackTable)
	averageMultiplier += averageMultiplier * spell.SpellCritChance(result.Target) * (spell.CritMultiplier - 1)

	result.Damage *= averageMultiplier
}

func (dot *Dot) OutcomeExpectedMagicSnapshotCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if dot.Spell.CritMultiplier == 0 {
		panic("Spell " + dot.Spell.ActionID.String() + " missing CritMultiplier")
	}

	averageMultiplier := 1.0
	averageMultiplier += dot.SnapshotCritChance * (dot.Spell.CritMultiplier - 1)

	result.Damage *= averageMultiplier
}
