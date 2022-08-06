package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

//T7 Naxx
var ItemSetFrostfireGarb = core.NewItemSet(core.ItemSet{
	Name: "Frostfire Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			//Implemented in mana gems
		},
		4: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()

			mage.bonusCritDamage += .05
		},
	},
})

//T8 Ulduar
var ItemSetKirinTorGarb = core.NewItemSet(core.ItemSet{
	Name: "Kirin Tor Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {

			applyProcAura := func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) bool {
				if !spell.Flags.Matches(BarrageSpells) {
					return false
				}

				return sim.RandomFloat("Mage2pT8") < .25

			}

			agent.GetCharacter().StatProcWithICD("Kirin Tor 2pc", core.ActionID{SpellID: 64867}, stats.Stats{stats.SpellPower: 350}, 15*time.Second, 45*time.Second, applyProcAura)

		},
		4: func(agent core.Agent) {
			//Implemented at 10% chance needs testing
		},
	},
})

var ItemSetKhadgarsRegalia = core.NewItemSet(core.ItemSet{
	Name: "Khadgar's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			//Implemented in initialization
		},
		4: func(agent core.Agent) {
			//Implemented in each spell
		},
	},
})

var ItemSetSunstridersRegalia = core.NewItemSet(core.ItemSet{
	Name: "Sunstrider's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			//Implemented in initialization
		},
		4: func(agent core.Agent) {
			//Implemented in each spell
		},
	},
})
