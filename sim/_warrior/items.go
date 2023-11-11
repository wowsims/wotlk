package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

/////////////////////////////////////////////////////////////////
// TBC Item set
/////////////////////////////////////////////////////////////////

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

var ItemSetGladiatorsBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases attack power by 50.
			// +100 resilience rating.
			agent.GetCharacter().AddStat(stats.Resilience, 100)
			agent.GetCharacter().AddStat(stats.AttackPower, 50)
		},
		4: func(agent core.Agent) {
			// Reduces the cooldown of your Intercept ability by 5 sec.
			// Increases attack power by 150.
			agent.GetCharacter().AddStat(stats.AttackPower, 150)
		},
	},
})

var ItemSetDreadnaughtPlate = core.NewItemSet(core.ItemSet{
	Name: "Dreadnaught Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the damage of your Shield Slam ability by 10%.
			// Handled in shield_slam.go.
		},
		4: func(agent core.Agent) {
			// Increases the duration of Shield Wall by 3 seconds.
			// NYI
		},
	},
})

var ItemSetSiegebreakerPlate = core.NewItemSet(core.ItemSet{
	Name: "Siegebreaker Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the critical strike chance of Devastate by 10%.
			// Handled in devastate.go
		},
		4: func(agent core.Agent) {
			// Shield Block grants 10% magic DR
			// NYI
		},
	},
})

var ItemSetWrynnsPlate = core.NewItemSet(core.ItemSet{
	Name:            "Wrynn's Plate",
	AlternativeName: "Hellscream's Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Decreases the cooldown on Taunt by 2sec
			// NYI

			// Increases damage done by Devastate by 5%
			// Handled in devastate.go
		},
		4: func(agent core.Agent) {
			// Decreases the cooldown of Shield Block by 10 sec
			// NYI
		},
	},
})

var ItemSetYmirjarLordsPlate = core.NewItemSet(core.ItemSet{
	Name: "Ymirjar Lord's Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Shield Slam and Shockwave deal 20% increased damage
			// Handled in shield_slam.go and shockwave.go
		},
		4: func(agent core.Agent) {
			// Bloodrage no longer costs health to use, and now causes you to absorb damage equal to 20% max HP
			// NYI
		},
	},
})

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
				Label:    "Dreadnaught Battlegear 4pc Proc",
				ActionID: core.ActionID{SpellID: 61571},
				Duration: time.Second * 30,
				OnGain: func(_ *core.Aura, sim *core.Simulation) {
					warrior.PseudoStats.CostReduction += 5
				},
				OnExpire: func(_ *core.Aura, sim *core.Simulation) {
					warrior.PseudoStats.CostReduction -= 5
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if !spell.ProcMask.Matches(core.ProcMaskMeleeSpecial) {
						return
					}

					// one-shot cost reductions cannot be reliably reset, since both OnCastComplete and OnSpellHit
					//  are too late (e.g. there might be a proc after Slam cast, but before either callback),
					//  or happen to often (e.g. internal 0-cost casts like Mutilate, or multiple Whirlwind hits, in case of OnSpellHit)

					// doesn't handle multiple dynamic cost reductions at once, or 0-cost default casts
					if actualGain := spell.DefaultCast.Cost - spell.CurCast.Cost; actualGain > 0 {
						rageMetrics.AddEvent(5, actualGain)
						aura.Deactivate(sim)
					}
				},
			})

			warrior.RegisterAura(core.Aura{
				Label:    "Dreadnaught Battlegear 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnPeriodicDamageDealt: func(_ *core.Aura, sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
					if result.Landed() && sim.RandomFloat("Dreadnaught Battlegear 4pc") < 0.1 {
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
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.ActionID.SpellID != 47450 && spell != warrior.Slam {
						return
					}
					if !result.Outcome.Matches(core.OutcomeCrit) {
						return
					}
					if result.Landed() && sim.RandomFloat("Siegebreaker Battlegear 2pc") < 0.4 {
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
					bonusAP = warrior.GetStat(stats.AttackPower) * 0.16
					aura.Unit.AddStatDynamic(sim, stats.AttackPower, bonusAP)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.AddStatDynamic(sim, stats.AttackPower, -bonusAP)
				},
			})

			warrior.RegisterAura(core.Aura{
				Label:    "Ymirjar Lord's Battlegear 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell != warrior.DeepWounds {
						return
					}
					if result.Landed() && sim.RandomFloat("Ymirjar 2pc") < 0.03 {
						procAura.Activate(sim)
					}
				},
			})
		},
		4: func(agent core.Agent) {
			// You have a 20% chance for your Bloodsurge and Sudden Death talents to grant 2 charges of their effect instead of 1,
			// reduce the global cooldown on Execute or Slam by 0.5 sec, and for the duration of the effect to be increased by 100%.

			// handled with specialized Auras for either Bloodsurge or Sudden Death
		},
	},
})

func init() {
}
