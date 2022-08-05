package mage

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerArcaneExplosionSpell() {
	baseCost := .22 * mage.BaseMana

	mage.ArcaneExplosion = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 42921},
		SpellSchool: core.SpellSchoolArcane,
		Flags:       SpellFlagMage,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncAOEDamageCapped(mage.Env, core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellHitRating:  float64(mage.Talents.ArcaneFocus) * 2 * core.SpellHitRatingPerHitChance,
			BonusSpellCritRating: float64(mage.Talents.SpellImpact) * 2 * core.CritRatingPerCritChance,

			DamageMultiplier: mage.spellDamageMultiplier,
			ThreatMultiplier: 1 - 0.2*float64(mage.Talents.ArcaneSubtlety),

			BaseDamage:     core.BaseDamageConfigMagic(538, 582, 0.214),
			OutcomeApplier: mage.OutcomeFuncMagicHitAndCrit(mage.SpellCritMultiplier(1, 0.25*float64(mage.Talents.SpellPower))),
		}),
	})
}
