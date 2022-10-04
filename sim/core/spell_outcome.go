package core

import (
	"github.com/wowsims/wotlk/sim/core/stats"
)

// This function should do 3 things:
//  1. Set the Outcome of the hit effect.
//  2. Update spell outcome metrics.
//  3. Modify the damage if necessary.
type OutcomeApplier func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable)
type NewOutcomeApplier func(sim *Simulation, result *SpellEffect, attackTable *AttackTable)

func (spell *Spell) OutcomeAlwaysHit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	result.Outcome = OutcomeHit
	spell.SpellMetrics[result.Target.UnitIndex].Hits++
}
func (spell *Spell) CalcAndDealDamageAlwaysHit(sim *Simulation, target *Unit, baseDamage float64) {
	result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
	spell.DealDamage(sim, &result)
}
func (unit *Unit) OutcomeFuncAlwaysHit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spell.OutcomeAlwaysHit(sim, spellEffect, attackTable)
	}
}

// A tick always hits, but we don't count them as hits in the metrics.
func (unit *Unit) OutcomeFuncTick() OutcomeApplier {
	return func(_ *Simulation, _ *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spellEffect.Outcome = OutcomeHit
	}
}

func (unit *Unit) OutcomeFuncTickHitAndCrit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if spell.CritMultiplier == 0 {
			panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
		}
		roll := sim.RandomFloat("Physical Tick Hit")
		chance := 0.0
		missChance := attackTable.BaseMissChance - spell.PhysicalHitChance(spellEffect.Target)
		chance = MaxFloat(0, missChance)
		if roll < chance {
			spellEffect.Outcome = OutcomeHit
		} else {
			if spellEffect.physicalCritRoll(sim, spell, attackTable) {
				spellEffect.Outcome = OutcomeCrit
				spellEffect.Damage *= spell.CritMultiplier
			} else {
				spellEffect.Outcome = OutcomeHit
			}
		}
	}
}

func (unit *Unit) OutcomeFuncTickMagicHitAndCrit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if spell.CritMultiplier == 0 {
			panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
		}
		if spell.MagicHitCheck(sim, attackTable) {
			if spellEffect.MagicCritCheck(sim, spell, attackTable) {
				spellEffect.Outcome = OutcomeCrit
				spellEffect.Damage *= spell.CritMultiplier
			} else {
				spellEffect.Outcome = OutcomeHit
			}
		} else {
			spellEffect.Outcome = OutcomeHit
		}
	}
}

func (spell *Spell) OutcomeMagicHitAndCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	if spell.MagicHitCheck(sim, attackTable) {
		if result.MagicCritCheck(sim, spell, attackTable) {
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
func (spell *Spell) CalcAndDealDamageMagicHitAndCrit(sim *Simulation, target *Unit, baseDamage float64) {
	result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	spell.DealDamage(sim, &result)
}
func (unit *Unit) OutcomeFuncMagicHitAndCrit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spell.OutcomeMagicHitAndCrit(sim, spellEffect, attackTable)
	}
}

func (spell *Spell) OutcomeMagicCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	if result.MagicCritCheck(sim, spell, attackTable) {
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
	spell.DealDamage(sim, &result)
}
func (unit *Unit) OutcomeFuncMagicCrit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spell.OutcomeMagicCrit(sim, spellEffect, attackTable)
	}
}

func (spell *Spell) OutcomeMagicHitAndCritBinary(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	if spell.MagicHitCheckBinary(sim, attackTable) {
		if result.MagicCritCheck(sim, spell, attackTable) {
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
func (spell *Spell) CalcAndDealDamageMagicHitAndCritBinary(sim *Simulation, target *Unit, baseDamage float64) {
	result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCritBinary)
	spell.DealDamage(sim, &result)
}

func (spell *Spell) OutcomeHealingCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	if result.HealingCritCheck(sim, spell, attackTable) {
		result.Outcome = OutcomeCrit
		result.Damage *= spell.CritMultiplier
		spell.SpellMetrics[result.Target.UnitIndex].Crits++
	} else {
		result.Outcome = OutcomeHit
		spell.SpellMetrics[result.Target.UnitIndex].Hits++
	}
}
func (spell *Spell) CalcAndDealHealingCrit(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)
	spell.DealHealing(sim, &result)
}
func (unit *Unit) OutcomeFuncHealingCrit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spell.OutcomeHealingCrit(sim, spellEffect, attackTable)
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

