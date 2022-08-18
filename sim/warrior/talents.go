package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) ApplyTalents() {
	warrior.AddStat(stats.Parry, core.ParryRatingPerParryChance*1*float64(warrior.Talents.Deflection))
	warrior.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(warrior.Talents.Cruelty))
	warrior.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*1*float64(warrior.Talents.Precision))
	warrior.AddStat(stats.Defense, core.DefenseRatingPerDefense*4*float64(warrior.Talents.Anticipation))
	warrior.AddStat(stats.Armor, warrior.Equip.Stats()[stats.Armor]*0.02*float64(warrior.Talents.Toughness))
	warrior.AddStat(stats.Expertise, core.ExpertisePerQuarterPercentReduction*2*float64(warrior.Talents.StrengthOfArms))
	warrior.PseudoStats.DodgeReduction += 0.01 * float64(warrior.Talents.WeaponMastery)

	if warrior.Talents.DualWieldSpecialization > 0 {
		warrior.AutoAttacks.OHEffect.BaseDamage.Calculator = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 0, 1+0.05*float64(warrior.Talents.DualWieldSpecialization), 1.0, true)
	}

	if warrior.Talents.StrengthOfArms > 0 {
		warrior.MultiplyStat(stats.Strength, 1.0+0.01*float64(warrior.Talents.ImprovedBerserkerStance))
	}

	// TODO: This should only be applied while berserker stance is active.
	if warrior.Talents.ImprovedBerserkerStance > 0 {
		warrior.MultiplyStat(stats.Strength, 1.0+0.04*float64(warrior.Talents.ImprovedBerserkerStance))
	}

	if warrior.Talents.ShieldMastery > 0 {
		warrior.MultiplyStat(stats.BlockValue, 1.0+0.1*float64(warrior.Talents.ShieldMastery))
	}

	if warrior.Talents.Vitality > 0 {
		warrior.MultiplyStat(stats.Stamina, 1.0+0.01*float64(warrior.Talents.Vitality))
		warrior.MultiplyStat(stats.Strength, 1.0+0.02*float64(warrior.Talents.Vitality))
	}

	warrior.applyAngerManagement()
	warrior.applyDeepWounds()
	warrior.applyTitansGrip()
	warrior.applyOneHandedWeaponSpecialization()
	warrior.applyTwoHandedWeaponSpecialization()
	warrior.applyWeaponSpecializations()
	warrior.applyTrauma()
	warrior.applyBloodFrenzy()
	warrior.applyUnbridledWrath()
	warrior.applyFlurry()
	warrior.applyWreckingCrew()
	warrior.applyShieldSpecialization()
	warrior.registerDeathWishCD()
	warrior.registerSweepingStrikesCD()
	warrior.registerLastStandCD()
	warrior.applyTasteForBlood()
	warrior.applyBloodsurge()
	warrior.applySuddenDeath()
}

func (warrior *Warrior) applyAngerManagement() {
	if !warrior.Talents.AngerManagement {
		return
	}

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 12296})

	warrior.RegisterResetEffect(func(sim *core.Simulation) {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period: time.Second * 3,
			OnAction: func(sim *core.Simulation) {
				warrior.AddRage(sim, 1, rageMetrics)
			},
		})
	})
}
func (warrior *Warrior) applyTasteForBlood() {
	if warrior.Talents.TasteForBlood == 0 {
		return
	}

	var procChance float64
	if warrior.Talents.TasteForBlood == 1 {
		procChance = 0.33
	} else if warrior.Talents.TasteForBlood == 2 {
		procChance = 0.66
	} else if warrior.Talents.TasteForBlood == 3 {
		procChance = 1
	}

	icd := core.Cooldown{
		Timer:    warrior.NewTimer(),
		Duration: time.Second * 6,
	}

	warrior.RegisterAura(core.Aura{
		Label:    "Taste for Blood",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell != warrior.Rend {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			if sim.RandomFloat("Taste for Blood") > procChance {
				return
			}
			icd.Use(sim)
			warrior.overpowerValidUntil = sim.CurrentTime + time.Second*9
		},
	})
}

func (warrior *Warrior) applyTrauma() {
	if warrior.Talents.Trauma == 0 {
		return
	}

	for i := int32(0); i < warrior.Env.GetNumTargets(); i++ {
		target := warrior.Env.GetTargetUnit(i)
		warrior.TraumaAuras = append(warrior.TraumaAuras, core.TraumaAura(target, int(warrior.Talents.Trauma)))
	}

	warrior.RegisterAura(core.Aura{
		Label:    "Trauma",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			proc := warrior.TraumaAuras[spellEffect.Target.Index]
			proc.Duration = time.Minute * 1
			aura.Activate(sim)
		},
	})
}

