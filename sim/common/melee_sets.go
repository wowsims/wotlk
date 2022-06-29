package common

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Keep these in alphabetical order.

var ItemSetBurningRage = core.NewItemSet(core.ItemSet{
	Name: "Burning Rage",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.MeleeHit, 20)
		},
	},
})

var ItemSetDesolationBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Desolation Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.MeleeHit, 35)
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			procAura := character.NewTemporaryStatsAura("Desolation Battlegear Proc", core.ActionID{SpellID: 37617}, stats.Stats{stats.AttackPower: 160, stats.RangedAttackPower: 160}, time.Second*15)

			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Second * 20,
			}
			const procChance = 0.01

			character.RegisterAura(core.Aura{
				Label:    "Desolation Battlegear",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if !spellEffect.Landed() {
						return
					}
					if !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
						return
					}
					if !icd.IsReady(sim) {
						return
					}
					if sim.RandomFloat("Desolation Battlegear") > procChance {
						return
					}
					icd.Use(sim)
					procAura.Activate(sim)
				},
			})
		},
	},
})

var ItemSetDoomplateBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Doomplate Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.MeleeHit, 35)
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			procAura := character.NewTemporaryStatsAura("Doomplate Battlegear Proc", core.ActionID{SpellID: 37611}, stats.Stats{stats.AttackPower: 160, stats.RangedAttackPower: 160}, time.Second*15)

			const procChance = 0.02
			character.RegisterAura(core.Aura{
				Label:    "Doomplate Battlegear",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if !spellEffect.Landed() {
						return
					}
					if !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
						return
					}
					if sim.RandomFloat("Doomplate Battlegear") > procChance {
						return
					}
					procAura.Activate(sim)
				},
			})
		},
	},
})

var ItemSetEbonNetherscale = core.NewItemSet(core.ItemSet{
	Name: "Netherscale Armor",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.MeleeHit, 20)
		},
	},
})

var ItemSetFaithInFelsteel = core.NewItemSet(core.ItemSet{
	Name: "Faith in Felsteel",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.Strength, 25)
		},
	},
})

var ItemSetFelstalker = core.NewItemSet(core.ItemSet{
	Name: "Felstalker Armor",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.MeleeHit, 20)
		},
	},
})

var ItemSetFistsOfFury = core.NewItemSet(core.ItemSet{
	Name: "The Fists of Fury",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()

			procSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 41989},
				SpellSchool: core.SpellSchoolFire,
				ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
					ProcMask:         core.ProcMaskEmpty,
					DamageMultiplier: 1,
					ThreatMultiplier: 1,

					BaseDamage:     core.BaseDamageConfigRoll(100, 150),
					OutcomeApplier: character.OutcomeFuncMagicHitAndCrit(character.DefaultSpellCritMultiplier()),
				}),
			})

			ppmm := character.AutoAttacks.NewPPMManager(2, core.ProcMaskMelee)

			character.RegisterAura(core.Aura{
				Label:    "Fists of Fury",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}
					if !ppmm.Proc(sim, spellEffect.ProcMask, "The Fists of Fury") {
						return
					}

					procSpell.Cast(sim, spellEffect.Target)
				},
			})
		},
	},
})

var ItemSetFlameGuard = core.NewItemSet(core.ItemSet{
	Name: "Flame Guard",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.Parry, 20)
		},
	},
})

var ItemSetPrimalstrike = core.NewItemSet(core.ItemSet{
	Name: "Primal Intent",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.AttackPower, 40)
			agent.GetCharacter().AddStat(stats.RangedAttackPower, 40)
		},
	},
})

var ItemSetStrengthOfTheClefthoof = core.NewItemSet(core.ItemSet{
	Name: "Strength of the Clefthoof",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.Strength, 20)
		},
	},
})

var ItemSetTwinBladesOfAzzinoth = core.NewItemSet(core.ItemSet{
	Name: "The Twin Blades of Azzinoth",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()

			if character.CurrentTarget.MobType == proto.MobType_MobTypeDemon {
				character.PseudoStats.MobTypeAttackPower += 200
			}
			procAura := character.NewTemporaryStatsAura("Twin Blade of Azzinoth Proc", core.ActionID{SpellID: 41435}, stats.Stats{stats.MeleeHaste: 450}, time.Second*10)

			ppmm := character.AutoAttacks.NewPPMManager(1.0, core.ProcMaskMelee)
			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Second * 45,
			}

			character.RegisterAura(core.Aura{
				Label:    "Twin Blades of Azzinoth",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if !spellEffect.Landed() {
						return
					}

					// https://tbc.wowhead.com/spell=41434/the-twin-blades-of-azzinoth, proc mask = 20.
					if !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}

					if !icd.IsReady(sim) {
						return
					}

					if !ppmm.Proc(sim, spellEffect.ProcMask, "Twin Blades of Azzinoth") {
						return
					}
					icd.Use(sim)
					procAura.Activate(sim)
				},
			})
		},
	},
})

var ItemSetWastewalkerArmor = core.NewItemSet(core.ItemSet{
	Name: "Wastewalker Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.MeleeHit, 35)
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			procAura := character.NewTemporaryStatsAura("Wastewalker Armor Proc", core.ActionID{SpellID: 37618}, stats.Stats{stats.AttackPower: 160, stats.RangedAttackPower: 160}, time.Second*15)

			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Second * 20,
			}
			const procChance = 0.02

			character.RegisterAura(core.Aura{
				Label:    "Wastewalker Armor",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					if !spellEffect.Landed() {
						return
					}
					if !spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
						return
					}
					if !icd.IsReady(sim) {
						return
					}
					if sim.RandomFloat("Wastewalker Armor") > procChance {
						return
					}
					icd.Use(sim)
					procAura.Activate(sim)
				},
			})
		},
	},
})
