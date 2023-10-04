package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (priest *Priest) registerShadowWordPainSpell() {
	twistedFaithMultiplier := 1 + 0.02*float64(priest.Talents.TwistedFaith)
	mentalAgility := []float64{0, .04, .07, .10}[priest.Talents.MentalAgility]
	shadowFocus := 0.02 * float64(priest.Talents.ShadowFocus)
	mindFlayMod := twistedFaithMultiplier +
		core.TernaryFloat64(priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfMindFlay)), 0.1, 0)

	var glyphManaMetrics *core.ResourceMetrics
	if priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadowWordPain)) {
		glyphManaMetrics = priest.NewManaMetrics(core.ActionID{SpellID: 56172})
	}

	priest.ShadowWordPain = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48125},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.22,
			Multiplier: 1 - (shadowFocus + mentalAgility),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusHitRating:  float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating: float64(priest.Talents.MindMelt)*3*core.CritRatingPerCritChance + core.TernaryFloat64(priest.HasSetBonus(ItemSetCrimsonAcolyte, 2), 5, 0)*core.CritRatingPerCritChance,
		DamageMultiplier: 1 +
			float64(priest.Talents.Darkness)*0.02 +
			float64(priest.Talents.TwinDisciplines)*0.01 +
			float64(priest.Talents.ImprovedShadowWordPain)*0.03,
		CritMultiplier:   priest.SpellCritMultiplier(1, 1),
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "ShadowWordPain",
				OnGain: func(_ *core.Aura, _ *core.Simulation) {
					priest.MindBlastModifier *= twistedFaithMultiplier
					priest.MindFlayModifier *= mindFlayMod
				},
				OnExpire: func(_ *core.Aura, _ *core.Simulation) {
					priest.MindBlastModifier /= twistedFaithMultiplier
					priest.MindFlayModifier /= mindFlayMod
				},
			},

			NumberOfTicks: 6 +
				core.TernaryInt32(priest.HasSetBonus(ItemSetAbsolution, 2), 1, 0),
			TickLength: time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = 1380/6 + 0.1833*dot.Spell.SpellPower()
				if !isRollover {
					dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if priest.Talents.Shadowform {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				}

				if glyphManaMetrics != nil {
					priest.AddMana(sim, priest.BaseMana*0.01, glyphManaMetrics)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				priest.AddShadowWeavingStack(sim)
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				if priest.Talents.Shadowform {
					return dot.CalcSnapshotDamage(sim, target, dot.OutcomeExpectedMagicSnapshotCrit)
				} else {
					return dot.CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicAlwaysHit)
				}
			} else {
				baseDamage := 1380/6 + 0.1833*spell.SpellPower()
				if priest.Talents.Shadowform {
					return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
				} else {
					return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
				}
			}
		},
	})
}
