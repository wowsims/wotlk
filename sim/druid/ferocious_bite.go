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
	refundPercent := 0.4 * float64(druid.Talents.PrimalPrecision)
	dmgPerComboPoint := 290.0 + core.TernaryFloat64(druid.Equip[items.ItemSlotRanged].ID == 25667, 14, 0)

	var excessEnergy float64
	var refundAmount float64

	biteBaseBonusCrit := core.TernaryFloat64(druid.HasT9FeralSetBonus(4), 5*core.CritRatingPerCritChance, 0.0)
	if druid.AssumeBleedActive {
		biteBaseBonusCrit += (5 * float64(druid.Talents.RendAndTear)) * core.CritRatingPerCritChance
	}

	druid.FerociousBite = druid.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				// Bite refunds based on base cost, which can change based on berserk
				// It also won't spend 'excess' energy on miss
				druid.ApplyClearcasting(sim, spell, cast)
				currentCost := druid.CurrentFerociousBiteCost()
				// fixup currentCost, will account for berserk but not clearcasting
				if cast.Cost == 0 {
					currentCost = 0
				}
				if refundPercent > 0.0 {
					refundAmount = currentCost * refundPercent
				}

				excessEnergy = core.MinFloat(spell.Unit.CurrentEnergy()-currentCost, 30)
			},
		},

		BonusCritRating: biteBaseBonusCrit,
		DamageMultiplier: (1 + 0.03*float64(druid.Talents.FeralAggression)) *
			core.TernaryFloat64(druid.HasSetBonus(ItemSetThunderheartHarness, 4), 1.15, 1.0),
		CritMultiplier:   druid.MeleeCritMultiplier(),
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
			},
			OutcomeApplier: druid.OutcomeFuncMeleeSpecialHitAndCrit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.SpendEnergy(sim, excessEnergy, spell.ResourceMetrics)
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
