package wotlk

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func init() {
	// TODO
	// Althor's Abacus
	// Bauble of True Blood
	// Glowing Twilight Scale
	// The General's Heart
	// Ephemeral Snowflake
	// Purified Shard 2pc bonus
	// Talisman of Troll Divinity
	// Val'anyr
	// Scarab Brooch

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
		actionID := core.ActionID{ItemID: 40258}
		hots := core.NewHotArray(
			&character.Unit,
			core.Dot{
				Spell: character.GetOrRegisterSpell(core.SpellConfig{
					ActionID:    actionID,
					SpellSchool: core.SpellSchoolPhysical,
					ProcMask:    core.ProcMaskSpellHealing,
					Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

					DamageMultiplier: 1,
					ThreatMultiplier: 1,
				}),
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
			core.Aura{
				Label:    "Forethought Talisman" + strconv.Itoa(int(character.Index)),
				ActionID: actionID,
			})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Forethought Talisman",
			Callback:   core.CallbackOnHealDealt,
			Outcome:    core.OutcomeCrit,
			ProcChance: 0.2,
			ICD:        time.Second * 45,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				hot := hots[result.Target.UnitIndex]
				hot.Apply(sim)
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
				character.AddMana(sim, 900, manaMetrics, false)
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

	core.NewItemEffect(45703, func(agent core.Agent) { // Spark of Hope
		character := agent.GetCharacter()

		if !character.HasManaBar() {
			return
		}

		for _, spell := range character.Spellbook {
			if spell.ResourceType == stats.Mana && spell.BaseCost > 0 {
				defaultCastRatio := spell.DefaultCast.Cost / spell.BaseCost

				spell.BaseCost = core.MaxFloat(spell.BaseCost-44, 0)
				spell.DefaultCast.Cost = spell.BaseCost * defaultCastRatio
			}
		}
	})

	core.NewItemEffect(50259, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 50259}

		activeAura := MakeStackingAura(character, StackingProcAura{
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
					character.AddMana(sim, amount, manaMetrics, false)
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

			procAura := MakeStackingAura(character, StackingProcAura{
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
						CleanUp: func(s *core.Simulation) {
							procAura.Deactivate(sim)
						},
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
						character.AddMana(sim, 195, manaMetrics, false)
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
				character.AddMana(sim, 500, manaMetrics, false)
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
}
