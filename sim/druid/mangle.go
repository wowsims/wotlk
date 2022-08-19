package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerMangleBearSpell() {
	if !druid.Talents.Mangle {
		return
	}

	druid.MangleAura = core.MangleAura(druid.CurrentTarget)

	cost := 20.0 - float64(druid.Talents.Ferocity)
	refundAmount := cost * 0.8
	durReduction := (0.5) * float64(druid.Talents.ImprovedMangle)

	druid.Mangle = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48564},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			ModifyCast:  druid.ApplyClearcasting,
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Duration(float64(time.Second) * (6 - durReduction)),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHSpecial,

			DamageMultiplier: 1 + 0.1*float64(druid.Talents.SavageFury),
			ThreatMultiplier: (1.5 / 1.15) *
				core.TernaryFloat64(druid.InForm(Bear) && druid.HasSetBonus(ItemSetThunderheartHarness, 2), 1.15, 1),

			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 299/1.15, 1.0, 1.15, true),
			OutcomeApplier: druid.OutcomeFuncMeleeSpecialHitAndCrit(druid.MeleeCritMultiplier()),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.MangleAura.Activate(sim)
				} else {
					druid.AddRage(sim, refundAmount, druid.RageRefundMetrics)
				}

				if druid.BerserkAura.IsActive() {
					spell.CD.Reset()
				}
			},
		}),
	})
}

func (druid *Druid) registerMangleCatSpell() {
	if !druid.Talents.Mangle {
		return
	}

	druid.MangleAura = core.MangleAura(druid.CurrentTarget)

	cost := 45.0 - (2.0 * float64(druid.Talents.ImprovedMangle)) - float64(druid.Talents.Ferocity) - core.TernaryFloat64(druid.HasSetBonus(ItemSetThunderheartHarness, 2), 5.0, 0)
	refundAmount := cost * 0.8

	druid.Mangle = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48566},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Energy,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  time.Second,
			},
			ModifyCast:  druid.ApplyClearcasting,
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHSpecial,

			DamageMultiplier: 1 + 0.1*float64(druid.Talents.SavageFury),
			ThreatMultiplier: 1,

			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 566/2.0, 1.0, 2.0, true),
			OutcomeApplier: druid.OutcomeFuncMeleeSpecialHitAndCrit(druid.MeleeCritMultiplier()),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())
					druid.MangleAura.Activate(sim)
				} else {
					druid.AddEnergy(sim, refundAmount, druid.EnergyRefundMetrics)
				}
			},
		}),
	})
}

func (druid *Druid) CanMangleBear(sim *core.Simulation) bool {
	return druid.Mangle != nil && druid.CurrentRage() >= druid.Mangle.DefaultCast.Cost && druid.Mangle.IsReady(sim)
}

func (druid *Druid) CanMangleCat() bool {
	return druid.Mangle != nil && druid.CurrentEnergy() >= druid.Mangle.DefaultCast.Cost
}

func (druid *Druid) ShouldMangle(sim *core.Simulation) bool {
	if druid.Mangle == nil {
		return false
	}

	if !druid.Mangle.IsReady(sim) {
		return false
	}

	return druid.CurrentTarget.ShouldRefreshAuraWithTagAtPriority(sim, core.BleedDamageAuraTag, druid.MangleAura.Priority, time.Second*3)
}
