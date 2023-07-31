package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var ItemSetThunderheartRegalia = core.NewItemSet(core.ItemSet{
	Name: "Thunderheart Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in moonfire.go
		},
		4: func(agent core.Agent) {
			// Implemented in starfire.go
		},
	},
})

var ItemSetThunderheartHarness = core.NewItemSet(core.ItemSet{
	Name: "Thunderheart Harness",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in mangle.go.
		},
		4: func(agent core.Agent) {
			// Implemented in swipe.go.
		},
	},
})

// T7 Balance
var ItemSetDreamwalkerGarb = core.NewItemSet(core.ItemSet{
	Name: "Dreamwalker Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			//Your Insect Swarm deals an additional 10% damage.
			// Implemented in insect_swarm.go.
		},
		4: func(agent core.Agent) {
			// Your Wrath and Starfire spells gain an additional 5% critical strike chance.
			// Implemented in spell files.
		},
	},
})

// T8 Balance
var ItemSetNightsongGarb = core.NewItemSet(core.ItemSet{
	Name: "Nightsong Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the bonus granted by Eclipse for Starfire and Wrath by 7%.
			// Implemented in spell files.
		},
		4: func(agent core.Agent) {
			// Each time your Insect Swarm deals damage, it has a chance to make your next Starfire cast within 10 sec instant.
			// Implemented in spell files.
		},
	},
})

// T9 Balance
var ItemSetMalfurionsRegalia = core.NewItemSet(core.ItemSet{
	Name:            "Malfurion's Regalia",
	AlternativeName: "Runetotem's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your Moonfire ability now has a chance for its periodic damage to be critical strikes.
			// Implemented in moonfire.go
		},
		4: func(agent core.Agent) {
			// Increases the damage done by your Starfire and Wrath spells by 4%.
			// Implemented in starfire.go and wrath.go
		},
	},
})

// T10 Balance
var ItemSetLasherweaveRegalia = core.NewItemSet(core.ItemSet{
	Name: "Lasherweave Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// When you gain Clearcasting from your Omen of Clarity talent, you deal 15% additional Nature and Arcane damage for 6 sec.
			// Implemented in talents.go
		},
		4: func(agent core.Agent) {
			// Your critical strikes from Starfire and Wrath cause the target to languish for an additional 7% of your spell's damage over 4 sec.
			druid := agent.(DruidAgent).GetDruid()

			druid.Languish = druid.RegisterSpell(core.SpellConfig{
				ActionID:         core.ActionID{SpellID: 71023},
				SpellSchool:      core.SpellSchoolNature,
				ProcMask:         core.ProcMaskEmpty,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				Dot: core.DotConfig{
					Aura: core.Aura{
						Label: "Languish",
					},
					NumberOfTicks: 2,
					TickLength:    time.Second * 2,

					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.Dot(target).ApplyOrReset(sim)
					spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
				},
			})

			druid.RegisterAura(core.Aura{
				Label:    "Languish proc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell != druid.Starfire && spell != druid.Wrath {
						return
					}
					if result.DidCrit() {
						dot := druid.Languish.Dot(result.Target)

						newDamage := result.Damage * 0.07
						outstandingDamage := core.TernaryFloat64(dot.IsActive(), dot.SnapshotBaseDamage*float64(dot.NumberOfTicks-dot.TickCount), 0)

						dot.SnapshotAttackerMultiplier = 1
						dot.SnapshotBaseDamage = (outstandingDamage + newDamage) / 2.0
						druid.Languish.Cast(sim, result.Target)
					}
				},
			})
		},
	},
})

var ItemSetGladiatorsWildhide = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Wildhide",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.AddStat(stats.SpellPower, 29)
			druid.AddStat(stats.Resilience, 100)
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.AddStat(stats.SpellPower, 88)

			percentReduction := float64(time.Millisecond*1500) / float64(druid.starfireCastTime())
			swiftStarfireAura := druid.RegisterAura(core.Aura{
				Label:    "Swift Starfire",
				ActionID: core.ActionID{SpellID: 46832},
				Duration: time.Second * 15,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					druid.Starfire.CastTimeMultiplier -= percentReduction
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					druid.Starfire.CastTimeMultiplier += percentReduction
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell == druid.Starfire {
						aura.Deactivate(sim)
					}
				},
			})

			druid.RegisterAura(core.Aura{
				Label:    "Swift Starfire trigger",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell == druid.Wrath && sim.RandomFloat("Swift Starfire proc") > 0.85 {
						swiftStarfireAura.Activate(sim)
					}
				},
			})
		},
	},
})

