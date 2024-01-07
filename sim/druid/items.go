package druid

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)


var ItemSetBlackfathomSlayerLeather = core.NewItemSet(core.ItemSet{
	Name: "Blackfathom Slayer's Leather",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.AddStat(stats.AttackPower, 12)
		},
		3: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.AddStat(stats.MeleeHit, 1)
			druid.AddStat(stats.SpellHit, 1)
		},
	},
})

func init() {
}
