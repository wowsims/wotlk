package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (rogue *Rogue) registerBackstabSpell() {
	baseCost := 60.0
	refundAmount := baseCost * 0.8

	rogue.Backstab = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 26863},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagBuilder,

		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
			ModifyCast:  rogue.CastModifier,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHSpecial,
			BonusCritRating: core.TernaryFloat64(rogue.HasSetBonus(ItemSetVanCleefs, 4), 5*core.CritRatingPerCritChance, 0) +
				10*core.CritRatingPerCritChance*float64(rogue.Talents.PuncturingWounds),
			// All of these use "Apply Aura: Modifies Damage/Healing Done", and stack additively (up to 142%).
			DamageMultiplier: 1 +
				0.1*float64(rogue.Talents.Opportunity) +
				0.02*float64(rogue.Talents.Aggression) +
				core.TernaryFloat64(rogue.Talents.SurpriseAttacks, 0.1, 0) +
				core.TernaryFloat64(rogue.HasSetBonus(ItemSetSlayers, 4), 0.06, 0),
			ThreatMultiplier: 1,

			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, true, 170, 1.0, 1.5+0.01*float64(rogue.Talents.SinisterCalling), true),
			OutcomeApplier: rogue.OutcomeFuncMeleeSpecialHitAndCrit(rogue.MeleeCritMultiplier(true, true)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				} else {
					rogue.AddEnergy(sim, refundAmount, rogue.EnergyRefundMetrics)
				}
			},
		}),
	})
}
