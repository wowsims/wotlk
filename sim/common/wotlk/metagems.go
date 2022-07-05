package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func init() {
	// Keep these in order by item ID.

	// TODO: Destructive Skyflare (1% spell reflect)
	// TODO: Revitalizing Skyflare (3% increased critical healing effect)
	// TODO: Invigorating Earthsiege (heal on crits)

	core.NewItemEffect(41333, func(agent core.Agent) {
		agent.GetCharacter().AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.Intellect,
			Modifier: func(intellect float64, _ float64) float64 {
				return intellect * 1.02
			},
		})
	})

	core.NewItemEffect(41377, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ArcaneDamageTakenMultiplier *= 0.98
		character.PseudoStats.FireDamageTakenMultiplier *= 0.98
		character.PseudoStats.FrostDamageTakenMultiplier *= 0.98
		character.PseudoStats.HolyDamageTakenMultiplier *= 0.98
		character.PseudoStats.NatureDamageTakenMultiplier *= 0.98
		character.PseudoStats.ShadowDamageTakenMultiplier *= 0.98
	})

	core.NewItemEffect(41380, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.AddStat(stats.Armor, character.Equip.Stats()[stats.Armor]*0.02)
	})

	core.NewItemEffect(41389, func(agent core.Agent) {
		agent.GetCharacter().AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Mana,
			ModifiedStat: stats.Mana,
			Modifier: func(mana float64, _ float64) float64 {
				return mana * 1.02
			},
		})
	})

	core.NewItemEffect(41395, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 0.98
	})

	core.NewItemEffect(41396, func(agent core.Agent) {
		agent.GetCharacter().AddStatDependency(stats.StatDependency{
			SourceStat:   stats.BlockValue,
			ModifiedStat: stats.BlockValue,
			Modifier: func(bv float64, _ float64) float64 {
				return bv * 1.05
			},
		})
	})

	core.NewItemEffect(41400, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Thundering Skyflare Diamond Proc", core.ActionID{SpellID: 55379}, stats.Stats{stats.MeleeHaste: 480, stats.SpellHaste: 480}, time.Second*6)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 40,
		}
		ppmm := character.AutoAttacks.NewPPMManager(1.5, core.ProcMaskMeleeOrRanged)

		character.RegisterAura(core.Aura{
			Label:    "Thundering Skyflare Diamond",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// Mask 68, melee or ranged auto attacks.
				if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskWhiteHit) {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if !ppmm.Proc(sim, spellEffect.ProcMask, "Thundering Skyflare Diamond") {
					return
				}
				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(41401, func(agent core.Agent) {
		character := agent.GetCharacter()
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 15,
		}
		manaMetrics := character.NewManaMetrics(core.ActionID{ItemID: 41401})

		character.RegisterAura(core.Aura{
			Label:    "Insightful Earthsiege Diamond",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !icd.IsReady(sim) || sim.RandomFloat("Insightful Earthsiege Diamond") > 0.05 {
					return
				}
				icd.Use(sim)
				character.AddMana(sim, 600, manaMetrics, false)
			},
		})
	})

	// These are handled in character.go, but create empty effects so they are included in tests.
	core.NewItemEffect(41285, func(_ core.Agent) {}) // Chaotic Skyflare Diamond
	core.NewItemEffect(41398, func(_ core.Agent) {}) // Relentless Earthsiege Diamond
}
