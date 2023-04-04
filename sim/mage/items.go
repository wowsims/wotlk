package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// T6 Sunwell
var ItemSetTempestRegalia = core.NewItemSet(core.ItemSet{
	Name: "Tempest Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the duration of your Evocation ability by 2 sec.
			// Implemented in evocation.go.
		},
		4: func(agent core.Agent) {
			// Increases the damage of your Fireball, Frostbolt, and Arcane Missles abilities by 5%.
			// Implemented in the files for those spells.
		},
	},
})

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
	Name:            "Kirin Tor Garb",
	AlternativeName: "Kirin'dor Garb", // Wowhead spells this incorrectly
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			procAura := mage.NewTemporaryStatsAura("Kirin Tor 2pc", core.ActionID{SpellID: 64868}, stats.Stats{stats.SpellPower: 350}, 15*time.Second)

			// Handle ICD ourselves since we use a custom check.
			icd := core.Cooldown{
				Timer:    agent.GetCharacter().NewTimer(),
				Duration: time.Second * 45,
			}

			core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
				Name:       "Mage2pT8",
				Callback:   core.CallbackOnCastComplete,
				ProcChance: 0.25,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if !icd.IsReady(sim) {
						return
					}

					if spell == mage.ArcaneBlast || spell == mage.Fireball || spell == mage.FrostfireBolt || spell == mage.Frostbolt {
						icd.Use(sim)
						procAura.Activate(sim)
					}
				},
			})
		},
		4: func(agent core.Agent) {
			//Implemented at 10% chance needs testing
		},
	},
})

const T84PcProcChance = 0.2

// T9
var ItemSetKhadgarsRegalia = core.NewItemSet(core.ItemSet{
	Name:            "Khadgar's Regalia",
	AlternativeName: "Sunstrider's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			//Implemented in initialization
		},
		4: func(agent core.Agent) {
			//Implemented in each spell
		},
	},
})

var ItemSetBloodmagesRegalia = core.NewItemSet(core.ItemSet{
	Name: "Bloodmage's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in each spell
		},
		4: func(agent core.Agent) {
			// Implemented in mirror_image.go
		},
	},
})

func (mage *Mage) BloodmagesRegalia2pcAura() *core.Aura {
	if !mage.HasSetBonus(ItemSetBloodmagesRegalia, 2) {
		return nil
	}

	return mage.GetOrRegisterAura(core.Aura{
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
}

var ItemSetGladiatorsRegalia = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Regalia",
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
