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
	warrior.AddStat(stats.Dodge, core.DodgeRatingPerDodgeChance*1*float64(warrior.Talents.Anticipation))
	warrior.AddStat(stats.Armor, warrior.Equip.Stats()[stats.Armor]*0.02*float64(warrior.Talents.Toughness))
	warrior.PseudoStats.DodgeReduction += 0.01 * float64(warrior.Talents.WeaponMastery)
	warrior.AutoAttacks.OHConfig.DamageMultiplier *= 1 + 0.05*float64(warrior.Talents.DualWieldSpecialization)

	if warrior.Talents.ArmoredToTheTeeth > 0 {
		coeff := float64(warrior.Talents.ArmoredToTheTeeth)
		warrior.AddStatDependency(stats.Armor, stats.AttackPower, coeff/108.0)
	}

	if warrior.Talents.StrengthOfArms > 0 {
		warrior.MultiplyStat(stats.Strength, 1.0+0.01*float64(warrior.Talents.StrengthOfArms))
		warrior.AddStat(stats.Expertise, core.ExpertisePerQuarterPercentReduction*2*float64(warrior.Talents.StrengthOfArms))
	}

	// TODO: This should only be applied while berserker stance is active.
	if warrior.Talents.ImprovedBerserkerStance > 0 {
		warrior.MultiplyStat(stats.Strength, 1.0+0.04*float64(warrior.Talents.ImprovedBerserkerStance))
	}

	if warrior.Talents.ShieldMastery > 0 {
		warrior.MultiplyStat(stats.BlockValue, 1.0+0.15*float64(warrior.Talents.ShieldMastery))
	}

	if warrior.Talents.Vitality > 0 {
		warrior.MultiplyStat(stats.Stamina, 1.0+0.01*float64(warrior.Talents.Vitality))
		warrior.MultiplyStat(stats.Strength, 1.0+0.02*float64(warrior.Talents.Vitality))
		warrior.AddStat(stats.Expertise, core.ExpertisePerQuarterPercentReduction*2*float64(warrior.Talents.Vitality))
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
	warrior.RegisterBladestormCD()
	warrior.applyDamageShield()
	warrior.applyCriticalBlock()
}

func (warrior *Warrior) applyCriticalBlock() {
	if warrior.Talents.CriticalBlock == 0 {
		return
	}

	dummyCriticalBlockSpell := warrior.GetOrRegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 47296},
		Flags:    core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,
	})

	warrior.AddDynamicDamageTakenModifier(func(sim *core.Simulation, _ *core.Spell, spellEffect *core.SpellEffect) {
		if spellEffect.Outcome.Matches(core.OutcomeBlock) && !spellEffect.Outcome.Matches(core.OutcomeMiss) && !spellEffect.Outcome.Matches(core.OutcomeParry) && !spellEffect.Outcome.Matches(core.OutcomeDodge) {
			procChance := 0.2 * float64(warrior.Talents.CriticalBlock)
			if sim.RandomFloat("Critical Block Roll") <= procChance {
				blockValue := warrior.GetStat(stats.BlockValue)
				spellEffect.Damage = core.MaxFloat(0, spellEffect.Damage-blockValue)
				dummyCriticalBlockSpell.Cast(sim, warrior.CurrentTarget)
			}
		}
	})
}

