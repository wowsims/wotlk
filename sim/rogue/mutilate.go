package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var MHOutcome = core.OutcomeHit
var OHOutcome = core.OutcomeHit

var MutilateSpellID int32 = 48666

func (rogue *Rogue) newMutilateHitSpell(isMH bool) *core.Spell {
	actionID := core.ActionID{SpellID: 48665}
	if !isMH {
		actionID = core.ActionID{SpellID: 48664}
	}

	procMask := core.ProcMaskMeleeMHSpecial
	effect := core.SpellEffect{
		BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, true, 181, false),
		OutcomeApplier: rogue.OutcomeFuncMeleeSpecialCritOnly(), // Crit/Hit, should include Block

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if isMH {
				MHOutcome = spellEffect.Outcome
			} else {
				OHOutcome = spellEffect.Outcome
			}
		},
	}
	if !isMH {
		procMask = core.ProcMaskMeleeOHSpecial
		effect.BaseDamage = core.BaseDamageConfigMeleeWeapon(core.OffHand, true, 181, false)
	}

	effect.BaseDamage = core.WrapBaseDamageConfig(effect.BaseDamage, func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
		return func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
			normalDamage := oldCalculator(sim, spellEffect, spell)
			// TODO: Add support for all poison effects
			if rogue.deadlyPoisonDots[spellEffect.Target.Index].IsActive() || rogue.woundPoisonDebuffAuras[spellEffect.Target.Index].IsActive() {
				return normalDamage * 1.2
			} else {
				return normalDamage
			}
		}
	})

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics,

		BonusCritRating: core.TernaryFloat64(rogue.HasSetBonus(ItemSetVanCleefs, 4), 5*core.CritRatingPerCritChance, 0) +
			5*core.CritRatingPerCritChance*float64(rogue.Talents.PuncturingWounds),
		DamageMultiplierAdditive: 1 +
			0.1*float64(rogue.Talents.Opportunity) +
			0.02*float64(rogue.Talents.FindWeakness) +
			core.TernaryFloat64(rogue.HasSetBonus(ItemSetSlayers, 4), 0.06, 0),
		DamageMultiplier: 1 *
			core.TernaryFloat64(!isMH, 1+0.1*float64(rogue.Talents.DualWieldSpecialization), 1),
		CritMultiplier:   rogue.MeleeCritMultiplier(isMH, true),
		ThreatMultiplier: 1,

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
		ActionID:     core.ActionID{SpellID: MutilateSpellID},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | SpellFlagBuilder,
		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
			ModifyCast:  rogue.CastModifier,
		},

		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			OutcomeApplier: rogue.OutcomeFuncMeleeSpecialHit(), // Miss/Dodge/Parry/Hit
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					rogue.AddEnergy(sim, refundAmount, rogue.EnergyRefundMetrics)
					return
				}

				rogue.AddComboPoints(sim, 2, spell.ComboPointMetrics())
				mhHitSpell.Cast(sim, spellEffect.Target)
				ohHitSpell.Cast(sim, spellEffect.Target)
				if MHOutcome == core.OutcomeCrit || OHOutcome == core.OutcomeCrit {
					spellEffect.Outcome = core.OutcomeCrit
				}
			},
		}),
	})
}