func (unit *Unit) OutcomeFuncTickMagicHit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if spell.MagicHitCheck(sim, attackTable) {
			spellEffect.Outcome = OutcomeHit
		} else {
			spellEffect.Outcome = OutcomeMiss
			spellEffect.Damage = 0
		}
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
	spell.DealDamage(sim, &result)
}
func (unit *Unit) OutcomeFuncMagicHit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spell.OutcomeMagicHit(sim, spellEffect, attackTable)
	}
}

func (spell *Spell) OutcomeMagicHitBinary(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	if spell.MagicHitCheckBinary(sim, attackTable) {
		result.Outcome = OutcomeHit
		spell.SpellMetrics[result.Target.UnitIndex].Hits++
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
		spell.SpellMetrics[result.Target.UnitIndex].Misses++
	}
}
func (spell *Spell) CalcAndDealDamageMagicHitBinary(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMagicHitBinary)
	spell.DealDamage(sim, &result)
}
func (unit *Unit) OutcomeFuncMagicHitBinary() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spell.OutcomeMagicHitBinary(sim, spellEffect, attackTable)
	}
}

func (spell *Spell) OutcomeMeleeWhite(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMiss(spell, unit, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, unit, attackTable, roll, &chance) &&
			!result.applyAttackTableParry(spell, unit, attackTable, roll, &chance) &&
			!result.applyAttackTableGlance(spell, unit, attackTable, roll, &chance) &&
			!result.applyAttackTableBlock(spell, unit, attackTable, roll, &chance) &&
			!result.applyAttackTableCrit(spell, unit, attackTable, roll, &chance) {
			result.applyAttackTableHit(spell)
		}
	} else {
		if !result.applyAttackTableMiss(spell, unit, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, unit, attackTable, roll, &chance) &&
			!result.applyAttackTableGlance(spell, unit, attackTable, roll, &chance) &&
			!result.applyAttackTableCrit(spell, unit, attackTable, roll, &chance) {
			result.applyAttackTableHit(spell)
		}
	}
}
func (spell *Spell) CalcAndDealDamageMeleeWhite(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMeleeWhite)
	spell.DealDamage(sim, &result)
}

func (spell *Spell) OutcomeMeleeSpecialHit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
			(spell.Flags.Matches(SpellFlagCannotBeDodged) || !result.applyAttackTableDodge(spell, unit, attackTable, roll, &chance)) &&
			!result.applyAttackTableParry(spell, unit, attackTable, roll, &chance) {
			result.applyAttackTableHit(spell)
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
			(spell.Flags.Matches(SpellFlagCannotBeDodged) || !result.applyAttackTableDodge(spell, unit, attackTable, roll, &chance)) {
			result.applyAttackTableHit(spell)
		}
	}
}
func (spell *Spell) CalcAndDealDamageMeleeSpecialHit(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMeleeSpecialHit)
	spell.DealDamage(sim, &result)
}
func (unit *Unit) OutcomeFuncMeleeSpecialHit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spell.OutcomeMeleeSpecialHit(sim, spellEffect, attackTable)
	}
}

func (spell *Spell) OutcomeMeleeSpecialHitAndCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
			(spell.Flags.Matches(SpellFlagCannotBeDodged) || !result.applyAttackTableDodge(spell, unit, attackTable, roll, &chance)) &&
			!result.applyAttackTableParry(spell, unit, attackTable, roll, &chance) {
			if result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
				result.applyAttackTableBlock(spell, unit, attackTable, roll, &chance)
			} else {
				if !result.applyAttackTableBlock(spell, unit, attackTable, roll, &chance) {
					result.applyAttackTableHit(spell)
				}
			}
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
			(spell.Flags.Matches(SpellFlagCannotBeDodged) || !result.applyAttackTableDodge(spell, unit, attackTable, roll, &chance)) &&
			!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
			result.applyAttackTableHit(spell)
		}
	}
}
func (spell *Spell) CalcAndDealDamageMeleeSpecialHitAndCrit(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMeleeSpecialHitAndCrit)
	spell.DealDamage(sim, &result)
}
func (unit *Unit) OutcomeFuncMeleeSpecialHitAndCrit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spell.OutcomeMeleeSpecialHitAndCrit(sim, spellEffect, attackTable)
	}
}

