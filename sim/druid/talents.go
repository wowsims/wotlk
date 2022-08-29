package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) ApplyTalents() {
	druid.AddStat(stats.SpellHit, float64(druid.Talents.BalanceOfPower)*2*core.SpellHitRatingPerHitChance)
	druid.AddStat(stats.SpellCrit, float64(druid.Talents.NaturalPerfection)*1*core.CritRatingPerCritChance)
	druid.AddStat(stats.SpellPower, (float64(druid.Talents.ImprovedMoonkinForm)*0.1)*druid.GetStat(stats.Spirit))
	druid.PseudoStats.CastSpeedMultiplier *= 1 + (float64(druid.Talents.CelestialFocus) * 0.01)
	druid.PseudoStats.DamageDealtMultiplier *= 1 + (float64(druid.Talents.EarthAndMoon) * 0.02)
	druid.PseudoStats.SpiritRegenRateCasting = float64(druid.Talents.Intensity) * (0.5 / 3)
	druid.PseudoStats.ThreatMultiplier *= 1 - 0.04*float64(druid.Talents.Subtlety)
	druid.PseudoStats.PhysicalDamageDealtMultiplier *= 1 + 0.02*float64(druid.Talents.Naturalist)

	if druid.InForm(Bear | Cat) {
		if druid.Talents.PredatoryStrikes > 0 {
			druid.AddStat(stats.AttackPower, float64(druid.Talents.PredatoryStrikes)*0.5*float64(core.CharacterLevel))

			weap := druid.GetMHWeapon()
			if weap != nil {
				weap := druid.Equip[items.ItemSlotMainHand]
				dps := (((weap.WeaponDamageMax - weap.WeaponDamageMin) / 2.0) + weap.WeaponDamageMin) / weap.SwingSpeed
				fap := (dps - 54.8) * 14

				druid.AddStat(stats.AttackPower, fap*((0.2/3)*float64(druid.Talents.PredatoryStrikes)))
			}
		}
		druid.AddStat(stats.MeleeCrit, float64(druid.Talents.SharpenedClaws)*2*core.CritRatingPerCritChance)
		druid.AddStat(stats.Dodge, core.DodgeRatingPerDodgeChance*2*float64(druid.Talents.FeralSwiftness))
	}
	if druid.InForm(Bear) {
		druid.AddStat(stats.Armor, druid.Equip.Stats()[stats.Armor]*(0.5/3)*float64(druid.Talents.ThickHide))
	} else {
		druid.AddStat(stats.Armor, druid.Equip.Stats()[stats.Armor]*(0.1/3)*float64(druid.Talents.ThickHide))
	}
	if druid.InForm(Moonkin) && druid.Talents.MoonkinForm {
		druid.MultiplyStat(stats.Intellect, 1+(0.02*float64(druid.Talents.Furor)))
		druid.PseudoStats.DamageDealtMultiplier *= 1 + (float64(druid.Talents.MasterShapeshifter) * 0.02)
	}

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

		if druid.InForm(Cat) {
			druid.MultiplyStat(stats.AttackPower, 1.0+0.5*bonus)
		} else if druid.InForm(Bear) {
			druid.MultiplyStat(stats.Stamina, 1.0+0.5*bonus)
		}
	}

	if druid.Talents.SurvivalOfTheFittest > 0 {
		bonus := 0.02 * float64(druid.Talents.SurvivalOfTheFittest)
		druid.MultiplyStat(stats.Stamina, 1.0+bonus)
		druid.MultiplyStat(stats.Strength, 1.0+bonus)
		druid.MultiplyStat(stats.Agility, 1.0+bonus)
		druid.MultiplyStat(stats.Intellect, 1.0+bonus)
		druid.MultiplyStat(stats.Spirit, 1.0+bonus)
		druid.PseudoStats.ReducedCritTakenChance += 0.02 * float64(druid.Talents.SurvivalOfTheFittest)
		if druid.InForm(Bear) {
			druid.AddStat(stats.Armor, druid.Equip.Stats()[stats.Armor]*(0.33/3)*float64(druid.Talents.ThickHide))
		}
	}

	if druid.Talents.ImprovedMarkOfTheWild > 0 {
		bonus := 0.01 * float64(druid.Talents.ImprovedMarkOfTheWild)
		druid.MultiplyStat(stats.Stamina, 1.0+bonus)
		druid.MultiplyStat(stats.Strength, 1.0+bonus)
		druid.MultiplyStat(stats.Agility, 1.0+bonus)
		druid.MultiplyStat(stats.Intellect, 1.0+bonus)
		druid.MultiplyStat(stats.Spirit, 1.0+bonus)
	}

	if druid.Talents.ProtectorOfThePack > 0 {
		bonus := 0.02 * float64(druid.Talents.ProtectorOfThePack)
		if druid.InForm(Bear) {
			druid.MultiplyStat(stats.AttackPower, 1.0+bonus)
			druid.PseudoStats.DamageTakenMultiplier -= 0.04 * float64(druid.Talents.ProtectorOfThePack)
		}
	}

	if druid.Talents.PrimalPrecision > 0 {
		druid.AddStat(stats.Expertise, 5.0*float64(druid.Talents.PrimalPrecision))
	}

	if druid.Talents.LivingSpirit > 0 {
		bonus := 0.05 * float64(druid.Talents.LivingSpirit)
		druid.MultiplyStat(stats.Spirit, 1.0+bonus)
	}

	if druid.Talents.MasterShapeshifter > 0 {
		bonus := 0.02 * float64(druid.Talents.MasterShapeshifter)
		if druid.InForm(Bear) {
			druid.PseudoStats.DamageDealtMultiplier += bonus
		} else if druid.InForm(Cat) {
			druid.AddStat(stats.MeleeCrit, 2*float64(druid.Talents.MasterShapeshifter)*core.CritRatingPerCritChance)
		}
	}

	druid.setupNaturesGrace()
	druid.registerNaturesSwiftnessCD()
	druid.applyPrimalFury()
	druid.applyOmenOfClarity()
	druid.applyEclipse()
}

