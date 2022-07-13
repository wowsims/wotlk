package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerChaosBoltSpell() {
	effect := core.SpellEffect{
		ProcMask:             core.ProcMaskSpellDamage,
		BonusSpellCritRating: core.TernaryFloat64(warlock.Talents.Devastation, 0, 1) * 5 * core.CritRatingPerCritChance,
		DamageMultiplier:     (1 + 0.03 * float64(warlock.Talents.Emberstorm)) *
			(1 + 0.02 * float64(warlock.Talents.FireAndBrimstone) * core.TernaryFloat64(warlock.ImmolateDot.IsActive(), 0, 1)),
		ThreatMultiplier:     1 - 0.1*float64(warlock.Talents.DestructiveReach),
		BaseDamage:           core.BaseDamageConfigMagic(1429.0, 1813.0, 0.7142*(1+0.04*float64(warlock.Talents.ShadowAndFlame))),
		OutcomeApplier:       warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin) / 5)),
	}

	baseCost := 0.07 * warlock.BaseMana
	costReduction := 0.0
	if float64(warlock.Talents.Cataclysm) > 0 {
		costReduction += 0.01 + 0.03 * float64(warlock.Talents.Cataclysm)
	}
	
	warlock.ChaosBolt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 59172},
		SpellSchool:  core.SpellSchoolFire,
		Flags:        core.SpellFlagIgnoreResists,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - costReduction),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*2500 - (time.Millisecond * 100 * time.Duration(warlock.Talents.Bane)),
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

}
