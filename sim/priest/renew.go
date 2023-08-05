package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (priest *Priest) registerRenewSpell() {
	actionID := core.ActionID{SpellID: 48068}
	spellCoeff := (1.88 + .05*float64(priest.Talents.EmpoweredRenew)) / 5

	if priest.Talents.EmpoweredRenew > 0 {
		priest.EmpoweredRenew = priest.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 63543},
			SpellSchool: core.SpellSchoolHoly,
			ProcMask:    core.ProcMaskSpellHealing,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

			BonusCritRating: float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
			DamageMultiplier: 1 *
				float64(priest.renewTicks()) *
				priest.renewHealingMultiplier() *
				.05 * float64(priest.Talents.EmpoweredRenew) *
				core.TernaryFloat64(priest.HasSetBonus(ItemSetZabrasRaiment, 4), 1.1, 1),
			CritMultiplier:   priest.DefaultHealingCritMultiplier(),
			ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseHealing := 280 + spellCoeff*spell.HealingPower(target)
				spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)
			},
		})
	}

	priest.Renew = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.17,
			Multiplier: 1 - []float64{0, .04, .07, .10}[priest.Talents.MentalAgility],
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: priest.renewHealingMultiplier(),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		Hot: core.DotConfig{
			Aura: core.Aura{
				Label: "Renew",
			},
			NumberOfTicks: priest.renewTicks(),
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = 280 + spellCoeff*dot.Spell.HealingPower(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SpellMetrics[target.UnitIndex].Hits++
			spell.Hot(target).Apply(sim)

			if priest.EmpoweredRenew != nil {
				priest.EmpoweredRenew.Cast(sim, target)
			}
		},
	})
}

func (priest *Priest) renewTicks() int32 {
	return 5 - core.TernaryInt32(priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfRenew), 1, 0)
}

func (priest *Priest) renewHealingMultiplier() float64 {
	return 1 *
		(1 + .02*float64(priest.Talents.SpiritualHealing)) *
		(1 + .01*float64(priest.Talents.BlessedResilience)) *
		(1 + .02*float64(priest.Talents.FocusedPower)) *
		(1 + .01*float64(priest.Talents.TwinDisciplines)) *
		(1 + .05*float64(priest.Talents.ImprovedRenew)) *
		core.TernaryFloat64(priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfRenew), 1.25, 1)
}
