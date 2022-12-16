package druid

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) ThickHideMultiplier() float64 {
	thickHideMulti := 1.0

	if druid.Talents.ThickHide > 0 {
		thickHideMulti += 0.04 + 0.03*float64(druid.Talents.ThickHide-1)
	}

	return thickHideMulti
}

func (druid *Druid) TotalBearArmorMultiplier() float64 {
	sotfMulti := 1.0 + 0.33/3.0*float64(druid.Talents.SurvivalOfTheFittest)
	return 4.7 * sotfMulti * druid.ThickHideMultiplier()
}

func (druid *Druid) ApplyTalents() {
	druid.AddStat(stats.SpellHit, float64(druid.Talents.BalanceOfPower)*2*core.SpellHitRatingPerHitChance)
	druid.AddStat(stats.SpellCrit, float64(druid.Talents.NaturalPerfection)*1*core.CritRatingPerCritChance)
	druid.PseudoStats.CastSpeedMultiplier *= 1 + (float64(druid.Talents.CelestialFocus) * 0.01)
	druid.PseudoStats.DamageDealtMultiplier *= 1 + (float64(druid.Talents.EarthAndMoon) * 0.02)
	druid.PseudoStats.SpiritRegenRateCasting = float64(druid.Talents.Intensity) * (0.5 / 3)
	druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1 + 0.02*float64(druid.Talents.Naturalist)
	druid.AddStat(stats.Armor, druid.ScaleBaseArmor(druid.ThickHideMultiplier()-1.0))

	if druid.Talents.LunarGuidance > 0 {
		bonus := 0.04 * float64(druid.Talents.LunarGuidance)
		druid.AddStatDependency(stats.Intellect, stats.SpellPower, bonus)
	}

	if druid.Talents.Dreamstate > 0 {
		bonus := 0.04 * float64(druid.Talents.Dreamstate)
		druid.AddStatDependency(stats.Intellect, stats.MP5, bonus)
	}

	if druid.Talents.HeartOfTheWild > 0 {
		bonus := 0.04 * float64(druid.Talents.HeartOfTheWild)
		druid.MultiplyStat(stats.Intellect, 1.0+bonus)
	}

	if druid.Talents.ImprovedFaerieFire > 0 && druid.CurrentTarget.HasAuraWithTag(core.FaerieFireAuraTag) {
		druid.AddStat(stats.SpellCrit, float64(druid.Talents.ImprovedFaerieFire)*1*core.CritRatingPerCritChance)
	}

	if druid.Talents.SurvivalOfTheFittest > 0 {
		bonus := 0.02 * float64(druid.Talents.SurvivalOfTheFittest)
		druid.MultiplyStat(stats.Stamina, 1.0+bonus)
		druid.MultiplyStat(stats.Strength, 1.0+bonus)
		druid.MultiplyStat(stats.Agility, 1.0+bonus)
		druid.MultiplyStat(stats.Intellect, 1.0+bonus)
		druid.MultiplyStat(stats.Spirit, 1.0+bonus)
		druid.PseudoStats.ReducedCritTakenChance += 0.02 * float64(druid.Talents.SurvivalOfTheFittest)
	}

	if druid.Talents.ImprovedMarkOfTheWild > 0 {
		bonus := 0.01 * float64(druid.Talents.ImprovedMarkOfTheWild)
		druid.MultiplyStat(stats.Stamina, 1.0+bonus)
		druid.MultiplyStat(stats.Strength, 1.0+bonus)
		druid.MultiplyStat(stats.Agility, 1.0+bonus)
		druid.MultiplyStat(stats.Intellect, 1.0+bonus)
		druid.MultiplyStat(stats.Spirit, 1.0+bonus)
	}

	if druid.Talents.PrimalPrecision > 0 {
		druid.AddStat(stats.Expertise, 5.0*float64(druid.Talents.PrimalPrecision)*core.ExpertisePerQuarterPercentReduction)
	}

	if druid.Talents.LivingSpirit > 0 {
		bonus := 0.05 * float64(druid.Talents.LivingSpirit)
		druid.MultiplyStat(stats.Spirit, 1.0+bonus)
	}

	druid.setupNaturesGrace()
	druid.registerNaturesSwiftnessCD()
	druid.applyMoonkinForm()
	druid.applyPrimalFury()
	druid.applyOmenOfClarity()
	druid.applyEclipse()
	druid.applyImprovedLotp()
	druid.applyPredatoryInstincts()
}

