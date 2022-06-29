package core

import (
	"github.com/wowsims/tbc/sim/core/stats"
)

// Modifies damage based on Armor or Magic resistances, depending on the damage type.
func (spellEffect *SpellEffect) applyResistances(sim *Simulation, spell *Spell, attackTable *AttackTable) {
	if spell.Flags.Matches(SpellFlagIgnoreResists) {
		return
	}

	if spell.SpellSchool.Matches(SpellSchoolPhysical) {
		// All physical dots (Bleeds) ignore armor.
		if spellEffect.IsPeriodic {
			return
		}

		// Physical resistance (armor).
		spellEffect.Damage *= attackTable.ArmorDamageReduction
	} else if !spell.Flags.Matches(SpellFlagBinary) {
		// Magical resistance.

		resistanceRoll := sim.RandomFloat("Partial Resist")

		threshold00, threshold25, threshold50 := attackTable.GetPartialResistThresholds(spell.SpellSchool)
		//if sim.Log != nil {
		//	sim.Log("Resist thresholds: %0.04f, %0.04f, %0.04f", threshold00, threshold25, threshold50)
		//}

		if resistanceRoll > threshold00 {
			// No partial resist.
		} else if resistanceRoll > threshold25 {
			spellEffect.Outcome |= OutcomePartial1_4
			spellEffect.Damage *= 0.75
		} else if resistanceRoll > threshold50 {
			spellEffect.Outcome |= OutcomePartial2_4
			spellEffect.Damage *= 0.5
		} else {
			spellEffect.Outcome |= OutcomePartial3_4
			spellEffect.Damage *= 0.25
		}
	}
}

// ArmorDamageReduction currently assumes a level 70 attacker
func (at *AttackTable) UpdateArmorDamageReduction() {
	effectiveArmor := MaxFloat(0, at.Defender.stats[stats.Armor]-at.Attacker.stats[stats.ArmorPenetration])
	at.ArmorDamageReduction = MaxFloat(0.25, 1-(effectiveArmor/(effectiveArmor+(float64(at.Attacker.Level)*467.5-22167.5))))
}

func (at *AttackTable) UpdatePartialResists() {
	at.PartialResistArcaneRollThreshold00, at.PartialResistArcaneRollThreshold25, at.PartialResistArcaneRollThreshold50, at.BinaryArcaneHitChance = at.Defender.partialResistRollThresholds(SpellSchoolArcane, at.Attacker)
	at.PartialResistHolyRollThreshold00, at.PartialResistHolyRollThreshold25, at.PartialResistHolyRollThreshold50, at.BinaryHolyHitChance = at.Defender.partialResistRollThresholds(SpellSchoolHoly, at.Attacker)
	at.PartialResistFireRollThreshold00, at.PartialResistFireRollThreshold25, at.PartialResistFireRollThreshold50, at.BinaryFireHitChance = at.Defender.partialResistRollThresholds(SpellSchoolFire, at.Attacker)
	at.PartialResistFrostRollThreshold00, at.PartialResistFrostRollThreshold25, at.PartialResistFrostRollThreshold50, at.BinaryFrostHitChance = at.Defender.partialResistRollThresholds(SpellSchoolFrost, at.Attacker)
	at.PartialResistNatureRollThreshold00, at.PartialResistNatureRollThreshold25, at.PartialResistNatureRollThreshold50, at.BinaryNatureHitChance = at.Defender.partialResistRollThresholds(SpellSchoolNature, at.Attacker)
	at.PartialResistShadowRollThreshold00, at.PartialResistShadowRollThreshold25, at.PartialResistShadowRollThreshold50, at.BinaryShadowHitChance = at.Defender.partialResistRollThresholds(SpellSchoolShadow, at.Attacker)
}