// Like OutcomeFuncMeleeSpecialHitAndCrit, but blocks prevent crits.
func (spell *Spell) OutcomeMeleeWeaponSpecialHitAndCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	if spell.Unit.PseudoStats.InFrontOfTarget {
		unit := spell.Unit
		roll := sim.RandomFloat("White Hit Table")
		chance := 0.0

		if !result.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
			(spell.Flags.Matches(SpellFlagCannotBeDodged) || !result.applyAttackTableDodge(spell, unit, attackTable, roll, &chance)) &&
			!result.applyAttackTableParry(spell, unit, attackTable, roll, &chance) &&
			!result.applyAttackTableBlock(spell, unit, attackTable, roll, &chance) &&
			!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
			result.applyAttackTableHit(spell)
		}
	} else {
		spell.OutcomeMeleeSpecialHitAndCrit(sim, result, attackTable)
	}
}
func (spell *Spell) CalcAndDealDamageMeleeWeaponSpecialHitAndCrit(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
	spell.DealDamage(sim, &result)
}
func (unit *Unit) OutcomeFuncMeleeWeaponSpecialHitAndCrit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spell.OutcomeMeleeWeaponSpecialHitAndCrit(sim, spellEffect, attackTable)
	}
}

func (spell *Spell) OutcomeMeleeWeaponSpecialNoCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
			(spell.Flags.Matches(SpellFlagCannotBeDodged) || !result.applyAttackTableDodge(spell, unit, attackTable, roll, &chance)) &&
			!result.applyAttackTableParry(spell, unit, attackTable, roll, &chance) &&
			!result.applyAttackTableBlock(spell, unit, attackTable, roll, &chance) {
			result.applyAttackTableHit(spell)
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
			spell.Flags.Matches(SpellFlagCannotBeDodged) || !result.applyAttackTableDodge(spell, unit, attackTable, roll, &chance) {
			result.applyAttackTableHit(spell)
		}
	}
}
func (spell *Spell) CalcAndDealDamageMeleeWeaponSpecialNoCrit(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMeleeWeaponSpecialNoCrit)
	spell.DealDamage(sim, &result)
}
func (unit *Unit) OutcomeFuncMeleeWeaponSpecialNoCrit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spell.OutcomeMeleeWeaponSpecialNoCrit(sim, spellEffect, attackTable)
	}
}

func (spell *Spell) OutcomeMeleeSpecialNoBlockDodgeParry(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
		!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
		result.applyAttackTableHit(spell)
	}
}
func (spell *Spell) CalcAndDealDamageMeleeSpecialNoBlockDodgeParry(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
	spell.DealDamage(sim, &result)
}
func (unit *Unit) OutcomeFuncMeleeSpecialNoBlockDodgeParry() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spell.OutcomeMeleeSpecialNoBlockDodgeParry(sim, spellEffect, attackTable)
	}
}

func (spell *Spell) OutcomeMeleeSpecialCritOnly(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	if !result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
		result.applyAttackTableHit(spell)
	}
}
func (spell *Spell) CalcAndDealDamageMeleeSpecialCritOnly(sim *Simulation, target *Unit, baseHealing float64) {
	result := spell.CalcDamage(sim, target, baseHealing, spell.OutcomeMeleeSpecialCritOnly)
	spell.DealDamage(sim, &result)
}

//func (unit *Unit) OutcomeFuncMeleeSpecialCritOnly() OutcomeApplier {
//	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
//		spell.OutcomeMeleeSpecialCritOnly(sim, spellEffect, attackTable)
//	}
//}

