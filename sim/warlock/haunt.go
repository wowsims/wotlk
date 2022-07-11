package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerHauntSpell() {
	warlock.HauntAura = warlock.RegisterAura(core.Aura{
		Label:     "Haunt Buff",
		ActionID:  core.ActionID{SpellID: 59164},
		Duration:  time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.PeriodicShadowDamageDealtMultiplier *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.PeriodicShadowDamageDealtMultiplier /= 1.2
		},
	})

	effect := core.SpellEffect{
		ProcMask:             core.ProcMaskSpellDamage,
		ThreatMultiplier: 	  1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		DamageMultiplier: 	  1,
		BaseDamage:           core.BaseDamageConfigMagic(645.0, 753.0, 0.4286),
		OutcomeApplier:       warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, core.TernaryFloat64(warlock.Talents.Pandemic, 0, 1))),
		OnSpellHitDealt:  	  func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			// Everlasting Affliction Refresh
			if warlock.CorruptionDot.IsActive() {
				if sim.RandomFloat("EverlastingAffliction") < 0.2 * float64(warlock.Talents.EverlastingAffliction) {
					 warlock.CorruptionDot.Refresh(sim)
				}
			}
			warlock.HauntAura.Activate(sim)
		},
	}

	baseCost := 0.12 * warlock.BaseMana
	warlock.Haunt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 59164},
		SpellSchool: core.SpellSchoolShadow,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*1500,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

}