func (warrior *Warrior) applyDamageShield() {
	if warrior.Talents.DamageShield == 0 {
		return
	}

	coeff := 0.1 * float64(warrior.Talents.DamageShield)
	damageShieldProcSpell := warrior.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 58874},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := coeff * warrior.GetStat(stats.BlockValue)
			spell.CalcAndDealDamageAlwaysHit(sim, target, baseDamage)
		},
	})

	core.MakePermanent(warrior.GetOrRegisterAura(core.Aura{
		Label:    "Damage Shield Trigger",
		Duration: core.NeverExpires,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() && !spellEffect.Outcome.Matches(core.OutcomeBlock) {
				return
			}

			if spell.SpellSchool != core.SpellSchoolPhysical {
				return
			}

			damageShieldProcSpell.Cast(sim, spell.Unit)
		},
	}))
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

			// Taste for Blood has 25% chance to not proc if ICD is expired during rend ticks. The chance is calculated from a controlled test here
			// https://classic.warcraftlogs.com/reports/2zcDnpNFXGaAPg34/#fight=last&type=damage-done&source=1
			if sim.RandomFloat("Taste for Blood bug") <= 0.25 && (sim.CurrentTime-warrior.lastTasteForBloodProc == time.Second*6) {
				return
			}
			icd.Use(sim)
			warrior.overpowerValidUntil = sim.CurrentTime + time.Second*9
			warrior.lastTasteForBloodProc = sim.CurrentTime
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
			proc.Activate(sim)
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

	Ymirjar4Set := warrior.HasSetBonus(ItemSetYmirjarLordsBattlegear, 4)

	warrior.BloodsurgeAura = warrior.RegisterAura(core.Aura{
		Label:    "Bloodsurge Proc",
		ActionID: core.ActionID{SpellID: 46916},
		Duration: time.Second * 5,
		// 2 stacks to accomodate T10 4 pc
		MaxStacks: 2,
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Bloodsurge",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell == warrior.Slam && warrior.BloodsurgeAura.IsActive() {
				warrior.BloodsurgeAura.RemoveStack(sim)
				return
			}

			if !spellEffect.Landed() {
				return
			}

			// Using heroic strike SpellID for now as Cleave and HS is a single spell variable
			if spell.ActionID.SpellID != 47450 && spell != warrior.Bloodthirst && spell != warrior.Whirlwind && spell != warrior.WhirlwindOH {
				return
			}

			if sim.RandomFloat("Bloodsurge") > procChance {
				return
			}

			warrior.BloodsurgeAura.Activate(sim)
			if Ymirjar4Set {
				if sim.RandomFloat("T10 4 set") < 0.2 {

					warrior.BloodsurgeAura.Duration = time.Second * 10
					warrior.BloodsurgeAura.SetStacks(sim, 2)
					warrior.Ymirjar4pcProcAura.Activate(sim)
					warrior.Ymirjar4pcProcAura.SetStacks(sim, 2)
					return
				}
			}

			if warrior.BloodsurgeAura.GetStacks() <= 1 {
				warrior.BloodsurgeAura.Duration = time.Second * 5
				warrior.BloodsurgeAura.SetStacks(sim, 1)
			}

			warrior.lastBloodsurgeProc = sim.CurrentTime
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
	if !warrior.AutoAttacks.IsDualWielding {
		return
	}
	if warrior.Equip[proto.ItemSlot_ItemSlotMainHand].HandType != proto.HandType_HandTypeTwoHand && warrior.Equip[proto.ItemSlot_ItemSlotOffHand].HandType != proto.HandType_HandTypeTwoHand {
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
	swordSpecMask := core.ProcMaskUnknown
	if weapon := warrior.Equip[proto.ItemSlot_ItemSlotMainHand]; weapon.ID != 0 {
		if weapon.WeaponType == proto.WeaponType_WeaponTypeAxe || weapon.WeaponType == proto.WeaponType_WeaponTypePolearm {
			warrior.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
					spell.BonusCritRating += 1 * core.CritRatingPerCritChance * float64(warrior.Talents.PoleaxeSpecialization)
				}
			})
		} else if weapon.WeaponType == proto.WeaponType_WeaponTypeMace {
			warrior.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
					spell.BonusArmorPenRating += 3 * core.ArmorPenPerPercentArmor * float64(warrior.Talents.MaceSpecialization)
				}
			})
		} else if weapon.WeaponType == proto.WeaponType_WeaponTypeSword {
			swordSpecMask |= core.ProcMaskMeleeMH
		}
	}
	if weapon := warrior.Equip[proto.ItemSlot_ItemSlotOffHand]; weapon.ID != 0 {
		if weapon.WeaponType == proto.WeaponType_WeaponTypeAxe || weapon.WeaponType == proto.WeaponType_WeaponTypePolearm {
			warrior.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusCritRating += 1 * core.CritRatingPerCritChance * float64(warrior.Talents.PoleaxeSpecialization)
				}
			})
		} else if weapon.WeaponType == proto.WeaponType_WeaponTypeMace {
			warrior.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusArmorPenRating += 3 * core.ArmorPenPerPercentArmor * float64(warrior.Talents.MaceSpecialization)
				}
			})
		} else if weapon.WeaponType == proto.WeaponType_WeaponTypeSword {
			swordSpecMask |= core.ProcMaskMeleeOH
		}
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
					ProcMask:    core.ProcMaskMeleeMHAuto,
					Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,

					DamageMultiplier: 1,
					CritMultiplier:   warrior.critMultiplier(mh),
					ThreatMultiplier: 1,

					ApplyEffects: warrior.AutoAttacks.MHConfig.ApplyEffects,
				})
			},
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}

				if !spell.ProcMask.Matches(swordSpecMask) {
					return
				}

				if spell == warrior.WhirlwindOH {
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

			if !spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			if !ppmm.Proc(sim, spell.ProcMask, "Unbrided Wrath") {
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
			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if spellEffect.Outcome.Matches(core.OutcomeCrit) {
				procAura.Activate(sim)
				procAura.SetStacks(sim, 3)
				return
			}

			// Remove a stack.
			if procAura.IsActive() && spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
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
			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
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
	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 29724})

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

	Ymirjar4Set := warrior.HasSetBonus(ItemSetYmirjarLordsBattlegear, 4)

	warrior.SuddenDeathAura = warrior.RegisterAura(core.Aura{
		Label:    "Sudden Death Proc",
		ActionID: core.ActionID{SpellID: 29724},
		Duration: time.Second * 10,
		// 2 stacks to accomodate T10 4 pc
		MaxStacks: 2,
	})
	warrior.RegisterAura(core.Aura{
		Label:    "Sudden Death",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if warrior.SuddenDeathAura.IsActive() && spell == warrior.Execute {
				warrior.SuddenDeathAura.RemoveStack(sim)
				warrior.AddRage(sim, rage_refund, rageMetrics)
			}

			if !spellEffect.Landed() {
				return
			}

			if spell.ProcMask.Matches(core.ProcMaskMelee) && sim.RandomFloat("Sudden Death") < procChance {
				warrior.SuddenDeathAura.Activate(sim)
				if Ymirjar4Set {
					if sim.RandomFloat("T10 4 set") < 0.2 {
						warrior.SuddenDeathAura.Activate(sim)
						warrior.SuddenDeathAura.Duration = time.Second * 20
						warrior.SuddenDeathAura.SetStacks(sim, 2)
						warrior.Ymirjar4pcProcAura.Activate(sim)
						warrior.Ymirjar4pcProcAura.SetStacks(sim, 2)
						return
					}
				}

				if warrior.SuddenDeathAura.GetStacks() <= 1 {
					warrior.SuddenDeathAura.Duration = time.Second * 10
					warrior.SuddenDeathAura.SetStacks(sim, 1)
				}
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
	rageAdded := float64(warrior.Talents.ShieldSpecialization)
	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 12727})

	warrior.RegisterAura(core.Aura{
		Label:    "Shield Specialization",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeBlock | core.OutcomeDodge | core.OutcomeParry) {
				if procChance == 1.0 || sim.RandomFloat("Shield Specialization") < procChance {
					warrior.AddRage(sim, rageAdded, rageMetrics)
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
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			lastStandAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: lastStandSpell,
		Type:  core.CooldownTypeSurvival,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return warrior.StanceMatches(DefensiveStance)
		},
	})
}