var ItemSetGladiatorsSanctuary = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Sanctuary",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.Resilience, 100)
			agent.GetCharacter().AddStat(stats.AttackPower, 50)
		},
		4: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.AttackPower, 150)
		},
	},
})

var ItemSetNightsongBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Nightsong Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			// The periodic damage dealt by your Rake, Rip, and Lacerate abilities
			// has a chance to cause you to enter a Clearcasting state.
			// (Proc chance: 2%, 15s cooldown)

			procChance := 0.02

			cca := druid.GetAura("Clearcasting")

			icd := core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 15,
			}

			druid.RegisterAura(core.Aura{
				Label:    "Nightsong 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
					cca = druid.GetAura("Clearcasting")
					cca.Icd = &icd
					if cca == nil {
						panic("no valid clearcasting aura")
					}
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell != druid.Rake && spell != druid.Rip && spell != druid.Lacerate {
						return
					}
					if !icd.IsReady(sim) {
						return
					}
					if sim.RandomFloat("Nightsong 2pc") < procChance {
						icd.Use(sim)
						cca.Activate(sim)
					}
				},
			})
		},
		4: func(agent core.Agent) {
			// Implemented in savage roar
		},
	},
})

var ItemSetLasherweaveBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Lasherweave Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// implemented in skills
		},
		4: func(agent core.Agent) {
			// implemented in skills
		},
	},
})

var ItemSetDreamwalkerBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Dreamwalker Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// implemented in skills
		},
		4: func(agent core.Agent) {
			// implemented in skills
		},
	},
})

var ItemSetMalfurionsBattlegear = core.NewItemSet(core.ItemSet{
	Name:            "Malfurion's Battlegear",
	AlternativeName: "Runetotem's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// implemented in skills
		},
		4: func(agent core.Agent) {
			// implemented in skills
		},
	},
})

