package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerFlashHealSpell() {
	baseCost := .18 * priest.BaseMana
	spellCoeff := 0.8057 + 0.04*float64(priest.Talents.EmpoweredHealing)

	priest.FlashHeal = priest.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48071},
		SpellSchool:  core.SpellSchoolHoly,
		ProcMask:     core.ProcMaskSpellHealing,
		Flags:        core.SpellFlagHelpful,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 -
					.05*float64(priest.Talents.ImprovedFlashHeal) -
					core.TernaryFloat64(priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfFlashHeal), .1, 0)),

				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		BonusCritRating: float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + .02*float64(priest.Talents.SpiritualHealing)) *
			(1 + .01*float64(priest.Talents.BlessedResilience)) *
			(1 + .02*float64(priest.Talents.FocusedPower)),
		CritMultiplier:   priest.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHealing := sim.Roll(1896, 2203) + spellCoeff*spell.HealingPower(target)
			spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)
		},
	})
}
