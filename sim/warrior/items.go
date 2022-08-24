package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

/////////////////////////////////////////////////////////////////
// TBC Item set
/////////////////////////////////////////////////////////////////

var ItemSetBoldArmor = core.NewItemSet(core.ItemSet{
	Name: "Bold Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// All of your shout abilities cost 2 less rage.
			// Handled in demoralizing_shout.go.
		},
		4: func(agent core.Agent) {
			// Your Charge ability generates an additional 5 rage.
		},
	},
})

var ItemSetWarbringerArmor = core.NewItemSet(core.ItemSet{
	Name: "Warbringer Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// You have a chance each time you parry to gain Blade Turning, absorbing 200 damage for 15 sec.
		},
		4: func(agent core.Agent) {
			// Your Revenge ability causes your next damaging ability to do 10% more damage.
			warrior := agent.(WarriorAgent).GetWarrior()

			// TODO: This needs to apply only to specific abilities, not any source of damage.
			procAura := warrior.RegisterAura(core.Aura{
				Label:    "Warbringer 4pc Proc",
				ActionID: core.ActionID{SpellID: 37516},
				Duration: core.NeverExpires,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.1
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spellEffect.Damage > 0 {
						aura.Deactivate(sim)
					}
				},
			})

			warrior.RegisterAura(core.Aura{
				Label:    "Warbringer 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spell == warrior.Revenge {
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})

var ItemSetWarbringerBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Warbringer Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your whirlwind ability costs 5 less rage.
			// Handled in whirlwind.go.
		},
		4: func(agent core.Agent) {
			// You gain an additional 2 rage each time one of your attacks is parried or dodged.
			warrior := agent.(WarriorAgent).GetWarrior()
			rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 37519})

			warrior.RegisterAura(core.Aura{
				Label:    "Warbringer 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spellEffect.Outcome.Matches(core.OutcomeDodge | core.OutcomeParry) {
						warrior.AddRage(sim, 2, rageMetrics)
					}
				},
			})
		},
	},
})

var ItemSetDestroyerArmor = core.NewItemSet(core.ItemSet{
	Name: "Destroyer Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Each time you use your Shield Block ability, you gain 100 block value against a single attack in the next 6 sec.
			warrior := agent.(WarriorAgent).GetWarrior()
			procAura := warrior.RegisterAura(core.Aura{
				Label:     "Destroyer 2pc Proc",
				ActionID:  core.ActionID{SpellID: 37522},
				Duration:  time.Second * 6,
				MaxStacks: 1,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warrior.AddStatDynamic(sim, stats.BlockValue, 100)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warrior.AddStatDynamic(sim, stats.BlockValue, -100)
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spellEffect.Outcome.Matches(core.OutcomeBlock) {
						aura.RemoveStack(sim)
					}
				},
			})

			warrior.RegisterAura(core.Aura{
				Label:    "Destroyer 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell == warrior.ShieldBlock {
						procAura.Activate(sim)
						procAura.SetStacks(sim, 1)
					}
				},
			})
		},
		4: func(agent core.Agent) {
			// You have a chance each time you are hit to gain 200 haste rating for 10 sec.
			warrior := agent.(WarriorAgent).GetWarrior()
			procAura := warrior.NewTemporaryStatsAura("Destroyer 4pc Proc", core.ActionID{SpellID: 37525}, stats.Stats{stats.MeleeHaste: 200}, time.Second*10)

			warrior.RegisterAura(core.Aura{
				Label:    "Destroyer 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spellEffect.Landed() && sim.RandomFloat("Destroyer 4pc") < 0.07 {
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})

var ItemSetDestroyerBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Destroyer Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your Overpower ability now grants you 100 attack power for 5 sec.
			warrior := agent.(WarriorAgent).GetWarrior()
			procAura := warrior.NewTemporaryStatsAura("Destroyer 2pc Proc", core.ActionID{SpellID: 37528}, stats.Stats{stats.AttackPower: 100}, time.Second*5)

			warrior.RegisterAura(core.Aura{
				Label:    "Destroyer 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell == warrior.Overpower {
						procAura.Activate(sim)
					}
				},
			})
		},
		4: func(agent core.Agent) {
			// Your Bloodthirst and Mortal Strike abilities cost 5 less rage.
			// Handled in bloodthirst.go and mortal_strike.go.
		},
	},
})

var ItemSetOnslaughtArmor = core.NewItemSet(core.ItemSet{
	Name: "Onslaught Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the health bonus from your Commanding Shout ability by 170.
		},
		4: func(agent core.Agent) {
			// Increases the damage of your Shield Slam ability by 10%.
			// Handled in shield_slam.go.
		},
	},
})

var ItemSetOnslaughtBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Onslaught Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Reduces the rage cost of your Execute ability by 3.
		},
		4: func(agent core.Agent) {
			// Increases the damage of your Mortal Strike and Bloodthirst abilities by 5%.
			// Handled in bloodthirst.go and mortal_strike.go.
		},
	},
})

/////////////////////////////////////////////////////////////////
// Wrath Item set
/////////////////////////////////////////////////////////////////

var ItemSetDreadnaughtBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Dreadnaught Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the damage of your Slam by 10%.
		},
		4: func(agent core.Agent) {
			// Your Bleed periodic effects have a chance to make your next ability cost 5 less rage.
			warrior := agent.(WarriorAgent).GetWarrior()
			rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 61571})

			procAura := warrior.RegisterAura(core.Aura{
				Label:    "Dreadnaught Battlegear 2pc Proc",
				ActionID: core.ActionID{SpellID: 61571},
				Duration: time.Second * 30,
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					warrior.AddRage(sim, 5, rageMetrics)
					aura.Deactivate(sim)
				},
			})

			warrior.RegisterAura(core.Aura{
				Label:    "Dreadnaught Battlegear 2pc",
				Duration: core.NeverExpires,
				ActionID: core.ActionID{SpellID: 60176},
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spellEffect.Landed() && sim.RandomFloat("Dreadnaught Battlegear 2pc") < 0.1 {
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})

var ItemSetSiegebreakerBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Siegebreaker Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Heroic Strike and Slam critical strikes have a chance to grant you 150 haste rating for 5 sec.
			warrior := agent.(WarriorAgent).GetWarrior()
			procAura := warrior.RegisterAura(core.Aura{
				Label:    "Siegebreaker Battlegear 2pc Proc",
				ActionID: core.ActionID{SpellID: 64937},
				Duration: time.Second * 5,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warrior.AddStatDynamic(sim, stats.MeleeHaste, 150)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warrior.AddStatDynamic(sim, stats.MeleeHaste, -150)
				},
			})

			warrior.RegisterAura(core.Aura{
				Label:    "Siegebreaker Battlegear 2pc",
				Duration: core.NeverExpires,
				ActionID: core.ActionID{SpellID: 64938},
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spell.ActionID.SpellID != 47450 && spell != warrior.Slam {
						return
					}
					if spellEffect.Landed() && sim.RandomFloat("Siegebreaker Battlegear 2pc") < 0.4 {
						procAura.Activate(sim)
					}
				},
			})
		},
		4: func(agent core.Agent) {
			// Increases the critical strike chance of Mortal Strike and Bloodthirst by 10%.
			// Handled in bloodthirst.go and mortal_strike.go.
		},
	},
})

var ItemSetWrynnsBattlegear = core.NewItemSet(core.ItemSet{
	Name:            "Wrynn's Battlegear",
	AlternativeName: "Hellscream's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Berserker Stance grants an additional 2% critical strike chance, and Battle Stance grants an additional 6% armor penetration.
			// Handled in stances.go.
		},
		4: func(agent core.Agent) {
			// Increases the critical strike chance of your Slam and Heroic Strike abilities by 5%.
			// Handled in slam.go and heroic_strike_cleave.go.
		},
	},
})

var ItemSetYmirjarLordsBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Ymirjar Lord's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// When your Deep Wounds ability deals damage you have a 3% chance to gain 16% attack power for 10 sec.
			warrior := agent.(WarriorAgent).GetWarrior()
			var bonusAP float64
			procAura := warrior.RegisterAura(core.Aura{
				Label:    "Ymirjar Lord's Battlegear 2pc Proc",
				ActionID: core.ActionID{SpellID: 70855},
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					bonusAP = float64(warrior.GetStat(stats.AttackPower)) * 0.16
					aura.Unit.AddStatDynamic(sim, stats.AttackPower, bonusAP)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.AddStatDynamic(sim, stats.AttackPower, -bonusAP)
				},
			})

			warrior.RegisterAura(core.Aura{
				Label:    "Ymirjar Lord's Battlegear 2pc",
				Duration: core.NeverExpires,
				ActionID: core.ActionID{SpellID: 70854},
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spell != warrior.DeepWounds {
						return
					}
					if spellEffect.Landed() && sim.RandomFloat("Ymirjar Lord's Battlegear 2pc") < 0.03 {
						procAura.Activate(sim)
					}
				},
			})
		},
		4: func(agent core.Agent) {
			// You have a 20% chance for your Bloodsurge and Sudden Death talents to grant 2 charges of their effect instead of 1,
			// reduce the global cooldown on Execute or Slam by 0.5 sec, and for the duration of the effect to be increased by 100%.
			// Handled in talents.go.
		},
	},
})

func init() {

	core.NewItemEffect(32485, func(agent core.Agent) {
		warrior := agent.(WarriorAgent).GetWarrior()
		procAura := warrior.NewTemporaryStatsAura("Ashtongue Talisman Proc", core.ActionID{ItemID: 32485}, stats.Stats{stats.Strength: 55}, time.Second*12)

		warrior.RegisterAura(core.Aura{
			Label:    "Ashtongue Talisman",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell != warrior.ShieldSlam && spell != warrior.Bloodthirst && spell != warrior.MortalStrike {
					return
				}

				if sim.RandomFloat("AshtongueTalismanOfValor") < 0.25 {
					procAura.Activate(sim)
				}
			},
		})
	})

}
