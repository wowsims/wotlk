package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerFireBlastSpell() {
	baseCost := 0.21 * mage.BaseMana

	mage.FireBlast = mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 42873},
		SpellSchool:  core.SpellSchoolFire,
		Flags:        SpellFlagMage | HotStreakSpells,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second*8 - time.Second*time.Duration(mage.Talents.ImprovedFireBlast),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskSpellDamage,

			BonusCritRating: 0 +
				float64(mage.Talents.CriticalMass+mage.Talents.Incineration)*2*core.CritRatingPerCritChance,

			DamageMultiplier: mage.spellDamageMultiplier * (1 + 0.02*float64(mage.Talents.SpellImpact)),

			ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

			BaseDamage:     core.BaseDamageConfigMagic(925, 1095, 1.5/3.5),
			OutcomeApplier: mage.fireSpellOutcomeApplier(mage.bonusCritDamage),
		}),
	})
}
