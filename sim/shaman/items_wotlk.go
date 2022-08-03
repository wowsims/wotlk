package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var itemSetT9Bonuses = map[int32]core.ApplyEffect{
	2: func(agent core.Agent) {
		//
	},
	4: func(agent core.Agent) {
		//
	},
}

var ItemSetThrallsRegalia = core.NewItemSet(core.ItemSet{
	Name:    "Thrall's Regalia",
	Bonuses: itemSetT9Bonuses,
})
var ItemSetNobundosRegalia = core.NewItemSet(core.ItemSet{
	Name:    "Nobundo's Regalia",
	Bonuses: itemSetT9Bonuses,
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
		},
		4: func(agent core.Agent) {
		},
	},
})

var ItemSetFrostWitchRegalia = core.NewItemSet(core.ItemSet{
	Name: "Frost Witch's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// TODO: Your Lightning Bolt and Chain Lightning spells reduce the remaining cooldown on your Elemental Mastery talent by 2 sec.
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.RegisterAura(core.Aura{
				Label:    "Shaman T10 Elemental 4P Bonus",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spell == shaman.LavaBurst && shaman.FlameShockDot.IsActive() { // Doesn't have to hit from tooltip
						// Modify dot to last 6 more seconds than it has left, and refresh aura
						shaman.FlameShockDot.Duration = shaman.FlameShockDot.RemainingDuration(sim) + time.Second*6
						shaman.FlameShockDot.Refresh(sim)
					}
				},
			})
		},
	},
})

func init() {
	core.NewItemEffect(40708, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		procAura := shaman.NewTemporaryStatsAura("Totem of the Elemental Plane Proc", core.ActionID{ItemID: 40708}, stats.Stats{stats.SpellHaste: 196, stats.MeleeHaste: 196}, time.Second*10)

		icd := core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: time.Second * 30,
		}
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
				if (spell == shaman.LightningBolt || spell == shaman.LightningBoltLO) && sim.RandomFloat("totem of elemental plane") < 0.15 {
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
				if (spell == shaman.LightningBolt || spell == shaman.LightningBoltLO) && sim.RandomFloat("totem of elemental plane") < 0.7 {
					procAura.Activate(sim)
					icd.Use(sim) // put on CD
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

var itemSetEnhanceT9Bonuses = map[int32]core.ApplyEffect{
	2: func(agent core.Agent) {
		// +3% increase to static shock proc rate
	},
	4: func(agent core.Agent) {
		// +25% shock damage
	},
}

var ItemSetThrallsBattlegear = core.NewItemSet(core.ItemSet{
	Name:    "Thrall's Battlegear",
	Bonuses: itemSetEnhanceT9Bonuses,
})
var ItemSetNobundosBattlegear = core.NewItemSet(core.ItemSet{
	Name:    "Nobundo's Battlegear",
	Bonuses: itemSetEnhanceT9Bonuses,
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
