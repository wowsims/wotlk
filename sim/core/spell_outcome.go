package core

import (
	"github.com/wowsims/wotlk/sim/core/stats"
)

// This function should do 3 things:
//  1. Set the Outcome of the hit effect.
//  2. Update spell outcome metrics.
//  3. Modify the damage if necessary.
type NewOutcomeApplier func(sim *Simulation, result *SpellEffect, attackTable *AttackTable)

func (spell *Spell) OutcomeAlwaysHit(_ *Simulation, result *SpellEffect, _ *AttackTable) {
	result.Outcome = OutcomeHit
	spell.SpellMetrics[result.Target.UnitIndex].Hits++
}
func (spell *Spell) OutcomeAlwaysMiss(_ *Simulation, result *SpellEffect, _ *AttackTable) {
	result.Outcome = OutcomeMiss
	result.Damage = 0
	spell.SpellMetrics[result.Target.UnitIndex].Misses++
}

// A tick always hits, but we don't count them as hits in the metrics.
func (dot *Dot) OutcomeTick(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	result.Outcome = OutcomeHit
}

func (dot *Dot) OutcomeTickPhysicalCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	if dot.Spell.PhysicalCritCheck(sim, result.Target, attackTable) {
		result.Outcome = OutcomeCrit
		result.Damage *= dot.Spell.CritMultiplier
	} else {
		result.Outcome = OutcomeHit
	}
}

func (dot *Dot) OutcomeSnapshotCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	if dot.Spell.CritMultiplier == 0 {
		panic("Spell " + dot.Spell.ActionID.String() + " missing CritMultiplier")
	}
	if sim.RandomFloat("Magical Crit Roll") < dot.SnapshotCritChance {
		result.Outcome = OutcomeCrit
		result.Damage *= dot.Spell.CritMultiplier
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Crits++
	} else {
		result.Outcome = OutcomeHit
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Hits++
	}
}

// TODO: Delete this, it's identical to OutcomeSnapshotCrit except it uses a different RNG string to preserve test results.
func (dot *Dot) OutcomeSnapshotCritPhysical(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	if dot.Spell.CritMultiplier == 0 {
		panic("Spell " + dot.Spell.ActionID.String() + " missing CritMultiplier")
	}
	if sim.RandomFloat("Physical Crit Roll") < dot.SnapshotCritChance {
		result.Outcome = OutcomeCrit
		result.Damage *= dot.Spell.CritMultiplier
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Crits++
	} else {
		result.Outcome = OutcomeHit
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Hits++
	}
}

// TODO: Remove this
func (dot *Dot) OutcomeTickSnapshotCritPhysical(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	if dot.Spell.CritMultiplier == 0 {
		panic("Spell " + dot.Spell.ActionID.String() + " missing CritMultiplier")
	}

	roll := sim.RandomFloat("Physical Tick Hit")
	chance := 0.0
	missChance := attackTable.BaseMissChance - dot.Spell.PhysicalHitChance(result.Target)
	chance = MaxFloat(0, missChance)
	if roll < chance {
		result.Outcome = OutcomeHit
	} else {
		if sim.RandomFloat("Physical Crit Roll") < dot.SnapshotCritChance {
			result.Outcome = OutcomeCrit
			result.Damage *= dot.Spell.CritMultiplier
		} else {
			result.Outcome = OutcomeHit
		}
	}
}

func (dot *Dot) OutcomeMagicHitAndSnapshotCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	if dot.Spell.CritMultiplier == 0 {
		panic("Spell " + dot.Spell.ActionID.String() + " missing CritMultiplier")
	}
	if dot.Spell.MagicHitCheck(sim, attackTable) {
		if sim.RandomFloat("Magical Crit Roll") < dot.SnapshotCritChance {
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

func (spell *Spell) OutcomeMagicHitAndCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
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

func (spell *Spell) OutcomeMagicCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
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
func (spell *Spell) CalcAndDealDamageMagicCrit(sim *Simulation, target *Unit, baseDamage float64) {
	result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
	spell.DealDamage(sim, result)
}

func (spell *Spell) OutcomeHealingCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
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

func (spell *Spell) OutcomeCritFixedChance(critChance float64) NewOutcomeApplier {
	return func(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
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

func (spell *Spell) OutcomeTickMagicHit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	if spell.MagicHitCheck(sim, attackTable) {
		result.Outcome = OutcomeHit
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
	}
}
func (spell *Spell) OutcomeMagicHit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	if spell.MagicHitCheck(sim, attackTable) {
		result.Outcome = OutcomeHit
		spell.SpellMetrics[result.Target.UnitIndex].Hits++
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
		spell.SpellMetrics[result.Target.UnitIndex].Misses++
	}
}
func (spell *Spell) CalcAndDealDamageMagicHit(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMagicHit)
	spell.DealDamage(sim, result)
}

func (spell *Spell) OutcomeMeleeWhite(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
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
func (spell *Spell) CalcAndDealDamageMeleeWhite(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMeleeWhite)
	spell.DealDamage(sim, result)
}

func (spell *Spell) OutcomeMeleeSpecialHit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
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
func (spell *Spell) CalcAndDealDamageMeleeSpecialHit(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMeleeSpecialHit)
	spell.DealDamage(sim, result)
}

