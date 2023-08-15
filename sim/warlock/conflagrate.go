package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warlock *Warlock) registerConflagrateSpell() {
	if !warlock.Talents.Conflagrate {
		return
	}

	hasGlyphOfConflag := warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate)
	warlock.Conflagrate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 17962},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

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
			core.TernaryFloat64(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfImmolate), 0.1, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDeathbringerGarb, 2), 0.1, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetGuldansRegalia, 4), 0.1, 0),
		CritMultiplier:   warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Conflagrate",
			},
			NumberOfTicks: 3,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = (314.0 / 3) + (0.4/3)*dot.Spell.SpellPower()
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)

				// DoT does not benefit from firestone and also not from spellstone
				dot.Spell.DamageMultiplierAdditive -= warlock.GrandFirestoneBonus()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				dot.Spell.DamageMultiplierAdditive += warlock.GrandFirestoneBonus()
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// takes the SP of the immolate (or shadowflame) dot on the target
			baseDamage := 471.0 + 0.6*warlock.Immolate.Dot(target).Spell.SpellPower()
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
