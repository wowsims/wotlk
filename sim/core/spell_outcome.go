package core

import (
	"github.com/wowsims/wotlk/sim/core/stats"
)

// This function should do 3 things:
//  1. Set the Outcome of the hit effect.
//  2. Update spell outcome metrics.
//  3. Modify the damage if necessary.
type OutcomeApplier func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable)

func (spell *Spell) ApplyOutcomeAlwaysHit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
	result.Outcome = OutcomeHit
	spell.SpellMetrics[result.Target.UnitIndex].Hits++
}
func (spell *Spell) CalcDamageAlwaysHit(sim *Simulation, target *Unit, baseDamage float64) SpellEffect {
	attackTable := spell.Unit.AttackTables[target.UnitIndex]
	result := spell.CalcDamagePreOutcome(sim, target, attackTable, baseDamage)
	spell.ApplyOutcomeAlwaysHit(sim, &result, attackTable)
	spell.ApplyPostOutcomeDamageModifiers(sim, &result)
	return result
}
func (spell *Spell) CalcAndDealDamageAlwaysHit(sim *Simulation, target *Unit, baseDamage float64) {
	result := spell.CalcDamageAlwaysHit(sim, target, baseDamage)
	spell.DealDamage(sim, &result)
}
func (unit *Unit) OutcomeFuncAlwaysHit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spell.ApplyOutcomeAlwaysHit(sim, spellEffect, attackTable)
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

func (spell *Spell) ApplyOutcomeMagicHitAndCrit(sim *Simulation, result *SpellEffect, attackTable *AttackTable) {
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
func (spell *Spell) CalcDamageMagicHitAndCrit(sim *Simulation, target *Unit, baseDamage float64) SpellEffect {
	attackTable := spell.Unit.AttackTables[target.UnitIndex]
	result := spell.CalcDamagePreOutcome(sim, target, attackTable, baseDamage)
	spell.ApplyOutcomeMagicHitAndCrit(sim, &result, attackTable)
	spell.ApplyPostOutcomeDamageModifiers(sim, &result)
	return result
}
func (spell *Spell) CalcAndDealDamageMagicHitAndCrit(sim *Simulation, target *Unit, baseDamage float64) {
	result := spell.CalcDamageMagicHitAndCrit(sim, target, baseDamage)
	spell.DealDamage(sim, &result)
}

func (unit *Unit) OutcomeFuncMagicHitAndCrit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spell.ApplyOutcomeMagicHitAndCrit(sim, spellEffect, attackTable)
	}
}

func (unit *Unit) OutcomeFuncMagicCrit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if spell.CritMultiplier == 0 {
			panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
		}
		if spellEffect.MagicCritCheck(sim, spell, attackTable) {
			spellEffect.Outcome = OutcomeCrit
			spell.SpellMetrics[spellEffect.Target.UnitIndex].Crits++
			spellEffect.Damage *= spell.CritMultiplier
		} else {
			spellEffect.Outcome = OutcomeHit
			spell.SpellMetrics[spellEffect.Target.UnitIndex].Hits++
		}
	}
}

func (unit *Unit) OutcomeFuncMagicHitAndCritBinary() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if spell.CritMultiplier == 0 {
			panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
		}
		if spell.MagicHitCheckBinary(sim, attackTable) {
			if spellEffect.MagicCritCheck(sim, spell, attackTable) {
				spellEffect.Outcome = OutcomeCrit
				spell.SpellMetrics[spellEffect.Target.UnitIndex].Crits++
				spellEffect.Damage *= spell.CritMultiplier
			} else {
				spellEffect.Outcome = OutcomeHit
				spell.SpellMetrics[spellEffect.Target.UnitIndex].Hits++
			}
		} else {
			spellEffect.Outcome = OutcomeMiss
			spell.SpellMetrics[spellEffect.Target.UnitIndex].Misses++
			spellEffect.Damage = 0
		}
	}
}

func (unit *Unit) OutcomeFuncHealingCrit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if spell.CritMultiplier == 0 {
			panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
		}
		if spellEffect.HealingCritCheck(sim, spell, attackTable) {
			spellEffect.Outcome = OutcomeCrit
			spell.SpellMetrics[spellEffect.Target.UnitIndex].Crits++
			spellEffect.Damage *= spell.CritMultiplier
		} else {
			spellEffect.Outcome = OutcomeHit
			spell.SpellMetrics[spellEffect.Target.UnitIndex].Hits++
		}
	}
}

