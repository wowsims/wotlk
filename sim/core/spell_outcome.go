package core

import (
	"github.com/wowsims/tbc/sim/core/stats"
)

// This function should do 3 things:
//  1. Set the Outcome of the hit effect.
//  2. Update spell outcome metrics.
//  3. Modify the damage if necessary.
type OutcomeApplier func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable)

func (unit *Unit) OutcomeFuncAlwaysHit() OutcomeApplier {
	return func(_ *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spellEffect.Outcome = OutcomeHit
		spell.SpellMetrics[spellEffect.Target.TableIndex].Hits++
	}
}

// A tick always hits, but we don't count them as hits in the metrics.
func (unit *Unit) OutcomeFuncTick() OutcomeApplier {
	return func(_ *Simulation, _ *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		spellEffect.Outcome = OutcomeHit
	}
}

func (unit *Unit) OutcomeFuncMagicHitAndCrit(critMultiplier float64) OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if spellEffect.magicHitCheck(sim, spell, attackTable) {
			if spellEffect.magicCritCheck(sim, spell, attackTable) {
				spellEffect.Outcome = OutcomeCrit
				spell.SpellMetrics[spellEffect.Target.TableIndex].Crits++
				spellEffect.Damage *= critMultiplier
			} else {
				spellEffect.Outcome = OutcomeHit
				spell.SpellMetrics[spellEffect.Target.TableIndex].Hits++
			}
		} else {
			spellEffect.Outcome = OutcomeMiss
			spell.SpellMetrics[spellEffect.Target.TableIndex].Misses++
			spellEffect.Damage = 0
		}
	}
}

func (unit *Unit) OutcomeFuncMagicCrit(critMultiplier float64) OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if spellEffect.magicCritCheck(sim, spell, attackTable) {
			spellEffect.Outcome = OutcomeCrit
			spell.SpellMetrics[spellEffect.Target.TableIndex].Crits++
			spellEffect.Damage *= critMultiplier
		} else {
			spellEffect.Outcome = OutcomeHit
			spell.SpellMetrics[spellEffect.Target.TableIndex].Hits++
		}
	}
}

func (unit *Unit) OutcomeFuncMagicHitAndCritBinary(critMultiplier float64) OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if spellEffect.magicHitCheckBinary(sim, spell, attackTable) {
			if spellEffect.magicCritCheck(sim, spell, attackTable) {
				spellEffect.Outcome = OutcomeCrit
				spell.SpellMetrics[spellEffect.Target.TableIndex].Crits++
				spellEffect.Damage *= critMultiplier
			} else {
				spellEffect.Outcome = OutcomeHit
				spell.SpellMetrics[spellEffect.Target.TableIndex].Hits++
			}
		} else {
			spellEffect.Outcome = OutcomeMiss
			spell.SpellMetrics[spellEffect.Target.TableIndex].Misses++
			spellEffect.Damage = 0
		}
	}
}

func (unit *Unit) OutcomeFuncCritFixedChance(critChance float64, critMultiplier float64) OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if spellEffect.fixedCritCheck(sim, critChance) {
			spellEffect.Outcome = OutcomeCrit
			spell.SpellMetrics[spellEffect.Target.TableIndex].Crits++
			spellEffect.Damage *= critMultiplier
		} else {
			spellEffect.Outcome = OutcomeHit
			spell.SpellMetrics[spellEffect.Target.TableIndex].Hits++
		}
	}
}

func (unit *Unit) OutcomeFuncMagicHit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if spellEffect.magicHitCheck(sim, spell, attackTable) {
			spellEffect.Outcome = OutcomeHit
			spell.SpellMetrics[spellEffect.Target.TableIndex].Hits++
		} else {
			spellEffect.Outcome = OutcomeMiss
			spell.SpellMetrics[spellEffect.Target.TableIndex].Misses++
			spellEffect.Damage = 0
		}
	}
}

func (unit *Unit) OutcomeFuncMagicHitBinary() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if spellEffect.magicHitCheckBinary(sim, spell, attackTable) {
			spellEffect.Outcome = OutcomeHit
			spell.SpellMetrics[spellEffect.Target.TableIndex].Hits++
		} else {
			spellEffect.Outcome = OutcomeMiss
			spell.SpellMetrics[spellEffect.Target.TableIndex].Misses++
			spellEffect.Damage = 0
		}
	}
}

