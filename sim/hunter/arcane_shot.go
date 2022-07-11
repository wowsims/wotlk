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
				GCD:  core.GCDDefault + hunter.latency,
			},
			IgnoreHaste: true,
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				if hunter.LockAndLoadAura.IsActive() {
					cast.Cost = 0
				} else if hunter.ImprovedSteadyShotAura.IsActive() {
					cast.Cost *= 0.8
				}
			},
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second*6 - time.Millisecond*200*time.Duration(hunter.Talents.ImprovedArcaneShot),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:        core.ProcMaskRangedSpecial,
			BonusCritRating: 2 * core.CritRatingPerCritChance * float64(hunter.Talents.SurvivalInstincts),
			DamageMultiplier: 1 *
				(1 + 0.05*float64(hunter.Talents.ImprovedArcaneShot)) *
				(1 + 0.03*float64(hunter.Talents.FerociousInspiration)) *
				(1 + 0.01*float64(hunter.Talents.MarkedForDeath)),
			ThreatMultiplier: 1,

			BaseDamage: hunter.talonOfAlarDamageMod(core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					damage := (hitEffect.RangedAttackPower(spell.Unit)+hitEffect.RangedAttackPowerOnTarget())*0.15 + 492
					if hunter.ImprovedSteadyShotAura.IsActive() {
						damage *= 1.15
						hunter.ImprovedSteadyShotAura.Deactivate(sim)
					}
					return damage
				},
				TargetSpellCoefficient: 1,
			}),
			OutcomeApplier:  hunter.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(true, true, hunter.CurrentTarget)),
			OnSpellHitDealt: onSpellHit,
		}),
	})
}