func (unit *Unit) OutcomeFuncCritFixedChance(critChance float64) OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if spell.CritMultiplier == 0 {
			panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
		}
		if spell.fixedCritCheck(sim, critChance) {
			spellEffect.Outcome = OutcomeCrit
			spell.SpellMetrics[spellEffect.Target.UnitIndex].Crits++
			spellEffect.Damage *= spell.CritMultiplier
		} else {
			spellEffect.Outcome = OutcomeHit
			spell.SpellMetrics[spellEffect.Target.UnitIndex].Hits++
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

func (unit *Unit) OutcomeFuncMagicHit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if spell.MagicHitCheck(sim, attackTable) {
			spellEffect.Outcome = OutcomeHit
			spell.SpellMetrics[spellEffect.Target.UnitIndex].Hits++
		} else {
			spellEffect.Outcome = OutcomeMiss
			spell.SpellMetrics[spellEffect.Target.UnitIndex].Misses++
			spellEffect.Damage = 0
		}
	}
}

func (unit *Unit) OutcomeFuncMagicHitBinary() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if spell.MagicHitCheckBinary(sim, attackTable) {
			spellEffect.Outcome = OutcomeHit
			spell.SpellMetrics[spellEffect.Target.UnitIndex].Hits++
		} else {
			spellEffect.Outcome = OutcomeMiss
			spell.SpellMetrics[spellEffect.Target.UnitIndex].Misses++
			spellEffect.Damage = 0
		}
	}
}

func (unit *Unit) OutcomeFuncMeleeWhite() OutcomeApplier {
	if unit.PseudoStats.InFrontOfTarget {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			unit := spell.Unit
			roll := sim.RandomFloat("White Hit Table")
			chance := 0.0

			if !spellEffect.applyAttackTableMiss(spell, unit, attackTable, roll, &chance) &&
				!spellEffect.applyAttackTableDodge(spell, unit, attackTable, roll, &chance) &&
				!spellEffect.applyAttackTableParry(spell, unit, attackTable, roll, &chance) &&
				!spellEffect.applyAttackTableGlance(spell, unit, attackTable, roll, &chance) &&
				!spellEffect.applyAttackTableBlock(spell, unit, attackTable, roll, &chance) &&
				!spellEffect.applyAttackTableCrit(spell, unit, attackTable, roll, &chance) {
				spellEffect.applyAttackTableHit(spell)
			}
		}
	} else {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			unit := spell.Unit
			roll := sim.RandomFloat("White Hit Table")
			chance := 0.0

			if !spellEffect.applyAttackTableMiss(spell, unit, attackTable, roll, &chance) &&
				!spellEffect.applyAttackTableDodge(spell, unit, attackTable, roll, &chance) &&
				!spellEffect.applyAttackTableGlance(spell, unit, attackTable, roll, &chance) &&
				!spellEffect.applyAttackTableCrit(spell, unit, attackTable, roll, &chance) {
				spellEffect.applyAttackTableHit(spell)
			}
		}
	}
}

func (unit *Unit) OutcomeFuncMeleeSpecialHit() OutcomeApplier {
	if unit.PseudoStats.InFrontOfTarget {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			unit := spell.Unit
			roll := sim.RandomFloat("White Hit Table")
			chance := 0.0

			if !spellEffect.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
				(spell.Flags.Matches(SpellFlagCannotBeDodged) || !spellEffect.applyAttackTableDodge(spell, unit, attackTable, roll, &chance)) &&
				!spellEffect.applyAttackTableParry(spell, unit, attackTable, roll, &chance) {
				spellEffect.applyAttackTableHit(spell)
			}
		}
	} else {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			unit := spell.Unit
			roll := sim.RandomFloat("White Hit Table")
			chance := 0.0

			if !spellEffect.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
				(spell.Flags.Matches(SpellFlagCannotBeDodged) || !spellEffect.applyAttackTableDodge(spell, unit, attackTable, roll, &chance)) {
				spellEffect.applyAttackTableHit(spell)
			}
		}
	}
}

func (unit *Unit) OutcomeFuncMeleeSpecialHitAndCrit() OutcomeApplier {
	if unit.PseudoStats.InFrontOfTarget {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			unit := spell.Unit
			roll := sim.RandomFloat("White Hit Table")
			chance := 0.0

			if !spellEffect.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
				(spell.Flags.Matches(SpellFlagCannotBeDodged) || !spellEffect.applyAttackTableDodge(spell, unit, attackTable, roll, &chance)) &&
				!spellEffect.applyAttackTableParry(spell, unit, attackTable, roll, &chance) {
				if spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
					spellEffect.applyAttackTableBlock(spell, unit, attackTable, roll, &chance)
				} else {
					if !spellEffect.applyAttackTableBlock(spell, unit, attackTable, roll, &chance) {
						spellEffect.applyAttackTableHit(spell)
					}
				}
			}
		}
	} else {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			unit := spell.Unit
			roll := sim.RandomFloat("White Hit Table")
			chance := 0.0

			if !spellEffect.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
				(spell.Flags.Matches(SpellFlagCannotBeDodged) || !spellEffect.applyAttackTableDodge(spell, unit, attackTable, roll, &chance)) &&
				!spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
				spellEffect.applyAttackTableHit(spell)
			}
		}
	}
}

