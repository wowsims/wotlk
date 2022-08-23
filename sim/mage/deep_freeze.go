package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerDeepFreezeSpell() {
	if !mage.Talents.DeepFreeze {
		return
	}

	baseCost := .09 * mage.BaseMana

	mage.DeepFreeze = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 44572},
		SpellSchool: core.SpellSchoolFrost,
		Flags:       SpellFlagMage,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskSpellDamage,

			DamageMultiplier: mage.spellDamageMultiplier,
			ThreatMultiplier: 1 - (0.1/3)*float64(mage.Talents.FrostChanneling),

			BaseDamage:     core.BaseDamageConfigMagic(1469, 1741, 7.5/3.5),
			OutcomeApplier: mage.OutcomeFuncMagicHitAndCrit(mage.SpellCritMultiplier(1, 0.25*float64(mage.Talents.SpellPower)+float64(mage.Talents.IceShards)/3)),
		}),
	})
}
