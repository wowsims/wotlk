package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (priest *Priest) registerPrayerOfHealingSpell() {
	var glyphSpell *core.Spell

	priest.PrayerOfHealing = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48072},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.48,
			Multiplier: 1 -
				.1*float64(priest.Talents.HealingPrayers) -
				core.TernaryFloat64(priest.HasSetBonus(ItemSetVestmentsOfAbsolution, 2), 0.1, 0),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},

		BonusCritRating: 0 +
			1*float64(priest.Talents.HolySpecialization)*core.CritRatingPerCritChance +
			core.TernaryFloat64(priest.HasSetBonus(ItemSetSanctificationRegalia, 2), 10*core.CritRatingPerCritChance, 0),
		DamageMultiplier: 1 *
			(1 + .02*float64(priest.Talents.SpiritualHealing)) *
			(1 + .01*float64(priest.Talents.BlessedResilience)) *
			(1 + .02*float64(priest.Talents.FocusedPower)) *
			(1 + .02*float64(priest.Talents.DivineProvidence)),
		CritMultiplier:   priest.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			targetAgent := target.Env.Raid.GetPlayerFromUnitIndex(target.UnitIndex)
			party := targetAgent.GetCharacter().Party

			for _, partyAgent := range party.PlayersAndPets {
				partyTarget := &partyAgent.GetCharacter().Unit
				baseHealing := sim.Roll(2109, 2228) + 0.526*spell.HealingPower(partyTarget)
				spell.CalcAndDealHealing(sim, partyTarget, baseHealing, spell.OutcomeHealingCrit)
				if glyphSpell != nil {
					glyphSpell.Hot(partyTarget).Apply(sim)
				}
			}
		},
	})

	if priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfPrayerOfHealing) {
		glyphSpell = priest.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 42409},
			SpellSchool: core.SpellSchoolHoly,
			ProcMask:    core.ProcMaskSpellHealing,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

			DamageMultiplier: priest.PrayerOfHealing.DamageMultiplier * 0.2 / 2,
			ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

			Hot: core.DotConfig{
				Aura: core.Aura{
					Label: "PoH Glyph",
				},
				NumberOfTicks: 2,
				TickLength:    time.Second * 3,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
					dot.SnapshotBaseDamage = sim.Roll(2109, 2228) + 0.526*dot.Spell.HealingPower(target)
					dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
				},
			},
		})
	}
}
