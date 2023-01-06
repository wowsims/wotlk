package mage

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (mage *Mage) registerArcaneExplosionSpell() {
	baseCost := .22 * mage.BaseMana

	mage.ArcaneExplosion = mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 42921},
		SpellSchool:  core.SpellSchoolArcane,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfArcaneExplosion), .9, 1),
				GCD: core.GCDDefault,
			},
		},

		BonusHitRating:   float64(mage.Talents.ArcaneFocus) * 2 * core.SpellHitRatingPerHitChance,
		BonusCritRating:  float64(mage.Talents.SpellImpact) * 2 * core.CritRatingPerCritChance,
		DamageMultiplier: mage.spellDamageMultiplier,
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.2*float64(mage.Talents.ArcaneSubtlety),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dmgFromSP := (1.5 / 3.5 / 2) * spell.SpellPower()
			for _, aoeTarget := range sim.Encounter.Targets {
				baseDamage := sim.Roll(538, 582) + dmgFromSP
				baseDamage *= sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, &aoeTarget.Unit, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}
