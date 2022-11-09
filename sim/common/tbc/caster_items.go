package tbc

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Keep these (and their functions) in alphabetical order.
func init() {
	// Proc effects. Keep these in order by item ID.
	core.AddEffectsToTest = false

	core.NewItemEffect(29305, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Band of the Eternal Sage Proc", core.ActionID{ItemID: 29305}, stats.Stats{stats.SpellPower: 95}, time.Second*10)

		// Your offensive spells have a chance on hit to increase your spell damage by 95 for 10 secs.
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 60,
		}
		const proc = 0.1

		character.RegisterAura(core.Aura{
			Label:    "Band of the Eternal Sage",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !result.Landed() {
					return
				}
				if !icd.IsReady(sim) || sim.RandomFloat("Band of the Eternal Sage") > proc { // can't activate if on CD or didn't proc
					return
				}
				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	core.AddEffectsToTest = true
}
