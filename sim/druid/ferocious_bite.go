package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (druid *Druid) registerFerociousBiteSpell() {
	dmgPerComboPoint := 290.0 + core.TernaryFloat64(druid.Ranged().ID == 25667, 14, 0)

	druid.FerociousBite = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48577},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:          35,
			Refund:        0.4 * float64(druid.Talents.PrimalPrecision),
			RefundMetrics: druid.PrimalPrecisionRecoveryMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return druid.ComboPoints() > 0
		},

		BonusCritRating: 0 +
			core.TernaryFloat64(druid.HasSetBonus(ItemSetMalfurionsBattlegear, 4), 5*core.CritRatingPerCritChance, 0.0) +
			core.TernaryFloat64(druid.AssumeBleedActive, 5*float64(druid.Talents.RendAndTear)*core.CritRatingPerCritChance, 0),
		DamageMultiplier: (1 + 0.03*float64(druid.Talents.FeralAggression)) *
			core.TernaryFloat64(druid.HasSetBonus(ItemSetThunderheartHarness, 4), 1.15, 1.0),
		CritMultiplier:   druid.MeleeCritMultiplier(Cat),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			comboPoints := float64(druid.ComboPoints())
			attackPower := spell.MeleeAttackPower()
			excessEnergy := min(druid.CurrentEnergy(), 30)

			baseDamage := 120.0 +
				sim.RandomFloat("Ferocious Bite")*140.0 +
				dmgPerComboPoint*comboPoints +
				excessEnergy*(9.4+attackPower/410) +
				attackPower*0.07*comboPoints

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				druid.SpendEnergy(sim, excessEnergy, spell.Cost.(*core.EnergyCost).ResourceMetrics)
				druid.SpendComboPoints(sim, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (druid *Druid) CurrentFerociousBiteCost() float64 {
	return druid.FerociousBite.ApplyCostModifiers(druid.FerociousBite.DefaultCast.Cost)
}
