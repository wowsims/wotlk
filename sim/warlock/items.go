package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var ItemSetOblivionRaiment = core.NewItemSet(core.ItemSet{
	Name: "Oblivion Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// in pet.go constructor
		},
		4: func(agent core.Agent) {
			// in seed.go
		},
	},
})

var ItemSetVoidheartRaiment = core.NewItemSet(core.ItemSet{
	Name: "Voidheart Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			shadowBonus := warlock.NewTemporaryStatsAura("Shadowflame", core.ActionID{SpellID: 37377}, stats.Stats{stats.ShadowSpellPower: 135}, time.Second*15)
			fireBonus := warlock.NewTemporaryStatsAura("Shadowflame Hellfire", core.ActionID{SpellID: 39437}, stats.Stats{stats.ShadowSpellPower: 135}, time.Second*15)

			warlock.RegisterAura(core.Aura{
				Label:    "Voidheart Raiment 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if sim.RandomFloat("cycl4p") > 0.05 {
						return
					}
					if spell.SpellSchool.Matches(core.SpellSchoolShadow) {
						shadowBonus.Activate(sim)
					}
					if spell.SpellSchool.Matches(core.SpellSchoolFire) {
						fireBonus.Activate(sim)
					}
				},
			})
		},
		4: func(agent core.Agent) {
			// implemented in immolate.go
			// TODO: add to corruption.go
		},
	},
})

var ItemSetCorruptorRaiment = core.NewItemSet(core.ItemSet{
	Name: "Corruptor Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// heals pet
		},
		4: func(agent core.Agent) {
			// TODO: increase corruption tick damage on target whenever shadowbolt hits.
		},
	},
})

var ItemSetMaleficRaiment = core.NewItemSet(core.ItemSet{
	Name: "Malefic Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// heals... not implemented yet
		},
		4: func(agent core.Agent) {
			// Increases damage done by shadowbolt and incinerate by 6%.
			// Implemented in shadowbolt.go and incinerate.go
		},
	},
})

func init() {
	core.NewItemEffect(19337, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()
		if warlock.Pet == nil {
			return
		}

		actionID := core.ActionID{ItemID: 19337}
		bbAura := warlock.Pet.NewTemporaryStatsAura("The Black Book", actionID, stats.Stats{stats.SpellPower: 200, stats.AttackPower: 325, stats.Armor: 1600}, time.Second*30)

		bbSpell := warlock.RegisterSpell(core.SpellConfig{
			ActionID: actionID,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    warlock.NewTimer(),
					Duration: time.Minute * 5,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				bbAura.Activate(sim)
			},
		})

		warlock.AddMajorCooldown(core.MajorCooldown{
			Spell: bbSpell,
			Type:  core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(30449, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()
		warlock.AddStat(stats.SpellPower, 48)
		if warlock.Pet != nil {
			warlock.Pet.AddStats(stats.Stats{
				stats.ArcaneResistance: 130,
				stats.FireResistance:   130,
				stats.FrostResistance:  130,
				stats.NatureResistance: 130,
				stats.ShadowResistance: 130,
			})
		}
	})

	core.NewItemEffect(32493, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()
		procAura := warlock.NewTemporaryStatsAura("Asghtongue Talisman Proc", core.ActionID{SpellID: 40478}, stats.Stats{stats.SpellPower: 220}, time.Second*5)

		warlock.RegisterAura(core.Aura{
			Label:    "Ashtongue Talisman",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spell == warlock.Corruption && sim.RandomFloat("Ashtongue Talisman of Insight") < 0.2 {
					procAura.Activate(sim)
				}
			},
		})
	})
}
