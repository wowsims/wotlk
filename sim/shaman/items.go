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
	core.NewItemEffect(19344, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		const dur = time.Second * 20
		actionID := core.ActionID{ItemID: 19344}

		activeAura := shaman.GetOrRegisterAura(core.Aura{
			Label:    "Natural Alignment Crystal",
			ActionID: actionID,
			Duration: dur,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.SpellPower, 250)
				shaman.PseudoStats.CostMultiplier *= 1.2
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.AddStatDynamic(sim, stats.SpellPower, -250)
				shaman.PseudoStats.CostMultiplier /= 1.2
			},
		})

		spell := shaman.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    shaman.NewTimer(),
					Duration: time.Minute * 5,
				},
				SharedCD: core.Cooldown{
					Timer:    shaman.GetOffensiveTrinketCD(),
					Duration: dur,
				},
			},
			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				activeAura.Activate(sim)
			},
		})

		shaman.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})

	// ActivateFathomBrooch adds an aura that has a chance on cast of nature spell
	//  to restore 335 mana. 40s ICD
	core.NewItemEffect(30663, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		icd := core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: time.Second * 40,
		}
		manaMetrics := shaman.NewManaMetrics(core.ActionID{ItemID: 30663})

		shaman.RegisterAura(core.Aura{
			Label:    "Fathom Brooch of the Tidewalker",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !icd.IsReady(sim) {
					return
				}
				if spell.SpellSchool != core.SpellSchoolNature {
					return
				}
				if sim.RandomFloat("Fathom-Brooch of the Tidewalker") > 0.15 {
					return
				}
				icd.Use(sim)
				shaman.AddMana(sim, 335, manaMetrics, false)
			},
		})
	})

	core.NewItemEffect(32491, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		procAura := shaman.NewTemporaryStatsAura("Ashtongue Talisman of Vision Proc", core.ActionID{ItemID: 32491}, stats.Stats{stats.AttackPower: 275}, time.Second*10)

		shaman.RegisterAura(core.Aura{
			Label:    "Ashtongue Talisman of Vision",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				// Note that shaman.Stormstrike is the first 'fake' SS hit.
				if spell != shaman.Stormstrike {
					return
				}
				if sim.RandomFloat("Ashtongue Talisman of Vision") > 0.5 {
					return
				}
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(33506, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		procAura := shaman.NewTemporaryStatsAura("Skycall Totem Proc", core.ActionID{ItemID: 33506}, stats.Stats{stats.SpellHaste: 101}, time.Second*10)

		icd := core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: time.Second * 30,
		}
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
				if (spell == shaman.LightningBolt || spell == shaman.LightningBoltLO) && sim.RandomFloat("Skycall Totem") < 0.15 {
					procAura.Activate(sim)
					icd.Use(sim)
				}
			},
		})
	})

	core.NewItemEffect(33507, func(agent core.Agent) {
		shaman := agent.(ShamanAgent).GetShaman()
		procAura := shaman.NewTemporaryStatsAura("Stonebreakers Totem Proc", core.ActionID{ItemID: 33507}, stats.Stats{stats.AttackPower: 110}, time.Second*10)

		icd := core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: time.Second * 10,
		}
		const procChance = 0.5

		shaman.RegisterAura(core.Aura{
			Label:    "Stonebreakers Totem",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
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

	// Even though these item effects are handled elsewhere, add them so they are
	// detected for automatic testing.
	core.NewItemEffect(TotemOfThePulsingEarth, func(core.Agent) {})
}
