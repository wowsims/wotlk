package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerShadowWordPainSpell() {
	actionID := core.ActionID{SpellID: 48125}
	baseCost := priest.BaseMana * 0.22

	glyphManaMetric := priest.NewManaMetrics(core.ActionID{SpellID: 56172})
	applier := priest.OutcomeFuncTick()
	if priest.Talents.Shadowform {
		applier = priest.OutcomeFuncMagicCrit(priest.SpellCritMultiplier(1, 1))
	}

	priest.ShadowWordPain = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(priest.Talents.MentalAgility)),
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskSpellDamage,
			BonusSpellHitRating: float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
			ThreatMultiplier:    1 - 0.08*float64(priest.Talents.ShadowAffinity),
			OutcomeApplier:      priest.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					priest.AddShadowWeavingStack(sim)
					priest.ShadowWordPainDot.Apply(sim)
				}
			},
		}),
	})

	target := priest.CurrentTarget
	priest.ShadowWordPainDot = core.NewDot(core.Dot{
		Spell: priest.ShadowWordPain,
		Aura: target.RegisterAura(core.Aura{
			Label:    "ShadowWordPain-" + strconv.Itoa(int(priest.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks: 6 +
			core.TernaryInt(priest.HasSetBonus(ItemSetAbsolution, 2), 1, 0),
		TickLength: time.Second * 3,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask: core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 +
				float64(priest.Talents.Darkness)*0.02 +
				float64(priest.Talents.TwinDisciplines)*0.01 +
				float64(priest.Talents.ImprovedShadowWordPain)*0.03,

			BonusSpellCritRating: float64(priest.Talents.MindMelt)*3*core.CritRatingPerCritChance + core.TernaryFloat64(priest.HasSetBonus(ItemSetCrimsonAcolyte, 2), 5, 0)*core.CritRatingPerCritChance,

			ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),
			IsPeriodic:       true,
			BaseDamage: core.WrapBaseDamageConfig(
				core.BaseDamageConfigMagicNoRoll(1380/6, 0.1833),
				func(oldCalculator core.BaseDamageCalculator) core.BaseDamageCalculator {
					return func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
						swMod := 1 + float64(priest.ShadowWeavingAura.GetStacks())*0.02
						dmg := oldCalculator(sim, spellEffect, spell)

						return dmg * swMod
					}
				}),
			OnPeriodicDamageDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadowWordPain)) {
					priest.AddMana(sim, priest.BaseMana*0.01, glyphManaMetric, false)
				}
			},
			OutcomeApplier: applier,
		}),
	})
}
