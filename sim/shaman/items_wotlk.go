package shaman

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var ItemSetThrallsRegalia = core.NewItemSet(core.ItemSet{
	Name:            "Thrall's Regalia",
	AlternativeName: "Nobundo's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// shocks.go
		},
		4: func(agent core.Agent) {
			// lavaburst.go
		},
	},
})

var ItemSetEarthShatterGarb = core.NewItemSet(core.ItemSet{
	Name: "Earthshatter Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Reduces LB cost by 5%
		},
		4: func(agent core.Agent) {
			// lavaburst crit strike dmg +10%
		},
	},
})
var ItemSetWorldbreakerGarb = core.NewItemSet(core.ItemSet{
	Name: "Worldbreaker Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// shocks.go
		},
		4: func(agent core.Agent) {
			// lightning_bolt.go
		},
	},
})

var ItemSetFrostWitchRegalia = core.NewItemSet(core.ItemSet{
	Name: "Frost Witch's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// This is implemented in talents.go so that the aura has easy access to the elemental mastery MCD.
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.RegisterAura(core.Aura{
				Label:    "Shaman T10 Elemental 4P Bonus",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					fsDot := shaman.FlameShock.Dot(shaman.CurrentTarget)
					if spell == shaman.LavaBurst && fsDot.IsActive() { // Doesn't have to hit from tooltip
						// 4p t10 immediately updates the tickPeriod based on current haste
						fsDot.RescheduleNextTick(sim)

						// Find the number of ticks whose duration is closest to 6s.
						// "our testing confirms that the 4pc t10 setbonus adds to FS the closest number of ticks to 6 seconds always"
						// https://web.archive.org/web/20100808192139/http://elitistjerks.com/f79/t76510-elemental_patch_3_3_now_more_fire_nova/p25/
						tickPeriod := fsDot.TickPeriod()
						numTicks := int32(math.Round(float64(time.Second) * 6 / float64(tickPeriod)))
						fsDot.NumberOfTicks += numTicks

						// Set duration to remaining ticks, minus the elapsed time since last tick
						fsDot.Aura.Duration = time.Duration(fsDot.MaxTicksRemaining())*tickPeriod - (tickPeriod - (fsDot.NextTickAt() - sim.CurrentTime))
						fsDot.Aura.Refresh(sim) // update aura's duration
					}
				},
			})
		},
	},
})

func init() {
	core.NewItemEffect(40708, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		procAura := shaman.NewTemporaryStatsAura("Totem of the Elemental Plane Proc", core.ActionID{SpellID: 60771}, stats.Stats{stats.SpellHaste: 196, stats.MeleeHaste: 196}, time.Second*10)

		icd := core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: time.Second * 30,
		}
		procAura.Icd = &icd
		shaman.RegisterAura(core.Aura{
			Label:    "Totem of the Elemental Plane",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !icd.IsReady(sim) {
					return
				}
				if spell.ActionID.SpellID == 49238 && sim.RandomFloat("totem of elemental plane") < 0.15 {
					procAura.Activate(sim)
					icd.Use(sim)
				}
			},
		})
	})

	core.NewItemEffect(47666, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		procAura := shaman.NewTemporaryStatsAura("ToEW - Energized", core.ActionID{SpellID: 67385}, stats.Stats{stats.SpellHaste: 200, stats.MeleeHaste: 200}, time.Second*12)

		icd := core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: time.Second * 6,
		}
		procAura.Icd = &icd
		shaman.RegisterAura(core.Aura{
			Label:    "Totem of Electrifying Wind",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !icd.IsReady(sim) {
					return
				}
				if spell.ActionID.SpellID == 49238 && sim.RandomFloat("totem of elemental plane") < 0.7 {
					procAura.Activate(sim)
					icd.Use(sim) // put on CD
				}
			},
		})
	})

	core.NewItemEffect(50463, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		procAura := shaman.RegisterAura(core.Aura{
			Label:     "Enraged",
			ActionID:  core.ActionID{SpellID: 71216},
			Duration:  time.Second * 15,
			MaxStacks: 3,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				shaman.AddStatDynamic(sim, stats.AttackPower, -146*float64(oldStacks))
				shaman.AddStatDynamic(sim, stats.AttackPower, 146*float64(newStacks))
			},
		})
		shaman.RegisterAura(core.Aura{
			Label:    "Totem of the Avalanche",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() { //TODO: verify it needs to land
					return
				}
				if spell == shaman.Stormstrike {
					procAura.Activate(sim)
					procAura.AddStack(sim)
				}
			},
		})
	})

	// Bizuri's Totem of Shattered Ice
	core.NewItemEffect(50458, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		procAura := shaman.RegisterAura(core.Aura{
			Label:     "Furious",
			ActionID:  core.ActionID{SpellID: 71199},
			Duration:  time.Second * 30,
			MaxStacks: 5,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				shaman.AddStatsDynamic(sim, stats.Stats{
					stats.SpellHaste: 44 * float64(newStacks-oldStacks),
					stats.MeleeHaste: 44 * float64(newStacks-oldStacks),
				})
			},
		})
		shaman.RegisterAura(core.Aura{
			Label:    "Bizuri's Totem of Shattered Ice Aura",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == shaman.FlameShock {
					procAura.Activate(sim)
					procAura.AddStack(sim)
				}
			},
		})
	})
}

var ItemSetEarthshatterBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Earthshatter Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// 10% damage to lightning shield. implemented in lightning_shield.go
		},
		4: func(agent core.Agent) {
			// +5% to flurry. implemented in talents.go
		},
	},
})

var ItemSetWorldbreakerBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Worldbreaker Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			//20% damage to stormstrike and lava lash
		},
		4: func(agent core.Agent) {
			//20% increase to maelstrom proc rate
		},
	},
})

var ItemSetThrallsBattlegear = core.NewItemSet(core.ItemSet{
	Name:            "Thrall's Battlegear",
	AlternativeName: "Nobundo's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// +3% increase to static shock proc rate
		},
		4: func(agent core.Agent) {
			// +25% shock damage
		},
	},
})

var ItemSetFrostWitchBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Frost Witch's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// TODO: add 12% damage buff to shamanistic rage
		},
		4: func(agent core.Agent) {
			// TODO: at 5 maelstrom stacks, 15% chance to gain +20% attack power for 10s
		},
	},
})

var ItemSetGladiatorsEarthshaker = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Earthshaker",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.AddStat(stats.AttackPower, 50)
			shaman.AddStat(stats.Resilience, 100)
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.AddStat(stats.AttackPower, 150)
			// also -2s on stormstrike CD
		},
	},
})

var ItemSetGladiatorsWartide = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Wartide",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.AddStat(stats.SpellPower, 29)
			shaman.AddStat(stats.Resilience, 100)
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.AddStat(stats.SpellPower, 88)
		},
	},
})
