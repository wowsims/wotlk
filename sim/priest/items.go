package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var ItemSetAbsolution = core.NewItemSet(core.ItemSet{
	Name: "Absolution Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in swp.go
		},
		4: func(agent core.Agent) {
			// Implemented in mindblast.go
		},
	},
})

var ItemSetVestmentsOfAbsolution = core.NewItemSet(core.ItemSet{
	Name: "Vestments of Absolution",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in prayer_of_healing.go
		},
		4: func(agent core.Agent) {
			// Implemented in greater_heal.go
		},
	},
})

var ItemSetValorous = core.NewItemSet(core.ItemSet{
	Name: "Garb of Faith",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// this is implemented in mind_blast.go
		},
		4: func(agent core.Agent) {
			// this is implemented in shadow_word_death.go
		},
	},
})

var ItemSetRegaliaOfFaith = core.NewItemSet(core.ItemSet{
	Name: "Regalia of Faith",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in prayer_of_mending.go
		},
		4: func(agent core.Agent) {
			// Implemented in greater_heal.go
		},
	},
})

var ItemSetConquerorSanct = core.NewItemSet(core.ItemSet{
	Name: "Sanctification Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in devouring_plague.go
		},
		4: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			procAura := priest.NewTemporaryStatsAura("Devious Mind", core.ActionID{SpellID: 64907}, stats.Stats{stats.SpellHaste: 240}, time.Second*4)

			priest.RegisterAura(core.Aura{
				Label:    "Devious Mind Proc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				// TODO: Does this affect the spell that procs it?
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell == priest.MindBlast {
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})

var ItemSetSanctificationRegalia = core.NewItemSet(core.ItemSet{
	Name: "Sanctification Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
		},
		4: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			procAura := priest.NewTemporaryStatsAura("Sanctification Reglia 4pc", core.ActionID{SpellID: 64912}, stats.Stats{stats.SpellPower: 250}, time.Second*5)

			priest.RegisterAura(core.Aura{
				Label:    "Sancitifcation Reglia 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				// TODO: Does this affect the spell that procs it?
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell == priest.PowerWordShield {
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})

var ItemSetZabras = core.NewItemSet(core.ItemSet{
	Name:            "Zabra's Regalia",
	AlternativeName: "Velen's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in vampiric_touch.go
		},
		4: func(agent core.Agent) {
			// Implemented in mind_flay.go
		},
	},
})

var ItemSetZabrasRaiment = core.NewItemSet(core.ItemSet{
	Name:            "Zabra's Raiment",
	AlternativeName: "Velen's Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in prayer_of_mending.go
		},
		4: func(agent core.Agent) {
			// Implemented in talents.go and renew.go
		},
	},
})

var ItemSetCrimsonAcolyte = core.NewItemSet(core.ItemSet{
	Name: "Crimson Acolyte's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in vampiric_touch.go/devouring_plague.go/swp.go
		},
		4: func(agent core.Agent) {
			// Implemented in mind_flay.go
		},
	},
})

var ItemSetCrimsonAcolytesRaiment = core.NewItemSet(core.ItemSet{
	Name: "Crimson Acolyte's Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()

			var curAmount float64
			procSpell := priest.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 70770},
				SpellSchool: core.SpellSchoolHoly,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers | core.SpellFlagHelpful,

				DamageMultiplier: 1,
				ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

				Hot: core.DotConfig{
					Aura: core.Aura{
						Label: "CrimsonAcolyteRaiment2pc",
					},
					NumberOfTicks: 3,
					TickLength:    time.Second * 3,
					OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
						dot.SnapshotBaseDamage = curAmount * 0.33
						dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
					},
					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
					},
				},
			})

			priest.RegisterAura(core.Aura{
				Label:    "Crimson Acolytes Raiment 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell != priest.FlashHeal || sim.RandomFloat("Crimson Acolytes Raiment 2pc") >= 0.33 {
						return
					}

					curAmount = result.Damage
					hot := procSpell.Hot(result.Target)
					hot.Apply(sim)
				},
			})
		},
		4: func(agent core.Agent) {
			// Implemented in power_word_shield.go and circle_of_healing.go
		},
	},
})

var ItemSetGladiatorsInvestiture = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Investiture",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.Resilience, 100)
			agent.GetCharacter().AddStat(stats.SpellPower, 29)
		},
		4: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.SpellPower, 88)
		},
	},
})
var ItemSetGladiatorsRaiment = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.Resilience, 100)
			agent.GetCharacter().AddStat(stats.SpellPower, 29)
		},
		4: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.SpellPower, 88)
		},
	},
})

func init() {
}
