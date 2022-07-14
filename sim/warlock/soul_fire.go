package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerSoulFireSpell() {
	effect := core.SpellEffect{
		ProcMask:             core.ProcMaskSpellDamage,
		BonusSpellCritRating: core.TernaryFloat64(warlock.Talents.Devastation, 0, 1) * 5 * core.CritRatingPerCritChance,
		DamageMultiplier:     (1 + 0.03*float64(warlock.Talents.Emberstorm)),
		ThreatMultiplier:     1 - 0.1*float64(warlock.Talents.DestructiveReach),
		BaseDamage:           warlock.soulFireDamage(),
		OutcomeApplier:       warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5)),
	}

	baseCost := 0.09 * warlock.BaseMana
	costReduction := 0.0
	if float64(warlock.Talents.Cataclysm) > 0 {
		costReduction += 0.01 + 0.03*float64(warlock.Talents.Cataclysm)
	}

	warlock.SoulFire = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 47825},
		SpellSchool:  core.SpellSchoolFire,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - costReduction),
				GCD:      core.GCDDefault,
				CastTime: warlock.soulFireCastTime(),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

}

func (warlock *Warlock) soulFireCastTime() time.Duration {
	baseCastTime := 6000 - 400*float64(warlock.Talents.Bane)
	if warlock.DecimationAura.IsActive() {
		baseCastTime *= 1.0 - 0.2*float64(warlock.Talents.Decimation)
	}
	return (time.Millisecond * time.Duration(baseCastTime))
}

func (warlock *Warlock) soulFireDamage() core.BaseDamageConfig {
	base := core.BaseDamageConfigMagic(1323.0, 1657.0, 1.15)

	return core.WrapBaseDamageConfig(base, func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
		return func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
			if warlock.MoltenCoreAura.IsActive() {
				hitEffect.BonusSpellCritRating += 5 * float64(warlock.Talents.MoltenCore) * core.CritRatingPerCritChance
			}
			normalDamage := oldCalculator(sim, hitEffect, spell)
			// Boost damage if immolate is ticking
			if warlock.MoltenCoreAura.IsActive() {
				normalDamage *= 1 + 0.06*float64(warlock.Talents.MoltenCore)
			}
			return normalDamage
		}
	})
}
