package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var ItemSetXXX = core.NewItemSet(core.ItemSet{
	Name: "",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// shaman := agent.(ShamanAgent).GetShaman()
		},
		4: func(agent core.Agent) {
			// Increases damage done by Lightning Bolt by 5%.
			// Implemented in lightning_bolt.go.
		},
	},
})

func init() {
	core.NewItemEffect(40708, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		procAura := shaman.NewTemporaryStatsAura("Totem of the Elemental Plane Proc", core.ActionID{ItemID: 40708}, stats.Stats{stats.SpellHaste: 196, stats.MeleeHaste: 196}, time.Second*10)

		shaman.RegisterAura(core.Aura{
			Label:    "Totem of the Elemental Plane",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if (spell == shaman.LightningBolt || spell == shaman.LightningBoltLO) && sim.RandomFloat("totem of elemental plane") < 0.15 {
					procAura.Activate(sim)
				}
			},
		})
	})
}

var ItemSetEarthshatterBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Earthshatter Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// 10% damage to lightning shield, handle in wherever its stored
		},
		4: func(agent core.Agent) {
			// +5% to flurry, handle in talents.go
		},
	},
})
