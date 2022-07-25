package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerIncinerateSpell() {
	fireAndBrimstoneBonus := 0.02 * float64(warlock.Talents.FireAndBrimstone)
	actionID := core.ActionID{SpellID: 47838}
	spellSchool := core.SpellSchoolFire
	baseAdditiveMultiplier := warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false)

	effect := core.SpellEffect{
		ProcMask: core.ProcMaskSpellDamage,
		BonusSpellCritRating: core.CritRatingPerCritChance * 5 * (core.TernaryFloat64(warlock.Talents.Devastation, 1, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDeathbringerGarb, 4), 1, 0) + core.TernaryFloat64(warlock.HasSetBonus(ItemSetDarkCovensRegalia, 2), 1, 0)),
		DamageMultiplier: baseAdditiveMultiplier,
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),
		BaseDamage:       warlock.incinerateDamage(),
		OutcomeApplier:   warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5)),
		OnInit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if warlock.ImmolateDot.IsActive() {
				spellEffect.DamageMultiplier = baseAdditiveMultiplier + fireAndBrimstoneBonus
			}
		},
	}

	baseCost := 0.14 * warlock.BaseMana
	costReductionFactor := 1.0
	if float64(warlock.Talents.Cataclysm) > 0 {
		costReductionFactor -= 0.01 + 0.03*float64(warlock.Talents.Cataclysm)
	}

	warlock.Incinerate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * costReductionFactor,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * time.Duration(2500-50*warlock.Talents.Emberstorm),
			},
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				cast.GCD = time.Duration(float64(cast.GCD) * warlock.backdraftModifier())
				cast.CastTime = time.Duration(float64(cast.CastTime) * warlock.moltenCoreIncinerateModifier() * warlock.backdraftModifier())
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

}

func (warlock *Warlock) moltenCoreIncinerateModifier() float64 {
	castTimeModifier := 1.0
	if warlock.MoltenCoreAura.IsActive() {
		castTimeModifier *= (1.0 - 0.1*float64(warlock.Talents.MoltenCore))
	}
	return castTimeModifier
}

func (warlock *Warlock) incinerateDamage() core.BaseDamageConfig {
	spellCoefficient := 0.713 * (1 + 0.04*float64(warlock.Talents.ShadowAndFlame))
	base := core.BaseDamageConfigMagic(582.0, 676.0, spellCoefficient)

	return core.WrapBaseDamageConfig(base, func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
		return func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
			normalDamage := oldCalculator(sim, hitEffect, spell)
			// Boost damage if immolate is ticking
			if warlock.ImmolateDot.IsActive() {
				normalDamage += 157 //  145 to 169 averages to 157
			}
			return normalDamage
		}
	})
}
