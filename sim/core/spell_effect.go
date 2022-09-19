package core

import (
	"fmt"
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core/stats"
)

// Callback for after a spell hits the target and after damage is calculated. Use it for proc effects
// or anything that comes from the final result of the spell.
type OnSpellHit func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect)
type EffectOnSpellHitDealt func(sim *Simulation, spell *Spell, spellEffect *SpellEffect)

// OnPeriodicDamage is called when dots tick, after damage is calculated. Use it for proc effects
// or anything that comes from the final result of a tick.
type OnPeriodicDamage func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect)
type EffectOnPeriodicDamageDealt func(sim *Simulation, spell *Spell, spellEffect *SpellEffect)

type SpellEffect struct {
	// Target of the spell.
	Target *Unit

	BaseDamage     BaseDamageConfig // Callback for calculating base damage.
	OutcomeApplier OutcomeApplier   // Callback for determining outcome.

	// Used only for dot snapshotting. Internal-only.
	snapshotDamageMultiplier float64
	snapshotMeleeCritRating  float64
	snapshotSpellCritRating  float64

	// Use snapshotted values for crit/damage rather than recomputing them.
	isSnapshot bool

	// Used in determining snapshot based damage from effect details (e.g. snapshot crit and % damage modifiers)
	IsPeriodic bool

	// Indicates this is a healing spell, rather than a damage spell.
	IsHealing bool

	// Controls which effects can proc from this effect.
	ProcMask ProcMask

	// Callbacks for providing additional custom behavior.
	OnSpellHitDealt       EffectOnSpellHitDealt
	OnPeriodicDamageDealt EffectOnPeriodicDamageDealt

	// Results
	Outcome HitOutcome
	Damage  float64 // Damage done by this cast.
	Threat  float64

	PreoutcomeDamage float64 // Damage done by this cast before Outcome is applied.
}

func (spellEffect *SpellEffect) Validate() {
	if spellEffect.ProcMask == ProcMaskUnknown {
		panic("SpellEffects must set a ProcMask!")
	}
	if spellEffect.ProcMask.Matches(ProcMaskEmpty) && spellEffect.ProcMask != ProcMaskEmpty {
		panic("ProcMaskEmpty must be exclusive!")
	}
}

func (spellEffect *SpellEffect) Landed() bool {
	return spellEffect.Outcome.Matches(OutcomeLanded)
}

func (spellEffect *SpellEffect) DidCrit() bool {
	return spellEffect.Outcome.Matches(OutcomeCrit)
}

func (spellEffect *SpellEffect) calcThreat(spell *Spell) float64 {
	if spellEffect.Landed() {
		flatBonus := spell.FlatThreatBonus
		if spell.DynamicThreatBonus != nil {
			flatBonus += spell.DynamicThreatBonus(spellEffect, spell)
		}

		return (spellEffect.Damage*spell.ThreatMultiplier + flatBonus) * spell.Unit.PseudoStats.ThreatMultiplier
	} else {
		return 0
	}
}

func (spellEffect *SpellEffect) MeleeAttackPower(unit *Unit) float64 {
	return unit.stats[stats.AttackPower] + unit.PseudoStats.MobTypeAttackPower
}

func (spellEffect *SpellEffect) RangedAttackPower(unit *Unit) float64 {
	return unit.stats[stats.RangedAttackPower] + unit.PseudoStats.MobTypeAttackPower
}

func (spellEffect *SpellEffect) RangedAttackPowerOnTarget() float64 {
	return spellEffect.Target.PseudoStats.BonusRangedAttackPowerTaken
}

func (spellEffect *SpellEffect) BonusWeaponDamage(unit *Unit) float64 {
	return unit.PseudoStats.BonusDamage
}

func (spellEffect *SpellEffect) ExpertisePercentage(unit *Unit) float64 {
	expertiseRating := unit.stats[stats.Expertise]
	if spellEffect.ProcMask.Matches(ProcMaskMeleeMH) {
		expertiseRating += unit.PseudoStats.BonusMHExpertiseRating
	} else if spellEffect.ProcMask.Matches(ProcMaskMeleeOH) {
		expertiseRating += unit.PseudoStats.BonusOHExpertiseRating
	}
	return math.Floor(expertiseRating/ExpertisePerQuarterPercentReduction) / 400
}