func (spell *Spell) OutcomeRangedHit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) {
		result.applyAttackTableHit(spell)
	}
}
func (spell *Spell) CalcAndDealDamageRangedHit(sim *Simulation, target *Unit, baseDamage float64) {
	result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHit)
	spell.DealDamage(sim, &result)
}
func (unit *Unit) OutcomeFuncRangedHit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spell.OutcomeRangedHit(sim, spellEffect, attackTable)
	}
}

func (spell *Spell) OutcomeRangedHitAndCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if spell.Unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) {
			if result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
				result.applyAttackTableBlock(spell, unit, attackTable, roll, &chance)
			} else {
				if !result.applyAttackTableBlock(spell, unit, attackTable, roll, &chance) {
					result.applyAttackTableHit(spell)
				}
			}
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
			!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
			result.applyAttackTableHit(spell)
		}
	}
}
func (spell *Spell) CalcAndDealDamageRangedHitAndCrit(sim *Simulation, target *Unit, baseDamage float64) {
	result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
	spell.DealDamage(sim, &result)
}
func (unit *Unit) OutcomeFuncRangedHitAndCrit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spell.OutcomeRangedHitAndCrit(sim, spellEffect, attackTable)
	}
}

func (spell *Spell) OutcomeRangedCritOnly(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	// Block already checks for this, but we can skip the RNG roll which is expensive.
	if spell.Unit.PseudoStats.InFrontOfTarget {
		unit := spell.Unit
		roll := sim.RandomFloat("White Hit Table")
		chance := 0.0

		if result.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
			result.applyAttackTableBlock(spell, unit, attackTable, roll, &chance)
		} else {
			if !result.applyAttackTableBlock(spell, unit, attackTable, roll, &chance) {
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
	spell.DealDamage(sim, &result)
}

func (unit *Unit) OutcomeFuncEnemyMeleeWhite() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		unit := spell.Unit
		roll := sim.RandomFloat("Enemy White Hit Table")
		chance := 0.0

		if !spellEffect.applyEnemyAttackTableMiss(spell, unit, attackTable, roll, &chance) &&
			!spellEffect.applyEnemyAttackTableDodge(spell, unit, attackTable, roll, &chance) &&
			!spellEffect.applyEnemyAttackTableParry(spell, unit, attackTable, roll, &chance) &&
			!spellEffect.applyEnemyAttackTableBlock(spell, unit, attackTable, roll, &chance) &&
			!spellEffect.applyEnemyAttackTableCrit(spell, unit, attackTable, roll, &chance) {
			spellEffect.applyAttackTableHit(spell)
		}
	}
}

// Calculates a hit check using the stats from this spell.
func (spell *Spell) MagicHitCheck(sim *Simulation, attackTable *AttackTable) bool {
	missChance := attackTable.BaseSpellMissChance - spell.SpellHitChance(attackTable.Defender)
	return sim.RandomFloat("Magical Hit Roll") > missChance
}
func (spell *Spell) MagicHitCheckBinary(sim *Simulation, attackTable *AttackTable) bool {
	baseHitChance := (1 - attackTable.BaseSpellMissChance) * attackTable.GetBinaryHitChance(spell.SpellSchool)
	missChance := 1 - baseHitChance - spell.SpellHitChance(attackTable.Defender)
	return sim.RandomFloat("Magical Hit Roll") > missChance
}

func (spell *Spell) fixedCritCheck(sim *Simulation, critChance float64) bool {
	return sim.RandomFloat("Fixed Crit Roll") < critChance
}

func (spellEffect *SpellEffect) MagicCritCheck(sim *Simulation, spell *Spell, attackTable *AttackTable) bool {
	critChance := spellEffect.SpellCritChance(spell.Unit, spell)
	return sim.RandomFloat("Magical Crit Roll") < critChance
}

func (spellEffect *SpellEffect) HealingCritCheck(sim *Simulation, spell *Spell, attackTable *AttackTable) bool {
	critChance := spellEffect.HealingCritChance(spell.Unit, spell)
	return sim.RandomFloat("Healing Crit Roll") < critChance
}

