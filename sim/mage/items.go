package mage

import (
	"time"

	"github.com/wowsims/tbc/sim/common"
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

const SerpentCoilBraidID = 30720

var ItemSetAldorRegalia = core.NewItemSet(core.ItemSet{
	Name: "Aldor Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Interruption avoidance.
		},
		4: func(agent core.Agent) {
			// Reduces the cooldown on PoM/Blast Wave/Ice Block.
		},
	},
})

var ItemSetTirisfalRegalia = core.NewItemSet(core.ItemSet{
	Name: "Tirisfal Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the damage and mana cost of Arcane Blast by 20%.
			// Implemented in arcane_blast.go.
		},
		4: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			// Your spell critical strikes grant you up to 70 spell damage for 6 sec.
			procAura := mage.NewTemporaryStatsAura("Tirisfal 4pc Proc", core.ActionID{SpellID: 37443}, stats.Stats{stats.SpellPower: 70}, time.Second*6)
			mage.RegisterAura(core.Aura{
				Label:    "Tirisfal 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
						return
					}
					if spellEffect.Outcome.Matches(core.OutcomeCrit) {
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})

var ItemSetTempestRegalia = core.NewItemSet(core.ItemSet{
	Name: "Tempest Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the duratoin of your Evocation ability by 2 sec.
			// Implemented in mage.go.
		},
		4: func(agent core.Agent) {
			// Increases the damage of your Fireball, Frostbolt, and Arcane Missles abilities by 5%.
			// Implemented in the files for those spells.
		},
	},
})

func init() {
	common.NewSimpleStatOffensiveTrinketEffect(19339, stats.Stats{stats.SpellHaste: 330}, time.Second*20, time.Minute*5) // MQG

	core.NewItemEffect(32488, func(agent core.Agent) {
		mage := agent.(MageAgent).GetMage()
		procAura := mage.NewTemporaryStatsAura("Asghtongue Talisman Proc", core.ActionID{SpellID: 32488}, stats.Stats{stats.SpellHaste: 150}, time.Second*5)

		mage.RegisterAura(core.Aura{
			Label:    "Ashtongue Talisman",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
					return
				}

				if sim.RandomFloat("Ashtongue Talisman of Insight") > 0.5 {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	// Even though these item effects are handled elsewhere, add them so they are
	// detected for automatic testing.
	core.NewItemEffect(SerpentCoilBraidID, func(core.Agent) {})
}
