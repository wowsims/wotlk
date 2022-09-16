package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerShieldOfRighteousnessSpell() {
	baseCost := paladin.BaseMana * 0.06

	baseModifiers := Multiplicative{}
	baseMultiplier := baseModifiers.Get()

	paladin.ShieldOfRighteousness = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 61411},
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.02*float64(paladin.Talents.Benediction)),
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		// TODO: Why is this here?
		BonusCritRating:  1,
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHSpecial,

			DamageMultiplier: baseMultiplier,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, _ *core.SpellEffect, _ *core.Spell) float64 {
					// TODO: Confirm this is an accurate calculation.
					bv := 2760.0
					if paladin.GetStat(stats.BlockValue) < bv {
						bv = paladin.GetStat(stats.BlockValue)
					}
					return 390 + bv
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: paladin.OutcomeFuncMeleeSpecialHitAndCrit(paladin.MeleeCritMultiplier()),
		}),
	})
}