func init() {

	core.NewItemEffect(32486, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		// Not in the game yet so cant test; this logic assumes that:
		// - does not affect the starfire which procs it
		// - can proc off of any completed cast, not just hits
		actionID := core.ActionID{ItemID: 32486}

		var procAura *core.Aura
		if druid.InForm(Moonkin) {
			procAura = druid.NewTemporaryStatsAura("Ashtongue Talisman Proc", actionID, stats.Stats{stats.SpellPower: 150}, time.Second*8)
		} else if druid.InForm(Bear | Cat) {
			procAura = druid.NewTemporaryStatsAura("Ashtongue Talisman Proc", actionID, stats.Stats{stats.Strength: 140}, time.Second*8)
		} else {
			return
		}

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Ashtongue Talisman",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}
				if spell == druid.Starfire {
					if sim.RandomFloat("Ashtongue Talisman") < 0.25 {
						procAura.Activate(sim)
					}
				} else if druid.IsMangle(spell) {
					if sim.RandomFloat("Ashtongue Talisman") < 0.4 {
						procAura.Activate(sim)
					}
				}
			},
		}))
	})

	core.NewItemEffect(32257, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		procAura := druid.NewTemporaryStatsAura("Idol of the White Stag Proc", core.ActionID{ItemID: 32257}, stats.Stats{stats.AttackPower: 94}, time.Second*20)

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Idol of the White Stag",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if druid.IsMangle(spell) {
					procAura.Activate(sim)
				}
			},
		}))
	})

	core.NewItemEffect(33510, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		procAura := druid.NewTemporaryStatsAura("Idol of the Unseen Moon Proc", core.ActionID{ItemID: 33510}, stats.Stats{stats.SpellPower: 140}, time.Second*10)

		druid.RegisterAura(core.Aura{
			Label:    "Idol of the Unseen Moon",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if spell == druid.Moonfire {
					if sim.RandomFloat("Idol of the Unseen Moon") > 0.5 {
						return
					}
					procAura.Activate(sim)
				}
			},
		})
	})

	// This Idol is badly listed on Wowhead, not accessible from UI
	core.NewItemEffect(50457, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		procAura := core.MakeStackingAura(agent.GetCharacter(), core.StackingStatAura{
			Aura: core.Aura{
				Label:     "Idol of the Lunar Eclipse proc",
				ActionID:  core.ActionID{ItemID: 50457},
				Duration:  time.Second * 15,
				MaxStacks: 5,
			},
			BonusPerStack: stats.Stats{stats.SpellCrit: 44},
		})

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Idol of the Lunar Eclipse",
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procAura.Activate(sim)
				procAura.AddStack(sim)
			},
		}))
	})

	core.NewItemEffect(32387, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label:      "Idol of the Raven Goddess",
			BuildPhase: core.CharacterBuildPhaseGear,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				// For now this assume we'll never leave main form
				if druid.StartingForm.Matches(Bear | Cat) {
					druid.AddStatDynamic(sim, stats.MeleeCrit, 40.0)
				} else if druid.StartingForm.Matches(Moonkin) {
					druid.AddStatDynamic(sim, stats.SpellCrit, 40.0)
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if druid.StartingForm.Matches(Bear | Cat) {
					druid.AddStatDynamic(sim, stats.MeleeCrit, -40.0)
				} else if druid.StartingForm.Matches(Moonkin) {
					druid.AddStatDynamic(sim, stats.SpellCrit, -40.0)
				}
			},
		}))
	})

	core.NewItemEffect(45509, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		procAura := druid.NewTemporaryStatsAura("Idol of the Corruptor Proc", core.ActionID{ItemID: 45509}, stats.Stats{stats.Agility: 162}, time.Second*12)

		// Proc chances based on testing by druid discord
		procChanceBear := 0.50
		procChanceCat := 1.0
		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Idol of the Corruptor",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				procChance := 0.0
				if spell == druid.MangleBear {
					procChance = procChanceBear
				} else if spell == druid.MangleCat {
					procChance = procChanceCat
				} else {
					return
				}

				if sim.RandomFloat("Idol of the Corruptor") > procChance {
					return
				}
				procAura.Activate(sim)
			},
		}))
	})

	core.NewItemEffect(47668, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		actionID := core.ActionID{ItemID: 47668}
		bearAura := druid.NewTemporaryStatsAura("Idol of Mutilation Bear Proc", actionID, stats.Stats{stats.Dodge: 200.0}, time.Second*9)
		catAura := druid.NewTemporaryStatsAura("Idol of Mutilation Cat Proc", actionID, stats.Stats{stats.Agility: 200.0}, time.Second*16)

		// Based off of wowhead tooltip
		icd := core.Cooldown{
			Timer:    druid.NewTimer(),
			Duration: time.Second * 8,
		}
		bearAura.Icd = &icd
		catAura.Icd = &icd
		procChance := 0.7

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Idol of Mutilation",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if !icd.IsReady(sim) {
					return
				}

				if sim.RandomFloat("Idol of Mutilation") > procChance {
					return
				}
				if spell == druid.SwipeBear || spell == druid.Lacerate {
					icd.Use(sim)
					bearAura.Activate(sim)
				}
				if spell == druid.MangleCat || spell == druid.Shred {
					icd.Use(sim)
					catAura.Activate(sim)
				}
			},
		}))
	})

	core.NewItemEffect(50456, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		procAura := core.MakeStackingAura(agent.GetCharacter(), core.StackingStatAura{
			Aura: core.Aura{
				Label:     "Idol of the Crying Moon Proc",
				ActionID:  core.ActionID{ItemID: 50456},
				Duration:  time.Second * 15,
				MaxStacks: 5,
			},
			BonusPerStack: stats.Stats{stats.Agility: 44},
		})

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Idol of the Crying Moon",
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell != druid.Rake && spell != druid.Lacerate {
					return
				}
				procAura.Activate(sim)
				procAura.AddStack(sim)
			},
		}))
	})

	core.NewItemEffect(33947, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		procAura := druid.NewTemporaryStatsAura("Vengeful Gladiator's Idol of Resolve Proc", core.ActionID{ItemID: 33947}, stats.Stats{stats.Resilience: 34}, time.Second*6)

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Vengeful Gladiator's Idol of Resolve",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !druid.IsMangle(spell) {
					return
				}
				procAura.Activate(sim)
			},
		}))
	})

	core.NewItemEffect(35019, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		procAura := druid.NewTemporaryStatsAura("Brutal Gladiator's Idol of Resolve Proc", core.ActionID{ItemID: 35019}, stats.Stats{stats.Resilience: 39}, time.Second*6)

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Brutal Gladiator's Idol of Resolve",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !druid.IsMangle(spell) {
					return
				}
				procAura.Activate(sim)
			},
		}))
	})

	core.NewItemEffect(42574, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		procAura := druid.NewTemporaryStatsAura("Savage Gladiator's Idol of Resolve Proc", core.ActionID{ItemID: 42574}, stats.Stats{stats.AttackPower: 94}, time.Second*6)

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Savage Gladiator's Idol of Resolve",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !druid.IsMangle(spell) {
					return
				}
				procAura.Activate(sim)
			},
		}))
	})

	core.NewItemEffect(42587, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		procAura := druid.NewTemporaryStatsAura("Hateful Gladiator's Idol of Resolve Proc", core.ActionID{ItemID: 42587}, stats.Stats{stats.AttackPower: 106}, time.Second*6)

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Hateful Gladiator's Idol of Resolve",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !druid.IsMangle(spell) {
					return
				}
				procAura.Activate(sim)
			},
		}))
	})

	core.NewItemEffect(42588, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		procAura := druid.NewTemporaryStatsAura("Deadly Gladiator's Idol of Resolve Proc", core.ActionID{ItemID: 42588}, stats.Stats{stats.AttackPower: 120}, time.Second*10)

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Deadly Gladiator's Idol of Resolve",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !druid.IsMangle(spell) {
					return
				}
				procAura.Activate(sim)
			},
		}))
	})

	core.NewItemEffect(42589, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		procAura := druid.NewTemporaryStatsAura("Furious Gladiator's Idol of Resolve Proc", core.ActionID{ItemID: 42589}, stats.Stats{stats.AttackPower: 152}, time.Second*10)

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Furious Gladiator's Idol of Resolve",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !druid.IsMangle(spell) {
					return
				}
				procAura.Activate(sim)
			},
		}))
	})

	core.NewItemEffect(47670, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		procAura := druid.NewTemporaryStatsAura("Lunar Fire", core.ActionID{SpellID: 67360}, stats.Stats{stats.MeleeCrit: 200, stats.SpellCrit: 200}, time.Second*12)
		icd := core.Cooldown{
			Timer:    druid.NewTimer(),
			Duration: time.Second * 6,
		}
		procAura.Icd = &icd

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Idol of Lunar Fury",
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == druid.Moonfire && icd.IsReady(sim) && sim.RandomFloat("lunar fire") < 0.7 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		}))
	})

	core.NewItemEffect(42591, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		procAura := druid.NewTemporaryStatsAura("Relentless Gladiator's Idol of Resolve Proc", core.ActionID{ItemID: 42591}, stats.Stats{stats.AttackPower: 172}, time.Second*10)

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Relentless Gladiator's Idol of Resolve",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !druid.IsMangle(spell) {
					return
				}
				procAura.Activate(sim)
			},
		}))
	})

	core.NewItemEffect(51429, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		procAura := druid.NewTemporaryStatsAura("Wrathful Gladiator's Idol of Resolve Proc", core.ActionID{ItemID: 51429}, stats.Stats{stats.AttackPower: 204}, time.Second*10)

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Wrathful Gladiator's Idol of Resolve",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !druid.IsMangle(spell) {
					return
				}
				procAura.Activate(sim)
			},
		}))
	})
}
