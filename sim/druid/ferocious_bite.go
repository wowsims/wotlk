package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerFerociousBiteSpell() {
	actionID := core.ActionID{SpellID: 48577}
	baseCost := 35.0
	refundAmount := baseCost * (0.4 * float64(druid.Talents.PrimalPrecision))
	dmgPerComboPoint := 290.0 + core.TernaryFloat64(druid.Equip[items.ItemSlotRanged].ID == 25667, 14, 0)

	var excessEnergy float64

	druid.FerociousBite = druid.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics,
		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				druid.ApplyClearcasting(sim, spell, cast)
				excessEnergy = core.MinFloat(spell.Unit.CurrentEnergy()-cast.Cost, 30)
				cast.Cost = baseCost + excessEnergy
			},
		},

		BonusCritRating: core.TernaryFloat64(druid.HasT9FeralSetBonus(4), 5*core.CritRatingPerCritChance, 0.0),
		DamageMultiplier: (1 + 0.03*float64(druid.Talents.FeralAggression)) *
			core.TernaryFloat64(druid.HasSetBonus(ItemSetThunderheartHarness, 4), 1.15, 1.0),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					comboPoints := float64(druid.ComboPoints())

					attackPower := spell.MeleeAttackPower()
					bonusDmg := excessEnergy * (9.4 + attackPower/410)
					base := 120.0 + dmgPerComboPoint*comboPoints + bonusDmg
					roll := sim.RandomFloat("Ferocious Bite") * 140.0
					return base + roll + attackPower*0.07*comboPoints
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: druid.OutcomeFuncMeleeSpecialHitAndCrit(druid.MeleeCritMultiplier()),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.SpendComboPoints(sim, spell.ComboPointMetrics())
				} else if refundAmount > 0 {
					druid.AddEnergy(sim, refundAmount, druid.PrimalPrecisionRecoveryMetrics)
				}
			},
		}),
	})
}

func (druid *Druid) CanFerociousBite() bool {
	return druid.InForm(Cat) && druid.ComboPoints() > 0 && ((druid.CurrentEnergy() >= druid.CurrentFerociousBiteCost()) || druid.ClearcastingAura.IsActive())
}

func (druid *Druid) CurrentFerociousBiteCost() float64 {
	return druid.FerociousBite.ApplyCostModifiers(druid.FerociousBite.BaseCost)
}
