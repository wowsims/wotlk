package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var MHOutcome = core.OutcomeHit
var OHOutcome = core.OutcomeHit

var MutilateSpellID int32 = 48666

func (rogue *Rogue) newMutilateHitSpell(isMH bool) *core.Spell {
	actionID := core.ActionID{SpellID: 48665}
	procMask := core.ProcMaskMeleeMHSpecial
	if !isMH {
		actionID = core.ActionID{SpellID: 48664}
		procMask = core.ProcMaskMeleeOHSpecial
	}

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagColdBlooded,

		BonusCritRating: core.TernaryFloat64(rogue.HasSetBonus(ItemSetVanCleefs, 4), 5*core.CritRatingPerCritChance, 0) +
			[]float64{0, 2, 4, 6}[rogue.Talents.TurnTheTables]*core.CritRatingPerCritChance +
			5*core.CritRatingPerCritChance*float64(rogue.Talents.PuncturingWounds),

		DamageMultiplierAdditive: 1 +
			0.1*float64(rogue.Talents.Opportunity) +
			0.02*float64(rogue.Talents.FindWeakness) +
			core.TernaryFloat64(rogue.HasSetBonus(ItemSetSlayers, 4), 0.06, 0),
		DamageMultiplier: 1 *
			core.TernaryFloat64(isMH, 1, rogue.dwsMultiplier()),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMH {
				baseDamage = 181 + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			} else {
				baseDamage = 181 + spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			}
			// TODO: Add support for all poison effects
			if rogue.DeadlyPoison.Dot(target).IsActive() || rogue.woundPoisonDebuffAuras[target.Index].IsActive() {
				baseDamage *= 1.2
			}

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)

			if isMH {
				MHOutcome = result.Outcome
			} else {
				OHOutcome = result.Outcome
			}
		},
	})
}

func (rogue *Rogue) registerMutilateSpell() {
	mhHitSpell := rogue.newMutilateHitSpell(true)
	ohHitSpell := rogue.newMutilateHitSpell(false)

	rogue.Mutilate = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: MutilateSpellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagBuilder,

		EnergyCost: core.EnergyCostOptions{
			Cost:   rogue.costModifier(60 - core.TernaryFloat64(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfMutilate), 5, 0)),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit) // Miss/Dodge/Parry/Hit
			if result.Landed() {
				rogue.AddComboPoints(sim, 2, spell.ComboPointMetrics())
				ohHitSpell.Cast(sim, target)
				mhHitSpell.Cast(sim, target)
				if MHOutcome == core.OutcomeCrit || OHOutcome == core.OutcomeCrit {
					result.Outcome = core.OutcomeCrit
				}
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
