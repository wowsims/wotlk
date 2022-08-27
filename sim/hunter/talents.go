package hunter

import (
	"strconv"
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
	hunter.PseudoStats.RangedDamageDealtMultiplier *= 1 + []float64{0, .01, .03, .05}[hunter.Talents.RangedWeaponSpecialization]
	hunter.PseudoStats.BonusRangedCritRating += 1 * float64(hunter.Talents.LethalShots) * core.CritRatingPerCritChance
	hunter.PseudoStats.DamageTakenMultiplier *= 1 - 0.02*float64(hunter.Talents.SurvivalInstincts)
	hunter.AutoAttacks.RangedEffect.DamageMultiplier *= hunter.markedForDeathMultiplier()

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
		// TODO: Activate replenishment
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

	hunter.registerReadinessCD()
}

func (hunter *Hunter) critMultiplier(isRanged bool, isMFDSpell bool, target *core.Unit) float64 {
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
				return
			}

			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			if procChance == 1 || sim.RandomFloat("Invigoration") < procChance {
				hunter.AddMana(sim, 0.01*hunter.MaxMana(), manaMetrics, false)
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
			if !hunter.pet.specialAbility.IsEmpty() {
				hunter.pet.specialAbility.BonusCritRating += 100 * core.CritRatingPerCritChance
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.pet.focusDump.BonusCritRating -= 100 * core.CritRatingPerCritChance
			if !hunter.pet.specialAbility.IsEmpty() {
				hunter.pet.specialAbility.BonusCritRating -= 100 * core.CritRatingPerCritChance
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
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

	actionID := core.ActionID{SpellID: 53238}
	dmgMultiplier := 0.1 * float64(hunter.Talents.PiercingShots)
	var psDot *core.Dot

	psSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagNoOnCastComplete,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			psDot.Apply(sim)
		},
	})

	target := hunter.CurrentTarget
	psDot = core.NewDot(core.Dot{
		Spell: psSpell,
		Aura: target.GetOrRegisterAura(core.Aura{
			Label:    "PiercingShots-" + strconv.Itoa(int(hunter.Index)),
			ActionID: actionID,
			Duration: time.Second * 8,
		}),
		NumberOfTicks: 8,
		TickLength:    time.Second * 1,
	})

	var currentTickDmg float64

	hunter.RegisterAura(core.Aura{
		Label:    "Piercing Shots Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			if spell != hunter.AimedShot && spell != hunter.SteadyShot && spell != hunter.ChimeraShot {
				return
			}

			totalDmg := spellEffect.Damage * dmgMultiplier
			if psDot.IsActive() {
				remainingTicks := 8 - psDot.TickCount
				totalDmg += currentTickDmg * float64(remainingTicks)
			}
			currentTickDmg = totalDmg / 8

			// Reassign tick effect to update the damage.
			psDot.TickEffects = core.TickFuncSnapshot(target, core.SpellEffect{
				ProcMask:         core.ProcMaskPeriodicDamage,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				IsPeriodic:       true,
				BaseDamage:       core.BaseDamageConfigFlat(currentTickDmg),
				OutcomeApplier:   hunter.OutcomeFuncTick(),
			})

			psSpell.Cast(sim, spellEffect.Target)
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
		Flags:       core.SpellFlagNoOnCastComplete,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskRangedAuto,
			DamageMultiplier: 0.8,
			ThreatMultiplier: 1,

			BaseDamage:     core.BaseDamageConfigRangedWeapon(0),
			OutcomeApplier: hunter.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(false, false, hunter.CurrentTarget)),
		}),
	})

	hunter.RegisterAura(core.Aura{
		Label:    "Wild Quiver Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell != hunter.AutoAttacks.RangedAuto {
				return
			}

			if sim.RandomFloat("Wild Quiver") < procChance {
				wqSpell.Cast(sim, spellEffect.Target)
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			if procChance == 1 || sim.RandomFloat("Frenzy") < procChance {
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
			aura.Unit.PseudoStats.CostMultiplier *= 0.5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
			aura.Unit.PseudoStats.CostMultiplier /= 0.5
		},
	})

	manaCost := hunter.BaseMana * 0.1

	bwSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.Mana,
		BaseCost:     manaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: manaCost,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: hunter.applyLongevity(time.Minute*2) - core.TernaryDuration(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfBestialWrath), time.Second*20, 0),
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
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return hunter.CurrentMana() >= manaCost
		},
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.ProcMask.Matches(core.ProcMaskRanged) || !spellEffect.Outcome.Matches(core.OutcomeCrit) {
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
			if !applied {
				for _, target := range hunter.Env.Encounter.Targets {
					switch target.MobType {
					case proto.MobType_MobTypeBeast, proto.MobType_MobTypeDemon,
						proto.MobType_MobTypeDragonkin, proto.MobType_MobTypeElemental,
						proto.MobType_MobTypeGiant, proto.MobType_MobTypeHumanoid,
						proto.MobType_MobTypeUndead:

						hunter.AttackTables[target.UnitIndex].DamageDealtMultiplier *= 1.0 + 0.01*float64(hunter.Talents.ImprovedTracking)
					}
				}
				applied = true
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
			if hunter.ExplosiveShot != nil {
				hunter.ExplosiveShot.CostMultiplier -= 1
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.ArcaneShot.CostMultiplier += 1
			if hunter.ExplosiveShot != nil {
				hunter.ExplosiveShot.CostMultiplier += 1
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell == hunter.ArcaneShot || spell == hunter.ExplosiveShot {
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
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			// mask 256
			if !spellEffect.ProcMask.Matches(core.ProcMaskRangedSpecial) {
				return
			}

			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			if procChance == 1 || sim.RandomFloat("ThrillOfTheHunt") < procChance {
				hunter.AddMana(sim, spell.CurCast.Cost*0.4, manaMetrics, false)
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.ProcMask.Matches(core.ProcMaskRanged) && spell != hunter.ExplosiveTrap {
				return
			}

			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			if procChance == 1 || sim.RandomFloat("ExposeWeakness") < procChance {
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.ProcMask.Matches(core.ProcMaskRanged) || !spellEffect.Landed() {
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
			hunter.SteadyShot.DamageMultiplier += dmgMod
			if hunter.AimedShot != nil {
				hunter.AimedShot.DamageMultiplier += dmgMod
			}
			if hunter.BlackArrow != nil {
				hunter.BlackArrow.DamageMultiplier += dmgMod
			}
			if hunter.ExplosiveShot != nil {
				hunter.ExplosiveShot.DamageMultiplier += dmgMod
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.SteadyShot.DamageMultiplier -= dmgMod
			if hunter.AimedShot != nil {
				hunter.AimedShot.DamageMultiplier -= dmgMod
			}
			if hunter.BlackArrow != nil {
				hunter.BlackArrow.DamageMultiplier -= dmgMod
			}
			if hunter.ExplosiveShot != nil {
				hunter.ExplosiveShot.DamageMultiplier -= dmgMod
			}
		},
	})

	core.ApplyFixedUptimeAura(stAura, uptime, time.Second*15)
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

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.RapidFire.CD.Reset()
			hunter.MultiShot.CD.Reset()
			hunter.ArcaneShot.CD.Reset()
			hunter.KillCommand.CD.Reset()
			hunter.RaptorStrike.CD.Reset()
			hunter.ExplosiveTrap.CD.Reset()
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
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: readinessSpell,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Don't use if there are no cooldowns to reset.
			return !hunter.RapidFire.IsReady(sim)
		},
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return !hunter.RapidFireAura.IsActive() || hunter.RapidFireAura.RemainingDuration(sim) < time.Second*10
		},
	})
}
