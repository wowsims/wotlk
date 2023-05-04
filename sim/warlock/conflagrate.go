package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warlock *Warlock) registerConflagrateSpell() {
	hasGlyphOfConflag := warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate)

	directFlatDamage := 0.6 * 785 / 5 * float64(warlock.Immolate.CurDot().NumberOfTicks)
	directSpellCoeff := 0.6 * 0.2 * float64(warlock.Immolate.CurDot().NumberOfTicks)
	dotFlatDamage := 0.4 / 3 * 785 / 5 * float64(warlock.Immolate.CurDot().NumberOfTicks)
	dotSpellCoeff := 0.4 / 3 * 0.2 * float64(warlock.Immolate.CurDot().NumberOfTicks)

	bonusPeriodicDamageMultiplier := 0 +
		warlock.GrandSpellstoneBonus() -
		warlock.GrandFirestoneBonus()

	warlock.Conflagrate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 17962},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.16,
			Multiplier: 1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm],
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 10,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.Immolate.Dot(target).IsActive()
		},

		BonusCritRating: 0 +
			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0) +
			5*float64(warlock.Talents.FireAndBrimstone)*core.CritRatingPerCritChance,
		DamageMultiplierAdditive: 1 +
			warlock.GrandFirestoneBonus() +
			0.03*float64(warlock.Talents.Emberstorm) +
			0.03*float64(warlock.Talents.Aftermath) +
			0.1*float64(warlock.Talents.ImprovedImmolate) +
			core.TernaryFloat64(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfImmolate), 0.1, 0),
		CritMultiplier:   warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Conflagrate",
			},
			NumberOfTicks: 3,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = dotFlatDamage + dotSpellCoeff*dot.Spell.SpellPower()
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)

				dot.Spell.DamageMultiplierAdditive += bonusPeriodicDamageMultiplier
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				dot.Spell.DamageMultiplierAdditive -= bonusPeriodicDamageMultiplier
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := directFlatDamage + directSpellCoeff*spell.SpellPower()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if !result.Landed() {
				return
			}

			spell.Dot(target).Apply(sim)

			if !hasGlyphOfConflag {
				warlock.Immolate.Dot(target).Deactivate(sim)
				//warlock.ShadowflameDot.Deactivate(sim)
			}
		},
	})
}
