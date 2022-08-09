package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var MHOutcome = core.OutcomeHit
var OHOutcome = core.OutcomeHit

func (rogue *Rogue) newMutilateHitSpell(isMH bool) *core.Spell {
	actionID := core.ActionID{SpellID: 48665}
	if !isMH {
		actionID = core.ActionID{SpellID: 48664}
	}

	effect := core.SpellEffect{
		ProcMask: core.ProcMaskMeleeMHSpecial,

		BonusCritRating: 5 * core.CritRatingPerCritChance * float64(rogue.Talents.PuncturingWounds),
		DamageMultiplier: 1 +
			0.1*float64(rogue.Talents.Opportunity) +
			core.TernaryFloat64(rogue.HasSetBonus(ItemSetSlayers, 4), 0.06, 0),
		ThreatMultiplier: 1,

		BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, true, 181, 1, 1, false),
		OutcomeApplier: rogue.OutcomeFuncMeleeSpecialCritOnly(rogue.MeleeCritMultiplier(isMH, true)),

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if isMH {
				MHOutcome = spellEffect.Outcome
			} else {
				OHOutcome = spellEffect.Outcome
			}
		},
	}
	if !isMH {
		effect.ProcMask = core.ProcMaskMeleeOHSpecial
		effect.BaseDamage = core.BaseDamageConfigMeleeWeapon(core.OffHand, true, 181, 1+0.1*float64(rogue.Talents.DualWieldSpecialization), 1, false)
	}

	effect.BaseDamage = core.WrapBaseDamageConfig(effect.BaseDamage, func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
		return func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
			normalDamage := oldCalculator(sim, spellEffect, spell)
			// TODO: Add support for all poison effects
			if rogue.DeadlyPoisonDots[spellEffect.Target.Index].IsActive() {
				return normalDamage * 1.2
			} else {
				return normalDamage
			}
		}
	})

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagRogueAbility,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (rogue *Rogue) registerMutilateSpell() {
	mhHitSpell := rogue.newMutilateHitSpell(true)
	ohHitSpell := rogue.newMutilateHitSpell(false)

	baseCost := 60.0
	if rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfMutilate) {
		baseCost -= 5
	}
	refundAmount := baseCost * 0.8

	rogue.Mutilate = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48666},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagBuilder,

		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
			ModifyCast:  rogue.builderCastModifier,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			ThreatMultiplier: 1,
			OutcomeApplier:   rogue.OutcomeFuncMeleeSpecialHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					rogue.AddEnergy(sim, refundAmount, rogue.EnergyRefundMetrics)
					return
				}

				rogue.AddComboPoints(sim, 2, spell.ComboPointMetrics())

				// TODO: while this is the most natural handling, the oh attack might have effects
				//  from the mh attack applied
				mhHitSpell.Cast(sim, spellEffect.Target)
				ohHitSpell.Cast(sim, spellEffect.Target)

				if MHOutcome == core.OutcomeCrit || OHOutcome == core.OutcomeCrit {
					//rogue.Mutilate.ApplyEffects.Outcome = core.OutcomeCrit
					spellEffect.Outcome = core.OutcomeCrit
				}
			},
		}),
	})
}