// Like OutcomeFuncMeleeSpecialHitAndCrit, but blocks prevent crits.
func (unit *Unit) OutcomeFuncMeleeWeaponSpecialHitAndCrit() OutcomeApplier {
	if unit.PseudoStats.InFrontOfTarget {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			unit := spell.Unit
			roll := sim.RandomFloat("White Hit Table")
			chance := 0.0

			if !spellEffect.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
				(spell.Flags.Matches(SpellFlagCannotBeDodged) || !spellEffect.applyAttackTableDodge(spell, unit, attackTable, roll, &chance)) &&
				!spellEffect.applyAttackTableParry(spell, unit, attackTable, roll, &chance) &&
				!spellEffect.applyAttackTableBlock(spell, unit, attackTable, roll, &chance) &&
				!spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
				spellEffect.applyAttackTableHit(spell)
			}
		}
	} else {
		return unit.OutcomeFuncMeleeSpecialHitAndCrit()
	}
}

func (unit *Unit) OutcomeFuncMeleeWeaponSpecialNoCrit() OutcomeApplier {
	if unit.PseudoStats.InFrontOfTarget {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			unit := spell.Unit
			roll := sim.RandomFloat("White Hit Table")
			chance := 0.0

			if !spellEffect.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
				(spell.Flags.Matches(SpellFlagCannotBeDodged) || !spellEffect.applyAttackTableDodge(spell, unit, attackTable, roll, &chance)) &&
				!spellEffect.applyAttackTableParry(spell, unit, attackTable, roll, &chance) &&
				!spellEffect.applyAttackTableBlock(spell, unit, attackTable, roll, &chance) {
				spellEffect.applyAttackTableHit(spell)
			}
		}
	} else {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			unit := spell.Unit
			roll := sim.RandomFloat("White Hit Table")
			chance := 0.0

			if !spellEffect.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
				spell.Flags.Matches(SpellFlagCannotBeDodged) || !spellEffect.applyAttackTableDodge(spell, unit, attackTable, roll, &chance) {
				spellEffect.applyAttackTableHit(spell)
			}
		}
	}
}

func (unit *Unit) OutcomeFuncMeleeSpecialNoBlockDodgeParry() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		unit := spell.Unit
		roll := sim.RandomFloat("White Hit Table")
		chance := 0.0

		if !spellEffect.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
			!spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
			spellEffect.applyAttackTableHit(spell)
		}
	}
}

func (unit *Unit) OutcomeFuncMeleeSpecialCritOnly() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if !spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
			spellEffect.applyAttackTableHit(spell)
		}
	}
}

func (unit *Unit) OutcomeFuncRangedHit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		unit := spell.Unit
		roll := sim.RandomFloat("White Hit Table")
		chance := 0.0

		if !spellEffect.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) {
			spellEffect.applyAttackTableHit(spell)
		}
	}
}

func (unit *Unit) OutcomeFuncRangedHitAndCrit() OutcomeApplier {
	if unit.PseudoStats.InFrontOfTarget {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			unit := spell.Unit
			roll := sim.RandomFloat("White Hit Table")
			chance := 0.0

			if !spellEffect.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) {
				if spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
					spellEffect.applyAttackTableBlock(spell, unit, attackTable, roll, &chance)
				} else {
					if !spellEffect.applyAttackTableBlock(spell, unit, attackTable, roll, &chance) {
						spellEffect.applyAttackTableHit(spell)
					}
				}
			}
		}
	} else {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			unit := spell.Unit
			roll := sim.RandomFloat("White Hit Table")
			chance := 0.0

			if !spellEffect.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
				!spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
				spellEffect.applyAttackTableHit(spell)
			}
		}
	}
}

func (unit *Unit) OutcomeFuncRangedCritOnly() OutcomeApplier {
	if unit.PseudoStats.InFrontOfTarget {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			unit := spell.Unit
			roll := sim.RandomFloat("White Hit Table")
			chance := 0.0

			if spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
				spellEffect.applyAttackTableBlock(spell, unit, attackTable, roll, &chance)
			} else {
				if !spellEffect.applyAttackTableBlock(spell, unit, attackTable, roll, &chance) {
					spellEffect.applyAttackTableHit(spell)
				}
			}
		}
	} else {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			if !spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable) {
				spellEffect.applyAttackTableHit(spell)
			}
		}
	}
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
