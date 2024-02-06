package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (rogue *Rogue) registerBackstabSpell() {
	flatDamageBonus := map[int32]int32{
		25: 32,
		40: 60,
		50: 90,
		60: 140,
	}[rogue.Level]

	spellID := map[int32]int32{
		25: 2590,
		40: 8721,
		50: 11279,
		60: 11281,
	}[rogue.Level]

	rogue.Backstab = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | SpellFlagColdBlooded | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   rogue.costModifier(60),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !rogue.PseudoStats.InFrontOfTarget && rogue.HasDagger(core.MainHand)
		},

		BonusCritRating: core.TernaryFloat64(rogue.HasSetBonus(Tier9, 4), 5*core.CritRatingPerCritChance, 0) +
			[]float64{0, 2, 4, 6}[rogue.Talents.TurnTheTables]*core.CritRatingPerCritChance +
			10*core.CritRatingPerCritChance*float64(rogue.Talents.PuncturingWounds),
		// All of these use "Apply Aura: Modifies Damage/Healing Done", and stack additively (up to 142%).
		DamageMultiplier: 1.5 * (1 +
			0.04*float64(rogue.Talents.Opportunity) +
			0.02*float64(rogue.Talents.Aggression)),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := float64(flatDamageBonus) +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
