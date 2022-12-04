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

var T84PcProcChance = 0.1

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
