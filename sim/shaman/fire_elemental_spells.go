package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (fireElemental *FireElemental) registerFireBlast() {
	var manaCost float64 = 276

	fireElemental.FireBlast = fireElemental.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 13339},
		SpellSchool:  core.SpellSchoolFire,
		ResourceType: stats.Mana,
		BaseCost:     manaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: manaCost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    fireElemental.NewTimer(),
				Duration: time.Second * 4, // TODO estimated from from log diggig,
			},
			OnCastComplete: func(sim *core.Simulation, _ *core.Spell) {
				fireElemental.AutoAttacks.DelayMeleeUntil(sim, sim.CurrentTime+fireElemental.AutoAttacks.MainhandSwingSpeed())
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagic(323, 459, 0.429), // TODO these are approximation, from base SP
			OutcomeApplier:   fireElemental.OutcomeFuncMagicHitAndCrit(fireElemental.DefaultSpellCritMultiplier()),
		}),
	})

}

func (fireElemental *FireElemental) registerFireNova() {
	var manaCost float64 = 207

	fireElemental.FireNova = fireElemental.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 12470},
		SpellSchool:  core.SpellSchoolFire,
		ResourceType: stats.Mana,
		BaseCost:     manaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     manaCost,
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    fireElemental.NewTimer(),
				Duration: time.Second * 4, // TODO estimated from from log diggig,
			},
			ModifyCast: func(sim *core.Simulation, _ *core.Spell, _ *core.Cast) {
				fireElemental.AutoAttacks.DelayMeleeUntil(sim, sim.CurrentTime+fireElemental.AutoAttacks.MainhandSwingSpeed()*2)
			},
		},

		// TODO is this the right affect should it be Capped?
		ApplyEffects: core.ApplyEffectFuncAOEDamageCapped(fireElemental.Env, core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagic(1, 150, 1.0071), // TODO these are approximation, from base SP
			OutcomeApplier:   fireElemental.OutcomeFuncMagicHitAndCrit(fireElemental.DefaultSpellCritMultiplier()),
		}),
	})

}

func (fireElemental *FireElemental) registerFireShieldAura() {
	actionID := core.ActionID{SpellID: 11350}

	//dummy spell
	spell := fireElemental.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
		},
	})

	target := fireElemental.CurrentTarget

	fireShieldDot := core.NewDot(core.Dot{
		Spell: spell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Fire Shield",
			ActionID: actionID,
		}),
		NumberOfTicks: 40,
		TickLength:    time.Second * 3,

		// TODO is this the right affect should it be Capped?
		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncAOEDamage(fireElemental.Env, core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagic(68, 70, 0.032), // TODO these are approximation, from base SP
			OutcomeApplier:   fireElemental.OutcomeFuncMagicCrit(fireElemental.DefaultSpellCritMultiplier()),
		})),
	})

	fireElemental.FireShieldAura = fireElemental.RegisterAura(core.Aura{
		Label:    "Fire Shield",
		ActionID: actionID,
		Duration: time.Minute * 2,
		OnGain: func(_ *core.Aura, sim *core.Simulation) {
			fireShieldDot.Apply(sim)
		},
	})
}
