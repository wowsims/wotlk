package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) ApplyTalents() {
	priest.setupSurgeOfLight()
	priest.registerInnerFocus()

	if priest.Talents.Shadowform {
		priest.PseudoStats.ShadowDamageDealtMultiplier *= 1.15
	}

	if priest.Talents.Meditation > 0 {
		priest.PseudoStats.SpiritRegenRateCasting = []float64{0.0, 0.17, 0.33, 0.5}[priest.Talents.Meditation]
	}

	if priest.Talents.SpiritualGuidance > 0 {
		bonus := (0.25 / 5) * float64(priest.Talents.SpiritualGuidance)
		priest.AddStatDependency2(stats.Spirit, stats.SpellPower, bonus)
	}

	if priest.Talents.MentalStrength > 0 {
		coeff := 0.02 * float64(priest.Talents.MentalStrength)
		priest.MultiplyStat(stats.Mana, 1.0+coeff)
	}

	// if priest.Talents.ForceOfWill > 0 {
	//coeff := 0.01 * float64(priest.Talents.ForceOfWill)
	//priest.AddStatDependency(stats.StatDependency{
	//	SourceStat:   stats.SpellPower,
	//ModifiedStat: stats.SpellPower,
	//	Modifier: func(spellPower float64, _ float64) float64 {
	//	return spellPower + spellPower*coeff
	//	},
	//})
	//priest.AddStat(stats.SpellCrit, float64(priest.Talents.ForceOfWill)*1*core.SpellCritRatingPerCritChance)
	//	}

	if priest.Talents.Enlightenment > 0 {
		coeff := 0.01 * float64(priest.Talents.Enlightenment)
		priest.MultiplyStat(stats.Intellect, 1.0+coeff)
		priest.MultiplyStat(stats.Stamina, 1.0+coeff)
		priest.MultiplyStat(stats.Spirit, 1.0+coeff)
	}

	if priest.Talents.SpiritOfRedemption {
		priest.MultiplyStat(stats.Spirit, 1.0+0.05)
	}

	if priest.Talents.TwistedFaith > 0 {
		priest.AddStatDependency2(stats.Spirit, stats.SpellPower, 0.04*float64(priest.Talents.TwistedFaith))
	}
}

func (priest *Priest) setupSurgeOfLight() {
	if priest.Talents.SurgeOfLight == 0 {
		return
	}

	priest.SurgeOfLightProcAura = priest.RegisterAura(core.Aura{
		Label:    "Surge of Light Proc",
		ActionID: core.ActionID{SpellID: 33151},
		Duration: core.NeverExpires,
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
			if spellEffect.Outcome.Matches(core.OutcomeCrit) {
				if procChance < sim.RandomFloat("SurgeOfLight") {
					priest.SurgeOfLightProcAura.Activate(sim)
					priest.SurgeOfLightProcAura.Prioritize()
				}
			}
		},
	})
}

func (priest *Priest) applySurgeOfLight(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
	if priest.SurgeOfLightProcAura.IsActive() {
		cast.CastTime = 0
		cast.Cost = 0
	}
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
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			priest.InnerFocusAura.Activate(sim)
		},
	})
}
