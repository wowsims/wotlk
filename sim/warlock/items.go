package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

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

			const bonusCrit = 10 * core.CritRatingPerCritChance
			warlock.DemonicSoulAura = warlock.RegisterAura(core.Aura{
				Label:    "Demonic Soul",
				ActionID: core.ActionID{SpellID: 61595},
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warlock.ShadowBolt.BonusCritRating += bonusCrit
					warlock.Incinerate.BonusCritRating += bonusCrit
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warlock.ShadowBolt.BonusCritRating -= bonusCrit
					warlock.Incinerate.BonusCritRating -= bonusCrit
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell == warlock.ShadowBolt || spell == warlock.Incinerate {
						warlock.DemonicSoulAura.Deactivate(sim)
					}
				},
			})

			warlock.RegisterAura(core.Aura{
				Label: "2pT7 Hidden Aura",
				// ActionID: core.ActionID{SpellID: 60170},
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if (spell == warlock.Corruption || spell == warlock.Immolate) && sim.Proc(0.15, "2pT7") {
						warlock.DemonicSoulAura.Activate(sim)
					}
				},
			})
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			warlock.SpiritsoftheDamnedAura = warlock.RegisterAura(core.Aura{
				Label:    "Spirits of the Damned",
				ActionID: core.ActionID{SpellID: 61082},
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.AddStatDynamic(sim, stats.Spirit, 300)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.AddStatDynamic(sim, stats.Spirit, -300)
				},
			})

			warlock.RegisterAura(core.Aura{
				Label:    "4pT7 Hidden Aura",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell == warlock.LifeTap {
						warlock.SpiritsoftheDamnedAura.Activate(sim)
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
	Name:            "Gul'dan's Regalia",
	AlternativeName: "Kel'Thuzad's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			if warlock.Pet != nil {
				warlock.Pet.AddStats(stats.Stats{
					stats.MeleeCrit: 10 * core.CritRatingPerCritChance,
					stats.SpellCrit: 10 * core.CritRatingPerCritChance,
				})
			}
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

			deviousMindsAura := warlock.RegisterAura(core.Aura{
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

			var petDeviousMindsAura *core.Aura
			if warlock.Pet != nil {
				petDeviousMindsAura = warlock.Pet.RegisterAura(core.Aura{
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
			}

			warlock.RegisterAura(core.Aura{
				Label:    "4pT10 Hidden Aura",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell == warlock.UnstableAffliction || spell == warlock.Immolate {
						if sim.Proc(0.15, "4pT10") {
							deviousMindsAura.Activate(sim)
							if petDeviousMindsAura != nil {
								petDeviousMindsAura.Activate(sim)
							}
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

// black book is only ever used pre fight, after which we switch to a real trinket. For this reason we implement it as a
// cooldown and only allow it being cast before combat starts during prepull actions.
func (warlock *Warlock) registerBlackBook() {
	if warlock.Options.Summon == proto.Warlock_Options_NoSummon {
		return
	}

	effectAura := warlock.Pet.NewTemporaryStatsAura("Blessing of the Black Book", core.ActionID{SpellID: 23720},
		stats.Stats{stats.SpellPower: 200, stats.AttackPower: 325, stats.Armor: 1600}, 30*time.Second)

	spell := warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 23720},
		SpellSchool: core.SpellSchoolShadow,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: 5 * time.Minute,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sim.CurrentTime < 0
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			effectAura.Activate(sim)
		},
	})

	warlock.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func init() {
	core.NewItemEffect(32493, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()
		procAura := warlock.NewTemporaryStatsAura("Ashtongue Talisman Proc", core.ActionID{SpellID: 40478}, stats.Stats{stats.SpellPower: 220}, time.Second*5)

		warlock.RegisterAura(core.Aura{
			Label:    "Ashtongue Talisman",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == warlock.Corruption && sim.Proc(0.2, "Ashtongue Talisman of Insight") {
					procAura.Activate(sim)
				}
			},
		})
	})
}