func (warrior *Warrior) RegisterBladestormCD() {
	if !warrior.Talents.Bladestorm {
		return
	}

	var bladestormDot *core.Dot
	actionID := core.ActionID{SpellID: 46924}
	cost := 25.0 - float64(warrior.Talents.FocusedRage)
	numHits := core.MinInt32(4, warrior.Env.GetNumTargets())

	var ohDamageEffects core.ApplySpellEffects
	if warrior.AutoAttacks.IsDualWielding {
		baseEffectOH := core.SpellEffect{
			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.OffHand, true, 0, true),
			OutcomeApplier: warrior.OutcomeFuncMeleeWeaponSpecialHitAndCrit(),
		}

		effects := make([]core.SpellEffect, 0, numHits)
		for i := int32(0); i < numHits; i++ {
			effect := baseEffectOH
			effect.Target = warrior.Env.GetTargetUnit(i)
			effects = append(effects, effect)
		}
		ohDamageEffects = core.ApplyEffectFuncDamageMultiple(effects)

		warrior.BladestormOH = warrior.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskMeleeOHSpecial,
			Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

			DamageMultiplier: 1 + 0.05*float64(warrior.Talents.DualWieldSpecialization),
			CritMultiplier:   warrior.critMultiplier(oh),
			ThreatMultiplier: 1.25,
		})
	}

	baseEffectMH := core.SpellEffect{
		BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, true, 0, true),
		OutcomeApplier: warrior.OutcomeFuncMeleeWeaponSpecialHitAndCrit(),
	}

	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		effect := baseEffectMH
		effect.Target = warrior.Env.GetTargetUnit(i)
		effects = append(effects, effect)
	}
	mhDamageEffects := core.ApplyEffectFuncDamageMultiple(effects)

	warrior.Bladestorm = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagChanneled | core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: core.TernaryDuration(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfBladestorm), time.Second*75, time.Second*90),
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1.25,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			bladestormDot.Apply(sim)
			bladestormDot.TickOnce()

			// Using regular cast/channel options would disable melee swings, so do it manually instead.
			warrior.SetGCDTimer(sim, sim.CurrentTime+time.Second*6)
			warrior.disableHsCleaveUntil = sim.CurrentTime + time.Second*6
		},
	})

	bladestormDot = core.NewDot(core.Dot{
		Spell: warrior.Bladestorm,
		Aura: warrior.RegisterAura(core.Aura{
			Label:    "Bladestorm",
			ActionID: actionID,
		}),
		NumberOfTicks: 6,
		TickLength:    time.Second * 1,
		TickEffects: core.TickFuncApplyEffects(func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			mhDamageEffects(sim, target, spell)
			if warrior.BladestormOH != nil {
				ohDamageEffects(sim, target, warrior.BladestormOH)
			}
		}),
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: warrior.Bladestorm,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return warrior.CurrentRage() >= cost
		},
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return true
		},
	})
}
