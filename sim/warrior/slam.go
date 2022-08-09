package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerSlamSpell() {
	cost := 15.0 - float64(warrior.Talents.FocusedRage)
	refundAmount := cost * 0.8

	warrior.Slam = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47475},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     cost,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*1500 - time.Millisecond*500*time.Duration(warrior.Talents.ImprovedSlam),
			},
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHSpecial,

			DamageMultiplier: 1 + 0.02*float64(warrior.Talents.UnendingFury),
			ThreatMultiplier: 1,
			FlatThreatBonus:  70,

			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 250, 1, 1, true),
			OutcomeApplier: warrior.OutcomeFuncMeleeWeaponSpecialHitAndCrit(warrior.critMultiplier(true)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
}
func (warrior *Warrior) CanSlam(sim *core.Simulation) bool {
	normalCastTime := warrior.Slam.DefaultCast.CastTime
	if warrior.BloodsurgeAura.IsActive() {
		warrior.Slam.DefaultCast.CastTime = 0
	} else {
		warrior.Slam.DefaultCast.CastTime = normalCastTime
	}

	return warrior.CurrentRage() >= warrior.Slam.DefaultCast.Cost && warrior.Slam.IsReady(sim) && (warrior.Talents.ImprovedSlam >= 1 || warrior.BloodsurgeAura.IsActive())
}