func (druid *Druid) setupNaturesGrace() {
	if druid.Talents.NaturesGrace == 0 {
		return
	}

	druid.NaturesGraceProcAura = druid.RegisterAura(core.Aura{
		Label:    "Natures Grace Proc",
		ActionID: core.ActionID{SpellID: 16886},
		Duration: time.Second * 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.MultiplyCastSpeed(1.2)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.MultiplyCastSpeed(1 / 1.2)
		},
	})

	procChance := []float64{0, .33, .66, 1}[druid.Talents.NaturesGrace]

	druid.RegisterAura(core.Aura{
		Label:    "Natures Grace",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			if spell.Flags.Matches(SpellFlagNaturesGrace) && sim.Proc(procChance, "Natures Grace") {
				druid.NaturesGraceProcAura.Activate(sim)
			}
		},
	})
}

func (druid *Druid) registerNaturesSwiftnessCD() {
	if !druid.Talents.NaturesSwiftness {
		return
	}
	actionID := core.ActionID{SpellID: 17116}

	var nsAura *core.Aura
	nsSpell := druid.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 3,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			nsAura.Activate(sim)
		},
	})

	nsAura = druid.RegisterAura(core.Aura{
		Label:    "Natures Swiftness",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.Starfire.CastTimeMultiplier -= 1
			druid.Wrath.CastTimeMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.Starfire.CastTimeMultiplier += 1
			druid.Wrath.CastTimeMultiplier += 1
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell != druid.Wrath && spell != druid.Starfire {
				return
			}

			// Remove the buff and put skill on CD
			aura.Deactivate(sim)
			nsSpell.CD.Use(sim)
			druid.UpdateMajorCooldowns()
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: nsSpell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Don't use NS unless we're casting a full-length starfire or wrath.
			return !character.HasTemporarySpellCastSpeedIncrease()
		},
	})
}

func (druid *Druid) applyEarthAndMoon() {
	if druid.Talents.EarthAndMoon == 0 {
		return
	}
	druid.EarthAndMoonAura = core.EarthAndMoonAura(druid.CurrentTarget, druid.Talents.EarthAndMoon)
}

func (druid *Druid) applyPrimalFury() {
	if druid.Talents.PrimalFury == 0 {
		return
	}

	procChance := []float64{0, 0.5, 1}[druid.Talents.PrimalFury]
	actionID := core.ActionID{SpellID: 37117}
	rageMetrics := druid.NewRageMetrics(actionID)
	cpMetrics := druid.NewComboPointMetrics(actionID)

	druid.RegisterAura(core.Aura{
		Label:    "Primal Fury",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if druid.InForm(Bear) {
				if result.Outcome.Matches(core.OutcomeCrit) {
					if sim.Proc(procChance, "Primal Fury") {
						druid.AddRage(sim, 5, rageMetrics)
					}
				}
			} else if druid.InForm(Cat) {
				if druid.IsMangle(spell) || spell == druid.Shred || spell == druid.Rake {
					if result.Outcome.Matches(core.OutcomeCrit) {
						if sim.Proc(procChance, "Primal Fury") {
							druid.AddComboPoints(sim, 1, cpMetrics)
						}
					}
				}
			}
		},
	})
}

// Modifies the Bleed aura to apply the bonus.
func (druid *Druid) applyRendAndTear(aura core.Aura) core.Aura {
	if druid.Talents.RendAndTear == 0 || druid.AssumeBleedActive {
		return aura
	}

	bonusCrit := 5.0 * float64(druid.Talents.RendAndTear) * core.CritRatingPerCritChance

	aura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
		if druid.BleedsActive == 0 {
			druid.FerociousBite.BonusCritRating += bonusCrit
		}
		druid.BleedsActive++
	})
	aura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
		druid.BleedsActive--
		if druid.BleedsActive == 0 {
			druid.FerociousBite.BonusCritRating -= bonusCrit
		}
	})

	return aura
}