func (druid *Druid) setupNaturesGrace() {
	if druid.Talents.NaturesGrace < 1 {
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

	druid.RegisterAura(core.Aura{
		Label:    "Natures Grace",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			if (spell == druid.Starfire || spell == druid.Wrath) && float64(druid.Talents.NaturesGrace)*(1.0/3.0) >= sim.RandomFloat("Natures Grace") {
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

	spell := druid.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 3,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.NaturesSwiftnessAura.Activate(sim)
		},
	})

	druid.NaturesSwiftnessAura = druid.GetOrRegisterAura(core.Aura{
		Label:    "Natures Swiftness",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell != druid.Wrath && spell != druid.Starfire {
				return
			}

			// Remove the buff and put skill on CD
			aura.Deactivate(sim)
			spell.CD.Use(sim)
			druid.UpdateMajorCooldowns()
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Don't use NS unless we're casting a full-length starfire or wrath.
			return !character.HasTemporarySpellCastSpeedIncrease()
		},
	})
}

func (druid *Druid) applyNaturesSwiftness(cast *core.Cast) {
	if druid.NaturesSwiftnessAura.IsActive() {
		cast.CastTime = 0
	}
}

func (druid *Druid) applyPrimalFury() {
	if druid.Talents.PrimalFury == 0 {
		return
	}

	procChance := 0.5 * float64(druid.Talents.PrimalFury)
	actionID := core.ActionID{SpellID: 37117}
	rageMetrics := druid.NewRageMetrics(actionID)
	cpMetrics := druid.NewComboPointMetrics(actionID)

	druid.RegisterAura(core.Aura{
		Label:    "Primal Fury",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if druid.InForm(Bear) {
				if spellEffect.Outcome.Matches(core.OutcomeCrit) {
					if procChance == 1 || sim.RandomFloat("Primal Fury") < procChance {
						druid.AddRage(sim, 5, rageMetrics)
					}
				}
			} else if druid.InForm(Cat) {
				if druid.IsMangle(spell) || spell == druid.Shred || spell == druid.Rake {
					if spellEffect.Outcome.Matches(core.OutcomeCrit) {
						if procChance == 1 || sim.RandomFloat("Primal Fury") < procChance {
							druid.AddComboPoints(sim, 1, cpMetrics)
						}
					}
				}
			}
		},
	})
}

