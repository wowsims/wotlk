package priest

import (
	"time"

	"github.com/wowsims/classic/sod/sim/core"
)

func (priest *Priest) registerPrayerOfHealingSpell() {
	var glyphSpell *core.Spell

	priest.PrayerOfHealing = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48072},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.48,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},

		BonusCritRating:  float64(priest.Talents.HolySpecialization),
		DamageMultiplier: 1 + .02*float64(priest.Talents.SpiritualHealing),
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
}
