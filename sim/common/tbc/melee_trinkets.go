package tbc

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false

	// Offensive trinkets. Keep these in order by item ID.
	core.NewSimpleStatOffensiveTrinketEffect(29383, stats.Stats{stats.AttackPower: 278, stats.RangedAttackPower: 278}, time.Second*20, time.Minute*2) // Bloodlust Brooch
	core.NewSimpleStatOffensiveTrinketEffect(32658, stats.Stats{stats.Agility: 150}, time.Second*20, time.Minute*2)                                   // Badge of Tenacity
	core.NewSimpleStatOffensiveTrinketEffect(33831, stats.Stats{stats.AttackPower: 360, stats.RangedAttackPower: 360}, time.Second*20, time.Minute*2) // Berserkers Call
	core.NewSimpleStatOffensiveTrinketEffect(38287, stats.Stats{stats.AttackPower: 278, stats.RangedAttackPower: 278}, time.Second*20, time.Minute*2) // Empty Direbrew Mug

	// Defensive trinkets. Keep these in order by item ID.
	core.NewSimpleStatDefensiveTrinketEffect(29387, stats.Stats{stats.BlockValue: 200}, time.Second*40, time.Minute*2) // Gnomeregan Auto-Blocker 600
	core.NewSimpleStatDefensiveTrinketEffect(32501, stats.Stats{stats.Health: 1750}, time.Second*20, time.Minute*3)    // Shadowmoon Insignia
	core.NewSimpleStatDefensiveTrinketEffect(38289, stats.Stats{stats.BlockValue: 200}, time.Second*40, time.Minute*2) // Coren's Lucky Coin

	// Proc effects. Keep these in order by item ID.

	core.NewItemEffect(11815, func(agent core.Agent) {
		character := agent.GetCharacter()
		if !character.AutoAttacks.AutoSwingMelee {
			return
		}

		var handOfJusticeSpell *core.Spell
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 2,
		}
		procChance := 0.013333

		character.RegisterAura(core.Aura{
			Label:    "Hand of Justice",
			Duration: core.NeverExpires,
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				config := *character.AutoAttacks.MHConfig()
				config.ActionID = core.ActionID{ItemID: 11815}
				handOfJusticeSpell = character.GetOrRegisterSpell(config)
			},
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				// https://wotlk.wowhead.com/spell=15600/hand-of-justice, proc mask = 20.
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				if !icd.IsReady(sim) {
					return
				}

				if sim.RandomFloat("HandOfJustice") > procChance {
					return
				}
				icd.Use(sim)

				aura.Unit.AutoAttacks.MaybeReplaceMHSwing(sim, handOfJusticeSpell).Cast(sim, result.Target)
			},
		})
	})

	core.NewItemEffect(28830, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Dragonspine Trophy Proc", core.ActionID{ItemID: 28830}, stats.Stats{stats.MeleeHaste: 325}, time.Second*10)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 20,
		}
		ppmm := character.AutoAttacks.NewPPMManager(1.0, core.ProcMaskMeleeOrRanged)

		character.RegisterAura(core.Aura{
			Label:    "Dragonspine Trophy",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if !icd.IsReady(sim) {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "dragonspine") {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(30627, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Tsunami Talisman Proc", core.ActionID{ItemID: 30627}, stats.Stats{stats.AttackPower: 340, stats.RangedAttackPower: 340}, time.Second*10)
		const procChance = 0.1

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 45,
		}

		character.RegisterAura(core.Aura{
			Label:    "Tsunami Talisman",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Outcome.Matches(core.OutcomeCrit) {
					return
				}
				if !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if sim.RandomFloat("Tsunami Talisman") > procChance {
					return
				}

				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(32505, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Madness of the Betrayer Proc", core.ActionID{ItemID: 32505}, stats.Stats{stats.ArmorPenetration: 42}, time.Second*10)

		ppmm := character.AutoAttacks.NewPPMManager(1.0, core.ProcMaskMeleeOrRanged)

		character.RegisterAura(core.Aura{
			Label:    "Madness of the Betrayer",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Madness of the Betrayer") {
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(34427, func(agent core.Agent) {
		character := agent.GetCharacter()

		var bonusPerStack stats.Stats
		procAura := character.RegisterAura(core.Aura{
			Label:     "Blackened Naaru Sliver Proc",
			ActionID:  core.ActionID{ItemID: 34427},
			Duration:  time.Second * 20,
			MaxStacks: 10,
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				bonusPerStack = stats.Stats{stats.AttackPower: 44, stats.RangedAttackPower: 44}
			},
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				character.AddStatsDynamic(sim, bonusPerStack.Multiply(float64(newStacks-oldStacks)))
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					aura.AddStack(sim)
				}
			},
		})

		const procChance = 0.1

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 45,
		}

		character.RegisterAura(core.Aura{
			Label:    "Blackened Naaru Sliver",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				// mask 340
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if sim.RandomFloat("Blackened Naaru Sliver") > procChance {
					return
				}

				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(34472, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura("Shard of Contempt Proc", core.ActionID{ItemID: 34472}, stats.Stats{stats.AttackPower: 230, stats.RangedAttackPower: 230}, time.Second*20)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 45,
		}
		const procChance = 0.1

		character.RegisterAura(core.Aura{
			Label:    "Shard of Contempt",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					return
				}
				if !icd.IsReady(sim) {
					return
				}
				if sim.RandomFloat("Shard of Contempt") > procChance {
					return
				}

				icd.Use(sim)
				procAura.Activate(sim)
			},
		})
	})

	core.AddEffectsToTest = true
}
