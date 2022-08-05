package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerFrostboltSpell() {
	baseCost := 330.0

	bonusCrit := 0.0
	if mage.MageTier.t9_4 {
		bonusCrit += 5 * core.CritRatingPerCritChance
	}

	mage.Frostbolt = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27072},
		SpellSchool: core.SpellSchoolFrost,
		Flags:       SpellFlagMage | core.SpellFlagBinary | BarrageSpells,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.05*float64(mage.Talents.FrostChanneling)),

				GCD:      core.GCDDefault,
				CastTime: time.Second*3 - time.Millisecond*100*time.Duration(mage.Talents.ImprovedFrostbolt),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskSpellDamage,
			// Frostbolt get 2x bonus from Elemental Precision because it's a binary spell.
			BonusSpellHitRating:  0,
			BonusSpellCritRating: float64(mage.Talents.EmpoweredFrostbolt) * 1 * core.CritRatingPerCritChance,

			DamageMultiplier: mage.spellDamageMultiplier *
				(1 + 0.02*float64(mage.Talents.PiercingIce)) *
				(1 + 0.01*float64(mage.Talents.ArcticWinds)),

			ThreatMultiplier: 1 - (0.1/3)*float64(mage.Talents.FrostChanneling),

			BaseDamage:     core.BaseDamageConfigMagic(600, 647, (3.0/3.5)*0.95+0.02*float64(mage.Talents.EmpoweredFrostbolt)),
			OutcomeApplier: mage.OutcomeFuncMagicHitAndCritBinary(mage.SpellCritMultiplier(1, 0.25*float64(mage.Talents.SpellPower)+0.2*float64(mage.Talents.IceShards))),
		}),
	})
}
