package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) ApplyTalents() {
	priest.applySurgeOfLight()
	priest.applyMisery()
	priest.applyShadowWeaving()
	priest.applyImprovedSpiritTap()
	priest.registerInnerFocus()

	priest.AddStat(stats.SpellCrit, 1*float64(priest.Talents.FocusedWill)*core.CritRatingPerCritChance)
	priest.PseudoStats.SpiritRegenRateCasting = []float64{0.0, 0.17, 0.33, 0.5}[priest.Talents.Meditation]
	priest.PseudoStats.ArcaneDamageTakenMultiplier *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.HolyDamageTakenMultiplier *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.FireDamageTakenMultiplier *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.FrostDamageTakenMultiplier *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.NatureDamageTakenMultiplier *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.ShadowDamageTakenMultiplier *= 1 - .02*float64(priest.Talents.SpellWarding)

	if priest.Talents.Shadowform {
		priest.PseudoStats.ShadowDamageDealtMultiplier *= 1.15
	}

	if priest.Talents.SpiritualGuidance > 0 {
		priest.AddStatDependency(stats.Spirit, stats.SpellPower, 0.05*float64(priest.Talents.SpiritualGuidance))
	}

	if priest.Talents.MentalStrength > 0 {
		priest.MultiplyStat(stats.Intellect, 1.0+0.03*float64(priest.Talents.MentalStrength))
	}

	if priest.Talents.ImprovedPowerWordFortitude > 0 {
		priest.MultiplyStat(stats.Stamina, 1.0+.02*float64(priest.Talents.ImprovedPowerWordFortitude))
	}

	if priest.Talents.Enlightenment > 0 {
		priest.MultiplyStat(stats.Spirit, 1+.02*float64(priest.Talents.Enlightenment))
		priest.PseudoStats.CastSpeedMultiplier *= 1 + .02*float64(priest.Talents.Enlightenment)
	}

	if priest.Talents.FocusedPower > 0 {
		priest.PseudoStats.DamageDealtMultiplier *= 1 + .02*float64(priest.Talents.FocusedPower)
		priest.PseudoStats.HealingDealtMultiplier *= 1 + .02*float64(priest.Talents.FocusedPower)
	}

	if priest.Talents.SpiritOfRedemption {
		priest.MultiplyStat(stats.Spirit, 1.05)
	}

	if priest.Talents.TwistedFaith > 0 {
		priest.AddStatDependency(stats.Spirit, stats.SpellPower, 0.04*float64(priest.Talents.TwistedFaith))
	}
}

func (priest *Priest) applySurgeOfLight() {
	if priest.Talents.SurgeOfLight == 0 {
		return
	}

	priest.SurgeOfLightProcAura = priest.RegisterAura(core.Aura{
		Label:    "Surge of Light Proc",
		ActionID: core.ActionID{SpellID: 33154},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			priest.Smite.CastTimeMultiplier -= 1
			priest.Smite.CostMultiplier -= 1
			priest.Smite.BonusCritRating -= 100 * core.CritRatingPerCritChance
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			priest.Smite.CastTimeMultiplier += 1
			priest.Smite.CostMultiplier += 1
			priest.Smite.BonusCritRating += 100 * core.CritRatingPerCritChance
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell == priest.Smite {
				aura.Deactivate(sim)
			}
		},
	})

	procChance := 0.25 * float64(priest.Talents.SurgeOfLight)

	priest.RegisterAura(core.Aura{
		Label:    "Surge of Light",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeCrit) && sim.RandomFloat("SurgeOfLight") < procChance {
				priest.SurgeOfLightProcAura.Activate(sim)
			}
		},
	})
}

func (priest *Priest) applyMisery() {
	if priest.Talents.Misery == 0 {
		return
	}

	priest.MiseryAura = core.MiseryAura(priest.CurrentTarget)
}

func (priest *Priest) applyShadowWeaving() {
	if priest.Talents.ShadowWeaving == 0 {
		return
	}

	priest.ShadowWeavingAura = priest.GetOrRegisterAura(core.Aura{
		Label:     "Shadow Weaving",
		ActionID:  core.ActionID{SpellID: 15258},
		Duration:  time.Second * 15,
		MaxStacks: 5,
		// TODO: This affects all spells not just direct damage. Dot damage should omit multipliers since it's snapshot at cast time.
		// OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
		// 	aura.Unit.PseudoStats.ShadowDamageDealtMultiplier /= 1.0 + 0.02*float64(oldStacks)
		// 	aura.Unit.PseudoStats.ShadowDamageDealtMultiplier *= 1.0 + 0.02*float64(newStacks)
		// },
	})
}

func (priest *Priest) applyImprovedSpiritTap() {
	if priest.Talents.ImprovedSpiritTap == 0 {
		return
	}

	increase := 1 + 0.05*float64(priest.Talents.ImprovedSpiritTap)
	statDep := priest.NewDynamicMultiplyStat(stats.Spirit, increase)

	priest.ImprovedSpiritTap = priest.GetOrRegisterAura(core.Aura{
		Label:    "Improved Spirit Tap",
		ActionID: core.ActionID{SpellID: 59000},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			priest.EnableDynamicStatDep(sim, statDep)
			priest.PseudoStats.SpiritRegenRateCasting += 0.33
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			priest.DisableDynamicStatDep(sim, statDep)
			priest.PseudoStats.SpiritRegenRateCasting -= 0.33
		},
	})
}

func (priest *Priest) registerInnerFocus() {
	if !priest.Talents.InnerFocus {
		return
	}

	actionID := core.ActionID{SpellID: 14751}

	priest.InnerFocusAura = priest.RegisterAura(core.Aura{
		Label:    "Inner Focus",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			priest.AddStatDynamic(sim, stats.SpellCrit, 25*core.CritRatingPerCritChance)
			priest.PseudoStats.NoCost = true
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			priest.AddStatDynamic(sim, stats.SpellCrit, -25*core.CritRatingPerCritChance)
			priest.PseudoStats.NoCost = false
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			// Remove the buff and put skill on CD
			aura.Deactivate(sim)
			priest.InnerFocus.CD.Use(sim)
		},
	})

	priest.InnerFocus = priest.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Duration(float64(time.Minute*3) * (1 - .1*float64(priest.Talents.Aspiration))),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			priest.InnerFocusAura.Activate(sim)
		},
	})
}
