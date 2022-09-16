package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerSoulFireSpell() {
	actionID := core.ActionID{SpellID: 47825}
	spellSchool := core.SpellSchoolFire
	baseAdditiveMultiplier := warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false)

	effect := core.SpellEffect{
		ProcMask: core.ProcMaskSpellDamage,

		DamageMultiplier: baseAdditiveMultiplier,

		BaseDamage:     core.BaseDamageConfigMagic(1323.0, 1657.0, 1.15),
		OutcomeApplier: warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5)),
	}

	baseCost := 0.09 * warlock.BaseMana
	costReductionFactor := 1.0
	if float64(warlock.Talents.Cataclysm) > 0 {
		costReductionFactor -= 0.01 + 0.03*float64(warlock.Talents.Cataclysm)
	}

	warlock.SoulFire = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * costReductionFactor,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * time.Duration(6000-400*warlock.Talents.Bane),
			},
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				cast.GCD = time.Duration(float64(cast.GCD) * warlock.backdraftModifier())
				cast.CastTime = time.Duration(float64(cast.CastTime) * warlock.backdraftModifier() * warlock.soulFireCastTime())
			},
		},

		BonusCritRating: 0 +
			warlock.masterDemonologistFireCrit() +
			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDarkCovensRegalia, 2), 5*core.CritRatingPerCritChance, 0),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (warlock *Warlock) soulFireCastTime() float64 {
	castTimeModifier := 1.0
	if warlock.DecimationAura.IsActive() {
		castTimeModifier *= 1.0 - 0.2*float64(warlock.Talents.Decimation)
	}
	return castTimeModifier
}