func (spellEffect *SpellEffect) physicalCritRoll(sim *Simulation, spell *Spell, attackTable *AttackTable) bool {
	return sim.RandomFloat("Physical Crit Roll") < spellEffect.PhysicalCritChance(spell.Unit, spell, attackTable)
}

func (spellEffect *SpellEffect) applyAttackTableMiss(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	missChance := attackTable.BaseMissChance - spell.PhysicalHitChance(spellEffect.Target)
	if unit.AutoAttacks.IsDualWielding && !unit.PseudoStats.DisableDWMissPenalty {
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

func (spellEffect *SpellEffect) applyAttackTableMissNoDWPenalty(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
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

func (spellEffect *SpellEffect) applyAttackTableBlock(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += attackTable.BaseBlockChance

	if roll < *chance {
		spellEffect.Outcome |= OutcomeBlock
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Blocks++
		spellEffect.Damage = MaxFloat(0, spellEffect.Damage-spellEffect.Target.GetStat(stats.BlockValue))
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableDodge(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += MaxFloat(0, attackTable.BaseDodgeChance-spell.ExpertisePercentage()-unit.PseudoStats.DodgeReduction)

	if roll < *chance {
		spellEffect.Outcome = OutcomeDodge
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Dodges++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableParry(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += MaxFloat(0, attackTable.BaseParryChance-spell.ExpertisePercentage())

	if roll < *chance {
		spellEffect.Outcome = OutcomeParry
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Parries++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableGlance(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += attackTable.BaseGlanceChance

	if roll < *chance {
		spellEffect.Outcome = OutcomeGlance
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Glances++
		// TODO glancing blow damage reduction is actually a range ([65%, 85%] vs. 73)
		spellEffect.Damage *= attackTable.GlanceMultiplier
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableCrit(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	*chance += spellEffect.PhysicalCritChance(unit, spell, attackTable)

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
	if spellEffect.physicalCritRoll(sim, spell, attackTable) {
		spellEffect.Outcome = OutcomeCrit
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Crits++
		spellEffect.Damage *= spell.CritMultiplier
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableHit(spell *Spell) {
	spellEffect.Outcome = OutcomeHit
	spell.SpellMetrics[spellEffect.Target.UnitIndex].Hits++
}

func (spellEffect *SpellEffect) applyEnemyAttackTableMiss(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	missChance := attackTable.BaseMissChance + unit.PseudoStats.IncreasedMissChance + spellEffect.Target.stats[stats.Defense]*DefenseRatingToChanceReduction
	if unit.AutoAttacks.IsDualWielding && !unit.PseudoStats.DisableDWMissPenalty {
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

func (spellEffect *SpellEffect) applyEnemyAttackTableBlock(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
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

func (spellEffect *SpellEffect) applyEnemyAttackTableDodge(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	dodgeChance := attackTable.BaseDodgeChance +
		spellEffect.Target.stats[stats.Dodge]/DodgeRatingPerDodgeChance/100 +
		spellEffect.Target.stats[stats.Defense]*DefenseRatingToChanceReduction -
		unit.PseudoStats.DodgeReduction
	*chance += MaxFloat(0, dodgeChance)

	if roll < *chance {
		spellEffect.Outcome = OutcomeDodge
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Dodges++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyEnemyAttackTableParry(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	if !spellEffect.Target.PseudoStats.CanParry {
		return false
	}

	parryChance := attackTable.BaseParryChance +
		spellEffect.Target.stats[stats.Parry]/ParryRatingPerParryChance/100 +
		spellEffect.Target.stats[stats.Defense]*DefenseRatingToChanceReduction
	*chance += MaxFloat(0, parryChance)

	if roll < *chance {
		spellEffect.Outcome = OutcomeParry
		spell.SpellMetrics[spellEffect.Target.UnitIndex].Parries++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyEnemyAttackTableCrit(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	critRating := unit.stats[stats.MeleeCrit] + spell.BonusCritRating
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
