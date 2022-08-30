package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerAimedShotSpell(timer *core.Timer) {
	if !hunter.Talents.AimedShot {
		return
	}
	baseCost := 0.08 * hunter.BaseMana

	hunter.AimedShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49050},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.03*float64(hunter.Talents.Efficiency)) *
					(1 - 0.05*float64(hunter.Talents.MasterMarksman)),
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second*10 - core.TernaryDuration(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfAimedShot), time.Second*2, 0),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskRangedSpecial,
			BonusCritRating: 4*core.CritRatingPerCritChance*float64(hunter.Talents.ImprovedBarrage) +
				core.TernaryFloat64(hunter.Talents.TrueshotAura && hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfTrueshotAura), 10*core.CritRatingPerCritChance, 0),
			DamageMultiplier: 1 *
				hunter.markedForDeathMultiplier(),
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (hitEffect.RangedAttackPower(spell.Unit)+hitEffect.RangedAttackPowerOnTarget())*0.2 +
						hunter.AutoAttacks.Ranged.BaseDamage(sim) +
						hunter.AmmoDamageBonus +
						hitEffect.BonusWeaponDamage(spell.Unit) +
						408
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: hunter.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(true, true, hunter.CurrentTarget)),
		}),

		InitialDamageMultiplier: 1 +
			.04*float64(hunter.Talents.Barrage),
	})
}
