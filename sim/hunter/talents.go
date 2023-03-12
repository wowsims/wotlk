package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) ApplyTalents() {
	if hunter.pet != nil {
		hunter.applyFocusedFire()
		hunter.applyFrenzy()
		hunter.registerBestialWrathCD()

		hunter.pet.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*2*float64(hunter.Talents.Ferocity))
		hunter.pet.AddStat(stats.SpellCrit, core.CritRatingPerCritChance*2*float64(hunter.Talents.Ferocity))
		hunter.pet.AddStat(stats.Dodge, 3*core.DodgeRatingPerDodgeChance*float64(hunter.Talents.CatlikeReflexes))
		hunter.pet.PseudoStats.DamageDealtMultiplier *= 1 + 0.03*float64(hunter.Talents.UnleashedFury)
		hunter.pet.PseudoStats.DamageDealtMultiplier *= 1 + 0.04*float64(hunter.Talents.KindredSpirits)
		hunter.pet.PseudoStats.MeleeSpeedMultiplier *= 1 + 0.04*float64(hunter.Talents.SerpentsSwiftness)

		if hunter.Talents.AnimalHandler != 0 {
			hunter.pet.MultiplyStat(stats.AttackPower, 1+(0.05*float64(hunter.Talents.AnimalHandler)))
		}
		hunter.pet.ApplyTalents()
	}

	hunter.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*1*float64(hunter.Talents.FocusedAim))
	hunter.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(hunter.Talents.KillerInstinct))
	hunter.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(hunter.Talents.MasterMarksman))
	hunter.AddStat(stats.Parry, core.ParryRatingPerParryChance*1*float64(hunter.Talents.Deflection))
	hunter.AddStat(stats.Dodge, 1*core.DodgeRatingPerDodgeChance*float64(hunter.Talents.CatlikeReflexes))
	hunter.PseudoStats.RangedSpeedMultiplier *= 1 + 0.04*float64(hunter.Talents.SerpentsSwiftness)
	hunter.PseudoStats.DamageTakenMultiplier *= 1 - 0.02*float64(hunter.Talents.SurvivalInstincts)
	hunter.AutoAttacks.RangedConfig.DamageMultiplier *= hunter.markedForDeathMultiplier()

	if hunter.Talents.LethalShots > 0 {
		hunter.AddBonusRangedCritRating(1 * float64(hunter.Talents.LethalShots) * core.CritRatingPerCritChance)
	}
	if hunter.Talents.RangedWeaponSpecialization > 0 {
		mult := 1 + []float64{0, .01, .03, .05}[hunter.Talents.RangedWeaponSpecialization]
		hunter.OnSpellRegistered(func(spell *core.Spell) {
			if spell.ProcMask.Matches(core.ProcMaskRanged) {
				spell.DamageMultiplier *= mult
			}
		})
	}

	if hunter.Talents.EnduranceTraining > 0 {
		healthBonus := 0.01 * float64(hunter.Talents.EnduranceTraining)
		hunter.MultiplyStat(stats.Health, 1.0+healthBonus)
		if hunter.pet != nil {
			hunter.pet.MultiplyStat(stats.Health, 1.0+2*healthBonus)
		}
	}

	if hunter.Talents.ThickHide > 0 {
		var hunterBonus, petBonus float64
		if hunter.Talents.ThickHide == 1 {
			hunterBonus = 0.04
			petBonus = 0.07
		} else if hunter.Talents.ThickHide == 2 {
			hunterBonus = 0.07
			petBonus = 0.14
		} else if hunter.Talents.ThickHide == 3 {
			hunterBonus = 0.1
			petBonus = 0.2
		}
		hunter.AddStat(stats.Armor, hunter.Equip.Stats()[stats.Armor]*hunterBonus)
		if hunter.pet != nil {
			hunter.pet.MultiplyStat(stats.Armor, 1.0+petBonus)
		}
	}

	if hunter.Talents.Survivalist > 0 {
		hunter.MultiplyStat(stats.Stamina, 1.0+0.02*float64(hunter.Talents.Survivalist))
	}

	if hunter.Talents.CombatExperience > 0 {
		bonus := 1.0 + (0.02 * float64(hunter.Talents.CombatExperience))
		hunter.MultiplyStat(stats.Agility, bonus)
		hunter.MultiplyStat(stats.Intellect, bonus)
	}
	if hunter.Talents.CarefulAim > 0 {
		hunter.AddStatDependency(stats.Intellect, stats.RangedAttackPower, (1.0/3.0)*float64(hunter.Talents.CarefulAim))
	}
	if hunter.Talents.HunterVsWild > 0 {
		bonus := 0.1 * float64(hunter.Talents.HunterVsWild)
		hunter.AddStatDependency(stats.Stamina, stats.AttackPower, bonus)
		hunter.AddStatDependency(stats.Stamina, stats.RangedAttackPower, bonus)
	}
	if hunter.Talents.LightningReflexes > 0 {
		agiBonus := 0.03 * float64(hunter.Talents.LightningReflexes)
		hunter.MultiplyStat(stats.Agility, 1.0+agiBonus)
	}
	if hunter.Talents.HuntingParty > 0 {
		agiBonus := 0.01 * float64(hunter.Talents.HuntingParty)
		hunter.MultiplyStat(stats.Agility, 1.0+agiBonus)
	}

	hunter.applySpiritBond()
	hunter.applyInvigoration()
	hunter.applyCobraStrikes()
	hunter.applyGoForTheThroat()
	hunter.applyPiercingShots()
	hunter.applyWildQuiver()
	hunter.applyImprovedTracking()
	hunter.applyThrillOfTheHunt()
	hunter.applyLockAndLoad()
	hunter.applyExposeWeakness()
	hunter.applyMasterTactician()
	hunter.applySniperTraining()
	hunter.applyHuntingParty()

	hunter.registerReadinessCD()
}

