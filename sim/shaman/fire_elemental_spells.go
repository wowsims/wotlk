package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (fireElemental *FireElemental) registerFireBlast() {

	// TODO : needs to be figured out
	manaCost := float64(120)

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
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagic(925, 1095, 1.5/3.5), // TODO need proper values just copied from mages fireblast
			OutcomeApplier:   fireElemental.OutcomeFuncMagicHitAndCrit(fireElemental.DefaultSpellCritMultiplier()),
		}),
	})

}

func (fireElemental *FireElemental) registerFireNova() {

	// TODO : needs to be figured out
	manaCost := float64(30)

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
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagic(925, 1095, 1.5/3.5), // TODO need proper values just copied from fireblast
			OutcomeApplier:   fireElemental.OutcomeFuncMagicHitAndCrit(fireElemental.DefaultSpellCritMultiplier()),
		}),
	})

}

func (fireElemental *FireElemental) registerFireShieldAura() {
	actionID := core.ActionID{SpellID: 11350}

	//dummy spell, dots require a spell
	fireShieldSpell := fireElemental.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
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
		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncAOEDamageCapped(fireElemental.Env, core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(371, 0.1), // TODO need proper values
			OutcomeApplier:   fireElemental.OutcomeFuncMagicHitBinary(),
		})),
	})
}
