package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) ApplyTalents() {
	// TODO:
	// Reflective Shield
	// Improved Flash Heal
	// Renewed Hope
	// Rapture
	// Pain Suppression
	// Test of Faith
	// Guardian Spirit

	priest.applyDivineAegis()
	priest.applyGrace()
	priest.applyBorrowedTime()
	priest.applyInspiration()
	priest.applyHolyConcentration()
	priest.applySerendipity()
	priest.applySurgeOfLight()
	priest.applyMisery()
	priest.applyShadowWeaving()
	priest.applyImprovedSpiritTap()
	priest.registerInnerFocus()

	priest.AddStat(stats.SpellCrit, 1*float64(priest.Talents.FocusedWill)*core.CritRatingPerCritChance)
	priest.PseudoStats.SpiritRegenRateCasting = []float64{0.0, 0.17, 0.33, 0.5}[priest.Talents.Meditation]
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= 1 - .02*float64(priest.Talents.SpellWarding)
	priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= 1 - .02*float64(priest.Talents.SpellWarding)

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

	if priest.Talents.Enlightenment > 0 {
		priest.MultiplyStat(stats.Spirit, 1+.02*float64(priest.Talents.Enlightenment))
		priest.PseudoStats.CastSpeedMultiplier *= 1 + .02*float64(priest.Talents.Enlightenment)
	}

	if priest.Talents.FocusedPower > 0 {
		priest.PseudoStats.DamageDealtMultiplier *= 1 + .02*float64(priest.Talents.FocusedPower)
	}

	if priest.Talents.SpiritOfRedemption {
		priest.MultiplyStat(stats.Spirit, 1.05)
	}

	if priest.Talents.TwistedFaith > 0 {
		priest.AddStatDependency(stats.Spirit, stats.SpellPower, 0.04*float64(priest.Talents.TwistedFaith))
	}
}

func (priest *Priest) applyDivineAegis() {
	if priest.Talents.DivineAegis == 0 {
		return
	}

	divineAegis := priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47515},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

		DamageMultiplier: 1 *
			(0.1 * float64(priest.Talents.DivineAegis)) *
			core.TernaryFloat64(priest.HasSetBonus(ItemSetZabrasRaiment, 4), 1.1, 1),
		ThreatMultiplier: 1,

		Shield: core.ShieldConfig{
			Aura: core.Aura{
				Label:    "Divine Aegis",
				Duration: time.Second * 12,
			},
		},
	})

	priest.RegisterAura(core.Aura{
		Label:    "Divine Aegis Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeCrit) {
				divineAegis.Shield(result.Target).Apply(sim, result.Damage)
			}
		},
	})
}

func (priest *Priest) applyGrace() {
	if priest.Talents.Grace == 0 {
		return
	}

	procChance := .5 * float64(priest.Talents.Grace)

	auras := make([]*core.Aura, len(priest.Env.AllUnits))
	for _, unit := range priest.Env.AllUnits {
		if !priest.IsOpponent(unit) {
			aura := unit.RegisterAura(core.Aura{
				Label:     "Grace" + strconv.Itoa(int(priest.Index)),
				ActionID:  core.ActionID{SpellID: 47517},
				Duration:  time.Second * 15,
				MaxStacks: 3,
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					priest.AttackTables[aura.Unit.UnitIndex].HealingDealtMultiplier /= 1 + .03*float64(oldStacks)
					priest.AttackTables[aura.Unit.UnitIndex].HealingDealtMultiplier *= 1 + .03*float64(newStacks)
				},
			})
			auras[unit.UnitIndex] = aura
		}
	}

	priest.RegisterAura(core.Aura{
		Label:    "Grace Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == priest.FlashHeal || spell == priest.GreaterHeal || spell == priest.PenanceHeal {
				if sim.Proc(procChance, "Grace") {
					aura := auras[result.Target.UnitIndex]
					aura.Activate(sim)
					aura.AddStack(sim)
				}
			}
		},
	})
}

