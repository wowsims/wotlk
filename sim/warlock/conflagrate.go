package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warlock *Warlock) registerConflagrateSpell() {
	actionID := core.ActionID{SpellID: 17962}
	hasGlyphOfConflag := warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate)

	directFlatDamage := 0.6 * 785 / 5 * float64(warlock.ImmolateDot.NumberOfTicks)
	directSpellCoeff := 0.6 * 0.2 * float64(warlock.ImmolateDot.NumberOfTicks)
	dotFlatDamage := 0.4 / 3 * 785 / 5 * float64(warlock.ImmolateDot.NumberOfTicks)
	dotSpellCoeff := 0.4 / 3 * 0.2 * float64(warlock.ImmolateDot.NumberOfTicks)

	warlock.Conflagrate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
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

		BonusCritRating: 0 +
			warlock.masterDemonologistFireCrit +
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

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if !warlock.ImmolateDot.IsActive() {
				panic("Conflagrate spell is cast while Immolate is not active.")
			}

			baseDamage := directFlatDamage + directSpellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				warlock.ConflagrateDot.Apply(sim)
			}
			spell.DealDamage(sim, result)

			if !hasGlyphOfConflag {
				warlock.ImmolateDot.Deactivate(sim)
				//warlock.ShadowflameDot.Deactivate(sim)
			}
		},
	})

	warlock.ConflagrateDot = core.NewDot(core.Dot{
		Spell: warlock.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskSpellDamage,

			BonusCritRating: warlock.Conflagrate.BonusCritRating,
			DamageMultiplierAdditive: 1 +
				warlock.GrandSpellstoneBonus() +
				0.03*float64(warlock.Talents.Emberstorm) +
				0.03*float64(warlock.Talents.Aftermath) +
				0.1*float64(warlock.Talents.ImprovedImmolate) +
				core.TernaryFloat64(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfImmolate), 0.1, 0),
			CritMultiplier:   warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
			ThreatMultiplier: warlock.Conflagrate.ThreatMultiplier,
		}),
		Aura: warlock.CurrentTarget.RegisterAura(core.Aura{
			Label:    "conflagrate-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 3,
		TickLength:    time.Second * 2,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
			dot.SnapshotBaseDamage = dotFlatDamage + dotSpellCoeff*dot.Spell.SpellPower()
			attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
			dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
		},
	})
}