func (unit *Unit) OutcomeFuncMeleeWhite(critMultiplier float64) OutcomeApplier {
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
				!spellEffect.applyAttackTableCrit(spell, unit, attackTable, roll, critMultiplier, &chance) {
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
				!spellEffect.applyAttackTableCrit(spell, unit, attackTable, roll, critMultiplier, &chance) {
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

func (unit *Unit) OutcomeFuncMeleeSpecialHitAndCrit(critMultiplier float64) OutcomeApplier {
	if unit.PseudoStats.InFrontOfTarget {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			unit := spell.Unit
			roll := sim.RandomFloat("White Hit Table")
			chance := 0.0

			if !spellEffect.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
				(spell.Flags.Matches(SpellFlagCannotBeDodged) || !spellEffect.applyAttackTableDodge(spell, unit, attackTable, roll, &chance)) &&
				!spellEffect.applyAttackTableParry(spell, unit, attackTable, roll, &chance) {
				if spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable, critMultiplier) {
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
				!spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable, critMultiplier) {
				spellEffect.applyAttackTableHit(spell)
			}
		}
	}
}

// Like OutcomeFuncMeleeSpecialHitAndCrit, but blocks prevent crits.
func (unit *Unit) OutcomeFuncMeleeWeaponSpecialHitAndCrit(critMultiplier float64) OutcomeApplier {
	if unit.PseudoStats.InFrontOfTarget {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			unit := spell.Unit
			roll := sim.RandomFloat("White Hit Table")
			chance := 0.0

			if !spellEffect.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
				(spell.Flags.Matches(SpellFlagCannotBeDodged) || !spellEffect.applyAttackTableDodge(spell, unit, attackTable, roll, &chance)) &&
				!spellEffect.applyAttackTableParry(spell, unit, attackTable, roll, &chance) &&
				!spellEffect.applyAttackTableBlock(spell, unit, attackTable, roll, &chance) &&
				!spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable, critMultiplier) {
				spellEffect.applyAttackTableHit(spell)
			}
		}
	} else {
		return unit.OutcomeFuncMeleeSpecialHitAndCrit(critMultiplier)
	}
}

func (unit *Unit) OutcomeFuncMeleeSpecialNoBlockDodgeParry(critMultiplier float64) OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		unit := spell.Unit
		roll := sim.RandomFloat("White Hit Table")
		chance := 0.0

		if !spellEffect.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) &&
			!spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable, critMultiplier) {
			spellEffect.applyAttackTableHit(spell)
		}
	}
}