// This one is called from healing priest sim initialization because it needs an input.
func (priest *Priest) ApplyRapture(ppm float64) {
	if priest.Talents.Rapture == 0 {
		return
	}

	if ppm <= 0 {
		return
	}

	raptureManaCoeff := []float64{0, .015, .020, .025}[priest.Talents.Rapture]
	raptureMetrics := priest.NewManaMetrics(core.ActionID{SpellID: 47537})

	priest.RegisterResetEffect(func(sim *core.Simulation) {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period: time.Minute / time.Duration(ppm),
			OnAction: func(sim *core.Simulation) {
				priest.AddMana(sim, raptureManaCoeff*priest.MaxMana(), raptureMetrics)
			},
		})
	})
}

func (priest *Priest) applyBorrowedTime() {
	if priest.Talents.BorrowedTime == 0 {
		return
	}

	multiplier := 1 + .05*float64(priest.Talents.BorrowedTime)

	procAura := priest.RegisterAura(core.Aura{
		Label:    "Borrowed Time",
		ActionID: core.ActionID{SpellID: 52800},
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			priest.MultiplyCastSpeed(multiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			priest.MultiplyCastSpeed(1 / multiplier)
		},
	})

	priest.RegisterAura(core.Aura{
		Label:    "Borrwed Time Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == priest.PowerWordShield {
				procAura.Activate(sim)
			} else if spell.CurCast.CastTime > 0 {
				procAura.Deactivate(sim)
			}
		},
	})
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

func (priest *Priest) applyHolyConcentration() {
	if priest.Talents.HolyConcentration == 0 {
		return
	}

	multiplier := 1 + []float64{0, .16, .32, .50}[priest.Talents.HolyConcentration]

	procAura := priest.RegisterAura(core.Aura{
		Label:    "Holy Concentration",
		ActionID: core.ActionID{SpellID: 34860},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			priest.PseudoStats.SpiritRegenMultiplier *= multiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			priest.PseudoStats.SpiritRegenMultiplier /= multiplier
		},
	})

	priest.RegisterAura(core.Aura{
		Label:    "Holy Concentration Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() &&
				(spell == priest.FlashHeal || spell == priest.GreaterHeal || spell == priest.EmpoweredRenew) {
				procAura.Activate(sim)
			}
		},
	})
}

func (priest *Priest) applySerendipity() {
	if priest.Talents.Serendipity == 0 {
		return
	}

	reductionPerStack := .04 * float64(priest.Talents.Serendipity)

	procAura := priest.RegisterAura(core.Aura{
		Label:     "Serendipity",
		ActionID:  core.ActionID{SpellID: 63737},
		Duration:  time.Second * 20,
		MaxStacks: 3,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			priest.PrayerOfHealing.CastTimeMultiplier += reductionPerStack * float64(oldStacks)
			priest.PrayerOfHealing.CastTimeMultiplier -= reductionPerStack * float64(newStacks)
			priest.GreaterHeal.CastTimeMultiplier += reductionPerStack * float64(oldStacks)
			priest.GreaterHeal.CastTimeMultiplier -= reductionPerStack * float64(newStacks)
		},
	})

	priest.RegisterAura(core.Aura{
		Label:    "Serendipity Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == priest.FlashHeal || spell == priest.BindingHeal {
				procAura.Activate(sim)
				procAura.AddStack(sim)
			} else if spell == priest.GreaterHeal || spell == priest.PrayerOfHealing {
				procAura.Deactivate(sim)
			}
		},
	})
}

