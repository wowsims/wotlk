package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var ItemSetTidefury = core.NewItemSet(core.ItemSet{
	Name: "Tidefury Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Handled in chain_lightning.go
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			if shaman.SelfBuffs.Shield == proto.ShamanShield_WaterShield {
				shaman.AddStat(stats.MP5, 3)
			}
		},
	},
})

var ItemSetSkyshatterRegalia = core.NewItemSet(core.ItemSet{
	Name: "Skyshatter Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			if shaman.Totems.Air == proto.AirTotem_NoAirTotem ||
				shaman.Totems.Water == proto.WaterTotem_NoWaterTotem ||
				shaman.Totems.Earth == proto.EarthTotem_NoEarthTotem ||
				shaman.Totems.Fire == proto.FireTotem_NoFireTotem {
				return
			}

			shaman.AddStat(stats.MP5, 19)
			shaman.AddStat(stats.SpellCrit, 35)
			shaman.AddStat(stats.SpellPower, 45)
		},
		4: func(agent core.Agent) {
			// Increases damage done by Lightning Bolt by 5%.
			// Implemented in lightning_bolt.go.
		},
	},
})

// Skyshatter Harness
// 2 pieces: Your Earth Shock, Flame Shock, and Frost Shock abilities cost 10% less mana.
// 4 pieces: Whenever you use Stormstrike, you gain 70 attack power for 12 sec.

var ItemSetSkyshatterHarness = core.NewItemSet(core.ItemSet{
	Name: "Skyshatter Harness",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// implemented in shocks.go
		},
		4: func(agent core.Agent) {
			// implemented in stormstrike.go
		},
	},
})

func init() {
	core.NewItemEffect(33506, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		procAura := shaman.NewTemporaryStatsAura("Skycall Totem Proc", core.ActionID{SpellID: 43751}, stats.Stats{stats.SpellHaste: 101}, time.Second*10)

		icd := core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: time.Second * 30,
		}
		procAura.Icd = &icd
		shaman.RegisterAura(core.Aura{
			Label:    "Skycall Totem",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !icd.IsReady(sim) {
					return
				}
				if spell.ActionID.SpellID == 49238 && sim.RandomFloat("Skycall Totem") < 0.15 {
					procAura.Activate(sim)
					icd.Use(sim)
				}
			},
		})
	})

	core.NewItemEffect(33507, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		procAura := shaman.NewTemporaryStatsAura("Stonebreakers Totem Proc", core.ActionID{SpellID: 43749}, stats.Stats{stats.AttackPower: 110}, time.Second*10)

		icd := core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: time.Second * 10,
		}
		procAura.Icd = &icd
		const procChance = 0.5

		shaman.RegisterAura(core.Aura{
			Label:    "Stonebreakers Totem",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if !spell.Flags.Matches(SpellFlagShock) {
					return
				}

				if !icd.IsReady(sim) {
					return
				}

				if sim.RandomFloat("Stonebreakers Totem") > procChance {
					return
				}

				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	registerSpellPVPTotem("Savage Gladiator's Totem of Survival", 42594, 60565, 52, 6)
	registerSpellPVPTotem("Hateful Gladiator's Totem of Survival", 42601, 60566, 62, 6)
	registerSpellPVPTotem("Deadly Gladiator's Totem of Survival", 42602, 60567, 70, 10)
	registerSpellPVPTotem("Furious Gladiator's Totem of Survival", 42603, 60568, 84, 10)
	registerSpellPVPTotem("Relentless Gladiator's Totem of Survival", 42604, 60569, 101, 10)
	registerSpellPVPTotem("Wrathful Gladiator's Totem of Survival", 51513, 60570, 119, 10)

	core.NewItemEffect(47667, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()

		statAura := shaman.NewTemporaryStatsAura("Volcanic Fury", core.ActionID{SpellID: 67391}, stats.Stats{stats.AttackPower: 400}, time.Second*18)

		triggerAura := core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
			Name:       "Totem of Quaking Earth Trigger",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskMeleeOHSpecial,
			ProcChance: .80,
			ICD:        time.Second * 9,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == shaman.LavaLash {
					statAura.Activate(sim)
				}
			},
		})
		statAura.Icd = triggerAura.Icd
	})
}

func registerSpellPVPTotem(name string, itemId int32, spellId int32, sp float64, seconds float64) {
	core.NewItemEffect(itemId, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		procAura := shaman.NewTemporaryStatsAura(name+" proc", core.ActionID{SpellID: spellId}, stats.Stats{stats.SpellPower: sp}, time.Second*time.Duration(seconds))

		shaman.RegisterAura(core.Aura{
			Label:    name,
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if !spell.Flags.Matches(SpellFlagShock) {
					return
				}

				procAura.Activate(sim)
			},
		})
	})
}
