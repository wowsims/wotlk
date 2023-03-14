package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) registerBackstabSpell() {
	hasGlyph := rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfBackstab)

	rogue.Backstab = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 26863},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | SpellFlagColdBlooded,

		EnergyCost: core.EnergyCostOptions{
			Cost:   rogue.costModifier(60 - 4*float64(rogue.Talents.SlaughterFromTheShadows)),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !rogue.PseudoStats.InFrontOfTarget && rogue.GetMHWeapon().WeaponType == proto.WeaponType_WeaponTypeDagger
		},

		BonusCritRating: core.TernaryFloat64(rogue.HasSetBonus(ItemSetVanCleefs, 4), 5*core.CritRatingPerCritChance, 0) +
			[]float64{0, 2, 4, 6}[rogue.Talents.TurnTheTables]*core.CritRatingPerCritChance +
			10*core.CritRatingPerCritChance*float64(rogue.Talents.PuncturingWounds),
		// All of these use "Apply Aura: Modifies Damage/Healing Done", and stack additively (up to 142%).
		DamageMultiplier: 1.5 * (1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			0.1*float64(rogue.Talents.Opportunity) +
			0.03*float64(rogue.Talents.Aggression) +
			0.05*float64(rogue.Talents.BladeTwisting) +
			core.TernaryFloat64(rogue.Talents.SurpriseAttacks, 0.1, 0) +
			core.TernaryFloat64(rogue.HasSetBonus(ItemSetSlayers, 4), 0.06, 0)) *
			(1 + 0.02*float64(rogue.Talents.SinisterCalling)),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 310 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				// FIXME: Extension of a Rupture Dot can occur up to 3 times
				ruptureDot := rogue.Rupture.Dot(target)
				if hasGlyph && ruptureDot.IsActive() {
					ruptureDot.NumberOfTicks += 1
					ruptureDot.RecomputeAuraDuration()
					ruptureDot.UpdateExpires(ruptureDot.ExpiresAt() + ruptureDot.TickLength)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