func (spellEffect *SpellEffect) PhysicalHitChance(unit *Unit, spell *Spell, attackTable *AttackTable) float64 {
	hitRating := unit.stats[stats.MeleeHit] +
		spell.BonusHitRating +
		spellEffect.Target.PseudoStats.BonusMeleeHitRatingTaken
	return hitRating / (MeleeHitRatingPerHitChance * 100)
}

func (spellEffect *SpellEffect) PhysicalCritChance(unit *Unit, spell *Spell, attackTable *AttackTable) float64 {
	critRating := 0.0
	if spellEffect.isSnapshot {
		// periodic spells apply crit from snapshot at time of initial cast if capable of a crit
		// ignoring units real time crit in this case
		critRating = spellEffect.snapshotMeleeCritRating
	} else {
		critRating = spellEffect.physicalCritRating(unit, spell)
	}

	return (critRating / (CritRatingPerCritChance * 100)) - attackTable.CritSuppression
}
func (spellEffect *SpellEffect) physicalCritRating(unit *Unit, spell *Spell) float64 {
	critRating := unit.stats[stats.MeleeCrit] +
		spell.BonusCritRating +
		spellEffect.Target.PseudoStats.BonusCritRatingTaken

	if spellEffect.ProcMask.Matches(ProcMaskMeleeMH) {
		critRating += unit.PseudoStats.BonusMHCritRating
	} else if spellEffect.ProcMask.Matches(ProcMaskMeleeOH) {
		critRating += unit.PseudoStats.BonusOHCritRating
	}
	return critRating
}

func (spellEffect *SpellEffect) SpellPower(unit *Unit, spell *Spell) float64 {
	return unit.GetStat(stats.SpellPower) + unit.GetStat(spell.SpellSchool.Stat()) + unit.PseudoStats.MobTypeSpellPower
}

func (spellEffect *SpellEffect) SpellHitChance(spell *Spell) float64 {
	hitRating := spell.Unit.stats[stats.SpellHit] +
		spell.BonusHitRating +
		spellEffect.Target.PseudoStats.BonusSpellHitRatingTaken

	return hitRating / (SpellHitRatingPerHitChance * 100)
}

func (spellEffect *SpellEffect) SpellCritChance(unit *Unit, spell *Spell) float64 {
	critRating := 0.0
	if spellEffect.isSnapshot {
		// periodic spells apply crit from snapshot at time of initial cast if capable of a crit
		// ignoring units real time crit in this case
		critRating = spellEffect.snapshotSpellCritRating
	} else {
		critRating = spellEffect.spellCritRating(unit, spell)
	}

	return critRating / (CritRatingPerCritChance * 100)
}
func (spellEffect *SpellEffect) spellCritRating(unit *Unit, spell *Spell) float64 {
	return unit.stats[stats.SpellCrit] +
		spell.BonusCritRating +
		spellEffect.Target.PseudoStats.BonusCritRatingTaken +
		spellEffect.Target.PseudoStats.BonusSpellCritRatingTaken
}

func (spellEffect *SpellEffect) HealingPower(unit *Unit, _ *Spell) float64 {
	return unit.GetStat(stats.HealingPower)
}
func (spellEffect *SpellEffect) HealingCritChance(unit *Unit, spell *Spell) float64 {
	critRating := 0.0
	if spellEffect.isSnapshot {
		// periodic spells apply crit from snapshot at time of initial cast if capable of a crit
		// ignoring units real time crit in this case
		critRating = spellEffect.snapshotSpellCritRating
	} else {
		critRating = spellEffect.healingCritRating(unit, spell)
	}

	return critRating / (CritRatingPerCritChance * 100)
}
func (spellEffect *SpellEffect) healingCritRating(unit *Unit, spell *Spell) float64 {
	return unit.GetStat(stats.SpellCrit) + spell.BonusCritRating
}

