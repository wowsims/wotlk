package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (priest *Priest) registerGreaterHealSpell() {
	spellCoeff := 1.6114 + 0.08*float64(priest.Talents.EmpoweredHealing)

	priest.GreaterHeal = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48063},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.32,
			Multiplier: 1 *
				(1 - .05*float64(priest.Talents.ImprovedHealing)) *
				core.TernaryFloat64(priest.HasSetBonus(ItemSetRegaliaOfFaith, 4), .95, 1),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second*3 - time.Millisecond*100*time.Duration(priest.Talents.DivineFury),
			},
		},

		BonusCritRating: float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + .02*float64(priest.Talents.SpiritualHealing)) *
			(1 + .01*float64(priest.Talents.BlessedResilience)) *
			(1 + .02*float64(priest.Talents.FocusedPower)) *
			core.TernaryFloat64(priest.HasSetBonus(ItemSetVestmentsOfAbsolution, 4), 1.05, 1),
		CritMultiplier:   priest.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHealing := sim.Roll(3980, 4621) + spellCoeff*spell.HealingPower(target)
			spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)
		},
	})
}