func (spell *Spell) OutcomeMeleeSpecialHitAndCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
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
func (spell *Spell) CalcAndDealDamageMeleeSpecialHitAndCrit(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMeleeSpecialHitAndCrit)
	spell.DealDamage(sim, result)
}

// Like OutcomeMeleeSpecialHitAndCrit, but blocks prevent crits (all weapon damage based attacks).
func (spell *Spell) OutcomeMeleeWeaponSpecialHitAndCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
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
func (spell *Spell) CalcAndDealDamageMeleeWeaponSpecialHitAndCrit(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
	spell.DealDamage(sim, result)
}

func (spell *Spell) OutcomeMeleeWeaponSpecialNoCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
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
func (spell *Spell) CalcAndDealDamageMeleeWeaponSpecialNoCrit(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMeleeWeaponSpecialNoCrit)
	spell.DealDamage(sim, result)
}

func (spell *Spell) OutcomeMeleeSpecialNoBlockDodgeParry(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
		!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
		result.applyAttackTableHit(spell)
	}
}
func (spell *Spell) CalcAndDealDamageMeleeSpecialNoBlockDodgeParry(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
	spell.DealDamage(sim, result)
}

func (spell *Spell) OutcomeMeleeSpecialNoBlockDodgeParryNoCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) {
		result.applyAttackTableHit(spell)
	}
}

func (spell *Spell) OutcomeMeleeSpecialCritOnly(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	if !result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
		result.applyAttackTableHit(spell)
	}
}
func (spell *Spell) CalcAndDealDamageMeleeSpecialCritOnly(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMeleeSpecialCritOnly)
	spell.DealDamage(sim, result)
}

func (spell *Spell) OutcomeRangedHit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) {
		result.applyAttackTableHit(spell)
	}
}
func (spell *Spell) CalcAndDealDamageRangedHit(sim *Simulation, target *Unit, baseDamage float64) {
	result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHit)
	spell.DealDamage(sim, result)
}

