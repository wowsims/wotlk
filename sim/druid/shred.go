package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerShredSpell() {
	baseCost := 60.0 - 9*float64(druid.Talents.ShreddingAttacks)
	refundAmount := baseCost * 0.8

	flatDamageBonus := 405 +
		core.TernaryFloat64(druid.HasSetBonus(ItemSetNordrassilHarness, 4), 75, 0) +
		core.TernaryFloat64(druid.Equip[items.ItemSlotRanged].ID == 29390, 88, 0)

	druid.Shred = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27002},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			ModifyCast:  druid.ApplyClearcasting,
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.WrapBaseDamageConfig(
				core.BaseDamageConfigMeleeWeapon(core.MainHand, false, flatDamageBonus/2.25, 1.0, 2.25, true),
				func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
					return func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
						normalDamage := oldCalculator(sim, spellEffect, spell)
						modifier := 1.0
						if druid.CurrentTarget.HasActiveAuraWithTag(core.BleedDamageAuraTag) {
							modifier += .3
						}
						if druid.RipDot.IsActive() || druid.RakeDot.IsActive() || druid.LacerateDot.IsActive() {
							modifier += (0.04 * float64(druid.Talents.RendAndTear))
						}

						return normalDamage * modifier
					}
				}),
			OutcomeApplier: druid.OutcomeFuncMeleeSpecialHitAndCrit(druid.MeleeCritMultiplier()),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				} else {
					druid.AddEnergy(sim, refundAmount, druid.EnergyRefundMetrics)
				}
			},
		}),
	})
}

func (druid *Druid) CanShred() bool {
	return !druid.PseudoStats.InFrontOfTarget && druid.CurrentEnergy() >= druid.Shred.DefaultCast.Cost
}
