package priest

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (priest *Priest) ApplyTalents() {
	priest.setupSurgeOfLight()
	priest.registerInnerFocus()

	if priest.Talents.Meditation > 0 {
		priest.PseudoStats.SpiritRegenRateCasting = float64(priest.Talents.Meditation) * 0.1
	}

	if priest.Talents.SpiritualGuidance > 0 {
		bonus := (0.25 / 5) * float64(priest.Talents.SpiritualGuidance)
		priest.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Spirit,
			ModifiedStat: stats.SpellPower,
			Modifier: func(spirit float64, spellPower float64) float64 {
				return spellPower + spirit*bonus
			},
		})
	}

	if priest.Talents.MentalStrength > 0 {
		coeff := 0.02 * float64(priest.Talents.MentalStrength)
		priest.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Mana,
			ModifiedStat: stats.Mana,
			Modifier: func(mana float64, _ float64) float64 {
				return mana + mana*coeff
			},
		})
	}

	if priest.Talents.ForceOfWill > 0 {
		coeff := 0.01 * float64(priest.Talents.ForceOfWill)
		priest.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.SpellPower,
			ModifiedStat: stats.SpellPower,
			Modifier: func(spellPower float64, _ float64) float64 {
				return spellPower + spellPower*coeff
			},
		})
		priest.AddStat(stats.SpellCrit, float64(priest.Talents.ForceOfWill)*1*core.SpellCritRatingPerCritChance)
	}

	if priest.Talents.Enlightenment > 0 {
		coeff := 0.01 * float64(priest.Talents.Enlightenment)
		priest.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.Intellect,
			Modifier: func(intellect float64, _ float64) float64 {
				return intellect + intellect*coeff
			},
		})

		priest.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Stamina,
			ModifiedStat: stats.Stamina,
			Modifier: func(stamina float64, _ float64) float64 {
				return stamina + stamina*coeff
			},
		})

		priest.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Spirit,
			ModifiedStat: stats.Spirit,
			Modifier: func(spirit float64, _ float64) float64 {
				return spirit + spirit*coeff
			},
		})
	}

	if priest.Talents.SpiritOfRedemption {
		priest.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Spirit,
			ModifiedStat: stats.Spirit,
			Modifier: func(spirit float64, _ float64) float64 {
				return spirit + spirit*0.05
			},
		})
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
	if priest.SurgeOfLightProcAura != nil && priest.SurgeOfLightProcAura.IsActive() {
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
			priest.AddStatDynamic(sim, stats.SpellCrit, 25*core.SpellCritRatingPerCritChance)
			priest.PseudoStats.NoCost = true
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			priest.AddStatDynamic(sim, stats.SpellCrit, -25*core.SpellCritRatingPerCritChance)
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
