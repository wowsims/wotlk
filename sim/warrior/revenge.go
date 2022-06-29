package warrior

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warrior *Warrior) registerRevengeSpell(cdTimer *core.Timer) {
	warrior.RegisterAura(core.Aura{
		Label:    "Revenge Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeBlock | core.OutcomeDodge | core.OutcomeParry) {
				warrior.RevengeValidUntil = sim.CurrentTime + time.Second*5
			}
		},
	})

	cost := 5.0 - float64(warrior.Talents.FocusedRage)
	refundAmount := cost * 0.8

	warrior.Revenge = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 30357},
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

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHSpecial,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			FlatThreatBonus:  200,

			BaseDamage:     core.BaseDamageConfigRoll(414, 506),
			OutcomeApplier: warrior.OutcomeFuncMeleeSpecialHitAndCrit(warrior.critMultiplier(true)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				warrior.RevengeValidUntil = 0
				if !spellEffect.Landed() {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
}

func (warrior *Warrior) CanRevenge(sim *core.Simulation) bool {
	return sim.CurrentTime < warrior.RevengeValidUntil &&
		warrior.StanceMatches(DefensiveStance) &&
		warrior.CurrentRage() >= warrior.Revenge.DefaultCast.Cost &&
		warrior.Revenge.IsReady(sim)
}