func (spellEffect *SpellEffect) calculateBaseDamage(sim *Simulation, spell *Spell) float64 {
	if spellEffect.BaseDamage.Calculator == nil {
		return 0
	} else {
		return spellEffect.BaseDamage.Calculator(sim, spellEffect, spell)
	}
}

func (spellEffect *SpellEffect) calcDamageSingle(sim *Simulation, spell *Spell, attackTable *AttackTable) {
	if sim.Log != nil {
		baseDmg := spellEffect.Damage
		spellEffect.applyAttackerModifiers(sim, spell)
		afterAttackMods := spellEffect.Damage
		spellEffect.applyResistances(sim, spell, attackTable)
		afterResistances := spellEffect.Damage
		spellEffect.applyTargetModifiers(sim, spell, attackTable)
		afterTargetMods := spellEffect.Damage
		spellEffect.PreoutcomeDamage = spellEffect.Damage
		spellEffect.OutcomeApplier(sim, spell, spellEffect, attackTable)
		afterOutcome := spellEffect.Damage
		spell.Unit.Log(
			sim,
			"%s %s [DEBUG] MAP: %0.01f, RAP: %0.01f, SP: %0.01f, BaseDamage:%0.01f, AfterAttackerMods:%0.01f, AfterResistances:%0.01f, AfterTargetMods:%0.01f, AfterOutcome:%0.01f",
			spellEffect.Target.LogLabel(), spell.ActionID, spell.Unit.GetStat(stats.AttackPower), spell.Unit.GetStat(stats.RangedAttackPower), spell.Unit.GetStat(stats.SpellPower), baseDmg, afterAttackMods, afterResistances, afterTargetMods, afterOutcome)
	} else {
		spellEffect.applyAttackerModifiers(sim, spell)
		spellEffect.applyResistances(sim, spell, attackTable)
		spellEffect.applyTargetModifiers(sim, spell, attackTable)
		spellEffect.PreoutcomeDamage = spellEffect.Damage
		spellEffect.OutcomeApplier(sim, spell, spellEffect, attackTable)
	}
}
func (spellEffect *SpellEffect) calcDamageTargetOnly(sim *Simulation, spell *Spell, attackTable *AttackTable) {
	spellEffect.applyResistances(sim, spell, attackTable)
	spellEffect.applyTargetModifiers(sim, spell, attackTable)
	spellEffect.OutcomeApplier(sim, spell, spellEffect, attackTable)
}

func (spellEffect *SpellEffect) finalize(sim *Simulation, spell *Spell) {
	if spell.MissileSpeed == 0 {
		if spellEffect.IsHealing {
			spellEffect.finalizeHealingInternal(sim, spell)
		} else {
			spellEffect.finalizeInternal(sim, spell)
		}
	} else {
		travelTime := time.Duration(float64(time.Second) * spell.Unit.DistanceFromTarget / spell.MissileSpeed)

		// We need to make a copy of this SpellEffect because some spells re-use the effect objects.
		effectCopy := *spellEffect

		StartDelayedAction(sim, DelayedActionOptions{
			DoAt: sim.CurrentTime + travelTime,
			OnAction: func(sim *Simulation) {
				if spellEffect.IsHealing {
					effectCopy.finalizeHealingInternal(sim, spell)
				} else {
					effectCopy.finalizeInternal(sim, spell)
				}
			},
		})
	}
}

