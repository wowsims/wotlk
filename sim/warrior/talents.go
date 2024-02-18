package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) ToughnessArmorMultiplier() float64 {
	return 1.0 + 0.02*float64(warrior.Talents.Toughness)
}

func (warrior *Warrior) ApplyTalents() {
	warrior.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(warrior.Talents.Cruelty))
	warrior.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*1*float64(warrior.Talents.Precision))
	warrior.ApplyEquipScaling(stats.Armor, warrior.ToughnessArmorMultiplier())
	warrior.PseudoStats.BaseDodge += 0.01 * float64(warrior.Talents.Anticipation)
	warrior.PseudoStats.BaseParry += 0.01 * float64(warrior.Talents.Deflection)
	warrior.PseudoStats.DodgeReduction += 0.01 * float64(warrior.Talents.WeaponMastery)
	warrior.AutoAttacks.OHConfig().DamageMultiplier *= 1 + 0.05*float64(warrior.Talents.DualWieldSpecialization)

	if warrior.Talents.ArmoredToTheTeeth > 0 {
		coeff := float64(warrior.Talents.ArmoredToTheTeeth)
		warrior.AddStatDependency(stats.Armor, stats.AttackPower, coeff/108.0)
	}

	if warrior.Talents.StrengthOfArms > 0 {
		warrior.MultiplyStat(stats.Strength, 1.0+0.02*float64(warrior.Talents.StrengthOfArms))
		warrior.MultiplyStat(stats.Stamina, 1.0+0.02*float64(warrior.Talents.StrengthOfArms))
		warrior.AddStat(stats.Expertise, core.ExpertisePerQuarterPercentReduction*2*float64(warrior.Talents.StrengthOfArms))
	}

	// Shield Mastery, Shield Block, Glyph of Blocking, Eternal Earthsiege treated as additive sources
	if warrior.Talents.ShieldMastery > 0 {
		warrior.PseudoStats.BlockValueMultiplier += 0.15 * float64(warrior.Talents.ShieldMastery)
	}

	if warrior.Talents.Vitality > 0 {
		warrior.MultiplyStat(stats.Stamina, 1.0+0.03*float64(warrior.Talents.Vitality))
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
	warrior.applySwordAndBoard()
}

// Multiplicative with all other modifiers and only applies to the block damage event
func (warrior *Warrior) applyCriticalBlock() {
	if warrior.Talents.CriticalBlock == 0 {
		return
	}

	dummyCriticalBlockSpell := warrior.GetOrRegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 47296},
		Flags:    core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,
	})

	warrior.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if result.Outcome.Matches(core.OutcomeBlock) && !result.Outcome.Matches(core.OutcomeMiss) && !result.Outcome.Matches(core.OutcomeParry) && !result.Outcome.Matches(core.OutcomeDodge) {
			procChance := 0.2 * float64(warrior.Talents.CriticalBlock)
			if sim.RandomFloat("Critical Block Roll") <= procChance {
				blockValue := warrior.BlockValue()
				result.Damage = max(0, result.Damage-blockValue)
				dummyCriticalBlockSpell.Cast(sim, spell.Unit)
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
			baseDamage := coeff * warrior.BlockValue()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})

	core.MakePermanent(warrior.GetOrRegisterAura(core.Aura{
		Label:    "Damage Shield Trigger",
		Duration: core.NeverExpires,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() && !result.Outcome.Matches(core.OutcomeBlock) {
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
				warrior.LastAMTick = sim.CurrentTime
			},
		})
	})
}
func (warrior *Warrior) applyTasteForBlood() {
	if warrior.Talents.TasteForBlood == 0 {
		return
	}

	procChance := []float64{0, 0.33, 0.66, 1}[warrior.Talents.TasteForBlood]

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
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
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
			warrior.OverpowerAura.Duration = time.Second * 9
			warrior.OverpowerAura.Activate(sim)
			warrior.OverpowerAura.Duration = time.Second * 5
			warrior.lastOverpowerProc = sim.CurrentTime
		},
	})
}

