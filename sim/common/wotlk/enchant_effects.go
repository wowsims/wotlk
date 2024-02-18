package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func CreateBlackMagicProcAura(character *core.Character) *core.Aura {
	return character.NewTemporaryStatsAura("Black Magic Proc", core.ActionID{SpellID: 59626}, stats.Stats{stats.MeleeHaste: 250, stats.SpellHaste: 250}, time.Second*10)
}

func init() {
	// Keep these in order by item ID.

	core.NewEnchantEffect(3251, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForEnchant(3251)
		ppmm := character.AutoAttacks.NewPPMManager(4.0, procMask)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 44622},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 237, spell.OutcomeMagicHitAndCrit)
			},
		})

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Giant Slayer",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if result.Target.MobType != proto.MobType_MobTypeGiant {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Giant Slayer") {
					procSpell.Cast(sim, result.Target)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(3251, 4.0, &ppmm, aura)
	})

	core.NewEnchantEffect(3239, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForEnchant(3239)
		ppmm := character.AutoAttacks.NewPPMManager(4.0, procMask)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 44525},
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(185, 215), spell.OutcomeMagicHitAndCrit)
			},
		})

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Icebreaker",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Icebreaker") {
					procSpell.Cast(sim, result.Target)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(3239, 4.0, &ppmm, aura)
	})

	core.NewEnchantEffect(3607, func(agent core.Agent) {
		character := agent.GetCharacter()
		// TODO: This should be ranged-only haste. For now just make it hunter-only.
		if character.Class == proto.Class_ClassHunter {
			character.AddStats(stats.Stats{stats.MeleeHaste: 40, stats.SpellHaste: 40})
		}
	})

	core.NewEnchantEffect(3608, func(agent core.Agent) {
		agent.GetCharacter().AddBonusRangedCritRating(40)
	})

	core.NewEnchantEffect(3748, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 42500}

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultMeleeCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(45, 67)
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			},
		})

		aura := character.RegisterAura(core.Aura{
			Label:    "Titanium Shield Spike",
			ActionID: actionID,
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee) {
					procSpell.Cast(sim, spell.Unit)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(3748, aura)
	})

	core.NewEnchantEffect(3247, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
			character.PseudoStats.MobTypeAttackPower += 140
		}
	})

	core.NewEnchantEffect(3253, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 1.02
	})

	core.NewEnchantEffect(3296, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 0.98
	})

	core.NewEnchantEffect(3789, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForEnchant(3789)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		// Modify only gear armor, including from agility
		fivePercentOfArmor := (character.EquipStats()[stats.Armor] + 2.0*character.EquipStats()[stats.Agility]) * 0.05
		procAuraMH := character.NewTemporaryStatsAura("Berserking MH Proc", core.ActionID{SpellID: 59620, Tag: 1}, stats.Stats{stats.AttackPower: 400, stats.RangedAttackPower: 400, stats.Armor: -fivePercentOfArmor}, time.Second*15)
		procAuraOH := character.NewTemporaryStatsAura("Berserking OH Proc", core.ActionID{SpellID: 59620, Tag: 2}, stats.Stats{stats.AttackPower: 400, stats.RangedAttackPower: 400, stats.Armor: -fivePercentOfArmor}, time.Second*15)

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Berserking (Enchant)",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Berserking") {
					if spell.IsMH() {
						procAuraMH.Activate(sim)
					} else {
						procAuraOH.Activate(sim)
					}
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(3789, 1.0, &ppmm, aura)
	})

	// TODO: These are stand-in values without any real reference.
	core.NewEnchantEffect(3241, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForEnchant(3241)
		ppmm := character.AutoAttacks.NewPPMManager(3.0, procMask)

		healthMetrics := character.NewHealthMetrics(core.ActionID{ItemID: 44494})

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Lifeward",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Lifeward") {
					character.GainHealth(sim, 300*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(3241, 3.0, &ppmm, aura)
	})

	core.NewEnchantEffect(3790, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := CreateBlackMagicProcAura(character)
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 35,
		}
		procAura.Icd = &icd

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Black Magic",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				// Special case for spells that aren't spells that can proc black magic.
				specialCaseSpell := spell.ActionID.SpellID == 47465 || spell.ActionID.SpellID == 12867

				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskSpellDamage|core.ProcMaskWeaponProc) && !specialCaseSpell {
					return
				}

				if icd.IsReady(sim) && sim.RandomFloat("Black Magic") < 0.35 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(3790, aura)
	})

	core.AddWeaponEffect(3843, func(agent core.Agent, _ proto.ItemSlot) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 15
		w.BaseDamageMax += 15
	})

	core.NewEnchantEffect(3603, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 54757}

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 45,
				},
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 10,
				},
			},

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(1654, 2020), spell.OutcomeMagicCrit)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow, // Use low prio so other actives get used first.
			Type:     core.CooldownTypeDPS,
		})
	})

	core.NewEnchantEffect(3604, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 54758}

		procAura := character.NewTemporaryStatsAura("Hyperspeed Acceleration", actionID, stats.Stats{stats.MeleeHaste: 340, stats.SpellHaste: 340}, time.Second*12)

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical | core.SpellSchoolMagic,
			Flags:       core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 60,
				},
				// Shared CD with Offensive trinkets has been removed.
				// https://twitter.com/AggrendWoW/status/1579664462843633664
				// Change possibly temporary, but developers have confirmed it was intended.

				// SharedCD: core.Cooldown{
				// 	Timer:    character.GetOffensiveTrinketCD(),
				// 	Duration: time.Second * 12,
				// },
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				procAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow, // Use low prio so other actives get used first.
			Type:     core.CooldownTypeDPS,
		})
	})

	core.NewEnchantEffect(3722, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Lightweave Embroidery Proc", core.ActionID{SpellID: 55637}, stats.Stats{stats.SpellPower: 295}, time.Second*15)
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 60,
		}
		procAura.Icd = &icd

		callback := func(_ *core.Aura, sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if icd.IsReady(sim) && sim.RandomFloat("Lightweave") < 0.35 {
				icd.Use(sim)
				procAura.Activate(sim)
			}
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Lightweave Embroidery",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnHealDealt:           callback,
			OnPeriodicDamageDealt: callback,
			OnSpellHitDealt:       callback,
		})
	})

	core.NewEnchantEffect(3728, func(agent core.Agent) {
		character := agent.GetCharacter()
		if !character.HasManaBar() {
			return
		}

		manaMetrics := character.NewManaMetrics(core.ActionID{SpellID: 55767})
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 45,
		}

		character.GetOrRegisterAura(core.Aura{
			Icd:      &icd,
			Label:    "Darkglow Embroidery",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if icd.IsReady(sim) && sim.RandomFloat("Darkglow") < 0.35 {
					icd.Use(sim)
					character.AddMana(sim, 400, manaMetrics)
				}
			},
		})
	})

	core.NewEnchantEffect(3730, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Swordguard Embroidery Proc", core.ActionID{SpellID: 55775}, stats.Stats{stats.AttackPower: 400, stats.RangedAttackPower: 400}, time.Second*15)
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 55,
		}
		procAura.Icd = &icd

		character.GetOrRegisterAura(core.Aura{
			Label:    "Swordguard Embroidery",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}

				if icd.IsReady(sim) && sim.RandomFloat("Swordguard") < 0.2 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewEnchantEffect(3870, func(agent core.Agent) {
		character := agent.GetCharacter()
		healthMetrics := character.NewHealthMetrics(core.ActionID{SpellID: 64569})

		bloodReserveAura := character.GetOrRegisterAura(core.Aura{
			Label:     "Blood Reserve",
			ActionID:  core.ActionID{SpellID: 64568},
			Duration:  time.Second * 20,
			MaxStacks: 5,
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if character.CurrentHealth()/character.MaxHealth() < 0.35 {
					amountHealed := float64(aura.GetStacks()) * (360. + sim.RandomFloat("Blood Reserve")*80.)
					character.GainHealth(sim, amountHealed*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
					aura.Deactivate(sim)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Blood Draining",
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
			ProcMask:   core.ProcMaskMelee,
			Harmful:    true,
			ProcChance: 0.5,
			ICD:        time.Second * 10,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				if bloodReserveAura.IsActive() {
					bloodReserveAura.Refresh(sim)
					bloodReserveAura.AddStack(sim)
				} else {
					bloodReserveAura.Activate(sim)
					bloodReserveAura.SetStacks(sim, 1)
				}
			},
		})
	})
}
