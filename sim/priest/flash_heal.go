package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerFlashHealSpell() {
	baseCost := .18 * priest.BaseMana

	priest.FlashHeal = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48071},
		SpellSchool: core.SpellSchoolHoly,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - .05*float64(priest.Talents.ImprovedFlashHeal)),

				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			IsHealing: true,
			ProcMask:  core.ProcMaskSpellDamage,

			BonusCritRating:  float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
			DamageMultiplier: 1,
			ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

			BaseDamage:     core.BaseDamageConfigHealing(1896, 2203, 0.8057+0.04*float64(priest.Talents.EmpoweredHealing)),
			OutcomeApplier: priest.OutcomeFuncMagicCrit(priest.DefaultSpellCritMultiplier()),
		}),
	})
}
