package rogue

import (
	"time"

	"github.com/wowsims/sim/core"
	"github.com/wowsims/sim/core/stats"
)

func (rogue) registerAmbushSpell() {
	baseCost := rogue.costModifier(60 - 4*float64(rogue.Talents.SlaughterFromTheShadows))
	refundAmount := baseCost * 0.8

	rogue.Ambush = rogue.RegisterSpell(core.SpellConfig{
		ActionID:	core.ActionID{SpellID: 48691},
		SpellSchool:	core.SpellSchoolPhysical,
		ProcMask:		core.ProcMaskMeleeMHSpecial,
		Flags:			core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder,
		ResourceType:	stats.Energy,
		BaseCost:		baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
		},

		BonusCritRating: core.TernaryFloat64(rogue.HasSetBonus(ItemSetVanCleefs, 4), 5*core.CritRatingPerCritChance, 0) +
			[]float64{0, 2, 4, 6}[rogue.Talents.TurnTheTables]*core.CritRatingPerCritChance +
			25*core.CritRatingPerCritChance*float64(rogue.Talents.ImprovedAmbush),
		// All of these use "Apply Aura: Modifies Damage/Healing Done", and stack additively (up to 142%).
		DamageMultiplier: 2.75 * (1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			0.1*float64(rogue.Talents.Opportunity)),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 908 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()
			
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				comboPoints = 2 + (1 * 1/3 * float64(rogue.Talents.Initiative))
				rogue.AddComboPoints(sim, comboPoints, spell.ComboPointMetrics())
			} else {
				rogue.AddEnergy(sim, refundAmount, rogue.EnergyRefundMetrics)
			}
		},
	})
}