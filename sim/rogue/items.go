package rogue

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

// TODO: Update Sets with Old Raid Tier sets and any new SoD sets

var Arena = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.Resilience, 100)
			agent.GetCharacter().AddStat(stats.AttackPower, 50)
		},
		4: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.AttackPower, 150)
			// 10 maximum energy added in rogue.go
		},
	},
})

func init() {
	core.NewItemEffect(32492, func(agent core.Agent) {
		rogue := agent.(RogueAgent).GetRogue()
		procAura := rogue.NewTemporaryStatsAura("Ashtongue Talisman Proc", core.ActionID{ItemID: 32492}, stats.Stats{stats.MeleeCrit: 145}, time.Second*10)

		var numPoints int32

		rogue.RegisterAura(core.Aura{
			Label:    "Ashtongue Talisman",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				numPoints = 0
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !spell.Flags.Matches(SpellFlagFinisher) {
					return
				}

				// Need to store the points because they get spent before OnSpellHit is called.
				numPoints = rogue.ComboPoints()

				if spell == rogue.SliceAndDice {
					// SND won't call OnSpellHit, so we have to add the effect now.
					if p := 0.2 * float64(numPoints); sim.Proc(p, "AshtongueTalismanOfLethality") {
						procAura.Activate(sim)
					}
				}
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !spell.Flags.Matches(SpellFlagFinisher) {
					return
				}

				if p := 0.2 * float64(numPoints); sim.Proc(p, "AshtongueTalismanOfLethality") {
					procAura.Activate(sim)
				}
			},
		})
	})

}