func (druid *Druid) applyOmenOfClarity() {
	if !druid.Talents.OmenOfClarity {
		return
	}

	// T10-2P
	var lasherweave2P *core.Aura
	if druid.HasSetBonus(ItemSetLasherweaveRegalia, 2) {
		lasherweave2P = druid.RegisterAura(core.Aura{
			Label:    "T10-2P proc",
			ActionID: core.ActionID{SpellID: 70718},
			Duration: time.Second * 6,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexArcane] *= 1.15
				druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] *= 1.15
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexArcane] /= 1.15
				druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] /= 1.15
			},
		})
	}

	var affectedSpells []*core.Spell
	druid.ClearcastingAura = druid.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: 16870},
		Duration: time.Second * 15,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice([]*core.Spell{
				// Balance
				druid.Hurricane,
				druid.Starfire,
				druid.Typhoon,
				druid.Wrath,

				// Feral
				druid.DemoralizingRoar,
				druid.FerociousBite,
				druid.Lacerate,
				druid.MangleBear,
				druid.MangleCat,
				druid.Maul,
				druid.Rake,
				druid.Rip,
				druid.Shred,
				druid.SwipeBear,
				druid.SwipeCat,
			}, func(spell *core.Spell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.CostMultiplier -= 1
			}

			if lasherweave2P != nil {
				lasherweave2P.Activate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.CostMultiplier += 1
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			for _, as := range affectedSpells {
				if as == spell {
					aura.Deactivate(sim)
					break
				}
			}
		},
	})

	druid.RegisterAura(core.Aura{
		Label:    "Omen of Clarity",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// https://github.com/JamminL/wotlk-classic-bugs/issues/66#issuecomment-1182017571
			if spell == druid.Hurricane {
				curCastTickSpeed := spell.CurCast.ChannelTime.Seconds() / 10
				hurricaneCoeff := 1.0 - (7.0 / 9.0)
				spellCoeff := hurricaneCoeff * curCastTickSpeed
				chanceToProc := ((1.5 / 60) * 3.5) * spellCoeff
				if sim.RandomFloat("Clearcasting") <= chanceToProc {
					druid.ClearcastingAura.Activate(sim)
				}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}
			// Not ideal to have new ppm manager here, but this needs to account for feral druid bear<->cat swaps
			ppmm := druid.AutoAttacks.NewPPMManager(3.5, core.ProcMaskMeleeWhiteHit)
			if ppmm.Proc(sim, spell.ProcMask, "Omen of Clarity") { // Melee
				druid.ClearcastingAura.Activate(sim)
			} else if spell.ProcMask.Matches(core.ProcMaskSpellDamage) { // Spells
				// Heavily based on comment here
				// https://github.com/JamminL/wotlk-classic-bugs/issues/66#issuecomment-1182017571
				// Instants are treated as 1.5
				castTime := spell.DefaultCast.CastTime.Seconds()
				if castTime == 0 {
					castTime = 1.5
				}

				chanceToProc := (castTime / 60) * 3.5
				if spell == druid.Typhoon { // Add Typhoon
					chanceToProc *= 0.25
				} else if spell == druid.Moonfire { // Add Moonfire
					chanceToProc *= 0.076
				} else {
					chanceToProc *= 0.666
				}
				if sim.RandomFloat("Clearcasting") <= chanceToProc {
					druid.ClearcastingAura.Activate(sim)
				}
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == druid.GiftOfTheWild {
				// Based on ingame testing by druid discord, subject to change or incorrectness
				chanceToProc := 1.0 - math.Pow(1.0-0.0875, float64(druid.RaidBuffTargets))
				if sim.RandomFloat("Clearcasting") <= chanceToProc {
					druid.ClearcastingAura.Activate(sim)
				}
			}
		},
	})
}