func (unit *Unit) OutcomeFuncMeleeSpecialCritOnly(critMultiplier float64) OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
		if !spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable, critMultiplier) {
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

func (unit *Unit) OutcomeFuncRangedHitAndCrit(critMultiplier float64) OutcomeApplier {
	if unit.PseudoStats.InFrontOfTarget {
		return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, attackTable *AttackTable) {
			unit := spell.Unit
			roll := sim.RandomFloat("White Hit Table")
			chance := 0.0

			if !spellEffect.applyAttackTableMissNoDWPenalty(spell, unit, attackTable, roll, &chance) {
				if spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable, critMultiplier) {
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
				!spellEffect.applyAttackTableCritSeparateRoll(sim, spell, attackTable, critMultiplier) {
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
			!spellEffect.applyEnemyAttackTableCrit(spell, unit, attackTable, roll, &chance) &&
			!spellEffect.applyEnemyAttackTableCrush(spell, unit, attackTable, roll, &chance) {
			spellEffect.applyAttackTableHit(spell)
		}
	}
}

// Calculates a hit check using the stats from this spell.
func (spellEffect *SpellEffect) magicHitCheck(sim *Simulation, spell *Spell, attackTable *AttackTable) bool {
	missChance := attackTable.BaseSpellMissChance - (spell.Unit.GetStat(stats.SpellHit)+spellEffect.BonusSpellHitRating)/(SpellHitRatingPerHitChance*100)
	missChance = MaxFloat(missChance, 0.01) // can't get away from the 1% miss

	return sim.RandomFloat("Magical Hit Roll") > missChance
}
func (spellEffect *SpellEffect) magicHitCheckBinary(sim *Simulation, spell *Spell, attackTable *AttackTable) bool {
	baseHitChance := (1 - attackTable.BaseSpellMissChance) * attackTable.GetBinaryHitChance(spell.SpellSchool)
	missChance := 1 - baseHitChance - (spell.Unit.GetStat(stats.SpellHit)+spellEffect.BonusSpellHitRating)/(SpellHitRatingPerHitChance*100)
	missChance = MaxFloat(missChance, 0.01) // can't get away from the 1% miss

	return sim.RandomFloat("Magical Hit Roll") > missChance
}

func (spellEffect *SpellEffect) fixedCritCheck(sim *Simulation, critChance float64) bool {
	return sim.RandomFloat("Fixed Crit Roll") < critChance
}

func (spellEffect *SpellEffect) magicCritCheck(sim *Simulation, spell *Spell, attackTable *AttackTable) bool {
	critChance := spellEffect.SpellCritChance(spell.Unit, spell)
	return sim.RandomFloat("Magical Crit Roll") < critChance
}

func (spellEffect *SpellEffect) physicalCritRoll(sim *Simulation, spell *Spell, attackTable *AttackTable) bool {
	return sim.RandomFloat("Physical Crit Roll") < spellEffect.PhysicalCritChance(spell.Unit, spell, attackTable)
}

func (spellEffect *SpellEffect) applyAttackTableMiss(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	missChance := attackTable.BaseMissChance - spellEffect.PhysicalHitChance(unit, attackTable)
	if unit.AutoAttacks.IsDualWielding && !unit.PseudoStats.DisableDWMissPenalty {
		missChance += 0.19
	}
	*chance = MaxFloat(0, missChance)

	if roll < *chance {
		spellEffect.Outcome = OutcomeMiss
		spell.SpellMetrics[spellEffect.Target.TableIndex].Misses++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableMissNoDWPenalty(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	missChance := attackTable.BaseMissChance - spellEffect.PhysicalHitChance(unit, attackTable)
	*chance = MaxFloat(0, missChance)

	if roll < *chance {
		spellEffect.Outcome = OutcomeMiss
		spell.SpellMetrics[spellEffect.Target.TableIndex].Misses++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableBlock(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += attackTable.BaseBlockChance

	if roll < *chance {
		spellEffect.Outcome |= OutcomeBlock
		spell.SpellMetrics[spellEffect.Target.TableIndex].Blocks++
		spellEffect.Damage = MaxFloat(0, spellEffect.Damage-spellEffect.Target.GetStat(stats.BlockValue))
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableDodge(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += MaxFloat(0, attackTable.BaseDodgeChance-spellEffect.ExpertisePercentage(unit)-unit.PseudoStats.DodgeReduction)

	if roll < *chance {
		spellEffect.Outcome = OutcomeDodge
		spell.SpellMetrics[spellEffect.Target.TableIndex].Dodges++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableParry(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += MaxFloat(0, attackTable.BaseParryChance-spellEffect.ExpertisePercentage(unit))

	if roll < *chance {
		spellEffect.Outcome = OutcomeParry
		spell.SpellMetrics[spellEffect.Target.TableIndex].Parries++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableGlance(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += attackTable.BaseGlanceChance

	if roll < *chance {
		spellEffect.Outcome = OutcomeGlance
		spell.SpellMetrics[spellEffect.Target.TableIndex].Glances++
		// TODO glancing blow damage reduction is actually a range ([65%, 85%] vs. 73)
		spellEffect.Damage *= attackTable.GlanceMultiplier
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableCrit(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, critMultiplier float64, chance *float64) bool {
	*chance += spellEffect.PhysicalCritChance(unit, spell, attackTable)

	if roll < *chance {
		spellEffect.Outcome = OutcomeCrit
		spell.SpellMetrics[spellEffect.Target.TableIndex].Crits++
		spellEffect.Damage *= critMultiplier
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableCritSeparateRoll(sim *Simulation, spell *Spell, attackTable *AttackTable, critMultiplier float64) bool {
	if spellEffect.physicalCritRoll(sim, spell, attackTable) {
		spellEffect.Outcome = OutcomeCrit
		spell.SpellMetrics[spellEffect.Target.TableIndex].Crits++
		spellEffect.Damage *= critMultiplier
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyAttackTableHit(spell *Spell) {
	spellEffect.Outcome = OutcomeHit
	spell.SpellMetrics[spellEffect.Target.TableIndex].Hits++
}

func (spellEffect *SpellEffect) applyEnemyAttackTableMiss(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	missChance := attackTable.BaseMissChance + unit.PseudoStats.IncreasedMissChance + spellEffect.Target.stats[stats.Defense]*DefenseRatingToChanceReduction
	if unit.AutoAttacks.IsDualWielding && !unit.PseudoStats.DisableDWMissPenalty {
		missChance += 0.19
	}
	*chance = MaxFloat(0, missChance)

	if roll < *chance {
		spellEffect.Outcome = OutcomeMiss
		spell.SpellMetrics[spellEffect.Target.TableIndex].Misses++
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
		spell.SpellMetrics[spellEffect.Target.TableIndex].Blocks++
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
		spell.SpellMetrics[spellEffect.Target.TableIndex].Dodges++
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
		spell.SpellMetrics[spellEffect.Target.TableIndex].Parries++
		spellEffect.Damage = 0
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyEnemyAttackTableCrit(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	critRating := unit.stats[stats.MeleeCrit] + spellEffect.BonusCritRating
	critChance := critRating / (MeleeCritRatingPerCritChance * 100)
	critChance -= spellEffect.Target.stats[stats.Defense] * DefenseRatingToChanceReduction
	critChance -= spellEffect.Target.stats[stats.Resilience] / ResilienceRatingPerCritReductionChance / 100
	critChance -= spellEffect.Target.PseudoStats.ReducedCritTakenChance
	*chance += MaxFloat(0, critChance)

	if roll < *chance {
		spellEffect.Outcome = OutcomeCrit
		spell.SpellMetrics[spellEffect.Target.TableIndex].Crits++
		resilCritMultiplier := 1 - spellEffect.Target.stats[stats.Resilience]/ResilienceRatingPerCritDamageReductionPercent/100
		spellEffect.Damage *= 2 * resilCritMultiplier
		return true
	}
	return false
}

func (spellEffect *SpellEffect) applyEnemyAttackTableCrush(spell *Spell, unit *Unit, attackTable *AttackTable, roll float64, chance *float64) bool {
	if !unit.PseudoStats.CanCrush {
		return false
	}

	*chance += CrushChance

	if roll < *chance {
		spellEffect.Outcome = OutcomeCrush
		spell.SpellMetrics[spellEffect.Target.TableIndex].Crushes++
		spellEffect.Damage *= 1.5
		return true
	}
	return false
}
