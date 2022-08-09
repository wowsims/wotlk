package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerShieldSlamSpell(cdTimer *core.Timer) {
	cost := 20.0 - float64(warrior.Talents.FocusedRage)
	refundAmount := cost * 0.8

	damageRollFunc := core.DamageRollFunc(420, 440)

	warrior.ShieldSlam = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 30356},
		SpellSchool: core.SpellSchoolPhysical,
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
				Duration: time.Second * 6,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHSpecial, // TODO: Is this right?

			DamageMultiplier: 1 * core.TernaryFloat64(warrior.HasSetBonus(ItemSetOnslaughtArmor, 4), 1.1, 1),
			ThreatMultiplier: 1,
			FlatThreatBonus:  305,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, _ *core.SpellEffect, _ *core.Spell) float64 {
					return damageRollFunc(sim) + warrior.GetStat(stats.BlockValue)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: warrior.OutcomeFuncMeleeSpecialHitAndCrit(warrior.critMultiplier(true)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
}

func (warrior *Warrior) HasEnoughRageForShieldSlam() bool {
	return warrior.CurrentRage() >= warrior.ShieldSlam.DefaultCast.Cost
}

func (warrior *Warrior) CanShieldSlam(sim *core.Simulation) bool {
	return warrior.PseudoStats.CanBlock && warrior.HasEnoughRageForShieldSlam() && warrior.ShieldSlam.IsReady(sim)
}
