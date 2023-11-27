package priest

import (
	"time"

	"github.com/wowsims/classic/sod/sim/core"
	"github.com/wowsims/classic/sod/sim/core/stats"
)

func (priest *Priest) ApplyTalents() {
	priest.applyInspiration()
	priest.applyShadowWeaving()
	priest.registerInnerFocus()

	priest.PseudoStats.SpiritRegenRateCasting = []float64{0.0, 0.17, 0.33, 0.5}[priest.Talents.Meditation]
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= 1 - .02*float64(priest.Talents.SpellWarding)

	priest.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1 + (0.2 * float64(priest.Talents.Darkness))

	if priest.Talents.Shadowform {
		priest.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.15
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

	if priest.Talents.SpiritOfRedemption {
		priest.MultiplyStat(stats.Spirit, 1.05)
	}

	if priest.Talents.SilentResolve > 0 {
		priest.PseudoStats.ThreatMultiplier *= 1 - []float64{0, .04, .08, .12, .16, .20}[priest.Talents.SilentResolve]
	}
}

func (priest *Priest) applyInspiration() {
	if priest.Talents.Inspiration == 0 {
		return
	}

	auras := make([]*core.Aura, len(priest.Env.AllUnits))
	for _, unit := range priest.Env.AllUnits {
		if !priest.IsOpponent(unit) {
			aura := core.InspirationAura(unit, priest.Talents.Inspiration)
			auras[unit.UnitIndex] = aura
		}
	}

	priest.RegisterAura(core.Aura{
		Label:    "Inspiration Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == priest.FlashHeal ||
				spell == priest.GreaterHeal ||
				spell == priest.BindingHeal ||
				spell == priest.PrayerOfMending ||
				spell == priest.PrayerOfHealing ||
				spell == priest.CircleOfHealing ||
				spell == priest.PenanceHeal {
				auras[result.Target.UnitIndex].Activate(sim)
			}
		},
	})
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
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= 1.0 + 0.03*float64(oldStacks)
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.0 + 0.03*float64(newStacks)
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
			aura.Unit.AddStatDynamic(sim, stats.SpellCrit, 25*core.CritRatingPerCritChance)
			aura.Unit.PseudoStats.CostMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stats.SpellCrit, -25*core.CritRatingPerCritChance)
			aura.Unit.PseudoStats.CostMultiplier += 1
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Remove the buff and put skill on CD
			aura.Deactivate(sim)
			priest.InnerFocus.CD.Use(sim)
			priest.UpdateMajorCooldowns()
		},
	})

	priest.InnerFocus = priest.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

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