func (hunter *Hunter) critMultiplier(isRanged bool, isMFDSpell bool) float64 {
	primaryModifier := 1.0
	secondaryModifier := 0.0

	if isRanged {
		secondaryModifier += 0.06 * float64(hunter.Talents.MortalShots)
		if isMFDSpell {
			secondaryModifier += 0.02 * float64(hunter.Talents.MarkedForDeath)
		}
	}

	return hunter.MeleeCritMultiplier(primaryModifier, secondaryModifier)
}

func (hunter *Hunter) markedForDeathMultiplier() float64 {
	if hunter.Options.UseHuntersMark || hunter.Env.GetTarget(0).HasAuraWithTag(core.HuntersMarkAuraTag) {
		return 1 + .01*float64(hunter.Talents.MarkedForDeath)
	} else {
		return 1
	}
}

func (hunter *Hunter) applySpiritBond() {
	if hunter.Talents.SpiritBond == 0 || hunter.pet == nil {
		return
	}

	hunter.PseudoStats.HealingTakenMultiplier *= 1 + 0.05*float64(hunter.Talents.SpiritBond)
	hunter.pet.PseudoStats.HealingTakenMultiplier *= 1 + 0.05*float64(hunter.Talents.SpiritBond)

	actionID := core.ActionID{SpellID: 20895}
	healthMultiplier := 0.01 * float64(hunter.Talents.SpiritBond)
	healthMetrics := hunter.NewHealthMetrics(actionID)
	petHealthMetrics := hunter.pet.NewHealthMetrics(actionID)

	hunter.RegisterResetEffect(func(sim *core.Simulation) {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period: time.Second * 10,
			OnAction: func(sim *core.Simulation) {
				hunter.GainHealth(sim, hunter.MaxHealth()*healthMultiplier, healthMetrics)
				hunter.pet.GainHealth(sim, hunter.pet.MaxHealth()*healthMultiplier, petHealthMetrics)
			},
		})
	})
}

func (hunter *Hunter) applyInvigoration() {
	if hunter.Talents.Invigoration == 0 || hunter.pet == nil {
		return
	}

	procChance := 0.5 * float64(hunter.Talents.Invigoration)
	manaMetrics := hunter.NewManaMetrics(core.ActionID{SpellID: 53253})

	hunter.pet.RegisterAura(core.Aura{
		Label:    "Invigoration",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
				return
			}

			if !result.DidCrit() {
				return
			}

			if sim.Proc(procChance, "Invigoration") {
				hunter.AddMana(sim, 0.01*hunter.MaxMana(), manaMetrics)
			}
		},
	})
}

