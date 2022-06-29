package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warlock *Warlock) registerIncinerateSpell() {
	baseCost := 355.0
	has4pMal := ItemSetMaleficRaiment.CharacterHasSetBonus(&warlock.Character, 4)

	effect := core.SpellEffect{
		ProcMask:             core.ProcMaskSpellDamage,
		BonusSpellCritRating: float64(warlock.Talents.Devastation) * 1 * core.SpellCritRatingPerCritChance,
		DamageMultiplier:     1 * core.TernaryFloat64(has4pMal, 1.06, 1.0),
		ThreatMultiplier:     1 - 0.05*float64(warlock.Talents.DestructiveReach),
		BaseDamage:           warlock.incinerateDamage(),
		OutcomeApplier:       warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, core.TernaryFloat64(warlock.Talents.Ruin, 1, 0))),
	}

	warlock.Incinerate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 32231},
		SpellSchool: core.SpellSchoolFire,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.01*float64(warlock.Talents.Cataclysm)),
				GCD:  core.GCDDefault,
				// Emberstorm reduces cast time by up to 10%
				CastTime: time.Duration(float64(time.Millisecond*2500) * (1.0 - (0.02 * float64(warlock.Talents.Emberstorm)))),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

}

func (warlock *Warlock) incinerateDamage() core.BaseDamageConfig {
	base := core.BaseDamageConfigMagic(444.0, 514.0, 0.714+0.04*float64(warlock.Talents.ShadowAndFlame))

	return core.WrapBaseDamageConfig(base, func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
		return func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
			normalDamage := oldCalculator(sim, hitEffect, spell)
			// Boost damage if immolate is ticking
			// TODO: in a raid simulator we need to be able to see which dots are ticking from other warlocks.
			if warlock.ImmolateDot.IsActive() { // TODO: use target.getaurabytag(immolatetag)
				return normalDamage + 119.5 //  111 to 128 averages to 119.5
			} else {
				return normalDamage
			}
		}
	})
}