func (spell *Spell) OutcomeRangedHitAndCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
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
func (dot *Dot) OutcomeRangedHitAndCritSnapshot(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
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
func (spell *Spell) CalcAndDealDamageRangedHitAndCrit(sim *Simulation, target *Unit, baseDamage float64) {
	result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
	spell.DealDamage(sim, result)
}

func (spell *Spell) OutcomeRangedCritOnly(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
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
func (spell *Spell) CalcAndDealDamageRangedCritOnly(sim *Simulation, target *Unit, baseDamage float64) {
	result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedCritOnly)
	spell.DealDamage(sim, result)
}

func (spell *Spell) OutcomeEnemyMeleeWhite(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
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

func (spellEffect *SpellEffect) applyAttackTableMiss(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	missChance := attackTable.BaseMissChance - spell.PhysicalHitChance(spellEffect.Target)
	if spell.Unit.AutoAttacks.IsDualWielding && !spell.Unit.PseudoStats.DisableDWMissPenalty {
		missChance += 0.19
	}
	*chance = MaxFloat(0, missChance)

	if roll < *chance {
		spellEffect.Outcome = OutcomeMiss
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Misses++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableMissNoDWPenalty(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	missChance := attackTable.BaseMissChance - spell.PhysicalHitChance(spellEffect.Target)
	*chance = MaxFloat(0, missChance)

	if roll < *chance {
		spellEffect.Outcome = OutcomeMiss
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Misses++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableBlock(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += attackTable.BaseBlockChance

	if roll < *chance {
		spellEffect.Outcome |= OutcomeBlock
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Blocks++
		spellEffect.Damage = MaxFloat(0, spellEffect.Damage-spellEffect.Target.GetStat(stats.BlockValue))
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableDodge(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	if spell.Flags.Matches(SpellFlagCannotBeDodged) {
		return false
	}

	*chance += MaxFloat(0, attackTable.BaseDodgeChance-spell.ExpertisePercentage()-spell.Unit.PseudoStats.DodgeReduction)

	if roll < *chance {
		spellEffect.Outcome = OutcomeDodge
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Dodges++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableParry(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += MaxFloat(0, attackTable.BaseParryChance-spell.ExpertisePercentage())

	if roll < *chance {
		spellEffect.Outcome = OutcomeParry
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Parries++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableGlance(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += attackTable.BaseGlanceChance

	if roll < *chance {
		spellEffect.Outcome = OutcomeGlance
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Glances++
		// TODO glancing blow damage reduction is actually a range ([65%, 85%] vs. +3, [80%, 90%] vs. +2, [91%, 99%] vs. +1 and +0)
		spellEffect.Damage *= attackTable.GlanceMultiplier
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableCrit(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	*chance += spell.PhysicalCritChance(spellEffect.Target, attackTable)

	if roll < *chance {
		spellEffect.Outcome = OutcomeCrit
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Crits++
		spellEffect.Damage *= spell.CritMultiplier
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableCritSeparateRoll(sim *Simulation, spell *Spell, attackTable *AttackTable) bool {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	if spell.PhysicalCritCheck(sim, spellEffect.Target, attackTable) {
		spellEffect.Outcome = OutcomeCrit
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Crits++
		spellEffect.Damage *= spell.CritMultiplier
		return true
	}
	return false
}
func (result *SpellEffect) applyAttackTableCritSeparateRollSnapshot(sim *Simulation, dot *Dot) bool {
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

func (spellEffect *SpellEffect) applyAttackTableHit(spell *Spell) {
	spellEffect.Outcome = OutcomeHit
	spell.SpellMetrics[spellEffect.Target.UnitIndex].Hits++
}

func (spellEffect *SpellEffect) applyEnemyAttackTableMiss(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	missChance := attackTable.BaseMissChance + spell.Unit.PseudoStats.IncreasedMissChance + spellEffect.Target.GetDiminishedMissChance()
	if spell.Unit.AutoAttacks.IsDualWielding && !spell.Unit.PseudoStats.DisableDWMissPenalty {
		missChance += 0.19
	}
	*chance = MaxFloat(0, missChance)

	if roll < *chance {
		spellEffect.Outcome = OutcomeMiss
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Misses++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyEnemyAttackTableBlock(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	if !spellEffect.Target.PseudoStats.CanBlock {
		return false
	}

	blockChance := attackTable.BaseBlockChance +
		spellEffect.Target.stats[stats.Block]/BlockRatingPerBlockChance/100 +
		spellEffect.Target.stats[stats.Defense]*DefenseRatingToChanceReduction
	*chance += MaxFloat(0, blockChance)

	if roll < *chance {
		spellEffect.Outcome |= OutcomeBlock
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Blocks++
		spellEffect.Damage = MaxFloat(0, spellEffect.Damage-spellEffect.Target.GetStat(stats.BlockValue))
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyEnemyAttackTableDodge(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	dodgeChance := attackTable.BaseDodgeChance +
		spellEffect.Target.PseudoStats.BaseDodge +
		spellEffect.Target.GetDiminishedDodgeChance() -
		spell.Unit.PseudoStats.DodgeReduction
	*chance += MaxFloat(0, dodgeChance)

	if roll < *chance {
		spellEffect.Outcome = OutcomeDodge
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Dodges++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyEnemyAttackTableParry(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	if !spellEffect.Target.PseudoStats.CanParry {
		return false
	}

	parryChance := attackTable.BaseParryChance +
		spellEffect.Target.PseudoStats.BaseParry +
		spellEffect.Target.GetDiminishedParryChance()
	*chance += MaxFloat(0, parryChance)

	if roll < *chance {
		spellEffect.Outcome = OutcomeParry
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Parries++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyEnemyAttackTableCrit(spell *Spell, _ *AttackTable, roll float64, chance *float64) bool {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	critRating := spell.Unit.stats[stats.MeleeCrit] + spell.BonusCritRating
	critChance := critRating / (CritRatingPerCritChance * 100)
	critChance -= spellEffect.Target.stats[stats.Defense] * DefenseRatingToChanceReduction
	critChance -= spellEffect.Target.stats[stats.Resilience] / ResilienceRatingPerCritReductionChance / 100
	critChance -= spellEffect.Target.PseudoStats.ReducedCritTakenChance
	*chance += MaxFloat(0, critChance)

	if roll < *chance {
		spellEffect.Outcome = OutcomeCrit
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Crits++
		resilCritMultiplier := 1 - spellEffect.Target.stats[stats.Resilience]/ResilienceRatingPerCritDamageReductionPercent/100
		spellEffect.Damage *= 2 * resilCritMultiplier
		return true
	}
	return false
}