func (hunter *Hunter) applyCobraStrikes() {
	if hunter.Talents.CobraStrikes == 0 || hunter.pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 53260}
	procChance := 0.2 * float64(hunter.Talents.CobraStrikes)

	hunter.pet.CobraStrikesAura = hunter.pet.RegisterAura(core.Aura{
		Label:     "Cobra Strikes",
		ActionID:  actionID,
		Duration:  time.Second * 10,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.pet.focusDump.BonusCritRating += 100 * core.CritRatingPerCritChance
			if hunter.pet.specialAbility != nil {
				hunter.pet.specialAbility.BonusCritRating += 100 * core.CritRatingPerCritChance
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.pet.focusDump.BonusCritRating -= 100 * core.CritRatingPerCritChance
			if hunter.pet.specialAbility != nil {
				hunter.pet.specialAbility.BonusCritRating -= 100 * core.CritRatingPerCritChance
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
				aura.RemoveStack(sim)
			}
		},
	})

	hunter.RegisterAura(core.Aura{
		Label:    "Cobra Strikes",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() {
				return
			}

			if spell != hunter.ArcaneShot && spell != hunter.SteadyShot && spell != hunter.KillShot {
				return
			}

			if sim.RandomFloat("Cobra Strikes") < procChance {
				hunter.pet.CobraStrikesAura.Activate(sim)
				hunter.pet.CobraStrikesAura.SetStacks(sim, 2)
			}
		},
	})
}

func (hunter *Hunter) applyPiercingShots() {
	if hunter.Talents.PiercingShots == 0 {
		return
	}

	psSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53238},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "PiercingShots",
				Duration: time.Second * 8,
			},
			NumberOfTicks: 8,
			TickLength:    time.Second * 1,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Specifically account for bleed modifiers, since it still affects the spell, but we're ignoring all modifiers.
				dot.SnapshotAttackerMultiplier = target.PseudoStats.PeriodicPhysicalDamageTakenMultiplier
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrReset(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
	})

	hunter.RegisterAura(core.Aura{
		Label:    "Piercing Shots Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() {
				return
			}
			if spell != hunter.AimedShot && spell != hunter.SteadyShot && spell != hunter.ChimeraShot {
				return
			}

			dot := psSpell.Dot(result.Target)
			outstandingDamage := core.TernaryFloat64(dot.IsActive(), dot.SnapshotBaseDamage*float64(dot.NumberOfTicks-dot.TickCount), 0)
			newDamage := result.Damage * 0.1 * float64(hunter.Talents.PiercingShots)

			dot.SnapshotBaseDamage = (outstandingDamage + newDamage) / float64(dot.NumberOfTicks)
			psSpell.Cast(sim, result.Target)
		},
	})
}

func (hunter *Hunter) applyWildQuiver() {
	if hunter.Talents.WildQuiver == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 53217}
	procChance := 0.04 * float64(hunter.Talents.WildQuiver)

	wqSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskRangedAuto,
		Flags:       core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 0.8,
		CritMultiplier:   hunter.critMultiplier(false, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.RangedWeaponDamage(sim, spell.RangedAttackPower(target)) +
				spell.BonusWeaponDamage()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
		},
	})

	hunter.RegisterAura(core.Aura{
		Label:    "Wild Quiver Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != hunter.AutoAttacks.RangedAuto {
				return
			}

			if sim.RandomFloat("Wild Quiver") < procChance {
				wqSpell.Cast(sim, result.Target)
			}
		},
	})
}

func (hunter *Hunter) applyFocusedFire() {
	if hunter.Talents.FocusedFire == 0 || hunter.pet == nil {
		return
	}

	hunter.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.01*float64(hunter.Talents.FocusedFire)
}

