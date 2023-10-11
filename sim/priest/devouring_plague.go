package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (priest *Priest) registerDevouringPlagueSpell() {
	actionID := core.ActionID{SpellID: 48300}
	mentalAgility := []float64{0, .04, .07, .10}[priest.Talents.MentalAgility]
	shadowFocus := 0.02 * float64(priest.Talents.ShadowFocus)
	priest.DpInitMultiplier = 8 * 0.1 * float64(priest.Talents.ImprovedDevouringPlague)
	hasGlyphOfShadow := priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadow))

	var impDevouringPlague *core.Spell = nil
	if priest.DpInitMultiplier != 0 {
		impDevouringPlague = priest.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 63675},
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskSuppressedProc, // this can proc things like Talisman of Volatile Power
			Flags:       core.SpellFlagDisease,

			BonusHitRating: float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
			BonusCritRating: 0 +
				3*float64(priest.Talents.MindMelt)*core.CritRatingPerCritChance +
				core.TernaryFloat64(priest.HasSetBonus(ItemSetCrimsonAcolyte, 2), 5, 0)*core.CritRatingPerCritChance,
			DamageMultiplier: 1 +
				0.02*float64(priest.Talents.Darkness) +
				0.01*float64(priest.Talents.TwinDisciplines) +
				0.05*float64(priest.Talents.ImprovedDevouringPlague) +
				core.TernaryFloat64(priest.HasSetBonus(ItemSetConquerorSanct, 2), 0.15, 0),
			CritMultiplier:   priest.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1 - 0.05*float64(priest.Talents.ShadowAffinity),

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := (1376/8 + 0.1849*spell.SpellPower()) * priest.DpInitMultiplier
				result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

				if result.DidCrit() && hasGlyphOfShadow {
					priest.ShadowyInsightAura.Activate(sim)
				}
			},
		})
	}

	priest.DevouringPlague = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagDisease | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.25,
			Multiplier: 1 - (shadowFocus + mentalAgility),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusHitRating: float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating: 0 +
			3*float64(priest.Talents.MindMelt)*core.CritRatingPerCritChance +
			core.TernaryFloat64(priest.HasSetBonus(ItemSetCrimsonAcolyte, 2), 5, 0)*core.CritRatingPerCritChance,
		DamageMultiplier: 1 +
			0.02*float64(priest.Talents.Darkness) +
			0.01*float64(priest.Talents.TwinDisciplines) +
			0.05*float64(priest.Talents.ImprovedDevouringPlague) +
			core.TernaryFloat64(priest.HasSetBonus(ItemSetConquerorSanct, 2), 0.15, 0),
		CritMultiplier:   priest.SpellCritMultiplier(1, 1),
		ThreatMultiplier: 1 - 0.05*float64(priest.Talents.ShadowAffinity),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "DevouringPlague",
			},

			NumberOfTicks:       8,
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: priest.Talents.Shadowform,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = 1376/8 + 0.1849*dot.Spell.SpellPower()
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if priest.Talents.Shadowform {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// calculate first, so that if imp. DP procs Shadowy Insight it doesn't influence the dot damage
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if impDevouringPlague != nil {
				impDevouringPlague.Cast(sim, target)
			}

			spell.DealOutcome(sim, result)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				priest.AddShadowWeavingStack(sim)
				spell.Dot(target).Apply(sim)
			}
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
				baseDamage := 1376/8 + 0.1849*spell.SpellPower()
				if priest.Talents.Shadowform {
					return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
				} else {
					return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
				}
			}
		},
	})
}