func (druid *Druid) applyOmenOfClarity() {
	if !druid.Talents.OmenOfClarity {
		return
	}

	ppmm := druid.AutoAttacks.NewPPMManager(3.5, core.ProcMaskMelee)

	druid.ClearcastingAura = druid.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: 16870},
		Duration: time.Second * 15,
	})
	// T10-2P
	lasherweave2P := druid.RegisterAura(core.Aura{
		Label:    "T10-2P proc",
		ActionID: core.ActionID{SpellID: 70718},
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.ArcaneDamageDealtMultiplier *= 1.15
			druid.PseudoStats.NatureDamageDealtMultiplier *= 1.15
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.ArcaneDamageDealtMultiplier /= 1.15
			druid.PseudoStats.NatureDamageDealtMultiplier /= 1.15
		},
	})

	druid.RegisterAura(core.Aura{
		Label:    "Omen of Clarity",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}
			if spellEffect.ProcMask.Matches(core.ProcMaskMeleeSpecial) && ppmm.Proc(sim, spellEffect.ProcMask, "Omen of Clarity") { // Melee Special
				druid.ClearcastingAura.Activate(sim)
			}
			if spellEffect.ProcMask.Matches(core.ProcMaskSpellDamage) && (spell == druid.Starfire || spell == druid.Wrath) { // Spells
				if sim.RandomFloat("Clearcasting") <= 1.75/(60/spell.CurCast.CastTime.Seconds()) { // 1.75 PPM emulation : https://github.com/JamminL/wotlk-classic-bugs/issues/66#issuecomment-1178282422
					druid.ClearcastingAura.Activate(sim)
					if druid.SetBonuses.balance_t10_2 {
						lasherweave2P.Activate(sim)
					}
				}
			}
		},
	})
}

func (druid *Druid) ApplyClearcasting(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
	if druid.ClearcastingAura.IsActive() {
		cast.Cost = 0
		druid.ClearcastingAura.Deactivate(sim)
	}
}

func (druid *Druid) applyEclipse() {
	druid.SolarICD = core.Cooldown{Timer: druid.NewTimer(), Duration: 0}
	druid.LunarICD = core.Cooldown{Timer: druid.NewTimer(), Duration: 0}
	if druid.Talents.Eclipse == 0 {
		return
	}

	// Solar
	solarProcChance := (1.0 / 3.0) * float64(druid.Talents.Eclipse)
	// TODO : make this proc a regular Aura
	solarProcAura := druid.NewTemporaryStatsAura("Solar Eclipse proc", core.ActionID{SpellID: 48517}, stats.Stats{}, time.Millisecond*15000)
	druid.SolarICD.Duration = time.Millisecond * 30000
	druid.RegisterAura(core.Aura{
		Label:    "Eclipse (Solar)",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
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
				solarProcAura.Activate(sim)
			}
		},
	})

	// Lunar
	lunarProcChance := 0.2 * float64(druid.Talents.Eclipse)
	// TODO : make this proc a regular Aura
	lunarProcAura := druid.NewTemporaryStatsAura("Lunar Eclipse proc", core.ActionID{SpellID: 48518}, stats.Stats{}, time.Millisecond*15000)
	druid.LunarICD.Duration = time.Millisecond * 30000
	druid.RegisterAura(core.Aura{
		Label:    "Eclipse (Lunar)",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
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
				lunarProcAura.Activate(sim)
			}
		},
	})
}

func (druid *Druid) ApplySwiftStarfireBonus(sim *core.Simulation, cast *core.Cast) {
	if druid.SwiftStarfireAura.IsActive() && druid.SetBonuses.balance_pvp_4 {
		cast.CastTime -= 1500 * time.Millisecond
		druid.SwiftStarfireAura.Deactivate(sim)
	}
}
