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

	var glyphManaMetrics *core.ResourceMetrics
	if priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadowWordPain)) {
		glyphManaMetrics = priest.NewManaMetrics(core.ActionID{SpellID: 56172})
	}

	priest.ShadowWordPain = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - []float64{0, .04, .07, .10}[priest.Talents.MentalAgility]),
				GCD:  core.GCDDefault,
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

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				priest.AddShadowWeavingStack(sim)
				priest.ShadowWordPainDot.Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
		ExpectedDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := 1380/6 + 0.1833*spell.SpellPower()
			//baseDamage *= float64(priest.ShadowWordPainDot.NumberOfTicks)

			if priest.Talents.Shadowform {
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
			} else {
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	})

	twistedFaithMultiplier := 1 + 0.02*float64(priest.Talents.TwistedFaith)
	mindFlayMod := twistedFaithMultiplier +
		core.TernaryFloat64(priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfMindFlay)), 0.1, 0)
	priest.ShadowWordPainDot = core.NewDot(core.Dot{
		Spell: priest.ShadowWordPain,
		Aura: priest.CurrentTarget.RegisterAura(core.Aura{
			Label:    "ShadowWordPain-" + strconv.Itoa(int(priest.Index)),
			ActionID: actionID,
			OnGain: func(_ *core.Aura, _ *core.Simulation) {
				priest.MindBlast.DamageMultiplier *= twistedFaithMultiplier
				for _, dot := range priest.MindFlayDot {
					if dot != nil {
						dot.Spell.DamageMultiplier *= mindFlayMod
					}
				}
			},
			OnExpire: func(_ *core.Aura, _ *core.Simulation) {
				priest.MindBlast.DamageMultiplier /= twistedFaithMultiplier
				for _, dot := range priest.MindFlayDot {
					if dot != nil {
						dot.Spell.DamageMultiplier /= mindFlayMod
					}
				}
			},
		}),

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
				priest.AddMana(sim, priest.BaseMana*0.01, glyphManaMetrics, false)
			}
		},
	})
}
