package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerBindingHealSpell() {
	baseCost := .27 * priest.BaseMana

	baseEffect := core.SpellEffect{
		IsHealing: true,
		ProcMask:  core.ProcMaskSpellHealing,

		DamageMultiplier: 1 *
			(1 + .02*float64(priest.Talents.DivineProvidence)),

		BaseDamage:     core.BaseDamageConfigHealing(1959, 2516, 0.8057+0.04*float64(priest.Talents.EmpoweredHealing)),
		OutcomeApplier: priest.OutcomeFuncHealingCrit(priest.DefaultHealingCritMultiplier()),
	}

	var effects []core.SpellEffect
	targets := []*core.Unit{&priest.Unit}
	if priest.CurrentTarget != &priest.Unit {
		targets = append(targets, priest.CurrentTarget)
	}
	for _, target := range targets {
		effect := baseEffect
		effect.Target = target
		effects = append(effects, effect)
	}

	priest.BindingHeal = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48120},
		SpellSchool: core.SpellSchoolHoly,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,

				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		BonusCritRating:  float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
		ThreatMultiplier: 0.5 * (1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve]),

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}
