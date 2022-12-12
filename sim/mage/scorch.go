package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerScorchSpell() {
	baseCost := .08 * mage.BaseMana

	hasImpScorch := mage.Talents.ImprovedScorch > 0
	procChance := float64(mage.Talents.ImprovedScorch) / 3.0
	if hasImpScorch {
		mage.ScorchAura = core.ImprovedScorchAura(mage.CurrentTarget)
	}

	mage.Scorch = mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 42859},
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | HotStreakSpells,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		BonusCritRating: 0 +
			float64(mage.Talents.Incineration+mage.Talents.CriticalMass)*2*core.CritRatingPerCritChance +
			float64(mage.Talents.ImprovedScorch)*1*core.CritRatingPerCritChance,
		DamageMultiplier: mage.spellDamageMultiplier *
			(1 + 0.02*float64(mage.Talents.SpellImpact)) *
			core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfScorch), 1.2, 1),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(382, 451) + (1.5/3.5)*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if hasImpScorch && result.Landed() {
				if sim.Proc(procChance, "Improved Scorch") {
					mage.ScorchAura.Activate(sim)
				}
			}
			spell.DealDamage(sim, result)
		},
	})
}