func (warrior *Warrior) applyBloodsurge() {
	if warrior.Talents.Bloodsurge == 0 {
		return
	}
	procChance := 0.0

	if warrior.Talents.Bloodsurge == 1 {
		procChance = 0.07
	} else if warrior.Talents.Bloodsurge == 2 {
		procChance = 0.13
	} else if warrior.Talents.Bloodsurge == 3 {
		procChance = 0.20
	}

	warrior.BloodsurgeAura = warrior.RegisterAura(core.Aura{
		Label:    "Bloodsurge Proc",
		ActionID: core.ActionID{SpellID: 46916},
		Duration: time.Second * 5,
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Bloodsurge",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell == warrior.Slam && warrior.BloodsurgeAura.IsActive() {
				warrior.BloodsurgeAura.Deactivate(sim)
				return
			}

			if !spellEffect.Landed() {
				return
			}

			// Using heroic strike SpellID for now as Cleave and HS is a single spell variable
			if spell.ActionID.SpellID != 47450 && spell != warrior.Bloodthirst && spell != warrior.Whirlwind {
				return
			}

			if sim.RandomFloat("Bloodsurge") > procChance {
				return
			}

			warrior.BloodsurgeAura.Activate(sim)
		},
	})
}
func (warrior *Warrior) applyBloodFrenzy() {
	if warrior.Talents.BloodFrenzy == 0 {
		return
	}

	for i := int32(0); i < warrior.Env.GetNumTargets(); i++ {
		target := warrior.Env.GetTargetUnit(i)
		warrior.BloodFrenzyAuras = append(warrior.BloodFrenzyAuras, core.BloodFrenzyAura(target, warrior.Talents.BloodFrenzy))
	}

	warrior.PseudoStats.MeleeSpeedMultiplier *= 1 + 0.05*float64(warrior.Talents.BloodFrenzy)
}

func (warrior *Warrior) procBloodFrenzy(sim *core.Simulation, effect *core.SpellEffect, dur time.Duration) {
	if warrior.Talents.BloodFrenzy == 0 {
		return
	}

	aura := warrior.BloodFrenzyAuras[effect.Target.Index]
	aura.Duration = dur
	aura.Activate(sim)
}

func (warrior *Warrior) applyTitansGrip() {
	if !warrior.Talents.TitansGrip {
		return
	}
	if warrior.Equip[proto.ItemSlot_ItemSlotOffHand].HandType != proto.HandType_HandTypeTwoHand {
		return
	}

	warrior.PseudoStats.PhysicalDamageDealtMultiplier *= 1 - 0.1
}

func (warrior *Warrior) applyTwoHandedWeaponSpecialization() {
	if warrior.Talents.TwoHandedWeaponSpecialization == 0 {
		return
	}
	if warrior.Equip[proto.ItemSlot_ItemSlotMainHand].HandType != proto.HandType_HandTypeTwoHand {
		return
	}

	warrior.PseudoStats.PhysicalDamageDealtMultiplier *= 1 + 0.02*float64(warrior.Talents.TwoHandedWeaponSpecialization)
}

func (warrior *Warrior) applyOneHandedWeaponSpecialization() {
	if warrior.Talents.OneHandedWeaponSpecialization == 0 {
		return
	}
	if warrior.Equip[proto.ItemSlot_ItemSlotMainHand].HandType == proto.HandType_HandTypeTwoHand {
		return
	}

	warrior.PseudoStats.PhysicalDamageDealtMultiplier *= 1 + 0.02*float64(warrior.Talents.OneHandedWeaponSpecialization)
}