func (hunter *Hunter) applyFrenzy() {
	if hunter.Talents.Frenzy == 0 {
		return
	}

	procChance := 0.2 * float64(hunter.Talents.Frenzy)

	procAura := hunter.pet.RegisterAura(core.Aura{
		Label:    "Frenzy Proc",
		ActionID: core.ActionID{SpellID: 19625},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.MeleeSpeedMultiplier *= 1.3
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.MeleeSpeedMultiplier /= 1.3
		},
	})

	hunter.pet.RegisterAura(core.Aura{
		Label:    "Frenzy",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() {
				return
			}
			if sim.Proc(procChance, "Frenzy") {
				procAura.Activate(sim)
			}
		},
	})
}

func (hunter *Hunter) applyLongevity(dur time.Duration) time.Duration {
	return time.Duration(float64(dur) * (1.0 - 0.1*float64(hunter.Talents.Longevity)))
}

func (hunter *Hunter) registerBestialWrathCD() {
	if !hunter.Talents.BestialWrath {
		return
	}
	if hunter.Talents.TheBeastWithin {
		hunter.PseudoStats.DamageDealtMultiplier *= 1.1
	}

	actionID := core.ActionID{SpellID: 19574}

	bestialWrathPetAura := hunter.pet.RegisterAura(core.Aura{
		Label:    "Bestial Wrath Pet",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.5
		},
	})

	bestialWrathAura := hunter.RegisterAura(core.Aura{
		Label:    "Bestial Wrath",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.1
			aura.Unit.PseudoStats.CostMultiplier -= 0.5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
			aura.Unit.PseudoStats.CostMultiplier += 0.5
		},
	})

	bwSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.1,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: hunter.applyLongevity(time.Minute*2 - core.TernaryDuration(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfBestialWrath), time.Second*20, 0)),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			bestialWrathPetAura.Activate(sim)

			if hunter.Talents.TheBeastWithin {
				bestialWrathAura.Activate(sim)
			}
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: bwSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (hunter *Hunter) applyGoForTheThroat() {
	if hunter.Talents.GoForTheThroat == 0 {
		return
	}
	if hunter.pet == nil {
		return
	}

	amount := 25.0 * float64(hunter.Talents.GoForTheThroat)

	hunter.RegisterAura(core.Aura{
		Label:    "Go for the Throat",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskRanged) || !result.DidCrit() {
				return
			}
			if !hunter.pet.IsEnabled() {
				return
			}
			hunter.pet.AddFocus(sim, amount, core.ActionID{SpellID: 34954})
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() {
				return
			}
			if !hunter.pet.IsEnabled() {
				return
			}
			hunter.pet.AddFocus(sim, amount, core.ActionID{SpellID: 34954})
		},
	})
}

func (hunter *Hunter) applyImprovedTracking() {
	if hunter.Talents.ImprovedTracking == 0 {
		return
	}

	var applied bool

	hunter.RegisterResetEffect(
		func(s *core.Simulation) {
			if applied {
				return
			}
			applied = true

			for _, target := range hunter.Env.Encounter.TargetUnits {
				switch target.MobType {
				case proto.MobType_MobTypeBeast, proto.MobType_MobTypeDemon,
					proto.MobType_MobTypeDragonkin, proto.MobType_MobTypeElemental,
					proto.MobType_MobTypeGiant, proto.MobType_MobTypeHumanoid,
					proto.MobType_MobTypeUndead:

					hunter.AttackTables[target.UnitIndex].DamageDealtMultiplier *= 1.0 + 0.01*float64(hunter.Talents.ImprovedTracking)
				}
			}
		},
	)
}

func (hunter *Hunter) applyLockAndLoad() {
	if hunter.Talents.LockAndLoad == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 56344}
	procChance := []float64{0, 0.02, 0.04, 0.20}[hunter.Talents.LockAndLoad]

	icd := core.Cooldown{
		Timer:    hunter.NewTimer(),
		Duration: time.Second * 22,
	}

	hunter.LockAndLoadAura = hunter.RegisterAura(core.Aura{
		Label:     "Lock and Load Proc",
		ActionID:  actionID,
		Duration:  time.Second * 12,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.ArcaneShot.CostMultiplier -= 1
			if hunter.ExplosiveShotR4 != nil {
				hunter.ExplosiveShotR4.CostMultiplier -= 1
				hunter.ExplosiveShotR3.CostMultiplier -= 1
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.ArcaneShot.CostMultiplier += 1
			if hunter.ExplosiveShotR4 != nil {
				hunter.ExplosiveShotR4.CostMultiplier += 1
				hunter.ExplosiveShotR3.CostMultiplier += 1
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == hunter.ArcaneShot || spell == hunter.ExplosiveShotR4 || spell == hunter.ExplosiveShotR3 {
				aura.RemoveStack(sim)
				hunter.ArcaneShot.CD.Reset() // Shares the CD with explosive shot.
			}
		},
	})

	hunter.RegisterAura(core.Aura{
		Label:    "Lock and Load Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != hunter.BlackArrow && spell != hunter.ExplosiveTrap {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			if sim.RandomFloat("Lock and Load") < procChance {
				icd.Use(sim)
				hunter.LockAndLoadAura.Activate(sim)
				hunter.LockAndLoadAura.SetStacks(sim, 2)
			}
		},
	})
}

