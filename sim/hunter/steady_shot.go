package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerSteadyShotSpell() {
	baseCost := 0.05 * hunter.BaseMana()

	impSSProcChance := 0.05 * float64(hunter.Talents.ImprovedSteadyShot)
	if hunter.Talents.ImprovedSteadyShot > 0 {
		hunter.ImprovedSteadyShotAura = hunter.RegisterAura(core.Aura{
			Label:    "Improved Steady Shot",
			ActionID: core.ActionID{SpellID: 53220},
			Duration: time.Second * 12,
		})
	}

	hunter.SteadyShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49052},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.03*float64(hunter.Talents.Efficiency)) *
					(1 - 0.05*float64(hunter.Talents.MasterMarksman)),
				GCD:      core.GCDDefault + hunter.latency,
				CastTime: 1, // Dummy value so core doesn't optimize the cast away
			},
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				cast.CastTime = hunter.SteadyShotCastTime()
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskRangedSpecial,

			BonusCritRating: 0 +
				2*core.CritRatingPerCritChance*float64(hunter.Talents.SurvivalInstincts) +
				core.TernaryFloat64(ItemSetRiftStalker.CharacterHasSetBonus(&hunter.Character, 4), 5*core.CritRatingPerCritChance, 0),
			DamageMultiplier: 1 *
				(1 + 0.03*float64(hunter.Talents.FerociousInspiration)) *
				(1 + 0.01*float64(hunter.Talents.MarkedForDeath)) *
				hunter.sniperTrainingMultiplier() *
				core.TernaryFloat64(ItemSetGronnstalker.CharacterHasSetBonus(&hunter.Character, 4), 1.1, 1),
			ThreatMultiplier: 1,

			BaseDamage: core.WrapBaseDamageConfig(
				core.BaseDamageConfig{
					Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
						return (hitEffect.RangedAttackPower(spell.Unit)+hitEffect.RangedAttackPowerOnTarget())*0.2 +
							hunter.AutoAttacks.Ranged.BaseDamage(sim)*2.8/hunter.AutoAttacks.Ranged.SwingSpeed +
							252
					},
					TargetSpellCoefficient: 1,
				},
				func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
					if hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfSteadyShot) {
						return func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
							normalDamage := oldCalculator(sim, hitEffect, spell)
							// TODO: Other hunters' stings should be allowed also
							if hunter.SerpentStingDot.IsActive() {
								return normalDamage * 1.1
							} else {
								return normalDamage
							}
						}
					} else {
						return oldCalculator
					}
				}),
			OutcomeApplier: hunter.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(true, true, hunter.CurrentTarget)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				hunter.rotation(sim, false)

				if spellEffect.Landed() && impSSProcChance > 0 && sim.RandomFloat("Imp Steady Shot") < impSSProcChance {
					hunter.ImprovedSteadyShotAura.Activate(sim)
				}
			},
		}),
	})
}

func (hunter *Hunter) SteadyShotCastTime() time.Duration {
	return time.Duration(float64(time.Millisecond*1500)/hunter.RangedSwingSpeed()) + hunter.latency
}
