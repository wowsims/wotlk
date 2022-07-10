package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	//"github.com/wowsims/wotlk/sim/core/proto"
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
		// TODO: Beast Mastery (in UI)

		if hunter.Talents.AnimalHandler != 0 {
			bonus := 1 + 0.05*float64(hunter.Talents.AnimalHandler)
			hunter.pet.AddStatDependency(stats.StatDependency{
				SourceStat:   stats.AttackPower,
				ModifiedStat: stats.AttackPower,
				Modifier: func(ap float64, _ float64) float64 {
					return ap * bonus
				},
			})
		}
	}

	hunter.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*1*float64(hunter.Talents.Surefooted))
	hunter.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(hunter.Talents.KillerInstinct))
	hunter.AddStat(stats.Parry, core.ParryRatingPerParryChance*1*float64(hunter.Talents.Deflection))
	hunter.pet.AddStat(stats.Dodge, 1*core.DodgeRatingPerDodgeChance*float64(hunter.Talents.CatlikeReflexes))
	hunter.PseudoStats.RangedSpeedMultiplier *= 1 + 0.04*float64(hunter.Talents.SerpentsSwiftness)
	hunter.PseudoStats.RangedDamageDealtMultiplier *= 1 + 0.01*float64(hunter.Talents.RangedWeaponSpecialization)
	hunter.PseudoStats.BonusRangedCritRating += 1 * float64(hunter.Talents.LethalShots) * core.CritRatingPerCritChance

	if hunter.Talents.EnduranceTraining > 0 {
		healthBonus := 0.01 * float64(hunter.Talents.EnduranceTraining)
		hunter.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Health,
			ModifiedStat: stats.Health,
			Modifier: func(health float64, _ float64) float64 {
				return health * (1 + healthBonus)
			},
		})
		if hunter.pet != nil {
			hunter.pet.AddStatDependency(stats.StatDependency{
				SourceStat:   stats.Health,
				ModifiedStat: stats.Health,
				Modifier: func(health float64, _ float64) float64 {
					return health * (1 + 2*healthBonus)
				},
			})
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
			hunter.pet.AddStatDependency(stats.StatDependency{
				SourceStat:   stats.Armor,
				ModifiedStat: stats.Armor,
				Modifier: func(armor float64, _ float64) float64 {
					return armor * (1 + petBonus)
				},
			})
		}
	}

	if hunter.Talents.Survivalist > 0 {
		healthBonus := 1 + 0.02*float64(hunter.Talents.Survivalist)
		hunter.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Health,
			ModifiedStat: stats.Health,
			Modifier: func(health float64, _ float64) float64 {
				return health * healthBonus
			},
		})
	}

	if hunter.Talents.CombatExperience > 0 {
		agiBonus := 1 + 0.01*float64(hunter.Talents.CombatExperience)
		hunter.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Agility,
			ModifiedStat: stats.Agility,
			Modifier: func(agility float64, _ float64) float64 {
				return agility * agiBonus
			},
		})
		intBonus := 1 + 0.03*float64(hunter.Talents.CombatExperience)
		hunter.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.Intellect,
			Modifier: func(intellect float64, _ float64) float64 {
				return intellect * intBonus
			},
		})
	}
	if hunter.Talents.CarefulAim > 0 {
		bonus := 0.15 * float64(hunter.Talents.CarefulAim)
		hunter.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.RangedAttackPower,
			Modifier: func(intellect float64, rap float64) float64 {
				return rap + intellect*bonus
			},
		})
	}
	if hunter.Talents.MasterMarksman > 0 {
		bonus := 1 + 0.02*float64(hunter.Talents.MasterMarksman)
		hunter.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.RangedAttackPower,
			ModifiedStat: stats.RangedAttackPower,
			Modifier: func(rap float64, _ float64) float64 {
				return rap * bonus
			},
		})
	}
	if hunter.Talents.SurvivalInstincts > 0 {
		apBonus := 1 + 0.02*float64(hunter.Talents.SurvivalInstincts)
		hunter.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.AttackPower,
			ModifiedStat: stats.AttackPower,
			Modifier: func(ap float64, _ float64) float64 {
				return ap * apBonus
			},
		})
		hunter.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.RangedAttackPower,
			ModifiedStat: stats.RangedAttackPower,
			Modifier: func(rap float64, _ float64) float64 {
				return rap * apBonus
			},
		})
	}
	if hunter.Talents.LightningReflexes > 0 {
		agiBonus := 1 + 0.03*float64(hunter.Talents.LightningReflexes)
		hunter.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Agility,
			ModifiedStat: stats.Agility,
			Modifier: func(agility float64, _ float64) float64 {
				return agility * agiBonus
			},
		})
	}

	hunter.applySpiritBond()
	hunter.applyInvigoration()
	hunter.applyCobraStrikes()

	hunter.applyGoForTheThroat()
	hunter.applySlaying()
	hunter.applyThrillOfTheHunt()
	hunter.applyExposeWeakness()
	hunter.applyMasterTactician()

	hunter.registerReadinessCD()
}

