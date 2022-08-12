package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerScorchSpell() {
	baseCost := .08 * mage.BaseMana

	var onSpellHitDealt core.EffectOnSpellHitDealt
	if mage.Talents.ImprovedScorch > 0 {
		mage.ScorchAura = mage.CurrentTarget.GetAura(core.ImprovedScorchAuraLabel)
		if mage.ScorchAura == nil {
			mage.ScorchAura = core.ImprovedScorchAura(mage.CurrentTarget)
		}

		procChance := float64(mage.Talents.ImprovedScorch) / 3.0
		onSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}

			if procChance != 1.0 || sim.RandomFloat("Improved Scorch") > procChance {
				return
			}

			mage.ScorchAura.Activate(sim)
		}
	}

	mage.Scorch = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 42859},
		SpellSchool: core.SpellSchoolFire,
		Flags:       SpellFlagMage,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskSpellDamage,
			BonusSpellHitRating: 0,

			BonusSpellCritRating: 0 +
				float64(mage.Talents.Incineration+mage.Talents.CriticalMass)*2*core.CritRatingPerCritChance +
				float64(mage.Talents.ImprovedScorch)*1*core.CritRatingPerCritChance,

			DamageMultiplier: mage.spellDamageMultiplier *
				(1 + 0.02*float64(mage.Talents.SpellImpact)) *
				core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfScorch), 1.2, 1),
			ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

			BaseDamage:      core.BaseDamageConfigMagic(382, 451, 1.5/3.5),
			OutcomeApplier:  mage.OutcomeFuncMagicHitAndCrit(mage.SpellCritMultiplier(1, mage.bonusCritDamage)),
			OnSpellHitDealt: onSpellHitDealt,
		}),
	})
}