func (warrior *Warrior) applyTrauma() {
	if warrior.Talents.Trauma == 0 {
		return
	}

	traumaAuras := warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.TraumaAura(target, int(warrior.Talents.Trauma))
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Trauma",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			if !spell.SpellSchool.Matches(core.SpellSchoolPhysical) || !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			proc := traumaAuras.Get(result.Target)
			proc.Duration = time.Minute * 1
			proc.Activate(sim)
		},
	})
}

func (warrior *Warrior) isBloodsurgeActive() bool {
	return warrior.BloodsurgeAura.IsActive() || (warrior.Talents.Bloodsurge > 0 && warrior.Ymirjar4pcProcAura.IsActive())
}

func (warrior *Warrior) applyBloodsurge() {
	if warrior.Talents.Bloodsurge == 0 {
		return
	}

	procChance := []float64{0, 0.07, 0.13, 0.2}[warrior.Talents.Bloodsurge]

	ymirjar4Set := warrior.HasSetBonus(ItemSetYmirjarLordsBattlegear, 4)

	warrior.BloodsurgeAura = warrior.RegisterAura(core.Aura{
		Label:    "Bloodsurge Proc",
		ActionID: core.ActionID{SpellID: 46916},
		Duration: time.Second * 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.Slam.DefaultCast.CastTime = 0
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.Slam.DefaultCast.CastTime = 1500 * time.Millisecond
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == warrior.Slam { // removed even if slam doesn't land
				aura.Deactivate(sim)
			}
		},
	})

	if ymirjar4Set {
		warrior.Ymirjar4pcProcAura = warrior.RegisterAura(core.Aura{
			Label:     "Ymirjar 4pc (Bloodsurge) Proc",
			ActionID:  core.ActionID{SpellID: 70847},
			Duration:  time.Second * 10,
			MaxStacks: 2,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if warrior.BloodsurgeAura.IsActive() {
					warrior.BloodsurgeAura.Deactivate(sim)
				}

				aura.SetStacks(sim, aura.MaxStacks)
				warrior.Slam.DefaultCast.CastTime = 0
				warrior.Slam.DefaultCast.GCD = core.GCDMin
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warrior.Slam.DefaultCast.CastTime = 1500 * time.Millisecond
				warrior.Slam.DefaultCast.GCD = core.GCDDefault
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == warrior.Slam {
					aura.RemoveStack(sim)
				}
			},
		})
	}

	warrior.RegisterAura(core.Aura{
		Label:    "Bloodsurge",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if !spell.Flags.Matches(SpellFlagBloodsurge) {
				return
			}

			if sim.RandomFloat("Bloodsurge") > procChance {
				return
			}

			warrior.lastBloodsurgeProc = sim.CurrentTime

			// as per https://www.wowhead.com/wotlk/spell=70847/item-warrior-t10-melee-4p-bonus#comments,
			//  the improved aura is not overwritten by the regular one, but simply refreshed
			if ymirjar4Set && (sim.RandomFloat("Ymirjar 4pc") < 0.2 || warrior.Ymirjar4pcProcAura.IsActive()) {
				warrior.BloodsurgeAura.Deactivate(sim)
				warrior.Ymirjar4pcProcAura.Activate(sim)

				warrior.BloodsurgeValidUntil = sim.CurrentTime + warrior.Ymirjar4pcProcAura.Duration
				return
			}

			warrior.BloodsurgeValidUntil = sim.CurrentTime + warrior.BloodsurgeAura.Duration
			warrior.BloodsurgeAura.Activate(sim)
		},
	})
}
func (warrior *Warrior) applyBloodFrenzy() {
	if warrior.Talents.BloodFrenzy == 0 {
		return
	}

	warrior.PseudoStats.MeleeSpeedMultiplier *= 1 + 0.05*float64(warrior.Talents.BloodFrenzy)

	bfAuras := warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.BloodFrenzyAura(target, warrior.Talents.BloodFrenzy)
	})
	warrior.Env.RegisterPreFinalizeEffect(func() {
		if warrior.Rend != nil {
			warrior.Rend.RelatedAuras = append(warrior.Rend.RelatedAuras, bfAuras)
		}
		if warrior.DeepWounds != nil {
			warrior.DeepWounds.RelatedAuras = append(warrior.DeepWounds.RelatedAuras, bfAuras)
		}
	})

	warrior.RegisterAura(core.Aura{
		Label:    "Blood Frenzy Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if spell == warrior.Rend || spell == warrior.DeepWounds {
				aura := bfAuras.Get(result.Target)
				dot := warrior.Rend.Dot(result.Target)
				aura.Duration = dot.TickLength * time.Duration(dot.NumberOfTicks)
				aura.Activate(sim)
			}
		},
	})
}

