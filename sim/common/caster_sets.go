package common

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Keep these in alphabetical order.

var ItemSetManaEtched = core.NewItemSet(core.ItemSet{
	Name: "Mana-Etched Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.SpellHit, 35)
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			procAura := character.NewTemporaryStatsAura("Mana-Etched Insight Proc", core.ActionID{SpellID: 37619}, stats.Stats{stats.SpellPower: 110}, time.Second*15)

			character.RegisterAura(core.Aura{
				Label:    "Mana-Etched Insight",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if sim.RandomFloat("Mana-Etched Insight") > 0.02 {
						return
					}
					procAura.Activate(sim)
				},
			})
		},
	},
})

var ItemSetNetherstrike = core.NewItemSet(core.ItemSet{
	Name: "Netherstrike Armor",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.SpellPower, 23)
		},
	},
})

var ItemSetSpellstrike = core.NewItemSet(core.ItemSet{
	Name: "Spellstrike Infusion",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			procAura := character.NewTemporaryStatsAura("Spellstrike Proc", core.ActionID{SpellID: 32106}, stats.Stats{stats.SpellPower: 92}, time.Second*10)

			character.RegisterAura(core.Aura{
				Label:    "Spellstrike",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if sim.RandomFloat("spellstrike") > 0.05 {
						return
					}
					procAura.Activate(sim)
				},
			})
		},
	},
})

var ItemSetTheTwinStars = core.NewItemSet(core.ItemSet{
	Name: "The Twin Stars",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.SpellPower, 15)
		},
	},
})

var ItemSetWindhawk = core.NewItemSet(core.ItemSet{
	Name: "Windhawk Armor",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.MP5, 8)
		},
	},
})

var ItemSetSpellfire = core.NewItemSet(core.ItemSet{
	Name: "Wrath of Spellfire",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStatDependency(stats.StatDependency{
				SourceStat:   stats.Intellect,
				ModifiedStat: stats.SpellPower,
				Modifier: func(intellect float64, spellPower float64) float64 {
					return spellPower + intellect*0.07 // 7% bonus to sp from int
				},
			})
		},
	},
})