// Applies the fully computed results from this SpellEffect to the sim.
func (spellEffect *SpellEffect) finalizeInternal(sim *Simulation, spell *Spell) {
	for i := range spellEffect.Target.DynamicDamageTakenModifiers {
		spellEffect.Target.DynamicDamageTakenModifiers[i](sim, spellEffect)
	}

	spell.SpellMetrics[spellEffect.Target.UnitIndex].TotalDamage += spellEffect.Damage
	spell.SpellMetrics[spellEffect.Target.UnitIndex].TotalThreat += spellEffect.calcThreat(spell)

	// Mark total damage done in raid so far for health based fights.
	// Don't include damage done by EnemyUnits to Players
	if spellEffect.Target.Type == EnemyUnit {
		sim.Encounter.DamageTaken += spellEffect.Damage
	}

	if sim.Log != nil {
		if spellEffect.IsPeriodic {
			spell.Unit.Log(sim, "%s %s tick %s. (Threat: %0.3f)", spellEffect.Target.LogLabel(), spell.ActionID, spellEffect, spellEffect.calcThreat(spell))
		} else {
			spell.Unit.Log(sim, "%s %s %s. (Threat: %0.3f)", spellEffect.Target.LogLabel(), spell.ActionID, spellEffect, spellEffect.calcThreat(spell))
		}
	}

	if !spellEffect.IsPeriodic {
		if spellEffect.OnSpellHitDealt != nil {
			spellEffect.OnSpellHitDealt(sim, spell, spellEffect)
		}
		spell.Unit.OnSpellHitDealt(sim, spell, spellEffect)
		spellEffect.Target.OnSpellHitTaken(sim, spell, spellEffect)
	} else {
		if spellEffect.OnPeriodicDamageDealt != nil {
			spellEffect.OnPeriodicDamageDealt(sim, spell, spellEffect)
		}
		spell.Unit.OnPeriodicDamageDealt(sim, spell, spellEffect)
		spellEffect.Target.OnPeriodicDamageTaken(sim, spell, spellEffect)
	}
}

func (spellEffect *SpellEffect) finalizeHealingInternal(sim *Simulation, spell *Spell) {
	spell.SpellMetrics[spellEffect.Target.UnitIndex].TotalThreat += spellEffect.calcThreat(spell)
	spell.SpellMetrics[spellEffect.Target.UnitIndex].TotalHealing += spellEffect.Damage
	if spellEffect.Target.HasHealthBar() {
		spellEffect.Target.GainHealth(sim, spellEffect.Damage, spell.HealthMetrics(spellEffect.Target))
	}

	if sim.Log != nil {
		if spellEffect.IsPeriodic {
			spell.Unit.Log(sim, "%s %s tick %s. (Threat: %0.3f)", spellEffect.Target.LogLabel(), spell.ActionID, spellEffect, spellEffect.calcThreat(spell))
		} else {
			spell.Unit.Log(sim, "%s %s %s. (Threat: %0.3f)", spellEffect.Target.LogLabel(), spell.ActionID, spellEffect, spellEffect.calcThreat(spell))
		}
	}

	if !spellEffect.IsPeriodic {
		if spellEffect.OnSpellHitDealt != nil {
			spellEffect.OnSpellHitDealt(sim, spell, spellEffect)
		}
		spell.Unit.OnHealDealt(sim, spell, spellEffect)
		spellEffect.Target.OnHealTaken(sim, spell, spellEffect)
	} else {
		if spellEffect.OnPeriodicDamageDealt != nil {
			spellEffect.OnPeriodicDamageDealt(sim, spell, spellEffect)
		}
		spell.Unit.OnPeriodicHealDealt(sim, spell, spellEffect)
		spellEffect.Target.OnPeriodicHealTaken(sim, spell, spellEffect)
	}
}

func (spellEffect *SpellEffect) String() string {
	outcomeStr := spellEffect.Outcome.String()
	if !spellEffect.Landed() {
		return outcomeStr
	}
	if spellEffect.IsHealing {
		return fmt.Sprintf("%s for %0.3f healing", outcomeStr, spellEffect.Damage)
	} else {
		return fmt.Sprintf("%s for %0.3f damage", outcomeStr, spellEffect.Damage)
	}
}

func (spellEffect *SpellEffect) applyAttackerModifiers(sim *Simulation, spell *Spell) {
	if spell.Flags.Matches(SpellFlagIgnoreAttackerModifiers) {
		// Even when ignoring attacker multipliers we still apply this one, because its
		// specific to the spell.
		spellEffect.Damage *= spell.DamageMultiplier * spell.DamageMultiplierAdditive
		return
	}

	// For dot snapshots, everything has already been stored in spellEffect.snapshotDamageMultiplier.
	if spellEffect.isSnapshot {
		spellEffect.Damage *= spellEffect.snapshotDamageMultiplier
		return
	}

	attacker := spell.Unit

	if spellEffect.IsHealing {
		spellEffect.Damage *= attacker.PseudoStats.HealingDealtMultiplier
		return
	}

	spellEffect.Damage *= spellEffect.snapshotAttackModifiers(spell)
}