func (at *AttackTable) GetPartialResistThresholds(ss SpellSchool) (float64, float64, float64) {
	switch ss {
	case SpellSchoolArcane:
		return at.PartialResistArcaneRollThreshold00, at.PartialResistArcaneRollThreshold25, at.PartialResistArcaneRollThreshold50
	case SpellSchoolHoly:
		return at.PartialResistHolyRollThreshold00, at.PartialResistHolyRollThreshold25, at.PartialResistHolyRollThreshold50
	case SpellSchoolFire:
		return at.PartialResistFireRollThreshold00, at.PartialResistFireRollThreshold25, at.PartialResistFireRollThreshold50
	case SpellSchoolFrost:
		return at.PartialResistFrostRollThreshold00, at.PartialResistFrostRollThreshold25, at.PartialResistFrostRollThreshold50
	case SpellSchoolNature:
		return at.PartialResistNatureRollThreshold00, at.PartialResistNatureRollThreshold25, at.PartialResistNatureRollThreshold50
	case SpellSchoolShadow:
		return at.PartialResistShadowRollThreshold00, at.PartialResistShadowRollThreshold25, at.PartialResistShadowRollThreshold50
	}
	return 0, 0, 0
}

func (at *AttackTable) GetBinaryHitChance(ss SpellSchool) float64 {
	switch ss {
	case SpellSchoolArcane:
		return at.BinaryArcaneHitChance
	case SpellSchoolHoly:
		return at.BinaryHolyHitChance
	case SpellSchoolFire:
		return at.BinaryFireHitChance
	case SpellSchoolFrost:
		return at.BinaryFrostHitChance
	case SpellSchoolNature:
		return at.BinaryNatureHitChance
	case SpellSchoolShadow:
		return at.BinaryShadowHitChance
	}
	return 0
}

// All of the following calculations are based on this guide:
// https://royalgiraffe.github.io/resist-guide

func (unit *Unit) resistCoeff(school SpellSchool, attacker *Unit, binary bool) float64 {
	resistanceCap := float64(unit.Level * 5)

	resistance := MaxFloat(0, unit.GetStat(school.ResistanceStat())-attacker.stats[stats.SpellPenetration])
	if school == SpellSchoolHoly {
		resistance = 0
	}

	effectiveResistance := resistance
	if !binary {
		levelBasedResistance := 0.0
		if unit.Type == EnemyUnit {
			levelBasedResistance = LevelBasedNPCSpellResistancePerLevel * float64(MaxInt32(0, unit.Level-attacker.Level))
		}
		effectiveResistance += levelBasedResistance
	}

	return MinFloat(resistanceCap, effectiveResistance) / resistanceCap
}

func (unit *Unit) binaryHitChance(school SpellSchool, attacker *Unit) float64 {
	resistCoeff := unit.resistCoeff(school, attacker, true)
	return 1 - 0.75*resistCoeff
}

// Roll threshold for each type of partial resist.
// Also returns binary miss chance as 4th value.
func (unit *Unit) partialResistRollThresholds(school SpellSchool, attacker *Unit) (float64, float64, float64, float64) {
	resistCoeff := unit.resistCoeff(school, attacker, false)

	// Based on the piecewise linear regression estimates at https://royalgiraffe.github.io/partial-resist-table.
	//partialResistChance00 := piecewiseLinear3(resistCoeff, 1, 0.24, 0.00, 0.00)
	partialResistChance25 := piecewiseLinear3(resistCoeff, 0, 0.55, 0.22, 0.04)
	partialResistChance50 := piecewiseLinear3(resistCoeff, 0, 0.18, 0.56, 0.16)
	partialResistChance75 := piecewiseLinear3(resistCoeff, 0, 0.03, 0.22, 0.80)

	return partialResistChance25 + partialResistChance50 + partialResistChance75,
		partialResistChance50 + partialResistChance75,
		partialResistChance75,
		unit.binaryHitChance(school, attacker)
}

// Interpolation for a 3-part piecewise linear function (which all the partial resist equations use).
func piecewiseLinear3(val float64, p0 float64, p1 float64, p2 float64, p3 float64) float64 {
	if val < 1.0/3.0 {
		return interpolate(val*3, p0, p1)
	} else if val < 2.0/3.0 {
		return interpolate((val-1.0/3.0)*3, p1, p2)
	} else {
		return interpolate((val-2.0/3.0)*3, p2, p3)
	}
}

func interpolate(val float64, p0 float64, p1 float64) float64 {
	return p0*(1-val) + p1*val
}
