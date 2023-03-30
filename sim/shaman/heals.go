package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (shaman *Shaman) registerAncestralHealingSpell() {
	shaman.AncestralAwakening = shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 52752},
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskSpellHealing,
		Flags:            core.SpellFlagHelpful,
		DamageMultiplier: 1 * (1 + .02*float64(shaman.Talents.Purification)),
		CritMultiplier:   1,
		ThreatMultiplier: 1 - (float64(shaman.Talents.HealingGrace) * 0.05),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealHealing(sim, target, shaman.ancestralHealingAmount, spell.OutcomeHealing)
		},
	})
}

func (shaman *Shaman) registerLesserHealingWaveSpell() {
	spellCoeff := 0.807
	bonusCoeff := 0.02 * float64(shaman.Talents.TidalWaves)
	impShieldChance := 0.2 * float64(shaman.Talents.ImprovedWaterShield)
	impShieldManaGain := 428.0 * (1 + 0.05*float64(shaman.Talents.ImprovedShields))

	shaman.LesserHealingWave = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49276},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.15,
			Multiplier: 1 *
				(1 - .01*float64(shaman.Talents.TidalFocus)),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		BonusCritRating: float64(shaman.Talents.TidalMastery)*1*core.CritRatingPerCritChance +
			float64(shaman.Talents.BlessingOfTheEternals)*2*core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + .02*float64(shaman.Talents.Purification)),
		CritMultiplier:   shaman.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1 - (float64(shaman.Talents.HealingGrace) * 0.05),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healPower := spell.HealingPower(target)
			baseHealing := sim.Roll(1624, 1852) + spellCoeff*healPower + bonusCoeff*healPower
			result := spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

			if result.Outcome.Matches(core.OutcomeCrit) {
				if impShieldChance > 0 {
					if sim.RandomFloat("imp water shield") > impShieldChance {
						shaman.AddMana(sim, impShieldManaGain, shaman.waterShieldManaMetrics)
					}
				}
				if shaman.Talents.AncestralAwakening > 0 {
					shaman.ancestralHealingAmount = result.Damage * 0.3

					// TODO: this should actually target the lowest health target in the raid.
					//  does it matter in a sim? We currently only simulate tanks taking damage (multiple tanks could be handled here though.)
					shaman.AncestralAwakening.Cast(sim, target)
				}
			}
		},
	})
}

func (shaman *Shaman) registerRiptideSpell() {
	spellCoeff := 0.402
	impShieldChance := []float64{0.33, 0.66, 1.0}[shaman.Talents.ImprovedWaterShield]
	impShieldManaGain := 428.0 * (1 + 0.05*float64(shaman.Talents.ImprovedShields))

	shaman.Riptide = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 61301},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.18,
			Multiplier: 1 *
				(1 - .01*float64(shaman.Talents.TidalFocus)),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		BonusCritRating: float64(shaman.Talents.TidalMastery)*1*core.CritRatingPerCritChance +
			float64(shaman.Talents.BlessingOfTheEternals)*2*core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + .02*float64(shaman.Talents.Purification)),
		CritMultiplier:   shaman.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1 - (float64(shaman.Talents.HealingGrace) * 0.05),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healPower := spell.HealingPower(target)
			baseHealing := sim.Roll(1604, 1736) + spellCoeff*healPower
			result := spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

			if result.Outcome.Matches(core.OutcomeCrit) {
				if impShieldChance > 0 {
					if impShieldChance > 0.9999 || sim.RandomFloat("imp water shield") > impShieldChance {
						shaman.AddMana(sim, impShieldManaGain, shaman.waterShieldManaMetrics)
					}
				}
				if shaman.Talents.AncestralAwakening > 0 {
					shaman.ancestralHealingAmount = result.Damage * 0.3
					// TODO: this should actually target the lowest health target in the raid.
					//  does it matter in a sim? We currently only simulate tanks taking damage (multiple tanks could be handled here though.)
					shaman.AncestralAwakening.Cast(sim, target)
				}
			}
		},
	})
}

func (shaman *Shaman) registerHealingWaveSpell() {
	// TODO: finish this.

	spellCoeff := 0.807
	bonusCoeff := 0.02 * float64(shaman.Talents.TidalWaves)
	impShieldChance := 0.2 * float64(shaman.Talents.ImprovedWaterShield)
	impShieldManaGain := 428.0 * (1 + 0.05*float64(shaman.Talents.ImprovedShields))

	shaman.HealingWave = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.15,
			Multiplier: 1 *
				(1 - .01*float64(shaman.Talents.TidalFocus)),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		BonusCritRating: float64(shaman.Talents.TidalMastery)*1*core.CritRatingPerCritChance +
			float64(shaman.Talents.BlessingOfTheEternals)*2*core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + .02*float64(shaman.Talents.Purification)),
		CritMultiplier:   shaman.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1 - (float64(shaman.Talents.HealingGrace) * 0.05),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healPower := spell.HealingPower(target)
			baseHealing := sim.Roll(1624, 1852) + spellCoeff*healPower + bonusCoeff*healPower
			result := spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

			if result.Outcome.Matches(core.OutcomeCrit) {
				if impShieldChance > 0 {
					if sim.RandomFloat("imp water shield") > impShieldChance {
						shaman.AddMana(sim, impShieldManaGain, shaman.waterShieldManaMetrics)
					}
				}
				if shaman.Talents.AncestralAwakening > 0 {
					shaman.ancestralHealingAmount = result.Damage * 0.3

					// TODO: this should actually target the lowest health target in the raid.
					//  does it matter in a sim? We currently only simulate tanks taking damage (multiple tanks could be handled here though.)
					shaman.AncestralAwakening.Cast(sim, target)
				}
			}
		},
	})
}
