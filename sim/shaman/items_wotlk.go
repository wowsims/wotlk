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
			// TODO: Your Lava Burst spell causes your Flame Shock effect on the target to deal at least two additional periodic damage ticks before expiring.
			//  This will actually just extend the FS by 6 seconds. This could be 2-4 more ticks depending on current haste.
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
			// 10% damage to lightning shield, handle in wherever its stored
		},
		4: func(agent core.Agent) {
			// +5% to flurry, handle in talents.go
		},
	},
})
