package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerSoulFireSpell() {
	effect := core.SpellEffect{
		ProcMask:             core.ProcMaskSpellDamage,
		BonusSpellCritRating: core.CritRatingPerCritChance * 5 * (core.TernaryFloat64(warlock.Talents.Devastation, 1, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDarkCovensRegalia, 2), 1, 0)),
		DamageMultiplier:     1,
		ThreatMultiplier:     1 - 0.1*float64(warlock.Talents.DestructiveReach),
		BaseDamage:           core.BaseDamageConfigMagic(1323.0, 1657.0, 1.15),
		OutcomeApplier:       warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5)),
		OnInit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			spellEffect.DamageMultiplier = warlock.spellDamageMultiplierHelper(sim, spell, spellEffect)
		},
	}

	baseCost := 0.09 * warlock.BaseMana
	costReductionFactor := 1.0
	if float64(warlock.Talents.Cataclysm) > 0 {
		costReductionFactor -= 0.01 + 0.03*float64(warlock.Talents.Cataclysm)
	}

	warlock.SoulFire = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 47825},
		SpellSchool:  core.SpellSchoolFire,
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
