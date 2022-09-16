package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerFlamestrikeSpell() {
	actionID := core.ActionID{SpellID: 42926}
	baseCost := .30 * mage.BaseMana

	applyAOEDamage := core.ApplyEffectFuncAOEDamageCapped(mage.Env, core.SpellEffect{
		ProcMask: core.ProcMaskSpellDamage,

		DamageMultiplier: mage.spellDamageMultiplier * (1 + 0.02*float64(mage.Talents.FirePower)),
		ThreatMultiplier: 1 - 0.05*float64(mage.Talents.BurningSoul),

		BaseDamage:     core.BaseDamageConfigMagic(876, 1071, 0.243),
		OutcomeApplier: mage.fireSpellOutcomeApplier(mage.bonusCritDamage),
	})

	mage.Flamestrike = mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		Flags:       SpellFlagMage,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.01*float64(mage.Talents.Pyromaniac)),

				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},

		BonusCritRating: 0 +
			float64(mage.Talents.CriticalMass)*2*core.CritRatingPerCritChance +
			float64(mage.Talents.Pyromaniac)*1*core.CritRatingPerCritChance,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			applyAOEDamage(sim, target, spell)
			mage.FlamestrikeDot.Apply(sim)
		},
	})

	mage.FlamestrikeDot = core.NewDot(core.Dot{
		Spell: mage.Flamestrike,
		Aura: mage.RegisterAura(core.Aura{
			Label:    "Flamestrike",
			ActionID: actionID,
		}),
		NumberOfTicks: 4,
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncAOESnapshot(mage.Env, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: mage.spellDamageMultiplier * (1 + 0.02*float64(mage.Talents.FirePower)),

			ThreatMultiplier: 1 - 0.05*float64(mage.Talents.BurningSoul),

			BaseDamage:     core.BaseDamageConfigMagicNoRoll(780/4, 0.122),
			OutcomeApplier: mage.OutcomeFuncTick(),
			IsPeriodic:     true,
		}),
	})
}
