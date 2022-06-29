package common

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func init() {
	//// Battlemasters trinkets
	//sharedBattlemasterCooldownID := core.NewCooldownID()
	//addBattlemasterEffect := func(itemID int32) {
	//	core.NewItemEffect(itemID, core.MakeTemporaryStatsOnUseCDRegistration(
	//		"BattlemasterTrinket-"+strconv.Itoa(int(itemID)),
	//		stats.Stats{stats.Health: 1750},
	//		time.Second*15,
	//		core.MajorCooldown{
	//			ActionID:         core.ActionID{ItemID: itemID},
	//			CooldownID:       sharedBattlemasterCooldownID,
	//			Cooldown:         time.Minute * 3,
	//			SharedCooldownID: core.DefensiveTrinketSharedCooldownID,
	//		},
	//	))
	//}
	//addBattlemasterEffect(33832)
	//addBattlemasterEffect(34049)
	//addBattlemasterEffect(34050)
	//addBattlemasterEffect(34162)
	//addBattlemasterEffect(34163)

	// Offensive trinkets. Keep these in order by item ID.
	NewSimpleStatOffensiveTrinketEffect(22954, stats.Stats{stats.MeleeHaste: 200}, time.Second*15, time.Minute*2)                                 // Kiss of the Spider
	NewSimpleStatOffensiveTrinketEffect(23041, stats.Stats{stats.AttackPower: 260, stats.RangedAttackPower: 260}, time.Second*20, time.Minute*2)  // Slayer's Crest
	NewSimpleStatOffensiveTrinketEffect(24128, stats.Stats{stats.AttackPower: 320, stats.RangedAttackPower: 320}, time.Second*12, time.Minute*3)  // Figurine Nightseye Panther
	NewSimpleStatOffensiveTrinketEffect(28041, stats.Stats{stats.AttackPower: 200, stats.RangedAttackPower: 200}, time.Second*15, time.Minute*2)  // Bladefists Breadth
	NewSimpleStatOffensiveTrinketEffect(28121, stats.Stats{stats.ArmorPenetration: 600}, time.Second*20, time.Minute*2)                           // Icon of Unyielding Courage
	NewSimpleStatOffensiveTrinketEffect(28288, stats.Stats{stats.MeleeHaste: 260}, time.Second*10, time.Minute*2)                                 // Abacus of Violent Odds
	NewSimpleStatOffensiveTrinketEffect(29383, stats.Stats{stats.AttackPower: 278, stats.RangedAttackPower: 278}, time.Second*20, time.Minute*2)  // Bloodlust Brooch
	NewSimpleStatOffensiveTrinketEffect(29776, stats.Stats{stats.AttackPower: 200, stats.RangedAttackPower: 200}, time.Second*20, time.Minute*2)  // Core of Arkelos
	NewSimpleStatOffensiveTrinketEffect(32658, stats.Stats{stats.Agility: 150}, time.Second*20, time.Minute*2)                                    // Badge of Tenacity
	NewSimpleStatOffensiveTrinketEffect(33831, stats.Stats{stats.AttackPower: 360, stats.RangedAttackPower: 360}, time.Second*20, time.Minute*2)  // Berserkers Call
	NewSimpleStatOffensiveTrinketEffect(35702, stats.Stats{stats.AttackPower: 320, stats.RangedAttackPower: 320}, time.Second*15, time.Second*90) // Figurine Shadowsong Panther
	NewSimpleStatOffensiveTrinketEffect(38287, stats.Stats{stats.AttackPower: 278, stats.RangedAttackPower: 278}, time.Second*20, time.Minute*2)  // Empty Direbrew Mug

	// Defensive trinkets. Keep these in order by item ID.
	NewSimpleStatDefensiveTrinketEffect(27891, stats.Stats{stats.Armor: 1280}, time.Second*20, time.Minute*2)                                                          // Adamantine Figurine
	NewSimpleStatDefensiveTrinketEffect(28528, stats.Stats{stats.Dodge: 300}, time.Second*10, time.Minute*2)                                                           // Moroes Lucky Pocket Watch
	NewSimpleStatDefensiveTrinketEffect(29387, stats.Stats{stats.BlockValue: 200}, time.Second*20, time.Minute*2)                                                      // Gnomeregan Auto-Blocker 600
	NewSimpleStatDefensiveTrinketEffect(30300, stats.Stats{stats.Block: 125}, time.Second*15, time.Second*90)                                                          // Dabiris Enigma
	NewSimpleStatDefensiveTrinketEffect(30629, stats.Stats{stats.Defense: 165, stats.AttackPower: -330, stats.RangedAttackPower: -330}, time.Second*15, time.Minute*3) // Scarab of Displacement
	NewSimpleStatDefensiveTrinketEffect(32501, stats.Stats{stats.Health: 1750}, time.Second*20, time.Minute*3)                                                         // Shadowmoon Insignia
	NewSimpleStatDefensiveTrinketEffect(32534, stats.Stats{stats.Health: 1250}, time.Second*15, time.Minute*5)                                                         // Brooch of the Immortal King
	NewSimpleStatDefensiveTrinketEffect(33830, stats.Stats{stats.Armor: 2500}, time.Second*20, time.Minute*2)                                                          // Ancient Aqir Artifact
	NewSimpleStatDefensiveTrinketEffect(38289, stats.Stats{stats.BlockValue: 200}, time.Second*20, time.Minute*2)                                                      // Coren's Lucky Coin

	// Proc effects. Keep these in order by item ID.

	core.NewItemEffect(11815, func(agent core.Agent) {
		character := agent.GetCharacter()
		if !character.AutoAttacks.IsEnabled() {
			return
		}

		var handOfJusticeSpell *core.Spell
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 2,
		}
		procChance := 0.013333

		character.RegisterAura(core.Aura{
			Label:    "Hand of Justice",
			Duration: core.NeverExpires,
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				handOfJusticeSpell = character.GetOrRegisterSpell(core.SpellConfig{
					ActionID:     core.ActionID{ItemID: 11815},
					SpellSchool:  core.SpellSchoolPhysical,
					Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,
					ApplyEffects: core.ApplyEffectFuncDirectDamage(character.AutoAttacks.MHEffect),
				})
			},
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// https://tbc.wowhead.com/spell=15600/hand-of-justice, proc mask = 20.
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if !icd.IsReady(sim) {
					return
				}

				if sim.RandomFloat("HandOfJustice") > procChance {
					return
				}
				icd.Use(sim)

				aura.Unit.AutoAttacks.MaybeReplaceMHSwing(sim, handOfJusticeSpell).Cast(sim, spellEffect.Target)
			},
		})
	})

	core.NewItemEffect(32654, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.BonusDamage += 7
		core.RegisterTemporaryStatsOnUseCD(
			character,
			"Crystalforged Trinket",
			stats.Stats{stats.AttackPower: 216, stats.RangedAttackPower: 216},
			time.Second*10,
			core.SpellConfig{
				ActionID: core.ActionID{ItemID: 32654},
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute,
					},
					SharedCD: core.Cooldown{
						Timer:    character.GetOffensiveTrinketCD(),
						Duration: time.Second * 10,
					},
				},
			},
		)
	})

	core.NewItemEffect(21670, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.RegisterAura(core.Aura{
			Label:     "Badge of the Swarmguard Proc",
			ActionID:  core.ActionID{SpellID: 26481},
			Duration:  core.NeverExpires,
			MaxStacks: 6,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				character.AddStatDynamic(sim, stats.ArmorPenetration, 200*float64(newStacks-oldStacks))
			},
		})

		actionID := core.ActionID{ItemID: 21670}
		ppmm := character.AutoAttacks.NewPPMManager(10.0, core.ProcMaskMeleeOrRanged)
		activeAura := character.RegisterAura(core.Aura{
			Label:    "Badge of the Swarmguard",
			ActionID: actionID,
			Duration: time.Second * 30,
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				procAura.Deactivate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}
				if !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}

				if !ppmm.ProcWithWeaponSpecials(sim, spellEffect.ProcMask, "Badge of the Swarmguard") {
					return
				}

				procAura.Activate(sim)
				procAura.AddStack(sim)
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				activeAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS | core.CooldownTypeUsableShapeShifted,
		})
	})

	core.NewItemEffect(23206, func(agent core.Agent) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDemon || character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
			character.PseudoStats.MobTypeAttackPower += 150
		}
	})

	core.NewItemEffect(28034, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Rage of the Unraveller", core.ActionID{ItemID: 28034}, stats.Stats{stats.AttackPower: 300, stats.RangedAttackPower: 300}, time.Second*10)
		const procChance = 0.1

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 50,
		}

		character.RegisterAura(core.Aura{
			Label:    "Hourglass of the Unraveller",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
					return
				}
				if !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if sim.RandomFloat("Hourglass of the Unraveller") > procChance {
					return
				}

				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(28579, func(agent core.Agent) {
		character := agent.GetCharacter()

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 28579},
			SpellSchool: core.SpellSchoolNature,
			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskEmpty,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				BaseDamage:     core.BaseDamageConfigRoll(222, 332),
				OutcomeApplier: character.OutcomeFuncMagicHitAndCrit(character.DefaultSpellCritMultiplier()),
			}),
		})

		ppmm := character.AutoAttacks.NewPPMManager(1.0, core.ProcMaskMeleeOrRanged)

		character.RegisterAura(core.Aura{
			Label:    "Romulos Poison Vial",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// mask 340
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !ppmm.ProcWithWeaponSpecials(sim, spellEffect.ProcMask, "RomulosPoisonVial") {
					return
				}

				procSpell.Cast(sim, spellEffect.Target)
			},
		})
	})

	core.NewItemEffect(28830, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Dragonspine Trophy Proc", core.ActionID{ItemID: 28830}, stats.Stats{stats.MeleeHaste: 325}, time.Second*10)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 20,
		}
		ppmm := character.AutoAttacks.NewPPMManager(1.0, core.ProcMaskMeleeOrRanged)

		character.RegisterAura(core.Aura{
			Label:    "Dragonspine Trophy",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// mask: 340
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if !ppmm.ProcWithWeaponSpecials(sim, spellEffect.ProcMask, "dragonspine") {
					return
				}
				icd.Use(sim)

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(30627, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Tsunami Talisman Proc", core.ActionID{ItemID: 30627}, stats.Stats{stats.AttackPower: 340, stats.RangedAttackPower: 340}, time.Second*10)
		const procChance = 0.1

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 45,
		}

		character.RegisterAura(core.Aura{
			Label:    "Tsunami Talisman",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
					return
				}
				if !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if sim.RandomFloat("Tsunami Talisman") > procChance {
					return
				}

				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(31857, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.RegisterAura(core.Aura{
			Label:     "DMC Wrath Proc",
			ActionID:  core.ActionID{ItemID: 31857},
			Duration:  time.Second * 10,
			MaxStacks: 1000,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				character.AddStatDynamic(sim, stats.MeleeCrit, 17*float64(newStacks-oldStacks))
				character.AddStatDynamic(sim, stats.SpellCrit, 17*float64(newStacks-oldStacks))
			},
		})

		character.RegisterAura(core.Aura{
			Label:    "DMC Wrath",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// mask 340
				if !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}

				if spellEffect.Outcome.Matches(core.OutcomeCrit) {
					procAura.Deactivate(sim)
				} else {
					procAura.Activate(sim)
					procAura.AddStack(sim)
				}
			},
		})
	})

	core.NewItemEffect(31858, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 31858}

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolHoly,
			ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
				ProcMask:         core.ProcMaskEmpty,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				BaseDamage:     core.BaseDamageConfigRoll(95, 115),
				OutcomeApplier: character.OutcomeFuncCritFixedChance(0.03, character.DefaultMeleeCritMultiplier()),
			}),
		})

		// Normal proc chance.
		procChance := 0.1

		// JoL and JoW procs can activate this effect. JoL and JoW both have a 50% chance
		// to proc so just add them.
		procChanceOnHitDealt := 0.0
		if character.CurrentTarget.HasAura(core.JudgementOfLightAuraLabel) {
			procChanceOnHitDealt += procChance / 2
		}
		if character.CurrentTarget.HasAura(core.JudgementOfWisdomAuraLabel) {
			procChanceOnHitDealt += procChance / 2
		}

		// Can also proc when the player's melee attacks trigger a Seal of Light proc.
		var onSpellHitDealt core.OnSpellHit
		if procChanceOnHitDealt > 0 {
			onSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() && spellEffect.ProcMask.Matches(core.ProcMaskMelee) && sim.RandomFloat("DMC Vengeance") < procChanceOnHitDealt {
					procSpell.Cast(sim, spellEffect.Target)
				}
			}
		}

		character.RegisterAura(core.Aura{
			Label:    "DMC Vengeance",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() && sim.RandomFloat("DMC Vengeance") < procChance {
					procSpell.Cast(sim, spell.Unit)
				}
			},
			OnSpellHitDealt: onSpellHitDealt,
		})
	})

	core.NewItemEffect(32505, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Madness of the Betrayer Proc", core.ActionID{ItemID: 32505}, stats.Stats{stats.ArmorPenetration: 300}, time.Second*10)

		ppmm := character.AutoAttacks.NewPPMManager(1.0, core.ProcMaskMeleeOrRanged)

		character.RegisterAura(core.Aura{
			Label:    "Madness of the Betrayer",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// mask 340
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !ppmm.ProcWithWeaponSpecials(sim, spellEffect.ProcMask, "Madness of the Betrayer") {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(34427, func(agent core.Agent) {
		character := agent.GetCharacter()

		var bonusPerStack stats.Stats
		procAura := character.RegisterAura(core.Aura{
			Label:     "Blackened Naaru Sliver Proc",
			ActionID:  core.ActionID{ItemID: 34427},
			Duration:  time.Second * 20,
			MaxStacks: 10,
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				bonusPerStack = character.ApplyStatDependencies(stats.Stats{stats.AttackPower: 44, stats.RangedAttackPower: 44})
			},
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				character.AddStatsDynamic(sim, bonusPerStack.Multiply(float64(newStacks-oldStacks)))
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() && spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					aura.AddStack(sim)
				}
			},
		})

		const procChance = 0.1

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 45,
		}

		character.RegisterAura(core.Aura{
			Label:    "Blackened Naaru Sliver",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// mask 340
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if sim.RandomFloat("Blackened Naaru Sliver") > procChance {
					return
				}

				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(34472, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Shard of Contempt Proc", core.ActionID{ItemID: 34472}, stats.Stats{stats.AttackPower: 230, stats.RangedAttackPower: 230}, time.Second*20)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 45,
		}
		const procChance = 0.1

		character.RegisterAura(core.Aura{
			Label:    "Shard of Contempt",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if sim.RandomFloat("Shard of Contempt") > procChance {
					return
				}

				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

}