func (warrior *Warrior) applyTitansGrip() {
	if !warrior.Talents.TitansGrip {
		return
	}
	if !warrior.AutoAttacks.IsDualWielding {
		return
	}
	if warrior.MainHand().HandType != proto.HandType_HandTypeTwoHand && warrior.OffHand().HandType != proto.HandType_HandTypeTwoHand {
		return
	}

	warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1 - 0.1
}

func (warrior *Warrior) applyTwoHandedWeaponSpecialization() {
	if warrior.Talents.TwoHandedWeaponSpecialization == 0 {
		return
	}
	if warrior.MainHand().HandType != proto.HandType_HandTypeTwoHand {
		return
	}

	warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1 + 0.02*float64(warrior.Talents.TwoHandedWeaponSpecialization)
}

func (warrior *Warrior) applyOneHandedWeaponSpecialization() {
	if warrior.Talents.OneHandedWeaponSpecialization == 0 {
		return
	}
	if warrior.MainHand().HandType == proto.HandType_HandTypeTwoHand {
		return
	}

	warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1 + 0.02*float64(warrior.Talents.OneHandedWeaponSpecialization)
}

func (warrior *Warrior) applyWeaponSpecializations() {
	if ss := warrior.Talents.SwordSpecialization; ss > 0 {
		if mask := warrior.GetProcMaskForTypes(proto.WeaponType_WeaponTypeSword); mask != core.ProcMaskUnknown {
			warrior.registerSwordSpecialization(mask)
		}
	}

	if pas := warrior.Talents.PoleaxeSpecialization; pas > 0 {
		// the default character pane displays critical strike chance for main hand only
		switch warrior.GetProcMaskForTypes(proto.WeaponType_WeaponTypeAxe, proto.WeaponType_WeaponTypePolearm) {
		case core.ProcMaskMelee:
			warrior.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance*float64(pas))
		case core.ProcMaskMeleeMH:
			warrior.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance*float64(pas))
			warrior.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusCritRating -= 1 * core.CritRatingPerCritChance * float64(pas)
				}
			})
		case core.ProcMaskMeleeOH:
			warrior.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusCritRating += 1 * core.CritRatingPerCritChance * float64(pas)
				}
			})
		}
	}

	if ms := warrior.Talents.MaceSpecialization; ms > 0 {
		if mask := warrior.GetProcMaskForTypes(proto.WeaponType_WeaponTypeMace); mask != core.ProcMaskEmpty {
			warrior.AddStat(stats.ArmorPenetration, 3*core.ArmorPenPerPercentArmor*float64(ms))
		}
	}
}

func (warrior *Warrior) registerSwordSpecialization(procMask core.ProcMask) {
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
			config := *warrior.AutoAttacks.MHConfig()
			config.ActionID = core.ActionID{SpellID: 12281}
			swordSpecializationSpell = warrior.GetOrRegisterSpell(config)
		},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}
			if !spell.ProcMask.Matches(procMask) {
				return
			}
			if spell == warrior.WhirlwindOH {
				return // OH WW hits can't proc this
			}
			if !icd.IsReady(sim) {
				return
			}
			if sim.RandomFloat("Sword Specialization") < procChance {
				icd.Use(sim)
				aura.Unit.AutoAttacks.MaybeReplaceMHSwing(sim, swordSpecializationSpell).Cast(sim, result.Target)
			}
		},
	})
}