// Returns the combined attacker modifiers. For snapshot dots, these are precomputed and stored.
func (spellEffect *SpellEffect) snapshotAttackModifiers(spell *Spell) float64 {
	if spell.Flags.Matches(SpellFlagIgnoreAttackerModifiers) {
		return 1.0
	}

	attacker := spell.Unit

	multiplier := attacker.PseudoStats.DamageDealtMultiplier

	multiplier *= spell.DamageMultiplier
	multiplier *= spell.DamageMultiplierAdditive

	if spellEffect.ProcMask.Matches(ProcMaskRanged) {
		multiplier *= attacker.PseudoStats.RangedDamageDealtMultiplier
	}
	if spell.Flags.Matches(SpellFlagDisease) {
		multiplier *= attacker.PseudoStats.DiseaseDamageDealtMultiplier
	}

	if spell.SpellSchool.Matches(SpellSchoolPhysical) {
		multiplier *= attacker.PseudoStats.PhysicalDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolArcane) {
		multiplier *= attacker.PseudoStats.ArcaneDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolFire) {
		multiplier *= attacker.PseudoStats.FireDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolFrost) {
		multiplier *= attacker.PseudoStats.FrostDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolHoly) {
		multiplier *= attacker.PseudoStats.HolyDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolNature) {
		multiplier *= attacker.PseudoStats.NatureDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolShadow) {
		multiplier *= attacker.PseudoStats.ShadowDamageDealtMultiplier
	}

	return multiplier
}

func (spellEffect *SpellEffect) applyTargetModifiers(sim *Simulation, spell *Spell, attackTable *AttackTable) {
	if spell.Flags.Matches(SpellFlagIgnoreTargetModifiers) {
		return
	}

	target := spellEffect.Target

	if spellEffect.IsHealing {
		spellEffect.Damage *= target.PseudoStats.HealingTakenMultiplier * attackTable.HealingDealtMultiplier
		return
	}

	spellEffect.Damage *= attackTable.DamageDealtMultiplier
	spellEffect.Damage *= target.PseudoStats.DamageTakenMultiplier
	spellEffect.Damage = MaxFloat(0, spellEffect.Damage+target.PseudoStats.BonusDamageTaken)

	if spell.Flags.Matches(SpellFlagDisease) {
		spellEffect.Damage *= target.PseudoStats.DiseaseDamageTakenMultiplier
	}

	if spell.SpellSchool.Matches(SpellSchoolPhysical) {
		if spellEffect.IsPeriodic {
			spellEffect.Damage *= target.PseudoStats.PeriodicPhysicalDamageTakenMultiplier
		}
		if spellEffect.BaseDamage.TargetSpellCoefficient > 0 {
			spellEffect.Damage += target.PseudoStats.BonusPhysicalDamageTaken
		}
		spellEffect.Damage *= target.PseudoStats.PhysicalDamageTakenMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolArcane) {
		spellEffect.Damage *= target.PseudoStats.ArcaneDamageTakenMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolFire) {
		spellEffect.Damage *= target.PseudoStats.FireDamageTakenMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolFrost) {
		spellEffect.Damage *= target.PseudoStats.FrostDamageTakenMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolHoly) {
		spellEffect.Damage += target.PseudoStats.BonusHolyDamageTaken * spellEffect.BaseDamage.TargetSpellCoefficient
		spellEffect.Damage *= target.PseudoStats.HolyDamageTakenMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolNature) {
		spellEffect.Damage *= target.PseudoStats.NatureDamageTakenMultiplier
		spellEffect.Damage *= attackTable.NatureDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolShadow) {
		spellEffect.Damage *= target.PseudoStats.ShadowDamageTakenMultiplier
		if spellEffect.IsPeriodic {
			spellEffect.Damage *= attackTable.PeriodicShadowDamageDealtMultiplier
		}
	}
}