func (priest *Priest) applySurgeOfLight() {
	if priest.Talents.SurgeOfLight == 0 {
		return
	}

	procHandler := func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if spell == priest.Smite || spell == priest.FlashHeal {
			aura.Deactivate(sim)
		}
	}

	priest.SurgeOfLightProcAura = priest.RegisterAura(core.Aura{
		Label:    "Surge of Light Proc",
		ActionID: core.ActionID{SpellID: 33154},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if priest.Smite != nil {
				priest.Smite.CastTimeMultiplier -= 1
				priest.Smite.CostMultiplier -= 1
				priest.Smite.BonusCritRating -= 100 * core.CritRatingPerCritChance
			}
			if priest.FlashHeal != nil {
				priest.FlashHeal.CastTimeMultiplier -= 1
				priest.FlashHeal.CostMultiplier -= 1
				priest.FlashHeal.BonusCritRating -= 100 * core.CritRatingPerCritChance
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if priest.Smite != nil {
				priest.Smite.CastTimeMultiplier += 1
				priest.Smite.CostMultiplier += 1
				priest.Smite.BonusCritRating += 100 * core.CritRatingPerCritChance
			}
			if priest.FlashHeal != nil {
				priest.FlashHeal.CastTimeMultiplier += 1
				priest.FlashHeal.CostMultiplier += 1
				priest.FlashHeal.BonusCritRating += 100 * core.CritRatingPerCritChance
			}
		},
		OnSpellHitDealt: procHandler,
		OnHealDealt:     procHandler,
	})

	procChance := 0.25 * float64(priest.Talents.SurgeOfLight)
	icd := core.Cooldown{
		Timer:    priest.NewTimer(),
		Duration: time.Second * 6,
	}
	priest.SurgeOfLightProcAura.Icd = &icd

	handler := func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if icd.IsReady(sim) && result.Outcome.Matches(core.OutcomeCrit) && sim.RandomFloat("SurgeOfLight") < procChance {
			icd.Use(sim)
			priest.SurgeOfLightProcAura.Activate(sim)
		}
	}

	priest.RegisterAura(core.Aura{
		Label:    "Surge of Light",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: handler,
		OnHealDealt:     handler,
	})
}

func (priest *Priest) applyMisery() {
	if priest.Talents.Misery == 0 {
		return
	}

	miseryAuras := priest.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.MiseryAura(target, priest.Talents.Misery)
	})
	priest.Env.RegisterPreFinalizeEffect(func() {
		priest.ShadowWordPain.RelatedAuras = append(priest.ShadowWordPain.RelatedAuras, miseryAuras)
		if priest.VampiricTouch != nil {
			priest.VampiricTouch.RelatedAuras = append(priest.VampiricTouch.RelatedAuras, miseryAuras)
		}
		if priest.MindFlay[1] != nil {
			priest.MindFlayAPL.RelatedAuras = append(priest.MindFlayAPL.RelatedAuras, miseryAuras)
			priest.MindFlay[1].RelatedAuras = append(priest.MindFlay[1].RelatedAuras, miseryAuras)
			priest.MindFlay[2].RelatedAuras = append(priest.MindFlay[2].RelatedAuras, miseryAuras)
			priest.MindFlay[3].RelatedAuras = append(priest.MindFlay[3].RelatedAuras, miseryAuras)
		}
	})

	priest.RegisterAura(core.Aura{
		Label:    "Priest Shadow Effects",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if spell == priest.ShadowWordPain || spell == priest.VampiricTouch || spell.ActionID.SpellID == priest.MindFlay[1].ActionID.SpellID {
				miseryAuras.Get(result.Target).Activate(sim)
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
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= 1.0 + 0.02*float64(oldStacks)
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.0 + 0.02*float64(newStacks)
		},
	})
}

func (priest *Priest) applyImprovedSpiritTap() {
	if priest.Talents.ImprovedSpiritTap == 0 {
		return
	}

	increase := 1 + 0.05*float64(priest.Talents.ImprovedSpiritTap)
	statDep := priest.NewDynamicMultiplyStat(stats.Spirit, increase)
	regen := []float64{0, 0.17, 0.33}[priest.Talents.ImprovedSpiritTap]

	priest.ImprovedSpiritTap = priest.GetOrRegisterAura(core.Aura{
		Label:    "Improved Spirit Tap",
		ActionID: core.ActionID{SpellID: 59000},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			priest.EnableDynamicStatDep(sim, statDep)
			priest.PseudoStats.SpiritRegenRateCasting += regen
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			priest.DisableDynamicStatDep(sim, statDep)
			priest.PseudoStats.SpiritRegenRateCasting -= regen
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
				Duration: time.Duration(float64(time.Minute*3) * (1 - .1*float64(priest.Talents.Aspiration))),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			priest.InnerFocusAura.Activate(sim)
		},
	})
}
