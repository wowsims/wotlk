package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (hunter *Hunter) ApplyTalents() {
	if hunter.pet != nil {
		hunter.applyFocusedFire()
		hunter.applyFrenzy()
		hunter.applyFerociousInspiration()
		hunter.registerBestialWrathCD()

		hunter.pet.AddStat(stats.MeleeCrit, core.MeleeCritRatingPerCritChance*2*float64(hunter.Talents.Ferocity))
		hunter.pet.AddStat(stats.SpellCrit, core.SpellCritRatingPerCritChance*2*float64(hunter.Talents.Ferocity))
		hunter.pet.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*2*float64(hunter.Talents.AnimalHandler))
		hunter.pet.AddStat(stats.SpellHit, core.SpellHitRatingPerHitChance*2*float64(hunter.Talents.AnimalHandler))
		hunter.pet.PseudoStats.DamageDealtMultiplier *= 1 + 0.04*float64(hunter.Talents.UnleashedFury)
		hunter.pet.PseudoStats.MeleeSpeedMultiplier *= 1 + 0.04*float64(hunter.Talents.SerpentsSwiftness)
	}

	hunter.applyGoForTheThroat()
	hunter.applySlaying()
	hunter.applyThrillOfTheHunt()
	hunter.applyExposeWeakness()
	hunter.applyMasterTactician()
	hunter.registerReadinessCD()

	hunter.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*1*float64(hunter.Talents.Surefooted))
	hunter.AddStat(stats.MeleeCrit, core.MeleeCritRatingPerCritChance*1*float64(hunter.Talents.KillerInstinct))
	hunter.AddStat(stats.Parry, core.ParryRatingPerParryChance*1*float64(hunter.Talents.Deflection))
	hunter.PseudoStats.RangedSpeedMultiplier *= 1 + 0.04*float64(hunter.Talents.SerpentsSwiftness)
	hunter.PseudoStats.RangedDamageDealtMultiplier *= 1 + 0.01*float64(hunter.Talents.RangedWeaponSpecialization)
	hunter.PseudoStats.BonusRangedCritRating += 1 * float64(hunter.Talents.LethalShots) * core.MeleeCritRatingPerCritChance

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

	hunter.applyKillCommand()
	hunter.registerRapidFireCD()
}

func (hunter *Hunter) critMultiplier(isRanged bool, target *core.Unit) float64 {
	primaryModifier := 1.0
	secondaryModifier := 0.0

	monsterMultiplier := 1.0 + 0.01*float64(hunter.Talents.MonsterSlaying)
	humanoidMultiplier := 1.0 + 0.01*float64(hunter.Talents.HumanoidSlaying)
	if target.MobType == proto.MobType_MobTypeBeast || target.MobType == proto.MobType_MobTypeGiant || target.MobType == proto.MobType_MobTypeDragonkin {
		primaryModifier *= monsterMultiplier
	} else if target.MobType == proto.MobType_MobTypeHumanoid {
		primaryModifier *= humanoidMultiplier
	}

	if isRanged {
		secondaryModifier += 0.06 * float64(hunter.Talents.MortalShots)
	}

	return hunter.MeleeCritMultiplier(primaryModifier, secondaryModifier)
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

func (hunter *Hunter) applyFerociousInspiration() {
	if hunter.pet == nil || hunter.Talents.FerociousInspiration == 0 {
		return
	}

	multiplier := 1.0 + 0.01*float64(hunter.Talents.FerociousInspiration)

	makeProcAura := func(character *core.Character) *core.Aura {
		return character.RegisterAura(core.Aura{
			Label:    "Ferocious Inspiration-" + strconv.Itoa(int(hunter.Index)),
			ActionID: core.ActionID{SpellID: 34460, Tag: int32(hunter.Index)},
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.DamageDealtMultiplier *= multiplier
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.DamageDealtMultiplier /= multiplier
			},
		})
	}

	var procAuras []*core.Aura
	hunter.RegisterAura(core.Aura{
		Label:    "Ferocious Inspiration",
		Duration: core.NeverExpires,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			procAuras = make([]*core.Aura, len(hunter.Party.PlayersAndPets))
			for i, playerOrPet := range hunter.Party.PlayersAndPets {
				procAuras[i] = makeProcAura(playerOrPet.GetCharacter())
			}
		},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			for _, procAura := range procAuras {
				procAura.Activate(sim)
			}
		},
	})
}

func (hunter *Hunter) registerBestialWrathCD() {
	if !hunter.Talents.BestialWrath {
		return
	}

	actionID := core.ActionID{SpellID: 19574}

	bestialWrathPetAura := hunter.pet.RegisterAura(core.Aura{
		Label:    "Bestial Wrath Pet",
		ActionID: actionID,
		Duration: time.Second * 18,
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
		Duration: time.Second * 18,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.1
			aura.Unit.PseudoStats.CostMultiplier *= 0.8
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
			aura.Unit.PseudoStats.CostMultiplier /= 0.8
		},
	})

	manaCost := hunter.BaseMana() * 0.1

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
				Duration: time.Minute * 2,
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
	if hunter.Talents.MonsterSlaying == 0 && hunter.Talents.HumanoidSlaying == 0 {
		return
	}

	monsterMultiplier := 1.0 + 0.01*float64(hunter.Talents.MonsterSlaying)
	humanoidMultiplier := 1.0 + 0.01*float64(hunter.Talents.HumanoidSlaying)

	switch hunter.CurrentTarget.MobType {
	case proto.MobType_MobTypeBeast, proto.MobType_MobTypeGiant, proto.MobType_MobTypeDragonkin:
		hunter.PseudoStats.DamageDealtMultiplier *= monsterMultiplier
	case proto.MobType_MobTypeHumanoid:
		hunter.PseudoStats.DamageDealtMultiplier *= humanoidMultiplier
	}
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
	critBonus := 2 * core.MeleeCritRatingPerCritChance * float64(hunter.Talents.MasterTactician)

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