func (warrior *Warrior) applyUnbridledWrath() {
	if warrior.Talents.UnbridledWrath == 0 {
		return
	}

	ppmm := warrior.AutoAttacks.NewPPMManager(3*float64(warrior.Talents.UnbridledWrath), core.ProcMaskMeleeWhiteHit)

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 13002})

	warrior.RegisterAura(core.Aura{
		Label:    "Unbridled Wrath",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Damage == 0 {
				return
			}

			if ppmm.Proc(sim, spell.ProcMask, "Unbrided Wrath") {
				warrior.AddRage(sim, 1, rageMetrics)
			}
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) && !spell.Flags.Matches(SpellFlagWhirlwindOH) {
				return
			}

			if result.Outcome.Matches(core.OutcomeCrit) {
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
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= bonus
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= bonus
		},
	})
	warrior.RegisterAura(core.Aura{
		Label:    "Wrecking Crew",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			procAura.Activate(sim)
		},
	})
	core.RegisterPercentDamageModifierEffect(procAura, bonus)
}

func (warrior *Warrior) IsSuddenDeathActive() bool {
	return warrior.SuddenDeathAura.IsActive() || (warrior.Talents.SuddenDeath > 0 && warrior.Ymirjar4pcProcAura.IsActive())
}

func (warrior *Warrior) applySuddenDeath() {
	if warrior.Talents.SuddenDeath == 0 {
		return
	}

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 29724})

	minRageKept := []float64{0, 3, 7, 10}[warrior.Talents.SuddenDeath]
	procChance := []float64{0, 0.03, 0.06, 0.09}[warrior.Talents.SuddenDeath]

	warrior.SuddenDeathAura = warrior.RegisterAura(core.Aura{
		Label:    "Sudden Death Proc",
		ActionID: core.ActionID{SpellID: 29724},
		Duration: time.Second * 10,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || spell != warrior.Execute { // removed only when landed
				return
			}
			if rageRefund := minRageKept - warrior.CurrentRage(); rageRefund > 0 { // refund only when below minRageKept
				warrior.AddRage(sim, rageRefund, rageMetrics)
			}
			aura.Deactivate(sim)
		},
	})

	ymirjar4Set := warrior.HasSetBonus(ItemSetYmirjarLordsBattlegear, 4)

	if ymirjar4Set {
		warrior.Ymirjar4pcProcAura = warrior.RegisterAura(core.Aura{
			Label:     "Ymirjar 4pc (Sudden Death) Proc",
			ActionID:  core.ActionID{SpellID: 70847},
			Duration:  time.Second * 20,
			MaxStacks: 2,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.SetStacks(sim, aura.MaxStacks)
				warrior.Execute.DefaultCast.GCD = core.GCDMin
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warrior.Execute.DefaultCast.GCD = core.GCDDefault
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || spell != warrior.Execute {
					return
				}
				if rageRefund := minRageKept - warrior.CurrentRage(); rageRefund > 0 {
					warrior.AddRage(sim, rageRefund, rageMetrics)
				}
				aura.RemoveStack(sim)
			},
		})
	}

	warrior.RegisterAura(core.Aura{
		Label:    "Sudden Death",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if sim.RandomFloat("Sudden Death") > procChance {
				return
			}

			// as per https://www.wowhead.com/wotlk/spell=70847/item-warrior-t10-melee-4p-bonus#comments,
			//  the improved aura is not overwritten by the regular one, but simply refreshed
			if ymirjar4Set && (warrior.Ymirjar4pcProcAura.IsActive() || sim.RandomFloat("Ymirjar 4pc") < 0.2) {
				warrior.SuddenDeathAura.Deactivate(sim)

				warrior.Ymirjar4pcProcAura.Activate(sim)
				return
			}

			warrior.SuddenDeathAura.Activate(sim)
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
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeBlock | core.OutcomeDodge | core.OutcomeParry) {
				if sim.Proc(procChance, "Shield Specialization") {
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

	deathWishAura := warrior.RegisterAura(core.Aura{
		Label:    "Death Wish",
		ActionID: actionID,
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.2
			warrior.PseudoStats.DamageTakenMultiplier *= 1.05
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.2
			warrior.PseudoStats.DamageTakenMultiplier /= 1.05
		},
	})
	core.RegisterPercentDamageModifierEffect(deathWishAura, 1.2)

	deathWishSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		RageCost: core.RageCostOptions{
			Cost: 10,
		},
		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: warrior.intensifyRageCooldown(time.Minute * 3),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			deathWishAura.Activate(sim)
			warrior.WaitUntil(sim, sim.CurrentTime+core.GCDDefault)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: deathWishSpell,
		Type:  core.CooldownTypeDPS,
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
				Duration: core.TernaryDuration(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfLastStand), time.Minute*3, time.Minute*2),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(DefensiveStance)
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

func (warrior *Warrior) RegisterBladestormCD() {
	if !warrior.Talents.Bladestorm {
		return
	}

	actionID := core.ActionID{SpellID: 46924}
	numHits := min(4, warrior.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	if warrior.AutoAttacks.IsDualWielding {
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

	warrior.Bladestorm = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagChanneled | core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		RageCost: core.RageCostOptions{
			Cost: 25 - float64(warrior.Talents.FocusedRage),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: core.TernaryDuration(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfBladestorm), time.Second*75, time.Second*90),
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1.25,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Bladestorm",
			},
			NumberOfTicks: 6,
			TickLength:    time.Second * 1,
			OnTick: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot) {
				target := warrior.CurrentTarget
				spell := dot.Spell

				curTarget := target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					baseDamage := 0 +
						spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
						spell.BonusWeaponDamage()
					results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}

				curTarget = target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					spell.DealDamage(sim, results[hitIndex])
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}

				if warrior.BladestormOH != nil {
					curTarget = target
					for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
						baseDamage := 0 +
							spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
							spell.BonusWeaponDamage()
						results[hitIndex] = warrior.BladestormOH.CalcDamage(sim, curTarget, baseDamage, warrior.BladestormOH.OutcomeMeleeWeaponSpecialHitAndCrit)

						curTarget = sim.Environment.NextTargetUnit(curTarget)
					}

					curTarget = target
					for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
						warrior.BladestormOH.DealDamage(sim, results[hitIndex])
						curTarget = sim.Environment.NextTargetUnit(curTarget)
					}
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			dot := spell.AOEDot()
			dot.Apply(sim)
			dot.TickOnce(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: warrior.Bladestorm,
		Type:  core.CooldownTypeDPS,
	})
}

func (warrior *Warrior) applySwordAndBoard() {
	if warrior.Talents.SwordAndBoard == 0 {
		return
	}

	sabAura := warrior.GetOrRegisterAura(core.Aura{
		Label:    "Sword And Board",
		ActionID: core.ActionID{SpellID: 46953},
		Duration: 5 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.ShieldSlam.CostMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.ShieldSlam.CostMultiplier += 1
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == warrior.ShieldSlam {
				aura.Deactivate(sim)
			}
		},
	})

	procChance := 0.1 * float64(warrior.Talents.SwordAndBoard)
	core.MakePermanent(warrior.GetOrRegisterAura(core.Aura{
		Label: "Sword And Board Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if !(spell == warrior.Revenge || spell == warrior.Devastate) {
				return
			}

			if sim.RandomFloat("Sword And Board") < procChance {
				sabAura.Activate(sim)
				warrior.ShieldSlam.CD.Reset()
			}
		},
	}))
}
