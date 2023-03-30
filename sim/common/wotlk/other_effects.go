package wotlk

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func init() {
	core.NewItemEffect(37220, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 37220}

		procAura := character.RegisterAura(core.Aura{
			Label:    "Essence of Gossamer",
			ActionID: actionID,
			Duration: time.Second * 10,
		})

		character.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if procAura.IsActive() {
				result.Damage = core.MaxFloat(0, result.Damage-140)
			}
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Essence of Gossamer Trigger",
			Callback:   core.CallbackOnSpellHitTaken,
			ProcMask:   core.ProcMaskMelee,
			Harmful:    true,
			ProcChance: 0.05,
			ICD:        time.Second * 50,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			},
		})
	})
	core.NewItemEffect(45507, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 45507}

		procAura := character.RegisterAura(core.Aura{
			Label:    "The General's Heart",
			ActionID: actionID,
			Duration: time.Second * 10,
		})

		character.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if procAura.IsActive() {
				result.Damage = core.MaxFloat(0, result.Damage-205)
			}
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "The General's Heart Trigger",
			Callback:   core.CallbackOnSpellHitTaken,
			ProcMask:   core.ProcMaskMelee,
			Harmful:    true,
			ProcChance: 0.05,
			ICD:        time.Second * 50,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(37734, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 37734}

		procAuras := make([]*core.Aura, len(character.Env.AllUnits))
		for _, target := range character.Env.AllUnits {
			if !character.IsOpponent(target) {
				procAuras[target.UnitIndex] = target.GetOrRegisterAura(core.Aura{
					Label:     "Touched by a Troll",
					ActionID:  core.ActionID{SpellID: 60518},
					Duration:  time.Second * 10,
					MaxStacks: 5,
					OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
						aura.Unit.PseudoStats.BonusHealingTaken += 58 * float64(newStacks-oldStacks)
					},
				})
			}
		}

		activeAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Talisman of Troll Divinity",
			ActionID: actionID,
			Duration: time.Second * 20,
			Callback: core.CallbackOnHealDealt,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				aura := procAuras[result.Target.UnitIndex]
				aura.Activate(sim)
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				activeAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	NewItemEffectWithHeroic(func(isHeroic bool) {
		name := "Deathbringer's Will"
		itemID := int32(50362)
		amount := 600.0
		auraIDs := []int32{
			71484,
			71485,
			71486,
			71491,
			71492,
		}
		if isHeroic {
			name += " H"
			itemID = 50363
			amount = 700
			auraIDs = []int32{
				71561,
				71556,
				71558,
				71559,
				71560,
			}
		}

		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			strAura := character.NewTemporaryStatsAura(name+" Strength Proc", core.ActionID{SpellID: auraIDs[0]}, stats.Stats{stats.Strength: amount}, time.Second*30)
			agiAura := character.NewTemporaryStatsAura(name+" Agility Proc", core.ActionID{SpellID: auraIDs[1]}, stats.Stats{stats.Agility: amount}, time.Second*30)
			apAura := character.NewTemporaryStatsAura(name+" AP Proc", core.ActionID{SpellID: auraIDs[2]}, stats.Stats{stats.AttackPower: amount * 2, stats.RangedAttackPower: amount * 2}, time.Second*30)
			critAura := character.NewTemporaryStatsAura(name+" Crit Proc", core.ActionID{SpellID: auraIDs[3]}, stats.Stats{stats.MeleeCrit: amount, stats.SpellCrit: amount}, time.Second*30)
			hasteAura := character.NewTemporaryStatsAura(name+" Haste Proc", core.ActionID{SpellID: auraIDs[4]}, stats.Stats{stats.MeleeHaste: amount, stats.SpellHaste: amount}, time.Second*30)

			var auras [3]*core.Aura
			switch character.Class {
			case proto.Class_ClassDeathknight:
				auras = [3]*core.Aura{strAura, critAura, hasteAura}
			case proto.Class_ClassDruid:
				auras = [3]*core.Aura{strAura, agiAura, hasteAura}
			case proto.Class_ClassHunter:
				auras = [3]*core.Aura{agiAura, critAura, apAura}
			case proto.Class_ClassPaladin:
				auras = [3]*core.Aura{strAura, critAura, hasteAura}
			case proto.Class_ClassRogue:
				auras = [3]*core.Aura{agiAura, apAura, hasteAura}
			case proto.Class_ClassShaman:
				auras = [3]*core.Aura{agiAura, apAura, hasteAura}
			case proto.Class_ClassWarrior:
				auras = [3]*core.Aura{strAura, critAura, hasteAura}
			default:
				return
			}

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       name,
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrRanged,
				Outcome:    core.OutcomeLanded,
				ProcChance: 0.35,
				ICD:        time.Second * 105,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					rand := sim.RandomFloat("Deathbringer's Will")
					if rand < 1.0/3.0 {
						auras[0].Activate(sim)
					} else if rand < 2.0/3.0 {
						auras[1].Activate(sim)
					} else {
						auras[2].Activate(sim)
					}
				},
			})
		})
	})

	core.NewItemEffect(40258, func(agent core.Agent) {
		character := agent.GetCharacter()

		healSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 40258},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskSpellHealing,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			Hot: core.DotConfig{
				Aura: core.Aura{
					Label: "Forethought Talisman",
				},
				NumberOfTicks: 4,
				TickLength:    time.Second * 3,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
					dot.SnapshotBaseDamage = 3752.0 / 4
					dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
				},
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Forethought Talisman",
			Callback:   core.CallbackOnHealDealt,
			Outcome:    core.OutcomeCrit,
			ProcChance: 0.2,
			ICD:        time.Second * 45,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				healSpell.Hot(result.Target).Apply(sim)
			},
		})
	})

	core.NewItemEffect(40382, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 40382}
		manaMetrics := character.NewManaMetrics(actionID)

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Soul of the Dead",
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
			ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
			Outcome:    core.OutcomeCrit,
			ProcChance: 0.25,
			ICD:        time.Second * 45,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				character.AddMana(sim, 900, manaMetrics)
			},
		})
	})

	core.NewItemEffect(41121, func(agent core.Agent) { // Gnomish Lightning Generator
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 41121}

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 1,
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
				spell.CalcAndDealDamage(sim, target, sim.Roll(1530, 1870), spell.OutcomeMagicHitAndCrit)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow, // Use low prio so other actives get used first.
			Type:     core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(46017, func(agent core.Agent) { // Val'anyr
		character := agent.GetCharacter()

		shieldID := core.ActionID{SpellID: 64413}
		shields := core.NewAllyShieldArray(
			&character.Unit,
			core.Shield{
				Spell: character.GetOrRegisterSpell(core.SpellConfig{
					ActionID:    shieldID,
					SpellSchool: core.SpellSchoolNature,
					ProcMask:    core.ProcMaskSpellHealing,
					Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

					DamageMultiplier: 1,
					ThreatMultiplier: 1,
				}),
			},
			core.Aura{
				Label:    "Val'anyr Shield",
				ActionID: shieldID,
				Duration: time.Second * 30,
			})

		activeAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Blessing of Ancient Kings",
			ActionID: core.ActionID{SpellID: 64411},
			Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
			Duration: time.Second * 15,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				// TODO: Shield needs to stack with itself up to 20k.
				shield := shields[result.Target.UnitIndex]
				shield.Apply(sim, result.Damage*0.15)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Val'anyr, Hammer of Ancient Kings Trigger",
			Callback:   core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
			Harmful:    true, // Better name for this would be, 'nonzero'
			ProcChance: 0.1,
			ICD:        time.Second * 45,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				activeAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(50259, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 50259}

		activeAura := core.MakeStackingAura(character, core.StackingStatAura{
			Aura: core.Aura{
				Label:     "Nevermelting Ice Crystal",
				ActionID:  actionID,
				Duration:  time.Second * 20,
				MaxStacks: 5,
			},
			BonusPerStack: stats.Stats{stats.MeleeCrit: 184, stats.SpellCrit: 184},
		})

		core.ApplyProcTriggerCallback(&character.Unit, activeAura, core.ProcTrigger{
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskSpellDamage,
			Outcome:  core.OutcomeCrit,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				activeAura.RemoveStack(sim)
			},
		})

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 20,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				activeAura.Activate(sim)
				activeAura.SetStacks(sim, 5)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})

	NewItemEffectWithHeroic(func(isHeroic bool) { // Sliver of Pure Ice
		itemID := int32(50339)
		amount := 1625.0
		if isHeroic {
			itemID = 50346
			amount = 1830.0
		}

		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{ItemID: itemID}
			manaMetrics := character.NewManaMetrics(actionID)

			spell := character.RegisterSpell(core.SpellConfig{
				ActionID: actionID,
				Flags:    core.SpellFlagNoOnCastComplete,
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
					character.AddMana(sim, amount, manaMetrics)
				},
			})
			character.AddMajorCooldown(core.MajorCooldown{
				Type:  core.CooldownTypeMana,
				Spell: spell,
			})
		})
	})

	NewItemEffectWithHeroic(func(isHeroic bool) {
		name := "Corpse Tongue Coin"
		itemID := int32(50352)
		amount := 5712.0
		if isHeroic {
			name += " H"
			itemID = 50349
			amount = 6426.0
		}

		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{ItemID: itemID}

			procAura := character.NewTemporaryStatsAura(name+" Proc", actionID, stats.Stats{stats.Armor: amount}, time.Second*10)

			// Handle ICD ourselves since we use a custom check.
			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Second * 30,
			}

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:     name + " Trigger",
				Callback: core.CallbackOnSpellHitTaken,
				ProcMask: core.ProcMaskMelee,
				Harmful:  true,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					if icd.IsReady(sim) && character.CurrentHealthPercent() < 0.35 {
						icd.Use(sim)
						procAura.Activate(sim)
					}
				},
			})
		})
	})

	NewItemEffectWithHeroic(func(isHeroic bool) {
		name := "Dislodged Foreign Object"
		itemID := int32(50353)
		amount := 105.0
		if isHeroic {
			name += " H"
			itemID = 50348
			amount = 121.0
		}

		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{ItemID: itemID}

			procAura := core.MakeStackingAura(character, core.StackingStatAura{
				Aura: core.Aura{
					Label:     name,
					ActionID:  actionID,
					Duration:  time.Second * 20,
					MaxStacks: 10,
				},
				BonusPerStack: stats.Stats{stats.SpellPower: amount},
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       name + " Trigger",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskSpellDamage,
				Harmful:    true,
				ProcChance: 0.10,
				ICD:        time.Second * 45,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					procAura.Activate(sim)
					core.StartPeriodicAction(sim, core.PeriodicActionOptions{
						NumTicks:        10,
						Period:          time.Second * 2,
						TickImmediately: true,
						OnAction: func(sim *core.Simulation) {
							if procAura.IsActive() {
								procAura.AddStack(sim)
							}
						},
					})
				},
			})
		})
	})

	core.NewItemEffect(42413, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 56186}
		manaMetrics := character.NewManaMetrics(actionID)

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 5,
				},
			},
			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					NumTicks: 12,
					Period:   time.Second * 1,
					OnAction: func(sim *core.Simulation) {
						character.AddMana(sim, 195, manaMetrics)
					},
				})
			},
		})
		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeMana,
			Spell: spell,
		})
	})

	core.NewItemEffect(47215, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 67667}
		manaMetrics := character.NewManaMetrics(actionID)

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Tears of the Vanquished Trigger",
			Callback:   core.CallbackOnCastComplete,
			SpellFlags: core.SpellFlagHelpful,
			ProcChance: 0.25,
			ICD:        time.Second * 45,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				character.AddMana(sim, 500, manaMetrics)
			},
		})
	})

	NewItemEffectWithHeroic(func(isHeroic bool) {
		name := "Heartpierce"
		itemID := int32(49982)
		numTicks := 5
		if isHeroic {
			name += " H"
			itemID = 50641
			numTicks = 6
		}

		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{ItemID: itemID}

			var resourceMetricsMana *core.ResourceMetrics
			var resourceMetricsRage *core.ResourceMetrics
			var resourceMetricsEnergy *core.ResourceMetrics
			if character.HasManaBar() {
				resourceMetricsMana = character.NewManaMetrics(actionID)
			}
			if character.HasRageBar() {
				resourceMetricsRage = character.NewRageMetrics(actionID)
			}
			if character.HasEnergyBar() {
				resourceMetricsEnergy = character.NewEnergyMetrics(actionID)
			}

			var pa *core.PendingAction
			procAura := character.GetOrRegisterAura(core.Aura{
				Label:    name,
				ActionID: actionID,
				Duration: time.Second * 2 * time.Duration(numTicks),
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					pa = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
						NumTicks: numTicks,
						Period:   time.Second * 2,
						OnAction: func(sim *core.Simulation) {
							cpb := character.GetCurrentPowerBar()
							if cpb == core.ManaBar {
								character.AddMana(sim, 120, resourceMetricsMana)
							} else if cpb == core.RageBar {
								character.AddRage(sim, 2, resourceMetricsRage)
							} else if cpb == core.EnergyBar {
								character.AddEnergy(sim, 4, resourceMetricsEnergy)
							}
						},
					})
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					pa.Cancel(sim)
				},
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:     name + " Trigger",
				Callback: core.CallbackOnSpellHitDealt,
				ProcMask: core.ProcMaskMelee,
				PPM:      1,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					// Deactivate first, to cancel old PA.
					procAura.Deactivate(sim)
					procAura.Activate(sim)
				},
			})
		})
	})

	core.NewSimpleStatOffensiveTrinketEffectWithOtherEffects(50260, stats.Stats{stats.MeleeHaste: 464, stats.SpellHaste: 464}, time.Second*20, time.Minute*2, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 50260}
		manaMetrics := character.NewManaMetrics(actionID)

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Ephemeral Snowflake",
			Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
			ICD:      time.Millisecond * 250,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				character.AddMana(sim, 11, manaMetrics)
			},
		})
	})

	core.NewItemEffect(50356, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 71586}

		var shield *core.Shield
		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolHoly,
			ProcMask:    core.ProcMaskSpellHealing,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				shield.Apply(sim, 6400)
			},
		})

		shield = core.NewShield(core.Shield{
			Spell: spell,
			Aura: character.GetOrRegisterAura(core.Aura{
				Label:    "Hardened Skin",
				ActionID: actionID,
				Duration: time.Second * 10,
			}),
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeSurvival,
			Spell: spell,
		})
	})

	NewItemEffectWithHeroic(func(isHeroic bool) {
		name := "Bauble of True Blood"
		itemID := int32(50354)
		if isHeroic {
			name += " H"
			itemID = 50726
		}

		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{ItemID: itemID}

			spell := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolHoly,
				ProcMask:    core.ProcMaskSpellHealing,
				Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				CritMultiplier:   character.DefaultHealingCritMultiplier(),

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseHealing := sim.Roll(7400, 8600)
					spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)
				},
			})

			character.AddMajorCooldown(core.MajorCooldown{
				Type:  core.CooldownTypeDPS,
				Spell: spell,
			})
		})
	})

	NewItemEffectWithHeroic(func(isHeroic bool) {
		name := "Althor's Abacus"
		itemID := int32(50359)
		minHeal := 5550.0
		maxHeal := 6450.0
		if isHeroic {
			name += " H"
			itemID = 50366
			minHeal = 6280.0
			maxHeal = 7298.0
		}

		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{ItemID: itemID}

			spell := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskSpellHealing,
				Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				CritMultiplier:   character.DefaultHealingCritMultiplier(),

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseHealing := sim.Roll(minHeal, maxHeal)
					spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)
				},
			})

			eligibleUnits := character.Env.Raid.AllUnits

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Althor's Abacus Trigger",
				Callback:   core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
				ProcChance: 0.3,
				ICD:        time.Second * 45,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					randIndex := int(math.Floor(sim.RandomFloat("Althor's Abacus") * float64(len(eligibleUnits))))
					healTarget := eligibleUnits[randIndex]
					spell.Cast(sim, healTarget)
				},
			})
		})
	})

	NewItemEffectWithHeroic(func(isHeroic bool) {
		name := "Glowing Twilight Scale"
		itemID := int32(54573)
		healingPerTick := 356.0
		if isHeroic {
			name += " H"
			itemID = 54589
			healingPerTick = 402.0
		}

		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{ItemID: itemID}

			healSpell := character.GetOrRegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolHoly,
				ProcMask:    core.ProcMaskSpellHealing,
				Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				Hot: core.DotConfig{
					Aura: core.Aura{
						Label: "Glowing Twilight Scale",
					},
					NumberOfTicks: 6,
					TickLength:    time.Second * 1,
					OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
						dot.SnapshotBaseDamage = healingPerTick
						dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
					},
					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
					},
				},
			})

			activeAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:     name + " Trigger",
				ActionID: actionID,
				Callback: core.CallbackOnHealDealt,
				Duration: time.Second * 15,
				Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
					healSpell.Hot(result.Target).Apply(sim)
				},
			})

			spell := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete,

				Cast: core.CastConfig{
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 2,
					},
				},

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					activeAura.Activate(sim)
				},
			})

			character.AddMajorCooldown(core.MajorCooldown{
				Type:  core.CooldownTypeDPS,
				Spell: spell,
			})
		})
	})

	NewItemEffectWithHeroic(func(isHeroic bool) {
		name := "Petrified Twilight Scale"
		itemID := int32(54571)
		amount := 733.0
		if isHeroic {
			name += " H"
			itemID = 54591
			amount = 828.0
		}

		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{ItemID: itemID}

			procAura := character.NewTemporaryStatsAura(name+" Proc", actionID, stats.Stats{stats.Dodge: amount}, time.Second*10)

			// Handle ICD ourselves since we use a custom check.
			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Second * 45,
			}

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:     name + " Trigger",
				Callback: core.CallbackOnSpellHitTaken,
				ProcMask: core.ProcMaskMelee,
				Harmful:  true,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					if icd.IsReady(sim) && character.CurrentHealthPercent() < 0.35 {
						icd.Use(sim)
						procAura.Activate(sim)
					}
				},
			})
		})
	})

	NewItemEffectWithHeroic(func(isHeroic bool) {
		name := "Bryntroll, the Bone Arbiter"
		itemID := int32(50415)
		minDmg := float64(2138)
		maxDmg := float64(2362)
		if isHeroic {
			name += " H"
			itemID = 50709
			minDmg = float64(2412)
			maxDmg = float64(2664)
		}

		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			mh, oh := character.GetWeaponHands(itemID)
			procMask := core.GetMeleeProcMaskForHands(mh, oh)
			ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

			procActionID := core.ActionID{ItemID: itemID}

			proc := character.RegisterSpell(core.SpellConfig{
				ActionID:    procActionID.WithTag(1),
				SpellSchool: core.SpellSchoolShadow,
				ProcMask:    core.ProcMaskEmpty,

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultSpellCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, sim.Roll(minDmg, maxDmg), spell.OutcomeMagicHitAndCrit)
				},
			})

			character.GetOrRegisterAura(core.Aura{
				Label:    name,
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() || !spell.ProcMask.Matches(procMask) {
						return
					}
					if !ppmm.Proc(sim, spell.ProcMask, name) {
						return
					}

					proc.Cast(sim, result.Target)
				},
			})
		})
	})

	NewItemEffectWithHeroic(func(isHeroic bool) {
		name := "Black Bruise"
		itemID := int32(50035)
		amount := 0.09
		spellID := int32(71875)
		if isHeroic {
			name += " H"
			itemID = 50692
			spellID = 71877
			amount = 0.10
		}

		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			var curDmg float64
			necrosisHit := character.RegisterSpell(core.SpellConfig{
				ActionID:         core.ActionID{SpellID: 51460},
				SpellSchool:      core.SpellSchoolShadow,
				ProcMask:         core.ProcMaskEmpty,
				Flags:            core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, curDmg*amount, spell.OutcomeAlwaysHit)
				},
			})

			procAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:     name + " Proc",
				ActionID: core.ActionID{SpellID: spellID},
				Callback: core.CallbackOnSpellHitDealt,
				ProcMask: core.ProcMaskMelee,
				Duration: time.Second * 10,
				Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
					curDmg = result.Damage
					necrosisHit.Cast(sim, result.Target)
				},
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       name + " Trigger",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: 0.03,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					procAura.Activate(sim)
				},
			})
		})
	})

	core.NewItemEffect(37111, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 37111}

		if !character.HasManaBar() {
			return
		}
		resourceMetricsMana := character.NewManaMetrics(actionID)

		procAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Soul Preserver",
			ActionID: actionID,
			Duration: time.Second * 15,
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !spell.ProcMask.Matches(core.ProcMaskSpellHealing) || spell.CurCast.Cost == 0 {
					return
				}
				amount := spell.CurCast.Cost
				if spell.CurCast.Cost > 800 {
					amount = 800
				}
				spell.Unit.AddMana(sim, amount, resourceMetricsMana)
				aura.Deactivate(sim)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Soul Preserver Trigger",
			Callback:   core.CallbackOnHealDealt,
			ProcMask:   core.ProcMaskSpellHealing,
			ProcChance: 0.02,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			},
		})
	})
}