func (druid *Druid) applyEclipse() {
	druid.SolarICD = core.Cooldown{Timer: druid.NewTimer(), Duration: 0}
	druid.LunarICD = core.Cooldown{Timer: druid.NewTimer(), Duration: 0}
	if druid.Talents.Eclipse == 0 {
		return
	}

	// Solar
	solarProcChance := (1.0 / 3.0) * float64(druid.Talents.Eclipse)
	solarProcMultiplier := 1.4 + core.TernaryFloat64(druid.HasSetBonus(ItemSetNightsongGarb, 2), 0.07, 0)
	druid.SolarICD.Duration = time.Millisecond * 30000
	druid.SolarEclipseProcAura = druid.RegisterAura(core.Aura{
		Label:    "Solar Eclipse proc",
		Duration: time.Millisecond * 15000,
		ActionID: core.ActionID{SpellID: 48517},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.Wrath.DamageMultiplier *= solarProcMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.Wrath.DamageMultiplier /= solarProcMultiplier
		},
	})

	druid.RegisterAura(core.Aura{
		Label:    "Eclipse (Solar)",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			if spell != druid.Starfire {
				return
			}
			if !druid.SolarICD.Timer.IsReady(sim) {
				return
			}
			if druid.LunarICD.Timer.TimeToReady(sim) > time.Millisecond*15000 {
				return
			}
			if sim.RandomFloat("Eclipse (Solar)") < solarProcChance {
				druid.SolarICD.Use(sim)
				druid.SolarEclipseProcAura.Activate(sim)
			}
		},
	})

	// Lunar
	lunarProcChance := 0.2 * float64(druid.Talents.Eclipse)
	lunarBonusCrit := (40 + core.TernaryFloat64(druid.HasSetBonus(ItemSetNightsongGarb, 2), 7, 0)) * core.CritRatingPerCritChance
	druid.LunarICD.Duration = time.Millisecond * 30000
	druid.LunarEclipseProcAura = druid.RegisterAura(core.Aura{
		Label:    "Lunar Eclipse proc",
		Duration: time.Millisecond * 15000,
		ActionID: core.ActionID{SpellID: 48518},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.Starfire.BonusCritRating += lunarBonusCrit
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.Starfire.BonusCritRating -= lunarBonusCrit
		},
	})
	druid.RegisterAura(core.Aura{
		Label:    "Eclipse (Lunar)",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			if spell != druid.Wrath {
				return
			}
			if !druid.LunarICD.Timer.IsReady(sim) {
				return
			}
			if druid.SolarICD.Timer.TimeToReady(sim) > time.Millisecond*15000 {
				return
			}
			if sim.RandomFloat("Eclipse (Lunar)") < lunarProcChance {
				druid.LunarICD.Use(sim)
				druid.LunarEclipseProcAura.Activate(sim)
			}
		},
	})
}

func (druid *Druid) applyImprovedLotp() {
	if druid.Talents.ImprovedLeaderOfThePack == 0 {
		return
	}

	manaMetrics := druid.NewManaMetrics(core.ActionID{SpellID: 34300})
	manaRestore := float64(druid.Talents.ImprovedLeaderOfThePack) * 0.04

	icd := core.Cooldown{
		Timer:    druid.NewTimer(),
		Duration: time.Second * 6,
	}

	druid.RegisterAura(core.Aura{
		Label:    "Improved Leader of the Pack",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}
			if !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) || !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			if !icd.IsReady(sim) {
				return
			}
			icd.Use(sim)
			druid.AddMana(sim, druid.MaxMana()*manaRestore, manaMetrics, false)
		},
	})
}

func (druid *Druid) applyPredatoryInstincts() {
	if druid.Talents.PredatoryInstincts == 0 {
		return
	}

	onGainMod := druid.MeleeCritMultiplier(Cat)
	onExpireMod := druid.MeleeCritMultiplier(Humanoid)

	druid.PredatoryInstinctsAura = druid.RegisterAura(core.Aura{
		Label:    "Predatory Instincts",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.LacerateDot.Spell.CritMultiplier = onGainMod
			druid.RipDot.Spell.CritMultiplier = onGainMod
			druid.RakeDot.Spell.CritMultiplier = onGainMod
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.LacerateDot.Spell.CritMultiplier = onExpireMod
			druid.RipDot.Spell.CritMultiplier = onExpireMod
			druid.RakeDot.Spell.CritMultiplier = onExpireMod
		},
	})
}
