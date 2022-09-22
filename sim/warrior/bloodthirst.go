package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerBloodthirstSpell(cdTimer *core.Timer) {
	cost := 20.0
	refundAmount := cost * 0.8

	warrior.Bloodthirst = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 23881},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: time.Second * 4,
			},
		},

		BonusCritRating:  core.TernaryFloat64(warrior.HasSetBonus(ItemSetSiegebreakerBattlegear, 4), 10, 0) * core.CritRatingPerCritChance,
		DamageMultiplier: 1 + 0.02*float64(warrior.Talents.UnendingFury) + core.TernaryFloat64(warrior.HasSetBonus(ItemSetOnslaughtBattlegear, 4), 0.05, 0),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return 0.5 * spell.MeleeAttackPower()
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: warrior.OutcomeFuncMeleeSpecialHitAndCrit(warrior.critMultiplier(mh)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
}

func (warrior *Warrior) CanBloodthirst(sim *core.Simulation) bool {
	return warrior.Talents.Bloodthirst && warrior.CurrentRage() >= warrior.Bloodthirst.DefaultCast.Cost && warrior.Bloodthirst.IsReady(sim)
}
