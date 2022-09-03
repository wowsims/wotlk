package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/common/wotlk"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var ItemSetThunderheartRegalia = core.NewItemSet(core.ItemSet{
	Name: "Thunderheart Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// handled in moonfire.go in template construction
		},
		4: func(agent core.Agent) {
			// handled in starfire.go in template construction
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

// T9 Balance Alliance
var ItemSetMalfurionsRegalia = core.NewItemSet(core.ItemSet{
	Name: "Malfurion's Regalia",
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

// T9 Balance Horde
var ItemSetRunetotemsRegalia = core.NewItemSet(core.ItemSet{
	Name: "Runetotem's Regalia",
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
			// Implemented in spell files.
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
			druid.SwiftStarfireAura = druid.RegisterAura(core.Aura{
				Label:    "Moonkin Starfire Bonus",
				ActionID: core.ActionID{SpellID: 46832},
				Duration: time.Second * 15,
			})
			// Rest implemented in spells
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

			cca := druid.GetOrRegisterAura(core.Aura{
				Label:    "Clearcasting",
				ActionID: core.ActionID{SpellID: 16870},
				Duration: time.Second * 15,
			})

			icd := core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 15,
			}

			druid.RegisterAura(core.Aura{
				Label:    "Nightsong 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spell != druid.Rake && spell != druid.Rip && spell != druid.Lacerate {
						return
					}
					if !icd.IsReady(sim) {
						return
					}
					if sim.RandomFloat("Nightsong 4pc") < procChance {
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

var ItemSetRunetotemsBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Runetotem's Battlegear",
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
	Name: "Malfurion's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// implemented in skills
		},
		4: func(agent core.Agent) {
			// implemented in skills
		},
	},
})

func (druid *Druid) HasT9FeralSetBonus(num int32) bool {
	return druid.HasSetBonus(ItemSetRunetotemsBattlegear, num) || druid.HasSetBonus(ItemSetMalfurionsBattlegear, num)
}

func init() {

	core.NewItemEffect(30664, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		var procAura *core.Aura
		if druid.InForm(Moonkin) {
			procAura = druid.NewTemporaryStatsAura("Living Root Moonkin Proc", core.ActionID{SpellID: 37343}, stats.Stats{stats.SpellPower: 209}, time.Second*15)
		} else if druid.InForm(Bear) {
			procAura = druid.NewTemporaryStatsAura("Living Root Bear Proc", core.ActionID{SpellID: 37340}, stats.Stats{stats.Armor: 4070}, time.Second*15)
		} else if druid.InForm(Cat) {
			procAura = druid.NewTemporaryStatsAura("Living Root Cat Proc", core.ActionID{SpellID: 37341}, stats.Stats{stats.Strength: 64}, time.Second*15)
		} else {
			return
		}

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Living Root of the Wildheart",
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if druid.InForm(Moonkin) && sim.RandomFloat("Living Root of the Wildheart") < 0.03 {
					procAura.Activate(sim)
				}
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}
				if sim.RandomFloat("Living Root of the Wildheart") > 0.03 {
					return
				}

				procAura.Activate(sim)
			},
		}))
	})

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
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
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

		actionID := core.ActionID{ItemID: 32257}
		procAura := druid.NewTemporaryStatsAura("Idol of the White Stag Proc", actionID, stats.Stats{stats.AttackPower: 94}, time.Second*20)

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Idol of the White Stag",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if druid.IsMangle(spell) {
					procAura.Activate(sim)
				}
			},
		}))
	})

	core.NewItemEffect(33509, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		actionID := core.ActionID{ItemID: 33509}
		procAura := druid.NewTemporaryStatsAura("Idol of Terror Proc", actionID, stats.Stats{stats.Agility: 65}, time.Second*10)

		procChance := 0.85
		icd := core.Cooldown{
			Timer:    druid.NewTimer(),
			Duration: time.Second * 10,
		}

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Idol of Terror",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !druid.IsMangle(spell) {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if sim.RandomFloat("Idol of Terror") > procChance {
					return
				}

				icd.Use(sim)
				procAura.Activate(sim)
			},
		}))
	})

	core.NewItemEffect(33510, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		actionID := core.ActionID{ItemID: 33510}
		procAura := druid.NewTemporaryStatsAura("Idol of the Unseen Moon Proc", actionID, stats.Stats{stats.SpellPower: 140}, time.Second*10)

		druid.RegisterAura(core.Aura{
			Label:    "Idol of the Unseen Moon",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
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

		actionID := core.ActionID{ItemID: 50457}

		procAura := wotlk.MakeStackingAura(agent.GetCharacter(), wotlk.StackingProcAura{
			Aura: core.Aura{
				Label:     "Idol of the Lunar Eclipse proc",
				ActionID:  actionID,
				Duration:  time.Second * 15,
				MaxStacks: 5,
			},
			BonusPerStack: stats.Stats{stats.SpellCrit: 44},
		})

		core.MakePermanent(druid.GetOrRegisterAura(core.Aura{
			Label:    "Idol of the Lunar Eclipse",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				procAura.Activate(sim)
				procAura.AddStack(sim)
			},
		}))
	})

	core.NewItemEffect(32387, func(agent core.Agent) {
		// Idol of the Raven Goddess
		// This should maybe be an Aura, but this way it changes stats on sheet
		druid := agent.(DruidAgent).GetDruid()

		if druid.InForm(Bear | Cat) {
			druid.AddStat(stats.MeleeCrit, 40.0)
		} else if druid.InForm(Moonkin) {
			druid.AddStat(stats.SpellCrit, 40.0)
		}
	})

	core.NewItemEffect(37573, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		actionID := core.ActionID{ItemID: 37573}
		procAura := druid.NewTemporaryStatsAura("Idol of the Plainstalker Proc", actionID, stats.Stats{stats.Agility: 55}, time.Second*10)
		procChance := 0.75
		icd := core.Cooldown{
			Timer:    druid.NewTimer(),
			Duration: time.Second * 10,
		}
		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Idol of the Plainstalker",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !druid.IsMangle(spell) {
					return
				}

				if !icd.IsReady(sim) {
					return
				}
				if sim.RandomFloat("Idol of the Plainstalker") > procChance {
					return
				}

				icd.Use(sim)
				procAura.Activate(sim)
			},
		}))
	})

	core.NewItemEffect(45509, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		actionID := core.ActionID{ItemID: 45509}
		procAura := druid.NewTemporaryStatsAura("Idol of the Corruptor Proc", actionID, stats.Stats{stats.Agility: 153}, time.Second*12)

		// This proc chance might be wrong, going off of wowhead notes
		procChance := 0.85
		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Idol of the Corruptor",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !druid.IsMangle(spell) {
					return
				}

				if sim.RandomFloat("Idol of the Corruptor") > procChance {
					return
				}
				procAura.Activate(sim)
			},
		}))
	})

	core.NewItemEffect(38295, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		actionID := core.ActionID{ItemID: 38295}
		procAura := druid.NewTemporaryStatsAura("Idol of the Wastes Proc", actionID, stats.Stats{stats.Strength: 61}, time.Second*10)
		procChance := 0.75
		icd := core.Cooldown{
			Timer:    druid.NewTimer(),
			Duration: time.Second * 10,
		}
		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Idol of the Wastes",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if druid.Shred != spell && druid.IsSwipeSpell(spell) {
					return
				}

				if !icd.IsReady(sim) {
					return
				}
				if sim.RandomFloat("Idol of the Wastes") > procChance {
					return
				}

				icd.Use(sim)
				procAura.Activate(sim)
			},
		}))
	})

	core.NewItemEffect(42574, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		actionID := core.ActionID{ItemID: 42574}
		procAura := druid.NewTemporaryStatsAura("Savage Gladiator's Idol of Resolve Proc", actionID, stats.Stats{stats.AttackPower: 94}, time.Second*6)

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Savage Gladiator's Idol of Resolve",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !druid.IsMangle(spell) {
					return
				}
				procAura.Activate(sim)
			},
		}))
	})

	core.NewItemEffect(42587, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		actionID := core.ActionID{ItemID: 42587}
		procAura := druid.NewTemporaryStatsAura("Hateful Gladiator's Idol of Resolve Proc", actionID, stats.Stats{stats.AttackPower: 106}, time.Second*6)

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Hateful Gladiator's Idol of Resolve",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !druid.IsMangle(spell) {
					return
				}
				procAura.Activate(sim)
			},
		}))
	})

	core.NewItemEffect(42588, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()
		actionID := core.ActionID{ItemID: 42588}
		procAura := druid.NewTemporaryStatsAura("Deadly Gladiator's Idol of Resolve Proc", actionID, stats.Stats{stats.AttackPower: 120}, time.Second*10)

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Deadly Gladiator's Idol of Resolve",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !druid.IsMangle(spell) {
					return
				}
				procAura.Activate(sim)
			},
		}))
	})
}

func (druid *Druid) registerLasherweaveDot() {
	if !druid.SetBonuses.balance_t10_4 {
		return
	}

	dotSpell := druid.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 71023},
		SpellSchool: core.SpellSchoolNature,
	})

	druid.LasherweaveDot = core.NewDot(core.Dot{
		Spell: dotSpell,
		Aura: druid.CurrentTarget.RegisterAura(core.Aura{
			Label:    "Languish",
			ActionID: core.ActionID{SpellID: 71023},
		}),
		NumberOfTicks: 2,
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncSnapshot(druid.CurrentTarget, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return druid.GetStat(stats.SpellPower) * 0.07
				},
			},
			OutcomeApplier: druid.OutcomeFuncTick(),
		}),
	})
}