func (hunter *Hunter) applyThrillOfTheHunt() {
	if hunter.Talents.ThrillOfTheHunt == 0 {
		return
	}

	procChance := float64(hunter.Talents.ThrillOfTheHunt) / 3
	manaMetrics := hunter.NewManaMetrics(core.ActionID{SpellID: 34499})

	hunter.RegisterAura(core.Aura{
		Label:    "Thrill of the Hunt",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// mask 256
			if !spell.ProcMask.Matches(core.ProcMaskRangedSpecial) {
				return
			}

			if !result.DidCrit() {
				return
			}

			if sim.Proc(procChance, "ThrillOfTheHunt") {
				hunter.AddMana(sim, spell.CurCast.Cost*0.4, manaMetrics)
			}
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() && (spell == hunter.ExplosiveShotR4 || spell == hunter.ExplosiveShotR3) {
				// Explosive shot ticks can proc TotH but with 1/3 the bonus.
				if sim.Proc(procChance, "ThrillOfTheHunt") {
					hunter.AddMana(sim, spell.CurCast.Cost*0.4/3, manaMetrics)
				}
			}
		},
	})
}

func (hunter *Hunter) applyExposeWeakness() {
	if hunter.Talents.ExposeWeakness == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 34503}
	procChance := float64(hunter.Talents.ExposeWeakness) / 3

	apDep := hunter.NewDynamicStatDependency(stats.Agility, stats.AttackPower, .25)
	rapDep := hunter.NewDynamicStatDependency(stats.Agility, stats.RangedAttackPower, .25)
	procAura := hunter.RegisterAura(core.Aura{
		Label:    "Expose Weakness Proc",
		ActionID: actionID,
		Duration: time.Second * 7,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.EnableDynamicStatDep(sim, apDep)
			aura.Unit.EnableDynamicStatDep(sim, rapDep)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.DisableDynamicStatDep(sim, apDep)
			aura.Unit.DisableDynamicStatDep(sim, rapDep)
		},
	})

	hunter.RegisterAura(core.Aura{
		Label:    "Expose Weakness Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskRanged) && spell != hunter.ExplosiveTrap {
				return
			}

			if !result.DidCrit() {
				return
			}

			if sim.Proc(procChance, "ExposeWeakness") {
				procAura.Activate(sim)
			}
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() {
				return
			}

			if sim.Proc(procChance, "ExposeWeakness") {
				procAura.Activate(sim)
			}
		},
	})
}

func (hunter *Hunter) applyMasterTactician() {
	if hunter.Talents.MasterTactician == 0 {
		return
	}

	procChance := 0.1
	critBonus := 2 * core.CritRatingPerCritChance * float64(hunter.Talents.MasterTactician)

	procAura := hunter.NewTemporaryStatsAura("Master Tactician Proc", core.ActionID{SpellID: 34839}, stats.Stats{stats.MeleeCrit: critBonus}, time.Second*8)

	hunter.RegisterAura(core.Aura{
		Label:    "Master Tactician",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskRanged) || !result.Landed() {
				return
			}

			if sim.RandomFloat("Master Tactician") > procChance {
				return
			}

			procAura.Activate(sim)
		},
	})
}

