package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerArcaneShotSpell(timer *core.Timer) {
	baseCost := 0.05 * hunter.BaseMana

	var onSpellHit func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)
	if hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfArcaneShot) {
		manaMetrics := hunter.NewManaMetrics(core.ActionID{ItemID: 42898})
		onSpellHit = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() && (hunter.SerpentStingDot.IsActive() || hunter.ScorpidStingAura.IsActive()) {
				hunter.AddMana(sim, 0.2*hunter.ArcaneShot.DefaultCast.Cost, manaMetrics, false)
			}
		}
	}

	hunter.ArcaneShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49045},
		SpellSchool: core.SpellSchoolArcane,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.03*float64(hunter.Talents.Efficiency)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second*6 - time.Millisecond*200*time.Duration(hunter.Talents.ImprovedArcaneShot),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:        core.ProcMaskRangedSpecial,
			BonusCritRating: 2 * core.CritRatingPerCritChance * float64(hunter.Talents.SurvivalInstincts),
			DamageMultiplier: 1 *
				hunter.markedForDeathMultiplier(),
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (hitEffect.RangedAttackPower(spell.Unit)+hitEffect.RangedAttackPowerOnTarget())*0.15 + 492
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier:  hunter.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(true, true, hunter.CurrentTarget)),
			OnSpellHitDealt: onSpellHit,
		}),

		InitialDamageMultiplier: 1 +
			.03*float64(hunter.Talents.FerociousInspiration) +
			.05*float64(hunter.Talents.ImprovedArcaneShot),
	})
}
