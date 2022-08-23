package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// TBC Dungeon Set
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

// T4
var ItemSetVoidheartRaiment = core.NewItemSet(core.ItemSet{
	Name: "Voidheart Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			shadowBonus := warlock.NewTemporaryStatsAura(
				"Shadowflame",
				core.ActionID{SpellID: 37377},
				stats.Stats{stats.ShadowSpellPower: 135},
				time.Second*15,
			)

			fireBonus := warlock.NewTemporaryStatsAura(
				"Shadowflame Hellfire",
				core.ActionID{SpellID: 39437},
				stats.Stats{stats.ShadowSpellPower: 135},
				time.Second*15,
			)

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
			// implemented in immolate.go and corruption.go
		},
	},
})

// T5
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

// T6
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

// T7
var ItemSetPlagueheartGarb = core.NewItemSet(core.ItemSet{
	Name: "Plagueheart Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			DemonicSoulAura := warlock.RegisterAura(core.Aura{
				Label:    "Demonic Soul",
				ActionID: core.ActionID{SpellID: 61595},
				Duration: time.Second * 10,
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell == warlock.ShadowBolt || spell == warlock.Incinerate {
						aura.Deactivate(sim)
					}
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warlock.ShadowBolt.BonusCritRating += 10 * core.CritRatingPerCritChance
					warlock.Incinerate.BonusCritRating += 10 * core.CritRatingPerCritChance
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warlock.ShadowBolt.BonusCritRating -= 10 * core.CritRatingPerCritChance
					warlock.Incinerate.BonusCritRating -= 10 * core.CritRatingPerCritChance
				},
			})

			warlock.RegisterAura(core.Aura{
				Label:    "2pT7 Hidden Aura",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spell == warlock.Corruption || spell == warlock.Immolate {
						if sim.RandomFloat("2pT7") < 0.15 {
							DemonicSoulAura.Activate(sim)
							DemonicSoulAura.Refresh(sim)
						}
					}
				},
			})
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			SpiritsoftheDamnedAura := warlock.RegisterAura(core.Aura{
				Label:    "Spirits of the Damned",
				ActionID: core.ActionID{SpellID: 61082},
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.AddStatDynamic(sim, stats.Spirit, 300.)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.AddStatDynamic(sim, stats.Spirit, -300.)
				},
			})

			warlock.RegisterAura(core.Aura{
				Label:    "4pT7 Hidden Aura",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spell == warlock.LifeTap {
						if SpiritsoftheDamnedAura.IsActive() {
							SpiritsoftheDamnedAura.Refresh(sim)
						} else {
							SpiritsoftheDamnedAura.Activate(sim)
						}
					}
				},
			})
		},
	},
})

// T8
var ItemSetDeathbringerGarb = core.NewItemSet(core.ItemSet{
	Name: "Deathbringer Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented
		},
		4: func(agent core.Agent) {
			// Implemented
		},
	},
})

// T9
var ItemSetGuldansRegalia = core.NewItemSet(core.ItemSet{
	Name: "Gul'dan's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			pet := warlock.Pets[0].GetCharacter()
			pet.AddStats(stats.Stats{
				stats.MeleeCrit: 10 * core.CritRatingPerCritChance,
				stats.SpellCrit: 10 * core.CritRatingPerCritChance,
			})
		},
		4: func(agent core.Agent) {
			// Implemented
		},
	},
})

// T10
var ItemSetDarkCovensRegalia = core.NewItemSet(core.ItemSet{
	Name: "Dark Coven's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			pet := warlock.Pets[0].GetCharacter()

			DeviousMindsAura := warlock.RegisterAura(core.Aura{
				Label:    "Devious Minds",
				ActionID: core.ActionID{SpellID: 70840},
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.1
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
				},
			})

			petDeviousMindsAura := pet.RegisterAura(core.Aura{
				Label:    "Devious Minds",
				ActionID: core.ActionID{SpellID: 70840},
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.1
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
				},
			})

			warlock.RegisterAura(core.Aura{
				Label:    "4pT10 Hidden Aura",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if spell == warlock.UnstableAffliction || spell == warlock.Immolate {
						if sim.RandomFloat("4pT10") < 0.15 {
							DeviousMindsAura.Activate(sim)
							DeviousMindsAura.Refresh(sim)
							petDeviousMindsAura.Activate(sim)
							petDeviousMindsAura.Refresh(sim)
						}
					}
				},
			})
		},
	},
})

var ItemSetGladiatorsFelshroud = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Felshroud",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.AddStat(stats.SpellPower, 29)
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.AddStat(stats.SpellPower, 88)
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
