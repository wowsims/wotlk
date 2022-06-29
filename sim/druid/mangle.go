package druid

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (druid *Druid) registerMangleBearSpell() {
	if !druid.Talents.Mangle {
		return
	}

	druid.MangleAura = core.MangleAura(druid.CurrentTarget)

	cost := 20.0 - float64(druid.Talents.Ferocity)
	refundAmount := cost * 0.8

	druid.Mangle = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 33987},
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
				Duration: time.Second * 6,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHSpecial,

			DamageMultiplier: 1,
			ThreatMultiplier: (1.5 / 1.15) *
				core.TernaryFloat64(druid.InForm(Bear) && ItemSetThunderheartHarness.CharacterHasSetBonus(&druid.Character, 2), 1.15, 1),

			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 155/1.15, 1.15, true),
			OutcomeApplier: druid.OutcomeFuncMeleeSpecialHitAndCrit(druid.MeleeCritMultiplier()),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.MangleAura.Activate(sim)
				} else {
					druid.AddRage(sim, refundAmount, druid.RageRefundMetrics)
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

	cost := 45.0 - float64(druid.Talents.Ferocity) - core.TernaryFloat64(ItemSetThunderheartHarness.CharacterHasSetBonus(&druid.Character, 2), 5.0, 0)
	refundAmount := cost * 0.8

	druid.Mangle = druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 33983},
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

			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 264/1.6, 1.6, true),
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
