package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerChimeraShotSpell() {
	if !hunter.Talents.ChimeraShot {
		return
	}
	baseCost := 0.12 * hunter.BaseMana

	hunter.ChimeraShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53209},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.03*float64(hunter.Talents.Efficiency)) *
					(1 - 0.05*float64(hunter.Talents.MasterMarksman)),
				GCD: core.GCDDefault + hunter.latency,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				if hunter.ImprovedSteadyShotAura.IsActive() {
					cast.Cost *= 0.8
				}
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second*10 - core.TernaryDuration(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfChimeraShot), time.Second*1, 0),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskRangedSpecial,

			DamageMultiplier: 1 + 0.01*float64(hunter.Talents.MarkedForDeath),
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					damage := (hitEffect.RangedAttackPower(spell.Unit)+hitEffect.RangedAttackPowerOnTarget())*0.2 +
						1.2*hunter.AutoAttacks.Ranged.BaseDamage(sim)*2.8/hunter.AutoAttacks.Ranged.SwingSpeed

					if hunter.ImprovedSteadyShotAura.IsActive() {
						damage *= 1.15
						hunter.ImprovedSteadyShotAura.Deactivate(sim)
					}
					return damage
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: hunter.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(true, true, hunter.CurrentTarget)),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}

				if hunter.SerpentStingDot.IsActive() {
					hunter.SerpentStingDot.Refresh(sim)
					// SS has 5 ticks, so 2 ticks is 40%
					hunter.SerpentStingDot.TickOnce()
					hunter.SerpentStingDot.TickOnce()
				} else if hunter.ScorpidStingAura.IsActive() {
					hunter.ScorpidStingAura.Refresh(sim)
				}
			},
		}),
	})
}

func (hunter *Hunter) ChimeraShotCastTime() time.Duration {
	return time.Duration(float64(time.Millisecond*1500)/hunter.RangedSwingSpeed()) + hunter.latency
}
