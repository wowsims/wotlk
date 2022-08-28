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
				Duration: time.Second * 4,
			},
			OnCastComplete: func(sim *core.Simulation, spell *core.Spell) {
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
				CastTime: time.Second * 2,
				Cost:     manaCost,
				GCD:      core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    fireElemental.NewTimer(),
				Duration: time.Second * 4,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				fireElemental.AutoAttacks.DelayMeleeUntil(sim, sim.CurrentTime+fireElemental.AutoAttacks.MainhandSwingSpeed())
			},
			OnCastComplete: func(sim *core.Simulation, spell *core.Spell) {
				fireElemental.AutoAttacks.DelayMeleeUntil(sim, sim.CurrentTime+fireElemental.AutoAttacks.MainhandSwingSpeed())
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagic(1, 150, 1.071), // TODO these are approximation, from base SP
			OutcomeApplier:   fireElemental.OutcomeFuncMagicHitAndCrit(fireElemental.DefaultSpellCritMultiplier()),
		}),
	})

}

func (fireElemental *FireElemental) registerFireShieldDot() {
	actionID := core.ActionID{SpellID: 11350}

	//dummy spell, dots require a spell
	fireShieldSpell := fireElemental.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
		},
	})

	//TODO Will need to account for mutliple targets
	target := fireElemental.CurrentTarget

	//TODO Dont no what the best way to handle this is.
	fireElemental.FireShieldDot = core.NewDot(core.Dot{
		Spell: fireShieldSpell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Fire Shield",
			ActionID: actionID,
		}),
		NumberOfTicks: 40,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagic(68, 70, 0.032), // TODO these are approximation, from base SP
			OutcomeApplier:   fireElemental.OutcomeFuncMagicCrit(fireElemental.DefaultSpellCritMultiplier()),
		})),
	})
}
