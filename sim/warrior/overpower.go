package warrior

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warrior *Warrior) registerOverpowerSpell(cdTimer *core.Timer) {
	warrior.RegisterAura(core.Aura{
		Label:    "Overpower Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeDodge) {
				warrior.overpowerValidUntil = sim.CurrentTime + time.Second*5
			}
		},
	})

	cost := 5 - float64(warrior.Talents.FocusedRage)
	refundAmount := cost * 0.8

	damageEffect := core.ApplyEffectFuncDirectDamage(core.SpellEffect{
		ProcMask: core.ProcMaskMeleeMHSpecial,

		DamageMultiplier: 1,
		ThreatMultiplier: 0.75,
		BonusCritRating:  25 * core.MeleeCritRatingPerCritChance * float64(warrior.Talents.ImprovedOverpower),

		BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, true, 35, 1, true),
		OutcomeApplier: warrior.OutcomeFuncMeleeSpecialNoBlockDodgeParry(warrior.critMultiplier(true)),

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
			}
		},
	})

	warrior.Overpower = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 11585},
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
				Duration: time.Second * 5,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warrior.overpowerValidUntil = 0
			damageEffect(sim, target, spell)
		},
	})
}

func (warrior *Warrior) ShouldOverpower(sim *core.Simulation) bool {
	return sim.CurrentTime < warrior.overpowerValidUntil &&
		warrior.Overpower.IsReady(sim) &&
		warrior.CurrentRage() >= warrior.Overpower.DefaultCast.Cost
}