func (hunter *Hunter) applySniperTraining() {
	if hunter.Talents.SniperTraining == 0 {
		return
	}

	uptime := hunter.Options.SniperTrainingUptime
	if uptime <= 0 {
		return
	}
	uptime = core.MinFloat(1, uptime)

	dmgMod := .02 * float64(hunter.Talents.SniperTraining)

	stAura := hunter.RegisterAura(core.Aura{
		Label:    "Sniper Training",
		ActionID: core.ActionID{SpellID: 53304},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.SteadyShot.DamageMultiplierAdditive += dmgMod
			if hunter.AimedShot != nil {
				hunter.AimedShot.DamageMultiplierAdditive += dmgMod
			}
			if hunter.BlackArrow != nil {
				hunter.BlackArrow.DamageMultiplierAdditive += dmgMod
			}
			if hunter.ExplosiveShotR4 != nil {
				hunter.ExplosiveShotR4.DamageMultiplierAdditive += dmgMod
				hunter.ExplosiveShotR3.DamageMultiplierAdditive += dmgMod
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.SteadyShot.DamageMultiplierAdditive -= dmgMod
			if hunter.AimedShot != nil {
				hunter.AimedShot.DamageMultiplierAdditive -= dmgMod
			}
			if hunter.BlackArrow != nil {
				hunter.BlackArrow.DamageMultiplierAdditive -= dmgMod
			}
			if hunter.ExplosiveShotR4 != nil {
				hunter.ExplosiveShotR4.DamageMultiplierAdditive -= dmgMod
				hunter.ExplosiveShotR3.DamageMultiplierAdditive -= dmgMod
			}
		},
	})

	core.ApplyFixedUptimeAura(stAura, uptime, time.Second*15, 1)
}

func (hunter *Hunter) applyHuntingParty() {
	if hunter.Talents.HuntingParty == 0 {
		return
	}

	procChance := float64(hunter.Talents.HuntingParty) / 3
	replSrc := hunter.Env.Raid.NewReplenishmentSource(core.ActionID{SpellID: 53292})

	hunter.RegisterAura(core.Aura{
		Label:    "Hunting Party",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() {
				return
			}

			if spell != hunter.SteadyShot && spell != hunter.ArcaneShot && spell != hunter.ExplosiveShotR4 && spell != hunter.ExplosiveShotR3 {
				return
			}

			if procChance == 1 || sim.RandomFloat("Hunting Party") < procChance {
				hunter.Env.Raid.ProcReplenishment(sim, replSrc)
			}
		},
	})
}

func (hunter *Hunter) registerReadinessCD() {
	if !hunter.Talents.Readiness {
		return
	}

	actionID := core.ActionID{SpellID: 23989}

	readinessSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 1,
			},
			IgnoreHaste: true, // Hunter GCD is locked
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute * 3,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Don't use if there are no cooldowns to reset.
			return !hunter.RapidFire.IsReady(sim)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.RapidFire.CD.Reset()
			hunter.MultiShot.CD.Reset()
			hunter.ArcaneShot.CD.Reset()
			hunter.KillShot.CD.Reset()
			hunter.RaptorStrike.CD.Reset()
			hunter.ExplosiveTrap.CD.Reset()
			if hunter.KillCommand != nil {
				hunter.KillCommand.CD.Reset()
			}
			if hunter.AimedShot != nil {
				hunter.AimedShot.CD.Reset()
			}
			if hunter.SilencingShot != nil {
				hunter.SilencingShot.CD.Reset()
			}
			if hunter.ChimeraShot != nil {
				hunter.ChimeraShot.CD.Reset()
			}
			if hunter.BlackArrow != nil {
				hunter.BlackArrow.CD.Reset()
			}

			// TODO: This is needed because there are edge cases where core doesn't re-use Rapid Fire.
			// Fix core so this isn't necessary.
			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + 1,
				OnAction: func(_ *core.Simulation) {
					hunter.UpdateMajorCooldowns()
				},
			})
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: readinessSpell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// If RF is about to become ready naturally, wait so we can get 2x usages.
			if !hunter.RapidFire.IsReady(sim) && hunter.RapidFire.TimeToReady(sim) < time.Second*10 {
				return false
			}
			return !hunter.RapidFireAura.IsActive() || hunter.RapidFireAura.RemainingDuration(sim) < time.Second*10
		},
	})
}
