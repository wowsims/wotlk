package priest

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var ItemSetIncarnate = core.NewItemSet(core.ItemSet{
	Name: "Incarnate Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your shadowfiend now has 75 more stamina and lasts 3 sec. longer.
			// Implemented in shadowfiend.go.
		},
		4: func(agent core.Agent) {
			// Your Mind Flay and Smite spells deal 5% more damage.
			// Implemented in mind_flay.go.
		},
	},
})

var ItemSetAvatar = core.NewItemSet(core.ItemSet{
	Name: "Avatar Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			manaMetrics := priest.NewManaMetrics(core.ActionID{SpellID: 37600})

			priest.RegisterAura(core.Aura{
				Label:    "Avatar Regalia 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if sim.RandomFloat("avatar 2p") > 0.06 {
						return
					}
					// This is a cheat...
					// easier than adding another aura the subtracts 150 mana from next cast.
					priest.AddMana(sim, 150, manaMetrics, false)
				},
			})
		},
		4: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()

			procAura := priest.NewTemporaryStatsAura("Avatar Regalia 4pc Proc", core.ActionID{SpellID: 37604}, stats.Stats{stats.SpellPower: 100}, time.Second*15)
			procAura.OnSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				aura.Deactivate(sim)
			}

			priest.RegisterAura(core.Aura{
				Label:    "Avatar Regalia 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spell != priest.ShadowWordPain {
						return
					}

					if sim.RandomFloat("avatar 4p") > 0.4 { // 60% chance of not activating.
						return
					}

					procAura.Activate(sim)
				},
			})
		},
	},
})

var ItemSetAbsolution = core.NewItemSet(core.ItemSet{
	Name: "Absolution Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// this is implemented in swp.go
		},
		4: func(agent core.Agent) {
			// this is implemented in mindblast.go
		},
	},
})

func init() {

	core.NewItemEffect(32490, func(agent core.Agent) {
		priest := agent.(PriestAgent).GetPriest()
		// Not in the game yet so cant test; this logic assumes that:
		// - procrate is 10%
		// - no ICD on proc
		const procrate = 0.1
		procAura := priest.NewTemporaryStatsAura("Ashtongue Talisman Proc", core.ActionID{ItemID: 32490}, stats.Stats{stats.SpellPower: 220}, time.Second*10)

		priest.RegisterAura(core.Aura{
			Label:    "Ashtongue Talisman",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell != priest.ShadowWordPain {
					return
				}

				if sim.RandomFloat("Ashtongue Talisman of Acumen") > procrate {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

}
