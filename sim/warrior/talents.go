package warrior

import (
	"time"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

func (warrior *Warrior) ToughnessArmorMultiplier() float64 {
	return 1.0 + 0.02*float64(warrior.Talents.Toughness)
}

func (warrior *Warrior) ApplyTalents() {
	warrior.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(warrior.Talents.Cruelty))
	warrior.ApplyEquipScaling(stats.Armor, warrior.ToughnessArmorMultiplier())
	warrior.PseudoStats.BaseDodge += 0.01 * float64(warrior.Talents.Anticipation)
	warrior.PseudoStats.BaseParry += 0.01 * float64(warrior.Talents.Deflection)
	warrior.AutoAttacks.OHConfig().DamageMultiplier *= 1 + 0.05*float64(warrior.Talents.DualWieldSpecialization)

	warrior.applyAngerManagement()
	warrior.applyDeepWounds()
	warrior.applyOneHandedWeaponSpecialization()
	warrior.applyTwoHandedWeaponSpecialization()
	warrior.applyWeaponSpecializations()
	warrior.applyUnbridledWrath()
	warrior.applyFlurry()
	warrior.applyShieldSpecialization()
	warrior.registerDeathWishCD()
	warrior.registerSweepingStrikesCD()
	warrior.registerLastStandCD()
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

func (warrior *Warrior) applyTwoHandedWeaponSpecialization() {
	if warrior.Talents.TwoHandedWeaponSpecialization == 0 {
		return
	}
	if warrior.MainHand().HandType != proto.HandType_HandTypeTwoHand {
		return
	}

	warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1 + 0.01*float64(warrior.Talents.TwoHandedWeaponSpecialization)
}

func (warrior *Warrior) applyOneHandedWeaponSpecialization() {
	if warrior.Talents.OneHandedWeaponSpecialization == 0 {
		return
	}
	if warrior.MainHand().HandType == proto.HandType_HandTypeTwoHand {
		return
	}

	warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1 + 0.01*float64(warrior.Talents.OneHandedWeaponSpecialization)
}

func (warrior *Warrior) applyWeaponSpecializations() {
	if ss := warrior.Talents.SwordSpecialization; ss > 0 {
		if mask := warrior.GetProcMaskForTypes(proto.WeaponType_WeaponTypeSword); mask != core.ProcMaskUnknown {
			warrior.registerSwordSpecialization(mask)
		}
	}

	if as := warrior.Talents.AxeSpecialization; as > 0 {
		// the default character panel displays critical strike chance for main hand only
		switch warrior.GetProcMaskForTypes(proto.WeaponType_WeaponTypeAxe) {
		case core.ProcMaskMelee:
			warrior.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance*float64(as))
		case core.ProcMaskMeleeMH:
			warrior.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance*float64(as))
			warrior.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusCritRating -= 1 * core.CritRatingPerCritChance * float64(as)
				}
			})
		case core.ProcMaskMeleeOH:
			warrior.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
					spell.BonusCritRating += 1 * core.CritRatingPerCritChance * float64(as)
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
		Duration: time.Millisecond * 200,
	}
	procChance := 0.01 * float64(warrior.Talents.SwordSpecialization)

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
			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
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
			warrior.PseudoStats.ArmorMultiplier *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.2
			warrior.PseudoStats.ArmorMultiplier /= 1.2
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
				Duration: time.Minute * 3,
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
				Duration: time.Minute * 10,
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
