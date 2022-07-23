package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerIncinerateSpell() {

	effect := core.SpellEffect{
		ProcMask:             core.ProcMaskSpellDamage,
		BonusSpellCritRating: core.CritRatingPerCritChance * 5 * (core.TernaryFloat64(warlock.Talents.Devastation, 1, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDeathbringerGarb, 4), 1, 0) + core.TernaryFloat64(warlock.HasSetBonus(ItemSetDarkCovensRegalia, 2), 1, 0)),
		DamageMultiplier:	  1,
		ThreatMultiplier: 	  1 - 0.1*float64(warlock.Talents.DestructiveReach),
		BaseDamage:           warlock.incinerateDamage(),
		OutcomeApplier:       warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5)),
	}

	baseCost := 0.14 * warlock.BaseMana
	costReductionFactor := 1.0
	if float64(warlock.Talents.Cataclysm) > 0 {
		costReductionFactor -= 0.01 + 0.03*float64(warlock.Talents.Cataclysm)
	}

	warlock.Incinerate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 47838},
		SpellSchool:  core.SpellSchoolFire,
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
	immoBonus:= 0.0
	if warlock.ImmolateDot.IsActive() {
		immoBonus += 157 //145 + (169-145) * sim.RandomFloat("immoBonus")
	}
	spellCoefficient:= 0.713*(1+0.04*float64(warlock.Talents.ShadowAndFlame))
	return core.BaseDamageConfigMagic(582.0 + immoBonus, 676.0 + immoBonus, spellCoefficient)
}
