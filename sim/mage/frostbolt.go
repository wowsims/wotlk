package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerFrostboltSpell() {
	baseCost := .11 * mage.BaseMana

	mage.Frostbolt = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 42842},
		SpellSchool: core.SpellSchoolFrost,
		Flags:       SpellFlagMage | BarrageSpells,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,

				GCD:      core.GCDDefault,
				CastTime: time.Second*3 - time.Millisecond*100*time.Duration(mage.Talents.ImprovedFrostbolt+mage.Talents.EmpoweredFrostbolt),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:       core.ProcMaskSpellDamage,
			BonusHitRating: 0,
			BonusCritRating: 0 +
				core.TernaryFloat64(mage.MageTier.t9_4, 5*core.CritRatingPerCritChance, 0),

			DamageMultiplier: mage.spellDamageMultiplier *
				(1 + .01*float64(mage.Talents.ChilledToTheBone)) *
				core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfFrostbolt), 1.05, 1),

			ThreatMultiplier: 1 - (0.1/3)*float64(mage.Talents.FrostChanneling),

			BaseDamage:     core.BaseDamageConfigMagic(799, 861, (3.0/3.5)*0.95+0.05*float64(mage.Talents.EmpoweredFrostbolt)),
			OutcomeApplier: mage.OutcomeFuncMagicHitAndCritBinary(mage.SpellCritMultiplier(1, 0.25*float64(mage.Talents.SpellPower)+float64(mage.Talents.IceShards)/3)),
		}),
	})
}
