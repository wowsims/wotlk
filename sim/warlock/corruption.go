package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerCorruptionSpell() {
	actionID := core.ActionID{SpellID: 47813}
	spellSchool := core.SpellSchoolShadow
	baseCost := 0.14 * warlock.BaseMana
	spellCoeff := 0.2 + 0.12*float64(warlock.Talents.EmpoweredCorruption)/6 + 0.01*float64(warlock.Talents.EverlastingAffliction)
	canCrit := warlock.Talents.Pandemic

	warlock.Corruption = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:  core.GCDDefault,
			},
		},

		BonusCritRating: 0 +
			warlock.masterDemonologistShadowCrit() +
			3*float64(warlock.Talents.Malediction)*core.CritRatingPerCritChance +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDarkCovensRegalia, 2), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, true),
		CritMultiplier:           warlock.SpellCritMultiplier(1, 1),
		ThreatMultiplier:         1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),

		// TODO: The application of the dot here is counting as a hit for 0 damage (not crit)
		// This messes with final dmg and crit rate metrics.
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.CorruptionDot.Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})

	warlock.CorruptionDot = core.NewDot(core.Dot{
		Spell: warlock.Corruption,
		Aura: warlock.CurrentTarget.RegisterAura(core.Aura{
			Label:    "Corruption-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks:       6,
		TickLength:          time.Second * 3,
		AffectedByCastSpeed: warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfQuickDecay),

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
			dot.SnapshotBaseDamage = 1080/6 + spellCoeff*dot.Spell.SpellPower()
			if !isRollover {
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
			}
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			if canCrit {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			} else {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			}
		},
	})
}

func applyDotOnLanded(dot **core.Dot) func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	return func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if result.Landed() {
			(*dot).Apply(sim)
		}
	}
}
