package sod

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetBlackfathomElementalistHide = core.NewItemSet(core.ItemSet{
	Name: "Blackfathom Elementalist's Hide",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellDamage, 9)
			c.AddStat(stats.Healing, 9)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
	},
})

var ItemSetBlackfathomInvokerVestaments = core.NewItemSet(core.ItemSet{
	Name: "Twilight Invoker's Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellDamage, 9)
			c.AddStat(stats.Healing, 9)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
	},
})
