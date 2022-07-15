package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warlock *Warlock) registerIncinerateSpell() {
	baseCost := 0.14 * warlock.BaseMana
	has4pMal := warlock.HasSetBonus(ItemSetMaleficRaiment, 4)

	effect := core.SpellEffect{
		ProcMask:             core.ProcMaskSpellDamage,
		BonusSpellCritRating: core.TernaryFloat64(warlock.Talents.Devastation, 1, 0) * 5 * core.CritRatingPerCritChance,
		DamageMultiplier:     (1 + 0.06 * core.TernaryFloat64(has4pMal, 1, 0)) * (1 + 0.03*float64(warlock.Talents.Emberstorm)) * 
			(1 + 0.05 * core.TernaryFloat64(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfIncinerate), 1, 0)),
		ThreatMultiplier:     1 - 0.1*float64(warlock.Talents.DestructiveReach),
		BaseDamage:           warlock.incinerateDamage(),
		OutcomeApplier:       warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin) / 5)),
	}

	costReduction := 0.0
	if float64(warlock.Talents.Cataclysm) > 0 {
		costReduction += 0.01 + 0.03*float64(warlock.Talents.Cataclysm)
	}

	warlock.Incinerate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 47838},
		SpellSchool:  core.SpellSchoolFire,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - costReduction),
				GCD:      core.GCDDefault,
				CastTime: warlock.incinerateCastTime(),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

}

func (warlock *Warlock) incinerateCastTime() time.Duration {
	baseCastTime := 2500 - 50*float64(warlock.Talents.Emberstorm)
	if warlock.MoltenCoreAura.IsActive() {
		baseCastTime *= 1.0 - 0.1*float64(warlock.Talents.MoltenCore)
	}
	return (time.Millisecond * time.Duration(baseCastTime))
}

func (warlock *Warlock) incinerateDamage() core.BaseDamageConfig {
	base := core.BaseDamageConfigMagic(582.0, 676.0, 0.713*(1+0.04*float64(warlock.Talents.ShadowAndFlame)))

	return core.WrapBaseDamageConfig(base, func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
		return func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
			normalDamage := oldCalculator(sim, hitEffect, spell)
			// Boost damage if immolate is ticking
			if warlock.ImmolateDot.IsActive() {
				normalDamage += 157 //  145 to 169 averages to 157
				normalDamage *= 1 + 0.02*float64(warlock.Talents.FireAndBrimstone)
			}
			if warlock.MoltenCoreAura.IsActive() {
				normalDamage *= 1 + 0.06*float64(warlock.Talents.MoltenCore)
			}
			return normalDamage
		}
	})
}
