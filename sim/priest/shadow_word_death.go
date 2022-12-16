package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerShadowWordDeathSpell() {
	baseCost := priest.BaseMana * 0.12

	if priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfShadowWordDeath) {
		priest.RegisterResetEffect(func(sim *core.Simulation) {
			sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int) {
				if isExecute == 35 {
					priest.ShadowWordDeath.DamageMultiplier *= 1.1
				}
			})
		})
	}

	priest.ShadowWordDeath = priest.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48158},
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - []float64{0, .04, .07, .10}[priest.Talents.MentalAgility]),
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		BonusHitRating: float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating: 0 +
			float64(priest.Talents.MindMelt)*2*core.CritRatingPerCritChance +
			core.TernaryFloat64(priest.HasSetBonus(ItemSetValorous, 4), 10, 0)*core.CritRatingPerCritChance, // might be 0.1?
		DamageMultiplier: 1 +
			0.02*float64(priest.Talents.Darkness) +
			0.01*float64(priest.Talents.TwinDisciplines),
		CritMultiplier:   priest.SpellCritMultiplier(1, float64(priest.Talents.ShadowPower)/5),
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(750, 870) + 0.429*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				priest.AddShadowWeavingStack(sim)
			}
			if result.DidCrit() && priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadow)) {
				priest.ShadowyInsightAura.Activate(sim)
			}
			if result.DidCrit() && priest.ImprovedSpiritTap != nil {
				priest.ImprovedSpiritTap.Activate(sim)
			}
			spell.DealDamage(sim, result)
		},
		ExpectedDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) *core.SpellResult {
			baseDamage := sim.Roll(750, 870) + 0.429*spell.SpellPower()
			return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicHitAndCrit)
		},
	})
}