func (warrior *Warrior) applyWeaponSpecializations() {
	maceSpecMask := core.ProcMaskUnknown
	swordSpecMask := core.ProcMaskUnknown
	if weapon := warrior.Equip[proto.ItemSlot_ItemSlotMainHand]; weapon.ID != 0 {
		if weapon.WeaponType == proto.WeaponType_WeaponTypeAxe || weapon.WeaponType == proto.WeaponType_WeaponTypePolearm {
			warrior.PseudoStats.BonusMHCritRating += 1 * core.CritRatingPerCritChance * float64(warrior.Talents.PoleaxeSpecialization)
		} else if weapon.WeaponType == proto.WeaponType_WeaponTypeMace {
			maceSpecMask |= core.ProcMaskMeleeMH
		} else if weapon.WeaponType == proto.WeaponType_WeaponTypeSword {
			swordSpecMask |= core.ProcMaskMeleeMH
		}
	}
	if weapon := warrior.Equip[proto.ItemSlot_ItemSlotOffHand]; weapon.ID != 0 {
		if weapon.WeaponType == proto.WeaponType_WeaponTypeAxe || weapon.WeaponType == proto.WeaponType_WeaponTypePolearm {
			warrior.PseudoStats.BonusOHCritRating += 1 * core.CritRatingPerCritChance * float64(warrior.Talents.PoleaxeSpecialization)
		} else if weapon.WeaponType == proto.WeaponType_WeaponTypeMace {
			maceSpecMask |= core.ProcMaskMeleeOH
		} else if weapon.WeaponType == proto.WeaponType_WeaponTypeSword {
			swordSpecMask |= core.ProcMaskMeleeOH
		}
	}

	if warrior.Talents.MaceSpecialization > 0 && maceSpecMask != core.ProcMaskUnknown {
		ppmm := warrior.AutoAttacks.NewPPMManager(1.5, maceSpecMask)
		rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 12704})

		warrior.RegisterAura(core.Aura{
			Label:    "Mace Specialization",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}

				if !spellEffect.ProcMask.Matches(maceSpecMask) {
					return
				}

				if spell == warrior.Whirlwind && spellEffect.ProcMask.Matches(core.ProcMaskMeleeOHSpecial) {
					// OH WW hits cant proc this
					return
				}

				if !ppmm.Proc(sim, spellEffect.ProcMask, "Mace Specialization") {
					return
				}

				warrior.AddRage(sim, 7, rageMetrics)
			},
		})
	}

	if warrior.Talents.SwordSpecialization > 0 && swordSpecMask != core.ProcMaskUnknown {
		var swordSpecializationSpell *core.Spell
		icd := core.Cooldown{
			Timer:    warrior.NewTimer(),
			Duration: time.Second * 6,
		}
		procChance := 0.02 * float64(warrior.Talents.SwordSpecialization)

		warrior.RegisterAura(core.Aura{
			Label:    "Sword Specialization",
			Duration: core.NeverExpires,
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				swordSpecializationSpell = warrior.GetOrRegisterSpell(core.SpellConfig{
					ActionID:    core.ActionID{SpellID: 12281},
					SpellSchool: core.SpellSchoolPhysical,
					Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

					ApplyEffects: core.ApplyEffectFuncDirectDamage(warrior.AutoAttacks.MHEffect),
				})
			},
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}

				if !spellEffect.ProcMask.Matches(swordSpecMask) {
					return
				}

				if spell == warrior.Whirlwind && spellEffect.ProcMask.Matches(core.ProcMaskMeleeOHSpecial) {
					// OH WW hits cant proc this
					return
				}

				if !icd.IsReady(sim) {
					return
				}

				if sim.RandomFloat("Sword Specialization") > procChance {
					return
				}
				icd.Use(sim)

				aura.Unit.AutoAttacks.MaybeReplaceMHSwing(sim, swordSpecializationSpell).Cast(sim, spellEffect.Target)
			},
		})
	}
}

func (warrior *Warrior) applyUnbridledWrath() {
	if warrior.Talents.UnbridledWrath == 0 {
		return
	}

	ppmm := warrior.AutoAttacks.NewPPMManager(3*float64(warrior.Talents.UnbridledWrath), core.ProcMaskMelee)
	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 13002})

	warrior.RegisterAura(core.Aura{
		Label:    "Unbridled Wrath",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Damage == 0 {
				return
			}

			if !spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			if !ppmm.Proc(sim, spellEffect.ProcMask, "Unbrided Wrath") {
				return
			}

			warrior.AddRage(sim, 1, rageMetrics)
		},
	})
}

func (warrior *Warrior) applyFlurry() {
	if warrior.Talents.Flurry == 0 {
		return
	}

	bonus := 1 + 0.05*float64(warrior.Talents.Flurry)
	inverseBonus := 1 / bonus

	procAura := warrior.RegisterAura(core.Aura{
		Label:     "Flurry Proc",
		ActionID:  core.ActionID{SpellID: 12974},
		Duration:  core.NeverExpires,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.MultiplyMeleeSpeed(sim, bonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.MultiplyMeleeSpeed(sim, inverseBonus)
		},
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Flurry",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if spellEffect.Outcome.Matches(core.OutcomeCrit) {
				procAura.Activate(sim)
				procAura.SetStacks(sim, 3)
				return
			}

			// Remove a stack.
			if procAura.IsActive() && spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				procAura.RemoveStack(sim)
			}
		},
	})
}

