package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func init() {
	core.NewSimpleStatOffensiveTrinketEffect(37873, stats.Stats{stats.SpellPower: 346}, time.Second*20, time.Minute*2) // Mark of the War Prisoner

	core.NewItemEffect(40682, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Now is the Time!", core.ActionID{SpellID: 60063}, stats.Stats{stats.SpellPower: 590}, time.Second*10)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 45,
		}

		character.RegisterAura(core.Aura{
			Label:    "Sundial of the Exiled",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.ProcMask.Matches(core.ProcMaskSpellDamage) {
					return
				}
				if !icd.IsReady(sim) || !spellEffect.Outcome.Matches(core.OutcomeCrit) {
					return
				}
				if sim.RandomFloat("Sundial of the Exiled") > 0.1 {
					return
				}
				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})
}
