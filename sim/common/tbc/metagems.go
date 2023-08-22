package tbc

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false
	// Keep these in order by item ID.

	core.NewItemEffect(25893, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Mystic Focus Proc", core.ActionID{ItemID: 25893}, stats.Stats{stats.SpellHaste: 320}, time.Second*4)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 35,
		}

		character.RegisterAura(core.Aura{
			Label:    "Mystical Skyfire Diamond",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !icd.IsReady(sim) || sim.RandomFloat("Mystical Skyfire Diamond") > 0.15 {
					return
				}
				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(25899, func(agent core.Agent) {
		agent.GetCharacter().PseudoStats.BonusDamage += 3
	})

	core.NewItemEffect(25901, func(agent core.Agent) {
		character := agent.GetCharacter()
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 15,
		}
		manaMetrics := character.NewManaMetrics(core.ActionID{ItemID: 25901})

		character.RegisterAura(core.Aura{
			Label:    "Insightful Earthstorm Diamond",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !icd.IsReady(sim) || sim.RandomFloat("Insightful Earthstorm Diamond") > 0.04 {
					return
				}
				icd.Use(sim)
				character.AddMana(sim, 300, manaMetrics)
			},
		})
	})

	core.NewItemEffect(32410, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Thundering Skyfire Diamond Proc", core.ActionID{ItemID: 32410}, stats.Stats{stats.MeleeHaste: 240}, time.Second*6)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 40,
		}
		ppmm := character.AutoAttacks.NewPPMManager(1.5, core.ProcMaskWhiteHit) // Mask 68, melee or ranged auto attacks.

		character.RegisterAura(core.Aura{
			Label:    "Thundering Skyfire Diamond",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if !icd.IsReady(sim) {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Thundering Skyfire Diamond") {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	// Eternal Earthstorm
	core.NewItemEffect(35501, func(agent core.Agent) {
		agent.GetCharacter().PseudoStats.BlockValueMultiplier += 0.05
	})

	core.NewItemEffect(35503, func(agent core.Agent) {
		agent.GetCharacter().MultiplyStat(stats.Intellect, 1.02)
	})

	// These are handled in character.go, but create empty effects, so they are included in tests.
	core.NewItemEffect(34220, func(_ core.Agent) {}) // Chaotic Skyfire Diamond
	core.NewItemEffect(32409, func(_ core.Agent) {}) // Relentless Earthstorm Diamond

	core.AddEffectsToTest = true
}
