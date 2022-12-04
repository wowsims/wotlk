package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// T7 Naxx
var ItemSetFrostfireGarb = core.NewItemSet(core.ItemSet{
	Name: "Frostfire Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			//Implemented in mana gems
		},
		4: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			mage.bonusCritDamage += .05
		},
	},
})

// T8 Ulduar
var ItemSetKirinTorGarb = core.NewItemSet(core.ItemSet{
	Name: "Kirin Tor Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			procAura := mage.NewTemporaryStatsAura("Kiron Tor 2pc", core.ActionID{SpellID: 64867}, stats.Stats{stats.SpellPower: 350}, 15*time.Second)
			core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
				Name:       "Mage2pT8",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcChance: 0.25,
				ICD:        time.Second * 45,
				SpellFlags: BarrageSpells,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					procAura.Activate(sim)
				},
			})
		},
		4: func(agent core.Agent) {
			//Implemented at 10% chance needs testing
		},
	},
})

// T9
var ItemSetKhadgarsRegalia = core.NewItemSet(core.ItemSet{
	Name: "Khadgar's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			//Implemented in initialization
		},
		4: func(agent core.Agent) {
			//Implemented in each spell
		},
	},
})

// T9 horde
var ItemSetSunstridersRegalia = core.NewItemSet(core.ItemSet{
	Name: "Sunstrider's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			//Implemented in initialization
		},
		4: func(agent core.Agent) {
			//Implemented in each spell
		},
	},
})

var bloodmageHasteAura *core.Aura
var bloodmageDamageAura *core.Aura
var ItemSetBloodmagesRegalia = core.NewItemSet(core.ItemSet{
	Name: "Bloodmage's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter()
			bloodmageHasteAura = agent.GetCharacter().RegisterAura(core.Aura{
				Label:    "Spec Based Haste T10 2PC",
				ActionID: core.ActionID{SpellID: 70752},
				Duration: time.Second * 5,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.MultiplyCastSpeed(1.12)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.MultiplyCastSpeed(1 / 1.12)
				},
			})
		},
		4: func(agent core.Agent) {
			bloodmageDamageAura = agent.GetCharacter().RegisterAura(core.Aura{
				Label:    "Mirror Image Bonus Damage T10 4PC",
				ActionID: core.ActionID{SpellID: 70748},
				Duration: time.Second * 30,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					agent.GetCharacter().PseudoStats.DamageDealtMultiplier *= 1.18
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					agent.GetCharacter().PseudoStats.DamageDealtMultiplier /= 1.18
				},
			})
		},
	},
})