func (warrior *Warrior) applyWreckingCrew() {
	if warrior.Talents.WreckingCrew == 0 {
		return
	}

	bonus := 1 + 0.02*float64(warrior.Talents.WreckingCrew)

	procAura := warrior.RegisterAura(core.Aura{
		Label:    "Enrage",
		ActionID: core.ActionID{SpellID: 57518},
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= bonus
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= bonus
		},
	})
	warrior.RegisterAura(core.Aura{
		Label:    "Wrecking Crew",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			procAura.Activate(sim)
		},
	})
}

func (warrior *Warrior) applySuddenDeath() {
	if warrior.Talents.SuddenDeath == 0 {
		return
	}

	var rage_refund float64
	var procChance float64

	if warrior.Talents.SuddenDeath == 1 {
		rage_refund = 3.0
		procChance = 0.03
	} else if warrior.Talents.SuddenDeath == 2 {
		rage_refund = 7.0
		procChance = 0.06
	} else if warrior.Talents.SuddenDeath == 3 {
		rage_refund = 10.0
		procChance = 0.09
	}

	warrior.SuddenDeathAura = warrior.RegisterAura(core.Aura{
		Label:    "Sudden Death Proc",
		ActionID: core.ActionID{SpellID: 29724},
		Duration: core.NeverExpires,
	})
	warrior.RegisterAura(core.Aura{
		Label:    "Sudden Death",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}

			if spellEffect.ProcMask.Matches(core.ProcMaskMelee) && sim.RandomFloat("Sudden Death") < procChance {
				warrior.SuddenDeathAura.Activate(sim)
			}

			if warrior.SuddenDeathAura.IsActive() && spell == warrior.Execute {
				warrior.SuddenDeathAura.Deactivate(sim)
				warrior.AddRage(sim, rage_refund, warrior.RageRefundMetrics)
			}
		},
	})
}

func (warrior *Warrior) applyShieldSpecialization() {
	if warrior.Talents.ShieldSpecialization == 0 {
		return
	}

	warrior.AddStat(stats.Block, core.BlockRatingPerBlockChance*1*float64(warrior.Talents.ShieldSpecialization))

	procChance := 0.2 * float64(warrior.Talents.ShieldSpecialization)
	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 12727})

	warrior.RegisterAura(core.Aura{
		Label:    "Shield Specialization",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeBlock) {
				if procChance == 1 || sim.RandomFloat("Shield Specialization") < procChance {
					warrior.AddRage(sim, 1, rageMetrics)
				}
			}
		},
	})
}

func (warrior *Warrior) registerDeathWishCD() {
	if !warrior.Talents.DeathWish {
		return
	}

	actionID := core.ActionID{SpellID: 12292}
	const hasteBonus = 1.2
	const inverseHasteBonus = 1 / 1.2

	deathWishAura := warrior.RegisterAura(core.Aura{
		Label:    "Death Wish",
		ActionID: actionID,
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.PhysicalDamageDealtMultiplier *= 1.2
			warrior.PseudoStats.DamageTakenMultiplier *= 1.05
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.PhysicalDamageDealtMultiplier /= 1.2
			warrior.PseudoStats.DamageTakenMultiplier /= 1.05
		},
	})

	cost := 10.0
	cooldownDur := time.Minute * 3
	if warrior.Talents.IntensifyRage == 1 {
		cooldownDur = time.Duration(float64(cooldownDur) * 0.89)
	} else if warrior.Talents.IntensifyRage == 2 {
		cooldownDur = time.Duration(float64(cooldownDur) * 0.78)
	} else if warrior.Talents.IntensifyRage == 3 {
		cooldownDur = time.Duration(float64(cooldownDur) * 0.67)
	}
	deathWishSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownDur,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			deathWishAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: deathWishSpell,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return warrior.CurrentRage() >= cost
		},
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return true
		},
	})
}

func (warrior *Warrior) registerLastStandCD() {
	if !warrior.Talents.LastStand {
		return
	}

	actionID := core.ActionID{SpellID: 12975}
	healthMetrics := warrior.NewHealthMetrics(actionID)

	var bonusHealth float64
	lastStandAura := warrior.RegisterAura(core.Aura{
		Label:    "Last Stand",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bonusHealth = warrior.MaxHealth() * 0.3
			warrior.AddStatsDynamic(sim, stats.Stats{stats.Health: bonusHealth})
			warrior.GainHealth(sim, bonusHealth, healthMetrics)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatsDynamic(sim, stats.Stats{stats.Health: -bonusHealth})
		},
	})

	lastStandSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute * 8,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			lastStandAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: lastStandSpell,
		Type:  core.CooldownTypeSurvival,
	})
}