func (hunter *Hunter) critMultiplier(isRanged bool, target *core.Unit) float64 {
	primaryModifier := 1.0
	secondaryModifier := 0.0

	//monsterMultiplier := 1.0 + 0.01*float64(hunter.Talents.MonsterSlaying)
	//humanoidMultiplier := 1.0 + 0.01*float64(hunter.Talents.HumanoidSlaying)
	//if target.MobType == proto.MobType_MobTypeBeast || target.MobType == proto.MobType_MobTypeGiant || target.MobType == proto.MobType_MobTypeDragonkin {
	//	primaryModifier *= monsterMultiplier
	//} else if target.MobType == proto.MobType_MobTypeHumanoid {
	//	primaryModifier *= humanoidMultiplier
	//}

	//if isRanged {
	//	secondaryModifier += 0.06 * float64(hunter.Talents.MortalShots)
	//}

	return hunter.MeleeCritMultiplier(primaryModifier, secondaryModifier)
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

			// TODO: Kill shot too
			if spell != hunter.ArcaneShot && spell != hunter.SteadyShot {
				return
			}

			if sim.RandomFloat("Cobra Strikes") < procChance {
				hunter.pet.CobraStrikesAura.Activate(sim)
				hunter.pet.CobraStrikesAura.SetStacks(sim, 2)
			}
		},
	})
}

func (hunter *Hunter) applyFocusedFire() {
	if hunter.Talents.FocusedFire == 0 || hunter.pet == nil {
		return
	}

	hunter.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.01*float64(hunter.Talents.FocusedFire)
	// TODO: Pet special crit %
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
				Duration: hunter.applyLongevity(time.Minute * 2),
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

func (hunter *Hunter) applySlaying() {
	//if hunter.Talents.MonsterSlaying == 0 && hunter.Talents.HumanoidSlaying == 0 {
	//	return
	//}

	//monsterMultiplier := 1.0 + 0.01*float64(hunter.Talents.MonsterSlaying)
	//humanoidMultiplier := 1.0 + 0.01*float64(hunter.Talents.HumanoidSlaying)

	//switch hunter.CurrentTarget.MobType {
	//case proto.MobType_MobTypeBeast, proto.MobType_MobTypeGiant, proto.MobType_MobTypeDragonkin:
	//	hunter.PseudoStats.DamageDealtMultiplier *= monsterMultiplier
	//case proto.MobType_MobTypeHumanoid:
	//	hunter.PseudoStats.DamageDealtMultiplier *= humanoidMultiplier
	//}
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

	var debuffAura *core.Aura
	procChance := float64(hunter.Talents.ExposeWeakness) / 3

	hunter.RegisterAura(core.Aura{
		Label:    "Expose Weakness Talent",
		Duration: core.NeverExpires,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			debuffAura = core.ExposeWeaknessAura(hunter.CurrentTarget, float64(hunter.Index), 1.0)
		},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.ProcMask.Matches(core.ProcMaskRanged) {
				return
			}

			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			if procChance == 1 || sim.RandomFloat("ExposeWeakness") < procChance {
				// TODO: Find a cleaner way to do this
				newBonus := hunter.GetStat(stats.Agility) * 0.25
				if !debuffAura.IsActive() {
					debuffAura.Priority = newBonus
					debuffAura.Activate(sim)
				} else if debuffAura.Priority == newBonus {
					debuffAura.Activate(sim)
				} else if debuffAura.Priority < newBonus {
					debuffAura.Deactivate(sim)
					debuffAura.Priority = newBonus
					debuffAura.Activate(sim)
				}
			}
		},
	})
}

func (hunter *Hunter) applyMasterTactician() {
	if hunter.Talents.MasterTactician == 0 {
		return
	}

	procChance := 0.06
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

func (hunter *Hunter) registerReadinessCD() {
	if !hunter.Talents.Readiness {
		return
	}

	actionID := core.ActionID{SpellID: 23989}

	readinessSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			//GCD:         time.Second * 1, TODO: GCD causes panic
			//IgnoreHaste: true, // Hunter GCD is locked
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute * 5,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.RapidFire.CD.Reset()
			hunter.MultiShot.CD.Reset()
			hunter.ArcaneShot.CD.Reset()
			hunter.KillCommand.CD.Reset()
			hunter.RaptorStrike.CD.Reset()
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: readinessSpell,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Don't use if there are no cooldowns to reset.
			return !hunter.RapidFire.IsReady(sim)
		},
	})
}
